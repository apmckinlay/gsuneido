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

func TestGetViewDefinition(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	d := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return d }
	query.DoAdmin(db, `create alpha (a, b) key(a)`, nil)
	query.DoAdmin(db, `view myview = alpha extend c = 123`, nil)

	// Test view definition
	viewDef := getViewDefinition("myview")
	assert.This(viewDef).Is("alpha extend c = 123")

	// Test non-existent view
	noView := getViewDefinition("nonexistent")
	assert.This(noView).Is("")
}

func TestSchemaToolWithView(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	d := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return d }
	query.DoAdmin(db, `create alpha (a, b) key(a)`, nil)
	query.DoAdmin(db, `view myview = alpha extend c = 123`, nil)

	// Test table schema
	schema := core.GetDbms().Schema("alpha")
	assert.This(schema).Is("alpha (a,b) key(a)")

	// Test view definition via schema tool
	viewDef := getViewDefinition("myview")
	assert.This(viewDef).Is("alpha extend c = 123")
}
