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
	var buf []byte
	const from = 100
	for to := 101; to < 199; to++ {
		st := stor.HeapStor(64 * 1024)
		b := Builder(st)
		b.splitSize = 40
		for i := from; i < to; i++ {
			assert.That(b.Add(strconv.Itoa(i), uint64(i)))
		}
		bt := b.Finish()
		// fmt.Println(to, "================================")
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
			assert.That(it.Next())
			buf = it.Key(buf)
			assert.That(string(buf) == strconv.Itoa(i))
			assert.This(it.Off()).Is(uint64(i))
		}
		assert.That(!it.Next())

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
	var buf []byte
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
		assert.That(it.Next())
		buf = it.Key(buf)
		assert.That(string(buf) == strconv.Itoa(i))
		assert.This(it.Off()).Is(uint64(i))
	}
	assert.That(!it.Next())

	// Lookup
	for i := from; i < to; i++ {
		assert.This(bt.Lookup(strconv.Itoa(i))).Is(uint64(i))
	}
}

// ------------------------------------------------------------------

func (bt *btree) print() {
	bt.print1(0, bt.root)
}

func (bt *btree) print1(depth int, offset uint64) {
	if depth < bt.treeLevels {
		nd := readTree(bt.stor, offset)
		fmt.Println(offset, "->", nd)
		for i := 0; i < nd.nkeys(); i++ {
			bt.print1(depth+1, nd.offset(i)) // RECURSE
			fmt.Println(strings.Repeat(" .", depth), string(nd.key(i)))
		}
		bt.print1(depth+1, nd.offset(nd.nkeys())) // RECURSE
	} else {
		nd := readLeaf(bt.stor, offset)
		fmt.Println(strings.Repeat(" .", depth), offset, "->", nd)
	}
}
