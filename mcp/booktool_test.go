// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"slices"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBookTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create mybook (name, path, text, order) key(name, path)", nil)

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	tran.Action(th, "insert { name: 'Introduction', path: '', text: 'intro text', order: 1 } into mybook")
	tran.Action(th, "insert { name: 'Reference', path: '', text: 'ref text', order: 2 } into mybook")
	tran.Action(th, "insert { name: 'res', path: '', text: 'res text', order: 3 } into mybook")
	tran.Action(th, "insert { name: 'Date', path: '/Reference', text: 'date text', order: 2 } into mybook")
	tran.Action(th, "insert { name: 'Array', path: '/Reference', text: 'array text', order: 1 } into mybook")
	tran.Action(th, "insert { name: 'FormatEn', path: '/Reference/Date', text: 'format text', order: 1 } into mybook")
	tran.Action(th, "insert { name: 'Images', path: '/res', text: 'image text', order: 1 } into mybook")
	tran.Complete()

	// root children (path not supplied)
	res, err := bookTool("mybook", "")
	assert.That(err == nil)
	assert.This(res.Book).Is("mybook")
	assert.This(res.Path).Is("")
	assert.This(res.Text).Is("")
	children := res.Children
	assert.This(len(children)).Is(2)
	assert.That(!slices.Contains(children, "res"))

	// root children (path supplied)
	res, err = bookTool("mybook", "/")
	assert.That(err == nil)
	assert.This(res.Book).Is("mybook")
	assert.This(res.Path).Is("")
	assert.This(res.Text).Is("")
	children = res.Children
	assert.This(len(children)).Is(2)
	assert.That(!slices.Contains(children, "res"))

	// text and children, sorted by order then name
	res, err = bookTool("mybook", "Reference")
	assert.That(err == nil)
	assert.This(res.Text).Is("ref text")
	children = res.Children
	assert.This(len(children)).Is(2)
	assert.This(children[0]).Is("Array")
	assert.This(children[1]).Is("Date")

	// also works with leading /
	res, err = bookTool("mybook", "/Reference")
	assert.That(err == nil)
	assert.This(res.Text).Is("ref text")
	children = res.Children
	assert.This(len(children)).Is(2)

	// deeper path
	res, err = bookTool("mybook", "Reference/Date")
	assert.That(err == nil)
	assert.This(res.Text).Is("date text")
	children = res.Children
	assert.This(len(children)).Is(1)
	assert.This(children[0]).Is("FormatEn")

	// leaf with no children
	res, err = bookTool("mybook", "Reference/Date/FormatEn")
	assert.That(err == nil)
	assert.This(res.Text).Is("format text")
	assert.This(len(res.Children)).Is(0)

	res, err = bookTool("mybook", "res")
	assert.That(err == nil)
	assert.This(res.Text).Is("")
	assert.This(len(res.Children)).Is(0)

	res, err = bookTool("mybook", "res/Images")
	assert.That(err == nil)
	assert.This(res.Text).Is("")
	assert.This(len(res.Children)).Is(0)
}

func TestSplitPath(t *testing.T) {
	assert := assert.T(t)
	dir, name := splitPath("/a/b/c")
	assert.This(dir).Is("/a/b")
	assert.This(name).Is("c")

	dir, name = splitPath("/abc")
	assert.This(dir).Is("")
	assert.This(name).Is("abc")

	dir, name = splitPath("")
	assert.This(dir).Is("")
	assert.This(name).Is("")
}
