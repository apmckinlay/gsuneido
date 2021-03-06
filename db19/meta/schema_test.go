// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func allSchema(*Schema) bool { return true }

func TestSchema(t *testing.T) {
	tbl := SchemaHamt{}.Mutable()
	const n = 900
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&Schema{Schema: schema.Schema{
			Table:   data[i],
			Columns: []string{"one", "two"},
			Indexes: []schema.Index{
				{Mode: 'k', Columns: []string{"one"}},
			},
		}})
	}
	st := stor.HeapStor(32 * 1024)
	st.Alloc(1) // avoid offset 0
	off := tbl.Write(st, 0, allSchema)

	test := func(i int, table string, ts *Schema) {
		t.Helper()
		assert := assert.T(t).This
		assert(ts.Table).Msg("table").Is(table)
		assert(ts.Columns).Msg("columns").Is([]string{"one", "two"})
		assert(ts.Indexes[0].Columns).Msg("indexes").Is([]string{"one"})
	}

	tbl,_ = ReadSchemaChain(st, off)
	for i, table := range data {
		test(i, table, tbl.MustGet(table))
	}
}
