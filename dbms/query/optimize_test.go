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
	test("tables project table",
		"tables^(table) PROJECT-COPY table")
	test("tables project tablename sort tablename",
		"tables^(tablename) PROJECT-COPY tablename")
	test("abc project a",
		"abc^(a) PROJECT-SEQ a")
	test("columns project column",
		"columns^(table,column) PROJECT-LOOKUP column")
	test("columns where table is 1 project column",
		"(columns^(table,column) WHERE table is 1) PROJECT-COPY column")

	mode = updateMode
	test("table rename b to bb sort c",
		"table^(a) TEMPINDEX(c) RENAME b to bb")

	mode = cursorMode
	assert.T(t).This(func() { test("table rename b to bb sort c", "") }).
		Panics("invalid query")
}

var testInfo = map[string]*meta.Info{
	"hist2":   {Nrows: 1000, Size: 100000},
}

func (testTran) GetInfo(table string) *meta.Info {
	if ti, ok := testInfo[table]; ok {
		return ti
	}
	return &meta.Info{Nrows: 100, Size: 10000}
}
