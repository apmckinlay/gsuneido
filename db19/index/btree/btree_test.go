// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

func TestBtreeIO(t *testing.T) {
	st := stor.HeapStor(128)
	before := node([]byte("helloworld"))
	off := before.putNode(st)
	after := readNode(st, off)
	cksum.MustCheck(after[:len(after)+cksum.Len])
	assert.T(t).This(string(after)).Is(string(before))
}

func TestBtreeBuilder(t *testing.T) {
	assert := assert.T(t)
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	bldr := Builder(stor.HeapStor(8192))
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
	bt := bldr.Finish()
	bt.Check(nil)

	// iterate
	iter := bt.Iterator()
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
		assert.This(bt.Lookup(key)).Is(i)
		assert.This(bt.Lookup(key + "0")).Is(0) // nonexistent
	}
}

func ExampleBuilder() {
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	bldr := Builder(stor.HeapStor(8192))
	bldr.Add("1000xxxx", 1000)
	bldr.Add("1001xxxx", 1001)
	bldr.Add("1002xxxx", 1002)
	bldr.Add("1003xxxx", 1003)
	bt := bldr.Finish()
	bt.print()
	// The important thing here is that the second known (1001)
	// is NOT "1" which would mean searches for 1000 would fail
	// and NOT "1001xxxx" which is longer than necessary.

	// Output:
	// <<<------------------------------
	// offset 1  LEAF
	// '' 1001 1002 1003
	// ------------------------------>>>
}

func Examplebtree_MergeAndSave() {
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	x := &ixbuf.T{}
	x.Insert("1000xxxx", 1000)
	x.Insert("1001xxxx", 1001)
	x.Insert("1002xxxx", 1002)
	x.Insert("1003xxxx", 1003)
	bt := CreateBtree(stor.HeapStor(8192), nil)
	bt = bt.MergeAndSave(x.Iter())
	bt.print()
	// The important thing here is that the second known (1001)
	// is NOT "1" which would mean searches for 1000 would fail
	// and NOT "1001xxxx" which is longer than necessary.

	// Output:
	// <<<------------------------------
	// offset 4  LEAF
	// '' 1001 1002 1003
	// ------------------------------>>>
}

func TestBtreeFracPos(t *testing.T) {
	var bt *btree
	makeBtree := func(size int) {
		b := Builder(stor.HeapStor(8192))
		for i := 0; i < size; i++ {
			key := fmt.Sprintf("%04d", i)
			b.Add(key, 1)
		}
		bt = b.Finish()
	}
	test := func(key string, expected float32) {
		t.Helper()
		fracPos := bt.fracPos(key)
		diff := expected - fracPos
		if diff < 0 {
			diff = -diff
		}
		if diff > .02 {
			t.Error("\nkey", fmt.Sprintf("%q", key),
				"got", fracPos, "expected", expected, "difference", diff)
		}
	}
	makeBtree(10) // single root leaf node
	test(ixkey.Min, 0)
	test(" ", 0)
	test("0000", 0)
	test("0001", .1)
	test("0009", .9)
	test("~", .9)
	test(ixkey.Max, 1)

	makeBtree(10000)
	// var diff float32
	for i := 0; i < 10000; i += 250 {
		key := fmt.Sprintf("%04d", i)
		// f := bt.fracPos(key)
		exp := float32(i) / 10000
		test(key, exp)
		// fmt.Printf("%d %.3f %+.3f\n", i, f, f - exp)
		// diff += f - exp
	}
	// fmt.Println("total diff", diff, "avg", diff / 40)
}
