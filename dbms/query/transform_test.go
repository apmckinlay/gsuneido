// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTransform(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	test := func(from, expected string) {
		t.Helper()
		if expected == "" {
			expected = from
		}
		q := ParseQuery(from, testTran{}, nil)
		q = q.Transform()
		actual := String(q)
		// *1 depends on whether optInit runs, e.g. if Nrows is called
		actual = strings.ReplaceAll(actual, "where*1", "where")
		assert.T(t).This(actual).Is(expected)
	}

	test("table", "")

	// TablesLookup
	test("tables where table is 'foo'", "tables(foo)")

	test("table rename a to x, c to y", "")
	test("table remove c, d, e",
		"table project a, b")
	test("table remove x, y, z",
		"table")
	test("table project a, b, c",
		"table")
	test("withdeps remove b",
		"withdeps project a, c, c_deps")
	test("withdeps remove b_deps, c_deps",
		"withdeps project a, b, c")
	test("withdeps rename b to bb, c to cc project a, bb",
		"withdeps project a, b, b_deps rename b to bb, b_deps to bb_deps")

	// combine extend's
	test("customer extend a = 5 extend b = 6",
		"customer extend a = 5, b = 6")
	test("customer extend a = 5 extend b = 6 where id > 5",
		"customer where id > 5 extend a = 5, b = 6")
	// combine project's
	test("table project a, b project b",
		"table project /*NOT UNIQUE*/ b")
	test("customer project id, name project id",
		"customer project id")
	test("customer project id, name where id > 5 project id",
		"customer where id > 5 project id")
	test("customer project id, name project id where id > 5",
		"customer where id > 5 project id")
	// combine rename's
	test("table rename a to x rename b to y rename c to z",
		"table rename a to x, b to y, c to z")
	test("table rename a to x rename x to y rename y to z",
		"table rename a to x, x to y, y to z")
	test("table rename a to aa, b to bb rename bb to b, aa to a",
		"table rename a to aa, b to bb, bb to b, aa to a")
	test("table rename a to x rename c to a",
		"table rename a to x, c to a")
	test("customer rename id to x rename name to y",
		"customer rename id to x, name to y")
	test("table rename a to x rename x to y",
		"table rename a to x, x to y")
	// combine where's
	test("customer where id is 5 where city is 6 where name is 7",
		"customer where id is 5 and city is 6 and name is 7")
	// leftjoin to join
	test("(cus leftjoin task) where cnum is 1 where tnum is 2",
		"cus where cnum is 1 join 1:1 by(cnum) (task where cnum is 1 and tnum is 2)")
	test("(cus leftjoin task) where tnum is 2 where cnum is 1",
		"cus where cnum is 1 join 1:1 by(cnum) (task where tnum is 2 and cnum is 1)")

	// remove projects of all fields
	test("customer project id, city, name", "customer")
	// remove disjoint difference
	test("(customer where id is 3) minus (customer where id is 5)",
		"customer where id is 3")
	// remove empty extends
	test("customer extend zone = 3 project id, city",
		"customer project id, city")
	// remove empty renames
	test("customer rename name to nom project id, city",
		"customer project id, city")

	// move project before rename
	test("customer rename id to num, name to nom project num, city",
		"customer project id, city rename id to num")
	test("customer rename id to id2, id2 to id3 remove id3",
		"customer project /*NOT UNIQUE*/ name, city")
	// move project before rename & remove empty rename
	test("customer rename id to num, name to nom project city",
		"customer project /*NOT UNIQUE*/ city")
	// move project before extend
	test("customer extend a = 5, b = 6 project id, a, name",
		"customer project id, name extend a = 5")
	// ... but not if extend uses fields not in project
	test("customer extend a = city, b = 6 project id, a, name",
		"customer extend a = city project id, a, name")
	// move project before extend & remove empty extend
	test("customer extend a = 5, b = 6 project id, name",
		"customer project id, name")
	test("table extend x = 123 project x",
		"project-none extend x = 123")
	test("table extend x = 123, y = Random() project x, y",
		"table extend y = Random() project /*NOT UNIQUE*/ y extend x = 123") // split
	test("table extend x = 123 project a, x",
		"table project a extend x = 123")
	test("table extend x = Random() project a, x",
		"table extend x = Random() project a, x")
	// remove unused constant extends
	test("table extend x = 123 project a, b",
		"table project a, b")
	test("table project a, b extend x = 123 project a, x",
		"table project a extend x = 123")
	// not unique so no deps
	test("withdeps rename b to bb, c to cc project bb, cc",
		"withdeps project /*NOT UNIQUE*/ b, c rename b to bb, c to cc")
	test("withdeps rename b to bb project a, bb",
		"withdeps project a, b, b_deps rename b to bb, b_deps to bb_deps")
	test("table rename a to aa extend x = 1 project aa, x",
		"table project a rename a to aa extend x = 1")

	// move where before project
	test("trans project id,cost where id is 5",
		"trans where id is 5 project /*NOT UNIQUE*/ id, cost")
	test("table project a,b where a is 5 project /*NOT UNIQUE*/ b",
		"table where a is 5 project b")
	// move where before rename
	test("trans where cost is 200 rename cost to x where id is 5",
		"trans where cost is 200 and id is 5 rename cost to x")
	// move where before extend
	test("trans where cost is 200 extend x = 1 where id is 5",
		"trans where cost is 200 and id is 5 extend x = 1")
	// move where before summarize
	test("hist where cost summarize id, total cost where id is 3",
		"hist where cost and id is 3 summarize id, total cost")
	test("hist where cost summarize id, total cost where id is 3 and total_cost > 10",
		"hist where cost and id is 3 summarize id, total cost "+
			"where total_cost > 10")

	// distribute where over intersect
	test("(hist intersect trans) where cost > 10",
		"hist where cost > 10 intersect (trans where cost > 10)")
	// distribute where over minus
	test("(hist minus trans) where cost > 10",
		"hist where cost > 10 minus (trans where cost > 10)")
	// distribute where over union
	test("(hist union trans) where cost > 10",
		"hist where cost > 10 union /*NOT DISJOINT*/ (trans where cost > 10)")
	// distribute where over times
	test("(customer times inven) where qty > 10 and city isnt 'toon'",
		"customer where city isnt 'toon' times (inven where qty > 10)")
	// distribute where over leftjoin
	test("(customer leftjoin trans) where id > 5",
		"customer where id > 5 leftjoin 1:n by(id) trans")
	// distribute where over leftjoin
	test("(customer leftjoin trans) where id > 5 and item =~ 'x'",
		"(customer where id > 5 leftjoin 1:n by(id) trans) where item =~ 'x'")
	// distribute where over join
	test("(customer join trans) where cost > 10 and city isnt 'toon'",
		"customer where city isnt 'toon' join 1:n by(id) (trans where cost > 10)")

	// convert leftjoin to join
	test("(tables leftjoin columns) where column isnt ''",
		"tables join 1:n by(table) (columns where column isnt '')")
	test("(tables leftjoin columns) where column is 123",
		// 1:1 because of `column is 123`
		"tables join 1:1 by(table) (columns where column is 123)")
	// same due to folding
	test("(tables leftjoin columns) where column in (123)",
		"tables join 1:1 by(table) (columns where column is 123)")
	test("(tables leftjoin columns) where table isnt ''",
		"tables where table isnt '' leftjoin 1:n by(table) columns")

	// distribute project over union
	test("(hist union trans) project item, cost",
		"hist project /*NOT UNIQUE*/ item, cost union /*NOT DISJOINT*/ "+
			"(trans project /*NOT UNIQUE*/ item, cost)")
	// split project over product
	test("(customer times inven) project city, item, id",
		"customer project city, id times (inven project item)")
	// split project over join
	test("(trans join customer) project city, item, id",
		"trans project /*NOT UNIQUE*/ item, id join n:1 by(id) (customer project city, id)")
	// ... but only if project includes join fields
	test("(trans join by(id) customer) project city, item",
		"(trans join n:1 by(id) customer) project /*NOT UNIQUE*/ city, item")
	// combine ... summarize ... project ...
	test("table summarize a, total b project a",
		"table project a")
	test("table summarize a, min b, max b project a, min_b",
		"table summarize a, min b")
	test("table summarize a, total b project total_b",
		"table summarize a, total b project /*NOT UNIQUE*/ total_b")
	// combine ... project ... summarize ...
	test("table project a, b summarize a, total b", // project-copy
		"table summarize a, total b")
	test("table project b, c summarize b, total c",
		"table project /*NOT UNIQUE*/ b, c summarize b, total c")

	test("trans where id is 1 join customer",
		"trans where id is 1 join n:1 by(id) (customer where id is 1)")
	test("trans join (customer where id is 1)",
		"trans where id is 1 join n:1 by(id) (customer where id is 1)")
	test("trans where id is 1 join (customer where id is 1)",
		"trans where id is 1 join n:1 by(id) (customer where id is 1)")
	test("trans where id is 1 join (customer where id is 2)",
		"nothing")
	test("(abc where b is 1) join (bcd where c is 2)",
		"abc where b is 1 and c is 2 join 1:1 by(b,c) (bcd where c is 2 and b is 1)")

	test("trans where id is 1 leftjoin customer",
		"trans where id is 1 leftjoin n:1 by(id) (customer where id is 1)")
	test("trans leftjoin (customer where id is 1)",
		"trans leftjoin n:1 by(id) (customer where id is 1)")
	test("trans where id is 1 leftjoin (customer where id is 1)",
		"trans where id is 1 leftjoin n:1 by(id) (customer where id is 1)")
	test("trans where id is 1 leftjoin (customer where id is 2)",
		"trans where id is 1 extend name = '', city = ''")
	test("(abc where b is 1) leftjoin (bcd where c is 2)",
		"abc where b is 1 leftjoin 1:1 by(b,c) (bcd where c is 2 and b is 1)")
}
