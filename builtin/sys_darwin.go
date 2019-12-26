// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"encoding/binary"
	"os"
	"syscall"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func Run() {
}

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

var _ = builtin0("OperatingSystem()", func() Value {
	return SuStr("darwin")
})

var _ = builtin1("GetDiskFreeSpace(dir = '.')", func(arg Value) Value {
	var stat syscall.Statfs_t
	syscall.Statfs(ToStr(arg), &stat)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	return Int64Val(int64(freeBytes))
})

var _ = builtin0("GetComputerName()", func() Value {
	name, err := os.Hostname()
	if err != nil {
		panic("GetComputerName " + err.Error())
	}
	return SuStr(name)
})
