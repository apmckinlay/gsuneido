package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	StringMethods = Methods{
		"Size": method0(func(this Value) Value {
			return SuInt(len(this.ToStr()))
		}),
		"Lower": method0(func(this Value) Value {
			return SuStr(strings.ToLower(this.ToStr()))
		}),
		"Upper": method0(func(this Value) Value {
			return SuStr(strings.ToUpper(this.ToStr()))
		}),
		"Has?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.Contains(this.ToStr(), arg.ToStr()))
		}),
		"Repeat": method1("(count)", func(this, arg Value) Value {
			return SuStr(strings.Repeat(this.ToStr(), arg.ToInt()))
		}),
		// TODO more methods
	}
}
