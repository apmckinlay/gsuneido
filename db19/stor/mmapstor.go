// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"os"

	"github.com/apmckinlay/gsuneido/db19/filelock"
	"github.com/apmckinlay/gsuneido/runtime"
)

type Mode int

const (
	READ Mode = iota
	CREATE
	UPDATE
)

type mmapStor struct {
	file *os.File
	mode Mode
	ptrs []uintptr // needed on windows
}

const MMAP_CHUNKSIZE = 64 * 1024 * 1024 // 64 mb

// MmapStor returns a memory mapped file stor.
func MmapStor(filename string, mode Mode) (*Stor, error) {
	var perm os.FileMode
	flags := os.O_RDONLY
	if mode == UPDATE {
		flags = os.O_RDWR
	} else if mode == CREATE {
		perm = 0666
		flags = os.O_CREATE | os.O_TRUNC | os.O_RDWR
	}
	file, err := os.OpenFile(filename, flags, perm)
	if err != nil {
		return nil, err
	}
	if mode == READ {
		err = filelock.RLock(file)
	} else {
		err = filelock.Lock(file)
	}
	if err != nil {
		runtime.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()
	nchunks := int(((size + MMAP_CHUNKSIZE - 1) / MMAP_CHUNKSIZE))
	impl := &mmapStor{file, mode, nil}
	chunks := make([][]byte, nchunks)

	for i := 0; i < nchunks; i++ {
		if i < nchunks-1 {
			impl.mode = READ // map full chunks as READ
		} else {
			impl.mode = mode // map last chunk with actual mode
		}
		chunks[i] = impl.Get(i)
	}
	last := nchunks - 1
	if mode == READ {
		remainder := size % MMAP_CHUNKSIZE
		if remainder > 0 {
			chunks[last] = chunks[last][:remainder] // last chunk not full
		}
	}
	// trim trailing zero bytes (from memory mapping)
	if size > 0 {
		buf := chunks[last]
		r := (size - 1) % MMAP_CHUNKSIZE
		b := r
		for ; b >= 0 && buf[b] == 0; b-- {
		}
		size -= int64(r - b)
	}

	ms := NewStor(impl, MMAP_CHUNKSIZE, uint64(size))
	ms.chunks.Store(chunks)
	return ms, nil
}

// Write writes directly to the file, not via memory map
func (ms *mmapStor) Write(off uint64, data []byte) {
	ms.file.WriteAt(data, int64(off))
}
