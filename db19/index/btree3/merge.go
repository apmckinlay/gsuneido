// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

// leafMerge represents the current leaf node in the merge operation.
// limit will be "" on the right hand edge i.e. no limit.
// If modified is true, node is an in-memory copy that has been modified.
type leafMerge struct {
	limit    string // upper bound of this node ("" = no limit, rightmost)
	leaf     leafNode
	off      uint64 // storage offset of this node
	modified bool
}

// treeMerge represents one tree node on the current path.
// limit will be "" on the right hand edge i.e. no limit.
// If modified is true, node is an in-memory copy that has been modified.
type treeMerge struct {
	limit    string
	tree     treeNode
	off      uint64 // storage offset of this node
	pos      int    // child index
	modified bool
}

// state stores a path through the tree
// so dense updates don't have to do a full lookup for each
type state struct {
	bt   *btree
	tree []treeMerge // path through tree, tree[0] is root, len(tree) == treeLevels
	leaf leafMerge   // current leaf node
}

// t should be used like: _ = t && trace(...)
const t = false

func trace(args ...any) bool {
	fmt.Println(args...)
	return true
}

/*
MergeAndSave combines a btree and an ixbuf iter.
This is the primary way that btrees are updated.

It is immutable persistent, returning a new btree
that usually shares some of the structure of the original btree.
Modified nodes are written to storage.
It path copies.
Normally iter will be small relative to the btree.
The ixbuf iter contains inserts, updates, and deletes
which are identified by tag bits on the offsets.

We need to efficiently handle both sparse (far apart) and
dense (close together) updates.
For sparse updates we could just search for the location of each update.
But for dense updates (e.g. same leaf node) that would be inefficient.
So we maintain a search path down the tree and reuse it.
For an update in the current leaf node, we can just update the leaf.
For an update farther away, we go up the tree as needed
(not necessarily to the root) and then descend to the new leaf.

We also need to handle splitting full leaf nodes and updating their parent(s)
and removing empty nodes by updating their parent(s).
*/
func (bt *btree) MergeAndSave(iter ixbuf.Iter) *btree {
	bt2 := *bt // copy
	st := state{bt: &bt2}
	st.tree = make([]treeMerge, 0, maxLevels)
	for {
		key, off, ok := iter()
		if !ok {
			break
		}
		st.Print()
		_ = t && trace("MERGE", key, offstr(off))
		st.advanceTo(key)
		st.Print()
		st.check()
		st.updateLeaf(key, off)
		st.check()
	}
	_ = t && trace("flush")
	for st.haveLeaf() || len(st.tree) > 0 {
		st.ascend()
	}
	return st.bt
}

// advanceTo traverses the tree to the leaf containing key
// ascending and then descending as necessary.
func (st *state) advanceTo(key string) {
	_ = t && trace("advanceTo", key)
	bt := st.bt
	// if already on correct leaf, avoid tree traversal
	// this optimizes dense updates
	if len(st.tree) == bt.treeLevels && st.haveLeaf() && st.leaf.contains(key) {
		_ = t && trace("advance: already on correct node")
		return
	}
	// ascend the tree only as far as necessary
	for st.haveLast() && !st.lastContains(key) {
		st.ascend()
	}
	st.descendToLeaf(key)
}

// ascend goes up the tree
// by removing the last leaf or tree from state
// saving it if it was modified.
func (st *state) ascend() {
	_ = t && trace("ascend")
	var off uint64
	if st.haveLeaf() {
		// leaf
		lm := st.leaf
		st.leaf = leafMerge{}
		if !lm.modified {
			return
		}
		off = lm.leaf.write(st.bt.stor)
		_ = t && trace("write", lm.leaf, "=>", off)
	} else {
		// tree
		tm := &st.tree[len(st.tree)-1]
		st.tree = st.tree[:len(st.tree)-1] // pop
		if !tm.modified {
			return
		}
		off = tm.tree.write(st.bt.stor)
		_ = t && trace("write", tm.tree, "=>", off)
	}
	st.updateParent(off)
}

func (st *state) haveLast() bool {
	return len(st.tree) > 0 || st.haveLeaf()
}

func (st *state) lastContains(key string) bool {
	if st.haveLeaf() {
		return st.leaf.contains(key)
	}
	if len(st.tree) > 0 {
		return st.tree[len(st.tree)-1].contains(key)
	}
	return false
}

