// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"fmt"
	"math"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
)

type suCOMObject struct {
	ValueBase[*suCOMObject]
	ptr   uintptr
	idisp bool // true if IDispatch
}

var _ Value = (*suCOMObject)(nil)

var _ = builtin(COMobject, "(progid)")

func COMobject(arg Value) Value {
	if n, ok := arg.ToInt(); ok {
		ptr := uintptr(n)
		if idisp := goc.QueryIDispatch(ptr); idisp != 0 {
			goc.Release(ptr)
			return &suCOMObject{ptr: idisp, idisp: true}
		}
		return &suCOMObject{ptr: ptr}
	}
	if s, ok := arg.ToStr(); ok {
		defer heap.FreeTo(heap.CurSize())
		idisp := goc.CreateInstance(uintptr(heap.CopyStr(s)))
		if idisp == 0 {
			return False
		}
		return &suCOMObject{ptr: idisp, idisp: true}
	}
	panic("COMobject requires integer or string")
}

var suComObjectMethods = methods()

var _ = method(com_DispatchQ, "()")

func com_DispatchQ(this Value) Value {
	return SuBool(this.(*suCOMObject).idisp)
}

var _ = method(com_Release, "()")

func com_Release(this Value) Value {
	return IntVal(goc.Release(this.(*suCOMObject).ptr))
}

func (sco *suCOMObject) Get(_ *Thread, mem Value) Value {
	if !sco.idisp {
		panic("COMobject can't get property of IUnknown")
	}
	return GetProperty(sco.ptr, ToStr(mem))
}

func (sco *suCOMObject) Put(_ *Thread, mem Value, val Value) {
	if !sco.idisp {
		panic("COMobject can't put property of IUnknown")
	}
	PutProperty(sco.ptr, ToStr(mem), val)

}

func (sco *suCOMObject) GetPut(th *Thread, m Value, v Value,
	op func(x, y Value) Value, retOrig bool) Value {
	orig := sco.Get(th, m)
	v = op(orig, v)
	sco.Put(th, m, v)
	if retOrig {
		return orig
	}
	return v
}

func (sco *suCOMObject) Equal(other any) bool {
	return sco == other
}

func (*suCOMObject) SetConcurrent() {
	// ok since immutable (assuming the COM object is thread safe)
}

func (sco *suCOMObject) Lookup(_ *Thread, method string) Value {
	if f, ok := suComObjectMethods[method]; ok {
		return f
	}
	return &SuBuiltinMethodRaw{
		Fn: func(th *Thread, as *ArgSpec, this Value, args []Value) Value {
			return sco.call(method, as, args)
		}}
}

func (sco *suCOMObject) call(method string, as *ArgSpec, args []Value) Value {
	if !sco.idisp {
		panic("COMobject can't call method of IUnknown")
	}
	if as.Spec != nil || as.Each != 0 {
		panic("COMobject invalid call arguments")
	}
	return CallMethod(sco.ptr, method, args[:as.Nargs])
}

//-------------------------------------------------------------------

const (
	DISPATCH_METHOD      = 1
	DISPATCH_PROPERTYGET = 2
	DISPATCH_PROPERTYPUT = 4
)

var flagnames = []string{1: "call", 2: "get", 4: "put"}

type stVariant struct {
	vt  uint16
	_   uint16
	_   uint16
	_   uint16
	val int64
	_   [8]byte
}

const nVariant = unsafe.Sizeof(stVariant{})

func GetProperty(idisp uintptr, name string) Value {
	return invoke(idisp, name, DISPATCH_PROPERTYGET)
}

func PutProperty(idisp uintptr, name string, value Value) {
	invoke(idisp, name, DISPATCH_PROPERTYPUT, value)
}

func CallMethod(idisp uintptr, name string, args []Value) Value {
	return invoke(idisp, name, DISPATCH_METHOD, args...)
}

func invoke(idisp uintptr, name string, flags uintptr, args ...Value) Value {
	defer heap.FreeTo(heap.CurSize())
	pargs := convertArgs(args)
	params := heap.Alloc(nDispParams)
	dp := (*stDispParams)(params)
	dp.cArgs = uint32(len(args))
	dp.rgvarg = uintptr(pargs)
	var result unsafe.Pointer
	if flags == DISPATCH_PROPERTYPUT {
		dp.cNamedArgs = 1
		dp.rgdispidNamedArgs = uintptr(unsafe.Pointer(&DISPID_PROPERTYPUT))
	} else {
		result = heap.Alloc(nVariant)
	}
	hr := goc.Invoke(idisp, uintptr(heap.CopyStr(name)), flags, uintptr(params),
		uintptr(result))
	if hr < 0 {
		panic(fmt.Sprintf("COMobject %s failed %s %x",
			flagnames[flags], name, uint32(hr)))
	}
	if flags == DISPATCH_PROPERTYPUT {
		return nil
	}
	return variantToSu((*stVariant)(result))
}

