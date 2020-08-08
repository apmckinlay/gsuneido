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
	tbl := SchemaHamt{}.Mutable()
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
	st := stor.HeapStor(32 * 1024)
	off := tbl.Write(st)

	test := func(i int, table string, ts *Schema) {
		Assert(t).That(ts.Table, Is(table).Comment("table"))
		Assert(t).That(ts.Columns[0].Name, Is("one").Comment("one"))
		Assert(t).That(ts.Columns[0].Field, Is(i).Comment("one field"))
		Assert(t).That(ts.Columns[1].Name, Is("two").Comment("two"))
		Assert(t).That(ts.Columns[0].Field, Is(i).Comment("two field"))
		Assert(t).That(ts.Indexes[0].Fields, Is([]int{i}).Comment("indexes"))
	}

	tbl = ReadSchemaHamt(st, off)

	packed := NewSchemaPacked(st, off)

	for i, table := range data {
		test(i, table, tbl.MustGet(table))
		test(i, table, packed.MustGet(table))
	}
}
