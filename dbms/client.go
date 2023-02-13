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
	"time"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

func ConnectClient(addr string, port string) (conn net.Conn, jserver bool) {
	conn, err := net.Dial("tcp", addr+":"+port)
	if err != nil {
		checkServerStatus(addr, port)
		cantConnect(err.Error())
	}
	jserver, errmsg := checkHello(conn)
	if errmsg != "" {
		cantConnect(errmsg)
	}
	return conn, jserver
}

func cantConnect(s string) {
	Fatal("client: connect failed:", s)
}

const helloTimeout = 250 * time.Millisecond

// checkHello is used by both the client and the server
func checkHello(conn net.Conn) (jserver bool, errmsg string) {
	var buf [helloSize]byte
	conn.SetReadDeadline(time.Now().Add(helloTimeout))
	n, err := io.ReadFull(conn, buf[:])
	var never time.Time
	conn.SetReadDeadline(never)
	if n == 0 {
		return false, "hello: timeout"
	}
	if n != helloSize || err != nil {
		return false, "hello: invalid response"
	}
	s := string(buf[:])
	if !strings.HasPrefix(s, "Suneido ") {
		return false, "hello: invalid response"
	}
	if strings.Contains(s, " (Java)") {
		return true, ""
	}
	s = strings.TrimPrefix(s, "Suneido ")
	if noTime(s) != noTime(options.BuiltDate) && !options.IgnoreVersion {
		return false, "hello: version mismatch"
	}
	return false, ""
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
