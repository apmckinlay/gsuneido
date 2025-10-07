// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// test for both ixbuf.Iterator and btree.Iterator

var itoa = strconv.Itoa

func TestIterRange(*testing.T) {
	start := 1000
	limit := 9999

	ib := &ixbuf.T{}
	testIterEmpty(ib.Iterator())
	for i := start; i <= limit; i++ {
		key := itoa(i)
		ib.Insert(key, uint64(i))
	}
	testIterRange(ib.Iterator())

	store := stor.HeapStor(8192)
	testIterEmpty(btree.CreateBtree(store, nil).Iterator())
	bldr := btree.Builder(store)
	for i := start; i <= limit; i++ {
		key := itoa(i)
		assert.That(bldr.Add(key, uint64(i)))
	}
	bt := bldr.Finish()
	btree.GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return itoa(int(i))
	}
	testIterRange(bt.Iterator())
}

func testIterEmpty(it iface.Iter) {
	it.Rewind()
	it.Next()
	assert.That(it.Eof())
	it.Prev()
	assert.That(it.Eof())

	it.Rewind()
	it.Prev()
	assert.That(it.Eof())
	it.Next()
	assert.That(it.Eof())
}

func testIterRange(it iface.Iter) {
	test := func(fn func(), expected, delta int) {
		for range 10 {
			fn()
			_, off := it.Cur()
			assert.This(off).Is(expected)
			if delta == 0 {
				return
			}
			expected += delta
		}
	}
	test(it.Next, 1000, +1)
	it.Rewind()
	test(it.Prev, 9999, -1)

	it.Range(Range{Org: itoa(0), End: itoa(99999)})
	test(it.Next, 1000, +1)
	it.Rewind()
	test(it.Prev, 9999, -1)

	for range 10000 {
		// org and end must be at least 10 apart
		org := 1000 + rand.Intn(4490)
		end := 9999 - rand.Intn(4490)
		it.Range(Range{Org: itoa(org), End: itoa(end)})
		test(it.Next, org, +1)
		it.Rewind()
		test(it.Prev, end-1, -1)
	}
}
