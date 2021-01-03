// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"bytes"
	"strings"
	"syscall"
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/runtime"
	reg "golang.org/x/sys/windows/registry"
)

var kernel32 = MustLoadDLL("kernel32.dll")

// NOTE: We want these functions to be available on secondary threads.
// Therefore we can't use goc.Syscall or heap.

// zbuf returns a zero terminated byte slice copy of ToStr(v)
func zbuf(v Value) []byte {
	s := ToStr(v)
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = 0
	return buf
}

// zstr returns an SuStr from up to the first zero byte in buf,
// or all of buf if there is no zero byte.
func zstr(buf []byte) SuStr {
	i := bytes.IndexByte(buf, 0)
	if i == -1 {
		i = len(buf)
	}
	return SuStr(string(buf[:i]))
}

// dll Kernel32:GetComputerName(buffer lpBuffer, LONG* lpnSize) bool
var getComputerName = kernel32.MustFindProc("GetComputerNameA").Addr()
var _ = builtin0("GetComputerName()", func() Value {
	const bufsize = 32
	buf := make([]byte, bufsize+1)
	n := uint32(bufsize)
	rtn, _, _ := syscall.Syscall(getComputerName, 2,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&n)),
		0)
	if rtn == 0 {
		return EmptyStr
	}
	return SuStr(string(buf[:n]))
})

// dll Kernel32:GetTempPath(DWORD nBufferLength, buffer lpBuffer) bool
var getTempPath = kernel32.MustFindProc("GetTempPathA").Addr()
var _ = builtin0("GetTempPath()", func() Value {
	buf := make([]byte, MAX_PATH+1)
	rtn, _, _ := syscall.Syscall(getTempPath, 2,
		MAX_PATH,
		uintptr(unsafe.Pointer(&buf[0])),
		0)
	s := string(buf[:rtn])
	return SuStr(strings.ReplaceAll(s, "\\", "/"))
})

// dll Kernel32:GetModuleHandle(instring name) pointer
var getModuleHandle = kernel32.MustFindProc("GetModuleHandleA").Addr()
var _ = builtin1("GetModuleHandle(unused)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(getModuleHandle, 1,
			0,
			0, 0)
		return intRet(rtn)
	})

// dll Kernel32:GetLocaleInfo(long locale, long lctype, string lpLCData, long cchData) long
var getLocaleInfo = kernel32.MustFindProc("GetLocaleInfoA").Addr()
var _ = builtin2("GetLocaleInfo(a,b)",
	func(a, b Value) Value {
		const bufsize = 255
		buf := make([]byte, bufsize+1)
		rtn, _, _ := syscall.Syscall6(getLocaleInfo, 4,
			intArg(a),
			intArg(b),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(bufsize),
			0, 0)
		return SuStr(string(buf[:rtn-1]))
	})

// dll Kernel32:GetProcAddress(pointer hModule, [in] string procName) pointer
var getProcAddress = kernel32.MustFindProc("GetProcAddress").Addr()
var _ = builtin2("GetProcAddress(hModule, procName)",
	func(a, b Value) Value {
		procName := zbuf(b)
		rtn, _, _ := syscall.Syscall(getProcAddress, 2,
			intArg(a),
			uintptr(unsafe.Pointer(&procName[0])),
			0)
		return intRet(rtn)
	})

// dll Kernel32:GetProcessHeap() pointer
var getProcessHeap = kernel32.MustFindProc("GetProcessHeap").Addr()
var _ = builtin0("GetProcessHeap()",
	func() Value {
		rtn, _, _ := syscall.Syscall(getProcessHeap, 0, 0, 0, 0)
		return intRet(rtn)
	})

// dll Kernel32:GlobalAlloc(long flags, long size) pointer
var globalAlloc = kernel32.MustFindProc("GlobalAlloc").Addr()
var _ = builtin2("GlobalAlloc(flags, size)",
	func(a, b Value) Value {
		rtn, _, _ := syscall.Syscall(globalAlloc, 2,
			intArg(a),
			intArg(b),
			0)
		return intRet(rtn)
	})

