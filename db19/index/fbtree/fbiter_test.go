// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

import (
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFbiterEmpty(*testing.T) {
	store := stor.HeapStor(256)
	fb := CreateFbtree(store, nil)
	it := fb.Iterator()
	it.Next()
	assert.That(it.Eof())
	it.Next()
	assert.That(it.Eof())
	it.Rewind()
	it.Next()
	assert.That(it.Eof())
}

func TestFbiter(*testing.T) {
	const n = 1000
	var data [n]string
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string { return data[i-1] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	randKey := str.UniqueRandomOf(3, 6, "abcde")
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	sort.Strings(data[:])
	store := stor.HeapStor(8192)
	bldr := Builder(store)
	for i, k := range data {
		bldr.Add(k, uint64(i+1)) // +1 to avoid zero
	}
	fb := bldr.Finish()

	// test Iterator Next
	it := fb.Iterator()
	for i := 0; i < n; i++ {
		it.Next()
		assert.This(it.curOff).Is(i + 1)
		assert.This(it.curKey).Is(data[i])
	}
	it.Next()
	assert.That(it.Eof())

	// test Iterator Prev
	it = fb.Iterator()
	for i := n - 1; i >= 0; i-- {
		it.Prev()
		assert.This(it.curOff).Is(i + 1)
		assert.This(it.curKey).Is(data[i])
	}
	it.Prev()
	assert.That(it.Eof())

	// test Seek between keys
	for i, k := range data {
		k += "0" // increment to nonexistent
		assert.That(!it.Seek(k))
		if i+1 < len(data) {
			assert.This(it.curOff).Is(i + 2)
			assert.This(it.curKey).Is(data[i+1])
		} else {
			assert.That(it.Eof())
		}
	}

	// test Seek & Next
	for i, k := range data {
		assert.That(it.Seek(k))
		assert.This(it.curOff).Is(i + 1)
		assert.This(it.curKey).Is(k)
		it.Next()
		if i+1 < len(data) {
			assert.This(it.curOff).Is(i + 2)
			assert.This(it.curKey).Is(data[i+1])
		} else {
			assert.That(it.Eof())
		}
	}

	// test Seek & Prev
	for i, k := range data {
		assert.That(it.Seek(k))
		assert.This(it.curOff).Is(i + 1)
		assert.This(it.curKey).Is(k)
		it.Prev()
		if i-1 >= 0 {
			assert.This(it.curOff).Is(i)
			assert.This(it.curKey).Is(data[i-1])
		} else {
			assert.That(it.Eof())
		}
	}

	assert.That(!it.Seek("")) // before first
	assert.This(it.curOff).Is(1)
	assert.This(it.curKey).Is(data[0])

	assert.That(!it.Seek("~")) // after last
	assert.That(it.Eof())
}

func TestFbiterToChunk(t *testing.T) {
	assert := assert.T(t).This
	data := []string{"ant", "cat", "dog"}
	b := fNodeBuilder{}
	for i, k := range data {
		b.Add(k, uint64(i+1), 1) // +1 to avoid zero
	}
	fn := b.Entries()
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string { return data[i-1] }

	fb := &fbtree{}
	fi := fn.iter()
	ci := fi.toChunk(fb, true).(*chunkIter) // leaf
	ci.next()
	assert(ci.i).Is(0)
	assert(ci.c).Is(chunk{{key: "ant", off: 1}, {key: "cat", off: 2},
		{key: "dog", off: 3}})
	ci = fi.toChunk(fb, false).(*chunkIter) // tree
	assert(ci.c).Is(chunk{{key: "", off: 1}, {key: "c", off: 2},
		{key: "d", off: 3}})

	fi.next()
	fi.next()
	assert(fi.offset).Is(2)
	ci = fi.toChunk(fb, false).(*chunkIter)
	assert(ci.off()).Is(2)
}
