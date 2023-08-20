// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package slc contains additions to the standard slices package
package slc

import (
	"slices"

	. "golang.org/x/exp/constraints"
)

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
// omitting any occurences of the given value,
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

// WithoutFn returns a new slice
// omitting any values where fn returns true,
// maintaining the existing order.
func WithoutFn[S ~[]E, E any](list S, fn func(E) bool) S {
	dest := make(S, 0, len(list))
	for _, s := range list {
		if !fn(s) {
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
// replaced by the corresponding value in to.
// If no replacements are done, it returns the original list.
func Replace[S ~[]E, E comparable](list, from, to S) S {
	cloned := false
	for i := range list {
		if j := slices.Index(from, list[i]); j != -1 {
			if !cloned {
				list = slices.Clone(list)
				cloned = true
			}
			list[i] = to[j]
		}
	}
	return list
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

// Grow grows the buffer to guarantee space for n more elements.
// NOTE: Unlike x/exp/slices.Grow it extends the length, not the capacity.
// Using append and make like x/exp/slices.Grow assuming that is optimized.
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

// With returns a copy of the list with the values appended
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
func MapFn[S ~[]E, E any](list S, fn func(E) E) S {
	if list == nil {
		return nil
	}
	dest := make(S, len(list))
	for i, s := range list {
		dest[i] = fn(s)
	}
	return dest
}

// Min returns the minimum value in the list
func Min[E Ordered](list []E) E {
	min := list[0]
	for _, x := range list {
		if x < min {
			min = x
		}
	}
	return min
}

// Max returns the maximum value in the list
func Max[E Ordered](list []E) E {
	max := list[0]
	for _, x := range list {
		if x > max {
			max = x
		}
	}
	return max
}
