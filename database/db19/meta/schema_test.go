// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestSchema(t *testing.T) {
	tbl := NewSchemaHtbl(0)
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&Schema{
			Table: data[i],
			Columns: []ColumnSchema{
				{Name: "one", Field: i},
				{Name: "two", Field: i * 2},
			},
			Indexes: []IndexSchema{
				{Fields: []int{i}},
			},
		})
	}
	st := stor.HeapStor(8192)
	off := tbl.Write(st)

	test := func(i int, table string, ts *Schema) {
		Assert(t).That(ts.Table, Equals(table).Comment("table"))
		Assert(t).That(ts.Columns[0].Name, Equals("one").Comment("one"))
		Assert(t).That(ts.Columns[0].Field, Equals(i).Comment("one field"))
		Assert(t).That(ts.Columns[1].Name, Equals("two").Comment("two"))
		Assert(t).That(ts.Columns[0].Field, Equals(i).Comment("two field"))
		Assert(t).That(ts.Indexes[0].Fields, Equals([]int{i}).Comment("indexes"))
	}

	tbl = ReadSchemaHtbl(st, off)

	packed := NewSchemaPacked(st, off)

	for i, table := range data {
		test(i, table, tbl.Get(table))
		test(i, table, packed.Get(table))
	}
}
