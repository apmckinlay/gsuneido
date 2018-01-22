package interp

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/tr"
)

// See interp.go for the rest of the Thread methods

type Thread struct {
	// frames are the Frame's making up the call stack.
	// The end of the slice is top of the stack (the current frame).
	frames []Frame
	// stack is the Value stack for arguments and expressions.
	// The end of the slice is the top of the stack.
	stack   []Value
	rxcache *regex.LruMapCache
	trcache *tr.LruMapCache
}

func NewThread() *Thread {
	return &Thread{
		rxcache: regex.NewLruMapCache(100, regex.Compile),
		trcache: tr.NewLruMapCache(100, tr.Set)}
}

func (t *Thread) Push(x Value) {
	t.stack = append(t.stack, x)
}

func (t *Thread) Pop() Value {
	last := len(t.stack) - 1
	x := t.stack[last]
	t.stack = t.stack[:last]
	return x
}

func (t *Thread) Top() Value {
	return t.stack[len(t.stack)-1]
}

func (t *Thread) SetTop(x Value) {
	t.stack[len(t.stack)-1] = x
}

func (t *Thread) Dup2() {
	t.stack = append(t.stack, t.stack[len(t.stack)-2], t.stack[len(t.stack)-1])
}

func (t *Thread) Dupx2() {
	n := len(t.stack)
	t.stack = append(t.stack, nil)
	copy(t.stack[n-2:], t.stack[n-3:])
	t.stack[n-3] = t.Top()
}
