// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

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
	CantConvert
	dnum.Dnum
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

func (dn SuDnum) AsStr() (string, bool) {
	return dn.Dnum.String(), true
}

func (dn SuDnum) String() string {
	return dn.Dnum.String()
}

func (SuDnum) Get(*Thread, Value) Value {
	panic("number does not support get")
}

func (SuDnum) Put(*Thread, Value, Value) {
	panic("number does not support put")
}

func (SuDnum) RangeTo(int, int) Value {
	panic("number does not support range")
}

func (SuDnum) RangeLen(int, int) Value {
	panic("number does not support range")
}

func (dn SuDnum) Hash() uint32 {
	if n, ok := dn.ToInt64(); ok && MinSuInt <= n && n <= MaxSuInt {
		return uint32(n) // for compatibility with SuInt
	}
	return dn.Dnum.Hash()
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

func (SuDnum) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Number")
}

// NumMethods is initialized by the builtin package
var NumMethods Methods

var gnNumbers = Global.Num("Numbers")

func (SuDnum) Lookup(t *Thread, method string) Callable {
	return Lookup(t, NumMethods, gnNumbers, method)
}

// Packable interface ===============================================

var _ Packable = SuDnum{}

// new format -------------------------------------------------------
// first byte is tag - PackPlus or PackMinus
// zero is just tag
// second byte is exponent
// following bytes encode two decimal digits (i.e. 0 to 99) per byte
// plus/minus infinite is PackPlus/Minus, (0xff), 0xff
//
// Encode decimal so that numbers with less digits are smaller.
// With coeficient maximized, binary encoding is longer.

const E14 = uint64(1e14)
const E12 = uint64(1e12)
const E10 = uint64(1e10)
const E8 = uint64(1e8)
const E6 = uint64(1e6)
const E4 = uint64(1e4)
const E2 = uint64(1e2)

func (dn SuDnum) PackSize(*int32) int {
	if dn.Sign() == 0 {
		return 1 // just tag
	}
	if dn.IsInf() {
		return 3
	}
	coef := dn.Coef()
	// unrolled, partly because mod by constant can be faster
	coef %= E14
	if coef == 0 {
		return 3
	}
	coef %= E12
	if coef == 0 {
		return 4
	}
	coef %= E10
	if coef == 0 {
		return 5
	}
	coef %= E8
	if coef == 0 {
		return 6
	}
	coef %= E6
	if coef == 0 {
		return 7
	}
	coef %= E4
	if coef == 0 {
		return 8
	}
	coef %= E2
	if coef == 0 {
		return 9
	}
	return 10
}

func (dn SuDnum) PackSize2(int32, packStack) int {
	return dn.PackSize(nil)
}

func (dn SuDnum) PackSize3() int {
	return dn.PackSize(nil)
}

func (dn SuDnum) Pack(_ int32, buf *pack.Encoder) {
	xor := byte(0)
	if dn.Sign() < 0 {
		xor = 0xff
		buf.Put1(PackMinus)
	} else {
		buf.Put1(PackPlus)
	}
	if dn.Sign() == 0 {
		return
	}
	if dn.IsInf() {
		buf.Put2(^xor, ^xor)
		return
	}

	// exponent
	buf.Put1(byte(dn.Exp()) ^ 0x80 ^ xor)

	// coefficient
	coef := dn.Coef()
	// unrolled, partly because div/mod by constant can be faster
	buf.Put1(byte(coef/E14) ^ xor)
	coef %= E14
	if coef == 0 {
		return
	}
	buf.Put1(byte(coef/E12) ^ xor)
	coef %= E12
	if coef == 0 {
		return
	}
	buf.Put1(byte(coef/E10) ^ xor)
	coef %= E10
	if coef == 0 {
		return
	}
	buf.Put1(byte(coef/E8) ^ xor)
	coef %= E8
	if coef == 0 {
		return
	}
	buf.Put1(byte(coef/E6) ^ xor)
	coef %= E6
	if coef == 0 {
		return
	}
	buf.Put1(byte(coef/E4) ^ xor)
	coef %= E4
	if coef == 0 {
		return
	}
	buf.Put1(byte(coef/E2) ^ xor)
	coef %= E2
	if coef == 0 {
		return
	}
	buf.Put1(byte(coef) ^ xor)
}

func UnpackNumber(s string) Value {
	if len(s) <= 1 {
		return Zero
	}
	sign := int8(s[0]-PackMinus)*2 - 1 // -1 or +1
	xor := byte(0)
	if sign < 0 {
		xor = 0xff
	}
	if s[2] == ^xor {
		return SuDnum{Dnum: dnum.Inf(sign)}
	}

	exp := s[1] ^ 0x80 ^ xor

	coef := uint64(0)
	switch len(s) {
	case 10:
		coef += uint64(s[9] ^ xor)
		fallthrough
	case 9:
		coef += uint64(s[8]^xor) * E2
		fallthrough
	case 8:
		coef += uint64(s[7]^xor) * E4
		fallthrough
	case 7:
		coef += uint64(s[6]^xor) * E6
		fallthrough
	case 6:
		coef += uint64(s[5]^xor) * E8
		fallthrough
	case 5:
		coef += uint64(s[4]^xor) * E10
		fallthrough
	case 4:
		coef += uint64(s[3]^xor) * E12
		fallthrough
	case 3:
		coef += uint64(s[2]^xor) * E14
	default:
		panic("invalid packed number length")
	}
	dn := dnum.Raw(sign, coef, int(exp))
	if n, ok := dn.ToInt(); ok && int(int16(n)) == n {
		return SuInt(n)
	}
	return SuDnum{Dnum: dn}
}

// old format -------------------------------------------------------

const maxShiftable = math.MaxUint16 / 10000

func UnpackNumberOld(s string) SuDnum {
	if len(s) <= 1 {
		return SuDnum{Dnum: dnum.Zero}
	}
	buf := pack.NewDecoder(s)
	sign := int8(+1)
	if buf.Get1() == PackMinus {
		sign = -1
	}
	exp := int8(buf.Get1())
	if exp == 0 {
		return SuDnum{Dnum: dnum.NegInf}
	}
	if exp == -1 {
		return SuDnum{Dnum: dnum.PosInf}
	}
	if sign < 0 {
		exp = ^exp
	}
	exp = exp ^ -128
	exp = exp - int8(buf.Remaining()/2)

	coef := unpackLongPartOld(buf, sign < 0)

	return SuDnum{Dnum: dnum.New(sign, coef, int(exp)*4+16)}
}

func unpackLongPartOld(buf *pack.Decoder, minus bool) uint64 {
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
