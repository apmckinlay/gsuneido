// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package slc contains additions to the standard slices package
package slc

import (
	"cmp"
	"unsafe"
)

func LastIndex[E comparable](list []E, e E) int {
	for i := len(list) - 1; i >= 0; i-- {
		if list[i] == e {
			return i
		}
	}
	return -1
}

func IndexFn[E, E2 any](list []E, e2 E2, eq func(E, E2) bool) int {
	for i, e := range list {
		if eq(e, e2) {
			return i
		}
	}
	return -1
}

func ContainsFn[E any, S ~[]E](s S, e E, eq func(E, E) bool) bool {
	for _, x := range s {
		if eq(x, e) {
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
// omitting any occurrences of the given value,
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

// WithoutFn returns a slice omitting any values where fn returns true,
// maintaining the existing order.
// If there are no changes, it returns the original slice.
func WithoutFn[S ~[]E, E any](list S, fn func(E) bool) S {
	result := list
	same := true
	for i, s := range list {
		if fn(s) {
			if same {
				result = make(S, 0, len(list))
				result = append(result, list[:i]...)
				same = false
			}
		} else if !same {
			result = append(result, s)
		}
	}
	return result
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

// Same returns true if x and y are the same slice
func Same[E any](x, y []E) bool {
	return len(x) == len(y) &&
		unsafe.SliceData(x) == unsafe.SliceData(y)
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
	for i := range len(data) {
		data[i] = value
	}
}

// Grow grows the buffer to guarantee space for n more elements.
// NOTE: Unlike x/exp/slices.Grow it extends the length, as well as the capacity
// Using append and make like slices.Grow assuming that is optimized.
func Grow[S ~[]E, E any](s S, n int) S {
	if n <= 0 {
		return s
	}
	return append(s, make(S, n)...)
}

// Allow increases the length of the slice to guarantee space for n elements.
// NOTE: Unlike x/exp/slices.Grow it extends the length, not the capacity.
// Using append and make like x/exp/slices.Grow assuming that is optimized.
func Allow[S ~[]E, E any](s S, n int) S {
	return Grow(s, n-len(s))
}

// Shrink returns a copy of the list with the excess capacity removed
// or the original if there is no excess capacity.
func Shrink[S ~[]E, E any](s S) S {
	if len(s) == cap(s) {
		return s
	}
	t := make(S, len(s))
	copy(t, s)
	return t
}

// With returns a copy of the list with the values appended.
// Unlike append, it allocates just the right size, no extra.
func With[S ~[]E, E any](s1 S, s2 ...E) S {
	result := make(S, len(s1)+len(s2))
	copy(result, s1)
	copy(result[len(s1):], s2)
	return result
}

func Swap[E any](data []E, i, j int) {
	data[i], data[j] = data[j], data[i]
}

// Empty returns true if the list is empty but not nil
func Empty[E any](s []E) bool {
	return s != nil && len(s) == 0
}

func Repeat[E any](x E, n int) []E {
	result := make([]E, n)
	Fill(result, x)
	return result
}

// MapFn returns a new slice
// with each element the result of calling fn on the original element.
func MapFn[E1, E2 any](list []E1, fn func(E1) E2) []E2 {
	if list == nil {
		return nil
	}
	dest := make([]E2, len(list))
	for i, s := range list {
		dest[i] = fn(s)
	}
	return dest
}

// Min returns the minimum value in the list
func Min[E cmp.Ordered](list []E) E {
	min := list[0]
	for _, x := range list {
		if x < min {
			min = x
		}
	}
	return min
}

// Max returns the maximum value in the list
func Max[E cmp.Ordered](list []E) E {
	max := list[0]
	for _, x := range list {
		if x > max {
			max = x
		}
	}
	return max
}

func StartsWith[E comparable](list []E, e E) bool {
	return len(list) > 0 && list[0] == e
}

// Clone is like slices.Clone except that
// it guarantees that the returned slice does not reference the original.
// slices.Clone for a zero length slice still references the original
// which can cause memory retention issues.
// Also, the result is the exact size, no extra.
func Clone[E any](list []E) []E {
	if list == nil {
		return nil
	}
	dup := make([]E, len(list))
	copy(dup, list)
	return dup
}

// CommonPrefixLen returns the length of the common prefix of two slices
func CommonPrefixLen[E comparable](s, t []E) int {
	for i := 0; ; i++ {
		if i >= len(s) || i >= len(t) || s[i] != t[i] {
			return i
		}
	}
}

// HasDup returns true if the slice contains any duplicate values
func HasDup[E comparable](list []E) bool {
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[i] == list[j] {
				return true
			}
		}
	}
	return false
}
