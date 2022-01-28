// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sync"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/dbms/commands"
	"github.com/apmckinlay/gsuneido/dbms/csio"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/str"
)

// serverConn is one connection to the server
type serverConn struct {
	ended bool
	// id is primarily used as a key to store the set of connections in a map
	id   int
	dbms *DbmsLocal
	conn net.Conn
	csio.ReadWrite
	nonce     string
	trans     map[int]ITran
	cursors   map[int]ICursor
	queries   map[int]IQuery
	lastNum   int
	thread    Thread
}

var serverConns = make(map[int]*serverConn)
var serverConnsLock sync.Mutex

func Server(dbms *DbmsLocal) {
	l, err := net.Listen("tcp", ":"+options.Port)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	lastId := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		lastId++
		go handler(lastId, dbms, conn)
	}
}

func handler(id int, dbms *DbmsLocal, conn net.Conn) {
	trace.Dbms.Println("connected")
	sc := &serverConn{id: id, dbms: dbms, conn: conn,
		trans:   make(map[int]ITran),
		cursors: make(map[int]ICursor),
		queries: make(map[int]IQuery)}
	sc.ReadWrite = *csio.NewReadWrite(conn, sc.error)
	serverConnsLock.Lock()
	serverConns[sc.id] = sc
	serverConnsLock.Unlock()
	defer func() {
		serverConnsLock.Lock()
		delete(serverConns, sc.id)
		serverConnsLock.Unlock()
	}()
	sc.serve()
}

func (sc *serverConn) serve() {
	sc.sendHello()
	addr := sc.conn.RemoteAddr().String()
	sc.thread.SetSession(str.BeforeLast(addr, ":"))
	for !sc.ended {
		sc.request()
	}
}

func (sc *serverConn) sendHello() {
	sc.conn.Write(hello)
}

var hello = func() []byte {
	buf := make([]byte, helloSize)
	copy(buf, "Suneido "+builtin.Built()+"\r\n")
	return buf
}()

func (sc *serverConn) request() {
	defer func() {
		if e := recover(); e != nil && !sc.ended {
			debug.PrintStack()
			sc.ResetWrite(sc.conn)
			sc.PutBool(false).PutStr(fmt.Sprint(e))
		}
	}()
	icmd := commands.Command(sc.GetByte())
	trace.ClientServer.Println("<<<", icmd)
	if int(icmd) >= len(cmds) {
		sc.close()
		log.Println("dbms server, closing connection: invalid command")
	}
	cmd := cmds[icmd]
	cmd(sc)
	sc.Flush()
}

func (sc *serverConn) error(err string) {
	sc.close()
	log.Panicln("dbms server, closing connection:", err)
}

func (sc *serverConn) close() {
	sc.ended = true
	sc.conn.Close()
	for _, tran := range sc.trans {
		tran.Abort()
	}
}

//-------------------------------------------------------------------

func cmdAbort(sc *serverConn) {
	tn := sc.GetInt()
	tran := sc.tran(tn)
	tran.Abort()
	delete(sc.trans, tn)
	sc.PutBool(true)
}

func (sc *serverConn) getTran() ITran {
	tn := sc.GetInt()
	if tn == 0 {
		return nil
	}
	return sc.tran(tn)
}

func (sc *serverConn) tran(tn int) ITran {
	tran, ok := sc.trans[tn]
	if !ok {
		sc.error("transaction not found")
	}
	return tran
}

func cmdAction(sc *serverConn) {
	tran := sc.getTran()
	action := sc.GetStr()
	n := tran.Action(&sc.thread, action)
	sc.PutBool(true).PutInt(n)
}

func cmdAdmin(sc *serverConn) {
	s := sc.GetStr()
	sc.dbms.Admin(s)
	sc.PutBool(true)
}

func cmdAuth(sc *serverConn) {
	s := sc.GetStr()
	result := sc.auth(s)
	sc.PutBool(true).PutBool(result)
}

func (sc *serverConn) auth(s string) bool {
	if AuthUser(&sc.thread, s, sc.nonce) {
		sc.nonce = ""
		return true
	}
	return AuthToken(s)
}

func cmdCheck(sc *serverConn) {
	s := sc.dbms.Check()
	sc.PutBool(true).PutStr(s)
}

func cmdClose(sc *serverConn) {
	n := sc.GetInt()
	if sc.GetByte() == 'q' {
		q := sc.queries[n]
		if q == nil {
			sc.error("query not found")
		}
		delete(sc.queries, n)
		q.Close()
	} else {
		c := sc.cursors[n]
		if c == nil {
			sc.error("cursor not found")
		}
		delete(sc.cursors, n)
		c.Close()
	}
}

func cmdCommit(sc *serverConn) {
	tn := sc.GetInt()
	tran := sc.tran(tn)
	result := tran.Complete()
	delete(sc.trans, tn)
	sc.PutBool(true)
	if result == "" {
		sc.PutBool(true)
	} else {
		sc.PutBool(false).PutStr(result)
	}
}

