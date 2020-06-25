// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInfo(t *testing.T) {
	tbl := NewTableInfoHtbl(0)
	tbl.Put(&TableInfo{
		Table: "one",
		Nrows: 100,
		Size:  1000,
		Indexes: []*btree.Overlay{},
	})
	tbl.Put(&TableInfo{
		Table:   "two",
		Nrows:   200,
		Size:    2000,
		Indexes: []*btree.Overlay{},
	})
	// over := NewTableInfoHtbl(0)
	// over.Put(&TableInfo{
	// 	Table: "two",
	// 	Nrows: 9,
	// 	Size:  99,
	// 	Indexes: []*btree.Overlay{},
	// })
	// merged := base.Merge(over)

	st := stor.HeapStor(blockSize)
	off := tbl.WriteInfo(st)

	packed := NewInfoPacked(st, off)
	Assert(t).That(*packed.Get("one"), Equals(*tbl.Get("one")))
	Assert(t).That(*packed.Get("two"), Equals(TableInfo{
		Table: "two",
		Nrows: 200,
		Size:  2000,
		Indexes: []*btree.Overlay{},
	}))

	reread := ReadInfo(st, off)
	Assert(t).That(*reread.Get("one"), Equals(*tbl.Get("one")))
	Assert(t).That(*reread.Get("two"), Equals(TableInfo{
		Table: "two",
		Nrows: 200,
		Size:  2000,
		Indexes: []*btree.Overlay{},
	}))
}

func TestMetadata2(t *testing.T) {
	tbl := NewTableInfoHtbl(0)
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&TableInfo{
			Table: data[i],
			Nrows: i,
			Size:  1000,
		})
	}
	st := stor.HeapStor(2 * blockSize)
	off := tbl.WriteInfo(st)
	packed := NewInfoPacked(st, off)
	for i, s := range data {
		ti := packed.Get(s)
		Assert(t).That(ti.Table, Equals(s).Comment("table"))
		Assert(t).That(ti.Nrows, Equals(i).Comment("nrows"))
	}
}
