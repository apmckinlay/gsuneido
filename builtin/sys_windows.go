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

// dll bool Kernel32:GetDiskFreeSpaceEx()
// [in] string			directoryName,
// ULARGE_INTEGER*		freeBytesAvailableToCaller,
// ULARGE_INTEGER*		totalNumberOfBytes,
// ULARGE_INTEGER*		totalNumberOfFreeBytes)
var getDiskFreeSpaceEx = kernel32.MustFindProc("GetDiskFreeSpaceExA").Addr()

var _ = builtin(GetDiskFreeSpace, "(dir = '.')")

func GetDiskFreeSpace(arg Value) Value {
	dir := zbuf(arg)
	var n int64
	syscall.SyscallN(getDiskFreeSpaceEx,
		uintptr(unsafe.Pointer(&dir[0])),
		uintptr(unsafe.Pointer(&n)),
		0,
		0)
	return Int64Val(n)
}

// zbuf returns a zero terminated byte slice copy of ToStr(v)
func zbuf(v Value) []byte {
	s := ToStr(v)
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = 0
	return buf
}

type stMemoryStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	unused       [6]uint64
}

const nMemoryStatusEx = unsafe.Sizeof(stMemoryStatusEx{})

var globalMemoryStatusEx = kernel32.MustFindProc("GlobalMemoryStatusEx").Addr()

var _ = builtin(SystemMemory, "()")

func SystemMemory() Value {
	buf := make([]byte, nMemoryStatusEx)
	(*stMemoryStatusEx)(unsafe.Pointer(&buf[0])).dwLength = uint32(nMemoryStatusEx)
	rtn, _, _ := syscall.SyscallN(globalMemoryStatusEx,
		uintptr(unsafe.Pointer(&buf[0])))
	if rtn == 0 {
		return Zero
	}
	return Int64Val(int64((*stMemoryStatusEx)(unsafe.Pointer(&buf[0])).ullTotalPhys))
}

var copyFile = kernel32.MustFindProc("CopyFileA").Addr()
var _ = builtin(CopyFile, "(from, to, failIfExists)")

func CopyFile(a, b, c Value) Value {
	from := zbuf(a)
	to := zbuf(b)
	rtn, _, _ := syscall.SyscallN(copyFile,
		uintptr(unsafe.Pointer(&from[0])),
		uintptr(unsafe.Pointer(&to[0])),
		boolArg(c))
	return boolRet(rtn)
}

func boolArg(arg Value) uintptr {
	if ToBool(arg) {
		return 1
	}
	return 0
}

func boolRet(rtn uintptr) Value {
	if rtn == 0 {
		return False
	}
	return True
}
