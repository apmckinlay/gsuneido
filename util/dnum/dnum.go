/*
Package dnum implements decimal floating point numbers.

Uses uint64 to hold the coefficient and int8 for exponent.

Value is sign * coef * 10^exp, zeroed value is 0.
*/
package dnum

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/util/bits"
)

type Dnum struct {
	coef uint64
	sign int8
	exp  int8
}

const (
	signPos  = +1
	signZero = 0
	signNeg  = -1
	expInf   = math.MaxInt8
)

var (
	Zero     = Dnum{}
	One      = Dnum{1, signPos, 0}
	Inf      = Dnum{math.MaxUint64, signPos, expInf}
	MinusInf = Dnum{math.MaxUint64, signNeg, expInf}

	OutsideRange = errors.New("outside range")
)

func NewDnum(neg bool, coef uint64, exp int8) Dnum {
	if neg {
		return Dnum{coef, signNeg, exp}
	} else {
		return Dnum{coef, signPos, exp}
	}
}

func abs(n int32) uint64 {
	n64 := int64(n)
	if n64 < 0 {
		n64 = -n64
	}
	return uint64(n64)
}

// Parse convert a string to a Dnum
func Parse(s string) (Dnum, error) {
	if len(s) < 1 {
		return Zero, errors.New("cannot convert empty string to Dnum")
	}
	if s == "0" {
		return Zero, nil
	}
	var dn Dnum
	dn.sign = signPos
	i := 0
	if s[i] == '+' {
		i++
	} else if s[i] == '-' {
		dn.sign = signNeg
		i++
	}
	before := spanDigits(s[i:])
	i += len(before)
	after := ""
	if i < len(s) && s[i] == '.' {
		i++
		after = spanDigits(s[i:])
		i += len(after)
	}
	after = strings.TrimRight(after, "0")
	coef, err := strconv.ParseUint(before+after, 10, 64)
	if err != nil {
		return Zero, errors.New("invalid number (" + s + ")")
	}
	dn.coef = coef

	exp := 0
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		e, err := strconv.ParseInt(s[i:], 10, 8)
		if err != nil {
			return Zero, errors.New("invalid exponent (" + s + ")")
		}
		exp = int(e)
	}
	if coef == 0 {
		return Zero, nil
	}
	exp -= len(after)
	if exp < -127 || exp >= 127 {
		return Zero, errors.New("exponent out of range (" + s + ")")
	}
	dn.exp = int8(exp)
	return dn, nil
}

// spanDigits returns the leading span of digits
func spanDigits(s string) string {
	i := 0
	for i < len(s) && '0' <= s[i] && s[i] <= '9' {
		i++
	}
	return s[:i]
}

// String converts a Dnum to a string.
// If the exponent is 0 it will print the number as an integer.
// Otherwise it will try to avoid scientific notation
// adding up to 4 zeroes at the end or 3 zeroes at the beginning.
func (dn Dnum) String() string {
	if dn == Zero {
		return "0"
	}
	sign := ""
	if dn.sign == signNeg {
		sign = "-"
	}
	if dn.IsInf() {
		return sign + "inf"
	}
	exp := int(dn.exp)
	digits := strconv.FormatUint(dn.coef, 10)
	if 0 <= exp && exp <= 4 {
		// decimal to the right
		digits += strings.Repeat("0", exp)
		return sign + digits
	}
	sexp := ""
	if -len(digits)-4 < exp && exp <= -len(digits) {
		// decimal to the left
		digits = "." + strings.Repeat("0", -exp-len(digits)) + digits
	} else if -len(digits) < exp && exp <= -1 {
		// decimal within
		i := len(digits) + exp
		digits = digits[:i] + "." + digits[i:]
	} else {
		// use scientific notation
		exp += len(digits) - 1
		digits = digits[:1] + "." + digits[1:]
		sexp = "e" + strconv.FormatInt(int64(exp), 10)
	}
	digits = strings.TrimRight(digits, "0")
	digits = strings.TrimRight(digits, ".")
	return sign + digits + sexp
}

func (dn Dnum) Float64() float64 {
	if dn.IsInf() {
		return math.Inf(int(dn.sign))
	}
	g := float64(dn.coef)
	if dn.sign == signNeg {
		g = -g
	}
	e := math.Pow10(int(dn.exp))
	return g * e
}

func FromFloat64(dn float64) Dnum {
	switch {
	case math.IsInf(dn, +1):
		return Inf
	case math.IsInf(dn, -1):
		return MinusInf
	case math.IsNaN(dn):
		panic("dnum.FromFloat64 can't convert NaN")
	}
	s := strconv.FormatFloat(dn, 'g', -1, 64)
	g, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return g
}

func (dn Dnum) Uint64() (uint64, error) {
	if dn.sign == signNeg {
		return 0, OutsideRange
	}
	return dn.toUint()
}

