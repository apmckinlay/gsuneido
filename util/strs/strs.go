// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package strs has miscellaneous functions for slices of strings
package strs

// Cow returns the slice with the capacity set to the length
// so that append will allocate a new slice. (copy on write)
func Cow(ss []string) []string {
	n := len(ss)
	return ss[:n:n]
}
