// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuStr is a string Value implemented as a thin wrapper around a Go string.
type SuStr string

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

func (ss SuStr) Display(th *Thread) string {
	q := 0
	if th != nil {
		q = th.Quote
	}
	return escapeStr(string(ss), q)
}

const singleQuote = '\''
const doubleQuote = '"'
const backQuote = '`'

func escapeStr(s string, which int) string {
	var q byte
	switch which {
	case 0:
		if DefaultSingleQuotes {
			q = singleQuote
		} else {
			q = bestQuote(s)
		}
	case 1:
		q = singleQuote
	case 2:
		q = doubleQuote
	default:
		panic("invalid quotes value")
	}
	if q == backQuote {
		buf := make([]byte, 2+len(s))
		buf[0] = backQuote
		copy(buf[1:], s)
		buf[1+len(s)] = backQuote
		return hacks.BStoS(buf)
	}
	return escape(s, q)
}

func bestQuote(s string) byte {
	badForSingle := 0
	badForDouble := 0
	canBack := true
	for _, c := range []byte(s) {
		switch c {
		case singleQuote:
			badForSingle++
		case doubleQuote:
			badForDouble++
		case '\\':
			badForSingle++
			badForDouble++
		default:
			if c == backQuote {
				canBack = false
			}
			if c < ' ' || '~' < c {
				canBack = false
				badForSingle++
				badForDouble++
			}
		}
	}
	if len(s) == 1 && (badForSingle == 0 ||
		(!canBack && badForSingle <= badForDouble)) {
		// prefer single quotes for single characters
		return singleQuote
	}
	if badForDouble == 0 {
		return doubleQuote
	}
	if badForSingle == 0 {
		return singleQuote
	}
	if canBack {
		return backQuote
	}
	if badForSingle < badForDouble {
		return singleQuote
	}
	return doubleQuote
}

func escape(s string, q byte) string {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c == q || c == '\\' || c < ' ' || '~' < c {
			break
		}
	}
	buf := make([]byte, 0, 2+len(s))
	buf = append(buf, q)
	buf = append(buf, s[:i]...)
	for ; i < len(s); i++ {
		switch c := s[i]; c {
		case q:
			buf = append(buf, '\\', q)
		case '\t':
			buf = append(buf, '\\', 't')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\\':
			buf = append(buf, '\\', '\\')
		default:
			if c < ' ' || '~' < c {
				buf = append(buf, '\\', 'x')
				if c < 16 {
					buf = append(buf, '0')
				}
				buf = strconv.AppendInt(buf, int64(c), 16)
			} else {
				buf = append(buf, c)
			}
		}
	}
	buf = append(buf, q)
	return hacks.BStoS(buf)
}

func (ss SuStr) Get(_ *Thread, key Value) Value {
	return strGet(string(ss), key)
}

// strGet is used by SuStr and SuConcat .Get
func strGet(s string, key Value) Value {
	i, ok := key.IfInt()
	if !ok {
		return nil
	}
	n := len(s)
	if i < -n || n <= i {
		return EmptyStr
	}
	if i < 0 {
		i += n
	}
	return StrVal1(s[i])
}

func (SuStr) Put(*Thread, Value, Value) {
	panic("string does not support put")
}

func (SuStr) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("string does not support update")
}

func (ss SuStr) RangeTo(from int, to int) Value {
	return strRangeTo(ss, string(ss), from, to)
}

func (ss SuStr) RangeLen(from int, n int) Value {
	return strRangeLen(ss, string(ss), from, n)
}

func (ss SuStr) Hash() uint64 {
	return hash.String(string(ss))
}

func (ss SuStr) Hash2() uint64 {
	return ss.Hash()
}

func (ss SuStr) Equal(other any) bool {
	if v, ok := other.(Value); ok {
		if s, ok := v.ToStr(); ok {
			return string(ss) == s
		}
	}
	return false
}

func (SuStr) Type() types.Type {
	return types.String
}

func (ss SuStr) Compare(other Value) int {
	if cmp := cmp.Compare(ordStr, Order(other)); cmp != 0 {
		return cmp * 2
	}
	return strings.Compare(string(ss), ToStr(other))
}

// Call implements s(ob, ...) being treated as ob[s](...)
func (ss SuStr) Call(th *Thread, _ Value, as *ArgSpec) Value {
	return stringCall(th, as, string(ss))
}

// StringMethods is initialized by the builtin package
var StringMethods Methods

var gnStrings = Global.Num("Strings")

func (SuStr) Lookup(th *Thread, method string) Value {
	return Lookup(th, StringMethods, gnStrings, method)
}

func (SuStr) SetConcurrent() {
	// immutable so ok
}

// Packable interface -----------------------------------------------

func (ss SuStr) PackSize(*uint64) int {
	if ss == "" {
		return 0
	}
	return 1 + len(ss)
}

func (ss SuStr) PackSize2(*uint64, PackStack) int {
	return ss.PackSize(nil)
}

func (ss SuStr) Pack(_ *uint64, buf *pack.Encoder) {
	if ss != "" {
		buf.Put1(PackString).PutStr(string(ss))
	}
}

// iterator ---------------------------------------------------------

func (ss SuStr) Iter() Iter {
	return &stringIter{s: string(ss)}
}
