package builtin

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("MessageLoop(hwnd = 0)", func(t *Thread, args []Value) Value {
	MessageLoop(ToInt(args[0]))
	return nil
})

type MSG struct {
	hwnd    HANDLE
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      POINT
	_       [4]byte // padding
}

var uiThreadId uintptr

const WM_USER = 0x400

var retChan = make(chan uintptr, 1)

func Init() {
	runtime.LockOSThread()
	uiThreadId = GetCurrentThreadId()
}

const GA_ROOT = 2
const WM_QUIT = 0x12
const GMEM_FIXED = 0

func MessageLoop(hdlg int) {
	// see: https://github.com/lxn/walk/pull/493
	msg := (*MSG)(GlobalAlloc(unsafe.Sizeof(MSG{})))
	defer GlobalFree(unsafe.Pointer(msg))
	for {
		if -1 == GetMessage(msg, 0, 0, 0) {
			continue // error
		}
		// if msg.message == WM_QUIT {
		// 	if hdlg != 0 {
		// 		PostQuitMessage(msg.wParam)
		// 		return
		// 	}
		// 	return //TODO shutdown(msg.wParam)
		// }
		// if msg.message == WM_USER && msg.hwnd == 0 {
		// 	// from SetTimer in another thread
		// 	rtn, _, _ := syscall.Syscall(setTimer, ?,0, 0, msg.wParam, msg.lParam)
		// 	retChan <- rtn
		// 	continue
		// }
		// if window := GetAncestor(msg.hwnd, GA_ROOT); window != 0 {
		// 	// if (HACCEL haccel = (HACCEL) GetWindowLong(window, GWL_USERDATA))
		// 	// 	if (TranslateAccelerator(window, haccel, &msg))
		// 	// 		continue;
		// 	if IsDialogMessage(window, msg) {
		// 		continue
		// 	}
		// }
		TranslateMessage(msg)
		DispatchMessage(msg)
	}
}

func GetCurrentThreadId() uintptr {
	ret, _, _ := syscall.Syscall(getCurrentThreadId, 0, 0, 0, 0)
	return ret
}

func GlobalAlloc(dwBytes uintptr) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(globalAlloc, 2,
		GMEM_FIXED,
		dwBytes,
		0)
	return unsafe.Pointer(ret)
}
func GlobalFree(hglb unsafe.Pointer) uintptr {
	ret, _, _ := syscall.Syscall(globalFree, 2,
		uintptr(hglb),
		0, 0)
	return ret
}

var getMessage = user32.MustFindProc("GetMessageA").Addr()

func GetMessage(msg *MSG, hwnd HANDLE, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := syscall.Syscall6(getMessage, 4,
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
		0, 0)
	return int(ret)
}

var translateMessage = user32.MustFindProc("TranslateMessage").Addr()

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := syscall.Syscall(translateMessage, 1,
		uintptr(unsafe.Pointer(msg)),
		0, 0)
	return ret != 0
}

var dispatchMessage = user32.MustFindProc("DispatchMessageA").Addr()

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := syscall.Syscall(dispatchMessage, 1,
		uintptr(unsafe.Pointer(msg)),
		0, 0)
	return ret
}

func GetAncestor(hwnd HANDLE, gaFlags uint32) uintptr {
	ret, _, _ := syscall.Syscall(getAncestor, 2,
		uintptr(hwnd),
		uintptr(gaFlags),
		0)
	return ret
}

func IsDialogMessage(hwnd HANDLE, msg *MSG) bool {
	ret, _, _ := syscall.Syscall(isDialogMessage, 2,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(msg)),
		0)
	return ret != 0
}

func PostQuitMessage(exitCode uintptr) uintptr {
	ret, _, _ := syscall.Syscall(postQuitMessage, 2,
		uintptr(exitCode),
		0, 0)
	return ret
}

//-------------------------------------------------------------------

// since we don't have the functions it calls
var _ = builtin0("WorkSpaceStatus()", func() Value {
	return EmptyStr
})

//-------------------------------------------------------------------

// NOT thread safe

const align uintptr = 8
const heapsize = 64 * 1024
const magicBefore = 0xb5
const magicAfter = 0xc3

var myheap = [heapsize]byte{248, 249, 250, 251, 252, 253, 254, 255}
var heapnext = align

func alloc(n uintptr) unsafe.Pointer {
	if n%align != 0 {
		panic("unaligned alloc")
	}
	heapcheck("alloc")
	p := &myheap[heapnext]
	heapnext += n + align
	for i := align; i > 0; i-- {
		myheap[heapnext-i] = byte(256 - i)
	}
	return unsafe.Pointer(p)
}

func free(p unsafe.Pointer, n uintptr) {
	heapcheck("free1")
	heapnext -= n + align
	if (*byte)(p) != &myheap[heapnext] {
		panic("non-stack free")
	}
	heapcheck("free2")
}

func heapcheck(s string) {
	for i := align; i > 0; i-- {
		if myheap[heapnext-i] != byte(256-i) {
			fmt.Println("heapnext", heapnext, "i", i, ",",
				myheap[heapnext-i], "should be", byte(256-i))
			panic("heap corrupt " + s)
		}
	}
}
