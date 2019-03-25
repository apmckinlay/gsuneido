package builtin

import (
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Client?()",
	func() Value {
		return SuBool(options.Client)
	})
