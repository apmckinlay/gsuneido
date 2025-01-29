// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
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
	for i := range 12 {
		buf[i] = byte(i)
	}
}

func TestMmapRead(t *testing.T) {
	ms, _ := MmapStor("stor_test.go", Read) // use code as test file
	buf := ms.Data(0)
	assert.T(t).This(string(buf[:12])).Is("// Copyright")
	ms.Close(true)
}

func TestMmapWrite(t *testing.T) {
	ms, _ := MmapStor("stor_test.tmp", Create)
	const N = 100
	_, buf := ms.Alloc(N)
	for i := range N {
		buf[i] = byte(i)
	}
	ms.Close(true)

	ms, _ = MmapStor("stor_test.tmp", Update)
	data := ms.Data(0)
	for i := range N {
		assert.T(t).This(data[i]).Is(byte(i))
	}
	ms.Close(true)

	ms, _ = MmapStor("stor_test.tmp", Read)
	data = ms.Data(0)
	for i := range N {
		assert.T(t).This(data[i]).Is(byte(i))
	}
	ms.Close(true)

	os.Remove("stor_test.tmp")
}

func TestLastOffset(t *testing.T) {
	assert := assert.T(t).This
	ms, _ := MmapStor("stor_test.tmp", Create)
	defer os.Remove("stor_test.tmp")
	defer ms.Close(true)

	const N = 10
	const magic = "helloworld"
	for i := range N {
		off, buf := ms.Alloc(10)
		assert(off).Is(i * 100)
		copy(buf, magic)
		ms.Alloc(90)
	}

	off := ms.Size()/2 + 10
	for i := N / 2; i >= 0; i-- {
		off = ms.LastOffset(off, magic, nil)
		assert(off).Is(i * 100)
	}
	assert(ms.LastOffset(off, magic, nil)).Is(0)
}

func TestLastOffset2(t *testing.T) {
	const magic = "xyz"
	for _, size := range []int{32, 64, 96, 128, 160} {
		for _, pos := range []int{0, size / 2, size - len(magic)} {
			hs := HeapStor(64)
			for range size / 32 {
				hs.Alloc(32)
			}
			copy(hs.Data(uint64(pos)), magic)
			for _, off := range []int{0, size / 3, size - 10, size} {
				expected := pos
				if off <= pos {
					expected = 0
				}
				fmt.Println("size", size, "pos", pos, "off", off, "expected", expected)
				assert.T(t).
					This(hs.LastOffset(uint64(off), magic, nil)).
					Is(expected)
			}
		}
	}
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
	s, err := MmapStor("stor.tmp", Create)
	defer os.Remove("stor.tmp")
	defer s.Close(true)
	if err != nil {
		panic(err.Error())
	}
	for range nThreads {
		wg.Add(1)
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			for range nIterations {
				n := r.Intn(allocSize) + 1
				s.Alloc(n)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestAcessAfterClose(t *testing.T) {
	s, err := MmapStor("stor.tmp", Create)
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove("stor.tmp")
	off, buf1 := s.Alloc(100)
	slc.Fill(buf1, 'a')
	_, buf2 := s.Alloc(100)
	s.Close(false)
	buf1 = s.Data(off)
	assert.T(t).This(buf1[0]).Is('a')
	slc.Fill(buf2, 'b')
}

func TestSize(*testing.T) {
	s, err := MmapStor("stor.tmp", Create)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		s.Close(true)
		os.Remove("stor.tmp")
	}()
	allocSize := 48 * 1024 * 1024 // more than half a chunk
	for i := 0; i < int(s.chunksize); i += allocSize {
		assert.This(s.Size()).Is(i)
		_, buf := s.Alloc(allocSize)
		slc.Fill(buf, 123)
		s.Close(false)
		s, err = MmapStor("stor.tmp", Update)
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkFlush(b *testing.B) {
	s, err := MmapStor("stor.tmp", Create)
	if err != nil {
		panic(err.Error())
	}
	s.impl.(*mmapStor).mode = Update // flush doesn't run for Create
	defer func() {
		s.Close(true)
		os.Remove("stor.tmp")
	}()
	var flushing atomic.Bool
	for range b.N {
		off, buf := s.Alloc(1)
		slc.Fill(buf, 123)
		if flushing.CompareAndSwap(false, true) {
			go func() {
				s.FlushTo(off)
				flushing.Store(false)
			}()
		}
	}
}

func TestFlag(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var flushing atomic.Bool
	f := func() {
		if flushing.CompareAndSwap(false, true) {
			go func() {
				fmt.Println("Flushing")
				time.Sleep(time.Millisecond)
				flushing.Store(false)
			}()
		}
	}
	f()
	f()
	time.Sleep(10 * time.Millisecond)
	f()
	f()
	time.Sleep(10 * time.Millisecond)
}

func TestFlush(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	s, err := MmapStor("stor.tmp", Create)
	if err != nil {
		panic(err.Error())
	}
	s.impl.(*mmapStor).mode = Update // flush doesn't run for Create
	for range 8 {
		go func() {
			for {
				off, buf := s.Alloc(8)
				buf[0] = 123
				s.FlushTo(off)
			}
		}()
	}
	time.Sleep(1 * time.Second)
	s.Close(false)
	time.Sleep(1 * time.Second)
}
