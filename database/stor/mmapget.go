// +build !windows

package stor

import (
	"syscall"
	"unsafe"
)

// NOTE: no provision for unmapping (same as Java)

func (ms mmapStor) Get(chunk int) []byte {
	prot := syscall.PROT_READ
	if ms.mode != READ {
		prot |= syscall.PROT_WRITE
		ms.file.Truncate(int64(chunk+1) * MMAP_CHUNKSIZE)
	}
	mmap, err := syscall.Mmap(int(ms.file.Fd()),
		int64(chunk)*MMAP_CHUNKSIZE, MMAP_CHUNKSIZE,
		prot, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return (*[MMAP_CHUNKSIZE]byte)(unsafe.Pointer(&mmap[0]))[:]
}

func (ms mmapStor) Close() {
	ms.file.Close()
}
