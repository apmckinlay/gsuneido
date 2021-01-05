// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFbtreeIO(t *testing.T) {
	store := stor.HeapStor(128)
	before := fnode([]byte("helloworld"))
	off := before.putNode(store)
	after := readNode(store, off)
	cksum.MustCheck(after[:len(after)+cksum.Len])
	assert.T(t).This(string(after)).Is(string(before))
}

func TestFbtreeIter(t *testing.T) {
	const n = 1000
	var data [n]string
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 440
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
	i := 0
	iter := fb.Iter(true)
	for k, o, ok := iter(); ok; k, o, ok = iter() {
		assert.T(t).That(strings.HasPrefix(data[i], k))
		assert.T(t).This(o).Is(i)
		i++
	}
	assert.T(t).This(i).Is(n)
}

func TestFbtreeBuilder(t *testing.T) {
	assert := assert.T(t)
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	store := stor.HeapStor(8192)
	bldr := Builder(store)
	limit := 599999
	if testing.Short() {
		limit = 199999
	}
	for i := 100000; i <= limit; i++ {
		key := strconv.Itoa(i)
		bldr.Add(key, uint64(i))
	}
	fb := bldr.Finish()
	fb.Check(nil)

	// iterate
	iter := fb.Iter(true)
	for i := 100000; i <= limit; i++ {
		key := strconv.Itoa(i)
		k, o, ok := iter()
		assert.True(ok)
		assert.True(strings.HasPrefix(key, k))
		assert.This(o).Is(i)
	}
	_, _, ok := iter()
	assert.False(ok)

	// search
	for i := 100000; i <= limit; i++ {
		key := strconv.Itoa(i)
		assert.This(fb.Search(key)).Is(i)
	}
}

func ExampleBuilder() {
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	store := stor.HeapStor(8192)
	bldr := Builder(store)
	bldr.Add("1000xxxx", 1000)
	bldr.Add("1001xxxx", 1001)
	bldr.Add("1002xxxx", 1002)
	bldr.Add("1003xxxx", 1003)
	fb := bldr.Finish()
	fb.print()
	// The important thing here is that the second known (1001)
	// is NOT "1" which would mean searches for 1000 would fail
	// and NOT "1001xxxx" which is longer than necessary.

	// Output:
	// <<<------------------------------
	// offset 0  LEAF
	// '' 1001 1002 1003
	// ------------------------------>>>
}

func Examplefbtree_MergeAndSave() {
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	store := stor.HeapStor(8192)
	x := &ixbuf.T{}
	x.Insert("1000xxxx", 1000)
	x.Insert("1001xxxx", 1001)
	x.Insert("1002xxxx", 1002)
	x.Insert("1003xxxx", 1003)
	fb := CreateFbtree(store, nil)
	fb = fb.MergeAndSave(x.Iter(false))
	fb.print()
	// The important thing here is that the second known (1001)
	// is NOT "1" which would mean searches for 1000 would fail
	// and NOT "1001xxxx" which is longer than necessary.

	// Output:
	// <<<------------------------------
	// offset 4  LEAF
	// '' 1001 1002 1003
	// ------------------------------>>>
}
