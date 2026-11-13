// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSchemaTableInfo(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	test := func(tran QueryTran, ntables, ncolumns, nindexes, nviews int) {
		t.Helper()
		check := func(table string, nrows int) {
			t.Helper()
			info := tran.GetInfo(table)
			a.This(info.Table).Is(table)
			a.This(info.Nrows).Is(nrows)
			a.This(info.Size).Is(nrows * 32)
		}
		check("tables", ntables)
		check("columns", ncolumns)
		check("indexes", nindexes)
		check("views", nviews)
	}
	empty := db.NewReadTran(AllPerms)
	test(empty, 6, 18, 6, 0)
	db.adm("create zz (a, b, c) key(a) index(c)")
	db.adm("alter zz create (G)")
	db.adm("alter zz drop (b)")
	db.adm("view foo = zz")
	db.adm("view bar = zz")
	rt := db.NewReadTran(AllPerms)
	test(rt, 7, 21, 8, 2)
	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	test(ut, 7, 21, 8, 2)
	DoAction(nil, ut, "insert {a: 1, c: 2} into zz")
	test(ut, 7, 21, 8, 2)
	ut.Abort()
	db.adm("drop zz")
	db.adm("drop foo")
	test(db.NewReadTran(AllPerms), 6, 18, 6, 1)
	test(rt, 7, 21, 8, 2)
	test(empty, 6, 18, 6, 0)
}

func TestSchemaTablesLookup(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a) key(a)")
	db.act("insert {a: 1} into zz")
	test := func(tran QueryTran) {
		t.Helper()
		tbl := NewTable(tran, "tables")
		a.This(tbl.Header().Columns).Is([]string{"table"})
		tbl.SetIndex([]string{"table"}, ReadMode)
		row := tbl.Lookup(nil, Sels{NewSel("table", Pack(SuStr("zz")))})
		a.This(row.GetVal(tbl.Header(), "table", nil, nil)).Is(SuStr("zz"))
		row = tbl.Lookup(nil, Sels{NewSel("table", Pack(SuStr("tables")))})
		a.This(row.GetVal(tbl.Header(), "table", nil, nil)).Is(SuStr("tables"))
		a.True(tbl.Lookup(nil, Sels{NewSel("table", Pack(SuStr("missing")))}).Eof())
		a.True(tbl.Lookup(nil, Sels{NewSel("table", Pack(One))}).Eof())

		// iteration over tables ordered by name
		all := []string{"columns", "dbstats", "history", "indexes", "tables", "views", "zz"}
		tbl.Rewind()
		for _, name := range all {
			iterRow := tbl.Get(nil, Next)
			a.This(ToStr(iterRow.GetVal(tbl.Header(), "table", nil, nil))).Is(name)
		}
		a.True(tbl.Get(nil, Next).Eof())
		tbl.Rewind()
		for i := len(all) - 1; i >= 0; i-- {
			a.This(ToStr(tbl.Get(nil, Prev).GetVal(tbl.Header(), "table", nil, nil))).Is(all[i])
		}
		a.True(tbl.Get(nil, Prev).Eof())

		q := ParseQuery("tables where table = 'zz'", tran, nil)
		q, _, _ = Setup(q, ReadMode, tran)
		w := q.(*Where)
		a.True(w.tbl != nil && w.idxSelBase != nil)
		a.This(q.Get(nil, Next).GetVal(q.Header(), "table", nil, nil)).Is(SuStr("zz"))
		a.True(q.Get(nil, Next).Eof())
	}
	rt := db.NewReadTran(AllPerms)
	test(rt)
	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	DoAction(nil, ut, "insert {a: 2} into zz")
	test(ut)
	test(rt)
	a.This(queryAll(db.Database, "tables where table = 'missing'")).Is("")
	a.This(queryAll(db.Database, "tables where table = 123")).Is("")
	a.This(queryAll(db.Database, "tables project table")).
		Is("table=columns | table=dbstats | table=history | table=indexes | table=tables | table=views | table=zz")
	a.This(queryAll(db.Database, "tables where table >= 'zz' project table")).Is("table=zz")
	a.This(queryAll(db.Database, "tables where table > 'zz' project table")).Is("")
	a.This(queryAll(db.Database, "tables where table = 'tables'")).Is("table=tables")
}

// TestSchemaTablesCursor verifies that schema tables cannot be used in a
// cursor, with or without a Where, since they are virtualized per transaction.
func TestSchemaTablesCursor(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a) key(a)")
	tran := db.NewReadTran(AllPerms)
	queries := []string{
		"tables",
		"tables where table is 'zz'",
		"tables where table = 'zz'",
		"columns where table is 'zz'",
		"indexes where table is 'zz'",
		"views where view_name is 'zz'",
	}
	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			q := ParseQuery(query, tran, nil)
			assert.T(t).This(func() { Setup(q, CursorMode, tran) }).
				Panics("invalid query")
		})
	}
}

