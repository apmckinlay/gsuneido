// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package strss has miscellaneous functions for slices of slices of strings
package strss

import "github.com/apmckinlay/gsuneido/util/strs"

// Index returns the position of the first occurrence of the given value,
// or -1 if not found.
func Index(list [][]string, x []string) int {
	for i, s := range list {
		if strs.Equal(s, x) {
			return i
		}
	}
	return -1
}
