package runtime

// Packable is the interface to packable values
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

// Pack is a helper that allocates a buffer and packs a value into it
func Pack(x Packable) []byte {
	buf := make([]byte, 0, x.PackSize())
	return x.Pack(buf)
}

/*
Unpack returns the decoded value.
NOTE: The correct buf slice length is required.
*/
func Unpack(buf []byte) Value {
	if len(buf) == 0 {
		return SuStr("")
	}
	switch buf[0] {
	case packFalse:
		return False
	case packTrue:
		return True
	case packString:
		return UnpackSuStr(buf[1:])
	case packDate:
		return UnpackDate(buf[1:])
	case packPlus, packMinus:
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

func packInt32(n int32, buf []byte) []byte {
	// complement leading bit to ensure correct unsigned compare
	buf = append(buf, byte(n>>24)^0x80, byte(n>>16), byte(n>>8), byte(n))
	return buf
}

func unpackInt32(b []byte) int32 {
	n := int32(b[0]^0x80)<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	return n
}

func packUint32(n uint32, buf []byte) []byte {
	buf = append(buf, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
	return buf
}

func unpackUint32(b []byte) uint32 {
	n := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return n
}

func PackSizeInt64(n int64) int {
	if n == 0 {
		return 1
	}
	if n < 0 {
		n = -n
	}
	for n%10000 == 0 {
		n /= 10000
	}
	if n < 10000 {
		return 4
	} else if n < e8 {
		return 6
	} else if n < e12 {
		return 8
	}
	return 10
}

// PackInt64 packs an int64 in the same format as a Dnum
func PackInt64(n int64, buf []byte) []byte {
	buf = buf[:1]
	xor := uint16(0)
	if n < 0 {
		buf[0] = packMinus
		xor = 0xffff
		n = -n
	} else {
		buf[0] = packPlus
	}
	if n == 0 {
		return buf
	}
	e := 0
	for n%10000 == 0 {
		n /= 10000
		e++
	}
	var x uint16
	switch {
	case n < 10000:
		buf = buf[:4]
		buf[1] = byte(e+1^0x80) ^ byte(xor) // exponent
		x = uint16(n) ^ xor
		buf[2] = byte(x >> 8)
		buf[3] = byte(x)
	case n < 100000000:
		buf = buf[:6]
		buf[1] = byte(e+2^0x80) ^ byte(xor) // exponent
		x = uint16(n/10000) ^ xor
		buf[2] = byte(x >> 8)
		buf[3] = byte(x)
		x = uint16(n%10000) ^ xor
		buf[4] = byte(x >> 8)
		buf[5] = byte(x)
	case n < 1000000000000:
		buf = buf[:8]
		buf[1] = byte(e+3^0x80) ^ byte(xor) // exponent
		x = uint16(n/100000000) ^ xor
		buf[2] = byte(x >> 8)
		buf[3] = byte(x)
		n %= 100000000
		x = uint16(n/10000) ^ xor
		buf[4] = byte(x >> 8)
		buf[5] = byte(x)
		x = uint16(n % 10000)
		buf[6] = byte(x >> 8)
		buf[7] = byte(x)
	default:
		buf = buf[:10]
		buf[1] = byte(e+4^0x80) ^ byte(xor) // exponent
		x = uint16(n/1000000000000) ^ xor
		buf[2] = byte(x >> 8)
		buf[3] = byte(x)
		n %= 1000000000000
		x = uint16(n/100000000) ^ xor
		buf[4] = byte(x >> 8)
		buf[5] = byte(x)
		n %= 100000000
		x = uint16(n/10000) ^ xor
		buf[6] = byte(x >> 8)
		buf[7] = byte(x)
		x = uint16(n%10000) ^ xor
		buf[8] = byte(x >> 8)
		buf[9] = byte(x)
	}
	return buf
}

// Ensure ensures the slice has sufficient room
func Ensure(buf []byte, n int) []byte {
	if cap(buf)-len(buf) >= n {
		return buf
	}
	newbuf := make([]byte, len(buf), 2*(len(buf)+n)) // double
	copy(newbuf, buf)
	return newbuf
}
