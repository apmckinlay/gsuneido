// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBuilderErrors(t *testing.T) {
	b := Builder(stor.HeapStor(256))
	assert.That(b.Add("", 1))
	assert.That(!b.Add("", 2))
	assert.That(b.Add("x", 1))
	assert.That(!b.Add("x", 2))
	assert.This(func() { b.Add("a", 2) }).Panics("out of order")
}

func TestBuilderSmall(t *testing.T) {
	bt := Builder(stor.HeapStor(256)).Finish()
	bt.Check(nil)

	const from = 100
	for to := 101; to < 199; to++ {
		st := stor.HeapStor(64 * 1024)
		b := Builder(st)
		b.shouldSplit = func(nd node) bool {
			return nd.noffs() >= 4
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
			assert.This(it.Key()).Is(strconv.Itoa(i))
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
		assert.This(it.Key()).Is(strconv.Itoa(i))
		assert.This(it.Offset()).Is(uint64(i))
	}
	it.Next()
	assert.That(it.Eof())

	// Lookup
	for i := from; i < to; i++ {
		assert.This(bt.Lookup(strconv.Itoa(i))).Is(uint64(i))
	}
}
