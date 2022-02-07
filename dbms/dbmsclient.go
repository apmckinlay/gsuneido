// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apmckinlay/gsuneido/dbms/commands"
	"github.com/apmckinlay/gsuneido/dbms/csio"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/trace"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/strs"
)

// token is to authorize the next connection
var token string

// tokenLock guards token
var tokenLock sync.Mutex

type dbmsClient struct {
	*csio.ReadWrite
	conn      net.Conn
	sessionId string
	jserver   bool
}

// helloSize is the size of the initial connection message from the server
// the size must match cSuneido and jSuneido
const helloSize = 50

func NewDbmsClient(addr string, port string) *dbmsClient {
	conn, err := net.Dial("tcp", addr+":"+port)
	if err != nil {
		checkServerStatus(addr, port)
		cantConnect(err.Error())
	}
	ok, jserver := checkHello(conn)
	if !ok {
		cantConnect("invalid response from server")
	}
	errfn := func(err string) {
		Fatal("client:", err)
	}
	rw := csio.NewReadWrite(conn, errfn)
	c := &dbmsClient{ReadWrite: rw, conn: conn, jserver: jserver}
	tokenLock.Lock()
	defer tokenLock.Unlock()
	if token != "" {
		c.auth(token)
		token = c.Token()
	}
	return c
}

func cantConnect(s string) {
	Fatal("Can't connect. " + s)
}

func checkHello(conn net.Conn) (ok, jserver bool) {
	var buf [helloSize]byte
	n, err := io.ReadFull(conn, buf[:])
	if n != helloSize || err != nil {
		return
	}
	s := string(buf[:])
	if !strings.HasPrefix(s, "Suneido ") {
		return
	}
	//TODO built date check
	if strings.Contains(s, "Java") {
		return true, true
	}
	return true, false
}

func checkServerStatus(addr string, port string) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return
	}
	url := "http://" + addr + ":" + strconv.Itoa(p+1) + "/"
	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 1024)
	io.ReadFull(resp.Body, buf)
	if bytes.Contains(buf, []byte("Checking database ...")) {
		cantConnect("Database is being checked, please try again later")
	}
	if bytes.Contains(buf, []byte("Rebuilding database ...")) {
		cantConnect("Database is being repaired, please try again later")
	}
}

// Dbms interface

var _ IDbms = (*dbmsClient)(nil)

func (dc *dbmsClient) Admin(admin string) {
	dc.PutCmd(commands.Admin).PutStr(admin).Request()
}

func (dc *dbmsClient) Auth(s string) bool {
	if !dc.auth(s) {
		return false
	}
	tokenLock.Lock()
	defer tokenLock.Unlock()
	if token == "" {
		token = dc.Token()
	}
	return true
}

func (dc *dbmsClient) auth(s string) bool {
	if s == "" {
		return false
	}
	dc.PutCmd(commands.Auth).PutStr(s).Request()
	return dc.GetBool()
}

func (dc *dbmsClient) Check() string {
	dc.PutCmd(commands.Check).Request()
	return dc.GetStr()
}

func (dc *dbmsClient) Close() {
	dc.conn.Close()
}

func (dc *dbmsClient) Connections() Value {
	dc.PutCmd(commands.Connections).Request()
	ob := dc.GetVal().(*SuObject)
	ob.SetReadOnly()
	return ob
}

func (dc *dbmsClient) Cursor(query string) ICursor {
	dc.PutCmd(commands.Cursor).PutStr(query).Request()
	cn := dc.GetInt()
	return newClientCursor(dc, cn)
}

func (dc *dbmsClient) Cursors() int {
	dc.PutCmd(commands.Cursors).Request()
	return dc.GetInt()
}

func (dc *dbmsClient) DisableTrigger(string) {
	panic("DoWithoutTriggers can't be used by a client")
}
func (dc *dbmsClient) EnableTrigger(string) {
	panic("shouldn't reach here")
}

