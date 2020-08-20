// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"os"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestAlloc(t *testing.T) {
	assert := assert.T(t).This
	hs := HeapStor(64)
	offset, _ := hs.Alloc(12)
	assert(offset).Is(Offset(0))
	offset, _ = hs.Alloc(8)
	assert(offset).Is(Offset(12))
	offset, _ = hs.Alloc(8)
	assert(offset).Is(Offset(20))
	offset, _ = hs.Alloc(48) // requires new chunk
	assert(offset).Is(Offset(64))
}

func TestData(t *testing.T) {
	hs := HeapStor(64)
	hs.Alloc(12)
	offset, buf := hs.Alloc(12)
	assert.T(t).This(len(buf)).Is(12)             // Alloc gives correct length
	assert.T(t).This(len(hs.Data(offset))).Is(52) // Data gives to end of chunk
	for i := 0; i < 12; i++ {
		buf[i] = byte(i)
	}
}

func TestMmapRead(t *testing.T) {
	ms, _ := MmapStor("stor_test.go", READ)
	buf := ms.Data(0)
	assert.T(t).This(string(buf[:12])).Is("// Copyright")
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
		assert.T(t).This(data[i]).Is(byte(i))
	}
	ms.Close()

	os.Remove("stor_test.tmp")
}
