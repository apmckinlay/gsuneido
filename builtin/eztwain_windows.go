// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

// import (
// 	. "github.com/apmckinlay/gsuneido/runtime"
// 	"golang.org/x/sys/windows"
// )

// var eztwain4 = MustLoadDLL("eztwain4.dll")

// // dll long Eztwain4:TWAIN_AcquireMultipageFile(pointer hwndApp, string pszFile)
// var twain_AcquireMultipageFile = eztwain4.MustFindProc("TWAIN_AcquireMultipageFile").Addr()
// var _ = builtin2("TWAIN_AcquireMultipageFile(hwndApp, pszFile)",
// 	func(a, b Value) Value {
// 		rtn, _, _ := syscall.Syscall(twain_AcquireMultipageFile, ?,
// 			intArg(a),
// 			uintptr(stringArg(b)))
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_GetHideUI()
// var twain_GetHideUI = eztwain4.MustFindProc("TWAIN_GetHideUI").Addr()
// var _ = builtin0("TWAIN_GetHideUI()",
// 	func() Value {
// 		rtn, _, _ := syscall.Syscall(twain_GetHideUI, 0)
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_GetSourceList()
// var twain_GetSourceList = eztwain4.MustFindProc("TWAIN_GetSourceList").Addr()
// var _ = builtin0("TWAIN_GetSourceList()",
// 	func() Value {
// 		rtn, _, _ := syscall.Syscall(twain_GetSourceList, 0)
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_HasControllableUI()
// var twain_HasControllableUI = eztwain4.MustFindProc("TWAIN_HasControllableUI").Addr()
// var _ = builtin0("TWAIN_HasControllableUI()",
// 	func() Value {
// 		rtn, _, _ := syscall.Syscall(twain_HasControllableUI, 0)
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_LastErrorCode()
// var twain_LastErrorCode = eztwain4.MustFindProc("TWAIN_LastErrorCode").Addr()
// var _ = builtin0("TWAIN_LastErrorCode()",
// 	func() Value {
// 		rtn, _, _ := syscall.Syscall(twain_LastErrorCode, 0)
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_OpenDefaultSource()
// var twain_OpenDefaultSource = eztwain4.MustFindProc("TWAIN_OpenDefaultSource").Addr()
// var _ = builtin0("TWAIN_OpenDefaultSource()",
// 	func() Value {
// 		rtn, _, _ := syscall.Syscall(twain_OpenDefaultSource, 0)
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_SelectImageSource(pointer hwnd)
// var twain_SelectImageSource = eztwain4.MustFindProc("TWAIN_SelectImageSource").Addr()
// var _ = builtin1("TWAIN_SelectImageSource(hwnd)",
// 	func(a Value) Value {
// 		rtn, _, _ := syscall.Syscall(twain_SelectImageSource, ?,
// 			intArg(a))
// 		return intRet(rtn)
// 	})

// // dll void Eztwain4:TWAIN_SetAppTitle(string title)
// var twain_SetAppTitle = eztwain4.MustFindProc("TWAIN_SetAppTitle").Addr()
// var _ = builtin1("TWAIN_SetAppTitle(title)",
// 	func(a Value) Value {
// 		syscall.Syscall(twain_SetAppTitle, ?,
// 			uintptr(stringArg(a)))
// 		return nil
// 	})

// // dll long Eztwain4:TWAIN_SetHideUI(long fHide)
// var twain_SetHideUI = eztwain4.MustFindProc("TWAIN_SetHideUI").Addr()
// var _ = builtin1("TWAIN_SetHideUI(fHide)",
// 	func(a Value) Value {
// 		rtn, _, _ := syscall.Syscall(twain_SetHideUI, ?,
// 			intArg(a))
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_SetPaperSize(long nPaper)
// var twain_SetPaperSize = eztwain4.MustFindProc("TWAIN_SetPaperSize").Addr()
// var _ = builtin1("TWAIN_SetPaperSize(nPaper)",
// 	func(a Value) Value {
// 		rtn, _, _ := syscall.Syscall(twain_SetPaperSize, ?,
// 			intArg(a))
// 		return intRet(rtn)
// 	})

// // dll long Eztwain4:TWAIN_SetPixelType(long nPixType)
// var twain_SetPixelType = eztwain4.MustFindProc("TWAIN_SetPixelType").Addr()
// var _ = builtin1("TWAIN_SetPixelType(nPixType)",
// 	func(a Value) Value {
// 		rtn, _, _ := syscall.Syscall(twain_SetPixelType, ?,
// 			intArg(a))
// 		return intRet(rtn)
// 	})

// // dll void Eztwain4:TWAIN_UniversalLicense(string pzVendorName, long nKey)
// var twain_UniversalLicense = eztwain4.MustFindProc("TWAIN_UniversalLicense").Addr()
// var _ = builtin2("TWAIN_UniversalLicense(pzVendorName, nKey)",
// 	func(a, b Value) Value {
// 		syscall.Syscall(twain_UniversalLicense, ?,
// 			uintptr(stringArg(a)),
// 			intArg(b))
// 		return nil
// 	})
