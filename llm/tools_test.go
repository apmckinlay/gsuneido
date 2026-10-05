// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func testToolThread() *core.Thread {
	th := core.NewThread(nil)
	th.SetPerms(core.AllPerms)
	return th
}

func testToolContext() context.Context {
	return context.WithValue(context.Background(), toolThreadKey{}, testToolThread())
}

func TestGetViewDefinition(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	d := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return d }
	query.DoAdminTest(db, `create alpha (a, b) key(a)`)
	query.DoAdminTest(db, `view myview = alpha extend c = 123`)

	// Test view definition
	viewDef, _ := getSchema(testToolContext(), "myview")
	assert.This(viewDef).Is(schemaOutput{Schema: "view myview = alpha extend c = 123"})

	// Test non-existent view
	noView, _ := getSchema(testToolContext(), "nonexistent")
	assert.This(noView).Is(schemaOutput{Schema: ""})
}

func TestSchemaToolWithView(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	d := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return d }
	query.DoAdminTest(db, `create alpha (a, b) key(a)`)
	query.DoAdminTest(db, `view myview = alpha extend c = 123`)

	// Test table schema
	schema, _ := getSchema(testToolContext(), "alpha")
	assert.This(schema).Is(schemaOutput{Schema: "alpha (a,b) key(a)"})

	// Test view definition via schema tool
	viewDef, _ := getSchema(testToolContext(), "myview")
	assert.This(viewDef).Is(schemaOutput{Schema: "view myview = alpha extend c = 123"})
}
