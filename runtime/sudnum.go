package runtime

import (
	"math"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuDnum wraps a Dnum and implements Value and Packable
type SuDnum struct {
	dnum.Dnum
	CantConvert
}

// Value interface --------------------------------------------------

var _ Value = (*SuDnum)(nil)

func (dn SuDnum) ToInt() (int, bool) {
	return dn.Dnum.ToInt()
}

func (dn SuDnum) IfInt() (int, bool) {
	return dn.Dnum.ToInt()
}

func (dn SuDnum) ToDnum() (dnum.Dnum, bool) {
	return dn.Dnum, true
}

func (dn SuDnum) ToStr() (string, bool) {
	return dn.Dnum.String(), true
}

func (dn SuDnum) String() string {
	return dn.Dnum.String()
}

func (SuDnum) Get(*Thread, Value) Value {
	panic("number does not support get")
}

func (SuDnum) Put(Value, Value) {
	panic("number does not support put")
}

func (SuDnum) RangeTo(int, int) Value {
	panic("number does not support range")
}

func (SuDnum) RangeLen(int, int) Value {
	panic("number does not support range")
}

func (dn SuDnum) Hash2() uint32 {
	return dn.Hash()
}

func (dn SuDnum) Equal(other interface{}) bool {
	if d2, ok := other.(SuDnum); ok {
		return dnum.Equal(dn.Dnum, d2.Dnum)
	} else if i, ok := SuIntToInt(other); ok {
		return dnum.Equal(dn.Dnum, dnum.FromInt(int64(i)))
	}
	return false
}

func (SuDnum) Type() types.Type {
	return types.Number
}

func (dn SuDnum) Compare(other Value) int {
	if cmp := ints.Compare(ordNum, Order(other)); cmp != 0 {
		return cmp
	}
	// now know other is a number and ToDnum won't panic
	return dnum.Compare(dn.Dnum, ToDnum(other))
}

func (SuDnum) Call(*Thread, *ArgSpec) Value {
	panic("can't call Number")
}

// NumMethods is initialized by the builtin package
var NumMethods Methods

var gnNumbers = Global.Num("Numbers")

func (SuDnum) Lookup(method string) Value {
	return Lookup(NumMethods, gnNumbers, method)
}

// Packing (old format) ---------------------------------------------

var _ Packable = SuDnum{}

var pow10 = [...]uint64{1, 10, 100, 1000}

// PackSize returns the packed size of an SuDnum
func (dn SuDnum) PackSize(int) int {
	if dn.IsZero() {
		return 1
	}
	if dn.IsInf() {
		return 2
	}
	e := int(dn.Exp())
	n := dn.Coef()
	var p int
	if e > 0 {
		p = 4 - (e % 4)
	} else {
		p = abs(e) % 4
	}
	if p != 4 {
		n /= pow10[p]
	}
	n %= 1000000000000
	if n == 0 {
		return 4
	}
	n %= 100000000
	if n == 0 {
		return 6
	}
	n %= 10000
	if n == 0 {
		return 8
	}
	return 10
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Pack packs the SuDnum into buf (which must be large enough)
func (dn SuDnum) Pack(buf *pack.Encoder) {
	sign := dn.Sign()
	if sign >= 0 {
		buf.Put1(packPlus)
	} else {
		buf.Put1(packMinus)
	}
	if sign == 0 {
		return
	}
	if dn.IsInf() {
		if sign < 0 {
			buf.Put1(0)
		} else {
			buf.Put1(255)
		}
		return
	}
	packDnum(sign < 0, dn.Coef(), dn.Exp(), buf)
}

func packDnum(neg bool, coef uint64, exp int, buf *pack.Encoder) {
	var p int
	if exp > 0 {
		p = 4 - (exp % 4)
	} else {
		p = abs(exp) % 4
	}
	if p != 4 {
		coef /= pow10[p] // may lose up to 3 digits of precision
		exp += p
	}
	exp /= 4
	buf.Put1(scale(exp, neg))
	packCoef(buf, coef, neg)
}

func scale(exp int, neg bool) byte {
	eb := (byte(exp) ^ 0x80) & 0xff
	if neg {
		eb = (^eb) & 0xff
	}
	return eb
}

const e12 = 1000000000000
const e8 = 100000000
const e4 = 10000

func packCoef(buf *pack.Encoder, coef uint64, neg bool) {
	flip := uint16(0)
	if neg {
		flip = 0xffff
	}
	buf.Uint16(uint16(coef/e12) ^ flip)
	coef %= e12
	if coef == 0 {
		return
	}
	buf.Uint16(uint16(coef/e8) ^ flip)
	coef %= e8
	if coef == 0 {
		return
	}
	buf.Uint16(uint16(coef/e4) ^ flip)
	coef %= e4
	if coef == 0 {
		return
	}
	buf.Uint16(uint16(coef) ^ flip)
}

const maxShiftable = math.MaxUint16 / 10000

// UnpackNumber unpacks an SuInt or SuDnum
func UnpackNumber(s string) Value {
	if len(s) <= 1 {
		return Zero
	}
	buf := pack.NewDecoder(s)
	sign := int8(+1)
	if buf.Get1() == packMinus {
		sign = -1
	}
	exp := int8(buf.Get1())
	if exp == 0 {
		return SuDnum{Dnum: dnum.NegInf}
	}
	if exp == -1 {
		return SuDnum{Dnum: dnum.Inf}
	}
	if sign < 0 {
		exp = ^exp
	}
	exp = exp ^ -128
	exp = exp - int8(buf.Remaining()/2)

	coef := unpackLongPart(buf, sign < 0)

	if exp == 1 && coef <= maxShiftable {
		coef *= 10000
		exp--
	}
	if exp == 0 && coef <= MaxSuInt {
		return SuInt(int(sign) * int(coef))
	}
	return SuDnum{Dnum: dnum.New(sign, coef, int(exp)*4+16)}
}

func unpackLongPart(buf *pack.Decoder, minus bool) uint64 {
	flip := uint16(0)
	if minus {
		flip = 0xffff
	}
	n := uint64(0)
	for buf.Remaining() > 0 {
		n = n*10000 + uint64(buf.Uint16()^flip)
	}
	return n
}
