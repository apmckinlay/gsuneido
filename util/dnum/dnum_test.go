// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dnum

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func Test_size(t *testing.T) {
	// due to allignment and padding, size is 16 bytes instead of 10
	Assert(t).That(int(unsafe.Sizeof(Dnum{})), Equals(16))
	var a [10]Dnum
	Assert(t).That(int(unsafe.Sizeof(a)), Equals(160))
}

func Test_inf(t *testing.T) {
	Assert(t).That(Inf(0), Equals(Zero))
	Assert(t).That(Inf(+1), Equals(PosInf))
	Assert(t).That(Inf(-1), Equals(NegInf))
}

func Test_ilog10(t *testing.T) {
	Assert(t).That(ilog10(0), Equals(0))
	Assert(t).That(ilog10(123), Equals(2))
}

func Test_New(t *testing.T) {
	Assert(t).That(New(signZero, 0, 0), Equals(Zero))
	Assert(t).That(New(signPos, 1, 999), Equals(PosInf)) // exponent overflow
	Assert(t).That(New(signNeg, 1, 999), Equals(NegInf)) // exponent overflow
	Assert(t).That(New(signPos, 1, -999), Equals(Zero))  // exponent underflow
	Assert(t).That(New(signNeg, 1, -999), Equals(Zero))  // exponent underflow
	Assert(t).That(New(signPos, 1, 0), Equals(Dnum{1000000000000000, 1, -15}))
	Assert(t).That(New(signPos, 123, 0), Equals(Dnum{1230000000000000, 1, -13}))
}

func Test_String(t *testing.T) {
	assert := Assert(t)
	assert.That(Zero.String(), Equals("0"))
	assert.That(One.String(), Equals("1"))
	assert.That(PosInf.String(), Equals("inf"))
	assert.That(NegInf.String(), Equals("-inf"))
	assert.That(FromInt(123).String(), Equals("123"))
	assert.That(FromInt(-123).String(), Equals("-123"))

	assert.That(New(signPos, 1234000000000000, -20).String(), Equals("1.234e-21"))
	assert.That(New(signPos, 1234000000000000, -2).String(), Equals(".001234"))
	assert.That(New(signPos, 1234000000000000, 0).String(), Equals(".1234"))
	assert.That(New(signPos, 1234000000000000, 2).String(), Equals("12.34"))
	assert.That(New(signPos, 1234000000000000, 4).String(), Equals("1234"))
	assert.That(New(signPos, 1234000000000000, 6).String(), Equals("123400"))
	assert.That(New(signPos, 1234000000000000, 20).String(), Equals("1.234e19"))
}

func Test_FromStr(t *testing.T) {
	assert := Assert(t)
	assert.That(FromStr("inf"), Equals(PosInf))
	assert.That(FromStr("+inf"), Equals(PosInf))
	assert.That(FromStr("-inf"), Equals(NegInf))
	assert.That(FromStr("0"), Equals(Zero))
	assert.That(FromStr("+0"), Equals(Zero))
	assert.That(FromStr("-0"), Equals(Zero))
	assert.That(FromStr("0e4"), Equals(Zero))
	assert.That(FromStr("0000"), Equals(Zero))
	assert.That(FromStr("0000."), Equals(Zero))
	assert.That(FromStr(".0000"), Equals(Zero))
	assert.That(FromStr("0000.0000"), Equals(Zero))
	assert.That(FromStr("1"), Equals(One))
	assert.That(FromStr("000000000000000000001"), Equals(One))
	assert.That(FromStr("1.0000000000000000000"), Equals(One))
	assert.That(FromStr("100000000000000000000"), Equals(FromStr("1e20")))
	assert.That(FromStr(".1234567890123456789"), Equals(FromStr(".1234567890123456")))
	assert.That(FromStr(".000000000000000000001"), Equals(FromStr(".1e-20")))
}

func Test_FromToStr(t *testing.T) {
	test := func(s string) {
		Assert(t).That(FromStr(s).String(), Equals(s))
	}
	test("inf")
	test("-inf")
	test("0")
	test("1")
	test("-1")
	test("123")
	test("-123")
	test("100")
	test(".1")
	test(".00001")
	test("1e20")
	test("-1e-20")
	test("1e18")
}

func Test_getExp(t *testing.T) {
	e := getExp(&reader{"e20", 0})
	Assert(t).That(e, Equals(20))
}

func Test_FromToInt(t *testing.T) {
	assert := Assert(t)
	test := func(x int64) {
		n, ok := FromInt(x).ToInt64()
		assert.True(ok)
		assert.That(n, Equals(x))
		n, ok = FromInt(-x).ToInt64()
		assert.True(ok)
		assert.That(n, Equals(-x))
	}
	test(0)
	test(1)
	test(100)
	test(123)
	test(coefMax)
	test(1e15)
	test(1e16)
	test(1e17)
	test(1e18)
}

