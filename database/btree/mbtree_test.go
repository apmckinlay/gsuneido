// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
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
		x.Insert(v.key, v.rec)
	}
	mCompare(t, x, data)
}

func mCompare(t *testing.T, x *mbtree, data mLeafSlots) {
	sort.Sort(data)
	iter := x.Iterator()
	for _, v := range data {
		key, rec, ok := iter()
		Assert(t).That(ok, Equals(true))
		Assert(t).That(key, Equals(v.key))
		Assert(t).That(rec, Equals(v.rec))
	}
	_, _, ok := iter()
	Assert(t).That(ok, Equals(false))
}

type mLeafSlots []mLeafSlot

func (a mLeafSlots) Len() int      { return len(a) }
func (a mLeafSlots) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a mLeafSlots) Less(i, j int) bool {
	return a[i].key < a[j].key ||
		(a[i].key == a[j].key && a[i].rec < a[j].rec)
}

func TestMbtreeRandom(t *testing.T) {
	const n = mSize * 87
	data := make(mLeafSlots, n)
	x := newMbtree()
	for i := uint64(0); i < n; i++ {
		key := str.Random(3, 10)
		x.Insert(key, i)
		data[i] = mLeafSlot{key, i}
	}
	mCompare(t, x, data)
}

func TestMbtreeOrdered(t *testing.T) {
	const n = mSize * 87
	data := make(mLeafSlots, n)
	x := newMbtree()
	for i := uint64(0); i < n; i++ {
		key := str.Random(2, 10)
		data[i] = mLeafSlot{key, i}
	}
	sort.Sort(data)
	for _, v := range data {
		x.Insert(v.key, v.rec)
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
