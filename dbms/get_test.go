// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"strings"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestGetQuery(t *testing.T) {
	test := func(ob *SuObject, expected string) {
		t.Helper()
		assert.T(t).This(getQuery(ob)).Is(expected)
	}

	// From list position 0
	ob := &SuObject{}
	ob.Add(SuStr("mytable"))
	test(ob, "mytable")

	ob2 := &SuObject{}
	ob2.Add(SuStr("customer join orders"))
	test(ob2, "customer join orders")

	// From named argument
	obj := &SuObject{}
	obj.Set(SuStr("query"), SuStr("table"))
	test(obj, "table")

	// Empty object
	test(&SuObject{}, "")

	// List takes precedence
	obj2 := &SuObject{}
	obj2.Add(SuStr("list_query"))
	obj2.Set(SuStr("query"), SuStr("named_query"))
	test(obj2, "list_query")
}

func TestGetWhere(t *testing.T) {
	// Empty object
	assert.T(t).This(getWhere(&SuObject{})).Is("")

	// Single field
	obj := &SuObject{}
	obj.Set(SuStr("name"), SuStr("John"))
	assert.T(t).This(getWhere(obj)).Is(`where name is "John"`)

	// With query argument (should be excluded)
	obj3 := &SuObject{}
	obj3.Set(SuStr("query"), SuStr("customer"))
	obj3.Set(SuStr("city"), SuStr("Calgary"))
	assert.T(t).This(getWhere(obj3)).Is(`where city is "Calgary"`)

	// Multiple fields - check structure
	obj4 := &SuObject{}
	obj4.Set(SuStr("a"), SuInt(1))
	obj4.Set(SuStr("b"), SuInt(2))
	w := getWhere(obj4)
	assert.T(t).True(strings.HasPrefix(w, "where "))
	assert.T(t).True(strings.Contains(w, "\nand "))
}

func TestFindUnique(t *testing.T) {
	indexes := [][]string{
		{"x"},
		{"a", "b"},
	}

	test := func(sels Sels, expected []string) {
		t.Helper()
		result := findUnique(indexes, sels)
		assert.T(t).This(result).Is(expected)
	}

	// non-empty value matches
	test(selVals("x", "5"), []string{"x"})

	// empty value does not match (Lookup would not be unique)
	test(selVals("x", ""), nil)

	// not all index columns selected
	test(selVals("a", "1"), nil)

	// exact match on multi-column index
	test(selVals("a", "1", "b", "2"), []string{"a", "b"})

	// superset of index columns
	test(selVals("a", "1", "b", "2", "z", "3"), []string{"a", "b"})

	// multi-column with one empty, one non-empty
	test(selVals("a", "", "b", "2"), []string{"a", "b"})

	// all values empty does not match
	test(selVals("a", "", "b", ""), nil)

	// no match
	test(selVals("z", "1"), nil)
}

func TestFindKey(t *testing.T) {
	keys := [][]string{
		{"id"},
		{"name", "date"},
		{"a", "b", "c"},
	}

	test := func(sels Sels, expected []string) {
		t.Helper()
		result := findKey(keys, sels)
		assert.T(t).This(result).Is(expected)
	}

	// Exact match
	test(selCols("id"), []string{"id"})

	// Superset match (sels must contain all of key)
	test(selCols("name", "date", "extra"), []string{"name", "date"})
	test(selCols("name", "date"), []string{"name", "date"})

	// No match (sels doesn't contain all of key)
	test(selCols("name"), nil)
	test(selCols("date"), nil)
	test(selCols("xyz"), nil)
	test(Sels{}, nil)

	// Multiple keys possible (should return first match)
	test(selCols("a", "b", "c"), []string{"a", "b", "c"})
}

func selCols(cols ...string) Sels {
	sels := make(Sels, len(cols))
	for i, col := range cols {
		sels[i] = qry.NewSel(col, "")
	}
	return sels
}

func selVals(colvals ...string) Sels {
	sels := make(Sels, len(colvals)/2)
	for i := 0; i < len(colvals); i += 2 {
		sels[i/2] = qry.NewSel(colvals[i], Pack(SuStr(colvals[i+1])))
	}
	return sels
}

