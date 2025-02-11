// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
)

type suSocketClient struct {
	ValueBase[*suSocketClient]
	conn    *net.TCPConn
	rdr     *bufio.Reader
	timeout time.Duration
}

var nSocketClient atomic.Int32
var _ = AddInfo("builtin.nSocketClient", &nSocketClient)

var _ = builtin(SocketClient,
	"(ipaddress, port, timeout=60, timeoutConnect=0, block=false)")

func SocketClient(th *Thread, args []Value) Value {
	ipaddr := ToStr(args[0])
	port := ToInt(args[1])
	ipaddr += ":" + strconv.Itoa(port)
	var c net.Conn
	var e error
	toc := time.Duration(ToInt(OpMul(args[3], SuInt(1000)))) * time.Millisecond
	if toc <= 0 {
		c, e = net.Dial("tcp", ipaddr)
	} else {
		c, e = net.DialTimeout("tcp", ipaddr, toc)
	}
	if e != nil {
		panic("SocketClient: " + e.Error())
	}
	sc := &suSocketClient{conn: c.(*net.TCPConn), rdr: bufio.NewReader(c),
		timeout: time.Duration(ToInt(args[2])) * time.Second}
	nSocketClient.Add(1)
	if args[4] == False {
		return sc
	}
	// block form
	defer sc.Close()
	return th.Call(args[4], sc)
}

func (sc *suSocketClient) Equal(other any) bool {
	return sc == other
}

func (*suSocketClient) SetConcurrent() {
	//FIXME concurrency
	// panic("SocketClient cannot be set to concurrent")
}

func (*suSocketClient) Lookup(_ *Thread, method string) Value {
	return suSocketClientMethods[method]
}

var crnl = []byte{'\r', '\n'}

var noDeadline time.Time

// socket client methods are also used by SocketServer

var suSocketClientMethods = methods("sock")

var _ = method(sock_Close, "()")

func sock_Close(this Value) Value {
	c := this.(interface{ Close() })
	c.Close()
	return nil
}

var _ = method(sock_Read, "(nbytes=false)")

func sock_Read(this, arg Value) Value {
	sc := scOpen(this)
	if sc.timeout > 0 {
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
	}
	return limitedRead("socket.Read", sc.rdr, arg)
}

var _ = method(sock_Readline, "()")

func sock_Readline(this Value) Value {
	sc := scOpen(this)
	if sc.timeout > 0 {
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
	}
	return Readline(sc.rdr, "socket.Readline: ")
}

var _ = method(sock_SetTimeout, "(seconds)")

func sock_SetTimeout(this, arg Value) Value {
	sc := scOpen(this)
	sc.timeout = time.Duration(ToInt(arg)) * time.Second
	return nil
}

var _ = method(sock_Write, "(string)")

func sock_Write(this, arg Value) Value {
	sc := scOpen(this)
	if sc.timeout > 0 {
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
	}
	s := AsStr(arg)
	_, e := io.WriteString(sc.conn, s)
	if e != nil {
		panic("socketClient.Write: " + e.Error())
	}
	return nil
}

var _ = method(sock_Writeline, "(string)")

func sock_Writeline(this, arg Value) Value {
	sc := scOpen(this)
	if sc.timeout > 0 {
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
	}
	s := AsStr(arg)
	_, e := io.WriteString(sc.conn, s)
	if e != nil {
		panic("socketClient.Writeline: " + e.Error())
	}
	_, e = sc.conn.Write(crnl)
	if e != nil {
		panic("socketClient.Writeline: " + e.Error())
	}
	return nil
}

var _ = method(sock_CopyTo, "(dest, nbytes = false)")

func sock_CopyTo(th *Thread, this Value, args []Value) Value {
	return CopyTo(th, scOpen(this).rdr, args[0], args[1])
}

func (sc *suSocketClient) writer() io.Writer {
	return scOpen(sc).conn
}

func (sc *suSocketClient) Close() {
	if sc.conn == nil {
		return
	}
	nSocketClient.Add(-1)
	sc.conn.Close()
	sc.conn = nil
}

func scOpen(this Value) *suSocketClient {
	sc, ok := this.(*suSocketClient)
	if !ok {
		sc = &this.(*suServerConnect).client
	}
	if sc.conn == nil {
		panic("can't use a closed SocketClient")
	}
	return sc
}
