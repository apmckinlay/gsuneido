package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	StringMethods = Methods{
		"Asc": method0(func(this Value) Value {
			return SuInt(int(this.ToStr()[0]))
		}),
		"Find": method1("(string)", func(this, arg Value) Value {
			s := this.ToStr()
			i := strings.Index(s, arg.ToStr())
			if i == -1 {
				i = len(s)
			}
			return IntToValue(i)
		}),
		"Has?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.Contains(this.ToStr(), arg.ToStr()))
		}),
		"Lower": method0(func(this Value) Value {
			return SuStr(strings.ToLower(this.ToStr()))
		}),
		"Prefix?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.HasPrefix(this.ToStr(), arg.ToStr()))
		}),
		"Suffix?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.HasSuffix(this.ToStr(), arg.ToStr()))
		}),
		"Repeat": method1("(count)", func(this, arg Value) Value {
			return SuStr(strings.Repeat(this.ToStr(), arg.ToInt()))
		}),
		"Size": method0(func(this Value) Value {
			return IntToValue(len(this.ToStr()))
		}),
		"Upper": method0(func(this Value) Value {
			return SuStr(strings.ToUpper(this.ToStr()))
		}),
		// TODO more methods
	}
}
