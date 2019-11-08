package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin3("Traccel(pointer, message, wParam)",
	func(a, b, c Value) Value {
		return IntVal(goc.Traccel(ToInt(a), ToInt(b), ToInt(c)))
	})
