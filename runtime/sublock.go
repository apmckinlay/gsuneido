package runtime

import "github.com/apmckinlay/gsuneido/runtime/types"

// SuBlock is an instance of a closure block
type SuBlock struct {
	SuFunc
	locals []Value
	this   Value
}

// Value interface

var _ Value = (*SuBlock)(nil)

func (b *SuBlock) String() string {
	return "/* block */"
}

func (b *SuBlock) Call(t *Thread, as *ArgSpec) Value {
	bf := &b.SuFunc

	// normally done by SuFunc Call
	args := t.Args(&b.ParamSpec, as)

	// copy args
	for i := 0; i < int(b.Nparams); i++ {
		b.locals[int(bf.Offset)+i] = args[i]
	}

	// normally done by Thread.Call
	this := b.this
	if t.this != nil {
		this = t.this
	}
	t.frames[t.fp] = Frame{fn: bf, locals: b.locals, this: this}
	t.this = nil
	return t.run()
}

// TypeName returns the Suneido name for the type (Value interface)
func (*SuBlock) Type() types.Type {
	return types.Block
}
