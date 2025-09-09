// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"sort"
	"strconv"
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

func TestToUnodeIter(t *testing.T) {
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
	ui := it.toUnodeIter(bt).(*unodeIter) // tree
	assert(ui.u).Is(unode{{key: "", off: 1}, {key: "c", off: 2},
		{key: "d", off: 3}})

	it.next()
	it.next()
	assert(it.offset).Is(2)
	ui = it.toUnodeIter(bt).(*unodeIter)
	assert(ui.off()).Is(2)
	
	it = nd.iter()
	ui = it.toUnodeIter(bt).(*unodeIter)
	assert(ui.prev()).Is(true)
	assert(ui.off()).Is(3)
}

// SeekAll

// buildTree builds a simple one-level btree with keys "1".."n"
// and offsets equal to the integer value.
func buildTree(n int) *btree {
	b := Builder(stor.HeapStor(8192))
	for i := 1; i <= n; i++ {
		k := strconv.Itoa(i)
		assert.That(b.Add(k, uint64(i)))
	}
	bt := b.Finish()
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, off uint64) string {
		return strconv.Itoa(int(off))
	}
	return bt
}

// SeekAll should position on the last item when searching ixkey.Max.
// This is essential to allow Prev() from rewind to begin from the end.
func TestSeekAllMaxPositionsAtLast(t *testing.T) {
	bt := buildTree(9)
	it := bt.Iterator()

	// rootU should be initialized for non-empty trees built via Builder.Finish
	assert.That(len(bt.rootUnode) > 0)

	it.SeekAll(ixkey.Max)
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("9")
	assert.This(off).Is(uint64(9))
}

// SeekAll on an exact key should position on that key (not the next).
func TestSeekAllExactKey(t *testing.T) {
	bt := buildTree(9)
	it := bt.Iterator()

	it.SeekAll("5")
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("5")
	assert.This(off).Is(uint64(5))
}

// SeekAll on a key between existing keys should land on the next greater key.
// e.g., between "5" and "6" should position to "6".
func TestSeekAllBetweenKeysGoesToNextGreater(t *testing.T) {
	bt := buildTree(9)
	it := bt.Iterator()

	it.SeekAll("5~") // between 5 and 6 (since "~" > "5" and < "6")
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("6")
	assert.This(off).Is(uint64(6))
}

// SeekAll on a key larger than the maximum (but not ixkey.Max) should remain at last.
// This aligns with SeekAll semantics to avoid setting EOF so Prev() can step back.
func TestSeekAllAboveMaxStaysAtLast(t *testing.T) {
	bt := buildTree(9)
	it := bt.Iterator()

	it.SeekAll("9zzzz") // greater than the last key
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("9")
	assert.This(off).Is(uint64(9))
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
// 		ci := it.stack[i].(*unodeIter)
// 		for j, s := range ci.u {
// 			if j == ci.i {
// 				fmt.Print("(")
// 			}
// 			fmt.Print(s.key, " ")
// 		}
// 		if ci.i >= len(ci.u) {
// 			fmt.Print("(")
// 		}
// 	}
// 	fmt.Println()
// }
