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

	// optimizeFor dispatches ReqLookup to optimizeLookup2,
	// all other uses to optimize2 (matching the design: Optimize2
	// never receives ReqLookup; lookup uses a separate interface).
	optimizeFor := func(tbl *Table, req Require, frac float64) (Cost, Cost, any) {
		if req.use == ReqLookup {
			return tbl.optimizeLookup2(ReadMode, req.cols, frac)
		}
		return tbl.optimize2(ReadMode, req, frac)
	}

	test := func(table string, req Require, frac float64, expected []string) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, app := optimizeFor(tbl, req, frac)
		assert.True(f+v < impossible)
		assert.This(app.(tableApproach).index).Is(expected)
	}
	assertImpossible := func(table string, req Require, frac float64) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, app := optimizeFor(tbl, req, frac)
		assert.False(f+v < impossible)
		assert.This(app).Is(nil)
	}

	// table with single index: key on {a}
	// indexes: [[a]], allKeys: [[a]]

	// ReqUnordered
	test("table", reqUnordered, 1, []string{"a"})

	// ReqOrdered — match
	test("table", Require{ReqOrdered, []string{"a"}}, 1, []string{"a"})

	// ReqOrdered — no match (b not a prefix of {a})
	assertImpossible("table", Require{ReqOrdered, []string{"b"}}, 1)

	// ReqGrouped — match
	test("table", Require{ReqGrouped, []string{"a"}}, 1, []string{"a"})

	// ReqGrouped — b not found in any index
	assertImpossible("table", Require{ReqGrouped, []string{"b"}}, 1)

	// ReqLookup — {a} is a physical index
	test("table", Require{ReqLookup, []string{"a"}}, 1, []string{"a"})

	// ReqLookup — {b} is not a physical index, no lookup-eligible index has {b}
	assertImpossible("table", Require{ReqLookup, []string{"b"}}, 1)

	// supplier: key on {supplier}, index on {city}
	// indexes: [[supplier], [city,supplier]], allKeys: [[supplier]]

	// ReqUnordered — primary index
	test("supplier", reqUnordered, 1, []string{"supplier"})

	// ReqOrdered — match key
	test("supplier", Require{ReqOrdered, []string{"supplier"}}, 1, []string{"supplier"})

	// ReqOrdered — match non-key index
	test("supplier", Require{ReqOrdered, []string{"city", "supplier"}}, 1, []string{"city", "supplier"})

	// ReqOrdered — no match (name not in any index)
	assertImpossible("supplier", Require{ReqOrdered, []string{"name"}}, 1)

	// ReqOrdered — city is a prefix of [city,supplier], so indexFor finds it
	test("supplier", Require{ReqOrdered, []string{"city"}}, 1, []string{"city", "supplier"})

	// ReqGrouped — {city} is grouped by [city,supplier]
	test("supplier", Require{ReqGrouped, []string{"city"}}, 1, []string{"city", "supplier"})

	// ReqGrouped — {name} not in any index
	assertImpossible("supplier", Require{ReqGrouped, []string{"name"}}, 1)

	// ReqLookup — {supplier} is a physical index
	test("supplier", Require{ReqLookup, []string{"supplier"}}, 1, []string{"supplier"})

	// ReqLookup — {city,supplier} is a physical index
	test("supplier", Require{ReqLookup, []string{"city", "supplier"}}, 1, []string{"city", "supplier"})

	// ReqLookup — {city} is not a physical index and no index is both
	// lookup-eligible (contains a key as a prefix) AND grouped by {city}.
	// [supplier] is lookup-eligible but not grouped by {city}.
	// [city,supplier] is grouped by {city} but NOT lookup-eligible
	// (key {supplier} is not a prefix of [city,supplier]).
	assertImpossible("supplier", Require{ReqLookup, []string{"city"}}, 1)

	// abc: key on {b}, index on {a}, key on {c}
	// indexes: [[b], [a,b], [c]], allKeys: [[b], [c]]

	// ReqUnordered — primary index {b}
	test("abc", reqUnordered, 1, []string{"b"})

	// ReqGrouped — {a} is grouped by [a,b]
	test("abc", Require{ReqGrouped, []string{"a"}}, 1, []string{"a", "b"})

	// ReqGrouped — {b} is grouped by [b] (cheapest)
	test("abc", Require{ReqGrouped, []string{"b"}}, 1, []string{"b"})

	// ReqGrouped — {c} is grouped by [c]
	test("abc", Require{ReqGrouped, []string{"c"}}, 1, []string{"c"})

	// ReqLookup — {b} is a physical index (and a key)
	test("abc", Require{ReqLookup, []string{"b"}}, 1, []string{"b"})

	// ReqLookup — {c} is a physical index (and a key)
	test("abc", Require{ReqLookup, []string{"c"}}, 1, []string{"c"})

	// ReqLookup — {a} is not a physical index. [a,b] is not
	// lookup-eligible (key {b} is not a prefix, key {c} is not in index).
	assertImpossible("abc", Require{ReqLookup, []string{"a"}}, 1)

	// comp: key on {a,b,c} — single index [a,b,c], keys [[a,b,c]]
	// ReqLookup {b,a} is not a physical index, but the fallback search
	// finds [a,b,c]: it is lookup-eligible (starts with key) and has
	// {b,a} grouped (first 2 cols {a,b} are in {b,a}).
	test("comp", Require{ReqLookup, []string{"b", "a"}}, 1,
		[]string{"a", "b", "c"})

	// singleton: table with empty key — all req types return indexes[0]
	// Set up manually (no SetTran since _singleton_ is not in testSchemas).
	singleton := &Table{}
	singleton.name = "_singleton_"
	singleton.indexes = [][]string{{"x"}}
	singleton.singleton = true
	singleton.info = &meta.Info{Nrows: 1, Size: 100}
	for _, req := range []Require{
		reqUnordered,
		{ReqOrdered, []string{"x"}},
		{ReqGrouped, []string{"x"}},
		{ReqLookup, []string{"x"}},
	} {
		f, v, app := optimizeFor(singleton, req, 1)
		assert.Msg(req).True(f+v < impossible)
		assert.Msg(req).This(app.(tableApproach).index).Is([]string{"x"})
	}
}