// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

// Iter is the internal type for iterators.
// builtin.SuIter wraps Iter and implements Value and methods
type Iter interface {
	Next() Value
	Infinite() bool
	Dup() Iter
}
