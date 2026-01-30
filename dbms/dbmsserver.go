// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !gui

package dbms

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/dbms/commands"
	"github.com/apmckinlay/gsuneido/dbms/mux"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/atomics"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/time/rate"
)

// This is the multiplexed server.
// It only works with the gSuneido multiplexed client.

var workers *mux.Workers

var serverConns = make(map[uint32]*serverConn)
var serverConnsLock sync.Mutex // guards serverConns and idleCount

var lastNum atomic.Int64 // used for queries, cursors

// serverConn is one client connection which handles multiple sessions
type serverConn struct {
	dbms       IDbms
	conn       net.Conn
	sessions   map[uint32]*serverSession // the sessions on this connection
	remoteAddr string
	Sviews
	idleCount    int          // guarded by serverConnsLock
	sessionsLock sync.Mutex   // guards sessions
	logSize      atomic.Int32 // cumulative size of logged data in bytes
	nonce        string       // for authentication, shared across sessions
	nonceOld     bool         // for two-phase expiration like tokens
	// id is primarily used as a key to store the set of connections in a map
	id uint32
}

// serverSession handles one client session.
// It should be thread contained, other than sessionId
type serverSession struct {
	sc            *serverConn
	*mux.WriteBuf         // set per request
	thread        *Thread // set per request
	trans         map[int]ITran
	cursors       map[int]ICursor
	queries       map[int]IQuery
	tranQueries   map[int]intSet // transaction id => query ids
	queryTrans    map[int]int    // query id => transaction id
	sessionId     atomics.String
	mux.ReadBuf
	// id is primarily used as a key to store the set of sessions in a map
	id uint32
}

type intSet map[int]struct{}

//go:embed server.key
var ServerKey []byte

