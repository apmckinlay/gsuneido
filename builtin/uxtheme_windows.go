package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var uxtheme = windows.NewLazyDLL("uxtheme.dll")

// dll uxtheme:DrawThemeBackground(pointer hTheme, pointer hdc, long iPartId,
// long iStateId, RECT* pRect, RECT* pClipRect) long
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
			rectArg(e, &r1),
			rectArg(f, &r2))
		return IntVal(int(rtn))
	})

// dll uxtheme:DrawThemeText(pointer hTheme, pointer hdc, long iPartId,
// long iStateId, string pszText, long iCharCount, long dwTextFlags,
// long dwTextFlags2, RECT* pRect) long
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
			stringArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			rectArg(a[8], &r))
		return IntVal(int(rtn))
	})

// dll uxtheme:SetWindowTheme(pointer hwnd, string appname, string idlist) long
var setWindowTheme = uxtheme.NewProc("SetWindowTheme")
var _ = builtin3("SetWindowTheme(hwnd, appname, idlist)",
	func(a, b, c Value) Value {
		rtn, _, _ := setWindowTheme.Call(
			intArg(a),
			stringArg(b),
			stringArg(c))
		return IntVal(int(rtn))
	})

// dll uxtheme:GetWindowTheme(pointer hwnd) pointer
var getWindowTheme = uxtheme.NewProc("GetWindowTheme")
var _ = builtin1("GetWindowTheme(hwnd)",
	func(a Value) Value {
		rtn, _, _ := getWindowTheme.Call(intArg(a))
		return IntVal(int(rtn))
	})

// dll uxtheme:IsAppThemed() bool
var isAppThemed = uxtheme.NewProc("IsAppThemed")
var _ = builtin0("IsAppThemed()",
	func() Value {
		rtn, _, _ := isAppThemed.Call()
		return boolRet(rtn)
	})