func TestSchemaTablesIterRecords(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a) key(a)")
	db.adm("create aa (a) key(a)")
	db.act("insert {a: 1} into zz")
	test := func(tran QueryTran) {
		t.Helper()
		tbl := NewTable(tran, "tables")
		tbl.SetIndex([]string{"table"}, ReadMode)
		check := func(row Row, name string) {
			t.Helper()
			a.This(ToStr(row.GetVal(tbl.Header(), "table", nil, nil))).Is(name)
			a.This(row.Get(0).Count()).Is(1)
		}
		all := []string{"aa", "columns", "dbstats", "history", "indexes", "tables", "views", "zz"}
		for _, name := range all {
			check(tbl.Get(nil, Next), name)
		}
		a.True(tbl.Get(nil, Next).Eof())
		tbl.Rewind()
		for i := len(all) - 1; i >= 0; i-- {
			check(tbl.Get(nil, Prev), all[i])
		}
		a.True(tbl.Get(nil, Prev).Eof())
		rows := tbl.Simple(nil)
		a.This(len(rows)).Is(len(all))
		for i, name := range all {
			check(rows[i], name)
		}
		tbl.Select(Sels{NewSel("table", Pack(SuStr("zz")))})
		check(tbl.Get(nil, Next), "zz")
		a.True(tbl.Get(nil, Next).Eof())
		testRange := func(org, end string, names ...string) {
			t.Helper()
			tbl.SelectRaw(org, end)
			for _, name := range names {
				check(tbl.Get(nil, Next), name)
			}
			a.True(tbl.Get(nil, Next).Eof())
			tbl.Rewind()
			for i := len(names) - 1; i >= 0; i-- {
				check(tbl.Get(nil, Prev), names[i])
			}
			a.True(tbl.Get(nil, Prev).Eof())
		}
		testRange("", "\xff", all...)
		testRange(Pack(SuStr("aa")), Pack(SuStr("zz")), all[:len(all)-1]...)
		testRange(Pack(SuStr("aa"))+"\x00", "\xff", all[1:]...)
		testRange(Pack(SuStr("zz"))+"\x00", "\xff")
		testRange(Pack(One), "\xff", all...)
		testRange("", Pack(One))
		testRange("\xff", "\xff")
	}
	test(db.NewReadTran(AllPerms))
	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	DoAction(nil, ut, "insert {a: 2} into zz")
	test(ut)
}

func TestSchemaColumnsLookup(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a, b) key(a)")
	db.adm("alter zz create (G)")
	db.act("insert {a: 1, b: 2} into zz")
	test := func(tran QueryTran) {
		t.Helper()
		tbl := NewTable(tran, "columns")
		tbl.SetIndex([]string{"table", "column"}, ReadMode)
		lookup := func(table, col string) Row {
			return tbl.Lookup(nil, Sels{
				NewSel("table", Pack(SuStr(table))),
				NewSel("column", Pack(SuStr(col)))})
		}
		testCol := func(table, col string, field int) {
			t.Helper()
			row := lookup(table, col)
			a.This(row.GetVal(tbl.Header(), "table", nil, nil)).Is(SuStr(table))
			a.This(row.GetVal(tbl.Header(), "column", nil, nil)).Is(SuStr(col))
			a.This(ToInt(row.GetVal(tbl.Header(), "field", nil, nil))).Is(field)
		}
		testCol("zz", "a", 0)
		testCol("zz", "b", 1)
		// derived columns are reported uncapitalized with field -1
		testCol("zz", "g", -1)
		testCol("columns", "column", 1)
		testCol("tables", "table", 0)
		a.True(lookup("zz", "missing").Eof())
		a.True(lookup("missing", "a").Eof())
		a.True(tbl.Lookup(nil, Sels{
			NewSel("table", Pack(SuStr("zz"))),
			NewSel("column", Pack(One))}).Eof())
		// iteration over columns for zz ordered by {table, column}
		tbl.Select(Sels{NewSel("table", Pack(SuStr("zz")))})
		iterCol := func(dir Dir, col string, field int) {
			t.Helper()
			row := tbl.Get(nil, dir)
			a.This(ToStr(row.GetVal(tbl.Header(), "table", nil, nil))).Is("zz")
			a.This(ToStr(row.GetVal(tbl.Header(), "column", nil, nil))).Is(col)
			a.This(ToInt(row.GetVal(tbl.Header(), "field", nil, nil))).Is(field)
		}
		iterCol(Next, "a", 0)
		iterCol(Next, "b", 1)
		iterCol(Next, "g", -1)
		a.True(tbl.Get(nil, Next).Eof())
		tbl.Rewind()
		iterCol(Prev, "g", -1)
		iterCol(Prev, "b", 1)
		iterCol(Prev, "a", 0)
		a.True(tbl.Get(nil, Prev).Eof())

		q := ParseQuery("columns where table = 'zz' and column = 'a'", tran, nil)
		q, _, _ = Setup(q, ReadMode, tran)
		w := q.(*Where)
		a.True(w.tbl != nil && w.idxSelBase != nil)
		a.This(ToInt(q.Get(nil, Next).GetVal(q.Header(), "field", nil, nil))).Is(0)
		a.True(q.Get(nil, Next).Eof())
	}
	rt := db.NewReadTran(AllPerms)
	test(rt)
	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	DoAction(nil, ut, "insert {a: 2, b: 3} into zz")
	test(ut)
	test(rt)
	a.This(queryAll(db.Database, "columns where table = 'zz' and column = 'b'")).
		Is("column=b field=1 table=zz")
	a.This(queryAll(db.Database, "columns where table = 'missing' and column = 'x'")).Is("")
	a.This(queryAll(db.Database, "columns where table = 123 and column = 'a'")).Is("")
}

