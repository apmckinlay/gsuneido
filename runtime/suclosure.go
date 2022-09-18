// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// SuClosure is an instance of a closure block
type SuClosure struct {
	SuFunc
	locals     []Value
	this       Value
	concurrent bool
	// parent is the Frame of the outer function that created this closure.
	// It is used by interp to handle block returns.
	parent *Frame
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
		b.locals[int(bf.Offset)+i] = args[i]
	}

	if this == nil {
		this = b.this
	}
	v := b.locals
	if b.concurrent {
		// make a mutable copy of the locals for the frame
		v = make([]Value, len(b.locals))
		copy(v, b.locals)
	}
	return t.run(Frame{fn: bf, this: this, blockParent: b.parent,
		locals: locals{v: v, onHeap: true}})
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
	v := make([]Value, len(b.locals))
	copy(v, b.locals)
	// make them concurrent
	for _, x := range v {
		if x != nil {
			x.SetConcurrent()
		}
	}
	b.locals = v
	if b.this != nil {
		b.this.SetConcurrent()
	}
}

func (b *SuClosure) IsConcurrent() Value {
	return SuBool(b.concurrent)
}
