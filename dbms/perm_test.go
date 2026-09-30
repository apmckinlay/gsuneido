// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	qry "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// TestPermCursor checks that a cursor enforces the permissions of the
// transaction it is used with, not the perms-less transaction it is created with.
func TestPermCursor(t *testing.T) {
	a := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	d := NewDbmsLocal(db)
	d.AdminTest("create tmp (k, data) key(k)")
	ut := db.NewUpdateTran(nil)
	qry.DoAction(&Thread{}, ut, "insert { k: 1, data: 'x' } into tmp")
	ut.Commit()

	cur := d.Cursor("tmp", nil, nil)

	// the tran used to fetch has no read permission
	deny := db.NewReadTran(&Perms{})
	a.This(func() {
		cur.Get(&Thread{}, &ReadTranLocal{ReadTran: deny}, Next)
	}).Panics("not authorized")

	// with read permission the fetch succeeds
	perms := &Perms{}
	perms.AddTable("tmp", "read")
	allow := db.NewReadTran(perms)
	row, _ := cur.Get(&Thread{}, &ReadTranLocal{ReadTran: allow}, Next)
	a.This(row == nil).Is(false)
}

// TestPermCursorBuild checks that the transaction a cursor is built with
// enforces permissions, so a cursor can't escape them by using its build
// transaction instead of the transaction it is fetched with.
func TestPermCursorBuild(t *testing.T) {
	a := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	defer db.Close()
	d := NewDbmsLocal(db)
	d.AdminTest("create tmp (k, data) key(k)")

	// the build transaction carries the connection's perms
	a.This(func() {
		d.Cursor("tmp", nil, &Perms{})
	}).Panics("not authorized: tmp")

	perms := &Perms{}
	perms.AddTable("tmp", "read")
	d.Cursor("tmp", nil, perms) // does not panic
}
