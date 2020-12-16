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
	even := from(0, 2, 4, 6, 8)
	odd := from(1, 3, 5, 7, 9)
	modCount := 0
	callback := func(mc int) (int, []iterator) {
		if mc == modCount {
			return mc, nil
		}
		return modCount, []iterator{even.Iterator(), odd.Iterator()}
	}
	it := NewMergeIter(callback)
	assert := assert.T(t)
	test := func(expected int) {
		t.Helper()
		if expected == -1 {
			assert.That(it.Eof())
		} else {
			key, off := it.Cur()
			assert.This(key).Is(strconv.Itoa(expected))
			assert.This(off).Is(uint64(expected))
		}
	}
	testNext := func(expected int) { it.Next(); t.Helper(); test(expected) }
	testPrev := func(expected int) { it.Prev(); t.Helper(); test(expected) }
	for i := 0; i < 10; i++ {
		testNext(i)
		if i == 5 {
			modCount++
		}
	}
	testNext(-1)

	it.Rewind()
	for i := 9; i >= 0; i-- {
		testPrev(i)
		if i == 5 {
			modCount++
		}
	}
	testPrev(-1)

	it.Rewind()
	testNext(0)
	testPrev(-1) // stick at eof
	testPrev(-1)
	testNext(-1)

	it.Rewind()
	testPrev(9)
	testPrev(8)
	testPrev(7)
	testNext(8)
	testNext(9) // last
	testPrev(8)

	it.Rewind()
	testNext(0)
	testNext(1)
	even.Insert("11", uint64(11))
	testNext(11)
	testNext(2)
	even.Delete("11", uint64(11))
	testPrev(1) // modified AND changed direction

	it.Rewind()
	testPrev(9)
	testPrev(8)
	odd.Insert("77", uint64(77))
	testPrev(77)
	modCount++
	testPrev(7)
}
