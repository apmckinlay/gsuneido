package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type memStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	unused       [6]uint64
}

var globalMemoryStatusEx = kernel32.MustFindProc("GlobalMemoryStatusEx").Addr()

var _ = builtin0("SystemMemory()", func() Value {
	msx := &memStatusEx{dwLength: 64}
	r, _, _ := syscall.Syscall(globalMemoryStatusEx, 1,
		uintptr(unsafe.Pointer(msx)),
		0, 0)
	if r == 0 {
		return Zero
	}
	return Int64Val(int64(msx.ullTotalPhys))
})

var _ = builtin0("OperatingSystem()", func() Value {
	return SuStr("Windows") //TODO version
})

// dll bool Kernel32:GetDiskFreeSpaceEx(
// 	[in] string			directoryName,
// 	ULARGE_INTEGER*		freeBytesAvailableToCaller,
// 	ULARGE_INTEGER*		totalNumberOfBytes,
// 	ULARGE_INTEGER*		totalNumberOfFreeBytes
// 	)
var getDiskFreeSpaceEx = kernel32.MustFindProc("GetDiskFreeSpaceExA").Addr()

var _ = builtin1("GetDiskFreeSpace(dir = '.')", func(arg Value) Value {
	var freeBytes int64
	syscall.Syscall6(getDiskFreeSpaceEx, 4,
		uintptr(stringArg(arg)),
		uintptr(unsafe.Pointer(&freeBytes)),
		0,
		0,
		0, 0)
	return Int64Val(freeBytes)
})
