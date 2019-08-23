package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type memStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	unused       [6]uint64
}

var globalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")

var _ = builtin0("SystemMemory()", func() Value {
	msx := &memStatusEx{dwLength: 64}
	r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return Zero
	}
	return Int64Val(int64(msx.ullTotalPhys))
})

var _ = builtin0("OperatingSystem()", func() Value {
	return SuStr("Windows") //TODO version
})

var getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExA")

var _ = builtin1("GetDiskFreeSpace(dir = '.')", func(arg Value) Value {
	var freeBytes int64
	getDiskFreeSpaceEx.Call(
		uintptr(stringArg(arg)),
		uintptr(unsafe.Pointer(&freeBytes)), 0, 0)
	return Int64Val(freeBytes)
})