func TestSchemaIndexesLookup(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create hdr (a) key(a)")
	db.adm("create zz (a, b, c) key(a) index(b) in hdr(a) index unique(c)")
	db.adm("create kk (a) key()")
	test := func(tran QueryTran) {
		t.Helper()
		tbl := NewTable(tran, "indexes")
		tbl.SetIndex([]string{"table", "columns"}, ReadMode)
		hdr := tbl.Header()
		lookup := func(table, cols string) Row {
			return tbl.Lookup(nil, Sels{
				NewSel("table", Pack(SuStr(table))),
				NewSel("columns", Pack(SuStr(cols)))})
		}
		val := func(row Row, col string) Value { return row.GetVal(hdr, col, nil, nil) }
		// key
		row := lookup("zz", "a")
		a.This(val(row, "table")).Is(SuStr("zz"))
		a.This(val(row, "columns")).Is(SuStr("a"))
		a.This(val(row, "key")).Is(True)
		a.True(IsEmptyStr(val(row, "fktable")))
		// foreign key index
		row = lookup("zz", "b")
		a.This(val(row, "key")).Is(False)
		a.This(val(row, "fktable")).Is(SuStr("hdr"))
		a.This(val(row, "fkcolumns")).Is(SuStr("a"))
		a.This(ToInt(val(row, "fkmode"))).Is(0)
		// unique index
		a.This(val(lookup("zz", "c"), "key")).Is(SuStr("u"))
		// schema table index
		a.This(val(lookup("indexes", "table,columns"), "key")).Is(True)
		// missing
		a.True(lookup("zz", "nosuch").Eof())
		a.True(lookup("missing", "a").Eof())
		a.True(tbl.Lookup(nil, Sels{
			NewSel("table", Pack(SuStr("zz"))),
			NewSel("columns", Pack(One))}).Eof())
		// key() index has empty columns; a one-part key must still find it
		row = lookup("kk", "")
		a.This(val(row, "table")).Is(SuStr("kk"))
		a.This(val(row, "columns")).Is(SuStr(""))
		a.This(val(row, "key")).Is(True)
		// a non-string columns must not match the empty key() columns
		a.True(tbl.Lookup(nil, Sels{
			NewSel("table", Pack(SuStr("kk"))),
			NewSel("columns", Pack(One))}).Eof())
		// iteration over indexes ordered by {table, columns}
		allIdx := [][2]string{
			{"columns", "table,column"},
			{"dbstats", "table,column"},
			{"hdr", "a"},
			{"history", "asof"},
			{"indexes", "table,columns"},
			{"kk", ""},
			{"tables", "table"},
			{"views", "view_name"},
			{"zz", "a"},
			{"zz", "b"},
			{"zz", "c"},
		}
		tbl.Rewind()
		iterIdx := func(dir Dir, table, cols string) {
			t.Helper()
			row := tbl.Get(nil, dir)
			a.This(ToStr(row.GetVal(tbl.Header(), "table", nil, nil))).Is(table)
			a.This(ToStr(row.GetVal(tbl.Header(), "columns", nil, nil))).Is(cols)
		}
		for _, ix := range allIdx {
			iterIdx(Next, ix[0], ix[1])
		}
		a.True(tbl.Get(nil, Next).Eof())
		tbl.Rewind()
		for i := len(allIdx) - 1; i >= 0; i-- {
			iterIdx(Prev, allIdx[i][0], allIdx[i][1])
		}
		a.True(tbl.Get(nil, Prev).Eof())

		q := ParseQuery("indexes where table = 'zz' and columns = 'b'", tran, nil)
		q, _, _ = Setup(q, ReadMode, tran)
		w := q.(*Where)
		a.True(w.tbl != nil && w.idxSelBase != nil)
		a.This(ToStr(q.Get(nil, Next).GetVal(q.Header(), "fktable", nil, nil))).Is("hdr")
		a.True(q.Get(nil, Next).Eof())
	}
	rt := db.NewReadTran(AllPerms)
	test(rt)
	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	test(ut)
	test(rt)
	a.This(queryAll(db.Database, "indexes where table = 'zz' and columns = 'c'")).
		Is("columns=c key=u table=zz")
	a.This(queryAll(db.Database, "indexes where table = 'kk' and columns = ''")).
		Is("key=true table=kk")
	a.This(queryAll(db.Database, "indexes where table = 'kk' and columns = 123")).Is("")
	a.This(queryAll(db.Database, "indexes where table = 'missing' and columns = 'x'")).Is("")
	a.This(queryAll(db.Database, "indexes where table = 123 and columns = 'a'")).Is("")
}

