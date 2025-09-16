// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/str"
)

// builder is used to bulk load an btree.
// Keys must be added in order.
// The btree is built bottom up with no splitting or inserting.
// All nodes will be "full" except for the right hand edge.
// The builder holds the right hand edge of the btree.
type builder struct {
	stor      *stor.Stor
	leaf      leafBuilder
	tree      []*treeBuilder // root is last (since tree grows up)
	splitSize int            // for testing
	cur       string
	prev      string
}

func Builder(st *stor.Stor) *builder {
	return &builder{stor: st, splitSize: 1024} // ???
}

// Add returns false for duplicate keys and panics for out of order
func (b *builder) Add(key string, off uint64) bool {
	if b.leaf.empty() && key == b.prev {
		return false // duplicate
	}
	b.cur = key
	b.addLeaf(key, off)
	b.prev = key
	return true
}

func (b *builder) addLeaf(key string, off uint64) {
	if b.leaf.size() >= b.splitSize {
		nd := b.leaf.finish()
		off2 := nd.write(b.stor)
		sep := b.sep(b.prev, key)
		b.addTree(0, off2, sep)
		b.leaf.reset() // reuse the memory
	}
	b.leaf.add(key, off)
}

func (b *builder) addTree(ti int, off uint64, sep string) {
	if ti >= len(b.tree) {
		b.tree = append(b.tree, &treeBuilder{}) // new root
	}
	tree := b.tree[ti]
	if tree.size() >= b.splitSize {
		nd := tree.finish(off)
		off2 := nd.write(b.stor)
		b.addTree(ti+1, off2, b.sep(b.prev, b.cur)) // RECURSE
		tree.reset()                                // reuse the memory
	} else {
		tree.add(off, sep)
	}
}

func (b *builder) sep(prev, key string) string {
	cp := str.CommonPrefixLen(prev, key)
	return key[:cp+1]
}

func (b *builder) Finish() *btree {
	nd := b.leaf.finish()
	off := nd.write(b.stor)
	for i := range b.tree {
		nd := b.tree[i].finish(off)
		off = nd.write(b.stor)
	}
	return &btree{stor: b.stor, root: off, treeLevels: len(b.tree)}
}
