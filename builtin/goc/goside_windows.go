// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package goc

// #cgo CFLAGS: -DWINVER=0x601 -D_WIN32_WINNT=0x0601
// #cgo LDFLAGS: -L . -lscintilla -lgdi32 -lcomdlg32 -lcomctl32 -limm32 -lmsimg32
// #cgo LDFLAGS: -lurlmon -lole32 -loleaut32 -luuid -lwininet -static
// #include "cside.h"
import "C"

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/apmckinlay/gsuneido/options"
	"golang.org/x/sys/windows"
)

const Ncb2s = C.ncb2s
const Ncb3s = C.ncb3s
const Ncb4s = C.ncb4s

var uiThreadId uint32

func Init() {
	runtime.LockOSThread()
	uiThreadId = windows.GetCurrentThreadId()
	C.start()
}

func Run() {
	C.args[0] = C.msg_result
	C.signalAndWait()
	interact()
	log.Fatalln("!!! should not reach here !!!")
}

func CThreadId() uintptr {
	return uintptr(C.threadid)
}

func CNotifyHwnd() uintptr {
	return uintptr(C.notifyHwnd)
}

func QueryIDispatch(iunk uintptr) uintptr {
	return interact(C.msg_queryidispatch, iunk)
}

func CreateInstance(progid uintptr) uintptr {
	return interact(C.msg_createinstance, progid)
}

func Invoke(idisp, name, flags, args, result uintptr) int {
	return int(interact(C.msg_invoke, idisp, name, flags, args, result))
}

func Release(iunk uintptr) {
	interact(C.msg_release, iunk)
}

func Traccel(ob int, msg unsafe.Pointer) int {
	return int(interact(C.msg_traccel, uintptr(ob), uintptr(msg)))
}

// Interrupt checks if control+break has been pressed.
// It is called regularly by Interp.
func Interrupt() bool {
	return interact(C.msg_interrupt) == 1
}

func GetCallback(nargs, i int) uintptr {
	switch nargs {
	case 2:
		return uintptr(C.cb2s[i])
	case 3:
		return uintptr(C.cb3s[i])
	case 4:
		return uintptr(C.cb4s[i])
	}
	log.Panicln("GetCallback unsupported nargs", nargs)
	return 0 // unreachable
}

func Alert(args ...interface{}) {
	s := fmt.Sprintln(args...)
	log.Print("Alert: ", s)
	if !options.Unattended {
		C.alert(C.CString(s))
	}
}

var fatalOnce sync.Once

func Fatal(args ...interface{}) {
	fatalOnce.Do(func() {
		s := fmt.Sprintln(args...)
		log.Print("FATAL: ", s)
		go func() {
			time.Sleep(10 * time.Second)
			os.Exit(1)
		}()
		if !options.Unattended {
			C.fatal(C.CString(s[:len(s)-1]))
		}
		os.Exit(1)
	})
}

// must be injected
var TimerId func(a uintptr)
var Callback2 func(i, a, b uintptr) uintptr
var Callback3 func(i, a, b, c uintptr) uintptr
var Callback4 func(i, a, b, c, d uintptr) uintptr
var RunOnGoSide func()
var SunAPP func(string) string

func interact(args ...uintptr) uintptr {
	if windows.GetCurrentThreadId() != uiThreadId {
		log.Println("illegal UI call from background thread")
		runtime.Goexit()
	}
	for i,a := range args {
		C.args[i] = C.uintptr(a)
	}
	for {
		// these are the messages sent from c-side to go-side
		switch C.args[0] {
		case C.msg_callback2:
			C.args[0] = C.msg_result
			C.args[1] = C.uintptr(Callback2(uintptr(C.args[1]),
				uintptr(C.args[2]), uintptr(C.args[3])))
		case C.msg_callback3:
			C.args[0] = C.msg_result
			C.args[1] = C.uintptr(Callback3(uintptr(C.args[1]),
				uintptr(C.args[2]), uintptr(C.args[3]), uintptr(C.args[4])))
		case C.msg_callback4:
			C.args[0] = C.msg_result
			C.args[1] = C.uintptr(Callback4(uintptr(C.args[1]),
				uintptr(C.args[2]), uintptr(C.args[3]), uintptr(C.args[4]),
				uintptr(C.args[5])))
		case C.msg_timerid:
			TimerId(uintptr(C.args[1]))
			C.args[0] = C.msg_result
		case C.msg_runongoside:
			RunOnGoSide()
			C.args[0] = C.msg_result
		case C.msg_sunapp:
			s := SunAPP(C.GoString((*C.char)(unsafe.Pointer(uintptr(C.args[1])))))
			C.args[0] = C.msg_result
			C.args[1] = (C.uintptr)(uintptr(unsafe.Pointer(C.CString(s))))
			// C.CString will malloc, sunapp will free
			C.args[2] = (C.uintptr)(len(s))
		case C.msg_result:
			return uintptr(C.args[1])
		}
		C.signalAndWait()
	}
}

func MessageLoop(hdlg uintptr) {
	interact(C.msg_msgloop, hdlg)
}

func Syscall0(adr uintptr) uintptr {
	return interact(C.msg_syscall, adr, 0)
}
func Syscall1(adr, a uintptr) uintptr {
	return interact(C.msg_syscall, adr, 1, a)
}
func Syscall2(adr, a, b uintptr) uintptr {
	return interact(C.msg_syscall, adr, 2, a, b)
}
func Syscall3(adr, a, b, c uintptr) uintptr {
	return interact(C.msg_syscall, adr, 3, a, b, c)
}
func Syscall4(adr, a, b, c, d uintptr) uintptr {
	return interact(C.msg_syscall, adr, 4, a, b, c, d)
}
func Syscall5(adr, a, b, c, d, e uintptr) uintptr {
	return interact(C.msg_syscall, adr, 5, a, b, c, d, e)
}
func Syscall6(adr, a, b, c, d, e, f uintptr) uintptr {
	return interact(C.msg_syscall, adr, 6, a, b, c, d, e, f)
}
func Syscall7(adr, a, b, c, d, e, f, g uintptr) uintptr {
	return interact(C.msg_syscall, adr, 7, a, b, c, d, e, f, g)
}
func Syscall8(adr, a, b, c, d, e, f, g, h uintptr) uintptr {
	return interact(C.msg_syscall, adr, 8, a, b, c, d, e, f, g, h)
}
func Syscall9(adr, a, b, c, d, e, f, g, h, i uintptr) uintptr {
	return interact(C.msg_syscall, adr, 9, a, b, c, d, e, f, g, h, i)
}
func Syscall10(adr, a, b, c, d, e, f, g, h, i, j uintptr) uintptr {
	return interact(C.msg_syscall, adr, 10, a, b, c, d, e, f, g, h, i, j)
}
func Syscall11(adr, a, b, c, d, e, f, g, h, i, j, k uintptr) uintptr {
	return interact(C.msg_syscall, adr, 11, a, b, c, d, e, f, g, h, i, j, k)
}
func Syscall12(adr, a, b, c, d, e, f, g, h, i, j, k, l uintptr) uintptr {
	return interact(C.msg_syscall, adr, 12, a, b, c, d, e, f, g, h, i, j, k, l)
}
func Syscall14(adr, a, b, c, d, e, f, g, h, i, j, k, l, m, n uintptr) uintptr {
	return interact(C.msg_syscall, adr, 14, a, b, c, d, e, f, g, h, i, j, k, l, m, n)
}
