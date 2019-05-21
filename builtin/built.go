package builtin

import (
	"runtime"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Built()", func() Value {
	return SuStr(options.BuiltDate + " (" + runtime.Version() + ")")
})
