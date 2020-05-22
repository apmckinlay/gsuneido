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