// TestGetOnlyUniqueIndex verifies that get Only uses a unique ('u') index
// for a Lookup when the selected values are non-empty.
func TestGetOnlyUniqueIndex(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	qry.DoAdmin(db, "create tmp (k, u, data) key(k) index unique(u)", nil)
	act := func(action string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		qry.DoAction(nil, ut, action)
	}
	act("insert { k: 1, u: '', data: 'first' } into tmp")
	act("insert { k: 2, u: '', data: 'second' } into tmp")
	act("insert { k: 3, u: 'x', data: 'third' } into tmp")

	tran := db.NewReadTran()
	defer tran.Complete()
	th := &Thread{}
	ob := &SuObject{}
	ob.Set(SuStr("query"), SuStr("tmp"))
	ob.Set(SuStr("u"), SuStr("x"))
	row, hdr, _ := get(th, tran, ob, Only)
	assert.T(t).Msg("get Only u=x").
		This(AsStr(row.GetVal(hdr, "u", nil, nil))).Is("x")
	assert.T(t).This(AsStr(row.GetVal(hdr, "k", nil, nil))).Is("3")

	// verify it uses the unique index for the Lookup
	tbl := qry.NewTable(tran, "tmp").(*qry.Table)
	sels := Sels{qry.NewSel("u", Pack(SuStr("x")))}
	single, strat, getfn := getIndex(th, tran, tbl, sels, Only)
	assert.T(t).Msg("single").True(single)
	assert.T(t).Msg("strategy").
		True(strings.Contains(strat, "unique index"))
	assert.T(t).Msg("lookup row").
		This(AsStr(getfn().GetVal(tbl.Header(), "data", nil, nil))).
		Is("third")
}

// TestGetOnlyUniqueIndexEmpty verifies that get Only with an empty value
// on a unique index does NOT use a Lookup (it falls back to a scan),
// and correctly reports multiple matching records as not unique.
func TestGetOnlyUniqueIndexEmpty(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	qry.DoAdmin(db, "create tmp (k, u, data) key(k) index unique(u)", nil)
	act := func(action string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		qry.DoAction(nil, ut, action)
	}
	act("insert { k: 1, u: '', data: 'first' } into tmp")
	act("insert { k: 2, u: '', data: 'second' } into tmp")

	tran := db.NewReadTran()
	defer tran.Complete()
	th := &Thread{}
	tbl := qry.NewTable(tran, "tmp").(*qry.Table)
	sels := Sels{qry.NewSel("u", Pack(SuStr("")))}
	single, strat, getfn := getIndex(th, tran, tbl, sels, Only)
	assert.T(t).Msg("not single").True(!single)
	assert.T(t).Msg("strategy").
		True(!strings.Contains(strat, "unique index"))
	assert.T(t).Msg("lookup row").This(getfn()).Isnt(nil)
}

func TestFindAll(t *testing.T) {
	indexes := [][]string{
		{"a", "b", "c"},
		{"x", "y"},
		{"p", "q", "r"},
	}

	test := func(sels Sels, expected []string) {
		t.Helper()
		result := findAll(indexes, sels)
		assert.T(t).This(result).Is(expected)
	}

	// Exact match on first index
	test(selCols("a", "b", "c"), []string{"a", "b", "c"})

	// Prefix match
	test(selCols("a", "b"), []string{"a", "b", "c"})
	test(selCols("a"), []string{"a", "b", "c"})

	// Match on second index
	test(selCols("x", "y"), []string{"x", "y"})
	test(selCols("x"), []string{"x", "y"})

	// No match
	test(selCols("z"), nil)
	test(selCols("a", "x"), nil)
	test(selCols("b", "c"), nil) // must start from beginning
}

func TestHasPrefix(t *testing.T) {
	test := func(idx []string, sels Sels, expected bool) {
		t.Helper()
		result := hasPrefix(idx, sels)
		assert.T(t).This(result).Is(expected)
	}

	test([]string{"a", "b"}, selCols("a"), true)
	test([]string{"a", "b"}, selCols("b"), false)
	test([]string{"a", "b"}, selCols("a", "b"), true)
	test([]string{"a", "b"}, selCols("x"), false)
	test([]string{"a"}, selCols("a"), true)
}

func TestUsableIndexes(t *testing.T) {
	indexes := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"c", "d"},
		{"a", "x"},
	}

	test := func(sels Sels, expected [][]string) {
		t.Helper()
		result := usableIndexes(indexes, sels)
		assert.T(t).This(result).Is(expected)
	}

	// Fields match first elements (any index whose first field is in sels)
	test(selCols("a"), [][]string{{"a", "b"}, {"a", "x"}})
	test(selCols("b"), [][]string{{"b", "c"}})
	test(selCols("c"), [][]string{{"c", "d"}})

	// Multiple fields - includes indexes starting with "a" or "b"
	test(selCols("a", "b"), [][]string{{"a", "b"}, {"b", "c"}, {"a", "x"}})

	// No matches
	test(selCols("x"), nil)
	test(selCols("d"), nil)
}

func TestFormatFieldsVals(t *testing.T) {
	test := func(sels Sels, expected string) {
		t.Helper()
		result := formatFieldsVals(sels)
		assert.T(t).This(result).Is(expected)
	}

	// Single field
	test(Sels{qry.NewSel("name", Pack(SuStr("John")))}, `name: "John"`)

	// Multiple fields
	test(Sels{qry.NewSel("name", Pack(SuStr("John"))), qry.NewSel("age", Pack(SuInt(30)))}, `name: "John", age: 30`)

	// Empty
	test(Sels{}, "")
}

