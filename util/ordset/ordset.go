// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ordset

import "github.com/apmckinlay/gsuneido/util/str"

// Set is an ordered set of strings.
// It uses a specialized in-memory btree with a max size and number of levels.
// Nodes are fixed size to reduce allocation and bounds checks.
// ranges uses a variation of this code.
type Set struct {
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
	slots [nodeSize]string
}

type treeNode struct {
	slots [nodeSize + 1]treeSlot
	size  int
}

type treeSlot struct {
	key  string
	leaf *leafNode
}

//-------------------------------------------------------------------

func (set *Set) Insert(key string) bool {
	//TODO ignore duplicates
	leaf := &set.leaf
	if set.tree != nil {
		i := set.tree.searchBinary(key)
		leaf = set.tree.slots[i-1].leaf
	}
	if leaf.size >= nodeSize {
		if !set.split(leaf, key) {
			return false
		}
		i := set.tree.searchBinary(key)
		leaf = set.tree.slots[i-1].leaf
	}
	return leaf.insert(key)
}

func (leaf *leafNode) insert(key string) bool {
	if leaf.size >= nodeSize {
		return false
	}
	i := leaf.searchBinary(key)
	copy(leaf.slots[i+1:], leaf.slots[i:])
	leaf.slots[i] = key
	leaf.size++
	return true
}

func (set *Set) split(leaf *leafNode, key string) bool {
	if set.tree == nil && set.leaf.size == nodeSize {
		set.tree = &treeNode{size: 1}
		set.tree.slots[0].leaf = &set.leaf
	} else if set.tree.size >= nodeSize {
		return false
	}
	var left int
	if key > leaf.slots[nodeSize-1] {
		left = (nodeSize * 3) / 4
	} else if key < leaf.slots[0] {
		left = nodeSize / 4
	} else {
		left = nodeSize / 2
	}
	right := nodeSize - left
	leaf2 := &leafNode{size: right}
	copy(leaf2.slots[:], leaf.slots[left:])
	leaf.size = left
	set.tree.insert(leaf2.slots[0], leaf2)
	return true
}

func (tree *treeNode) insert(key string, leaf *leafNode) {
	i := tree.searchBinary(key)
	copy(tree.slots[i+1:], tree.slots[i:])
	tree.slots[i].key, tree.slots[i].leaf = key, leaf
	tree.size++
}

//-------------------------------------------------------------------

func (set *Set) AnyInRange(from, to string) bool {
	if set == nil {
		return false
	}
	ti, leaf, li := set.search(from)
	if li >= leaf.size {
		if set.tree == nil || ti >= set.tree.size {
			return false
		}
		// advance to next leaf
		leaf = set.tree.slots[ti+1].leaf
		li = 0
	}
	return leaf.slots[li] <= to
}

//-------------------------------------------------------------------

func (set *Set) Contains(key string) bool {
	if set == nil {
		return false
	}
	_, leaf, li := set.search(key)
	return li < leaf.size && leaf.slots[li] == key
}

func (set *Set) search(key string) (int, *leafNode, int) {
	if set.tree != nil {
		return set.tree.search(key)
	}
	li := set.leaf.searchBinary(key)
	return 0, &set.leaf, li
}

func (tree *treeNode) search(key string) (int, *leafNode, int) {
	ti := tree.searchBinary(key) - 1
	leaf := tree.slots[ti].leaf
	li := leaf.searchBinary(key)
	return ti, leaf, li
}

func (tree *treeNode) searchBinary(key string) int {
	i, j := 0, tree.size
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if key >= tree.slots[h].key {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// searchBinary returns the index of the first item >= val, or leaf.size
func (leaf *leafNode) searchBinary(key string) int {
	i, j := 0, leaf.size
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if key > leaf.slots[h] {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

//-------------------------------------------------------------------

func (set *Set) Empty() bool {
	return set.tree == nil && set.leaf.size == 0
}

//-------------------------------------------------------------------

func (set *Set) String() string {
	if set.tree != nil {
		return "Set too big"
	}
	return str.Join(",", set.leaf.slots[:set.leaf.size])
}
