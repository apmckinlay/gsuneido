// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/strp"
)

// SuStrP is a small string optimization for strings up to 255 bytes.
// "P" can be read as "pointer" or as "Pascal" due to the leading size
// *SuStrP is a pointer so it can be stored in the Value interface
// without incurring an additional boxing allocation.
type SuStrP byte

// toStr returns the string contents without quotes or escaping.
func (sp *SuStrP) toStr() string {
	return strp.String((*byte)(sp))
}

// Len returns the length of the string.
func (sp *SuStrP) Len() int {
	return strp.Len((*byte)(sp))
}

func NewSuStrP(s string) *SuStrP {
	return (*SuStrP)(strp.From(s))
}

func CatSuStrP(parts ...string) *SuStrP {
	return (*SuStrP)(strp.Concat(parts...))
}

func SuStrP1(c byte) *SuStrP {
	return (*SuStrP)(strp.Single(c))
}

// Value implementation ---------------------------------------------

var _ Value = (*SuStrP)(nil)

func (sp *SuStrP) String() string {
	return escapeStr(sp.toStr(), 0)
}

func (sp *SuStrP) Display(th *Thread) string {
	q := 0
	if th != nil {
		q = th.Quote
	}
	return escapeStr(sp.toStr(), q)
}

func (sp *SuStrP) AsStr() (string, bool) {
	return sp.toStr(), true
}

func (sp *SuStrP) ToStr() (string, bool) {
	return sp.toStr(), true
}

func (sp *SuStrP) ToInt() (int, bool) {
	return 0, sp.Len() == 0
}

func (*SuStrP) IfInt() (int, bool) {
	return 0, false
}

func (sp *SuStrP) ToDnum() (dnum.Dnum, bool) {
	return dnum.Zero, sp.Len() == 0
}

func (*SuStrP) ToContainer() (Container, bool) {
	return nil, false
}

func (sp *SuStrP) Get(_ *Thread, key Value) Value {
	return strGet(sp.toStr(), key)
}

func (*SuStrP) Put(*Thread, Value, Value) {
	panic("string does not support put")
}

func (*SuStrP) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("string does not support update")
}

func (sp *SuStrP) RangeTo(from int, to int) Value {
	size := sp.Len()
	from = prepFrom(from, size)
	to = prepTo(from, to, size)
	if from == 0 && to == size {
		return sp
	}
	return StrVal(sp.toStr()[from:to])
}

func (sp *SuStrP) RangeLen(from int, n int) Value {
	size := sp.Len()
	from = prepFrom(from, size)
	n = prepLen(n, size-from)
	if from == 0 && n == size {
		return sp
	}
	return StrVal(sp.toStr()[from : from+n])
}

func (sp *SuStrP) Hash() uint64 {
	return hash.String(sp.toStr())
}

func (sp *SuStrP) Hash2() uint64 {
	return sp.Hash()
}

func (sp *SuStrP) Equal(other any) bool {
	// because SuStrP is the majority it's worth checking first
	if op, ok := other.(*SuStrP); ok {
		return sp.toStr() == op.toStr() // usual fast path
	}
	if v, ok := other.(Value); ok {
		if s, ok := v.ToStr(); ok {
			return sp.toStr() == s
		}
	}
	return false
}

func (*SuStrP) Type() types.Type {
	return types.String
}

func (sp *SuStrP) Compare(other Value) int {
	if cmp := cmp.Compare(ordStr, Order(other)); cmp != 0 {
		return cmp * 2
	}
	return strings.Compare(sp.toStr(), ToStr(other))
}

func (sp *SuStrP) Call(th *Thread, this Value, as *ArgSpec) Value {
	return stringCall(th, as, sp.toStr())
}

func (*SuStrP) Lookup(th *Thread, method string) Value {
	return Lookup(th, StringMethods, gnStrings, method)
}

func (*SuStrP) SetConcurrent() {
	// immutable so ok
}

// Packable implementation ------------------------------------------

var _ Packable = (*SuStrP)(nil)

func (sp *SuStrP) PackSize(*uint64) int {
	if sp.Len() == 0 {
		return 0
	}
	return 1 + sp.Len()
}

func (sp *SuStrP) PackSize2(*uint64, PackStack) int {
	return sp.PackSize(nil)
}

func (sp *SuStrP) Pack(_ *uint64, buf *pack.Encoder) {
	if sp.Len() > 0 {
		buf.Put1(PackString).PutStr(sp.toStr())
	}
}

// iterator ---------------------------------------------------------

func (sp *SuStrP) Iter() Iter {
	return &stringIter{s: sp.toStr()}
}
