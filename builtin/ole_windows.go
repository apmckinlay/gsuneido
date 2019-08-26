package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var ole32 = windows.NewLazyDLL("ole32.dll")

// dll long Ole32:CreateStreamOnHGlobal(
// 	pointer hGlobal,
// 	bool fDeleteOnRelease,
// 	POINTER* ppstm)
var createStreamOnHGlobal = ole32.NewProc("CreateStreamOnHGlobal")
var _ = builtin3("CreateStreamOnHGlobal(hGlobal, fDeleteOnRelease, ppstm)",
	func(a, b, c Value) Value {
		var p uintptr
		rtn, _, _ := createStreamOnHGlobal.Call(
			intArg(a),
			boolArg(b),
			uintptr(unsafe.Pointer(&p)))
		c.Put(nil, SuStr("x"), IntVal(int(p)))
		return intRet(rtn)
	})
