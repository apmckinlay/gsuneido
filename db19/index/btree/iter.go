// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Iterator is a Suneido style iterator with Next/Prev/Rewind
type Iterator struct {
	bt    *btree
	state iterState
	// stack is the path to the current position
	// stack[0] is the leaf, stack[treeLevels] is the root
	stack  [maxlevels + 1]iNodeIter
	curKey string
	curOff uint64
	// rng is the Range of the iterator
	rng Range
}

type Range = iterator.Range

var _ iterator.T = (*Iterator)(nil)

func (bt *btree) Iterator() *Iterator {
	return &Iterator{bt: bt, state: rewound, rng: iterator.All}
}

type iterState byte

const (
	rewound iterState = iota
	within
	eof
)

func (it *Iterator) Range(rng Range) {
	it.rng = rng
	it.Rewind()
}

func (it *Iterator) Eof() bool {
	return it.state == eof
}

func (it *Iterator) Cur() (string, uint64) {
	assert.Msg("Cur when not within").That(it.state == within)
	return it.curKey, it.curOff
}

func (it *Iterator) Rewind() {
	it.state = rewound
}

func (it *Iterator) Modified() bool {
	return false
}

// Next -------------------------------------------------------------

func (it *Iterator) Next() {
	if it.Eof() {
		return // stick at eof
	}
	if it.state == rewound {
		it.SeekAll(it.rng.Org)
	} else {
		it.next()
	}
	it.checkRange()
}

func (it *Iterator) next() {
	for {
		if it.stack[0].next() {
			it.curOff = it.stack[0].off()
			it.curKey = it.bt.getLeafKey(it.curOff)
			return
		} else if !it.nextLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) nextLeaf() bool {
	bt := it.bt
	i := 1
	var nodeOff uint64
	// end of leaf, go up the tree as necessary
	for ; i <= bt.treeLevels; i++ {
		if it.stack[i].next() {
			nodeOff = it.stack[i].off()
			break
		} // else end of tree node, keep going up
	}
	if i > bt.treeLevels {
		return false // eof
	}
	// then descend back down
	for {
		i--
		it.stack[i] = bt.getNode(nodeOff).iter()
		if i == 0 {
			return true
		}
		it.stack[i].next()
		nodeOff = it.stack[i].off()
	}
}

// Prev -------------------------------------------------------------

func (it *Iterator) Prev() {
	if it.state == eof {
		return // stick at eof
	}
	if it.state == rewound {
		it.SeekAll(it.rng.End)
		if it.Eof() {
			return
		}
		if it.curKey < it.rng.End {
			it.checkRange()
			return
		}
		it.state = within
		// Seek went to >= so fallthrough to do previous
	}
	it.prev()
	it.checkRange()
}

func (it *Iterator) prev() {
	for {
		it.stack[0] = it.stack[0].toChunk(it.bt, true)
		if it.stack[0].prev() {
			it.curOff = it.stack[0].off()
			it.curKey = it.bt.getLeafKey(it.curOff)
			return
		} else if !it.prevLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) prevLeaf() bool {
	bt := it.bt
	i := 1
	var nodeOff uint64
	// go up the tree as necessary
	for ; i <= bt.treeLevels; i++ {
		it.stack[i] = it.stack[i].toChunk(it.bt, i == 0)
		if it.stack[i].prev() {
			nodeOff = it.stack[i].off()
			break
		} // else end of tree node, keep going up
	}
	if i > bt.treeLevels {
		return false // eof
	}
	// then descend back down
	for {
		i--
		it.stack[i] = bt.getNode(nodeOff).iter().toChunk(bt, i == 0)
		if i == 0 {
			return true
		}
		it.stack[i].prev()
		nodeOff = it.stack[i].off()
	}
}

// Seek -------------------------------------------------------------

func (it *Iterator) Seek(key string) {
	it.SeekAll(key)
	it.checkRange()
}

func (it *Iterator) SeekAll(key string) {
	bt := it.bt
	off := bt.root
	nd := bt.getNode(off)
	if len(nd) == 0 {
		it.state = eof
		return
	}
	rightEdge := true
	it.state = within
	for i := bt.treeLevels; ; i-- {
		ni := nd.seek(key)
		rightEdge = rightEdge && ni.eof()
		it.stack[i] = ni
		off = it.stack[i].off()
		if i == 0 {
			k := bt.getLeafKey(off)
			if key > k && !rightEdge {
				it.next()
				return
			}
			it.curOff = off
			it.curKey = k
			return
		}
		nd = bt.getNode(off)
	}
}

func (nd node) seek(key string) *nodeIter {
	// similar to node.search
	it := nd.iter()
	itPrev := nodeIter{node: nd}
	for it.next() && key >= string(it.known) {
		itPrev.copyFrom(it)
	}
	return &itPrev
}

func (it *Iterator) checkRange() {
	if it.curKey < it.rng.Org || it.rng.End <= it.curKey {
		it.curOff = 0
		it.curKey = ""
		it.state = eof
	}
}

//-------------------------------------------------------------------

type iNodeIter interface {
	// next returns false if it hits the end
	next() bool
	// prev returns false if it hits the start
	prev() bool
	// off returns the current offset
	off() uint64
	// toChunk converts nodeIter to chunkIter to allow Prev
	toChunk(bt *btree, leaf bool) iNodeIter
	// eof returns true if on the last slot
	eof() bool
}

func (it *nodeIter) off() uint64 {
	return it.offset
}

func (it *nodeIter) prev() bool {
	assert.ShouldNotReachHere()
	return false
}

// toChunk converts a nodeIter to a chunkIter to allow prev
func (it *nodeIter) toChunk(bt *btree, leaf bool) iNodeIter {
	nd := it.node
	var c chunk
	i := -1
	var key string
	it2 := nd.iter()
	for it2.next() {
		if it2.pos == it.pos {
			i = len(c)
		}
		if leaf {
			key = bt.getLeafKey(it2.offset)
		} else {
			key = string(it2.known)
		}
		c = append(c, slot{key: key, off: it2.offset})
	}
	return &chunkIter{c: c, i: i}
}

type chunk []slot

type slot struct {
	key string
	off uint64
}

type chunkIter struct {
	c chunk
	i int
}

func (ci *chunkIter) eof() bool {
	return ci.i+1 >= len(ci.c)
}

func (ci *chunkIter) next() bool {
	ci.i++
	return ci.i < len(ci.c)
}

func (ci *chunkIter) prev() bool {
	if ci.i == -1 {
		ci.i = len(ci.c)
	}
	ci.i--
	return ci.i >= 0
}

func (ci *chunkIter) off() uint64 {
	return ci.c[ci.i].off
}

func (ci *chunkIter) toChunk(*btree, bool) iNodeIter {
	return ci
}
