// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

// SuBool is a boolean Value
type SuBool bool

// NOTE: converting bool/SuBool to any doesn't seem to allocate
// e.g. v Value = SuBool(b)

// Value interface --------------------------------------------------

var _ Value = (*SuBool)(nil)

func (b SuBool) ToInt() (int, bool) {
	return 0, b == false //nolint:gosimple
}

func (b SuBool) IfInt() (int, bool) {
	return 0, false
}

func (b SuBool) ToDnum() (dnum.Dnum, bool) {
	return dnum.Zero, b == false //nolint:gosimple
}

func (SuBool) ToContainer() (Container, bool) {
	return nil, false
}

func (b SuBool) AsStr() (string, bool) {
	return b.String(), true
}

func (b SuBool) ToStr() (string, bool) {
	return "", false
}

func (b SuBool) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (SuBool) Get(*Thread, Value) Value {
	panic("boolean does not support get")
}

func (SuBool) Put(*Thread, Value, Value) {
	panic("boolean does not support put")
}

func (SuBool) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("boolean does not support update")
}

func (SuBool) RangeTo(int, int) Value {
	panic("boolean does not support range")
}

func (SuBool) RangeLen(int, int) Value {
	panic("boolean does not support range")
}

func (b SuBool) Hash() uint64 {
	if !b {
		return 0x11111111
	}
	return 0x22222222
}

func (b SuBool) Hash2() uint64 {
	return b.Hash()
}

func (b SuBool) Equal(other any) bool {
	return b == other
}

func (SuBool) Type() types.Type {
	return types.Boolean
}

func (b SuBool) Compare(other Value) int {
	if cmp := cmp.Compare(ordBool, Order(other)); cmp != 0 {
		return cmp * 2
	}
	if b == other {
		return 0
	} else if b {
		return 1
	} else {
		return -1
	}
}

func (b SuBool) Not() SuBool {
	return SuBool(!bool(b))
}

func (SuBool) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Boolean")
}

func (SuBool) Lookup(*Thread, string) Value {
	return nil
}

func (SuBool) SetConcurrent() {
}

// Packable interface -----------------------------------------------

var _ Packable = SuBool(true)

func (SuBool) PackSize(*packing) int {
	return 1
}

func (b SuBool) Pack(pk *packing) {
	if b {
		pk.Put1(PackTrue)
	} else {
		pk.Put1(PackFalse)
	}
}
