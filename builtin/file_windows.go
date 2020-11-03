// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"io"
	"os"
	"syscall"
)

type myfile struct {
	h syscall.Handle
}

var _ iFile = (*myfile)(nil)

// appendFile is specifically for append on Windows
// to handle bug with RDP shared folders
func appendFile(name string) (iFile,error) {
	h,err := syscall.Open(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	return &myfile{syscall.Handle(h)}, err
}

func (mf *myfile) Read([]byte) (int, error) {
	panic("not implemented")
}

func (mf *myfile) Write(buf []byte) (int, error) {
	var done uint32
	err := syscall.WriteFile(mf.h, buf, &done, nil)
	return int(done), err
}

func (mf *myfile) Seek(int64, int) (int64, error) {
	panic("not implemented")
}

func fileSize(f iFile) int64 {
	var off int64
	var err error
	if mf, ok := f.(*myfile); ok {
		// only used for append, so ok to seek to end
		off, err = syscall.Seek(mf.h, 0, io.SeekEnd)
	} else {
		var info os.FileInfo
		info, err = f.(*os.File).Stat()
		off = info.Size()
	}
	if err != nil {
		panic("File: " + err.Error())
	}
	return off
}

func (mf *myfile) Close() error {
	return syscall.Close(mf.h)
}
