// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDeleteCodeTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create stdlib (name, text, lib_before_text, lib_modified, group, num, parent) key(num) key(name, group)", nil)

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	n := tran.Action(th, "insert { name: 'Foo', text: 'function(){}', lib_before_text: '', lib_modified: #20200101, group: -1, num: 1, parent: 0 } into stdlib")
	assert.This(n).Is(1)
	tran.Complete()
	th.Close()

	ctx := context.WithValue(context.Background(), approvalFnKey{}, func(before, after string) (bool, error) {
		return true, nil
	})

	// invalid library
	_, err := deleteCodeTool(ctx, "nonexistent", "Foo")
	assert.That(err != nil)
	assert.This(err.Error()).Is("library not found: nonexistent")

	// invalid name
	_, err = deleteCodeTool(ctx, "stdlib", "lowercase")
	assert.That(err != nil)
	assert.This(err.Error()).Is("invalid name: lowercase")

	// not found
	_, err = deleteCodeTool(ctx, "stdlib", "Bar")
	assert.That(err != nil)
	assert.This(err.Error()).Is("code not found for: Bar in stdlib")

	// delete success
	res, err := deleteCodeTool(ctx, "stdlib", "Foo")
	if err != nil {
		t.Fatal(err)
	}
	assert.This(res.Library).Is("stdlib")
	assert.This(res.Name).Is("Foo")
	assert.This(res.Action).Is("deleted")

	// verify delete
	th2 := core.NewThread(core.MainThread)
	defer th2.Close()
	tran2 := dbmsLocal.Transaction(false)
	q := tran2.Query("stdlib where group = -1 and name = 'Foo'", nil)
	row, _ := q.Get(th2, core.Next)
	assert.That(row == nil)
	tran2.Complete()
	th2.Close()
}

func TestDeleteCodeTool_SoftDeleteCommitted(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create stdlib (name, text, path, lib_before_text, lib_before_path, lib_modified, lib_committed, group, num, parent) key(num) key(name, group)", nil)

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	n := tran.Action(th, "insert { name: 'Foo', text: 'function(){}', path: 'A/B', lib_before_text: '', lib_before_path: '', lib_modified: #20200101, lib_committed: #20240203, group: -1, num: 1, parent: 0 } into stdlib")
	assert.This(n).Is(1)
	tran.Complete()
	th.Close()

	ctx := context.WithValue(context.Background(), approvalFnKey{}, func(before, after string) (bool, error) {
		return true, nil
	})

	res, err := deleteCodeTool(ctx, "stdlib", "Foo")
	if err != nil {
		t.Fatal(err)
	}
	assert.This(res.Library).Is("stdlib")
	assert.This(res.Name).Is("Foo")
	assert.This(res.Action).Is("soft-deleted")

	th2 := core.NewThread(core.MainThread)
	defer th2.Close()
	tran2 := dbmsLocal.Transaction(false)

	q1 := tran2.Query("stdlib where group = -1 and name = 'Foo'", nil)
	row1, _ := q1.Get(th2, core.Next)
	assert.That(row1 == nil)

	q2 := tran2.Query("stdlib where group = -2 and name = 'Foo'", nil)
	hdr2 := q2.Header()
	row2, _ := q2.Get(th2, core.Next)
	assert.That(row2 != nil)

	st := core.NewSuTran(tran2, false)
	assert.This(core.ToStr(row2.GetVal(hdr2, "text", th2, st))).Is("function(){}")
	assert.This(core.ToStr(row2.GetVal(hdr2, "path", th2, st))).Is("A/B")
	assert.This(core.ToStr(row2.GetVal(hdr2, "lib_before_text", th2, st))).Is("function(){}")
	assert.This(core.ToStr(row2.GetVal(hdr2, "lib_before_path", th2, st))).Is("A/B")
	assert.That(row2.GetVal(hdr2, "lib_modified", th2, st) != nil)

	tran2.Complete()
}
