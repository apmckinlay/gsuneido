// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var kernel32 = windows.MustLoadDLL("kernel32.dll")

// dll bool Kernel32:GetDiskFreeSpaceEx(
// 	[in] string			directoryName,
// 	ULARGE_INTEGER*		freeBytesAvailableToCaller,
// 	ULARGE_INTEGER*		totalNumberOfBytes,
// 	ULARGE_INTEGER*		totalNumberOfFreeBytes)
var getDiskFreeSpaceEx = kernel32.MustFindProc("GetDiskFreeSpaceExA").Addr()

var _ = builtin1("GetDiskFreeSpace(dir = '.')", func(arg Value) Value {
	dir := zbuf(arg)
	var n int64
	syscall.SyscallN(getDiskFreeSpaceEx,
		uintptr(unsafe.Pointer(&dir[0])),
		uintptr(unsafe.Pointer(&n)),
		0,
		0)
	return Int64Val(n)
})

// zbuf returns a zero terminated byte slice copy of ToStr(v)
func zbuf(v Value) []byte {
	s := ToStr(v)
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = 0
	return buf
}

type MEMORYSTATUSEX struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	unused       [6]uint64
}

const nMEMORYSTATUSEX = unsafe.Sizeof(MEMORYSTATUSEX{})

var globalMemoryStatusEx = kernel32.MustFindProc("GlobalMemoryStatusEx").Addr()

var _ = builtin0("SystemMemory()", func() Value {
	buf := make([]byte, nMEMORYSTATUSEX)
	(*MEMORYSTATUSEX)(unsafe.Pointer(&buf[0])).dwLength = uint32(nMEMORYSTATUSEX)
	rtn, _, _ := syscall.SyscallN(globalMemoryStatusEx,
		uintptr(unsafe.Pointer(&buf[0])))
	if rtn == 0 {
		return Zero
	}
	return Int64Val(int64((*MEMORYSTATUSEX)(unsafe.Pointer(&buf[0])).ullTotalPhys))
})
