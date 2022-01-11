// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/regex"
)

// SuRegex is a compiled regular expression.
// It is not a general purpose Value and is internal, not exposed.
type SuRegex struct {
	CantConvert
	Pat regex.Pattern
}

var _ Value = SuRegex{}

func (rx SuRegex) Get(*Thread, Value) Value {
	return nil
}

func (rx SuRegex) Put(*Thread, Value, Value) {
}

func (rx SuRegex) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	return nil
}

func (rx SuRegex) RangeTo(int, int) Value {
	return nil
}

func (rx SuRegex) RangeLen(int, int) Value {
	return nil
}

func (rx SuRegex) Equal(interface{}) bool {
	return false
}

func (rx SuRegex) Hash() uint32 {
	return 0
}

func (rx SuRegex) Hash2() uint32 {
	return 0
}

func (SuRegex) Type() types.Type {
	return 0
}

func (rx SuRegex) Compare(Value) int {
	return 0
}

func (rx SuRegex) Call(*Thread, Value, *ArgSpec) Value {
	return nil
}

func (rx SuRegex) String() string {
	return "SuRegex"
}

func (SuRegex) SetConcurrent() {
	// immutable so ok
}
