// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

import (
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Iterator is a Suneido style iterator with Next/Prev/Rewind
type Iterator struct {
	fb    *fbtree
	state iterState
	// stack is the path to the current position
	// stack[0] is the leaf, stack[treeLevels] is the root
	stack  [maxlevels]nodeIter
	curKey string
	curOff uint64
	// rng is the Range of the iterator
	rng Range
}

type Range = iterator.Range

var _ iterator.T = (*Iterator)(nil)

func (fb *fbtree) Iterator() *Iterator {
	return &Iterator{fb: fb, state: rewound, rng: iterator.All}
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
	if it.state == eof {
		return // stick at eof
	}
	if it.state == rewound {
		it.Seek(it.rng.Org)
		return
	}
	it.next()
}

func (it *Iterator) next() {
	for {
		if it.stack[0].next() {
			it.curOff = it.stack[0].off()
			it.curKey = it.fb.getLeafKey(it.curOff)
			return
		} else if !it.nextLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) nextLeaf() bool {
	fb := it.fb
	i := 1
	var nodeOff uint64
	// end of leaf, go up the tree as necessary
	for ; i <= fb.treeLevels; i++ {
		if it.stack[i].next() {
			nodeOff = it.stack[i].off()
			break
		} // else end of tree node, keep going up
	}
	if i > fb.treeLevels {
		return false // eof
	}
	// then descend back down
	for {
		i--
		it.stack[i] = fb.getNode(nodeOff).iter()
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
		it.Seek(it.rng.End)
		if it.Eof() {
			return
		}
		key, _ := it.Cur()
		if key < it.rng.End {
			return
		}
		it.state = within
		// Seek went to >= so fallthrough to do previous
	}
	it.prev()
}

func (it *Iterator) prev() {
	for {
		it.stack[0] = it.stack[0].toChunk(it.fb, true)
		if it.stack[0].prev() {
			it.curOff = it.stack[0].off()
			it.curKey = it.fb.getLeafKey(it.curOff)
			return
		} else if !it.prevLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) prevLeaf() bool {
	fb := it.fb
	i := 1
	var nodeOff uint64
	// go up the tree as necessary
	for ; i <= fb.treeLevels; i++ {
		it.stack[i] = it.stack[i].toChunk(it.fb, i == 0)
		if it.stack[i].prev() {
			nodeOff = it.stack[i].off()
			break
		} // else end of tree node, keep going up
	}
	if i > fb.treeLevels {
		return false // eof
	}
	// then descend back down
	for {
		i--
		it.stack[i] = fb.getNode(nodeOff).iter().toChunk(fb, i == 0)
		if i == 0 {
			return true
		}
		it.stack[i].prev()
		nodeOff = it.stack[i].off()
	}
}

// Seek -------------------------------------------------------------

func (it *Iterator) Seek(key string) {
	fb := it.fb
	off := fb.root
	node := fb.getNode(off)
	if len(node) == 0 {
		it.state = eof
		return
	}
	rightEdge := true
	it.state = within
	for i := fb.treeLevels; ; i-- {
		fit := node.seek(key)
		rightEdge = rightEdge && fit.eof()
		it.stack[i] = fit
		off = it.stack[i].off()
		if i == 0 {
			k := fb.getLeafKey(off)
			if key > k && !rightEdge {
				it.next()
				return
			}
			it.curOff = off
			it.curKey = k
			return
		}
		node = fb.getNode(off)
	}
}

func (fn fnode) seek(key string) *fnIter {
	// similar to fnode.search
	it := fn.iter()
	itPrev := fnIter{fn: fn}
	for it.next() && key >= string(it.known) {
		itPrev.copyFrom(it)
	}
	return &itPrev
}

//-------------------------------------------------------------------

type nodeIter interface {
	// next returns false if it hits the end
	next() bool
	// prev returns false if it hits the start
	prev() bool
	// off returns the current offset
	off() uint64
	// toChunk converts fnIter to chunkIter, to allow Prev
	toChunk(fb *fbtree, leaf bool) nodeIter
	// eof returns true if on the last slot
	eof() bool
}

func (fi *fnIter) off() uint64 {
	return fi.offset
}

func (fi *fnIter) prev() bool {
	panic("shouldn't get here")
}

// toChunk converts an fnIter to a chunkIter (to allow prev)
func (fi *fnIter) toChunk(fb *fbtree, leaf bool) nodeIter {
	fn := fi.fn
	var c chunk
	i := -1
	var key string
	fit := fn.iter()
	for fit.next() {
		if fit.fi == fi.fi {
			i = len(c)
		}
		if leaf {
			key = fb.getLeafKey(fit.offset)
		} else {
			key = string(fit.known)
		}
		c = append(c, slot{key: key, off: fit.offset})
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

func (ci *chunkIter) toChunk(*fbtree, bool) nodeIter {
	return ci
}
