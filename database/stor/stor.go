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

// Stor is the externally visible storage
type Stor struct {
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
// Returning data here allows slicing to the correct length and capacity
// to prevent erroneously writing too far.
// If insufficient room in the current chunk, advance to next
// (allocations may not straddle chunks)
func (s *Stor) Alloc(n int) (Offset, []byte) {
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
			return offset, s.Data(offset)[:n:n] // fast path
		}
		// another thread beat us, loop and try again
	}
}

// Data returns a byte slice starting at the given offset
// and extending to the end of the chunk
// since we don't know the size of the original alloc.
func (s *Stor) Data(offset Offset) []byte {
	// The existing chunks must be mapped initially
	// since lazily mapping would require locking.
	chunk := s.offsetToChunk(offset)
	chunks := s.chunks.Load().([][]byte)
	c := chunks[chunk]
	return c[offset&(s.chunksize-1):]
}

func (s *Stor) offsetToChunk(offset Offset) int {
	return int(offset / s.chunksize)
}

func (s *Stor) Close() {
	s.impl.Close(int64(s.size))
}
