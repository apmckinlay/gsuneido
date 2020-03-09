// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var atl = MustLoadDLL("atl.dll")

// dll bool atl:AtlAxWinInit()
var atlAxWinInit = atl.MustFindProc("AtlAxWinInit").Addr()
var _ = builtin0("AtlAxWinInit()",
	func() Value {
		rtn := goc.Syscall0(atlAxWinInit)
		return boolRet(rtn)
	})

// dll long Atl:AtlAxGetHost(pointer hwnd, LONG* iunk)
var atlAxGetHost = atl.MustFindProc("AtlAxGetHost").Addr()
var _ = builtin2("AtlAxGetHost(hwnd, iunk)",
	func(a Value, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		iunk := heap.Alloc(int32Size)
		rtn := goc.Syscall2(atlAxGetHost,
			intArg(a),
			uintptr(iunk))
		b.Put(nil, SuStr("x"), IntVal(int(*(*int32)(iunk))))
		return intRet(rtn)
	})

// dll long Atl:AtlAxGetControl(pointer hwnd, LONG* iunk)
var atlAxGetControl = atl.MustFindProc("AtlAxGetControl").Addr()
var _ = builtin2("AtlAxGetControl(hwnd, iunk)",
	func(a Value, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		iunk := heap.Alloc(int32Size)
		rtn := goc.Syscall2(atlAxGetControl,
			intArg(a),
			uintptr(iunk))
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
		rtn := goc.Syscall3(atlAxAttachControl,
			uintptr(iunk),
			intArg(b),
			uintptr(unkContainer))
		b.Put(nil, SuStr("x"), IntVal(int(*(*int32)(unkContainer))))
		return intRet(rtn)
	})
