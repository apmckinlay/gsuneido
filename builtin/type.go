package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Type(value)",
	func(_ *Thread, args ...Value) Value {
		return SuStr(args[0].TypeName())
	})
