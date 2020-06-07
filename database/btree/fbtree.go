// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import "github.com/apmckinlay/gsuneido/database/stor"

// fbtree is a btree designed to be stored immutable in a file.
type fbtree struct {
	// treeLevels is how many levels of tree nodes there are.
	// Initially a tree consists of a single leaf node with treeLevels = 0
	treeLevels int
	// root is the offset of the root node
	root uint64
	// store is where the btree is stored
	store *stor.Stor
}

func (fb *fbtree) Search(key string) uint64 {
	nodeOff := fb.root
	for i := 0; i <= fb.treeLevels; i++ {
		node := fb.getNode(nodeOff)
		nodeOff, _, _ = node.search(key)
	}
	return nodeOff
}

// putNode stores the node with a leading uint16 size
func (fb *fbtree) putNode(node fNode) uint64 {
	off, buf := fb.store.Alloc(2 + len(node))
	size := len(node)
	buf[0] = byte(size)
	buf[1] = byte(size >> 8)
	copy(buf[2:], node)
	return off
}

func (fb *fbtree) getNode(off uint64) fNode {
	buf := fb.store.Data(off)
	size := int(buf[0]) + int(buf[1])<<8 //TODO validate
	return fNode(buf[2 : 2+size])
}
