package runtime

import (
	"encoding/binary"
	"math"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// SuDnum wraps a Dnum and implements Value and Packable
type SuDnum struct {
	dnum.Dnum
	// use an anonymous member in a struct to include the methods from Dnum
	// i.e. Hash, String
}

var _ Value = SuDnum{}
var _ Packable = SuDnum{}

// Value interface --------------------------------------------------

// ToInt converts a SuDnum to an integer (Value interface)
func (dn SuDnum) ToInt() int {
	n, ok := dn.Dnum.ToInt()
	if !ok {
		panic("can't convert number to integer " + dn.String())
	}
	return n
}

// ToDnum returns the wrapped Dnum (Value interface)
func (dn SuDnum) ToDnum() dnum.Dnum {
	return dn.Dnum
}

// ToStr converts the Dnum to a string (Value interface)
func (dn SuDnum) ToStr() string {
	return dn.Dnum.String()
}

// String returns a string representation of the Dnum (Value interface)
func (dn SuDnum) String() string {
	return dn.Dnum.String()
}

// Get is not applicable to SuDnum (Value interface)
func (SuDnum) Get(Value) Value {
	panic("number does not support get")
}

// Put is not applicable to SuDnum (Value interface)
func (SuDnum) Put(Value, Value) {
	panic("number does not support put")
}

func (SuDnum) RangeTo(int, int) Value {
	panic("number does not support range")
}

func (SuDnum) RangeLen(int, int) Value {
	panic("number does not support range")
}

// Hash2 is used to hash nested values (Value interface)
func (dn SuDnum) Hash2() uint32 {
	return dn.Hash()
}

// Equals returns true if other is an equal SuDnum or integer (Value interface)
func (dn SuDnum) Equal(other interface{}) bool {
	if d2, ok := other.(SuDnum); ok {
		return dnum.Equal(dn.Dnum, d2.Dnum)
	} else if i, ok := SmiToInt(other); ok {
		return dnum.Equal(dn.Dnum, dnum.FromInt(int64(i)))
	}
	return false
}

// TypeName returns the name of this type (Value interface)
func (SuDnum) TypeName() string {
	return "Number"
}

// Order returns the ordering of SuDnum (Value interface)
func (SuDnum) Order() Ord {
	return ordNum
}

// Compare compares an SuDnum to another Value (Value interface)
func (dn SuDnum) Compare(other Value) int {
	if cmp := ints.Compare(dn.Order(), other.Order()); cmp != 0 {
		return cmp
	}
	if y, ok := other.(SuDnum); ok {
		return dnum.Compare(dn.Dnum, y.Dnum)
	}
	return dnum.Compare(dn.Dnum, other.ToDnum())
}

func (SuDnum) Call0(_ *Thread) Value {
	panic("can't call Number")
}
func (SuDnum) Call1(_ *Thread, _ Value) Value {
	panic("can't call Number")
}
func (SuDnum) Call2(_ *Thread, _, _ Value) Value {
	panic("can't call Number")
}
func (SuDnum) Call3(_ *Thread, _, _, _ Value) Value {
	panic("can't call Number")
}
func (SuDnum) Call4(_ *Thread, _, _, _, _ Value) Value {
	panic("can't call Number")
}
func (SuDnum) Call(*Thread, *ArgSpec) Value {
	panic("can't call Number")
}

func (SuDnum) Lookup(string) Callable {
	return nil // TODO
}

// Packing (old format) ---------------------------------------------

var pow10 = [...]uint64{1, 10, 100, 1000}

// PackSize returns the packed size of an SuDnum
func (dn SuDnum) PackSize() int {
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
func (dn SuDnum) Pack(buf []byte) []byte {
	// for performance we avoid append
	buf = buf[:1]
	sign := dn.Sign()
	if sign >= 0 {
		buf[0] = packPlus
	} else {
		buf[0] = packMinus
	}
	if sign == 0 {
		return buf
	}
	buf = buf[0:2]
	if dn.IsInf() {
		if sign < 0 {
			buf[1] = 0
		} else {
			buf[1] = 255
		}
		return buf
	}
	return packDnum(sign < 0, dn.Coef(), dn.Exp(), buf)
}

func packDnum(neg bool, coef uint64, exp int, buf []byte) []byte {
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
	buf[1] = scale(exp, neg)
	buf = packCoef(buf, coef, neg)
	return buf
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

func packCoef(buf []byte, coef uint64, neg bool) []byte {
	flip := uint16(0)
	if neg {
		flip = 0xffff
	}
	buf = buf[:4]
	binary.BigEndian.PutUint16(buf[2:], uint16(coef/e12)^flip)
	coef %= e12
	if coef == 0 {
		return buf
	}
	buf = buf[:6]
	binary.BigEndian.PutUint16(buf[4:], uint16(coef/e8)^flip)
	coef %= e8
	if coef == 0 {
		return buf
	}
	buf = buf[:6]
	binary.BigEndian.PutUint16(buf[6:], uint16(coef/e4)^flip)
	coef %= e4
	if coef == 0 {
		return buf
	}
	buf = buf[:8]
	binary.BigEndian.PutUint16(buf[8:], uint16(coef)^flip)
	return buf
}

const maxShiftable = math.MaxUint16 / 10000

// UnpackNumber unpacks an SuInt or SuDnum
func UnpackNumber(buf rbuf) Value {
	sign := int8(+1)
	if buf.get() == packMinus {
		sign = -1
	}
	if buf.remaining() == 0 {
		return SuInt(0)
	}
	exp := int8(buf.get())
	if exp == 0 {
		return SuDnum{dnum.NegInf}
	}
	if exp == -1 {
		return SuDnum{dnum.Inf}
	}
	if sign < 0 {
		exp = ^exp
	}
	exp = exp ^ -128
	exp = exp - int8(buf.remaining()/2)

	coef := unpackLongPart(buf, sign < 0)

	if exp == 1 && coef <= maxShiftable {
		coef *= 10000
		exp--
	}
	if exp == 0 && coef <= MaxSuInt {
		return SuInt(int(sign) * int(coef))
	}
	return SuDnum{dnum.New(sign, coef, int(exp)*4+16)}
}

func unpackLongPart(buf rbuf, minus bool) uint64 {
	flip := uint16(0)
	if minus {
		flip = 0xffff
	}
	n := uint64(0)
	for buf.remaining() > 0 {
		n = n*10000 + uint64(buf.getUint16()^flip)
	}
	return n
}
