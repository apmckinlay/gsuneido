package dbms

import (
	"io"
	"net"
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

func (dc *DbmsClient) Timestamp() SuDate {
	dc.PutCmd(commands.Timestamp).Request()
	return dc.GetVal().(SuDate)
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
