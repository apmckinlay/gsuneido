package runtime

import (
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

	rxcache *regex.LruMapCache
	TrCache *tr.LruMapCache
}

// NewThread creates a new thread
// zero value does not handle rxcache and trcache
func NewThread() *Thread {
	return &Thread{
		rxcache: regex.NewLruMapCache(100, regex.Compile),
		TrCache: tr.NewLruMapCache(100, tr.Set)}
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

// CallMethod calls a Suneido method
// arguments (including "this") should be on the stack
func (t *Thread) CallMethod(method string, argSpec *ArgSpec) Value {
	base := t.sp - int(argSpec.Nargs) - 1
	ob := t.stack[base]
	f := ob.Lookup(method)
	if f == nil {
		panic("method not found " + ob.TypeName() + "." + method)
	}
	t.this = ob
	result := f.Call(t, argSpec)
	t.sp = base
	return result
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
		call.Put(SuStr("fn"), fr.fn)
		locals := &SuObject{}
		if fr.this != nil {
			locals.Put(SuStr("this"), fr.this)
		}
		for i, v := range fr.locals {
			if v != nil {
				locals.Put(SuStr(fr.fn.Names[i]), v)
			}
		}
		call.Put(SuStr("locals"), locals)
		cs.Add(call)
	}
	return cs
}
