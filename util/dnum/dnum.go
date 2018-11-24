/*
Package dnum implements decimal floating point numbers.

Uses uint64 to hold the coefficient and int8 for exponent.
Only uses 16 decimal digits.

Value is sign * .coef * 10^exp, i.e. assumed decimal to left
Coefficient is kept "maximized" in 16 decimal digits.
Zeroed value is 0.
*/
package dnum

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
)

// Dnum is a decimal floating point number
type Dnum struct {
	coef uint64
	sign int8
	exp  int8
}

const (
	signPosInf = +2
	signPos    = +1
	signZero   = 0
	signNeg    = -1
	signNegInf = -2
	expMin     = math.MinInt8
	expMax     = math.MaxInt8
	coefMin    = 1000000000000000
	coefMax    = 9999999999999999
	digitsMax  = 16
	shiftMax   = digitsMax - 1
)

// common values
var (
	Zero   = Dnum{}
	One    = Dnum{1000000000000000, signPos, 1}
	NegOne = Dnum{1000000000000000, signNeg, 1}
	Inf    = Dnum{1, signPosInf, 0}
	NegInf = Dnum{1, signNegInf, 0}
)

var pow10 = [...]uint64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000}

var halfpow10 = [...]uint64{
	0,
	5,
	50,
	500,
	5000,
	50000,
	500000,
	5000000,
	50000000,
	500000000,
	5000000000,
	50000000000,
	500000000000,
	5000000000000,
	50000000000000,
	500000000000000,
	5000000000000000,
	50000000000000000,
	500000000000000000,
	5000000000000000000}

// NOTE: comment out body in production
func check(_ bool) {
	// if !cond {
	// 	panic("verify failed")
	// }
}

// FromInt returns a Dnum for an int
func FromInt(n int64) Dnum {
	if n == 0 {
		return Zero
	}
	sign := int8(signPos)
	if n < 0 {
		n = -n
		sign = signNeg
	}
	coef := uint64(n)
	p := maxShift(coef)
	coef *= pow10[p]
	check(coefMin <= coef && coef <= coefMax)
	exp := int8(digitsMax - p)
	return Dnum{coef, sign, exp}
}

// FromFloat converts a float64 to a Dnum
func FromFloat(f float64) Dnum {
	switch {
	case math.IsInf(f, +1):
		return Inf
	case math.IsInf(f, -1):
		return NegInf
	case math.IsNaN(f):
		panic("dnum.FromFloat64 can't convert NaN")
	}
	n := int64(f)
	if f == float64(n) {
		return FromInt(n)
	}
	s := strconv.FormatFloat(f, 'g', -1, 64)
	return FromStr(s)
}

// New constructs a Dnum, maximizing coef and handling exp out of range
// Used to normalize results of operations
func New(sign int8, coef uint64, exp int) Dnum {
	if sign == 0 || coef == 0 || exp < expMin {
		return Zero
	} else if sign == signPosInf {
		return Inf
	} else if sign == signNegInf {
		return NegInf
	} else {
		atmax := false
		for coef > coefMax {
			coef = (coef + 5) / 10 // drop/round least significant digit
			exp++
			atmax = true
		}
		if !atmax {
			p := maxShift(coef)
			coef *= pow10[p]
			exp -= p
		}
		if exp > expMax {
			return inf(sign)
		}
		return Dnum{coef, sign, int8(exp)}
	}
}

func maxShift(x uint64) int {
	i := ilog10(x)
	if i > shiftMax {
		return 0
	}
	return shiftMax - i
}

func ilog10(x uint64) int {
	// based on Hacker's Delight
	if x == 0 {
		return 0
	}
	y := (19 * (63 - bits.LeadingZeros64(x))) >> 6
	if y < 18 && x >= pow10[y+1] {
		y++
	}
	return y
}

func inf(sign int8) Dnum {
	switch {
	case sign < 0:
		return NegInf
	case sign > 0:
		return Inf
	default:
		return Zero
	}
}

