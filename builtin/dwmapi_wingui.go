// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
)

var dwmapi = MustLoadDLL("dwmapi.dll")

// dll pointer Dwmapi:DwmGetWindowAttribute(pointer hwnd, long dwAttribute,
// RECT* pvAttribute, long cbAttribute)
var dwmGetWindowAttribute = dwmapi.MustFindProc("DwmGetWindowAttribute").Addr()
var _ = builtin(DwmGetWindowAttributeRect,
	"(hwnd, dwAttribute, pvAttribute, cbAttribute)")

func DwmGetWindowAttributeRect(a, b, c, d Value) Value {
	r := toRect(c)
	rtn, _, _ := syscall.SyscallN(dwmGetWindowAttribute,
		intArg(a),
		intArg(b),
		uintptr(unsafe.Pointer(r)),
		intArg(d))
	fromRect(r, c)
	return intRet(rtn)
}
