package builtin

import (
	"syscall"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
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
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn, _, _ := syscall.Syscall6(dwmGetWindowAttribute, 4,
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, r)),
			intArg(d),
			0, 0)
		urectToOb(r, c)
		return intRet(rtn)
	})
