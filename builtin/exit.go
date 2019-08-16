package builtin

import (
	"os"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Exit(code = 0)",
	func(arg Value) Value {
		os.Exit(IfInt(arg))
		return nil
	})
