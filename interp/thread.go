package interp

import v "github.com/apmckinlay/gsuneido/value"

type Thread struct {
	// frames are the Frame's making up the call stack.
	// The end of the slice is top of the stack (the current frame).
	frames []Frame
	// stack is the Value stack for arguments and expressions.
	// The end of the slice is the top of the stack.
	stack []v.Value
}

func (t *Thread) Push(x v.Value) {
	t.stack = append(t.stack, x)
}

func (t *Thread) Pop() v.Value {
	last := len(t.stack) - 1
	x := t.stack[last]
	t.stack = t.stack[:last]
	return x
}

func (t *Thread) Top() v.Value {
	return t.stack[len(t.stack)-1]
}

func (t *Thread) Dup2() {
	t.stack = append(t.stack, t.stack[len(t.stack)-1], t.stack[len(t.stack)-2])
}

// Call executes a SuFunc and returns the result.
// The arguments must be already on the stack as per the ArgSpec.
// On return, the arguments are removed from the stack.
func (t *Thread) Call(fn *v.SuFunc, as ArgSpec) v.Value {
	defer func(sp int) { t.stack = t.stack[:sp] }(len(t.stack) - as.Nargs())
	t.args(fn, as)
	base := len(t.stack) - fn.Nparams
	for i := fn.Nparams; i < fn.Nlocals; i++ {
		t.Push(nil)
	}
	frame := Frame{fn: fn, ip: 0, locals: t.stack[base:]}
	t.frames = append(t.frames, frame)
	defer func(fp int) { t.frames = t.frames[:fp] }(len(t.frames) - 1)
	return t.Interp()
}

// args converts the arguments on the stack as per the ArgSpec
// into the parameters expected by the function.
// On return, the stack is guaranteed to match the SuFunc.
func (t *Thread) args(fn *v.SuFunc, as ArgSpec) {
	if fn.Nparams == as.N_unnamed() {
		return // simple fast path
	}
	panic("not implemented") // TODO
}
