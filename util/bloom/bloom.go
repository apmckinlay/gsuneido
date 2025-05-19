// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package bloom implements a simple Bloom filter
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

type Bloom struct {
	k    int
	bits bitset
}

func New(m, k int) *Bloom {
	return &Bloom{k: k, bits: newBitset(m)}
}

func (b *Bloom) Add(h uint64) {
	h1 := int(uint32(h))
	h2 := int(h >> 32)
	for i := range b.k {
		b.bits.Set((h1 + i*h2) % (len(b.bits) * 64))
	}
}

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

func (b *Bloom) Size() int {
	return len(b.bits) * 64 / 8 // bytes
}

func Calc(n int, p float64) (m int, k int) {
	ln2 := math.Log(2)
	mdivn := -math.Log(p) / (ln2 * ln2)
	m = int(math.Ceil(float64(n) * mdivn))
	k = int(math.Ceil(mdivn * ln2))
	return
}
