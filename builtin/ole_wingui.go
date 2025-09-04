// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
)

var ole32 = MustLoadDLL("ole32.dll")

// dll long Ole32:CreateStreamOnHGlobal(
// pointer hGlobal,
// bool fDeleteOnRelease,
// POINTER* ppstm)
var createStreamOnHGlobal = ole32.MustFindProc("CreateStreamOnHGlobal").Addr()
var _ = builtin(CreateStreamOnHGlobal, "(hGlobal, fDeleteOnRelease, ppstm)")

func CreateStreamOnHGlobal(a, b, c Value) Value {
	var x uintptr
	rtn, _, _ := syscall.SyscallN(createStreamOnHGlobal,
		intArg(a),
		boolArg(b),
		uintptr(unsafe.Pointer(&x)))
	c.Put(nil, SuStr("x"), IntVal(int(x)))
	return intRet(rtn)
}

var oleaut32 = MustLoadDLL("oleaut32.dll")

// dll long OleAut32:OleLoadPicture(
// pointer lpstream,
// long lSize,
// bool fRunmode,
// GUID* riid,
// POINTER* lplpvObj)
var oleLoadPicture = oleaut32.MustFindProc("OleLoadPicture").Addr()
var _ = builtin(OleLoadPicture, "(lpstream, lSize, fRunmode, riid, lplpvobj)")

func OleLoadPicture(a, b, c, d, e Value) Value {
	var p uintptr
	guid := &GUID{
		Data1: getInt32(d, "Data1"),
		Data2: int16(getInt(d, "Data2")),
		Data3: int16(getInt(d, "Data3")),
	}
	data4 := d.Get(nil, SuStr("Data4"))
	for i := range 8 {
		guid.Data4[i] = byte(ToInt(data4.Get(nil, SuInt(i))))
	}
	rtn, _, _ := syscall.SyscallN(oleLoadPicture,
		intArg(a),
		intArg(b),
		boolArg(c),
		uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(&p)))
	e.Put(nil, SuStr("x"), IntVal(int(p)))
	return intRet(rtn)
}

type GUID struct {
	Data1 int32
	Data2 int16
	Data3 int16
	Data4 [8]byte
}
