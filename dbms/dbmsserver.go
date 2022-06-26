// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/dbms/commands"
	"github.com/apmckinlay/gsuneido/dbms/mux"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// This is the multiplexed server.
// It only works with the gSuneido mulitplexed client.

var workers *mux.Workers

var serverConns = make(map[uint32]*serverConn)
var serverConnsLock sync.Mutex // guards serverConns and idleCount

// serverConn is one client connection which handles multiple sessions
type serverConn struct {
	// id is primarily used as a key to store the set of connections in a map
	id           uint32
	remoteAddr   string
	ended        bool
	dbms         IDbms
	conn         net.Conn
	sessionsLock sync.Mutex                // guards sessions
	sessions     map[uint32]*serverSession // the sessions on this connection
	Sviews
	idleCount int // guarded by serverConnsLock
}

// serverSession is one client session (thread)
type serverSession struct {
	sc *serverConn
	mux.ReadBuf
	*mux.WriteBuf
	thread *runtime.Thread
	// id is primarily used as a key to store the set of sessions in a map
	id        uint32
	sessionId string
	nonce     string
	trans     map[int]ITran
	cursors   map[int]ICursor
	queries   map[int]IQuery
	lastNum   int // used for queries, cursors, and transactions
}

// Server listens and accepts connections. It never returns.
func Server(dbms *DbmsLocal) {
	l, err := net.Listen("tcp", ":"+options.Port)
	if err != nil {
		Fatal(err)
	}
	defer l.Close()
	go idleTimeout()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := l.Accept()
		if err != nil {
			// error handling based on Go net/http
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Println("ERROR server accept:", err)
				time.Sleep(tempDelay)
				continue
			}
			Fatal(err)
		}
		tempDelay = 0
		newServerConn(dbms, conn)
	}
}

func idleTimeout() {
	for {
		time.Sleep(idleCheckInterval)
		idleCheck()
	}
}

const idleCheckInterval = time.Minute
const maxIdleCount = 2 * 60 // 2 hours if interval is one minute

func idleCheck() {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	for _, sc := range serverConns {
		sc.idleCount++
		if sc.idleCount > maxIdleCount {
			log.Println("closing idle connection", sc.remoteAddr)
			sc.close()
			delete(serverConns, sc.id)
		}
	}
}

func newServerConn(dbms *DbmsLocal, conn net.Conn) {
	trace.ClientServer.Println("server connection")
	workers = mux.NewWorkers(doRequest)
	sendHello(conn)
	addr := str.BeforeFirst(conn.RemoteAddr().String(), ":")
	connId := mux.NewServerConn(conn, workers.Submit)
	sc := &serverConn{dbms: dbms, id: connId, conn: conn, remoteAddr: addr,
		sessions: make(map[uint32]*serverSession)}
	if dbms.db.HaveUsers() {
		sc.dbms = &DbmsUnauth{dbms: dbms}
	}
	serverConnsLock.Lock()
	serverConns[connId] = sc
	serverConnsLock.Unlock()
}

func sendHello(conn net.Conn) {
	conn.Write(hello)
}

var hello = func() []byte {
	buf := make([]byte, helloSize)
	copy(buf, "Suneido "+builtin.Built()+"\r\n")
	return buf
}()

// doRequest is called by workers
func doRequest(wb *mux.WriteBuf, th *runtime.Thread, id uint64, req []byte) {
	connId := uint32(id >> 32)
	serverConnsLock.Lock()
	if req == nil { // closing
		delete(serverConns, connId)
		serverConnsLock.Unlock()
		return
	}
	sc := serverConns[connId]
	sc.idleCount = 0
	serverConnsLock.Unlock()

	sid := uint32(id)
	sc.sessionsLock.Lock()
	ss := sc.sessions[sid]
	if ss == nil { // new session
		trace.ClientServer.Println("server session", sid)
		ss = &serverSession{
			id:        sid,
			sc:        sc,
			sessionId: sc.remoteAddr,
			trans:     make(map[int]ITran),
			cursors:   make(map[int]ICursor),
			queries:   make(map[int]IQuery),
		}
		sc.sessions[sid] = ss
	}
	sc.sessionsLock.Unlock()

	ss.ReadBuf.SetBuf(req)
	ss.WriteBuf = wb
	th.SetSession(ss.sessionId)
	ss.thread = th
	ss.request()
}

func (ss *serverSession) request() {
	defer func() {
		if e := recover(); e != nil {
			LogInternalError(ss.thread, ss.sessionId, e)
			ss.ResetWrite()
			ss.PutBool(false).PutStr(fmt.Sprint(e)).EndMsg()
		}
	}()
	icmd := ss.GetCmd()
	if icmd == commands.Eof {
		ss.close()
		return
	}
	if int(icmd) >= len(cmds) {
		ss.close()
		log.Println("dbms server, closing connection: invalid command")
	}
	cmd := cmds[icmd]
	cmd(ss)
	assert.That(ss.Remaining() == 0) // should consume entire message
	if icmd != commands.EndSession {
		ss.EndMsg()
	}
}

