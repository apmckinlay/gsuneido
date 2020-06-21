// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/stor"
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
			table: data[i],
			schema: &TableSchema{
				table: data[i],
				columns: []ColumnSchema{
					{name: "one", field: i},
					{name: "two", field: i*2},
				},
				indexes: []IndexSchema{
					{fields: []int{i}},
				},
			},
		})
	}
	st := stor.HeapStor(2 * blockSize)
	off := tbl.WriteSchema(st)

	test := func (i int, table string, ts *TableSchema) {
		Assert(t).That(ts.table, Equals(table).Comment("table"))
		Assert(t).That(ts.columns[0].name, Equals("one").Comment("one"))
		Assert(t).That(ts.columns[0].field, Equals(i).Comment("one field"))
		Assert(t).That(ts.columns[1].name, Equals("two").Comment("two"))
		Assert(t).That(ts.columns[0].field, Equals(i).Comment("two field"))
		Assert(t).That(ts.indexes[0].fields, Equals([]int{i}).Comment("indexes"))
	}

	for _,table := range data {
		tbl.Get(table).schema = nil
	}
	tbl.ReadSchema(st, off)

	packed := NewSchemaPacked(st, off)

	for i, table := range data {
		test(i, table, tbl.Get(table).schema)
		test(i, table, packed.Get(table))
	}
}
