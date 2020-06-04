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
	bt fbtree
	// fstor is where the immutable base btree is stored
	fstor stor.Stor
	// moffs assigns temporary offsets to new and updated nodes
	moffs memOffsets
	// getLeafKey returns the key associated with a leaf data offset
	getLeafKey func(uint64) string
	// maxNodeSize is the maximum node size in bytes, split if larger
	maxNodeSize int
}

const maxNodeSize = 1536 // * .75 ~ 1k

func (up *fbupdate) Insert(key string, off uint64) {
	const maxlevels = 8
	var stack [maxlevels]*fNode
	var stackOff [maxlevels]uint64

	// search down the tree to the appropriate leaf
	nodeOff := up.bt.root
	for i := 0; i < up.bt.treeLevels; i++ {
		stackOff[i] = nodeOff
		stack[i] = up.getNode(nodeOff)
		node := *stack[i]
		nodeOff, _, _ = node.search(key)
	}

	// insert into leaf
	pnode := up.getNode(nodeOff)
	//TODO copy on write from fstor to moffs
	*pnode = pnode.insert(key, off, up.getLeafKey)
	node := *pnode
	size := len(node)
	if size <= up.maxNodeSize {
		return // fast path, just insert into leaf
	}

	// split leaf
	splitKey, rightOff := up.split(pnode)

	// insert up the tree
	for i := up.bt.treeLevels - 1; i >= 0; i-- {
		//TODO copy on write from fstor to moffs
		pnode = stack[i]
		*pnode = pnode.insert(splitKey, rightOff, nil)
		node = *pnode
		size := len(node)
		if size <= up.maxNodeSize {
			return // finished
		}
		splitKey, rightOff = up.split(pnode)
	}

	// split all the way up, create new root
	newRoot := make(fNode, 0, 24)
	newRoot = fAppend(newRoot, uint64(up.bt.root), 0, "")
	newRoot = fAppend(newRoot, uint64(rightOff), 0, splitKey)
	up.bt.root = up.moffs.add(newRoot)
	up.bt.treeLevels++
}

func (up *fbupdate) getNode(off uint64) *fNode {
	return up.moffs.get(off)
}

func (up *fbupdate) split(pnode *fNode) (splitKey string, rightOff uint64) {
	node := *pnode
	size := len(node)
	half := size / 2
	it := node.Iter()
	for it.next() && it.fi < half {
	}
	splitKey = it.known
	//TODO uneven split
	left := node[:it.fi]
	*pnode = left
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

//-------------------------------------------------------------------

// memOffsets is use to assign offsets to temporary, in-memory, mutable nodes.
// zero value is ready to go.
type memOffsets struct {
	nodes []*fNode
}

// add stores the node and returns its assigned offset
func (mo *memOffsets) add(node fNode) uint64 {
	i := len(mo.nodes)
	mo.nodes = append(mo.nodes, &node)
	return stor.MaxSmallOffset - uint64(i)
}

// get retrieves a node given its offset.
// It returns *fNode because fNode is a slice and we may need to extend it.
func (mo *memOffsets) get(off uint64) *fNode {
	i := stor.MaxSmallOffset - off
	if i > uint64(len(mo.nodes)) {
		return nil
	}
	return mo.nodes[i]
}

// ------------------------------------------------------------------

func (up *fbupdate) print() {
	up.print1(0, up.bt.root)
}

func (up *fbupdate) print1(depth int, offset uint64) {
	node := *up.getNode(offset)
	for it := node.Iter(); it.next(); {
		offset := node.offset(it.fi)
		if depth < up.bt.treeLevels {
			// tree
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known)
			up.print1(depth+1, offset)
		} else {
			// leaf
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known, "("+up.getLeafKey(offset)+")")
		}
	}
}

// check verifies that the keys are in order and returns the number of keys
func (up *fbupdate) check() (count, size, nnodes int) {
	return up.check1(0, up.bt.root, "")
}

func (up *fbupdate) check1(depth int, offset uint64, key string) (count, size, nnodes int) {
	node := *up.getNode(offset)
	size += len(node)
	nnodes++
	for it := node.Iter(); it.next(); {
		offset := node.offset(it.fi)
		if depth < up.bt.treeLevels {
			// tree
			if it.fi > 0 && key > it.known {
				panic("keys out of order " + key + " " + it.known)
			}
			c,s,n := up.check1(depth+1, offset, it.known)
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
