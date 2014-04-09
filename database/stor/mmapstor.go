package stor

import "os"

type Mode int

const (
	READ Mode = iota
	CREATE
	UPDATE
)

type mmapStor struct {
	file *os.File
	mode Mode
	ptrs []uintptr // needed on windows to unmap
}

const MMAP_CHUNKSIZE = 64 * 1024 * 1024 // 64 mb

// MmapStor returns a memory mapped file stor.
func MmapStor(filename string, mode Mode) (*stor, error) {
	var perm os.FileMode
	flags := os.O_RDONLY
	if mode == UPDATE {
		flags = os.O_RDWR
	} else if mode == CREATE {
		perm = 0666
		flags = os.O_CREATE | os.O_RDWR
	}
	file, err := os.OpenFile(filename, flags, perm)
	if err != nil {
		return nil, err
	}
	return &stor{chunksize: MMAP_CHUNKSIZE,
		impl: &mmapStor{file, mode, nil}}, nil
}
