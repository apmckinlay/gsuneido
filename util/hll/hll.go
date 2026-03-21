// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package hll implements HyperLogLog cardinality estimation for strings.
package hll

import (
	"hash/maphash"
	"math"
	"math/bits"
)

// defaultPrecision of 14 gives standard error of about .8 %
const defaultPrecision = 14

// HLL is a HyperLogLog sketch for approximate distinct counting.
//
// It uses 64-bit hashes from maphash and keeps a single maphash.Hash
// instance to avoid per-Add allocations.
type HLL struct {
	p         uint8
	m         uint32
	registers []uint8
	seed         maphash.Seed
}

// New creates a new HyperLogLog with default precision.
func New() *HLL {
	return NewWithPrecision(defaultPrecision)
}

func NewWithPrecision(p uint8) *HLL {
	if p < 4 || p > 18 {
		panic("hll precision out of range")
	}
	m := uint32(1) << p
	return &HLL{p: p, m: m, registers: make([]uint8, m), seed: maphash.MakeSeed()}
}

// Add hashes s and adds it to the sketch.
func (h *HLL) Add(s string) {
	x := maphash.String(h.seed, s)

	i := x >> (64 - h.p)
	rank := bits.LeadingZeros64(x<<h.p) + 1
	maxRank := int(64-h.p) + 1
	if rank > maxRank {
		rank = maxRank
	}

	if uint8(rank) > h.registers[i] {
		h.registers[i] = uint8(rank)
	}
}

// Count returns the estimated number of distinct values.
func (h *HLL) Count() uint64 {
	m := float64(h.m)
	sum := 0.0
	zeros := 0
	for _, r := range h.registers {
		sum += math.Ldexp(1, -int(r))
		if r == 0 {
			zeros++
		}
	}

	estimate := alpha(h.m) * m * m / sum

	if estimate <= 2.5*m && zeros > 0 {
		estimate = m * math.Log(m/float64(zeros))
	} else if estimate > (1.0/30.0)*twoTo64 {
		estimate = -twoTo64 * math.Log(1-estimate/twoTo64)
	}

	if estimate < 0 {
		return 0
	}
	return uint64(estimate + 0.5)
}

const twoTo64 = 18446744073709551616.0

func alpha(m uint32) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		mf := float64(m)
		return 0.7213 / (1 + 1.079/mf)
	}
}
