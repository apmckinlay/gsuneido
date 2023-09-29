// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"net"
	"strconv"
	"strings"

	"slices"

	"github.com/apmckinlay/gsuneido/dbms/commands"
	"github.com/apmckinlay/gsuneido/dbms/mux"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// This is the mux client that matches DbmsServer.
// See jsunClient for the version that matches jSuneido.

type dbmsClient struct {
	cc *mux.ClientConn
}

func NewDbmsClient(conn net.Conn) *dbmsClient {
	conn.Write(hello())
	cc := mux.NewClientConn(conn)
	return &dbmsClient{cc: cc}
}

type muxSession struct {
	*mux.ClientSession
}

func (dc *dbmsClient) NewSession() *muxSession {
	cs := dc.cc.NewClientSession()
	return &muxSession{ClientSession: cs}
}

// Dbms interface

var _ IDbms = (*muxSession)(nil)

func (ms *muxSession) Admin(admin string, _ *Sviews) {
	ms.PutCmd(commands.Admin).PutStr(admin)
	ms.Request()
}

func (ms *muxSession) Auth(_ *Thread, s string) bool {
	return ms.auth(s)
}

func (ms *muxSession) auth(s string) bool {
	if s == "" {
		return false
	}
	ms.PutCmd(commands.Auth).PutStr(s)
	ms.Request()
	return ms.GetBool()
}

func (ms *muxSession) Check() string {
	ms.PutCmd(commands.Check)
	ms.Request()
	return ms.GetStr()
}

func (ms *muxSession) Close() {
	ms.PutCmd(commands.EndSession)
	ms.EndMsg()
}

func (ms *muxSession) Connections() Value {
	ms.PutCmd(commands.Connections)
	ms.Request()
	ob := ms.GetVal().(*SuObject)
	ob.SetReadOnly()
	return ob
}

func (ms *muxSession) Cursor(query string, _ *Sviews) ICursor {
	ms.PutCmd(commands.Cursor).PutStr(query)
	ms.Request()
	cn := ms.GetInt()
	return ms.newClientCursor(cn)
}

func (ms *muxSession) Cursors() int {
	ms.PutCmd(commands.Cursors)
	ms.Request()
	return ms.GetInt()
}

func (ms *muxSession) DisableTrigger(string) {
	panic("DoWithoutTriggers can't be used by a client")
}
func (ms *muxSession) EnableTrigger(string) {
	assert.ShouldNotReachHere()
}

func (ms *muxSession) Dump(table string) string {
	ms.PutCmd(commands.Dump).PutStr(table)
	ms.Request()
	return ms.GetStr()
}

func (ms *muxSession) Exec(_ *Thread, args Value) Value {
	packed := PackValue(args) // do this first because it could panic
	if trace.ClientServer.On() {
		if len(packed) < 100 {
			trace.ClientServer.Println("    ->", args)
		}
	}
	ms.PutCmd(commands.Exec).PutStr_(packed)
	ms.Request()
	return ms.ValueResult()
}

func (ms *muxSession) Final() int {
	ms.PutCmd(commands.Final)
	ms.Request()
	return ms.GetInt()
}

func (ms *muxSession) Get(_ *Thread, query string, dir Dir,
	_ *Sviews) (Row, *Header, string) {
	return ms.get(0, query, dir)
}

func (ms *muxSession) get(tn int, query string, dir Dir) (Row, *Header, string) {
	ms.PutCmd(commands.GetOne).PutByte(byte(dir)).PutInt(tn).PutStr(query)
	ms.Request()
	if !ms.GetBool() {
		return nil, nil, ""
	}
	off := ms.GetInt()
	hdr := ms.getHdr()
	tbl := ms.GetStr()
	row := ms.getRow(off)
	return row, hdr, tbl
}

func (ms *muxSession) Info() Value {
	ms.PutCmd(commands.Info)
	ms.Request()
	return ms.GetVal()
}

func (ms *muxSession) Kill(sessionid string) int {
	ms.PutCmd(commands.Kill).PutStr(sessionid)
	ms.Request()
	return ms.GetInt()
}

func (ms *muxSession) Load(table string) int {
	ms.PutCmd(commands.Load).PutStr(table)
	ms.Request()
	return ms.GetInt()
}

func (ms *muxSession) Log(s string) {
	ms.PutCmd(commands.Log).PutStr(s)
	ms.Request()
}

