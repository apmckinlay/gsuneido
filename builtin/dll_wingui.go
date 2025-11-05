// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"golang.org/x/exp/constraints"
	"golang.org/x/sys/windows"
)

// ...Arg - convert from Suneido Value to native
//			e.g. intArg, zstrArg
// ...Ret - converts from native to a Suneido Value
//			e.g. intRet, int32Ret,
// to...  - converts a Suneido object to a native struct
//			e.g. toPoint, toRect
// from... - updates a Suneido object from a native struct
//			e.g. fromPoint, fromRect
// get... - gets a member from a Suneido Object and converts it to native
//			e.g. getBool, getInt

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

const uintptrMinusOne = ^uintptr(0) // -1

func truncToInt(x Value) int {
	if dn, ok := x.(SuDnum); ok {
		// seems like rounding would make more sense but cSuneido truncates
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

// see also zstrArg in sys_windows.go

// getZstr returns a nul terminated copy of a string member.
func getZstr(ob Value, mem string) *byte {
	x := ob.Get(nil, SuStr(mem))
	if x == nil || x.Equal(Zero) || x.Equal(False) {
		return nil
	}
	s := ToStr(x)
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	return &buf[0]
}

// bufZstr converts a byte slice containing a nul terminated string to SuStr.
func bufZstr(buf []byte) SuStr {
	if i := bytes.IndexByte(buf, 0); i >= 0 {
		buf = buf[:i]
	}
	return SuStr(string(buf))
}

// ptrZstr copies a nul terminated string from an unsafe.Pointer.
// If nul is not found, then the entire length is returned.
func ptrZstr(p unsafe.Pointer, n int) Value {
	if p == nil || n == 0 {
		return False
	}
	srcSlice := unsafe.Slice((*byte)(p), n)
	i := slices.Index(srcSlice, 0)
	if i == -1 {
		i = n // No null terminator found, use entire length
	}
	buf := make([]byte, i)
	copy(buf, srcSlice[:i])
	return SuStr(hacks.BStoS(buf))
}

// ptrNstr copies a string of a given length from an unsafe.Pointer
func ptrNstr(p unsafe.Pointer, n uintptr) Value {
	if p == nil || n == 0 {
		return EmptyStr
	}
	buf := make([]byte, n)
	srcSlice := unsafe.Slice((*byte)(p), n)
	copy(buf, srcSlice)
	return SuStr(hacks.BStoS(buf))
}

// getZstrBs copies the string into the byte slice and adds a nul terminator.
// If the string is too long, the excess is ignored
func getZstrBs(ob Value, mem string, dst []byte) {
	src := ToStr(ob.Get(nil, SuStr(mem)))
	if len(src) > len(dst)-1 {
		src = src[:len(dst)-1]
	}
	copy(dst, src)
	dst[min(len(src), len(dst)-1)] = 0
}

// rect -------------------------------------------------------------

// toRect converts an object to *stRect
func toRect(ob Value) *stRect {
	if ob.Equal(Zero) {
		return nil
	}
	return &stRect{
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
	return *toRect(x)
}

// fromRect updates an object from a *stRect
func fromRect(r *stRect, ob Value) Value {
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

func fromPoint(pt *stPoint, ob Value) Value {
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
	if x == nil || x.Equal(Zero) {
		return stPoint{}
	}
	return *toPoint(x)
}

func toPoint(ob Value) *stPoint {
	return &stPoint{
		x: getInt32(ob, "x"),
		y: getInt32(ob, "y"),
	}
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

func toptr[T constraints.Integer](x T) unsafe.Pointer {
	if ((uintptr)(x) < 0x10000) || ((uintptr)(x) > 0x00007FFFFFFFFFFF) {
		panic(fmt.Sprintf("invalid pointer value %x", x))
	}
	return unsafe.Pointer(uintptr(x)) //nolint
}