func (ss *serverSession) error(err string) {
	ss.close()
	log.Panicln("dbms server, closing connection:", err)
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

func (sc *serverConn) close() {
	trace.ClientServer.Println("closing connection")
	sc.conn.Close()
	for _, ss := range sc.sessions {
		ss.abort()
	}
}

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
			sessions = append(sessions, ss.sessionId)
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

//-------------------------------------------------------------------

// NOTE: as soon as we send the response we may get a new request
// which will take over the serverSession (ss)
// so we can't use the serverSession after sending the response

func cmdAbort(ss *serverSession) {
	tn := ss.GetInt()
	tran := ss.tran(tn)
	tran.Abort()
	delete(ss.trans, tn)
	ss.PutBool(true)
}

func (ss *serverSession) getTran() ITran {
	tn := ss.GetInt()
	if tn == 0 {
		return nil
	}
	return ss.tran(tn)
}

func (ss *serverSession) tran(tn int) ITran {
	tran, ok := ss.trans[tn]
	if !ok {
		ss.error("transaction not found")
	}
	return tran
}

func cmdAction(ss *serverSession) {
	tran := ss.getTran()
	action := ss.GetStr()
	n := tran.Action(ss.thread, action, &ss.sc.Sviews)
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
	if AuthUser(ss.thread, s, ss.nonce) {
		ss.nonce = ""
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
	s := ss.sc.dbms.Check()
	ss.PutBool(true).PutStr(s)
}

func cmdClose(ss *serverSession) {
	n := ss.GetInt()
	switch ss.GetChar() {
	case 'q':
		q := ss.queries[n]
		if q == nil {
			ss.error("query not found")
		}
		delete(ss.queries, n)
		q.Close()
	case 'c':
		c := ss.cursors[n]
		if c == nil {
			ss.error("cursor not found")
		}
		delete(ss.cursors, n)
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
	delete(ss.trans, tn)
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
			list.Add(SuStr(ss.sessionId))
		}
		sc.sessionsLock.Unlock()
	}
	return list
}

func cmdCursor(ss *serverSession) {
	query := ss.GetStr()
	q := ss.sc.dbms.Cursor(query, &ss.sc.Sviews)
	ss.lastNum++
	ss.cursors[ss.lastNum] = q
	ss.PutBool(true).PutInt(ss.lastNum)
}

func cmdCursors(ss *serverSession) {
	ss.PutBool(true).PutInt(len(ss.cursors))
}

func cmdDump(ss *serverSession) {
	table := ss.GetStr()
	s := ss.sc.dbms.Dump(table)
	ss.PutBool(true).PutStr(s)
}

func cmdEndSession(ss *serverSession) {
	ss.close()
	// no response
}

func cmdErase(ss *serverSession) {
	tran := ss.getTran()
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
	tn := ss.GetInt()
	qn := ss.GetInt()
	if tn == 0 {
		q := ss.queries[qn]
		hdr = q.Header()
		row, tbl = q.Get(ss.thread, dir)
	} else {
		t := ss.trans[tn]
		c := ss.cursors[qn]
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

func rowToRecord(row Row, hdr *Header) (rec Record, fields []string) {
	if len(row) == 1 {
		assert.That(len(hdr.Fields) == 1)
		return row[0].Record, hdr.Fields[0]
	}
	var rb RecordBuilder
	// need to include deleted "-" in case it's used for updates
	fields = hdr.GetFields()
	for _, fld := range fields {
		rb.AddRaw(row.GetRaw(hdr, fld))
	}
	return rb.Trim().Build(), fields
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
	default:
		ss.error("dbms server: expected + - 1")
	}
	tran := ss.getTran()
	query := ss.GetStr()
	var g func(*Thread, string, Dir, *Sviews) (Row, *Header, string)
	if tran == nil {
		g = ss.sc.dbms.Get
	} else {
		g = tran.Get
	}
	row, hdr, tbl := g(ss.thread, query, dir, &ss.sc.Sviews)
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

func kill(remoteAddr string) int {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	nkilled := 0
	for id, sc := range serverConns {
		if sc.remoteAddr == remoteAddr {
			delete(serverConns, id)
			sc.conn.Close()
			nkilled++
		}
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

func cmdLoad(ss *serverSession) {
	table := ss.GetStr()
	n := ss.sc.dbms.Load(table)
	ss.PutBool(true).PutInt(n)
}

func cmdLog(ss *serverSession) {
	s := ss.GetStr()
	ss.sc.dbms.Log(s)
	ss.PutBool(true)
}

func cmdNonce(ss *serverSession) {
	ss.nonce = Nonce()
	ss.PutBool(true).PutStr_(ss.nonce)
}

func cmdOrder(ss *serverSession) {
	order := ss.getQorC().Order()
	ss.PutBool(true).PutStrs(order)
}

func cmdOutput(ss *serverSession) {
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

func cmdQuery(ss *serverSession) {
	tran := ss.getTran()
	query := ss.GetStr()
	q := tran.Query(query, &ss.sc.Sviews)
	ss.lastNum++
	ss.queries[ss.lastNum] = q
	ss.PutBool(true).PutInt(ss.lastNum)
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
		ss.sessionId = s
	}
	ss.PutBool(true).PutStr(ss.sessionId)
}

func cmdSize(ss *serverSession) {
	n := ss.sc.dbms.Size()
	ss.PutBool(true).PutInt64(int64(n))
}

func cmdStrategy(ss *serverSession) {
	qc := ss.getQorC()
	strategy := qc.Strategy()
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
	tn := ss.nextNum(update)
	ss.trans[tn] = tran
	ss.PutBool(true).PutInt(tn)
}

func (ss *serverSession) nextNum(update bool) int {
	ss.lastNum++
	// update tran# are odd, read-only are even
	if ((ss.lastNum % 2) == 1) != update {
		ss.lastNum++
	}
	return ss.lastNum
}

func cmdTransactions(ss *serverSession) {
	list := make([]int, 0, len(ss.trans))
	for tn := range ss.trans {
		list = append(list, tn)
	}
	ss.PutBool(true).PutInts(list)
}

func cmdUpdate(ss *serverSession) {
	tran := ss.getTran()
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
	cmdDump,
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
	cmdLoad,
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
}
