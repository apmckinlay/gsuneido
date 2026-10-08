// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Sequence, "(iter) :sequence")

func Sequence(th *Thread, args []Value) Value {
	return NewSuSequence(&wrapIter{it: args[0], th: th, ctx: th.Context()})
}

// wrapIter adapts a Suneido iterator (a class with Next,Dup,Infinite)
// to the runtime.Iter interface. For the reverse see SuIter.
// No locking since not mutable.
type wrapIter struct {
	it Value
	// When not concurrent we use the creating thread,
	// when concurrent we use a temporary thread with this ctx
	th         *Thread
	ctx        ThreadContext
	concurrent bool
}

func (wi *wrapIter) Next() Value {
	x := wi.call("Next")
	if x == wi.it {
		return nil
	}
	return x
}

func (wi *wrapIter) Infinite() (result bool) {
	return wi.call("Infinite?") == True
}

func (wi *wrapIter) Dup() Iter {
	it := wi.call("Dup")
	return &wrapIter{it: it, th: wi.th, ctx: wi.ctx,
		concurrent: wi.concurrent}
}

func (wi *wrapIter) SetConcurrent() {
	if !wi.concurrent {
		wi.concurrent = true
		wi.th = nil
		wi.it.SetConcurrent()
		wi.ctx.SetConcurrent()
	}
}

func (wi *wrapIter) IsConcurrent() Value {
	return SuBool(wi.concurrent)
}

func (wi *wrapIter) call(method string) Value {
	th := wi.th
	if wi.concurrent { // concurrent
		th = wi.ctx.NewThread()
		th.Name = "*internal*"
		defer th.Close()
	}
	return th.CallLookup(wi.it, method)
}

func (wi *wrapIter) Instantiate() *SuObject {
	return InstantiateIter(wi)
}

var _ Iter = (*wrapIter)(nil)

// for SuSequence

var _ = exportMethods(&SequenceMethods, "seq")

var _ = method(seq_Copy, "() :object")

func seq_Copy(this Value) Value {
	return this.(*SuSequence).Copy()
}

var _ = method(seq_InfiniteQ, "() :boolean")

func seq_InfiniteQ(this Value) Value {
	return SuBool(this.(*SuSequence).Infinite())
}

var _ = method(seq_InstantiatedQ, "() :boolean")

func seq_InstantiatedQ(this Value) Value {
	return SuBool(this.(*SuSequence).Instantiated())
}

var _ = method(seq_Iter, "() :object")

func seq_Iter(this Value) Value {
	iter := this.(*SuSequence).Iter()
	if wi, ok := iter.(*wrapIter); ok {
		return wi.it
	}
	return SuIter{Iter: iter}
}

var _ = method(seq_Join, "(separator :string ='') :string")

func seq_Join(this, arg Value) Value {
	iter := this.(*SuSequence).Iter()
	separator := ToStr(arg)
	sep := ""
	var sb strings.Builder
	for {
		val := iter.Next()
		if val == nil {
			break
		}
		sb.WriteString(sep)
		sep = separator
		if s, ok := val.ToStr(); ok {
			sb.WriteString(s)
		} else {
			sb.WriteString(val.String())
		}
	}
	return SuStr(sb.String())
}
