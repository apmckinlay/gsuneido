// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
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
	// object/record set hash to detect nested changes.
	PackSize(hash *uint64) int
	// PackSize2 is used by object/record to handle nesting
	PackSize2(hash *uint64, stack packStack) int
	// Pack appends the value to the Encoder
	Pack(hash *uint64, buf *pack.Encoder)
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

// Note: no pop required because of passing slice by value

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
	return Pack2(x).String()
}

func Pack2(x Packable) *pack.Encoder {
	hash1 := uint64(17)
	size := x.PackSize(&hash1)
	buf := pack.NewEncoder(size)
	hash2 := uint64(17)
	x.Pack(&hash2, buf)
	if hash1 != hash2 || len(buf.Buffer()) != size {
		panic("object modified during packing")
	}
	return buf
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
		return UnpackDate(s)
	case PackPlus, PackMinus:
		return UnpackNumber(s)
	case PackObject:
		return UnpackObject(s)
	case PackRecord:
		return UnpackRecord(s)
	default:
		panic("invalid pack tag " + strconv.Itoa(int(s[0])))
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
