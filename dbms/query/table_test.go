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
		return tbl.optimize(ReadMode, req)
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
	test("table", UnorderedReq(1), []string{"a"})

	// ReqOrdered — match
	test("table", OrderedReq([]string{"a"}, 1), []string{"a"})

	// ReqOrdered — no match (b not a prefix of {a})
	assertImpossible("table", OrderedReq([]string{"b"}, 1))

	// ReqGrouped — match
	test("table", GroupedReq([]string{"a"}, 1, 1), []string{"a"})

	// ReqGrouped — b not found in any index
	assertImpossible("table", GroupedReq([]string{"b"}, 1, 1))

	// ReqLookup — {a} is a physical index
	test("table", LookupReq([]string{"a"}, 1), []string{"a"})

	// ReqLookup — {b} is not a physical index, no lookup-eligible index has {b}
	assertImpossible("table", LookupReq([]string{"b"}, 1))

	// supplier: key on {supplier}, index on {city}
	// indexes: [[supplier], [city,supplier]], allKeys: [[supplier]]

	// ReqUnordered — primary index
	test("supplier", UnorderedReq(1), []string{"supplier"})

	// ReqOrdered — match key
	test("supplier", OrderedReq([]string{"supplier"}, 1), []string{"supplier"})

	// ReqOrdered — match non-key index
	test("supplier", OrderedReq([]string{"city", "supplier"}, 1), []string{"city", "supplier"})

	// ReqOrdered — no match (name not in any index)
	assertImpossible("supplier", OrderedReq([]string{"name"}, 1))

	// ReqOrdered — city is a prefix of [city,supplier], so indexFor finds it
	test("supplier", OrderedReq([]string{"city"}, 1), []string{"city", "supplier"})

	// ReqGrouped — {city} is grouped by [city,supplier]
	test("supplier", GroupedReq([]string{"city"}, 1, 1), []string{"city", "supplier"})

	// ReqGrouped — {name} not in any index
	assertImpossible("supplier", GroupedReq([]string{"name"}, 1, 1))

	// ReqLookup — {supplier} is a physical index
	test("supplier", LookupReq([]string{"supplier"}, 1), []string{"supplier"})

	// ReqLookup — {city,supplier} is a physical index
	test("supplier", LookupReq([]string{"city", "supplier"}, 1), []string{"city", "supplier"})

	// ReqLookup — {city} is not a physical index and no index is both
	// lookup-eligible (contains a key as a prefix) AND grouped by {city}.
	// [supplier] is lookup-eligible but not grouped by {city}.
	// [city,supplier] is grouped by {city} but NOT lookup-eligible
	// (key {supplier} is not a prefix of [city,supplier]).
	assertImpossible("supplier", LookupReq([]string{"city"}, 1))

	// abc: key on {b}, index on {a}, key on {c}
	// indexes: [[b], [a,b], [c]], allKeys: [[b], [c]]

	// ReqUnordered — primary index {b}
	test("abc", UnorderedReq(1), []string{"b"})

	// ReqGrouped — {a} is grouped by [a,b]
	test("abc", GroupedReq([]string{"a"}, 1, 1), []string{"a", "b"})

	// ReqGrouped — {b} is grouped by [b] (cheapest)
	test("abc", GroupedReq([]string{"b"}, 1, 1), []string{"b"})

	// ReqGrouped — {c} is grouped by [c]
	test("abc", GroupedReq([]string{"c"}, 1, 1), []string{"c"})

	// ReqLookup — {b} is a physical index (and a key)
	test("abc", LookupReq([]string{"b"}, 1), []string{"b"})

	// ReqLookup — {c} is a physical index (and a key)
	test("abc", LookupReq([]string{"c"}, 1), []string{"c"})

	// ReqLookup — {a} is not a physical index. [a,b] is not
	// lookup-eligible (key {b} is not a prefix, key {c} is not in index).
	assertImpossible("abc", LookupReq([]string{"a"}, 1))

	// comp: key on {a,b,c} — single index [a,b,c], keys [[a,b,c]]
	// ReqLookup {b,a} is not a physical index.
	// [a,b,c] is lookup-eligible and has {b,a} grouped,
	// but is not indexCovered (c is not in {b,a}).
	assertImpossible("comp", LookupReq([]string{"b", "a"}, 1))

	// singleton: table with empty key — all req types return indexes[0]
	// Set up manually (no SetTran since _singleton_ is not in testSchemas).
	singleton := &Table{}
	singleton.name = "_singleton_"
	singleton.indexes = [][]string{{"x"}}
	singleton.singleton = true
	singleton.info = &meta.Info{Nrows: 1, Size: 100}
	for _, req := range []Require{
		UnorderedReq(1),
		OrderedReq([]string{"x"}, 1),
		GroupedReq([]string{"x"}, 1, 1),
		LookupReq([]string{"x"}, 1),
	} {
		f, v, app := optimizeFor(singleton, req)
		assert.Msg(req).True(f+v < impossible)
		assert.Msg(req).This(app.(tableApproach).index).Is([]string{"x"})
	}
}