func cmdConnections(sc *serverConn) {
	sc.PutBool(true).PutVal(connections())
}

func connections() *SuObject {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	list := &SuObject{}
	for _, sc := range serverConns {
		list.Add(SuStr(sc.thread.Session()))
	}
	return list
}

func cmdCursor(sc *serverConn) {
	query := sc.GetStr()
	q := sc.dbms.Cursor(query)
	sc.lastNum++
	sc.cursors[sc.lastNum] = q
	sc.PutBool(true).PutInt(sc.lastNum)
}

func cmdCursors(sc *serverConn) {
	sc.PutBool(true).PutInt(len(sc.cursors))
}

func cmdDelete(sc *serverConn) {
	tran := sc.getTran()
	table := sc.GetStr()
	off := uint64(sc.GetInt64())
	tran.Delete(&sc.thread, table, off)
	sc.PutBool(true)
}

func cmdDump(sc *serverConn) {
	table := sc.GetStr()
	s := sc.dbms.Dump(table)
	sc.PutBool(true).PutStr(s)
}

func cmdErase(sc *serverConn) {
	sc.error("gSuneido server requires gSuneido client")
}

func cmdExec(sc *serverConn) {
	ob := sc.GetVal()
	v := sc.dbms.Exec(&sc.thread, ob)
	sc.PutBool(true).PutVal(v)
}

func cmdFinal(sc *serverConn) {
	final := sc.dbms.Final()
	sc.PutBool(true).PutInt(final)
}

func cmdGet(sc *serverConn) {
	sc.error("gSuneido server requires gSuneido client")
}

func cmdGet2(sc *serverConn) {
	tbl, hdr, row := sc.getQorTC()
	sc.rowResult(tbl, hdr, false, row)
}

func (sc *serverConn) getQorTC() (tbl string, hdr *Header, row Row) {
	dir := sc.getDir()
	tn := sc.GetInt()
	qn := sc.GetInt()
	if tn == 0 {
		q := sc.queries[qn]
		hdr = q.Header()
		row, tbl = q.Get(&sc.thread, dir)

	} else {
		t := sc.trans[tn]
		c := sc.cursors[qn]
		hdr = c.Header()
		row, tbl = c.Get(&sc.thread, t, dir)
	}
	return
}

func (sc *serverConn) getDir() Dir {
	b := sc.GetByte()
	if b == '-' {
		return Prev
	}
	return Next
}

const maxRec = 1024 * 1024 // 1 mb

func (sc *serverConn) rowResult(tbl string, hdr *Header, sendHdr bool, row Row) {
	if row == nil {
		sc.PutBool(true).PutBool(false)
	} else {
		rec := rowToRecord(row, hdr)
		if len(rec) > maxRec {
			panic("result too large")
		}
		off := int64(row[0].Off)
		sc.PutBool(true).PutBool(true).PutInt64(off)
		if sendHdr {
			sc.PutStrs(hdr.Schema())
		}
		if tbl != "-" {
			sc.PutStr(tbl)
		}
		sc.PutRec(rec)
	}
}

func rowToRecord(row Row, hdr *Header) Record {
	if len(row) == 1 {
		return row[0].Record
	}
	var rb RecordBuilder
	for _, flds := range hdr.Fields {
		for _, fld := range flds {
			rb.AddRaw(row.GetRaw(hdr, fld))
		}
	}
	return rb.Trim().Build()
}

func cmdGetOne(sc *serverConn) {
	sc.error("gSuneido server requires gSuneido client")
}

func cmdGetOne2(sc *serverConn) {
	var dir Dir
	switch sc.GetByte() {
	case '+':
		dir = Next
	case '-':
		dir = Prev
	case '1':
		dir = Only
	default:
		sc.error("dbms server: expected + - 1")
	}
	tran := sc.getTran()
	query := sc.GetStr()
	var g func(*Thread, string, Dir) (Row, *Header, string)
	if tran == nil {
		g = sc.dbms.Get
	} else {
		g = tran.Get
	}
	row, hdr, tbl := g(&sc.thread, query, dir)
	sc.rowResult(tbl, hdr, true, row)
}

func cmdHeader(sc *serverConn) {
	hdr := sc.getQorC().Header()
	sc.PutBool(true).PutStrs(hdr.Schema())
}

func (sc *serverConn) getQorC() (qc IQueryCursor) {
	n := sc.GetInt()
	switch sc.GetByte() {
	case 'q':
		qc = sc.queries[n]
	case 'c':
		qc = sc.cursors[n]
	default:
		sc.error("dbms server expected q or c")
	}
	if qc == nil {
		sc.error("dbms server: query/cursor not found")
	}
	return qc
}

func cmdInfo(sc *serverConn) {
	info := sc.dbms.Info()
	sc.PutBool(true).PutVal(info)
}

func cmdKeys(sc *serverConn) {
	keys := sc.getQorC().Keys()
	sc.PutBool(true).PutStrs(keys)
}

