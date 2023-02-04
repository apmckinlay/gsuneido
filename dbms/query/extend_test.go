// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExtendInit(t *testing.T) {
	test := func(query string) {
		t.Helper()
		ParseQuery(query, testTran{}, nil)
	}
	test("hist extend price = cost")
	test("columns extend a = 1, b = 2, c = a + b")

	xtest := func(query, expected string) {
		t.Helper()
		assert.T(t).Msg(query).
			This(func() { ParseQuery(query, testTran{}, nil) }).Panics(expected)
	}
	xtest("inven extend qty = 1",
		"extend: column(s) already exist")
	xtest("inven extend price = cost",
		"extend: invalid column(s) in expressions: cost")
	xtest("columns extend c = a + b, a = 1, b = 2",
		"extend: invalid column(s) in expressions: a, b")
}

func TestExtendSelect(t *testing.T) {
	db := testDb()
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	rt := db.NewReadTran()
	ex := []string{"ex"}

	q := ParseQuery("cus extend ex=1", rt, nil) // constant
	zero := []string{Pack(Zero.(Packable))}
	q.Select(ex, zero)
	assert.T(t).This(q.Get(nil, Next)).Is(nil)
	one := []string{Pack(One.(Packable))}
	q.Select(ex, one)
	assert.T(t).That(q.Get(nil, Next) != nil)

	// where singleton
	q = ParseQuery("cus where cnum=1 extend ex=cnum+1", rt, nil) // expression
	q, _, _ = Setup(q, ReadMode, rt)
	q.Select(ex, zero)
	assert.T(t).This(q.Get(nil, Next)).Is(nil)
	two := []string{Pack(IntVal(2))}
	q.Select(ex, two)
	assert.T(t).That(q.Get(nil, Next) != nil)

	assert.T(t).That(q.Lookup(nil, ex, two) != nil)
}
