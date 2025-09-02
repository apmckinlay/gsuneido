// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOptimize(t *testing.T) {
	MakeSuTran = func(qt QueryTran) *core.SuTran { return nil }
	var mode = ReadMode
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)

		query2 := format(0, q, 0)
		q2 := ParseQuery(query2, testTran{}, nil)
		assert.This(format(0, q2, 0)).Is(query2)

		q, _, _ = Setup(q, mode, testTran{})
		// fmt.Println("-----------------------------")
		// fmt.Println(Format(q))
		assert.T(t).Msg(query).This(String(q)).Like(expected)
	}
	// trace.Set(int(trace.QueryOpt))
	// test("table rename a to x, x to y sort y",
	// 	"table^(a) rename a to x, x to y")
	// t.SkipNow()

	test("table",
		"table^(a)")
	test("table sort a",
		"table^(a)")
	test("trans sort date",
		"trans^(date,item,id)")
	test("table sort c",
		"table^(a) tempindex(c)")
	test("hist where date is 5",
		"hist^(date,item,id) where date is 5") // not where*1
	test("comp where a=1 sort c, a, b",
		"comp^(a,b,c) where a is 1 tempindex(c,b)")
	test("comp where a=1 and b=2 and c=3 sort c, a, b",
		"comp^(a,b,c) where*1 a is 1 and b is 2 and c is 3")

	test("supplier",
		"supplier^(supplier)")
	test("supplier where city is 5",
		"supplier^(city) where city is 5")
	test("supplier where supplier is 5",
		"supplier^(supplier) where*1 supplier is 5")
	test("supplier where name is 5",
		"supplier^(supplier) where name is 5")
	test("supplier where name is 5 and city is 5",
		"supplier^(city) where name is 5 and city is 5")
	test("supplier where supplier is 3 and city is 5",
		"supplier^(supplier) where*1 supplier is 3 and city is 5")
	test("supplier where String?(city)",
		"supplier^(city) where String?(city)")
	test("supplier where false",
		"nothing(supplier)")

	test("supplier where Func(name)",
		"supplier^(supplier) where Func(name)")
	test("supplier where city is 5 and Func(name)",
		"supplier^(city) where city is 5 and Func(name)") // previous bug
	test("supplier where Func(name) and city is 5",
		"supplier^(city) where Func(name) and city is 5")
	test("supplier where Func(name).x",
		"supplier^(supplier) where Func(name).x")

	test("table rename b to bb sort c",
		"table^(a) tempindex(c) rename b to bb")
	test("table rename a to x, x to y sort y",
		"table^(a) rename a to x, x to y")

	test("table extend x = F() sort c",
		"table^(a) tempindex(c) extend x = F()")
	test("table extend x = F() sort x",
		"table^(a) extend x = F() tempindex(x)")

	test("table minus table",
		"table^(a) minus(a) table^(a)")
	test("(table extend x = 1) minus hist",
		"table^(a) extend x = 1")

	test("hist intersect hist2",
		"hist^(item) intersect(date) hist2^(date)")
	test("hist2 intersect hist",
		"hist^(item) intersect(date) hist2^(date)")

	test("hist union hist2",
		"hist^(item) union-lookup(date) hist2^(date)")
	test("hist2 union hist",
		"hist^(item) union-lookup(date) hist2^(date)")
	test("hist union hist sort date",
		"hist^(date,item,id) union-merge(date,item,id) hist^(date,item,id)")
	test("table union table",
		"table^(a) union-merge(a) table^(a)")
	test("(table where a is 1) union (table where a is 2)",
		"table^(a) where*1 a is 1 "+
			"union-disjoint(a) (table^(a) where*1 a is 2)")
	test("supplier where supplier > 1 sort city",
		"supplier^(city) where supplier > 1")
	test("supplier where supplier > 9 sort city",
		"supplier^(supplier) where supplier > 9 tempindex(city)")

	test("table project a",
		"table^(a) project-copy a")
	test("table project a sort a",
		"table^(a) project-copy a")
	test("abc project a",
		"abc^(a) project-seq a")
	test("comp project b",
		"comp^(a,b,c) project-map b")
	test("comp where a is 1 and b is 2 project c",
		"comp^(a,b,c) where a is 1 and b is 2 project-copy c")
	test("customer project id, name",
		"customer^(id) project-copy id, name")
	test("trans project item",
		"trans^(item) project-seq item")
	test("trans project item,id,cost,date project item",
		"trans^(item) project-seq item")
	test("trans project item,id,cost project item,id project item",
		"trans^(item) project-seq item")
	test("hist project date,item",
		"hist^(date,item,id) project-seq date, item")
	test("customer project city",
		"customer^(id) project-map city")
	test("customer project id,city project city",
		"customer^(id) project-map city")
	test("customer project city sort city",
		"customer^(id) project-map city tempindex(city)")

	test("trans summarize total cost", // by is empty
		"trans^(date,item,id) summarize-seq total cost")
	test("trans summarize total cost sort total_cost", // ignore sort
		"trans^(date,item,id) summarize-seq total cost")
	test("trans summarize item, total cost",
		"trans^(item) summarize-seq item, total cost")
	test("trans summarize item, total cost sort total_cost",
		"trans^(item) summarize-seq item, total cost"+
			" tempindex(total_cost)")
	test("supplier summarize max supplier", // key
		"supplier^(supplier) summarize-idx* max supplier")
	test("supplier summarize max supplier sort name", // ignore sort
		"supplier^(supplier) summarize-idx* max supplier")
	test("supplier summarize max city", // index
		"supplier^(city) summarize-idx max city")
	// hints
	test("hist summarize id, total cost",
		"hist^(item) summarize-map id, total cost")
	test("hist summarize/*small*/ id, total cost",
		"hist^(item) summarize-map id, total cost")
	test("hist summarize/*large*/ id, total cost",
		"hist^(item) tempindex(id) summarize-seq id, total cost")
	test("trans summarize id, count",
		"trans^(date,item,id) tempindex(id) summarize-seq id, count")
	test("trans summarize/*large*/ id, count",
		"trans^(date,item,id) tempindex(id) summarize-seq id, count")
	test("trans summarize/*small*/ id, count",
		"trans^(date,item,id) summarize-map id, count")

	test("customer times inven",
		"customer^(id) times inven^(item)")
	test("inven times customer sort id",
		"customer^(id) times inven^(item)")
	test("inven times (customer where false)",
		"nothing")

	test("hist join customer",
		"hist^(item) join n:1 by(id) customer^(id)")
	test("customer join hist",
		"hist^(item) join n:1 by(id) customer^(id)")
	test("trans join inven",
		"inven^(item) join 1:n by(item) trans^(item)")
	test("task join co",
		"task^(tnum) join 1:1 by(tnum) co^(tnum)")
	test("customer join alias",
		"alias^(id) join 1:1 by(id) customer^(id)")
	test("(inven join trans) union (inven join trans)",
		"(inven^(item) join 1:n by(item) trans^(item)) "+
			"union-lookup(date,item,id) "+
			"(trans^(date,item,id) join n:1 by(item) inven^(item))")
	test("task join co join cus",
		"(task^(tnum) join 1:1 by(tnum) co^(tnum)) "+
			"join n:1 by(cnum) cus^(cnum)")
	test("trans join inven",
		"inven^(item) join 1:n by(item) trans^(item)")

	test("(trans union trans) join (inven union inven)",
		"(trans^(date,item,id) union-merge(date,item,id) trans^(date,item,id)) "+
			"join n:n by(item) "+
			"(inven^(item) tempindex(item) "+
			"union-merge(item) (inven^(item) tempindex(item)))")

	test("inven leftjoin trans",
		"inven^(item) leftjoin 1:n by(item) trans^(item)")
	test("customer leftjoin hist2",
		"customer^(id) leftjoin 1:n by(id) hist2^(id)")
	test("customer leftjoin hist2 sort date",
		"(customer^(id) leftjoin 1:n by(id) hist2^(id)) tempindex(date)")

	test("hist2 where date > 1 sort id",
		"hist2^(id) where date > 1")
	test("hist2 where date is 1 sort id",
		"hist2^(date) where*1 date is 1")

	test("comp where a = 1 sort b",
		"comp^(a,b,c) where a is 1")

	test("((inven extend x=1) where x is 2) union inven",
		`inven^(item) extend x = ""`)

	test("inven extend x = 0, c = item[1..2] where c > 0",
		"inven^(item) where item[1..2] > 0 extend x = 0, c = item[1..2]")

	test("inven extend x = 0, c = item[..2] where c > 0",
		"inven^(item) where item[..2] > 0 extend x = 0, c = item[..2]")

	test("inven extend x = 0, c = item[1..] where c > 0",
		"inven^(item) where item[1..] > 0 extend x = 0, c = item[1..]")
		
	test("comp2 where a = 1 and b = 2 and c = 3",
		"comp2^(a,b,c) where*1 a is 1 and b is 2 and c is 3")

	mode = CursorMode
	test("(trans union trans) join (inven union inven)",
		"(trans^(date,item,id) union-merge(date,item,id) trans^(date,item,id)) "+
			"join n:n by(item) "+
			"(inven^(item) union-merge(item) inven^(item))")
	test("(inven join trans) union (inven join trans)",
		"(inven^(item) join 1:n by(item) trans^(item)) "+
			"union-lookup(date,item,id) "+
			"(trans^(date,item,id) join n:1 by(item) inven^(item))")
	test("trans join customer",
		"trans^(date,item,id) join n:1 by(id) customer^(id)")
	test("trans join inven join customer",
		"(inven^(item) join 1:n by(item) trans^(item)) "+
			"join n:1 by(id) customer^(id)")
	assert.T(t).This(func() { test("table rename b to bb sort c", "") }).
		Panics("invalid query")
}
