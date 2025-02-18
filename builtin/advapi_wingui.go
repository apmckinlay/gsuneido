// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
)

var advapi32 = MustLoadDLL("advapi32.dll")

// RegOpenKeyEx
var regOpenKeyEx = advapi32.MustFindProc("RegOpenKeyExA").Addr()
var _ = builtin(RegOpenKeyEx, "(hKey, lpSubKey, ulOptions, samDesired, phkResult)")

func RegOpenKeyEx(a, b, c, d, e Value) Value {
	var result uintptr
	rtn, _, _ := syscall.SyscallN(regOpenKeyEx,
		intArg(a),
		uintptr(zstrArg(b)),
		intArg(c),
		intArg(d),
		uintptr(unsafe.Pointer(&result)))
	e.Put(nil, SuStr("x"), IntVal(int(result)))
	return intRet(rtn)
}

// RegCloseKey
var regCloseKey = advapi32.MustFindProc("RegCloseKey").Addr()
var _ = builtin(RegCloseKey, "(hKey)")

func RegCloseKey(a Value) Value {
	rtn, _, _ := syscall.SyscallN(regCloseKey,
		intArg(a))
	return intRet(rtn)
}

// RegCreateKeyEx
var regCreateKeyEx = advapi32.MustFindProc("RegCreateKeyExA").Addr()
var _ = builtin(RegCreateKeyEx, "(hKey, lpSubKey, reserved/*unused*/, lpClass, "+
	"dwOptions, samDesired, lpSecurityAttributes, phkResult, lpdwDisposition)")

func RegCreateKeyEx(_ *Thread, a []Value) Value {
	var result uintptr
	rtn, _, _ := syscall.SyscallN(regCreateKeyEx,
		intArg(a[0]),
		uintptr(zstrArg(a[1])),
		0, // Reserved - must be 0
		uintptr(zstrArg(a[3])),
		intArg(a[4]),
		intArg(a[5]),
		0, // lpSecurityAttributes - always null
		uintptr(unsafe.Pointer(&result)),
		0) // lpdwDisposition - always null
	a[7].Put(nil, SuStr("x"), IntVal(int(result)))
	return intRet(rtn)
}

// RegQueryValueEx - hard coded for 4 byte data
var regQueryValueEx = advapi32.MustFindProc("RegQueryValueExA").Addr()
var _ = builtin(RegQueryValueEx, "(hKey, lpValueName, lpReserved/*unused*/, "+
	"lpType/*unused*/, lpData, lpcbData/*unused*/)")

func RegQueryValueEx(a, b, c, d, e, f Value) Value {
	var data int32
	cbData := int32(int32Size)
	rtn, _, _ := syscall.SyscallN(regQueryValueEx,
		intArg(a),
		uintptr(zstrArg(b)),
		0, // lpReserved - must be 0
		0, // lpType - NULL
		uintptr(unsafe.Pointer(&data)),
		uintptr(unsafe.Pointer(&cbData)))
	e.Put(nil, SuStr("x"), IntVal(int(data)))
	return intRet(rtn)
}

const REG_DWORD = 4

// RegSetValueEx - hard coded for 4 byte data
var regSetValueEx = advapi32.MustFindProc("RegSetValueExA").Addr()
var _ = builtin(RegSetValueEx, "(hKey, lpValueName, reserved/*unused*/, "+
	"dwType/*unused*/, lpData, cbData/*unused*/)")

func RegSetValueEx(a, b, c, d, e, f Value) Value {
	data := getInt32(e, "x")
	rtn, _, _ := syscall.SyscallN(regSetValueEx,
		intArg(a),
		uintptr(zstrArg(b)),
		0,         // reserved - must be 0
		REG_DWORD, // dwType
		uintptr(unsafe.Pointer(&data)),
		int32Size) // cbData
	return intRet(rtn)
}
