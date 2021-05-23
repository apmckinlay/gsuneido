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

// Contains returns whether the list contains the string
// or -1 if not found.
func Contains(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

// Index returns the position of the first occurrence of the given string,
// or -1 if not found.
func Index(list []string, str string) int {
	for i, s := range list {
		if s == str {
			return i
		}
	}
	return -1
}

// HasPrefix returns true if list2 is a prefix of list
func HasPrefix(list, list2 []string) bool {
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

// Without returns a new slice of strings
// with any occurences of a given string removed,
// maintaining the existing order.
func Without(list []string, str string) []string {
	dest := make([]string, 0, len(list))
	for _, s := range list {
		if s != str {
			dest = append(dest, s)
		}
	}
	return dest
}

// Replace returns a new list with occurrences of from
// replaced by the corresponding value from to
func Replace(list, from, to []string) []string {
	list2 := make([]string, len(list))
	for i := range list {
		if j := Index(from, list[i]); j != -1 {
			list2[i] = to[j]
		} else {
			list2[i] = list[i]
		}
	}
	return list2
}