func TestSchemaViewsLookup(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("view foo = bar baz")
	test := func(tran QueryTran) {
		t.Helper()
		tbl := NewTable(tran, "views")
		tbl.SetIndex([]string{"view_name"}, ReadMode)
		hdr := tbl.Header()
		lookup := func(name string) Row {
			return tbl.Lookup(nil, Sels{NewSel("view_name", Pack(SuStr(name)))})
		}
		row := lookup("foo")
		a.This(row.GetVal(hdr, "view_name", nil, nil)).Is(SuStr("foo"))
		a.This(row.GetVal(hdr, "view_definition", nil, nil)).Is(SuStr("bar baz"))
		a.True(lookup("missing").Eof())
		a.True(tbl.Lookup(nil, Sels{NewSel("view_name", Pack(One))}).Eof())
		// iteration over views ordered by name
		tbl.Rewind()
		iterView := func(dir Dir) {
			t.Helper()
			row := tbl.Get(nil, dir)
			a.This(ToStr(row.GetVal(hdr, "view_name", nil, nil))).Is("foo")
			a.This(ToStr(row.GetVal(hdr, "view_definition", nil, nil))).Is("bar baz")
		}
		iterView(Next)
		a.True(tbl.Get(nil, Next).Eof())
		tbl.Rewind()
		iterView(Prev)
		a.True(tbl.Get(nil, Prev).Eof())

		q := ParseQuery("views where view_name = 'foo'", tran, nil)
		q, _, _ = Setup(q, ReadMode, tran)
		w := q.(*Where)
		a.True(w.tbl != nil && w.idxSelBase != nil)
		a.This(ToStr(q.Get(nil, Next).GetVal(q.Header(), "view_definition", nil, nil))).Is("bar baz")
		a.True(q.Get(nil, Next).Eof())
	}
	rt := db.NewReadTran(AllPerms)
	test(rt)
	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	test(ut)
	test(rt)
	a.This(queryAll(db.Database, "views where view_name = 'foo'")).
		Is("view_definition=bar baz view_name=foo")
	a.This(queryAll(db.Database, "views where view_name = 'missing'")).Is("")
	a.This(queryAll(db.Database, "views where view_name = 123")).Is("")
}

// TestSchemaTablesNotUpdateable verifies that the schema tables
// reject output, update, and delete.
func TestSchemaTablesNotUpdateable(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a) key(a)")
	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	schemaTables :=
		[]string{"tables", "columns", "indexes", "views", "history", "dbstats"}
	for _, name := range schemaTables {
		t.Run(name, func(t *testing.T) {
			a := assert.T(t)
			q := NewTable(ut, name)
			a.This(q.Updateable()).Is("")
			a.This(func() { q.Output(nil, new(RecordBuilder).Build()) }).
				Panics("can't output to schema table")
		})
	}
	a.This(func() { DoAction(nil, ut, "update tables set table = 0") }).
		Panics("update: query not updateable")
	a.This(func() { DoAction(nil, ut, "delete tables") }).
		Panics("delete: query not updateable")
}

func TestSchemaColumnsSkipScan(t *testing.T) {
	testSchemaSkipScan(t, "columns", "column")
}

