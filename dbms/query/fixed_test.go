// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestFixed(t *testing.T) {
	runtime.DefaultSingleQuotes = true
	defer func() { runtime.DefaultSingleQuotes = false }()
	test := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query, testTran{})
		assert.T(t).This(fixedStr(q.Fixed())).Is(expected)
	}
	test("table", "[]")

	test("table extend f=1", "[f=(1)]")
	test("table extend f=1, g='s'", "[f=(1), g=('s')]")
	test("table extend f=1 extend g=2", "[g=(2), f=(1)]")
	test("table extend f=1 where f is 1", "[f=(1)]")
	assert.T(t).
		This(func() { test("table extend f=1 where f is 2", "[f=(2)]") }).
		Panics("conflict")
	test("table extend f=1, g=2 where f is 1", "[f=(1), g=(2)]")
	assert.T(t).
		This(func() { test("table extend f=1, g=2 where f is 3", "[f=(3), g=(2)]") }).
		Panics("conflict")

	test("table where a is 1", "[a=(1)]")
	test("table where a is 1 and b is 's' and a is b", "[a=(1), b=('s')]")

	test("table union (table extend f=1)", "[f=(1,'')]")
	test("(table extend f=2) union (table extend f=1)", "[f=(2,1)]")

	test("(table extend f=1, g=2) project a,b", "[]")
	test("(table extend f=1, g=2) project a,g", "[g=(2)]")
	test("(table extend f=1, g=2) project a,f,g", "[f=(1), g=(2)]")

	test("(table extend f=1) join (table extend f=1, g=2)", "[f=(1), g=(2)]")

	test("table extend f=1, g=2 rename g to h", "[f=(1), h=(2)]")
}
