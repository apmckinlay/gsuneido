// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"math/rand"
	"testing"
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
		blockSize := uint64(512 + rand.Intn(1024))
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
