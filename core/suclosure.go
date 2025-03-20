// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

// SuClosure is an instance of a closure block
type SuClosure struct {
	this Value
	// parent is the Frame of the outer function that created this closure.
	// It is used by interp to handle block returns.
	parent *Frame
	locals []Value // if concurrent, then read-only
	*SuFunc
	concurrent bool
}

// Value interface

var _ Value = (*SuClosure)(nil)

func (b *SuClosure) String() string {
	return strings.Replace(b.SuFunc.String(), "block */", "closure */", 1)
}

func (b *SuClosure) Equal(other any) bool {
	return b == other
}

func (b *SuClosure) Call(th *Thread, this Value, as *ArgSpec) Value {
	bf := b.SuFunc

	v := b.locals
	if b.concurrent {
		// make a mutable copy of the locals for the frame
		v = slc.Clone(b.locals)
	}

	// normally done by SuFunc Call
	args := th.Args(&b.ParamSpec, as)

	// copy args
	for i := range int(b.Nparams) {
		v[int(bf.Offset)+i] = args[i]
	}

	if this == nil {
		this = b.this
	}
	return th.run(Frame{fn: bf, this: this, blockParent: b.parent,
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
	v := slc.Clone(b.locals)
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
