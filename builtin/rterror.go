package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

// cause a Go runtime error (for testing)
var _ = builtin0("RuntimeError()",
	func() Value {
		var x []Value
		return x[123]
	})
