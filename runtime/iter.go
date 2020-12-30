// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

// Iter is the internal type for iterators.
// See also: SuIter and wrapIter.
type Iter interface {
	// Next returns nil when there are no more values
	Next() Value
	// Dup returns a copy of this Iter that starts at the beginning
	Dup() Iter
	Infinite() bool
	SetConcurrent()
	IsConcurrent() bool
}
