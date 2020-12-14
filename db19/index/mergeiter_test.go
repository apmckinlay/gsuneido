// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestMergeIter(t *testing.T) {
	from := func(args ...int) *ixbuf.T {
		ib := &ixbuf.T{}
		for _, n := range args {
			ib.Insert(strconv.Itoa(n), uint64(n))
		}
		return ib
	}
	ib1 := from(0, 2, 4, 6, 8)
	ib2 := from(1, 3, 5, 7, 9)
	it := NewMergeIter([]iterator{ib1.Iterator(), ib2.Iterator()})
	assert := assert.T(t)
	test := func(expected int) {
		if expected == -1 {
			assert.That(it.Eof())
		} else {
			key, off := it.Cur()
			assert.This(key).Is(strconv.Itoa(expected))
			assert.This(off).Is(uint64(expected))
		}
	}
	testNext := func(expected int) { it.Next(); test(expected) }
	testPrev := func(expected int) { it.Prev(); test(expected) }
	for i := 0; i < 10; i++ {
		testNext(i)
	}
	testNext(-1)

	it.Rewind()
	for i := 9; i >= 0; i-- {
		testPrev(i)
	}
	testPrev(-1)
}
