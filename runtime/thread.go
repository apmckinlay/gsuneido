// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"log"
	"strconv"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tr"
)

// MainThread is injected by gsuneido.go to use for debugging
var MainThread *Thread

// See interp.go and args.go for the rest of the Thread methods

// maxStack is the size of the value stack, fixed size for performance
const maxStack = 1024

// maxFrames is the size of the frame stack, fixed size for performance
const maxFrames = 256

type Thread struct {
	// frames are the Frame's making up the call stack.
	frames [maxFrames]Frame
	// fp is the frame pointer, top is frames[fp-1]
	fp int
	// fpMax is the "high water" mark for fp
	fpMax int

	// stack is the Value stack for arguments and expressions.
	// The end of the slice is the top of the stack.
	stack [maxStack]Value
	// sp is the stack pointer, top is stack[sp-1]
	sp int
	// spMax is the "high water" mark for sp
	spMax int

	// blockReturnFrame is the parent frame of the block that is returning
	blockReturnFrame *Frame

	// RxCache is per thread so no locking is required
	RxCache regex.Cache
	// TrCache is per thread so no locking is required
	TrCache tr.Cache

	// rules is a stack of the currently running rules, used by SuRecord
	rules activeRules

	// dbms is the database (client or local) for this Thread
	dbms IDbms

	// Num is a unique number assigned to the thread
	Num int32

	// Name is the name of the thread (default is Thread-#)
	Name string

	// UIThread is only set for the main UI thread.
	// It controls whether interp checks for UI requests from other threads.
	UIThread bool

	// OpCount counts op codes in interp, for polling
	OpCount int

	// Quote is used by Display to request specific quotes
	Quote int

	// InHandler is used to detect nested handler calls
	InHandler bool

	// Session is the name of the database session for clients and standalone.
	// Server tracks client session names separately.
	// Needs atomic because we access MainThread from other threads.
	session atomic.Value

	// Suneido is a per-thread SuneidoObject that overrides the global one
	Suneido *SuneidoObject

	profile profile
}

var nThread int32

// NewThread creates a new thread.
// It is primarily used for user initiated threads.
// Internal threads can just use a zero Thread.
func NewThread(parent *Thread) *Thread {
	n := atomic.AddInt32(&nThread, 1)
	name := "Thread-" + strconv.Itoa(int(n))
	t := &Thread{Num: n, Name: name}
	if parent != nil && parent.Suneido != nil {
		parent.Suneido.SetConcurrent()
		t.Suneido = parent.Suneido
	}
	mts := ""
	if MainThread != nil {
		mts = MainThread.Session()
	}
	t.session.Store(str.Opt(mts, ":") + name)
	return t
}

func (t *Thread) Session() string {
	if v := t.session.Load(); v != nil {
		return v.(string)
	}
	return ""
}

func (t *Thread) SetSession(s string) {
	t.session.Store(s)
}

// Push pushes a value onto the value stack
func (t *Thread) Push(x Value) {
	if t.sp >= maxStack {
		Fatal("value stack overflow")
	}
	t.stack[t.sp] = x
	t.sp++
}

// Pop pops a value off the value stack
func (t *Thread) Pop() Value {
	t.sp--
	return t.stack[t.sp]
}

// Top returns the top of the value stack (without modifying the stack)
func (t *Thread) Top() Value {
	return t.stack[t.sp-1]
}

// Swap exchanges the top two values on the stack
func (t *Thread) Swap() {
	t.stack[t.sp-1], t.stack[t.sp-2] = t.stack[t.sp-2], t.stack[t.sp-1]
}

// Reset sets a Thread back to its initial state
func (t *Thread) Reset() {
	t.fp = 0
	t.sp = 0
	t.Name = ""
	t.blockReturnFrame = nil
	t.InHandler = false
	t.Suneido = nil
}

// GetState and RestoreState are used by callbacks_windows.go

type state struct {
	fp int
	sp int
}

