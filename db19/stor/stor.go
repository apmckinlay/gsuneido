// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package stor is used to access physical storage,
normally by memory mapped file access.

Storage is chunked. Allocations may not straddle chunks.

Stor has an impl(ementation) of either heapstor or mmapstor.
OS dependent parts of mmapstor are in mmap_nonwin.go and mmap_windows.go.
*/
package stor

import (
	"bytes"
	"log"
	"math"
	"math/bits"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/exit"
)

// Offset is an offset within storage
type Offset = uint64

// storage is the interface to different kinds of storage,
// either mmapstor or heapstor (for tests).
type storage interface {
	// Get returns the i'th chunk of storage
	Get(chunk int) []byte
	// Flush writes changes to disk in the background
	Flush(chunk []byte)
	// Close closes the storage (if necessary)
	Close(size int64, unmap bool)
}

// Stor is the externally visible storage
type Stor struct {
	impl storage
	// chunks must be initialized up to size,
	// with at least one chunk if size is 0
	chunks atomic.Value // [][]byte
	// chunksize must be a power of two and must be initialized
	chunksize uint64
	// shift must be initialized to match chunksize
	shift int
	// size is the currently used amount.
	size atomic.Uint64
	// allocChunk is the chunk we're currently allocating in
	allocChunk atomic.Int64
	lock       sync.Mutex // guards extending the storage

	prevFlushChunk atomic.Int32
	curFlushChunk  atomic.Int32

	// flushes is incremented each time a flush is requested
	// and reset to zero by flusher
	flushes atomic.Uint32

	// flushChan is used to track whether flush is running or not
	// in a way that we can check without blocking (in FlushTo for persist)
	// or wait with a timeout (in flushWait for Close)
	flushChan chan struct{}
}

// closedSize needs to allow room to be incremented.
// NOTE: check for >= closedSize, not equals.
const closedSize = math.MaxUint64 / 2

func NewStor(impl storage, chunksize uint64, size uint64, chunks [][]byte) *Stor {
	shift := bits.TrailingZeros(uint(chunksize))
	assert.That(1<<shift == chunksize) // chunksize must be power of 2
	stor := &Stor{impl: impl, chunksize: chunksize, shift: shift}
	stor.size.Store(size)
	stor.chunks.Store(chunks)
	last := int64(len(chunks) - 1)
	stor.allocChunk.Store(last)
	stor.prevFlushChunk.Store(int32(max(0, last)))
	stor.flushChan = make(chan struct{}, 1)
	return stor
}

// Alloc allocates n bytes of storage and returns its Offset and byte slice
// Returning data here allows slicing to the correct length and capacity
// to prevent erroneously writing too far.
// If insufficient room in the current chunk, advance to next
// (allocations may not straddle chunks).
// Usual case is just an atomic increment,
// and an atomic load to check if hit the end of the chunk.
// Locking is only used to extend the storage.
func (s *Stor) Alloc(n int) (Offset, []byte) {
	assert.That(0 < n && n <= int(s.chunksize))
	const maxRetries = 3
	for range maxRetries {
		allocChunk := s.allocChunk.Load()
		// another thread could Alloc at this point and advance to the next chunk
		// which will be caught by the endchunk check (and retry)
		newsize := s.size.Add(uint64(n)) // serializable
		if newsize >= closedSize {
			log.Println("stor: use after close")
			runtime.Goexit()
		}
		endChunk := s.offsetToChunk(newsize - 1)
		offset := newsize - uint64(n)
		// this check catches two cases:
		// - allocation straddled a chunk boundary
		// - another thread bumped us into the next chunk
		if endChunk == int(allocChunk) {
			return offset, s.Data(offset)[:n:n]
		}
		// once one thread goes into the next chunk
		// following threads will also get here
		// (until extend increments allocChunk)
		s.extend(allocChunk)
		// loop to realloc
	}
	panic("Stor.Alloc too many retries")
}