func TestSchemaIndexesSkipScan(t *testing.T) {
	testSchemaSkipScan(t, "indexes", "columns")
}

func testSchemaSkipScan(t *testing.T, table, suffix string) {
	t.Helper()
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create aa (a, b, d) key(a) index(b) index(d)")
	db.adm("create ab (a, c, d) key(a) index(c) index(d)")
	db.adm("create ba (a, b, c) key(a) index(b) index(c)")
	db.adm("create ca (a, b, d) key(a) index(b) index(d)")
	test := func(where string, expected ...string) {
		t.Helper()
		tran := db.NewReadTran(AllPerms)
		q := ParseQuery(table+" where "+where, tran, nil)
		q, _, _ = Setup(q, ReadMode, tran)
		w := q.(*Where)
		a.True(w.tbl != nil && w.idxSelBase != nil)
		a.This(w.idxSelBase.skipStart).Is(1)
		a.This(w.idxSelBase.skipLen).Is(1)
		check := func(dir Dir, expected string) {
			t.Helper()
			row := q.Get(nil, dir)
			a.False(row.Eof())
			a.This(ToStr(row.GetVal(q.Header(), "table", nil, nil)) + "/" +
				ToStr(row.GetVal(q.Header(), suffix, nil, nil))).Is(expected)
		}
		for _, expected := range expected {
			check(Next, expected)
		}
		a.True(q.Get(nil, Next).Eof())
		q.Rewind()
		for i := len(expected) - 1; i >= 0; i-- {
			check(Prev, expected[i])
		}
		a.True(q.Get(nil, Prev).Eof())
	}
	test(suffix+" = 'b'", "aa/b", "ba/b", "ca/b")
	rangeExpected := []string{"aa/b", "ab/c", "ba/b", "ba/c", "ca/b"}
	if table == "columns" {
		rangeExpected = append(rangeExpected,
			"columns/column", "dbstats/column", "dbstats/common", "indexes/columns")
	}
	test(suffix+" >= 'b' and "+suffix+" < 'd'", rangeExpected...)
	test("table >= 'ab' and table < 'ca' and "+suffix+" = 'b'", "ba/b")
	test("table >= 'ab' and table < 'ca' and "+suffix+" >= 'b' and "+suffix+" < 'd'",
		"ab/c", "ba/b", "ba/c")
	test(suffix + " = 'missing'")
	test("table >= 'ab' and table < 'ca' and " + suffix + " > 'z'")
}

// TestSchemaRangeFrac verifies that RangeFrac for the "_" schema tables
// counts the rows within the range instead of using btree statistics.
func TestSchemaRangeFrac(t *testing.T) {
	a := assert.T(t)
	db := heapDb()
	defer db.Close()
	db.adm("create aa (a, b) key(a) index(b)")
	db.adm("create bb (a) key(a)")
	db.adm("view foo = aa")
	rt := db.NewReadTran(AllPerms)
	frac := func(table, org, end string) float64 {
		return rt.RangeFrac(table, 0, org, end)
	}
	for _, table := range []string{"tables", "columns", "indexes", "views"} {
		a.This(frac(table, ixkey.Min, ixkey.Max)).Is(1)
	}
	aa, bb := Pack(SuStr("aa")), Pack(SuStr("bb"))
	a.This(frac("tables", aa, bb)).Is(1.0 / 8)
	a.This(frac("tables", aa, aa+"\x00")).Is(1.0 / 8)
	a.This(frac("tables", aa, aa)).Is(0)
	a.This(frac("tables", bb, ixkey.Max)).Is(7.0 / 8)
	a.This(frac("tables", "\xff", ixkey.Max)).Is(0)
	// composite index keys are encoded and directly comparable
	caa, cbb := ixkey.CompKey(aa), ixkey.CompKey(bb)
	a.This(frac("columns", caa, cbb)).Is(2.0 / 21)
	a.This(frac("indexes", caa, cbb)).Is(2.0 / 9)
	caaColA := ixkey.CompKey(aa, Pack(SuStr("a")))
	caaColB := ixkey.CompKey(aa, Pack(SuStr("b")))
	a.This(frac("columns", caaColA, caaColB)).Is(1.0 / 21)
	a.This(frac("columns", caaColB, ixkey.Max)).Is(20.0 / 21)
	// views has a single row
	a.This(frac("views", Pack(SuStr("foo")), Pack(SuStr("foo")))).Is(0)
	// history is only an estimate, dbstats is empty
	a.This(frac("history", ixkey.Min, ixkey.Max)).Is(.5)
	a.This(frac("dbstats", ixkey.Min, ixkey.Max)).Is(0)
}
