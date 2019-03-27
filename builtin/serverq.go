package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Server?()",
	func() Value {
		return False
	})
