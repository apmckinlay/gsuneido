// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/str"
)

var VersionMismatch func(string) // injected by gsuneido.go

func ConnectClient(addr string, port string) net.Conn {
	conn, err := net.Dial("tcp", addr+":"+port)
	if err != nil {
		checkServerStatus(addr, port)
		cantConnect(err.Error())
	}
	conn.Write(hello())
	errmsg := checkHello(conn)
	if errmsg != "" {
		if strings.HasPrefix(errmsg, "version mismatch") {
			clientVersionMismatch(conn)
		}
		cantConnect(errmsg)
	}
	return conn
}

func clientVersionMismatch(conn net.Conn) {
	buf := make([]byte, 2)
	io.ReadFull(conn, buf)
	n := int(buf[0])<<8 | int(buf[1])
	if n == 0 {
		return
	}
	buf = make([]byte, n)
	io.ReadFull(conn, buf)
	VersionMismatch(string(buf))
}

func cantConnect(s string) {
	Fatal("client: connect failed:", s)
}

const helloTimeout = 500 * time.Millisecond

// checkHello is used by both the client and the server
func checkHello(conn net.Conn) string {
	var buf [helloSize]byte
	conn.SetReadDeadline(time.Now().Add(helloTimeout))
	n, err := io.ReadFull(conn, buf[:])
	var never time.Time
	conn.SetReadDeadline(never)
	if n == 0 {
		return "hello: timeout"
	}
	if n != helloSize || err != nil {
		return "hello: invalid response"
	}
	s := string(buf[:])
	if !strings.HasPrefix(s, "Suneido ") {
		return "hello: invalid response"
	}
	s = strings.TrimPrefix(s, "Suneido ")
	if noTime(s) != noTime(options.BuiltDate) && !options.IgnoreVersion {
		return fmt.Sprintf("version mismatch (got %s, want %s)",
			noTime(s), noTime(options.BuiltDate))
	}
	return ""
}

func noTime(s string) string {
	s = str.BeforeFirst(s, ":")
	return str.BeforeLast(s, " ")
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
	if bytes.Contains(buf, []byte("Rebuilding database ...")) ||
		bytes.Contains(buf, []byte("Repairing database ...")) {
		cantConnect("Database is being repaired, please try again later")
	}
}
