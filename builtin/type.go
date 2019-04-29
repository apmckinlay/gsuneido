package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

var _ = builtin1("Type(value)",
	func(arg Value) Value {
		return SuStr(arg.Type().String())
	})

var _ = builtin1("Boolean?(value)",
	func(arg Value) Value {
		return SuBool(arg.Type() == types.Boolean)
	})

var _ = builtin1("Number?(value)",
	func(arg Value) Value {
		return SuBool(arg.Type() == types.Number)
	})

var _ = builtin1("String?(value)",
	func(arg Value) Value {
		t := arg.Type()
		return SuBool(t == types.String || t == types.Except)
	})

var _ = builtin1("Date?(value)",
	func(arg Value) Value {
		return SuBool(arg.Type() == types.Date)
	})

var _ = builtin1("Object?(value)",
	func(arg Value) Value {
		return SuBool(arg.Type() == types.Object)
	})

var _ = builtin1("Record?(value)",
	func(arg Value) Value {
		return SuBool(arg.Type() == types.Record)
	})

var _ = builtin1("Class?(value)",
	func(arg Value) Value {
		return SuBool(arg.Type() == types.Class)
	})

var _ = builtin1("Instance?(value)",
	func(arg Value) Value {
		return SuBool(arg.Type() == types.Instance)
	})

var _ = builtin1("Function?(value)",
	func(arg Value) Value {
		return SuBool(isFunc(arg))
	})

func isFunc(v Value) bool {
	switch v.Type() {
	case types.Function, types.Block, types.Method, types.BuiltinFunction:
		return true
	}
	return false
}
