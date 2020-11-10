// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

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
	// object/record set clock to detect nested changes.
	PackSize(clock *int32) int
	// PackSize2 is used by object/record to handle nesting
	PackSize2(clock int32, stack packStack) int
	// PackSize3 is used by object/record during Pack
	PackSize3() int
	// Pack appends the value to the Encoder
	Pack(clock int32, buf *pack.Encoder)
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
)

var packClock = int32(0)

type packStack []Value

func newPackStack() packStack {
	// initialSize should handle almost all cases without further allocation
	const initialSize = 16
	return make([]Value, 0, initialSize)
}

func (ps *packStack) push(x Value) {
	for _, v := range *ps {
		if x == v { // NOTE: == not Equals
			panic("can't pack object/record containing itself")
		}
	}
	*ps = append(*ps, x)
}

// Note: no pop required because of passing slice by value

// Pack is a convenience function that packs a single Packable
func Pack(x Packable) string {
	var clock int32
	buf := pack.NewEncoder(x.PackSize(&clock))
	x.Pack(clock, buf)
	return buf.String()
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

func UnpackOld(s string) Value {
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
		return UnpackNumberOld(s)
	case PackObject:
		return UnpackObjectOld(s)
	case PackRecord:
		return UnpackRecordOld(s)
	default:
		panic("invalid pack tag " + strconv.Itoa(int(s[0])))
	}
}
