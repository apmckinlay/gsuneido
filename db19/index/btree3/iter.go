// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Iterator traverses a range of a btree.
type Iterator struct {
	bt      *btree
	rng     Range
	tree    [maxLevels]treeIter // tree[0] is root
	leaf    leafIter
	state   iterState
	buf     []byte
	noRange bool // true if rng is iterator.All, bypasses checkRange
}

// var _ iface.Iter = (*Iterator)(nil)

const maxLevels = 8

type Range = iface.Range

type iterState byte

const (
	rewound iterState = iota
	within
	eof
)

func (bt *btree) Iterator() iface.Iter {
	return &Iterator{bt: bt, state: rewound, rng: iface.All, noRange: true}
}

// Key returns the current key or nil.
// It returns an internal buffer which must be copied to be retained.
func (it *Iterator) Keybs() []byte {
	if it.state != within {
		return nil
	}
	it.buf = it.leaf.key(it.buf)
	return it.buf
}

// Key returns the current key or an empty string. It allocates.
func (it *Iterator) Key() string {
	if it.state != within {
		return ""
	}
	return string(it.leaf.key(it.buf))
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
	return it.Key(), it.Offset()
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
	it.noRange = (rng == iface.All)
	it.Rewind()
}

//-------------------------------------------------------------------

// Next advances the iterator to the next key in the range or sets eof.
func (it *Iterator) Next() {
	switch it.state {
	case rewound:
		it.SeekAll(it.rng.Org)
		it.checkRange() // need to check both bounds after seek
	case within:
		it.next()
		it.checkRangeEnd() // only need to check end when moving forward
	case eof: // stick at eof
		return
	}
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
			nodeOff = it.tree[i].offset()
			break
		} // else end of tree node, keep going up
	}
	if i < 0 {
		return false // end of root = eof
	}
	// then descend back down
	for i++; i < bt.treeLevels; i++ {
		it.tree[i] = bt.readTree(nodeOff).iter()
		assert.That(it.tree[i].next())
		nodeOff = it.tree[i].offset()
	}
	it.leaf = bt.readLeaf(nodeOff).iter()
	return true
}

//-------------------------------------------------------------------

// Prev moves the iterator to the previous key in the range or sets eof.
func (it *Iterator) Prev() {
	switch it.state {
	case rewound:
		it.SeekAll(it.rng.End)
		if it.Eof() {
			return // empty tree
		}
		if string(it.Key()) >= it.rng.End {
			it.prev()
		}
		it.checkRange() // need to check both bounds after seek
	case within:
		it.prev()
		it.checkRangeOrg() // only need to check org when moving backward
	case eof: // stick at eof
		return
	}
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
			nodeOff = it.tree[i].offset()
			break
		} // else beginning of tree node, keep going up
	}
	if i < 0 {
		return false // beginning of root = eof
	}
	// then descend back down to rightmost
	for i++; i < bt.treeLevels; i++ {
		it.tree[i] = bt.readTree(nodeOff).iter()
		it.tree[i].i = it.tree[i].nd.noffs() - 1 // position at end
		nodeOff = it.tree[i].offset()
	}
	it.leaf = bt.readLeaf(nodeOff).iter()
	it.leaf.i = it.leaf.nd.nkeys() // position at end
	return true
}

//-------------------------------------------------------------------

// Seek moves the iterator to the first position >= key.
// If the key is outside the current range, eof will be set.
func (it *Iterator) Seek(key string) {
	it.SeekAll(key)
	it.checkRange()
}

// SeekAll moves the iterator to the first position >= key.
// If the key is larger than the largest key,
// it will be positioned at the largest key.
// The state will be set to within
// unless the btree is empty in which case it will be set to eof.
// It does *not* apply the current range.
func (it *Iterator) SeekAll(key string) {
	bt := it.bt
	off := bt.root
	rightEdge := true
	for i := range bt.treeLevels {
		it.tree[i] = bt.readTree(off).seek(key)
		off = it.tree[i].offset()
		rightEdge = rightEdge && it.tree[i].i >= it.tree[i].nd.nkeys()
	}
	leaf := bt.readLeaf(off)
	if leaf.nkeys() == 0 {
		assert.That(bt.treeLevels == 0) // only root can be empty
		it.state = eof
		return
	}
	it.leaf = leaf.seek(key)
	if it.leaf.eof() {
		if rightEdge {
			it.prev()
		} else {
			it.next()
		}
	}
	it.state = within
}

// checkRange changes state from within to eof
// if the current key is outside the range
func (it *Iterator) checkRange() {
	if it.noRange || it.state != within {
		return
	}
	prefix := it.leaf.prefix()
	suffix := it.leaf.suffix()
	if gte(prefix, suffix, it.rng.End) || !gte(prefix, suffix, it.rng.Org) {
		it.state = eof
	}
}

// checkRangeEnd changes state from within to eof
// if the current key is >= rng.End (used for Next)
func (it *Iterator) checkRangeEnd() {
	if it.noRange || it.state != within {
		return
	}
	prefix := it.leaf.prefix()
	suffix := it.leaf.suffix()
	if gte(prefix, suffix, it.rng.End) {
		it.state = eof
	}
}

// checkRangeOrg changes state from within to eof
// if the current key is < rng.Org (used for Prev)
func (it *Iterator) checkRangeOrg() {
	if it.noRange || it.state != within {
		return
	}
	prefix := it.leaf.prefix()
	suffix := it.leaf.suffix()
	if !gte(prefix, suffix, it.rng.Org) {
		it.state = eof
	}
}

// gte returns true if prefix+suffix >= target
// without concatenating prefix and suffix
func gte(prefix, suffix []byte, bound string) bool {
	plen := len(prefix)
	slen := len(suffix)
	tlen := len(bound)

	// Compare the prefix portion first
	cmpLen := min(tlen, plen)
	for i := 0; i < cmpLen; i++ {
		if prefix[i] < bound[i] {
			return false
		}
		if prefix[i] > bound[i] {
			return true
		}
	}

	// If bound is entirely within prefix length, check if we have more data
	if tlen <= plen {
		return plen+slen >= tlen
	}

	// Compare the suffix portion
	boundOffset := plen
	remainingBound := tlen - plen
	cmpLen = min(remainingBound, slen)
	for i := 0; i < cmpLen; i++ {
		if suffix[i] < bound[boundOffset+i] {
			return false
		}
		if suffix[i] > bound[boundOffset+i] {
			return true
		}
	}

	return plen+slen >= tlen
}
