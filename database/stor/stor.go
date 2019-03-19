/*
Package stor is used to access physical storage,
normally by memory mapped file access.

Storage is chunked. Allocations may not straddle chunks.
*/
package stor

import (
	"sync"

	"github.com/apmckinlay/gsuneido/util/verify"
)

// Offset is an offset within storage
type Offset = uint64

// storage is the interface to different kinds of storage.
// The main implementation accesses memory mapped files.
// There is also an in memory version for testing.
type storage interface {
	// Get returns the i'th chunk of storage
	Get(chunk int) []byte
	// Close closes the storage (if necessary)
	Close(size int64)
}

// stor is the externally visible storage
type stor struct {
	impl storage
	// chunksize must be a power of two and must be initialized
	chunksize uint64
	size      uint64
	// chunks must be initialized up to size,
	// with at least one chunk if size is 0
	chunks [][]byte
	lock   sync.Mutex
}

// Alloc allocates n bytes of storage and returns its Offset and data
// Returning data here allows slicing to correct length and capacity
// to prevent erroneously writing too far.
// Alloc is threadsafe, guarded by s.lock
func (s *stor) Alloc(n int) (Offset, []byte) {
	verify.That(n < int(s.chunksize))
	s.lock.Lock()

	// if insufficient room in this chunk, advance to next
	// (allocations may not straddle chunks)
	remaining := s.chunksize - s.size&(s.chunksize-1)
	if n > int(remaining) {
		s.size += remaining
		chunk := s.offsetToChunk(s.size)
		s.chunks = append(s.chunks, s.impl.Get(chunk))
	}

	offset := s.size
	s.size += uint64(n)
	s.lock.Unlock()
	return offset, s.Data(offset)[:n:n]
}

// Data returns the slice of bytes at the given offset.
// The slice extends to the end of the chunk,
// since we don't know the size of the original alloc.
// NOTE: There is no locking.
// Code using this must ensure it locks at some point
// prior to accessing "new" chunks from another thread's Alloc.
// This requires mapping the existing chunks initially
// since lazily mapping would require locking.
func (s *stor) Data(offset Offset) []byte {
	chunk := s.offsetToChunk(offset)
	c := s.chunks[chunk]
	return c[offset&(s.chunksize-1):]
}

func (s *stor) offsetToChunk(offset Offset) int {
	return int(offset / s.chunksize)
}

func (s *stor) Close() {
	s.impl.Close(int64(s.size))
}
