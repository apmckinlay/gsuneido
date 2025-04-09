// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestPack(t *testing.T) {
	cv := NewSuConcat().Add("foo").Add("bar")
	values := []Packable{SuBool(false), SuBool(true), SuStr(""), SuStr("foo"), cv,
		SuInt(0), SuInt(1), SuInt(-1), dv("123.456"), dv(".1"), dv("-1e22"),
		dv("1234"), dv("12345678"), dv("123456789012"), dv("1234567890123456")}
	for _, v := range values {
		assert.T(t).This(Unpack(Pack(v))).Is(v)
	}
}

func TestPackSuInt(t *testing.T) {
	test := func(n int, expected ...byte) {
		t.Helper()
		v := IntVal(n)
		s := Pack(v)
		assert.T(t).This([]byte(s)).Is(expected)
		num := UnpackNumber(s)
		x, ok := SuIntToInt(num)
		assert.T(t).That(ok)
		assert.T(t).This(x).Is(n)
	}
	test(0, PackPlus)
	test(1, PackPlus, 129, 10)
	test(100, PackPlus, 0x83, 0x0a)
	test(10000, PackPlus, 133, 10)
	test(10002, PackPlus, 133, 10, 0, 20)
	test(-1, PackMinus, 126, 10^0xff)
	test(-10002, PackMinus, 122, 10^0xff, 0^0xff, 20^0xff)
}

func TestPackNum(t *testing.T) {
	test := func(s string, b ...byte) {
		t.Helper()
		p := Pack(dv(s))
		assert.T(t).This([]byte(p)).Is(b)
		assert.T(t).This(UnpackNumber(p).String()).Is(s)
	}
	test("0", PackPlus)
	test("1", PackPlus, 129, 10)
	test("-1", PackMinus, 126, 10^0xff)
	test(".1", PackPlus, 128, 10)
	test("20000", PackPlus, 133, 20)
	test("123.456", PackPlus, 131, 12, 34, 56)
	test("12345678.87654321", PackPlus, 136, 12, 34, 56, 78, 87, 65, 43, 21)
	test("1e23", PackPlus, 152, 10)
	test("-1e23", PackMinus, 152^0xff, 10^0xff)
	test("1e-23", PackPlus, 106, 10)
	test("inf", PackPlus, 0xff, 0xff)
	test("-inf", PackMinus, 0, 0)
}

func dv(s string) SuDnum {
	return SuDnum{Dnum: dnum.FromStr(s)}
}

func TestPackedToLower(t *testing.T) {
	same := func(v Value) {
		packed := Pack(v.(Packable))
		assert.T(t).Msg(v).This(packed).Is(PackedToLower(packed))
	}
	same(EmptyStr)
	same(True)
	same(False)
	same(Zero)
	same(IntVal(12345678))

	s := "Hello World!"
	ls := str.ToLower(s)
	assert.T(t).Msg(s).This(PackedToLower(Pack(SuStr(s)))).Is(Pack(SuStr(ls)))
}

func TestPackedCmpLower(t *testing.T) {
	values := []Value{EmptyStr, False, True, IntVal(-123), Zero, IntVal(12345678),
		SuStr("ant"), SuStr("Bug"), SuStr("cow")}
	packed := make([]string, len(values))
	for i, v := range values {
		packed[i] = Pack(v.(Packable))
	}
	for i, p1 := range packed {
		for j := i + 1; j < len(packed); j++ {
			p2 := packed[j]
			assert.T(t).Msg(values[i], "<=>", values[j]).
				This(PackedCmpLower(p1, p2)).Is(-1)
			assert.T(t).Msg(values[j], "<=>", values[i]).
				This(PackedCmpLower(p2, p1)).Is(+1)
		}
	}
	p1 := Pack(SuStr("hello world"))
	p2 := Pack(SuStr("Hello World"))
	assert.T(t).This(PackedCmpLower(p1, p2)).Is(0)
}

func BenchmarkPack(b *testing.B) {
	for range b.N {
		bench = Pack(emptyStr)
	}
}

var bench string

func TestPackV2(t *testing.T) {
	var v2 bool
	test := func(x Value) int {
		t.Helper()
		// fmt.Println("\ntest", x)
		// n := PackSize(x)
		// fmt.Println("PackSize:", n)
		s := packv(x.(Packable), v2)
		// fmt.Println("size:", len(s))
		// fmt.Printf("packed: %d %x\n", len(s), s)
		y := Unpack(s)
		assert.T(t).This(y).Is(x)
		return len(s)
	}
	for _, v2 = range []bool{false, true} {
		test(True)
		test(False)
		test(Zero)
		test(One)
		test(IntVal(12345678))
		test(SuStr("hello world"))
		test(&SuObject{})
		test(NewSuRecord())

		ob := &SuObject{}
		ob.Add(False)
		test(ob)
		ob.Add(True)
		test(ob)
		ob.Set(SuStr("hello"), SuInt(0x11))
		test(ob)
		ob.Set(SuStr("world"), SuInt(0x22))
		test(ob)

		x := &SuObject{}
		x.Set(SuStr("val"), SuInt(0x11))
		y := &SuObject{}
		y.Set(SuStr("val"), SuInt(0x22))
		data := &SuObject{}
		// data.Set(SuStr("x"), x)
		// data.Set(SuStr("y"), y)
		data.Add(x)
		data.Add(y)
		test(data)

		small := &SuObject{}
		small.Set(SuStr("abracadabra"), SuInt(123))
		large := &SuObject{}
		for range 6000 {
			large.Add(small)
		}
		n := test(large)
		assert.That(n > 64*1024)
		
		outer := &SuObject{}
		outer.Add(large)
		test(outer)
	}
}
