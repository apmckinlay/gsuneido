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

func CreateLexer(name uintptr) uintptr {
	return uintptr(C.createLexer(C.uintptr(name)))
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

func CreateInstance(progid uintptr) uintptr {
	return uintptr(C.createInstance(C.uintptr(progid)))
}

func Invoke(idisp, name, flags, args, result uintptr) int {
	return int(C.invoke(C.uintptr(idisp), C.uintptr(name), C.uintptr(flags),
		C.uintptr(args), C.uintptr(result)))
}

func Release(iunk uintptr) int {
	return int(C.release(C.uintptr(iunk)))
}

// browser

//export SuneidoAPP
func SuneidoAPP(buf *C.buf_t) {
	url := C.GoString(buf.buf)
	result := SunAPP(url)
	buf.buf = C.CString(result) // WARNING: allocates
	buf.size = C.int(len(result))
}

func EmbedBrowserObject(hwnd, iunk, pptr uintptr) uintptr {
	return uintptr(C.EmbedBrowserObject(C.uintptr(hwnd), C.uintptr(iunk),
		C.uintptr(pptr)))
}

func UnEmbedBrowserObject(iunk, ptr uintptr) {
	C.UnEmbedBrowserObject(C.uintptr(iunk), C.uintptr(ptr))
}

func WebBrowser2(op, arg1, arg2, arg3, arg4, arg5 uintptr) uintptr {
	return uintptr(C.WebView2(C.uintptr(op), C.uintptr(arg1), C.uintptr(arg2),
		C.uintptr(arg3), C.uintptr(arg4), C.uintptr(arg5)))
}
