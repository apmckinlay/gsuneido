package builtin

import (
	"runtime"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("MessageLoop()", func(t *Thread, _ []Value) Value {
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

var uiThreadId uintptr

const WM_USER = 0x400

var retChan = make(chan uintptr, 1)

func Init() {
	runtime.LockOSThread()
	uiThreadId, _, _ = getCurrentThreadId.Call()
}

func MessageLoop() {
	var msg MSG
	for 0 != GetMessage(&msg, 0, 0, 0) {
		if msg.hwnd == 0 && msg.message == WM_USER {
			rtn, _, _ := setTimer.Call(0, 0, msg.wParam, msg.lParam)
			retChan <- rtn
			continue
		}
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
