// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"bytes"
	"syscall"
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
	reg "golang.org/x/sys/windows/registry"
)

// NOTE: We want these functions to be available on secondary threads.
// Therefore we can't use goc.Syscall or heap.

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
var _ = builtin(GetComputerName, "()")

func GetComputerName() Value {
	const bufsize = 32
	buf := make([]byte, bufsize+1)
	n := uint32(bufsize)
	rtn, _, _ := syscall.SyscallN(getComputerName,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&n)))
	if rtn == 0 {
		return EmptyStr
	}
	return SuStr(string(buf[:n]))
}

// dll Kernel32:GetModuleHandle(instring name) pointer
var getModuleHandle = kernel32.MustFindProc("GetModuleHandleA").Addr()
var _ = builtin(GetModuleHandle, "(unused)")

func GetModuleHandle(a Value) Value {
	rtn, _, _ := syscall.SyscallN(getModuleHandle, 0)
	return intRet(rtn)
}

// dll Kernel32:GetLocaleInfo(long locale, long lctype, string lpLCData, long cchData) long
var getLocaleInfo = kernel32.MustFindProc("GetLocaleInfoA").Addr()
var _ = builtin(GetLocaleInfo, "(a,b)")

func GetLocaleInfo(a, b Value) Value {
	const bufsize = 255
	buf := make([]byte, bufsize+1)
	rtn, _, _ := syscall.SyscallN(getLocaleInfo,
		intArg(a),
		intArg(b),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(bufsize))
	return SuStr(string(buf[:rtn-1]))
}

// dll Kernel32:GetProcAddress(pointer hModule, [in] string procName) pointer
var getProcAddress = kernel32.MustFindProc("GetProcAddress").Addr()
var _ = builtin(GetProcAddress, "(hModule, procName)")

func GetProcAddress(a, b Value) Value {
	procName := zbuf(b)
	rtn, _, _ := syscall.SyscallN(getProcAddress,
		intArg(a),
		uintptr(unsafe.Pointer(&procName[0])))
	return intRet(rtn)
}

// dll Kernel32:GetProcessHeap() pointer
var getProcessHeap = kernel32.MustFindProc("GetProcessHeap").Addr()
var _ = builtin(GetProcessHeap, "()")

func GetProcessHeap() Value {
	rtn, _, _ := syscall.SyscallN(getProcessHeap)
	return intRet(rtn)
}

// Global -----------------------------------------------------------

// dll Kernel32:GlobalAlloc(long flags, long size) pointer
var globalAlloc = kernel32.MustFindProc("GlobalAlloc").Addr()

func globalalloc(flags, n uintptr) HANDLE {
	rtn, _, _ := syscall.SyscallN(globalAlloc, flags, n)
	return rtn
}

var _ = builtin(GlobalAlloc, "(flags, size)")

func GlobalAlloc(a, b Value) Value {
	return intRet(globalalloc(intArg(a), intArg(b)))
}

// dll Kernel32:GlobalLock(pointer handle) pointer
var globalLock = kernel32.MustFindProc("GlobalLock").Addr()

func globallock(handle HANDLE) unsafe.Pointer {
	rtn, _, _ := syscall.SyscallN(globalLock, handle)
	return unsafe.Pointer(rtn)
}

var _ = builtin(GlobalLock, "(hMem)")

func GlobalLock(a Value) Value {
	return intRet(uintptr(globallock(intArg(a))))
}

// dll Kernel32:GlobalSize(pointer handle) bool
var globalSize = kernel32.MustFindProc("GlobalSize").Addr()

func globalsize(handle HANDLE) uintptr {
	rtn, _, _ := syscall.SyscallN(globalSize, handle)
	return rtn
}

var _ = builtin(GlobalSize, "(hMem)")

func GlobalSize(a Value) Value {
	return intRet(globalsize(intArg(a)))
}

const GMEM_MOVEABLE = 2

var _ = builtin(GlobalAllocData, "(s)")

func GlobalAllocData(a Value) Value {
	s := ToStr(a)
	handle := globalalloc(GMEM_MOVEABLE, uintptr(len(s)))
	if len(s) > 0 {
		p := globallock(handle)
		assert.That(p != nil)
		defer globalunlock(handle)
		bufToPtr(s, p)
	}
	return intRet(handle) // caller must GlobalFree
}

