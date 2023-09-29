// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

func ExampleSuInt() {
	v := SuInt(123)
	fmt.Printf("%d %s\n", v.toInt(), v.String())
	// Output: 123 123
}

func TestStrConvert(t *testing.T) {
	assert.T(t).This(AsStr(SuStr("123"))).Is("123")
}

func TestStringGet(t *testing.T) {
	var v Value = SuStr("hello")
	v = v.Get(nil, SuInt(1))
	assert.T(t).This(v).Is(Value(SuStr("e")))
}

func TestPanics(t *testing.T) {
	v := SuInt(123)
	assert.T(t).This(func() { v.Get(nil, v) }).Panics("number does not support get")
	ob := &SuObject{}
	assert.T(t).This(func() { ToInt(ob) }).Panics("can't convert object to integer")
}

func TestCompare(t *testing.T) {
	assert := assert.T(t)
	vals := []Value{False, True,
		SuDnum{Dnum: dnum.NegInf},
		SuInt(-1), SuInt(0), SuInt(+1), SuDnum{Dnum: dnum.PosInf},
		SuStr(""), SuStr("abc"),
		&SuExcept{SuStr: "bar"},
		NewSuConcat().Add("foo"),
		SuStr("world"),
		&SuExcept{SuStr: "zoo"}}
	for i := 1; i < len(vals); i++ {
		assert.This(vals[i].Compare(vals[i])).Is(0)
		assert.That(vals[i-1].Compare(vals[i]) < 0)
		assert.That(vals[i].Compare(vals[i-1]) > 0)
	}
}

func TestIfStr(t *testing.T) {
	xtest := func(v Value) {
		t.Helper()
		_, ok := v.ToStr()
		assert.T(t).False(ok)
	}
	xtest(True)
	xtest(False)
	xtest(Zero)   // SuInt
	xtest(MaxInt) // SuDnum
	xtest(&SuObject{})

	test := func(s string) {
		t.Helper()
		assert.T(t).This(ToStr(SuStr(s))).Is(s)
	}
	test("")
	test("hello")
}

func TestToStr(t *testing.T) {
	test := func(v Value, expected string) {
		t.Helper()
		assert.T(t).This(AsStr(v)).Is(expected)
	}
	test(EmptyStr, "")
	test(SuStr("hello"), "hello")
	test(True, "true")
	test(False, "false")
	test(Zero, "0")
	test(MaxInt, "2147483647")

	xtest := func(v Value) {
		t.Helper()
		_, ok := v.AsStr()
		assert.T(t).False(ok)
	}
	xtest(&SuObject{})
}

func TestIfInt(t *testing.T) {
	assert := assert.T(t)
	xtest := func(v Value) {
		t.Helper()
		_, ok := v.IfInt()
		assert.False(ok)
	}
	xtest(True)
	xtest(False)
	xtest(EmptyStr)
	xtest(&SuObject{})

	test := func(v Value, expected int) {
		t.Helper()
		got, ok := v.IfInt()
		assert.True(ok)
		assert.This(got).Is(expected)
	}
	test(Zero, 0)            // SuInt
	test(MaxInt, 2147483647) // SuDnum
}

func TestToInt(t *testing.T) {
	assert := assert.T(t)
	xtest := func(v Value) {
		t.Helper()
		_, ok := v.ToInt()
		assert.False(ok)
	}
	xtest(True)
	xtest(SuStr("hello"))
	xtest(&SuObject{})

	test := func(v Value, expected int) {
		t.Helper()
		got, ok := v.ToInt()
		assert.True(ok)
		assert.This(got).Is(expected)
	}
	test(Zero, 0)            // SuInt
	test(MaxInt, 2147483647) // SuDnum
	test(False, 0)
	test(EmptyStr, 0)
}

func TestIntVal(t *testing.T) {
	test := func(n int, expected string) {
		t.Helper()
		v := IntVal(n)
		typ := fmt.Sprintf("%T", v)
		assert.T(t).This(str.AfterFirst(typ, ".")).Is(expected)
		assert.T(t).This(ToInt(v)).Is(n)
	}
	test(0, "smi")
	test(123, "smi")
	test(-123, "smi")
	test(math.MaxInt16, "smi")
	test(math.MinInt16, "smi")
	test(math.MaxInt16+1, "SuDnum")
	test(math.MinInt16-1, "SuDnum")
	test(math.MaxInt32, "SuDnum")
	test(math.MinInt32, "SuDnum")
}
