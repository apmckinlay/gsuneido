// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import "strconv"

// Iter is the internal type for iterators.
// See also: SuIter and wrapIter.
type Iter interface {
	// Next returns nil when there are no more values
	Next() Value
	// Dup returns a copy of this Iter that starts at the beginning
	Dup() Iter
	Infinite() bool
	SetConcurrent()
	IsConcurrent() Value
	Instantiate() *SuObject
}

const MaxInstantiate = MaxSuInt

func InstantiateMax(n int) {
	if n >= MaxInstantiate {
		panic("can't instantiate sequence larger than " +
			strconv.Itoa(MaxInstantiate))
	}
}

func InstantiateIter(iter Iter) *SuObject {
	if iter.Infinite() {
		panic("can't instantiate infinite sequence")
	}
	list := make([]Value, 0, 8)
	for x := iter.Next(); x != nil; x = iter.Next() {
		list = append(list, x)
		InstantiateMax(len(list))
	}
	return NewSuObject(list)
}
