// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBtreeFracPos(t *testing.T) {
	var n int
	var bt *btree
	key := func(i int) string {
		return fmt.Sprintf("%05d", i)
	}
	makeBtree := func(m int) {
		n = m
		b := Builder(heapstor(8192))
		for i := range n {
			assert.That(b.Add(key(i), 1))
		}
		bt = b.Finish().(*btree)
	}
	test := func(key string, expected float64) {
		t.Helper()
		fracPos := bt.fracPos(key)
		diff := fracPos - expected
		// fmt.Println(key, "expected", expected, "got", fracPos, "diff", diff)
		if math.Abs(diff) > .025 {
			t.Error("\nkey", key,
				"got", fracPos, "expected", expected, "difference", diff)
		}
	}
	
	// single root (leaf) node
	makeBtree(75)
	assert.Msg("tree levels").This(bt.treeLevels).Is(0)
	test(ixkey.Min, 0)
	for i := range n {
		test(key(i), float64(i)/float64(n))
	}
	test(ixkey.Max, 1)

	// small root, treeLevels 1 
	makeBtree(120)
	assert.Msg("tree levels").This(bt.treeLevels).Is(1)
	assert.This(bt.readTree(bt.root).noffs()).Is(2)
	test(ixkey.Min, 0)
	for i := range n {
		test(key(i), float64(i)/float64(n))
	}
	test(ixkey.Max, 1)

	// large root, treeLevels 1
	makeBtree(1200)
	assert.Msg("tree levels").This(bt.treeLevels).Is(1)
	assert.This(bt.readTree(bt.root).noffs()).Is(12)
	test(ixkey.Min, 0)
	for i := range n {
		test(key(i), float64(i)/float64(n))
	}
	test(ixkey.Max, 1)

	// small root, treeLevels 2
	makeBtree(12000)
	assert.Msg("tree levels").This(bt.treeLevels).Is(2)
	assert.This(bt.readTree(bt.root).noffs()).Is(2)
	test(ixkey.Min, 0)
	for i := 1000; i < n; i += 100 {
		exp := float64(i) / float64(n)
		test(key(i), exp)
	}
	test(ixkey.Max, 1)
}
