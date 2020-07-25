// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInfo(t *testing.T) {
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
	off := tbl.Write(st)

	packed := NewInfoPacked(st, off)

	Assert(t).That(*packed.MustGet("one"), Equals(*tbl.MustGet("one")))
	Assert(t).That(*packed.MustGet("two"), Equals(Info{
		Table:   "two",
		Nrows:   200,
		Size:    2000,
		Indexes: []*btree.Overlay{},
	}))

	reread := ReadInfoHamt(st, off)
	Assert(t).That(*reread.MustGet("one"), Equals(*tbl.MustGet("one")))
	Assert(t).That(*reread.MustGet("two"), Equals(Info{
		Table:   "two",
		Nrows:   200,
		Size:    2000,
		Indexes: []*btree.Overlay{},
	}))
}

func TestInfo2(t *testing.T) {
	tbl := InfoHamt{}.Mutable()
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&Info{
			Table: data[i],
			Nrows: i,
			Size:  1000,
		})
	}
	st := stor.HeapStor(32 * 1024)
	off := tbl.Write(st)
	packed := NewInfoPacked(st, off)
	for i, s := range data {
		ti := packed.MustGet(s)
		Assert(t).That(ti.Table, Equals(s).Comment("table"))
		Assert(t).That(ti.Nrows, Equals(i).Comment("nrows"))
	}
}
