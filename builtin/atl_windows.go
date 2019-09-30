package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var atl = windows.MustLoadDLL("atl.dll")

// dll bool atl:AtlAxWinInit()
var atlAxWinInit = atl.MustFindProc("AtlAxWinInit").Addr()
var _ = builtin0("AtlAxWinInit()",
	func() Value {
		rtn, _, _ := syscall.Syscall(atlAxWinInit, 0, 0, 0, 0)
		return boolRet(rtn)
	})

// dll long Atl:AtlAxGetHost(pointer hwnd, LONG* iunk)
var atlAxGetHost = atl.MustFindProc("AtlAxGetHost").Addr()
var _ = builtin2("AtlAxGetHost(hwnd, iunk)",
	func(a Value, b Value) Value {
		var iunk int32
		rtn, _, _ := syscall.Syscall(atlAxGetHost, 2,
			intArg(a),
			uintptr(unsafe.Pointer(&iunk)),
			0)
		b.Put(nil, SuStr("x"), IntVal(int(iunk)))
		return intRet(rtn)
	})

// dll long Atl:AtlAxGetControl(pointer hwnd, LONG* iunk)
var atlAxGetControl = atl.MustFindProc("AtlAxGetControl").Addr()
var _ = builtin2("AtlAxGetControl(hwnd, iunk)",
	func(a Value, b Value) Value {
		var iunk int32
		rtn, _, _ := syscall.Syscall(atlAxGetControl, 2,
			intArg(a),
			uintptr(unsafe.Pointer(&iunk)),
			0)
		b.Put(nil, SuStr("x"), IntVal(int(iunk)))
		return intRet(rtn)
	})

// dll long Atl:AtlAxAttachControl(LONG* iunk, pointer hwnd, LONG* unkContainer)
var atlAxAttachControl = atl.MustFindProc("AtlAxAttachControl").Addr()
var _ = builtin3("AtlAxAttachControl(iunk, hwnd, unkContainer)",
	func(a, b, c Value) Value {
		iunk := getInt32(a, "x")
		var unkContainer int32
		rtn, _, _ := syscall.Syscall(atlAxAttachControl, 3,
			uintptr(unsafe.Pointer(&iunk)),
			intArg(b),
			uintptr(unsafe.Pointer(&unkContainer)))
		b.Put(nil, SuStr("x"), IntVal(int(unkContainer)))
		return intRet(rtn)
	})
