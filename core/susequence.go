// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuSequence wraps an Iter and instantiates it lazily
// the Iter is either built-in e.g. Seq or object.Members,
// or user defined via Sequence
type SuSequence struct {
	ValueBase[SuSequence]
	// iter is the iterator we're wrapping
	iter Iter
	// ob is nil until the sequence is instantiated
	ob *SuObject
	MayLock
	// duped tracks whether the sequence has been duplicated.
	// This is set by Iter() and used by asSeq() to decide to instantiate
	duped bool
}

func NewSuSequence(it Iter) *SuSequence {
	return &SuSequence{iter: it}
}

func (seq *SuSequence) Iter() Iter {
	iter, ob := seq.iter2()
	if ob != nil {
		return ob.Iter()
	}
	return iter.Dup() // may lock
}

func (seq *SuSequence) iter2() (Iter, *SuObject) {
	if seq.Lock() {
		defer seq.Unlock()
	}
	if seq.ob != nil {
		return nil, seq.ob
	}
	seq.duped = true
	return seq.iter, nil
}

func (seq *SuSequence) Instantiated() bool {
	if seq.Lock() {
		defer seq.Unlock()
	}
	return seq.ob != nil
}

func (seq *SuSequence) snapshot() (Iter, *SuObject) {
	if seq.Lock() {
		defer seq.Unlock()
	}
	return seq.iter, seq.ob
}

func (seq *SuSequence) Infinite() bool {
	iter, ob := seq.snapshot()
	return ob == nil && iter.Infinite() // may lock
}

func (seq *SuSequence) Copy() Value {
	iter, ob := seq.snapshot()
	if ob != nil {
		return ob.Copy()
	}
	return iter.Dup().Instantiate() // may lock
}

func (seq *SuSequence) instantiate() *SuObject {
	iter, ob := seq.snapshot()
	if ob == nil {
		ob = iter.Instantiate() // may lock
		if seq.concurrent {
			ob.SetConcurrent()
		}
		seq.setOb(ob) // race, but should be benign/idempotent
	}
	return ob
}

func (seq *SuSequence) setOb(ob *SuObject) {
	if seq.Lock() {
		defer seq.Unlock()
	}
	seq.ob = ob
	seq.iter = nil
}

// Value interface --------------------------------------------------

var _ Value = (*SuSequence)(nil)

func (seq *SuSequence) String() string {
	if seq.Infinite() {
		return "infiniteSequence"
	}
	return seq.instantiate().String()
}

func (seq *SuSequence) ToContainer() (Container, bool) {
	return seq.instantiate(), true
}

func (seq *SuSequence) Get(th *Thread, key Value) Value {
	return seq.instantiate().Get(th, key)
}

func (seq *SuSequence) Put(th *Thread, key Value, val Value) {
	seq.instantiate().Put(th, key, val)
}

func (seq *SuSequence) GetPut(th *Thread, key Value, val Value,
	op func(x, y Value) Value, retOrig bool) Value {
	return seq.instantiate().GetPut(th, key, val, op, retOrig)
}

func (seq *SuSequence) RangeTo(i int, j int) Value {
	return seq.instantiate().RangeTo(i, j)
}

func (seq *SuSequence) RangeLen(i int, n int) Value {
	return seq.instantiate().RangeLen(i, n)
}

func (seq *SuSequence) Equal(other any) bool {
	x := seq.instantiate()
	if y, ok := other.(*SuSequence); ok {
		other = y.instantiate()
	}
	return x.Equal(other)
}

func (seq *SuSequence) Hash() uint64 {
	return seq.instantiate().Hash()
}

func (seq *SuSequence) Hash2() uint64 {
	return seq.instantiate().Hash2()
}

func (*SuSequence) Type() types.Type {
	return types.Object
}

func (seq *SuSequence) Compare(other Value) int {
	return seq.instantiate().Compare(other)
}

// SequenceMethods is initialized by the builtin package
var SequenceMethods Methods

var gnSequences = Global.Num("Sequences")

func (seq *SuSequence) Lookup(th *Thread, method string) Callable {
	if seq.asSeq(method) {
		if m := Lookup(th, SequenceMethods, gnSequences, method); m != nil {
			return m
		}
	}
	return seq.instantiate().Lookup(th, method)
}

func (seq *SuSequence) asSeq(method string) bool {
	if method == "Instantiated?" || seq.Infinite() {
		return true
	}
	if seq.Instantiated() {
		return false
	}
	if seq.Lock() {
		defer seq.Unlock()
	}
	return !seq.duped
}

func (seq *SuSequence) SetConcurrent() {
	if seq.SetConc() {
		if seq.ob != nil {
			seq.ob.SetConcurrent()
		} else {
			seq.iter.SetConcurrent()
		}
	}
}

// Packable ---------------------------------------------------------

var _ Packable = (*SuSequence)(nil)

func (seq *SuSequence) PackSize(hash *uint64) int {
	return seq.instantiate().PackSize(hash)
}

func (seq *SuSequence) Pack(hash *uint64, buf *pack.Encoder) {
	seq.instantiate().Pack(hash, buf)
}

func (seq *SuSequence) PackSize2(hash *uint64, stack packStack) int {
	return seq.instantiate().PackSize2(hash, stack)
}
