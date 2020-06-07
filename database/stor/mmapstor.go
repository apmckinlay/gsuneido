// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"os"
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
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()
	nchunks := int(((size - 1) / MMAP_CHUNKSIZE) + 1)
	impl := &mmapStor{file, mode, nil}
	chunks := make([][]byte, nchunks)
	for i := 0; i < nchunks; i++ {
		chunks[i] = impl.Get(0)
	}
	ms := &Stor{impl: impl, chunksize: MMAP_CHUNKSIZE, size: uint64(size)}
	ms.chunks.Store(chunks)
	return ms, nil
}

//TODO map all but partial last chunk as read-only
