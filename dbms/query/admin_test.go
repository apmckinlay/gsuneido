// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const tmpschema = "(a,b,c,d) key(a) index(b,c)"

func doAdmin(db *db19.Database, cmd string) {
	DoAdmin(db, cmd, nil)
}

func createTestDb() *db19.Database {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	doAdmin(db, "create tmp "+tmpschema)
	return db
}

func TestAdminCreate(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(db.Schema("tmp")).Is("tmp " + tmpschema)
	xtest := func(cmd, err string) {
		t.Helper()
		assert.T(t).This(func() { doAdmin(db, "create "+cmd) }).Panics(err)
		if !strings.Contains(err, "create") {
			assert.T(t).This(func() { doAdmin(db, "ensure "+cmd) }).Panics(err)
		}
	}

	xtest("tables (a) key(a)",
		"can't modify system table: tables")
	xtest("tmp (a) key(a)",
		"can't create existing table: tmp")
	xtest("xtmp () key(foo)",
		"invalid index column: foo")
	xtest("xtmp (one,two,three) index(one)",
		"key required in xtmp")
	xtest("xtmp (one,two,three) key(bar)",
		"invalid index column: bar")
	xtest("xtmp (one,two,three_lower!) key(one)",
		"_lower! nonexistent column: three")
	db.MustCheck()
}

func TestAdminEnsure(t *testing.T) {
	db := createTestDb()
	defer db.Close()

	// nothing to do
	doAdmin(db, "ensure tmp "+tmpschema)
	assert.T(t).This(db.Schema("tmp")).Is("tmp " + tmpschema)

	// modify
	doAdmin(db, "ensure tmp (a, c, e, f, G) index(b,c) index(e,f)")
	assert.T(t).This(db.Schema("tmp")).
		Is("tmp (a,b,c,d,e,f,G) key(a) index(b,c) index(e,f)")

	// create
	doAdmin(db, "ensure tmp2 "+tmpschema)
	assert.T(t).This(db.Schema("tmp2")).Is("tmp2 " + tmpschema)

	doAdmin(db, "ensure tmp3 (a) key(a) index(a_lower!)")

	assert.T(t).This(func() { doAdmin(db, "ensure tmp (z, z)") }).
		Panics("duplicate column")

	assert.T(t).This(func() { doAdmin(db, "ensure tmp key(x) key(x)") }).
		Panics("duplicate index")

	assert.T(t).This(func() { doAdmin(db, "ensure tmp key(x) index(x)") }).
		Panics("duplicate index")

	// existing index but different
	assert.T(t).This(func() { doAdmin(db, "ensure tmp index unique(b,c)") }).
		Panics(("ensure: index exists but is different"))

	doAdmin(db, "ensure tmp key(d_lower!)")
	assert.T(t).This(db.Schema("tmp")).
		Is("tmp (a,b,c,d,e,f,G) key(a) index(b,c) index(e,f) key(d_lower!)")

	doAdmin(db, "create tmp4 "+tmpschema)
	doAdmin(db, "ensure tmp4 (X)")
	assert.T(t).This(db.Schema("tmp4")).Is("tmp4 (a,b,c,d,X) key(a) index(b,c)")
	db.MustCheck()
}

func TestAdminEnsure2(*testing.T) {
	db := createTestDb()
	defer db.Close()
	act(db, "insert { a: 1 } into tmp")
	doAdmin(db, "ensure tmp (x, y) index unique(x)")
	db.MustCheck()
}

