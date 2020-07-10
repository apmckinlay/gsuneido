// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package hamt implements a hash array mapped trie.
// 		http://lampwww.epfl.ch/papers/idealhashtrees.pdf
// It is persistent in the functional sense.
// ItemHmat is immutable, read-only.
// ItemHmatUpdate groups updates to reduce path copying.
//
// The Item type must have:
// 		func (*item) Key() KeyType
// 		func ItemHash(key KeyType) uint32
//
// The KeyType must be comparable with ==
package hamt
