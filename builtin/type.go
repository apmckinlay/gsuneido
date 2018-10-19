package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin1("Type(value)",
	func(arg Value) Value {
		return SuStr(arg.TypeName())
	})
