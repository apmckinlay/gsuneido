package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Object(@args)",
	func(_ *Thread, args ...Value) Value {
		return args[0]
	})
