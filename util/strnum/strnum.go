// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package strnum provides a helper function for debugging that maps strings to sequential numbers
package strnum

import "sync"

// Helper for debugging - maps strings to sequential numbers
var (
	strNumMap     = make(map[string]int)
	strNumCounter int
	strNumMutex   sync.Mutex
)

// Num returns a sequential number for a string.
// The same string always returns the same number.
// Useful for shortening debugging output.
func Num(s string) int {
	strNumMutex.Lock()
	defer strNumMutex.Unlock()

	if num, ok := strNumMap[s]; ok {
		return num
	}

	strNumCounter++
	strNumMap[s] = strNumCounter
	return strNumCounter
}