func (ms *muxSession) LibGet(name string) []string {
	ms.PutCmd(commands.LibGet).PutStr(name)
	ms.Request()
	n := ms.GetSize()
	v := make([]string, 2*n)
	sizes := make([]int, n)
	for i := 0; i < 2*n; i += 2 {
		v[i] = ms.GetStr() // library
		sizes[i/2] = ms.GetSize()
	}
	for i := 1; i < 2*n; i += 2 {
		v[i] = ms.GetN(sizes[i/2]) // text
	}
	return v
}

func (ms *muxSession) Libraries() []string {
	ms.PutCmd(commands.Libraries)
	ms.Request()
	return ms.GetStrs()
}

func (ms *muxSession) Nonce(*Thread) string {
	ms.PutCmd(commands.Nonce)
	ms.Request()
	return ms.GetStr_()
}

func (ms *muxSession) Run(_ *Thread, code string) Value {
	ms.PutCmd(commands.Run).PutStr(code)
	ms.Request()
	return ms.ValueResult()
}

func (ms *muxSession) Schema(string) string {
	panic("Schema only available standalone")
}

func (ms *muxSession) SessionId(th *Thread, id string) string {
	if s := th.Session(); s != "" && id == "" {
		return s // use cached value
	}
	ms.PutCmd(commands.SessionId).PutStr(id)
	ms.Request()
	s := ms.GetStr()
	th.SetSession(s)
	return s
}

func (ms *muxSession) Size() uint64 {
	ms.PutCmd(commands.Size)
	ms.Request()
	return uint64(ms.GetInt64())
}

func (ms *muxSession) Timestamp() SuDate {
	ms.PutCmd(commands.Timestamp)
	ms.Request()
	return ms.GetVal().(SuDate)
}

func (ms *muxSession) Token() string {
	ms.PutCmd(commands.Token)
	ms.Request()
	return ms.GetStr()
}

func (ms *muxSession) Transaction(update bool) ITran {
	ms.PutCmd(commands.Transaction).PutBool(update)
	ms.Request()
	tn := ms.GetInt()
	return &muxTran{muxSession: ms, tn: tn}
}

func (ms *muxSession) Transactions() *SuObject {
	ms.PutCmd(commands.Transactions)
	ms.Request()
	ob := &SuObject{}
	for n := ms.GetInt(); n > 0; n-- {
		ob.Add(IntVal(ms.GetInt()))
	}
	return ob
}

func (ms *muxSession) Unuse(lib string) bool {
	panic("can't Unuse('" + lib + "')\n" +
		"When client-server, only the server can Unuse")
}

func (ms *muxSession) Use(lib string) bool {
	if slices.Contains(ms.Libraries(), lib) {
		return false
	}
	panic("can't Use('" + lib + "')\n" +
		"When client-server, only the server can Use")
}

func (ms *muxSession) Unwrap() IDbms {
	return ms
}

func (ms *muxSession) getHdr() *Header {
	n := ms.GetInt()
	fields := make([]string, 0, n)
	columns := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s := ms.GetStr()
		if ascii.IsUpper(s[0]) {
			s = str.UnCapitalize(s)
		} else if !strings.HasSuffix(s, "_lower!") {
			fields = append(fields, s)
		}
		if s != "-" {
			columns = append(columns, s)
		}
	}
	return NewHeader([][]string{fields}, columns)
}

func (ms *muxSession) getRow(off int) Row {
	return Row([]DbRec{{Record: ms.GetRec(), Off: uint64(off)}})
}

// ------------------------------------------------------------------

type muxTran struct {
	*muxSession
	conflict string
	tn       int
	ended    bool
}

var _ ITran = (*muxTran)(nil)

func (tc *muxTran) Abort() string {
	tc.ended = true
	tc.PutCmd(commands.Abort).PutInt(tc.tn)
	tc.Request()
	return ""
}

func (tc *muxTran) Asof(asof int64) int64 {
	tc.PutCmd(commands.Asof).PutInt(tc.tn).PutInt64(asof)
	tc.Request()
	return tc.GetInt64()
}

func (tc *muxTran) Complete() string {
	tc.ended = true
	tc.PutCmd(commands.Commit).PutInt(tc.tn)
	tc.Request()
	if tc.GetBool() {
		return ""
	}
	tc.conflict = tc.GetStr()
	return tc.conflict
}

func (tc *muxTran) Conflict() string {
	return tc.conflict
}

func (tc *muxTran) Ended() bool {
	return tc.ended
}

func (tc *muxTran) Delete(_ *Thread, table string, off uint64) {
	tc.PutCmd(commands.Erase).PutInt(tc.tn).PutStr(table).PutInt(int(off))
	tc.Request()
}

func (tc *muxTran) Get(_ *Thread, query string, dir Dir,
	_ *Sviews) (Row, *Header, string) {
	return tc.get(tc.tn, query, dir)
}

