// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package roaring implements a subset of Roaring Bitmaps.
// It only supports bitmap and array containers. (not runs)
// It handles 48 bit integers.
package roaring

import (
	"slices"
	"sync"
)

// Bitmap is a roaring bitmap.
// A zero value is ready to use.
type Bitmap struct {
	data []container
}

// blocks is used to recycle full array blocks after converting to bitmap
var blocks sync.Pool

const maxValue = 1 << 48

// container holds a range of 64k bits.
// either as a bitmap (if dense) or an array (if sparse)
// For simplicity we use uint16 for both.
// base << 16 | data[i] = value
type container struct {
	data   []uint16
	base   uint32
	bitmap bool // else array
}

// Add inserts a value i.e. sets its bit
func (b *Bitmap) Add(x uint64) {
	if x >= maxValue {
		panic("roaring: value too large")
	}
	pos, found := b.contPos(x)
	val := uint16(x & 0xFFFF)
	if !found {
		// create new array container and insert it at pos
		cont := container{
			bitmap: false,
			base:   uint32(x >> 16),
			data:   []uint16{val},
		}
		b.data = slices.Insert(b.data, pos, cont)
		return
	}
	cont := &b.data[pos]
	if cont.bitmap {
		addBit(cont.data, val)
	} else { // array container
		// optimize for adding on the end
		if val > cont.data[len(cont.data)-1] {
			if len(cont.data) < 4096 {
				cont.data = append(cont.data, val)
				return
			}
		} else {
			insertPos, found := slices.BinarySearch(cont.data, val)
			if found {
				return // already exists
			}
			if len(cont.data) < 4096 {
				cont.data = slices.Insert(cont.data, insertPos, val)
				return
			}
		}
		// convert to bitmap
		newBlk := b.alloc()
		for _, v := range cont.data {
			addBit(newBlk, v)
		}
		addBit(newBlk, val)
		cont.bitmap = true
		blocks.Put(cont.data) //nolint allocates, but only 24 bytes, not 8192
		cont.data = newBlk
	}
}

func (b *Bitmap) alloc() []uint16 {
	if blk := blocks.Get(); blk != nil {
		return blk.([]uint16)
	}
	return make([]uint16, 4096)
}

// Has returns whether a value is in the bitmap i.e. its bit is set
func (b *Bitmap) Has(x uint64) bool {
	if x >= maxValue {
		panic("roaring: value too large")
	}
	pos, found := b.contPos(x)
	if !found {
		return false
	}
	val := uint16(x & 0xffff)
	cont := &b.data[pos]
	if cont.bitmap {
		return hasBit(cont.data, val)
	} else {
		_, found := slices.BinarySearch(cont.data, val)
		return found
	}
}

// contPos returns the position of the container for a value
// and whether it exists
func (b *Bitmap) contPos(x uint64) (int, bool) {
	base := uint32(x >> 16)
	// based on Go slices.BinarySearch
	n := len(b.data)
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		if b.data[h].base < base {
			i = h + 1
		} else {
			j = h
		}
	}
	return i, i < n && b.data[i].base == base
}

// addBit adds a bit to a bitmap block
func addBit(blk []uint16, x uint16) {
	idx := x >> 4
	bit := uint16(1) << (x & 15)
	blk[idx] |= bit
}

// hasBit checks if a bit is set in a bitmap block
func hasBit(blk []uint16, x uint16) bool {
	idx := x >> 4
	bit := uint16(1) << (x & 15)
	return (blk[idx] & bit) != 0
}
