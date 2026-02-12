// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/str"
)

// builder is used to bulk load a btree.
// Keys must be added in order.
// The btree is built bottom up with no splitting or inserting.
// All nodes will be "full" except for the right hand edge.
// The builder holds the right hand edge of the btree.
type builder struct {
	stor     *stor.Stor
	leaf     leafBuilder
	tree     []*treeBuilder // root is last (since tree grows up)
	prev     string
	havePrev bool
	count    int
}

func Builder(st *stor.Stor) *builder {
	return &builder{stor: st}
}

// Add returns false for duplicate keys and panics for out of order
func (b *builder) Add(key string, off uint64) bool {
	if b.havePrev {
		if key == b.prev {
			return false // duplicate
		}
		if key < b.prev {
			panic("btree builder: keys out of order")
		}
	}
	b.addLeaf(key, off)
	b.prev = key
	b.havePrev = true
	b.count++
	return true
}

func (b *builder) addLeaf(key string, off uint64) {
	if !b.leaf.tryAdd(key, off) {
		off2 := b.leaf.finishTo(b.stor)
		sep := b.sep(b.prev, key)
		b.addTree(0, off2, sep)
		b.leaf.reset()       // reuse
		b.leaf.add(key, off) // Will always succeed on empty builder
	}
}

func (b *builder) addTree(ti int, off uint64, sep string) {
	if ti >= len(b.tree) {
		b.tree = append(b.tree, &treeBuilder{}) // new root
	}
	tree := b.tree[ti]
	newSize := tree.size() + len(sep) + 7
	if tree.noffs() >= splitCount || newSize > maxNodeSize {
		off2 := tree.finishTo(b.stor, off)
		b.addTree(ti+1, off2, sep) // RECURSE
		tree.reset()               // reuse the memory
	} else {
		tree.add(off, sep)
	}
}

func (b *builder) sep(prev, key string) string {
	cp := str.CommonPrefixLen(prev, key)
	return key[:cp+1]
}

func (b *builder) Finish() iface.Btree {
	off := b.leaf.finishTo(b.stor)

	for i := range b.tree {
		off = b.tree[i].finishTo(b.stor, off)
	}
	return &btree{stor: b.stor, root: off, treeLevels: len(b.tree), count: b.count}
}