func (t *Thread) GetState() state {
	return state{fp: t.fp, sp: t.sp}
}

func (t *Thread) RestoreState(st any) {
	s := st.(state)
	t.fp = s.fp
	t.sp = s.sp
}

// Callstack captures the call stack
func (t *Thread) Callstack() *SuObject {
	// NOTE: it might be more efficient
	// to capture the call stack in an internal format
	// and only build the SuObject if required
	cs := &SuObject{}
	for i := t.fp - 1; i >= 0; i-- {
		fr := t.frames[i]
		call := &SuObject{}
		call.Set(SuStr("fn"), fr.fn)
		call.Set(SuStr("srcpos"), IntVal(fr.fn.CodeToSrcPos(fr.ip-1)))
		call.Set(SuStr("locals"), t.locals(i))
		cs.Add(call)
	}
	return cs
}

func (t *Thread) Locals(i int) *SuObject {
	return t.locals(t.fp - 1 - i)
}

func (t *Thread) locals(i int) *SuObject {
	fr := t.frames[i]
	locals := &SuObject{}
	if fr.this != nil {
		locals.Set(SuStr("this"), fr.this)
	}
	for i, v := range fr.locals.v {
		if v != nil && fr.fn != nil && i < len(fr.fn.Names) {
			locals.Set(SuStr(fr.fn.Names[i]), v)
		}
	}
	return locals
}

// PrintStack prints the thread's call stack
func (t *Thread) PrintStack() {
	PrintStack(t.Callstack())
}

func PrintStack(cs *SuObject) {
	if cs == nil {
		return
	}
	// toStr := func(x Value) (s string) {
	// 	defer func() {
	// 		if e := recover(); e != nil {
	// 			s = fmt.Sprint(e)
	// 		}
	// 	}()
	// 	s = x.String()
	// 	if len(s) > 230 {
	// 		s = s[:230] + "..."
	// 	}
	// 	return s
	// }
	for i := 0; i < cs.ListSize(); i++ {
		frame := cs.ListGet(i)
		fn := frame.Get(nil, SuStr("fn"))
		log.Println(fn)
		// locals := frame.Get(nil, SuStr("locals"))
		// log.Println("   " + toStr(locals))
	}
}

// SetDbms is used to set up the main thread initially
func (t *Thread) SetDbms(dbms IDbms) {
	t.dbms = dbms
}

// GetDbms requires dependency injection
var GetDbms func() IDbms

func (t *Thread) Dbms() IDbms {
	if t.dbms == nil {
		t.dbms = GetDbms()
		if s := t.session.Load(); s != nil {
			// session id was set before connecting
			t.dbms.SessionId(t, s.(string))
		}
	}
	return t.dbms
}

// Close closes the thread's dbms connection (if it has one)
func (t *Thread) Close() {
	if t.dbms != nil && options.Action == "client" {
		t.dbms.Close()
		t.dbms = nil
	}
}

// SubThread is a NewThread with the same dbms as this thread.
// This is used for the RunOnGoSide and SuneidoAPP threads.
// We want a new thread for isolation e.g. for exceptions or dynamic variables
// but we don't need the overhead of another dbms connection.
// WARNING: This should only be used where it is guaranteed
// that the Threads will NOT be used concurrently.
func (t *Thread) SubThread() *Thread {
	t2 := NewThread(t)
	t2.dbms = t.dbms
	return t2
}

func (t *Thread) Cat(x, y Value) Value {
	return OpCat(t, x, y)
}

func (t *Thread) SessionId(id string) string {
	if t.dbms == nil {
		// don't create a connection just to get/set the session id
		if id != "" {
			t.SetSession(id)
		}
		return t.Session()
	}
	return t.dbms.SessionId(t, id)
}

func (t *Thread) RunWithMainSuneido(fn func() Value) Value {
	defer func(orig *SuneidoObject) {
		t.Suneido = orig
	}(t.Suneido)
	t.Suneido = nil
	return fn()
}
