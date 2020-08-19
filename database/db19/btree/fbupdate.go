// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func (fb *fbtree) Update(fn func(*fbtree)) *fbtree {
	mfb := fb.makeMutable()
	fn(mfb)
	return mfb.freeze()
}

func (fb *fbtree) makeMutable() *fbtree {
	mfb := *fb // copy
	mfb.mutable = true
	mfb.redirs.tbl = fb.redirs.tbl.Mutable()
	mfb.redirs.paths = fb.redirs.paths.Mutable()
	return &mfb
}

func (fb *fbtree) freeze() *fbtree {
	fb.redirs.generation++ // make new stuff immutable
	fb.redirs.tbl = fb.redirs.tbl.Freeze()
	fb.mutable = false
	return fb
}

//-------------------------------------------------------------------

const maxlevels = 8

func (fb *fbtree) Insert(key string, off uint64) {
	assert.That(fb.mutable)
	var stack [maxlevels]uint64

	// search down the tree to the appropriate leaf
	nodeOff := fb.root
	for i := 0; i < fb.treeLevels; i++ {
		stack[i] = nodeOff
		node := fb.getNode(nodeOff)
		fb.redirs.paths.Put(nodeOff)
		nodeOff, _, _ = node.search(key)
	}

	// insert into leaf
	node, where := fb.insert(nodeOff, key, off, fb.getLeafKey)
	fb.redirs.set(nodeOff, node)
	size := len(node)
	if size <= MaxNodeSize {
		return // fast path, just insert into leaf
	}

	// split leaf
	splitKey, rightOff := fb.split(node, nodeOff, where)
	// fmt.Println("split", splitKey, "old/left", OffStr(nodeOff), "new/right", OffStr(rightOff))

	// insert up the tree
	for i := fb.treeLevels - 1; i >= 0; i-- {
		// fmt.Println("tree insert", OffStr(stack[i]), "splitKey", splitKey, "rightOff", OffStr(rightOff))
		node, where = fb.insert(stack[i], splitKey, rightOff, nil)
		fb.redirs.set(stack[i], node)
		size := len(node)
		if size <= MaxNodeSize {
			return // finished
		}
		splitKey, rightOff = fb.split(node, stack[i], where)
		// fmt.Println("split", splitKey, "old/left", OffStr(stack[i]), "new/right", OffStr(rightOff))
	}

	// split all the way up, create new root
	newRoot := make(fNode, 0, 24)
	newRoot = fAppend(newRoot, uint64(fb.root), 0, "")
	newRoot = fAppend(newRoot, uint64(rightOff), 0, splitKey)
	fb.root = fb.redirs.add(newRoot)
	fb.treeLevels++
}

func (fb *fbtree) insert(nodeOff uint64, key string, off uint64,
	get func(uint64) string) (fNode, int) {
	node := fb.getMutableNode(nodeOff)
	return node.insert(key, off, get)
}

func (fb *fbtree) getMutableNode(off uint64) fNode {
	var roNode fNode
	if r, ok := fb.redirs.tbl.Get(off); ok {
		if r.newOffset != 0 {
			roNode = fb.readNode(r.newOffset)
		} else if r.generation == fb.redirs.generation {
			return r.mnode
		} else {
			roNode = r.mnode
		}
	} else {
		roNode = fb.readNode(off)
	}
	node := append(roNode[:0:0], roNode...)
	fb.redirs.tbl.Put(&redir{offset: off,
		mnode: node, generation: fb.redirs.generation})
	return node
}

func (fb *fbtree) split(node fNode, nodeOff uint64, where int) (
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
	fb.redirs.set(nodeOff, left)
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
	rightOff = fb.redirs.add(right)
	return
}

//-------------------------------------------------------------------

