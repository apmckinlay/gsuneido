// SuStr is a string Value

package value

import (
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/verify"
)

type SuDnum struct {
	dnum.Dnum
	// use an anonymous member in a struct to include the methods from Dnum
	// i.e. Hash, String
}

var _ Value = SuDnum{}
var _ Packable = SuDnum{}

func ParseNum(s string) (Value, error) {
	dn, err := dnum.Parse(s)
	return DnumToValue(dn), err
}

func (dn SuDnum) ToInt() int32 {
	n, _ := dn.Int32()
	return n
}

func (dn SuDnum) ToDnum() dnum.Dnum {
	return dn.Dnum
}

func (dn SuDnum) ToStr() string {
	return dn.Dnum.String()
}

func (dn SuDnum) Get(key Value) Value {
	panic("number does not support get")
}

func (dn SuDnum) Put(key Value, val Value) {
	panic("number does not support put")
}

func (dn SuDnum) hash2() uint32 {
	return dn.Hash()
}

func (dn SuDnum) Equals(other interface{}) bool {
	if d2, ok := other.(SuDnum); ok {
		return 0 == dnum.Cmp(dn.Dnum, d2.Dnum)
	} else if si, ok := other.(SuInt); ok {
		return 0 == dnum.Cmp(dn.Dnum, si.ToDnum())
	}
	return false
}

func (dn SuDnum) PackSize() int {
	return packSizeNum(dn)
}

func (dn SuDnum) Pack(buf []byte) []byte {
	return packNum(dn, buf)
}

// packing - ugly because of compatibility with cSuneido

type num interface {
	Sign() int
	Exp() int
	Coef() uint64
}

func packSizeNum(n num) int {
	coef := n.Coef()
	if coef == 0 {
		return 1
	}
	exp := n.Exp()
	exp, coef = adjustNum(exp, coef)
	ps := packshorts(coef)
	exp = exp/4 + ps
	if exp >= math.MaxInt8 {
		return 2
	}
	return 2 /* tag and exponent */ + 2*ps
}

func adjustNum(exp int, coef uint64) (int, uint64) {
	// 16 digits - maximum precision that cSuneido handles
	const MAX_PREC = 9999999999999999
	const MAX_PREC_DIV_10 = 999999999999999

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

func packNum(n num, buf []byte) []byte {
	sign := n.Sign()
	if sign >= 0 {
		buf = append(buf, packPlus)
	} else {
		buf = append(buf, packMinus)
	}
	coef := n.Coef()
	if coef == 0 {
		return buf
	}
	exp := n.Exp()
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
	neg := buf.get() == packMinus
	if buf.remaining() == 0 {
		return SuInt(0)
	}
	e := buf.get()
	if e == 0 {
		return SuDnum{dnum.MinusInf}
	}
	if e == 255 {
		return SuDnum{dnum.Inf}
	}
	if neg {
		e = ^e
	}
	exp := int(e ^ 0x80)
	exp = (exp - buf.remaining()/2) * 4
	coef := unpackCoef(buf, neg)
	if coef == 0 {
		return SuInt(0)
	}
	if 10 >= exp && exp > 0 {
		for ; exp > 0 && coef < math.MaxUint64/10; exp-- {
			coef *= 10
		}
	}
	if exp == 0 && 0 <= coef && coef <= math.MaxInt32 {
		if neg {
			return SuInt(-int32(coef))
		} else {
			return SuInt(int32(coef))
		}
	} else {
		return SuDnum{dnum.NewDnum(neg, coef, int8(exp))}
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

func (_ SuDnum) TypeName() string {
	return "Number"
}

func (_ SuDnum) order() ord {
	return ordNum
}

func (x SuDnum) cmp(other Value) int {
	if y, ok := other.(SuDnum); ok {
		return dnum.Cmp(x.Dnum, y.Dnum)
	} else {
		return dnum.Cmp(x.Dnum, y.ToDnum())
	}
}

func (_ SuDnum) Call(t AThread) Value {
	panic("can't call number")
}
