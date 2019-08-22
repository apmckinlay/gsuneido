package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var UIThread *Thread

var _ = builtin("MessageLoop()", func(t *Thread, _ []Value) Value {
	UIThread = t
	MessageLoop()
	return nil
})

type MSG struct {
	hwnd    HANDLE
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      POINT
}

func MessageLoop() {
	var msg MSG
	for 0 != GetMessage(&msg, 0, 0, 0) {
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
}

var getMessage = user32.NewProc("GetMessageA")

func GetMessage(msg *MSG, hwnd HANDLE, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := getMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

var translateMessage = user32.NewProc("TranslateMessage")

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := translateMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

var dispatchMessage = user32.NewProc("DispatchMessageA")

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := dispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret
}
