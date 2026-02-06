// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSearchTool(t *testing.T) {
	assert := assert.T(t)
	dd := db19.CreateDb(stor.HeapStor(8192))
	dd.CheckerSync()
	db19.StartConcur(dd, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(dd)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create stdlib (name, text, group, parent, num) key(num) key(name, group) index(parent, name) index(group)", nil)
	dbmsLocal.Admin("create app (name, text, group, parent, num) key(num) key(name, group) index(parent, name) index(group)", nil)
	assert.That(dbmsLocal.Use("app"))

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	tran.Action(th, "insert { name: 'Foo', text: 'function(){return 1}', group: -1, parent: 0, num: 1 } into stdlib")
	tran.Action(th, "insert { name: 'Bar', text: 'function(){return 2}', group: -1, parent: 0, num: 2 } into stdlib")
	tran.Action(th, "insert { name: 'FooApp', text: 'function(){return 3}', group: -1, parent: 0, num: 3 } into app")
	tran.Action(th, "insert { name: 'Folder', group: 0, parent: 0, num: 4 } into app")
	tran.Action(th, "insert { name: 'SearchTarget', text: 'function(){return \"hello\"}', group: -1, parent: 4, num: 5 } into app")
	tran.Complete()

	res, err := searchTool("std.*", "Foo", "return 1", false)
	assert.That(err == nil)
	assert.This(res.Matches).Is([]codeMatch{{Library: "stdlib", Name: "Foo", Path: "", Line: "0001: function(){return 1}"}})

	res, err = searchTool("STDLIB", "FOO", "RETURN 1", false)
	assert.That(err == nil)
	assert.This(res.Matches).Is([]codeMatch{{Library: "stdlib", Name: "Foo", Path: "", Line: "0001: function(){return 1}"}})

	res, err = searchTool("STDLIB", "FOO", "RETURN 1", true)
	assert.That(err == nil)
	assert.This(res.Matches).Is([]codeMatch{})

	res, err = searchTool("app", "Foo.*", "return 3", false)
	assert.That(err == nil)
	assert.This(res.Matches).Is([]codeMatch{{Library: "app", Name: "FooApp", Path: "", Line: "0001: function(){return 3}"}})

	res, err = searchTool("", "", "return", false)
	assert.That(err == nil)
	assert.This(res.Matches).Is([]codeMatch{
		{Library: "stdlib", Name: "Bar", Path: "", Line: "0001: function(){return 2}"},
		{Library: "stdlib", Name: "Foo", Path: "", Line: "0001: function(){return 1}"},
		{Library: "app", Name: "FooApp", Path: "", Line: "0001: function(){return 3}"},
		{Library: "app", Name: "SearchTarget", Path: "Folder", Line: "0001: function(){return \"hello\"}"},
	})
	assert.That(!res.HasMore)

	_, err = searchTool("", "", "", false)
	assert.That(err != nil)
}

func TestSearchQuery_InvalidRegex(t *testing.T) {
	assert := assert.T(t)
	_, err := searchQuery("[", "", false)
	assert.That(err != nil)
	assert.That(strings.Contains(err.Error(), "invalid name regex"))
}

func TestFilterLibraries_InvalidRegex(t *testing.T) {
	assert := assert.T(t)
	_, err := filterLibraries([]string{"stdlib"}, "[", false)
	assert.That(err != nil)
}
