// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ranges

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Ranges is an ordered set of non-overlapping ranges of strings.
// The zero value is ready to use.
// It uses a specialized in-memory btree with a max size and number of levels.
// Nodes are fixed size to reduce allocation and bounds checks.
// ordset uses a variation of this code.
type Ranges struct {
	// tree is not embedded since it's not needed when small
	tree *treeNode
	// leaf is embedded to reduce indirection and optimize when small
	leaf leafNode
}

// nodeSize of 128 means capacity of 128 * 128 = 16k
// with splitting giving an average of 3/4 full, that gives average max of 12k
// which is comparable to jSuneido's 10,000 update limit
const nodeSize = 128

type leafNode struct {
	size  int
	slots [nodeSize]leafSlot
}

type leafSlot struct {
	from string
	to   string
}

type treeNode struct {
	slots [nodeSize + 1]treeSlot
	size  int
}

type treeSlot struct {
	val  string
	leaf *leafNode
}

//-------------------------------------------------------------------

func (rs *Ranges) Insert(from, to string) {
	var ti int
	leaf := &rs.leaf
	if rs.tree != nil {
		ti = rs.tree.searchBinary(from) - 1
		leaf = rs.tree.slots[ti].leaf
	}
	if leaf.size >= nodeSize {
		rs.split(leaf, from)
		ti = rs.tree.searchBinary(from) - 1
		leaf = rs.tree.slots[ti].leaf
	}

	li := leaf.insert(from, to)
	if li == -1 {
		return // from,to is already contained in an existing range
	}

	// coalesce / merge
	it := rs.iter(ti, leaf, li)
	it.prev()
	if it.eof() || it.cur().to < from {
		it = rs.iter(ti, leaf, li)
	}
	prev := it.cur()
	it.next()
	for !it.eof() {
		next := it.cur()
		if !merge(prev, next) {
			break
		}
		it.remove()
	}
}

type iter struct {
	tree *treeNode
	ti   int
	leaf *leafNode // == tree.slots[ti].leaf
	li   int
}

func (rs *Ranges) iter(ti int, leaf *leafNode, li int) *iter {
	assert.That(rs.tree == nil || leaf == rs.tree.slots[ti].leaf)
	return &iter{tree: rs.tree, ti: ti, leaf: leaf, li: li}
}

func (it *iter) cur() *leafSlot {
	return &it.leaf.slots[it.li]
}

func (it *iter) remove() {
	copy(it.leaf.slots[it.li:], it.leaf.slots[it.li+1:])
	it.leaf.size--
	if it.tree != nil {
		if it.leaf.size == 0 {
			// leaf empty, remove from tree
			copy(it.tree.slots[it.ti:], it.tree.slots[it.ti+1:])
			it.tree.size--
			if it.ti < it.tree.size {
				it.leaf = it.tree.slots[it.ti].leaf
				it.li = 0
			}
			return
		}
		if it.li == 0 {
			// first slot, update separator
			assert.That(it.leaf == it.tree.slots[it.ti].leaf)
			it.tree.slots[it.ti].val = it.leaf.slots[0].from
		}
	}
	it.next2()
}

func (it *iter) eof() bool {
	return it.li < 0 || it.leaf.size <= it.li
}

func (it *iter) next() {
	it.li++
	it.next2()
}

func (it *iter) next2() {
	if it.li < it.leaf.size || it.tree == nil {
		return
	}
	it.ti++
	if it.ti >= it.tree.size {
		return
	}
	it.leaf = it.tree.slots[it.ti].leaf
	it.li = 0
}

func (it *iter) prev() {
	it.li--
	if it.li >= 0 || it.tree == nil {
		return
	}
	it.ti--
	if it.ti < 0 {
		return
	}
	it.leaf = it.tree.slots[it.ti].leaf
	it.li = it.leaf.size - 1
}

func merge(ls1 *leafSlot, ls2 *leafSlot) bool {
	if overlap(ls1, ls2) {
		ls1.from = str.Min(ls1.from, ls2.from)
		ls1.to = str.Max(ls1.to, ls2.to)
		return true
	}
	return false
}

func overlap(ls1, ls2 *leafSlot) bool {
	return ls1.to >= ls2.from && ls2.to >= ls1.from
}

func (ls *leafSlot) contains(from, to string) bool {
	return ls.from <= from && to <= ls.to
}

//-------------------------------------------------------------------

func (leaf *leafNode) insert(from, to string) int {
	i := leaf.searchBinary(from)
	if leaf.slots[i].contains(from, to) ||
		(i > 0 && leaf.slots[i-1].contains(from, to)) {
		return -1
	}
	copy(leaf.slots[i+1:], leaf.slots[i:])
	leaf.slots[i] = leafSlot{from: from, to: to}
	leaf.size++
	return i
}

func (rs *Ranges) split(leaf *leafNode, val string) {
	if rs.tree == nil {
		rs.tree = &treeNode{size: 1}
		rs.tree.slots[0].leaf = &rs.leaf
	}
	var left int
	if val > leaf.slots[nodeSize-1].from {
		left = (nodeSize * 3) / 4
	} else if val < leaf.slots[0].to {
		left = nodeSize / 4
	} else {
		left = nodeSize / 2
	}
	right := nodeSize - left
	leaf2 := &leafNode{size: right}
	copy(leaf2.slots[:], leaf.slots[left:])
	leaf.size = left
	rs.tree.insert(leaf2.slots[0].from, leaf2)
}

func (tree *treeNode) insert(val string, leaf *leafNode) {
	i := tree.searchBinary(val)
	copy(tree.slots[i+1:], tree.slots[i:])
	tree.slots[i].val, tree.slots[i].leaf = val, leaf
	tree.size++
}

//-------------------------------------------------------------------

func (rs *Ranges) Contains(val string) bool {
	if rs == nil {
		return false
	}
	_, leaf, li := rs.search(val)
	if li < leaf.size && leaf.slots[li].from == val {
		return true
	}
	if li > 0 {
		li--
		ls := leaf.slots[li]
		return ls.from <= val && val <= ls.to
	}
	return false
}

func (rs *Ranges) search(val string) (int, *leafNode, int) {
	if rs.tree != nil {
		return rs.tree.search(val)
	}
	li := rs.leaf.searchBinary(val)
	return 0, &rs.leaf, li
}

func (tree *treeNode) search(val string) (int, *leafNode, int) {
	ti := tree.searchBinary(val) - 1
	leaf := tree.slots[ti].leaf
	li := leaf.searchBinary(val)
	return ti, leaf, li
}

func (tree *treeNode) searchBinary(val string) int {
	i, j := 0, tree.size
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if val >= tree.slots[h].val {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// searchBinary returns the index of the first item >= val, or leaf.size
func (leaf *leafNode) searchBinary(val string) int {
	i, j := 0, leaf.size
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if val > leaf.slots[h].from {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

//-------------------------------------------------------------------

func (rs *Ranges) String() string {
	if rs.tree == nil {
		return rs.leaf.String()
	}
	var b strings.Builder
	for i := 0; i < rs.tree.size; i++ {
		b.WriteString(rs.tree.slots[i].val + " =>\n")
		b.WriteString(rs.tree.slots[i].leaf.String())
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}

func (leaf *leafNode) String() string {
	var b strings.Builder
	for i := 0; i < leaf.size; i++ {
		ls := leaf.slots[i]
		b.WriteString(ls.from + "->" + ls.to + " ")
	}
	return strings.TrimSpace(b.String())
}

func (ls *leafSlot) String() string {
	return ls.from + "->" + ls.to
}
