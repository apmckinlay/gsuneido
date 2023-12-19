// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"testing"
	"unsafe"
)

func TestSmallOffset(t *testing.T) {
	test := func(n uint64) {
		t.Helper()
		so := newSmallOffset(n)
		n2 := so.offset()
		if n2 != n {
			t.Error("expected", n, "got", n2)
		}
	}
	test(0)
	test(1)
	test(123)
	test(123456789)
	test(MaxSmallOffset)

	var v [5]SmallOffset
	if unsafe.Sizeof(v) != 25 {
		t.Error("expected 25 got ", unsafe.Sizeof(v))
	}
}

func newSmallOffset(offset uint64) SmallOffset {
	var so SmallOffset
	so[0] = byte(offset)
	so[1] = byte(offset >> 8)
	so[2] = byte(offset >> 16)
	so[3] = byte(offset >> 24)
	so[4] = byte(offset >> 32)
	return so
}

func (so SmallOffset) offset() uint64 {
	return uint64(so[0]) +
		uint64(so[1])<<8 +
		uint64(so[2])<<16 +
		uint64(so[3])<<24 +
		uint64(so[4])<<32
}
