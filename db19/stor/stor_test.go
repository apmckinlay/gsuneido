// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
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
	ms, _ := MmapStor("stor_test.go", READ) // use code as test file
	buf := ms.Data(0)
	assert.T(t).This(string(buf[:12])).Is("// Copyright")
	ms.Close(true)
}

func TestMmapWrite(t *testing.T) {
	ms, _ := MmapStor("stor_test.tmp", CREATE)
	const N = 100
	_, buf := ms.Alloc(N)
	for i := 0; i < N; i++ {
		buf[i] = byte(i)
	}
	ms.Close(true)

	ms, _ = MmapStor("stor_test.tmp", UPDATE)
	data := ms.Data(0)
	for i := 0; i < N; i++ {
		assert.T(t).This(data[i]).Is(byte(i))
	}
	ms.Close(true)

	ms, _ = MmapStor("stor_test.tmp", READ)
	data = ms.Data(0)
	for i := 0; i < N; i++ {
		assert.T(t).This(data[i]).Is(byte(i))
	}
	ms.Close(true)

	os.Remove("stor_test.tmp")
}

func TestLastOffset(t *testing.T) {
	assert := assert.T(t).This
	ms, _ := MmapStor("stor_test.tmp", CREATE)
	defer os.Remove("stor_test.tmp")
	defer ms.Close(true)

	const N = 10
	const magic = "helloworld"
	for i := 0; i < N; i++ {
		off, buf := ms.Alloc(10)
		assert(off).Is(i * 100)
		copy(buf, magic)
		ms.Alloc(90)
	}

	off := ms.Size()/2 + 10
	for i := N / 2; i >= 0; i-- {
		off = ms.LastOffset(off, magic)
		assert(off).Is(i * 100)
	}
	assert(ms.LastOffset(off, magic)).Is(0)
}

func TestStress(*testing.T) {
	var nThreads = 101
	var nIterations = 1_000_000
	const allocSize = 1024
	if testing.Short() {
		nThreads = 2
		nIterations = 10000
	}
	var wg sync.WaitGroup
	s, err := MmapStor("stor.tmp", CREATE)
	defer os.Remove("stor.tmp")
	defer s.Close(true)
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			for i := 0; i < nIterations; i++ {
				n := r.Intn(allocSize) + 1
				s.Alloc(n)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestAcessAfterClose(t *testing.T) {
	s, err := MmapStor("stor.tmp", CREATE)
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove("stor.tmp")
	off, buf1 := s.Alloc(100)
	slc.Fill(buf1, 'a')
	_, buf2 := s.Alloc(100)
	s.Close(true)
	buf1 = s.Data(off)
	assert.T(t).This(buf1[0]).Is('a')
	slc.Fill(buf2, 'b')
}
