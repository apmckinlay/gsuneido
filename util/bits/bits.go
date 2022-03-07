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
