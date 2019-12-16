// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ints supplies miscellaneous functions for integers
package ints

const (
	MinUint = uint(0)           // binary: all zeroes
	MaxUint = ^MinUint          // binary: all ones
	MaxInt  = int(MaxUint >> 1) // binary: all ones except high bit
	MinInt  = ^MaxInt           // binary: all zeroes except high bit
)

// Fill sets all the elements of data to value
func Fill(data []int, value int) {
	for i := 0; i < len(data); i++ {
		data[i] = value
	}
}

// Index returns the index of the first occurrence of value
// or else -1 if the value is not found.
func Index(data []int, value int) int {
	for i, v := range data {
		if v == value {
			return i
		}
	}
	return -1
}

// Compare returns -1 if x < y, 0 if x == y, and +1 if x > y
func Compare(x int, y int) int {
	if x < y {
		return -1
	} else if x > y {
		return +1
	} else {
		return 0
	}
}

// Min returns the smaller of two int's
func Min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the larger of two int's
func Max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
