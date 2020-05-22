// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package stor is used to access physical storage,
normally by memory mapped file access.

Storage is chunked. Allocations may not straddle chunks.
*/
package stor

import (
	"sync"
	"sync/atomic"
	"unsafe"

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
	// size is the currently used amount.
	// It must be accessed in a thread safe way.
	size uint64
	// chunks must be initialized up to size,
	// with at least one chunk if size is 0
	chunks atomic.Value // [][]byte
	lock   sync.Mutex
}

// Alloc allocates n bytes of storage and returns its Offset and byte slice
// Returning data here allows slicing to correct length and capacity
// to prevent erroneously writing too far.
// If insufficient room in the current chunk, advance to next
// (allocations may not straddle chunks)
func (s *stor) Alloc(n int) (Offset, []byte) {
	verify.That(0 < n && n < int(s.chunksize))

	for {
		oldsize := atomic.LoadUint64(&s.size)
		offset := oldsize
		newsize := offset + uint64(n)
		if newsize&(s.chunksize-1) < uint64(n) {
			// straddle, need to get another chunk
			s.lock.Lock() // note: lock does not prevent concurrent allocations
			chunk := s.offsetToChunk(newsize)
			chunks := s.chunks.Load().([][]byte)
			if chunk >= len(chunks) {
				// no one else beat us to it
				chunks = append(chunks, s.impl.Get(chunk))
				s.chunks.Store(chunks)
			}
			s.lock.Unlock()

			offset = uint64(chunk) * s.chunksize
			newsize = offset + uint64(n)
		}
		// attempt to confirm our allocation
		if atomic.CompareAndSwapUint64(&s.size, oldsize, newsize) {
			return offset, s.data(offset)[:n:n] // fast path
		}
		// another thread beat us, loop and try again
	}
}

// Data returns the bytes at the given offset as a string.
// The string extends to the end of the chunk,
// since we don't know the size of the original alloc.
// The existing chunks must be mapped initially
// since lazily mapping would require locking.
func (s *stor) Data(offset Offset) string {
	b := s.data(offset)
	return *(*string)(unsafe.Pointer(&b)) // no alloc converstion to string
}

func (s *stor) data(offset Offset) []byte {
	chunk := s.offsetToChunk(offset)
	chunks := s.chunks.Load().([][]byte)
	c := chunks[chunk]
	return c[offset&(s.chunksize-1):]
}

func (s *stor) offsetToChunk(offset Offset) int {
	return int(offset / s.chunksize)
}

func (s *stor) Close() {
	s.impl.Close(int64(s.size))
}
