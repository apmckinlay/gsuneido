// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTableOptimize(t *testing.T) {
	assert := assert.T(t)
	tran := testTran{}
	test := func(table string, req Require, expected []string) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, app := tbl.optimize(ReadMode, req)
		assert.True(f+v < impossible)
		assert.This(app.(tableApproach).index).Is(expected)
	}
	assertImpossible := func(table string, req Require) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, app := tbl.optimize(ReadMode, req)
		assert.False(f+v < impossible)
		assert.This(app).Is(nil)
	}

	// table with single index: key on {a}
	// indexes: [[a]], allKeys: [[a]]

	// ReqNone
	test("table", NoneReq(1), []string{"a"})

	// ReqOrder — match
	test("table", OrderReq([]string{"a"}, 1), []string{"a"})

	// ReqOrder — no match (b not a prefix of {a})
	assertImpossible("table", OrderReq([]string{"b"}, 1))

	// ReqGroup — match
	test("table", GroupReq([]string{"a"}, 1, 1), []string{"a"})

	// ReqGroup — b not found in any index
	assertImpossible("table", GroupReq([]string{"b"}, 1, 1))

	// ReqUnique — {a} is a physical index
	test("table", UniqueReq([]string{"a"}, 1), []string{"a"})

	// ReqUnique — {b} is not a physical index, no lookup-eligible index has {b}
	assertImpossible("table", UniqueReq([]string{"b"}, 1))

	// supplier: key on {supplier}, index on {city}
	// indexes: [[supplier], [city,supplier]], allKeys: [[supplier]]

	// ReqNone — primary index
	test("supplier", NoneReq(1), []string{"supplier"})

	// ReqOrder — match key
	test("supplier", OrderReq([]string{"supplier"}, 1), []string{"supplier"})

	// ReqOrder — match non-key index
	test("supplier", OrderReq([]string{"city", "supplier"}, 1), []string{"city", "supplier"})

	// ReqOrder — no match (name not in any index)
	assertImpossible("supplier", OrderReq([]string{"name"}, 1))

	// ReqOrder — city is a prefix of [city,supplier], so indexFor finds it
	test("supplier", OrderReq([]string{"city"}, 1), []string{"city", "supplier"})

	// ReqGroup — {city} is grouped by [city,supplier]
	test("supplier", GroupReq([]string{"city"}, 1, 1), []string{"city", "supplier"})

	// ReqGroup — {name} not in any index
	assertImpossible("supplier", GroupReq([]string{"name"}, 1, 1))

	// ReqUnique — {supplier} is a physical index
	test("supplier", UniqueReq([]string{"supplier"}, 1), []string{"supplier"})

	// ReqUnique — {supplier} is the shortest key
	test("supplier", UniqueReq([]string{"city", "supplier"}, 1), []string{"supplier"})

	// ReqUnique — {city} is not a physical index and no index is both
	// lookup-eligible (contains a key as a prefix) AND grouped by {city}.
	// [supplier] is lookup-eligible but not grouped by {city}.
	// [city,supplier] is grouped by {city} but NOT lookup-eligible
	// (key {supplier} is not a prefix of [city,supplier]).
	assertImpossible("supplier", UniqueReq([]string{"city"}, 1))

	// abc: key on {b}, index on {a}, key on {c}
	// indexes: [[b], [a,b], [c]], allKeys: [[b], [c]]

	// ReqNone — primary index {b}
	test("abc", NoneReq(1), []string{"b"})

	// ReqGroup — {a} is grouped by [a,b]
	test("abc", GroupReq([]string{"a"}, 1, 1), []string{"a", "b"})

	// ReqGroup — {b} is grouped by [b] (cheapest)
	test("abc", GroupReq([]string{"b"}, 1, 1), []string{"b"})

	// ReqGroup — {c} is grouped by [c]
	test("abc", GroupReq([]string{"c"}, 1, 1), []string{"c"})

	// ReqUnique — {b} is a physical index (and a key)
	test("abc", UniqueReq([]string{"b"}, 1), []string{"b"})

	// ReqUnique — {c} is a physical index (and a key)
	test("abc", UniqueReq([]string{"c"}, 1), []string{"c"})

	// ReqUnique — {a} is not a physical index. [a,b] is not
	// lookup-eligible (key {b} is not a prefix, key {c} is not in index).
	assertImpossible("abc", UniqueReq([]string{"a"}, 1))

	// comp: key on {a,b,c} — single index [a,b,c], keys [[a,b,c]]
	// ReqUnique {b,a} is not a physical index.
	// [a,b,c] is lookup-eligible and has {b,a} grouped,
	// but is not indexCovered (c is not in {b,a}).
	assertImpossible("comp", UniqueReq([]string{"b", "a"}, 1))

	// singleton: table with empty key — all req types return indexes[0]
	// Set up manually (no SetTran since _singleton_ is not in testSchemas).
	singleton := &Table{}
	singleton.indexes = [][]string{{"x"}}
	singleton.singleton = true
	singleton.info = &meta.Info{Nrows: 1, Size: 100}
	for _, req := range []Require{
		NoneReq(1),
		OrderReq([]string{"x"}, 1),
		GroupReq([]string{"x"}, 1, 1),
		UniqueReq([]string{"x"}, 1),
	} {
		f, v, app := singleton.optimize(ReadMode, req)
		assert.Msg(req).True(f+v < impossible)
		assert.Msg(req).This(app.(tableApproach).index).Is([]string{"x"})
	}
}

