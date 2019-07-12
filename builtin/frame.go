package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin1("Frame(i)",
	func(arg Value) Value {
		return False //TODO
	})
