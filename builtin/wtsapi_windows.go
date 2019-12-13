package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var wtsapi32 = windows.MustLoadDLL("wtsapi32.dll")

// dll void WTSAPI32:WTSFreeMemory(pointer adr)
var wtsFreeMemory = wtsapi32.MustFindProc("WTSFreeMemory").Addr()

func WTSFreeMemory(adr uintptr) {
	goc.Syscall1(wtsFreeMemory,
		adr)
}

// dll bool WTSAPI32:WTSQuerySessionInformation(pointer hServer, long SessionId,
//		long WTSInfoClass, POINTER* ppBuffer, LONG* pBytesReturned)
var wtsQuerySessionInformation = wtsapi32.MustFindProc("WTSQuerySessionInformationA").Addr()

const WTS_CURRENT_SERVER_HANDLE = 0
const WTS_CURRENT_SESSION = uintptrMinusOne
const WTS_ClientProtocolType = 16
const WTS_SessionId = 4

func WTS_GetClientProtocolType() int {
	defer heap.FreeTo(heap.CurSize())
	pbuf := heap.Alloc(uintptrSize)
	psize := heap.Alloc(int32Size)
	rtn := goc.Syscall5(wtsQuerySessionInformation,
		WTS_CURRENT_SERVER_HANDLE,
		WTS_CURRENT_SESSION,
		WTS_ClientProtocolType,
		uintptr(pbuf),
		uintptr(psize))
	buf := *(*uintptr)(pbuf)
	size := *(*int32)(psize)
	if rtn == 0 || size != 2 || buf == 0 {
		return 0
	}
	data := *(*int16)(unsafe.Pointer(buf))
	WTSFreeMemory(buf)
	return int(data)
}

var _ = builtin0("WTS_GetSessionId()",
	func() Value {
		if WTS_GetClientProtocolType() == 0 {
			return Zero
		}
		defer heap.FreeTo(heap.CurSize())
		pbuf := heap.Alloc(uintptrSize)
		psize := heap.Alloc(int32Size)
		rtn := goc.Syscall5(wtsQuerySessionInformation,
			WTS_CURRENT_SERVER_HANDLE,
			WTS_CURRENT_SESSION,
			WTS_SessionId,
			uintptr(pbuf),
			uintptr(psize))
		buf := *(*uintptr)(pbuf)
		size := *(*int32)(psize)
		if rtn == 0 || size != 4 || buf == 0 {
			return Zero
		}
		data := *(*int32)(unsafe.Pointer(buf))
		WTSFreeMemory(buf)
		return IntVal(int(data))
	})
