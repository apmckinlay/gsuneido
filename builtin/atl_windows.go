package builtin

import (
	"syscall"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
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
		defer heap.FreeTo(heap.CurSize())
		iunk := heap.Alloc(int32Size)
		rtn, _, _ := syscall.Syscall(atlAxGetHost, 2,
			intArg(a),
			uintptr(iunk),
			0)
		b.Put(nil, SuStr("x"), IntVal(int(*(*int32)(iunk))))
		return intRet(rtn)
	})

// dll long Atl:AtlAxGetControl(pointer hwnd, LONG* iunk)
var atlAxGetControl = atl.MustFindProc("AtlAxGetControl").Addr()
var _ = builtin2("AtlAxGetControl(hwnd, iunk)",
	func(a Value, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		iunk := heap.Alloc(int32Size)
		rtn, _, _ := syscall.Syscall(atlAxGetControl, 2,
			intArg(a),
			uintptr(iunk),
			0)
		b.Put(nil, SuStr("x"), IntVal(int(*(*int32)(iunk))))
		return intRet(rtn)
	})

// dll long Atl:AtlAxAttachControl(LONG* iunk, pointer hwnd, LONG* unkContainer)
var atlAxAttachControl = atl.MustFindProc("AtlAxAttachControl").Addr()
var _ = builtin3("AtlAxAttachControl(iunk, hwnd, unkContainer)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		iunk := heap.Alloc(int32Size)
		*(*int32)(iunk) = getInt32(a, "x")
		unkContainer := heap.Alloc(int32Size)
		rtn, _, _ := syscall.Syscall(atlAxAttachControl, 3,
			uintptr(iunk),
			intArg(b),
			uintptr(unkContainer))
		b.Put(nil, SuStr("x"), IntVal(int(*(*int32)(unkContainer))))
		return intRet(rtn)
	})
