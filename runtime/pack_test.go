package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestPack(t *testing.T) {
	cv := NewSuConcat().Add("foo").Add("bar")
	values := []Packable{SuBool(false), SuBool(true), SuStr(""), SuStr("foo"), cv,
		SuInt(0), SuInt(1), SuInt(-1), dv("123.456"), dv(".1"), dv("-1e22"),
		dv("1234"), dv("12345678"), dv("123456789012"), dv("1234567890123456")}
	for _, v := range values {
		Assert(t).That(Unpack(Pack(v)), Equals(v))
	}
}

func TestPackSuInt(t *testing.T) {
	test := func(n int, expected ...byte) {
		t.Helper()
		s := Pack(SuInt(n))
		Assert(t).That([]byte(s), Equals(expected))
		num := UnpackNumber(s)
		x, ok := SuIntToInt(num)
		Assert(t).True(ok)
		Assert(t).That(x, Equals(n))
	}
	test(0, packPlus)
	test(1, packPlus, 129, 10)
	test(10000, packPlus, 133, 10)
	test(10002, packPlus, 133, 10, 0, 20)
	test(-1, packMinus, 126, 10 ^ 0xff)
	test(-10002, packMinus, 122, 10 ^ 0xff, 0 ^ 0xff, 20 ^ 0xff)
}

func TestPackNum(t *testing.T) {
	test := func(s string, b ...byte) {
		t.Helper()
		p := Pack(dv(s))
		Assert(t).That([]byte(p), Equals(b))
		Assert(t).That(UnpackNumber(p).String(), Equals(s))
	}
	test("0", packPlus)
	test("1", packPlus, 129, 10)
	test("-1", packMinus, 126, 10 ^ 0xff)
	test(".1", packPlus, 128, 10)
	test("20000", packPlus, 133, 20)
	test("123.456", packPlus, 131, 12, 34, 56)
	test("12345678.87654321", packPlus, 136, 12, 34, 56, 78, 87, 65, 43, 21)
	test("1e23", packPlus, 152, 10)
	test("-1e23", packMinus, 152 ^ 0xff, 10 ^ 0xff)
	test("1e-23", packPlus, 106, 10)
}

func dv(s string) SuDnum {
	return SuDnum{Dnum: dnum.FromStr(s)}
}
