// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

// mbtree is a specialized in-memory btree with a max size and number of levels.
// Nodes are fixed size to reduce allocation and bounds checks.
type mbtree struct {
	// tree is not embedded since it's not needed when small
	tree *mTree
	// leaf is embedded to reduce indirection and optimize when small
	leaf mLeaf
}

// mSize of 128 means capacity of 128 * 128 = 16k
// with splitting giving an average of 3/4 full, that gives average max of 12k
// which is comparable to jSuneido's 10,000 update limit
const mSize = 128

type mLeaf struct {
	size  int
	slots [mSize]mLeafSlot
}

type mLeafSlot struct {
	key string
	off uint64
}

type mTree struct {
	slots [mSize + 1]mTreeSlot
	size  int
}

type mTreeSlot struct {
	key  string
	leaf *mLeaf
}

func newMbtree() *mbtree {
	return &mbtree{}
}

func (m *mbtree) Insert(key string, off uint64) {
	if m.tree == nil && m.leaf.size == mSize {
		m.tree = &mTree{size: 1}
		m.tree.slots[0].leaf = &m.leaf
	}
	if m.tree != nil {
		m.tree.insert(key, off)
	} else {
		m.leaf.insert(m.tree, key, off)
	}
}

func (leaf *mLeaf) insert(tree *mTree, key string, off uint64) {
	if leaf.size < mSize {
		leaf.insert2(key, off)
	} else {
		leaf.split(tree, key)
		tree.insert(key, off)
	}
}

func (leaf *mLeaf) split(tree *mTree, key string) {
	var left int
	if key > leaf.slots[mSize-1].key {
		left = (mSize * 3) / 4
	} else if key < leaf.slots[0].key {
		left = mSize / 4
	} else {
		left = mSize / 2
	}
	right := mSize - left
	leaf2 := &mLeaf{size: right}
	copy(leaf2.slots[:], leaf.slots[left:])
	leaf.size = left
	tree.insert2(leaf2.slots[0].key, leaf2)
}

func (leaf *mLeaf) insert2(key string, off uint64) {
	i := leaf.searchBinary(key)
	// i is either leaf.size or points to first slot > key
	copy(leaf.slots[i+1:], leaf.slots[i:])
	leaf.slots[i].key, leaf.slots[i].off = key, off
	leaf.size++
}

func (tree *mTree) insert(key string, off uint64) {
	i := tree.searchBinary(key)
	tree.slots[i-1].leaf.insert(tree, key, off)
}

// insert2 inserts a key & leaf into the tree node
func (tree *mTree) insert2(key string, leaf *mLeaf) {
	i := tree.searchBinary(key)
	copy(tree.slots[i+1:], tree.slots[i:])
	tree.slots[i].key, tree.slots[i].leaf = key, leaf
	tree.size++
}

func (m *mbtree) Search(key string) uint64 {
	if m.tree != nil {
		return m.tree.search(key)
	}
	return m.leaf.search(key)
}

func (tree *mTree) search(key string) uint64 {
	i := tree.searchBinary(key)
	return tree.slots[i-1].leaf.search(key)
}

func (tree *mTree) searchBinary(key string) int {
	i, j := 0, tree.size
	for i < j {
		h := int(uint(i+j) >> 1)
		if key >= tree.slots[h].key {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

func (leaf *mLeaf) search(key string) uint64 {
	i := leaf.searchBinary(key)
	if i >= leaf.size || leaf.slots[i].key != key {
		return 0
	}
	return leaf.slots[i].off
}

func (leaf *mLeaf) searchBinary(key string) int {
	i, j := 0, leaf.size
	for i < j {
		h := int(uint(i+j) >> 1)
		if key > leaf.slots[h].key {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

//-------------------------------------------------------------------

type mIter func() (string, uint64, bool)

func (m *mbtree) Iterator() mIter {
	tree := m.tree
	ti := 0
	var leaf *mLeaf
	if tree == nil {
		leaf = &m.leaf
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
