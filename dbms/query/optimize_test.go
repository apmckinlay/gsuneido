// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOptimize(t *testing.T) {
	var mode Mode
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query)
		Setup(q, mode, testTran{})
		assert.T(t).This(q.String()).Is(expected)
	}
	mode = readMode
	test("tables", "tables^(table)")
	test("tables sort tablename", "tables^(tablename)")
	test("table rename b to bb sort c",
		"table^(a) TEMPINDEX(c) RENAME b to bb")

	mode = updateMode
	test("table rename b to bb sort c",
		"table^(a) TEMPINDEX(c) RENAME b to bb")

	mode = cursorMode
	assert.T(t).This(func() { test("table rename b to bb sort c", "") }).
		Panics("invalid query")
}

var testInfo = map[string]*meta.Info{
	"tables": {Nrows: 100, Size: 10000},
	"table":  {Nrows: 100, Size: 10000},
}

func (testTran) GetInfo(table string) *meta.Info {
	return testInfo[table]
}
