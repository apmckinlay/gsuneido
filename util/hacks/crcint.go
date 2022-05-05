// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hacks

import (
	"hash/crc32"
	"unsafe"
)

func CrcUint32(crc, n uint32) uint32 {
	p := unsafe.Pointer(&n)
	a := (*[4]byte)(p)
	return crc32.Update(crc, crc32.IEEETable, a[:])
}
