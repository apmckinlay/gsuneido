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
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string { return data[i] }
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
		bldr.Add(k, uint64(i))
	}
	fb := bldr.Finish()

	// test Iterator Next
	it := fb.Iterator()
	for i := 0; i < n; i++ {
		it.Next()
		assert.This(it.curOff).Is(i)
		assert.This(it.curKey).Is(data[i])
	}
	it.Next()
	assert.That(it.Eof())

	// test Seek
	for i, k := range data {
		assert.That(it.Seek(k))
		assert.This(it.curOff).Is(i)
		assert.This(it.curKey).Is(k)
		it.Next()
		if i+1 < len(data) {
			assert.This(it.curOff).Is(i + 1)
			assert.This(it.curKey).Is(data[i+1])
		} else {
			assert.That(it.Eof())
		}

		k += "0" // increment to nonexistent
		assert.That(!it.Seek(k))
		if i+1 < len(data) {
			assert.This(it.curOff).Is(i + 1)
			assert.This(it.curKey).Is(data[i+1])
		} else {
			assert.That(it.Eof())
		}
	}
	assert.That(!it.Seek("~")) // after last
	assert.That(it.Eof())
}

func TestFnodeToChunk(t *testing.T) {
	data := []string{"ant", "cat", "dog"}
	b := fNodeBuilder{}
	for i, k := range data {
		b.Add(k, uint64(i), 1)
	}
	fn := b.Entries()
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string { return data[i] }

	it := Iterator{fb: &fbtree{}}
	c := it.fnodeToChunk(fn, true)
	assert.T(t).This(c).Is(chunk{{key: "ant", off: 0}, {key: "cat", off: 1},
		{key: "dog", off: 2}})
	c = it.fnodeToChunk(fn, false)
	assert.T(t).This(c).Is(chunk{{key: "", off: 0}, {key: "c", off: 1},
		{key: "d", off: 2}})
}
