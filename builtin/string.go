// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math"
	"strings"

	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/regex"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/tabs"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/tr"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

var _ = exportMethods(&StringMethods)

var _ = method(string_AlphaQ, "()")

func string_AlphaQ(this Value) Value {
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
}

var _ = method(string_AlphaNumQ, "()")

func string_AlphaNumQ(this Value) Value {
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
}

var _ = method(string_Asc, "()")

func string_Asc(this Value) Value {
	s := ToStr(this)
	if s == "" {
		return Zero
	}
	return SuInt(int(s[0]))
}

var _ = method(string_Compile, "(errob = false)")

func string_Compile(t *Thread, this Value, args []Value) Value {
	if args[0] == False {
		return compile.Constant(ToStr(this))
	}
	ob := ToContainer(args[0])
	val, checks := compile.Checked(t, ToStr(this))
	for _, w := range checks {
		ob.Add(SuStr(w))
	}
	return val
}

var _ = method(string_Count, "(string)")

func string_Count(this, arg Value) Value {
	return IntVal(strings.Count(ToStr(this), ToStr(arg)))
}

var _ = method(string_Detab, "()")

func string_Detab(this Value) Value {
	return SuStr(tabs.Detab(ToStr(this)))
}

var _ = method(string_Entab, "()")

func string_Entab(this Value) Value {
	return SuStr(tabs.Entab(ToStr(this)))
}

var _ = method(string_Eval, "()")

func string_Eval(t *Thread, this Value, args []Value) Value {
	result := compile.EvalString(t, ToStr(this))
	if result == nil {
		return EmptyStr
	}
	return result
}

var _ = method(string_Eval2, "()")

func string_Eval2(t *Thread, this Value, args []Value) Value {
	ob := &SuObject{}
	if result := compile.EvalString(t, ToStr(this)); result != nil {
		ob.Add(result)
	}
	return ob
}

var _ = method(string_Extract, "(pattern, part=false)")

func string_Extract(this, arg1, arg2 Value) Value {
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
}

var _ = method(string_Find, "(string, pos=0)")

func string_Find(this, arg1, arg2 Value) Value {
	s := ToStr(this)
	pos := position(arg2, len(s))
	i := strings.Index(s[pos:], ToStr(arg1))
	if i == -1 {
		return IntVal(len(s))
	}
	return IntVal(pos + i)
}

var _ = method(string_Find1of, "(string, pos=0)")

func string_Find1of(this, arg1, arg2 Value) Value {
	s := ToStr(this)
	pos := position(arg2, len(s))
	i := strings.IndexAny(s[pos:], ToStr(arg1))
	if i == -1 {
		return IntVal(len(s))
	}
	return IntVal(pos + i)
}

var _ = method(string_Findnot1of, "(string, pos=0)")

func string_Findnot1of(this, arg1, arg2 Value) Value {
	s := ToStr(this)
	pos := position(arg2, len(s))
	i := str.IndexNotAny(s[pos:], ToStr(arg1))
	if i == -1 {
		return IntVal(len(s))
	}
	return IntVal(pos + i)
}

var _ = method(string_FindLast, "(string, pos=false)")

func string_FindLast(this, arg1, arg2 Value) Value {
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
}

var _ = method(string_FindLast1of, "(string, pos=false)")

func string_FindLast1of(this, arg1, arg2 Value) Value {
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
}

var _ = method(string_FindLastnot1of, "(string, pos=false)")

func string_FindLastnot1of(this, arg1, arg2 Value) Value {
	s := ToStr(this)
	set := ToStr(arg1)
	end := last1ofEnd(s, arg2)
	if end < 0 || set == "" {
		return False
	}
	return intOrFalse(str.LastIndexNotAny(s[:end], set))
}

var _ = method(string_FromUtf8, "()")

func string_FromUtf8(this Value) Value {
	utf8 := ToStr(this)
	encoder := charmap.Windows1252.NewEncoder()
	encoder = encoding.ReplaceUnsupported(encoder)
	s, err := encoder.String(utf8)
	if err != nil {
		panic("string.FromUtf8 " + err.Error())
	}
	return SuStr(s)
}

var _ = method(string_HasQ, "(string)")

func string_HasQ(this, arg Value) Value {
	return SuBool(strings.Contains(ToStr(this), ToStr(arg)))
}

var _ = method(string_Iter, "()")

func string_Iter(this Value) Value {
	iterable := this.(interface{ Iter() Iter })
	return SuIter{Iter: iterable.Iter()}
}

var _ = method(string_Lower, "()")

func string_Lower(this Value) Value {
	return SuStr(str.ToLower(ToStr(this)))
}

var _ = method(string_LowerQ, "()")

func string_LowerQ(this Value) Value {
	result := false
	for _, c := range []byte(ToStr(this)) {
		if ascii.IsUpper(c) {
			return False
		} else if ascii.IsLower(c) {
			result = true
		}
	}
	return SuBool(result)
}

var _ = method(string_MapN, "(n, block)")