func Test_FromInt(t *testing.T) {
	Assert(t).That(FromInt(0), Equals(Zero))
	Assert(t).That(FromInt(1), Equals(Dnum{1000000000000000, +1, 1}))
	Assert(t).That(FromInt(100), Equals(Dnum{1000000000000000, +1, 3}))
	Assert(t).That(FromInt(123), Equals(Dnum{1230000000000000, +1, 3}))
	Assert(t).That(FromInt(-123), Equals(Dnum{1230000000000000, -1, 3}))
	Assert(t).That(FromInt(coefMax), Equals(Dnum{coefMax, +1, 16}))
	Assert(t).That(FromInt(-coefMax), Equals(Dnum{coefMax, -1, 16}))
	Assert(t).That(FromInt(1000000000000000000),
		Equals(Dnum{1000000000000000, +1, 19}))
}

func Test_ToInt(t *testing.T) {
	test := func(n int) {
		t.Helper()
		n2, ok := FromInt(int64(n)).ToInt()
		if !ok {
			t.Error("ToInt", n, FromInt(int64(n)), "failed")
		} else if n2 != n {
			t.Error("expected:", n, "got:", n2)
		}
	}
	test(0)
	test(1)
	test(-1)
	test(math.MinInt32)
	test(math.MaxInt32)
}

func Test_FromToFloat(t *testing.T) {
	assert := Assert(t)
	cvt := func(f float64) {
		t.Helper()
		assert.That(FromFloat(f).ToFloat(), Equals(f))
		assert.That(FromFloat(-f).ToFloat(), Equals(-f))
	}
	// special cases
	cvt(math.Inf(1))
	cvt(math.Inf(-1))

	// integer conversion
	cvt(0.0)
	cvt(123.0)
	cvt(123e3)

	// float conversion
	cvt(.1234)
	cvt(1.0 / 3.0)
	cvt(123456789e99)
	cvt(123456789e-99)
	cvt(1234567890123456e10)
	cvt(1000000000000001)
	for f := 1e15; f < 1e25; f *= 10 {
		cvt(f)
	}

	assert.That(FromFloat(1e200), Equals(PosInf))
	assert.That(FromFloat(-1e200), Equals(NegInf))
	assert.That(FromFloat(1e-200), Equals(Zero))
	assert.That(FromFloat(-1e-200), Equals(Zero))
}

func Test_Neg(t *testing.T) {
	assert := Assert(t)
	Neg := func(x string, y string) {
		xn := FromStr(x)
		yn := FromStr(y)
		assert.That(xn.Neg(), Equals(yn))
		assert.That(yn.Neg(), Equals(xn))
	}
	Neg("0", "0")
	Neg("123", "-123")
	Neg("inf", "-inf")
}

func Test_Compare(t *testing.T) {
	assert := Assert(t)
	data := []string{
		"-inf", "-1e9", "-123", "-1e-9", "0", "1e-9", "123", "1e9", "inf"}
	for i, xs := range data {
		x := FromStr(xs)
		assert.That(Compare(x, x), Equals(0).Comment(fmt.Sprint(x, " >< ", x)))
		for _, ys := range data[i+1:] {
			y := FromStr(ys)
			assert.That(Compare(x, y), Equals(-1).Comment(fmt.Sprint(x, " >< ", y)))
			assert.That(Compare(y, x), Equals(1).Comment(fmt.Sprint(y, " >< ", x)))
		}
	}
}

