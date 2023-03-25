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
	Read Mode = iota
	Create
	Update
)

type mmapStor struct {
	file *os.File
	ptrs []uintptr // needed on windows
	mode Mode
}

const mmapChunkSize = 64 * 1024 * 1024 // 64 mb

// MmapStor returns a memory mapped file stor.
func MmapStor(filename string, mode Mode) (*Stor, error) {
	var perm os.FileMode
	flags := os.O_RDONLY
	if mode == Update {
		flags = os.O_RDWR
	} else if mode == Create {
		perm = 0666
		flags = os.O_CREATE | os.O_TRUNC | os.O_RDWR
	}
	file, err := os.OpenFile(filename, flags, perm)
	if err != nil {
		return nil, err
	}
	if mode == Read {
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
	nchunks := int(((size + mmapChunkSize - 1) / mmapChunkSize))
	impl := &mmapStor{file: file, mode: mode}
	chunks := make([][]byte, nchunks)

	last := nchunks - 1
	for i := 0; i < nchunks; i++ {
		if i < last {
			impl.mode = Read // map full chunks as READ
		} else {
			impl.mode = mode // map last chunk with actual mode
		}
		chunks[i] = impl.Get(i)
	}
	if mode == Read {
		remainder := size % mmapChunkSize
		if remainder > 0 {
			chunks[last] = chunks[last][:remainder] // last chunk not full
		}
	}
	// ignore trailing zero bytes (from memory mapping, if truncate failed)
	if size > 0 {
		buf := chunks[last]
		r := (size - 1) % mmapChunkSize
		b := r
		for ; b >= 0 && buf[b] == 0; b-- {
		}
		size -= r - b
	}

	ms := NewStor(impl, mmapChunkSize, uint64(size), chunks)
	return ms, nil
}

// Write writes directly to the file, not via memory map
func (ms *mmapStor) Write(off uint64, data []byte) {
	ms.file.WriteAt(data, int64(off))
}
