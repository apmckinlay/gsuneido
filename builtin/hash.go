package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Hash(value)",
	func(arg Value) Value {
		return IntVal(int(arg.Hash()))
	})
