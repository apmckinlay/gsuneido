// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

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
	"golang.org/x/text/encoding/charmap"
)

func init() {
	StringMethods = Methods{
		"Alpha?": method0(func(this Value) Value {
			s := ToStr(this)
			if s == "" {
				return False
			}
			for _, c := range []byte(ToStr(this)) {
				if !ascii.IsLetter(c) {
					return False
				}
			}
			return True
		}),
		"AlphaNum?": method0(func(this Value) Value {
			s := ToStr(this)
			if s == "" {
				return False
			}
			for _, c := range []byte(ToStr(this)) {
				if !ascii.IsLetter(c) && !ascii.IsDigit(c) {
					return False
				}
			}
			return True
		}),
		"Asc": method0(func(this Value) Value {
			s := ToStr(this)
			if s == "" {
				return Zero
			}
			return SuInt(int(s[0]))
		}),
		"Compile": method("(errob = false)",
			func(t *Thread, this Value, args []Value) Value {
				if args[0] == False {
					return compile.Constant(ToStr(this))
				}
				ob := ToContainer(args[0])
				val, checks := compile.Checked(t, ToStr(this))
				for _, w := range checks {
					ob.Add(SuStr(w))
				}
				return val
			}),
		"Count": method1("(string)", func(this, arg Value) Value {
			return IntVal(strings.Count(ToStr(this), ToStr(arg)))
		}),
		"Detab": method0(func(this Value) Value {
			return SuStr(tabs.Detab(ToStr(this)))
		}),
		"Entab": method0(func(this Value) Value {
			return SuStr(tabs.Entab(ToStr(this)))
		}),
		"Eval": method("()", func(t *Thread, this Value, args []Value) Value {
			result := EvalString(t, ToStr(this))
			if result == nil {
				return EmptyStr
			}
			return result
		}),
		"Eval2": method("()", func(t *Thread, this Value, args []Value) Value {
			ob := &SuObject{}
			if result := EvalString(t, ToStr(this)); result != nil {
				ob.Add(result)
			}
			return ob
		}),
		"Extract": method2("(pattern, part=false)", func(this, arg1, arg2 Value) Value {
			s := ToStr(this)
			pat := regex.Compile(ToStr(arg1))
			var res regex.Result
			if pat.FirstMatch(s, 0, &res) == -1 {
				return False
			}
			var pos, end int
			if arg2 == False {
				pos, end = res[1].Range()
				if pos == -1 {
					pos, end = res[0].Range()
				}
			} else {
				part := ToInt(arg2)
				pos, end = res[part].Range()
			}
			if pos == -1 {
				return EmptyStr
			}
			return SuStr(s[pos:end])
		}),
		"Find": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := ToStr(this)
			pos := position(arg2, len(s))
			i := strings.Index(s[pos:], ToStr(arg1))
			if i == -1 {
				return IntVal(len(s))
			}
			return IntVal(pos + i)
		}),
		"Find1of": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := ToStr(this)
			pos := position(arg2, len(s))
			i := strings.IndexAny(s[pos:], ToStr(arg1))
			if i == -1 {
				return IntVal(len(s))
			}
			return IntVal(pos + i)
		}),
		"Findnot1of": method2("(string, pos=0)", func(this, arg1, arg2 Value) Value {
			s := ToStr(this)
			pos := position(arg2, len(s))
			i := str.IndexNotAny(s[pos:], ToStr(arg1))
			if i == -1 {
				return IntVal(len(s))
			}
			return IntVal(pos + i)
		}),
		"FindLast": method2("(string, pos=false)", func(this, arg1, arg2 Value) Value {
			s := ToStr(this)
			substr := ToStr(arg1)
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
			set := ToStr(arg1)
			if set == "" {
				return False
			}
			s := ToStr(this)
			end := last1ofEnd(s, arg2)
			if end < 0 {
				return False
			}
			return intOrFalse(strings.LastIndexAny(s[:end], set))
		}),
		"FindLastnot1of": method2("(string, pos=false)", func(this, arg1, arg2 Value) Value {
			s := ToStr(this)
			set := ToStr(arg1)
			end := last1ofEnd(s, arg2)
			if end < 0 || set == "" {
				return False
			}
			return intOrFalse(str.LastIndexNotAny(s[:end], set))
		}),
		"FromUtf8": method0(func(this Value) Value {
			utf8 := ToStr(this)
			encoder := charmap.Windows1252.NewEncoder()
			s, err := encoder.String(utf8)
			if err != nil {
				panic("string.FromUtf8 " + err.Error())
			}
			return SuStr(s)
		}),
		"Has?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.Contains(ToStr(this), ToStr(arg)))
		}),
		"Iter": method0(func(this Value) Value {
			iterable := this.(interface{ Iter() Iter })
			return SuIter{Iter: iterable.Iter()}
		}),
		"Lower": method0(func(this Value) Value {
			return SuStr(str.ToLower(ToStr(this)))
		}),
		"Lower?": method0(func(this Value) Value {
			result := false
			for _, c := range []byte(ToStr(this)) {
				if ascii.IsUpper(c) {
					return False
				} else if ascii.IsLower(c) {
					result = true
				}
			}
			return SuBool(result)
		}),
		"MapN": method("(n, block)",
			func(t *Thread, this Value, args []Value) Value {
				s := ToStr(this)
				n := IfInt(args[0])
				block := args[1]
				var buf strings.Builder
				for i := 0; i < len(s); i += n {
					end := ints.Min(i+n, len(s))
					val := t.Call(block, SuStr(s[i:end]))
					if val != nil {
						buf.WriteString(AsStr(val))
					}
				}
				return SuStr(buf.String())
			}),
		"Match": method3("(pattern, pos=false, prev=false)",
			func(this, arg1, arg2, arg3 Value) Value {
				s := ToStr(this)
				pat := regex.Compile(ToStr(arg1))
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
						ob.Set(SuInt(i), p)
					}
				}
				return ob
			}),
		"NthLine": method1("(n)", func(this, arg Value) Value {
			s := ToStr(this)
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
			return SuBool(numberPat.Matches(ToStr(this)))
		}),
		"Numeric?": method0(func(this Value) Value {
			s := ToStr(this)
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
			s := ToStr(this)
			pre := ToStr(arg1)
			pos := position(arg2, len(s))
			return SuBool(strings.HasPrefix(s[pos:], pre))
		}),
		"Repeat": method1("(count)", func(this, arg Value) Value {
			return SuStr(strings.Repeat(ToStr(this), ints.Max(0, ToInt(arg))))
		}),
		"Replace": method("(pattern, block = '', count = false)",
			func(t *Thread, this Value, args []Value) Value {
				count := ints.MaxInt
				if args[2] != False {
					count = ToInt(args[2])
				}
				return replace(t, ToStr(this), ToStr(args[0]), args[1], count)
			}),
		"Reverse": method0(func(this Value) Value {
			s := []byte(ToStr(this))
			lo := 0
			hi := len(s) - 1
			for lo < hi {
				s[lo], s[hi] = s[hi], s[lo]
				lo++
				hi--
			}
			return SuStr(string(s))
		}),
		"ServerEval": method("()",
			func(t *Thread, this Value, args []Value) Value {
				return t.Dbms().Run(ToStr(this))
			}),
		"Size": method0(func(this Value) Value {
			// avoid calling ToStr so we don't have to convert concats
			return IntVal(this.(interface{ Len() int }).Len())
			// "this" should always have Len
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
			return NewSuObject(vals)
		}),
		"Suffix?": method1("(string)", func(this, arg Value) Value {
			return SuBool(strings.HasSuffix(ToStr(this), ToStr(arg)))
		}),
		"ToUtf8": method0(func(this Value) Value {
			s := ToStr(this)
			decoder := charmap.Windows1252.NewDecoder()
			utf8, err := decoder.String(s)
			if err != nil {
				panic("string.ToUtf8 " + err.Error())
			}
			return SuStr(utf8)
		}),
		"Tr": method("(from, to='')",
			func(t *Thread, this Value, args []Value) Value {
				from := t.TrCache.Get(ToStr(args[0]))
				to := t.TrCache.Get(ToStr(args[1]))
				return SuStr(tr.Replace(ToStr(this), from, to))
			}),
		"Unescape": method0(func(this Value) Value {
			s := ToStr(this)
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
			return SuStr(str.ToUpper(ToStr(this)))
		}),
		"Upper?": method0(func(this Value) Value {
			result := false
			for _, c := range []byte(ToStr(this)) {
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

func replace(t *Thread, s string, patarg string, reparg Value, count int) Value {
	if count <= 0 || (patarg == "" && reparg == EmptyStr) {
		return SuStr(s)
	}
	pat := t.RxCache.Get(patarg)
	rep := ""
	if !isFunc(reparg) {
		rep = AsStr(reparg)
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
			v := t.Call(reparg, SuStr(r))
			if v != nil {
				r = AsStr(v)
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
