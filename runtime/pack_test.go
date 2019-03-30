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
		s := Pack(SuInt(n))
		Assert(t).That(s, Equals(string(expected)))
		num := UnpackNumber(s)
		x, _ := SuIntToInt(num)
		Assert(t).That(x, Equals(n))
	}
	test(0, packPlus)
	test(1, packPlus, 129, 0, 1)
	test(10000, packPlus, 130, 0, 1)
	test(10002, packPlus, 130, 0, 1, 0, 2)
	test(-1, packMinus, 126, 255, 254)
	test(-10002, packMinus, 125, 255, 254, 255, 253)
}

func TestPackNum(t *testing.T) {
	test := func(s string, b ...byte) {
		t.Helper()
		p := Pack(dv(s))
		Assert(t).That(p, Equals(string(b)))
		Assert(t).That(UnpackNumber(p).String(), Equals(s))
	}
	test("0", 3)
	test("1", 3, 129, 0, 1)
	test("-1", 2, 126, 255, 254)
	test(".1", 3, 128, 3, 232)
	test("20000", 3, 130, 0, 2)
	test("123.456", 3, 129, 0, 123, 17, 208)
	test("1e23", 3, 134, 3, 232)
	test("-1e-23", 2, 132, 255, 245)
}

func dv(s string) SuDnum {
	return SuDnum{Dnum: dnum.FromStr(s)}
}
