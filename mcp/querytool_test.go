// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package mcp

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestQueryTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	dbms := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbms }

	result, err := queryTool("tables join columns")
	assert.That(err == nil)
	assert.That(strings.Contains(result.Results, `["table", "nrows", "totalsize", "column", "field"]`))
	assert.That(strings.Contains(result.Results, `["tables", 4, 0, "table", 0]`))
}

func TestFormatQueryResult(t *testing.T) {
	assert := assert.T(t)
	cols := []string{"col1", "col2", "col3"}
	rows := [][]core.Value{
		{core.IntVal(1), core.SuStr("hello"), nil},
		{core.DateFromLiteral("#20260203"), core.True, core.SuObjectOf(core.IntVal(7))},
	}
	got := formatQueryResult(cols, rows, false)
	assert.This(got).Is("[\n" +
		"[\"col1\", \"col2\", \"col3\"]\n" +
		"[1, \"hello\", null]\n" +
		"[#20260203, true, #(7)]\n" +
		"]\n")
}

func TestFormatQueryResult_Truncated(t *testing.T) {
	assert := assert.T(t)
	cols := []string{"col1"}
	rows := [][]core.Value{{core.IntVal(1)}}
	got := formatQueryResult(cols, rows, true)
	assert.This(got).Is("[\n" +
		"[\"col1\"]\n" +
		"[1]\n" +
		"// truncated\n" +
		"]\n")
}
