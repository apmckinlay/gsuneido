// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

// mbtree is a specialized in-memory btree with a max size and number of levels.
// Nodes are fixed size to reduce allocation and bounds checks.
type mbtree struct {
	// tree is not embedded since it's not needed when small
	tree *mbTreeNode
	// leaf is embedded to reduce indirection and optimize when small
	leaf mbLeaf
}

// mSize of 128 means capacity of 128 * 128 = 16k
// with splitting giving an average of 3/4 full, that gives average max of 12k
// which is comparable to jSuneido's 10,000 update limit
const mSize = 128

type mbLeaf struct {
	size  int
	slots [mSize]mbLeafSlot
}

type mbLeafSlot struct {
	key string
	off uint64
}

type mbTreeNode struct {
	slots [mSize + 1]mbTreeSlot
	size  int
}

type mbTreeSlot struct {
	key  string
	leaf *mbLeaf
}

func newMbtree() *mbtree {
	return &mbtree{}
}

//-------------------------------------------------------------------

func (mb *mbtree) Insert(key string, off uint64) {
	if mb.tree == nil && mb.leaf.size == mSize {
		mb.tree = &mbTreeNode{size: 1}
		mb.tree.slots[0].leaf = &mb.leaf
	}
	if mb.tree != nil {
		mb.tree.insert(key, off)
	} else {
		mb.leaf.insert(mb.tree, key, off)
	}
}

func (leaf *mbLeaf) insert(tree *mbTreeNode, key string, off uint64) {
	if leaf.size < mSize {
		leaf.insert2(key, off)
	} else {
		leaf.split(tree, key)
		tree.insert(key, off)
	}
}

func (leaf *mbLeaf) split(tree *mbTreeNode, key string) {
	var left int
	if key > leaf.slots[mSize-1].key {
		left = (mSize * 3) / 4
	} else if key < leaf.slots[0].key {
		left = mSize / 4
	} else {
		left = mSize / 2
	}
	right := mSize - left
	leaf2 := &mbLeaf{size: right}
	copy(leaf2.slots[:], leaf.slots[left:])
	leaf.size = left
	tree.insert2(leaf2.slots[0].key, leaf2)
}

func (leaf *mbLeaf) insert2(key string, off uint64) {
	i := leaf.searchBinary(key)
	// i is either leaf.size or points to first slot > key
	copy(leaf.slots[i+1:], leaf.slots[i:])
	leaf.slots[i].key, leaf.slots[i].off = key, off
	leaf.size++
}

func (tree *mbTreeNode) insert(key string, off uint64) {
	i := tree.searchBinary(key)
	tree.slots[i-1].leaf.insert(tree, key, off)
}

// insert2 inserts a key & leaf into the tree node
func (tree *mbTreeNode) insert2(key string, leaf *mbLeaf) {
	i := tree.searchBinary(key)
	copy(tree.slots[i+1:], tree.slots[i:])
	tree.slots[i].key, tree.slots[i].leaf = key, leaf
	tree.size++
}

//-------------------------------------------------------------------

func (mb *mbtree) Delete(key string, off uint64) bool {
	leaf, i := mb.search(key)
	if leaf == nil || leaf.slots[i].off != off {
		return false
	}
	copy(leaf.slots[i:], leaf.slots[i+1:])
	leaf.size--
	return true
}

//-------------------------------------------------------------------

func (mb *mbtree) Search(key string) uint64 {
	leaf, i := mb.search(key)
	if leaf == nil {
		return 0
	}
	return leaf.slots[i].off
}

func (mb *mbtree) search(key string) (*mbLeaf, int) {
	if mb.tree != nil {
		return mb.tree.search(key)
	}
	return mb.leaf.search(key)
}

func (tree *mbTreeNode) search(key string) (*mbLeaf, int) {
	i := tree.searchBinary(key)
	return tree.slots[i-1].leaf.search(key)
}

func (tree *mbTreeNode) searchBinary(key string) int {
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

func (leaf *mbLeaf) search(key string) (*mbLeaf, int) {
	i := leaf.searchBinary(key)
	if i >= leaf.size || leaf.slots[i].key != key {
		return nil, 0
	}
	return leaf, i
}

func (leaf *mbLeaf) searchBinary(key string) int {
	i, j := 0, leaf.size
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if key > leaf.slots[h].key {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

//-------------------------------------------------------------------

type visitor func(key string, off uint64)

func (mb *mbtree) ForEach(fn visitor) {
	if mb.tree == nil {
		mb.leaf.forEach(fn)
	} else {
		for i := 0; i < mb.tree.size; i++ {
			mb.tree.slots[i].leaf.forEach(fn)
		}
	}
}

func (leaf *mbLeaf) forEach(fn visitor) {
	for i := 0; i < leaf.size; i++ {
		fn(leaf.slots[i].key, leaf.slots[i].off)
	}
}

//-------------------------------------------------------------------

type mbIter = func() (string, uint64, bool)

func (mb *mbtree) Iter() mbIter {
	tree := mb.tree
	ti := 0
	var leaf *mbLeaf
	if tree == nil {
		leaf = &mb.leaf
	} else {
		leaf = tree.slots[ti].leaf
	}
	i := -1
	return func() (string, uint64, bool) {
		i++
		if i >= leaf.size {
			if tree == nil || ti+1 >= tree.size {
				return "", 0, false
			}
			ti++
			leaf = tree.slots[ti].leaf
			i = 0
		}
		slot := leaf.slots[i]
		return slot.key, slot.off, true
	}
}
