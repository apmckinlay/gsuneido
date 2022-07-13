// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !windows

package stor

import (
	"syscall"

	"github.com/apmckinlay/gsuneido/db19/filelock"
)

// NOTE: no provision for unmapping (same as Java)

// Get returns a memory mapped portion of a file.
// It panics on error.
func (ms *mmapStor) Get(chunk int) []byte {
	prot := syscall.PROT_READ
	if ms.mode != READ {
		prot |= syscall.PROT_WRITE
		ms.file.Truncate(int64(chunk+1) * MMAP_CHUNKSIZE) // extend file
	}
	mmap, err := syscall.Mmap(int(ms.file.Fd()),
		int64(chunk)*MMAP_CHUNKSIZE, MMAP_CHUNKSIZE,
		prot, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return mmap
}

func (ms *mmapStor) Close(size int64, _ bool) {
	// could Munmap but doesn't seem necessary, at least on Mac
	ms.file.Truncate(size)
	filelock.Unlock(ms.file)
	ms.file.Close()
}
