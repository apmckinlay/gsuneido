package builtin

import (
	"strings"

	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var CmdlineOverride string

var _ = builtin0("Cmdline()", func() Value {
	if CmdlineOverride != "" {
		return SuStr(CmdlineOverride)
	}
	var sb strings.Builder
	sep := ""
	for _, arg := range options.Args {
		sb.WriteString(sep)
		sep = " "
		if strings.ContainsAny(arg, " '\"") {
			arg = SuStr(arg).String()
		}
		sb.WriteString(arg)
	}
	return SuStr(sb.String())
})
