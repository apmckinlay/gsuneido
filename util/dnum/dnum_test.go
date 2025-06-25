// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dnum

import (
	"fmt"
	"math"
	// "math/rand/v2"
	"testing"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ptest"
)

func Test_size(t *testing.T) {
	// due to allignment and padding, size is 16 bytes instead of 10
	assert.T(t).This(int(unsafe.Sizeof(Dnum{}))).Is(16)
	var a [10]Dnum
	assert.T(t).This(int(unsafe.Sizeof(a))).Is(160)
}

func Test_inf(t *testing.T) {
	assert := assert.T(t).This
	assert(Inf(0)).Is(Zero)
	assert(Inf(+1)).Is(PosInf)
	assert(Inf(-1)).Is(NegInf)
}

func Test_ilog10(t *testing.T) {
	assert.T(t).This(ilog10(0)).Is(0)
	assert.T(t).This(ilog10(123)).Is(2)
}

func Test_New(t *testing.T) {
	assert := assert.T(t).This
	assert(New(signZero, 0, 0)).Is(Zero)
	assert(New(signPos, 1, 999)).Is(PosInf) // exponent overflow
	assert(New(signNeg, 1, 999)).Is(NegInf) // exponent overflow
	assert(New(signPos, 1, -999)).Is(Zero)  // exponent underflow
	assert(New(signNeg, 1, -999)).Is(Zero)  // exponent underflow
	assert(New(signPos, 1, 0)).Is(Dnum{1000000000000000, 1, -15})
	assert(New(signPos, 123, 0)).Is(Dnum{1230000000000000, 1, -13})
}

func Test_String(t *testing.T) {
	assert := assert.T(t).This
	assert(Zero.String()).Is("0")
	assert(One.String()).Is("1")
	assert(PosInf.String()).Is("inf")
	assert(NegInf.String()).Is("-inf")
	assert(FromInt(123).String()).Is("123")
	assert(FromInt(-123).String()).Is("-123")

	assert(New(signPos, 1234000000000000, -20).String()).Is("1.234e-21")
	assert(New(signPos, 1234000000000000, -2).String()).Is(".001234")
	assert(New(signPos, 1234000000000000, 0).String()).Is(".1234")
	assert(New(signPos, 1234000000000000, 2).String()).Is("12.34")
	assert(New(signPos, 1234000000000000, 4).String()).Is("1234")
	assert(New(signPos, 1234000000000000, 6).String()).Is("123400")
	assert(New(signPos, 1234000000000000, 20).String()).Is("1.234e19")
}

func Test_FromStr(t *testing.T) {
	assert := assert.T(t).This
	assert(FromStr("inf")).Is(PosInf)
	assert(FromStr("+inf")).Is(PosInf)
	assert(FromStr("-inf")).Is(NegInf)
	assert(FromStr("0")).Is(Zero)
	assert(FromStr("+0")).Is(Zero)
	assert(FromStr("-0")).Is(Zero)
	assert(FromStr("0e4")).Is(Zero)
	assert(FromStr("0000")).Is(Zero)
	assert(FromStr("0000.")).Is(Zero)
	assert(FromStr(".0000")).Is(Zero)
	assert(FromStr("0000.0000")).Is(Zero)
	assert(FromStr("1")).Is(One)
	assert(FromStr("000000000000000000001")).Is(One)
	assert(FromStr("1.0000000000000000000")).Is(One)
	assert(FromStr("100000000000000000000")).Is(FromStr("1e20"))
	assert(FromStr(".1234567890123456789")).Is(FromStr(".1234567890123456"))
	assert(FromStr(".000000000000000000001")).Is(FromStr(".1e-20"))
}

