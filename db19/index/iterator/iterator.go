// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package iterator

// T is the interface for a Suneido style iterator
type T interface {
	// Eof returns true if the index is empty,
	// Next hit the end, or Prev hit the beginning
	Eof() bool

	// Modified returns true if the index has been modified.
	// Seek resets modified.
	Modified() bool

	// Cur returns the current key & offset
	// as of the most recent Next, Prev, or Seek
	Cur() (key string, off uint64)

	// Next advances to the first key > cur
	Next()

	// Prev moves backwards to the first key < cur
	Prev()

	// Rewind resets the iterator
	// so Next gives the first and Prev gives the last
	Rewind()

	// Seek leaves Cur at the first item >= the given key.
	// It returns true if the key was found.
	// After Seek, Modified returns false.
	Seek(key string) bool
}