// String returns a string representation of the Dnum
func (dn Dnum) String() string {
	if dn.sign == 0 {
		return "0"
	}
	const maxLeadingZeros = 7
	sign := ""
	if dn.sign < 0 {
		sign = "-"
	}
	if dn.IsInf() {
		return sign + "inf"
	}
	digits := getDigits(dn.coef)
	nd := len(digits)
	e := int(dn.exp) - nd
	if -maxLeadingZeros <= dn.exp && dn.exp <= 0 {
		// decimal to the left
		return sign + "." + strings.Repeat("0", -e-nd) + digits
	} else if -nd < e && e <= -1 {
		// decimal within
		dec := nd + e
		return sign + digits[:dec] + "." + digits[dec:]
	} else if 0 < dn.exp && dn.exp <= digitsMax {
		// decimal to the right
		return sign + digits + strings.Repeat("0", e)
	} else {
		// scientific notation
		after := ""
		if nd > 1 {
			after = "." + digits[1:]
		}
		return sign + digits[:1] + after + "e" + strconv.Itoa(int(dn.exp-1))
	}
}

func getDigits(coef uint64) string {
	var digits [digitsMax]byte
	i := shiftMax
	nd := 0
	for coef != 0 {
		digits[nd] = byte('0' + (coef / pow10[i]))
		coef %= pow10[i]
		nd++
		i--
	}
	return string(digits[:nd])
}

// FromStr parses a numeric string and returns a Dnum representation.
func FromStr(s string) Dnum {
	r := &reader{s, 0}
	sign := getSign(r)
	if r.matchStr("inf") {
		return inf(sign)
	}
	coef, exp := getCoef(r)
	exp += getExp(r)
	if r.len() != 0 { // didn't consume entire string
		panic("invalid number")
	} else if coef == 0 || exp < math.MinInt8 {
		return Zero
	} else if exp > math.MaxInt8 {
		return inf(sign)
	}
	check(coefMin <= coef && coef <= coefMax)
	return Dnum{coef, sign, int8(exp)}
}

type reader struct {
	s string
	i int
}

func (r *reader) cur() byte {
	if r.i >= len(r.s) {
		return 0
	}
	return byte(r.s[r.i])
}

func (r *reader) prev() byte {
	if r.i == 0 {
		return 0
	}
	return byte(r.s[r.i-1])
}

func (r *reader) len() int {
	return len(r.s) - r.i
}

func (r *reader) match(c byte) bool {
	if r.cur() == c {
		r.i++
		return true
	}
	return false
}

func (r *reader) matchDigit() bool {
	c := r.cur()
	if '0' <= c && c <= '9' {
		r.i++
		return true
	}
	return false
}

func (r *reader) matchStr(pre string) bool {
	if strings.HasPrefix(r.s[r.i:], pre) {
		r.i += len(pre)
		return true
	}
	return false
}

func getSign(r *reader) int8 {
	if r.match('-') {
		return int8(signNeg)
	}
	r.match('+')
	return int8(signPos)
}

func getCoef(r *reader) (uint64, int) {
	digits := false
	beforeDecimal := true
	for r.match('0') {
		digits = true
	}
	if r.cur() == '.' && r.len() > 1 {
		digits = false
	}
	n := uint64(0)
	exp := 0
	p := shiftMax
	for {
		c := r.cur()
		if r.matchDigit() {
			digits = true
			// ignore extra decimal places
			if c != '0' && p >= 0 {
				n += uint64(c-'0') * pow10[p]
			}
			p--
		} else if beforeDecimal {
			// decimal point or end
			exp = shiftMax - p
			if !r.match('.') {
				break
			}
			beforeDecimal = false
			if !digits {
				for r.match('0') {
					digits = true
					exp--
				}
			}
		} else {
			break
		}
	}
	if !digits {
		panic("numbers require at least one digit")
	}
	return n, exp
}

func getExp(r *reader) int {
	e := 0
	if r.match('e') || r.match('E') {
		esign := getSign(r)
		for r.matchDigit() {
			e = e*10 + int(r.prev()-'0')
		}
		e *= int(esign)
	}
	return e
}

// end of FromStr ---------------------------------------------------

// IsInf returns true if a Dnum is positive or negative infinite
func (dn Dnum) IsInf() bool {
	return dn.sign == signPosInf || dn.sign == signNegInf
}

// IsZero returns true if a Dnum is positive or negative infinite
func (dn Dnum) IsZero() bool {
	return dn.sign == signZero
}

