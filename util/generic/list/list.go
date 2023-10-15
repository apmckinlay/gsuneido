// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package list

type Equalable interface {
	Equal(x any) bool
}

// Equalable could be generic and typed
// but Value has Equal(any)

// List is a list of values
type List[T Equalable] struct {
	List []T
}

// Push adds a T to the end of the list
func (il *List[T]) Push(v T) {
	il.List = append(il.List, v)
}

// Pop removes the last element of the list
func (il *List[T]) Pop() {
	var zero T
	il.List[len(il.List)-1] = zero // for gc
	il.List = il.List[:len(il.List)-1]
}

// Has returns true if the list contains
func (il *List[T]) Has(v T) bool {
	for _, x := range il.List {
		if x.Equal(v) {
			return true
		}
	}
	return false
}

// Remove deletes the first occurrence of a value
// and returns true if the value was found, otherwise false.
func (il *List[T]) Remove(v T) bool {
	for i, x := range il.List {
		if x.Equal(v) {
			var zero T
			copy(il.List[i:], il.List[i+1:])
			il.List[len(il.List)-1] = zero // for gc
			il.List = il.List[:len(il.List)-1]
			return true
		}
	}
	return false
}
