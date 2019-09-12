package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var eztwain4 = windows.NewLazyDLL("eztwain4.dll")

// dll long Eztwain4:TWAIN_AcquireMultipageFile(pointer hwndApp, string pszFile)
var twain_AcquireMultipageFile = eztwain4.NewProc("TWAIN_AcquireMultipageFile")
var _ = builtin2("TWAIN_AcquireMultipageFile(hwndApp, pszFile)",
	func(a, b Value) Value {
		rtn, _, _ := twain_AcquireMultipageFile.Call(
			intArg(a),
			uintptr(stringArg(b)))
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_GetHideUI()
var twain_GetHideUI = eztwain4.NewProc("TWAIN_GetHideUI")
var _ = builtin0("TWAIN_GetHideUI()",
	func() Value {
		rtn, _, _ := twain_GetHideUI.Call()
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_GetSourceList()
var twain_GetSourceList = eztwain4.NewProc("TWAIN_GetSourceList")
var _ = builtin0("TWAIN_GetSourceList()",
	func() Value {
		rtn, _, _ := twain_GetSourceList.Call()
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_HasControllableUI()
var twain_HasControllableUI = eztwain4.NewProc("TWAIN_HasControllableUI")
var _ = builtin0("TWAIN_HasControllableUI()",
	func() Value {
		rtn, _, _ := twain_HasControllableUI.Call()
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_LastErrorCode()
var twain_LastErrorCode = eztwain4.NewProc("TWAIN_LastErrorCode")
var _ = builtin0("TWAIN_LastErrorCode()",
	func() Value {
		rtn, _, _ := twain_LastErrorCode.Call()
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_OpenDefaultSource()
var twain_OpenDefaultSource = eztwain4.NewProc("TWAIN_OpenDefaultSource")
var _ = builtin0("TWAIN_OpenDefaultSource()",
	func() Value {
		rtn, _, _ := twain_OpenDefaultSource.Call()
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_SelectImageSource(pointer hwnd)
var twain_SelectImageSource = eztwain4.NewProc("TWAIN_SelectImageSource")
var _ = builtin1("TWAIN_SelectImageSource(hwnd)",
	func(a Value) Value {
		rtn, _, _ := twain_SelectImageSource.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll void Eztwain4:TWAIN_SetAppTitle(string title)
var twain_SetAppTitle = eztwain4.NewProc("TWAIN_SetAppTitle")
var _ = builtin1("TWAIN_SetAppTitle(title)",
	func(a Value) Value {
		twain_SetAppTitle.Call(
			uintptr(stringArg(a)))
		return nil
	})

// dll long Eztwain4:TWAIN_SetHideUI(long fHide)
var twain_SetHideUI = eztwain4.NewProc("TWAIN_SetHideUI")
var _ = builtin1("TWAIN_SetHideUI(fHide)",
	func(a Value) Value {
		rtn, _, _ := twain_SetHideUI.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_SetPaperSize(long nPaper)
var twain_SetPaperSize = eztwain4.NewProc("TWAIN_SetPaperSize")
var _ = builtin1("TWAIN_SetPaperSize(nPaper)",
	func(a Value) Value {
		rtn, _, _ := twain_SetPaperSize.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_SetPixelType(long nPixType)
var twain_SetPixelType = eztwain4.NewProc("TWAIN_SetPixelType")
var _ = builtin1("TWAIN_SetPixelType(nPixType)",
	func(a Value) Value {
		rtn, _, _ := twain_SetPixelType.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll void Eztwain4:TWAIN_UniversalLicense(string pzVendorName, long nKey)
var twain_UniversalLicense = eztwain4.NewProc("TWAIN_UniversalLicense")
var _ = builtin2("TWAIN_UniversalLicense(pzVendorName, nKey)",
	func(a, b Value) Value {
		twain_UniversalLicense.Call(
			uintptr(stringArg(a)),
			intArg(b))
		return nil
	})
