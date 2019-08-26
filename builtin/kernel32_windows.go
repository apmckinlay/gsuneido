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
			uintptr(stringArg(b)))
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

var _ = builtin1("GlobalUnlockString(hMem)",
	func(a Value) Value {
		rtn, _, _ := globalUnlock.Call(intArg(a))
		return SuStr(strFromAddr(rtn))
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

// dll Kernel32:CopyMemory(pointer destination, [in] string source,
//		long length) void
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

// dll bool Kernel32:AllocConsole()
var allocConsole = kernel32.NewProc("AllocConsole")
var _ = builtin0("AllocConsole()",
	func() Value {
		rtn, _, _ := allocConsole.Call()
		return boolRet(rtn)
	})

// dll bool Kernel32:CloseHandle(pointer handle)
var closeHandle = kernel32.NewProc("CloseHandle")
var _ = builtin1("CloseHandle(handle)",
	func(a Value) Value {
		rtn, _, _ := closeHandle.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool Kernel32:CopyFile(
//		[in] string from,
//		[in] string to,
//		bool failIfExists)
var copyFile = kernel32.NewProc("CopyFile")
var _ = builtin3("CopyFile(from, to, failIfExists)",
	func(a, b, c Value) Value {
		rtn, _, _ := copyFile.Call(
			uintptr(stringArg(a)),
			uintptr(stringArg(b)),
			boolArg(c))
		return boolRet(rtn)
	})

// dll bool Kernel32:DeleteFile([in] string filename)
var deleteFile = kernel32.NewProc("DeleteFile")
var _ = builtin1("DeleteFile(filename)",
	func(a Value) Value {
		rtn, _, _ := deleteFile.Call(
			uintptr(stringArg(a)))
		return boolRet(rtn)
	})

// dll bool Kernel32:FindClose(pointer hFindFile)
var findClose = kernel32.NewProc("FindClose")
var _ = builtin1("FindClose(hFindFile)",
	func(a Value) Value {
		rtn, _, _ := findClose.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool Kernel32:FlushFileBuffers(handle hFile)
var flushFileBuffers = kernel32.NewProc("FlushFileBuffers")
var _ = builtin1("FlushFileBuffers(hFile)",
	func(a Value) Value {
		rtn, _, _ := flushFileBuffers.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool Kernel32:FreeConsole()
var freeConsole = kernel32.NewProc("FreeConsole")
var _ = builtin0("FreeConsole()",
	func() Value {
		rtn, _, _ := freeConsole.Call()
		return boolRet(rtn)
	})

// dll pointer Kernel32:GetCurrentProcess()
var getCurrentProcess = kernel32.NewProc("GetCurrentProcess")
var _ = builtin0("GetCurrentProcess()",
	func() Value {
		rtn, _, _ := getCurrentProcess.Call()
		return intRet(rtn)
	})

// dll long Kernel32:GetCurrentProcessId()
var getCurrentProcessId = kernel32.NewProc("GetCurrentProcessId")
var _ = builtin0("GetCurrentProcessId()",
	func() Value {
		rtn, _, _ := getCurrentProcessId.Call()
		return intRet(rtn)
	})

// dll long Kernel32:GetCurrentThreadId()
var getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
var _ = builtin0("GetCurrentThreadId()",
	func() Value {
		rtn, _, _ := getCurrentThreadId.Call()
		return intRet(rtn)
	})

// dll long Kernel32:GetFileAttributes(
// 		[in] string lpFileName)
var getFileAttributes = kernel32.NewProc("GetFileAttributes")
var _ = builtin1("GetFileAttributes(lpFileName)",
	func(a Value) Value {
		rtn, _, _ := getFileAttributes.Call(
			uintptr(stringArg(a)))
		return intRet(rtn)
	})

// dll long Kernel32:GetLastError()
var getLastError = kernel32.NewProc("GetLastError")
var _ = builtin0("GetLastError()",
	func() Value {
		rtn, _, _ := getLastError.Call()
		return intRet(rtn)
	})

// dll pointer Kernel32:GetStdHandle(long nStdHandle)
var getStdHandle = kernel32.NewProc("GetStdHandle")
var _ = builtin1("GetStdHandle(nStdHandle)",
	func(a Value) Value {
		rtn, _, _ := getStdHandle.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll int64 Kernel32:GetTickCount64()
var getTickCount64 = kernel32.NewProc("GetTickCount64")
var _ = builtin0("GetTickCount64()",
	func() Value {
		rtn, _, _ := getTickCount64.Call()
		return intRet(rtn)
	})

// dll long Kernel32:GetWindowsDirectory(string lpBuffer, long size)
var getWindowsDirectory = kernel32.NewProc("GetWindowsDirectory")
var _ = builtin2("GetWindowsDirectory(lpBuffer, size)",
	func(a, b Value) Value {
		rtn, _, _ := getWindowsDirectory.Call(
			uintptr(stringArg(a)),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer Kernel32:GlobalFree(pointer hglb)
var globalFree = kernel32.NewProc("GlobalFree")
var _ = builtin1("GlobalFree(hglb)",
	func(a Value) Value {
		rtn, _, _ := globalFree.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll long Kernel32:GlobalSize(pointer handle)
var globalSize = kernel32.NewProc("GlobalSize")
var _ = builtin1("GlobalSize(handle)",
	func(a Value) Value {
		rtn, _, _ := globalSize.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll pointer Kernel32:LoadLibrary([in] string library)
var loadLibrary = kernel32.NewProc("LoadLibrary")
var _ = builtin1("LoadLibrary(library)",
	func(a Value) Value {
		rtn, _, _ := loadLibrary.Call(
			uintptr(stringArg(a)))
		return intRet(rtn)
	})

// dll pointer Kernel32:LoadResource(pointer module, pointer res)
var loadResource = kernel32.NewProc("LoadResource")
var _ = builtin2("LoadResource(module, res)",
	func(a, b Value) Value {
		rtn, _, _ := loadResource.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll bool Kernel32:MoveFileEx([in] string from, [in] string to, long flags)
var moveFileEx = kernel32.NewProc("MoveFileEx")
var _ = builtin3("MoveFileEx(from, to, flags)",
	func(a, b, c Value) Value {
		rtn, _, _ := moveFileEx.Call(
			uintptr(stringArg(a)),
			uintptr(stringArg(b)),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool Kernel32:RemoveDirectory([in] string lpPathName)
var removeDirectory = kernel32.NewProc("RemoveDirectory")
var _ = builtin1("RemoveDirectory(lpPathName)",
	func(a Value) Value {
		rtn, _, _ := removeDirectory.Call(
			uintptr(stringArg(a)))
		return boolRet(rtn)
	})

// dll bool Kernel32:SetConsoleTitle(string title)
var setConsoleTitle = kernel32.NewProc("SetConsoleTitle")
var _ = builtin1("SetConsoleTitle(title)",
	func(a Value) Value {
		rtn, _, _ := setConsoleTitle.Call(
			uintptr(stringArg(a)))
		return boolRet(rtn)
	})

// dll bool Kernel32:SetCurrentDirectory(string lpPathName)
var setCurrentDirectory = kernel32.NewProc("SetCurrentDirectory")
var _ = builtin1("SetCurrentDirectory(lpPathName)",
	func(a Value) Value {
		rtn, _, _ := setCurrentDirectory.Call(
			uintptr(stringArg(a)))
		return boolRet(rtn)
	})

// dll bool Kernel32:SetFileAttributes(
//		[in] string lpFileName, long dwFileAttributes)
var setFileAttributes = kernel32.NewProc("SetFileAttributes")
var _ = builtin2("SetFileAttributes(lpFileName, dwFileAttributes)",
	func(a, b Value) Value {
		rtn, _, _ := setFileAttributes.Call(
			uintptr(stringArg(a)),
			intArg(b))
		return boolRet(rtn)
	})
