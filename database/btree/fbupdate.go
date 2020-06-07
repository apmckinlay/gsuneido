// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/database/stor"
)

// fbupdate handles updating an immutable btree.
// It creates new and updated nodes in memory.
type fbupdate struct {
	fb fbtree
	// moffs assigns temporary offsets to new and updated nodes
	moffs memOffsets
	// getLeafKey returns the key associated with a leaf data offset
	getLeafKey func(uint64) string
	// maxNodeSize is the maximum node size in bytes, split if larger
	maxNodeSize int
}

func (up *fbupdate) Search(key string) uint64 {
	nodeOff := up.fb.root
	for i := 0; i <= up.fb.treeLevels; i++ {
		node := up.getNode(nodeOff)
		nodeOff, _, _ = node.search(key)
	}
	return nodeOff
}

const maxNodeSize = 1536 // * .75 ~ 1k

func (up *fbupdate) Insert(key string, off uint64) {
	const maxlevels = 8
	var stack [maxlevels]uint64

	// search down the tree to the appropriate leaf
	nodeOff := up.fb.root
	for i := 0; i < up.fb.treeLevels; i++ {
		stack[i] = nodeOff
		node := up.getMutableNode(nodeOff) // mutable to mark path for save
		nodeOff, _, _ = node.search(key)
	}

	// insert into leaf
	node := up.getMutableNode(nodeOff)
	var where int
	node, where = node.insert(key, off, up.getLeafKey)
	up.moffs.set(nodeOff, node)
	size := len(node)
	if size <= up.maxNodeSize {
		return // fast path, just insert into leaf
	}

	// split leaf
	splitKey, rightOff := up.split(node, nodeOff, where)

	// insert up the tree
	for i := up.fb.treeLevels - 1; i >= 0; i-- {
		node = up.getMutableNode(stack[i])
		node, where = node.insert(splitKey, rightOff, nil)
		up.moffs.set(stack[i], node)
		size := len(node)
		if size <= up.maxNodeSize {
			return // finished
		}
		splitKey, rightOff = up.split(node, stack[i], where)
	}

	// split all the way up, create new root
	newRoot := make(fNode, 0, 24)
	newRoot = fAppend(newRoot, uint64(up.fb.root), 0, "")
	newRoot = fAppend(newRoot, uint64(rightOff), 0, splitKey)
	up.fb.root = up.moffs.add(newRoot)
	up.fb.treeLevels++
}

func (up *fbupdate) getNode(off uint64) fNode {
	if node := up.moffs.get(off); node != nil {
		return node
	}
	return up.fb.getNode(off)
}

func (up *fbupdate) getMutableNode(off uint64) fNode {
	if node := up.moffs.get(off); node != nil {
		return node
	}
	imu := up.fb.getNode(off)
	node := append(imu[:0:0], imu...)
	up.moffs.set(off, node)
	return node
}

func (up *fbupdate) split(node fNode, nodeOff uint64, where int) (
	splitKey string, rightOff uint64) {
	size := len(node)
	splitSize := size / 2
	if where == insStart {
		splitSize = size / 4
	} else if where == insEnd {
		splitSize = (size * 3) / 4
	}
	it := node.Iter()
	for it.next() && it.fi < splitSize {
	}
	splitKey = it.known

	left := node[:it.fi]
	up.moffs.set(nodeOff, left)
	right := make(fNode, 0, len(node)-it.fi+8)
	// first entry becomes 0, ""
	right = fAppend(right, node.offset(it.fi), 0, "")
	if it.next() {
		// second entry becomes 0, known
		right = fAppend(right, node.offset(it.fi), 0, it.known)
		if it.next() {
			right = append(right, node[it.fi:]...)
		}
	}
	rightOff = up.moffs.add(right)
	return
}

// save writes the btree (changes) to the stor and returns the new root offset.
// It works bottom (leaves) up so it can replace memOffsets with stor offsets.
func (up *fbupdate) save() {
	up.fb.root = up.save2(0, up.fb.root)
	up.moffs.clear()
}

func (up *fbupdate) save2(depth int, nodeOff uint64) uint64 {
	// only traverse modified paths, not entire (possibly huge) tree
	if up.inStor(nodeOff) {
		return nodeOff
	}
	node := up.getNode(nodeOff) // mutable since not in stor
	if depth < up.fb.treeLevels {
		// tree node, need to update any memOffsets
		for it := node.Iter(); it.next(); {
			off := node.offset(it.fi)
			off2 := up.save2(depth+1, off) // recurse
			// bottom up
			if off2 != off {
				node.setOffset(it.fi, off2)
			}
		}
	}
	return up.fb.putNode(node)
}

func (up *fbupdate) inStor(off uint64) bool {
	_, ok := up.moffs.nodes[off]
	return !ok
}

//-------------------------------------------------------------------

// memOffsets is use to map offsets to temporary, in-memory, mutable nodes.
// It is used for updated versions of existing nodes, using their old offset,
// and for new nodes with fake offsets.
type memOffsets struct {
	nodes   map[uint64]fNode
	nextOff uint64
}

func newMemOffsets() memOffsets {
	return memOffsets{nodes: map[uint64]fNode{}, nextOff: stor.MaxSmallOffset}
}

// add returns a fake offset for a new node
func (mo *memOffsets) add(node fNode) uint64 {
	off := mo.nextOff
	mo.nextOff--
	mo.nodes[off] = node
	return off
}

// set updates the node for an offset
func (mo *memOffsets) set(off uint64, node fNode) {
	mo.nodes[off] = node
}

// get retrieves a node given its offset.
func (mo *memOffsets) get(off uint64) fNode {
	return mo.nodes[off]
}

func (mo *memOffsets) clear() {
	*mo = newMemOffsets()
}

// ------------------------------------------------------------------

func (up *fbupdate) print() {
	up.print1(0, up.fb.root)
}

func (up *fbupdate) print1(depth int, offset uint64) {
	node := up.getNode(offset)
	for it := node.Iter(); it.next(); {
		offset := node.offset(it.fi)
		if depth < up.fb.treeLevels {
			// tree
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known)
			up.print1(depth+1, offset) // recurse
		} else {
			// leaf
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known, "("+up.getLeafKey(offset)+")")
		}
	}
}

// check verifies that the keys are in order and returns the number of keys
func (up *fbupdate) check() (count, size, nnodes int) {
	return up.check1(0, up.fb.root, "")
}

func (up *fbupdate) check1(depth int, offset uint64, key string) (count, size, nnodes int) {
	node := up.getNode(offset)
	size += len(node)
	nnodes++
	for it := node.Iter(); it.next(); {
		offset := node.offset(it.fi)
		if depth < up.fb.treeLevels {
			// tree
			if it.fi > 0 && key > it.known {
				panic("keys out of order " + key + " " + it.known)
			}
			c, s, n := up.check1(depth+1, offset, it.known) // recurse
			count += c
			size += s
			nnodes += n
		} else {
			// leaf
			itkey := up.getLeafKey(offset)
			if key > itkey {
				panic("keys out of order " + key + " " + itkey)
			}
			count++
		}
	}
	return
}
