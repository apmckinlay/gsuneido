package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	StringMethods = Methods{
		"Size": method("()", func(t *Thread, self Value, args ...Value) Value {
			return SuInt(len(self.ToStr()))
		}),
		"Lower": method("()", func(t *Thread, self Value, args ...Value) Value {
			return SuStr(strings.ToLower(self.ToStr()))
		}),
		"Upper": method("()", func(t *Thread, self Value, args ...Value) Value {
			return SuStr(strings.ToUpper(self.ToStr()))
		}),
		"Has?": method("(string)", func(t *Thread, self Value, args ...Value) Value {
			return SuBool(strings.Contains(self.ToStr(), args[0].ToStr()))
		}),
		// TODO
	}
}
