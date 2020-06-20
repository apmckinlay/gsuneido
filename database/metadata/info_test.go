// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/database/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestInfo(t *testing.T) {
	base := NewTableInfoHtbl(0)
	base.Put(&TableInfo{
		table: 1,
		nrows: 100,
		size:  1000,
		indexes: []IndexInfo{
			{root: 11111111111, treeLevels: 0},
			{root: 111111111111, treeLevels: 1},
		},
	})
	base.Put(&TableInfo{
		table: 2,
		nrows: 200,
		size:  2000,
		indexes: []IndexInfo{
			{root: 22222222222, treeLevels: 0},
			{root: 222222222222, treeLevels: 2},
		},
	})
	over := NewTableInfoHtbl(0)
	over.Put(&TableInfo{
		table: 2,
		nrows: 9,
		size:  99,
		indexes: []IndexInfo{
			{root: 22222222220, treeLevels: 1},
			{root: 222222222220, treeLevels: 2},
		},
	})
	merged := base.Merge(over)

	st := stor.HeapStor(blockSize)
	off := merged.Write(st)

	packed := NewTableInfoPacked(st, off)
	Assert(t).That(*packed.Get(1), Equals(*base.Get(1)))
	Assert(t).That(*packed.Get(2), Equals(TableInfo{
		table: 2,
		nrows: 209,
		size:  2099,
		indexes: []IndexInfo{
			{root: 22222222220, treeLevels: 1},
			{root: 222222222220, treeLevels: 2},
		},
	}))

	reread := ReadTablesInfo(st, off)
	Assert(t).That(*reread.Get(1), Equals(*base.Get(1)))
	Assert(t).That(*reread.Get(2), Equals(TableInfo{
		table: 2,
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
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = rand.Intn(1<<24 - 1)
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
	off := tbl.Write(st)
	packed := NewTableInfoPacked(st, off)
	for i, n := range data {
		ti := packed.Get(n)
		Assert(t).That(ti.table, Equals(n).Comment("table"))
		Assert(t).That(ti.nrows, Equals(i).Comment("nrows"))
	}
}
