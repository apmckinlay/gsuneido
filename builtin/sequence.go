// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(Sequence, "(iter)")

func Sequence(t *Thread, args []Value) Value {
	return NewSuSequence(&wrapIter{it: args[0], t: t})
}

// wrapIter adapts a Suneido iterator (a class with Next,Dup,Infinite)
// to the runtime.Iter interface. For the reverse see SuIter.
// No locking since not mutable.
type wrapIter struct {
	it Value
	// t is nil when concurrent.
	// When not concurrent we use the creating thread.
	t *Thread
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
	return &wrapIter{it: it, t: wi.t}
}

func (wi *wrapIter) SetConcurrent() {
	wi.t = nil
	wi.it.SetConcurrent()
}

func (wi *wrapIter) IsConcurrent() Value {
	return SuBool(wi.t == nil)
}

func (wi *wrapIter) call(method string) Value {
	t := wi.t
	if t == nil {
		t = &Thread{}
		t.Name = "*internal*"
		defer t.Close()
	}
	return t.CallLookup(wi.it, method)
}

func (wi *wrapIter) Instantiate() *SuObject {
	return InstantiateIter(wi)
}

var _ Iter = (*wrapIter)(nil)

// for SuSequence

var _ = exportMethods(&SequenceMethods)

var _ = method(seq_Copy, "()")

func seq_Copy(this Value) Value {
	return this.(*SuSequence).Copy()
}

var _ = method(seq_InfiniteQ, "()")

func seq_InfiniteQ(this Value) Value {
	return SuBool(this.(*SuSequence).Infinite())
}

var _ = method(seq_InstantiatedQ, "()")

func seq_InstantiatedQ(this Value) Value {
	return SuBool(this.(*SuSequence).Instantiated())
}

var _ = method(seq_Iter, "()")

func seq_Iter(this Value) Value {
	iter := this.(*SuSequence).Iter()
	if wi, ok := iter.(*wrapIter); ok {
		return wi.it
	}
	return SuIter{Iter: iter}
}

var _ = method(seq_Join, "(separator='')")

func seq_Join(this, arg Value) Value {
	iter := this.(*SuSequence).Iter()
	separator := ToStr(arg)
	sep := ""
	var buf strings.Builder
	for {
		val := iter.Next()
		if val == nil {
			break
		}
		buf.WriteString(sep)
		sep = separator
		if s, ok := val.ToStr(); ok {
			buf.WriteString(s)
		} else {
			buf.WriteString(val.String())
		}
	}
	return SuStr(buf.String())
}
