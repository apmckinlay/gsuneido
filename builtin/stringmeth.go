package builtin

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/regex"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/tr"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	StringMethods = Methods{
		"Asc": method0(func(this Value) Value {
			return SuInt(int(IfStr(this)[0]))
		}),
		"Compile": method1("(errob = false)", func(this, _ Value) Value {
			return compile.Constant(IfStr(this))
		}),
		"Count": method1("(string)", func(this, arg Value) Value {
			return IntToValue(strings.Count(IfStr(this), IfStr(arg)))
		}),
		//TODO Detab
		//TODO Entab
		"Eval": methodRaw("(string)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				t.Args(&ParamSpec0, as)
				result := EvalString(t, IfStr(this))
				if result == nil {
					return EmptyStr
				}
				return result
			}),
		"Eval2": methodRaw("(string)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				t.Args(&ParamSpec0, as)
				ob := &SuObject{}
				if result := EvalString(t, IfStr(this)); result != nil {
					ob.Add(result)
				}
				return ob
			}),
		"Find": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			pos := pos(arg2, len(s))
			i := strings.Index(s[pos:], IfStr(arg1))
			if i == -1 {
				i = len(s)
			}
			return IntToValue(i)
		}),
		//TODO FindLast
		//TODO Find1of
		//TODO FindLast1of
		//TODO Findnot1of
		//TODO FindLastnot1of
		"Has?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.Contains(IfStr(this), IfStr(arg)))
		}),
		"Iter": method0(func(this Value) Value { // TODO sequence
			iterable := this.(interface{ Iter() Iter })
			return SuIter{Iter: iterable.Iter()}
		}),
		"Lower": method0(func(this Value) Value {
			return SuStr(strings.ToLower(IfStr(this)))
		}),
		"Lower?": method0(func(this Value) Value {
			result := false
			for _, c := range []byte(IfStr(this)) {
				if ascii.IsUpper(c) {
					return False
				} else if ascii.IsLower(c) {
					result = true
				}
			}
			return SuBool(result)
		}),
		//TODO MapN
		//TODO Match
		"NthLine": method1("(n)", func(this, arg Value) Value {
			s := IfStr(this)
			n := len(s)
			nth := ToInt(arg)
			i := 0
			for ; i < n && nth > 0; i++ {
				if s[i] == '\n' {
					nth--
				}
			}
			j := i
			for j < n && s[j] != '\n' {
				j++
			}
			for j > i && s[j-1] == '\r' {
				j--
			}
			return SuStr(s[i:j])
		}),
		//TODO Number?
		//TODO Numeric
		"Prefix?": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			pre := IfStr(arg1)
			pos := pos(arg2, len(s))
			return SuBool(strings.HasPrefix(s[pos:], pre))
		}),
		"Repeat": method1("(count)", func(this, arg Value) Value {
			return SuStr(strings.Repeat(IfStr(this), ints.Max(0, ToInt(arg))))
		}),
		"Replace": methodRaw("(string)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&paramSpecReplace, as)
				count := ints.MaxInt
				if args[2] != False {
					count = ToInt(args[2])
				}
				return replace(t, IfStr(this), IfStr(args[0]), args[1], count)
			}),
		"Reverse": method0(func(this Value) Value {
			s := []byte(IfStr(this))
			lo := 0
			hi := len(s) - 1
			for lo < hi {
				s[lo], s[hi] = s[hi], s[lo]
				lo++
				hi--
			}
			return SuStr(string(s))
		}),
		"ServerEval": methodRaw("(string)",
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				result := EvalString(t, IfStr(this))
				if result == nil {
					return EmptyStr
				}
				return result
			}),
		"Size": method0(func(this Value) Value {
			// TODO handle Concat without converting
			return IntToValue(len(IfStr(this)))
		}),
		"Sort!": methodRaw("(block = false)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&ParamSpecOptionalBlock, as)
				ToObject(this).Sort(t, args[0])
				return this
			}),
		"Split": method1("(separator)", func(this, arg Value) Value {
			sep := IfStr(arg)
			if sep == "" {
				panic("string.Split separator must not be empty string")
			}
			strs := strings.Split(IfStr(this), sep)
			if strs[len(strs)-1] == "" {
				strs = strs[:len(strs)-1]
			}
			vals := make([]Value, len(strs))
			for i, str := range strs {
				vals[i] = SuStr(str)
			}
			return NewSuObject(vals...)
		}),
		"Substr": method2("(i, n=false)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			sn := len(s)
			i := Index(arg1)
			if i < 0 {
				i += sn
				if i < 0 {
					i = 0
				}
			}
			n := sn - i
			if arg2 != False {
				n = ToInt(arg2)
				if n < 0 {
					n += sn - i
					if n < 0 {
						n = 0
					}
				}
			}
			return SuStr(s[i : i+n])
		}),
		"Suffix?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.HasSuffix(IfStr(this), IfStr(arg)))
		}),
		"Tr": methodRaw("(from, to='')", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&paramSpecTr, as)
				from := t.TrCache.Get(IfStr(args[0]))
				to := t.TrCache.Get(IfStr(args[1]))
				return SuStr(tr.Replace(IfStr(this), from, to))
			}),
		//TODO Unescape
		"Upper": method0(func(this Value) Value {
			return SuStr(strings.ToUpper(IfStr(this)))
		}),
		"Upper?": method0(func(this Value) Value {
			result := false
			for _, c := range []byte(IfStr(this)) {
				if ascii.IsLower(c) {
					return False
				} else if ascii.IsUpper(c) {
					result = true
				}
			}
			return SuBool(result)
		}),
	}
}

var paramSpecReplace = params("(pattern, replacement = '', count = false)")
var paramSpecTr = params("(from, to='')")

func replace(t *Thread, s string, patarg string, reparg Value, count int) Value {
	if count <= 0 {
		return SuStr(s)
	}
	pat := t.RxCache.Get(patarg)
	rep := IfStr(reparg) //TODO block && backrefs
	from := 0
	nsubs := 0
	var buf strings.Builder
	pat.ForEachMatch(s, func(result *regex.Result) bool {
		buf.WriteString(s[from:result.Pos()])
		buf.WriteString(rep)
		from = result.End()
		nsubs++
		return nsubs < count
	})
	buf.WriteString(s[from:])
	return SuStr(buf.String())
}

func pos(arg Value, n int) int {
	pos := ToInt(arg)
	if pos >= n {
		return n
	}
	if pos < 0 {
		pos += n
		if pos < 0 {
			pos = 0
		}
	}
	return pos
}
