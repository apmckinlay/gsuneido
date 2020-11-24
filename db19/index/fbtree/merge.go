// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

import (
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/util/assert"
)

//TODO handle delete tombstones

// merge is one node on the current path.
// limit will be "" on the right hand edge i.e. no limit.
// If modified is true, node is an in-memory copy that has been modified.
type merge struct {
	off      uint64
	node     fnode
	fi       int
	limit    string
	modified bool
}

type state struct {
	fb   *fbtree
	path []merge
}

// Merge combines an fbtree and an iter.
// It is immutable persistent, returning a new fbtree
// that usually shares some of the structure of the original fbtree.
// Modified nodes are written to storage.
// It path copies.
func (fb *fbtree) MergeAndSave(iter ixbuf.Iter) *fbtree {
	fb2 := *fb // copy
	st := state{fb: &fb2}
	for {
		key, off, ok := iter()
		if !ok {
			break
		}
		st.advanceTo(key)
		st.insertInLeaf(key, off)
	}
	for len(st.path) > 0 {
		st.ascend()
	}
	return st.fb
}

// advanceTo traverses the tree to the leaf for key
// ascending and then descending as necessary.
func (st *state) advanceTo(key string) {
	fb := st.fb
	if len(st.path) == 0 {
		// first time
		st.push(fb.root, fb.getNode(fb.root), "")
	} else {
		// if on the right node, just return
		if len(st.path) == fb.treeLevels+1 &&
			(fb.treeLevels == 0 || st.last().contains(key)) {
			_ = t && trace("advance: already on correct node")
			m := st.last()
			m.fi, _, _ = m.node.search2(key) //TODO continue from last
			return
		}
		// ascend tree as necessary
		for len(st.path) > 1 && !st.last().contains(key) {
			_ = t && trace("advance: ascend")
			st.ascend()
		}
	}
	// descend to appropriate leaf
	for len(st.path) <= fb.treeLevels {
		_ = t && trace("advance: descend")
		m := st.last()
		fi, off, limit := m.node.search2(key) //TODO continue from last
		m.fi = fi
		node := fb.getNode(off)
		if limit == "" {
			limit = m.limit
		}
		st.push(off, node, limit)
	}
	// path now goes from root to leaf
	assert.That(len(st.path) == fb.treeLevels+1)
}

// ascend moves up the tree one level.
// It saves the node at the current level and updates its parent.
// It splits if the node is too large.
// Splitting may propagate up the tree, including making a new root.
func (st *state) ascend() {
	// _ = T && trace("ascend")
	m := st.last()
	st.pop()
	if !m.modified {
		return
	}
	fb := st.fb
	insertOff := uint64(0)
	insertKey := ""
	if len(m.node) > MaxNodeSize {
		_ = t && trace("split", m)
		left, right, splitKey := m.split()
		m.node = left
		insertKey = splitKey
		insertOff = right.putNode(fb.store)
	}
	off := m.node.putNode(fb.store)
	if len(st.path) > 0 {
		parent := st.last()
		parent.getMutableNode()
		_ = t && trace("update", m, "set", m.fi, off)
		parent.node.setOffset(parent.fi, off)
		if insertOff != 0 {
			get := fb.getLeafKey
			if len(st.path) <= fb.treeLevels {
				get = nil
			}
			assert.That(parent.contains(insertKey))
			parent.insertInNode(insertKey, insertOff, get)
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
			newRoot := make(fnode, 0, 24)
			newRoot = newRoot.append(uint64(off), 0, "")
			newRoot = newRoot.append(uint64(insertOff), 0, insertKey)
			off = newRoot.putNode(fb.store)
			st.push(off, newRoot, "")
			fb.treeLevels++
		}
		fb.root = off
	}
}

func (m *merge) contains(key string) bool {
	return m.limit == "" || key < m.limit
}

func (fn fnode) search2(key string) (fi int, off uint64, known string) {
	it := fn.iter()
	for it.next() && key >= string(it.known) {
		fi = it.fi
		off = it.offset
	}
	return fi, off, string(it.known)
}

func (st *state) last() *merge {
	return &st.path[len(st.path)-1]
}

func (st *state) pop() {
	st.path = st.path[:len(st.path)-1]
}

func (st *state) push(off uint64, node fnode, limit string) {
	st.path = append(st.path,
		merge{off: off, node: node, fi: -1, limit: limit})
}

func (st *state) insertInLeaf(key string, off uint64) {
	m := st.last()
	if m.limit != "" {
		assert.Msg("key > limit").That(key < m.limit)
	}
	m.insertInNode(key, off, st.fb.getLeafKey)
	if len(m.node) >= (MaxNodeSize*3)/2 {
		// if it gets too big, leave the node so it will be split
		_ = t && trace("overflow - ascend")
		st.ascend()
	}
}

func (m *merge) insertInNode(key string, off uint64, get func(uint64) string) {
	node := m.getMutableNode()
	m.node = node.insert(key, off, get)
	_ = t && trace("after insert", m.node.knowns())
}

func (m *merge) getMutableNode() fnode {
	if !m.modified {
		node := make(fnode, len(m.node))
		copy(node, m.node)
		m.node = node
		m.modified = true
	}
	return m.node
}

func (m *merge) split() (left, right fnode, splitKey string) {
	node := m.getMutableNode()
	size := len(node)
	splitSize := size / 2
	it := node.iter()
	for it.next() && it.fi < splitSize {
	}
	splitKey = string(it.known)

	left = node[:it.fi]

	right = make(fnode, 0, len(node)-it.fi+8)
	// first entry becomes 0, ""
	right = right.append(it.offset, 0, "")
	if it.next() {
		// second entry becomes 0, known
		right = right.append(it.offset, 0, string(it.known))
		if it.next() {
			right = append(right, node[it.fi:]...)
		}
	}
	if t {
		trace("split at", splitKey)
		trace("    ", node.knowns())
		trace("    left:", left.knowns())
		trace("    right:", right.knowns())
	}
	return
}
