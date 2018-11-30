package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin1("Name(value)",
	func(arg Value) Value {
		if named,ok := arg.(Named); ok {
			return SuStr(named.GetName())
		}
		return EmptyStr
	})
