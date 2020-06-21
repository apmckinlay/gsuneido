// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestInfo(t *testing.T) {
	base := NewTableInfoHtbl(0)
	base.Put(&TableInfo{
		table: "one",
		nrows: 100,
		size:  1000,
		indexes: []IndexInfo{
			{root: 11111111111, treeLevels: 0},
			{root: 111111111111, treeLevels: 1},
		},
	})
	base.Put(&TableInfo{
		table: "two",
		nrows: 200,
		size:  2000,
		indexes: []IndexInfo{
			{root: 22222222222, treeLevels: 0},
			{root: 222222222222, treeLevels: 2},
		},
	})
	over := NewTableInfoHtbl(0)
	over.Put(&TableInfo{
		table: "two",
		nrows: 9,
		size:  99,
		indexes: []IndexInfo{
			{root: 22222222220, treeLevels: 1},
			{root: 222222222220, treeLevels: 2},
		},
	})
	merged := base.Merge(over)

	st := stor.HeapStor(blockSize)
	off := merged.WriteInfo(st)

	packed := NewTableInfoPacked(st, off)
	Assert(t).That(*packed.Get("one"), Equals(*base.Get("one")))
	Assert(t).That(*packed.Get("two"), Equals(TableInfo{
		table: "two",
		nrows: 209,
		size:  2099,
		indexes: []IndexInfo{
			{root: 22222222220, treeLevels: 1},
			{root: 222222222220, treeLevels: 2},
		},
	}))

	reread := ReadTablesInfo(st, off)
	Assert(t).That(*reread.Get("one"), Equals(*base.Get("one")))
	Assert(t).That(*reread.Get("two"), Equals(TableInfo{
		table: "two",
		nrows: 209,
		size:  2099,
		indexes: []IndexInfo{
			{root: 22222222220, treeLevels: 1},
			{root: 222222222220, treeLevels: 2},
		},
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
			table: data[i],
			nrows: i,
			size:  1000,
			indexes: []IndexInfo{
				{root: 11111111111, treeLevels: 0},
				{root: 111111111111, treeLevels: 1},
			},
		})
	}
	st := stor.HeapStor(2 * blockSize)
	off := tbl.WriteInfo(st)
	packed := NewTableInfoPacked(st, off)
	for i, s := range data {
		ti := packed.Get(s)
		Assert(t).That(ti.table, Equals(s).Comment("table"))
		Assert(t).That(ti.nrows, Equals(i).Comment("nrows"))
	}
}
