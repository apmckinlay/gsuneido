package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var wtsapi32 = windows.NewLazyDLL("wtsapi32.dll")

// dll void WTSAPI32:WTSFreeMemory(pointer adr)
var wtsFreeMemory = wtsapi32.NewProc("WTSFreeMemory")
var _ = builtin1("WTSFreeMemory(adr)",
	func(a Value) Value {
		wtsFreeMemory.Call(
			intArg(a))
		return nil
	})

// dll bool WTSAPI32:WTSQuerySessionInformation(pointer hServer, long SessionId,
//		long WTSInfoClass, POINTER* ppBuffer, LONG* pBytesReturned)
var wtsQuerySessionInformation = wtsapi32.NewProc("WTSQuerySessionInformationA")
var _ = builtin5("WTSQuerySessionInformation(hServer, SessionId, WTSInfoClass,"+
	" ppBuffer, pBytesReturned)",
	func(a, b, c, d, e Value) Value {
		var ppBuffer uintptr
		var pBytesReturned int32
		rtn, _, _ := wtsQuerySessionInformation.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&ppBuffer)),
			uintptr(unsafe.Pointer(&pBytesReturned)))
		d.Put(nil, SuStr("x"), IntVal(int(ppBuffer)))
		e.Put(nil, SuStr("x"), IntVal(int(pBytesReturned)))
		return boolRet(rtn)
	})
