// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
)

var wtsapi32 = MustLoadDLL("wtsapi32.dll")

// dll void WTSAPI32:WTSFreeMemory(pointer adr)

var wtsFreeMemory = wtsapi32.MustFindProc("WTSFreeMemory").Addr()

func WTSFreeMemory(adr uintptr) {
	syscall.SyscallN(wtsFreeMemory,
		adr)
}

// dll bool WTSAPI32:WTSQuerySessionInformation(pointer hServer, long SessionId,
// long WTSInfoClass, POINTER* ppBuffer, LONG* pBytesReturned)
var wtsQuerySessionInformation = wtsapi32.MustFindProc("WTSQuerySessionInformationA").Addr()

const WTS_CURRENT_SERVER_HANDLE = 0
const WTS_CURRENT_SESSION = uintptrMinusOne
const WTS_ClientProtocolType = 16
const WTS_SessionId = 4

func WTS_GetClientProtocolType() int {
	var buf uintptr
	var size int32
	rtn, _, _ := syscall.SyscallN(wtsQuerySessionInformation,
		WTS_CURRENT_SERVER_HANDLE,
		WTS_CURRENT_SESSION,
		WTS_ClientProtocolType,
		uintptr(unsafe.Pointer(&buf)),
		uintptr(unsafe.Pointer(&size)))
	if rtn == 0 || size != 2 || buf == 0 {
		return 0
	}
	data := *(*int16)(unsafe.Pointer(buf))
	WTSFreeMemory(buf)
	return int(data)
}

var _ = builtin(WTS_GetSessionId, "()")

func WTS_GetSessionId() Value {
	if WTS_GetClientProtocolType() == 0 {
		return Zero
	}
	var buf uintptr
	var size int32
	rtn, _, _ := syscall.SyscallN(wtsQuerySessionInformation,
		WTS_CURRENT_SERVER_HANDLE,
		WTS_CURRENT_SESSION,
		WTS_SessionId,
		uintptr(unsafe.Pointer(&buf)),
		uintptr(unsafe.Pointer(&size)))
	if rtn == 0 || size != 4 || buf == 0 {
		return Zero
	}
	data := *(*int32)(unsafe.Pointer(buf))
	WTSFreeMemory(buf)
	return IntVal(int(data))
}
