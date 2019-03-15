package builtin

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/regex"
)

// Eval executes Suneido code
func Eval(t *Thread, s string) Value {
	s = strings.Trim(s, " \t\r\n")
	if isGlobal(s) {
		return Global.Get(Global.Num(s))
	}
	s = "function () {\n" + s + "\n}"
	fn := compile.NamedConstant("eval", s).(*SuFunc)
	return t.Call(fn)
}

var rxGlobal = regex.Compile("^[A-Z][_a-zA-Z0-9]*?[!?]?$")

func isGlobal(s string) bool {
	return rxGlobal.Matches(s)
}
