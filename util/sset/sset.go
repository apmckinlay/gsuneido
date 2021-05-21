// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package sset

//go:generate genny -in ../../genny/set/set.go -out sset2.go -pkg sset gen "T=string"

func eq(x, y string) bool {
	return x == y
}

func StartsWithSet(x, set []string) bool {
	if len(x) < len(set) {
		return false
	}
	for _, xs := range x[:len(set)] {
		if !Contains(set, xs) {
			return false
		}
	}
	return true
}

// Unique returns a list with duplicates removed, retaining the original order.
//
// NOTE: If there are no duplicate it returns the *original* list.
func Unique(x []string) []string {
	i := 0
	n := len(x)
	for i := 1; i < n && !Contains(x[:i], x[i]); i++ {
	}
	if i == n {
		return x // no duplicates
	}
	z := make([]string, 0, len(x))
	copy(z, x[:i])
	for _, s := range x[i:] {
		z = AddUnique(z, s)
	}
	return z
}
