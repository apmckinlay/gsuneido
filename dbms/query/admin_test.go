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
}

func TestAdminEnsure(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "ensure tables (a) key(a)") }).
		Panics("can't ensure system table: tables")
	DoAdmin(db, "ensure tmp " + tmpschema)
	assert.T(t).This(db.Schema("tmp")).Is("tmp " + tmpschema)
	DoAdmin(db, "ensure tmp (a, c, e, f) index(b,c) index(e,f)")
	assert.T(t).This(db.Schema("tmp")).
		Is("tmp (a,b,c,d,e,f) key(a) index(b,c) index(e,f)")
}

func TestAdminRename(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "rename tmp to indexes") }).
		Panics("can't rename to system table: indexes")
	DoAdmin(db, "rename tmp to foo")
	assert.T(t).This(db.Schema("foo")).Is("foo " + tmpschema)
}

func TestAdminAlterCreate(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "alter tables create (x)") }).
		Panics("can't alter system table: tables")
	DoAdmin(db, "alter tmp create (x) index(x)")
	assert.T(t).This(db.Schema("tmp")).
		Is("tmp (a,b,c,d,x) key(a) index(b,c) index(x)")
}

func TestAdminAlterRename(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "alter tables rename table to foo") }).
		Panics("can't alter system table: tables")
	DoAdmin(db, "alter tmp rename b to x")
	assert.T(t).This(db.Schema("tmp")).Is("tmp (a,x,c,d) key(a) index(x,c)")
}

func TestAdminAlterDrop(t *testing.T) {
	db := createTestDb()
	defer db.Close()
	assert.T(t).This(func() { DoAdmin(db, "alter tables drop (table)") }).
		Panics("can't alter system table: tables")
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
	DoAdmin(db, "drop tmp")
	assert.T(t).This(db.Schema("tmp")).Is("")
}
