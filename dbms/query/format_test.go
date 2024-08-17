// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	. "github.com/apmckinlay/gsuneido/core"
)

func TestFormat(t *testing.T) {
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	tran := &testTran{}
	test := func(query, expected string) {
		actual := Format(tran, query)
		assert.T(t).This(actual).Is(strings.Trim(expected, "\r\n"))
	}
	test("tables", `
tables`)
	test("myview", `
myview /*view*/`)
	test("tables extend b = 123, r", `
tables
extend b = 123, r`)
	test("tables project table, nrows", `
tables
project table, nrows`)
	test("tables project nrows", `
tables
project /*NOT UNIQUE*/ nrows`)
	test("tables remove nrows, totalsize", `
tables
remove nrows, totalsize`)
	test("columns rename column to col", `
columns
rename column to col`)
	test("tables sort table", `
tables
sort table`)
	test("columns summarize table, count", `
columns
summarize table, count`)
	test("tables where nrows > 5", `
tables
where nrows > 5`)
	test("(tables extend x = 1) union columns", `
    (tables
    extend x = 1)
union
    columns`)
	test("tables union (columns extend x = 1)", `
    tables
union
    (columns
    extend x = 1)`)
	test("tables union columns", `
    tables
union /*NOT DISJOINT*/
    columns`)
	test("tables intersect tables", `
    tables
intersect
    tables`)
	test("tables minus tables", `
    tables
minus
    tables`)
	test("tables times history", `
    tables
times
    history`)
	test("tables join columns", `
    tables
join by(table)
    columns`)
	test("tables leftjoin columns", `
    tables
leftjoin by(table)
    columns`)
}