func TestTableOptimize2_ReqLookup_indexCovered(t *testing.T) {
	assert := assert.T(t)
	optimizeFor := func(tbl *Table, req Require) (Cost, Cost, any) {
		return tbl.optimize(ReadMode, req)
	}
	test := func(tbl *Table, req Require, expected []string) {
		t.Helper()
		f, v, app := optimizeFor(tbl, req)
		assert.True(f+v < impossible)
		assert.This(app.(tableApproach).index).Is(expected)
	}
	assertImpossible := func(tbl *Table, req Require) {
		t.Helper()
		f, v, app := optimizeFor(tbl, req)
		assert.False(f+v < impossible)
		assert.This(app).Is(nil)
	}
	newTable := func(name string, indexes, keys [][]string, nrows int) *Table {
		tbl := &Table{name: name}
		tbl.indexes = indexes
		tbl.allKeys = keys
		tbl.info = &meta.Info{Nrows: nrows, Size: int64(nrows * 100)}
		return tbl
	}

	// by=(x,y), key=(y): partial lookup — indexCovered passes (y is in by)
	tbl := newTable("xy", [][]string{{"y"}}, [][]string{{"y"}}, 100)
	test(tbl, LookupReq([]string{"x", "y"}, 1), []string{"y"})

	// by=(x,y), key=(x,y): full lookup — exact index match
	tbl = newTable("xy2", [][]string{{"x", "y"}}, [][]string{{"x", "y"}}, 100)
	test(tbl, LookupReq([]string{"x", "y"}, 1), []string{"x", "y"})

	// by=(x,y), key=(x,y,z): should fail — indexCovered fails (z not in by)
	tbl = newTable("xyz", [][]string{{"x", "y", "z"}}, [][]string{{"x", "y", "z"}}, 100)
	assertImpossible(tbl, LookupReq([]string{"x", "y"}, 1))

	// by=(x,y), index=(x,c), key=(x): fails — indexCovered fails (c not in by)
	tbl = newTable("xc", [][]string{{"x", "c"}}, [][]string{{"x"}}, 100)
	assertImpossible(tbl, LookupReq([]string{"x", "y"}, 1))

	// by=(x,y), index=(x), no key: fails — lookupIndexEligible fails
	tbl = newTable("xnokey", [][]string{{"x"}}, [][]string{}, 100)
	assertImpossible(tbl, LookupReq([]string{"x", "y"}, 1))

	// by=(x,y), indexes=[x],[y], key=(y): picks [y] (smallest eligible index)
	tbl = newTable("twoidx", [][]string{{"x"}, {"y"}}, [][]string{{"y"}}, 100)
	test(tbl, LookupReq([]string{"x", "y"}, 1), []string{"y"})

	// by=(x,y,z), key=(y,z): picks [y,z] — indexCovered passes (y,z both in by)
	tbl = newTable("xyz2", [][]string{{"y", "z"}}, [][]string{{"y", "z"}}, 100)
	test(tbl, LookupReq([]string{"x", "y", "z"}, 1), []string{"y", "z"})
}
