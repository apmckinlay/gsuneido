// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestIteratorEmpty(*testing.T) {
	bt := CreateBtree(stor.HeapStor(256), nil)
	it := bt.Iterator()
	it.Next()
	assert.That(it.Eof())
	it.Next()
	assert.That(it.Eof())
	it.Rewind()
	it.Next()
	assert.That(it.Eof())

	it.Rewind()
	it.Prev()
	assert.That(it.Eof())
	it.Prev()
	assert.That(it.Eof())
	it.Rewind()
	it.Prev()
	assert.That(it.Eof())
}

func TestIterator(*testing.T) {
	const n = 1000
	var data [n]string
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string { return data[i-1] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	randKey := str.UniqueRandomOf(4, 6, "abcde")
	for i := range n {
		data[i] = randKey()
	}
	sort.Strings(data[:])
	bldr := Builder(stor.HeapStor(8192))
	for i, k := range data {
		assert.That(bldr.Add(k, uint64(i+1))) // +1 to avoid zero
	}
	bt := bldr.Finish()

	it := bt.Iterator()
	test := func(i int) {
		assert.Msg("eof ", i).That(!it.Eof())
		assert.This(it.curOff - 1).Is(i)
		assert.This(it.curKey).Is(data[i])
	}

	// test Iterator Next
	for i := range n {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())

	// test Iterator Prev
	it = bt.Iterator()
	for i := n - 1; i >= 0; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	// test Seek between keys
	for i, k := range data {
		k += "0" // increment to nonexistent
		it.Seek(k)
		if i+1 < len(data) {
			test(i + 1)
		} else {
			test(n - 1)
		}
	}

	// test Seek & Next
	for i, k := range data {
		it.Seek(k)
		test(i)
		it.Next()
		if i+1 < len(data) {
			test(i + 1)
		} else {
			assert.That(it.Eof())
		}
	}

	// test Seek & Prev
	for i, k := range data {
		it.Seek(k)
		test(i)
		it.Prev()
		if i-1 >= 0 {
			test(i - 1)
		} else {
			assert.That(it.Eof())
		}
	}

	it.Seek("") // before first
	test(0)

	it.Seek("~") // after last
	test(n - 1)

	org := n / 4
	it.Range(Range{Org: data[org], End: ixkey.Max})
	for i := org; i < n; i++ {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())

	end := n / 2
	it.Range(Range{End: data[end]})
	for i := range end {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())

	it.Range(Range{Org: data[org], End: data[end]})
	for i := org; i < end; i++ {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())
	it.Seek(data[0])
	assert.That(it.Eof())
	it.Seek(data[end])
	assert.That(it.Eof())

	it.Range(Range{Org: data[org], End: ixkey.Max})
	for i := n - 1; i >= org; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	it.Range(Range{End: data[end]})
	for i := end - 1; i >= 0; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	it.Range(Range{Org: data[org], End: data[end]})
	for i := end - 1; i >= org; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	it.Range(Range{Org: data[123], End: data[123] + "\x00"})
	it.Next()
	test(123)
	it.Next()
	assert.That(it.Eof())
}

func TestToChunk(t *testing.T) {
	assert := assert.T(t).This
	data := []string{"ant", "cat", "dog"}
	b := nodeBuilder{}
	for i, k := range data {
		b.Add(k, uint64(i+1), 1) // +1 to avoid zero
	}
	nd := b.Entries()
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string { return data[i-1] }

	bt := &btree{}
	it := nd.iter()
	ci := it.toChunk(bt, true).(*chunkIter) // leaf
	ci.next()
	assert(ci.i).Is(0)
	assert(ci.c).Is(chunk{{key: "ant", off: 1}, {key: "cat", off: 2},
		{key: "dog", off: 3}})
	ci = it.toChunk(bt, false).(*chunkIter) // tree
	assert(ci.c).Is(chunk{{key: "", off: 1}, {key: "c", off: 2},
		{key: "d", off: 3}})

	it.next()
	it.next()
	assert(it.offset).Is(2)
	ci = it.toChunk(bt, false).(*chunkIter)
	assert(ci.off()).Is(2)
}

//-------------------------------------------------------------------

// func (it *Iterator) printStack() {
// 	for i := 0; i <= it.bt.treeLevels; i++ {
// 		it.printLevel(i)
// 	}
// }

// func (it *Iterator) printLevel(i int) {
// 	if ni, ok := it.stack[i].(*nodeIter); ok {
// 		fmt.Print(i, " | ")
// 		ni2 := ni.node.iter()
// 		for ni2.next() {
// 			if ni2.pos == ni.pos {
// 				fmt.Print("(")
// 			}
// 			if i == 0 {
// 				fmt.Print(it.bt.getLeafKey(ni2.offset), " ")
// 			} else if len(ni2.known) == 0 {
// 				fmt.Print("'' ")
// 			} else {
// 				fmt.Print(string(ni2.known), " ")
// 			}
// 		}
// 		if ni.pos >= ni2.pos {
// 			fmt.Print("(")
// 		}
// 	} else {
// 		fmt.Print(i, " + ")
// 		ci := it.stack[i].(*chunkIter)
// 		for j, s := range ci.c {
// 			if j == ci.i {
// 				fmt.Print("(")
// 			}
// 			fmt.Print(s.key, " ")
// 		}
// 		if ci.i >= len(ci.c) {
// 			fmt.Print("(")
// 		}
// 	}
// 	fmt.Println()
// }
