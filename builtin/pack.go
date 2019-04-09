package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Pack(value)",
	func(arg Value) Value {
		if p,ok := arg.(Packable); ok {
			return SuStr(Pack(p))
		}
		panic("can't pack " + arg.Type().String())
	})

var _ = builtin1("Unpack(string)",
	func(arg Value) Value {
		return Unpack(IfStr(arg))
	})
