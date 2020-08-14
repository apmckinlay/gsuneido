// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuStr is a string Value
type SuStr string

var EmptyStr Value = SuStr("")

func (ss SuStr) Len() int {
	return len(ss)
}

// Value interface --------------------------------------------------

func (ss SuStr) ToInt() (int, bool) {
	return 0, ss == ""
}

func (ss SuStr) IfInt() (int, bool) {
	return 0, false
}

func (ss SuStr) ToDnum() (dnum.Dnum, bool) {
	return dnum.Zero, ss == ""
}

func (SuStr) ToContainer() (Container, bool) {
	return nil, false
}

func (ss SuStr) AsStr() (string, bool) {
	return string(ss), true
}

func (ss SuStr) ToStr() (string, bool) {
	return string(ss), true
}

var DefaultSingleQuotes = false

// String returns a human readable string with quotes and escaping
func (ss SuStr) String() string {
	return escapeStr(string(ss), 0)
}

func (ss SuStr) Display(t *Thread) string {
	q := 0
	if t != nil {
		q = t.Quote
	}
	return escapeStr(string(ss), q)
}

func escapeStr(s string, q int) string {
	if q == 0 &&
		strings.ContainsRune(s, '\\') && !strings.ContainsAny(s, "`\x00") {
		return "`" + s + "`"
	}
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "\x00", `\x00`)
	if q == 0 && !strings.ContainsRune(s, '\'') &&
		(DefaultSingleQuotes || strings.ContainsRune(s, '"')) {
		return "'" + s + "'"
	}
	if q == 1 {
		return "'" + strings.ReplaceAll(s, "'", `\'`) + "'"
	}
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

func (ss SuStr) Get(_ *Thread, key Value) Value {
	return strGet(string(ss), key)
}

// strGet is used by SuStr and SuConcat .Get
func strGet(s string, key Value) Value {
	i := ToIndex(key)
	n := len(s)
	if i < -n || n <= i {
		return EmptyStr
	}
	if i < 0 {
		i += n
	}
	return SuStr(s[i : i+1])
}

func (SuStr) Put(*Thread, Value, Value) {
	panic("strings do not support put")
}

func (ss SuStr) RangeTo(from int, to int) Value {
	size := len(ss)
	from = prepFrom(from, size)
	to = prepTo(from, to, size)
	return SuStr(string(ss)[from:to])
}

func (ss SuStr) RangeLen(from int, n int) Value {
	size := len(ss)
	from = prepFrom(from, size)
	n = prepLen(n, size-from)
	return SuStr(string(ss)[from : from+n])
}

func (ss SuStr) Hash() uint32 {
	return hash.HashString(string(ss))
}

func (ss SuStr) Hash2() uint32 {
	return ss.Hash()
}

func (ss SuStr) Equal(other interface{}) bool {
	if s2, ok := other.(SuStr); ok {
		return ss == s2
	}
	if cv, ok := other.(SuConcat); ok {
		return cv.n == len(ss) && string(ss) == cv.toStr()
	}
	if se, ok := other.(*SuExcept); ok {
		return ss == se.SuStr
	}
	return false
}

func (SuStr) Type() types.Type {
	return types.String
}

func (ss SuStr) Compare(other Value) int {
	if cmp := ints.Compare(ordStr, Order(other)); cmp != 0 {
		return cmp
	}
	return strings.Compare(string(ss), ToStr(other))
}

// Call implements s(ob, ...) being treated as ob[s](...)
func (ss SuStr) Call(t *Thread, _ Value, as *ArgSpec) Value {
	base := t.sp - int(as.Nargs)
	args := t.stack[base : base+int(as.Nargs)]
	k, v := NewArgsIter(as, args)()
	if v == nil || k != nil {
		panic("string call requires 'this' argument")
	}
	method := string(ss)
	fn := v.Lookup(t, method)
	if fn == nil {
		panic("method not found " + ErrType(v) + "." + method)
	}
	return fn.Call(t, v, as.DropFirst())
}

// StringMethods is initialized by the builtin package
var StringMethods Methods

var gnStrings = Global.Num("Strings")

func (SuStr) Lookup(t *Thread, method string) Callable {
	return Lookup(t, StringMethods, gnStrings, method)
}

func (SuStr) SetConcurrent() {
}

// Packable interface -----------------------------------------------

var _ Packable = SuStr("")

func (ss SuStr) PackSize(*int32) int {
	if ss == "" {
		return 0
	}
	return 1 + len(ss)
}

func (ss SuStr) PackSize2(int32, packStack) int {
	return ss.PackSize(nil)
}

func (ss SuStr) PackSize3() int {
	return ss.PackSize(nil)
}

func (ss SuStr) Pack(_ int32, buf *pack.Encoder) {
	if ss != "" {
		buf.Put1(PackString).PutStr(string(ss))
	}
}

// iterator ---------------------------------------------------------

type stringIter struct {
	s string
	i int
}

func (ss SuStr) Iter() Iter {
	return &stringIter{s: string(ss)}
}

func (si *stringIter) Next() Value {
	si.i++
	if si.i > len(si.s) {
		return nil
	}
	// can't use SuStr(si.s[si.i-1])
	// because > 127 turns into two byte string
	return SuStr(si.s[si.i-1:si.i])
}

func (si *stringIter) Dup() Iter {
	return &stringIter{s: si.s}
}

func (si *stringIter) Infinite() bool {
	return false
}
