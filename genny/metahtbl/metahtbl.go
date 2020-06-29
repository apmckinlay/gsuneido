// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metahtbl

import (
	"math/bits"
	"sort"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/verify"
	"github.com/cheekybits/genny/generic"
)

type Item generic.Type

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
	key := item.Table
	for i := h.hashToIndex(hash.HashString(key)); ; i = (i + 1) & h.mask {
		if h.slots[i] == nil {
			h.nitems++
		} else if h.slots[i].Table != key {
			continue
		}
		h.slots[i] = item
		return
	}
}

func (h *ItemHtbl) Get(key string) *Item {
	if h == nil {
		return nil
	}
	for i := h.hashToIndex(hash.HashString(key)); ; i = (i + 1) & h.mask {
		if h.slots[i] == nil {
			return nil
		}
		if key == h.slots[i].Table {
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
		panic("metahtbl grow failed")
	}
}

func (h *ItemHtbl) Dup() *ItemHtbl {
	h2 := *h
	h2.slots = append([]*Item(nil), h.slots...)
	return &h2
}

// List returns a list of the keys in the table
func (h *ItemHtbl) List() []string {
	keys := make([]string, 0, h.nitems)
	for _, slot := range h.slots {
		if slot != nil {
			keys = append(keys, slot.Table)
		}
	}
	if len(keys) != h.nitems {
		panic("metahtbl List failed")
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

const blockSizeItem = 2000
const perFingerItem = 16

func (h *ItemHtbl) Write(st *stor.Stor) uint64 {
	if h.nitems == 0 {
		off, buf := st.Alloc(2)
		stor.NewWriter(buf).Put2(0)
		return off
	}
	nfingers := 1 + h.nitems/perFingerItem
	size := 2 + 3*nfingers
	iter := h.Iter()
	for it := iter(); it != nil; it = iter() {
		size += it.storSize()
	}
	off, buf := st.Alloc(size)
	w := stor.NewWriter(buf)
	w.Put2(h.nitems)

	keys := h.List()
	sort.Strings(keys)
	w2 := *w
	for i := 0; i < nfingers; i++ {
		w.Put3(0) // leave room
	}
	fingers := make([]int, 0, nfingers)
	for i, k := range keys {
		if i%16 == 0 {
			fingers = append(fingers, w.Len())
		}
		h.Get(k).Write(w)
	}
	verify.That(len(fingers) == nfingers)
	for _, f := range fingers {
		w2.Put3(f) // update with actual values
	}
	return off
}

func ReadItemHtbl(st *stor.Stor, off uint64) *ItemHtbl {
	r := st.Reader(off)
	nitems := r.Get2()
	t := NewItemHtbl(nitems)
	if nitems == 0 {
		return t
	}
	nfingers := 1 + nitems/perFingerItem
	for i := 0; i < nfingers; i++ {
		r.Get3() // skip the fingers
	}
	for i := 0; i < nitems; i++ {
		t.Put(ReadItem(st, r))
	}
	return t
}

//-------------------------------------------------------------------

type ItemPacked struct {
	stor    *stor.Stor
	off     uint64
	buf     []byte
	fingers []ItemFinger
}

type ItemFinger struct {
	table string
	pos   int
}

func NewItemPacked(st *stor.Stor, off uint64) *ItemPacked {
	buf := st.Data(off)
	r := stor.NewReader(buf)
	nitems := r.Get2()
	nfingers := 1 + nitems/perFingerItem
	fingers := make([]ItemFinger, nfingers)
	for i := 0; i < nfingers; i++ {
		fingers[i].pos = r.Get3()
	}
	for i := 0; i < nfingers; i++ {
		fingers[i].table = stor.NewReader(buf[fingers[i].pos:]).GetStr()
	}
	return &ItemPacked{stor: st, off: off, buf: buf, fingers: fingers}
}

func (p ItemPacked) Get(key string) *Item {
	pos := p.binarySearch(key)
	r := stor.NewReader(p.buf[pos:])
	count := 0
	for {
		item := ReadItem(p.stor, r)
		if item.Table == key {
			return item
		}
		count++
		if count > 20 {
			panic("linear search too long")
		}
	}
}

// binarySearch does a binary search of the fingers
func (p ItemPacked) binarySearch(table string) int {
	i, j := 0, len(p.fingers)
	count := 0
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if table >= p.fingers[h].table {
			i = h + 1
		} else {
			j = h
		}
		count++
		if count > 20 {
			panic("binary search too long")
		}
	}
	// i is first one greater, so we want i-1
	return int(p.fingers[i-1].pos)
}

func (p ItemPacked) Offset() uint64 {
	return p.off
}
