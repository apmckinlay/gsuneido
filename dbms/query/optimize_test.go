// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOptimize(t *testing.T) {
	var mode Mode
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query)
		Setup(q, mode, testTran{})
		assert.T(t).Msg(query).This(q.String()).Is(expected)
	}
	mode = readMode
	test("tables",
		"tables^(table)")
	test("tables sort tablename",
		"tables^(tablename)")
	test("table sort c",
		"table^(a) TEMPINDEX(c)")

	test("table rename b to bb sort c",
		"table^(a) TEMPINDEX(c) RENAME b to bb")

	test("table extend x = F() sort c",
		"table^(a) TEMPINDEX(c) EXTEND x = F()")
	test("table extend x = F() sort x",
		"table^(a) EXTEND x = F() TEMPINDEX(x)")

	test("table minus table",
		"table^(a) MINUS table^(a)")

	test("hist intersect hist2",
		"hist^(date) INTERSECT hist2^(date)")
	test("hist2 intersect hist",
		"hist^(date) INTERSECT hist2^(date)")

	test("hist union hist2",
		"hist^(date) UNION-LOOKUP hist2^(date)")
	test("hist2 union hist",
		"hist^(date) UNION-LOOKUP hist2^(date)")
	test("hist union hist sort date",
		"hist^(date,item,id) UNION-MERGE hist^(date,item,id)")
	test("table union table",
		"table^(a) UNION-MERGE table^(a)")
	// test("(table where a is 1) union (table where a is 2)",
	// 	"table^(a) WHERE UNION-FOLLOW-DISJOINT(a) (table^(a) WHERE)")

	test("tables project table",
		"tables^(table) PROJECT-COPY table")
	test("tables project tablename sort tablename",
		"tables^(tablename) PROJECT-COPY tablename")
	test("abc project a",
		"abc^(a) PROJECT-SEQ a")
	test("columns project column",
		"columns^(table,column) PROJECT-LOOKUP column")
	// test("columns where table is 1 project column",
	// 	"columns^(table,column) WHERE PROJECT-COPY column")
	test("customer project id,name",
		"customer^(id) PROJECT-COPY id, name")
	test("trans project item",
		"trans^(item) PROJECT-SEQ item")
	test("trans project item,id,cost,date project item",
		"trans^(item) PROJECT-SEQ item")
	test("trans project item,id,cost project item,id project item",
		"trans^(item) PROJECT-SEQ item")
	test("hist project date,item",
		"hist^(date,item,id) PROJECT-SEQ date, item")
	test("customer project city",
		"customer^(id) PROJECT-LOOKUP city")
	test("customer project id,city project city",
		"customer^(id) PROJECT-LOOKUP city")

	test("trans summarize total cost", // by is empty
		"trans^(date,item,id) SUMMARIZE-SEQ total_cost = total cost")
	test("trans summarize total cost sort total_cost", // ignore sort
		"trans^(date,item,id) SUMMARIZE-SEQ total_cost = total cost")
	test("trans summarize item, total cost",
		"trans^(item) SUMMARIZE-SEQ item, total_cost = total cost")
	test("trans summarize id, total cost",
		"trans^(date,item,id) SUMMARIZE-MAP id, total_cost = total cost")
	test("supplier summarize max supplier", // key
		"supplier^(supplier) SUMMARIZE-IDX max_supplier = max supplier")
	test("supplier summarize max supplier sort name", // ignore sort
		"supplier^(supplier) SUMMARIZE-IDX max_supplier = max supplier")
	test("supplier summarize max city", // index
		"supplier^(city) SUMMARIZE-IDX max_city = max city")

	test("customer times inven",
		"customer^(id) TIMES inven^(item)")
	test("inven times customer sort id",
		"customer^(id) TIMES inven^(item)")

	test("hist join customer",
		"hist^(date) JOIN n:1 by(id) customer^(id)")
	test("customer join hist",
		"hist^(date) JOIN n:1 by(id) customer^(id)")
	test("trans join inven",
		"inven^(item) JOIN 1:n by(item) trans^(item)")
	test("task join co",
		"co^(tnum) JOIN 1:1 by(tnum) task^(tnum)")
	test("customer join alias",
		"alias^(id) JOIN 1:1 by(id) customer^(id)")
	// test("(trans join customer) union (hist join customer)",
	// 	"(trans^(date,item,id) JOIN n:1 by(id) customer^(id)) "+
	// 		"UNION-MERGE^(date,item,id) (hist^(date,item,id) "+
	// 		"JOIN n:1 by(id) customer^(id)")
	// test("(trans join customer) intersect (hist join customer)",
	// 	"(trans^(date,item,id) JOIN n:1 by(id) customer^(id)) "+
	// 		"INTERSECT (hist^(date,item,id) "+
	// 		"JOIN n:1 by(id) customer^(id))")
	test("task join co join cus",
		"(co^(tnum) JOIN 1:1 by(tnum) task^(tnum)) "+
			"JOIN n:1 by(cnum) cus^(cnum)")
	test("inven leftjoin trans",
		"inven^(item) LEFTJOIN 1:n by(item) trans^(item)")
	test("customer leftjoin hist2",
		"customer^(id) LEFTJOIN 1:n by(id) hist2^(id)")
	test("customer leftjoin hist2 sort date",
		"(customer^(id) LEFTJOIN 1:n by(id) hist2^(id)) TEMPINDEX(date)")

	mode = updateMode
	test("table rename b to bb sort c",
		"table^(a) TEMPINDEX(c) RENAME b to bb")

	mode = cursorMode
	test("trans join customer",
		"trans^(date,item,id) JOIN n:1 by(id) customer^(id)")
	test("trans join inven join customer",
		"(inven^(item) JOIN 1:n by(item) trans^(item)) "+
			"JOIN n:1 by(id) customer^(id)")
	assert.T(t).This(func() { test("table rename b to bb sort c", "") }).
		Panics("invalid query")
}