func convertArgs(args []Value) unsafe.Pointer {
	pargs := heap.Alloc(nVariant * uintptr(len(args)))
	p := pargs
	for i := len(args) - 1; i >= 0; i-- {
		suToVariant(args[i], (*stVariant)(p))
		p = unsafe.Pointer(uintptr(p) + nVariant)
	}
	return pargs
}

var DISPID_PROPERTYPUT int32 = -3

func suToVariant(x Value, v *stVariant) {
	if x == True {
		v.vt = VT_BOOL
		v.val = -1
	} else if x == False {
		v.vt = VT_BOOL
		v.val = 0
	} else if n, ok := x.ToInt(); ok {
		if math.MinInt32 <= n && n <= math.MaxInt32 {
			v.vt = VT_I4
		} else {
			v.vt = VT_I8
		}
		v.val = int64(n)
	} else if _, ok := x.ToStr(); ok {
		v.vt = VT_BSTR
		v.val = int64(uintptr(stringArg(x))) // C side converts
	} else if sco, ok := x.(*suCOMObject); ok {
		if sco.idisp {
			v.vt = VT_DISPATCH
		} else {
			v.vt = VT_UNKNOWN
		}
		v.val = int64(sco.ptr)
	} else {
		panic("COMobject can't convert " + ErrType(x))
	}
}

type stDispParams struct {
	rgvarg            uintptr
	rgdispidNamedArgs uintptr
	cArgs             uint32
	cNamedArgs        uint32
}

const nDispParams = unsafe.Sizeof(stDispParams{})

const (
	VT_EMPTY    = 0
	VT_NULL     = 1
	VT_I2       = 2
	VT_I4       = 3
	VT_BSTR     = 8
	VT_DISPATCH = 9
	VT_BOOL     = 11
	VT_UNKNOWN  = 13
	VT_UI2      = 18
	VT_UI4      = 19
	VT_I8       = 20
	VT_UI8      = 21
)

func variantToSu(v *stVariant) Value {
	var result Value
	switch v.vt {
	case VT_NULL, VT_EMPTY:
		result = Zero
	case VT_BOOL:
		result = SuBool(v.val != 0)
	case VT_I2:
		result = IntVal(int(int16(v.val)))
	case VT_I4:
		result = IntVal(int(int32(v.val)))
	case VT_I8:
		result = IntVal(int(v.val))
	case VT_UI2:
		result = IntVal(int(uint16(v.val)))
	case VT_UI4:
		result = IntVal(int(uint32(v.val)))
	case VT_UI8:
		result = IntVal(int(v.val))
	case VT_BSTR:
		result = SuStr(bstrToString(v))
		VariantClear(v)
	case VT_DISPATCH:
		result = &suCOMObject{ptr: uintptr(v.val), idisp: true}
	case VT_UNKNOWN:
		iunk := uintptr(v.val)
		if idisp := goc.QueryIDispatch(iunk); idisp != 0 {
			goc.Release(iunk)
			result = &suCOMObject{ptr: idisp, idisp: true}
		} else {
			result = &suCOMObject{ptr: iunk}
		}
	default:
		panic("COMobject: can't convert to Suneido value")
	}
	return result
}

func bstrToString(v *stVariant) string {
	if v.val == 0 {
		return ""
	}
	p := uintptr(v.val)
	length := SysStringLen(p)
	a := make([]uint16, length)
	for i := range int(length) {
		a[i] = *(*uint16)(unsafe.Pointer(p))
		p += 2
	}
	return string(utf16.Decode(a))
}

var variantClear = oleaut32.MustFindProc("VariantClear").Addr()

func VariantClear(v *stVariant) {
	syscall.SyscallN(variantClear, uintptr(unsafe.Pointer(v)))
}

var sysStringLen = oleaut32.MustFindProc("SysStringLen").Addr()

func SysStringLen(s uintptr) int {
	rtn, _, _ := syscall.SyscallN(sysStringLen, s)
	return int(rtn)
}
