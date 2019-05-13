// +build windows

package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var (
	user32      = windows.NewLazyDLL("user32.dll")
	getSysColor = user32.NewProc("GetSysColor")
	messageBox  = user32.NewProc("MessageBoxA")

	kernel32            = windows.NewLazyDLL("kernel32.dll")
	getCurrentDirectory = kernel32.NewProc("GetCurrentDirectoryA")
)

const maxpath = 1024

var _ = builtin0("GetCurrentDirectory()", func() Value {
	var buf [maxpath]byte
	n, _, _ := getCurrentDirectory.Call(maxpath, uintptr(unsafe.Pointer(&buf)))
	if n == 0 {
		return False
	}
	return SuStr(string(buf[:n]))
})

var _ = builtin1("GetSysColor(index)", func(arg Value) Value {
	index := ToInt(arg)
	n, _, _ := getSysColor.Call(uintptr(index))
	return IntVal(int(n))
})

var _ = builtin4("MessageBox(hwnd, text, caption, flags)",
	func(a, b, c, d Value) Value {
		hwnd := ToInt(a)
		s1, _ := windows.BytePtrFromString(IfStr(b))
		s2, _ := windows.BytePtrFromString(IfStr(c))
		flags := ToInt(d)
		n, _, _ := messageBox.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(s1)),
			uintptr(unsafe.Pointer(s2)),
			uintptr(flags))
		return IntVal(int(n))
	})
