// StrVal is a string Value

package value

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/verify"
)

type DnumVal struct {
	dnum.Dnum
	// use an anonymous member in a struct to include the methods from Dnum
}

var _ Value = DnumVal{} // confirm it implements Value

func ParseNum(s string) (Value, error) {
	dn, err := dnum.Parse(s)
	return DnumToValue(dn), err
}

func (dn DnumVal) ToInt() int32 {
	n, _ := dn.Int32()
	return n
}

func (dn DnumVal) ToDnum() dnum.Dnum {
	return dn.Dnum
}

func (dn DnumVal) ToStr() string {
	return dn.Dnum.String()
}

func (dn DnumVal) String() string {
	return dn.Dnum.String()
}

func (dn DnumVal) Get(key Value) Value {
	panic("number does not support get")
}

func (dn DnumVal) Put(key Value, val Value) {
	panic("number does not support put")
}

func (dn DnumVal) hash2() uint32 {
	return dn.Hash()
}

func (dn DnumVal) Equals(other interface{}) bool {
	if d2, ok := other.(DnumVal); ok {
		return 0 == dnum.Cmp(dn.Dnum, d2.Dnum)
	}
	return false
}

func (dn DnumVal) PackSize() int {
	return packSizeNum(dn)
}

func (dn DnumVal) Pack(buf []byte) []byte {
	return packNum(dn, buf)
}

// packing - ugly because of compatibility with cSuneido

//
type Num interface {
	Sign() int
	Exp() int
	Coef() uint64
}

func packSizeNum(num Num) int {
	coef := num.Coef()
	if coef == 0 {
		return 1
	}
	exp := num.Exp()
	exp, coef = adjustNum(exp, coef)
	ps := packshorts(coef)
	exp = exp/4 + ps
	if exp >= math.MaxInt8 {
		return 2
	}
	return 2 /* tag and exponent */ + 2*ps
}

// 16 digits - maximum precision that cSuneido handles
const MAX_PREC = 9999999999999999
const MAX_PREC_DIV_10 = 999999999999999

func adjustNum(exp int, coef uint64) (int, uint64) {
	// strip trailing zeroes
	for (coef % 10) == 0 {
		coef /= 10
		exp++
	}
	// adjust exp to a multiple of 4 (to match cSuneido)
	for (exp%4) != 0 && coef < MAX_PREC_DIV_10 {
		coef *= 10
		exp--
	}
	for (exp%4) != 0 || coef > MAX_PREC {
		coef /= 10
		exp++
	}
	return exp, coef
}

func packshorts(n uint64) int {
	i := 0
	for ; n != 0; i++ {
		n /= 10000
	}
	verify.That(i <= 4) // cSuneido limit
	return i
}

func packNum(num Num, buf []byte) []byte {
	sign := num.Sign()
	if sign >= 0 {
		buf = append(buf, PLUS)
	} else {
		buf = append(buf, MINUS)
	}
	coef := num.Coef()
	if coef == 0 {
		return buf
	}
	exp := num.Exp()
	exp, coef = adjustNum(exp, coef)
	exp = exp/4 + packshorts(coef)
	if exp >= math.MaxInt8 {
		if sign < 0 {
			return append(buf, 0)
		} else {
			return append(buf, 255)
		}
	}
	buf = append(buf, scale(sign, exp))
	return packCoef(sign, coef, buf)
}

func scale(sign int, exp int) byte {
	eb := byte(exp) ^ 0x80
	if sign < 0 {
		eb = (byte)((^eb) & 0xff)
	}
	return eb
}

func packCoef(sign int, coef uint64, buf []byte) []byte {
	var sh [4]uint16
	i := 0
	for ; coef != 0; i++ {
		sh[i] = digit(sign, coef)
		coef /= 10000
	}
	for i--; i >= 0; i-- {
		buf = append(buf, byte(sh[i]>>8), byte(sh[i]))
	}
	return buf
}

func digit(sign int, coef uint64) uint16 {
	n := uint16(coef % 10000)
	if sign < 0 {
		n = ^n
	}
	return n
}

func UnpackNumber(buf rbuf) Value {
	neg := buf.get() == MINUS
	if buf.remaining() == 0 {
		return IntVal(0)
	}
	e := buf.get()
	if e == 0 {
		return DnumVal{dnum.MinusInf}
	}
	if e == 255 {
		return DnumVal{dnum.Inf}
	}
	if neg {
		e = ^e
	}
	exp := int(e ^ 0x80)
	exp = (exp - buf.remaining()/2) * 4
	coef := unpackCoef(buf, neg)
	if coef == 0 {
		return IntVal(0)
	}
	if 10 >= exp && exp > 0 {
		for ; exp > 0 && coef < math.MaxUint64/10; exp-- {
			coef *= 10
		}
	}
	if exp == 0 && 0 <= coef && coef <= math.MaxInt32 {
		if neg {
			return IntVal(-int32(coef))
		} else {
			return IntVal(int32(coef))
		}
	} else {
		return DnumVal{dnum.NewDnum(neg, coef, int8(exp))}
	}
}

func unpackCoef(buf rbuf, neg bool) uint64 {
	var n uint64
	for buf.remaining() > 0 {
		x := buf.getUint16()
		if neg {
			x = ^x
		}
		n = n*10000 + uint64(x)
	}
	return n
}
