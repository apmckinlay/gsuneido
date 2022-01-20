// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !portable
// +build !portable

package builtin

import (
	"fmt"
	"math"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type suComObject struct {
	CantConvert
	ptr   uintptr
	idisp bool // true if IDispatch
}

var _ Value = (*suComObject)(nil)

var _ = builtin1("COMobject(progid)",
	func(arg Value) Value {
		if n, ok := arg.ToInt(); ok {
			ptr := uintptr(n)
			if idisp := goc.QueryIDispatch(ptr); idisp != 0 {
				goc.Release(ptr)
				return &suComObject{ptr: idisp, idisp: true}
			}
			return &suComObject{ptr: ptr}
		}
		if s, ok := arg.ToStr(); ok {
			defer heap.FreeTo(heap.CurSize())
			idisp := goc.CreateInstance(uintptr(heap.CopyStr(s)))
			if idisp == 0 {
				return False
			}
			return &suComObject{ptr: idisp, idisp: true}
		}
		panic("COMobject requires integer or string")
	})

var suComObjectMethods = Methods{
	"Dispatch?": method0(func(this Value) Value {
		return SuBool(this.(*suComObject).idisp)
	}),
	"Release": method0(func(this Value) Value {
		return IntVal(goc.Release(this.(*suComObject).ptr))
	}),
}

var _ Value = (*suComObject)(nil)

func (sco *suComObject) Get(_ *Thread, mem Value) Value {
	if !sco.idisp {
		panic("COMobject can't get property of IUnknown")
	}
	return GetProperty(sco.ptr, ToStr(mem))
}

func (sco *suComObject) Put(_ *Thread, mem Value, val Value) {
	if !sco.idisp {
		panic("COMobject can't put property of IUnknown")
	}
	PutProperty(sco.ptr, ToStr(mem), val)

}

func (sco *suComObject) GetPut(t *Thread, m Value, v Value,
	op func(x, y Value) Value, retOrig bool) Value {
	orig := sco.Get(t, m)
	v = op(orig, v)
	sco.Put(t, m, v)
	if retOrig {
		return orig
	}
	return v
}

func (*suComObject) RangeTo(int, int) Value {
	panic("COMobject does not support range")
}

func (*suComObject) RangeLen(int, int) Value {
	panic("COMobject does not support range")
}

func (*suComObject) Hash() uint32 {
	panic("COMobject hash not implemented")
}

func (*suComObject) Hash2() uint32 {
	panic("COMobject hash not implemented")
}

func (sco *suComObject) Equal(other interface{}) bool {
	sco2, ok := other.(*suComObject)
	return ok && sco == sco2
}

func (*suComObject) Compare(Value) int {
	panic("COMobject compare not implemented")
}

func (*suComObject) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call a COMobject instance")
}

func (*suComObject) String() string {
	return "COMobject"
}

func (*suComObject) Type() types.Type {
	return types.BuiltinClass
}

func (*suComObject) SetConcurrent() {
	// ok since immutable (assuming the COM object is thread safe)
}

func (sco *suComObject) Lookup(_ *Thread, method string) Callable {
	if f, ok := suComObjectMethods[method]; ok {
		return f
	}
	return &SuBuiltinMethodRaw{
		Fn: func(t *Thread, as *ArgSpec, this Value, args []Value) Value {
			return sco.call(method, as, args)
		}}
}

func (sco *suComObject) call(method string, as *ArgSpec, args []Value) Value {
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

type VARIANT struct {
	vt  uint16
	_   uint16
	_   uint16
	_   uint16
	val int64
	_   [8]byte
}

const nVARIANT = unsafe.Sizeof(VARIANT{})

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
	params := heap.Alloc(nDISPPARAMS)
	dp := (*DISPPARAMS)(params)
	dp.cArgs = uint32(len(args))
	dp.rgvarg = uintptr(pargs)
	var result unsafe.Pointer
	if flags == DISPATCH_PROPERTYPUT {
		dp.cNamedArgs = 1
		dp.rgdispidNamedArgs = uintptr(unsafe.Pointer(&DISPID_PROPERTYPUT))
	} else {
		result = heap.Alloc(nVARIANT)
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
	return variantToSu((*VARIANT)(result))
}

func convertArgs(args []Value) unsafe.Pointer {
	pargs := heap.Alloc(nVARIANT * uintptr(len(args)))
	p := pargs
	for i := len(args) - 1; i >= 0; i-- {
		suToVariant(args[i], (*VARIANT)(p))
		p = unsafe.Pointer(uintptr(p) + nVARIANT)
	}
	return pargs
}

var DISPID_PROPERTYPUT int32 = -3

func suToVariant(x Value, v *VARIANT) {
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
	} else if sco, ok := x.(*suComObject); ok {
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

type DISPPARAMS struct {
	rgvarg            uintptr
	rgdispidNamedArgs uintptr
	cArgs             uint32
	cNamedArgs        uint32
}

const nDISPPARAMS = unsafe.Sizeof(DISPPARAMS{})

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

func variantToSu(v *VARIANT) Value {
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
		result = &suComObject{ptr: uintptr(v.val), idisp: true}
	case VT_UNKNOWN:
		iunk := uintptr(v.val)
		if idisp := goc.QueryIDispatch(iunk); idisp != 0 {
			goc.Release(iunk)
			result = &suComObject{ptr: idisp, idisp: true}
		} else {
			result = &suComObject{ptr: iunk}
		}
	default:
		panic("COMobject: can't convert to Suneido value")
	}
	return result
}

func bstrToString(v *VARIANT) string {
	if v.val == 0 {
		return ""
	}
	p := uintptr(v.val)
	length := SysStringLen(p)
	a := make([]uint16, length)
	for i := 0; i < int(length); i++ {
		a[i] = *(*uint16)(unsafe.Pointer(p))
		p += 2
	}
	return string(utf16.Decode(a))
}

var variantClear = oleaut32.MustFindProc("VariantClear").Addr()

func VariantClear(v *VARIANT) {
	syscall.Syscall(variantClear, 1, uintptr(unsafe.Pointer(v)), 0, 0)
}

var sysStringLen = oleaut32.MustFindProc("SysStringLen").Addr()

func SysStringLen(s uintptr) int {
	rtn, _, _ := syscall.Syscall(sysStringLen, 1, s, 0, 0)
	return int(rtn)
}
