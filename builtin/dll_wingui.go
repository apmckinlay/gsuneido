// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"bytes"
	"log"
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"golang.org/x/sys/windows"
)

// MustLoadDLL is like windows.MustLoadDLL but uses log.Fatalln instead of panic
func MustLoadDLL(name string) *mydll {
	d, err := windows.LoadDLL(name)
	if err != nil {
		log.Fatalln("FATAL:", err)
	}
	return (*mydll)(d)
}

type mydll windows.DLL

// MustFindProc is like windows.MustFindProc but uses log.Fatalln instead of panic
func (d *mydll) MustFindProc(name string) *windows.Proc {
	p, err := (*windows.DLL)(d).FindProc(name)
	if err != nil {
		log.Fatalln("FATAL:", err)
	}
	return p
}

var user32 = MustLoadDLL("user32.dll")

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

func int32Ret(rtn uintptr) Value {
	return IntVal(int(int32(rtn)))
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
	return heap.CopyStr(ToStr(v))
}

// bufToPtr copies data to an unsafe.Pointer
// WARNING: p must point to at least len(s) bytes
func bufToPtr(s string, p unsafe.Pointer) {
	for i := range len(s) {
		*(*byte)(unsafe.Pointer(uintptr(p) + uintptr(i))) = s[i]
	}
}

// strToPtr copies a nul terminated string to an unsafe.Pointer
// WARNING: p must point to at least len(s)+1 bytes
func strToPtr(s string, p unsafe.Pointer) {
	for i := range len(s) {
		*(*byte)(unsafe.Pointer(uintptr(p) + uintptr(i))) = s[i]
	}
	*(*byte)(unsafe.Pointer(uintptr(p) + uintptr(len(s)))) = 0
}

// bsStrZ copies a nul terminated string from a byte slice
func bsStrZ(buf []byte) Value {
	return SuStr(string(buf[:bytes.IndexByte(buf, 0)]))
}

// bufStrZ copies a nul terminated string from an unsafe.Pointer.
// If nul is not found, then the entire length is returned.
func bufStrZ(p unsafe.Pointer, n uintptr) Value {
	if p == nil || n == 0 {
		return False
	}
	i := uintptr(0)
	for ; i < n; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(p) + i)) == 0 {
			break
		}
	}
	return bufStrN(p, i)
}

// bufStrZ2 copies a *double* nul terminated string from a heap buffer.
// It includes the nuls in the result.
// If nuls are not found, then the entire length is returned.
func bufStrZ2(p unsafe.Pointer, n uintptr) Value {
	if p == nil || n == 0 {
		return EmptyStr
	}
	i := uintptr(2)
	for ; i < n; i++ {
		if *(*byte)(unsafe.Pointer(uintptr(p) + i - 2)) == 0 &&
			*(*byte)(unsafe.Pointer(uintptr(p) + i - 1)) == 0 {
			break
		}
	}
	return bufStrN(p, i)
}

// bufStrN copies a string of a given length from an unsafe.Pointer
func bufStrN(p unsafe.Pointer, n uintptr) Value {
	if p == nil || n == 0 {
		return EmptyStr
	}
	buf := make([]byte, n)
	for i := uintptr(0); i < n; i++ {
		buf[i] = *(*byte)(unsafe.Pointer(uintptr(p) + i))
	}
	return SuStr(hacks.BStoS(buf))
}

// getStrZbs copies the string into the byte slice and adds a nul terminator.
// If the string is too long, the excess is ignored
func getStrZbs(ob Value, mem string, dst []byte) {
	src := ToStr(ob.Get(nil, SuStr(mem)))
	if len(src) > len(dst)-1 {
		src = src[:len(dst)-1]
	}
	copy(dst, src)
	dst[len(src)] = 0
}

func zstrArg(v Value) *byte {
	// NOTE: don't change this to return uintptr
	// uintptr(unsafe.Pointer(x)) must be in the actual SyscallN arguments
	// Then it will be kept alive until the syscall returns.
	if v.Equal(Zero) {
		return nil
	}
	s := ToStr(v)
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	return &buf[0]
}

// rect -------------------------------------------------------------

func rectArg(ob Value, r unsafe.Pointer) unsafe.Pointer {
	//TODO if r is nil, alloc it
	if ob.Equal(Zero) {
		return nil
	}
	*(*stRect)(r) = obToRect(ob)
	return r
}

func obToRect(ob Value) stRect {
	return stRect{
		left:   getInt32(ob, "left"),
		top:    getInt32(ob, "top"),
		right:  getInt32(ob, "right"),
		bottom: getInt32(ob, "bottom"),
	}
}

func getRect(ob Value, mem string) stRect {
	if ob == nil {
		return stRect{}
	}
	x := ob.Get(nil, SuStr(mem))
	if x == nil {
		return stRect{}
	}
	return obToRect(x)
}

func urectToOb(p unsafe.Pointer, ob Value) Value {
	return rectToOb((*stRect)(p), ob)
}

func rectToOb(r *stRect, ob Value) Value {
	if ob == nil {
		ob = &SuObject{}
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

func obToPoint(ob Value) stPoint {
	if ob.Equal(Zero) {
		return stPoint{}
	}
	return stPoint{
		x: getInt32(ob, "x"),
		y: getInt32(ob, "y"),
	}
}

func upointToOb(p unsafe.Pointer, ob Value) Value {
	return pointToOb((*stPoint)(p), ob)
}

func pointToOb(pt *stPoint, ob Value) Value {
	if ob == nil {
		ob = &SuObject{}
	}
	ob.Put(nil, SuStr("x"), IntVal(int(pt.x)))
	ob.Put(nil, SuStr("y"), IntVal(int(pt.y)))
	return ob
}

func getPoint(ob Value, mem string) stPoint {
	if ob == nil {
		return stPoint{}
	}
	x := ob.Get(nil, SuStr(mem))
	if x == nil {
		return stPoint{}
	}
	return obToPoint(x)
}

func pointArg(ob Value, p unsafe.Pointer) unsafe.Pointer {
	pt := (*stPoint)(p)
	pt.x = getInt32(ob, "x")
	pt.y = getInt32(ob, "y")
	return p
}

//-------------------------------------------------------------------

func getUintptr(ob Value, mem string) uintptr {
	return uintptr(getInt(ob, mem))
}

func getCallback(th *Thread, ob Value, mem string, nargs int) uintptr {
	fn := ob.Get(nil, SuStr(mem))
	if fn == nil {
		return 0
	}
	return NewCallback(th, fn, nargs)
}
