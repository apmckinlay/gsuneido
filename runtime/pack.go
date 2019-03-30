package runtime

import "github.com/apmckinlay/gsuneido/util/pack"

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
	packFalse = iota
	packTrue
	packMinus
	packPlus
	packString
	packDate
	packObject
	packRecord
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
	case packFalse:
		return False
	case packTrue:
		return True
	case packString:
		return SuStr(s[1:])
	case packDate:
		return UnpackDate(s)
	case packPlus, packMinus:
		return UnpackNumber(s)
	case packObject:
		return UnpackObject(s)
	default:
		panic("invalid pack tag")
	}
}

