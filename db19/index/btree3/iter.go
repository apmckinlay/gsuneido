// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Iterator traverses a range of a btree.
type Iterator struct {
	bt    *btree
	rng   Range
	tree  [maxLevels]*treeIter // tree[0] is root
	leaf  *leafIter
	state iterState
	buf   []byte
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

// Key returns the current key or nil.
// It returns an internal buffer which must be copied to be retained.
func (it *Iterator) Key() []byte {
	if it.state != within {
		return nil
	}
	return it.leaf.key(it.buf)
}

// Offset returns the current offset or 0.
func (it *Iterator) Offset() uint64 {
	if it.state != within {
		return 0
	}
	return it.leaf.offset()
}

func (it *Iterator) Eof() bool {
	return it.state == eof
}

func (it *Iterator) Modified() bool {
	return false
}

// Cur returns the current key and offset. It allocates the key.
func (it *Iterator) Cur() (string, uint64) {
	if it.state != within {
		return "", 0
	}
	return string(it.Key()), it.Offset()
}

// HasCur returns true if the iterator has a current item
func (it *Iterator) HasCur() bool {
	return it.state == within
}

// Rewind sets the iterator so Next goes to the first key in the range
// and Prev goes to the last key in the range
func (it *Iterator) Rewind() {
	it.state = rewound
}

// Range sets the range and rewinds the iterator
func (it *Iterator) Range(rng Range) {
	it.rng = rng
	it.Rewind()
}

//-------------------------------------------------------------------

// Next advances the iterator to the next key in the range or sets eof.
func (it *Iterator) Next() {
	switch it.state {
	case rewound:
		it.seek(it.rng.Org)
	case within:
		it.next()
	case eof: // stick at eof
		return
	}
	it.checkRange()
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

//-------------------------------------------------------------------

// Prev moves the iterator to the previous key in the range or sets eof.
func (it *Iterator) Prev() {
	switch it.state {
	case rewound:
		it.seek(it.rng.End)
		if it.Eof() {
			return // empty tree
		}
		if string(it.Key()) >= it.rng.End {
			it.prev()
		}
	case within:
		it.prev()
	case eof: // stick at eof
		return
	}
	it.checkRange()
}

func (it *Iterator) prev() {
	for {
		if it.leaf.prev() {
			it.state = within
			return
		} else if !it.prevLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) prevLeaf() bool {
	bt := it.bt
	i := bt.treeLevels - 1 // closest to leaf
	var nodeOff uint64
	// go up the tree until we can go back
	for ; i >= 0; i-- {
		if it.tree[i].prev() {
			nodeOff = it.tree[i].off()
			break
		} // else beginning of tree node, keep going up
	}
	if i < 0 {
		return false // beginning of root = eof
	}
	// then descend back down to rightmost
	for i++; i < bt.treeLevels; i++ {
		it.tree[i] = bt.getTree(nodeOff).iter()
		it.tree[i].i = it.tree[i].nd.noffs() - 1 // position at end
		nodeOff = it.tree[i].off()
	}
	it.leaf = bt.getLeaf(nodeOff).iter()
	it.leaf.i = it.leaf.nd.nkeys() // position at end
	return true
}

//-------------------------------------------------------------------

// Seek moves the iterator to the first position >= key.
// If the key is outside the current range, eof will be set.
func (it *Iterator) Seek(key string) {
	it.seek(key)
	it.checkRange()
}

// seek moves the iterator to the first position >= key.
// If the key is larger than the largest key,
// it will be positioned at the largest key.
// The state will be set to within
// unless the btree is empty in which case it will be set to eof.
// It does *not* apply the current range.
func (it *Iterator) seek(key string) {
	bt := it.bt
	off := bt.root
	for i := range bt.treeLevels {
		it.tree[i] = bt.getTree(off).seek(key)
		off = it.tree[i].off()
	}
	leaf := bt.getLeaf(off)
	it.leaf = leaf.seek(key)
	if leaf.nkeys() == 0 {
		assert.That(bt.treeLevels == 0) // only root can be empty
		it.state = eof
		return
	}
	if it.leaf.eof() {
		it.prev()
	}
	it.state = within
}

// checkRange changes state from within to eof
// if the current key is outside the range
func (it *Iterator) checkRange() {
	if it.state != within {
		return
	}
	curKey := string(it.Key())
	if curKey < it.rng.Org || it.rng.End <= curKey {
		it.state = eof
	}
}
