// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamt

import (
	"sort"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

// list returns a list of the keys in the table
func (ht ItemHamt) list() []string {
	keys := make([]string, 0, 16)
	ht.ForEach(func(it Item) {
		keys = append(keys, ItemKey(it))
	})
	return keys
}

const blockSizeItem = 2000
const perFingerItem = 16

func (ht ItemHamt) Write(st *stor.Stor) uint64 {
	nitems := 0
	size := 3 + 2
	ht.ForEach(func(it Item) {
		size += it.storSize()
		nitems++
	})
	if nitems == 0 {
		return 0
	}
	nfingers := 1 + nitems/perFingerItem
	size += 3 * nfingers
	off, buf := st.Alloc(size + cksum.Len)
	w := stor.NewWriter(buf)
	w.Put3(size + cksum.Len)
	w.Put2(nitems)

	keys := ht.list()
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
		it, _ := ht.Get(k)
		it.Write(w)
	}
	assert.That(len(fingers) == nfingers)
	assert.That(w.Len() == size)
	for _, f := range fingers {
		w2.Put3(f) // update with actual values
	}
	cksum.Update(buf)
	return off
}

func ReadItemHamt(st *stor.Stor, off uint64) ItemHamt {
	if off == 0 {
		return ItemHamt{}
	}
	buf := st.Data(off)
	r := stor.NewReader(buf)
	size := r.Get3()
	if size == 0 {
		return ItemHamt{}
	}
	cksum.MustCheck(buf[:size])
	nitems := r.Get2()
	t := ItemHamt{}.Mutable()
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
	return t.Freeze()
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
	if off == 0 {
		return &ItemPacked{}
	}
	buf := st.Data(off)
	r := stor.NewReader(buf)
	size := r.Get3()
	cksum.MustCheck(buf[:size])
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

func (p *ItemPacked) MustGet(key string) Item {
	if item, ok := p.Get(key); ok {
		return item
	}
	panic("item not found")
}

func (p *ItemPacked) Get(key string) (Item, bool) {
	var zero Item
	if p.buf == nil {
		return zero, false
	}
	pos := p.binarySearch(key)
	r := stor.NewReader(p.buf[pos:])
	for n := 0; n <= perFingerItem; n++ {
		item := ReadItem(p.stor, r)
		if item.Table == key {
			return item, true
		}
	}
	return zero, false
}

// binarySearch does a binary search of the fingers
func (p *ItemPacked) binarySearch(table string) int {
	i, j := 0, len(p.fingers)
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if table >= p.fingers[h].table {
			i = h + 1
		} else {
			j = h
		}
	}
	// i is first one greater, so we want i-1
	return p.fingers[i-1].pos
}

func (p *ItemPacked) Offset() uint64 {
	return p.off
}
