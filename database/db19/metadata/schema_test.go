// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestSchema(t *testing.T) {
	tbl := NewTableInfoHtbl(0)
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&TableInfo{
			Table: data[i],
			Schema: &TableSchema{
				Table: data[i],
				Columns: []ColumnSchema{
					{Name: "one", Field: i},
					{Name: "two", Field: i*2},
				},
				Indexes: []IndexSchema{
					{Fields: []int{i}},
				},
			},
		})
	}
	st := stor.HeapStor(2 * blockSize)
	off := tbl.WriteSchema(st)

	test := func (i int, table string, ts *TableSchema) {
		Assert(t).That(ts.Table, Equals(table).Comment("table"))
		Assert(t).That(ts.Columns[0].Name, Equals("one").Comment("one"))
		Assert(t).That(ts.Columns[0].Field, Equals(i).Comment("one field"))
		Assert(t).That(ts.Columns[1].Name, Equals("two").Comment("two"))
		Assert(t).That(ts.Columns[0].Field, Equals(i).Comment("two field"))
		Assert(t).That(ts.Indexes[0].Fields, Equals([]int{i}).Comment("indexes"))
	}

	for _,table := range data {
		tbl.Get(table).Schema = nil
	}
	tbl.ReadSchema(st, off)

	packed := NewSchemaPacked(st, off)

	for i, table := range data {
		test(i, table, tbl.Get(table).Schema)
		test(i, table, packed.Get(table))
	}
}
