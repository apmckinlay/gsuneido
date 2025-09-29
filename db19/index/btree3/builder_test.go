// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBuilderSmall(t *testing.T) {
	bt := Builder(stor.HeapStor(256)).Finish()
	bt.Check(nil)
	
	const from = 100
	for to := 101; to < 199; to++ {
		st := stor.HeapStor(64 * 1024)
		b := Builder(st)
		b.shouldSplit = func(nd splitable) bool {
			return nd.nkeys() >= 4
		}
		for i := from; i < to; i++ {
			assert.That(b.Add(strconv.Itoa(i), uint64(i)))
		}
		bt := b.Finish()
		// bt.print()

		// Check
		i := from
		count, _, _ := bt.Check(func(off uint64) {
			assert.This(off).Is(i)
			i++
		})
		assert.This(count).Is(to - from)
		assert.This(i).Is(to)

		// Iterator
		it := bt.Iterator()
		for i := from; i < to; i++ {
			it.Next()
			assert.That(!it.Eof())
			assert.That(string(it.Key()) == strconv.Itoa(i))
			assert.This(it.Offset()).Is(uint64(i))
		}
		it.Next()
		assert.That(it.Eof())

		// Lookup
		for i := from; i < to; i++ {
			assert.This(bt.Lookup(strconv.Itoa(i))).Is(uint64(i))
		}
	}
}

func TestBuilderBig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping big test")
	}
	const from = 1_000_000
	const to = 10_000_000
	st := stor.HeapStor(64 * 1024)
	b := Builder(st)
	for i := from; i < to; i++ {
		assert.That(b.Add(strconv.Itoa(i), uint64(i)))
	}
	bt := b.Finish()

	// Check
	i := from
	count, _, _ := bt.Check(func(off uint64) {
		assert.This(off).Is(i)
		i++
	})
	assert.This(count).Is(to - from)
	assert.This(i).Is(to)

	// Iterator
	it := bt.Iterator()
	for i := from; i < to; i++ {
		it.Next()
		assert.That(!it.Eof())
	assert.That(string(it.Key()) == strconv.Itoa(i))
		assert.This(it.Offset()).Is(uint64(i))
	}
	it.Next()
	assert.That(it.Eof())

	// Lookup
	for i := from; i < to; i++ {
		assert.This(bt.Lookup(strconv.Itoa(i))).Is(uint64(i))
	}
}

// ------------------------------------------------------------------

func (bt *btree) print() {
	fmt.Println("-----------------------------")
	bt.print1(0, bt.root)
}

func (bt *btree) print1(depth int, offset uint64) {
	indent := strings.Repeat(" .", depth)
	if depth < bt.treeLevels {
		nd := readTree(bt.stor, offset)
		fmt.Println(indent, offset, "->", nd)
		for i := 0; i < nd.nkeys(); i++ {
			bt.print1(depth+1, nd.offset(i)) // RECURSE
			fmt.Println(indent, "<" + string(nd.key(i)) + ">")
		}
		bt.print1(depth+1, nd.offset(nd.nkeys())) // RECURSE
	} else {
		nd := readLeaf(bt.stor, offset)
		fmt.Println(indent, offset, "->", nd)
	}
}
