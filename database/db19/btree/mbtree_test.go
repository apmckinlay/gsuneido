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

func TestMbtreeRandom(t *testing.T) {
	var nGenerate = 5
	var nShuffle = 5
	if testing.Short() {
		nGenerate = 2
		nShuffle = 2
	}
	const n = mSize * 80
	data := make([]string, n)
	for g := 0; g < nGenerate; g++ {
		randKey := str.UniqueRandom(3, 10)
		for i := 0; i < n; i++ {
			data[i] = randKey()
		}
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			x := newMbtree(0)
			for i, k := range data {
				x.Insert(k, uint64(i))
			}
			x.checkData(t, data)
		}
	}
}

func TestMbtreeUnevenSplit(t *testing.T) {
	const n = mSize * 87 // won't fit without uneven splits
	data := make([]string, n)
	randKey := str.UniqueRandom(3, 10)
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	sort.Strings(data)
	m := newMbtree(0)
	for i, k := range data {
		m.Insert(k, uint64(i))
	}
	m.checkData(t, data)
	m = newMbtree(0)
	for i := len(data) - 1; i >= 0; i-- {
		m.Insert(data[i], uint64(i))
	}
	m.checkData(t, data)
}

//-------------------------------------------------------------------

func (m *mbtree) check() {
	prev := ""
	m.ForEach(func(key string, off uint64) {
		if key <= prev {
			panic("keys out of order " + prev + " " + key)
		}
		prev = key
	})
}

func (m *mbtree) checkData(t *testing.T, data []string) {
	m.check()
	for i, key := range data {
		Assert(t).That(m.Search(key), Is(i))
		Assert(t).That(m.Search(bigger(key)), Is(0))
		Assert(t).That(m.Search(smaller(key)), Is(0))
	}
}

func bigger(s string) string {
	return s + " "
}

func smaller(s string) string {
	return s[:len(s)-1] + "~"
}

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

func (tree *mbTreeNode) print() int {
	n := 0
	for i := 0; i < tree.size; i++ {
		fmt.Println(i, tree.slots[i].key, "leaf size", tree.slots[i].leaf.size)
		n += tree.slots[i].leaf.print()
	}
	return n
}

func (leaf *mbLeaf) print() int {
	for i := 0; i < leaf.size; i++ {
		// fmt.Println("   ", leaf.slots[i])
	}
	return leaf.size
}
