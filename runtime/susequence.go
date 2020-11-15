// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuSequence wraps an Iter and instantiates it lazily
// the Iter is either built-in e.g. Seq or object.Members,
// or user defined via Sequence
type SuSequence struct {
	CantConvert
	// iter is the iterator we're wrapping
	iter Iter
	// ob is nil until the sequence is instantiated
	ob *SuObject
	// duped tracks whether the sequence has been duplicated
	// this is used to decide to instantiate
	duped bool
}

func NewSuSequence(it Iter) *SuSequence {
	return &SuSequence{iter: it}
}

func (seq *SuSequence) Iter() Iter {
	if seq.Instantiated() {
		return seq.ob.Iter()
	}
	seq.duped = true
	return seq.iter.Dup()
}

func (seq *SuSequence) Instantiated() bool {
	return seq.ob != nil
}

func (seq *SuSequence) Infinite() bool {
	return seq.iter.Infinite()
}

func (seq *SuSequence) Copy() *SuObject {
	return iterToObject(seq.iter.Dup())
}

func (seq *SuSequence) instantiate() {
	if seq.ob == nil {
		seq.ob = iterToObject(seq.iter)
	}
}

const max_instantiate = 16000

func iterToObject(iter Iter) *SuObject {
	if iter.Infinite() {
		panic("can't instantiate infinite sequence")
	}
	ob := &SuObject{}
	for x := iter.Next(); x != nil; x = iter.Next() {
		ob.Add(x)
		if ob.Size() >= max_instantiate {
			panic("can't instantiate sequence larger than 16000")
		}
	}
	return ob
}

// Value interface --------------------------------------------------

var _ Value = (*SuSequence)(nil)

func (seq *SuSequence) String() string {
	if seq.iter.Infinite() {
		return "infiniteSequence"
	}
	seq.instantiate()
	return seq.ob.String()
}

func (seq *SuSequence) ToContainer() (Container, bool) {
	seq.instantiate()
	return seq.ob, true
}

func (seq *SuSequence) Get(t *Thread, key Value) Value {
	seq.instantiate()
	return seq.ob.Get(t, key)
}

func (seq *SuSequence) Put(t *Thread, key Value, val Value) {
	seq.instantiate()
	seq.ob.Put(t, key, val)
}

func (seq *SuSequence) RangeTo(i int, j int) Value {
	seq.instantiate()
	return seq.ob.RangeTo(i, j)
}

func (seq *SuSequence) RangeLen(i int, n int) Value {
	seq.instantiate()
	return seq.ob.RangeLen(i, n)
}

func (seq *SuSequence) Equal(other interface{}) bool {
	seq.instantiate()
	if y, ok := other.(*SuSequence); ok {
		y.instantiate()
		other = y.ob
	}
	return seq.ob.Equal(other)
}

func (seq *SuSequence) Hash() uint32 {
	seq.instantiate()
	return seq.ob.Hash()
}

func (seq *SuSequence) Hash2() uint32 {
	seq.instantiate()
	return seq.ob.Hash2()
}

func (*SuSequence) Type() types.Type {
	return types.Object
}

func (seq *SuSequence) Compare(other Value) int {
	seq.instantiate()
	return seq.ob.Compare(other)
}

func (*SuSequence) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Object")
}

// SequenceMethods is initialized by the builtin package
var SequenceMethods Methods

var gnSequences = Global.Num("Sequences")

func (seq *SuSequence) Lookup(t *Thread, method string) Callable {
	if seq.asSeq(method) {
		if m := Lookup(t, SequenceMethods, gnSequences, method); m != nil {
			return m
		}
	}
	seq.instantiate()
	return seq.ob.Lookup(t, method)
}

func (seq *SuSequence) asSeq(method string) bool {
	return method == "Instantiated?" ||
		(!seq.Instantiated() && (!seq.duped || seq.Infinite()))
}

// Packable ---------------------------------------------------------

var _ Packable = (*SuSequence)(nil)

func (seq *SuSequence) PackSize(clock *int32) int {
	seq.instantiate()
	return seq.ob.PackSize(clock)
}

func (seq *SuSequence) Pack(clock int32, buf *pack.Encoder) {
	seq.instantiate()
	seq.ob.Pack(clock, buf)
}

func (seq *SuSequence) PackSize2(clock int32, stack packStack) int {
	seq.instantiate()
	return seq.ob.PackSize2(clock, stack)
}

func (seq *SuSequence) PackSize3() int {
	seq.instantiate()
	return seq.ob.PackSize3()
}
