// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows

package builtin

import (
	"encoding/binary"
	"os"
	"syscall"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("SystemMemory()", func() Value {
	s, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		panic(err)
	}
	var buf [8]byte
	copy(buf[:], s)
	m := binary.LittleEndian.Uint64(buf[:])
	return Int64Val(int64(m))
})

var _ = builtin1("GetDiskFreeSpace(dir = '.')", func(arg Value) Value {
	var stat syscall.Statfs_t
	syscall.Statfs(ToStr(arg), &stat)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	return Int64Val(int64(freeBytes))
})

// appendFile is specialized for Windows
func appendFile(name string) (iFile, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}

func fileSize(f iFile) int64 {
	info, _ := f.(*os.File).Stat()
	return info.Size()
}
