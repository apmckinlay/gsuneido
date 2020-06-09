// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMbtree(t *testing.T) {
	x := newMbtree()
	mCompare(t, x, mLeafSlots{})
	data := mLeafSlots{
		{"hello", 456},
		{"andrew", 123},
		{"zorro", 789},
	}
	for _, v := range data {
		x.Insert(v.key, v.off)
	}
	mCompare(t, x, data)
}

func mCompare(t *testing.T, x *mbtree, data mLeafSlots) {
	sort.Sort(data)
	iter := x.Iterator()
	for _, v := range data {
		key, off, ok := iter()
		Assert(t).That(ok, Equals(true))
		Assert(t).That(key, Equals(v.key))
		Assert(t).That(off, Equals(v.off).Comment(key))
	}
	_, _, ok := iter()
	Assert(t).That(ok, Equals(false))
}

type mLeafSlots []mLeafSlot

func (a mLeafSlots) Len() int      { return len(a) }
func (a mLeafSlots) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a mLeafSlots) Less(i, j int) bool {
	return a[i].key < a[j].key ||
		(a[i].key == a[j].key && a[i].off < a[j].off)
}

func TestMbtreeRandom(t *testing.T) {
	var nGenerate = 8
	var nShuffle = 8
	if testing.Short() {
		nGenerate = 2
		nShuffle = 2
	}
	const n = mSize * 80
	for gi := 0; gi < nGenerate; gi++ {
		data := make(mLeafSlots, n)
		randKey := str.UniqueRandom(3, 10)
		for i := uint64(0); i < n; i++ {
			data[i] = mLeafSlot{randKey(), i}
		}
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			x := newMbtree()
			for _, v := range data {
				x.Insert(v.key, v.off)
			}
			mCompare(t, x, data)
		}
	}
}

func TestMbtreeUnevenSplit(t *testing.T) {
	const n = mSize * 87 // won't fit without uneven splits
	data := make(mLeafSlots, n)
	randKey := str.UniqueRandom(3, 10)
	for i := uint64(0); i < n; i++ {
		data[i] = mLeafSlot{randKey(), i}
	}
	sort.Sort(data)
	x := newMbtree()
	for _, v := range data {
		x.Insert(v.key, v.off)
	}
	mCompare(t, x, data)
	x = newMbtree()
	for i := len(data) - 1; i >= 0; i-- {
		x.Insert(data[i].key, data[i].off)
	}
	mCompare(t, x, data)
}

//-------------------------------------------------------------------

func (m *mbtree) print() int {
	var n int
	if m.tree != nil {
		n = m.tree.print()
		fmt.Println("total size", n,
			"average leaf occupancy", float32(n)/float32(m.tree.size)/float32(mSize))
	} else {
		fmt.Println("no tree, single leaf")
		n = m.leaf.print()
	}
	return n
}

func (tree *mTree) print() int {
	n := 0
	for i := 0; i < tree.size; i++ {
		fmt.Println(i, tree.slots[i].key, "leaf size", tree.slots[i].leaf.size)
		n += tree.slots[i].leaf.print()
	}
	return n
}

func (leaf *mLeaf) print() int {
	for i := 0; i < leaf.size; i++ {
		// fmt.Println("   ", leaf.slots[i])
	}
	return leaf.size
}
