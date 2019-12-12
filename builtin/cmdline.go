package builtin

import (
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Cmdline()", func() Value {
	return SuStr(options.CmdLine)
})
