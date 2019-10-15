package builtin

import (
	"bytes"
	"unsafe"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var _ = windows.MustLoadDLL("scilexer.dll")

var user32 = windows.MustLoadDLL("user32.dll")

type HANDLE = uintptr
type BOOL = int32

const int32Size = unsafe.Sizeof(int32(0))
const int64Size = unsafe.Sizeof(int64(0))
const uintptrSize = unsafe.Sizeof(uintptr(0))

const uintptrMinusOne = ^uintptr(0) // -1

func truncToInt(x Value) int {
	// seems like rounding would make more sense but cSuneido truncates
	if dn, ok := x.ToDnum(); ok {
		if n, ok := dn.Trunc().ToInt(); ok {
			return n
		}
	}
	return ToInt(x)
}

func boolArg(arg Value) uintptr {
	if ToBool(arg) {
		return 1
	}
	return 0
}

func boolRet(rtn uintptr) Value {
	if rtn == 0 {
		return False
	}
	return True
}

func intArg(arg Value) uintptr {
	if arg.Equal(True) {
		return 1
	}
	if arg.Equal(False) {
		return 0
	}
	return uintptr(truncToInt(arg))
}

func intRet(rtn uintptr) Value {
	return IntVal(int(rtn))
}

func getBool(ob Value, mem string) BOOL {
	x := ob.Get(nil, SuStr(mem))
	if x == nil || !ToBool(x) {
		return 0
	}
	return 1
}

func getInt(ob Value, mem string) int {
	x := ob.Get(nil, SuStr(mem))
	if x == nil || x.Equal(False) {
		return 0
	}
	if x.Equal(True) {
		return 1
	}
	return truncToInt(x)
}

func getInt32(ob Value, mem string) int32 {
	if ob == nil {
		return 0
	}
	return int32(getInt(ob, mem))
}

func getUint32(ob Value, mem string) uint32 {
	return uint32(getInt(ob, mem))
}

func getInt16(ob Value, mem string) int16 {
	if ob == nil {
		return 0
	}
	return int16(getInt(ob, mem))
}

// string -----------------------------------------------------------

// getStr returns a nul terminated heap copy of a string member.
// Callers should defer heap.Free
func getStr(ob Value, mem string) *byte {
	x := ob.Get(nil, SuStr(mem))
	if x == nil || x.Equal(Zero) || x.Equal(False) {
		return nil
	}
	return (*byte)(stringArg(x))
}

// stringArg returns a nul terminated heap copy of a string.
// Callers should defer heap.Free
func stringArg(v Value) unsafe.Pointer {
	if v.Equal(Zero) {
		return nil
	}
	s := ToStr(v)
	return strToBuf(s, len(s)+1)
}

// strToBuf returns a nul terminated heap copy of a string of a specified size.
// Callers should defer heap.Free
func strToBuf(s string, n int) unsafe.Pointer {
	if len(s)+1 > n {
		panic("string too long")
	}
	p := heap.Alloc(uintptr(n))
	for i := 0; i < len(s); i++ {
		*(*byte)(unsafe.Pointer(uintptr(p) + uintptr(i))) = s[i]
	}
	*(*byte)(unsafe.Pointer(uintptr(p) + uintptr(len(s)))) = 0
	return p
}

func strRet(buf []byte) Value {
	return SuStr(string(buf[:bytes.IndexByte(buf, 0)]))
}

// bufToStr copies a nul terminated string from an unsafe.Pointer.
// If nul is not found, then the entire length is returned.
func bufToStr(p unsafe.Pointer, n uintptr) Value {
	i := uintptr(0)
	for ; i < n; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(p) + i)) == 0 {
			break
		}
	}
	return bufRet(p, i)
}

// bufToStr2 copies a *double* nul terminated string from a heap buffer.
// It includes the nuls in the result.
// If nuls are not found, then the entire length is returned.
func bufToStr2(p unsafe.Pointer, n uintptr) Value {
	i := uintptr(2)
	for ; i < n; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(p) + i - 2)) == 0 &&
			*(*byte)(unsafe.Pointer(uintptr(p) + i - 1)) == 0 {
			break
		}
	}
	return bufRet(p, i)
}

// bufRet returns a string from a pointer and length (not nul terminated)
func bufRet(p unsafe.Pointer, n uintptr) Value {
	if p == nil {
		return EmptyStr
	}
	// same approach as strings.Builder to reuse []byte as string (without copy)
	buf := make([]byte, n)
	for i := uintptr(0); i < n; i++ {
		buf[i] = *(*byte)(unsafe.Pointer(uintptr(p) + i))
	}
	return SuStr(*(*string)(unsafe.Pointer(&buf)))
}

// copyStr copies the string into the byte slice and adds a nul terminator
func copyStr(dst []byte, ob Value, mem string) {
	src := ToStr(ob.Get(nil, SuStr(mem)))
	copy(dst[:], src)
	dst[len(src)] = 0
}

// rect -------------------------------------------------------------

func rectArg(ob Value, r unsafe.Pointer) unsafe.Pointer {
	//TODO if r is nil, alloc it
	if ob.Equal(Zero) {
		return nil
	}
	*(*RECT)(r) = obToRect(ob)
	return r
}

func obToRect(ob Value) RECT {
	return RECT{
		left:   getInt32(ob, "left"),
		top:    getInt32(ob, "top"),
		right:  getInt32(ob, "right"),
		bottom: getInt32(ob, "bottom"),
	}
}

func getRect(ob Value, mem string) RECT {
	x := ob.Get(nil, SuStr(mem))
	if x == nil {
		return RECT{}
	}
	return obToRect(x)
}

func urectToOb(p unsafe.Pointer, ob Value) Value {
	return rectToOb((*RECT)(p), ob)
}

func rectToOb(r *RECT, ob Value) Value {
	if ob == nil {
		ob = NewSuObject()
	} else if ob.Equal(Zero) {
		return ob
	}
	ob.Put(nil, SuStr("left"), IntVal(int(r.left)))
	ob.Put(nil, SuStr("top"), IntVal(int(r.top)))
	ob.Put(nil, SuStr("right"), IntVal(int(r.right)))
	ob.Put(nil, SuStr("bottom"), IntVal(int(r.bottom)))
	return ob
}

// point ------------------------------------------------------------

func obToPoint(ob Value) POINT {
	if ob.Equal(Zero) {
		return POINT{}
	}
	return POINT{
		x: getInt32(ob, "x"),
		y: getInt32(ob, "y"),
	}
}

func upointToOb(p unsafe.Pointer, ob Value) Value {
	return pointToOb((*POINT)(p), ob)
}

func pointToOb(pt *POINT, ob Value) Value {
	if ob == nil {
		ob = NewSuObject()
	}
	ob.Put(nil, SuStr("x"), IntVal(int(pt.x)))
	ob.Put(nil, SuStr("y"), IntVal(int(pt.y)))
	return ob
}

func getPoint(ob Value, mem string) POINT {
	x := ob.Get(nil, SuStr(mem))
	if x == nil {
		return POINT{}
	}
	return obToPoint(x)
}

func pointArg(ob Value, p unsafe.Pointer) unsafe.Pointer {
	pt := (*POINT)(p)
	pt.x = getInt32(ob, "x")
	pt.y = getInt32(ob, "y")
	return p
}

//-------------------------------------------------------------------

func getHandle(ob Value, mem string) HANDLE {
	return HANDLE(getInt(ob, mem))
}

func getCallback(ob Value, mem string, nargs byte) uintptr {
	fn := ob.Get(nil, SuStr(mem))
	if fn == nil {
		return 0
	}
	return NewCallback(fn, nargs)
}
