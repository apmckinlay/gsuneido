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

func TestCodeFoldersTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	db19.StartConcur(db, 50*time.Millisecond)
	dbmsLocal := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbmsLocal }

	dbmsLocal.Admin("create stdlib (name, text, group, parent, num) key(num) key(name, group) index(parent, name) index(group)", nil)

	th := core.NewThread(core.MainThread)
	tran := dbmsLocal.Transaction(true)
	tran.Action(th, "insert { name: 'Folder1', group: 0, parent: 0, num: 1 } into stdlib")
	tran.Action(th, "insert { name: 'Folder2', group: 0, parent: 0, num: 2 } into stdlib")
	tran.Action(th, "insert { name: 'Sub', group: 1, parent: 1, num: 3 } into stdlib")
	tran.Action(th, "insert { name: 'Leaf', group: -1, parent: 0, num: 4, text: 'function(){}' } into stdlib")
	tran.Action(th, "insert { name: 'Child', group: -1, parent: 1, num: 5, text: 'function(){}' } into stdlib")
	tran.Complete()

	res, err := codeFoldersTool("stdlib", "")
	assert.That(err == nil)
	assert.This(res.Library).Is("stdlib")
	assert.This(res.Path).Is("")
	assert.This(res.Children).Is([]string{"Folder1/", "Folder2/", "Leaf"})

	res, err = codeFoldersTool("stdlib", "Folder1")
	assert.That(err == nil)
	assert.This(res.Path).Is("Folder1")
	assert.This(res.Children).Is([]string{"Child", "Sub/"})

	res, err = codeFoldersTool("stdlib", "Folder1/")
	assert.That(err == nil)
	assert.This(res.Path).Is("Folder1")
	assert.This(res.Children).Is([]string{"Child", "Sub/"})

	_, err = codeFoldersTool("stdlib", "Leaf")
	assert.That(err != nil)
	assert.This(err.Error()).Is("path segment is not a folder: Leaf")

	_, err = codeFoldersTool("stdlib", "Missing")
	assert.That(err != nil)
	assert.This(err.Error()).Is("path not found: Missing")
}
