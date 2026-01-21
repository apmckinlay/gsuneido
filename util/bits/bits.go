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
	// Linear Congruential Generator (LCG) with full period for 32-bit values.
	// For modulus m = 2^32, full period requires:
	// - multiplier a ≡ 1 (mod 4)
	// - increment c is odd
	// These parameters ensure a single cycle of length 2^32.
	const (
		a = 0x9e3779b9 // 2654435769, satisfies a ≡ 1 (mod 4)
		c = 0x9e3779b9 // 2654435769, odd
	)
	n = n*a + c
	return n
}

func Shuffle16(n uint16) uint16 {
	// Linear Congruential Generator (LCG) with full period for 16-bit values.
	// For modulus m = 2^16, full period requires:
	// - multiplier a ≡ 1 (mod 4)
	// - increment c is odd
	// These parameters ensure a single cycle of length 65536.
	const (
		a = 0x9e35  // 40501, satisfies a ≡ 1 (mod 4)
		c = 0x9e37  // 40503, odd
	)
	n = n*a + c
	return n
}
