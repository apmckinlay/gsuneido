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

func TestBuilderErrors(t *testing.T) {
	b := Builder(stor.HeapStor(256))
	assert.That(b.Add("", 1))
	assert.That(!b.Add("", 2))
	assert.That(b.Add("x", 1))
	assert.That(!b.Add("x", 2))
	assert.This(func() { b.Add("a", 2) }).Panics("out of order")
}

func TestBuilderSmall(t *testing.T) {
	defer SetSplit(SetSplit(4))
	bt := Builder(stor.HeapStor(256)).Finish()
	bt.Check(nil)

	const from = 100
	for to := 101; to < 199; to++ {
		st := stor.HeapStor(64 * 1024)
		b := Builder(st)
		for i := from; i < to; i++ {
			assert.That(b.Add(strconv.Itoa(i), uint64(i)))
		}
		bt := b.Finish()
		// bt.(*btree).Print()

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

func TestBuilderLargeKeys1(t *testing.T) {
	st := stor.HeapStor(64 * 1024)
	large := strings.Repeat("a", 1500)
	b := Builder(st)
	for i := range 99 {
		b.Add(fmt.Sprintf("%02d", i)+large, uint64(i))
	}
	bt := b.Finish()
	bt.Check(nil)
}

func TestBuilderLargeKeys2(t *testing.T) {
	st := stor.HeapStor(64 * 1024)
	large := strings.Repeat("a", 1500)
	b := Builder(st)
	for i := range 99 {
		b.Add(large+fmt.Sprintf("%02d", i), uint64(i))
	}
	bt := b.Finish()
	bt.Check(nil)
}

func TestBuilderLargeKeys3(t *testing.T) {
	st := stor.HeapStor(64 * 1024)
	large := strings.Repeat("a", 5000)
	b := Builder(st)
	for i := 0; i < 99; i += 3 {
		b.Add(fmt.Sprintf("%02d", i), uint64(i))
		b.Add(fmt.Sprintf("%02d", i+1)+large, uint64(i+1))
		b.Add(fmt.Sprintf("%02d", i+2), uint64(i+2))
	}
	bt := b.Finish()
	bt.Check(nil)
}
