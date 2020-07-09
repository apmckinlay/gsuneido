// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMerge(t *testing.T) {
	randKey := str.UniqueRandomOf(3, 10, "abcdef")
	var data []string
	randMbtree := func() *mbtree {
		const n = mSize * 3
		mb := newMbtree(0)
		for i := 0; i < n; i++ {
			key := randKey()
			off := uint64(len(data))
			data = append(data, key)
			mb.Insert(key, off)
		}
		return mb
	}
	mb := randMbtree()
	mb.checkData(t, data)
	GetLeafKey = func(_ *stor.Stor, _ interface{}, i uint64) string { return data[i] }
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	fb := CreateFbtree(nil)
	fb = Merge(fb, mb)
	fb.checkData(t, data)

	mb = randMbtree()
	fb = Merge(fb, mb)
	fb.checkData(t, data)
}
