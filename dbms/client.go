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

	. "github.com/apmckinlay/gsuneido/core"
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

func cantConnect(s string) {
	Fatal("client: connect failed:", s)
}