// Server listens and accepts connections. It never returns.
func Server(dbms *DbmsLocal) {
	workers = mux.NewWorkers(doRequest)
	cert, err := tls.X509KeyPair(ServerCert, ServerKey)
	if err != nil {
		Fatal("Failed to load embedded key pair:", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	// Listen for plain TCP connection to handle version mismatch
	l, err := net.Listen("tcp", ":"+options.Port)
	if err != nil {
		Fatal(err)
	}
	go background()
	limiter := rate.NewLimiter(rate.Limit(8), 4)
	context := context.Background()
	for {
		limiter.Wait(context)
		conn, err := l.Accept()
		if err != nil {
			Fatal("DbmsServer:", err)
		}
		// start a new goroutine to avoid blocking
		go newServerConn(dbms, conn, config)
	}
}

func background() {
	for {
		time.Sleep(backgroundInterval)
		idleCheck()
		expireTokens()
		expireNonces()
	}
}

const backgroundInterval = time.Minute

func idleCheck() {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	for _, sc := range serverConns {
		sc.idleCount++ // reset by doRequest
		if sc.idleCount > options.TimeoutMinutes {
			sc.serverLog("closing idle connection")
			sc.close()
		}
	}
}

func expireTokens() {
	tokensLock.Lock()
	defer tokensLock.Unlock()
	for token, old := range tokens {
		if old {
			delete(tokens, token)
		} else {
			tokens[token] = true // mark it as old
		}
	}
}

func expireNonces() {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	for _, sc := range serverConns {
		if sc.nonceOld {
			sc.nonce = ""
			sc.nonceOld = false
		} else if sc.nonce != "" {
			sc.nonceOld = true
		}
	}
}

func (sc *serverConn) serverLog(args ...any) {
	args = append([]any{"dbms server:", sc.remoteAddr + ":"}, args...)
	log.Println(args...)
}

func newServerConn(dbms *DbmsLocal, conn net.Conn, config *tls.Config) {
	trace.ClientServer.Println("server connection")
	conn.Write(hello())
	if errmsg := checkHello(conn); errmsg != "" {
		if strings.HasPrefix(errmsg, "version mismatch") {
			serverVersionMismatch(dbms, conn)
		}
		conn.Close()
		return
	}
	// Upgrade to TLS after successful hello
	tlsConn := tls.Server(conn, config)
	if err := tlsConn.Handshake(); err != nil {
		log.Println("ERROR: dbms server: TLS handshake failed:", err)
		conn.Close()
		return
	}
	addr := str.BeforeLast(tlsConn.RemoteAddr().String(), ":") // strip port
	msc := mux.NewServerConn(tlsConn)
	sc := &serverConn{dbms: dbms, id: msc.Id(), conn: tlsConn, remoteAddr: addr,
		sessions: make(map[uint32]*serverSession)}
	if dbms.db.HaveUsers() {
		sc.dbms = &DbmsUnauth{dbms: dbms}
	}
	serverConnsLock.Lock()
	serverConns[sc.id] = sc
	serverConnsLock.Unlock()
	msc.Run(workers.Submit)
}

func serverVersionMismatch(dbms *DbmsLocal, conn net.Conn) {
	rt := dbms.db.NewReadTran()
	def := dbms.LibGet1(rt, "stdlib", "VersionMismatch", nil)
	if len(def) > 2 {
		Fatal("VersionMismatch must have a single definition")
	}
	s := def[1]
	n := len(s)
	conn.Write([]byte{byte(n >> 8), byte(n)})
	io.WriteString(conn, s)
}

// doRequest is called by workers (multi-threaded)
func doRequest(wb *mux.WriteBuf, th *Thread, id uint64, req []byte) {
	connId := uint32(id >> 32)
	serverConnsLock.Lock()
	if req == nil { // closing, e.g. error or lost connection in mux reader
		// log.Println("dbms server: nil request (closing)")
		delete(serverConns, connId)
		serverConnsLock.Unlock()
		return
	}
	sc := serverConns[connId]
	if sc == nil {
		serverConnsLock.Unlock()
		log.Println("dbms server doRequest: no such connection")
		return
	}
	sc.idleCount = 0
	serverConnsLock.Unlock()

	sid := uint32(id)
	sc.sessionsLock.Lock()
	ss := sc.sessions[sid]
	if ss == nil { // new session
		trace.ClientServer.Println("server session", sid)
		ss = &serverSession{
			id:          sid,
			sc:          sc,
			trans:       make(map[int]ITran),
			cursors:     make(map[int]ICursor),
			queries:     make(map[int]IQuery),
			tranQueries: make(map[int]intSet),
			queryTrans:  make(map[int]int),
		}
		ss.sessionId.Store(sc.remoteAddr)
		sc.sessions[sid] = ss
	}
	sc.sessionsLock.Unlock()

	ss.ReadBuf.SetBuf(req)
	ss.WriteBuf = wb
	th.SetSession(ss.sessionId.Load())
	th.SetSviews(&sc.Sviews)
	ss.thread = th
	ss.request()
}

func (ss *serverSession) request() {
	var icmd commands.Command
	defer func() {
		if e := recover(); e != nil {
			LogInternalError(ss.thread, ss.sessionId.Load(), e)
			ss.ResetWrite()
			ss.PutBool(false).PutStr(errToStr(e)).EndMsg()
		}
	}()
	icmd = ss.GetCmd()
	if int(icmd) >= len(cmds) {
		serverConnsLock.Lock()
		defer serverConnsLock.Unlock()
		ss.sc.close()
		ss.sc.serverLog("closed connection: invalid command")
		return
	}
	cmd := cmds[icmd]
	cmd(ss)
	assert.That(ss.Remaining() == 0) // should consume entire message
	if icmd != commands.EndSession {
		ss.EndMsg()
	}
}

func errToStr(e any) string {
	if t, ok := e.(interface{ ToStr() (string, bool) }); ok {
		if s, ok := t.ToStr(); ok {
			return s
		}
	}
	return fmt.Sprint(e)
}

func (ss *serverSession) error(err string) {
	ss.close()
	log.Panicln("dbms server, closing session:", err)
}

func (ss *serverSession) close() {
	ss.abort()
	ss.sc.sessionsLock.Lock()
	delete(ss.sc.sessions, ss.id)
	ss.sc.sessionsLock.Unlock()
}

func (ss *serverSession) abort() {
	for _, tran := range ss.trans {
		tran.Abort()
	}
}

// close is called by idleCheck and bad request. MUST hold serverConnsLock
func (sc *serverConn) close() {
	trace.ClientServer.Println("closing connection")
	sc.conn.Close()
	delete(serverConns, sc.id)
	// intentionally don't close the sessions and their transactions
	// because that would require additional locking
	// read-only transactions don't need to be closed
	// and update transactions will time out
}

// Conns is used by HttpStatus
func Conns() string {
	var sb strings.Builder
	sb.WriteString("<p>Connections:</p>\r\n<ul>\r\n")
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	var conns []*serverConn
	for _, sc := range serverConns {
		conns = append(conns, sc)
	}
	sort.Slice(conns,
		func(i, j int) bool { return conns[i].remoteAddr < conns[j].remoteAddr })
	for _, sc := range conns {
		sb.WriteString("<li>")
		sb.WriteString(sc.remoteAddr)
		sb.WriteString("</li>\r\n<ul>\r\n")
		sc.sessionsLock.Lock()
		var sessions []string
		for _, ss := range sc.sessions {
			sessions = append(sessions, ss.sessionId.Load())
		}
		sort.Strings(sessions)
		for _, sid := range sessions {
			sb.WriteString("<li>")
			sb.WriteString(sid)
			sb.WriteString("</li>\r\n")
		}
		sc.sessionsLock.Unlock()
		sb.WriteString("</ul>\r\n")
	}
	sb.WriteString("</ul>\r\n")
	return sb.String()
}

func StopServer() {
	defer func() {
		if e := recover(); e != nil {
			log.Println("StopServer", e)
		}
	}()
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	for _, sc := range serverConns {
		sc.conn.Close()
	}
	serverConns = nil
}

// commands ---------------------------------------------------------

// NOTE: as soon as we send the response we may get a new request
// which will take over the serverSession (ss)
// so we can't use the serverSession after sending the response

func cmdAbort(ss *serverSession) {
	tn := ss.GetInt()
	tran := ss.tran(tn)
	tran.Abort()
	ss.deleteTran(tn)
	ss.PutBool(true)
}

func (ss *serverSession) deleteTran(tn int) {
	delete(ss.trans, tn)
	for qn := range ss.tranQueries[tn] {
		delete(ss.queries, qn)
		delete(ss.queryTrans, qn)
		if len(ss.queries) != len(ss.queryTrans) {
			log.Println("ERROR: deleteTran", len(ss.queries), "!=", len(ss.queryTrans))
		}
	}
	delete(ss.tranQueries, tn)
}

func (ss *serverSession) getTran() (ITran, int) {
	tn := ss.GetInt()
	if tn == 0 {
		return nil, 0
	}
	return ss.tran(tn), tn
}

func (ss *serverSession) tran(tn int) ITran {
	tran, ok := ss.trans[tn]
	if !ok {
		ss.error("transaction not found")
	}
	return tran
}

func cmdAction(ss *serverSession) {
	tran, _ := ss.getTran()
	action := ss.GetStr()
	n := tran.Action(ss.thread, action)
	ss.PutBool(true).PutInt(n)
}

func cmdAdmin(ss *serverSession) {
	s := ss.GetStr()
	ss.sc.dbms.Admin(s, &ss.sc.Sviews)
	ss.PutBool(true)
}

func cmdAuth(ss *serverSession) {
	s := ss.GetStr()
	if _, ok := ss.sc.dbms.(*DbmsUnauth); !ok {
		panic("already authorized")
	}
	result := ss.auth(s)
	if result {
		ss.sc.dbms = ss.sc.dbms.(*DbmsUnauth).dbms // remove DbmsUnauth
	}
	ss.PutBool(true).PutBool(result)
}

func (ss *serverSession) auth(s string) bool {
	nonce := ss.sc.nonce
	ss.sc.nonce = ""
	ss.sc.nonceOld = false
	if AuthUser(ss.thread, s, nonce) {
		return true
	}
	return AuthToken(s)
}

func cmdAsof(ss *serverSession) {
	tn := ss.GetInt()
	asof := ss.GetInt64()
	tran := ss.tran(tn)
	ss.PutBool(true).PutInt64(tran.Asof(asof))
}

func cmdCheck(ss *serverSession) {
	full := ss.GetBool()
	s := ss.sc.dbms.Check(full)
	ss.PutBool(true).PutStr(s)
}

func cmdClose(ss *serverSession) {
	qn := ss.GetInt()
	switch ss.GetChar() {
	case 'q':
		q := ss.queries[qn]
		if q == nil {
			ss.error("query not found")
		}
		delete(ss.queries, qn)
		tn := ss.queryTrans[qn]
		delete(ss.queryTrans, qn)
		if len(ss.queries) != len(ss.queryTrans) {
			log.Println("ERROR: cmdClose", len(ss.queries), "!=", len(ss.queryTrans))
		}
		delete(ss.tranQueries[tn], qn)
		q.Close()
	case 'c':
		c := ss.cursors[qn]
		if c == nil {
			ss.error("cursor not found")
		}
		delete(ss.cursors, qn)
		c.Close()
	default:
		ss.error("dbms server expected q or c")
	}
	ss.PutBool(true)
}

func cmdCommit(ss *serverSession) {
	tn := ss.GetInt()
	tran := ss.tran(tn)
	result := tran.Complete()
	ss.deleteTran(tn)
	ss.PutBool(true)
	if result == "" {
		ss.PutBool(true)
	} else {
		ss.PutBool(false).PutStr(result)
	}
}

func cmdConnections(ss *serverSession) {
	ss.PutBool(true).PutVal(connections())
}

func connections() *SuObject {
	list := &SuObject{}
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	for _, sc := range serverConns {
		sc.sessionsLock.Lock()
		for _, ss := range sc.sessions {
			list.Add(SuStr(ss.sessionId.Load()))
		}
		sc.sessionsLock.Unlock()
	}
	return list
}

func cmdCursor(ss *serverSession) {
	query := ss.GetStr()
	q := ss.sc.dbms.Cursor(query, &ss.sc.Sviews)
	num := int(lastNum.Add(1))
	ss.cursors[num] = q
	ss.PutBool(true).PutInt(num)
}

func cmdCursors(ss *serverSession) {
	ss.PutBool(true).PutInt(len(ss.cursors))
}

func cmdEndSession(ss *serverSession) {
	// ss.sc.serverLog("closing connection: received EndSession")
	ss.close()
	// no response
}

func cmdErase(ss *serverSession) {
	defer ss.thread.Suneido.Store(ss.thread.Suneido.Load())
	ss.thread.Suneido.Store(nil) // use main Suneido object
	tran, _ := ss.getTran()
	table := ss.GetStr()
	off := uint64(ss.GetInt64())
	tran.Delete(ss.thread, table, off)
	ss.PutBool(true)
}

func cmdExec(ss *serverSession) {
	ob := ss.GetVal()
	v := ss.sc.dbms.Exec(ss.thread, ob)
	ss.PutResult(v)
}

func cmdFinal(ss *serverSession) {
	final := ss.sc.dbms.Final()
	ss.PutBool(true).PutInt(final)
}

func cmdGet(ss *serverSession) {
	tbl, hdr, row := ss.getQorTC()
	ss.rowResult(tbl, hdr, false, row)
}

func (ss *serverSession) getQorTC() (tbl string, hdr *Header, row Row) {
	dir := ss.getDir()
	t, _ := ss.getTran()
	if t == nil {
		q := ss.getQuery()
		hdr = q.Header()
		row, tbl = q.Get(ss.thread, dir)
	} else {
		c := ss.getCursor()
		hdr = c.Header()
		row, tbl = c.Get(ss.thread, t, dir)
	}
	return
}

func (ss *serverSession) getDir() Dir {
	dir := Dir(ss.GetByte())
	trace.ClientServer.Println("    <-", string(dir))
	return Dir(dir)
}

const maxRec = 1024 * 1024 // 1 mb

func (ss *serverSession) rowResult(tbl string, hdr *Header, sendHdr bool, row Row) {
	if row == nil {
		ss.PutBool(true).PutBool(false)
	} else {
		rec, flds := rowToRecord(row, hdr)
		if len(rec) > maxRec {
			panic("result too large")
		}
		ss.PutBool(true).PutBool(true).PutInt(int(row[0].Off))
		if sendHdr {
			ss.PutStrs(hdr.AppendDerived(flds))
		}
		ss.PutStr(tbl)
		ss.PutRec(rec)
	}
}

func cmdGetOne(ss *serverSession) {
	var dir Dir
	switch ss.GetChar() {
	case '+':
		dir = Next
	case '-':
		dir = Prev
	case '1':
		dir = Only
	case '@':
		dir = Any
	case '?':
		dir = Strat
	default:
		ss.error("dbms server: expected + - 1 @")
	}
	tran, _ := ss.getTran()
	query := ss.GetVal()
	var g func(*Thread, Value, Dir) (Row, *Header, string)
	if tran == nil {
		g = ss.sc.dbms.Get
	} else {
		g = tran.Get
	}
	row, hdr, tbl := g(ss.thread, query, dir)
	ss.rowResult(tbl, hdr, true, row)
}

func cmdHeader(ss *serverSession) {
	hdr := ss.getQorC().Header()
	ss.PutBool(true).PutStrs(hdr.Schema())
}

func (ss *serverSession) getQorC() (qc IQueryCursor) {
	n := ss.GetInt()
	switch ss.GetChar() {
	case 'q':
		qc = ss.queries[n]
	case 'c':
		qc = ss.cursors[n]
	default:
		ss.error("dbms server expected q or c")
	}
	if qc == nil {
		ss.error("dbms server: query/cursor not found")
	}
	return qc
}

func cmdInfo(ss *serverSession) {
	info := ss.sc.dbms.Info()
	ss.PutBool(true).PutVal(info)
}

func cmdKeys(ss *serverSession) {
	keys := ss.getQorC().Keys()
	ss.PutBool(true).PutStrs(keys)
}

func cmdKill(ss *serverSession) {
	sessionId := ss.GetStr()
	n := kill(sessionId)
	ss.PutBool(true).PutInt(n)
}

func kill(sid string) int {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	nkilled := 0
	for id, sc := range serverConns {
		func() {
			sc.sessionsLock.Lock()
			defer sc.sessionsLock.Unlock()
			for _, ss := range sc.sessions {
				if ss.sessionId.Load() == sid {
					sc.serverLog("dbms server: kill:", sid)
					delete(serverConns, id)
					sc.conn.Close()
					nkilled++
					break
				}
			}
		}()
	}
	return nkilled
}

func cmdLibGet(ss *serverSession) {
	name := ss.GetStr()
	defs := ss.sc.dbms.LibGet(name)
	ss.PutBool(true).PutInt(len(defs) / 2)
	for i := 0; i < len(defs); i += 2 {
		ss.PutStr(defs[i]).PutInt(len(defs[i+1]))
	}
	for i := 1; i < len(defs); i += 2 {
		ss.PutBuf(defs[i])
	}
}

func cmdLibraries(ss *serverSession) {
	libs := ss.sc.dbms.Libraries()
	ss.PutBool(true).PutStrs(libs)
}

func cmdLog(ss *serverSession) {
	s := ss.GetStr()
	if msg := ss.sc.limitLog(s); msg != "" {
		ss.sc.dbms.Log(msg)
	}
	ss.PutBool(true) // return true regardless
}

const logLimit = 10 * 1024 // ???

// limitLog limits log size per connection.
// Returns the message to log, or empty string if ignored.
func (sc *serverConn) limitLog(s string) string {
	logBytes := int32(len(s) + 1) // +1 for newline added by log.Println
	newSize := sc.logSize.Add(logBytes)
	if newSize <= logLimit {
		return s
	} else if newSize-logBytes <= logLimit {
		// First time over limit - log warning exactly once
		return "log size limit exceeded (10KB), ignoring further logs"
	}
	return ""
}

func cmdNonce(ss *serverSession) {
	ss.sc.nonce = Nonce()
	ss.sc.nonceOld = false
	ss.PutBool(true).PutStr_(ss.sc.nonce)
}

func cmdOrder(ss *serverSession) {
	order := ss.getQorC().Order()
	ss.PutBool(true).PutStrs(order)
}

func cmdOutput(ss *serverSession) {
	defer ss.thread.Suneido.Store(ss.thread.Suneido.Load())
	ss.thread.Suneido.Store(nil) // use main Suneido object
	q := ss.getQuery()
	rec := ss.GetRec()
	q.Output(ss.thread, rec)
	ss.PutBool(true)
}

func (ss *serverSession) getQuery() IQuery {
	qn := ss.GetInt()
	q := ss.queries[qn]
	if q == nil {
		ss.error("dbms server: query not found")
	}
	return q
}

func (ss *serverSession) getCursor() ICursor {
	qn := ss.GetInt()
	c := ss.cursors[qn]
	if c == nil {
		ss.error("dbms server: query not found")
	}
	return c
}

func cmdQuery(ss *serverSession) {
	tran, tn := ss.getTran()
	query := ss.GetStr()
	q := tran.Query(query, &ss.sc.Sviews)
	qn := int(lastNum.Add(1))
	ss.queries[qn] = q
	ss.queryTrans[qn] = tn
	if ss.tranQueries[tn] == nil {
		ss.tranQueries[tn] = make(intSet)
	}
	ss.tranQueries[tn][qn] = struct{}{}
	if len(ss.queries) != len(ss.queryTrans) {
		log.Println("ERROR: cmdQuery", len(ss.queries), "!=", len(ss.queryTrans))
	}
	ss.PutBool(true).PutInt(qn)
}

func cmdReadCount(ss *serverSession) {
	ss.getTran()
	ss.PutBool(true).PutInt(0) //TODO
}

func cmdRewind(ss *serverSession) {
	qc := ss.getQorC()
	qc.Rewind()
	ss.PutBool(true)
}

func cmdRun(ss *serverSession) {
	s := ss.GetStr()
	v := ss.sc.dbms.Run(ss.thread, s)
	ss.PutResult(v)
}

func cmdSessionId(ss *serverSession) {
	s := ss.GetStr()
	if s != "" {
		ss.sessionId.Store(s)
	}
	ss.PutBool(true).PutStr(ss.sessionId.Load())
}

func cmdSize(ss *serverSession) {
	n := ss.sc.dbms.Size()
	ss.PutBool(true).PutInt64(int64(n))
}

func cmdStrategy(ss *serverSession) {
	qc := ss.getQorC()
	formatted := ss.GetBool()
	strategy := qc.Strategy(formatted)
	ss.PutBool(true).PutStr(strategy)
}

func cmdTimestamp(ss *serverSession) {
	ts := ss.sc.dbms.Timestamp()
	ss.PutBool(true).PutVal(ts)
}

func cmdToken(ss *serverSession) {
	tok := Token()
	ss.PutBool(true).PutStr(tok)
}

func cmdTransaction(ss *serverSession) {
	update := ss.GetBool()
	tran := ss.sc.dbms.Transaction(update)
	tn := tran.Num()
	ss.trans[tn] = tran
	ss.PutBool(true).PutInt(tn)
}

func cmdTransactions(ss *serverSession) {
	list := ss.sc.dbms.Transactions()
	ss.PutBool(true).PutVal(list)
}

func cmdUpdate(ss *serverSession) {
	tran, _ := ss.getTran()
	table := ss.GetStr()
	off := uint64(ss.GetInt64())
	rec := ss.GetRec()
	newoff := tran.Update(ss.thread, table, off, rec)
	ss.PutBool(true).PutInt(int(newoff))
}

func cmdWriteCount(ss *serverSession) {
	ss.getTran()
	ss.PutBool(true).PutInt(0) //TODO
}

type command func(ss *serverSession)

var cmds = []command{ // order must match commmands.go
	cmdAbort,
	cmdAdmin,
	cmdAuth,
	cmdCheck,
	cmdClose,
	cmdCommit,
	cmdConnections,
	cmdCursor,
	cmdCursors,
	cmdErase,
	cmdExec,
	cmdStrategy,
	cmdFinal,
	cmdGet,
	cmdGetOne,
	cmdHeader,
	cmdInfo,
	cmdKeys,
	cmdKill,
	cmdLibGet,
	cmdLibraries,
	cmdLog,
	cmdNonce,
	cmdOrder,
	cmdOutput,
	cmdQuery,
	cmdReadCount,
	cmdAction,
	cmdRewind,
	cmdRun,
	cmdSessionId,
	cmdSize,
	cmdTimestamp,
	cmdToken,
	cmdTransaction,
	cmdTransactions,
	cmdUpdate,
	cmdWriteCount,
	cmdEndSession,
	cmdAsof,
	nil,
}

func init() {
	assert.That(cmds[commands.Asof] != nil && cmds[commands.Asof+1] == nil)
}
