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
	test("tables",
		"tables^(table)")
	test("tables sort tablename",
		"tables^(tablename)")
	test("table rename b to bb sort c",
		"table^(a) TEMPINDEX(c) RENAME b to bb")
	test("table extend x = F() sort c",
		"table^(a) TEMPINDEX(c) EXTEND x = F()")
	test("table extend x = F() sort x",
		"table^(a) EXTEND x = F() TEMPINDEX(x)")
	test("table minus table",
		"table^(a) MINUS table^(a)")
	test("hist intersect hist2",
		"hist^(date,item,id) INTERSECT hist2^(date)")
	test("hist2 intersect hist",
		"hist^(date,item,id) INTERSECT hist2^(date)")
	test("hist union hist2",
		"hist^(date,item,id) UNION-LOOKUP hist2^(date)")
	test("hist2 union hist",
		"hist^(date,item,id) UNION-LOOKUP hist2^(date)")
	test("hist union hist sort date",
		"hist^(date,item,id) UNION-MERGE hist^(date,item,id)")
	test("table union table",
		"table^(a) UNION-MERGE table^(a)")
	test("(table where a is 1) union (table where a is 2)",
		"(table^(a) WHERE a is 1) UNION-FOLLOW-DISJOINT(a) (table^(a) WHERE a is 2)")

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
	"hist":   {Nrows: 100, Size: 10000},
	"hist2":  {Nrows: 1000, Size: 100000},
}

func (testTran) GetInfo(table string) *meta.Info {
	return testInfo[table]
}
