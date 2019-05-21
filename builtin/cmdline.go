package builtin

import (
	"flag"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("Cmdline()", func() Value {
	var sb strings.Builder
	sep := ""
	for _, arg := range flag.Args() {
		sb.WriteString(sep)
		sep = " "
		if strings.ContainsAny(arg, " '\"") {
			arg = SuStr(arg).String()
		}
		sb.WriteString(arg)
	}
	return SuStr(sb.String())
})
