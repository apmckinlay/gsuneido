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

var _ = builtin("SocketClient(ipaddress, port, timeout=60, timeoutConnect=0, block=false)",
	func(t *Thread, args []Value) Value {
		ipaddr := ToStr(args[0])
		port := ToInt(args[1])
		ipaddr += ":" + strconv.Itoa(port)
		var c net.Conn
		var e error
		toc := ToInt(args[3])
		const sec = 1000 * 1000 * 1000
		if toc == 0 {
			c, e = net.Dial("tcp", ipaddr)
		} else {
			c, e = net.DialTimeout("tcp", ipaddr, time.Duration(toc*sec))
		}
		if e != nil {
			panic("SocketClient: " + e.Error())
		}
		sc := &suSocketClient{conn: c.(*net.TCPConn), rdr: bufio.NewReader(c),
			timeout: time.Duration(ToInt(args[2]) * sec)}
		if args[4] == False {
			return sc
		}
		// block form
		defer sc.conn.Close()
		return t.Call(args[4], sc)
	})

type suSocketClient struct {
	CantConvert
	conn    *net.TCPConn
	rdr     *bufio.Reader
	timeout time.Duration
}

var _ Value = (*suSocketClient)(nil)

func (*suSocketClient) Get(*Thread, Value) Value {
	panic("SocketClient does not support get")
}

func (*suSocketClient) Put(*Thread, Value, Value) {
	panic("SocketClient does not support put")
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

func (sf *suSocketClient) String() string {
	return "SocketClient"
}

func (*suSocketClient) Type() types.Type {
	return types.BuiltinClass
}

func (sf *suSocketClient) Equal(other interface{}) bool {
	if sf2, ok := other.(*suSocketClient); ok {
		return sf == sf2
	}
	return false
}

func (*suSocketClient) Lookup(_ *Thread, method string) Callable {
	return suSocketClientMethods[method]
}

var newline = []byte{'\r', '\n'}

var noDeadline time.Time

var suSocketClientMethods = Methods{
	"Close": method0(func(this Value) Value {
		this.(*suSocketClient).conn.Close()
		return nil
	}),
	"Read": method1("(n)", func(this, arg Value) Value {
		ssc := this.(*suSocketClient)
		n := ToInt(arg)
		buf := make([]byte, n)
		ssc.conn.SetDeadline(time.Now().Add(ssc.timeout))
		defer ssc.conn.SetDeadline(noDeadline)
		n, e := io.ReadFull(ssc.rdr, buf)
		if e != nil && e != io.ErrUnexpectedEOF {
			panic("socketClient.Read: " + e.Error())
		}
		return SuStr(string(buf[:n]))
	}),
	"Readline": method0(func(this Value) Value {
		ssc := this.(*suSocketClient)
		ssc.conn.SetDeadline(time.Now().Add(ssc.timeout))
		defer ssc.conn.SetDeadline(noDeadline)
		return Readline(this.(*suSocketClient).rdr, "file.Readline: ")
	}),
	"Write": method1("(string)", func(this, arg Value) Value {
		ssc := this.(*suSocketClient)
		ssc.conn.SetDeadline(time.Now().Add(ssc.timeout))
		defer ssc.conn.SetDeadline(noDeadline)
		s := AsStr(arg)
		_, e := io.WriteString(ssc.conn, s)
		if e != nil {
			panic("socketClient.Write: " + e.Error())
		}
		return nil
	}),
	"Writeline": method1("(string)", func(this, arg Value) Value {
		ssc := this.(*suSocketClient)
		ssc.conn.SetDeadline(time.Now().Add(ssc.timeout))
		defer ssc.conn.SetDeadline(noDeadline)
		s := AsStr(arg)
		_, e := io.WriteString(ssc.conn, s)
		if e != nil {
			panic("socketClient.Writeline: " + e.Error())
		}
		_, e = ssc.conn.Write(newline)
		if e != nil {
			panic("socketClient.Writeline: " + e.Error())
		}
		return nil
	}),
}
