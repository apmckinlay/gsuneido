// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"syscall"

	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
)

var uxtheme = MustLoadDLL("uxtheme.dll")

// dll uxtheme:DrawThemeBackground(pointer hTheme, pointer hdc, long iPartId,
// long iStateId, RECT* pRect, RECT* pClipRect) long
var drawThemeBackground = uxtheme.MustFindProc("DrawThemeBackground").Addr()
var _ = builtin(DrawThemeBackground,
	"(hTheme, hdc, iPartId, iStateId, pRect, pClipRect)")

func DrawThemeBackground(a, b, c, d, e, f Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r1 := heap.Alloc(nRect)
	r2 := heap.Alloc(nRect)
	rtn, _, _ := syscall.SyscallN(drawThemeBackground,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		uintptr(rectArg(e, r1)),
		uintptr(rectArg(f, r2)))
	return intRet(rtn)
}

// DrawThemeText(pointer hTheme, pointer hdc, long iPartId,
// long iStateId, string pszText, long iCharCount, long dwTextFlags,
// long dwTextFlags2, RECT* pRect) long
var drawThemeText = uxtheme.MustFindProc("DrawThemeText").Addr()
var _ = builtin(DrawThemeText, "(hTheme, hdc, iPartId, iStateId, pszText, "+
	"iCharCount, dwTextFlags, dwTextFlags2, pRect)")

func DrawThemeText(_ *Thread, a []Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn, _, _ := syscall.SyscallN(drawThemeText,
		intArg(a[0]),
		intArg(a[1]),
		intArg(a[2]),
		intArg(a[3]),
		uintptr(stringArg(a[4])),
		intArg(a[5]),
		intArg(a[6]),
		intArg(a[7]),
		uintptr(rectArg(a[8], r)))
	return intRet(rtn)
}

// dll uxtheme:SetWindowTheme(pointer hwnd, string appname, string idlist) long
var setWindowTheme = uxtheme.MustFindProc("SetWindowTheme").Addr()
var _ = builtin(SetWindowTheme, "(hwnd, appname, idlist)")

func SetWindowTheme(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn, _, _ := syscall.SyscallN(setWindowTheme,
		intArg(a),
		uintptr(stringArg(b)),
		uintptr(stringArg(c)))
	return intRet(rtn)
}

// dll uxtheme:GetWindowTheme(pointer hwnd) pointer
var getWindowTheme = uxtheme.MustFindProc("GetWindowTheme").Addr()
var _ = builtin(GetWindowTheme, "(hwnd)")

func GetWindowTheme(a Value) Value {
	rtn, _, _ := syscall.SyscallN(getWindowTheme,
		intArg(a))
	return intRet(rtn)
}

// dll uxtheme:IsAppThemed() bool
var isAppThemed = uxtheme.MustFindProc("IsAppThemed").Addr()
var _ = builtin(IsAppThemed, "()")

func IsAppThemed() Value {
	rtn, _, _ := syscall.SyscallN(isAppThemed)
	return boolRet(rtn)
}

// dll long uxtheme:CloseThemeData(pointer hTheme)
var closeThemeData = uxtheme.MustFindProc("CloseThemeData").Addr()
var _ = builtin(CloseThemeData, "(hTheme)")

func CloseThemeData(a Value) Value {
	rtn, _, _ := syscall.SyscallN(closeThemeData,
		intArg(a))
	return intRet(rtn)
}

// dll bool UxTheme:IsThemeBackgroundPartiallyTransparent(
// pointer hTheme, long iPartId, long iStateId)
var isThemeBackgroundPartiallyTransparent = uxtheme.MustFindProc("IsThemeBackgroundPartiallyTransparent").Addr()
var _ = builtin(IsThemeBackgroundPartiallyTransparent, "(hTheme, iPartId, iStateId)")

func IsThemeBackgroundPartiallyTransparent(a, b, c Value) Value {
	rtn, _, _ := syscall.SyscallN(isThemeBackgroundPartiallyTransparent,
		intArg(a),
		intArg(b),
		intArg(c))
	return boolRet(rtn)
}

// dll pointer uxtheme:OpenThemeData(pointer hwnd, string pszClassList)
var openThemeData = uxtheme.MustFindProc("OpenThemeData").Addr()
var _ = builtin(OpenThemeData, "(hwnd, pszClassList)")

func OpenThemeData(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn, _, _ := syscall.SyscallN(openThemeData,
		intArg(a),
		uintptr(stringArg(b)))
	return intRet(rtn)
}

// dll long UxTheme:DrawThemeParentBackground(pointer hwnd, pointer hdc,
// RECT* prc)
var drawThemeParentBackground = uxtheme.MustFindProc("DrawThemeParentBackground").Addr()
var _ = builtin(DrawThemeParentBackground, "(hwnd, hdc, prc)")

func DrawThemeParentBackground(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn, _, _ := syscall.SyscallN(drawThemeParentBackground,
		intArg(a),
		intArg(b),
		uintptr(rectArg(c, r)))
	return intRet(rtn)
}