func (dn Dnum) Int64() (int64, error) {
	ui, err := dn.toUint()
	if err != nil {
		return 0, err
	}
	if dn.sign == signPos && ui > math.MaxInt64 {
		return math.MaxInt32, OutsideRange
	}
	if dn.sign == signNeg && ui > -math.MinInt64 {
		return math.MinInt32, OutsideRange
	}
	n := int64(ui)
	if dn.sign == signNeg {
		n = -n
	}
	return n, nil
}

func (dn Dnum) Int32() (int32, error) {
	ui, err := dn.toUint()
	if err != nil {
		return 0, err
	}
	if dn.sign == signPos && ui > math.MaxInt32 {
		return math.MaxInt32, OutsideRange
	}
	if dn.sign == signNeg && ui > -math.MinInt32 {
		return math.MinInt32, OutsideRange
	}
	n := int32(ui)
	if dn.sign == signNeg {
		n = -n
	}
	return n, nil
}

// try to make the exponent zero
// if exponent is too small return 0
// if exponent is too large return error
// result does not include sign
func (dn Dnum) toUint() (uint64, error) {
	for dn.exp > 0 && dn.shiftLeft() {
	}
	roundup := false
	for dn.exp < 0 && dn.shiftRight(&roundup) {
	}
	if roundup {
		dn.coef++
	}

	if dn.exp > 0 {
		return 0, OutsideRange
	} else if dn.exp < 0 {
		return 0, nil
	} else {
		return dn.coef, nil
	}
}

func FromInt64(n int64) Dnum {
	switch {
	case n > 0:
		return Dnum{uint64(n), signPos, 0}
	case n < 0:
		return Dnum{uint64(-n), signNeg, 0}
	default:
		return Zero
	}
}

// Sign returns -1 for negative, 0 for zero, and +1 for positive
func (dn Dnum) Sign() int {
	return int(dn.sign)
}

func (dn Dnum) Coef() uint64 {
	return dn.coef
}

func (dn Dnum) Exp() int {
	return int(dn.exp)
}

func (dn Dnum) IsInt() bool {
	coef := dn.coef
	exp := dn.exp
	for exp < 0 && coef%10 == 0 {
		coef /= 10
		exp++
	}
	return exp >= 0
}

// arithmetic operations -------------------------------------------------------

func (dn Dnum) Neg() Dnum {
	return Dnum{dn.coef, -dn.sign, dn.exp}
}

func (dn Dnum) Abs() Dnum {
	if dn == Zero {
		return Zero
	} else {
		return Dnum{dn.coef, signPos, dn.exp}
	}
}

func Cmp(x, y Dnum) int {
	switch {
	case x.sign < y.sign:
		return -1
	case x.sign > y.sign:
		return 1
	case x == y:
		return 0
	}
	return int(Sub(x, y).sign)
}

func Add(x, y Dnum) Dnum {
	switch {
	case x == Zero:
		return y
	case y == Zero:
		return x
	case x == Inf:
		if y == MinusInf {
			return Zero
		} else {
			return Inf
		}
	case x == MinusInf:
		if y == Inf {
			return Zero
		} else {
			return MinusInf
		}
	case y == Inf:
		return Inf
	case y == MinusInf:
		return MinusInf
	case x.sign != y.sign:
		return usub(x, y)
	default:
		return uadd(x, y)
	}
}

func Sub(x, y Dnum) Dnum {
	switch {
	case x == Zero:
		return y.Neg()
	case y == Zero:
		return x
	case x == Inf:
		if y == Inf {
			return Zero
		} else {
			return Inf
		}
	case x == MinusInf:
		if y == MinusInf {
			return Zero
		} else {
			return MinusInf
		}
	case y == Inf:
		return MinusInf
	case y == MinusInf:
		return Inf
	case x.sign != y.sign:
		return uadd(x, y)
	default:
		return usub(x, y)
	}
}

func uadd(x, y Dnum) Dnum {
	align(&x, &y)
	// align may make coef 0 if exp is too different
	coef := x.coef + y.coef
	if coef < x.coef || coef < y.coef { // overflow
		roundup := false
		x.shiftRight(&roundup)
		if roundup {
			x.coef++
		}
		y.shiftRight(&roundup)
		if roundup {
			y.coef++
		}
		coef = x.coef + y.coef
	}
	return result(coef, x.sign, int(x.exp))
}

func align(x, y *Dnum) {
	if x.exp == y.exp {
		return
	}
	if x.exp > y.exp {
		x, y = y, x // swap
	}
	for y.exp > x.exp && y.shiftLeft() {
	}
	roundup := false
	for y.exp > x.exp && x.shiftRight(&roundup) {
	}
	if x.exp != y.exp {
		x.exp = y.exp
	} else if roundup {
		x.coef++
	}
}

