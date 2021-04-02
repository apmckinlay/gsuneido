// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestParseAct(t *testing.T) {
	test := func(s string) {
		t.Helper()
		act := ParseAction(s)
		assert.T(t).This(str.ToLower(act.String())).Is(s)
	}
	test("insert [a: 1, b: 3] into table")
	test("insert table into table1")
	test("insert table where x is 1 into table1")

	test("update table set x = 1")
	test("update table set x = 1, y = 2")

	test("delete table")
	test("delete table where x > 1")

	assert.This(func() { ParseAction("foo bar") }).Panics("action must")
}
