// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// SuClosure is an instance of a closure block
type SuClosure struct {
	SuFunc
	locals     Locals
	this       Value
	concurrent bool
}

// Value interface

var _ Value = (*SuClosure)(nil)

func (b *SuClosure) String() string {
	return "/* block */"
}

func (b *SuClosure) Call(t *Thread, this Value, as *ArgSpec) Value {
	bf := &b.SuFunc

	// normally done by SuFunc Call
	args := t.Args(&b.ParamSpec, as)

	// copy args
	for i := 0; i < int(b.Nparams); i++ {
		b.locals.v[int(bf.Offset)+i] = args[i]
	}

	if this == nil {
		this = b.this
	}
	fr := Frame{fn: bf, locals: b.locals, this: this}
	if b.concurrent {
		// make a mutable copy of the locals for the frame
		v := make([]Value, len(b.locals.v))
		copy(v, b.locals.v)
		fr.locals.v = v
	}
	t.frames[t.fp] = fr
	return t.run()
}

func (*SuClosure) Type() types.Type {
	return types.Block
}

func (b *SuClosure) SetConcurrent() {
	if b.concurrent {
		return
	}
	b.concurrent = true
	// make a copy of the locals - read-only since it will be shared
	v := make([]Value, len(b.locals.v))
	copy(v, b.locals.v)
	for _, x := range v {
		if x != nil {
			x.SetConcurrent()
		}
	}
	b.locals.v = v
}

func (b *SuClosure) IsConcurrent() bool {
	return b.concurrent
}
