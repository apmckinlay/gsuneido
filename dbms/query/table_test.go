// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type tableSetTran struct {
	QueryTran
	schema *Schema
	info   *meta.Info
	checks int
	deny   bool
}

func (tr *tableSetTran) CheckPerm(string, core.TablePerm) {
	tr.checks++
	if tr.deny {
		panic("not authorized")
	}
}

func (tr *tableSetTran) GetSchema(string) *Schema  { return tr.schema }
func (tr *tableSetTran) GetInfo(string) *meta.Info { return tr.info }

func TestTable_scanBusy(t *testing.T) {
	info := &meta.Info{Nrows: 1000, Size: smallTable}
	tran := &scanTran{QueryTran: &tableSetTran{
		schema: &Schema{Columns: []string{"a", "b", "c"}, Indexes: []Index{
			{Mode: 'k', Columns: []string{"a"}, Fields: []string{"a"}},
			{Mode: 'i', Columns: []string{"b"}, Fields: []string{"b", "a"}},
			{Mode: 'i', Columns: []string{"b", "c"}, Fields: []string{"b", "c", "a"}},
		}},
		info: info,
	}, scans: make(map[int]int)}
	tbl := NewTable(tran, "t")
	test := func(req Require, weight int, indexes ...int) {
		t.Helper()
		clear(tran.scans)
		tbl.optimize(CursorMode, req)
		a := assert.T(t).Msg(req)
		a.This(len(tran.scans)).Is(len(indexes))
		for _, i := range indexes {
			a.This(tran.table).Is("t")
			a.This(tran.scans[i]).Is(weight)
		}
	}
	test(OrderReq([]string{"b"}, 1), 1000, 1, 2)
	test(GroupReq([]string{"a", "b"}, 1, 1), 1000, 1)
	test(OrderReq([]string{"b"}, .25), 250, 1, 2)
	test(OrderReq([]string{"b"}, .002), 0) // two rows
	test(OrderReq([]string{"b"}, .003), 3, 1, 2)
	test(OrderReq([]string{"c"}, 1), 0) // no matching index
	test(NoneReq(1), 0)
	test(UniqueReq([]string{"a"}, 1), 0)
	info.Size = smallTable - 1
	test(OrderReq([]string{"b"}, 1), 0)
	info.Size = smallTable
	tbl.indexes = tbl.indexes[:1]
	test(OrderReq([]string{"a"}, 1), 0)
}

func TestTableGetFilter(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create filtertest (a) key(a)")
	db.act("insert {a: 1} into filtertest")
	db.act("insert {a: 2} into filtertest")
	tran := db.NewReadTran(core.AllPerms)
	defer tran.Complete()
	tbl := NewTable(tran, "filtertest")
	tbl.SetIndex([]string{"a"}, ReadMode)
	test := func(dir core.Dir) {
		t.Helper()
		tbl.Rewind()
		first := tbl.Get(nil, dir)
		off := tbl.iter.CurOff()
		offFilter := func(o uint64) bool { return o == off }
		tbl.Rewind()
		a.This(tbl.GetFilter(dir, offFilter, nil)).Is(first)
		a.True(tbl.GetFilter(dir, offFilter, nil).Eof())

		tbl.Rewind()
		keyCalls := 0
		keyFilter := func(_ string, o uint64) bool {
			keyCalls++
			a.This(o).Is(off)
			return false
		}
		a.True(tbl.GetFilter(dir, offFilter, keyFilter).Eof())
		a.This(keyCalls).Is(1)
	}
	test(core.Next)
	test(core.Prev)
}

func TestTableSetTran(t *testing.T) {
	a := assert.T(t)
	tr := &tableSetTran{
		schema: &Schema{Columns: []string{"a"},
			Indexes: []Index{{Mode: 'k'}}},
		info: &meta.Info{Nrows: 1, Size: 100},
	}
	tbl := &Table{name: "tmp"}
	tbl.SetTran(tr)
	header := tbl.header
	a.This(tr.checks).Is(1)
	a.True(tbl.singleton)
	a.This(tbl.rowSize()).Is(100)

	check := func(info *meta.Info, rowSize int) {
		t.Helper()
		tr.info = info
		tbl.SetTran(tr)
		a.This(tr.checks).Is(1)
		a.True(tbl.header == header)
		a.True(tbl.info == info)
		a.This(tbl.rowSize()).Is(rowSize)
	}
	check(&meta.Info{Nrows: 2, Size: 300}, 150)
	tr.info.Size = 400
	check(tr.info, 200)
	tr.info.Nrows = 0
	check(tr.info, 0)

	tr2 := &tableSetTran{schema: tr.schema, info: tr.info, deny: true}
	a.This(func() { tbl.SetTran(tr2) }).Panics("not authorized")
	a.True(tbl.tran == tr)
	tr2.deny = false
	tbl.SetTran(tr2)
	a.This(tr2.checks).Is(2)
	a.True(tbl.tran == tr2)
	a.True(tbl.header == header)

	tr2.schema = &Schema{Columns: []string{"a", "b"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a"}, Fields: []string{"a"}}}}
	tbl.SetTran(tr2)
	a.This(tr2.checks).Is(2)
	a.True(tbl.header != header)
	a.This(tbl.header.Columns).Is([]string{"a", "b"})
	a.This(tbl.keys).Is([][]string{{"a"}})
	a.False(tbl.singleton)

	tr2.schema = nil
	a.This(func() { tbl.SetTran(tr2) }).Panics("nonexistent table: tmp")
}

func TestTableOptimize(t *testing.T) {
	assert := assert.T(t)
	tran := testTran{}
	test := func(table string, req Require, expected []string) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, ap := tbl.optimize(ReadMode, req)
		assert.True(f+v < impossible)
		assert.This(ap.(tableApproach).index).Is(expected)
	}
	assertImpossible := func(table string, req Require) {
		t.Helper()
		tbl := &Table{name: table}
		tbl.SetTran(tran)
		f, v, ap := tbl.optimize(ReadMode, req)
		assert.False(f+v < impossible)
		assert.This(ap).Is(nil)
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

	// ReqNone — physical (first) index
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
		f, v, ap := singleton.optimize(ReadMode, req)
		assert.Msg(req).True(f+v < impossible)
		assert.Msg(req).This(ap.(tableApproach).index).Is([]string{"x"})
	}
}

func TestTableOptimize_UniqueReq(t *testing.T) {
	newTable := func(indexes, keys [][]string, nrows int) *Table {
		tbl := &Table{name: "foo"}
		tbl.indexes = indexes
		tbl.allKeys = keys
		tbl.info = &meta.Info{Nrows: nrows, Size: nrows * 100}
		schemaIndexes := make([]Index, len(indexes))
		for i, idx := range indexes {
			schemaIndexes[i] = Index{Columns: idx}
		}
		tbl.schema = &Schema{Table: "foo", Indexes: schemaIndexes}
		return tbl
	}
	test := func(tbl *Table, req Require, expected []string) {
		t.Helper()
		f, v, ap := tbl.optimize(ReadMode, req)
		assert.True(f+v < impossible)
		assert.This(ap.(tableApproach).index).Is(expected)
	}
	fail := func(tbl *Table, req Require) {
		t.Helper()
		f, v, ap := tbl.optimize(ReadMode, req)
		assert.False(f+v < impossible)
		assert.This(ap).Is(nil)
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
