package builtin

import (
	"syscall"
	"unsafe"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var wtsapi32 = windows.MustLoadDLL("wtsapi32.dll")

// dll void WTSAPI32:WTSFreeMemory(pointer adr)
var wtsFreeMemory = wtsapi32.MustFindProc("WTSFreeMemory").Addr()

func WTSFreeMemory(adr uintptr) {
	syscall.Syscall(wtsFreeMemory, 1,
		adr,
		0, 0)
}

// dll bool WTSAPI32:WTSQuerySessionInformation(pointer hServer, long SessionId,
//		long WTSInfoClass, POINTER* ppBuffer, LONG* pBytesReturned)
var wtsQuerySessionInformation = wtsapi32.MustFindProc("WTSQuerySessionInformationA").Addr()

const WTS_CURRENT_SERVER_HANDLE = 0
const WTS_CURRENT_SESSION = ^uintptr(0) // -1
const WTS_ClientProtocolType = 16
const WTS_SessionId = 4

var _ = builtin0("WTS_GetClientProtocolType()",
	func() Value {
		defer heap.FreeTo(heap.CurSize())
		pbuf := heap.Alloc(uintptrSize)
		psize := heap.Alloc(int32Size)
		rtn, _, _ := syscall.Syscall6(wtsQuerySessionInformation, 5,
			WTS_CURRENT_SERVER_HANDLE,
			WTS_CURRENT_SESSION,
			WTS_ClientProtocolType,
			uintptr(pbuf),
			uintptr(psize),
			0)
		buf := *(*uintptr)(pbuf)
		size := *(*int32)(psize)
		if rtn == 0 || size != 2 || buf == 0 {
			return Zero
		}
		data := *(*int16)(unsafe.Pointer(buf))
		WTSFreeMemory(buf)
		return IntVal(int(data))
	})

var _ = builtin0("WTS_GetSessionId()",
	func() Value {
		defer heap.FreeTo(heap.CurSize())
		pbuf := heap.Alloc(uintptrSize)
		psize := heap.Alloc(int32Size)
		rtn, _, _ := syscall.Syscall6(wtsQuerySessionInformation, 5,
			WTS_CURRENT_SERVER_HANDLE,
			WTS_CURRENT_SESSION,
			WTS_SessionId,
			uintptr(pbuf),
			uintptr(psize),
			0)
		buf := *(*uintptr)(pbuf)
		size := *(*int32)(psize)
		if rtn == 0 || size != 4 || buf == 0 {
			return Zero
		}
		data := *(*int32)(unsafe.Pointer(buf))
		WTSFreeMemory(buf)
		return IntVal(int(data))
	})
