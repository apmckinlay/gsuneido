// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFormat(t *testing.T) {
	MakeSuTran = func(qt QueryTran) *core.SuTran { return nil }
	tran := &testTran{}
	test := func(query, expected string) {
		t.Helper()
		actual := Format(tran, query)
		assert.T(t).This(actual).Is(strings.Trim(expected, "\r\n"))
	}
	test("table", `
table`)
	test("myview", `
myview /*view*/`)
	test("table extend d = 123, e", `
table
extend d = 123, e`)
	test("table project a, b", `
table
project a, b`)
	test("table project b", `
table
project /*NOT UNIQUE*/ b`)
	test("table remove b, c", `
table
remove b, c`)
	test("table rename b to d", `
table
rename b to d`)
	test("table sort a", `
table
sort a`)
	test("table summarize a, count", `
table
summarize a, count`)
	test("table where b > 5", `
table
where b > 5`)
	test("(table extend d = 1) union customer", `
    (table
    extend d = 1)
union
    customer`)
	test("table union (customer extend x = 1)", `
    table
union
    (customer
    extend x = 1)`)
	test("table union customer", `
    table
union /*NOT DISJOINT*/
    customer`)
	test("table intersect table2", `
    table
intersect
    table2`)
	test("table minus table2", `
    table
minus
    table2`)
	test("table times task", `
    table
times
    task`)
	test("cus join task", `
    cus
join by(cnum)
    task`)
	test("cus leftjoin task", `
    cus
leftjoin by(cnum)
    task`)
	test("cus semijoin task", `
    cus
semijoin by(cnum)
    task`)
}
