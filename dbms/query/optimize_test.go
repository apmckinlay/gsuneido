// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOptimize(t *testing.T) {
	var mode = ReadMode
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)
		q, _, _ = Setup(q, mode, testTran{})
		// fmt.Println("-----------------------------")
		// fmt.Println(Format(q))
		assert.T(t).Msg(query).This(q.String()).Like(expected)
	}
	// trace.Set(int(trace.QueryOpt))
	// test("table sort c",
	// 	"")
	// t.SkipNow()

	test("table",
		"table^(a)")
	test("table sort a",
		"table^(a)")
	test("trans sort date",
		"trans^(date,item,id)")
	test("table sort c",
		"table^(a) TEMPINDEX(c)")
	test("hist where date is 5",
		"hist^(date) WHERE date is 5") // not WHERE*1

	test("table rename b to bb sort c",
		"table^(a) TEMPINDEX(c) RENAME b to bb")

	test("table extend x = F() sort c",
		"table^(a) TEMPINDEX(c) EXTEND x = F()")
	test("table extend x = F() sort x",
		"table^(a) EXTEND x = F() TEMPINDEX(x)")

	test("table minus table",
		"table^(a) MINUS table^(a)")
	test("(table extend x = 1) minus hist",
		"table^(a) EXTEND x = 1")

	test("hist intersect hist2",
		"hist^(date) INTERSECT hist2^(date)")
	test("hist2 intersect hist",
		"hist^(date) INTERSECT hist2^(date)")

	test("hist union hist2",
		"hist^(date) UNION-LOOKUP(date) hist2^(date)")
	test("hist2 union hist",
		"hist^(date) UNION-LOOKUP(date) hist2^(date)")
	test("hist union hist sort date",
		"hist^(date,item,id) UNION-MERGE(date,item,id) hist^(date,item,id)")
	test("table union table",
		"table^(a) UNION-MERGE(a) table^(a)")
	test("(table where a is 1) union (table where a is 2)",
		"table^(a) WHERE*1 a is 1 "+
			"UNION-DISJOINT(a) (table^(a) WHERE*1 a is 2)")
	test("supplier where supplier > 1 sort city",
		"supplier^(city) WHERE supplier > 1")
	test("supplier where supplier > 9 sort city",
		"supplier^(supplier) WHERE supplier > 9 TEMPINDEX(city)")

	test("table project a",
		"table^(a) PROJECT-COPY a")
	test("table project a sort a",
		"table^(a) PROJECT-COPY a")
	test("abc project a",
		"abc^(a) PROJECT-SEQ a")
	test("comp project b",
		"comp^(a,b,c) PROJECT-MAP b")
	test("comp where a is 1 and b is 2 project c",
		"comp^(a,b,c) WHERE a is 1 and b is 2 PROJECT-COPY c")
	test("customer project id,name",
		"customer^(id) PROJECT-COPY id,name")
	test("trans project item",
		"trans^(item) PROJECT-SEQ item")
	test("trans project item,id,cost,date project item",
		"trans^(item) PROJECT-SEQ item")
	test("trans project item,id,cost project item,id project item",
		"trans^(item) PROJECT-SEQ item")
	test("hist project date,item",
		"hist^(date,item,id) PROJECT-SEQ date,item")
	test("customer project city",
		"customer^(id) PROJECT-MAP city")
	test("customer project id,city project city",
		"customer^(id) PROJECT-MAP city")
	test("customer project city sort city",
		"customer^(id) PROJECT-MAP city TEMPINDEX(city)")

	test("trans summarize total cost", // by is empty
		"trans^(date,item,id) SUMMARIZE-SEQ total_cost = total cost")
	test("trans summarize total cost sort total_cost", // ignore sort
		"trans^(date,item,id) SUMMARIZE-SEQ total_cost = total cost")
	test("trans summarize item, total cost",
		"trans^(item) SUMMARIZE-SEQ item, total_cost = total cost")
	test("trans summarize item, total cost sort total_cost",
		"trans^(item) SUMMARIZE-SEQ item, total_cost = total cost"+
			" TEMPINDEX(total_cost)")
	test("trans summarize id, total cost",
		"trans^(date,item,id) SUMMARIZE-MAP id, total_cost = total cost")
	test("supplier summarize max supplier", // key
		"supplier^(supplier) SUMMARIZE-IDX* max_supplier = max supplier")
	test("supplier summarize max supplier sort name", // ignore sort
		"supplier^(supplier) SUMMARIZE-IDX* max_supplier = max supplier")
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
		"task^(tnum) JOIN 1:1 by(tnum) co^(tnum)")
	test("customer join alias",
		"alias^(id) JOIN 1:1 by(id) customer^(id)")
	test("(inven join trans) union (inven join trans)",
		"(inven^(item) JOIN 1:n by(item) trans^(item)) "+
			"UNION-LOOKUP(date,item,id) "+
			"(trans^(date,item,id) JOIN n:1 by(item) inven^(item))")
	test("task join co join cus",
		"(task^(tnum) JOIN 1:1 by(tnum) co^(tnum)) "+
			"JOIN n:1 by(cnum) cus^(cnum)")
	test("trans join inven",
		"inven^(item) JOIN 1:n by(item) trans^(item)")

	test("(trans union trans) join (inven union inven)",
		"(trans^(date,item,id) UNION-MERGE(date,item,id) trans^(date,item,id)) "+
			"JOIN n:n by(item) "+
			"(inven^(item) TEMPINDEX(item) "+
			"UNION-MERGE(item) (inven^(item) TEMPINDEX(item)))")

	test("inven leftjoin trans",
		"inven^(item) LEFTJOIN 1:n by(item) trans^(item)")
	test("customer leftjoin hist2",
		"customer^(id) LEFTJOIN 1:n by(id) hist2^(id)")
	test("customer leftjoin hist2 sort date",
		"(customer^(id) LEFTJOIN 1:n by(id) hist2^(id)) TEMPINDEX(date)")

	test("hist2 where date > 1 sort id",
		"hist2^(id) WHERE date > 1")
	test("hist2 where date is 1 sort id",
		"hist2^(date) WHERE*1 date is 1")

	test("comp where a = 1 sort b",
		"comp^(a,b,c) WHERE a is 1")

	mode = CursorMode
	test("(trans union trans) join (inven union inven)",
		"(trans^(date,item,id) UNION-MERGE(date,item,id) trans^(date,item,id)) "+
			"JOIN n:n by(item) "+
			"(inven^(item) UNION-MERGE(item) inven^(item))")
	test("(inven join trans) union (inven join trans)",
		"(inven^(item) JOIN 1:n by(item) trans^(item)) "+
			"UNION-LOOKUP(date,item,id) "+
			"(trans^(date,item,id) JOIN n:1 by(item) inven^(item))")
	test("trans join customer",
		"trans^(date,item,id) JOIN n:1 by(id) customer^(id)")
	test("trans join inven join customer",
		"(inven^(item) JOIN 1:n by(item) trans^(item)) "+
			"JOIN n:1 by(id) customer^(id)")
	assert.T(t).This(func() { test("table rename b to bb sort c", "") }).
		Panics("invalid query")
}
