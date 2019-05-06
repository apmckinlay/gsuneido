package dbms

import (
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/database/dbms/commands"
	"github.com/apmckinlay/gsuneido/database/dbms/csio"
	. "github.com/apmckinlay/gsuneido/runtime"
)

type DbmsClient struct {
	*csio.ReadWrite
	conn net.Conn
}

// helloSize is the size of the initial connection message from the server
// the size must match cSuneido and jSuneido
const helloSize = 50

func NewDbmsClient(addr string) *DbmsClient {
	conn, err := net.Dial("tcp", addr)
	if err != nil || !checkHello(conn) {
		panic("can't connect to " + addr + " " + err.Error())
	}
	return &DbmsClient{ReadWrite: csio.NewReadWrite(conn), conn: conn}
}

func checkHello(conn net.Conn) bool {
	var buf [helloSize]byte
	n, err := io.ReadFull(conn, buf[:])
	if n != helloSize || err != nil {
		return false
	}
	s := string(buf[:])
	if !strings.HasPrefix(s, "Suneido ") {
		return false
	}
	//TODO built date check
	return true
}

// Dbms interface

var _ IDbms = (*DbmsClient)(nil)

func (dc *DbmsClient) Admin(request string) {
	dc.PutCmd(commands.Admin).PutStr(request).Request()
}

func (dc *DbmsClient) Auth(s string) bool {
	if s == "" {
		return false
	}
	dc.PutCmd(commands.Auth).PutStr(s).Request()
	return dc.GetBool()
}

func (dc *DbmsClient) Check() string {
	dc.PutCmd(commands.Check).Request()
	return dc.GetStr()
}

func (dc *DbmsClient) Connections() Value {
	dc.PutCmd(commands.Connections).Request()
	ob := dc.GetVal().(*SuObject)
	ob.SetReadOnly()
	return ob
}

func (dc *DbmsClient) Cursors() int {
	dc.PutCmd(commands.Cursors).Request()
	return int(dc.GetInt())
}

func (dc *DbmsClient) Dump(table string) string {
	dc.PutCmd(commands.Dump).PutStr(table).Request()
	return dc.GetStr()
}

func (dc *DbmsClient) Exec(_ *Thread, args Value) Value {
	dc.PutCmd(commands.Exec).PutVal(args).Request()
	return dc.ValueResult()
}

func (dc *DbmsClient) Final() int {
	dc.PutCmd(commands.Final).Request()
	return int(dc.GetInt())
}

func (dc *DbmsClient) Info() Value {
	dc.PutCmd(commands.Info).Request()
	return dc.GetVal()
}

func (dc *DbmsClient) Kill(sessionid string) int {
	dc.PutCmd(commands.Kill).PutStr(sessionid).Request()
	return int(dc.GetInt())
}

func (dc *DbmsClient) Load(table string) int {
	dc.PutCmd(commands.Load).PutStr(table).Request()
	return int(dc.GetInt())
}

func (dc *DbmsClient) Log(s string) {
	dc.PutCmd(commands.Log).PutStr(s).Request()
}

func (dc *DbmsClient) LibGet(name string) []string {
	dc.PutCmd(commands.LibGet).PutStr(name).Request()
	n := dc.GetSize()
	v := make([]string, 2*n)
	sizes := make([]int, n)
	for i := 0; i < 2*n; i += 2 {
		v[i] = dc.GetStr() // library
		sizes[i/2] = dc.GetSize()
	}
	for i := 1; i < 2*n; i += 2 {
		v[i] = string(dc.Get(sizes[i/2])) // text
	}
	return v
}

func (dc *DbmsClient) Libraries() *SuObject {
	dc.PutCmd(commands.Libraries).Request()
	n := dc.GetInt()
	ob := NewSuObject()
	for ; n > 0; n-- {
		ob.Add(SuStr(dc.GetStr()))
	}
	return ob
}

func (dc *DbmsClient) Nonce() string {
	dc.PutCmd(commands.Nonce).Request()
	return dc.GetStr()
}

func (dc *DbmsClient) Run(code string) Value {
	dc.PutCmd(commands.Run).PutStr(code).Request()
	return dc.ValueResult()
}

func (dc *DbmsClient) SessionId(id string) string {
	dc.PutCmd(commands.SessionId).PutStr(id).Request()
	return dc.GetStr()
}

func (dc *DbmsClient) Size() int64 {
	dc.PutCmd(commands.Size).Request()
	return dc.GetInt()
}

func (dc *DbmsClient) Timestamp() SuDate {
	dc.PutCmd(commands.Timestamp).Request()
	return dc.GetVal().(SuDate)
}

func (dc *DbmsClient) Token() string {
	dc.PutCmd(commands.Token).Request()
	return dc.GetStr()
}

func (dc *DbmsClient) Transaction(update bool) ITran {
	dc.PutCmd(commands.Transaction).PutBool(update).Request()
	tn := int(dc.GetInt())
	return &TranClient{dc: dc, tn: tn, readonly: !update}
}

func (dc *DbmsClient) Transactions() *SuObject {
	dc.PutCmd(commands.Transactions).Request()
	ob := NewSuObject()
	for n := dc.GetInt(); n > 0; n-- {
		ob.Add(IntVal(int(dc.GetInt())))
	}
	return ob
}

func (dc *DbmsClient) Unuse(lib string) bool {
	panic("can't Unuse('" + lib + "')\n" +
		"When client-server, only the server can Unuse")
}

func (dc *DbmsClient) Use(lib string) bool {
	if _, ok := dc.Libraries().Find(SuStr(lib)); ok {
		return false
	}
	panic("can't Use('" + lib + "')\n" +
		"When client-server, only the server can Use")
}

func (dc *DbmsClient) Close() {
	dc.conn.Close()
}

// ------------------------------------------------------------------

type TranClient struct {
	dc       *DbmsClient
	tn       int
}

var _ ITran = (*TranClient)(nil)

func (tc *TranClient) Abort() {
	tc.dc.PutCmd(commands.Abort).PutInt(int64(tc.tn)).Request()
}

func (tc *TranClient) Complete() string {
	tc.dc.PutCmd(commands.Commit).PutInt(int64(tc.tn)).Request()
	if tc.dc.GetBool() {
		return ""
	}
	return tc.dc.GetStr()
}

func (tc *TranClient) String() string {
	return "Transaction" + strconv.Itoa(tc.tn)
}
