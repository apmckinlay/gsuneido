// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// SuClosure is an instance of a closure block
type SuClosure struct {
	SuFunc
	locals Locals
	this   Value
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
	b.locals.Lock()
	for i := 0; i < int(b.Nparams); i++ {
		b.locals.v[int(bf.Offset)+i] = args[i]
	}
	b.locals.Unlock()

	if this == nil {
		this = b.this
	}
	t.frames[t.fp] = Frame{fn: bf, locals: b.locals, this: this}
	if bf.cover != nil {
		coverage(bf, 0)
	}
	return t.run()
}

func (*SuClosure) Type() types.Type {
	return types.Block
}

func (b *SuClosure) SetConcurrent() {
	b.locals.SetConcurrent()
	if b.this != nil {
		b.this.SetConcurrent()
	}
}
