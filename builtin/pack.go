package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Pack(value)",
	func(arg Value) Value {
		return SuStr(PackValue(arg))
	})

var _ = builtin1("Unpack(string)",
	func(arg Value) Value {
		return Unpack(IfStr(arg))
	})