var _ = builtin(GlobalAllocString, "(s)")

func GlobalAllocString(a Value) Value {
	s := ToStr(a)
	s = str.BeforeFirst(s, "\x00")
	handle := globalalloc(GMEM_MOVEABLE, uintptr(len(s))+1)
	p := globallock(handle)
	assert.That(p != nil)
	defer globalunlock(handle)
	strToPtr(s, p)
	return intRet(handle) // caller must GlobalFree
}

var _ = builtin(GlobalData, "(hMem)")

func GlobalData(a Value) Value {
	hm := intArg(a)
	n := globalsize(hm)
	if n == 0 {
		return EmptyStr
	}
	p := globallock(hm)
	assert.That(p != nil)
	defer globalunlock(hm)
	return bufStrN(p, n)
}

// dll Kernel32:GlobalUnlock(pointer handle) bool
var globalUnlock = kernel32.MustFindProc("GlobalUnlock").Addr()

func globalunlock(handle HANDLE) uintptr {
	rtn, _, _ := syscall.SyscallN(globalUnlock, handle)
	return rtn
}

var _ = builtin(GlobalUnlock, "(hMem)")

func GlobalUnlock(a Value) Value {
	return boolRet(globalunlock(intArg(a)))
}

// dll pointer Kernel32:GlobalFree(pointer hglb)
var globalFree = kernel32.MustFindProc("GlobalFree").Addr()
var _ = builtin(GlobalFree, "(hglb)")

func GlobalFree(a Value) Value {
	rtn, _, _ := syscall.SyscallN(globalFree, intArg(a))
	return intRet(rtn)
}

//-------------------------------------------------------------------

// dll Kernel32:HeapAlloc(pointer hHeap, long dwFlags, long dwBytes) pointer
var heapAlloc = kernel32.MustFindProc("HeapAlloc").Addr()
var _ = builtin(HeapAlloc, "(hHeap, dwFlags, dwBytes)")

