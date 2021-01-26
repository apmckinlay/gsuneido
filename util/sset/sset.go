// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package sset (string set) provides set operations on []string.
// If the inputs are sorted then more efficient algorithms are used,
// and the output will also be sorted.
// Operations (other than Optim) do not modify their inputs
// but may return an input unmodified.
package sset

import (
	"sort"

	"github.com/apmckinlay/gsuneido/util/ints"
)

// small is the threshold for switching to sorted algorithms
const small = 32 // roughly the breakeven for Equal

// Optim returns the "optimal" form of a list set - sorted with no duplicates.
//
// WARNING: modifies the original.
func Optim(x []string) []string {
	if len(x) == 0 {
		return x
	}
	sort.Strings(x)
	// remove duplicates (now adjacent)
	dst := 1
	for i := 1; i < len(x); i++ {
		if x[i] != x[i-1] {
			x[dst] = x[i]
			dst++
		}
	}
	return x[:dst]
}

// IsOptim returns true if the list is sorted and has no duplicates.
// It has to do a complete scan
// so it should only be used for larger lists where there's a bigger benefit.
func IsOptim(x []string) bool {
	for i := 1; i < len(x); i++ {
		if x[i] <= x[i-1] {
			return false
		}
	}
	return true
}

// Contains returns true if x contains s
func Contains(x []string, s string) bool {
	// no point checking for sorted because that will do a full scan
	for _, xs := range x {
		if xs == s {
			return true
		}
	}
	return false
}

func Copy(x []string) []string {
	z := make([]string, len(x))
	copy(z, x)
	return z
}

// Equal returns true if x and y contain the same set of strings.
//
// WARNING: unsorted approach requires that x does not contain duplicates
func Equal(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	if Same(x, y) {
		return true
	}
	if len(x) > small && Sorted2(x, y) {
		// faster, single pass approach if sorted
		for i := 0; i < len(x); i++ {
			if x[i] != y[i] {
				return false
			}
		}
		return true
	}
	// slower, unsorted approach
outer:
	for _, xs := range x {
		for _, ys := range y {
			if xs == ys {
				continue outer
			}
		}
		return false // xs wasn't found in y
	}
	return true
}

func Same(x, y []string) bool {
	return len(x) > 0 && len(y) > 0 && len(x) == len(y) && &x[0] == &y[0]
}

// Sorted2 returns true if both lists are sorted.
// It has to do a complete scan
// so it should only be used for larger lists where there's a bigger benefit.
func Sorted2(x, y []string) bool {
	// check both in parallel to fail faster
	n := ints.Max(len(x), len(y))
	for i := 1; i < n; i++ {
		if (i < len(x) && x[i] < x[i-1]) ||
			(i < len(y) && y[i] < y[i-1]) {
			return false
		}
	}
	return true
}

// Union returns a combined list.
// If either input contains duplicates, so will the output.
//
// WARNING: If either x or y is empty, it returns the *original* of the other.
func Union(x, y []string) []string {
	if len(x) == 0 {
		return y
	}
	if len(y) == 0 {
		return x
	}
	z := make([]string, 0, len(x)+len(y))
	if len(x)+len(y) > 2*small && Sorted2(x, y) {
		// faster, sorted single pass approach
		i, j := 0, 0
		for i < len(x) && j < len(y) {
			if x[i] < y[j] {
				z = append(z, x[i])
				i++
			} else if y[j] < x[i] {
				z = append(z, y[j])
				j++

			} else { // x[i] == y[j]
				z = append(z, y[j])
				i++
				j++
			}
		}
		if i < len(x) {
			z = append(z, x[i:]...)
		} else if j < len(y) {
			z = append(z, y[j:]...)
		}
	} else { // slower, unsorted approach
		z = append(z, x...) // copy larger
	outer:
		for _, ys := range y {
			for _, xs := range x {
				if ys == xs {
					continue outer
				}
			}
			z = append(z, ys)
		}
	}
	return z
}

// Difference returns the elements of x that are not in y.
//
// WARNING: If y is empty, it returns the *original* x.
//
// WARNING: duplicates the inputs may give duplicates in the result
func Difference(x, y []string) []string {
	if len(x) == 0 {
		return []string{}
	}
	if len(y) == 0 {
		return x
	}
	if Same(x, y) {
		return []string{}
	}
	z := make([]string, 0, len(x))
	if len(x)+len(y) > 2*small && Sorted2(x, y) {
		// faster, sorted single pass approach
		i, j := 0, 0
		for i < len(x) && j < len(y) {
			if x[i] < y[j] {
				z = append(z, x[i])
				i++
			} else if y[j] < x[i] {
				j++

			} else { // x[i] == y[j]
				i++
				j++
			}
		}
		if i < len(x) {
			z = append(z, x[i:]...)
		}
	} else {
		// slower, unsorted approach
	outer:
		for _, xs := range x {
			for _, ys := range y {
				if ys == xs {
					continue outer
				}
			}
			z = append(z, xs)
		}
	}
	return z

}

// Intersect returns a list of the strings common to the inputs
//
// WARNING: If x and y are the same list, it returns the *original*.
//
// WARNING: duplicates the inputs may give duplicates in the result
func Intersect(x, y []string) []string {
	if Same(x, y) {
		return x
	}
	if len(x) == 0 || len(y) == 0 {
		return []string{}
	}
	z := make([]string, 0, ints.Min(len(x), len(y))/2) // ???
	if len(x)+len(y) > 2*small && Sorted2(x, y) {
		// faster, sorted single pass approach
		i, j := 0, 0
		for i < len(x) && j < len(y) {
			if x[i] < y[j] {
				i++
			} else if y[j] < x[i] {
				j++
			} else { // x[i] == y[j]
				z = append(z, x[i])
				i++
				j++
			}
		}
	} else {
		// slower, unsorted approach
	outer:
		for _, xs := range x {
			for _, ys := range y {
				if xs == ys {
					z = append(z, xs)
					continue outer
				}
			}
		}
	}
	return z
}
