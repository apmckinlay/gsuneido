// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

// SmallOffset is used to store database offsets in the database file
// to save space rather than using int64.
// 5 bytes = 40 bits = 1tb
// In memory we use int64.
// Beware of padding or you won't actually save space.
type SmallOffset [5]byte

const MaxSmallOffset = 1<<40 - 1
const SmallOffsetLen = 5

func NewSmallOffset(offset uint64) SmallOffset {
	var so SmallOffset
	so[0] = byte(offset)
	so[1] = byte(offset >> 8)
	so[2] = byte(offset >> 16)
	so[3] = byte(offset >> 24)
	so[4] = byte(offset >> 32)
	return so
}

func (so SmallOffset) Offset() uint64 {
	return uint64(so[0]) +
		uint64(so[1])<<8 +
		uint64(so[2])<<16 +
		uint64(so[3])<<24 +
		uint64(so[4])<<32
}

func WriteSmallOffset(buf []byte, offset uint64) {
	buf[0] = byte(offset)
	buf[1] = byte(offset >> 8)
	buf[2] = byte(offset >> 16)
	buf[3] = byte(offset >> 24)
	buf[4] = byte(offset >> 32)
}

func AppendSmallOffset(buf []byte, offset uint64) []byte {
	return append(buf,
		byte(offset),
		byte(offset>>8),
		byte(offset>>16),
		byte(offset>>24),
		byte(offset>>32))
}

func ReadSmallOffset(buf []byte) ([]byte, uint64) {
	return buf[5:],
		uint64(buf[0]) +
		uint64(buf[1])<<8 +
		uint64(buf[2])<<16 +
		uint64(buf[3])<<24 +
		uint64(buf[4])<<32
}
