// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"slices"
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

func TestColumnsSchemaRefresh(t *testing.T) {
	test := func(cols, position, change string, dir Dir, expected string,
		moreChanges ...string) {
		t.Helper()
		db := heapDb()
		defer db.Close()
		db.adm("create zz " + cols + " key(a)")
		q := NewTable(db.NewReadTran(AllPerms), "columns")
		hdr := q.Header()
		for {
			q.SetTran(db.NewReadTran(AllPerms))
			row := q.Get(nil, Next)
			if row.Eof() {
				t.Fatal("cursor position not found:", position)
			}
			if ToStr(row.GetVal(hdr, "table", nil, nil)) == "zz" &&
				ToStr(row.GetVal(hdr, "column", nil, nil)) == position {
				break
			}
		}
		db.adm(change)
		q.SetTran(db.NewReadTran(AllPerms))
		for _, change := range moreChanges {
			db.adm(change)
			q.SetTran(db.NewReadTran(AllPerms))
		}
		// A cursor may change transactions more than once before fetching.
		q.SetTran(db.NewReadTran(AllPerms))
		assert.T(t).This(row2str(hdr, q.Get(nil, dir))).Is(expected)
	}

	test("(a, B, c)", "a", "alter zz drop (B)", Next,
		"column=c field=1 table=zz")
	test("(a, b, c)", "a", "alter zz drop (b)", Next,
		"column=c field=2 table=zz")
	test("(a, B, c, d)", "c", "alter zz drop (B)", Next,
		"column=d field=2 table=zz")
	test("(a, B, c)", "b", "alter zz drop (B)", Next,
		"column=c field=1 table=zz")
	test("(a, B, c)", "c", "alter zz drop (B)", Prev,
		"column=a field=0 table=zz")
	test("(a, b, c)", "c", "alter zz drop (b)", Prev,
		"column=a field=0 table=zz")
	test("(a, B, c)", "b", "alter zz drop (B)", Prev,
		"column=a field=0 table=zz")
	test("(a, C)", "a", "alter zz create (b)", Next,
		"column=b field=1 table=zz")
	test("(a, C)", "c", "alter zz create (b)", Prev,
		"column=b field=1 table=zz")
	test("(a, b)", "a", "create aa (x) key(x)", Next,
		"column=b field=1 table=zz")
	test("(a, b)", "a", "drop zz", Next, "nil")
	test("(a, b)", "b", "drop zz", Prev,
		"column=view_name field=0 table=views")
	test("(a, B, c)", "b", "alter zz drop (B)", Next,
		"column=c field=1 table=zz", "alter zz create (B)")
	test("(a, B, c)", "b", "alter zz drop (B)", Prev,
		"column=a field=0 table=zz", "alter zz create (B)")
	test("(a, b)", "a", "drop zz", Next,
		"column=b field=1 table=zz", "create zz (a, b) key(a)")
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

func TestSchemaTablesRefresh(t *testing.T) {
	test := func(table, position string, dir Dir, expected string, changes ...string) {
		t.Helper()
		db := heapDb()
		defer db.Close()
		for _, name := range []string{"aa", "bb", "cc"} {
			db.adm("create " + name + " (a,b,c) key(a) index(c) index(b)")
			db.adm("view v" + name[:1] + " = " + name)
		}
		q := NewTable(db.NewReadTran(AllPerms), table)
		hdr := q.Header()
		cols := []string{"table"}
		if table == "indexes" {
			cols = append(cols, "columns")
		} else if table == "views" {
			cols = []string{"view_name"}
		}
		key := func(row Row) string {
			if row.Eof() {
				return ""
			}
			parts := make([]string, len(cols))
			for i, col := range cols {
				parts[i] = ToStr(row.GetVal(hdr, col, nil, nil))
			}
			return str.Join("/", parts)
		}
		for {
			q.SetTran(db.NewReadTran(AllPerms))
			row := q.Get(nil, Next)
			if row.Eof() {
				t.Fatal(table, "cursor position not found:", position)
			}
			if key(row) == position {
				break
			}
		}
		for _, change := range changes {
			db.adm(change)
			q.SetTran(db.NewReadTran(AllPerms))
		}
		q.SetTran(db.NewReadTran(AllPerms))
		assert.T(t).This(key(q.Get(nil, dir))).Is(expected)

		fresh := NewTable(db.NewReadTran(AllPerms), table)
		rows := getAllRows(fresh, Next, nil)
		if expected != "" {
			i := slices.IndexFunc(rows, func(row Row) bool { return key(row) == expected })
			assert.T(t).That(i >= 0)
			step := 1
			if dir == Prev {
				step = -1
			}
			for i += step; i >= 0 && i < len(rows); i += step {
				q.SetTran(db.NewReadTran(AllPerms))
				assert.T(t).This(row2str(hdr, q.Get(nil, dir))).Is(row2str(hdr, rows[i]))
			}
		}
		assert.T(t).That(q.Get(nil, dir).Eof())
		q.SetTran(db.NewReadTran(AllPerms))
		assert.T(t).That(q.Get(nil, Next).Eof())
		assert.T(t).That(q.Get(nil, Prev).Eof())
		q.Rewind()
		for _, row := range rows {
			assert.T(t).This(row2str(hdr, q.Get(nil, Next))).Is(row2str(hdr, row))
		}
		assert.T(t).That(q.Get(nil, Next).Eof())
		q.Rewind()
		for i := len(rows) - 1; i >= 0; i-- {
			q.SetTran(db.NewReadTran(AllPerms))
			assert.T(t).This(row2str(hdr, q.Get(nil, Prev))).Is(row2str(hdr, rows[i]))
		}
		assert.T(t).That(q.Get(nil, Prev).Eof())
	}

	test("tables", "bb", Next, "cc", "create ab (a) key(a)")
	test("tables", "bb", Prev, "ab", "create ab (a) key(a)")
	test("tables", "bb", Next, "cc", "drop aa")
	test("tables", "bb", Prev, "", "drop aa")
	test("tables", "bb", Next, "cc", "drop bb")
	test("tables", "bb", Prev, "aa", "drop bb")
	test("tables", "bb", Next, "cc", "drop bb", "create bb (a) key(a)")
	test("tables", "bb", Prev, "ab", "drop bb", "create ab (a) key(a)")

	test("indexes", "bb/a", Next, "bb/c", "create ab (a) key(a)")
	test("indexes", "bb/a", Prev, "ab/a", "create ab (a) key(a)")
	test("indexes", "bb/b", Next, "cc/a", "drop aa")
	test("indexes", "bb/b", Prev, "bb/c", "drop aa")
	test("indexes", "bb/b", Next, "cc/a", "alter bb drop index(b)")
	test("indexes", "bb/b", Prev, "bb/c", "alter bb drop index(b)")
	test("indexes", "bb/c", Next, "cc/a", "alter bb drop index(b)")
	test("indexes", "bb/c", Prev, "bb/a", "alter bb drop index(b)")
	test("indexes", "bb/b", Next, "bb/a,b", "alter bb create index(a,b)")
	test("indexes", "bb/b", Prev, "bb/c", "alter bb create index(a,b)")
	test("indexes", "bb/a", Next, "bb/c", "alter bb create index(a,b)")
	test("indexes", "bb/a", Prev, "aa/b", "alter bb create key()")
	test("indexes", "bb/b", Next, "bb/", "alter bb create key()")
	test("indexes", "bb/c", Next, "cc/a", "drop bb")
	test("indexes", "bb/a", Prev, "aa/b", "drop bb")
	test("indexes", "cc/c", Next, "", "drop cc")
	test("indexes", "cc/c", Prev, "bb/b", "drop cc")
	test("indexes", "bb/b", Next, "cc/a", "alter bb drop index(b)", "alter bb drop index(c)")
	test("indexes", "bb/b", Next, "cc/a", "alter bb drop index(b)", "alter bb create index(b)")
	test("indexes", "bb/c", Next, "bb/b", "alter bb drop index(c)")
	test("indexes", "bb/c", Prev, "bb/a", "alter bb drop index(c)")
	test("indexes", "bb/c", Next, "cc/a", "alter bb drop index(c)", "alter bb drop index(b)")
	test("indexes", "bb/b", Next, "bb/a,b", "alter bb drop index(b)", "alter bb create index(a,b)")
	test("indexes", "bb/b", Next, "", "drop aa", "drop bb", "drop cc")
	test("indexes", "bb/b", Prev, "", "drop aa", "drop bb", "drop cc")

	test("views", "vb", Next, "vc", "view vab = aa")
	test("views", "vb", Prev, "vab", "view vab = aa")
	test("views", "vb", Next, "vc", "drop va")
	test("views", "vb", Prev, "", "drop va")
	test("views", "vb", Next, "vc", "drop vb")
	test("views", "vb", Prev, "va", "drop vb")
	test("views", "vc", Next, "", "drop vc")
	test("views", "vc", Prev, "vb", "drop vc")
	test("views", "vb", Next, "vc", "drop vb", "view vb = cc")
	test("views", "vb", Next, "", "drop va", "drop vb", "drop vc")
	test("views", "vb", Prev, "", "drop va", "drop vb", "drop vc")
}
