// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package bits provides bit manipulation functions
package bits

import "math/bits"

// NextPow2 returns the smallest power of 2 >= n
func NextPow2(n uint) int {
	return 1 << (bits.UintSize - bits.LeadingZeros(n-1))
}

// TrailingOnes returns the number of trailing one bits in x
func TrailingOnes(x int) int {
	return bits.TrailingZeros(^uint(x))
}

func Shuffle32(n uint32) uint32 {
	// Add a large, odd constant to break the fixed point at 0.
	// 0x9e3779b9 is derived from the golden ratio and is a common choice.
	n += 0x9e3779b9
	n ^= n >> 16
	n *= 0x85ebca6b
	n ^= n >> 13
	n *= 0xc2b2ae35
	n ^= n >> 16
	return n
}

func Shuffle16(n uint16) uint16 {
	// Add a constant to the initial state to move it away from 0.
	// This can be any number, but a large one is good practice.
	n += 0xda79
	n ^= n >> 7
	n *= 0x6955
	n ^= n >> 9
	n *= 0xde59
	n ^= n >> 8
	return n
}
