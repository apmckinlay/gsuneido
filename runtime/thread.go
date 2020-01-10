// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"log"
	"strconv"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/tr"
)

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

	// RxCache is per thread so no locking is required
	RxCache *regex.PatternCache
	// TrCache is per thread so no locking is required
	TrCache *tr.TrsetCache

	// rules is a stack of the currently running rules, used by SuRecord
	rules activeRules

	// dbms is the database (client or local) for this Thread
	dbms IDbms

	// Token is an authorization token
	Token string

	// Num is a unique number assigned to the thread
	Num int32

	// Name is the name of the thread (default is Thread-#)
	Name string

	// Poll is whether or not interp should poll (only set for main thread)
	Poll bool

	// OpCount counts op codes in interp, for polling
	OpCount int
}

var nThread int32

// NewThread creates a new thread
// zero value does not handle rxcache and trcache
func NewThread() *Thread {
	n := atomic.AddInt32(&nThread, 1)
	return &Thread{
		RxCache: regex.NewPatternCache(100, regex.Compile),
		TrCache: tr.NewTrsetCache(100, tr.Set),
		Num:     n,
		Name:    "Thread-" + strconv.Itoa(int(n))}
}

// Push pushes a value onto the value stack
func (t *Thread) Push(x Value) {
	if t.sp >= maxStack {
		log.Fatalln("value stack overflow")
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

// Dup2 duplicates the top two values on the stack i.e. a,b => a,b,a,b
func (t *Thread) Dup2() {
	t.stack[t.sp] = t.stack[t.sp-2]
	t.stack[t.sp+1] = t.stack[t.sp-1]
	t.sp += 2
}

// Dupx2 inserts a copy of the top value under the next two i.e. 1,2,3 => 3,1,2,3
func (t *Thread) Dupx2() {
	t.stack[t.sp] = t.stack[t.sp-1]
	t.stack[t.sp-1] = t.stack[t.sp-2]
	t.stack[t.sp-2] = t.stack[t.sp-3]
	t.stack[t.sp-3] = t.stack[t.sp]
	t.sp++
}

// Reset sets sp and fp to 0, only used by tests
func (t *Thread) Reset() {
	t.fp = 0
	t.sp = 0
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
		call.Set(SuStr("srcpos"), IntVal(fr.fn.CodeToSrcPos(fr.ip)))
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
	for i, v := range fr.locals {
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
	for i := 0; i < cs.ListSize(); i++ {
		frame := cs.ListGet(i)
		fn := frame.Get(nil, SuStr("fn"))
		log.Println(fn)
		locals := frame.Get(nil, SuStr("locals"))
		s := locals.String()
		if len(s) > 230 {
			s = s[:230] + "..."
		}
		log.Println(s)
	}
}

// GetDbms requires dependency injection
var GetDbms func() IDbms

func (t *Thread) Dbms() IDbms {
	if t.dbms == nil {
		t.dbms = GetDbms()
		if t.Token != "" {
			t.dbms.Auth(t.Token)
		}
	}
	return t.dbms
}

func (t *Thread) Close() {
	if t.dbms != nil {
		t.dbms.Close()
	}
}

// SubThread is a NewThread with the same dbms as this thread.
// This is used for the UpdateUI and SuneidoAPP threads.
// We want a new thread for isolation e.g. for exceptions or dynamic variables
// but we don't need the overhead of another dbms connection.
// WARNING: This should only be used where it is guaranteed
// that the Threads will NOT be used concurrently.
func (t *Thread) SubThread() *Thread {
	t2 := NewThread()
	t2.dbms = t.dbms
	return t2
}
