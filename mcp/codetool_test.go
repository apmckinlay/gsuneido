package mcp

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCodeTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	// Create stdlib table
	dbmsLocal.Admin("create stdlib (name, text, lib_before_text, group) key(name, group)", nil)

	// Insert a record
	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	n := tran.Action(th, "insert { name: 'Foo', text: 'function(){}', lib_before_text: '', group: -1 } into stdlib")
	assert.This(n).Is(1)
	tran.Complete()

	// Verify insert
	rt := dbmsLocal.Transaction(false)
	q := rt.Query("stdlib", nil)
	row, _ := q.Get(th, core.Next)
	assert.That(row != nil)
	rt.Complete()

	// Test codeTool
	res, err := codeTool("stdlib", "Foo", 1, true)
	if err != nil {
		t.Fatal(err)
	}
	assert.This(res.Library).Is("stdlib")
	assert.This(res.Name).Is("Foo")
	assert.This(res.Text).Is("function(){}")
	assert.This(res.StartLine).Is(1)
	assert.This(res.TotalLines).Is(1)
	assert.That(!res.HasMore)
	assert.This(res.Diff).Is(nil)

	// Test start_line past end
	res, err = codeTool("stdlib", "Foo", 2, true)
	if err != nil {
		t.Fatal(err)
	}
	assert.This(res.Text).Is("")
	assert.This(res.StartLine).Is(2)
	assert.This(res.TotalLines).Is(1)
	assert.That(res.HasMore)
	assert.That(res.Diff == nil)

	// Test invalid library
	_, err = codeTool("nonexistent", "Foo", 1, true)
	assert.That(err != nil)
	assert.This(err.Error()).Is("library not found: nonexistent")

	// Test invalid name
	_, err = codeTool("stdlib", "invalid name", 1, true)
	assert.That(err != nil)
	assert.This(err.Error()).Is("invalid name: invalid name")

	// Test invalid start_line
	_, err = codeTool("stdlib", "Foo", 0, true)
	assert.That(err != nil)
	assert.This(err.Error()).Is("start_line must be >= 1")

	// Test not found
	_, err = codeTool("stdlib", "NonExistent", 1, true)
	assert.That(err != nil)
	assert.This(err.Error()).Is("code not found for: NonExistent in stdlib")
}

func TestIsValidName(t *testing.T) {
	assert := assert.T(t)
	assert.True(isValidName("Object"))
	assert.True(isValidName("MyClass_1"))
	assert.True(isValidName("Is_Empty?"))
	assert.True(isValidName("Do_Something!"))

	assert.False(isValidName("lowerCase"))
	assert.False(isValidName("1Number"))
	assert.False(isValidName("_Underscore"))
	assert.False(isValidName("Space Name"))
	assert.False(isValidName("Hyphen-Name"))
}

func TestSliceCode(t *testing.T) {
	assert := assert.T(t)
	text := "a\nb\nc"
	snippet, total, hasMore := sliceCode(text, 1, 2)
	assert.This(snippet).Is("a\nb")
	assert.This(total).Is(3)
	assert.That(hasMore)

	snippet, total, hasMore = sliceCode(text, 2, 2)
	assert.This(snippet).Is("b\nc")
	assert.This(total).Is(3)
	assert.That(hasMore)

	snippet, total, hasMore = sliceCode(text, 4, 2)
	assert.This(snippet).Is("")
	assert.This(total).Is(3)
	assert.That(hasMore)

	snippet, total, hasMore = sliceCode(text, 1, 5)
	assert.This(snippet).Is("a\nb\nc")
	assert.This(total).Is(3)
	assert.That(!hasMore)
}

func TestAddLineNumbers(t *testing.T) {
	assert := assert.T(t)
	result := addLineNumbers("a\nb\nc", 1)
	assert.This(result).Is("0001: a\n0002: b\n0003: c")

	result = addLineNumbers("a\nb\nc", 10)
	assert.This(result).Is("0010: a\n0011: b\n0012: c")

	result = addLineNumbers("", 1)
	assert.This(result).Is("")
}
