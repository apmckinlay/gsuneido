// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// fbupdate handles updating an immutable btree.
// New and updated nodes are createing in memory.
type fbupdate struct {
	// fb is a copy of the fbtree (so we can update it)
	fb fbtree
	// moffs (mutable) maps offsets to mutable new and updated in-memory nodes
	moffs memOffsets
}

func (fb *fbtree) Update(fn func(*fbupdate)) *fbtree {
	up := newFbupdate(fb)
	fn(up)
	return up.freeze()
}

func newFbupdate(fb *fbtree) *fbupdate {
	moffs := memOffsets{
		nextOff:    fb.moffs.nextOff,
		generation: fb.moffs.generation,
		redirs:     fb.moffs.redirs.Mutable()}
	return &fbupdate{fb: *fb, moffs: moffs}
}

// freeze moves the changes to the fbtree.
// It will still reference in-memory new and updated nodes
// but they are no longer mutable.
func (up *fbupdate) freeze() *fbtree {
	up.moffs.generation++ // make new stuff immutable
	up.moffs.redirs = up.moffs.redirs.Freeze()
	up.fb.moffs = up.moffs
	return &up.fb
}

//-------------------------------------------------------------------

func (up *fbupdate) Search(key string) uint64 {
	off := up.fb.root
	for i := 0; i <= up.fb.treeLevels; i++ {
		node := up.getNode(off)
		off, _, _ = node.search(key)
	}
	return off
}

const maxlevels = 8

func (up *fbupdate) Insert(key string, off uint64) {
	var stack [maxlevels]uint64

	// search down the tree to the appropriate leaf
	nodeOff := up.fb.root
	for i := 0; i < up.fb.treeLevels; i++ {
		stack[i] = nodeOff
		node := up.getNode(nodeOff)
		nodeOff, _, _ = node.search(key)
	}

	// insert into leaf
	node, where := up.insert(nodeOff, key, off, up.fb.getLeafKey)
	up.moffs.set(nodeOff, node)
	size := len(node)
	if size <= MaxNodeSize {
		return // fast path, just insert into leaf
	}

	// split leaf
	splitKey, rightOff := up.split(node, nodeOff, where)

	// insert up the tree
	for i := up.fb.treeLevels - 1; i >= 0; i-- {
		node, where = up.insert(stack[i], splitKey, rightOff, nil)
		up.moffs.set(stack[i], node)
		size := len(node)
		if size <= MaxNodeSize {
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
	if r := up.moffs.redirs.Get(off); r != nil {
		if r.mnode != nil {
			return r.mnode
		}
		off = r.newOffset
	}
	return up.fb.readNode(off)
}

func (up *fbupdate) getMutableNode(off uint64) fNode {
	var roNode fNode
	if r := up.moffs.redirs.Get(off); r != nil {
		if r.newOffset != 0 {
			roNode = up.fb.readNode(r.newOffset)
		} else if r.generation == up.moffs.generation {
			return r.mnode
		} else {
			roNode = r.mnode
		}
	} else {
		roNode = up.fb.readNode(off)
	}
	node := append(roNode[:0:0], roNode...)
	up.moffs.redirs.Put(&redir{offset: off,
		mnode: node, generation: up.moffs.generation})
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
		node := up.getNode(nodeOff)
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

//go:generate genny -in ../../../genny/hamt/hamt.go -out redirs.go -pkg btree gen "Item=redir KeyType=uint64"

// memOffsets is use to map offsets to temporary, in-memory, mutable nodes.
// It is used for updated versions of existing nodes, using their old offset,
// and for new nodes with fake offsets.
type memOffsets struct {
	redirs     RedirHamt
	nextOff    uint64
	generation uint
}

// redir is a single redirection.
// Only one of mnode or newOffset is used at a time, the other should be nil/0.
// Save converts mnode to newOffset.
type redir struct {
	// offset is the "old" offset referenced by old immutable nodes
	offset uint64
	// mnode is an in-memory node, mutable until freeze, shadowing an immutable node.
	mnode fNode
	// newOffset is the new storage location for a node.
	newOffset uint64
	// path is the path through the btree to the node containing the old offset.
	// It is used for eventually flushing/flattening the redirs.
	path [4]uint8
	// generation is used to determine mutability,
	// the current generation is mutable, previous generations are immutable.
	generation uint
}

func (r *redir) Key() uint64 {
	return r.offset
}

func RedirHash(key uint64) uint32 {
	const phi64 = 11400714819323198485
	return uint32(key * phi64)
}

func nilMemOffsets() memOffsets {
	return memOffsets{nextOff: stor.MaxSmallOffset, generation: 1}
}

// add returns a fake offset for a new node (from split)
func (mo *memOffsets) add(node fNode) uint64 {
	off := mo.nextOff
	mo.nextOff--
	verify.That(mo.redirs.Get(off) == nil)
	mo.redirs.Put(&redir{offset: off, mnode: node, generation: mo.generation})
	return off
}

// set updates the node for an offset
func (mo *memOffsets) set(off uint64, node fNode) {
	mo.redirs.Put(&redir{offset: off, mnode: node, generation: mo.generation})

	r := mo.redirs.Get(off)
	verify.That(r.offset == off)
	verify.That(r.newOffset == 0)
	verify.That(len(node) == 0 || &r.mnode[0] == &node[0])
}

// isMem returns true for temporary in-memory offsets
func (mo *memOffsets) isMem(off uint64) bool {
	return off > mo.nextOff
}

// isMut returns true if mutable
func (mo *memOffsets) isMut(off uint64) bool {
	r := mo.redirs.Get(off)
	return r != nil && r.generation == mo.generation
}

func OffStr(off uint64) string {
	if off > 0xffff000000 {
		return strconv.Itoa(int(off - stor.MaxSmallOffset - 1))
	}
	return strconv.Itoa(int(off))
}

func (mo *memOffsets) String() string {
	s := "{"
	mo.redirs.ForEach(func(r *redir) {
		if r.offset > mo.nextOff {
			s += strconv.Itoa(int(mo.nextOff - r.offset))
		} else {
			s += strconv.Itoa(int(r.offset))
		}
		s += ": "
		if r.newOffset != 0 {
			s += strconv.Itoa(int(r.newOffset))
		} else {
			s += r.mnode.String()
		}
		s += " "
	})
	return strings.TrimSpace(s) + "}"
}

func (mo *memOffsets) Len() int {
	n := 0
	mo.redirs.ForEach(func(r *redir) { n++ })
	return n
}

// ------------------------------------------------------------------

func (up *fbupdate) print() {
	up.print1(0, up.fb.root)
}

func (up *fbupdate) print1(depth int, offset uint64) {
	print(strings.Repeat("    ", depth)+"offset", OffStr(offset), "=", offset)
	node := up.getNode(offset)
	for it := node.iter(); it.next(); {
		offset := it.offset
		if depth < up.fb.treeLevels {
			// tree
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known)
			up.print1(depth+1, offset) // recurse
		} else {
			// leaf
			print(strings.Repeat("    ", depth)+strconv.Itoa(it.fi)+":",
				OffStr(offset)+",", it.npre, it.diff, "=", it.known,
				"("+up.fb.getLeafKey(offset)+")")
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
