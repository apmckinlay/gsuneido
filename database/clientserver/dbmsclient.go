package clientserver

import (
	"io"
	"net"
	"strings"

	"github.com/apmckinlay/gsuneido/database/clientserver/csio"

	"github.com/apmckinlay/gsuneido/database/clientserver/commands"
)

type DbmsClient struct {
	*csio.ReadWrite
	conn net.Conn
}

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

var _ Dbms = (*DbmsClient)(nil)

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