func cmdKill(sc *serverConn) {
	sessionId := sc.GetStr()
	n := kill(sessionId)
	sc.PutBool(true).PutInt(n)
}

func kill(sessionId string) int {
	serverConnsLock.Lock()
	defer serverConnsLock.Unlock()
	nkilled := 0
	for id, sc := range serverConns {
		if sc.thread.Session() == sessionId {
			delete(serverConns, id)
			sc.conn.Close()
			nkilled++
		}
	}
	return nkilled
}

func cmdLibGet(sc *serverConn) {
	name := sc.GetStr()
	defs := sc.dbms.LibGet(name)
	sc.PutBool(true).PutInt(len(defs) / 2)
	for i := 0; i < len(defs); i += 2 {
		sc.PutStr(defs[i]).PutInt(len(defs[i+1]))
	}
	for i := 1; i < len(defs); i += 2 {
		sc.PutBuf(defs[i])
	}
}

func cmdLibraries(sc *serverConn) {
	libs := sc.dbms.Libraries()
	sc.PutBool(true).PutStrs(libs)
}

func cmdLoad(sc *serverConn) {
	table := sc.GetStr()
	n := sc.dbms.Load(table)
	sc.PutBool(true).PutInt(n)
}

func cmdLog(sc *serverConn) {
	s := sc.GetStr()
	sc.dbms.Log(s)
	sc.PutBool(true)
}

func cmdNonce(sc *serverConn) {
	sc.nonce = Nonce()
	sc.PutBool(true).PutStr(sc.nonce)
}

func cmdOrder(sc *serverConn) {
	order := sc.getQorC().Order()
	sc.PutBool(true).PutStrs(order)
}

func cmdOutput(sc *serverConn) {
	q := sc.getQuery()
	rec := sc.GetRec()
	q.Output(&sc.thread, rec)
	sc.PutBool(true)
}

func (sc *serverConn) getQuery() IQuery {
	qn := sc.GetInt()
	q := sc.queries[qn]
	if q == nil {
		sc.error("dbms server: query not found")
	}
	return q
}

func cmdQuery(sc *serverConn) {
	tran := sc.getTran()
	query := sc.GetStr()
	q := tran.Query(query)
	sc.lastNum++
	sc.queries[sc.lastNum] = q
	sc.PutBool(true).PutInt(sc.lastNum)
}

func cmdReadCount(sc *serverConn) {
	sc.getTran()
	sc.PutBool(true).PutInt(0) //TODO
}

func cmdRewind(sc *serverConn) {
	qc := sc.getQorC()
	qc.Rewind()
	sc.PutBool(true)
}

func cmdRun(sc *serverConn) {
	s := sc.GetStr()
	v := sc.dbms.Run(&sc.thread, s)
	sc.PutBool(true).PutVal(v)
}

func cmdSessionId(sc *serverConn) {
	s := sc.GetStr()
	if s != "" {
		sc.thread.SetSession(s)
	}
	sc.PutBool(true).PutStr(sc.thread.Session())
}

func cmdSize(sc *serverConn) {
	n := sc.dbms.Size()
	sc.PutBool(true).PutInt64(int64(n))
}

func cmdStrategy(sc *serverConn) {
	qc := sc.getQorC()
	strategy := qc.Strategy()
	sc.PutBool(true).PutStr(strategy)
}

func cmdTimestamp(sc *serverConn) {
	ts := sc.dbms.Timestamp()
	sc.PutBool(true).PutVal(ts)
}

func cmdToken(sc *serverConn) {
	tok := Token()
	sc.PutBool(true).PutStr(tok)
}

func cmdTransaction(sc *serverConn) {
	update := sc.GetBool()
	tran := sc.dbms.Transaction(update)
	tn := sc.nextNum(update)
	sc.trans[tn] = tran
	sc.PutBool(true).PutInt(tn)
}

func (sc *serverConn) nextNum(update bool) int {
	sc.lastNum++
	// update tran# are odd, read-only are even
	if ((sc.lastNum % 2) == 1) != update {
		sc.lastNum++
	}
	return sc.lastNum
}

func cmdTransactions(sc *serverConn) {
	list := make([]int, 0, len(sc.trans))
	for tn := range sc.trans {
		list = append(list, tn)
	}
	sc.PutBool(true).PutInts(list)
}

func cmdUpdate(sc *serverConn) {
	sc.error("gSuneido server does not work with jSuneido client")
}

func cmdUpdate2(sc *serverConn) {
	tran := sc.getTran()
	table := sc.GetStr()
	off := uint64(sc.GetInt64())
	rec := sc.GetRec()
	tran.Update(&sc.thread, table, off, rec)
	sc.PutBool(true)
}

func cmdWriteCount(sc *serverConn) {
	sc.getTran()
	sc.PutBool(true).PutInt(0) //TODO
}

type command func(sc *serverConn)

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
	cmdDelete,
	cmdUpdate2,
	cmdGet2,
	cmdGetOne2,
}
