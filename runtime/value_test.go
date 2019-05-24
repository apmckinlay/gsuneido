package runtime

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func ExampleSuInt() {
	v := SuInt(123)
	fmt.Printf("%d %s\n", v.toInt(), v.String())
	// Output: 123 123
}

func TestStrConvert(t *testing.T) {
	Assert(t).That(AsStr(SuStr("123")), Equals("123"))
}

func TestStringGet(t *testing.T) {
	var v Value = SuStr("hello")
	v = v.Get(nil, SuInt(1))
	Assert(t).That(v, Equals(Value(SuStr("e"))))
}

func TestPanics(t *testing.T) {
	v := SuInt(123)
	Assert(t).That(func() { v.Get(nil, v) }, Panics("number does not support get"))

	ob := &SuObject{}
	Assert(t).That(func() { ToInt(ob) }, Panics("can't convert object to integer"))
}

func TestCompare(t *testing.T) {
	vals := []Value{False, True, SuDnum{Dnum: dnum.NegInf},
		SuInt(-1), SuInt(0), SuInt(+1), SuDnum{Dnum: dnum.PosInf},
		SuStr(""), SuStr("abc"), NewSuConcat().Add("foo"), SuStr("world")}
	for i := 1; i < len(vals); i++ {
		Assert(t).That(vals[i].Compare(vals[i]), Equals(0))
		Assert(t).That(vals[i-1].Compare(vals[i]), Equals(-1).Comment(vals[i-1], vals[i]))
		Assert(t).That(vals[i].Compare(vals[i-1]), Equals(+1))
	}
}

func TestIfStr(t *testing.T) {
	xtest := func(v Value) {
		t.Helper()
		_, ok := v.ToStr()
		Assert(t).False(ok)
	}
	xtest(True)
	xtest(False)
	xtest(Zero)   // SuInt
	xtest(MaxInt) // SuDnum
	xtest(&SuObject{})

	test := func(s string) {
		t.Helper()
		Assert(t).That(ToStr(SuStr(s)), Equals(s))
	}
	test("")
	test("hello")
}

func TestToStr(t *testing.T) {
	test := func(v Value, expected string) {
		t.Helper()
		Assert(t).That(AsStr(v), Equals(expected))
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
		Assert(t).False(ok)
	}
	xtest(&SuObject{})
}

func TestIfInt(t *testing.T) {
	xtest := func(v Value) {
		t.Helper()
		_, ok := v.IfInt()
		Assert(t).False(ok)
	}
	xtest(True)
	xtest(False)
	xtest(EmptyStr)
	xtest(&SuObject{})

	test := func(v Value, expected int) {
		t.Helper()
		got, ok := v.IfInt()
		Assert(t).True(ok)
		Assert(t).That(got, Equals(expected))
	}
	test(Zero, 0)            // SuInt
	test(MaxInt, 2147483647) // SuDnum
}

func TestToInt(t *testing.T) {
	xtest := func(v Value) {
		t.Helper()
		_, ok := v.ToInt()
		Assert(t).False(ok)
	}
	xtest(True)
	xtest(SuStr("hello"))
	xtest(&SuObject{})

	test := func(v Value, expected int) {
		t.Helper()
		got, ok := v.ToInt()
		Assert(t).True(ok)
		Assert(t).That(got, Equals(expected))
	}
	test(Zero, 0)            // SuInt
	test(MaxInt, 2147483647) // SuDnum
	test(False, 0)
	test(EmptyStr, 0)
}
