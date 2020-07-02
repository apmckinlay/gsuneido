// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ordset

import "github.com/apmckinlay/gsuneido/util/str"

// Set is an ordered set of strings.
// It uses a specialized in-memory btree with a max size and number of levels.
// Nodes are fixed size to reduce allocation and bounds checks.
// ranges and mbtree use variations of this code.
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

func (set *Set) Insert(key string) {
	//TODO ignore duplicates
	leaf := &set.leaf
	if set.tree != nil {
		i := set.tree.searchBinary(key)
		leaf = set.tree.slots[i-1].leaf
	}
	if leaf.size >= nodeSize {
		set.split(leaf, key)
		i := set.tree.searchBinary(key)
		leaf = set.tree.slots[i-1].leaf
	}
	leaf.insert(key)
}

func (leaf *leafNode) insert(key string) {
	i := leaf.searchBinary(key)
	// i is either leaf.size or points to first slot > key
	copy(leaf.slots[i+1:], leaf.slots[i:])
	leaf.slots[i] = key
	leaf.size++
}

func (set *Set) split(leaf *leafNode, key string) {
	if set.tree == nil && set.leaf.size == nodeSize {
		set.tree = &treeNode{size: 1}
		set.tree.slots[0].leaf = &set.leaf
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
}

func (tree *treeNode) insert(key string, leaf *leafNode) {
	i := tree.searchBinary(key)
	copy(tree.slots[i+1:], tree.slots[i:])
	tree.slots[i].key, tree.slots[i].leaf = key, leaf
	tree.size++
}

//-------------------------------------------------------------------

func (set *Set) Delete(key string) bool {
	leaf, i := set.search(key)
	if leaf == nil || leaf.slots[i] != key {
		return false
	}
	copy(leaf.slots[i:], leaf.slots[i+1:])
	leaf.size--
	return true
}

//-------------------------------------------------------------------

func (set *Set) Contains(key string) bool {
	if set == nil {
		return false
	}
	leaf, _ := set.search(key)
	return leaf != nil
}

func (set *Set) search(key string) (*leafNode, int) {
	if set.tree != nil {
		return set.tree.search(key)
	}
	return set.leaf.search(key)
}

func (tree *treeNode) search(key string) (*leafNode, int) {
	i := tree.searchBinary(key)
	return tree.slots[i-1].leaf.search(key)
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

func (leaf *leafNode) search(key string) (*leafNode, int) {
	i := leaf.searchBinary(key)
	if i >= leaf.size || leaf.slots[i] != key {
		return nil, 0
	}
	return leaf, i
}

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

func (set *Set) String() string {
	if set.tree != nil {
		return "Set too big"
	}
	var b str.CommaBuilder
	for i := 0; i < set.leaf.size; i++ {
		b.Add(set.leaf.slots[i])
	}
	return b.String()
}
