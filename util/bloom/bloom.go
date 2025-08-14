// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package bloom implements a simple Bloom filter.
// It uses double hashing from a single 64-bit hash value, splitting it
// into two 32-bit halves to generate k hash positions.
package bloom

import "math"

type bitset []int64

func newBitset(size int) bitset {
	return make(bitset, (size+63)/64)
}

func (bs bitset) Set(n int) {
	bs[n/64] |= 1 << uint(n%64)
}

func (bs bitset) Get(n int) bool {
	return bs[n/64]&(1<<uint(n%64)) != 0
}

// Bloom is a Bloom filter with k hash functions over a bitset of m bits.
// The hash functions are derived via double hashing from a single 64-bit input.
type Bloom struct {
	k    int
	bits bitset
}

// New creates a Bloom filter with m bits and k hash functions.
// m is the number of bits, not bytes.
func New(m, k int) *Bloom {
	return &Bloom{k: k, bits: newBitset(m)}
}

// Add inserts an item, specified by a 64-bit hash h, into the filter.
// The 64-bit hash is split into two 32-bit values and combined for k positions.
func (b *Bloom) Add(h uint64) {
	h1 := int(uint32(h))
	h2 := int(h >> 32)
	for i := range b.k {
		b.bits.Set((h1 + i*h2) % (len(b.bits) * 64))
	}
}

// Test reports whether an item, specified by 64-bit hash h, may be in the set.
// It can return false positives but never false negatives.
func (b *Bloom) Test(h uint64) bool {
	h1 := int(uint32(h))
	h2 := int(h >> 32)
	for i := range b.k {
		if !b.bits.Get((h1 + i*h2) % (len(b.bits) * 64)) {
			return false
		}
	}
	return true
}

// Size returns the memory usage of the filter in bytes.
func (b *Bloom) Size() int {
	return len(b.bits) * 64 / 8 // bytes
}

// Calc computes recommended parameters for a Bloom filter.
// Given expected number of items n and target false positive rate p (0 < p < 1),
// it returns m (number of bits) and k (number of hash functions).
func Calc(n int, p float64) (m int, k int) {
	ln2 := math.Log(2)
	mdivn := -math.Log(p) / (ln2 * ln2)
	m = int(math.Ceil(float64(n) * mdivn))
	k = int(math.Ceil(mdivn * ln2))
	return
}
