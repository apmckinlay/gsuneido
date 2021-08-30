// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const tmpschema = "(a,b,c,d) key(a) index(b,c)"

func createTestDb() *db19.Database {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	db19.StartConcur(db, 50*time.Millisecond)
	DoAdmin(db, "create tmp "+tmpschema)
	return db
}

func TestAdminCreate(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(db.Schema("tmp")).Is("tmp " + tmpschema)
	assert.T(t).This(func() { DoAdmin(db, "create tables (a) key(a)") }).
		Panics("can't create system table: tables")
	assert.T(t).This(func() { DoAdmin(db, "create tmp (a) key(a)") }).
		Panics("can't create existing table: tmp")
}

func TestAdminEnsure(t *testing.T) {
	db := createTestDb()
	defer db.Close()

	assert.T(t).This(func() { DoAdmin(db, "ensure tables (a) key(a)") }).
		Panics("can't ensure system table: tables")

	// nothing to do
	DoAdmin(db, "ensure tmp "+tmpschema)
	assert.T(t).This(db.Schema("tmp")).Is("tmp " + tmpschema)

	// modify
	DoAdmin(db, "ensure tmp (a, c, e, f) index(b,c) index(e,f)")
	assert.T(t).This(db.Schema("tmp")).
		Is("tmp (a,b,c,d,e,f) key(a) index(b,c) index(e,f)")

	// create
	DoAdmin(db, "ensure tmp2 "+tmpschema)
	assert.T(t).This(db.Schema("tmp2")).Is("tmp2 " + tmpschema)
}

func TestAdminRename(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "rename tmp to indexes") }).
		Panics("can't rename to system table: indexes")
	assert.T(t).This(func() { DoAdmin(db, "rename nonex to foo") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { DoAdmin(db, "rename tmp to tmp") }).
		Panics("existing table: tmp")
	DoAdmin(db, "rename tmp to foo")
	assert.T(t).This(db.Schema("foo")).Is("foo " + tmpschema)
}

func TestAdminAlterCreate(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "alter tables create (x)") }).
		Panics("can't alter system table: tables")
	assert.T(t).This(func() { DoAdmin(db, "alter nonex create (x)") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { DoAdmin(db, "alter tmp create (b)") }).
		Panics("can't create existing column(s): b")
	assert.T(t).This(func() { DoAdmin(db, "alter tmp create index(x)") }).
		Panics("can't create index on nonexistent column(s): x")
	DoAdmin(db, "alter tmp create (x) index(x)")
	assert.T(t).This(db.Schema("tmp")).
		Is("tmp (a,b,c,d,x) key(a) index(b,c) index(x)")
}

func TestAdminAlterRename(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "alter tables rename table to foo") }).
		Panics("can't alter system table: tables")
	assert.T(t).This(func() { DoAdmin(db, "alter nonex rename x to y") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { DoAdmin(db, "alter tmp rename x to y") }).
		Panics("can't rename nonexistent column(s): x")
	assert.T(t).This(func() { DoAdmin(db, "alter nonex rename b to a") }).
		Panics("can't alter nonexistent table: nonex")
	DoAdmin(db, "alter tmp rename b to x")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (a,x,c,d) key(a) index(x,c)")
}

func TestAdminAlterDrop(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "alter tables drop (table)") }).
		Panics("can't alter system table: tables")
	assert.T(t).This(func() { DoAdmin(db, "alter nonex drop (table)") }).
		Panics("nonexistent table: nonex")
	assert.T(t).This(func() { DoAdmin(db, "alter tmp drop (x)") }).
		Panics("can't drop nonexistent column(s): x")
	assert.T(t).This(func() { DoAdmin(db, "alter tmp drop index(x)") }).
		Panics("can't drop nonexistent index: x")
	DoAdmin(db, "alter tmp drop (d)")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (a,b,c,-) key(a) index(b,c)")
	DoAdmin(db, "alter tmp drop (b) index(b,c)")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (a,-,c,-) key(a)")
}

func TestAdminDrop(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "drop columns") }).
		Panics("can't drop system table: columns")
	assert.T(t).This(func() { DoAdmin(db, "drop nonex") }).
		Panics("can't drop nonexistent table: nonex")
	DoAdmin(db, "drop tmp")
	assert.T(t).This(db.Schema("tmp")).Is("")
}

func TestView(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(db.GetView("nonexistent")).Is("")
	assert.T(t).This(func() { DoAdmin(db, "view columns = def") }).
		Panics("can't create view: system table: columns")
	DoAdmin(db, "view foo = bar baz")
	assert.T(t).This(db.GetView("foo")).Is("bar baz")
	assert.T(t).This(func() { DoAdmin(db, "view foo = dup def") }).
		Panics("view: 'foo' already exists")
	DoAdmin(db, "drop foo")
	assert.T(t).This(db.GetView("foo")).Is("")
	DoAdmin(db, "view tmp = over ride")
	assert.T(t).This(db.GetView("tmp")).Is("over ride")
}

func TestFkey(t *testing.T) {
	store := stor.HeapStor(8192)
	db, err := db19.CreateDb(store)
	ck(err)
	DoAdmin(db, "create hdr (a,b) key(a)")
	DoAdmin(db, "create lin (c,d) key(c) index(d) in hdr(a)")
	DoAdmin(db, "create two (e,a) key(e) index(a) in hdr")
	db.Close()
	db, err = db19.OpenDbStor(store, stor.READ, false)
	ck(err)
	db.CheckerSync()
	assert.T(t).This(db.Schema("hdr")).Is("hdr (a,b) key(a) from two(a) from lin(d)")

	DoAdmin(db, "alter two create (f) index(f) in hdr(a)")
	assert.T(t).This(db.Schema("hdr")).
		Is("hdr (a,b) key(a) from two(a) from lin(d) from two(f)")

	DoAdmin(db, "alter two drop index(a)")
	assert.T(t).This(db.Schema("hdr")).
		Is("hdr (a,b) key(a) from lin(d) from two(f)")
}
