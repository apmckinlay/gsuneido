package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	StringMethods = Methods{
		"Asc": method0(func(this Value) Value {
			return SuInt(int(ToStr(this)[0]))
		}),
		"Find": method1("(string)", func(this, arg Value) Value {
			s := ToStr(this)
			i := strings.Index(s, ToStr(arg))
			if i == -1 {
				i = len(s)
			}
			return IntToValue(i)
		}),
		"Has?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.Contains(ToStr(this), ToStr(arg)))
		}),
		"Iter": method0(func(this Value) Value { // TODO sequence
			iterable := this.(interface { Iter() Iter })
			return SuIter{Iter: iterable.Iter()}
		}),
		"Lower": method0(func(this Value) Value {
			return SuStr(strings.ToLower(ToStr(this)))
		}),
		"Prefix?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.HasPrefix(ToStr(this), ToStr(arg)))
		}),
		"Suffix?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.HasSuffix(ToStr(this), ToStr(arg)))
		}),
		"Repeat": method1("(count)", func(this, arg Value) Value {
			return SuStr(strings.Repeat(ToStr(this), ToInt(arg)))
		}),
		"Size": method0(func(this Value) Value {
			// TODO handle Concat without converting
			return IntToValue(len(ToStr(this)))
		}),
		"Upper": method0(func(this Value) Value {
			return SuStr(strings.ToUpper(ToStr(this)))
		}),
		// TODO more methods
	}
}
