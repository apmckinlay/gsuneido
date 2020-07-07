// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
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
	// ixspec is an opaque value passed to GetLeafKey
	// normally it specifies which fields make up the key, based on the schema
	ixspec interface{}
}

// MaxNodeSize is the maximum node size in bytes, split if larger.
// Overridden by tests.
var MaxNodeSize = 1536 // * .75 ~ 1k

// GetLeafKey returns the key for a data offset. (e.g. extract key from record)
// It is a dependency that must be injected
var GetLeafKey func(st *stor.Stor, ixspec interface{}, off uint64) string

func CreateFbtree(st *stor.Stor) *fbtree {
	mo := memOffsets{nextOff: stor.MaxSmallOffset, nodes: map[uint64]fNode{}}
	root := mo.add(fNode{})
	return &fbtree{root: root, moffs: mo, store: st}
}

func OpenFbtree(st *stor.Stor, root uint64, treeLevels int) *fbtree {
	mo := nilMemOffsets()
	return &fbtree{root: root, treeLevels: treeLevels, moffs: mo, store: st}
}

func (fb *fbtree) getLeafKey(off uint64) string {
	return GetLeafKey(fb.store, fb.ixspec, off)
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
	off := node.putNode(up.fb.store)
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
func (node fNode) putNode(store *stor.Stor) uint64 {
	off, buf := store.AllocSized(len(node))
	copy(buf, node)
	return off
}

func (fb *fbtree) getNode(off uint64) fNode {
	if node := fb.moffs.get(off); node != nil {
		return node
	}
	buf := fb.store.DataSized(off)
	return fNode(buf)
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

// Iter returns a function that can be called to return consecutive entries.
// NOTE: The returned key is only the known prefix.
// (unlike mbtree which returns the actual key)
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
				return iter.known, iter.offset, true // most common path
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
				it.npre, it.diff, "=", it.known,
				"("+fb.getLeafKey(offset)+")")
		}
	}
}

// ------------------------------------------------------------------

// fbtreeBuilder is used to bulk load an fbtree.
// Keys must be added in order.
// The fbtree is built bottom up with no splitting or inserting.
// All nodes will be "full" except for the right hand edge.
type fbtreeBuilder struct {
	levels []*level // leaf is [0]
	prev   string
	store  *stor.Stor
}

type level struct {
	first   string
	builder fNodeBuilder
}

func newFbtreeBuilder(store *stor.Stor) *fbtreeBuilder {
	return &fbtreeBuilder{store: store, levels: []*level{{}}}
}

func (fb *fbtreeBuilder) Add(key string, off uint64) {
	if key <= fb.prev {
		panic("fbtreeBuilder keys must be inserted in order, without duplicates")
	}
	fb.insert(0, key, off)
	fb.prev = key
}

var nput = 0

func (fb *fbtreeBuilder) insert(li int, key string, off uint64) {
	if li >= len(fb.levels) {
		fb.levels = append(fb.levels, &level{})
	}
	lev := fb.levels[li]
	if len(lev.builder.fe) > MaxNodeSize {
		// flush full node to stor
		nput++
		offNode := lev.builder.fe.putNode(fb.store)
		fb.insert(li+1, lev.first, offNode) // recurse
		*lev = level{}
	}
	if len(lev.builder.fe) == 0 {
		lev.first = key
	}
	lev.builder.Add(key, off)
}

func (fb *fbtreeBuilder) Finish() (off uint64, treeLevels int) {
	var key string
	for li := 0; li < len(fb.levels); li++ {
		if li > 0 {
			// allow node to slightly exceed max size
			fb.levels[li].builder.Add(key, off)
		}
		key = fb.levels[li].first
		off = fb.levels[li].builder.fe.putNode(fb.store)
	}
	return off, len(fb.levels) - 1
}
