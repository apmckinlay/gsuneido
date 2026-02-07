// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
	"golang.org/x/sys/windows"
)

// zstrArg returns a nul terminated copy of a string as *byte
func zstrArg(v Value) *byte {
	// NOTE: don't change this to return uintptr or unsafe.Pointer
	// Conversion from pointer to uintptr must be in the SyscallN argument.
	// Then it will be kept alive until the syscall returns.
	if v.Equal(Zero) {
		return nil
	}
	s := ToStr(v)
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	return &buf[0]
}

var kernel32 = windows.MustLoadDLL("kernel32.dll")

var getDiskFreeSpaceEx = kernel32.MustFindProc("GetDiskFreeSpaceExA").Addr()

var _ = builtin(GetDiskFreeSpace, "(dir = '.')")

func GetDiskFreeSpace(arg Value) Value {
	var n int64
	rtn, _, e := syscall.SyscallN(getDiskFreeSpaceEx,
		uintptr(unsafe.Pointer(zstrArg(arg))),
		uintptr(unsafe.Pointer(&n)),
		0,
		0)
	if rtn == 0 {
		panic("GetDiskFreeSpace: " + e.Error())
	}
	return Int64Val(n)
}

type stMemoryStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	unused       [6]uint64
}

const nMemoryStatusEx = unsafe.Sizeof(stMemoryStatusEx{})

var globalMemoryStatusEx = kernel32.MustFindProc("GlobalMemoryStatusEx").Addr()

func systemMemory() uint64 {
	mse := stMemoryStatusEx{dwLength: uint32(nMemoryStatusEx)}
	rtn, _, e := syscall.SyscallN(globalMemoryStatusEx,
		uintptr(unsafe.Pointer(&mse)))
	if rtn == 0 {
		panic("SystemMemory: " + e.Error())
	}
	return mse.ullTotalPhys
}

var copyFile = kernel32.MustFindProc("CopyFileA").Addr()
var _ = builtin(CopyFile, "(from, to, failIfExists)")

func CopyFile(th *Thread, args []Value) Value {
	if sandboxed() {
		from := ToStr(args[0])
		to := ToStr(args[1])
		fromPath, err := sandboxPath("CopyFile", from)
		if err != nil {
			th.ReturnThrow = true
			return SuStr("CopyFile: " + err.Error())
		}
		toPath, err := sandboxPath("CopyFile", to)
		if err != nil {
			th.ReturnThrow = true
			return SuStr("CopyFile: " + err.Error())
		}
		args = []Value{SuStr(fromPath), SuStr(toPath), args[2]}
	}
	rtn, _, e := syscall.SyscallN(copyFile,
		uintptr(unsafe.Pointer(zstrArg(args[0]))),
		uintptr(unsafe.Pointer(zstrArg(args[1]))),
		boolArg(args[2]))
	if rtn == 0 {
		th.ReturnThrow = true
		return SuStr("CopyFile: " + e.Error())
	}
	return True
}

func boolArg(arg Value) uintptr {
	if ToBool(arg) {
		return 1
	}
	return 0
}

func boolRet(rtn uintptr) Value { //nolint
	if rtn == 0 {
		return False
	}
	return True
}

var deleteFileA = kernel32.MustFindProc("DeleteFileA").Addr()
var setFileAttributesA = kernel32.MustFindProc("SetFileAttributesA").Addr()

const access_denied = 5

func deleteFile(filename string) error {
	file := zstrArg(SuStr(filename))
	rtn, _, e := syscall.SyscallN(deleteFileA,
		uintptr(unsafe.Pointer(file)))
	if rtn == 0 && e == access_denied {
		// retry after removing the read-only attribute
		r, _, _ := syscall.SyscallN(setFileAttributesA,
			uintptr(unsafe.Pointer(file)),
			windows.FILE_ATTRIBUTE_NORMAL)
		if r != 0 {
			rtn, _, e = syscall.SyscallN(deleteFileA,
				uintptr(unsafe.Pointer(file)))
		}
	}
	if rtn == 0 {
		return e
	}
	return nil
}
