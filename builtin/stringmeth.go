package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	StringMethods = Methods{
		"Size": method0(func(self Value) Value {
			return SuInt(len(self.ToStr()))
		}),
		"Lower": method0(func(self Value) Value {
			return SuStr(strings.ToLower(self.ToStr()))
		}),
		"Upper": method0(func(self Value) Value {
			return SuStr(strings.ToUpper(self.ToStr()))
		}),
		"Has?": method1("(string)", func(self, arg Value) Value {
			return SuBool(strings.Contains(self.ToStr(), arg.ToStr()))
		}),
		"Repeat": method1("(count)", func(self, arg Value) Value {
			return SuStr(strings.Repeat(self.ToStr(), arg.ToInt()))
		}),
		// TODO more methods
	}
}
