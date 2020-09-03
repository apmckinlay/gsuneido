// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
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
	st.Alloc(1) // don't use offset 0
	off := tbl.Write(st)

	test := func(i int, table string, ts *Schema) {
		t.Helper()
		assert := assert.T(t).This
		assert(ts.Table).Msg("table").Is(table)
		assert(ts.Columns[0].Name).Msg("one").Is("one")
		assert(ts.Columns[0].Field).Msg("one field").Is(i)
		assert(ts.Columns[1].Name).Msg("two").Is("two")
		assert(ts.Columns[0].Field).Msg("two field").Is(i)
		assert(ts.Indexes[0].Fields).Msg("indexes").Is([]int{i})
	}

	tbl = ReadSchemaHamt(st, off)

	packed := NewSchemaPacked(st, off)

	for i, table := range data {
		test(i, table, tbl.MustGet(table))
		test(i, table, packed.MustGet(table))
	}
}
