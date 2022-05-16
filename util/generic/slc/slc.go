// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package slc contains additions to x/exp/slices
package slc

import "golang.org/x/exp/slices"

func IndexFn[E any](list []E, e E, eq func(E, E) bool) int {
	for i, e2 := range list {
		if eq(e2, e) {
			return i
		}
	}
	return -1
}

func ContainsFn[E any, S ~[]E](s S, e E, eq func(E, E) bool) bool {
	for _, x := range s {
		if eq(x,e) {
			return true
		}
	}
	return false
}

// HasPrefix returns true if list2 is a prefix of list
func HasPrefix[E comparable](list, list2 []E) bool {
	if len(list) < len(list2) {
		return false
	}
	for i := range list2 {
		if list[i] != list2[i] {
			return false
		}
	}
	return true
}

// Without returns a new slice
// with any occurences of a given value removed,
// maintaining the existing order.
func Without[S ~[]E, E comparable](list S, str E) S {
	dest := make(S, 0, len(list))
	for _, s := range list {
		if s != str {
			dest = append(dest, s)
		}
	}
	return dest
}

// Replace1 returns a new list with occurrences of from replaced by to
func Replace1[S ~[]E, E comparable](list S, from, to E) S {
	list2 := make(S, len(list))
	for i := range list {
		if list[i] == from {
			list2[i] = to
		} else {
			list2[i] = list[i]
		}
	}
	return list2
}

// Replace returns a new list with occurrences of from
// replaced by the corresponding value from to
func Replace[S ~[]E, E comparable](list, from, to S) S {
	list2 := make(S, len(list))
	for i := range list {
		if j := slices.Index(from, list[i]); j != -1 {
			list2[i] = to[j]
		} else {
			list2[i] = list[i]
		}
	}
	return list2
}

// Same returns true if x and y are the same slice
func Same[E any](x, y []E) bool {
	return len(x) > 0 && len(y) > 0 && len(x) == len(y) && &x[0] == &y[0]
}

func Reverse[S ~[]E, E any](x S) S {
	for lo, hi := 0, len(x)-1; lo < hi; {
		tmp := x[lo]
		x[lo] = x[hi]
		x[hi] = tmp
		lo++
		hi--
	}
	return x
}

// Fill sets all the elements of data to value
func Fill[E any](data []E, value E) {
	for i := 0; i < len(data); i++ {
		data[i] = value
	}
}

// Grow grows the buffer to guarantee space for n more bytes.
// NOTE: Unlike x/exp/slices.Grow it extends the length.
// Using append and make like x/exp/slices.Grow assuming that is optimized.
func Grow[S ~[]E, E any](s S, n int) S {
	if n <= 0 {
		return s
	}
	return append(s, make(S, n)...)
}
