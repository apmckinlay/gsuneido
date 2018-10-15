package runtime

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestPack(t *testing.T) {
	cv := NewSuConcat().Add("foo").Add("bar")
	values := []Packable{SuBool(false), SuBool(true), SuStr(""), SuStr("foo"), cv,
		SuInt(0), SuInt(1), SuInt(-1), dv("123.456"), dv(".1"), dv("-1e22")}
	for _, v := range values {
		Assert(t).That(Unpack(Pack(v)), EqVal(v))
	}
}

func EqVal(expected interface{}) Tester {
	return func(actual interface{}) string {
		if actual.(Value).Equal(expected.(Value)) {
			return ""
		}
		return fmt.Sprintf("expected %v but got %v", expected, actual)
	}
}

func TestPackInt32(t *testing.T) {
	buf := make([]byte, 0, 4)
	for _, n := range []int32{0, -1, 1, 999999, -999999} {
		Assert(t).That(unpackInt32(packInt32(n, buf)), Equals(n))
	}
}

func TestPackNum(t *testing.T) {
	test := func(s string, b ...byte) {
		t.Helper()
		Assert(t).That(Pack(dv(s)), Equals(b))
		v := NumFromString(s)
		Assert(t).That(Pack(v.(Packable)), Equals(b))
	}
	test("0", 3)
	test("1", 3, 129, 0, 1)
	test("-1", 2, 126, 255, 254)
	test(".1", 3, 128, 3, 232)
	test("20000", 3, 130, 0, 2)
	test("123.456", 3, 129, 0, 123, 17, 208)
	test("1e23", 3, 134, 3, 232)
	test("-1e-23", 2, 132, 255, 245)
	test("-99e7", 2, 124, 255, 246, 220, 215)
}

func dv(s string) SuDnum {
	return SuDnum{dnum.FromStr(s)}
}

func TestPackInt64(t *testing.T) {
	test := func(n int64, expected ...byte) {
		buf := make([]byte, 0, PackSizeInt64(n))
		buf = PackInt64(n, buf)
		Assert(t).That(buf, Equals(expected))
		num := UnpackNumber(rbuf{buf})
		x := num.(*smi)
		Assert(t).That(int64(x.ToInt()), Equals(n))
	}
	test(0, packPlus)
	test(1, packPlus, 129, 0, 1)
	test(10000, packPlus, 130, 0, 1)
	test(10002, packPlus, 130, 0, 1, 0, 2)
	test(-1, packMinus, 126, 255, 254)
	test(-10002, packMinus, 125, 255, 254, 255, 253)
}
