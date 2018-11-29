package runtime

import (
	"strings"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// SuStr is a string Value
type SuStr string

var EmptyStr Value = SuStr("")
var _ Packable = SuStr("")

func (ss SuStr) ToInt() int {
	if ss.IsEmpty() {
		return 0
	}
	panic("can't convert String to integer")
}

func (ss SuStr) ToDnum() dnum.Dnum {
	if ss.IsEmpty() {
		return dnum.Zero
	}
	panic("can't convert String to number")
}

func (ss SuStr) ToStr() string {
	return string(ss)
}

var DefaultSingleQuotes = false

// String returns a human readable string with quotes and escaping
// TODO: handle escaping
func (ss SuStr) String() string {
	q := "\""
	if DefaultSingleQuotes {
		q = "'"
	}
	return q + string(ss) + q
}

func (ss SuStr) Get(key Value) Value {
	i := Index(key)
	n := len(ss)
	if i < -n || n <= i {
		return EmptyStr
	}
	if i < 0 {
		i += n
	}
	return SuStr(ss[i])
}

func (SuStr) Put(Value, Value) {
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
		// according to benchmark, this doesn't allocate
		return cv.n == len(ss) && string(ss) == cv.ToStr()
	}
	return false
}

func (ss SuStr) PackSize() int {
	if ss.IsEmpty() {
		return 0
	}
	return 1 + len(ss)
}

func (ss SuStr) Pack(buf []byte) []byte {
	if ss.IsEmpty() {
		return buf
	}
	buf = append(buf, packString)
	buf = append(buf, string(ss)...)
	return buf
}

func UnpackSuStr(buf []byte) Value {
	return SuStr(string(buf))
}

func (SuStr) TypeName() string {
	return "String"
}

func (SuStr) Order() Ord {
	return ordStr
}

func (ss SuStr) Compare(other Value) int {
	if cmp := ints.Compare(ss.Order(), other.Order()); cmp != 0 {
		return cmp
	}
	return strings.Compare(ss.ToStr(), other.ToStr())
}

// Call implements s(ob, ...) being treated as ob[s](...)
func (ss SuStr) Call(t *Thread, as *ArgSpec) Value {
	// TODO @args
	if as.Nargs < 1 {
		panic("string call requires 'this' argument")
	}
	ob := t.stack[t.sp-int(as.Nargs)]
	method := string(ss)
	fn := ob.Lookup(method)
	as2 := *as
	as2.Nargs--
	t.this = ob
	return fn.Call(t, &as2)
}

// StringMethods is initialized by the builtin package
var StringMethods Methods

func (SuStr) Lookup(method string) Value {
	return StringMethods[method]
}

func (ss SuStr) IsEmpty() bool {
	return len(ss) == 0
}
