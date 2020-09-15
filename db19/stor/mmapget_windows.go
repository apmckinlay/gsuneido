// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"syscall"
	"unsafe"
)

func (ms *mmapStor) Get(chunk int) []byte {
	handle := syscall.Handle(ms.file.Fd())
	prot := uint32(syscall.PAGE_READWRITE)
	if ms.mode == READ {
		prot = syscall.PAGE_READONLY
	}
	end := int64((chunk + 1) * MMAP_CHUNKSIZE)
	if ms.mode == READ {
		fi, err := ms.file.Stat()
		if err != nil {
			panic(err)
		}
		end = fi.Size()
	}
	fm, err := syscall.CreateFileMapping(handle,
		nil, // no security attributes
		prot,
		uint32(end>>32), uint32(end&0xffffffff),
		nil) // no name for mapping
	if err != nil {
		panic(err)
	}

	access := uint32(syscall.FILE_MAP_WRITE)
	if ms.mode == READ {
		access = syscall.FILE_MAP_READ
	}
	offset := int64(chunk) * MMAP_CHUNKSIZE
	size := uintptr(MMAP_CHUNKSIZE)
	if ms.mode == READ {
		size = 0 // to end of file
	}
	ptr, err := syscall.MapViewOfFile(fm, access,
		uint32(offset>>32), uint32(offset), size)
	if err != nil {
		panic(err)
	}
	syscall.CloseHandle(fm)
	ms.ptrs = append(ms.ptrs, ptr)
	return (*[MMAP_CHUNKSIZE]byte)(unsafe.Pointer(ptr))[:]
}

func (ms mmapStor) Close(size int64) {
	// MSDN: Although an application may close the file handle used to create
	// a file mapping object, the system holds the corresponding file open
	// until the last view of the file is unmapped.
	for _, ptr := range ms.ptrs {
		syscall.UnmapViewOfFile(ptr)
	}
	ms.file.Truncate(size)
	ms.file.Close()
}
