// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package flathash

import (
	"math/bits"

	"github.com/cheekybits/genny/generic"
)

type Item generic.Type
type Key generic.Type

type ItemHtbl struct {
	slots []*Item
	// nitems is the current number of key/values in the map
	nitems int
	// cap is the limit for nitems before resizing
	cap int
	// shift is used by hashToIndex
	shift int
	// mask is used to wrap around
	mask int
}

// NewItemHtbl creates a new ItemHtbl with roughly the specified capacity
func NewItemHtbl(cap int) *ItemHtbl {
	nextPow2 := func(n int) (int, int) {
		bl := bits.Len(uint(n))
		return 1 << bl, 32 - bl
	}
	const loadPercent = 70
	if cap < 9 {
		cap = 9
	}
	size, shift := nextPow2(cap * 100 / loadPercent)
	return &ItemHtbl{slots: make([]*Item, size),
		shift: shift, mask: size - 1, cap: size * loadPercent / 100}
}

func (h *ItemHtbl) Put(item *Item) {
	if h.nitems >= h.cap {
		h.grow()
	}
	key := h.keyOf(item)
	for i := h.hashToIndex(h.hash(key)); ; i = (i + 1) & h.mask {
		if h.slots[i] == nil {
			h.nitems++
		} else if h.keyOf(h.slots[i]) != key {
			continue
		}
		h.slots[i] = item
		return
	}
}

func (h *ItemHtbl) Get(key Key) *Item {
	for i := h.hashToIndex(h.hash(key)); ; i = (i + 1) & h.mask {
		if h.slots[i] == nil {
			return nil
		}
		if key == h.keyOf(h.slots[i]) {
			return h.slots[i]
		}
	}
}

func (h *ItemHtbl) hashToIndex(hash uint32) int {
	const phi32 = 2654435769
	x := (hash * phi32)
	return int(x >> h.shift)
}

func (h *ItemHtbl) grow() {
	old := *h
	h.cap *= 2
	h.shift--
	h.mask = h.mask<<1 + 1
	h.nitems = 0
	h.slots = make([]*Item, 2*len(old.slots))
	for _, slot := range old.slots {
		if slot != nil {
			h.Put(slot)
		}
	}
	if h.nitems != old.nitems {
		panic("flathash grow failed")
	}
}

func (h *ItemHtbl) Dup() *ItemHtbl {
	h2 := *h
	h2.slots = append([]*Item(nil), h.slots...)
	return &h2
}

// List returns a list of the keys in the table
func (h *ItemHtbl) List() []Key {
	keys := make([]Key, 0, h.nitems)
	for _, slot := range h.slots {
		if slot != nil {
			keys = append(keys, h.keyOf(slot))
		}
	}
	if len(keys) != h.nitems {
		panic("flathash Keys failed")
	}
	return keys
}

func (h *ItemHtbl) Iter() func() *Item {
	i := -1
	return func() *Item {
		for i++; i < len(h.slots); i++ {
			if h.slots[i] != nil {
				return h.slots[i]
			}
		}
		return nil // end
	}
}
