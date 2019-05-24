// +build windows

package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazyDLL("user32.dll")
	kernel32 = windows.NewLazyDLL("kernel32.dll")
)

type RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

// const maxpath = 1024

// var getCurrentDirectory = kernel32.NewProc("GetCurrentDirectoryA")

// var _ = builtin0("GetCurrentDirectory()", func() Value {
// 	var buf [maxpath]byte
// 	n, _, _ := getCurrentDirectory.Call(maxpath, uintptr(unsafe.Pointer(&buf)))
// 	if n == 0 {
// 		return False
// 	}
// 	return SuStr(string(buf[:n]))
// })

var getDesktopWindow = user32.NewProc("GetDesktopWindow")

var _ = builtin0("GetDesktopWindow()", func() Value {
	n, _, _ := getDesktopWindow.Call()
	return IntVal(int(n))
})

var getSysColor = user32.NewProc("GetSysColor")

var _ = builtin1("GetSysColor(index)", func(arg Value) Value {
	index := ToInt(arg)
	n, _, _ := getSysColor.Call(uintptr(index))
	return IntVal(int(n))
})

var getWindowRect = user32.NewProc("GetWindowRect")

var _ = builtin1("GetWindowRect(hwnd)", func(arg Value) Value {
	hwnd := ToInt(arg)
	var r RECT
	b, _, _ := getWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	if b == 0 {
		panic("GetWindowRect failed")
	}
	ob := &SuObject{}
	ob.Put(nil, SuStr("left"), IntVal(int(r.left)))
	ob.Put(nil, SuStr("top"), IntVal(int(r.top)))
	ob.Put(nil, SuStr("right"), IntVal(int(r.right)))
	ob.Put(nil, SuStr("bottom"), IntVal(int(r.bottom)))
	return ob
})

var messageBox = user32.NewProc("MessageBoxA")

var _ = builtin4("MessageBox(hwnd, text, caption, flags)",
	func(a, b, c, d Value) Value {
		hwnd := ToInt(a)
		s1, _ := windows.BytePtrFromString(ToStr(b))
		s2, _ := windows.BytePtrFromString(ToStr(c))
		flags := ToInt(d)
		n, _, _ := messageBox.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(s1)),
			uintptr(unsafe.Pointer(s2)),
			uintptr(flags))
		return IntVal(int(n))
	})
