// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type suSocketClient struct {
	CantConvert
	conn    *net.TCPConn
	rdr     *bufio.Reader
	timeout time.Duration
}

var nSocketClient = 0

var _ = builtin("SocketClient(ipaddress, port, timeout=60, timeoutConnect=0, block=false)",
	func(t *Thread, args []Value) Value {
		ipaddr := ToStr(args[0])
		port := ToInt(args[1])
		ipaddr += ":" + strconv.Itoa(port)
		var c net.Conn
		var e error
		toc := time.Duration(ToInt(OpMul(args[3], SuInt(1000)))) * 1000 * 1000
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
		nSocketClient++
		if args[4] == False {
			return sc
		}
		// block form
		defer sc.close()
		return t.Call(args[4], sc)
	})

func (*suSocketClient) Get(*Thread, Value) Value {
	panic("SocketClient does not support get")
}

func (*suSocketClient) Put(*Thread, Value, Value) {
	panic("SocketClient does not support put")
}

func (*suSocketClient) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("SocketClient does not support update")
}

func (*suSocketClient) RangeTo(int, int) Value {
	panic("SocketClient does not support range")
}

func (*suSocketClient) RangeLen(int, int) Value {
	panic("SocketClient does not support range")
}

func (*suSocketClient) Hash() uint32 {
	panic("SocketClient hash not implemented")
}

func (*suSocketClient) Hash2() uint32 {
	panic("SocketClient hash not implemented")
}

func (*suSocketClient) Compare(Value) int {
	panic("SocketClient compare not implemented")
}

func (*suSocketClient) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call a SocketClient instance")
}

func (*suSocketClient) String() string {
	return "SocketClient"
}

func (*suSocketClient) Type() types.Type {
	return types.BuiltinClass
}

func (sc *suSocketClient) Equal(other interface{}) bool {
	sc2, ok := other.(*suSocketClient)
	return ok && sc == sc2
}

func (*suSocketClient) Lookup(_ *Thread, method string) Callable {
	return suSocketClientMethods[method]
}

var crnl = []byte{'\r', '\n'}

var noDeadline time.Time

var suSocketClientMethods = Methods{
	"Close": method0(func(this Value) Value {
		scOpen(this).close()
		return nil
	}),
	"Read": method1("(n)", func(this, arg Value) Value {
		sc := scOpen(this)
		n := ToInt(arg)
		buf := make([]byte, n)
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
		n, e := io.ReadFull(sc.rdr, buf)
		if e != nil && e != io.ErrUnexpectedEOF {
			panic("socketClient.Read: " + e.Error())
		}
		return SuStr(string(buf[:n]))
	}),
	"Readline": method0(func(this Value) Value {
		sc := scOpen(this)
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
		line := Readline(sc.rdr, "socket.Readline: ")
		if line == False {
			panic("socket Readline lost connection or timeout")
		}
		return line
	}),
	"Write": method1("(string)", func(this, arg Value) Value {
		sc := scOpen(this)
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
		s := AsStr(arg)
		_, e := io.WriteString(sc.conn, s)
		if e != nil {
			panic("socketClient.Write: " + e.Error())
		}
		return nil
	}),
	"Writeline": method1("(string)", func(this, arg Value) Value {
		sc := scOpen(this)
		sc.conn.SetDeadline(time.Now().Add(sc.timeout))
		defer sc.conn.SetDeadline(noDeadline)
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
	}),
}

func (sc *suSocketClient) close() {
	if sc.conn == nil {
		return
	}
	nSocketClient--
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
