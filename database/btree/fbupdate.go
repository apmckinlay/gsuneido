// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/database/stor"
)

// fbupdate handles updating an immutable btree.
// It creates new and updated nodes in memory.
type fbupdate struct {
	// fb is a copy of the fbtree (so we can update it)
	fb fbtree
	// moffs (mutable) maps offsets to mutable new and updated in-memory nodes
	moffs memOffsets
	// paths is a set of nodes that will need to be path copied
	paths map[uint64]set
}

type set struct{}

var mark set

func (fb *fbtree) Update(fn func(*fbupdate)) *fbtree {
	up := newFbupdate(fb)
	fn(up)
	up.freeze()
	return &up.fb
}

func newFbupdate(fb *fbtree) *fbupdate {
	moffs := memOffsets{nextOff: fb.moffs.nextOff,
		nodes: make(map[uint64]fNode, len(fb.moffs.nodes))}
	return &fbupdate{fb: *fb, moffs: moffs,
		paths: make(map[uint64]set, len(fb.paths))}
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
const maxlevels = 8

func (up *fbupdate) Insert(key string, off uint64) {
	var stack [maxlevels]uint64

	// search down the tree to the appropriate leaf
	nodeOff := up.fb.root
	for i := 0; i < up.fb.treeLevels; i++ {
		stack[i] = nodeOff
		node := up.getPathNode(nodeOff)
		nodeOff, _, _ = node.search(key)
	}

	// insert into leaf
	node, where := up.insert(nodeOff, key, off, up.fb.getLeafKey)
	up.moffs.set(nodeOff, node)
	size := len(node)
	if size <= up.fb.maxNodeSize {
		return // fast path, just insert into leaf
	}

	// split leaf
	splitKey, rightOff := up.split(node, nodeOff, where)

	// insert up the tree
	for i := up.fb.treeLevels - 1; i >= 0; i-- {
		node, where = up.insert(stack[i], splitKey, rightOff, nil)
		up.moffs.set(stack[i], node)
		size := len(node)
		if size <= up.fb.maxNodeSize {
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

func (up *fbupdate) insert(nodeOff uint64, key string, off uint64,
	get func(uint64) string) (fNode, int) {
	node := up.getMutableNode(nodeOff)
	return node.insert(key, off, get)
}

func (up *fbupdate) getNode(off uint64) fNode {
	if node := up.moffs.get(off); node != nil {
		return node
	}
	return up.fb.getNode(off)
}

func (up *fbupdate) getPathNode(off uint64) fNode {
	up.paths[off] = mark
	return up.getNode(off)
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
	it := node.iter()
	for it.next() && it.fi < splitSize {
	}
	splitKey = it.known

	left := node[:it.fi]
	up.moffs.set(nodeOff, left)
	right := make(fNode, 0, len(node)-it.fi+8)
	// first entry becomes 0, ""
	right = fAppend(right, it.offset, 0, "")
	if it.next() {
		// second entry becomes 0, known
		right = fAppend(right, it.offset, 0, it.known)
		if it.next() {
			right = append(right, node[it.fi:]...)
		}
	}
	rightOff = up.moffs.add(right)
	return
}

//-------------------------------------------------------------------

func (up *fbupdate) Delete(key string, off uint64) bool {
	var stack [maxlevels]uint64

	// search down the tree to the appropriate leaf
	nodeOff := up.fb.root
	for i := 0; i < up.fb.treeLevels; i++ {
		stack[i] = nodeOff
		node := up.getPathNode(nodeOff)
		nodeOff, _, _ = node.search(key)
	}

	// delete from leaf
	node, ok := up.delete(nodeOff, off)
	if !ok {
		return false
	}
	if len(node) != 0 || up.fb.treeLevels == 0 {
		return true // usual fast path
	}

	// delete up the tree
	for i := up.fb.treeLevels - 1; i >= 0; i-- {
		node, ok = up.delete(stack[i], nodeOff)
		if !ok {
			panic("leaf node not found in tree")
		}
		if (i > 0 || up.fb.treeLevels == 0) && len(node) != 0 {
			return true
		}
		nodeOff = stack[i]
	}

	// remove empty root(s)
	for up.fb.treeLevels > 0 && len(node) == 7 {
		up.fb.treeLevels--
		up.fb.root = stor.ReadSmallOffset(node)
		node = up.getNode(up.fb.root)
	}

	return true
}

func (up *fbupdate) delete(nodeOff uint64, off uint64) (fNode, bool) {
	node := up.getMutableNode(nodeOff)
	node, ok := node.delete(off)
	if ok {
		up.moffs.set(nodeOff, node)
	}
	return node, ok
}

//-------------------------------------------------------------------

// freeze moves the changes to the fbtree.
// It will still reference in-memory new and updated nodes
// but they are no longer mutable.
func (up *fbupdate) freeze() *fbtree {
	for k, v := range up.fb.moffs.nodes {
		if _, ok := up.moffs.nodes[k]; !ok {
			up.moffs.nodes[k] = v
		}
	}
	up.fb.moffs = up.moffs
	up.moffs.nodes = nil
	for k := range up.fb.paths {
		up.paths[k] = mark
	}
	up.fb.paths = up.paths
	up.paths = nil
	return &up.fb
}

//-------------------------------------------------------------------

// memOffsets is use to map offsets to temporary, in-memory, mutable nodes.
// It is used for updated versions of existing nodes, using their old offset,
// and for new nodes with fake offsets.
type memOffsets struct {
	nodes   map[uint64]fNode
	nextOff uint64
}

func nilMemOffsets() memOffsets {
	return memOffsets{nextOff: stor.MaxSmallOffset}
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

// ------------------------------------------------------------------

func (up *fbupdate) print() {
	up.fb.print()
}

// check verifies that the keys are in order and returns the number of keys
func (up *fbupdate) check() (count, size, nnodes int) {
	return up.check1(0, up.fb.root, "")
}

func (up *fbupdate) check1(depth int, offset uint64, key string) (count, size, nnodes int) {
	node := up.getNode(offset)
	size += len(node)
	nnodes++
	for it := node.iter(); it.next(); {
		offset := it.offset
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
			itkey := up.fb.getLeafKey(offset)
			if key > itkey {
				panic("keys out of order " + key + " " + itkey)
			}
			count++
		}
	}
	return
}
