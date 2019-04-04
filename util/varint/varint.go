/*
Package varint encodes and decodes integers as in Google Protocol Buffers

Similar to encoding/binary but with a slightly different interface.
*/
package varint

import "math/bits"

// Len returns the number of bytes required to varint encode.
// signed and unsigned are the same size
func Len(n uint64) int {
	if n == 0 {
		return 1
	}
	return (bits.Len64(n) + 6) / 7
}

// EncodeUint32 encodes a uint32 into buf
func EncodeUint32(n uint32, buf []byte) []byte {
	for n >= 1<<7 {
		buf = append(buf, uint8(n&0x7f|0x80))
		n >>= 7
	}
	return append(buf, uint8(n))
}

// EncodeInt32 encodes an int32 using zigzag encoding
func EncodeInt32(n int32, buf []byte) []byte {
	un := uint32(n) << 1
	if n < 0 {
		un = ^un
	}
	return EncodeUint32(un, buf)
}

// DecodeInt32 decodes an int32 from buf[start:]
// and returns the number and the end position.
func DecodeInt32(buf []byte, start int) (n int32, i int) {
	ux, i := DecodeUint32(buf, start)
	n = int32(ux >> 1)
	if ux&1 != 0 {
		n = ^n
	}
	return
}

// DecodeInt32 decodes a uint32 from buf[start:]
// and returns the number and the end position.
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

// DecodeUint64 decodes a uint32 from buf[start:]
// and returns the number and the end position.
func DecodeUint64(s string) (n uint64, i int) {
	i = 0
	for shift := uint(0); shift < 64; shift += 7 {
		b := s[i]
		i++
		n |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			return
		}
	}
	panic("varint.DecodeUint64 overflow")
}
