// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestTransform(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(from, expected string) {
		t.Helper()
		if expected == "" {
			expected = from
		}
		q := ParseQuery(from, testTran{}, nil)
		q = q.Transform()
		actual := str.ToLower(q.String())
		// *1 depends on whether optInit runs, e.g. if Nrows is called
		actual = strings.ReplaceAll(actual, "where*1", "where")
		assert.T(t).This(actual).Is(str.ToLower(expected))
	}

	test("table", "")

	// TablesLookup
	test("tables where table is 'foo'", "tables(foo)")

	test("table rename a to x, c to y", "")
	test("table remove c, d, e",
		"table project a,b")
	test("table remove x, y, z",
		"table")
	test("table project a, b, c",
		"table")
	test("withdeps remove b",
		"withdeps project a,c,c_deps")
	test("withdeps remove b_deps, c_deps",
		"withdeps project a,b,c")
	test("withdeps rename b to bb, c to cc project a, bb",
		"withdeps project a,b,b_deps rename b to bb, b_deps to bb_deps")

	// combine extend's
	test("customer extend a = 5 extend b = 6",
		"customer EXTEND a = 5, b = 6")
	test("customer extend a = 5 extend b = 6 where id > 5",
		"customer WHERE id > 5 EXTEND a = 5, b = 6")
	// combine project's
	test("table project a, b project b",
		"table project b")
	test("customer project id, name project id",
		"customer PROJECT id")
	test("customer project id, name where id > 5 project id",
		"customer WHERE id > 5 PROJECT id")
	test("customer project id, name project id where id > 5",
		"customer WHERE id > 5 PROJECT id")
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
		"customer WHERE id is 5 and city is 6 and name is 7")
	// leftjoin to join
	test("(cus leftjoin task) where cnum is 1 where tnum is 2",
		"cus where cnum is 1 join 1:1 by(cnum) (task where tnum is 2 and cnum is 1)")
	test("(cus leftjoin task) where tnum is 2 where cnum is 1",
		"cus where cnum is 1 join 1:1 by(cnum) (task where tnum is 2 and cnum is 1)")

	// remove projects of all fields
	test("customer project id, city, name", "customer")
	// remove disjoint difference
	test("(customer where id is 3) minus (customer where id is 5)",
		"customer WHERE id is 3")
	// remove empty extends
	test("customer extend zone = 3 project id, city",
		"customer PROJECT id,city")
	// remove empty renames
	test("customer rename name to nom project id, city",
		"customer PROJECT id,city")

	// move project before rename
	test("customer rename id to num, name to nom project num, city",
		"customer PROJECT id,city RENAME id to num")
	test("customer rename id to id2, id2 to id3 remove id3",
		"customer PROJECT name,city")
	// move project before rename & remove empty rename
	test("customer rename id to num, name to nom project city",
		"customer PROJECT city")
	// move project before extend
	test("customer extend a = 5, b = 6 project id, a, name",
		"customer PROJECT id,name EXTEND a = 5")
	// ... but not if extend uses fields not in project
	test("customer extend a = city, b = 6 project id, a, name",
		"customer EXTEND a = city PROJECT id,a,name")
	// move project before extend & remove empty extend
	test("customer extend a = 5, b = 6 project id, name",
		"customer PROJECT id,name")
	test("table extend x = 123 project x",
		"project-none extend x = 123")
	test("table extend x = 123, y = Random() project x, y",
		"table extend y = random() project y extend x = 123") // split
	test("table extend x = 123 project a, x",
		"table project a extend x = 123")
	test("table extend x = Random() project a, x",
		"table extend x = random() project a,x")
	// remove unused constant extends
	test("table extend x = 123 project a, b",
		"table project a,b")
	test("table project a,b extend x = 123 project a, x",
		"table project a extend x = 123")
	test("withdeps rename b to bb, c to cc project bb, cc",
		"withdeps project b,c rename b to bb, c to cc") // not unique so no deps
	test("withdeps rename b to bb project a, bb",
		"withdeps project a,b,b_deps rename b to bb, b_deps to bb_deps")
	test("table rename a to aa extend x = 1 project aa, x",
		"table project a rename a to aa extend x = 1")

	// move where before project
	test("trans project id,cost where id is 5",
		"trans WHERE id is 5 PROJECT id,cost")
	test("table project a,b where a is 5 project b",
		"table where a is 5 project b")
	// move where before rename
	test("trans where cost is 200 rename cost to x where id is 5",
		"trans WHERE cost is 200 and id is 5 RENAME cost to x")
	// move where before extend
	test("trans where cost is 200 extend x = 1 where id is 5",
		"trans WHERE cost is 200 and id is 5 EXTEND x = 1")
	// move where before summarize
	test("hist where cost summarize id, total cost where id is 3",
		"hist WHERE cost and id is 3 SUMMARIZE id, total cost")
	test("hist where cost summarize id, total cost where id is 3 and total_cost > 10",
		"hist WHERE cost and id is 3 SUMMARIZE id, total cost "+
			"WHERE total_cost > 10")

	// distribute where over intersect
	test("(hist intersect trans) where cost > 10",
		"hist WHERE cost > 10 INTERSECT (trans WHERE cost > 10)")
	// distribute where over minus
	test("(hist minus trans) where cost > 10",
		"hist WHERE cost > 10 MINUS (trans WHERE cost > 10)")
	// distribute where over union
	test("(hist union trans) where cost > 10",
		"hist WHERE cost > 10 UNION (trans WHERE cost > 10)")
	// distribute where over times
	test("(customer times inven) where qty > 10 and city isnt 'toon'",
		"customer WHERE city isnt 'toon' TIMES (inven WHERE qty > 10)")
	// distribute where over leftjoin
	test("(customer leftjoin trans) where id > 5",
		"customer WHERE id > 5 LEFTJOIN 1:n by(id) trans")
	// distribute where over leftjoin
	test("(customer leftjoin trans) where id > 5 and item =~ 'x'",
		"(customer WHERE id > 5 LEFTJOIN 1:n by(id) trans) WHERE item =~ 'x'")
	// distribute where over join
	test("(customer join trans) where cost > 10 and city isnt 'toon'",
		"customer WHERE city isnt 'toon' JOIN 1:n by(id) (trans WHERE cost > 10)")

	// convert LEFTJOIN to JOIN
	test("(tables leftjoin columns) where column isnt ''",
		"tables JOIN 1:n by(table) (columns WHERE column isnt '')")
	test("(tables leftjoin columns) where column is 123",
		// 1:1 because of `column is 123`
		"tables JOIN 1:1 by(table) (columns WHERE column is 123)")
	// same due to folding
	test("(tables leftjoin columns) where column in (123)",
		"tables JOIN 1:1 by(table) (columns WHERE column is 123)")
	test("(tables leftjoin columns) where table isnt ''",
		"tables WHERE table isnt '' LEFTJOIN 1:n by(table) columns")

	// distribute project over union
	test("(hist union trans) project item, cost",
		"hist PROJECT item,cost UNION (trans PROJECT item,cost)")
	// split project over product
	test("(customer times inven) project city, item, id",
		"customer PROJECT city,id TIMES (inven PROJECT item)")
	// split project over join
	test("(trans join customer) project city, item, id",
		"trans PROJECT item,id JOIN n:1 by(id) (customer PROJECT city,id)")
	// ... but only if project includes join fields
	test("(trans join by(id) customer) project city, item",
		"(trans JOIN n:1 by(id) customer) PROJECT city,item")
	// combine ... summarize ... project ...
	test("table summarize a, total b project a",
		"table project a")
	test("table summarize a, min b, max b project a, min_b",
		"table summarize a, min b")
	test("table summarize a, total b project total_b",
		"table summarize a, total b project total_b")
	// combine ... project ... summarize ...
	test("table project a, b summarize a, total b", // project-copy
		"table summarize a, total b")
	test("table project b, c summarize b, total c",
		"table project b,c summarize b, total c")
}
