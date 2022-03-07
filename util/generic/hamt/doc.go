// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package hamt implements a hash array mapped trie.
// 		http://lampwww.epfl.ch/papers/idealhashtrees.pdf
// It is persistent in the functional sense.
//
// Because Go doesn't have unions, and interfaces are twice as big as pointers,
// we use separate bitmaps and arrays (slices) for values and pointers.
// This has the side benefit of allowing both a value and a pointer in a node.
// e.g. in the case of a collision, only one value must be pushed down
// to a new child node instead of both the new and the old.
//
// Unlike a conventional hash map, it just stores items.
// This works well when items already contain the key, or for a set.
// To make a map, the item type should be a struct with the key and value.
//
// It has two modes, mutable and immutable.
// This allows "batching" updates to share path copying.
//
// - Mutable returns an updateable copy.
//
// - Freeze returns an immutable copy.
//
// Put and Delete can only be used when mutable.
// When mutable it is NOT thread safe, it should be thread contained.
//
// If items are large, the element (E) type should probably be a pointer.
// However, to maintain immutability, items should not be modified via pointer.
// The code is not written to use *E
// because that's not what you want for e.g. string or int
package hamt
