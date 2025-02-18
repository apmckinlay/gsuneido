// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package goc

// #cgo LDFLAGS: -L . -lscintilla -llexilla -lgdi32 -lcomdlg32 -lcomctl32 -limm32 -lmsimg32
// #cgo LDFLAGS: -lurlmon -lole32 -loleaut32 -luuid -lwininet -lshlwapi -static
// #include "cside.h"
import "C"

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/exit"
	"golang.org/x/sys/windows"
)

var SunAPP func(string) string

var uiThreadId uint32

func init() {
	runtime.LockOSThread()
	uiThreadId = windows.GetCurrentThreadId()
	C.setup()
}

func Run() int {
	return int(C.run())
}

func Interrupt() bool {
	return C.interrupt() == 1
}

func CreateLexer(name unsafe.Pointer) uintptr {
	return uintptr(C.createLexer((*C.char)(name)))
}

func MessageLoop(hdlg uintptr) {
	C.message_loop(C.uintptr(hdlg))
}

func Alert(args ...any) {
	s := fmt.Sprintln(args...)
	log.Print("Alert: ", s)
	if options.Mode == "gui" {
		C.alert(C.CString(s), 0)
	}
}

func AlertCancel(args ...any) bool {
	s := fmt.Sprintln(args...)
	log.Print("Alert: ", s)
	if options.Mode == "gui" {
		if 2 == C.alert(C.CString(s), 1) {
			return false // cancel
		}
	}
	return true // ok
}

var fatalOnce sync.Once

func Fatal(s string) {
	go func() {
		time.Sleep(10 * time.Second)
		exit.Exit(1) // failsafe
	}()
	if options.Mode == "gui" {
		fatalOnce.Do(func() {
			C.fatal(C.CString(s[:len(s)-1]))
		})
	}
}

// COM

func QueryIDispatch(iunk uintptr) uintptr {
	return uintptr(C.queryIDispatch(C.uintptr(iunk)))
}

func CreateInstance(progid string) uintptr {
	buf := make([]byte, len(progid)+2) // +2 to ensure double nul termination
	copy(buf, progid)
	return uintptr(C.createInstance((*C.char)(unsafe.Pointer(&buf[0]))))
}

func Invoke(idisp uintptr, name string, flags uintptr,
	args unsafe.Pointer, result unsafe.Pointer) int {
	return int(C.invoke(C.uintptr(idisp),
		(*C.char)(zstr(name)),
		C.uintptr(flags), args, result))
}

func Release(iunk uintptr) int {
	return int(C.release(C.uintptr(iunk)))
}

// browser

//export SuneidoAPP
func SuneidoAPP(buf *C.buf_t) {
	url := C.GoString(buf.data)
	result := SunAPP(url)
	buf.data = C.CString(result) // WARNING: allocates
	buf.size = C.int(len(result))
}

func EmbedBrowserObject(hwnd uintptr) (ret, iunk, ptr uintptr) {
	ret = uintptr(C.EmbedBrowserObject(C.uintptr(hwnd),
		unsafe.Pointer(&iunk),
		unsafe.Pointer(&ptr)))
	return
}

func UnEmbedBrowserObject(iunk, ptr uintptr) {
	C.UnEmbedBrowserObject(C.uintptr(iunk), C.uintptr(ptr))
}

func WebView2_Create(hwnd uintptr, pBrowserObject unsafe.Pointer, dllPath string,
	userDataFolder string, cb uintptr) uintptr {
	return uintptr(C.WebView2_Create(C.uintptr(hwnd), pBrowserObject,
		(*C.char)(zstr(dllPath)),
		(*C.char)(zstr(userDataFolder)),
		C.uintptr(cb)))
}

func zstr(s string) unsafe.Pointer {
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	return unsafe.Pointer(&buf[0])
}

func WebView2_Resize(pBrowserObject uintptr, w uintptr, h uintptr) uintptr {
	return uintptr(C.WebView2_Resize(C.uintptr(pBrowserObject), C.long(w), C.long(h)))
}

func WebView2_Navigate(pBrowserObject uintptr, s string) uintptr {
	return uintptr(C.WebView2_Navigate(C.uintptr(pBrowserObject),
		(*C.char)(zstr(s))))
}

func WebView2_NavigateToString(pBrowserObject uintptr, s string) uintptr {
	return uintptr(C.WebView2_NavigateToString(C.uintptr(pBrowserObject),
		(*C.char)(zstr(s))))
}

func WebView2_ExecuteScript(pBrowserObject uintptr, script string) uintptr {
	return uintptr(C.WebView2_ExecuteScript(C.uintptr(pBrowserObject),
		(*C.char)(zstr(script))))
}

func WebView2_GetSource(pBrowserObject uintptr, dst *byte) uintptr {
	return uintptr(C.WebView2_GetSource(C.uintptr(pBrowserObject),
		(*C.char)(unsafe.Pointer(dst))))
}

func WebView2_Print(pBrowserObject uintptr) uintptr {
	return uintptr(C.WebView2_Print(C.uintptr(pBrowserObject)))
}

func WebView2_SetFocus(pBrowserObject uintptr) uintptr {
	return uintptr(C.WebView2_SetFocus(C.uintptr(pBrowserObject)))
}

func WebView2_Close(pBrowserObject uintptr) uintptr {
	return uintptr(C.WebView2_Close(C.uintptr(pBrowserObject)))
}
