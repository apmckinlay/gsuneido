// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

func init() {
	ss := &SuClass{Lib: "builtin", Name: "SocketServer", MemBase: NewMemBase()}
	ss.Data["CallClass"] = &SuBuiltinMethodRaw{Fn: ssCallClass,
		ParamSpec: params("(name = nil, port = nil, exit = false)")}
	Global.Builtin("SocketServer", ss)
}

func ssCallClass(t *Thread, as *ArgSpec, this Value, args []Value) Value {
	if OnUIThread() {
		panic("SocketServer not allowed on UI thread")
	}
	name, port, as2 := ssArgs(t, as, this, args)
	class := this.(*SuClass)
	sm := suServerMaster{SuInstance: class.New(t, as2)}
	sm.listen(ToStr(name), ToInt(port))
	return nil
}

func ssArgs(t *Thread, as *ArgSpec, this Value, args []Value) (
	name, port Value, as2 *ArgSpec) {
	name = this.Get(t, SuStr("Name"))
	port = this.Get(t, SuStr("Port"))
	ai := NewArgsIter(as, args)
	k, v := ai()
	if v != nil && k == nil {
		name = v
		k, v = ai()
		if v != nil && k == nil {
			port = v
			k, v = ai()
		}
	}
	n := len(as.Spec)
	names := make([]Value, 0, n)
	spec := make([]byte, 0, n)
	nargs := byte(0)
	for ; v != nil; k, v = ai() { // copy remaining args, extracting name & port
		if SuStr("name").Equal(k) {
			name = v
		} else if SuStr("port").Equal(k) {
			port = v
		} else {
			nargs++
			t.Push(v)
			if k != nil {
				spec = append(spec, byte(len(names)))
				names = append(names, k)
			}
		}
	}
	if name != nil {
		nargs++
		t.Push(name)
		spec = append(spec, byte(len(names)))
		names = append(names, SuStr("Name"))
	}
	if port != nil {
		nargs++
		t.Push(port)
		spec = append(spec, byte(len(names)))
		names = append(names, SuStr("Port"))
	}
	as2 = &ArgSpec{Nargs: nargs, Names: names, Spec: spec}
	if port == nil {
		panic("SocketServer: no port specified")
	}
	if name == nil {
		name = EmptyStr
	}
	return
}

type suServerMaster struct {
	*SuInstance
}

func (sm *suServerMaster) String() string {
	return "SocketServer master"
}

var nSocketServerClient atomic.Int32

const ssmax = 500 // for all SocketServer's

func (sm *suServerMaster) listen(name string, port int) {
	addr := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		if nSocketServerClient.Load() > ssmax {
			log.Printf("SocketServer: too many connections, stopping (%d %s)",
				port, name)
			return
		}
		go sm.connect(name, conn)
	}
}

func (sm *suServerMaster) connect(name string, conn net.Conn) {
	nSocketServerClient.Add(1)
	client := suSocketClient{
		conn: conn.(*net.TCPConn), rdr: bufio.NewReader(conn),
		// no timeout to match jSuneido
	}
	sc := &suServerConnect{
		SuInstance: sm.SuInstance.Copy(),
		client:     client,
	}
	defer sc.close()
	t := NewThread(nil)
	t.Name = str.BeforeFirst(t.Name, " ") + " " + name
	if f := sc.Lookup(t, "Run"); f != nil {
		threads.add(t)
		defer func() {
			t.Close()
			threads.remove(t.Num)
			if e := recover(); e != nil {
				LogUncaught(t, "SocketServer", e)
			}
		}()
		f.Call(t, sc, &ArgSpec0)
	}
}

type suServerConnect struct {
	*SuInstance
	client      suSocketClient
	manualClose bool
}

func (sc *suServerConnect) Lookup(t *Thread, method string) Callable {
	if f, ok := socketServerMethods[method]; ok {
		return f
	}
	if f := sc.client.Lookup(t, method); f != nil {
		return f
	}
	return sc.SuInstance.Lookup(t, method)
}

func (sc *suServerConnect) close() {
	if !sc.manualClose {
		sc.Close()
	}
}

func (sc *suServerConnect) Close() {
	if sc.client.conn != nil {
		sc.client.conn.Close()
		sc.client.conn = nil
	}
	nSocketServerClient.Add(-1)
}

var socketServerMethods = methods()

var _ = method(sock_RemoteUser, "()")

func sock_RemoteUser(this Value) Value {
	sc := this.(*suServerConnect)
	addr := sc.client.conn.RemoteAddr().String()
	return SuStr(str.BeforeLast(addr, ":"))
}

var _ = method(sock_ManualClose, "()")

func sock_ManualClose(this Value) Value {
	this.(*suServerConnect).manualClose = true
	return nil
}
