package runtime

import (
	"strconv"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/tr"
	"github.com/apmckinlay/gsuneido/util/verify"
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

	// this is used to pass "this" from interp to method
	// it is only temporary, Frame.this is the real "this"
	this Value

	// RxCache is per thread so no locking is required
	RxCache *regex.PatternCache
	// TrCache is per thread so no locking is required
	TrCache *tr.TrsetCache

	// rules is a stack of the currently running rules, used by SuRecord
	rules activeRules

	// dbms is the database (client or local) for this Thread
	dbms IDbms

	// Name is the name of the thread (default is Thread-#)
	Name string
}

var nThread int32

// NewThread creates a new thread
// zero value does not handle rxcache and trcache
func NewThread() *Thread {
	n := atomic.AddInt32(&nThread, 1)
	return &Thread{
		RxCache: regex.NewPatternCache(100, regex.Compile),
		TrCache: tr.NewTrsetCache(100, tr.Set),
		Name:    "Thread-" + strconv.Itoa(int(n))}
}

// Push pushes a value onto the value stack
func (t *Thread) Push(x Value) {
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

// CallWithArgs pushes the arguments onto the stack and calls the function
func (t *Thread) CallWithArgs(fn Value, args ...Value) Value {
	verify.That(len(args) < AsEach)
	as := StdArgSpecs[len(args)]
	base := t.sp
	for _, x := range args {
		t.Push(x)
	}
	result := fn.Call(t, as)
	t.sp = base
	return result
}

// CallWithArgSpec pushes the arguments onto the stack and calls the function
func (t *Thread) CallWithArgSpec(fn Value, as *ArgSpec, args ...Value) Value {
	base := t.sp
	for _, x := range args {
		t.Push(x)
	}
	result := fn.Call(t, as)
	t.sp = base
	return result
}

// CallMethod calls a *named* method.
// Arguments (including "this") should be on the stack
func (t *Thread) CallMethod(method string, argSpec *ArgSpec) Value {
	base := t.sp - int(argSpec.Nargs) - 1
	ob := t.stack[base]
	f := ob.Lookup(t, method)
	if f == nil {
		panic("method not found: " + ob.Type().String() + "." + method)
	}
	result := CallMethod(t, ob, f, argSpec)
	t.sp = base
	return result
}

// CallAsMethod runs a function as if it were a method of an object.
// Implements object.Eval
func (t *Thread) CallAsMethod(ob, fn Value, args ...Value) Value {
	if m, ok := fn.(*SuMethod); ok {
		fn = m.GetFn()
	}
	t.this = ob
	t.Push(ob)
	return t.CallWithArgs(fn, args...)
}

func (t *Thread) CallMethodWithArgSpec(
	this Value, method string, as *ArgSpec, args ...Value) Value {
	t.Push(this)
	for _, x := range args {
		t.Push(x)
	}
	return t.CallMethod(method, as)
}

// Callstack captures the call stack
func (t *Thread) CallStack() *SuObject {
	// NOTE: it might be more efficient
	// to capture the call stack in an internal format
	// and only build the SuObject if required
	cs := &SuObject{}
	for i := t.fp - 1; i >= 0; i-- {
		fr := t.frames[i]
		call := &SuObject{}
		call.Set(SuStr("fn"), fr.fn)
		call.Set(SuStr("srcpos"), IntVal(fr.fn.CodeToSrcPos(fr.ip)))
		locals := &SuObject{}
		if fr.this != nil {
			locals.Set(SuStr("this"), fr.this)
		}
		for i, v := range fr.locals {
			if v != nil && fr.fn != nil && i < len(fr.fn.Names) {
				locals.Set(SuStr(fr.fn.Names[i]), v)
			}
		}
		call.Set(SuStr("locals"), locals)
		cs.Add(call)
	}
	return cs
}

// GetDbms requires dependency injection
var GetDbms func() IDbms

func (t *Thread) Dbms() IDbms {
	if t.dbms == nil {
		t.dbms = GetDbms()
	}
	return t.dbms
}

func (t *Thread) Close() {
	if t.dbms != nil {
		t.dbms.Close()
	}
}

func (t *Thread) takeThis() Value {
	tmp := t.this
	t.this = nil
	return tmp
}