func (tc *muxTran) Query(query string, _ *Sviews) IQuery {
	tc.PutCmd(commands.Query).PutInt(tc.tn).PutStr(query)
	tc.Request()
	qn := tc.GetInt()
	return tc.muxSession.newClientQuery(qn)
}

func (tc *muxTran) ReadCount() int {
	tc.PutCmd(commands.ReadCount).PutInt(tc.tn)
	tc.Request()
	return tc.GetInt()
}

func (tc *muxTran) Action(_ *Thread, action string, _ *Sviews) int {
	tc.PutCmd(commands.Action).PutInt(tc.tn).PutStr(action)
	tc.Request()
	return tc.GetInt()
}

func (tc *muxTran) Update(_ *Thread, table string, off uint64, rec Record) uint64 {
	tc.PutCmd(commands.Update).
		PutInt(tc.tn).PutStr(table).PutInt(int(off)).PutRec(rec)
	tc.Request()
	return uint64(tc.GetInt())
}

func (tc *muxTran) WriteCount() int {
	tc.PutCmd(commands.WriteCount).PutInt(tc.tn)
	tc.Request()
	return tc.GetInt()
}

func (tc *muxTran) String() string {
	return "Transaction" + strconv.Itoa(tc.tn)
}

// ------------------------------------------------------------------

// muxQueryCursor is the common stuff for muxQuery and muxCursor
type muxQueryCursor struct {
	*muxSession
	hdr  *Header
	keys []string // cache
	id   int
	qc   qcType
}

func (qc *muxQueryCursor) Close() {
	qc.PutCmd(commands.Close).PutInt(qc.id).PutByte(byte(qc.qc))
	qc.Request()
}

func (qc *muxQueryCursor) Header() *Header {
	if qc.hdr == nil { // cached
		qc.PutCmd(commands.Header).PutInt(qc.id).PutByte(byte(qc.qc))
		qc.Request()
		qc.hdr = qc.getHdr()
	}
	return qc.hdr
}

func (qc *muxQueryCursor) Keys() []string {
	if qc.keys == nil { // cached
		qc.PutCmd(commands.Keys).PutInt(qc.id).PutByte(byte(qc.qc))
		qc.Request()
		qc.keys = qc.GetStrs()
	}
	return qc.keys
}

func (qc *muxQueryCursor) Order() []string {
	qc.PutCmd(commands.Order).PutInt(qc.id).PutByte(byte(qc.qc))
	qc.Request()
	return qc.GetStrs()
}

func (qc *muxQueryCursor) Rewind() {
	qc.PutCmd(commands.Rewind).PutInt(qc.id).PutByte(byte(qc.qc))
	qc.Request()
}

func (qc *muxQueryCursor) Strategy(formatted bool) string {
	qc.PutCmd(commands.Strategy).PutInt(qc.id).PutByte(byte(qc.qc)).PutBool(formatted)
	qc.Request()
	return qc.GetStr()
}

// muxQuery implements IQuery ------------------------------------
type muxQuery struct {
	muxQueryCursor
}

func (ms *muxSession) newClientQuery(qn int) *muxQuery {
	return &muxQuery{muxQueryCursor{muxSession: ms, id: qn, qc: query}}
}

var _ IQuery = (*muxQuery)(nil)

func (q *muxQuery) Get(_ *Thread, dir Dir) (Row, string) {
	q.PutCmd(commands.Get).PutByte(byte(dir)).PutInt(0).PutInt(q.id)
	q.Request()
	if !q.GetBool() {
		return nil, ""
	}
	off := q.GetInt()
	table := q.GetStr()
	row := q.getRow(off)
	return row, table
}

func (q *muxQuery) Output(_ *Thread, rec Record) {
	q.PutCmd(commands.Output).PutInt(q.id).PutRec(rec)
	q.Request()
}

// muxCursor implements IQuery ------------------------------------
type muxCursor struct {
	muxQueryCursor
}

func (ms *muxSession) newClientCursor(cn int) *muxCursor {
	return &muxCursor{muxQueryCursor{muxSession: ms, id: cn, qc: cursor}}
}

var _ ICursor = (*muxCursor)(nil)

func (q *muxCursor) Get(_ *Thread, tran ITran, dir Dir) (Row, string) {
	t := tran.(*muxTran)
	q.PutCmd(commands.Get).PutByte(byte(dir)).PutInt(t.tn).PutInt(q.id)
	q.Request()
	if !q.GetBool() {
		return nil, ""
	}
	off := q.GetInt()
	table := q.GetStr()
	row := q.getRow(off)
	return row, table
}
