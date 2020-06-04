// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

// fbtree is a btree designed to be stored immutable in a file.
type fbtree struct {
	// treeLevels is how many levels of tree nodes there are.
	// Initially a tree consists of a single leaf node with treeLevels = 0
	treeLevels int
	// root is the offset of the root node
	root uint64
}
