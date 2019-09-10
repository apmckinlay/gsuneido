package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var uxtheme = windows.NewLazyDLL("uxtheme.dll")

// dll uxtheme:DrawThemeBackground(pointer hTheme, pointer hdc, long iPartId,
//		long iStateId, RECT* pRect, RECT* pClipRect) long
var drawThemeBackground = uxtheme.NewProc("DrawThemeBackground")
var _ = builtin6("DrawThemeBackground(hTheme, hdc, iPartId, iStateId, pRect,"+
	" pClipRect)",
	func(a, b, c, d, e, f Value) Value {
		var r1 RECT
		var r2 RECT
		rtn, _, _ := drawThemeBackground.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			uintptr(rectArg(e, &r1)),
			uintptr(rectArg(f, &r2)))
		return intRet(rtn)
	})

// dll uxtheme:DrawThemeText(pointer hTheme, pointer hdc, long iPartId,
//		long iStateId, string pszText, long iCharCount, long dwTextFlags,
//		long dwTextFlags2, RECT* pRect) long
var drawThemeText = uxtheme.NewProc("DrawThemeText")
var _ = builtin("DrawThemeText(hTheme, hdc, iPartId, iStateId, pszText,"+
	" iCharCount, dwTextFlags, dwTextFlags2, pRect)",
	func(_ *Thread, a []Value) Value {
		var r RECT
		rtn, _, _ := drawThemeText.Call(
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			uintptr(stringArg(a[4])),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			uintptr(rectArg(a[8], &r)))
		return intRet(rtn)
	})

// dll uxtheme:SetWindowTheme(pointer hwnd, string appname, string idlist) long
var setWindowTheme = uxtheme.NewProc("SetWindowTheme")
var _ = builtin3("SetWindowTheme(hwnd, appname, idlist)",
	func(a, b, c Value) Value {
		rtn, _, _ := setWindowTheme.Call(
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(stringArg(c)))
		return intRet(rtn)
	})

// dll uxtheme:GetWindowTheme(pointer hwnd) pointer
var getWindowTheme = uxtheme.NewProc("GetWindowTheme")
var _ = builtin1("GetWindowTheme(hwnd)",
	func(a Value) Value {
		rtn, _, _ := getWindowTheme.Call(intArg(a))
		return intRet(rtn)
	})

// dll uxtheme:IsAppThemed() bool
var isAppThemed = uxtheme.NewProc("IsAppThemed")
var _ = builtin0("IsAppThemed()",
	func() Value {
		rtn, _, _ := isAppThemed.Call()
		return boolRet(rtn)
	})

// dll long uxtheme:CloseThemeData(pointer hTheme)
var closeThemeData = uxtheme.NewProc("CloseThemeData")
var _ = builtin1("CloseThemeData(hTheme)",
	func(a Value) Value {
		rtn, _, _ := closeThemeData.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll bool UxTheme:IsThemeBackgroundPartiallyTransparent(
//		pointer hTheme,
//		long iPartId,
//		long iStateId)
var isThemeBackgroundPartiallyTransparent = uxtheme.NewProc("IsThemeBackgroundPartiallyTransparent")
var _ = builtin3("IsThemeBackgroundPartiallyTransparent(hTheme, iPartId, iStateId)",
	func(a, b, c Value) Value {
		rtn, _, _ := isThemeBackgroundPartiallyTransparent.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll pointer uxtheme:OpenThemeData(pointer hwnd, string pszClassList)
var openThemeData = uxtheme.NewProc("OpenThemeData")
var _ = builtin2("OpenThemeData(hwnd, pszClassList)",
	func(a, b Value) Value {
		rtn, _, _ := openThemeData.Call(
			intArg(a),
			uintptr(stringArg(b)))
		return intRet(rtn)
	})

// dll long UxTheme:DrawThemeParentBackground(pointer hwnd, pointer hdc,
//		RECT* prc)
var drawThemeParentBackground = uxtheme.NewProc("DrawThemeParentBackground")
var _ = builtin3("DrawThemeParentBackground(hwnd, hdc, prc)",
	func(a, b, c Value) Value {
		var r RECT
		rtn, _, _ := drawThemeParentBackground.Call(
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, &r)))
		return intRet(rtn)
	})
