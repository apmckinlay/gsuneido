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
	q := OpDiv(SuInt(999), SuInt(3))
	xi, xok := SuIntToInt(q)
	assert.T(t).This(xok).Is(true)
	assert.T(t).This(xi).Is(333)
	q = OpDiv(SuInt(1), SuInt(3))
	_ = q.(SuDnum)
}

func TestBool(t *testing.T) {
	assert.T(t).That(SuBool(true) == True)
	assert.T(t).That(SuBool(false) == False)
}
func TestIndex(t *testing.T) {
	assert.T(t).This(ToIndex(SuInt(123))).Is(123)
	assert.T(t).This(ToIndex(SuDnum{Dnum: dnum.FromInt(123)})).Is(123)
}

var abc = SuStr("abc")
var G Value

func BenchmarkCat(b *testing.B) {
	for b.Loop() {
		s := EmptyStr
		for range 10000 {
			s = OpCat(nil, s, abc)
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
	test := func(values ...string) {
		t.Helper()
		th := &Thread{}
		expected := ""
		for _, v := range values {
			th.Push(SuStr(v))
			expected += v
		}
		result := OpCatN(th, len(values))
		assert.T(t).This(result).Is(SuStr(expected))
	}
	test("hello", " ", "world")
	test("a", "b", "c", "d")
	test("x", "", "y", "", "z")
	test("a", "b", "c", "d", "e", "f")
}
