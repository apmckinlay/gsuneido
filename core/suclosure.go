// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
)

// SuClosure is an instance of a closure block
type SuClosure struct {
	this Value
	// parent is the Frame of the outer function that created this closure.
	// It is used by interp to handle block returns.
	parent *Frame
	shared *Shared // captured shared variables from parent frame
	*SuFunc
}

// Value interface

var _ Value = (*SuClosure)(nil)

func (b *SuClosure) String() string {
	return strings.Replace(b.SuFunc.String(), "block */", "closure */", 1)
}

func (b *SuClosure) Equal(other any) bool {
	return b == other
}

// Call sets up a Frame and runs a closure.
// Thread.invoke (interp.go) does similar setup, it should be kept in sync.
func (b *SuClosure) Call(th *Thread, this Value, as *ArgSpec) Value {
	fn := b.SuFunc

	// normally done by SuFunc Call
	th.Args(&b.ParamSpec, as)
	for expand := fn.Nstack - fn.Nparams; expand > 0; expand-- {
		th.Push(nil)
	}

	if this == nil {
		this = b.this
	}
	if th.fp >= len(th.frames) {
		panic("function call overflow")
	}
	fr := &th.frames[th.fp]
	fr.fn = fn
	fr.this = this
	fr.blockParent = b.parent
	fr.locals = th.stack[th.sp-int(fn.Nstack) : th.sp]
	fr.shared = b.shared
	fr.moveLocalsToShared()
	return th.run()
}

func (*SuClosure) Type() types.Type {
	return types.Block
}

func (b *SuClosure) SetConcurrent() {
	if b.this != nil {
		b.this.SetConcurrent()
	}
	if b.shared == nil || b.shared.concurrent {
		return
	}
	b.shared.concurrent = true
	// make shared values concurrent
	for _, x := range b.shared.values {
		if x != nil {
			x.SetConcurrent()
		}
	}
}

func (b *SuClosure) IsConcurrent() Value {
	return SuBool(b.shared != nil && b.shared.concurrent)
}
