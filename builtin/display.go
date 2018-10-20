package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin1("Display(value)",
	func(arg Value) Value {
		return SuStr(arg.String())
	})