func TestAdminRename(t *testing.T) {
	db := createTestDb()
	defer db.Close()

	assert.T(t).This(func() { doAdmin(db, "rename tmp to indexes") }).
		Panics("can't modify system table: indexes")
	assert.T(t).This(func() { doAdmin(db, "rename nonex to foo") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { doAdmin(db, "rename tmp to tmp") }).
		Panics("existing table: tmp")
	doAdmin(db, "rename tmp to foo")
	assert.T(t).This(db.Schema("foo")).Is("foo " + tmpschema)
	db.MustCheck()
}

func TestAdminAlterCreate(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	act(db, "insert { a: 1, b: 2, c: 3, d: 4 } into tmp")

	assert.T(t).This(func() { doAdmin(db, "alter tables create (x)") }).
		Panics("can't modify system table: tables")
	assert.T(t).This(func() { doAdmin(db, "alter nonex create (x)") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { doAdmin(db, "alter tmp create (b)") }).
		Panics("can't create existing column(s): b")
	assert.T(t).This(func() { doAdmin(db, "alter tmp create index(x)") }).
		Panics("invalid index column: x in tmp")

	doAdmin(db, "alter tmp create (x,Y) index(x)")
	assert.T(t).This(db.Schema("tmp")).
		Is("tmp (a,b,c,d,x,Y) key(a) index(b,c) index(x)")

	assert.T(t).This(func() { doAdmin(db, "alter tmp create (z, z)") }).
		Panics("duplicate column")
	assert.T(t).This(func() { doAdmin(db, "alter tmp create index(c) index(c)") }).
		Panics("duplicate index")
	assert.T(t).This(func() { doAdmin(db, "alter tmp create index(x)") }).
		Panics("duplicate index")
	assert.T(t).This(func() { doAdmin(db, "alter tmp create key(x)") }).
		Panics("duplicate index")
	db.MustCheck()
}

func TestAdminAlterRename(t *testing.T) {
	db := createTestDb()
	defer db.Close()

	assert.T(t).This(func() { doAdmin(db, "alter tables rename table to foo") }).
		Panics("can't modify system table: tables")
	assert.T(t).This(func() { doAdmin(db, "alter nonex rename x to y") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { doAdmin(db, "alter tmp rename x to y") }).
		Panics("can't rename nonexistent column: x")
	assert.T(t).This(func() { doAdmin(db, "alter nonex rename b to a") }).
		Panics("can't alter nonexistent table: nonex")
	doAdmin(db, "alter tmp rename b to x")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (a,x,c,d) key(a) index(x,c)")
	assert.T(t).This(func() { doAdmin(db, "alter tmp rename c to d") }).
		Panics("existing column")
	assert.T(t).This(func() { doAdmin(db, "alter tmp rename c to z, d to z") }).
		Panics("can't rename to existing column: z")
	doAdmin(db, "alter tmp rename a to b, b to z, x to b")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (z,b,c,d) key(z) index(b,c)")
	db.MustCheck()
}

func TestAdminAlterDrop(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { doAdmin(db, "alter tables drop (table)") }).
		Panics("can't modify system table: tables")
	assert.T(t).This(func() { doAdmin(db, "alter nonex drop (table)") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { doAdmin(db, "alter tmp drop (x)") }).
		Panics("can't drop nonexistent column: x")
	assert.T(t).This(func() { doAdmin(db, "alter tmp drop index(x)") }).
		Panics("can't drop nonexistent index: tmp (x)")
	doAdmin(db, "alter tmp drop (d)")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (a,b,c,-) key(a) index(b,c)")
	doAdmin(db, "alter tmp drop (b) index(b,c)")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (a,-,c,-) key(a)")

	doAdmin(db, "create tmp2 (a,b,C,D,a_lower!) key(a)")
	doAdmin(db, "alter tmp2 drop (C,d,a_lower!)")
	assert.T(t).This(db.Schema("tmp2")).Is("tmp2 (a,b) key(a)")

	doAdmin(db, "create tmp3 (a,b,c) key(a)")
	assert.T(t).This(func() { doAdmin(db, "alter tmp3 drop key(a)") }).
		Panics("can't drop all keys: tmp")

	doAdmin(db, "create tmp4 (a,b,c) key(a) key(b,c) index(c)")
	assert.T(t).This(func() { doAdmin(db, "alter tmp4 drop key(a)") }).
		Panics("can't drop key used to make index unique: tmp4 (a)")

	doAdmin(db, "create tmp5 (a,b) key(a) key(a,b)") // key(a) will be primary
	doAdmin(db, "alter tmp5 drop key(a)")            // key(a,b) now primary
	act(db, "insert { a: 1 } into tmp5")
	assert.T(t).This(func() { act(db, "insert { a: 1 } into tmp5") }).
		Panics("duplicate")
	db.MustCheck()
}

func TestAdminDrop(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { doAdmin(db, "drop columns") }).
		Panics("can't modify system table: columns")
	assert.T(t).This(func() { doAdmin(db, "drop nonex") }).
		Panics("can't drop nonexistent table: nonex")
	doAdmin(db, "drop tmp")
	assert.T(t).This(db.Schema("tmp")).Is("")
	db.MustCheck()
}

func TestView(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(db.GetView("nonexistent")).Is("")
	assert.T(t).This(func() { doAdmin(db, "view columns = def") }).
		Panics("can't modify system table: columns")
	doAdmin(db, "view foo = bar baz")
	assert.T(t).This(db.GetView("foo")).Is("bar baz")
	assert.T(t).This(func() { doAdmin(db, "view foo = dup def") }).
		Panics("view: 'foo' already exists")
	doAdmin(db, "drop foo")
	assert.T(t).This(db.GetView("foo")).Is("")
	doAdmin(db, "view tmp = over ride")
	assert.T(t).This(db.GetView("tmp")).Is("over ride")
	db.MustCheck()
}

func TestFkey(t *testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db.CheckerSync()

	schemas := map[string]string{}
	check := func() {
		t.Helper()
		rt := db.NewReadTran()
		for table, schema := range schemas {
			assert.T(t).This(db.Schema(table)).Is(schema)
			if schema == "" {
				continue
			}
			sch := rt.GetSchema(table)
			for _, ix := range sch.Indexes {
				if ix.Fk.Table != "" && rt.GetInfo(ix.Fk.Table) != nil {
					sch2 := rt.GetSchema(ix.Fk.Table)
					assert.T(t).
						Msg(table, ix.Columns, "Fk", ix.Fk).
						This(sch2.Indexes[ix.Fk.IIndex].Columns).Is(ix.Fk.Columns)
				}
				for _, fk := range ix.FkToHere {
					sch2 := rt.GetSchema(fk.Table)
					assert.T(t).Msg(table, ix.Columns, "FkToHere", fk).
						This(sch2.Indexes[fk.IIndex].Columns).Is(fk.Columns)
				}
			}
		}
	}

	doAdmin(db, "create hdr (a,b) key(a)")
	schemas["hdr"] = "hdr (a,b) key(a)"
	check()

	doAdmin(db, "create lin (c,d) key(c) index(d) in hdr(a)")
	schemas["lin"] = "lin (c,d) key(c) index(d) in hdr(a)"
	schemas["hdr"] = "hdr (a,b) key(a) from lin(d)"
	check()

	doAdmin(db, "create two (e,a) key(e) index(a) in hdr")
	schemas["two"] = "two (e,a) key(e) index(a) in hdr"
	schemas["hdr"] = "hdr (a,b) key(a) from lin(d) from two(a)"
	check()

	doAdmin(db, "alter two create (f) index(f) in hdr(a)")
	schemas["two"] = "two (e,a,f) key(e) index(a) in hdr index(f) in hdr(a)"
	schemas["hdr"] = "hdr (a,b) key(a) from lin(d) from two(a) from two(f)"
	check()

	doAdmin(db, "alter two drop index(a)")
	schemas["two"] = "two (e,a,f) key(e) index(f) in hdr(a)"
	schemas["hdr"] = "hdr (a,b) key(a) from lin(d) from two(f)"
	check()

	assert.T(t).This(func() { doAdmin(db, "alter hdr drop key(a)") }).
		Panics("can't drop index used by foreign keys")

	doAdmin(db, "create three (f,e) key(f) index(e) in two")
	schemas["three"] = "three (f,e) key(f) index(e) in two"
	schemas["two"] = "two (e,a,f) key(e) from three(e) index(f) in hdr(a)"
	check()

	doAdmin(db, "create four (g) key(g)")
	doAdmin(db, "ensure four (h) index(h) in lin(c)")
	schemas["four"] = "four (g,h) key(g) index(h) in lin(c)"
	schemas["lin"] = "lin (c,d) key(c) from four(h) index(d) in hdr(a)"
	check()

	doAdmin(db, "rename four to newfour")
	schemas["newfour"] = "newfour (g,h) key(g) index(h) in lin(c)"
	schemas["four"] = ""
	schemas["lin"] = "lin (c,d) key(c) from newfour(h) index(d) in hdr(a)"
	check()

	doAdmin(db, "alter newfour rename h to hh")
	schemas["newfour"] = "newfour (g,hh) key(g) index(hh) in lin(c)"
	schemas["lin"] = "lin (c,d) key(c) from newfour(hh) index(d) in hdr(a)"
	check()

	assert.T(t).This(func() { doAdmin(db, "drop hdr") }).
		Panics("can't drop table used by foreign keys")

	// recursive foreign key
	doAdmin(db, "create recur (a,b) key(a) index(b) in recur(a)")
	schemas["recur"] = "recur (a,b) key(a) from recur(b) index(b) in recur(a)"
	check()
	doAdmin(db, "drop recur") // has FkToHere, but only to itself
	delete(schemas, "recur")
	check()

	doAdmin(db, "create head (a,b) key(a) key(b)")
	schemas["head"] = "head (a,b) key(a) key(b)"
	check()
	doAdmin(db, "create line (c,d) key(c)")
	doAdmin(db, "alter line create index(d) in head(b)")
	schemas["line"] = "line (c,d) key(c) index(d) in head(b)"
	schemas["head"] = "head (a,b) key(a) key(b) from line(d)"
	check()

	db.MustCheck()
	db.Close()
	db, err = db19.OpenDbStor(store, stor.Read, false)
	ck(err)
	check()
}

func TestCreateIndexOnExistingTable(*testing.T) {
	db := createTestDb()
	act(db, "insert { a: 1, b: 2, c: 3, d: 4 } into tmp")
	act(db, "insert { a: 3, b: 4 } into tmp")
	time.Sleep(100 * time.Millisecond) // ensure persisted
	db.MustCheck()
	doAdmin(db, "ensure tmp index(d)")
	db.MustCheck()
	doAdmin(db, "ensure tmp key(c,d)")
	db.MustCheck()
}

func TestNoColumns(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	doAdmin(db, "create nocols () key()")
	db.MustCheck()
	db.Close()
	db, err = db19.OpenDbStor(store, stor.Read, false)
	ck(err)
	db.MustCheck()
}

func act(db *db19.Database, act string) {
	ut := db.NewUpdateTran()
	defer ut.Commit()
	n := DoAction(nil, ut, act)
	assert.This(n).Is(1)
}
