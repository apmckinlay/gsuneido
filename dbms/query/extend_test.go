// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
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
	assert := assert.T(t)
	db := testDb()
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	rt := db.NewReadTran()
	ex := []string{"ex"}

	q := ParseQuery("cus extend ex=1", rt, nil) // constant
	zero := []string{Pack(Zero.(Packable))}
	one := []string{Pack(One.(Packable))}

	assert.This(q.Lookup(nil, ex, zero)).Is(nil)
	q.Select(ex, zero)
	assert.This(q.Get(nil, Next)).Is(nil)
	q.Select(nil, nil)

	q.Select(ex, one)
	assert.That(q.Get(nil, Next) != nil)
	q.Select(nil, nil)

	// where singleton
	q = ParseQuery("cus where cnum=1 extend ex=cnum+1", rt, nil) // expression
	q, _, _ = Setup(q, ReadMode, rt)
	assert.That(q.fastSingle())

	q.Select(ex, zero)
	assert.This(q.Get(nil, Next)).Is(nil)
	q.Select(nil, nil)

	two := []string{Pack(IntVal(2))}
	q.Select(ex, two)
	assert.That(q.Get(nil, Next) != nil)
	q.Select(nil, nil)

	assert.That(q.Lookup(nil, ex, two) != nil)
}

func TestExtendRuleBug(t *testing.T) {
	db := heapDb()
	db.adm("create cus (ck, c4) key(ck)")
	db.act("insert { ck: 1, c4: 2  } into cus")
	db.adm("create ivc (ik, ck, i3) key(ik)")
	db.act("insert { ck: 1, ik: 3, i3: 2 } into ivc")
	assert.This(queryAll(db.Database,
		"(cus where ck is 1 extend r1, i3 = c4) join by(ck,i3) ivc")).
		Is("ik=3 ck=1 i3=2 c4=2")
}

func BenchmarkExtend(b *testing.B) {
	MakeSuTran = func(qt QueryTran) *SuTran { return nil }
	q := ParseQuery("table extend z=1, aa=a, y=2, bb=b, x=3, cc=c", testTran{}, nil)
	e := q.(*Extend)
	e.ctx.Hdr = e.getHeader()
	var rb RecordBuilder
	rb.Add(SuStr("eh"))
	rb.Add(SuStr("be"))
	rb.Add(SuStr("see"))
	row := []DbRec{{Record: rb.Build()}}
	for range b.N {
		e.extendRow(nil, row)
	}
}
