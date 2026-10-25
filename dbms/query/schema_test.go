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
