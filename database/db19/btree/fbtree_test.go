// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"

	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFbtreeIter(t *testing.T) {
	const n = 1000
	var data [n]string
	GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
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
		Assert(t).True(strings.HasPrefix(data[i], k))
		Assert(t).That(o, Equals(i))
		i++
	}
	Assert(t).That(i, Equals(n))
}

func TestFbtreeBuilder(t *testing.T) {
	store := stor.HeapStor(8192)
	bldr := newFbtreeBuilder(store)
	limit := 999999
	if testing.Short() {
		limit = 199999
	}
	for i := 100000; i <= limit; i++ {
		key := strconv.Itoa(i)
		bldr.Add(key, uint64(i))
	}
	root, treeLevels := bldr.Finish()
	fmt.Println("treeLevels", treeLevels)

	fb := OpenFbtree(store, root, treeLevels, 0)
	iter := fb.Iter()
	for i := 100000; i <= limit; i++ {
		key := strconv.Itoa(i)
		k, o, ok := iter()
		Assert(t).True(ok)
		Assert(t).True(strings.HasPrefix(key, k))
		Assert(t).That(o, Equals(i))
	}
	_, _, ok := iter()
	Assert(t).False(ok)
}
