// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInfo(t *testing.T) {
	assert := assert.T(t).This
	tbl := InfoHamt{}.Mutable()
	tbl.Put(&Info{
		Table:   "one",
		Nrows:   100,
		Size:    1000,
		Indexes: []*btree.Overlay{},
	})
	tbl.Put(&Info{
		Table:   "two",
		Nrows:   200,
		Size:    2000,
		Indexes: []*btree.Overlay{},
	})

	st := stor.HeapStor(8192)
	st.Alloc(1) // avoid offset 0
	off := tbl.Write(st)

	tbl = InfoHamt{}.Mutable().Read(st, off)
	assert(*tbl.MustGet("one")).Is(*tbl.MustGet("one"))
	assert(*tbl.MustGet("two")).Is(Info{
		Table:   "two",
		Nrows:   200,
		Size:    2000,
		Indexes: []*btree.Overlay{},
	})
}

func TestInfo2(t *testing.T) {
	tbl := InfoHamt{}.Mutable()
	const n = 1000
	data := mkdata(tbl, n)
	st := stor.HeapStor(32 * 1024)
	st.Alloc(1) // avoid offset 0
	off := tbl.Write(st)

	tbl = InfoHamt{}.Mutable().Read(st, off).Freeze()
	for i, s := range data {
		ti := tbl.MustGet(s)
		assert.T(t).Msg("table").This(ti.Table).Is(s)
		assert.T(t).Msg("nrows").This(ti.Nrows).Is(i)
		_, ok := tbl.Get(s + "Z")
		assert.T(t).That(!ok)
	}
}

func mkdata(tbl InfoHamt, n int) []string {
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&Info{Table: data[i], Nrows: i})
	}
	return data
}