func Test_Add(t *testing.T) {
	assert := Assert(t)
	add := func(x string, y string, expected string) {
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
		assert.That(Add(xn, yn), Equals(zn))
		assert.That(Add(yn, xn), Equals(zn))
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
	add("1e12", "1e14", "1.01e14")
	add("1111111111111111", "2222222222222222e-4", "1111333333333333")
	add("1111111111111111", "6666666666666666e-4", "1111777777777778")
	// exceeds alignment
	add("123", "1e-99", "123")
	add("1e-99", "123", "123")
}

func Test_Sub(t *testing.T) {
	assert := Assert(t)
	sub := func(x string, y string, expected string) {
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
		assert.That(Sub(xn, yn), Equals(zn))
		if expected != "0" {
			assert.That(Sub(yn, xn), Equals(zn.Neg()))
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
	sub("1e50", "123", "1e50")
	sub("1e14", "1e12", "9.9e13")
	sub("12222222222222222222", "11111111111111111111e-4", "12221111111111111111")
}

func Test_Mul(t *testing.T) {
	assert := Assert(t)
	mul := func(x, y, expected string) {
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
		mul2 := func(x, y, zn Dnum) {
			assert.That(Mul(xn, yn), Equals(zn).Comment(fmt.Sprint(xn, " * ", yn)))
		}
		mul2(xn, yn, zn)
		mul2(yn, xn, zn)
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

	mul("2e9", "333e-9", "666")
	mul("2e3", "3e3", "6e6")
	mul("1.00000001", "1.00000001", "1.00000002")
	mul("1.000000001", "1.000000001", "1.000000002")
	mul(".4294967295", ".4294967295", ".1844674406511962")
	mul("1.12233445566", "1.12233445566", "1.259634630361628")
	mul("1.111111111111111", "1.111111111111111", "1.234567901234568")
	mul("1.23456789", "1.23456789", "1.524157875019052")
	mul("1.234567899", "1.234567899", "1.524157897241274")

	mul("2e99", "2e99", "inf") // exp overflow
}

func Test_Div(t *testing.T) {
	assert := Assert(t)
	div := func(x string, y string, expected string) {
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
		assert.That(Div(xn, yn), Equals(zn))
	}
	// special cases (no actual math)
	div("0", "0", "0")
	div("123", "0", "inf")
	div("123", "inf", "0")
	div("inf", "123", "inf")
	div("inf", "inf", "1")
	div("123", "123", "1")
	// exp overflow
	div("1e99", "1e-99", "inf")
	div("1e-99", "1e99", "0")
	// divides evenly
	div("4444", "2222", "2")
	div("2222", "4444", ".5")
	div("123000", ".000123", "1e9")
	// long division
	div("1", "3", ".3333333333333333333")
	div("2", "3", ".6666666666666666666")
	div("11", "17", ".6470588235294117647")
	div("1234567890123456", "9876543210123456", ".12499999887187493")
}

func Test_Format(t *testing.T) {
	test := func(s, mask, expected string) {
		t.Helper()
		dn := FromStr(s)
		Assert(t).That(dn.Format(mask), Equals(expected))
	}
	test("0", "#", "0")
	test("inf", "#", "#")
	test("0", "#.##", ".00")
	test("-1", "#", "-")
	test("1234", "##", "#")
	test("-123", "-###", "-123")
	test("-123", "(###)", "(123)")
	test("123", "(###)", "123 ")
	test("1234567", "###,###,###", "1,234,567")
	test(".8", "Foo", "#")
	// see also: suneido_tests/number.test
}

// benchmarks (for 1000 operations) ---------------------------------
/*
func BenchmarkAdd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 1; i < len(nums); i++ {
			Add(nums[i-1], nums[i])
		}
	}
}

func BenchmarkMul(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 1; i < len(nums); i++ {
			Mul(nums[i-1], nums[i])
		}
	}
}

func BenchmarkDiv(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 1; i < len(nums); i++ {
			Div(nums[i-1], nums[i])
		}
	}
}

var nums [1000]Dnum

func init() {
	for i := 0; i < len(nums); i++ {
		nums[i] = New(signPos, uint64(rand.Intn(1000000)), rand.Intn(9)-5)
	}
}
*/

var Bff Dnum

func BenchmarkFromFloat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Bff = FromFloat(123456e-99)
	}
}

// portable tests ---------------------------------------------------

func ptAdd(args []string, _ []bool) bool {
	xn := FromStr(args[0])
	yn := FromStr(args[1])
	zn := FromStr(args[2])
	return Add(xn, yn) == zn && Add(yn, xn) == zn
}

var _ = ptest.Add("dnum_add", ptAdd)

func ptSub(args []string, _ []bool) bool {
	xn := FromStr(args[0])
	yn := FromStr(args[1])
	zn := FromStr(args[2])
	return Sub(xn, yn) == zn &&
		(args[2] == "0" || Sub(yn, xn) == zn.Neg())
}

var _ = ptest.Add("dnum_sub", ptSub)

func ptMul(args []string, _ []bool) bool {
	xn := FromStr(args[0])
	yn := FromStr(args[1])
	zn := FromStr(args[2])
	return Mul(xn, yn) == zn && Mul(yn, xn) == zn
}

var _ = ptest.Add("dnum_mul", ptMul)

func ptDiv(args []string, _ []bool) bool {
	xn := FromStr(args[0])
	yn := FromStr(args[1])
	zn := FromStr(args[2])
	ok := Div(xn, yn) == zn
	if !ok {
		fmt.Println("got:", Div(xn, yn))
	}
	return ok
}

var _ = ptest.Add("dnum_div", ptDiv)

func ptCompare(args []string, _ []bool) bool {
	for i, xs := range args {
		x := FromStr(xs)
		if Compare(x, x) != 0 {
			return false
		}
		for _, ys := range args[i+1:] {
			y := FromStr(ys)
			if Compare(x, y) != -1 || Compare(y, x) != +1 {
				return false
			}
		}
	}
	return true
}

var _ = ptest.Add("dnum_cmp", ptCompare)

func TestPtest(t *testing.T) {
	if !ptest.RunFile("dnum.test") {
		t.Fail()
	}
}

/*
func closeTo(x Dnum) Tester {
	return func(actual interface{}) string {
		y := actual.(Dnum)
		if x.sign == y.sign && x.exp == y.exp &&
			(x.coef/10) == (y.coef/10) {
			return ""
		}
		return fmt.Sprintf("expected: %v but got: %v", x, y)
	}
}

func Test_ToUint(t *testing.T) {
	assert := Assert(t)
	test := func(x string, expected string) {
		dn := FromStr(x)
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
*/
