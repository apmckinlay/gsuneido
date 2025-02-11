// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = exportMethods(&IterMethods, "iter")

var _ = method(iter_Dup, "()")

func iter_Dup(this Value) Value {
	it := this.(SuIter)
	return SuIter{Iter: it.Dup()}
}

var _ = method(iter_InfiniteQ, "()")

func iter_InfiniteQ(this Value) Value {
	it := this.(SuIter)
	return SuBool(it.Infinite())
}

var _ = method(iter_Next, "()")

func iter_Next(this Value) Value {
	it := this.(SuIter)
	next := it.Next()
	if next == nil {
		return this
	}
	return next
}