func string_MapN(t *Thread, this Value, args []Value) Value {
	s := ToStr(this)
	n := IfInt(args[0])
	block := args[1]
	var buf strings.Builder
	for i := 0; i < len(s); i += n {
		end := ord.Min(i+n, len(s))
		val := t.Call(block, SuStr(s[i:end]))
		if val != nil {
			buf.WriteString(AsStr(val))
		}
	}
	return SuStr(buf.String())
}

var _ = method(string_Match, "(pattern, pos=false, prev=false)")

func string_Match(this, arg1, arg2, arg3 Value) Value {
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
}

var _ = method(string_NthLine, "(n)")

func string_NthLine(this, arg Value) Value {
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
}

var _ = method(string_NumberQ, "()")

func string_NumberQ(this Value) Value {
	return SuBool(numberPat.Matches(ToStr(this)))
}

var numberPat = regex.Compile(`\A[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d\d?)?\Z`)
// benchmark shows Suneido regex is faster than Go regexp for this

var _ = method(string_NumericQ, "()")

func string_NumericQ(this Value) Value {
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
}

var _ = method(string_PrefixQ, "(string, pos=0)")

func string_PrefixQ(this, arg1, arg2 Value) Value {
	s := ToStr(this)
	pre := ToStr(arg1)
	pos := position(arg2, len(s))
	return SuBool(strings.HasPrefix(s[pos:], pre))
}

var _ = method(string_Repeat, "(count)")

func string_Repeat(this, arg Value) Value {
	return SuStr(strings.Repeat(ToStr(this), ord.Max(0, ToInt(arg))))
}

var _ = method(string_Replace, "(pattern, block = '', count = false)")

func string_Replace(t *Thread, this Value, args []Value) Value {
	count := math.MaxInt
	if args[2] != False {
		count = ToInt(args[2])
	}
	return replace(t, ToStr(this), ToStr(args[0]), args[1], count)
}

var _ = method(string_Reverse, "()")

func string_Reverse(this Value) Value {
	s := []byte(ToStr(this))
	lo := 0
	hi := len(s) - 1
	for lo < hi {
		s[lo], s[hi] = s[hi], s[lo]
		lo++
		hi--
	}
	return SuStr(string(s))
}

var _ = method(string_ServerEval, "()")

func string_ServerEval(t *Thread, this Value, args []Value) Value {
	return t.Dbms().Run(t, ToStr(this))
}

var _ = method(string_Size, "()")

func string_Size(this Value) Value {
	// avoid calling ToStr so we don't have to convert concats
	return IntVal(this.(interface{ Len() int }).Len())
	// "this" should always have Len
}

var _ = method(string_Split, "(separator)")

func string_Split(this, arg Value) Value {
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
}

var _ = method(string_SuffixQ, "(string)")

func string_SuffixQ(this, arg Value) Value {
	return SuBool(strings.HasSuffix(ToStr(this), ToStr(arg)))
}

var _ = method(string_ToUtf8, "()")

func string_ToUtf8(this Value) Value {
	s := ToStr(this)
	decoder := charmap.Windows1252.NewDecoder()
	utf8, err := decoder.String(s)
	if err != nil {
		panic("string.ToUtf8 " + err.Error())
	}
	return SuStr(utf8)
}

var _ = method(string_Tr, "(from, to='')")

func string_Tr(t *Thread, this Value, args []Value) Value {
	from := t.TrCache.Get(ToStr(args[0]))
	to := t.TrCache.Get(ToStr(args[1]))
	return SuStr(tr.Replace(ToStr(this), from, to))
}

var _ = method(string_Unescape, "()")

func string_Unescape(this Value) Value {
	s := ToStr(this)
	var buf strings.Builder
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ {
		var c byte
		c, i = str.Doesc(s, i)
		buf.WriteByte(c)
	}
	return SuStr(buf.String())
}

var _ = method(string_Upper, "()")

func string_Upper(this Value) Value {
	return SuStr(str.ToUpper(ToStr(this)))
}

var _ = method(string_UpperQ, "()")

func string_UpperQ(this Value) Value {
	result := false
	for _, c := range []byte(ToStr(this)) {
		if ascii.IsLower(c) {
			return False
		} else if ascii.IsUpper(c) {
			result = true
		}
	}
	return SuBool(result)
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
		// use Go strings.Replace if literal
		if p, ok := pat.Literal(); ok {
			if r, ok := regex.LiteralRep(rep); ok {
				return SuStr(strings.Replace(s, p, r, count))
			}
		}
	}

	from := 0
	nreps := 0
	var buf strings.Builder
	pat.ForEachMatch(s, func(result *regex.Result) bool {
		pos, end := result[0].Range()
		buf.WriteString(s[from:pos])
		if reparg == nil {
			t := regex.Replacement(s, rep, result)
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
		nreps++
		return nreps < count
	})
	if nreps == 0 {
		// avoid copy if no replacements
		return SuStr(s)
	}
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