func (dc *dbmsClient) Dump(table string) string {
	dc.PutCmd(commands.Dump).PutStr(table).Request()
	return dc.GetStr()
}

func (dc *dbmsClient) Exec(_ *Thread, args Value) Value {
	packed := PackValue(args) // do this first because it could panic
	if trace.ClientServer.On() {
		if len(packed) < 100 {
			trace.ClientServer.Println("    ->", args)
		}
	}
	dc.PutCmd(commands.Exec)
	dc.PutStr_(packed).Request()
	return dc.ValueResult()
}

func (dc *dbmsClient) Final() int {
	dc.PutCmd(commands.Final).Request()
	return dc.GetInt()
}

func (dc *dbmsClient) Get(_ *Thread, query string, dir Dir) (Row, *Header, string) {
	return dc.get(0, query, dir)
}

func (dc *dbmsClient) get(tn int, query string, dir Dir) (Row, *Header, string) {
	cmd := commands.GetOne2
	if dc.jserver {
		cmd = commands.GetOne
	}
	dc.PutCmd(cmd).PutByte(byte(dir)).PutInt(tn).PutStr(query).Request()
	if !dc.GetBool() {
		return nil, nil, ""
	}
	off := dc.GetInt()
	hdr := dc.getHdr()
	tbl := "updateable"
	if !dc.jserver {
		tbl = dc.GetStr()
	}
	row := dc.getRow(off)
	return row, hdr, tbl
}

func (dc *dbmsClient) Info() Value {
	dc.PutCmd(commands.Info).Request()
	return dc.GetVal()
}

func (dc *dbmsClient) Kill(sessionid string) int {
	dc.PutCmd(commands.Kill).PutStr(sessionid).Request()
	return dc.GetInt()
}

func (dc *dbmsClient) Load(table string) int {
	dc.PutCmd(commands.Load).PutStr(table).Request()
	return dc.GetInt()
}

func (dc *dbmsClient) Log(s string) {
	dc.PutCmd(commands.Log).PutStr(s).Request()
}

func (dc *dbmsClient) LibGet(name string) []string {
	dc.PutCmd(commands.LibGet).PutStr(name).Request()
	n := dc.GetSize()
	v := make([]string, 2*n)
	sizes := make([]int, n)
	for i := 0; i < 2*n; i += 2 {
		v[i] = dc.GetStr() // library
		sizes[i/2] = dc.GetSize()
	}
	for i := 1; i < 2*n; i += 2 {
		v[i] = dc.GetN(sizes[i/2]) // text
	}
	return v
}

func (dc *dbmsClient) Libraries() []string {
	dc.PutCmd(commands.Libraries).Request()
	return dc.GetStrs()
}

func (dc *dbmsClient) Nonce() string {
	dc.PutCmd(commands.Nonce).Request()
	return dc.GetStr()
}

func (dc *dbmsClient) Run(_ *Thread, code string) Value {
	dc.PutCmd(commands.Run).PutStr(code).Request()
	return dc.ValueResult()
}

func (dc *dbmsClient) Schema(string) string {
	panic("Schema only available standalone")
}

func (dc *dbmsClient) SessionId(t *Thread, id string) string {
	if s := t.Session(); s != "" && id == "" {
		return s // use cached value
	}
	dc.PutCmd(commands.SessionId).PutStr(id).Request()
	s := dc.GetStr()
	t.SetSession(s)
	return s
}

func (dc *dbmsClient) Size() uint64 {
	dc.PutCmd(commands.Size).Request()
	return uint64(dc.GetInt64())
}

func (dc *dbmsClient) Timestamp() SuDate {
	dc.PutCmd(commands.Timestamp).Request()
	return dc.GetVal().(SuDate)
}

func (dc *dbmsClient) Token() string {
	dc.PutCmd(commands.Token).Request()
	return dc.GetStr()
}

func (dc *dbmsClient) Transaction(update bool) ITran {
	dc.PutCmd(commands.Transaction).PutBool(update).Request()
	tn := dc.GetInt()
	return &TranClient{dc: dc, tn: tn}
}

