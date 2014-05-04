/*
Package varint encodes and decodes integers as in Google Protocol Buffers
*/
package varint

// EncodeUint32 encodes an int32 to a series of bytes
func EncodeUint32(n uint32, buf []byte) []byte {
	for n >= 1<<7 {
		buf = append(buf, uint8(n&0x7f|0x80))
		n >>= 7
	}
	return append(buf, uint8(n))
}

// EncodeInt32 encodes an int32 using zigzag encoding
func EncodeInt32(n int32, buf []byte) []byte {
	return EncodeUint32((uint32(n)<<1)^uint32(n>>31), buf)
}

func DecodeInt32(buf []byte, start int) (n int32, i int) {
	x, i := DecodeUint32(buf, start)
	n = int32((x >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
}

func DecodeUint32(buf []byte, start int) (n uint32, i int) {
	i = start
	for shift := uint(0); shift < 32; shift += 7 {
		b := buf[i]
		i++
		n |= (uint32(b) & 0x7F) << shift
		if b < 0x80 {
			return
		}
	}
	panic("varint.DecodeUint32 overflow")
}
