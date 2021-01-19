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
	// After Seek, Modified returns false.
	Seek(key string)

	// Range sets the range for the iterator.
	// It also does Rewind.
	Range(Range)
}

// Range specifies (key >= org && key < end)
// For key > org, increment org.
// For key <= end, incrment end
type Range struct {
	Org string
	End string
}

const Min = ""
const Max = "\xff\xff\xff\xff\xff\xff\xff\xff"
// Technically there is no maximum key string.
// However, in practice keys are packed values, encoded when composite.
// Packed values start with a type byte from 0 to 7 so 0xff will be larger.
// And 0xff will be larger than any ascii strings.

var All = Range{Org: Min, End: Max}
