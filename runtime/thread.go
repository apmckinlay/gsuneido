package runtime

import (
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/tr"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// See interp.go and args.go for the rest of the Thread methods

const maxStack = 1024
const maxFrames = 256

type Thread struct {
	// frames are the Frame's making up the call stack.
	frames [maxFrames]Frame
	// fp is the frame pointer, top is frames[fp-1]
	fp int

	// stack is the Value stack for arguments and expressions.
	// The end of the slice is the top of the stack.
	stack [maxStack]Value
	// sp is the stack pointer, top is stack[sp-1]
	sp int

	// this is used to pass "this" from interp to method
	// it is only temporary, Frame.this is the real "this"
	this Value

	rxcache *regex.LruMapCache
	trcache *tr.LruMapCache
}

func NewThread() *Thread {
	return &Thread{
		rxcache: regex.NewLruMapCache(100, regex.Compile),
		trcache: tr.NewLruMapCache(100, tr.Set)}
}

func (t *Thread) Push(x Value) {
	t.stack[t.sp] = x
	t.sp++
}

func (t *Thread) Pop() Value {
	t.sp--
	return t.stack[t.sp]
}

func (t *Thread) Top() Value {
	return t.stack[t.sp-1]
}

func (t *Thread) Dup2() {
	t.stack[t.sp] = t.stack[t.sp-2]
	t.stack[t.sp+1] = t.stack[t.sp-1]
	t.sp += 2
}

// Dupx2 inserts a copy of the top value under the top three
// e.g. 0,1,2,3 => 0,3,1,2,3
func (t *Thread) Dupx2() {
	t.stack[t.sp] = t.stack[t.sp-1]
	t.stack[t.sp-1] = t.stack[t.sp-2]
	t.stack[t.sp-2] = t.stack[t.sp-3]
	t.stack[t.sp-3] = t.stack[t.sp]
	t.sp++
}

func (t *Thread) Reset() {
	t.fp = 0
	t.sp = 0
}

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
