package builtin

import (
	"syscall"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var dwmapi = windows.MustLoadDLL("dwmapi.dll")

// dll pointer Dwmapi:DwmGetWindowAttribute(pointer hwnd, long dwAttribute,
//		RECT* pvAttribute, long cbAttribute)
var dwmGetWindowAttribute = dwmapi.MustFindProc("DwmGetWindowAttribute").Addr()
var _ = builtin4("DwmGetWindowAttributeRect(hwnd, dwAttribute, pvAttribute,"+
	" cbAttribute)",
	func(a, b, c, d Value) Value {
		var r RECT
		rtn, _, _ := syscall.Syscall6(dwmGetWindowAttribute, 4,
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, &r)),
			intArg(d),
			0, 0)
		rectToOb(&r, c)
		return intRet(rtn)
	})
