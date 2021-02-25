// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExtendInit(t *testing.T) {
	test := func(query string) {
		t.Helper()
		q := ParseQuery(query)
		q.SetTran(testTran{})
		q.Init()
	}
	test("hist extend price = cost")
	test("columns extend a = 1, b = 2, c = a + b")

	xtest := func(query, expected string) {
		t.Helper()
		q := ParseQuery(query)
		q.SetTran(testTran{})
		assert.T(t).Msg(query).
			This(func() { q.Init() }).Panics(expected)
	}
	xtest("inven extend qty = 1",
		"extend: column(s) already exist")
	xtest("inven extend price = cost",
		"extend: invalid column(s) in expressions: cost")
	xtest("columns extend c = a + b, a = 1, b = 2",
		"extend: invalid column(s) in expressions: a, b")
}