// returns true if it was able to shift (losslessly)
func (dn *Dnum) shiftLeft() bool {
	if !mul10safe(dn.coef) {
		return false
	}
	dn.coef *= 10
	// don't decrement past min
	if dn.exp > math.MinInt8 {
		dn.exp--
	}
	return true
}

func mul10safe(n uint64) bool {
	const HI4 = 0xf << 60
	return (n & HI4) == 0
}

// NOTE: may lose precision and round
// returns false only if coef is 0
func (dn *Dnum) shiftRight(roundup *bool) bool {
	*roundup = false
	if dn.coef == 0 {
		return false
	}
	*roundup = (dn.coef % 10) >= 5
	dn.coef /= 10
	// don't increment past max
	if dn.exp < math.MaxInt8 {
		dn.exp++
	}
	return true
}

func result(coef uint64, sign int8, exp int) Dnum {
	switch {
	case exp >= expInf:
		return inf(sign)
	case exp < math.MinInt8 || coef == 0:
		return Zero
	default:
		return Dnum{coef, sign, int8(exp)}
	}
}

func usub(x, y Dnum) Dnum {
	align(&x, &y)
	sign := x.sign
	if x.coef < y.coef {
		x, y = y, x
		sign *= -1 // flip sign
	}
	return result(x.coef-y.coef, sign, int(x.exp))
}

func Mul(x, y Dnum) Dnum {
	sign := x.sign * y.sign
	switch {
	case x == One:
		return y
	case y == One:
		return x
	case x == Zero || y == Zero:
		return Zero
	case x.IsInf() || y.IsInf():
		return inf(sign)
	}
	x.minCoef()
	y.minCoef()
	if bits.Nlz(x.coef)+bits.Nlz(y.coef) >= 64 {
		// coef won't overflow
		return result(x.coef*y.coef, sign, int(x.exp)+int(y.exp))
	}
	// drop 5 least significant bits off x and y
	// 59 bits < 18 decimal digits
	// 32 bits > 9 decimal digits
	// so we can split x & y into two pieces
	// and multiply separately guaranteed not to overflow
	xlo, xhi := x.split()
	ylo, yhi := y.split()
	exp := int(x.exp) + int(y.exp)
	r1 := result(xlo*ylo, sign, exp)
	r2 := result(xlo*yhi, sign, exp+9)
	r3 := result(xhi*ylo, sign, exp+9)
	r4 := result(xhi*yhi, sign, exp+18)
	return Add(r1, Add(r2, Add(r3, r4)))
}

// makes coef as small as possible (losslessly)
// i.e. trim trailing zero decimal digits
func (dn *Dnum) minCoef() {
	roundup := false
	for dn.coef > 0 && dn.coef%10 == 0 {
		dn.shiftRight(&roundup)
	}
	if roundup {
		dn.coef++
	}
}

func (dn *Dnum) split() (lo, hi uint64) {
	const HI5 = 0x1f << 59
	roundup := false
	for dn.coef&HI5 != 0 {
		dn.shiftRight(&roundup)
	}
	if roundup {
		dn.coef++
	}
	const NINE = 1000000000
	return dn.coef % NINE, dn.coef / NINE
}

func Div(x, y Dnum) Dnum {
	sign := x.sign * y.sign
	switch {
	case x == Zero:
		return Zero
	case y == Zero:
		return inf(x.sign)
	case x.IsInf():
		if y.IsInf() {
			return Dnum{1, sign, 0}
		}
		return inf(sign)
	case y.IsInf():
		return Zero
	}
	coef, exp := div2(x.coef, y.coef)
	return result(coef, sign, int(x.exp)-int(y.exp)+exp)
}

func div2(x, y uint64) (uint64, int) {
	exp := 0
	// strip trailing zeroes from y i.e. shift right as far as possible
	for y%10 == 0 {
		y /= 10
		exp--
	}
	var z uint64
	for x > 0 {
		// shift x left until divisible or as far as possible
		for x%y != 0 && mul10safe(x) && mul10safe(z) {
			x *= 10
			z *= 10
			exp--
		}
		for x < y {
			if !mul10safe(z) {
				return z, exp
			}
			y /= 10
			z *= 10
			exp--
		}
		q := (x / y)
		if q == 0 {
			break
		}
		z += q
		x %= y
	}
	return z, exp
}

func (dn Dnum) IsInf() bool {
	return dn.exp == expInf
}

func inf(sign int8) Dnum {
	switch sign {
	case signPos:
		return Inf
	case signNeg:
		return MinusInf
	default:
		panic("invalid sign")
	}
}

func (dn Dnum) Hash() uint32 {
	return uint32(dn.coef>>32) ^ uint32(dn.coef) ^
		uint32(dn.sign)<<16 ^ uint32(dn.exp)<<8
}