func (fb *fbtree) Delete(key string, off uint64) bool {
	assert.That(fb.mutable)
	var stack [maxlevels]uint64

	// search down the tree to the appropriate leaf
	nodeOff := fb.root
	for i := 0; i < fb.treeLevels; i++ {
		stack[i] = nodeOff
		node := fb.getNode(nodeOff)
		fb.redirs.paths.Put(nodeOff)
		nodeOff, _, _ = node.search(key)
	}

	// delete from leaf
	node, ok := fb.delete(nodeOff, off)
	if !ok {
		return false
	}
	if len(node) != 0 || fb.treeLevels == 0 {
		return true // usual fast path
	}

	// delete up the tree
	for i := fb.treeLevels - 1; i >= 0; i-- {
		node, ok = fb.delete(stack[i], nodeOff)
		if !ok {
			panic("leaf node not found in tree")
		}
		if (i > 0 || fb.treeLevels == 0) && len(node) != 0 {
			return true
		}
		nodeOff = stack[i]
	}

	// remove empty root(s)
	for fb.treeLevels > 0 && len(node) == 7 {
		fb.treeLevels--
		fb.root = stor.ReadSmallOffset(node)
		node = fb.getNode(fb.root)
	}

	return true
}

func (fb *fbtree) delete(nodeOff uint64, off uint64) (fNode, bool) {
	node := fb.getMutableNode(nodeOff)
	node, ok := node.delete(off)
	if ok {
		fb.redirs.set(nodeOff, node)
	}
	return node, ok
}

//-------------------------------------------------------------------

//go:generate genny -in ../../../genny/hamt/hamt.go -out redirhamt.go -pkg btree gen "Item=*redir KeyType=uint64"
//go:generate genny -in ../../../genny/hamt/hamt.go -out pathhamt.go -pkg btree gen "Item=path KeyType=uint64"

// redirs is use to redirect offsets to new nodes
// to reduce write amplification e.g. from path copying.
// It is used for updated versions of existing nodes, using their old offset,
// and for new nodes with fake offsets.
type redirs struct {
	tbl        RedirHamt
	paths      PathHamt
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
	// generation is used to determine mutability,
	// the current generation is mutable, previous generations are immutable.
	generation uint
}

func RedirKey(r *redir) uint64 {
	return r.offset
}

const phi64 = 11400714819323198485

func RedirHash(key uint64) uint32 {
	return uint32(key * phi64)
}

type path = uint64

func PathKey(n uint64) uint64 {
	return n
}

func PathHash(key uint64) uint32 {
	return uint32(key * phi64)
}

func newRedirs() redirs {
	return redirs{nextOff: stor.MaxSmallOffset, generation: 1}
}

// add returns a fake offset for a new node (from split)
func (re *redirs) add(node fNode) uint64 {
	off := re.nextOff
	re.nextOff--
	re.tbl.Put(&redir{offset: off, mnode: node, generation: re.generation})
	return off
}

// set updates the node for an offset
func (re *redirs) set(off uint64, node fNode) {
	re.tbl.Put(&redir{offset: off, mnode: node, generation: re.generation})

	r, ok := re.tbl.Get(off)
	assert.That(ok)
	assert.That(r.offset == off)
	assert.That(r.newOffset == 0)
	assert.That(len(node) == 0 || &r.mnode[0] == &node[0])
}

// isFake returns true for temporary in-memory offsets
func (re *redirs) isFake(off uint64) bool {
	return off > re.nextOff
}

func OffStr(off uint64) string {
	s := strconv.Itoa(int(off))
	if off > 0xffff000000 {
		s = "^" + s[len(s)-4:]
	}
	return s
}

func (re *redirs) String() string {
	s := "{"
	re.tbl.ForEach(func(r *redir) {
		s += OffStr(r.offset) + ": "
		if r.newOffset != 0 {
			s += strconv.Itoa(int(r.newOffset))
		} else {
			s += r.mnode.String()
		}
		s += " "
	})
	return strings.TrimSpace(s) + "}"
}

func (re *redirs) Len() int {
	n := 0
	re.tbl.ForEach(func(r *redir) { n++ })
	return n
}
