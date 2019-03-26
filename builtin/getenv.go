package builtin

import (
	"os"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Getenv(string)",
	func(arg Value) Value {
		return SuStr(os.Getenv(IfStr(arg)))
	})
