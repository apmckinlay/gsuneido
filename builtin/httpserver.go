// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = builtin(HttpServer, `(port, app, stop = false) :void`)

func HttpServer(th *Thread, args []Value) Value {
	guardSandbox("HttpServer")
	port := ToInt(args[0])
	addr := fmt.Sprint(":", port)
	server := &http.Server{
		Addr:    addr,
		Handler: newHttpHandler(th, args[1])}
	if ob, ok := args[2].ToContainer(); ok {
		ob.Put(th, SuStr("stop"), &suStopper{server: server})
	}
	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprint("HttpServer:", err))
	}
	return nil
}

// HttpHandler serves requests, each in its own thread, since
// net/http handles each connection in a separate goroutine
// and a Thread must not be shared between goroutines.
type HttpHandler struct {
	ctx ThreadContext
	app Value
}

// newHttpHandler captures the calling thread's context and marks the
// shared app and Suneido state as concurrent before it is shared with
// the per-request threads started by net/http.
func newHttpHandler(th *Thread, app Value) *HttpHandler {
	ctx := th.Context()
	ctx.SetConcurrent()
	app.SetConcurrent()
	return &HttpHandler{ctx: ctx, app: app}
}

var argSpecEnv = &ArgSpec{Nargs: 1,
	Spec: []byte{0}, Names: []Value{SuStr("env")}}

func (h *HttpHandler) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	th := h.ctx.NewThread()
	th.Name += " HttpServer " + rq.Method + " " + rq.URL.Path
	threads.add(th)
	defer func() {
		th.Close()
		threads.remove(th.Num)
		if r := recover(); r != nil {
			http.Error(rw, fmt.Sprint("ERROR: ", r),
				http.StatusInternalServerError) // 500
		}
	}()

	rw.Header().Set("Server", "Suneido")
	env := &suHttpEnv{rq: rq, rw: rw}
	result := th.PushCall(h.app, nil, argSpecEnv, env) // CALL THE APP
	if env.done {
		return
	}
	// result is one of:
	// - content (string)
	// - [status, content]
	// - [status, headers, content]
	if s, ok := result.ToStr(); ok {
		io.WriteString(rw, s)
	} else if ob, ok := result.ToContainer(); ok {
		status := ToInt(ob.ListGet(0))
		n := ob.ListSize()
		if n > 2 {
			addHeader(rw.Header(), ob.ListGet(1))
		}
		body := ToStr(ob.ListGet(n - 1))
		rw.WriteHeader(status)
		io.WriteString(rw, body)
	}
}

type suHttpEnv struct {
	ValueBase[suHttpEnv]
	rq   *http.Request
	rw   http.ResponseWriter
	qv   Value  // cached
	body string // cached
	done bool
}

var _ Value = (*suHttpEnv)(nil)

func (e *suHttpEnv) Get(_ *Thread, k Value) Value {
	key := ToStr(k)
	switch key {
	case "method":
		return SuStr(e.rq.Method)
	case "path":
		return SuStr(e.rq.URL.Path)
	case "query":
		return SuStr(e.rq.URL.RawQuery)
	case "queryvalues":
		if e.qv == nil {
			qv := &SuObject{}
			for k, v := range e.rq.URL.Query() {
				if len(v) == 1 {
					if v[0] == "" {
						qv.Add(SuStr(k))
					} else {
						qv.Set(SuStr(k), SuStr(v[0]))
					}
				} else {
					qv.Set(SuStr(k), SuObjectOfStrs(v))
				}
			}
			e.qv = qv
		}
		return e.qv
	case "body", "content":
		if e.body == "" {
			var sb strings.Builder
			_, err := io.Copy(&sb, e.rq.Body)
			if err != nil {
				panic(fmt.Sprint("HttpServer: ", err))
			}
			e.body = sb.String()
		}
		return SuStr(e.body)
	case "socket": // temporary backwards compatible with RackServer
		return e // not really the socket, but handles Read, Write, and CopyTo
	case "remote_user":
		addr := e.rq.RemoteAddr
		return SuStr(str.BeforeLast(addr, ":"))
	default:
		key = strings.ReplaceAll(key, "_", "-")
		v := e.rq.Header[key]
		if len(v) == 1 {
			return SuStr(v[0])
		} else if len(v) > 1 {
			return SuObjectOfStrs(v)
		} else {
			return EmptyStr
		}
	}
}

func (e *suHttpEnv) Equal(other any) bool {
	return e == other
}

func (e *suHttpEnv) Lookup(_ *Thread, method string) Value {
	return httpEnvMethods[method]
}

// writer is for CopyTo
func (e *suHttpEnv) writer() io.Writer {
	return e.rw
}

var httpEnvMethods = methods("httpEnv")

// Read allows reading the body
var _ = method(httpEnv_Read, "(n) :false|string")

func httpEnv_Read(this Value, a Value) Value {
	rd := this.(*suHttpEnv).rq.Body
	n := ToInt(a)
	if n > readMax {
		panic("HttpServer.Read: too large")
	}
	buf := make([]byte, n)
	nr, err := rd.Read(buf)
	if nr > 0 {
		return SuStr(hacks.BStoS(buf[:nr]))
	}
	if err != nil {
		if err == io.EOF {
			return False
		}
		panic(fmt.Sprint("HttpServer.Read: ", err))
	}
	return EmptyStr
}

var _ = method(httpEnv_CopyTo, "(dest, nbytes = false) :number")

func httpEnv_CopyTo(th *Thread, this Value, args []Value) Value {
	rd := this.(*suHttpEnv).rq.Body
	return CopyTo(th, rd, args[0], args[1])
}

// WriteHeader sets header values and writes the header
var _ = method(httpEnv_WriteHeader, "(status, header=false) :void")

func httpEnv_WriteHeader(this Value, a, b Value) Value {
	e := this.(*suHttpEnv)
	addHeader(e.rw.Header(), b)
	e.rw.WriteHeader(ToInt(a))
	e.done = true
	return nil
}

var _ = method(httpEnv_Write, "(s :string) :void")

func httpEnv_Write(this Value, a Value) Value {
	e := this.(*suHttpEnv)
	rw := e.rw
	_, err := io.WriteString(rw, ToStr(a))
	if err != nil {
		panic(fmt.Sprint("HttpServer.Write: ", err))
	}
	e.done = true
	return nil
}

var _ = method(httpEnv_Hijack, "() :unknown")

func httpEnv_Hijack(this Value) Value {
	e := this.(*suHttpEnv)
	hj, ok := e.rw.(http.Hijacker)
	if !ok {
		panic("HttpServer cannot Hijack")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		panic("HttpServer Hijack failed: " + err.Error())
	}
	conn.SetDeadline(noDeadline)
	return &suSocketClient{conn: conn.(*net.TCPConn), rdr: bufrw.Reader, timeout: 0}
}

//-------------------------------------------------------------------

// @immutable
type suStopper struct {
	ValueBase[suStopper]
	server *http.Server
}

var _ Value = (*suStopper)(nil)

func (st *suStopper) Equal(other any) bool {
	return st == other
}

func (st *suStopper) SetConcurrent() {
	// safe
}

var stopperParams = params("(timeout = 5)")

func (st *suStopper) Call(th *Thread, this Value, as *ArgSpec) Value {
	args := th.Args(&stopperParams, as)
	timeout := time.Duration(ToInt(args[0])) * time.Second
	server := st.server
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic("HttpServer stop: " + err.Error())
	}
	return nil
}
