// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
)

var shlwapi = MustLoadDLL("shlwapi.dll")

// dll long SHCreateStreamOnFileA(string pszFile, int32 grfMode, POINTER* ppstm)
var shCreateStreamOnFile = shlwapi.MustFindProc("SHCreateStreamOnFileA").Addr()
var _ = builtin(SHCreateStreamOnFile, "(pszFile, grfMode, ppstm)")

func SHCreateStreamOnFile(a, b, c Value) Value {
	var p uintptr
	rtn, _, _ := syscall.SyscallN(shCreateStreamOnFile,
		uintptr(unsafe.Pointer(zstrArg(a))),
		intArg(b),
		uintptr(unsafe.Pointer(&p)))
	c.Put(nil, SuStr("x"), IntVal(int(p)))
	return intRet(rtn)
}
