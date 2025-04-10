// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"io"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Packable is the interface to packable values
// PackSize should be called prior to Pack
// since Pack methods assume capacity is sufficient
// and because PackSize does nesting limit check
type Packable interface {
	// PackSize returns the size (in bytes) of the packed value.
	// object/record update hash to detect nested changes.
	PackSize(*packing) int
	// Pack appends the value to the Encoder in packing
	Pack(*packing)
}

type packing struct {
	pack.Encoder
	hash  uint64
	stack packStack
	v2    bool
}

func newPacking(size int) *packing {
	return &packing{Encoder: pack.NewEncoder(size)}
}

// Packed values start with one of the following type tags,
// except for the special case of a zero length string
// which is encoded as a zero length buffer.
// NOTE: this order is significant, it determines sorting
const (
	PackFalse = iota
	PackTrue
	PackMinus
	PackPlus
	PackString
	PackDate
	PackObject
	PackRecord
	PackForward // for query extend
)

type packStack []Value

func newPackStack() packStack {
	// initialSize should handle almost all cases without further allocation
	const initialSize = 16
	return make([]Value, 0, initialSize)
}

const nestingLimit = 16

func (ps *packStack) push(x Value) {
	if len(*ps) >= nestingLimit {
		panic("object nesting overflow")
	}
	for _, v := range *ps {
		if x == v { // NOTE: == not Equals
			panic("can't pack object/record containing itself")
		}
	}
	*ps = append(*ps, x)
}

func (ps *packStack) pop() {
	(*ps)[len(*ps)-1] = nil
	*ps = (*ps)[:len(*ps)-1]
}

var emptyStr = EmptyStr.(SuStr)
var boolTrue = True.(SuBool)
var boolFalse = False.(SuBool)
var zeroNum = Zero.(*smi)

var PackedTrue = string([]byte{PackTrue})
var PackedFalse = string([]byte{PackFalse})
var packedZero = string([]byte{PackPlus})

// Pack is a convenience function that packs a single Packable.
//
// WARNING: It's possible to get a buffer overflow if a mutable value
// (e.g. object) is modified between/during PackSize and Pack.
func Pack(x Packable) string {
	return packv(x, false)
}

func packv(x Packable, v2 bool) string {
	switch x {
	case emptyStr:
		return ""
	case boolTrue:
		return PackedTrue
	case boolFalse:
		return PackedFalse
	case zeroNum:
		return packedZero
	}

	pk := &packing{v2: v2}
	size := x.PackSize(pk)
	hash := pk.hash
	*pk = *newPacking(size)
	pk.v2 = v2
	x.Pack(pk)
	if hash != pk.hash || size != pk.Len() {
		panic("object modified during packing")
	}
	return pk.String()
}

func PackTo(v Value, w io.Writer) error {
	p, ok := v.(Packable)
	if !ok {
		panic("can't pack " + ErrType(v))
	}
	pk := &packing{Encoder: pack.NewEncoder2(w), v2: true}
	p.Pack(pk)
	return pk.Flush()
}

func UnpackFrom(r io.Reader) Value {
	d := pack.NewDecoder2(r)
    return unpack(d)
}

// Unpack returns the decoded value
func Unpack(s string) Value {
	if len(s) == 0 {
		return EmptyStr
	}
	switch s[0] {
	case PackFalse:
		return False
	case PackTrue:
		return True
	case PackString:
		return SuStr(s[1:])
	case PackDate:
		return UnpackDate(pack.NewDecoder(s))
	case PackPlus, PackMinus:
		return UnpackNumber(s)
	case PackObject:
		return UnpackObject(pack.NewDecoder(s))
	case PackRecord:
		return UnpackRecord(pack.NewDecoder(s))
	default:
		panic("invalid pack tag " + strconv.Itoa(int(s[0])))
	}
}

// Unpack returns the decoded value
func unpack(d pack.Decoder) Value {
	if !d.Remaining() {
		return EmptyStr
	}
	switch d.Peek() {
	case PackFalse:
		return False
	case PackTrue:
		return True
	case PackString:
		d.Skip(1)
		return SuStr(d.Remainder())
	case PackDate:
		return UnpackDate(d)
	case PackPlus, PackMinus:
		return UnpackNumber(d.TempRemainder())
	case PackObject:
		return UnpackObject(d)
	case PackRecord:
		return UnpackRecord(d)
	default:
		panic("invalid pack tag " + strconv.Itoa(int(d.Peek())))
	}
}

func unpackLen(d pack.Decoder, n int) Value {
	if n <= 0 {
		return EmptyStr
	}
	switch d.Peek() {
	case PackFalse:
		d.Skip(1)
		return False
	case PackTrue:
		d.Skip(1)
		return True
	case PackString:
		d.Skip(1)
		return SuStr(d.Get(n-1))
	case PackDate:
		return UnpackDate(pack.NewDecoder(d.TempStr(n)))
	case PackPlus, PackMinus:
		return UnpackNumber(d.TempStr(n))
	default:
		panic("invalid pack tag " + strconv.Itoa(int(d.Peek())))
	}
}

// PackedToLower applies str.ToLower to packed strings.
// Other types of values are unchanged.
func PackedToLower(s string) string {
	if len(s) == 0 || s[0] != PackString {
		return s
	}
	return str.ToLower(s) // ToLower shouldn't change PackString (4)
}

// PackedCmpLower compares strings with str.CmpLower
// and other values with strings.Compare
func PackedCmpLower(s1, s2 string) int {
	if len(s1) == 0 || s1[0] != PackString || len(s2) == 0 || s2[0] != PackString {
		return strings.Compare(s1, s2)
	}
	return str.CmpLower(s1, s2)
}

func PackedOrd(s string) Ord {
	if s == "" {
		return ordStr
	}
	switch s[0] {
	case PackFalse:
		return ordBool
	case PackTrue:
		return ordBool
	case PackMinus:
		return ordNum
	case PackPlus:
		return ordNum
	case PackString:
		return ordStr
	case PackDate:
		return ordDate
	case PackObject:
		return ordObject
	case PackRecord:
		return ordObject
	}
	panic("unknown")
}

func PackBool(b bool) string {
	if b {
		return PackedTrue
	}
	return PackedFalse
}

func UnpackBool(s string) Value {
	if s == PackedTrue {
		return True
	}
	if s == PackedFalse {
		return False
	}
	panic("can't convert to boolean")
}
