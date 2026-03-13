// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestApplyLineEdit(t *testing.T) {
	assert := assert.T(t)

	oldText := strings.Join([]string{
		"function()",
		"\t{",
		"\treturn 1",
		"\t}",
	}, "\n")

	newText, err := applyLineEdit(oldText, "replace_lines", 3, 1, "\treturn 2")
	assert.That(err == nil)
	assert.This(newText).Is(strings.Join([]string{
		"function()",
		"\t{",
		"\treturn 2\r",
		"\t}",
	}, "\n"))

	newText, err = applyLineEdit(oldText, "insert_before", 3, 0, "\t// inserted\n")
	assert.That(err == nil)
	assert.This(newText).Is(strings.Join([]string{
		"function()",
		"\t{",
		"\t// inserted\r",
		"\treturn 1",
		"\t}",
	}, "\n"))

	newText, err = applyLineEdit(oldText, "replace_lines", 3, 1, "")
	assert.That(err == nil)
	assert.This(newText).Is(strings.Join([]string{
		"function()",
		"\t{",
		"\t}",
	}, "\n"))

	_, err = applyLineEdit(oldText, "replace_lines", 0, 1, "x")
	assert.That(err != nil)
	assert.This(err.Error()).Is("line must be >= 1")

	newText2, err := applyLineEdit(oldText, "insert_after", 3, 0, "\t// after\n")
	assert.That(err == nil)
	assert.This(newText2).Is(strings.Join([]string{
		"function()",
		"\t{",
		"\treturn 1",
		"\t// after\r",
		"\t}",
	}, "\n"))

	_, err = applyLineEdit(oldText, "replace_lines", 1, 6, "x")
	assert.That(err != nil)
	assert.This(err.Error()).Is("line 1 out of bounds for 4 lines")
}

func TestPatchCodeTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create stdlib (name, text, lib_before_text, lib_modified, group, num, parent) key(num) key(name, group)", nil)

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	n := tran.Action(th, "insert { name: 'Foo', text: 'function()\n\t{\n\treturn 1\n\t}', lib_before_text: '', lib_modified: #20200101, group: -1, num: 1, parent: 0 } into stdlib")
	assert.This(n).Is(1)
	tran.Complete()
	th.Close()

	ctx := context.WithValue(context.Background(), approvalFnKey{}, func() (bool, error) {
		return true, nil
	})

	_, err := patchCodeTool(ctx, "nonexistent", "Foo", "replace_lines", 1, 1, "x")
	assert.That(err != nil)
	assert.This(err.Error()).Is("library not found: nonexistent")

	_, err = patchCodeTool(ctx, "stdlib", "lowercase", "replace_lines", 1, 1, "x")
	assert.That(err != nil)
	assert.This(err.Error()).Is("invalid name: lowercase")

	_, err = patchCodeTool(ctx, "stdlib", "Bar", "replace_lines", 1, 1, "x")
	assert.That(err != nil)
	assert.This(err.Error()).Is("code not found for: Bar in stdlib")

	oldText := "function()\n\t{\n\treturn 1\n\t}"
	newText := "function()\n\t{\n\treturn 2\r\n\t}"
	res, err := patchCodeTool(ctx, "stdlib", "Foo", "replace_lines", 3, 1, "\treturn 2")
	if err != nil {
		t.Fatal(err)
	}
	assert.This(res.Library).Is("stdlib")
	assert.This(res.Name).Is("Foo")

	th2 := core.NewThread(core.MainThread)
	tran2 := dbmsLocal.Transaction(false)
	q2 := tran2.Query("stdlib where group = -1 and name = 'Foo'", nil)
	hdr2 := q2.Header()
	row2, _ := q2.Get(th2, core.Next)
	assert.That(row2 != nil)
	st2 := core.NewSuTran(tran2, false)
	assert.This(core.ToStr(row2.GetVal(hdr2, "text", th2, st2))).Is(newText)
	assert.This(core.ToStr(row2.GetVal(hdr2, "lib_before_text", th2, st2))).Is(oldText)
	assert.That(row2.GetVal(hdr2, "lib_modified", th2, st2) != nil)
	tran2.Complete()
	th2.Close()

	_, err = patchCodeTool(ctx, "stdlib", "Foo", "replace_lines", 3, 1, "\tnot valid {{{ code")
	assert.That(err != nil)

	_, err = patchCodeTool(ctx, "stdlib", "Foo", "replace_lines", 1, 10, "x")
	assert.That(err != nil)
	assert.This(err.Error()).Is("line 1 out of bounds for 4 lines")

	nextText := "function()\n\t{\n\treturn 3\r\n\t}"
	_, err = patchCodeTool(ctx, "stdlib", "Foo", "replace_lines", 3, 1, "\treturn 3")
	if err != nil {
		t.Fatal(err)
	}

	th3 := core.NewThread(core.MainThread)
	tran3 := dbmsLocal.Transaction(false)
	q3 := tran3.Query("stdlib where group = -1 and name = 'Foo'", nil)
	hdr3 := q3.Header()
	row3, _ := q3.Get(th3, core.Next)
	assert.That(row3 != nil)
	st3 := core.NewSuTran(tran3, false)
	assert.This(core.ToStr(row3.GetVal(hdr3, "text", th3, st3))).Is(nextText)
	assert.This(core.ToStr(row3.GetVal(hdr3, "lib_before_text", th3, st3))).Is(oldText)
	tran3.Complete()
	th3.Close()
}

func TestValidatePatchModeArgs(t *testing.T) {
	assert := assert.T(t)

	err := validatePatchModeArgs("replace_lines", 1)
	assert.That(err == nil)

	err = validatePatchModeArgs("insert_before", 0)
	assert.That(err == nil)

	err = validatePatchModeArgs("insert_after", 0)
	assert.That(err == nil)

	err = validatePatchModeArgs("invalid_mode", 0)
	assert.That(err != nil)
	assert.This(err.Error()).Is("invalid mode: invalid_mode (must be 'insert_before', 'insert_after', or 'replace_lines')")

	err = validatePatchModeArgs("replace_lines", 0)
	assert.That(err != nil)
	assert.This(err.Error()).Is("count must be >= 1 for replace_lines")

	err = validatePatchModeArgs("insert_before", 1)
	assert.That(err != nil)
	assert.This(err.Error()).Is("count is only valid for replace_lines")
}
