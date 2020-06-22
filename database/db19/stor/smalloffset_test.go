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
		so := NewSmallOffset(n)
		n2 := so.Offset()
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
