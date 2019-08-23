package builtin

import (
	"bytes"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var kernel32 = windows.NewLazyDLL("kernel32.dll")

type OSVERSIONINFOEX struct {
	dwOSVersionInfoSize int32
	dwMajorVersion      int32
	dwMinorVersion      int32
	dwBuildNumber       int32
	dwPlatformId        int32
	szCSDVersion        *byte
	wServicePackMajor   int
	wServicePackMinor   int
	wSuiteMask          int
	wProductType        *byte
	wReserved           *byte
}

// dll Kernel32:GetComputerName(buffer lpBuffer, LONG* lpnSize) bool
var getComputerName = kernel32.NewProc("GetComputerNameA")
var _ = builtin0("GetComputerName()", func() Value {
	const bufsize = 255
	var buf [bufsize + 1]byte
	n := int32(bufsize)
	getComputerName.Call(
		uintptr(unsafe.Pointer(&buf)),
		uintptr(unsafe.Pointer(&n)))
	return SuStr(string(buf[:n]))
})

// dll Kernel32:GetModuleHandle(instring name) pointer
var getModuleHandle = kernel32.NewProc("GetModuleHandleA")
var _ = builtin1("GetModuleHandle(unused)",
	func(a Value) Value {
		rtn, _, _ := getModuleHandle.Call(0)
		return intRet(rtn)
	})

// dll Kernel32:GetLocaleInfo(long locale, long lctype, string lpLCData, long cchData) long
var getLocaleInfo = kernel32.NewProc("GetLocaleInfoA")
var _ = builtin2("GetLocaleInfo(a,b)",
	func(a, b Value) Value {
		const bufsize = 255
		var buf [bufsize + 1]byte
		getLocaleInfo.Call(
			intArg(a),
			intArg(b),
			uintptr(unsafe.Pointer(&buf)),
			uintptr(bufsize))
		return SuStr(string(buf[:bytes.IndexByte(buf[:], 0)]))
	})

// dll Kernel32:GetProcAddress(pointer hModule, instring procName) pointer
var getProcAddress = kernel32.NewProc("GetProcAddress")
var _ = builtin2("GetProcAddress(a,b)",
	func(a, b Value) Value {
		rtn, _, _ := getProcAddress.Call(
			intArg(a),
			stringArg(b))
		return intRet(rtn)
	})

// dll Kernel32:GetProcessHeap() pointer
var getProcessHeap = kernel32.NewProc("GetProcessHeap")
var _ = builtin0("GetProcessHeap()",
	func() Value {
		rtn, _, _ := getProcessHeap.Call()
		return intRet(rtn)
	})

// dll Kernel32:GetVersionEx(OSVERSIONINFOEX* lpVersionInfo) bool
var getVersionEx = kernel32.NewProc("GetVersionEx")
var _ = builtin1("GetVersionEx(a)",
	func(a Value) Value {
		ovi := OSVERSIONINFOEX{
			dwOSVersionInfoSize: getInt32(a, "dwOSVersionInfoSize"),
			dwMajorVersion:      getInt32(a, "dwMajorVersion"),
			dwMinorVersion:      getInt32(a, "dwMinorVersion"),
			dwBuildNumber:       getInt32(a, "dwBuildNumber"),
			dwPlatformId:        getInt32(a, "dwPlatformId"),
			szCSDVersion:        getStr(a, "szCSDVersion"),
			wServicePackMajor:   getInt(a, "wServicePackMajor"),
			wServicePackMinor:   getInt(a, "wServicePackMinor"),
			wSuiteMask:          getInt(a, "wSuiteMask"),
			wProductType:        getStr(a, "wProductType"),
			wReserved:           getStr(a, "wReserved"),
		}
		rtn, _, _ := getVersionEx.Call(uintptr(unsafe.Pointer(&ovi)))
		return boolRet(rtn)
	})

// dll Kernel32:GlobalAlloc(long flags, long size) pointer
var globalAlloc = kernel32.NewProc("GlobalAlloc")
var _ = builtin2("GlobalAlloc(flags, size)",
	func(a, b Value) Value {
		rtn, _, _ := globalAlloc.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Kernel32:GlobalLock(pointer handle) pointer
var globalLock = kernel32.NewProc("GlobalLock")
var _ = builtin1("GlobalLock(hMem)",
	func(a Value) Value {
		rtn, _, _ := globalLock.Call(intArg(a))
		return intRet(rtn)
	})

// dll Kernel32:GlobalUnlock(pointer handle) void
var globalUnlock = kernel32.NewProc("GlobalUnlock")
var _ = builtin1("GlobalUnlock(hMem)",
	func(a Value) Value {
		rtn, _, _ := globalUnlock.Call(intArg(a))
		return intRet(rtn)
	})

// dll Kernel32:HeapAlloc(pointer hHeap, long dwFlags, long dwBytes) pointer
var heapAlloc = user32.NewProc("HeapAlloc")
var _ = builtin3("HeapAlloc(hHeap, dwFlags, dwBytes)",
	func(a, b, c Value) Value {
		rtn, _, _ := heapAlloc.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll Kernel32:HeapFree(pointer hHeap, long dwFlags, pointer lpMem) bool
var heapFree = user32.NewProc("HeapFree")
var _ = builtin3("HeapFree(hHeap, dwFlags, lpMem)",
	func(a, b, c Value) Value {
		rtn, _, _ := heapFree.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll Kernel32:MulDiv(long x, long y, long z) long
var _ = builtin3("MulDiv(x, y, z)",
	func(a, b, c Value) Value {
		return IntVal(int(int64(ToInt(a)) * int64(ToInt(b)) / int64(ToInt(c))))
	})

// dll Kernel32:CopyMemory(pointer destination, instring source,
// long length) void
var _ = builtin3("CopyMemory(destination, source, length)",
	func(a, b, c Value) Value {
		dst := uintptr(ToInt(a))
		src := ToStr(b)
		n := ToInt(c)
		for i := 0; i < n; i++ {
			*(*byte)(unsafe.Pointer(dst + uintptr(i))) = src[i]
		}
		return nil
	})
