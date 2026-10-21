// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package multi implements an immutable fixed size slice
// that is optimized for the case of a single element.
// It is like a "small string optimization" for slices.
package multi

import (
	"unsafe"
)

// Ptrs holds either a single pointer (n=1) or multiple pointers.
// @immutable
type Ptrs[T any] struct {
	p unsafe.Pointer
	n int
}

// New returns a Ptrs with a single element.
// No allocation.
func New[T any](p *T) Ptrs[T] {
	return Ptrs[T]{p: unsafe.Pointer(p), n: 1}
}

// Len returns the number of elements in the Ptrs.
func (ps Ptrs[T]) Len() int {
	return ps.n
}

// Get returns the i'th element in the Ptrs.
// It panics if i is out of range.
func (ps Ptrs[T]) Get(i int) *T {
	if uint(i) >= uint(ps.n) {
		panic("Ptrs.Get: index out of range")
	}
	if ps.n == 1 {
		return (*T)(ps.p)
	}
	return *(**T)(unsafe.Add(ps.p, uintptr(i)*unsafe.Sizeof((*T)(nil))))
}

// With returns a new Ptrs with the elements of ps followed by p.
// It allocates a new slice unless ps is empty.
func With[T any](ps Ptrs[T], p *T) Ptrs[T] {
	if ps.n == 0 {
		return New(p)
	}
	z := make([]*T, ps.n+1)
	ps.copyTo(z)
	z[ps.n] = p
	return Ptrs[T]{p: unsafe.Pointer(unsafe.SliceData(z)), n: len(z)}
}

// Join returns a new Ptrs with the elements of x followed by the elements of y.
// It allocates a new slice unless one of the inputs is empty.
func Join[T any](x, y Ptrs[T]) Ptrs[T] {
	if x.n == 0 {
		return y
	}
	if y.n == 0 {
		return x
	}
	z := make([]*T, x.n+y.n)
	x.copyTo(z)
	y.copyTo(z[x.n:])
	return Ptrs[T]{p: unsafe.Pointer(unsafe.SliceData(z)), n: len(z)}
}

func (ps Ptrs[T]) copyTo(dst []*T) {
	if ps.n == 1 {
		dst[0] = (*T)(ps.p)
	} else {
		copy(dst, unsafe.Slice((**T)(ps.p), ps.n))
	}
}

func Nils[T any](n int) Ptrs[T] {
	if n == 1 {
		return New[T](nil)
	}
	ptrs := make([]*T, n)
	return Ptrs[T]{p: unsafe.Pointer(unsafe.SliceData(ptrs)), n: len(ptrs)}
}
