// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTablesTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	dbms := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbms }
	query.DoAdmin(db, `create alpha (a, b) key(a)`, nil)
	query.DoAdmin(db, `create beta (x, y) key(x)`, nil)
	query.DoAdmin(db, `create gamma (m, n) key(m)`, nil)

	output, err := tablesTool("")
	assert.That(err == nil)
	assert.This(output.Tables).Is([]string{"alpha", "beta", "columns", "gamma", "indexes", "tables", "views"})
	assert.That(output.HasMore == false)

	output, err = tablesTool("b")
	assert.That(err == nil)
	assert.This(output.Tables).Is([]string{"beta"})
	assert.That(output.HasMore == false)
}
