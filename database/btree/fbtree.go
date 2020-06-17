// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/database/stor"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// fbtree is an immutable btree designed to be stored in a file
type fbtree struct {
	// treeLevels is how many levels of tree nodes there are (initially 0)
	treeLevels int
	// root is the offset of the root node
	root uint64
	// store is where the btree is stored
	store *stor.Stor
	// moffs (reaonly) maps offsets to updated in-memory nodes (not persisted)
	moffs memOffsets
	// paths is a set of nodes that will need to be path copied
	paths map[uint64]set
	// getLeafKey returns the key associated with a leaf data offset
	getLeafKey func(uint64) string
	// maxNodeSize is the maximum node size in bytes, split if larger
	maxNodeSize int
}

func CreateFbtree(st *stor.Stor,
	getLeafKey func(uint64) string, maxNodeSize int) *fbtree {
	mo := memOffsets{nextOff: stor.MaxSmallOffset, nodes: map[uint64]fNode{}}
	root := mo.add(fNode{})
	return &fbtree{root: root, moffs: mo, store: st,
		getLeafKey: getLeafKey, maxNodeSize: maxNodeSize}
}

func OpenFbtree(st *stor.Stor, root uint64,
	getLeafKey func(uint64) string, maxNodeSize int) *fbtree {
	mo := nilMemOffsets()
	return &fbtree{root: root, moffs: mo, store: st,
		getLeafKey: getLeafKey, maxNodeSize: maxNodeSize}
}

func (fb *fbtree) Search(key string) uint64 {
	nodeOff := fb.root
	for i := 0; i <= fb.treeLevels; i++ {
		node := fb.getNode(nodeOff)
		nodeOff, _, _ = node.search(key)
	}
	return nodeOff
}

// save writes the btree (changes) to the stor
// and returns a new fbtree with no in-memory nodes.
func (fb *fbtree) save() *fbtree {
	fb = fb.Update(func(up *fbupdate) {
		up.fb.root = up.save2(0, fb.root)
	})
	fb.moffs = nilMemOffsets()
	fb.paths = nil
	return fb
}

func (up *fbupdate) save2(depth int, nodeOff uint64) uint64 {
	// only traverse modified paths, not entire (possibly huge) tree
	if up.canSkip(nodeOff) {
		return nodeOff
	}
	node := up.getMutableNode(nodeOff)
	if depth < up.fb.treeLevels {
		// tree node, need to update any memOffsets
		for it := node.iter(); it.next(); {
			off := it.offset
			off2 := up.save2(depth+1, off) // recurse
			// bottom up
			if off2 != off {
				node.setOffset(it.fi, off2)
			}
		}
	}
	off := up.fb.putNode(node)
	verify.That(up.canSkip(off))
	return off
}

func (up *fbupdate) canSkip(off uint64) bool {
	if off > up.moffs.nextOff {
		return false
	}
	if _, ok := up.moffs.nodes[off]; ok {
		return false
	}
	if _, ok := up.fb.moffs.nodes[off]; ok {
		return false
	}
	if _, ok := up.fb.paths[off]; ok {
		return false
	}
	return true
}

// putNode stores the node with a leading uint16 size
func (fb *fbtree) putNode(node fNode) uint64 {
	off, buf := fb.store.Alloc(2 + len(node))
	size := len(node)
	buf[0] = byte(size)
	buf[1] = byte(size >> 8)
	copy(buf[2:], node)
	return off
}

func (fb *fbtree) getNode(off uint64) fNode {
	if node := fb.moffs.get(off); node != nil {
		return node
	}
	buf := fb.store.Data(off)
	size := int(buf[0]) + int(buf[1])<<8
	verify.That(7 <= size && size <= fb.maxNodeSize)
	return fNode(buf[2 : 2+size])
}

//-------------------------------------------------------------------

// check verifies that the keys are in order and returns the number of keys
func (fb *fbtree) check() (count, size, nnodes int) {
	return fb.check1(0, fb.root, "")
}

func (fb *fbtree) check1(depth int, offset uint64, key string) (count, size, nnodes int) {
	node := fb.getNode(offset)
	size += len(node)
	nnodes++
	for it := node.iter(); it.next(); {
		offset := it.offset
		if depth < fb.treeLevels {
			// tree
			if it.fi > 0 && key > it.known {
				panic("keys out of order " + key + " " + it.known)
			}
			c, s, n := fb.check1(depth+1, offset, it.known) // recurse
			count += c
			size += s
			nnodes += n
		} else {
			// leaf
			itkey := fb.getLeafKey(offset)
			if key > itkey {
				panic("keys out of order " + key + " " + itkey)
			}
			count++
		}
	}
	return
}

// ------------------------------------------------------------------

type fbIter = func() (string, uint64, bool)

func (fb *fbtree) Iter() fbIter {
	var stack [maxlevels]*fnIter

	// traverse down the tree to the leftmost leaf, making a stack of iterators
	nodeOff := fb.root
	for i := 0; i < fb.treeLevels; i++ {
		stack[i] = fb.getNode(nodeOff).iter()
		stack[i].next()
		nodeOff = stack[i].offset
	}
	iter := fb.getNode(nodeOff).iter()

	return func() (string, uint64, bool) {
		for {
			if iter.next() {
				off := iter.offset
				return fb.getLeafKey(off), off, true // most common path
			}
			// end of leaf, go up the tree
			i := fb.treeLevels - 1
			for ; i >= 0; i-- {
				if stack[i].next() {
					nodeOff = stack[i].offset
					break
				}
			}
			if i == -1 {
				return "", 0, false // eof
			}
			// and then back down to the next leaf
			for i++; i < fb.treeLevels; i++ {
				stack[i] = fb.getNode(nodeOff).iter()
				stack[i].next()
				nodeOff = stack[i].offset
			}
			iter = fb.getNode(nodeOff).iter()
		}
	}
}

// ------------------------------------------------------------------

func (fb *fbtree) print() {
	fb.print1(0, fb.root)
}

func (fb *fbtree) print1(depth int, offset uint64) {
	node := fb.getNode(offset)
	for it := node.iter(); it.next(); {
		offset := it.offset
		if depth < fb.treeLevels {
			// tree
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known)
			fb.print1(depth+1, offset) // recurse
		} else {
			// leaf
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known, "("+fb.getLeafKey(offset)+")")
		}
	}
}
