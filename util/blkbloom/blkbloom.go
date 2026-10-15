// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package blkbloom implements a simple blocked Bloom filter.
// The data is divided into cache-line sized blocks. Each item hashes to
// one block and sets k=4 bits within it, so a lookup touches a single block.
package blkbloom

import (
	"math/bits"
	"unsafe"
)

const (
	blockBytes  = 64 // cache line
	blockBits   = blockBytes * 8
	bitsPerItem = 10 // gives about 1.3% false positive rate
)

type block [8]uint64

// Bloom is a blocked Bloom filter sized for a hint of items.
type Bloom []block

// New creates a filter sized for approximately sizeHint items,
// using about 10 bits per item.
func New(sizeHint int) Bloom {
	nbits := sizeHint * bitsPerItem
	nblocks := max(1, (nbits+blockBits-1)/blockBits)
	data := alignedSlice[block](nblocks, blockBytes)
	return data
}

func alignedSlice[T any](n, align int) []T {
	var zero T
	size := int(unsafe.Sizeof(zero))
	buf := make([]T, n+align/size)
	addr := uintptr(unsafe.Pointer(&buf[0]))
	off := int(-addr&uintptr(align-1)) / size
	return buf[off : off+n : off+n]
}

// Add inserts an item
func (b Bloom) Add(h uint64) {
	w, h := b.locate(h)
	h1, h2 := uint32(h), uint32(h>>32)|1
	p := h1 & (blockBits - 1)
	w[p/64] |= 1 << (p % 64)
	p = (h1 + h2) & (blockBits - 1)
	w[p/64] |= 1 << (p % 64)
	p = (h1 + 2*h2) & (blockBits - 1)
	w[p/64] |= 1 << (p % 64)
	p = (h1 + 3*h2) & (blockBits - 1)
	w[p/64] |= 1 << (p % 64)
}

// Has reports whether the item may be in the filter.
// It can return false positives but never false negatives.
func (b Bloom) Has(h uint64) bool {
	w, h := b.locate(h)
	h1, h2 := uint32(h), uint32(h>>32)|1
	p0 := h1 & (blockBits - 1)
	p1 := (h1 + h2) & (blockBits - 1)
	p2 := (h1 + 2*h2) & (blockBits - 1)
	p3 := (h1 + 3*h2) & (blockBits - 1)
	return (w[p0/64]>>(p0%64))&
		(w[p1/64]>>(p1%64))&
		(w[p2/64]>>(p2%64))&
		(w[p3/64]>>(p3%64))&1 != 0
}

const phi64 = 0x9e3779b97f4a7c15

// locate selects a block and returns it
// along with a hash for the bit positions within the block.
func (b Bloom) locate(h uint64) (*block, uint64) {
	h *= phi64
	// use Lemire's fastrange to avoid `%`
	hi, _ := bits.Mul64(h, uint64(len(b)))
	return &b[hi], (h ^ h>>32) * 0xbf58476d1ce4e5b9
}
