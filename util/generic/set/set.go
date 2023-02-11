// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package set

// this provides set operations on []T.
// Operations do not modify their inputs
// but may return an input unmodified, rather than a copy.
// Where applicable, the order of the values is maintained.
// WARNING: these operations will be slow on large sets.

import (
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"golang.org/x/exp/slices"
)

// AddUnique appends s unless it is already in the set
//
// WARNING: it does append so usage must be: x = AddUnique(x, y)
func AddUnique[E comparable, S ~[]E](s S, e E) S {
	if slices.Contains(s, e) {
		return s
	}
	return append(s, e)
}

func AddUniqueFn[E any, S ~[]E](s S, e E, eq func(E, E) bool) S {
	if slc.ContainsFn(s, e, eq) {
		return s
	}
	return append(s, e)
}

// Equal returns true if x and y contain the same set of values.
func Equal[E comparable](x, y []E) bool {
	if len(x) != len(y) {
		return false
	}
	if slc.Same(x, y) {
		return true
	}
	// used is needed to handle duplicates
	used := make([]bool, len(y)) // hopefully on the stack (no alloc)
outer:
	for _, xe := range x {
		for i, ye := range y {
			if xe == ye && !used[i] {
				used[i] = true
                continue outer
            }
		}
		return false // xs wasn't found in y
	}
	return true
}

// Union returns a combined list.
// If either input contains duplicates, so will the output.
// The order of x and y are preserved.
//
// WARNING: If either x or y is empty, it returns the *original* of the other.
func Union[E comparable, S ~[]E](x, y S) S {
	if len(x) == 0 {
		return slices.Clip(y) // so append won't share
	}
	if len(y) == 0 {
		return slices.Clip(x) // so append won't share
	}
	z := make(S, 0, len(x)+len(y))
	z = append(z, x...)
outer:
	for _, ye := range y {
		for _, xe := range x {
			if xe == ye {
				continue outer
			}
		}
		z = append(z, ye)
	}
	return z
}

func UnionFn[E any, S ~[]E](x, y S, eq func(E, E) bool) S {
	if len(x) == 0 {
		return slices.Clip(y) // so append won't share
	}
	if len(y) == 0 {
		return slices.Clip(x) // so append won't share
	}
	z := make(S, 0, len(x)+len(y))
	z = append(z, x...)
outer:
	for _, ye := range y {
		for _, xe := range x {
			if eq(xe, ye) {
				continue outer
			}
		}
		z = append(z, ye)
	}
	return z
}

// Difference returns the elements of x that are not in y.
//
// WARNING: If y is empty, it returns the *original* x.
//
// WARNING: duplicates in the inputs may give duplicates in the result
func Difference[E comparable, S ~[]E](x, y S) S {
	if len(x) == 0 {
		return S{}
	}
	if len(y) == 0 {
		return slices.Clip(x) // so append won't share
	}
	if slc.Same(x, y) {
		return S{}
	}
	z := make(S, 0, len(x))
outer:
	for _, xe := range x {
		for _, ye := range y {
			if xe == ye {
				continue outer
			}
		}
		z = append(z, xe)
	}
	return z
}

// Intersect returns a list of the values common to the inputs,
// the result is in the same order as the first argument (x).
//
// WARNING: If x and y are the same list, it returns the *original*.
//
// WARNING: duplicates in the inputs may give duplicates in the result
func Intersect[E comparable, S ~[]E](x, y S) S {
	if len(x) == 0 || len(y) == 0 {
		return S{}
	}
	if slc.Same(x, y) {
		return slices.Clip(x) // so append won't share
	}
	z := make(S, 0) //, ord.Min(len(x), len(y))/2) // ???
outer:
	for _, xe := range x {
		for _, ye := range y {
			if xe == ye {
				z = append(z, xe)
				continue outer
			}
		}
	}
	return z
}

func IntersectFn[E any, S ~[]E](x, y S, eq func(E, E) bool) S {
	if len(x) == 0 || len(y) == 0 {
		return S{}
	}
	if slc.Same(x, y) {
		return slices.Clip(x) // so append won't share
	}
	z := make(S, 0) //, ord.Min(len(x), len(y))/2) // ???
outer:
	for _, xe := range x {
		for _, ye := range y {
			if eq(xe, ye) {
				z = append(z, xe)
				continue outer
			}
		}
	}
	return z
}

// Subset returns true is y is a subset of x
// i.e. x contains all of y
func Subset[E comparable](x, y []E) bool {
outer:
	for _, ye := range y {
		for _, xe := range x {
			if xe == ye {
				continue outer
			}
		}
		return false // ys wasn't found in x
	}
	return true
}

// Disjoint returns true if x and y have no elements in common.
// i.e. Intersect(x, y) is empty
func Disjoint[E comparable](x, y []E) bool {
	for _, xe := range x {
		for _, ye := range y {
			if xe == ye {
				return false
			}
		}
	}
	return true
}

func StartsWithSet[E comparable](x, y []E) bool {
	if len(x) < len(y) {
		return false
	}
	for _, xs := range x[:len(y)] {
		if !slices.Contains(y, xs) {
			return false
		}
	}
	return true
}

// Unique returns a list with duplicates removed, retaining the original order.
//
// NOTE: If there are no duplicates it returns the *original* list.
func Unique[E comparable, S ~[]E](x S) S {
	i := 0
	n := len(x)
	for i := 1; i < n && !slices.Contains(x[:i], x[i]); i++ {
	}
	if i == n {
		return slices.Clip(x) // no duplicates
	}
	z := make(S, 0, len(x))
	copy(z, x[:i])
	for _, xe := range x[i:] {
		z = AddUnique(z, xe)
	}
	return z
}
