// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package check

import "golang.org/x/exp/slices"

// set is a specialized set of strings in a slice.
// By only ever appending, we make it an immutable data structure.
// It is designed to be used "stack-wise", added to by called functions.
// Warning: operations are O(N) so should only be used with small sizes.
type set []string

// with extends a set with a new string.
// BEWARE: start = set(); a = start.with('x'); b = start.with('y')
// will overwrite, both a and b will end up as ['y']
// Instead, do: a = start.with('x'); b = start.cow().with('y')
func (ls set) with(s string) set {
	if ls.has(s) {
		return ls
	}
	return append(ls, s)
}

// has returns whether the set contains a string
func (ls set) has(s string) bool {
	for _, t := range ls {
		if s == t {
			return true
		}
	}
	return false
}

// union returns ls extended with elements of ls2 not already in ls
func (ls set) union(ls2 set) set {
	result := ls
	for _, s := range ls2 {
		if !ls.has(s) {
			result = append(result, s)
		}
	}
	return result
}

// intersect returns ls with elements of ls2 removed.
// Should be used like:
//
//	x = x.intersect(y)
//
// WARNING: don't use this on a shared set (use copy, cow is not sufficient)
func (ls set) intersect(ls2 set) set {
	dst := 0
	for _, s := range ls {
		if ls2.has(s) {
			ls[dst] = s
			dst++
		}
	}
	return ls[:dst]
}

// unionIntersect returns ls extended with the intersection of ls1 and ls2.
// Assumes that ls1 and ls2 are extensions of ls
func (ls set) unionIntersect(ls1, ls2 set) set {
	// only need to look at what was added
	ls1 = ls1[len(ls):]
	ls2 = ls2[len(ls):]
	for _, s := range ls1 {
		if ls2.has(s) {
			ls = append(ls, s)
		}
	}
	return ls
}

// cow (copy on write) returns the set with the capacity set to the length
// so that any extension (append) will be forced to re-alloc
// making a separate copy
func (ls set) cow() set {
	return ls[:len(ls):len(ls)]
}

func (ls set) copy() set {
	return slices.Clone(ls)
}
