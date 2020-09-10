// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/ixspec"
	"github.com/apmckinlay/gsuneido/database/db19/stor"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFbtreeIter(t *testing.T) {
	const n = 1000
	var data [n]string
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 440
	randKey := str.UniqueRandomOf(3, 6, "abcde")
	for i := 0; i < n; i++ {
		data[i] = randKey()
	}
	sort.Strings(data[:])
	fb := CreateFbtree(nil, nil)
	fb = fb.Update(func(mfb *fbtree) {
		for i, k := range data {
			mfb.Insert(k, uint64(i))
		}
	})
	i := 0
	iter := fb.Iter()
	for k, o, ok := iter(); ok; k, o, ok = iter() {
		assert.T(t).That(strings.HasPrefix(data[i], k))
		assert.T(t).This(o).Is(i)
		i++
	}
	assert.T(t).This(i).Is(n)
}

func TestFbtreeBuilder(t *testing.T) {
	assert := assert.T(t)
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string {
		return strconv.Itoa(int(i))
	}
	store := stor.HeapStor(8192)
	bldr := NewFbtreeBuilder(store)
	limit := 599999
	if testing.Short() {
		limit = 199999
	}
	for i := 100000; i <= limit; i++ {
		key := strconv.Itoa(i)
		bldr.Add(key, uint64(i))
	}
	fb := bldr.Finish().base()
	fb.check(nil)
	iter := fb.Iter()
	for i := 100000; i <= limit; i++ {
		key := strconv.Itoa(i)
		k, o, ok := iter()
		assert.True(ok)
		assert.True(strings.HasPrefix(key, k))
		assert.This(o).Is(i)
	}
	_, _, ok := iter()
	assert.False(ok)
}