// extend adds another chunk to the storage.
// Multiple concurrent threads may call extend
// but only one (the first) does the actual extending.
// The others wait on the lock.
func (s *Stor) extend(allocChunk int64) {
	s.lock.Lock() // note: lock does not prevent concurrent allocations
	defer s.lock.Unlock()
	chunks := s.chunks.Load().([][]byte)
	if int(allocChunk)+1 < len(chunks) {
		return // another thread beat us to it
	}
	chunks = append(chunks, s.impl.Get(int(allocChunk+1))) // potentially slow
	s.chunks.Store(chunks)
	// set size to start of chunk, to handle straddle
	s.size.Store(uint64(allocChunk+1) << s.shift)
	// NOTE: if another thread calls Alloc at this point
	// it will loop and realloc, wasting the first increment
	// but this should be rare and relatively harmless
	s.allocChunk.Add(1)
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
	return int(offset >> s.shift)
}

// Size returns the current (allocated) size of the data.
// The actual file size will be rounded up to the next chunk size.
func (s *Stor) Size() uint64 {
	size := s.size.Load()
	if size >= closedSize {
		log.Println("stor: use after close")
		runtime.Goexit()
	}
	return size
}

// FirstOffset searches forewards from a given offset for a given byte slice
// and returns the offset, or 0 if not found
func (s *Stor) FirstOffset(off uint64, str string) uint64 {
	b := []byte(str)
	chunks := s.chunks.Load().([][]byte)
	c := s.offsetToChunk(off)
	n := off & (s.chunksize - 1)
	for ; c < len(chunks); c++ {
		buf := chunks[c][n:]
		if i := bytes.Index(buf, b); i != -1 {
			return uint64(c)*s.chunksize + n + uint64(i)
		}
		n = 0
	}
	return 0
}

// LastOffset searches backwards from a given offset for a given string
// and returns the offset, or 0 if not found.
// It is used by repair and by asof/history.
func (s *Stor) LastOffset(off uint64, str string, stop *uint32) uint64 {
	b := []byte(str)
	chunks := s.chunks.Load().([][]byte)
	c := s.offsetToChunk(off)
	n := off & (s.chunksize - 1)
	if n == 0 {
		n = s.chunksize
		c--
	}
	for ; c >= 0; c-- {
		buf := chunks[c][:n]
		if i := bytes.LastIndex(buf, b); i != -1 {
			return uint64(c)*s.chunksize + uint64(i)
		}
		n = s.chunksize
		if stop != nil && atomic.LoadUint32(stop) == 1 {
			break
		}
	}
	return 0
}

type writable interface {
	Write(off uint64, data []byte)
}

func (s *Stor) Write(off uint64, data []byte) {
	if w, ok := s.impl.(writable); ok {
		w.Write(off, data)
	} else {
		copy(s.Data(off), data) // for testing with heap stor
	}
}

// Flush writes change to disk in the background.
// It is called by persist (e.g. once per minute).
func (s *Stor) FlushTo(offset uint64) {
	chunk := s.offsetToChunk(offset)
	s.curFlushChunk.Store(int32(chunk))
	s.flushes.Add(1) // flush requested
	select {
	case s.flushChan <- struct{}{}:
		go s.flusher()
	default:
		// flush already running
	}
}

func (s *Stor) flusher() {
	f := s.flushes.Swap(0)
	for f > 0 {
		s.flush()
		f = s.flushes.Swap(0)
		// loop if there have been flush requests during this flush
	}
	<-s.flushChan
}

func (s *Stor) flush() {
	chunks := s.chunks.Load().([][]byte)
	curFlushChunk := s.curFlushChunk.Load()
	prevFlushChunk := s.prevFlushChunk.Load()
	for c := prevFlushChunk; c <= curFlushChunk; c++ {
		// log.Println("flush", c)
		s.impl.Flush(chunks[c])
	}
	s.prevFlushChunk.Store(curFlushChunk)
}

// flushWait waits up to 5 seconds for flushes to finish.
func (s *Stor) flushWait() {
	select {
	case s.flushChan <- struct{}{}:
		<-s.flushChan
		exit.Progress("    store flushed")
		return
	case <-time.After(5 * time.Second):
		log.Println("ERROR: stor: FlushWait timed out")
		return
	}
}

func (s *Stor) Close(unmap bool, callback ...func(uint64)) {
	exit.Progress("  stor closing")
	var size uint64
	if _, ok := s.impl.(*heapStor); ok {
		size = s.size.Load() // for tests
	} else {
		size = s.size.Swap(closedSize)
	}
	if size < closedSize {
		for _, f := range callback {
			f(size)
		}
		s.flushWait()
		s.impl.Close(int64(size), unmap)
	}
	exit.Progress("  stor closed")
}
