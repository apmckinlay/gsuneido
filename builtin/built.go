package builtin

import (
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Built()", func() Value {
	return SuStr("Built: " + options.BuiltDate + " (Go)")
})
