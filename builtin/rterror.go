package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin0("RuntimeError()",
	func() Value {
		var x []Value
		return x[123]
	})
