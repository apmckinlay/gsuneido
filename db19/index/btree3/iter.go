// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type Iterator struct {
	bt    *btree
	rng   Range
	tree  [maxLevels]*treeIter // tree[0] is root
	leaf  *leafIter
	state iterState
}

const maxLevels = 8

type Range = iterator.Range

type iterState byte

const (
	rewound iterState = iota
	within
	eof
)

func (bt *btree) Iterator() *Iterator {
	return &Iterator{bt: bt, state: rewound, rng: iterator.All}
}

func (it *Iterator) Key(buf []byte) []byte {
	return it.leaf.key(buf)
}

func (it *Iterator) Off() uint64 {
	return it.leaf.offset()
}

func (it *Iterator) Next() bool {
	switch it.state {
	case rewound:
		it.SeekAll(it.rng.Org)
	case within:
		it.next()
	case eof:
		return false
	}
	return it.state != eof
}

func (it *Iterator) next() {
	for {
		if it.leaf.next() {
			return
		} else if !it.nextLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) nextLeaf() bool {
	bt := it.bt
	i := bt.treeLevels - 1 // closest to leaf
	var nodeOff uint64
	// go up the tree until we can advance
	for ; i >= 0; i-- {
		if it.tree[i].next() {
			nodeOff = it.tree[i].off()
			break
		} // else end of tree node, keep going up
	}
	if i < 0 {
		return false // end of root = eof
	}
	// then descend back down
	for i++; i < bt.treeLevels; i++ {
		it.tree[i] = bt.getTree(nodeOff).iter()
		assert.That(it.tree[i].next())
		nodeOff = it.tree[i].off()
	}
	it.leaf = bt.getLeaf(nodeOff).iter()
	return true
}

func (it *Iterator) SeekAll(string) {
	off := it.bt.root
	for i := range it.bt.treeLevels {
		it.tree[i] = it.bt.getTree(off).iter()
		it.tree[i].next()
		off = it.tree[i].off()
	}
	it.leaf = it.bt.getLeaf(off).iter()
	assert.That(it.leaf.next())
	it.state = within
}
