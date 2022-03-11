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
		act := ParseAction(s, testTran{}, nil)
		assert.T(t).This(str.ToLower(act.String())).Is(s)
	}
	test("insert [a: 1, b: 3] into table")
	test("insert table into table1")
	test("insert table where a is 1 into table1")

	test("update table set a = 1")
	test("update table set a = 1, b = 2")

	test("delete table")
	test("delete table where a > 1")

	assert.This(func() {
		ParseAction("foo bar", testTran{}, nil)
	}).Panics("action must")
	assert.This(func() {
		ParseAction("update table set a = b = 2", testTran{}, nil)
	}).Panics("assignment operators are not allowed")
	assert.This(func() {
		ParseAction("update table set a = b *= 2", testTran{}, nil)
	}).Panics("assignment operators are not allowed")
	assert.This(func() {
		ParseAction("update table set a = ++b", testTran{}, nil)
	}).Panics("increment/decrement operators are not allowed")
}