func TestTableOptimize_UniqueReq(t *testing.T) {
	newTable := func(indexes, keys [][]string, nrows int) *Table {
		tbl := &Table{name: "foo"}
		tbl.indexes = indexes
		tbl.allKeys = keys
		tbl.info = &meta.Info{Nrows: nrows, Size: int64(nrows * 100)}
		return tbl
	}
	test := func(tbl *Table, req Require, expected []string) {
		t.Helper()
		f, v, app := tbl.optimize(ReadMode, req)
		assert.True(f+v < impossible)
		assert.This(app.(tableApproach).index).Is(expected)
	}
	fail := func(tbl *Table, req Require) {
		t.Helper()
		f, v, app := tbl.optimize(ReadMode, req)
		assert.False(f+v < impossible)
		assert.This(app).Is(nil)
	}
	var tbl *Table

	// index(x) key(x) - req(x)
	tbl = newTable([][]string{{"x"}}, [][]string{{"x"}}, 100)
	test(tbl, UniqueReq([]string{"x"}, 1), []string{"x"})

	// index(x,y) key(x,y) - req(x,y)
	tbl = newTable([][]string{{"x", "y"}}, [][]string{{"x", "y"}}, 100)
	test(tbl, UniqueReq([]string{"x", "y"}, 1), []string{"x", "y"})

	// index(y) key(y) - req(x,y) - FAIL
	tbl = newTable([][]string{{"y"}}, [][]string{{"y"}}, 100)
	test(tbl, UniqueReq([]string{"x", "y"}, 1), []string{"y"})

	// index(x,y,z) key(x,y,z) - req(x,y) - FAIL
	tbl = newTable([][]string{{"x", "y", "z"}}, [][]string{{"x", "y", "z"}}, 100)
	fail(tbl, UniqueReq([]string{"x", "y"}, 1))

	// index(x,c), key=(x) - req(x,y) - FAIL
	tbl = newTable([][]string{{"x", "c"}}, [][]string{{"x"}}, 100)
	fail(tbl, UniqueReq([]string{"x", "y"}, 1))

	// index(y,z) key(y,z) - req(x,y,z) - FAIL
	tbl = newTable([][]string{{"y", "z"}}, [][]string{{"y", "z"}}, 100)
	test(tbl, UniqueReq([]string{"x", "y", "z"}, 1), []string{"y", "z"})
}
