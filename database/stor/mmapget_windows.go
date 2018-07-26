package stor

import (
	"reflect"
	"syscall"
	"unsafe"
)

// NOTE: no provision for unmapping (same as Java)

func (ms *mmapStor) Get(chunk int) []byte {
	handle := syscall.Handle(ms.file.Fd())
	prot := uint32(syscall.PAGE_READONLY)
	if ms.mode != READ {
		prot = syscall.PAGE_READWRITE
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

	access := uint32(syscall.FILE_MAP_READ)
	if ms.mode != READ {
		access = syscall.FILE_MAP_WRITE
	}
	offset := int64(chunk) * MMAP_CHUNKSIZE
	size := uintptr(MMAP_CHUNKSIZE)
	if ms.mode == READ {
		size = 0 // to end of file
	}
	ptr, err := syscall.MapViewOfFile(fm, access,
		uint32(offset>>32), uint32(offset&0xffffffff), size)
	if err != nil {
		panic(err)
	}
	syscall.CloseHandle(fm)
	ms.ptrs = append(ms.ptrs, ptr)
	// this seems simpler, and is used by golang mmap_windows.go
	// but gives "possible misuse of unsafe.Pointer"
	//return (*[MMAP_CHUNKSIZE]byte)(unsafe.Pointer(ptr))[:]
	var slice []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdr.Data = ptr
	hdr.Len = MMAP_CHUNKSIZE
	hdr.Cap = MMAP_CHUNKSIZE
	return slice
}

func (ms mmapStor) Close() {
	for _, ptr := range ms.ptrs {
		syscall.UnmapViewOfFile(ptr)
	}
	ms.file.Close()
}