func TestGetSchema(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	qry.DoAdmin(db, "create tmp (k, u, data, Foo) key(k) index unique(u)", nil)
	qry.DoAdmin(db, "view myview = tmp", nil)

	tran := db.NewReadTran()
	defer tran.Complete()
	th := &Thread{}

	// lookup a schema table by its key, returning the row and header
	test := func(table string, sels ...[2]string) (Row, *Header) {
		t.Helper()
		ob := &SuObject{}
		ob.Add(SuStr(table))
		for _, sel := range sels {
			ob.Set(SuStr(sel[0]), SuStr(sel[1]))
		}
		row, hdr, _ := get(th, tran, ob, Only)
		return row, hdr
	}
	val := func(row Row, hdr *Header, col string) string {
		t.Helper()
		if row == nil {
			return "<nil>"
		}
		return AsStr(row.GetVal(hdr, col, nil, nil))
	}

	// tables - real table, virtual schema table, and nonexistent
	row, hdr := test("tables", [2]string{"table", "tmp"})
	assert.T(t).Msg("tables tmp").This(val(row, hdr, "table")).Is("tmp")
	row, hdr = test("tables", [2]string{"table", "tables"})
	assert.T(t).Msg("tables tables").This(val(row, hdr, "table")).Is("tables")
	row, _ = test("tables", [2]string{"table", "nosuch"})
	assert.T(t).Msg("tables nosuch").This(row).Is(nil)

	// columns - regular, derived (uncapitalized), virtual, and nonexistent
	row, hdr = test("columns", [2]string{"table", "tmp"}, [2]string{"column", "u"})
	assert.T(t).Msg("columns tmp u").This(val(row, hdr, "field")).Is("1")
	row, hdr = test("columns", [2]string{"table", "tmp"}, [2]string{"column", "foo"})
	assert.T(t).Msg("columns tmp foo (derived)").
		This(val(row, hdr, "column")).Is("foo")
	assert.T(t).Msg("columns tmp foo field").This(val(row, hdr, "field")).Is("-1")
	row, _ = test("columns", [2]string{"table", "tmp"}, [2]string{"column", "nosuch"})
	assert.T(t).Msg("columns tmp nosuch").This(row).Is(nil)
	row, hdr = test("columns", [2]string{"table", "tables"}, [2]string{"column", "nrows"})
	assert.T(t).Msg("columns tables nrows (virtual)").
		This(val(row, hdr, "field")).Is("1")
	row, _ = test("columns", [2]string{"table", "nosuch"}, [2]string{"column", "x"})
	assert.T(t).Msg("columns nosuch (no panic)").This(row).Is(nil)

	// indexes - key, unique index, and nonexistent table
	row, hdr = test("indexes", [2]string{"table", "tmp"}, [2]string{"columns", "k"})
	assert.T(t).Msg("indexes tmp k").This(val(row, hdr, "key")).Is("true")
	row, hdr = test("indexes", [2]string{"table", "tmp"}, [2]string{"columns", "u"})
	assert.T(t).Msg("indexes tmp u").This(val(row, hdr, "key")).Is("u")
	row, _ = test("indexes", [2]string{"table", "nosuch"}, [2]string{"columns", "x"})
	assert.T(t).Msg("indexes nosuch (no panic)").This(row).Is(nil)

	// views
	row, hdr = test("views", [2]string{"view_name", "myview"})
	assert.T(t).Msg("views myview").This(val(row, hdr, "view_definition")).Is("tmp")
	row, _ = test("views", [2]string{"view_name", "nosuch"})
	assert.T(t).Msg("views nosuch").This(row).Is(nil)
}

// TestGetSchemaNoSels verifies that querying a schema table with no sels
// does not panic (eg. indexing into a nil Indexes) and instead falls back
// to the normal query path.
func TestGetSchemaNoSels(t *testing.T) {
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	qry.DoAdmin(db, "create tmp (k, u, data) key(k) index unique(u)", nil)

	tran := db.NewReadTran()
	defer tran.Complete()
	th := &Thread{}

	test := func(table string, dir Dir, exists bool) {
		t.Helper()
		ob := &SuObject{}
		ob.Add(SuStr(table))
		row, hdr, _ := get(th, tran, ob, dir)
		assert.T(t).Msg(table, " ", dir, " hdr").This(hdr != nil).Is(exists)
		assert.T(t).Msg(table, " ", dir, " row").This(row != nil).Is(exists)
	}
	for _, table := range []string{"tables", "columns", "indexes", "views", "dbstats"} {
		test(table, Any, table != "views" && table != "dbstats")
		test(table, Strat, true)
	}
}
