// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTablesTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	dbmsLocal := dbms.NewDbmsLocal(db)
	query.DoAdminTest(db, `create alpha (a, b) key(a)`)
	query.DoAdminTest(db, `create beta (x, y) key(x)`)
	query.DoAdminTest(db, `create gamma (m, n) key(m)`)

	output, err := tablesTool(testToolContext(dbmsLocal), "")
	assert.That(err == nil)
	assert.This(output.Tables).Is([]string{"alpha", "beta", "columns", "dbstats",
		"gamma", "indexes", "tables", "views"})
	assert.That(output.HasMore == false)

	output, err = tablesTool(testToolContext(dbmsLocal), "b")
	assert.That(err == nil)
	assert.This(output.Tables).Is([]string{"beta"})
	assert.That(output.HasMore == false)
}
