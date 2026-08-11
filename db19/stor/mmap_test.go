// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"fmt"
	"hash/crc32"
	"math/rand/v2"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// BenchmarkAligned1024 reads 1024-byte aligned blocks from suneido.db
func BenchmarkAligned1024(b *testing.B) {
	s, err := MmapStor("../../suneido.db", Read)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close(true)

	size := s.Size()
	blockSize := uint64(1024)
	offset := uint64(0)
	for b.Loop() {
		if offset+blockSize > size {
			offset = 0
		}
		data := s.Data(offset)
		// access all 1024 bytes to prevent optimization
		sum := byte(0)
		for i := 0; i < int(blockSize); i++ {
			sum += data[i]
		}
		_ = sum
	}
}

// BenchmarkRandom512to1536 reads random blocks of 512 to 1536 bytes from suneido.db
func BenchmarkRandom512to1536(b *testing.B) {
	s, err := MmapStor("../../suneido.db", Read)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close(true)

	size := s.Size()
	offset := uint64(0)
	for b.Loop() {
		// random size between 512 and 1536
		blockSize := uint64(512 + rand.IntN(1024))
		if offset+blockSize > size {
			offset = 0
		}
		data := s.Data(offset)
		// access data to prevent optimization
		sum := byte(0)
		for i := 0; i < int(blockSize); i++ {
			sum += data[i]
		}
		_ = sum
	}
}

// TestMmapWrite1GB creates a 1 GB file using MmapStor,
// with 8 goroutines writing random sized (1 to 4k byte) allocations.
func TestMmapWrite1GB(t *testing.T) {
	assert.TestOnlyIndividually(t)
	const nThreads = 8
	const target = 32 << 30 // 1 GB
	const targetPer = target / nThreads
	const maxSize = 4 * 1024
	const name = "stor.tmp"
	os.Remove(name)
	s, err := MmapStor(name, Create)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		s.Close(true)
		os.Remove(name)
	}()
	var wg sync.WaitGroup
	for range nThreads {
		wg.Go(func() {
			r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
			for size := 0; size < targetPer; {
				n := r.IntN(maxSize) + 1
				_, buf := s.Alloc(n)
				for i := range buf {
					buf[i] = byte(i)
				}
				size += n
			}
		})
	}
	wg.Wait()
	fmt.Println(s.Size())
}

// TestMmapWrite1GB64k creates a 1 GB file using MmapStor,
// with 8 goroutines allocating 64 kb at a time and then writing.
func TestMmapWrite1GB64k(t *testing.T) {
	assert.TestOnlyIndividually(t)
	const nThreads = 8
	const target = 32 << 30 // 1 GB
	const targetPer = target / nThreads
	const allocSize = 64 * 1024
	const maxSize = 4 * 1024
	const name = "stor.tmp"
	os.Remove(name)
	s, err := MmapStor(name, Create)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		s.Close(true)
		os.Remove(name)
	}()
	var wg sync.WaitGroup
	for range nThreads {
		wg.Go(func() {
			r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
			for total := 0; total < targetPer; {
				_, buf := s.Alloc(allocSize)
				for len(buf) > 0 {
					n := min(len(buf), r.IntN(maxSize)+1)
					for i := range buf[:n] {
						buf[i] = byte(i)
					}
					buf = buf[n:]
				}
				total += allocSize
			}
		})
	}
	wg.Wait()
	fmt.Println(s.Size())
}

func TestProcs(t *testing.T) {
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	fmt.Println("NumCPU", runtime.NumCPU())
}

// BenchmarkWrite1GBReadMiddle128MB creates a 1 GB file using MmapStor and
// then benchmarks reading 128 MB from the middle in random sized blocks
// between 128 and 1024 bytes.
func BenchmarkReadMmapContiguous(b *testing.B) {
	const name = "stor.bench.tmp"
	const fileSize = 1 << 30    // 1 GB
	const readSize = 128 << 20  // 128 MB
	const writeChunk = 64 << 20 // 64 MB
	readStart := uint64((fileSize - readSize) / 2)

	os.Remove(name)
	s, err := MmapStor(name, Create)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < fileSize/writeChunk; i++ {
		_, buf := s.Alloc(writeChunk)
		for j := range buf {
			buf[j] = byte(j)
		}
	}
	defer func() {
		s.Close(true)
		os.Remove(name)
	}()

	b.SetBytes(readSize)
	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	for b.Loop() {
		read := 0
		offset := readStart
		for read < readSize {
			size := 128 + r.IntN(1024-128+1)
			if remaining := readSize - read; size > remaining {
				size = remaining
			}
			if chunkRemain := writeChunk - int(offset&uint64(writeChunk-1)); size > chunkRemain {
				size = chunkRemain
			}
			data := s.Data(offset)[:size]
			_ = crc32.ChecksumIEEE(data)
			offset += uint64(size)
			read += size
		}
	}
}

// BenchmarkWrite1GBReadScattered128MB creates a 1 GB file using MmapStor and
// then benchmarks reading 128 MB scattered through it in random sized blocks
// between 128 and 1024 bytes, advancing by 7 times the block size each time.
func BenchmarkReadMmapInterleaved(b *testing.B) {
	const name = "stor.bench.tmp"
	const fileSize = 1 << 30    // 1 GB
	const readSize = 128 << 20  // 128 MB
	const writeChunk = 64 << 20 // 64 MB

	os.Remove(name)
	s, err := MmapStor(name, Create)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < fileSize/writeChunk; i++ {
		_, buf := s.Alloc(writeChunk)
		for j := range buf {
			buf[j] = byte(j)
		}
	}
	defer func() {
		s.Close(true)
		os.Remove(name)
	}()

	b.SetBytes(readSize)
	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	for b.Loop() {
		read := 0
		offset := uint64(0)
		for read < readSize {
			size := 128 + r.IntN(1024-128+1)
			if remaining := readSize - read; size > remaining {
				size = remaining
			}
			if fileRemain := int(fileSize - offset); fileRemain < size {
				size = fileRemain
			}
			if chunkRemain := writeChunk - int(offset&uint64(writeChunk-1)); size > chunkRemain {
				size = chunkRemain
			}
			data := s.Data(offset)[:size]
			_ = crc32.ChecksumIEEE(data)
			read += size
			offset = (offset + uint64(7*size)) % fileSize
		}
	}
}

// BenchmarkSliceReadScattered128MB fills a 128 MB slice and benchmarks
// reading it in random sized blocks between 128 and 1024 bytes,
// advancing by 7 times the block size each time.
func BenchmarkReadArray(b *testing.B) {
	const size = 128 << 20 // 128 MB
	var buf [size]byte
	for i := range buf {
		buf[i] = byte(i)
	}

	b.SetBytes(size)
	for b.Loop() {
		_ = crc32.ChecksumIEEE(buf[:])
	}
}
