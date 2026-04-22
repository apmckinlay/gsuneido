// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package iface

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
)

// IterFn is the function type for iterating over ixbuf entries.
type IterFn = func() (key string, off uint64, ok bool)

// Iter is the interface for a Suneido style iterator
// implemented by btree and ixbuf
type Iter interface {
	// Eof returns true if Next hit the end, or Prev hit the beginning.
	// Returns false when rewound, even if the index is empty.
	Eof() bool

	// Modified returns true if the index has been modified.
	// Seek resets modified.
	Modified() bool

	// Cur returns the current key & offset
	// as of the most recent Next, Prev, or Seek
	Cur() (key string, off uint64)

	// Key returns the current key.
	// Returns ixkey.Max when eof (regardless of direction).
	Key() string
	Offset() uint64

	HasCur() bool

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

	// SeekAll is the same as Seek
	// except it does NOT use the range to set eof
	SeekAll(key string)

	// Range sets the range for the iterator.
	// It also does Rewind.
	Range(Range)

	// SkipScan enables skip-scan mode.
	// prefixRng restricts visited prefix groups; iface.All means unrestricted.
	// suffixRng applies to suffix fields (excluding prefix fields).
	// skipStart must be >= 1
	SkipScan(prefixRng Range, suffixRng Range, skipStart int)
}

// Range specifies (key >= org && key < end)
// For key > org, increment org.
// For key <= end, increment end
type Range struct {
	Org string
	End string
}

var All = Range{Org: ixkey.Min, End: ixkey.Max}

func (r Range) String() string {
	if r.Org == ixkey.Min && r.End == ixkey.Max {
		return "{all}"
	}
	return fmt.Sprintf("%s=>%s", r.Org, r.End)
}
