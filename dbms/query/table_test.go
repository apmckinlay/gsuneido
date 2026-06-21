// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTableOptimize2(t *testing.T) {
	assert := assert.T(t)
	tran := testTran{}

	optimizeFor := func(tbl *Table, req Require) (Cost, Cost, any) {
		return tbl.optimize2(ReadMode, req)
	}

	test := func(table string, req Require, expected []string) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, app := optimizeFor(tbl, req)
		assert.True(f+v < impossible)
		assert.This(app.(tableApproach).index).Is(expected)
	}
	assertImpossible := func(table string, req Require) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, app := optimizeFor(tbl, req)
		assert.False(f+v < impossible)
		assert.This(app).Is(nil)
	}

	// table with single index: key on {a}
	// indexes: [[a]], allKeys: [[a]]

	// ReqUnordered
	test("table", Require{frac: 1}, []string{"a"})

	// ReqOrdered — match
	test("table", Require{cols: []string{"a"}, frac: 1}, []string{"a"})

	// ReqOrdered — no match (b not a prefix of {a})
	assertImpossible("table", Require{cols: []string{"b"}, frac: 1})

	// ReqGrouped — match
	test("table", Require{cols: []string{"a"}, frac: 1, nlookups: 1}, []string{"a"})

	// ReqGrouped — b not found in any index
	assertImpossible("table", Require{cols: []string{"b"}, frac: 1, nlookups: 1})

	// ReqLookup — {a} is a physical index
	test("table", Require{cols: []string{"a"}, frac: 0, nlookups: 1}, []string{"a"})

	// ReqLookup — {b} is not a physical index, no lookup-eligible index has {b}
	assertImpossible("table", Require{cols: []string{"b"}, frac: 0, nlookups: 1})

	// supplier: key on {supplier}, index on {city}
	// indexes: [[supplier], [city,supplier]], allKeys: [[supplier]]

	// ReqUnordered — primary index
	test("supplier", Require{frac: 1}, []string{"supplier"})

	// ReqOrdered — match key
	test("supplier", Require{cols: []string{"supplier"}, frac: 1}, []string{"supplier"})

	// ReqOrdered — match non-key index
	test("supplier", Require{cols: []string{"city", "supplier"}, frac: 1}, []string{"city", "supplier"})

	// ReqOrdered — no match (name not in any index)
	assertImpossible("supplier", Require{cols: []string{"name"}, frac: 1})

	// ReqOrdered — city is a prefix of [city,supplier], so indexFor finds it
	test("supplier", Require{cols: []string{"city"}, frac: 1}, []string{"city", "supplier"})

	// ReqGrouped — {city} is grouped by [city,supplier]
	test("supplier", Require{cols: []string{"city"}, frac: 1, nlookups: 1}, []string{"city", "supplier"})

	// ReqGrouped — {name} not in any index
	assertImpossible("supplier", Require{cols: []string{"name"}, frac: 1, nlookups: 1})

	// ReqLookup — {supplier} is a physical index
	test("supplier", Require{cols: []string{"supplier"}, frac: 0, nlookups: 1}, []string{"supplier"})

	// ReqLookup — {city,supplier} is a physical index
	test("supplier", Require{cols: []string{"city", "supplier"}, frac: 0, nlookups: 1}, []string{"city", "supplier"})

	// ReqLookup — {city} is not a physical index and no index is both
	// lookup-eligible (contains a key as a prefix) AND grouped by {city}.
	// [supplier] is lookup-eligible but not grouped by {city}.
	// [city,supplier] is grouped by {city} but NOT lookup-eligible
	// (key {supplier} is not a prefix of [city,supplier]).
	assertImpossible("supplier", Require{cols: []string{"city"}, frac: 0, nlookups: 1})

	// abc: key on {b}, index on {a}, key on {c}
	// indexes: [[b], [a,b], [c]], allKeys: [[b], [c]]

	// ReqUnordered — primary index {b}
	test("abc", Require{frac: 1}, []string{"b"})

	// ReqGrouped — {a} is grouped by [a,b]
	test("abc", Require{cols: []string{"a"}, frac: 1, nlookups: 1}, []string{"a", "b"})

	// ReqGrouped — {b} is grouped by [b] (cheapest)
	test("abc", Require{cols: []string{"b"}, frac: 1, nlookups: 1}, []string{"b"})

	// ReqGrouped — {c} is grouped by [c]
	test("abc", Require{cols: []string{"c"}, frac: 1, nlookups: 1}, []string{"c"})

	// ReqLookup — {b} is a physical index (and a key)
	test("abc", Require{cols: []string{"b"}, frac: 0, nlookups: 1}, []string{"b"})

	// ReqLookup — {c} is a physical index (and a key)
	test("abc", Require{cols: []string{"c"}, frac: 0, nlookups: 1}, []string{"c"})

	// ReqLookup — {a} is not a physical index. [a,b] is not
	// lookup-eligible (key {b} is not a prefix, key {c} is not in index).
	assertImpossible("abc", Require{cols: []string{"a"}, frac: 0, nlookups: 1})

	// comp: key on {a,b,c} — single index [a,b,c], keys [[a,b,c]]
	// ReqLookup {b,a} is not a physical index, but the fallback search
	// finds [a,b,c]: it is lookup-eligible (starts with key) and has
	// {b,a} grouped (first 2 cols {a,b} are in {b,a}).
	test("comp", Require{cols: []string{"b", "a"}, frac: 0, nlookups: 1},
		[]string{"a", "b", "c"})

	// singleton: table with empty key — all req types return indexes[0]
	// Set up manually (no SetTran since _singleton_ is not in testSchemas).
	singleton := &Table{}
	singleton.name = "_singleton_"
	singleton.indexes = [][]string{{"x"}}
	singleton.singleton = true
	singleton.info = &meta.Info{Nrows: 1, Size: 100}
	for _, req := range []Require{
		{frac: 1},
		{cols: []string{"x"}, frac: 1},
		{cols: []string{"x"}, frac: 1, nlookups: 1},
		{cols: []string{"x"}, frac: 0, nlookups: 1},
	} {
		f, v, app := optimizeFor(singleton, req)
		assert.Msg(req).True(f+v < impossible)
		assert.Msg(req).This(app.(tableApproach).index).Is([]string{"x"})
	}
}