// ToFloat converts a Dnum to float64
func (dn Dnum) ToFloat() float64 {
	if dn.IsInf() {
		return math.Inf(int(dn.sign))
	}
	g := float64(dn.coef)
	if dn.sign == signNeg {
		g = -g
	}
	e := math.Pow10(int(dn.exp) - digitsMax)
	return g * e
}

// ToInt converts a Dnum to an int64, returning whether it was convertible
func (dn Dnum) ToInt64() (int64, bool) {
	if dn.sign == 0 {
		return 0, true
	}
	if dn.sign != signNegInf && dn.sign != signPosInf {
		if 0 < dn.exp && dn.exp < digitsMax &&
			(dn.coef%pow10[digitsMax-dn.exp]) == 0 { // usual case
			return int64(dn.sign) * int64(dn.coef/pow10[digitsMax-dn.exp]), true
		}
		if dn.exp == digitsMax {
			return int64(dn.sign) * int64(dn.coef), true
		}
		if dn.exp == digitsMax+1 {
			return int64(dn.sign) * (int64(dn.coef) * 10), true
		}
		if dn.exp == digitsMax+2 {
			return int64(dn.sign) * (int64(dn.coef) * 100), true
		}
		if dn.exp == digitsMax+3 && dn.coef < math.MaxInt64/1000 {
			return int64(dn.sign) * (int64(dn.coef) * 1000), true
		}
	}
	return 0, false
}

func (dn Dnum) ToInt() (int, bool) {
	n, ok := dn.ToInt64()
	if !ok || int64(int(n)) != n {
		return 0, false
	}
	return int(n), true
}

// Sign returns -1 for negative, 0 for zero, and +1 for positive
func (dn Dnum) Sign() int {
	return int(dn.sign)
}

// Coef returns the coefficient
func (dn Dnum) Coef() uint64 {
	return dn.coef
}

// Exp returns the exponent
func (dn Dnum) Exp() int {
	return int(dn.exp)
}

// Frac returns the fractional portion, i.e. x - x.Int()
func (dn Dnum) Frac() Dnum {
	if dn.sign == 0 || dn.sign == signNegInf || dn.sign == signPosInf ||
		dn.exp >= digitsMax {
		return Zero
	}
	if dn.exp <= 0 {
		return dn
	}
	frac := dn.coef % pow10[digitsMax-dn.exp]
	if frac == dn.coef {
		return dn
	}
	return Dnum{frac, dn.sign, dn.exp}
}

type RoundingMode int

const (
	UP RoundingMode = iota
	DOWN
	HALF_UP
)

// Int returns the integer portion (truncating any fractional part)
func (dn Dnum) Int() Dnum {
	return dn.integer(DOWN)
}

func (dn Dnum) integer(mode RoundingMode) Dnum {
	if dn.sign == 0 || dn.sign == signNegInf || dn.sign == signPosInf ||
		dn.exp >= digitsMax {
		return dn
	}
	if dn.exp <= 0 {
		if mode == UP ||
			(mode == HALF_UP && dn.exp == 0 && dn.coef >= One.coef*5) {
			return New(dn.sign, One.coef, int(dn.exp)+1)
		}
		return Zero
	}
	e := digitsMax - dn.exp
	frac := dn.coef % pow10[e]
	if frac == 0 {
		return dn
	}
	i := dn.coef - frac
	if (mode == UP && frac > 0) || (mode == HALF_UP && frac >= halfpow10[e]) {
		return New(dn.sign, i+pow10[e], int(dn.exp)) // normalize
	}
	return Dnum{i, dn.sign, dn.exp} // TODO doesn't need to normalize
}

// arithmetic operations -------------------------------------------------------

// Neg returns the Dnum negated i.e. sign reversed
func (dn Dnum) Neg() Dnum {
	return Dnum{dn.coef, -dn.sign, dn.exp}
}

// Abs returns the Dnum with a positive sign
func (dn Dnum) Abs() Dnum {
	if dn.sign < 0 {
		return Dnum{dn.coef, -dn.sign, dn.exp}
	}
	return dn
}

// Equal returns true if two Dnum's are equal
func Equal(x, y Dnum) bool {
	return x.sign == y.sign && x.exp == y.exp && x.coef == y.coef
}

