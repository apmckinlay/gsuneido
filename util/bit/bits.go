// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package bits provides bit manipulation functions
package bit

import (
	"math/bits"
	"math/rand/v2"

	"github.com/apmckinlay/gsuneido/util/assert"
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

// Mix scrambles the low bitLen bits of x. Higher bits of x are ignored.
// The result is a bijection on [0, 2^bitLen)
// distinct masked inputs give distinct outputs, all within that range.
// bitLen must be in [0, 64].
func Mix(x uint64, bitLen int) uint64 {
	mask := uint64(1)<<bitLen - 1
	x &= mask
	shift := max(1, bitLen/2)
	x ^= x >> shift
	x = (x * 0xbf58476d1ce4e5b9) & mask
	x ^= x >> shift
	return x
}

//-------------------------------------------------------------------

type Gen struct {
	state      uint64
	rangeLimit uint64
	bitLen     int
	mask       uint64
	inc        uint64
	key        uint64
}

// NewGen returns a generator for a pseudo-random sequence in [0, rangeLimit)
// It uses a Linear Congruential Generator (LCG) and cycle walking.
// The sequence visits each value in [0, rangeLimit) exactly once before repeating.
func NewGen(rnd *rand.Rand, rangeLimit uint64) *Gen {
	assert.That(rangeLimit > 0)
	bitLen := bits.Len64(rangeLimit - 1)
	mask := uint64(1)<<bitLen - 1
	return &Gen{
		state:      rnd.Uint64() & mask,
		rangeLimit: rangeLimit,
		bitLen:     bitLen,
		mask:       mask,
		inc:        rnd.Uint64() | 1,
		key:        rnd.Uint64() & mask,
	}
}

func (g *Gen) Next() uint64 {
	// Advance a full period LCG (state) and apply Mix to the output
	// separately, then cycle walk to stay within rangeLimit.
	const a uint64 = 6364136223846793005
	for {
		x := Mix(g.state, g.bitLen) ^ g.key
		g.state = (a*g.state + g.inc) & g.mask
		if x < g.rangeLimit {
			return x
		}
	}
}
