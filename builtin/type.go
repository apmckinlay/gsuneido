// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
)

var _ = builtin(Type, "(value)")

func Type(arg Value) Value {
	return SuStr(arg.Type().String())
}

var _ = builtin(BooleanQ, "(value)")

func BooleanQ(arg Value) Value {
	return SuBool(arg.Type() == types.Boolean)
}

var _ = builtin(NumberQ, "(value)")

func NumberQ(arg Value) Value {
	return SuBool(arg.Type() == types.Number)
}

var _ = builtin(StringQ, "(value)")

func StringQ(arg Value) Value {
	t := arg.Type()
	return SuBool(t == types.String || t == types.Except)
}

var _ = builtin(DateQ, "(value)")

func DateQ(arg Value) Value {
	return SuBool(arg.Type() == types.Date)
}

var _ = builtin(ObjectQ, "(value)")

func ObjectQ(arg Value) Value {
	switch arg.Type() {
	case types.Object, types.Record:
		return True
	}
	return False
}

var _ = builtin(RecordQ, "(value)")

func RecordQ(arg Value) Value {
	return SuBool(arg.Type() == types.Record)
}

var _ = builtin(ClassQ, "(value)")

func ClassQ(arg Value) Value {
	return SuBool(arg.Type() == types.Class)
}

var _ = builtin(InstanceQ, "(value)")

func InstanceQ(arg Value) Value {
	return SuBool(arg.Type() == types.Instance)
}

var _ = builtin(FunctionQ, "(value)")

func FunctionQ(arg Value) Value {
	return SuBool(isFunc(arg))
}

func isFunc(v Value) bool {
	switch v.Type() {
	case types.Function, types.Block, types.Method, types.BuiltinFunction:
		return true
	}
	return false
}
