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
	ss.Data["CallClass"] =
		methodRaw("(name = nil, port = nil, exit = false)", ssCallClass)
	Global.Builtin("SocketServer", ss)
}

func ssCallClass(t *Thread, as *ArgSpec, this Value, args []Value) Value {
	name, port, as2 := ssArgs(t, as, this, args)
	class := this.(*SuClass)
	sm := suServerMaster{SuInstance: class.New(t, as2)}
	if OnUiThread() {
		// don't block UI thread
		go sm.listen(ToStr(name), ToInt(port))
	} else {
		sm.listen(ToStr(name), ToInt(port))
	}
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

var nSocketServerClient = int32(0)

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
		if atomic.LoadInt32(&nSocketServerClient) > ssmax {
			log.Printf("SocketServer: too many connections, stopping (%d %s)",
				port, name)
			return
		}
		go sm.connect(name, conn)
	}
}

func (sm *suServerMaster) connect(name string, conn net.Conn) {
	atomic.AddInt32(&nSocketServerClient, 1)
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
			if e := recover(); e != nil {
				LogInternalError("in SocketServer thread:", e)
				t.PrintStack()
			}
			threads.remove(t.Num)
			t.Close()
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
	switch method {
	case "RemoteUser":
		return remoteUser
	case "ManualClose":
		return manualClose
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
	sc.client.conn.Close()
	sc.client.conn = nil
	atomic.AddInt32(&nSocketServerClient, -1)
}

var remoteUser = method0(func(this Value) Value {
	sc := this.(*suServerConnect)
	addr := sc.client.conn.RemoteAddr().String()
	return SuStr(str.BeforeLast(addr, ":"))
})

var manualClose = method0(func(this Value) Value {
	this.(*suServerConnect).manualClose = true
	return nil
})