func Test_FromToStr(t *testing.T) {
	test := func(s string) {
		assert.T(t).This(FromStr(s).String()).Is(s)
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
	assert.T(t).This(e).Is(20)
}

func Test_FromToInt(t *testing.T) {
	assert := assert.T(t)
	test := func(x int64) {
		n, ok := FromInt(x).ToInt64()
		assert.True(ok)
		assert.This(n).Is(x)
		n, ok = FromInt(-x).ToInt64()
		assert.True(ok)
		assert.This(n).Is(-x)
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
	assert := assert.T(t).This
	assert(FromInt(0)).Is(Zero)
	assert(FromInt(1)).Is(Dnum{1000000000000000, +1, 1})
	assert(FromInt(100)).Is(Dnum{1000000000000000, +1, 3})
	assert(FromInt(123)).Is(Dnum{1230000000000000, +1, 3})
	assert(FromInt(-123)).Is(Dnum{1230000000000000, -1, 3})
	assert(FromInt(coefMax)).Is(Dnum{coefMax, +1, 16})
	assert(FromInt(-coefMax)).Is(Dnum{coefMax, -1, 16})
	assert(FromInt(1000000000000000000)).Is(Dnum{1000000000000000, +1, 19})
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
	assert := assert.T(t).This
	cvt := func(f float64) {
		t.Helper()
		assert(FromFloat(f).ToFloat()).Is(f)
		assert(FromFloat(-f).ToFloat()).Is(-f)
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

	assert(FromFloat(1e200)).Is(PosInf)
	assert(FromFloat(-1e200)).Is(NegInf)
	assert(FromFloat(1e-200)).Is(Zero)
	assert(FromFloat(-1e-200)).Is(Zero)
}

func Test_Neg(t *testing.T) {
	Neg := func(x string, y string) {
		xn := FromStr(x)
		yn := FromStr(y)
		assert.T(t).This(xn.Neg()).Is(yn)
		assert.T(t).This(yn.Neg()).Is(xn)
	}
	Neg("0", "0")
	Neg("123", "-123")
	Neg("inf", "-inf")
}

func Test_Compare(t *testing.T) {
	assert := assert.T(t).This
	data := []string{
		"-inf", "-1e9", "-123", "-1e-9", "0", "1e-9", "123", "1e9", "inf"}
	for i, xs := range data {
		x := FromStr(xs)
		assert(Compare(x, x)).Msg(x, " >< ", x).Is(0)
		for _, ys := range data[i+1:] {
			y := FromStr(ys)
			assert(Compare(x, y)).Msg(x, " >< ", y).Is(-1)
			assert(Compare(y, x)).Msg(y, " >< ", x).Is(1)
		}
	}
}

func Test_Add(t *testing.T) {
	add := func(x string, y string, expected string) {
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
		assert.T(t).This(Add(xn, yn)).Is(zn)
		assert.T(t).This(Add(yn, xn)).Is(zn)
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
	add("1e-128", "1", "1")
}

func Test_Sub(t *testing.T) {
	sub := func(x string, y string, expected string) {
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
		assert.T(t).This(Sub(xn, yn)).Is(zn)
		if expected != "0" {
			assert.T(t).This(Sub(yn, xn)).Is(zn.Neg())
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
	mul := func(x, y, expected string) {
		mul2 := func(x, y, z Dnum) {
			assert.T(t).Msg(x, " * ", y).This(Mul(x, y)).Is(z)
		}
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
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
	div := func(x string, y string, expected string) {
		xn := FromStr(x)
		yn := FromStr(y)
		zn := FromStr(expected)
		assert.T(t).This(Div(xn, yn)).Is(zn)
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
		assert.T(t).This(dn.Format(mask)).Is(expected)
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
	for b.Loop() {
		for i := 1; i < len(nums); i++ {
			Add(nums[i-1], nums[i])
		}
	}
}

func BenchmarkMul(b *testing.B) {
	for b.Loop() {
		for i := 1; i < len(nums); i++ {
			Mul(nums[i-1], nums[i])
		}
	}
}

func BenchmarkDiv(b *testing.B) {
	for b.Loop() {
		for i := 1; i < len(nums); i++ {
			Div(nums[i-1], nums[i])
		}
	}
}

var nums [1000]Dnum

func init() {
	for i := range len(nums) {
		nums[i] = New(signPos, uint64(rand.Intn(1000000)), rand.Intn(9)-5)
	}
}
*/

var Bff Dnum

func BenchmarkFromFloat(b *testing.B) {
	for b.Loop() {
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
	return func(actual any) string {
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

func FuzzAdd(f *testing.F) {
	f.Fuzz(func(t *testing.T, xcoef uint64, xexp int8, xsign bool,
		ycoef uint64, yexp int8, ysign bool) {
		xn := mknum(xcoef, xexp, xsign)
		yn := mknum(ycoef, yexp, ysign)
		result := Add(xn, yn)
		result2 := Add(yn, xn)
		assert.T(t).Msg(xn, "+", yn).This(result).Is(result2)
	})
}

func mknum(coef uint64, exp int8, neg bool) Dnum {
	sign := int8(signPos)
	if neg {
		sign = signNeg
	}
	coef = coef%(coefMax-coefMin) + coefMin
	assert.That(coefMin <= coef && coef <= coefMax)
	return Dnum{coef: coef, exp: exp, sign: sign}
}

// func TestMknum(t *testing.T) {
// 	for range 100 {
// 		coef := rand.Uint64()
// 		exp := int8(rand.Int32() & 0xff)
// 		neg := rand.IntN(2) == 0
// 		n := mknum(coef, exp, neg)
// 		fmt.Println(n)
// 	}
// }
