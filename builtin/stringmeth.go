package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	StringMethods = Methods{
		method("Size()", func(t *Thread, self Value, args ...Value) Value {
			return SuInt(len(self.ToStr()))
		}),
		method("Lower()", func(t *Thread, self Value, args ...Value) Value {
			return SuStr(strings.ToLower(self.ToStr()))
		}),
		method("Upper()", func(t *Thread, self Value, args ...Value) Value {
			return SuStr(strings.ToUpper(self.ToStr()))
		}),
		method("Has?(string)", func(t *Thread, self Value, args ...Value) Value {
			return SuBool(strings.Contains(self.ToStr(), args[0].ToStr()))
		}),
		// TODO
	}
}
