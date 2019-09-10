package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var dwmapi = windows.NewLazyDLL("dwmapi.dll")

// dll pointer Dwmapi:DwmGetWindowAttribute(pointer hwnd, long dwAttribute,
//		RECT* pvAttribute, long cbAttribute)
var dwmGetWindowAttribute = dwmapi.NewProc("DwmGetWindowAttribute")
var _ = builtin4("DwmGetWindowAttributeRect(hwnd, dwAttribute, pvAttribute,"+
	" cbAttribute)",
	func(a, b, c, d Value) Value {
		var r RECT
		rtn, _, _ := dwmGetWindowAttribute.Call(
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, &r)),
			intArg(d))
		rectToOb(&r, c)
		return intRet(rtn)
	})
