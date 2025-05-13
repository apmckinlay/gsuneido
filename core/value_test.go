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
	assert.T(t).This(v.Get(nil, v)).Is(nil)
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
	test(math.MaxInt16+1, "SuInt64")
	test(math.MinInt16-1, "SuInt64")
	test(math.MaxInt32, "SuInt64")
	test(math.MinInt32, "SuInt64")
}

func TestStringEquals(t *testing.T) {
	sustr := SuStr("hello world")
	suexcept := &SuExcept{SuStr: "hello world"}
	suconcat := NewSuConcat().Add("hello").Add(" world")
	for _, x := range []Value{sustr, suexcept, suconcat} {
		for _, y := range []Value{sustr, suexcept, suconcat} {
			assert.T(t).Msg(fmt.Sprintf("%T .Equal %T", x, y)).
				That(x.Equal(y))
        }
	}
	
}
func TestNumFromString(t *testing.T) {
    test := func(s string, expectedInt bool, expectedVal string) {
        t.Helper()
        v := NumFromString(s)
        if expectedInt {
            _, ok := v.(*smi)
            if !ok {
                _, ok = v.(SuInt64)
            }
            assert.T(t).Msg(s).This(ok).Is(true)
        } else {
            _, ok := v.(SuDnum)
            assert.T(t).Msg(s).This(ok).Is(true)
        }
        assert.T(t).Msg(s).This(v.String()).Is(expectedVal)
    }

    // Integer values
    test("0", true, "0")
    test("123", true, "123")
    test("-123", true, "-123")
    test("2147483647", true, "2147483647")    // MaxInt32
    test("-2147483648", true, "-2147483648")  // MinInt32
    
    // Hex values
    test("0x1a", true, "26")
	test("0xffff", true, "65535")
	test("0xffffffff", true, "4294967295") // MaxUint32
	test("0xffffffffffffffff", true, "-1")
    test("-0xff", true, "-255")
    
    // Decimal values
    test("123.456", false, "123.456")
    test("-123.456", false, "-123.456")
    test("0.123", false, ".123")
    test("-0.123", false, "-.123")
    
    // Scientific notation
    test("1e3", false, "1000")
    test("1.23e2", false, "123")
    test("1.23e-2", false, ".0123")
    
    // Large numbers
    test("9223372036854775807", true, "9223372036854775807")  // MaxInt64
    test("9999999999999999999", false, "9.999999999999999e18") // Beyond int64
}