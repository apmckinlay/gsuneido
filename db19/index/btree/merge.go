// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// merge is one node on the current path.
// limit will be "" on the right hand edge i.e. no limit.
// If modified is true, node is an in-memory copy that has been modified.
type merge struct {
	off      uint64
	node     node
	pos      int
	limit    string
	modified bool
}

type state struct {
	bt   *btree
	path []merge
}

// MergeAndSave combines an btree and an iter.
// It is immutable persistent, returning a new btree
// that usually shares some of the structure of the original btree.
// Modified nodes are written to storage.
// It path copies.
// Normally iter will be small relative to the btree.
func (bt *btree) MergeAndSave(iter ixbuf.Iter) *btree {
	bt2 := *bt // copy
	st := state{bt: &bt2}
	for {
		key, off, ok := iter()
		if !ok {
			break
		}
		_ = t && trace("merge", key, offstr(off))
		st.advanceTo(key)
		st.updateLeaf(key, off)
	}
	for len(st.path) > 0 {
		st.ascend()
	}
	return st.bt
}

func offstr(off uint64) string {
	pre := ""
	if off&ixbuf.Delete != 0 {
		pre = "-"
	}
	if off&ixbuf.Update != 0 {
		pre += "="
	}
	return pre + strconv.Itoa(int(off&0xffffffffff))
}

// advanceTo traverses the tree to the leaf for key
// ascending and then descending as necessary.
func (st *state) advanceTo(key string) {
	bt := st.bt
	if len(st.path) == 0 {
		// first time
		st.push(bt.root, bt.getNode(bt.root), "")
	} else {
		// if on the right node, just return
		if len(st.path) == bt.treeLevels+1 &&
			(bt.treeLevels == 0 || st.last().contains(key)) {
			_ = t && trace("advance: already on correct node")
			m := st.last()
			m.pos, _, _ = m.node.search2(key) //TODO continue from last
			return
		}
		// ascend tree as necessary
		for len(st.path) > 1 && !st.last().contains(key) {
			_ = t && trace("advance: ascend")
			st.ascend()
		}
	}
	// descend to appropriate leaf
	for len(st.path) <= bt.treeLevels {
		_ = t && trace("advance: descend")
		m := st.last()
		pos, off, limit := m.node.search2(key) //TODO continue from last
		m.pos = pos
		nd := bt.getNode(off)
		if limit == "" {
			limit = m.limit
		}
		st.push(off, nd, limit)
	}
	// path now goes from root to leaf
	assert.That(len(st.path) == bt.treeLevels+1)
}

// ascend moves up the tree, normally one level.
// It saves the node at the current level and updates its parent.
// It splits if the node is too large.
// Splitting may propagate up the tree, possibly making a new root.
// It removes empty nodes, which may also propagate up the tree.
func (st *state) ascend() {
	// _ = T && trace("ascend")
	m := st.last()
	st.pop()
	if !m.modified {
		return
	}
	bt := st.bt

	if len(m.node) == 0 {
		// empty node
		if len(st.path) == 0 {
			bt.treeLevels = 0
		} else {
			// delete empty non-root node
			parent := st.last()
			parent.getMutableNode()
			nd, ok := parent.node.delete(m.off)
			assert.That(ok)
			parent.node = nd
			if len(st.path) > 1 {
				st.ascend() // tail recurse
			}
			return
		}
	}

	insertOff := uint64(0)
	insertKey := ""
	if len(m.node) > MaxNodeSize {
		_ = t && trace("split", m)
		left, right, splitKey := m.split()
		m.node = left
		insertKey = splitKey
		insertOff = right.putNode(bt.stor)
	}
	off := m.node.putNode(bt.stor)
	if len(st.path) > 0 {
		parent := st.last()
		parent.getMutableNode()
		_ = t && trace("update", m, "set", m.pos, off)
		parent.node.setOffset(parent.pos, off)
		if insertOff != 0 {
			get := bt.getLeafKey
			if len(st.path) <= bt.treeLevels {
				get = nil
			}
			assert.That(parent.contains(insertKey))
			parent.updateNode(insertKey, insertOff, get)
			if len(parent.node) >= (MaxNodeSize*3)/2 {
				// if it gets too big, leave the node so it will be split
				_ = t && trace("split - ascend")
				st.ascend() // tail recurse
			}
		}
	} else {
		if insertOff != 0 {
			// split root, create new root
			_ = t && trace("new root")
			newRoot := make(node, 0, 24)
			newRoot = newRoot.append(uint64(off), 0, "")
			newRoot = newRoot.append(uint64(insertOff), 0, insertKey)
			off = newRoot.putNode(bt.stor)
			st.push(off, newRoot, "")
			bt.treeLevels++
		}
		bt.root = off
	}
}

func (m *merge) contains(key string) bool {
	return m.limit == "" || key < m.limit
}

func (nd node) search2(key string) (pos int, off uint64, known string) {
	it := nd.iter()
	for it.next() && key >= string(it.known) {
		pos = it.pos
		off = it.offset
	}
	return pos, off, string(it.known)
}

func (st *state) last() *merge {
	return &st.path[len(st.path)-1]
}

func (st *state) pop() {
	st.path = st.path[:len(st.path)-1]
}

func (st *state) push(off uint64, nd node, limit string) {
	st.path = append(st.path,
		merge{off: off, node: nd, pos: -1, limit: limit})
}

func (st *state) updateLeaf(key string, off uint64) {
	m := st.last()
	if m.limit != "" {
		assert.Msg("key > limit").That(key < m.limit)
	}
	m.updateNode(key, off, st.bt.getLeafKey)
	if len(m.node) >= (MaxNodeSize*3)/2 {
		// if it gets too big, leave the node so it will be split
		_ = t && trace("overflow - ascend")
		st.ascend()
	}
}

func (m *merge) updateNode(key string, off uint64, get func(uint64) string) {
	nd := m.getMutableNode()
	m.node = nd.update(key, off, get) // handles updates and deletes
	_ = t && trace("after update", m.node.knowns())
}

func (m *merge) getMutableNode() node {
	if !m.modified {
		nd := make(node, len(m.node))
		copy(nd, m.node)
		m.node = nd
		m.modified = true
	}
	return m.node
}

func (m *merge) split() (left, right node, splitKey string) {
	nd := m.getMutableNode()
	size := len(nd)
	splitSize := size / 2
	it := nd.iter()
	for it.next() && it.pos < splitSize {
	}
	splitKey = string(it.known)

	left = nd[:it.pos]

	right = make(node, 0, len(nd)-it.pos+8)
	// first entry becomes 0, ""
	right = right.append(it.offset, 0, "")
	if it.next() {
		// second entry becomes 0, known
		right = right.append(it.offset, 0, string(it.known))
		if it.next() {
			right = append(right, nd[it.pos:]...)
		}
	}
	if t {
		trace("split at", splitKey)
		trace("    ", nd.knowns())
		trace("    left:", left.knowns())
		trace("    right:", right.knowns())
	}
	return
}
