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
	return SuStr1s[s[i]]
}

func (SuStr) Put(*Thread, Value, Value) {
	panic("string does not support put")
}

func (SuStr) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("string does not support update")
}

func (ss SuStr) RangeTo(from int, to int) Value {
	size := len(ss)
	from = prepFrom(from, size)
	to = prepTo(from, to, size)
	if from == 0 && to == size {
		return ss
	}
	return SuStr1(string(ss)[from:to])
}

func (ss SuStr) RangeLen(from int, n int) Value {
	size := len(ss)
	from = prepFrom(from, size)
	n = prepLen(n, size-from)
	if from == 0 && n == size {
		return ss
	}
	return SuStr1(string(ss)[from : from+n])
}

func (ss SuStr) Hash() uint64 {
	return hash.String(string(ss))
}

func (ss SuStr) Hash2() uint64 {
	return ss.Hash()
}

func (ss SuStr) Equal(other any) bool {
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
	if cmp := cmp.Compare(ordStr, Order(other)); cmp != 0 {
		return cmp * 2
	}
	return strings.Compare(string(ss), ToStr(other))
}

// Call implements s(ob, ...) being treated as ob[s](...)
func (ss SuStr) Call(th *Thread, _ Value, as *ArgSpec) Value {
	base := th.sp - int(as.Nargs)
	args := th.stack[base : base+int(as.Nargs)]
	k, v := NewArgsIter(as, args)()
	if v == nil || k != nil {
		panic("string call requires 'this' argument")
	}
	method := string(ss)
	fn := v.Lookup(th, method)
	if fn == nil {
		panic("method not found: " + ErrType(v) + "." + method)
	}
	return fn.Call(th, v, as.DropFirst())
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

var _ Packable = SuStr("")

func (ss SuStr) PackSize(*packing) int {
	if ss == "" {
		return 0
	}
	return 1 + len(ss)
}

func (ss SuStr) Pack(pk *packing) {
	if ss != "" {
		pk.Put1(PackString).PutStr(string(ss))
	}
}

// iterator ---------------------------------------------------------

type stringIter struct {
	s string
	i int
	MayLock
}

func (ss SuStr) Iter() Iter {
	return &stringIter{s: string(ss)}
}

func (si *stringIter) Next() Value {
	if si.Lock() {
		defer si.Unlock()
	}
	si.i++
	if si.i > len(si.s) {
		return nil
	}
	return SuStr1s[si.s[si.i-1]]
}

func (si *stringIter) Dup() Iter {
	return &stringIter{s: si.s}
}

func (si *stringIter) Infinite() bool {
	return false
}

func (si *stringIter) Instantiate() *SuObject {
	InstantiateMax(len(si.s))
	list := make([]Value, len(si.s))
	for i := range len(si.s) {
		list[i] = SuStr1s[si.s[i]]
	}
	return NewSuObject(list)
}

//-------------------------------------------------------------------

// Storing a Go string in an interface (like Value)
// requires allocating a 16 byte string structure (pointer and length)
// (in addition to the string content itself)
// To avoid this allocation for single byte strings
// we keep a prefabricated array of values.

func SuStr1(s string) Value {
	if len(s) == 1 {
		return SuStr1s[s[0]]
	}
	return SuStr(s)
}

var SuStr1s = [256]Value{
	SuStr("\x00"), SuStr("\x01"), SuStr("\x02"), SuStr("\x03"),
	SuStr("\x04"), SuStr("\x05"), SuStr("\x06"), SuStr("\x07"),
	SuStr("\x08"), SuStr("\x09"), SuStr("\x0a"), SuStr("\x0b"),
	SuStr("\x0c"), SuStr("\x0d"), SuStr("\x0e"), SuStr("\x0f"),
	SuStr("\x10"), SuStr("\x11"), SuStr("\x12"), SuStr("\x13"),
	SuStr("\x14"), SuStr("\x15"), SuStr("\x16"), SuStr("\x17"),
	SuStr("\x18"), SuStr("\x19"), SuStr("\x1a"), SuStr("\x1b"),
	SuStr("\x1c"), SuStr("\x1d"), SuStr("\x1e"), SuStr("\x1f"),
	SuStr("\x20"), SuStr("\x21"), SuStr("\x22"), SuStr("\x23"),
	SuStr("\x24"), SuStr("\x25"), SuStr("\x26"), SuStr("\x27"),
	SuStr("\x28"), SuStr("\x29"), SuStr("\x2a"), SuStr("\x2b"),
	SuStr("\x2c"), SuStr("\x2d"), SuStr("\x2e"), SuStr("\x2f"),
	SuStr("\x30"), SuStr("\x31"), SuStr("\x32"), SuStr("\x33"),
	SuStr("\x34"), SuStr("\x35"), SuStr("\x36"), SuStr("\x37"),
	SuStr("\x38"), SuStr("\x39"), SuStr("\x3a"), SuStr("\x3b"),
	SuStr("\x3c"), SuStr("\x3d"), SuStr("\x3e"), SuStr("\x3f"),
	SuStr("\x40"), SuStr("\x41"), SuStr("\x42"), SuStr("\x43"),
	SuStr("\x44"), SuStr("\x45"), SuStr("\x46"), SuStr("\x47"),
	SuStr("\x48"), SuStr("\x49"), SuStr("\x4a"), SuStr("\x4b"),
	SuStr("\x4c"), SuStr("\x4d"), SuStr("\x4e"), SuStr("\x4f"),
	SuStr("\x50"), SuStr("\x51"), SuStr("\x52"), SuStr("\x53"),
	SuStr("\x54"), SuStr("\x55"), SuStr("\x56"), SuStr("\x57"),
	SuStr("\x58"), SuStr("\x59"), SuStr("\x5a"), SuStr("\x5b"),
	SuStr("\x5c"), SuStr("\x5d"), SuStr("\x5e"), SuStr("\x5f"),
	SuStr("\x60"), SuStr("\x61"), SuStr("\x62"), SuStr("\x63"),
	SuStr("\x64"), SuStr("\x65"), SuStr("\x66"), SuStr("\x67"),
	SuStr("\x68"), SuStr("\x69"), SuStr("\x6a"), SuStr("\x6b"),
	SuStr("\x6c"), SuStr("\x6d"), SuStr("\x6e"), SuStr("\x6f"),
	SuStr("\x70"), SuStr("\x71"), SuStr("\x72"), SuStr("\x73"),
	SuStr("\x74"), SuStr("\x75"), SuStr("\x76"), SuStr("\x77"),
	SuStr("\x78"), SuStr("\x79"), SuStr("\x7a"), SuStr("\x7b"),
	SuStr("\x7c"), SuStr("\x7d"), SuStr("\x7e"), SuStr("\x7f"),
	SuStr("\x80"), SuStr("\x81"), SuStr("\x82"), SuStr("\x83"),
	SuStr("\x84"), SuStr("\x85"), SuStr("\x86"), SuStr("\x87"),
	SuStr("\x88"), SuStr("\x89"), SuStr("\x8a"), SuStr("\x8b"),
	SuStr("\x8c"), SuStr("\x8d"), SuStr("\x8e"), SuStr("\x8f"),
	SuStr("\x90"), SuStr("\x91"), SuStr("\x92"), SuStr("\x93"),
	SuStr("\x94"), SuStr("\x95"), SuStr("\x96"), SuStr("\x97"),
	SuStr("\x98"), SuStr("\x99"), SuStr("\x9a"), SuStr("\x9b"),
	SuStr("\x9c"), SuStr("\x9d"), SuStr("\x9e"), SuStr("\x9f"),
	SuStr("\xa0"), SuStr("\xa1"), SuStr("\xa2"), SuStr("\xa3"),
	SuStr("\xa4"), SuStr("\xa5"), SuStr("\xa6"), SuStr("\xa7"),
	SuStr("\xa8"), SuStr("\xa9"), SuStr("\xaa"), SuStr("\xab"),
	SuStr("\xac"), SuStr("\xad"), SuStr("\xae"), SuStr("\xaf"),
	SuStr("\xb0"), SuStr("\xb1"), SuStr("\xb2"), SuStr("\xb3"),
	SuStr("\xb4"), SuStr("\xb5"), SuStr("\xb6"), SuStr("\xb7"),
	SuStr("\xb8"), SuStr("\xb9"), SuStr("\xba"), SuStr("\xbb"),
	SuStr("\xbc"), SuStr("\xbd"), SuStr("\xbe"), SuStr("\xbf"),
	SuStr("\xc0"), SuStr("\xc1"), SuStr("\xc2"), SuStr("\xc3"),
	SuStr("\xc4"), SuStr("\xc5"), SuStr("\xc6"), SuStr("\xc7"),
	SuStr("\xc8"), SuStr("\xc9"), SuStr("\xca"), SuStr("\xcb"),
	SuStr("\xcc"), SuStr("\xcd"), SuStr("\xce"), SuStr("\xcf"),
	SuStr("\xd0"), SuStr("\xd1"), SuStr("\xd2"), SuStr("\xd3"),
	SuStr("\xd4"), SuStr("\xd5"), SuStr("\xd6"), SuStr("\xd7"),
	SuStr("\xd8"), SuStr("\xd9"), SuStr("\xda"), SuStr("\xdb"),
	SuStr("\xdc"), SuStr("\xdd"), SuStr("\xde"), SuStr("\xdf"),
	SuStr("\xe0"), SuStr("\xe1"), SuStr("\xe2"), SuStr("\xe3"),
	SuStr("\xe4"), SuStr("\xe5"), SuStr("\xe6"), SuStr("\xe7"),
	SuStr("\xe8"), SuStr("\xe9"), SuStr("\xea"), SuStr("\xeb"),
	SuStr("\xec"), SuStr("\xed"), SuStr("\xee"), SuStr("\xef"),
	SuStr("\xf0"), SuStr("\xf1"), SuStr("\xf2"), SuStr("\xf3"),
	SuStr("\xf4"), SuStr("\xf5"), SuStr("\xf6"), SuStr("\xf7"),
	SuStr("\xf8"), SuStr("\xf9"), SuStr("\xfa"), SuStr("\xfb"),
	SuStr("\xfc"), SuStr("\xfd"), SuStr("\xfe"), SuStr("\xff"),
}
