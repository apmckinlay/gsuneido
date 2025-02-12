// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"syscall"

	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
)

var dwmapi = MustLoadDLL("dwmapi.dll")

// dll pointer Dwmapi:DwmGetWindowAttribute(pointer hwnd, long dwAttribute,
// RECT* pvAttribute, long cbAttribute)
var dwmGetWindowAttribute = dwmapi.MustFindProc("DwmGetWindowAttribute").Addr()
var _ = builtin(DwmGetWindowAttributeRect,
	"(hwnd, dwAttribute, pvAttribute, cbAttribute)")

func DwmGetWindowAttributeRect(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn, _, _ := syscall.SyscallN(dwmGetWindowAttribute,
		intArg(a),
		intArg(b),
		uintptr(rectArg(c, r)),
		intArg(d))
	urectToOb(r, c)
	return intRet(rtn)
}
