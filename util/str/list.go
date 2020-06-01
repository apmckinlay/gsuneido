// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package str

import "sort"

// ListHas returns true if the list contains the string, false otherwise
func ListHas(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

// ListRemove returns the list with the string removed (if present)
func ListRemove(list []string, str string) []string {
	for i, s := range list {
		if s == str {
			copy(list[i:], list[i+1:])
			list[len(list)-1] = "" // for gc
			return list[:len(list)-1]
		}
	}
	return list
}

// ListReverse reverses the elements of a list of strings
func ListReverse(list []string) []string {
	// could generalize by passing in swap function, like rand.Shuffle
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

// ListUnique
func ListUnique(in []string) []string {
	sort.Strings(in)
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] != in[i] {
			j++
			in[j] = in[i]
		}
	}
	return in[:j+1]
}
