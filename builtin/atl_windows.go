package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var atl = windows.NewLazyDLL("atl.dll")

// dll bool atl:AtlAxWinInit()
var atlAxWinInit = atl.NewProc("AtlAxWinInit")
var _ = builtin0("AtlAxWinInit()",
	func() Value {
		rtn, _, _ := atlAxWinInit.Call()
		return boolRet(rtn)
	})

// dll long Atl:AtlAxGetHost(pointer hwnd, LONG* iunk)
var atlAxGetHost = atl.NewProc("AtlAxGetHost")
var _ = builtin2("AtlAxGetHost(hwnd, iunk)",
	func(a Value, b Value) Value {
		var iunk int32
		rtn, _, _ := atlAxGetHost.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&iunk)))
		b.Put(nil, SuStr("x"), IntVal(int(iunk)))
		return intRet(rtn)
	})

// dll long Atl:AtlAxGetControl(pointer hwnd, LONG* iunk)
var atlAxGetControl = atl.NewProc("AtlAxGetControl")
var _ = builtin2("AtlAxGetControl(hwnd, iunk)",
	func(a Value, b Value) Value {
		var iunk int32
		rtn, _, _ := atlAxGetControl.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&iunk)))
		b.Put(nil, SuStr("x"), IntVal(int(iunk)))
		return intRet(rtn)
	})

// dll long Atl:AtlAxAttachControl(LONG* iunk, pointer hwnd, LONG* unkContainer)
var atlAxAttachControl = atl.NewProc("AtlAxAttachControl")
var _ = builtin3("AtlAxAttachControl(iunk, hwnd, unkContainer)",
	func(a, b, c Value) Value {
		iunk := getInt32(a, "x")
		var unkContainer int32
		rtn, _, _ := atlAxAttachControl.Call(
			uintptr(unsafe.Pointer(&iunk)),
			intArg(b),
			uintptr(unsafe.Pointer(&unkContainer)))
		b.Put(nil, SuStr("x"), IntVal(int(unkContainer)))
		return intRet(rtn)
	})
