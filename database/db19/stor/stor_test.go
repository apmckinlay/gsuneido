// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"os"
	"testing"
)

func TestAlloc(t *testing.T) {
	hs := HeapStor(64)
	offset, _ := hs.Alloc(12)
	Assert(t).That(offset, Is(Offset(0)))
	offset, _ = hs.Alloc(8)
	Assert(t).That(offset, Is(Offset(12)))
	offset, _ = hs.Alloc(8)
	Assert(t).That(offset, Is(Offset(20)))
	offset, _ = hs.Alloc(48) // requires new chunk
	Assert(t).That(offset, Is(Offset(64)))
}

func TestData(t *testing.T) {
	hs := HeapStor(64)
	hs.Alloc(12)
	offset, buf := hs.Alloc(12)
	Assert(t).That(len(buf), Is(12))             // Alloc gives correct length
	Assert(t).That(len(hs.Data(offset)), Is(52)) // Data gives to end of chunk
	for i := 0; i < 12; i++ {
		buf[i] = byte(i)
	}
}

func TestMmapRead(t *testing.T) {
	ms, _ := MmapStor("stor_test.go", READ)
	buf := ms.Data(0)
	Assert(t).That(string(buf[:12]), Is("// Copyright"))
	ms.Close()
}

func TestMmapWrite(t *testing.T) {
	ms, _ := MmapStor("stor_test.tmp", CREATE)
	const N = 100
	_, buf := ms.Alloc(N)
	for i := 0; i < N; i++ {
		buf[i] = byte(i)
	}
	ms.Close()

	ms, _ = MmapStor("stor_test.tmp", READ)
	data := ms.Data(0)
	for i := 0; i < N; i++ {
		Assert(t).That(data[i], Is(byte(i)))
	}
	ms.Close()

	os.Remove("stor_test.tmp")
}

