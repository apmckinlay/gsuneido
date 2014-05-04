/*
Package dnum implements decimal floating point numbers
using uint64 to hold the coefficient.
*/
package dnum

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/util/bits"
)

// value is -1^sign * coef * 10^exp
// zeroed value = 0
type Dnum struct {
	coef uint64
	sign int8
	exp  int8
}

const (
	POSITIVE = 0
	NEGATIVE = 1
	INF_EXP  = math.MaxInt8
)

var (
	Zero     = Dnum{}
	One      = Dnum{1, POSITIVE, 0}
	Inf      = Dnum{exp: INF_EXP}
	MinusInf = Dnum{sign: NEGATIVE, exp: INF_EXP}

	OutsideRange = errors.New("outside range")
)

func NewDnum(neg bool, coef uint64, exp int8) Dnum {
	if neg {
		return Dnum{coef, NEGATIVE, exp}
	} else {
		return Dnum{coef, POSITIVE, exp}
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
	i := 0
	if s[i] == '+' {
		i++
	} else if s[i] == '-' {
		dn.sign = NEGATIVE
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

// spanDigits returns the number of leading digits
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
	if dn.sign == NEGATIVE {
		sign = "-"
	}
	if dn.exp == INF_EXP {
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
	digits = strings.TrimRight(digits, "0.")
	return sign + digits + sexp
}

func (dn Dnum) Float64() float64 {
	if dn.IsInf() {
		return math.Inf(-int(dn.sign))
	}
	g := float64(dn.coef)
	if dn.sign == NEGATIVE {
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
	if dn.sign == NEGATIVE {
		return 0, OutsideRange
	}
	return dn.toUint()
}

func (dn Dnum) Int64() (int64, error) {
	ui, err := dn.toUint()
	if err != nil {
		return 0, err
	}
	if dn.sign == POSITIVE && ui > math.MaxInt64 {
		return math.MaxInt32, OutsideRange
	}
	if dn.sign == NEGATIVE && ui > -math.MinInt64 {
		return math.MinInt32, OutsideRange
	}
	n := int64(ui)
	if dn.sign == NEGATIVE {
		n = -n
	}
	return n, nil
}

func (dn Dnum) Int32() (int32, error) {
	ui, err := dn.toUint()
	if err != nil {
		return 0, err
	}
	if dn.sign == POSITIVE && ui > math.MaxInt32 {
		return math.MaxInt32, OutsideRange
	}
	if dn.sign == NEGATIVE && ui > -math.MinInt32 {
		return math.MinInt32, OutsideRange
	}
	n := int32(ui)
	if dn.sign == NEGATIVE {
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
		return Dnum{uint64(n), POSITIVE, 0}
	case n < 0:
		return Dnum{uint64(-n), NEGATIVE, 0}
	default:
		return Zero
	}
}

// Sign returns -1 for negative, 0 for zero, and +1 for positive
func (dn Dnum) Sign() int {
	switch {
	case dn == Zero:
		return 0
	case dn.sign == NEGATIVE:
		return -1
	default:
		return +1
	}
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
	if dn == Zero {
		return Zero
	} else {
		return Dnum{dn.coef, dn.sign ^ 1, dn.exp}
	}
}

func Cmp(x, y Dnum) int {
	switch {
	case x == y:
		return 0
	case x == MinusInf, y == Inf:
		return -1
	case x == Inf, y == MinusInf:
		return 1
	}
	switch d := Sub(x, y); {
	case d == Zero:
		return 0
	case d.sign == NEGATIVE:
		return -1
	default:
		return +1
	}
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
	sign := x.sign
	align(&x, &y)
	if x.coef == 0 {
		return Dnum{y.coef, sign, y.exp}
	}
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
	return result(coef, sign, int(x.exp))
}

func align(x, y *Dnum) (flipped int8) {
	if x.exp == y.exp {
		return
	}
	if x.exp > y.exp {
		*x, *y = *y, *x // swap
		flipped = 1
	}
	for y.exp > x.exp && y.shiftLeft() {
	}
	roundup := false
	for y.exp > x.exp && x.shiftRight(&roundup) {
	}
	if roundup {
		x.coef++
	}
	return
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

const HI4 = 0xf << 60

func mul10safe(n uint64) bool {
	return (n & HI4) == 0
}

// BUG rounds incorrectly if used repeatedly
// e.g. 123.456 will round to 124

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
	case exp >= INF_EXP:
		return inf(sign)
	case exp < math.MinInt8 || coef == 0:
		return Zero
	default:
		return Dnum{coef, sign, int8(exp)}
	}
}

func usub(x, y Dnum) Dnum {
	sign := x.sign
	sign ^= align(&x, &y)
	if x.coef < y.coef {
		x, y = y, x
		sign ^= 1 // flip sign
	}
	return result(x.coef-y.coef, sign, int(x.exp))
}

func Mul(x, y Dnum) Dnum {
	sign := x.sign ^ y.sign
	switch {
	case x == Zero || y == Zero:
		return Zero
	case x.IsInf() || y.IsInf():
		return result(0, sign, INF_EXP)
	}
	x.minCoef()
	y.minCoef()
	if bits.Nlz(x.coef)+bits.Nlz(y.coef) >= 64 {
		// coef won't overflow
		if int(x.exp)+int(y.exp) >= INF_EXP {
			return result(0, sign, INF_EXP)
		}
		return result(x.coef*y.coef, sign, int(x.exp)+int(y.exp))
	}
	// drop 5 least significant bits off x and y
	// 59 bits < 18 decimal digits
	// 32 bits > 9 decmal digits
	// so we can split x & y into two pieces
	// and multiply separately guaranteed not to overflow
	xlo, xhi := x.split()
	ylo, yhi := y.split()
	exp := x.exp + y.exp
	r1 := result(xlo*ylo, sign, int(exp))
	r2 := result(xlo*yhi, sign, int(exp)+9)
	r3 := result(xhi*ylo, sign, int(exp)+9)
	r4 := result(xhi*yhi, sign, int(exp)+18)
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

// makes coef as large as possible (losslessly)
// i.e. trim leading zero decimal digits
func (dn *Dnum) maxCoef() {
	for dn.shiftLeft() {
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
	sign := x.sign ^ y.sign
	switch {
	case x == Zero:
		return Zero
	case y == Zero || x.IsInf():
		return inf(sign)
	case y.IsInf():
		return Zero
	}
	if x.coef%y.coef == 0 {
		// divides evenly
		return result(x.coef/y.coef, sign, int(x.exp)-int(y.exp))
	}
	x.maxCoef()
	y.minCoef()
	if x.coef%y.coef == 0 {
		// divides evenly
		return result(x.coef/y.coef, sign, int(x.exp)-int(y.exp))
	}
	return longDiv(x, y)
}

func longDiv(x, y Dnum) Dnum {
	// shift y so it is just less than x
	xdiv10 := x.coef / 10
	for y.coef < xdiv10 && y.shiftLeft() {
	}
	exp := int(x.exp) - int(y.exp)
	rem := x.coef
	ycoef := y.coef
	coef := uint64(0)
	// each iteration calculates one digit of the result
	for rem != 0 && mul10safe(coef) {
		// shift so y is less than the remainder
		for ycoef > rem {
			rem, ycoef = shift(rem, ycoef)
			coef *= 10
			exp--
		}
		if ycoef == 0 {
			break
		}
		// subtract repeatedly
		for rem >= ycoef {
			rem -= ycoef
			coef++
		}
	}
	// round final digit
	if 2*rem >= ycoef {
		coef++
	}
	return result(coef, x.sign^y.sign, exp)
}

// shift x left (preferably) or y right
func shift(x, y uint64) (x2, y2 uint64) {
	if mul10safe(x) {
		x *= 10
	} else {
		roundup := (y % 10) >= 5
		y /= 10
		if roundup {
			y++
		}
	}
	return x, y
}

func (dn Dnum) IsInf() bool {
	return dn.exp == INF_EXP
}

func inf(sign int8) Dnum {
	switch sign {
	case POSITIVE:
		return Inf
	case NEGATIVE:
		return MinusInf
	default:
		panic("invalid sign")
	}
}

func (dn Dnum) Hash() uint32 {
	return uint32(dn.coef>>32) ^ uint32(dn.coef) ^
		uint32(dn.sign)<<16 ^ uint32(dn.exp)<<8
}
