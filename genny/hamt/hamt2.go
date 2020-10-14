// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamt

import (
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

func (ht ItemHamt) Write(st *stor.Stor, prevOff uint64,
	filter func(it Item) bool) uint64 {
	size := 0
	ht.ForEach(func(it Item) {
		if filter(it) {
			size += it.storSize()
		}
	})
	if size == 0 {
		return prevOff
	}
	size += 3 + 5 + cksum.Len
	off, buf := st.Alloc(size)
	w := stor.NewWriter(buf)
	w.Put3(size)
	w.Put5(prevOff)
	ht.ForEach(func(it Item) {
		if filter(it) {
			it.Write(w)
		}
	})
	assert.That(w.Len() == size-cksum.Len)
	cksum.Update(buf)
	return off
}

func ReadItemChain(st *stor.Stor, off uint64) ItemHamt {
	ht := ItemHamt{}.Mutable()
	for off != 0 {
		off = ht.read(st, off)
	}
	return ht.Freeze()
}

func (ht ItemHamt) read(st *stor.Stor, off uint64) uint64 {
	buf := st.Data(off)
	size := stor.NewReader(buf).Get3()
	cksum.MustCheck(buf[:size])
	r := stor.NewReader(buf[3 : size-cksum.Len])
	prevOff := r.Get5()
	for r.Remaining() > 0 {
		it := ReadItem(st, r)
		if _, ok := ht.Get(ItemKey(it)); !ok {
			ht.Put(it)
		}
	}
	return prevOff
}
