package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin1("Locals(i)",
	func(arg Value) Value {
		return &SuObject{} //TODO
	})
