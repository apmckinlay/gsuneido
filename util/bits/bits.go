// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package bits provides bit manipulation functions
package bits

import (
	"math/bits"
	"math/rand/v2"
)

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
		a = 0x9e35 // 40501, satisfies a ≡ 1 (mod 4)
		c = 0x9e37 // 40503, odd
	)
	n = n*a + c
	return n
}

// Mix shuffles x with b bits (bijective)
// This ensures that every input maps to exactly one unique output.
func Mix(x uint64, bitLen int) uint64 {
	m := uint64(1)<<bitLen - 1
	x ^= x >> (bitLen / 2)
	x = (x * 0xbf58476d1ce4e5b9) & (m) // A bijective multiplier
	x ^= x >> (bitLen / 2)
	return x
}

// Cycle uses a Linear Congruential Generator (LCG)
// and cycle walking to produce a pseudo-random sequence in [0, rangeLimit)
func Cycle(currentVal, rangeLimit uint64) uint64 {
	// LCG parameters
	bitLen := bits.Len64(rangeLimit)
	m := uint64(1) << bitLen
	const a uint64 = 6364136223846793005
	const c uint64 = 1442695040888963407

	// LCG
	step := func(x uint64) uint64 {
		return Mix((a*x+c)&(m-1), bitLen)
	}

	// Cycle Walking
	nextState := step(currentVal)
	for nextState >= rangeLimit {
		nextState = step(nextState)
	}

	return nextState
}

type Gen struct {
	cur        uint64
	rangeLimit uint64
	bitLen     int
}

func NewGen(rnd *rand.Rand, rangeLimit uint64) *Gen {
	return &Gen{
		cur:        rnd.Uint64N(rangeLimit),
		rangeLimit: rangeLimit,
		bitLen:     bits.Len64(rangeLimit)}
}

func (g *Gen) Next() uint64 {
	g.cur = Cycle(g.cur, g.rangeLimit)
	return g.cur
}
