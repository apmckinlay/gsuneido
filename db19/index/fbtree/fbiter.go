// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

import "github.com/apmckinlay/gsuneido/db19/index/iterator"

// Iterator is a Suneido style iterator with Next/Prev/Rewind
type Iterator struct {
	fb    *fbtree
	state iterState
	// stack is the path to the current position
	// stack[0] is the leaf, stack[treeLevels] is the root
	stack  [maxlevels]*fnIter
	curKey string
	curOff uint64
}

var _ iterator.T = (*Iterator)(nil)

func (fb *fbtree) Iterator() *Iterator {
	return &Iterator{fb: fb, state: rewound}
}

type chunk []slot

type slot struct {
	key string
	off uint64
}

type iterState byte

const (
	rewound iterState = iota
	within
	eof
)

func (it *Iterator) Eof() bool {
	return it.state == eof
}

func (it *Iterator) Cur() (string, uint64) {
	return it.curKey, it.curOff
}

func (it *Iterator) Rewind() {
	it.state = rewound
}

func (it *Iterator) Modified() bool {
	return false //TODO ???
}

func (it *Iterator) Prev() {
}

func (it *Iterator) Next() {
	if it.state == eof {
		return // stick at eof
	}
	if it.state == rewound {
		it.descendLeft()
		it.state = within
	}
	it.next()
}

func (it *Iterator) next() {
	for {
		if it.stack[0].next() {
			it.curOff = it.stack[0].offset
			it.curKey = it.fb.getLeafKey(it.curOff)
			return
		} else if !it.nextLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) descendLeft() { // maybe use Seek ???
	fb := it.fb
	nodeOff := fb.root
	for i := fb.treeLevels; ; i-- {
		it.stack[i] = fb.getNode(nodeOff).iter()
		if i == 0 {
			return
		}
		it.stack[i].next()
		nodeOff = it.stack[i].offset
	}
}

func (it *Iterator) nextLeaf() bool {
	fb := it.fb
	i := 1
	var nodeOff uint64
	// end of leaf, go up the tree as necessary
	for ; i <= fb.treeLevels; i++ {
		if it.stack[i].next() {
			nodeOff = it.stack[i].offset
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
		nodeOff = it.stack[i].offset
	}
}

// Seek returns true if the key was found
func (it *Iterator) Seek(key string) bool {
	fb := it.fb
	off := fb.root
	node := fb.getNode(off)
	if len(node) == 0 {
		it.state = eof
		return false
	}
	it.state = within
	for i := fb.treeLevels; ; i-- {
		it.stack[i] = node.seek(key)
		off = it.stack[i].offset
		if i == 0 {
			k := fb.getLeafKey(off)
			if key > k {
				it.next()
				return false
			}
			it.curOff = off
			it.curKey = k
			return key == k
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

// fnodeToChunk is used when we need to iterate reverse (Prev)
// since we can't do that on an fnode.
func (it *Iterator) fnodeToChunk(fn fnode, leaf bool) chunk {
	var c chunk
	var key string
	fit := fn.iter()
	for fit.next() {
		if leaf {
			key = it.fb.getLeafKey(fit.offset)
		} else {
			key = string(fit.known)
		}
		c = append(c, slot{key: key, off: fit.offset})
	}
	return c
}
