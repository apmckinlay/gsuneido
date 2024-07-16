// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/str"
)

const helloSize = 50

var helloBuf [helloSize]byte
var helloOnce sync.Once

// hello returns the initial connection message.
// Both the client and the server send and receive/check this message.
func hello() []byte {
	helloOnce.Do(func() {
		copy(helloBuf[:], "Suneido "+options.BuiltStr()+"\r\n")
	})
	return helloBuf[:]
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
