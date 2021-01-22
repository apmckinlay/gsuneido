// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

func TestFbtreeIO(t *testing.T) {
	store := stor.HeapStor(128)
	before := fnode([]byte("helloworld"))
	off := before.putNode(store)
	after := readNode(store, off)
	cksum.MustCheck(after[:len(after)+cksum.Len])
	assert.T(t).This(string(after)).Is(string(before))
}

func TestFbtreeBuilder(t *testing.T) {
	assert := assert.T(t)
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	store := stor.HeapStor(8192)
	bldr := Builder(store)
	start := 100000
	limit := 999999
	if testing.Short() {
		start = 10000
		limit = 99999
	}
	for i := start; i <= limit; i++ {
		key := strconv.Itoa(i)
		bldr.Add(key, uint64(i))
	}
	fb := bldr.Finish()
	fb.Check(nil)

	// iterate
	iter := fb.Iterator()
	for i := start; i <= limit; i++ {
		key := strconv.Itoa(i)
		iter.Next()
		k, o := iter.Cur()
		assert.True(strings.HasPrefix(key, k))
		assert.This(o).Is(i)
	}
	iter.Next()
	assert.That(iter.Eof())

	// Lookup
	for i := start; i <= limit; i++ {
		key := strconv.Itoa(i)
		assert.This(fb.Lookup(key)).Is(i)
		assert.This(fb.Lookup(key + "0")).Is(0) // nonexistent
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
	fb = fb.MergeAndSave(x.Iter())
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
