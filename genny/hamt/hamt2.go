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

func (ht ItemHamt) Write(st *stor.Stor) uint64 {
	size := 0
	ht.ForEach(func(it Item) {
		size += it.storSize()
	})
	if size == 0 {
		return 0
	}
	size += 3 + cksum.Len
	off, buf := st.Alloc(size)
	w := stor.NewWriter(buf)
	w.Put3(size)
	ht.ForEach(func(it Item) {
		it.Write(w)
	})
	assert.That(w.Len() == size-cksum.Len)
	cksum.Update(buf)
	return off
}

func (ht ItemHamt) Read(st *stor.Stor, off uint64) ItemHamt {
	if off == 0 {
		return ht
	}
	buf := st.Data(off)
	size := stor.NewReader(buf).Get3()
	cksum.MustCheck(buf[:size])
	r := stor.NewReader(buf[3 : size-cksum.Len])
	for r.Remaining() > 0 {
		ht.Put(ReadItem(st, r))
	}
	return ht
}
