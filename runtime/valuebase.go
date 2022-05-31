// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"reflect"
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

// ValueBase is embedded in Value types to supply default methods.
// The parameter type is only used to get a type name for errors.
type ValueBase[E any] struct{}

func (ValueBase[E]) TypeName() string {
	return typeName[E]()
}

func typeName[E any]() string {
	var z E
	t := reflect.TypeOf(z)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	s := t.Name()
	if strings.HasPrefix(s, "Su") || strings.HasPrefix(s, "su") {
		s = s[2:]
	}
	return s
}

func (ValueBase[E]) String() string {
	return str.UnCapitalize(typeName[E]())
}

func (ValueBase[E]) Type() types.Type {
	return types.BuiltinClass
}

func (ValueBase[E]) ToInt() (int, bool) {
	return 0, false
}

func (ValueBase[E]) IfInt() (int, bool) {
	return 0, false
}

func (ValueBase[E]) ToDnum() (dnum.Dnum, bool) {
	return dnum.Zero, false
}

func (ValueBase[E]) ToContainer() (Container, bool) {
	return nil, false
}

func (ValueBase[E]) AsStr() (string, bool) {
	return "", false
}

func (ValueBase[E]) ToStr() (string, bool) {
	return "", false
}

func (ValueBase[E]) Lookup(*Thread, string) Callable {
	return nil // no methods
}

func (ValueBase[E]) Get(*Thread, Value) Value {
	panic(typeName[E]() + " does not support get")
}

func (ValueBase[E]) Put(*Thread, Value, Value) {
	panic(typeName[E]() + " does not support put")
}

func (ValueBase[E]) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic(typeName[E]() + " does not support update")
}

func (ValueBase[E]) RangeTo(int, int) Value {
	panic(typeName[E]() + " does not support range")
}

func (ValueBase[E]) RangeLen(int, int) Value {
	panic(typeName[E]() + " does not support range")
}

func (ValueBase[E]) Hash() uint32 {
	panic(typeName[E]() + " hash not implemented")
}

func (ValueBase[E]) Hash2() uint32 {
	panic(typeName[E]() + " hash not implemented")
}

func (ValueBase[E]) Compare(Value) int {
	panic(typeName[E]() + " compare not implemented")
}

func (ValueBase[E]) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call " + typeName[E]())
}

func (ValueBase[E]) SetConcurrent() {
	panic(typeName[E]() + " can not be shared between threads")
}

func (ValueBase[E]) Equal(other any) bool {
	panic(typeName[E]() + " equal not implemented")
}
