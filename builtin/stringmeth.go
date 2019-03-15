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
		"Sort!": methodRaw("(block = false)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&ParamSpecOptionalBlock, as)
				ToObject(this).Sort(t, args[0])
				return this
			}),
		"Eval": methodRaw("(string)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				result := Eval(t, ToStr(this))
				if result == nil {
					return EmptyStr
				}
				return result
			}),
		"Eval2": methodRaw("(string)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				ob := &SuObject{}
				if result := Eval(t, ToStr(this)); result != nil {
					ob.Add(result)
				}
				return ob
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
			iterable := this.(interface{ Iter() Iter })
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
		"Split": method1("(separator)", func(this, arg Value) Value {
			sep := ToStr(arg)
			if sep == "" {
				panic("string.Split separator must not be empty string")
			}
			strs := strings.Split(ToStr(this), sep)
			if strs[len(strs)-1] == "" {
				strs = strs[:len(strs)-1]
			}
			vals := make([]Value, len(strs))
			for i, str := range strs {
				vals[i] = SuStr(str)
			}
			return NewSuObject(vals...)
		}),
		"Upper": method0(func(this Value) Value {
			return SuStr(strings.ToUpper(ToStr(this)))
		}),
		// TODO more methods
	}
}
