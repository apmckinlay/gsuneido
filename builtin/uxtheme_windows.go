package builtin

import (
	"syscall"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var uxtheme = windows.MustLoadDLL("uxtheme.dll")

// dll uxtheme:DrawThemeBackground(pointer hTheme, pointer hdc, long iPartId,
//		long iStateId, RECT* pRect, RECT* pClipRect) long
var drawThemeBackground = uxtheme.MustFindProc("DrawThemeBackground").Addr()
var _ = builtin6("DrawThemeBackground(hTheme, hdc, iPartId, iStateId, pRect,"+
	" pClipRect)",
	func(a, b, c, d, e, f Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r1 := heap.Alloc(nRECT)
		r2 := heap.Alloc(nRECT)
		rtn, _, _ := syscall.Syscall6(drawThemeBackground, 6,
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
		rtn, _, _ := syscall.Syscall9(drawThemeText, 9,
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
		rtn, _, _ := syscall.Syscall(setWindowTheme, 3,
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(stringArg(c)))
		return intRet(rtn)
	})

// dll uxtheme:GetWindowTheme(pointer hwnd) pointer
var getWindowTheme = uxtheme.MustFindProc("GetWindowTheme").Addr()
var _ = builtin1("GetWindowTheme(hwnd)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(getWindowTheme, 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll uxtheme:IsAppThemed() bool
var isAppThemed = uxtheme.MustFindProc("IsAppThemed").Addr()
var _ = builtin0("IsAppThemed()",
	func() Value {
		rtn, _, _ := syscall.Syscall(isAppThemed, 0, 0, 0, 0)
		return boolRet(rtn)
	})

// dll long uxtheme:CloseThemeData(pointer hTheme)
var closeThemeData = uxtheme.MustFindProc("CloseThemeData").Addr()
var _ = builtin1("CloseThemeData(hTheme)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(closeThemeData, 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll bool UxTheme:IsThemeBackgroundPartiallyTransparent(
//		pointer hTheme,
//		long iPartId,
//		long iStateId)
var isThemeBackgroundPartiallyTransparent = uxtheme.MustFindProc("IsThemeBackgroundPartiallyTransparent").Addr()
var _ = builtin3("IsThemeBackgroundPartiallyTransparent(hTheme, iPartId, iStateId)",
	func(a, b, c Value) Value {
		rtn, _, _ := syscall.Syscall(isThemeBackgroundPartiallyTransparent, 3,
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
		rtn, _, _ := syscall.Syscall(openThemeData, 2,
			intArg(a),
			uintptr(stringArg(b)),
			0)
		return intRet(rtn)
	})

// dll long UxTheme:DrawThemeParentBackground(pointer hwnd, pointer hdc,
//		RECT* prc)
var drawThemeParentBackground = uxtheme.MustFindProc("DrawThemeParentBackground").Addr()
var _ = builtin3("DrawThemeParentBackground(hwnd, hdc, prc)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn, _, _ := syscall.Syscall(drawThemeParentBackground, 3,
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, r)))
		return intRet(rtn)
	})