// Compare compares two Dnum's returning -1 for <, 0 for ==, +1 for >
func Compare(x, y Dnum) int {
	switch {
	case x.sign < y.sign:
		return -1
	case x.sign > y.sign:
		return 1
	case x == y:
		return 0
	}
	sign := int(x.sign)
	switch {
	case sign == 0 || sign == signNegInf || sign == signPosInf:
		return 0
	case x.exp < y.exp:
		return -sign
	case x.exp > y.exp:
		return +sign
	case x.coef < y.coef:
		return -sign
	case x.coef > y.coef:
		return +sign
	default:
		return 0
	}
}

// Sub returns the difference of two Dnum's
func Sub(x, y Dnum) Dnum {
	return Add(x, y.Neg())
}

// Add returns the sum of two Dnum's
func Add(x, y Dnum) Dnum {
	switch {
	case x.sign == signZero:
		return y
	case y.sign == signZero:
		return x
	case x.IsInf():
		if y.sign == -x.sign {
			return Zero
		}
		return x
	case y.IsInf():
		return y
	}
	if !align(&x, &y) {
		return x
	}
	if x.sign != y.sign {
		return usub(x, y)
	}
	return uadd(x, y)
}

func uadd(x, y Dnum) Dnum {
	return New(x.sign, x.coef+y.coef, int(x.exp))
}

func usub(x, y Dnum) Dnum {
	if x.coef < y.coef {
		return New(-x.sign, y.coef-x.coef, int(x.exp))
	}
	return New(x.sign, x.coef-y.coef, int(x.exp))
}

func align(x, y *Dnum) bool {
	if x.exp == y.exp {
		return true
	}
	if x.exp < y.exp {
		*x, *y = *y, *x // swap
	}
	yshift := ilog10(y.coef)
	e := int(x.exp - y.exp)
	if e > yshift {
		return false
	}
	yshift = e
	check(0 <= yshift && yshift <= 20)
	y.coef = (y.coef + halfpow10[yshift]) / pow10[yshift]
	check(int(y.exp)+yshift == int(x.exp))
	return true
}

const e7 = 10000000

// Mul returns the product of two Dnum's
func Mul(x, y Dnum) Dnum {
	sign := x.sign * y.sign
	switch {
	case sign == signZero:
		return Zero
	case x.IsInf() || y.IsInf():
		return inf(sign)
	}
	e := int(x.exp) + int(y.exp)

	// split unevenly to use full 64 bit range to get more precision
	// and avoid needing xlo * ylo
	xhi := x.coef / e7 // 9 digits
	xlo := x.coef % e7 // 7 digits
	yhi := y.coef / e7 // 9 digits
	ylo := y.coef % e7 // 7 digits

	c := xhi * yhi
	if xlo != 0 || ylo != 0 {
		c += (xlo*yhi + ylo*xhi) / e7
	}
	return New(sign, c, e-2)
}

// Div returns the quotient of two Dnum's
func Div(x, y Dnum) Dnum {
	sign := x.sign * y.sign
	switch {
	case x.sign == signZero:
		return x
	case y.sign == signZero:
		return inf(x.sign)
	case x.IsInf():
		if y.IsInf() {
			if sign < 0 {
				return NegOne
			}
			return One
		}
		return inf(sign)
	case y.IsInf():
		return Zero
	}
	coef := div128(x.coef, y.coef)
	return New(sign, coef, int(x.exp)-int(y.exp))
}

// Hash returns a hash value for a Dnum
func (dn Dnum) Hash() uint32 {
	return uint32(dn.coef>>32) ^ uint32(dn.coef) ^
		uint32(dn.sign)<<16 ^ uint32(dn.exp)<<8
}

/*
// makes coef as small as possible (losslessly)
// i.e. trim trailing zero decimal digits
func (dn *Dnum) minCoef() {
	roundup := false
	for dn.coef > 0 && dn.coef%10 == 0 {
		dn.shiftRight(&roundup)
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
		for x < y { // max once, but not necessarily first iteration
			if !mul10safe(z) {
				return z, exp
			}
			y /= 10 // should this round ???
			z *= 10
			exp--
		}
		q := (x / y)
		if q == 0 {
			break
		}
		z += q
		x %= y // or x -= q * y ???
	}
	return z, exp
}
*/
