// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found ax the LICENSE file.

package core

import (
	"cmp"
	"math"
	"strconv"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuInt64 is a 64-bit signed integer Value
type SuInt64 struct {
	ValueBase[SuInt64]
	int64
}

// Value interface

var _ Value = (*SuInt64)(nil)

func (si SuInt64) AsStr() (string, bool) {
	return si.String(), true
}

func (si SuInt64) Compare(other Value) int {
	if cmp := cmp.Compare(ordNum, Order(other)); cmp != 0 {
		return cmp * 2
	}
	if i2, ok := SuIntToInt(other); ok {
		return cmp.Compare(si.int64, int64(i2))
	}
	dn, _ := si.ToDnum()
	return dnum.Compare(dn, other.(SuDnum).Dnum)
}

func (si SuInt64) Equal(other any) bool {
	if i2, ok := SuIntToInt(other); ok {
		return si.int64 == int64(i2)
	}
	if dn, ok := other.(SuDnum); ok {
		if i2, ok := dn.IfInt(); ok {
			return si.int64 == int64(i2)
		}
	}
	return false
}

func (si SuInt64) Hash() uint64 {
	return uint64(si.int64) * phi64
}

func (si SuInt64) Hash2() uint64 {
	return si.Hash()
}

func (si SuInt64) IfInt() (int, bool) {
	return int(si.int64), true
}

func (si SuInt64) Lookup(th *Thread, method string) Value {
	return intLookup(th, method)
}

func (si SuInt64) SetConcurrent() {
	// immutable so ok
}

func (si SuInt64) String() string {
	return strconv.FormatInt(si.int64, 10)
}

func (si SuInt64) ToDnum() (dnum.Dnum, bool) {
	return dnum.FromInt(si.int64), true
}

func (si SuInt64) ToInt() (int, bool) {
	return int(si.int64), true
}

func (si SuInt64) ToStr() (string, bool) {
	return "", false
}

func (si SuInt64) Type() types.Type {
	return types.Number
}

// Packable interface -----------------------------------------------

var _ Packable = SuInt(0)

func (si SuInt64) PackSize(*uint64) int {
	return packSizeInt(si.int64)
}

func (si SuInt64) PackSize2(*uint64, packStack) int {
	return packSizeInt(si.int64)
}

func (si SuInt64) Pack(_ *uint64, enc *pack.Encoder) {
	packInt(si.int64, enc)
}

func packSizeInt(n int64) int {
	u := uint64(n)
	if u == 0 {
		return 1
	}
	if n < 0 {
		u = -u
	}
	size := 2 // tag and exponent
	var buf [20]byte
	i := 20
	for u > 0 {
		i--
		buf[i] = byte(u % 10)
		u /= 10
	}
	lim := 19
	for ; lim > 0 && buf[lim] == 0; lim-- {
	}
	for ; i < lim; i += 2 {
		size++
	}
	if i <= lim {
		size++
	}
	return size
}

func packInt(n int64, enc *pack.Encoder) {
	u := uint64(n)
	xor := byte(0)
	if n < 0 {
		u = -u
		xor = 0xff
		enc.Put1(PackMinus)
	} else {
		enc.Put1(PackPlus)
	}
	if u == 0 {
		return
	}
	var buf [20]byte
	i := 20
	for u > 0 {
		i--
		buf[i] = byte(u % 10)
		u /= 10
	}
	exp := 20 - i
	enc.Put1(byte(exp) ^ 0x80 ^ xor)
	lim := 19
	for ; lim > 0 && buf[lim] == 0; lim-- {
	}
	for ; i < lim; i += 2 {
		enc.Put1((buf[i]*10 + buf[i+1]) ^ xor)
	}
	if i <= lim {
		enc.Put1((buf[i] * 10) ^ xor)
	}
}

func unpackInt(s string, sign, exp int8, xor byte) Value {
	last := int8(2*(len(s)-2) - 1)
	i := 2
	u := uint64(0)
	for ; i < len(s)-1; i++ {
		d := s[i] ^ xor
		u = u*100 + uint64(d)
	}

	d := s[i] ^ xor
	if exp == last {
		u = u*10 + uint64(d)/10
	} else {
		u = u*100 + uint64(d)
		for j := last + 1; j < exp; j++ {
			u *= 10
		}
	}
	n := int(u)
	if sign == -1 {
		n = -n
	}
	return IntVal(n)
}

var PackedMinInt64 = Pack(SuInt64{int64: math.MinInt64})
var PackedMaxInt64 = Pack(SuInt64{int64: math.MaxInt64})
