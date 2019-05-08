package builtin

import (
	"os"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Exit(code)",
	func(arg Value) Value {
		os.Exit(IfInt(arg))
		return nil
	})
