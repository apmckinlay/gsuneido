// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
db19/index/btree3 is a new, optimized btree implementation.
It is not used yet.

Two things create or modify btrees:
- bulk in-order loading from load or compact (Builder)
- add, update, and delete from an ixbuf (MergeAndSave)

A btree consists of two kinds of nodes: leafNode and treeNode.
There are treeLevels of treeNodes.
The root is a leafNode if treeLevels == 0, otherwise it is a treeNode.

leafNode and treeNode have the same basic representations.
leafNode has prefix compression.
treeNode has an extra offset since the fields are separators.
*/
package btree

import (
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// MaxNodeSize is the maximum node size in bytes, split if larger.
// var rather than const because it is overridden by tests.
const minSplit = 1024 // ???
const maxSplit = 8192 // ???

// EntrySize is the estimated average entry size
const EntrySize = 10

// TreeHeight is the estimated average tree height.
// It is used by Table.lookupCost
const TreeHeight = 3

// Fanout is the estimated average number of children per node.
const fanout = 100

type btree struct {
	stor        *stor.Stor
	root        uint64
	treeLevels  int
	shouldSplit func(splitable) bool
}

// Lookup returns the offset for a key, or 0 if not found.
func (bt *btree) Lookup(key string) uint64 {
	off := bt.root
	for range bt.treeLevels {
		nd := bt.getTree(off)
		_, off = nd.search(key)
	}
	nd := bt.getLeaf(off)
	i, found := nd.search(key)
	if !found {
		return 0 // not found
	}
	return nd.offset(i)
}

func (bt *btree) getTree(off uint64) treeNode {
	return readTree(bt.stor, off)
}

func (bt *btree) getLeaf(off uint64) leafNode {
	return readLeaf(bt.stor, off)
}

// Check verifies that the keys are in order and returns the number of keys.
// If the supplied function is not nil, it is applied to each leaf offset.
func (bt *btree) Check(fn func(uint64)) (count, size, nnodes int) {
	var prev []byte // updated by leaf
	var check1 func(int, uint64)
	check1 = func(depth int, offset uint64) {
		nnodes++
		if depth < bt.treeLevels {
			// tree
			nd := readTree(bt.stor, offset)
			if depth == 0 {
				assert.That(nd.nkeys() >= 1)
			} else {
				assert.That(nd.noffs() >= 1)
			}
			size = len(nd)
			for i := 0; i < nd.nkeys(); i++ {
				sep := nd.key(i)
				assert.That(string(sep) > string(prev))
				check1(depth+1, nd.offset(i)) // RECURSE
			}
			check1(depth+1, nd.offset(nd.nkeys())) // RECURSE
		} else {
			// leaf
			nd := readLeaf(bt.stor, offset)
			if nd.nkeys() == 0 {
				assert.That(bt.treeLevels == 0)
				return
			}
			k := nd.key(0)
			assert.That(k > string(prev))
			size = len(nd)
			var prevSuffix []byte
			first := true
			for it := nd.iter(); it.next(); count++ {
				k := it.suffix()
				if !first {
					assert.That(string(k) > string(prevSuffix))
				}
				first = false
				prevSuffix = k
				if fn != nil {
					fn(it.offset())
				}
			}
			prev = append(append((prev)[:0], nd.prefix()...),
				nd.suffix(nd.nkeys()-1)...)
		}
	}
	check1(0, bt.root)
	return
}