// dll Kernel32:GlobalLock(pointer handle) pointer
var globalLock = kernel32.MustFindProc("GlobalLock").Addr()
var _ = builtin1("GlobalLock(hMem)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(globalLock, 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

var globalSize = kernel32.MustFindProc("GlobalSize").Addr()
var _ = builtin1("GlobalLockString(hMem)",
	func(a Value) Value {
		// NOTE: assumes string takes up entire globalSize
		n, _, _ := syscall.Syscall(globalSize, 1,
			intArg(a),
			0, 0)
		rtn, _, _ := syscall.Syscall(globalLock, 1,
			intArg(a),
			0, 0)
		return bufStrZ(unsafe.Pointer(rtn), n)
	})

// dll Kernel32:GlobalUnlock(pointer handle) bool
var globalUnlock = kernel32.MustFindProc("GlobalUnlock").Addr()
var _ = builtin1("GlobalUnlock(hMem)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(globalUnlock, 1,
			intArg(a),
			0, 0)
		return boolRet(rtn)
	})

// dll pointer Kernel32:GlobalFree(pointer hglb)
var globalFree = kernel32.MustFindProc("GlobalFree").Addr()
var _ = builtin1("GlobalFree(hglb)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(globalFree, 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll Kernel32:HeapAlloc(pointer hHeap, long dwFlags, long dwBytes) pointer
var heapAlloc = kernel32.MustFindProc("HeapAlloc").Addr()
var _ = builtin3("HeapAlloc(hHeap, dwFlags, dwBytes)",
	func(a, b, c Value) Value {
		rtn, _, _ := syscall.Syscall(heapAlloc, 3,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll Kernel32:HeapFree(pointer hHeap, long dwFlags, pointer lpMem) bool
var heapFree = kernel32.MustFindProc("HeapFree").Addr()
var _ = builtin3("HeapFree(hHeap, dwFlags, lpMem)",
	func(a, b, c Value) Value {
		rtn, _, _ := syscall.Syscall(heapFree, 3,
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

var _ = builtin3("MulDiv(x, y, z)",
	func(a, b, c Value) Value {
		return IntVal(int(int64(ToInt(a)) * int64(ToInt(b)) / int64(ToInt(c))))
	})

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

// dll bool Kernel32:CloseHandle(pointer handle)
var closeHandle = kernel32.MustFindProc("CloseHandle").Addr()
var _ = builtin1("CloseHandle(handle)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(closeHandle, 1,
			intArg(a),
			0, 0)
		return boolRet(rtn)
	})

// dll bool Kernel32:CopyFile(string from, string to, bool failIfExists)
var copyFile = kernel32.MustFindProc("CopyFileA").Addr()
var _ = builtin3("CopyFile(from, to, failIfExists)",
	func(a, b, c Value) Value {
		from := zbuf(a)
		to := zbuf(b)
		rtn, _, _ := syscall.Syscall(copyFile, 3,
			uintptr(unsafe.Pointer(&from[0])),
			uintptr(unsafe.Pointer(&to[0])),
			boolArg(c))
		return boolRet(rtn)
	})

// dll pointer Kernel32:GetCurrentProcess()
var getCurrentProcess = kernel32.MustFindProc("GetCurrentProcess").Addr()
var _ = builtin0("GetCurrentProcess()",
	func() Value {
		rtn, _, _ := syscall.Syscall(getCurrentProcess, 0, 0, 0, 0)
		return intRet(rtn)
	})

// dll long Kernel32:GetCurrentProcessId()
var getCurrentProcessId = kernel32.MustFindProc("GetCurrentProcessId").Addr()
var _ = builtin0("GetCurrentProcessId()",
	func() Value {
		rtn, _, _ := syscall.Syscall(getCurrentProcessId, 0, 0, 0, 0)
		return intRet(rtn)
	})

// dll long Kernel32:GetCurrentThreadId()
var getCurrentThreadId = kernel32.MustFindProc("GetCurrentThreadId").Addr()
var _ = builtin0("GetCurrentThreadId()",
	func() Value {
		// NOTE: always returns cside thread id
		return intRet(goc.CThreadId()) // thread safe
	})

// dll long Kernel32:GetFileAttributes([in] string lpFileName)
var getFileAttributes = kernel32.MustFindProc("GetFileAttributesA").Addr()
var _ = builtin1("GetFileAttributes(lpFileName)",
	func(a Value) Value {
		name := zbuf(a)
		rtn, _, _ := syscall.Syscall(getFileAttributes, 1,
			uintptr(unsafe.Pointer(&name[0])),
			0, 0)
		return intRet(rtn)
	})

// dll pointer Kernel32:GetStdHandle(long nStdHandle)
var getStdHandle = kernel32.MustFindProc("GetStdHandle").Addr()
var _ = builtin1("GetStdHandle(nStdHandle)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(getStdHandle, 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll int64 Kernel32:GetTickCount64()
var getTickCount64 = kernel32.MustFindProc("GetTickCount64").Addr()
var _ = builtin0("GetTickCount()",
	func() Value {
		rtn, _, _ := syscall.Syscall(getTickCount64, 0, 0, 0, 0)
		return intRet(rtn)
	})

// dll pointer Kernel32:LoadLibrary([in] string library)
var loadLibrary = kernel32.MustFindProc("LoadLibraryA").Addr()
var _ = builtin1("LoadLibrary(library)",
	func(a Value) Value {
		lib := zbuf(a)
		rtn, _, _ := syscall.Syscall(loadLibrary, 1,
			uintptr(unsafe.Pointer(&lib[0])),
			0, 0)
		return intRet(rtn)
	})

// dll pointer Kernel32:LoadResource(pointer module, pointer res)
var loadResource = kernel32.MustFindProc("LoadResource").Addr()
var _ = builtin2("LoadResource(module, res)",
	func(a, b Value) Value {
		rtn, _, _ := syscall.Syscall(loadResource, 2,
			intArg(a),
			intArg(b),
			0)
		return intRet(rtn)
	})

// dll bool Kernel32:SetCurrentDirectory(string lpPathName)
var setCurrentDirectory = kernel32.MustFindProc("SetCurrentDirectoryA").Addr()
var _ = builtin1("SetCurrentDirectory(lpPathName)",
	func(a Value) Value {
		name := zbuf(a)
		rtn, _, _ := syscall.Syscall(setCurrentDirectory, 1,
			uintptr(unsafe.Pointer(&name[0])),
			0, 0)
		return boolRet(rtn)
	})

// dll bool Kernel32:SetFileAttributes(
//		[in] string lpFileName, long dwFileAttributes)
var setFileAttributes = kernel32.MustFindProc("SetFileAttributesA").Addr()
var _ = builtin2("SetFileAttributes(lpFileName, dwFileAttributes)",
	func(a, b Value) Value {
		name := zbuf(a)
		rtn, _, _ := syscall.Syscall(setFileAttributes, 2,
			uintptr(unsafe.Pointer(&name[0])),
			intArg(b),
			0)
		return boolRet(rtn)
	})

// dll handle Kernel32:CreateFile([in] string lpFileName, long dwDesiredAccess,
//		long dwShareMode, SECURITY_ATTRIBUTES* lpSecurityAttributes,
//		long dwCreationDistribution, long dwFlagsAndAttributes,
//		pointer hTemplateFile)
var createFile = kernel32.MustFindProc("CreateFileA").Addr()
var _ = builtin7("CreateFile(lpFileName, dwDesiredAccess, dwShareMode,"+
	"lpSecurityAttributes, dwCreationDistribution, dwFlagsAndAttributes,"+
	"hTemplateFile)",
	func(a, b, c, d, e, f, g Value) Value {
		name := zbuf(a)
		rtn, _, _ := syscall.Syscall9(createFile, 7,
			uintptr(unsafe.Pointer(&name[0])),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f),
			intArg(g),
			0, 0)
		return intRet(rtn)
	})

// dll bool Kernel32:WriteFile(
//		handle hFile,
//		buffer lpBuffer,
//		long nNumberOfBytesToWrite,
//		LONG* lpNumberOfBytesWritten,
//		pointer lpOverlapped)
var writeFile = kernel32.MustFindProc("WriteFile").Addr()
var _ = builtin5("WriteFile(hFile, lpBuffer, nNumberOfBytesToWrite, "+
	"lpNumberOfBytesWritten, lpOverlapped/*unused*/)",
	func(a, b, c, d, e Value) Value {
		s := ToStr(b)
		n := ToInt(c)
		buf := ([]byte)(s[:n]) // n <= len(s)
		return WriteFile(a, unsafe.Pointer(&buf[0]), c, d)
	})

var _ = builtin5("WriteFilePtr(hFile, lpBuffer, nNumberOfBytesToWrite, "+
	"lpNumberOfBytesWritten, lpOverlapped/*unused*/)",
	func(a, b, c, d, e Value) Value {
		buf := unsafe.Pointer(uintptr(ToInt(b)))
		return WriteFile(a, buf, c, d)
	})

func WriteFile(f Value, buf unsafe.Pointer, size Value, written Value) Value {
	n := ToInt(size)
	if n == 0 {
		return False
	}
	var w int32
	rtn, _, _ := syscall.Syscall6(writeFile, 5,
		intArg(f),
		uintptr(buf),
		uintptr(n),
		uintptr(unsafe.Pointer(&w)),
		0,
		0)
	written.Put(nil, SuStr("x"), IntVal(int(w)))
	return boolRet(rtn)
}

// dll bool Kernel32:GetVolumeInformation([in] string lpRootPathName,
//		string lpVolumeNameBuffer, long nVolumeNameSize, LONG* lpVolumeSerialNumber,
//		LONG* lpMaximumComponentLength, LONG* lpFileSystemFlags,
//		string lpFileSystemNameBuffer, long nFileSystemNameSize)
var getVolumeInformation = kernel32.MustFindProc("GetVolumeInformationA").Addr()
var _ = builtin1("GetVolumeName(vol = 'c:\\\\')",
	func(a Value) Value {
		vol := zbuf(a)
		const bufsize = 255
		buf := make([]byte, bufsize+1)
		rtn, _, _ := syscall.Syscall9(getVolumeInformation, 8,
			uintptr(unsafe.Pointer(&vol[0])),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(bufsize),
			0,
			0,
			0,
			0,
			0,
			0)
		if rtn == 0 {
			return EmptyStr
		}
		return zstr(buf)
	})

type MEMORYSTATUSEX struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	unused       [6]uint64
}

const nMEMORYSTATUSEX = unsafe.Sizeof(MEMORYSTATUSEX{})

var globalMemoryStatusEx = kernel32.MustFindProc("GlobalMemoryStatusEx").Addr()

var _ = builtin0("SystemMemory()", func() Value {
	buf := make([]byte, nMEMORYSTATUSEX)
	(*MEMORYSTATUSEX)(unsafe.Pointer(&buf[0])).dwLength = uint32(nMEMORYSTATUSEX)
	rtn, _, _ := syscall.Syscall(globalMemoryStatusEx, 1,
		uintptr(unsafe.Pointer(&buf[0])),
		0, 0)
	if rtn == 0 {
		return Zero
	}
	return Int64Val(int64((*MEMORYSTATUSEX)(unsafe.Pointer(&buf[0])).ullTotalPhys))
})

// dll bool Kernel32:GetDiskFreeSpaceEx(
// 	[in] string			directoryName,
// 	ULARGE_INTEGER*		freeBytesAvailableToCaller,
// 	ULARGE_INTEGER*		totalNumberOfBytes,
// 	ULARGE_INTEGER*		totalNumberOfFreeBytes)
var getDiskFreeSpaceEx = kernel32.MustFindProc("GetDiskFreeSpaceExA").Addr()

var _ = builtin1("GetDiskFreeSpace(dir = '.')", func(arg Value) Value {
	dir := zbuf(arg)
	var n int64
	syscall.Syscall6(getDiskFreeSpaceEx, 4,
		uintptr(unsafe.Pointer(&dir[0])),
		uintptr(unsafe.Pointer(&n)),
		0,
		0,
		0, 0)
	return Int64Val(n)
})

//-------------------------------------------------------------------

var _ = builtin0("OperatingSystem()", OSName) // deprecated
var _ = builtin0("OSName()", OSName)

func OSName() Value {
	k, err := reg.OpenKey(reg.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, reg.QUERY_VALUE)
	if err != nil {
		return EmptyStr
	}
	defer k.Close()

	s, _, err := k.GetStringValue("ProductName")
	if err != nil {
		return EmptyStr
	}
	return SuStr(s)
}

//-------------------------------------------------------------------

var versionApi = MustLoadDLL("version.dll")

var getFileVersionInfo = versionApi.MustFindProc("GetFileVersionInfoA").Addr()
var getFileVersionInfoSize = versionApi.MustFindProc("GetFileVersionInfoSizeA").Addr()
var verQueryValue = versionApi.MustFindProc("VerQueryValueA").Addr()

var verFile = []byte("kernel32\x00")
var verFileW = []byte("\\\x00")

var _ = builtin0("OSVersion()", func() Value {
	var dummy int32
	size, _, _ := syscall.Syscall(getFileVersionInfoSize, 2,
		uintptr(unsafe.Pointer(&verFile[0])),
		uintptr(unsafe.Pointer(&dummy)),
		0)
	if size == 0 {
		return False
	}
	buf := make([]byte, size)
	rtn, _, _ := syscall.Syscall6(getFileVersionInfo, 4,
		uintptr(unsafe.Pointer(&verFile[0])),
		0,
		size,
		uintptr(unsafe.Pointer(&buf[0])),
		0, 0)
	if rtn == 0 {
		return False
	}
	var pffi *VS_FIXEDFILEINFO
	var len int32
	rtn, _, _ = syscall.Syscall6(verQueryValue, 4,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&verFileW[0])),
		uintptr(unsafe.Pointer(&pffi)),
		uintptr(unsafe.Pointer(&len)),
		0, 0)
	if rtn == 0 {
		return False
	}
	ob := &SuObject{}
	ob.Add(IntVal(hiword(pffi.dwFileVersionMS)))
	ob.Add(IntVal(loword(pffi.dwFileVersionMS)))
	ob.Add(IntVal(hiword(pffi.dwFileVersionLS)))
	ob.Add(IntVal(loword(pffi.dwFileVersionLS)))
	return ob
})

func loword(n int32) int {
	return int(n & 0xffff)
}

func hiword(n int32) int {
	return int(n >> 16)
}

type VS_FIXEDFILEINFO struct {
	dwSignature        int32
	dwStrucVersion     int32
	dwFileVersionMS    int32
	dwFileVersionLS    int32
	dwProductVersionMS int32
	dwProductVersionLS int32
	dwFileFlagsMask    int32
	dwFileFlags        int32
	dwFileOS           int32
	dwFileType         int32
	dwFileSubtype      int32
	dwFileDateMS       int32
	dwFileDateLS       int32
}

const nVS_FIXEDFILEINFO = unsafe.Sizeof(VS_FIXEDFILEINFO{})