func HeapAlloc(a, b, c Value) Value {
	rtn, _, _ := syscall.SyscallN(heapAlloc,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll Kernel32:HeapFree(pointer hHeap, long dwFlags, pointer lpMem) bool
var heapFree = kernel32.MustFindProc("HeapFree").Addr()
var _ = builtin(HeapFree, "(hHeap, dwFlags, lpMem)")

func HeapFree(a, b, c Value) Value {
	rtn, _, _ := syscall.SyscallN(heapFree,
		intArg(a),
		intArg(b),
		intArg(c))
	return boolRet(rtn)
}

var _ = builtin(MulDiv, "(x, y, z)")

func MulDiv(a, b, c Value) Value {
	return IntVal(int(int64(ToInt(a)) * int64(ToInt(b)) / int64(ToInt(c))))
}

var _ = builtin(CopyMemory, "(destination, source, length)")

func CopyMemory(a, b, c Value) Value {
	dst := uintptr(ToInt(a))
	src := ToStr(b)
	n := ToInt(c)
	for i := 0; i < n; i++ {
		*(*byte)(unsafe.Pointer(dst + uintptr(i))) = src[i]
	}
	return nil
}

// dll bool Kernel32:CloseHandle(pointer handle)
var closeHandle = kernel32.MustFindProc("CloseHandle").Addr()
var _ = builtin(CloseHandle, "(handle)")

func CloseHandle(a Value) Value {
	rtn, _, _ := syscall.SyscallN(closeHandle, intArg(a))
	return boolRet(rtn)
}

// dll pointer Kernel32:GetCurrentProcess()
var getCurrentProcess = kernel32.MustFindProc("GetCurrentProcess").Addr()
var _ = builtin(GetCurrentProcess, "()")

func GetCurrentProcess() Value {
	rtn, _, _ := syscall.SyscallN(getCurrentProcess)
	return intRet(rtn)
}

// dll long Kernel32:GetCurrentProcessId()
var getCurrentProcessId = kernel32.MustFindProc("GetCurrentProcessId").Addr()
var _ = builtin(GetCurrentProcessId, "()")

func GetCurrentProcessId() Value {
	rtn, _, _ := syscall.SyscallN(getCurrentProcessId)
	return intRet(rtn)
}

var _ = builtin(GetCurrentThreadId, "()")

func GetCurrentThreadId() Value {
	// NOTE: always returns cside thread id
	return intRet(goc.CThreadId()) // thread safe
}

// dll long Kernel32:GetFileAttributes([in] string lpFileName)
var getFileAttributes = kernel32.MustFindProc("GetFileAttributesA").Addr()
var _ = builtin(GetFileAttributes, "(lpFileName)")

func GetFileAttributes(a Value) Value {
	name := zbuf(a)
	rtn, _, _ := syscall.SyscallN(getFileAttributes,
		uintptr(unsafe.Pointer(&name[0])))
	return intRet(rtn)
}

// dll pointer Kernel32:GetStdHandle(long nStdHandle)
var getStdHandle = kernel32.MustFindProc("GetStdHandle").Addr()
var _ = builtin(GetStdHandle, "(nStdHandle)")

func GetStdHandle(a Value) Value {
	rtn, _, _ := syscall.SyscallN(getStdHandle, intArg(a))
	return intRet(rtn)
}

// dll int64 Kernel32:GetTickCount64()
var getTickCount64 = kernel32.MustFindProc("GetTickCount64").Addr()
var _ = builtin(GetTickCount, "()")

func GetTickCount() Value {
	rtn, _, _ := syscall.SyscallN(getTickCount64)
	return intRet(rtn)
}

// dll pointer Kernel32:LoadLibrary([in] string library)
var loadLibrary = kernel32.MustFindProc("LoadLibraryA").Addr()
var _ = builtin(LoadLibrary, "(library)")

func LoadLibrary(a Value) Value {
	lib := zbuf(a)
	rtn, _, _ := syscall.SyscallN(loadLibrary,
		uintptr(unsafe.Pointer(&lib[0])))
	return intRet(rtn)
}

// dll pointer Kernel32:LoadResource(pointer module, pointer res)
var loadResource = kernel32.MustFindProc("LoadResource").Addr()
var _ = builtin(LoadResource, "(module, res)")

func LoadResource(a, b Value) Value {
	rtn, _, _ := syscall.SyscallN(loadResource,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll bool Kernel32:SetCurrentDirectory(string lpPathName)
var setCurrentDirectory = kernel32.MustFindProc("SetCurrentDirectoryA").Addr()
var _ = builtin(SetCurrentDirectory, "(lpPathName)")

func SetCurrentDirectory(a Value) Value {
	name := zbuf(a)
	rtn, _, _ := syscall.SyscallN(setCurrentDirectory,
		uintptr(unsafe.Pointer(&name[0])))
	return boolRet(rtn)
}

// dll bool Kernel32:SetFileAttributes(
// [in] string lpFileName, long dwFileAttributes)
var setFileAttributes = kernel32.MustFindProc("SetFileAttributesA").Addr()
var _ = builtin(SetFileAttributes, "(lpFileName, dwFileAttributes)")

func SetFileAttributes(a, b Value) Value {
	name := zbuf(a)
	rtn, _, _ := syscall.SyscallN(setFileAttributes,
		uintptr(unsafe.Pointer(&name[0])),
		intArg(b))
	return boolRet(rtn)
}

// dll handle Kernel32:CreateFile([in] string lpFileName, long dwDesiredAccess,
// long dwShareMode, SECURITY_ATTRIBUTES* lpSecurityAttributes,
// long dwCreationDistribution, long dwFlagsAndAttributes,
// pointer hTemplateFile)
var createFile = kernel32.MustFindProc("CreateFileA").Addr()
var _ = builtin(CreateFile, "(lpFileName, dwDesiredAccess, dwShareMode, "+
	"lpSecurityAttributes, dwCreationDistribution, dwFlagsAndAttributes, "+
	"hTemplateFile)")

func CreateFile(a, b, c, d, e, f, g Value) Value {
	name := zbuf(a)
	rtn, _, _ := syscall.SyscallN(createFile,
		uintptr(unsafe.Pointer(&name[0])),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e),
		intArg(f),
		intArg(g))
	return intRet(rtn)
}

// dll bool Kernel32:WriteFile(
// handle hFile,
// buffer lpBuffer,
// long nNumberOfBytesToWrite,
// LONG* lpNumberOfBytesWritten,
// pointer lpOverlapped)
var writeFile = kernel32.MustFindProc("WriteFile").Addr()
var _ = builtin(WriteFile, "(hFile, lpBuffer, nNumberOfBytesToWrite, "+
	"lpNumberOfBytesWritten, lpOverlapped/*unused*/)")

func WriteFile(a, b, c, d, e Value) Value {
	s := ToStr(b)
	n := ToInt(c)
	buf := ([]byte)(s[:n]) // n <= len(s)
	return writefile(a, unsafe.Pointer(&buf[0]), c, d)
}

var _ = builtin(WriteFilePtr, "(hFile, lpBuffer, nNumberOfBytesToWrite, "+
	"lpNumberOfBytesWritten, lpOverlapped/*unused*/)")

func WriteFilePtr(a, b, c, d, e Value) Value {
	buf := unsafe.Pointer(uintptr(ToInt(b)))
	return writefile(a, buf, c, d)
}

func writefile(f Value, buf unsafe.Pointer, size Value, written Value) Value {
	n := ToInt(size)
	if n == 0 {
		return False
	}
	var w int32
	rtn, _, _ := syscall.SyscallN(writeFile,
		intArg(f),
		uintptr(buf),
		uintptr(n),
		uintptr(unsafe.Pointer(&w)),
		0)
	written.Put(nil, SuStr("x"), IntVal(int(w)))
	return boolRet(rtn)
}

// dll bool Kernel32:GetVolumeInformation([in] string lpRootPathName,
// string lpVolumeNameBuffer, long nVolumeNameSize, LONG* lpVolumeSerialNumber,
// LONG* lpMaximumComponentLength, LONG* lpFileSystemFlags,
// string lpFileSystemNameBuffer, long nFileSystemNameSize)
var getVolumeInformation = kernel32.MustFindProc("GetVolumeInformationA").Addr()
var _ = builtin(GetVolumeName, "(vol = `c:\\`)")

func GetVolumeName(a Value) Value {
	vol := zbuf(a)
	const bufsize = 255
	buf := make([]byte, bufsize+1)
	rtn, _, _ := syscall.SyscallN(getVolumeInformation,
		uintptr(unsafe.Pointer(&vol[0])),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(bufsize),
		0,
		0,
		0,
		0,
		0)
	if rtn == 0 {
		return EmptyStr
	}
	return zstr(buf)
}

//-------------------------------------------------------------------

var _ = builtin(OSName, "()")

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

var _ = builtin(OSVersion, "()")

func OSVersion() Value {
	var dummy int32
	size, _, _ := syscall.SyscallN(getFileVersionInfoSize,
		uintptr(unsafe.Pointer(&verFile[0])),
		uintptr(unsafe.Pointer(&dummy)))
	if size == 0 {
		return False
	}
	buf := make([]byte, size)
	rtn, _, _ := syscall.SyscallN(getFileVersionInfo,
		uintptr(unsafe.Pointer(&verFile[0])),
		0,
		size,
		uintptr(unsafe.Pointer(&buf[0])))
	if rtn == 0 {
		return False
	}
	var pffi *stVsFixedFileInfo
	var len int32
	rtn, _, _ = syscall.SyscallN(verQueryValue,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&verFileW[0])),
		uintptr(unsafe.Pointer(&pffi)),
		uintptr(unsafe.Pointer(&len)))
	if rtn == 0 {
		return False
	}
	ob := &SuObject{}
	ob.Add(IntVal(hiword(pffi.dwFileVersionMS)))
	ob.Add(IntVal(loword(pffi.dwFileVersionMS)))
	ob.Add(IntVal(hiword(pffi.dwFileVersionLS)))
	ob.Add(IntVal(loword(pffi.dwFileVersionLS)))
	return ob
}

func loword(n int32) int {
	return int(n & 0xffff)
}

func hiword(n int32) int {
	return int(n >> 16)
}

//lint:ignore U1000 Windows struct
type stVsFixedFileInfo struct {
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
