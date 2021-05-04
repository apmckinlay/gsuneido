// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package strs has miscellaneous functions for slices of strings
package strs

// Equal returns true if list2 is equal to list
func Equal(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// Cow returns the slice with the capacity set to the length
// so that append will allocate a new slice. (copy on write)
func Cow(ss []string) []string {
	n := len(ss)
	return ss[:n:n]
}
