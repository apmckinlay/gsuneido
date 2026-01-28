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
	thread     *Thread
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
	if th != b.thread {
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
	if th.fp >= len(th.frames) {
		panic("function call overflow")
	}
	fr := &th.frames[th.fp]
	fr.fn = bf
	fr.this = this
	fr.blockParent = b.parent
	fr.locals = locals{v: v, onHeap: true}
	if th != b.thread {
		defer func() {
			for i := range v {
				if (i < int(b.Offset) || i > int(b.Offset+b.Nparams)) &&
					v[i] != b.locals[i] {
					panic("closure changes from other thread")
				}
			}
		}()
	}
	return th.run()
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
