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

func TestCreateCodeTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create stdlib (name, text, lib_before_text, lib_modified, group, num, parent) key(num) key(name, group)", nil)

	ctx := context.WithValue(context.Background(), approvalFnKey{}, func(before, after string) (bool, error) {
		return true, nil
	})

	// Test invalid library
	_, err := createCodeTool(ctx, "nonexistent", "", "Foo", "function(){}")
	assert.That(err != nil)
	assert.This(err.Error()).Is("library not found: nonexistent")

	// Test invalid name
	_, err = createCodeTool(ctx, "stdlib", "", "lowercase", "function(){}")
	assert.That(err != nil)
	assert.This(err.Error()).Is("invalid name: lowercase")

	// Test invalid code
	_, err = createCodeTool(ctx, "stdlib", "", "Foo", "not valid {{{ code")
	assert.That(err != nil)

	// Test create
	res, err := createCodeTool(ctx, "stdlib", "", "Foo", "function() { return 1 }")
	if err != nil {
		t.Fatal(err)
	}
	assert.This(res.Library).Is("stdlib")
	assert.This(res.Name).Is("Foo")

	// Verify create via direct query: num, parent, text, lib_modified
	th0 := core.NewThread(core.MainThread)
	tran0 := dbmsLocal.Transaction(false)
	q0 := tran0.Query("stdlib where group = -1 and name = 'Foo'", nil)
	hdr0 := q0.Header()
	row0, _ := q0.Get(th0, core.Next)
	assert.That(row0 != nil)
	st0 := core.NewSuTran(tran0, false)
	assert.This(core.ToStr(row0.GetVal(hdr0, "text", th0, st0))).Is("function() { return 1 }")
	n, _ := row0.GetVal(hdr0, "parent", th0, st0).IfInt()
	assert.This(n).Is(0)
	assert.That(row0.GetVal(hdr0, "num", th0, st0) != nil)
	assert.That(row0.GetVal(hdr0, "lib_modified", th0, st0) != nil)
	tran0.Complete()
	th0.Close()

	// Verify create via codeTool
	cr, err := readCodeTool("stdlib", "Foo", 1, codeLineLimit, true)
	if err != nil {
		t.Fatal(err)
	}
	assert.This(cr.Text).Is("function() { return 1 }")
	assert.That(cr.Modified != "")

	// Test create duplicate should fail
	_, err = createCodeTool(ctx, "stdlib", "", "Foo", "function() { return 2 }")
	assert.That(err != nil)
	assert.This(err.Error()).Is("definition already exists: Foo")

	// Test path create creates intermediate folders and sets leaf parent
	res, err = createCodeTool(ctx, "stdlib", "A/B", "Bar", "function() { return 9 }")
	if err != nil {
		t.Fatal(err)
	}

	th3 := core.NewThread(core.MainThread)
	tran3 := dbmsLocal.Transaction(false)
	st3 := core.NewSuTran(tran3, false)

	qf1 := tran3.Query("stdlib where group = 0 and name = 'A'", nil)
	hf1 := qf1.Header()
	rf1, _ := qf1.Get(th3, core.Next)
	assert.That(rf1 != nil)
	f1num, _ := rf1.GetVal(hf1, "num", th3, st3).ToInt()
	assert.This(core.ToStr(rf1.GetVal(hf1, "text", th3, st3))).Is("")
	assert.This(rf1.GetVal(hf1, "group", th3, st3)).Is(core.IntVal(0))

	qf2 := tran3.Query("stdlib where group = "+core.IntVal(f1num).String()+" and name = 'B'", nil)
	hf2 := qf2.Header()
	rf2, _ := qf2.Get(th3, core.Next)
	assert.That(rf2 != nil)
	f2num, _ := rf2.GetVal(hf2, "num", th3, st3).ToInt()
	assert.This(core.ToStr(rf2.GetVal(hf2, "text", th3, st3))).Is("")
	assert.This(rf2.GetVal(hf2, "group", th3, st3)).Is(core.IntVal(f1num))

	qb := tran3.Query("stdlib where group = -1 and name = 'Bar'", nil)
	hb := qb.Header()
	rb, _ := qb.Get(th3, core.Next)
	assert.That(rb != nil)
	bparent, _ := rb.GetVal(hb, "parent", th3, st3).ToInt()
	assert.This(bparent).Is(f2num)

	tran3.Complete()
	th3.Close()
}

func TestCreateCodeTool_RestoreSoftDeleted(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create stdlib (name, text, path, lib_before_text, lib_before_path, lib_modified, lib_committed, group, num, parent) key(num) key(name, group)", nil)

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	n := tran.Action(th, "insert { name: 'Foo', text: 'function(){}', path: 'A/B', lib_before_text: '', lib_before_path: '', lib_modified: #20200101, lib_committed: #20240203, group: -1, num: 42, parent: 7 } into stdlib")
	assert.This(n).Is(1)
	tran.Complete()
	th.Close()

	ctx := context.WithValue(context.Background(), approvalFnKey{}, func(before, after string) (bool, error) {
		return true, nil
	})

	// Soft-delete the record
	_, err := deleteCodeTool(ctx, "stdlib", "Foo")
	if err != nil {
		t.Fatal(err)
	}

	// Verify it's soft-deleted (group = -2)
	th2 := core.NewThread(core.MainThread)
	tran2 := dbmsLocal.Transaction(false)
	q := tran2.Query("stdlib where group = -2 and name = 'Foo'", nil)
	row, _ := q.Get(th2, core.Next)
	assert.That(row != nil)
	tran2.Complete()
	th2.Close()

	// Re-create the definition
	res, err := createCodeTool(ctx, "stdlib", "", "Foo", "function() { return 1 }")
	if err != nil {
		t.Fatal(err)
	}
	assert.This(res.Library).Is("stdlib")
	assert.This(res.Name).Is("Foo")

	// Verify restoration: fields preserved, only text/lib_modified updated
	th3 := core.NewThread(core.MainThread)
	defer th3.Close()
	tran3 := dbmsLocal.Transaction(false)
	q2 := tran3.Query("stdlib where group = -1 and name = 'Foo'", nil)
	hdr := q2.Header()
	row2, _ := q2.Get(th3, core.Next)
	assert.That(row2 != nil)

	st := core.NewSuTran(tran3, false)
	assert.This(core.ToStr(row2.GetVal(hdr, "text", th3, st))).Is("function() { return 1 }")
	assert.This(row2.GetVal(hdr, "group", th3, st)).Is(core.IntVal(-1))
	n, _ = row2.GetVal(hdr, "num", th3, st).IfInt()
	// num preserved
	assert.This(n).Is(42)
	p, _ := row2.GetVal(hdr, "parent", th3, st).IfInt()
	// parent preserved
	assert.This(p).Is(7)
	// modified updated
	assert.That(row2.GetVal(hdr, "lib_modified", th3, st) != nil)
	// committed preserved
	assert.This(row2.GetVal(hdr, "lib_committed", th3, st).String()).Is("#20240203")
	// before_text set from old text
	assert.This(core.ToStr(row2.GetVal(hdr, "lib_before_text", th3, st))).Is("function(){}")
	// path preserved
	assert.This(core.ToStr(row2.GetVal(hdr, "path", th3, st))).Is("A/B")

	// Verify no duplicate (soft-deleted record should be gone)
	q3 := tran3.Query("stdlib where group = -2 and name = 'Foo'", nil)
	row3, _ := q3.Get(th3, core.Next)
	assert.That(row3 == nil)

	tran3.Complete()
}
