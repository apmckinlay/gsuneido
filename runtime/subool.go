// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuBool is a boolean Value
type SuBool bool

// NOTE: converting bool/SuBool to interface{} doesn't seem to allocate
// e.g. v Value = SuBool(b)

// Value interface --------------------------------------------------

var _ Value = (*SuBool)(nil)

func (b SuBool) ToInt() (int, bool) {
	return 0, b == false
}

func (b SuBool) IfInt() (int, bool) {
	return 0, false
}

func (b SuBool) ToDnum() (dnum.Dnum, bool) {
	return dnum.Zero, b == false
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
	if b == true {
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

func (SuBool) RangeTo(int, int) Value {
	panic("boolean does not support range")
}

func (SuBool) RangeLen(int, int) Value {
	panic("boolean does not support range")
}

func (b SuBool) Hash() uint32 {
	if b == false {
		return 0x11111111
	}
	return 0x22222222
}

func (b SuBool) Hash2() uint32 {
	return b.Hash()
}

func (b SuBool) Equal(other interface{}) bool {
	if b2, ok := other.(SuBool); ok {
		return b == b2
	}
	return false
}

func (SuBool) Type() types.Type {
	return types.Boolean
}

func (b SuBool) Compare(other Value) int {
	if cmp := ints.Compare(ordBool, Order(other)); cmp != 0 {
		return cmp
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

func (SuBool) Lookup(*Thread, string) Callable {
	return nil
}

func (SuBool) SetConcurrent() {
}

// Packable interface -----------------------------------------------

var _ Packable = SuBool(true)

func (SuBool) PackSize(*int32) int {
	return 1
}

func (SuBool) PackSize2(int32, packStack) int {
	return 1
}

func (SuBool) PackSize3() int {
	return 1
}

func (b SuBool) Pack(_ int32, buf *pack.Encoder) {
	if b == true {
		buf.Put1(PackTrue)
	} else {
		buf.Put1(PackFalse)
	}
}
