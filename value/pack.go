package value

type Packable interface {
	// PackSize returns the size (in bytes) of the packed value
	PackSize() int
	// Pack writes the value starting at len(buf)
	// and returns a slice with the len extended by the number of bytes used
	Pack(buf []byte) []byte
}

// Packed values start with one of the following type tags,
// except for the special case of a zero length string
// which is encoded as a zero length buffer.
const (
	FALSE = iota
	TRUE
	MINUS
	PLUS
	STRING
	DATE
	OBJECT
	RECORD
	FUNCTION
	CLASS
	// NOTE: this order is significant, it determines sorting
)

// Pack is a helper that allocates a buffer and packs a value into it
func Pack(x Packable) []byte {
	buf := make([]byte, 0, x.PackSize())
	return x.Pack(buf)
}

/*
Unpack returns the decoded value.

NOTE: The correct buffer length is required.
*/
func Unpack(buf []byte) Value {
	if len(buf) == 0 {
		return StrVal("")
	}
	switch buf[0] {
	case FALSE:
		return False
	case TRUE:
		return True
	case STRING:
		return UnpackStrVal(buf[1:])
	case DATE:
		return UnpackDate(buf[1:])
	case PLUS, MINUS:
		return UnpackNumber(rbuf{buf})
	default:
		panic("invalid pack tag")
	}
}

type rbuf struct {
	buf []byte
}

func (rb *rbuf) get() byte {
	b := rb.buf[0]
	rb.buf = rb.buf[1:]
	return b
}

func (rb *rbuf) getUint16() uint16 {
	n := uint16(rb.buf[0])<<8 | uint16(rb.buf[1])
	rb.buf = rb.buf[2:]
	return n
}

func (rb *rbuf) remaining() int {
	return len(rb.buf)
}

// support functions -----------------------------------------------------------

func packInt32(n int32, b []byte) []byte {
	// complement leading bit to ensure correct unsigned compare
	i := len(b)
	b = b[:i+4]
	b[i] = byte(n)
	b[i+1] = byte(n >> 8)
	b[i+2] = byte(n >> 16)
	b[i+3] = byte(n>>24) ^ 0x80
	return b
}

func unpackInt32(b []byte) int32 {
	n := int32(b[3]^0x80)<<24 | int32(b[2])<<16 | int32(b[1])<<8 | int32(b[0])
	return n
}

func packUint32(n uint32, b []byte) []byte {
	i := len(b)
	b = b[:i+4]
	b[i] = byte(n)
	b[i+1] = byte(n >> 8)
	b[i+2] = byte(n >> 16)
	b[i+3] = byte(n >> 24)
	return b
}

func unpackUint32(b []byte) uint32 {
	n := uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
	return n
}
