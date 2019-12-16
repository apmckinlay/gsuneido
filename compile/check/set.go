// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package check

// set is a set of strings in a slice.
// By only ever appending, we make it an immutable data structure.
// Warning: operations are O(N) so should only be used with small sizes.
type set []string

// with extends a set with a new string
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

// unionIntersect returns ls extended with the intersection of ls1 and ls2
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
