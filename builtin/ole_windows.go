// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var ole32 = MustLoadDLL("ole32.dll")

// dll long Ole32:CreateStreamOnHGlobal(
// 	pointer hGlobal,
// 	bool fDeleteOnRelease,
// 	POINTER* ppstm)
var createStreamOnHGlobal = ole32.MustFindProc("CreateStreamOnHGlobal").Addr()
var _ = builtin3("CreateStreamOnHGlobal(hGlobal, fDeleteOnRelease, ppstm)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(uintptrSize)
		rtn := goc.Syscall3(createStreamOnHGlobal,
			intArg(a),
			boolArg(b),
			uintptr(p))
		c.Put(nil, SuStr("x"), IntVal(int(*(*uintptr)(p))))
		return intRet(rtn)
	})

var oleaut32 = MustLoadDLL("oleaut32.dll")

// dll long OleAut32:OleLoadPicture(
//		pointer lpstream,
//		long lSize,
//		bool fRunmode,
//		GUID* riid,
//		POINTER* lplpvObj)
var oleLoadPicture = oleaut32.MustFindProc("OleLoadPicture").Addr()
var _ = builtin5("OleLoadPicture(lpstream, lSize, fRunmode, riid, lplpvobj)",
	func(a, b, c, d, e Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(uintptrSize)
		g := heap.Alloc(nGUID)
		guid := (*GUID)(g)
		*guid = GUID{
			Data1: getInt32(d, "Data1"),
			Data2: int16(getInt(d, "Data2")),
			Data3: int16(getInt(d, "Data3")),
		}
		data4 := d.Get(nil, SuStr("Data4"))
		for i := 0; i < 8; i++ {
			guid.Data4[i] = byte(ToInt(data4.Get(nil, SuInt(i))))
		}
		rtn := goc.Syscall5(oleLoadPicture,
			intArg(a),
			intArg(b),
			boolArg(c),
			uintptr(g),
			uintptr(p))
		e.Put(nil, SuStr("x"), IntVal(int(*(*uintptr)(p))))
		return intRet(rtn)
	})

type GUID struct {
	Data1 int32
	Data2 int16
	Data3 int16
	Data4 [8]byte
}

const nGUID = unsafe.Sizeof(GUID{})