func (dc *dbmsClient) Transactions() *SuObject {
	dc.PutCmd(commands.Transactions).Request()
	ob := &SuObject{}
	for n := dc.GetInt(); n > 0; n-- {
		ob.Add(IntVal(dc.GetInt()))
	}
	return ob
}

func (dc *dbmsClient) Unuse(lib string) bool {
	panic("can't Unuse('" + lib + "')\n" +
		"When client-server, only the server can Unuse")
}

func (dc *dbmsClient) Use(lib string) bool {
	if strs.Contains(dc.Libraries(), lib) {
		return false
	}
	panic("can't Use('" + lib + "')\n" +
		"When client-server, only the server can Use")
}

func (dc *dbmsClient) getHdr() *Header {
	n := dc.GetInt()
	fields := make([]string, 0, n)
	columns := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s := dc.GetStr()
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

func (dc *dbmsClient) getRow(off int) Row {
	return Row([]DbRec{{Record: dc.GetRec(), Off: uint64(off)}})
}

// ------------------------------------------------------------------

type TranClient struct {
	dc       *dbmsClient
	tn       int
	conflict string
	ended    bool
}

var _ ITran = (*TranClient)(nil)

func (tc *TranClient) Abort() string {
	tc.ended = true
	tc.dc.PutCmd(commands.Abort).PutInt(tc.tn).Request()
	return ""
}

func (tc *TranClient) Complete() string {
	tc.ended = true
	tc.dc.PutCmd(commands.Commit).PutInt(tc.tn).Request()
	if tc.dc.GetBool() {
		return ""
	}
	tc.conflict = tc.dc.GetStr()
	return tc.conflict
}

func (tc *TranClient) Conflict() string {
	return tc.conflict
}

func (tc *TranClient) Ended() bool {
	return tc.ended
}

func (tc *TranClient) Delete(_ *Thread, table string, off uint64) {
	if tc.dc.jserver {
		tc.dc.PutCmd(commands.Erase).PutInt(tc.tn).PutInt(int(off)).Request()
	} else {
		tc.dc.PutCmd(commands.Delete).
			PutInt(tc.tn).PutStr(table).PutInt(int(off)).Request()
	}
}

func (tc *TranClient) Get(_ *Thread, query string, dir Dir) (Row, *Header, string) {
	return tc.dc.get(tc.tn, query, dir)
}

func (tc *TranClient) Query(query string) IQuery {
	tc.dc.PutCmd(commands.Query).PutInt(tc.tn).PutStr(query).Request()
	qn := tc.dc.GetInt()
	return newClientQuery(tc.dc, qn)
}

func (tc *TranClient) ReadCount() int {
	tc.dc.PutCmd(commands.ReadCount).PutInt(tc.tn).Request()
	return tc.dc.GetInt()
}

func (tc *TranClient) Action(_ *Thread, action string) int {
	tc.dc.PutCmd(commands.Action).PutInt(tc.tn).PutStr(action).Request()
	return tc.dc.GetInt()
}

func (tc *TranClient) Update(_ *Thread, table string, off uint64, rec Record) uint64 {
	if tc.dc.jserver {
		tc.dc.PutCmd(commands.Update).
			PutInt(tc.tn).PutInt(int(off)).PutRec(rec).Request()
	} else {
		tc.dc.PutCmd(commands.Update2).
			PutInt(tc.tn).PutStr(table).PutInt(int(off)).PutRec(rec).Request()
	}
	return uint64(tc.dc.GetInt())
}

func (tc *TranClient) WriteCount() int {
	tc.dc.PutCmd(commands.WriteCount).PutInt(tc.tn).Request()
	return tc.dc.GetInt()
}

func (tc *TranClient) String() string {
	return "Transaction" + strconv.Itoa(tc.tn)
}

// ------------------------------------------------------------------

// clientQueryCursor is the common stuff for clientQuery and clientCursor
type clientQueryCursor struct {
	dc   *dbmsClient
	id   int
	qc   qcType
	hdr  *Header
	keys []string // cache
}

type qcType byte

const (
	query  qcType = 'q'
	cursor qcType = 'c'
)

func (qc *clientQueryCursor) Close() {
	qc.dc.PutCmd(commands.Close).PutInt(qc.id).PutByte(byte(qc.qc)).Request()
}

func (qc *clientQueryCursor) Header() *Header {
	if qc.hdr == nil { // cached
		qc.dc.PutCmd(commands.Header).PutInt(qc.id).PutByte(byte(qc.qc)).Request()
		qc.hdr = qc.dc.getHdr()
	}
	return qc.hdr
}

func (qc *clientQueryCursor) Keys() []string {
	if qc.keys == nil { // cached
		qc.dc.PutCmd(commands.Keys).PutInt(qc.id).PutByte(byte(qc.qc)).Request()
		if !qc.dc.jserver {
			qc.keys = qc.dc.GetStrs()
		} else {
			nk := qc.dc.GetInt()
			qc.keys = make([]string, nk)
			for i := range qc.keys {
				cb := str.CommaBuilder{}
				n := qc.dc.GetInt()
				for ; n > 0; n-- {
					cb.Add(qc.dc.GetStr())
				}
				qc.keys[i] = cb.String()
			}
		}
	}
	return qc.keys
}

func (qc *clientQueryCursor) Order() []string {
	qc.dc.PutCmd(commands.Order).PutInt(qc.id).PutByte(byte(qc.qc)).Request()
	return qc.dc.GetStrs()
}

func (qc *clientQueryCursor) Rewind() {
	qc.dc.PutCmd(commands.Rewind).PutInt(qc.id).PutByte(byte(qc.qc)).Request()
}

func (qc *clientQueryCursor) Strategy() string {
	qc.dc.PutCmd(commands.Strategy).PutInt(qc.id).PutByte(byte(qc.qc)).Request()
	return qc.dc.GetStr()
}

// clientQuery implements IQuery ------------------------------------
type clientQuery struct {
	clientQueryCursor
}

func newClientQuery(dc *dbmsClient, qn int) *clientQuery {
	return &clientQuery{clientQueryCursor{dc: dc, id: qn, qc: query}}
}

var _ IQuery = (*clientQuery)(nil)

func (q *clientQuery) Get(_ *Thread, dir Dir) (Row, string) {
	cmd := commands.Get2
	if q.dc.jserver {
		cmd = commands.Get
	}
	q.dc.PutCmd(cmd).PutByte(byte(dir)).PutInt(0).PutInt(q.id).Request()
	if !q.dc.GetBool() {
		return nil, ""
	}
	off := q.dc.GetInt()
	table := "updateable"
	if !q.dc.jserver {
		table = q.dc.GetStr()
	}
	row := q.dc.getRow(off)
	return row, table
}

func (q *clientQuery) Output(_ *Thread, rec Record) {
	q.dc.PutCmd(commands.Output).PutInt(q.id).PutRec(rec).Request()
}

// clientCursor implements IQuery ------------------------------------
type clientCursor struct {
	clientQueryCursor
}

func newClientCursor(dc *dbmsClient, cn int) *clientCursor {
	return &clientCursor{clientQueryCursor{dc: dc, id: cn, qc: cursor}}
}

var _ ICursor = (*clientCursor)(nil)

func (q *clientCursor) Get(_ *Thread, tran ITran, dir Dir) (Row, string) {
	cmd := commands.Get2
	if q.dc.jserver {
		cmd = commands.Get
	}
	t := tran.(*TranClient)
	q.dc.PutCmd(cmd).PutByte(byte(dir)).PutInt(t.tn).PutInt(q.id).Request()
	if !q.dc.GetBool() {
		return nil, ""
	}
	off := q.dc.GetInt()
	table := "updateable"
	if !q.dc.jserver {
		table = q.dc.GetStr()
	}
	row := q.dc.getRow(off)
	return row, table
}
