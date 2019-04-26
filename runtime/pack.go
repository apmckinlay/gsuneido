package runtime

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/pack"
)

// Packable is the interface to packable values
// PackSize should be called prior to Pack
// since Pack methods assume capacity is sufficient
// and because PackSize does nesting limit check
type Packable interface {
	// PackSize returns the size (in bytes) of the packed value
	PackSize(nest int) int
	// Pack appends the value to the Encoder
	Pack(buf *pack.Encoder)
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

// Pack is a convenience function that packs a single Packable
func Pack(x Packable) string {
	buf := pack.NewEncoder(x.PackSize(0))
	x.Pack(buf)
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
		panic("invalid pack tag")
	}
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
