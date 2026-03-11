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

	ctx := context.Background()

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

