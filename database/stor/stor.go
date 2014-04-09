// stor is used to access physical storage.
// Storage is chunked. Allocations may not straddle chunks.
package stor

import . "verify"

// Adr is an int32 value specifying an int64 offset within storage
// int64 values are aligned and shifted to fit into int32
type Adr uint32

// StorImpl is the lowest level interface to storage.
// The main implementation accesses memory mapped files.
// There is also an in memory version for testing.
// Get returns the i'th chunk of storage
type storage interface {
	Get(chunk int) []byte
	Close()
}

type stor struct {
	impl      storage
	chunksize int64
	size      int64
	chunks    [][]byte
}

const ALIGN = 8

// Alloc allocates n bytes of storage and returns its Adr
func (s *stor) Alloc(n int) Adr {
	Verify(n < int(s.chunksize))
	n = Align(n)

	// if insufficient room in this chunk, advance to next
	remaining := s.chunksize - s.size%s.chunksize
	if n > int(remaining) {
		s.size += remaining
	}

	offset := s.size
	s.size += int64(n)
	return offsetToAdr(offset)
}

func Align(n int) int {
	// requires ALIGN to be power of 2
	return ((n - 1) | (ALIGN - 1)) + 1
}

func (s *stor) offsetToChunk(offset int64) int {
	return int(offset / s.chunksize)
}

func offsetToAdr(offset int64) Adr {
	return Adr((offset / ALIGN) + 1) // +1 to avoid 0
}

func adrToOffset(adr Adr) int64 {
	return int64(adr-1) * ALIGN
}

// Data returns the slice of bytes at the given adr.
// The slice extends to the end of the chunk,
// since we don't know the size of the original alloc.
func (s *stor) Data(adr Adr) []byte {
	offset := adrToOffset(adr)
	chunk := s.offsetToChunk(offset)
	for chunk >= len(s.chunks) {
		s.chunks = append(s.chunks, nil)
	}
	if s.chunks[chunk] == nil {
		s.chunks[chunk] = s.impl.Get(chunk)
	}
	return s.chunks[chunk][offset%s.chunksize:]
}

func (s *stor) Close() {
	s.impl.Close()
}