func (st *state) updateParent(off uint64) {
	if len(st.tree) == 0 {
		st.bt.root = off
		_ = t && trace("set root", off)
	} else {
		tm := &st.tree[len(st.tree)-1]
		tm.makeMutable()
		_ = t && trace("update tree set", tm.pos, "to", off)
		tm.tree.update(tm.pos, off)
		_ = t && trace("=>", tm.tree)
	}
}

// descendToLeaf descends to the leaf for key.
// It starts with a partial path after ascend
// and extends it fully to the leaf.
func (st *state) descendToLeaf(key string) {
	_ = t && trace("descendToLeaf", key)
	bt := st.bt
	for len(st.tree) < bt.treeLevels {
		var off uint64
		limit := ""
		if len(st.tree) == 0 {
			off = bt.root
		} else {
			tm := &st.tree[len(st.tree)-1]
			tm.pos, off = tm.tree.search(key)
			if tm.pos < tm.tree.nkeys() {
				limit = string(tm.tree.key(tm.pos))
			}
			if limit == "" {
				limit = tm.limit
			}
		}
		st.tree = append(st.tree, treeMerge{
			off: off, tree: bt.getTree(off), limit: limit, pos: -1,
		})
	}
	// get the leaf
	if bt.treeLevels == 0 {
		st.setLeaf(bt.root, bt.getLeaf(bt.root), "")
	} else {
		tm := &st.tree[len(st.tree)-1]
		pos, off := tm.tree.search(key)
		tm.pos = pos
		limit := ""
		if pos < tm.tree.nkeys() {
			limit = string(tm.tree.key(pos))
		}
		if limit == "" {
			limit = tm.limit
		}
		st.setLeaf(off, bt.getLeaf(off), limit)
	}
	// path now goes from root to the leaf containing key
	assert.That(st.lastContains(key))
}

func (st *state) haveLeaf() bool {
	return st.leaf.off != 0
}

//-------------------------------------------------------------------

func (st *state) updateLeaf(key string, off uint64) {
	_ = t && trace("updateLeaf", key, offstr(off))
	assert.That(st.leaf.limit == "" || key < st.leaf.limit)
	st.leaf.makeMutable()
	st.leaf.leaf = st.leaf.leaf.modify(key, off)
	_ = t && trace("=>", st.leaf.leaf)

	if st.leaf.leaf.nkeys() == 0 {
		st.dropLeaf()
		st.Print()
	} else if st.bt.shouldSplit(st.leaf.leaf) {
		st.split()
		st.Print()
	}
}

func (st *state) dropLeaf() {
	_ = t && trace("dropLeaf")
	assert.That(len(st.tree) == st.bt.treeLevels)
	st.leaf = leafMerge{} // clear leaf
	// propagate the delete up the tree
	// i.e. delete the leaf from its parent
	// then if it becomes empty, delete it from its parent, and so on
	// Note: we shorten tree, but treeLevels doesn't change
	// descendToLeaf (in updateLeaf) will refill the state
	// There are three ending cases:
	// - non-empty non-root tree node - we're done, return
	// - root becomes empty - btree = empty
	// - root has a single child - pull up the child
	for len(st.tree) > 0 {
		i := len(st.tree) - 1
		tm := &st.tree[i]
		_ = t && trace("before delete", tm.pos, tm.tree)
		tm.makeMutable()
		tm.tree = tm.tree.delete(tm.pos)
		_ = t && trace("delete", tm.pos, "=>", tm.tree)
		if tm.pos >= tm.tree.noffs() {
			tm.pos = tm.tree.noffs() - 1
		}
		if len(st.tree) == 1 {
			break
		}
		if tm.tree.noffs() > 0 {
			return
		}
		st.tree = st.tree[:i]
	}
	if st.bt.treeLevels == 0 || // dropped root leaf
		st.tree[0].tree.noffs() == 0 { // root tree node became empty
		st.setBtreeEmpty()
		return
	}
	// pop single child root(s)
	for st.bt.treeLevels > 0 {
		if len(st.tree) == 0 {
			st.tree = append(st.tree,
				treeMerge{off: st.bt.root, tree: st.bt.getTree(st.bt.root)})
		}
		tm := &st.tree[0] // root
		if tm.tree.noffs() > 1 {
			break
		}
		_ = t && trace("single child root")
		st.bt.root = tm.tree.offset(0)
		st.tree = slices.Delete(st.tree, 0, 1)
		st.bt.treeLevels--
	}
}

func (st *state) setBtreeEmpty() {
	_ = t && trace("set btree empty")
	st.bt.treeLevels = 0
	st.tree = st.tree[:0]
	rootLeaf := leafNode(emptyLeaf)
	st.bt.root = rootLeaf.write(st.bt.stor)
	st.setLeaf(st.bt.root, rootLeaf, "")
}

