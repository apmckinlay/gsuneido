// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"syscall"
	"unsafe"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var eztwain4 = windows.NewLazyDLL("eztwain4.dll")

// dll long Eztwain4:TWAIN_AcquireMultipageFile(pointer hwndApp, string pszFile)
var twain_AcquireMultipageFile = eztwain4.NewProc("TWAIN_AcquireMultipageFile")
var _ = builtin2("TWAIN_AcquireMultipageFile(hwndApp, pszFile)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn, _, _ := syscall.Syscall(twain_AcquireMultipageFile.Addr(), 2,
			intArg(a),
			uintptr(stringArg(b)),
			0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_GetDuplexSupport()
var twain_GetDuplexSupport = eztwain4.NewProc("TWAIN_GetDuplexSupport")
var _ = builtin0("TWAIN_GetDuplexSupport()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_GetDuplexSupport.Addr(), 0,
			0, 0, 0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_EnableDuplex(long i)
var twain_EnableDuplex = eztwain4.NewProc("TWAIN_EnableDuplex")
var _ = builtin1("TWAIN_EnableDuplex(enable)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(twain_EnableDuplex.Addr(), 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_GetHideUI()
var twain_GetHideUI = eztwain4.NewProc("TWAIN_GetHideUI")
var _ = builtin0("TWAIN_GetHideUI()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_GetHideUI.Addr(), 0,
			0, 0, 0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_GetSourceList()
var twain_GetSourceList = eztwain4.NewProc("TWAIN_GetSourceList")
var _ = builtin0("TWAIN_GetSourceList()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_GetSourceList.Addr(), 0,
			0, 0, 0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_HasControllableUI()
var twain_HasControllableUI = eztwain4.NewProc("TWAIN_HasControllableUI")
var _ = builtin0("TWAIN_HasControllableUI()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_HasControllableUI.Addr(), 0,
			0, 0, 0)
		return intRet(rtn)
	})

// dll bool Eztwain4:TWAIN_IsDuplexEnabled()
var twain_IsDuplexEnabled = eztwain4.NewProc("TWAIN_IsDuplexEnabled")
var _ = builtin0("TWAIN_IsDuplexEnabled()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_IsDuplexEnabled.Addr(), 0,
			0, 0, 0)
		return boolRet(rtn)
	})

// dll long Eztwain4:TWAIN_LastErrorCode()
var twain_LastErrorCode = eztwain4.NewProc("TWAIN_LastErrorCode")
var _ = builtin0("TWAIN_LastErrorCode()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_LastErrorCode.Addr(), 0,
			0, 0, 0)
		return intRet(rtn)
	})

// dll string Eztwain4:TWAIN_LastErrorText()
var twain_LastErrorText = eztwain4.NewProc("TWAIN_LastErrorText")
var _ = builtin0("TWAIN_LastErrorText()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_LastErrorText.Addr(), 0,
			0, 0, 0)
		return bufToStr(unsafe.Pointer(rtn), 256)
	})

// dll string Eztwain4:TWAIN_NextSourceName()
var twain_NextSourceName = eztwain4.NewProc("TWAIN_NextSourceName")
var _ = builtin0("TWAIN_NextSourceName()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_LastErrorText.Addr(), 0,
			0, 0, 0)
		return bufToStr(unsafe.Pointer(rtn), 256)
	})

// dll long Eztwain4:TWAIN_OpenDefaultSource()
var twain_OpenDefaultSource = eztwain4.NewProc("TWAIN_OpenDefaultSource")
var _ = builtin0("TWAIN_OpenDefaultSource()",
	func() Value {
		rtn, _, _ := syscall.Syscall(twain_OpenDefaultSource.Addr(), 0,
			0, 0, 0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_SelectImageSource(pointer hwnd)
var twain_SelectImageSource = eztwain4.NewProc("TWAIN_SelectImageSource")
var _ = builtin1("TWAIN_SelectImageSource(hwnd)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(twain_SelectImageSource.Addr(), 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll void Eztwain4:TWAIN_SetAppTitle(string title)
var twain_SetAppTitle = eztwain4.NewProc("TWAIN_SetAppTitle")
var _ = builtin1("TWAIN_SetAppTitle(title)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		syscall.Syscall(twain_SetAppTitle.Addr(), 1,
			uintptr(stringArg(a)),
			0, 0)
		return nil
	})

// dll long Eztwain4:TWAIN_SetHideUI(long fHide)
var twain_SetHideUI = eztwain4.NewProc("TWAIN_SetHideUI")
var _ = builtin1("TWAIN_SetHideUI(fHide)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(twain_SetHideUI.Addr(), 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_SetPaperSize(long nPaper)
var twain_SetPaperSize = eztwain4.NewProc("TWAIN_SetPaperSize")
var _ = builtin1("TWAIN_SetPaperSize(nPaper)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(twain_SetPaperSize.Addr(), 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll long Eztwain4:TWAIN_SetPixelType(long nPixType)
var twain_SetPixelType = eztwain4.NewProc("TWAIN_SetPixelType")
var _ = builtin1("TWAIN_SetPixelType(nPixType)",
	func(a Value) Value {
		rtn, _, _ := syscall.Syscall(twain_SetPixelType.Addr(), 1,
			intArg(a),
			0, 0)
		return intRet(rtn)
	})

// dll void Eztwain4:TWAIN_UniversalLicense(string pzVendorName, long nKey)
var twain_UniversalLicense = eztwain4.NewProc("TWAIN_UniversalLicense")
var _ = builtin2("TWAIN_UniversalLicense(pzVendorName, nKey)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		syscall.Syscall(twain_UniversalLicense.Addr(), 2,
			uintptr(stringArg(a)),
			intArg(b),
			0)
		return nil
	})

// dll int32 Eztwain4:TWAIN_SetResolution(double nRes)
var twain_SetResolution = eztwain4.NewProc("TWAIN_SetResolution")
var _ = builtin1("TWAIN_SetResolution(nRes)",
	func(a Value) Value {
		n := ToDnum(a).ToFloat()
		rtn, _, _ := syscall.Syscall(twain_SetResolution.Addr(), 1,
			// assuming Go and Windows use same float representation ???
			*(*uintptr)(unsafe.Pointer(&n)),
			0, 0)
		return intRet(rtn)
	})
