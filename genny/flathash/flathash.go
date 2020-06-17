// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package flathash implements a simple hash map
// with linear probing open addressing.
// The flat structure means we can duplicate efficiently for copy-on-write.
// The key type must support ==.
// The zero value of keys is used to identify empty slots.
// The zero value of values is returned for failed searched.
// Users must define h.hash(k K) uint32 to complete the generated code.
package flathash

import (
	"math/bits"

	"github.com/cheekybits/genny/generic"
)

type K generic.Type
type V generic.Type

type KVMap struct {
	slots []slot
	// nitems is the current number of key/values in the map
	nitems int
	// cap is the limit for nitems before resizing
	cap int
	// shift is used by hashToIndex
	shift int
	// mask is used to wrap around
	mask int
}

type slot struct {
	key K
	val V
}

var empty slot

// loadPercent is deliberately low to reduce probing
const loadPercent = 60

// NewKVMap creates a new KVMap with roughly the specified capacity
func NewKVMap(cap int) *KVMap {
	if cap < 9 {
		cap = 9
	}
	size, shift := nextPow2(cap * 100 / loadPercent)
	return &KVMap{slots: make([]slot, size),
		shift: shift, mask: size - 1, cap: size * loadPercent / 100}
}

func nextPow2(n int) (int, int) {
	bl := bits.Len(uint(n))
	return 1 << bl, 32 - bl
}

var setn, setp, getn, getp int

func (h *KVMap) Set(k K, v V) {
	if h.nitems >= h.cap {
		h.grow()
	}
	setn++
	for i := h.hashToIndex(h.hash(k)); ; i = (i + 1) & h.mask {
		setp++
		if h.slots[i].key == empty.key {
			h.nitems++
		} else if h.slots[i].key != k {
			continue
		}
		h.slots[i] = slot{key: k, val: v}
		return
	}
}

func (h *KVMap) Get(k K) V {
	getn++
	for i := h.hashToIndex(h.hash(k)); ; i = (i + 1) & h.mask {
		getp++
		if h.slots[i].key == k {
			return h.slots[i].val
		}
		if h.slots[i].key == empty.key {
			return empty.val
		}
	}
}

func (h *KVMap) hashToIndex(hash uint32) int {
	const phi32 = 2654435769
	x := (hash * phi32)
	return int(x >> h.shift)
}

func (h *KVMap) grow() {
	old := *h
	h.cap *= 2
	h.shift--
	h.mask = h.mask<<1 + 1
	h.nitems = 0
	h.slots = make([]slot, 2*len(old.slots))
	for _, slot := range old.slots {
		if slot.key != empty.key {
			h.Set(slot.key, slot.val)
		}
	}
	if h.nitems != old.nitems {
		panic("flathash grow failed")
	}
}

func (h *KVMap) Dup() *KVMap {
	h2 := *h
	h2.slots = append([]slot(nil), h.slots...)
	return &h2
}

func (h *KVMap) Keys() []K {
	keys := make([]K, h.nitems, 0)
	for _, slot := range h.slots {
		if slot.key != empty.key {
			keys = append(keys, slot.key)
		}
	}
	if len(keys) != h.nitems {
		panic("flathash Keys failed")
	}
	return keys
}
