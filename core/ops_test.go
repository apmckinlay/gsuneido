// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

func TestDiv(t *testing.T) {
	q := OpDiv(SuInt16(999), SuInt16(3))
	xi, xok := SuIntToInt(q)
	assert.T(t).This(xok).Is(true)
	assert.T(t).This(xi).Is(333)
	q = OpDiv(SuInt16(1), SuInt16(3))
	_ = q.(SuDnum)
}

func TestBool(t *testing.T) {
	assert.T(t).That(SuBool(true) == True)
	assert.T(t).That(SuBool(false) == False)
}
func TestIndex(t *testing.T) {
	assert.T(t).This(ToIndex(SuInt16(123))).Is(123)
	assert.T(t).This(ToIndex(SuDnum{Dnum: dnum.FromInt(123)})).Is(123)
}

var abc = SuStr("abc")
var G Value

func BenchmarkCat(b *testing.B) {
	for b.Loop() {
		s := EmptyStr
		for range 10000 {
			s = OpCat(s, abc)
		}
		G = s
	}
}

func BenchmarkJoin(b *testing.B) {
	for b.Loop() {
		ob := &SuObject{}
		for range 10000 {
			ob.Add(abc)
		}
		G = join(ob, EmptyStr)
	}
}

func join(this Value, arg Value) Value {
	ob := ToContainer(this)
	separator := AsStr(arg)
	sb := strings.Builder{}
	sep := ""
	iter := ob.ArgsIter()
	for {
		k, v := iter()
		if k != nil || v == nil {
			break
		}
		sb.WriteString(sep)
		sep = separator
		sb.WriteString(ToStrOrString(v))
	}
	return SuStr(sb.String())
}

func TestOpCatN(t *testing.T) {
	test := func(values ...Value) {
		t.Helper()
		th := &Thread{}
		expected := ""
		for _, v := range values {
			th.Push(v)
			expected += AsStr(v)
		}
		result := OpCatN(th, len(values))
		assert.T(t).This(ToStr(result)).Is(expected)
	}
	test(SuStr("hello"), SuStr(" "), SuStr("world"))
	test(SuStr("a"), SuStr("b"), SuStr("c"), SuStr("d"))
	test(SuStr("x"), SuStr(""), SuStr("y"), SuStr(""), SuStr("z"))
	test(SuStr("count="), SuInt16(42), SuStr(" name="), SuStr("Fred"), SuStr(" flag="), True)
	test(SuStr("a"), SuStr("b"), SuStr("c"), SuStr("d"), SuStr("e"), SuStr("f"),
		SuStr("g"), SuStr("h"), SuStr("i"), SuStr("j"), SuStr("k"), SuStr("l"))
}

func BenchmarkOpCatN(b *testing.B) {
	cases := []struct {
		name   string
		values []Value
	}{
		{"strings", []Value{SuStr("hello"), SuStr("cruel"), SuStr("world"),
			SuStr("and"), SuStr("all"), SuStr("who"), SuStr("live")}},
		{"mixed", []Value{SuStr("count="), SuInt16(42), SuStr(" name="),
			SuStr("Fred"), SuStr(" flag="), True}},
	}
	run := func(b *testing.B, op func(*Thread, int) Value, values []Value) {
		b.ReportAllocs()
		for b.Loop() {
			th := &Thread{}
			for _, v := range values {
				th.Push(v)
			}
			G = op(th, len(values))
		}
	}
	for _, c := range cases {
		values := c.values
		b.Run(c.name+"/OpCatN", func(b *testing.B) {
			run(b, OpCatN, values)
		})
	}
}
