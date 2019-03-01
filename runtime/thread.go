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
	trcache *tr.LruMapCache
}

// NewThread creates a new thread
// zero value does not handle rxcache and trcache
func NewThread() *Thread {
	return &Thread{
		rxcache: regex.NewLruMapCache(100, regex.Compile),
		trcache: tr.NewLruMapCache(100, tr.Set)}
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

// callMethod is used by ITER and FORIN
// arguments should be on the stack
func (t *Thread) callMethod(method string, argSpec *ArgSpec) Value {
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
