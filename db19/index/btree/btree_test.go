// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"math"
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
	prefix := strings.Repeat("helloworld", 30)
	keyfn := func(i int) string {
		return prefix + strconv.Itoa(i)
	}
	assert := assert.T(t)
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return keyfn(int(i))
	}
	bldr := Builder(stor.HeapStor(8192))
	start := 100000
	limit := 999999
	if testing.Short() {
		start = 10000
		limit = 99999
	}
	for i := start; i <= limit; i++ {
		key := keyfn(i)
		assert.That(bldr.Add(key, uint64(i)))
	}
	bt := bldr.Finish()
	bt.Check(nil)

	// iterate
	iter := bt.Iterator()
	for i := start; i <= limit; i++ {
		key := keyfn(i)
		iter.Next()
		k, o := iter.Cur()
		assert.True(strings.HasPrefix(key, k))
		assert.This(o).Is(i)
	}
	iter.Next()
	assert.That(iter.Eof())

	// Lookup
	for i := start; i <= limit; i++ {
		key := keyfn(i)
		assert.This(bt.Lookup(key)).Is(i)
		assert.This(bt.Lookup(key + "0")).Is(0) // nonexistent
	}
}

func TestBuilder(t *testing.T) {
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	bldr := Builder(stor.HeapStor(8192))
	assert.That(bldr.Add("1000xxxx", 1000))
	assert.That(bldr.Add("1001xxxx", 1001))
	assert.That(bldr.Add("1002xxxx", 1002))
	assert.That(bldr.Add("1003xxxx", 1003))
	bt := bldr.Finish()
	nd := bt.getNode(bt.root)
	assert.T(t).This(nodeKnowns(nd)).Is([]string{"", "1001", "1002", "1003"})
}

func nodeKnowns(nd node) []string {
	keys := []string{}
	for it := nd.iter(); it.next(); {
		keys = append(keys, string(it.known))
	}
	return keys
}

func TestBtreeMergeAndSave(t *testing.T) {
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
	nd := bt.getNode(bt.root)
	assert.T(t).This(nodeKnowns(nd)).Is([]string{"", "1001", "1002", "1003"})
}

func TestBtreeFracPos(t *testing.T) {
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 192
	Fanout = 20
	var bt *btree
	key := func(i int) string {
		return fmt.Sprintf("%05d", i)
	}
	makeBtree := func(n int) {
		// for consistent results we need the root to be quite full
		// since Builder splits unevenly due to building in order
		stor := stor.HeapStor(8192)
		stor.Alloc(8)
		b := Builder(stor)
		for i := 0; i < n; i++ {
			assert.That(b.Add(key(i), 1))
		}
		//		assert.That(len(b.levels[len(b.levels)-1].nb.node) > 190)
		bt = b.Finish()
	}
	test := func(key string, expected float32) {
		t.Helper()
		fracPos := bt.fracPos(key)
		diff := fracPos - expected
		// fmt.Println(key, "expected", expected, "got", fracPos, "diff", diff)
		if math.Abs(float64(diff)) > .025 {
			t.Error("\nkey", fmt.Sprintf("%q", key),
				"got", fracPos, "expected", expected, "difference", diff)
		}
	}
	n := 24 // full single root node
	makeBtree(n)
	assert.Msg("tree levels").This(bt.treeLevels).Is(0)
	test(ixkey.Min, 0)
	for i := 0; i < n; i++ {
		test(key(i), float32(i)/float32(n))
	}
	test(ixkey.Max, 1)

	n = 9900
	makeBtree(n)
	assert.Msg("tree levels").This(bt.treeLevels).Is(3)
	assert.This(bt.getNode(bt.root).Size()).Is(2)
	test(ixkey.Min, 0)
	for i := 0; i < n; i += 100 {
		exp := float32(i) / float32(n)
		test(key(i), exp)
	}
	test(ixkey.Max, 1)
}
