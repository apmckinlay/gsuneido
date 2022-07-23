// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
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

	sc := hamt.ReadChain[string](st, off, ReadSchema)
	assert.T(t).This(sc.Ages[0]).Is(sc.MustGet(data[0]).lastMod)
	for i, table := range data {
		test(i, table, sc.MustGet(table))
	}
}

func TestFindPrimaryKeys(t *testing.T) {
	assert := assert.T(t)
	ts := &Schema{Schema: schema.Schema{}}
	key := func(cols string) schema.Index {
		return schema.Index{Mode: 'k', Columns: str.Split(cols, ",")}
	}
	// index := func(cols string, mode int) schema.Index {
	// 	return schema.Index{Mode: mode, Columns: str.Split(cols, ",")}
	// }
	primary := func() string {
		ts.findPrimaryKeys()
		s := ""
		for i, ix := range ts.Indexes {
			if ix.Primary {
				s += "," + strconv.Itoa(i)
			}
		}
		return s[1:]
	}
	ts.Indexes = []schema.Index{key("")}
	assert.This(primary()).Is("0")

	ts.Indexes = []schema.Index{key(""), key("a")}
	assert.This(primary()).Is("0")

	ts.Indexes = []schema.Index{key("a,b"), key("b,a"), key("b,c,a")}
	assert.This(primary()).Is("0")
}

func TestOptimizeIndexes(t *testing.T) {
	assert := assert.T(t)
	idx := func(mode byte, cols string) schema.Index {
		return schema.Index{Mode: mode, Columns: str.Split(cols, ",")}
	}
	str := func(ts *Schema) string {
		s := ""
		for _, ix := range ts.Indexes {
			s += " " + string(ix.Mode) + str.Join("(,)", ix.Columns)
		}
		return s[1:]
	}
	ts := &Schema{Schema: schema.Schema{}}
	ts.Indexes = []schema.Index{idx('k', "a"), idx('k', "z,x"),
		idx('i', "b"), idx('u', "c"), idx('i', "b,a"), idx('u', "c,a"),
		idx('i', "x,y,z")}
	ts.OptimizeIndexes()
	assert.This(str(ts)).Is("k(a) k(z,x) i(b) u(c) I(b,a) U(c,a) I(x,y,z)")
}
