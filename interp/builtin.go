package interp

import (
	"reflect"

	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

type Builtin func(argSpec *ArgSpec, args ...Value) Value

var _ Value = Builtin(nil)

func (b Builtin) ToInt() int {
	panic("can't convert builtin function to integer")
}

func (b Builtin) ToDnum() dnum.Dnum {
	panic("can't convert builtin function to number")
}

func (b Builtin) String() string {
	return "builtin-function"
}

func (b Builtin) ToStr() string {
	return b.String()
}

func (Builtin) Get(Value) Value {
	panic("builtin function does not support get")
}

func (Builtin) Put(Value, Value) {
	panic("builtin function does not support put")
}

func (Builtin) RangeTo(int, int) Value {
	panic("builtin function does not support range")
}

func (Builtin) RangeLen(int, int) Value {
	panic("builtin function does not support range")
}

func (b Builtin) Hash() uint32 {
	panic("can't use builtin function as object key")
}

func (b Builtin) Hash2() uint32 {
	return b.Hash()
}

func (b Builtin) Equal(other interface{}) bool {
	if b2, ok := other.(Builtin); ok {
		return reflect.ValueOf(b).Pointer() == reflect.ValueOf(b2).Pointer()
	}
	return false
}

func (Builtin) TypeName() string {
	return "BuiltinFunction"
}

func (Builtin) Order() Ord {
	panic("can't compare builtin functions")
}

func (b Builtin) Compare(Value) int {
	panic("can't compare builtin functions")
}
