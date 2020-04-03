// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var uxtheme = MustLoadDLL("uxtheme.dll")

// dll uxtheme:DrawThemeBackground(pointer hTheme, pointer hdc, long iPartId,
//		long iStateId, RECT* pRect, RECT* pClipRect) long
var drawThemeBackground = uxtheme.MustFindProc("DrawThemeBackground").Addr()
var _ = builtin6("DrawThemeBackground(hTheme, hdc, iPartId, iStateId, pRect,"+
	" pClipRect)",
	func(a, b, c, d, e, f Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r1 := heap.Alloc(nRECT)
		r2 := heap.Alloc(nRECT)
		rtn := goc.Syscall6(drawThemeBackground,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			uintptr(rectArg(e, r1)),
			uintptr(rectArg(f, r2)))
		return intRet(rtn)
	})

// dll uxtheme:DrawThemeText(pointer hTheme, pointer hdc, long iPartId,
//		long iStateId, string pszText, long iCharCount, long dwTextFlags,
//		long dwTextFlags2, RECT* pRect) long
var drawThemeText = uxtheme.MustFindProc("DrawThemeText").Addr()
var _ = builtin("DrawThemeText(hTheme, hdc, iPartId, iStateId, pszText,"+
	" iCharCount, dwTextFlags, dwTextFlags2, pRect)",
	func(_ *Thread, a []Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall9(drawThemeText,
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
	})

// dll uxtheme:SetWindowTheme(pointer hwnd, string appname, string idlist) long
var setWindowTheme = uxtheme.MustFindProc("SetWindowTheme").Addr()
var _ = builtin3("SetWindowTheme(hwnd, appname, idlist)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall3(setWindowTheme,
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(stringArg(c)))
		return intRet(rtn)
	})

// dll uxtheme:GetWindowTheme(pointer hwnd) pointer
var getWindowTheme = uxtheme.MustFindProc("GetWindowTheme").Addr()
var _ = builtin1("GetWindowTheme(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(getWindowTheme,
			intArg(a))
		return intRet(rtn)
	})

// dll uxtheme:IsAppThemed() bool
var isAppThemed = uxtheme.MustFindProc("IsAppThemed").Addr()
var _ = builtin0("IsAppThemed()",
	func() Value {
		rtn := goc.Syscall0(isAppThemed)
		return boolRet(rtn)
	})

// dll long uxtheme:CloseThemeData(pointer hTheme)
var closeThemeData = uxtheme.MustFindProc("CloseThemeData").Addr()
var _ = builtin1("CloseThemeData(hTheme)",
	func(a Value) Value {
		rtn := goc.Syscall1(closeThemeData,
			intArg(a))
		return intRet(rtn)
	})

// dll bool UxTheme:IsThemeBackgroundPartiallyTransparent(
//		pointer hTheme,
//		long iPartId,
//		long iStateId)
var isThemeBackgroundPartiallyTransparent = uxtheme.MustFindProc("IsThemeBackgroundPartiallyTransparent").Addr()
var _ = builtin3("IsThemeBackgroundPartiallyTransparent(hTheme, iPartId, iStateId)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(isThemeBackgroundPartiallyTransparent,
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll pointer uxtheme:OpenThemeData(pointer hwnd, string pszClassList)
var openThemeData = uxtheme.MustFindProc("OpenThemeData").Addr()
var _ = builtin2("OpenThemeData(hwnd, pszClassList)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall2(openThemeData,
			intArg(a),
			uintptr(stringArg(b)))
		return intRet(rtn)
	})

// dll long UxTheme:DrawThemeParentBackground(pointer hwnd, pointer hdc,
//		RECT* prc)
var drawThemeParentBackground = uxtheme.MustFindProc("DrawThemeParentBackground").Addr()
var _ = builtin3("DrawThemeParentBackground(hwnd, hdc, prc)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall3(drawThemeParentBackground,
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, r)))
		return intRet(rtn)
	})
