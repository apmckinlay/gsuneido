// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// TestColumns covers the edge cases that fuzzColumns' random tables cannot
// produce: a hidden "-" column and a derived column that sorts between stored
// columns (NewFuzzTable's derived columns are named Cr0.. and so always sort
// last).
func TestColumns(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create withhidden (a, -, b) key(a)")
	db.adm("create withderived (a, B, c) key(a)")

	// hidden "-" excluded; b's field number still reflects its physical
	// position after the "-" placeholder
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withhidden'")).
		Is("column=a field=0 table=withhidden | column=b field=2 table=withhidden")
	// derived reported uncapitalized (b) with field -1, sorted between a and c
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withderived'")).
		Is("column=a field=0 table=withderived | column=b field=-1 table=withderived |" +
			" column=c field=1 table=withderived")

	// lookups by (table, column)
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withderived' and column = 'b'")).
		Is("column=b field=-1 table=withderived")
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withhidden' and column = '-'")).
		Is("")
}

// TestTablesLookup verifies that looking up a real table is a keyed GetInfo
// (no ensure), matching Get, while virtual schema tables still resolve
// (which requires ensure to build ts.info).
func TestTablesLookup(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a, b) key(a)")

	ts := NewTable(db.NewReadTran(AllPerms), "tables").(*Tables)
	row := ts.Lookup(nil, Sels{NewSel("table", Pack(SuStr("zz")))})
	assert.T(t).That(!row.Eof())
	assert.T(t).Msg("ensure must not run for a real table").That(ts.info == nil)

	row = ts.Lookup(nil, Sels{NewSel("table", Pack(SuStr("columns")))})
	assert.T(t).Msg("virtual table").That(!row.Eof())
	assert.T(t).That(ts.info != nil)

	row = ts.Lookup(nil, Sels{NewSel("table", Pack(SuStr("nosuch")))})
	assert.T(t).Msg("nonexistent").That(row.Eof())
}

// TestTablesLookupUncommitted verifies that a lookup in tables during an
// update transaction reflects that transaction's inserted and deleted rows
// (via difInfo) rather than the pre-transaction nrows and totalsize.
func TestTablesLookupUncommitted(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a) key(a)")
	db.act("insert {a: 1} into zz")

	hdr := NewTable(db.NewReadTran(AllPerms), "tables").Header()
	lookup := func(tran QueryTran) Row {
		return NewTable(tran, "tables").(*Tables).Lookup(nil, Sels{NewSel("table", Pack(SuStr("zz")))})
	}
	nrows := func(row Row) int { return ToInt(row.GetVal(hdr, "nrows", nil, nil)) }
	size := func(row Row) int { return ToInt(row.GetVal(hdr, "totalsize", nil, nil)) }

	rt := db.NewReadTran(AllPerms)
	assert.T(t).This(nrows(lookup(rt))).Is(1)
	before := size(lookup(rt))

	ut := db.NewUpdateTran(AllPerms)
	defer ut.Abort()
	assert.T(t).This(DoAction(nil, ut, "insert {a: 2} into zz")).Is(1)
	assert.T(t).This(nrows(lookup(ut))).Is(2)
	assert.T(t).Msg("totalsize should include uncommitted insert").
		That(size(lookup(ut)) > before)

	assert.T(t).This(DoAction(nil, ut, "delete zz where a = 1")).Is(1)
	assert.T(t).This(nrows(lookup(ut))).Is(1)
}

// TestTablesMinusTables reproduces a schema table Lookup disagreeing with Get.
// Get scans GetAllInfo while Lookup is keyed via GetInfo, so if testTran
// returns inconsistent info from the two (its GetInfo previously dropped the
// Table name), "tables minus tables" fails to match and returns every table
// instead of none.
func TestTablesMinusTables(t *testing.T) {
	MakeSuTran = func(QueryTran) *SuTran { return nil }
	tran := testTran{}
	q := ParseQuery("tables minus tables", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	rows := getAllRows(q, Next, &Thread{})
	assert.T(t).Msg("tables minus tables").This(len(rows)).Is(0)
}

// TestSchemaJoinNonStringKey verifies that joining a real table to a schema
// table on a non-string key yields an unmatched row rather than panicking in
// the schema table's Lookup (which must not assume string keys).
func TestSchemaJoinNonStringKey(t *testing.T) {
	db := heapDb()
	defer db.Close()

	db.adm("create badtable (table, x) key(table)")
	db.act("insert { table: 123, x: 1 } into badtable")
	db.adm("create badcolumn (table, column, y) key(table, column)")
	db.act("insert { table: 123, column: 456, y: 1 } into badcolumn")
	db.adm("create badview (view_name, z) key(view_name)")
	db.act("insert { view_name: 123, z: 1 } into badview")

	assert.T(t).This(queryAll(db.Database, "badtable leftjoin tables")).
		Is("table=123 x=1")
	assert.T(t).This(queryAll(db.Database, "badcolumn leftjoin columns")).
		Is("column=456 table=123 y=1")
	assert.T(t).This(queryAll(db.Database, "badview leftjoin views")).
		Is("view_name=123 z=1")
}

func TestIndexesPhysicalOrder(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create zz (a,b,c) index(c) key(a) index(b)")
	test := func(dir Dir, expected string) {
		t.Helper()
		q := NewTable(db.NewReadTran(AllPerms), "indexes")
		var cols []string
		for _, row := range getAllRows(q, dir, nil) {
			cols = append(cols, ToStr(row.GetVal(q.Header(), "columns", nil, nil)))
		}
		assert.T(t).This(str.Join("/", cols)).Is(expected)
	}
	test(Next, "c/a/b")
	test(Prev, "b/a/c")
	db.adm("alter zz create index(a,b)")
	test(Next, "c/a/b/a,b")
	test(Prev, "a,b/b/a/c")
	db.adm("alter zz drop index(c)")
	test(Next, "a/b/a,b")
	test(Prev, "a,b/b/a")
}
