// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestSchema(t *testing.T) {
	tbl := SchemaHamt{}.Mutable()
	const n = 900
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := range n {
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
	off := tbl.Freeze().Write(st, 0, hamt.All)

	test := func(_ int, table string, ts *Schema) {
		t.Helper()
		assert := assert.T(t).This
		assert(ts.Table).Msg("table").Is(table)
		assert(ts.Columns).Msg("columns").Is([]string{"one", "two"})
		assert(ts.Indexes[0].Columns).Msg("indexes").Is([]string{"one"})
	}

	sc := hamt.ReadChain(st, off, ReadSchema)
	assert.T(t).This(sc.Ages[0]).Is(sc.MustGet(data[0]).lastMod)
	for i, table := range data {
		test(i, table, sc.MustGet(table))
	}
}

func TestSetPrimary(t *testing.T) {
	assert := assert.T(t)
	ts := &Schema{Schema: schema.Schema{}}
	key := func(cols string) schema.Index {
		return schema.Index{Mode: 'k', Columns: str.Split(cols, ",")}
	}
	// index := func(cols string, mode int) schema.Index {
	// 	return schema.Index{Mode: mode, Columns: str.Split(cols, ",")}
	// }
	primary := func() string {
		ts.setPrimary()
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

	ts.Indexes = []schema.Index{key("a"), key("b"), key("a_lower!")}
	assert.This(primary()).Is("1,2")
}

func TestOptimizeIndexes(t *testing.T) {
	assert := assert.T(t)
	idx := func(mode byte, cols string) schema.Index {
		return schema.Index{Mode: mode, Columns: str.Split(cols, ",")}
	}
	str := func(ts *Schema) string {
		s := ""
		for _, ix := range ts.Indexes {
			mode := ix.Mode
			if ix.Primary {
				mode = 'K'
			}
			if ix.ContainsKey {
				mode = ascii.ToUpper(mode)
			}
			s += " " + string(mode) + "(" + str.Join(",", ix.Columns)
			add := difference(ix.BestKey, ix.Columns)
			if len(add) > 0 {
				s += "+" + str.Join(",", add)
			}
			s += ")"
		}
		return s[1:]
	}
	ts := &Schema{Schema: schema.Schema{}}
	ts.Indexes = []schema.Index{idx('k', "a"), idx('k', "z,x"),
		idx('i', "b"), idx('u', "c"), idx('i', "b,a"), idx('u', "c,a"),
		idx('i', "x,y,z")}
	ts.SetBestKeys(0)
	ts.setPrimary()
	ts.setContainsKey()
	assert.This(str(ts)).Is("K(a) K(z,x) i(b+a) u(c+a) i(b,a) U(c,a) i(x,y,z)")

	ts.Indexes = []schema.Index{idx('k', "a_lower!"), idx('i', "b")}
	ts.SetBestKeys(0)
	assert.This(str(ts)).Is("k(a_lower!) i(b+a)")

	ts.Indexes = []schema.Index{idx('k', "a"), idx('k', "a_lower!")}
	ts.setPrimary()
	assert.This(str(ts)).Is("k(a) K(a_lower!)")

	ts.Indexes = []schema.Index{idx('k', "a_lower!"), idx('k', "x"),
		idx('u', "b,a"), idx('u', "b,a_lower!"), idx('u', "x_lower!")}
	ts.setContainsKey()
	assert.This(str(ts)).Is("k(a_lower!) k(x) U(b,a) U(b,a_lower!) u(x_lower!)")
}
