package interp

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

// TODO test comparison/ordering of packed values

func TestPack(t *testing.T) {
	cv := NewSuConcat().Add("foo").Add("bar")
	values := []Packable{False, True, SuStr(""), SuStr("foo"), cv,
		SuInt(0), SuInt(1), SuInt(-1), dv("123.456"), dv(".1"), dv("-1e22")}
	for _, v := range values {
		Assert(t).That(Unpack(Pack(v)), EqVal(v))
	}
}

func EqVal(expected interface{}) Tester {
	return func(actual interface{}) string {
		if actual.(Value).Equals(expected.(Value)) {
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
		Assert(t).That(Pack(dv(s)), Equals(b))
		v, err := ParseNum(s)
		if err == nil {
			Assert(t).That(Pack(v.(Packable)), Equals(b))
		}
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
	dn, _ := dnum.Parse(s)
	return SuDnum{dn}
}
