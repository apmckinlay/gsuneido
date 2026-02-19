// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

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

func TestSearchBookTool(t *testing.T) {
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
	tran.Action(th, "insert { name: 'Date', path: '/Reference', text: 'date functions', order: 2 } into mybook")
	tran.Action(th, "insert { name: 'Array', path: '/Reference', text: 'array functions', order: 1 } into mybook")
	tran.Action(th, "insert { name: 'FormatEn', path: '/Reference/Date', text: 'format text', order: 1 } into mybook")
	tran.Action(th, "insert { name: 'MultiMatch', path: '', text: 'line one has text\nline two has text\nline three has other', order: 3 } into mybook")
	tran.Complete()

	// search by text - matches "text" in intro text, ref text, format text, and MultiMatch
	res, err := searchBook("mybook", "", "text", false)
	assert.That(err == nil)
	assert.This(len(res.Matches)).Is(4)
	assert.This(res.Matches[0].Path).Is("/Introduction")
	assert.This(res.Matches[0].Lines).Is([]string{"0001: intro text"})
	assert.This(res.Matches[1].Path).Is("/MultiMatch")
	assert.This(res.Matches[1].Lines).Is([]string{"0001: line one has text", "0002: line two has text"})
	assert.This(res.Matches[2].Path).Is("/Reference")
	assert.This(res.Matches[2].Lines).Is([]string{"0001: ref text"})
	assert.This(res.Matches[3].Path).Is("/Reference/Date/FormatEn")
	assert.This(res.Matches[3].Lines).Is([]string{"0001: format text"})

	// search by path - sorted by path, name
	res, err = searchBook("mybook", "Reference", "", false)
	assert.That(err == nil)
	paths := make([]string, len(res.Matches))
	for i, m := range res.Matches {
		paths[i] = m.Path
	}
	assert.This(paths).Is([]string{"/Reference", "/Reference/Array", "/Reference/Date", "/Reference/Date/FormatEn"})

	// search by both path and text - sorted by path, name
	res, err = searchBook("mybook", "Reference", "functions", false)
	assert.That(err == nil)
	paths = make([]string, len(res.Matches))
	for i, m := range res.Matches {
		paths[i] = m.Path
	}
	assert.This(paths).Is([]string{"/Reference/Array", "/Reference/Date"})

	// case insensitive (default)
	res, err = searchBook("mybook", "", "TEXT", false)
	assert.That(err == nil)
	assert.This(len(res.Matches)).Is(4)

	// case sensitive
	res, err = searchBook("mybook", "", "TEXT", true)
	assert.That(err == nil)
	assert.This(len(res.Matches)).Is(0)

	// both path and text required
	_, err = searchBook("mybook", "", "", false)
	assert.That(err != nil)

	// multiple matching lines in same document
	res, err = searchBook("mybook", "MultiMatch", "text", false)
	assert.That(err == nil)
	assert.This(len(res.Matches)).Is(1)
	assert.This(res.Matches[0].Path).Is("/MultiMatch")
	assert.This(res.Matches[0].Lines).Is([]string{"0001: line one has text", "0002: line two has text"})
	assert.This(res.Matches[0].HasMore).Is(false)
}

func TestSearchBookLinesLimit(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create testbook (name, path, text, order) key(name, path)", nil)

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	// Create a document with 8 matching lines (more than linesLimit of 5)
	tran.Action(th, "insert { name: 'ManyMatches', path: '', text: 'line1 match\nline2 match\nline3 match\nline4 match\nline5 match\nline6 match\nline7 match\nline8 match', order: 1 } into testbook")
	tran.Complete()

	res, err := searchBook("testbook", "", "match", false)
	assert.That(err == nil)
	assert.This(len(res.Matches)).Is(1)
	assert.This(len(res.Matches[0].Lines)).Is(5)
	assert.This(res.Matches[0].HasMore).Is(true)
	assert.This(res.Matches[0].Lines[0]).Is("0001: line1 match")
	assert.This(res.Matches[0].Lines[4]).Is("0005: line5 match")
}

func TestBookMatchLine(t *testing.T) {
	assert := assert.T(t)
	// Test addLineNumbers directly
	result := addLineNumbers("matching line here", 5)
	assert.This(result).Is("0005: matching line here")
}