func (st *state) split() {
	leftOff, rightOff, splitKey := st.leaf.leaf.splitTo(st.bt.stor)
	_ = t && trace("write left leaf =>", leftOff)
	_ = t && trace("write right leaf =>", rightOff)
	st.leaf = leafMerge{} // clear leaf
	// propagate insert up the tree
	for len(st.tree) > 0 {
		i := len(st.tree) - 1
		tm := &st.tree[i]
		tm.makeMutable()
		_ = t && trace("before update", tm.tree)
		tm.tree.update(tm.pos, rightOff)
		tm.tree = tm.tree.insert(tm.pos, leftOff, splitKey)
		_ = t && trace("after update", tm.tree)
		if !st.bt.shouldSplit(tm.tree) {
			return
		}
		leftOff, rightOff, splitKey = tm.tree.splitTo(st.bt.stor)
		_ = t && trace("write left tree =>", leftOff)
		_ = t && trace("write right tree =>", rightOff)
		st.tree = st.tree[:i]
	}
	// new root
	var b treeBuilder
	b.add(leftOff, splitKey)
	newRoot := b.finish(rightOff)
	st.bt.treeLevels++
	st.bt.root = newRoot.write(st.bt.stor)
	_ = t && trace("write new root", newRoot, "=>", st.bt.root)
	st.tree = append(st.tree, treeMerge{off: st.bt.root, tree: newRoot})
}

//-------------------------------------------------------------------

func (st *state) setLeaf(off uint64, nd leafNode, limit string) {
	st.leaf = leafMerge{off: off, leaf: nd, limit: limit}
}

func (lm *leafMerge) contains(key string) bool {
	return lm.limit == "" || key < lm.limit
}

func (lm *leafMerge) makeMutable() {
	// copy the node when we start modifying
	if !lm.modified {
		size := len(lm.leaf)
		m := make(leafNode, size, size*6/5) // allow 20% for growth
		copy(m, lm.leaf)
		lm.leaf = m
		lm.modified = true
	}
}

func (tm *treeMerge) contains(key string) bool {
	return tm.limit == "" || key < tm.limit
}

func (tm *treeMerge) makeMutable() {
	// clone the node when we start modifying
	if !tm.modified {
		tm.tree = slc.Clone(tm.tree)
		tm.modified = true
	}
}

//-------------------------------------------------------------------

func offstr(off uint64) string {
	pre := ""
	if off&ixbuf.Delete != 0 {
		pre = "-"
	}
	if off&ixbuf.Update != 0 {
		pre += "="
	}
	return pre + strconv.Itoa(int(off&ixbuf.Mask))
}

func (lm leafMerge) String() string {
	limit := lm.limit
	if limit == "" {
		limit = `""`
	}
	mod := ""
	if lm.modified {
		mod = " modified"
	}
	return fmt.Sprint("leaf off: ", lm.off, " limit ", limit, mod)
}

func (tm treeMerge) String() string {
	limit := tm.limit
	if limit == "" {
		limit = `""`
	}
	mod := ""
	if tm.modified {
		mod = " modified"
	}
	return fmt.Sprint("tree off: ", tm.off, " pos: ", tm.pos, " limit: ", limit, mod)
}

func (st *state) Print() {
	// _ = t && trace("state root", st.bt.root, "treeLevels:", st.bt.treeLevels)
	// for i, tm := range st.tree {
	// 	_ = t && trace("   ", i, ":", &tm)
	// 	_ = t && trace("       ", tm.tree.String())
	// }
	// if st.haveLeaf() {
	// 	_ = t && trace("    leaf:", &st.leaf)
	// 	_ = t && trace("       ", st.leaf.leaf.String())
	// }
}

func (st *state) check() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		st.print()
	// 		panic(r)
	// 	}
	// }()
	// prevLimit := ""
	// off := st.bt.root
	// for _, tm := range st.tree {
	// 	assert.That(off == tm.off || tm.modified)
	// 	assert.That(0 <= tm.pos && tm.pos < tm.tree.noffs())
	// 	off = tm.tree.offset(tm.pos)
	// 	assert.This(tm.limit).Is(prevLimit)
	// 	if tm.pos < tm.tree.nkeys() {
	// 		prevLimit = string(tm.tree.key(tm.pos))
	// 	}
	// }
	// if st.haveLeaf() {
	// 	assert.That(off == st.leaf.off)
	// 	assert.This(st.leaf.limit).Is(prevLimit)
	// }
}
