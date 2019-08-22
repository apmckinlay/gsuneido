package builtin

import (
	"runtime"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Built()", func() Value {
	if options.BuiltDate == "" {
		options.BuiltDate = "Aug 19 2019 14:31:46"
	}
	return SuStr(options.BuiltDate + " (" + runtime.Version() + ")")
})
