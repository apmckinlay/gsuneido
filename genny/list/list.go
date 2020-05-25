package list

import "github.com/cheekybits/genny/generic"

// V is the value type.
// It must have an Equal method.
// (Which means it can't be used with raw primitive types.)
type V generic.Type

// VList is a list of values
type VList struct {
	list []V
}

// Push adds a V to the end of the list
func (il *VList) Push(v V) {
	il.list = append(il.list, v)
}

// Pop removes the last element of the list
func (il *VList) Pop() {
	var zero V
	il.list[len(il.list)-1] = zero // for gc
	il.list = il.list[:len(il.list)-1]
}

// Has returns true if the list contains 
func (il *VList) Has(v V) bool {
	for _, x := range il.list {
		if x.Equal(v) {
			return true
		}
	}
	return false
}

// Remove deletes the first occurence of a value
// and returns true if the value was found, otherwise false.
func (il *VList) Remove(v V) bool {
	for i, x := range il.list {
		if x.Equal(v) {
			var zero V
			copy(il.list[i:], il.list[i+1:])
			il.list[len(il.list)-1] = zero // for gc
			il.list = il.list[:len(il.list)-1]
			return true
		}
	}
	return false
}
