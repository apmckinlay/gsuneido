// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestOptimize2(t *testing.T) {
	MakeSuTran = func(qt QueryTran) *core.SuTran { return nil }
	var mode = ReadMode
	test2 := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{}, nil)
		q = q.Transform()

		fixcost, varcost := Optimize2(q, mode, reqUnordered, 1)
		if fixcost+varcost >= impossible {
			panic("invalid query: " + String(q))
		}
		q = SetApproach2(q, reqUnordered, 1, testTran{})

		assert.T(t).Msg(query).This(String(q)).Like(expected)
	}

	// leaf / pass-through: v2 fallback produces same results as v1

	test2("table",
		"table^(a)")
	test2("table sort a",
		"table^(a)")
	test2("table sort c",
		"table^(a) tempindex(c)")

	test2("supplier",
		"supplier^(supplier)")
	test2("supplier where false",
		"nothing(supplier)")

	test2("supplier where supplier is 5",
		"supplier^(supplier) where*1 supplier is 5")
	test2("supplier where city is 5",
		"supplier^(city,supplier) where city is 5")
	test2("supplier where name is 5",
		"supplier^(supplier) where name is 5")

	test2("hist where date is 5",
		"hist^(date,item,id) where date is 5")

	test2("comp where a=1 sort c, a, b",
		"comp^(a,b,c) where a is 1 tempindex(c,b)")

	test2("table project a",
		"table^(a) project-copy a")
	test2("abc project a",
		"abc^(a,b) project-seq a")
	test2("comp project b",
		"comp^(a,b,c) project-map b")

	test2("trans summarize total cost",
		"trans^(date,item,id) summarize-seq total cost")
	test2("trans summarize item, total cost",
		"trans^(item,date,id) summarize-seq item, total cost")
	test2("supplier summarize max supplier",
		"supplier^(supplier) summarize-idx* max supplier")
	test2("supplier summarize max city",
		"supplier^(city,supplier) summarize-idx max city")

	test2("customer times inven",
		"customer^(id) times inven^(item)")
	test2("inven times customer sort id",
		"customer^(id) times inven^(item)")

	test2("table minus table",
		"table^(a) minus(a) table^(a)")
	test2("hist intersect hist2",
		"hist^(date,item,id) intersect(date) hist2^(date)")
	test2("hist union hist2",
		"hist^(date,item,id) union-lookup(date) hist2^(date)")

	test2("hist join customer",
		"hist^(date,item,id) join n:1 by(id) customer^(id)")
	test2("customer join hist",
		"hist^(date,item,id) join n:1 by(id) customer^(id)")
	test2("customer join alias",
		"alias^(id) join 1:1 by(id) customer^(id)")

	test2("inven leftjoin trans",
		"inven^(item) leftjoin 1:n by(item) trans^(item,date,id)")
	test2("customer leftjoin alias",
		"customer^(id) leftjoin 1:1 by(id) alias^(id)")

	test2("inven semijoin trans",
		"inven^(item) semijoin by(item) trans^(item,date,id)")

	test2("table rename b to bb sort c",
		"table^(a) tempindex(c) rename b to bb")
	test2("table rename a to x, x to y sort y",
		"table^(a) rename a to x, x to y")

	test2("table extend x = F() sort c",
		"table^(a) tempindex(c) extend x = F()")
	test2("table extend x = F() sort x",
		"table^(a) extend x = F() tempindex(x)")
	test2("((inven extend x=1) where x is 2) union inven",
		`inven^(item) extend x = ""`)
}