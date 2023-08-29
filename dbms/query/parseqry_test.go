// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestParseQuery(t *testing.T) {
	test := func(args ...string) {
		t.Helper()
		query := args[0]
		expected := args[0]
		if len(args) > 1 {
			expected = args[1]
		}
		q := ParseQuery(query, testTran{}, nil)
		qs := str.ToLower(q.String())
		assert.T(t).This(qs).Is(expected)
	}
	test("table")
	test("table sort a")
	test("table sort reverse a, b")
	test("table project a")
	test("table project a,b,c")
	test("table rename a to aa")
	test("table rename a to aa, c to cc")
	test("table intersect table2")
	test("table minus table2")
	test("table times cus")
	test("table union table2")
	test("cus join task",
		"cus join 1:n by(cnum) task")
	test("cus join by(cnum) task",
		"cus join 1:n by(cnum) task")
	test("cus leftjoin task",
		"cus leftjoin 1:n by(cnum) task")
	test("cus leftjoin by(cnum) task",
		"cus leftjoin 1:n by(cnum) task")
	test("table summarize count",
		"table summarize count")
	test("table summarize n = count")
	test("table summarize total a",
		"table summarize total a")
	test("table summarize t1 = total a")
	test("table summarize count, total a, max b",
		"table summarize count, total a, max b")
	test("table summarize a, b, count",
		"table summarize a, b, count")

	test("(table union table2) join table2",
		"(table union table2) join n:1 by(c,d,e) table2")
	test("cus join task sort tnum",
		"cus join 1:n by(cnum) task sort tnum")
	test("(cus join task) project cnum, abbrev, tnum rename cnum to c sort tnum, c",
		"(cus join 1:n by(cnum) task) project cnum,abbrev,tnum"+
			" rename cnum to c sort tnum, c")
	test("cus extend x = function(){123}",
		"cus extend x = /* function */")
	test("cus extend x = function(){123}()",
		"cus extend x = /* function */()")
	test("cus extend x = cnum.Map(function(){123})",
		"cus extend x = cnum.map(/* function */)") // (test does lower)

	xtest := func(s, err string) {
		fn := func() { ParseQuery(s, testTran{}, nil) }
		assert.T(t).This(fn).Panics(err)
	}
	xtest("table project", "expecting identifier")
	xtest("table remove", "expecting identifier")
	xtest("table rename", "expecting identifier")
	xtest("cus join by() task", "invalid empty join by")
	xtest("table summarize a, b", "expecting Comma")
	xtest("table summarize total", "expecting identifier")

	xtest("cus extend x = y = 1",
		"assignment operators are not allowed")
	xtest("cus extend x = cnum.Map({ it.Size() })",
		"queries do not support blocks")

	xtest("table rename a to x, a to y",
		"rename: nonexistent column: a")
	xtest("table rename a to x, b to x",
		"rename: column already exists: x")
}

func TestParseQuery2(t *testing.T) {
	test := func(s string) {
		t.Helper()
		q := ParseQuery(s, testTran{}, nil)
		assert.T(t).This(str.ToLower(q.String())).Is(s)
	}

	test("table extend one")
	test("table extend one, two = a + b")

	test("table where a > 1")
	test("table where a and b and c")
	test("table where a in (1, 2, 3)")

	s := "table where (((a > 1)))"
	q := ParseQuery(s, testTran{}, nil)
	assert.T(t).This(q.String()).Is("table WHERE a > 1")
}

func TestParseQueryView(t *testing.T) {
	q := ParseQuery("table union myview", testTran{}, nil)
	assert.T(t).This(q.String()).Is("table UNION (cus JOIN 1:n by(cnum) task)")
}
