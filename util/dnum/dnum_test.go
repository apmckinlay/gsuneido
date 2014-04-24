package dnum

import (
	"fmt"
	"math"
	"strconv"
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func Test_String(t *testing.T) {
	assert := Assert(t)
	assert.That(Zero.String(), Equals("0"))
	assert.That(Inf.String(), Equals("inf"))
	assert.That(Dnum{123, 0, 0}.String(), Equals("123"))
	assert.That(Dnum{123000, 0, -3}.String(), Equals("123"))
}

func Test_Parse(t *testing.T) {
	assert := Assert(t)
	test := func(s string, expected Dnum) {
		dn := parse(s)
		assert.That(dn, Equals(expected))
	}
	test("0e4", Zero)
	test("-0", Zero)
}

// for testing - accepts "inf" and "-inf", panics on error
func parse(s string) Dnum {
	switch s {
	case "inf":
		return Inf
	case "-inf":
		return MinusInf
	case "-0":
		return Zero
	default:
		n, err := Parse(s)
		if err != nil {
			panic("parse failed")
		}
		return n
	}
}

func TestConvert(t *testing.T) {
	assert := Assert(t)
	test := func(s string, dn Dnum) {
		g := parse(s)
		assert.That(g, Equals(dn).Comment("from "+s))
		assert.That(dn.String(), Equals(s))
	}
	test("0", Zero)

	test("123", Dnum{123, 0, 0})
	test("-123", Dnum{123, 1, 0})

	test("10000", Dnum{10000, 0, 0})
	test("1e5", Dnum{1, 0, 5})

	test(".1234", Dnum{1234, 0, -4})
	test(".0001", Dnum{1, 0, -4})
	test("1e-5", Dnum{1, 0, -5})

	test("123.4", Dnum{1234, 0, -1})
	test("1.234", Dnum{1234, 0, -3})

	test("12345678912345678912", Dnum{12345678912345678912, 0, 0})
}

func Test_Neg(t *testing.T) {
	assert := Assert(t)
	Neg := func(x string, expected string) {
		xn := parse(x)
		zn := xn.Neg()
		assert.That(zn.String(), Equals(expected))
	}
	Neg("0", "0")
	Neg("123", "-123")
	Neg("-123", "123")
	Neg("inf", "-inf")
	Neg("-inf", "inf")
}

func Test_Cmp(t *testing.T) {
	assert := Assert(t)
	data := []string{"-inf", "-1e9", "-1e-9", "0", "1e-9", "1e9", "inf"}
	for i, xs := range data {
		x := parse(xs)
		assert.That(Cmp(x, x), Equals(0).Comment(fmt.Sprint(x, " >< ", x)))
		for _, ys := range data[i+1:] {
			y := parse(ys)
			assert.That(Cmp(x, y), Equals(-1).Comment(fmt.Sprint(x, " >< ", y)))
			assert.That(Cmp(y, x), Equals(1).Comment(fmt.Sprint(y, " >< ", x)))
		}
	}
}

func Test_Add(t *testing.T) {
	assert := Assert(t)
	add := func(x string, y string, expected string) {
		xn := parse(x)
		yn := parse(y)
		zn := Add(xn, yn)
		assert.That(zn.String(), Equals(expected))
		zn = Add(yn, xn)
		assert.That(zn.String(), Equals(expected))
	}
	// special cases (no actual math)
	add("123", "0", "123")
	add("inf", "-inf", "0")
	add("inf", "inf", "inf")
	add("-inf", "-inf", "-inf")
	add("inf", "123", "inf")
	add("-inf", "123", "-inf")
	// aligned
	add("123", "456", "579")
	add("-123", "-456", "-579")
	add("1.23e9", "4.56e9", "5.79e9")
	add("123", "-456", "-333")
	add("-123", "456", "333")
	// need aligning
	add("123", "1e-99", "123")
	add("1e12", "1e14", "1.01e14")
	add("11111111111111111111", "2222222222222222222e-4", "11111333333333333333")
	add("11111111111111111111", "6666666666666666666e-4", "11111777777777777778")
	// int64 overflow
	add("18446744073709551615", "11", "18446744073709551630")
}

func Test_Sub(t *testing.T) {
	assert := Assert(t)
	sub := func(x string, y string, expected string) {
		xn := parse(x)
		yn := parse(y)
		zn := Sub(xn, yn)
		assert.That(zn.String(), Equals(expected))
		if expected != "0" {
			zn = Sub(yn, xn)
			assert.That(zn.String(), Equals("-"+expected))
		}
	}
	// special cases (no actual math)
	sub("123", "0", "123")
	sub("inf", "-inf", "inf")
	sub("inf", "inf", "0")
	sub("-inf", "-inf", "0")
	sub("inf", "123", "inf")
	// aligned
	sub("456", "123", "333")
	sub("-123", "-456", "333")
	sub("4.56e9", "1.23e9", "3.33e9")
	sub("123", "-456", "579")
	sub("456", "-123", "579")
	sub("123", "-456", "579")
	// need aligning
	sub("123", "1e-99", "123")
	sub("1e14", "1e12", "9.9e13")
	sub("12222222222222222222", "11111111111111111111e-4", "12221111111111111111")
}

func Test_Mul(t *testing.T) {
	assert := Assert(t)
	mul := func(x string, y string, expected string) {
		mul2 := func(x string, y string, expected string) {
			xn := parse(x)
			yn := parse(y)
			zn := Mul(xn, yn)
			assert.That(zn.String(), Equals(expected).Comment(fmt.Sprint(xn, " * ", yn)))
		}
		mul2(x, y, expected)
		if expected != "0" {
			mul2("-"+x, y, "-"+expected)
			mul2("-"+x, y, "-"+expected)
		}
		mul2("-"+x, "-"+y, expected)
	}
	// special cases (no actual math)
	mul("0", "0", "0")
	mul("123", "0", "0")
	mul("123", "inf", "inf")
	mul("inf", "inf", "inf")
	// result fits in uint64
	mul("2", "333", "666")
	mul("2e9", "333e-9", "666")
	mul("2e3", "3e3", "6e6")
	mul("123456789000000000", "123456789000000000", "1.5241578750190521e34")
	mul("2e99", "2e99", "inf") // exp overflow
	// result too big for uint64
	mul("1234567890123456", "1234567890123456", "1.524157875323881728e30")
}

func Test_split(t *testing.T) {
	dn := Dnum{123456789987654321, 0, 0}
	lo, hi := dn.split()
	Assert(t).That(lo, Equals(uint64(987654321)))
	Assert(t).That(hi, Equals(uint64(123456789)))
}

func Test_Div(t *testing.T) {
	assert := Assert(t)
	div := func(x string, y string, expected string) {
		xn := parse(x)
		yn := parse(y)
		zn := Div(xn, yn)
		assert.That(zn.String(), Equals(expected))
	}
	// special cases (no actual math)
	div("0", "0", "0")
	div("123", "0", "inf")
	div("123", "inf", "0")
	div("inf", "inf", "inf")
	div("1e99", "1e-99", "inf") // exp overflow
	div("1e-99", "1e99", "0")   // exp underflow
	// divides evenly
	div("4444", "2222", "2")
	div("2222", "4444", ".5")
	// long division
	div("2", "3", ".6666666666666666667")
	div("1", "3", ".3333333333333333333")
	div("11", "17", ".6470588235294117647")
	div("1234567890123456", "9876543210123456", ".12499999887187493003")
}

func Test_float64_convert(t *testing.T) {
	assert := Assert(t)
	cvt := func(dn float64) {
		f10 := FromFloat64(dn)
		f2 := f10.Float64()
		assert.That(f2, Equals(dn))
	}
	cvt(0.0)
	cvt(123.0)
	cvt(1.0 / 3.0)
	cvt(123e3)
	cvt(-123e-44)
	cvt(math.Inf(1))
	cvt(math.Inf(-1))
}

func Test_toInt(t *testing.T) {
	assert := Assert(t)
	test := func(x string, expected uint64) {
		dn := parse(x)
		z, err := dn.toUint()
		if err != nil {
			z = math.MaxUint64
		}
		assert.That(z, Equals(expected))
	}
	test("123", 123)
	test("123e3", 123000)
	test("1.23e2", 123)
	test(".000123e6", 123)
	test("1e-99", 0)
	test("1e99", math.MaxUint64)
}

func Test_ToInt(t *testing.T) {
	assert := Assert(t)
	test := func(x string, expected string) {
		dn := parse(x)
		z, err := dn.Int64()
		if err != nil {
			assert.That(err.Error(), Equals(expected))
			return
		}
		nexpected, err := strconv.ParseInt(expected, 10, 64)
		if err != nil {
			panic("bad test data!")
		}
		assert.That(z, Equals(nexpected))
	}
	test("123", "123")
	test("-123", "-123")
	test("1e99", "outside range")
	test("9223372036854775807", "9223372036854775807") // max int64
	test("18446744073709551615", "outside range")      // max uint64
}

func Test_ToUint(t *testing.T) {
	assert := Assert(t)
	test := func(x string, expected string) {
		dn := parse(x)
		z, err := dn.Uint64()
		if err != nil {
			assert.That(err.Error(), Equals(expected))
			return
		}
		nexpected, err := strconv.ParseUint(expected, 10, 64)
		if err != nil {
			panic("bad test data!")
		}
		assert.That(z, Equals(nexpected))
	}
	test("123", "123")
	test("1e99", "outside range")
	test("-123", "outside range")
	test("9223372036854775807", "9223372036854775807")   // max int64
	test("18446744073709551615", "18446744073709551615") // max uint64
}

var bench Dnum

func BenchmarkAdd(b *testing.B) {
	x := parse("11111111111111111111")
	y := parse("2222222222222222222e-4")
	var z Dnum
	for n := 0; n < b.N; n++ {
		z = Add(x, y)
	}
	bench = z
}

func BenchmarkDiv(b *testing.B) {
	x := parse("11")
	y := parse("17")
	var z Dnum
	for n := 0; n < b.N; n++ {
		z = Div(x, y)
	}
	bench = z
}
