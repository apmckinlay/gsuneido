// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stor

import (
	"syscall"
	"unsafe"

	"github.com/apmckinlay/gsuneido/db19/filelock"
)

func (ms *mmapStor) Get(chunk int) []byte {
	handle := syscall.Handle(ms.file.Fd())
	prot := uint32(syscall.PAGE_READWRITE)
	if ms.mode == Read {
		prot = syscall.PAGE_READONLY
	}
	end := int64((chunk + 1) * mmapChunkSize)
	if ms.mode == Read {
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
	if ms.mode == Read {
		access = syscall.FILE_MAP_READ
	}
	offset := int64(chunk) * mmapChunkSize
	size := uintptr(mmapChunkSize)
	if ms.mode == Read {
		size = 0 // to end of file
	}
	ptr, err := syscall.MapViewOfFile(fm, access,
		uint32(offset>>32), uint32(offset), size)
	if err != nil {
		panic(err)
	}
	syscall.CloseHandle(fm)
	ms.ptrs = append(ms.ptrs, ptr)
	return (*[mmapChunkSize]byte)(unsafe.Pointer(ptr))[:]
}

func (ms mmapStor) Close(size int64, unmap bool) {
	// Things like -load need to unmap in order to close the file
	// in order to rename it. But for the server we do NOT want to unmap
	// because then threads get access violations during shutdown.
	// MSDN: Although an application may close the file handle used to create
	// a file mapping object, the system holds the corresponding file open
	// until the last view of the file is unmapped.
	if unmap {
		for _, ptr := range ms.ptrs {
			syscall.UnmapViewOfFile(ptr)
		}
	}
	ms.file.Truncate(size)
	filelock.Unlock(ms.file)
	ms.file.Close()
}
