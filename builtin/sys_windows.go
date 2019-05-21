package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
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
	return SuDnum{Dnum: dnum.FromInt(int64(msx.ullTotalPhys))}
})

var _ = builtin0("OperatingSystem()", func() Value {
	return SuStr("windows") //TODO version
})
