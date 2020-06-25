// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"

	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFbtreeIter(t *testing.T) {
	const n = 1000
	var data [n]string
	GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string { return data[i] }
	MaxNodeSize = 440
	fb := CreateFbtree(nil)
	up := newFbupdate(fb)
	randKey := str.UniqueRandomOf(3, 6, "abcde")
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	sort.Strings(data[:])
	for i, k := range data {
		up.Insert(k, uint64(i))
	}
	fb = up.freeze()
	i := 0
	iter := fb.Iter()
	for k, o, ok := iter(); ok; k, o, ok = iter() {
		Assert(t).That(k, Equals(data[i]))
		Assert(t).That(o, Equals(i))
		i++
	}
	Assert(t).That(i, Equals(n))
}
