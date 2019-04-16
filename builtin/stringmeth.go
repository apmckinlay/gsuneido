package builtin

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tabs"

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
			return IntVal(strings.Count(IfStr(this), IfStr(arg)))
		}),
		"Detab": method0(func(this Value) Value {
			return SuStr(tabs.Detab(IfStr(this)))
		}),
		"Entab": method0(func(this Value) Value {
			return SuStr(tabs.Entab(IfStr(this)))
		}),
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
		"Extract": method2("(pattern, part=0)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			pat := regex.Compile(IfStr(arg1))
			var res regex.Result
			if pat.FirstMatch(s, 0, &res) != -1 {
				return False
			}
			pos, end := res[1].Range()
			if pos == -1 {
				pos, end = res[0].Range()
			}
			return SuStr(s[pos:end])
		}),
		"Find": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			pos := position(arg2, len(s))
			i := strings.Index(s[pos:], IfStr(arg1))
			if i == -1 {
				return IntVal(len(s))
			}
			return IntVal(pos + i)
		}),
		"Find1of": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			pos := position(arg2, len(s))
			i := strings.IndexAny(s[pos:], IfStr(arg1))
			if i == -1 {
				return IntVal(len(s))
			}
			return IntVal(pos + i)
		}),
		"Findnot1of": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			pos := position(arg2, len(s))
			i := str.IndexNotAny(s[pos:], IfStr(arg1))
			if i == -1 {
				return IntVal(len(s))
			}
			return IntVal(pos + i)
		}),
		"FindLast": method2("(string, pos=false)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			substr := IfStr(arg1)
			end := len(s)
			if arg2 != False {
				end = ToInt(arg2) + len(substr)
				if end > len(s) {
					end = len(s)
				}
			}
			if end < 0 {
				return False
			}
			if substr == "" {
				return IntVal(end)
			}
			return intOrFalse(strings.LastIndex(s[:end], substr))
		}),
		"FindLast1of": method2("(string, pos=false)", func(this, arg1, arg2 Value) Value {
			set := IfStr(arg1)
			if set == "" {
				return False
			}
			s := IfStr(this)
			end := last1ofEnd(s, arg2)
			if end < 0 {
				return False
			}
			return intOrFalse(strings.LastIndexAny(s[:end], set))
		}),
		"FindLastnot1of": method2("(string, pos=false)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			set := IfStr(arg1)
			end := last1ofEnd(s, arg2)
			if end < 0 || set == "" {
				return False
			}
			return intOrFalse(str.LastIndexNotAny(s[:end], set))
		}),
		"Has?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.Contains(IfStr(this), IfStr(arg)))
		}),
		"Iter": method0(func(this Value) Value {
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
		"MapN": methodRaw("(n, default)", // methodRaw to get thread
			func(t *Thread, as *ArgSpec, this Value, args ...Value) Value {
				args = t.Args(&paramSpecMapN, as)
				s := IfStr(this)
				n := IfInt(args[0])
				block := args[1]
				var buf strings.Builder
				for i := 0; i < len(s); i += n {
					end := ints.Min(i+n, len(s))
					val := t.CallWithArgs(block, SuStr(s[i:end]))
					if val != nil {
						buf.WriteString(ToStr(val))
					}
				}
				return SuStr(buf.String())
			}),
		"Match": method3("(pattern, pos=false, prev=false)",
			func(this, arg1, arg2, arg3 Value) Value {
				s := IfStr(this)
				pat := regex.Compile(IfStr(arg1))
				prev := ToBool(arg3)
				pos := 0
				if arg2 != False {
					pos = IfInt(arg2)
				} else if prev {
					pos = len(s)
				}
				method := pat.FirstMatch
				if prev {
					method = pat.LastMatch
				}
				var res regex.Result
				if method(s, pos, &res) == -1 {
					return False
				}
				ob := &SuObject{}
				for i, part := range res {
					pos, end := part.Range()
					if pos >= 0 {
						p := &SuObject{}
						p.Add(IntVal(pos))
						p.Add(IntVal(end - pos))
						ob.Put(SuInt(i), p)
					}
				}
				return ob
			}),
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
		"Number?": method0(func(this Value) Value {
			return SuBool(numberPat.Matches(IfStr(this)))
		}),
		"Numeric?": method0(func(this Value) Value {
			s := IfStr(this)
			if len(s) == 0 {
				return False
			}
			for i := 0; i < len(s); i++ {
				if !ascii.IsDigit(s[i]) {
					return False
				}
			}
			return True
		}),
		"Prefix?": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := IfStr(this)
			pre := IfStr(arg1)
			pos := position(arg2, len(s))
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
			// avoid calling IfStr so we don't have to convert concats
			return IntVal(this.(interface{ Len() int }).Len())
			// "this" should always have Len
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
		"Unescape": method0(func(this Value) Value {
			s := IfStr(this)
			var buf strings.Builder
			buf.Grow(len(s))
			for i := 0; i < len(s); i++ {
				var c byte
				c, i = str.Doesc(s, i)
				buf.WriteByte(c)
			}
			return SuStr(buf.String())
		}),
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

var paramSpecMapN = params("(n, block)")
var paramSpecReplace = params("(pattern, replacement = '', count = false)")
var paramSpecTr = params("(from, to='')")

func replace(t *Thread, s string, patarg string, reparg Value, count int) Value {
	if count <= 0 {
		return SuStr(s)
	}
	pat := t.RxCache.Get(patarg)
	rep := ""
	if !isFunc(reparg) {
		rep = ToStr(reparg)
		reparg = nil
	}
	from := 0
	nsubs := 0
	var buf strings.Builder
	pat.ForEachMatch(s, func(result *regex.Result) bool {
		pos, end := result[0].Range()
		buf.WriteString(s[from:pos])
		if reparg == nil {
			t := regex.Replace(s, rep, result)
			buf.WriteString(t)
		} else {
			r := result[0].Part(s)
			v := t.CallWithArgs(reparg, SuStr(r))
			if v != nil {
				r = ToStr(v)
			}
			buf.WriteString(r)
		}
		from = end
		nsubs++
		return nsubs < count
	})
	buf.WriteString(s[from:])
	return SuStr(buf.String())
}

func position(arg Value, n int) int {
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

var numberPat = regex.Compile(`\A[+-]?(\d+\.?|\.\d)\d*?([eE][+-]?\d\d?)?\Z`)

func last1ofEnd(s string, arg2 Value) int {
	end := len(s)
	if arg2 != False {
		end = ToInt(arg2) + 1
		if end > len(s) {
			end = len(s)
		}
	}
	return end
}
func intOrFalse(i int) Value {
	if i == -1 {
		return False
	}
	return IntVal(i)
}
