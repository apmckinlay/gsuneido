// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

// Container is an interface to objects and records
type Container interface {
	Value
	Add(val Value)
	ListSize() int
	ListGet(i int) Value
	NamedSize() int
	Copy() Container
	Slice(n int) Container
	DeleteAll()
	Insert(at int, val Value)
	Erase(t *Thread, key Value) bool
	Delete(t *Thread, key Value) bool
	GetIfPresent(t *Thread, key Value) Value
	IsReadOnly() bool
	SetReadOnly()
	ArgsIter() func() (Value, Value)
	Iter2(list bool, named bool) func() (Value, Value)
	HasKey(key Value) bool
	ToObject() *SuObject
	ToRecord(t *Thread, hdr *Header) Record
	IsConcurrent() Value
}

// iterators

func IterValues(ob Container, list, named bool) Iter {
	return &obIter{ob: ob, list: list, named: named,
		iter:   ob.Iter2(list, named),
		result: func(k, v Value) Value { return v }}
}

func IterMembers(ob Container, list, named bool) Iter {
	return &obIter{ob: ob, list: list, named: named,
		iter:   ob.Iter2(list, named),
		result: func(k, v Value) Value { return k }}
}

func IterAssocs(ob Container, list, named bool) Iter {
	return &obIter{ob: ob, list: list, named: named,
		iter:   ob.Iter2(list, named),
		result: func(k, v Value) Value { return SuObjectOf(k, v) }}
}

type obIter struct {
	ob     Container
	list   bool
	named  bool
	iter   func() (Value, Value)
	result func(Value, Value) Value
}

func (it *obIter) Next() Value {
	k, v := it.iter()
	if v == nil {
		return nil
	}
	return it.result(k, v)
}

func (it *obIter) Dup() Iter {
	oi := *it
	oi.iter = it.ob.Iter2(it.list, it.named)
	return &oi
}

func (it *obIter) Infinite() bool {
	return false
}

func (it *obIter) SetConcurrent() {
	it.ob.SetConcurrent()
}

func (it *obIter) IsConcurrent() Value {
	return it.ob.IsConcurrent()
}

func (it *obIter) Instantiate() *SuObject {
	n := 0
	if it.list && !it.named {
		n = it.ob.ListSize()
	} else if !it.list && it.named {
		n = it.ob.NamedSize()
	} else {
		n = it.ob.ListSize() + it.ob.NamedSize()
	}
	InstantiateMax(n)
	list := make([]Value, n)
	i := 0
	for k, v := it.iter(); v != nil; k, v = it.iter() {
		list[i] = it.result(k, v)
		InstantiateMax(len(list))
		i++
	}
	return NewSuObject(list)
}

var _ Iter = (*obIter)(nil)
