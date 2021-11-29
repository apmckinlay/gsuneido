// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var comctl32 = MustLoadDLL("comctl32.dll")

var initCommonControlsEx = comctl32.MustFindProc("InitCommonControlsEx").Addr()
var _ = builtin1("InitCommonControlsEx(picce)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nINITCOMMONCONTROLSEX)
		*(*INITCOMMONCONTROLSEX)(p) = INITCOMMONCONTROLSEX{
			dwSize: uint32(nINITCOMMONCONTROLSEX),
			dwICC:  int32(getInt(a, "dwICC")),
		}
		rtn := goc.Syscall1(initCommonControlsEx,
			uintptr(p))
		return boolRet(rtn)
	})

type INITCOMMONCONTROLSEX struct {
	dwSize uint32
	dwICC  int32
}

const nINITCOMMONCONTROLSEX = unsafe.Sizeof(INITCOMMONCONTROLSEX{})

// dll Comctl32:ImageList_Create(
//		long x, long y, long flags, long initial, long grow) pointer
var imageList_Create = comctl32.MustFindProc("ImageList_Create").Addr()
var _ = builtin5("ImageList_Create(cx, cy, flags, cInitial, cGrow)",
	func(a, b, c, d, e Value) Value {
		rtn := goc.Syscall5(imageList_Create,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return intRet(rtn)
	})

// dll Comctl32:ImageList_Destroy(pointer himl) bool
var imageList_Destroy = comctl32.MustFindProc("ImageList_Destroy").Addr()
var _ = builtin1("ImageList_Destroy(himl)",
	func(a Value) Value {
		rtn := goc.Syscall1(imageList_Destroy,
			intArg(a))
		return boolRet(rtn)
	})

// dll Comctl32:ImageList_ReplaceIcon(pointer imagelist, long i, pointer hicon) long
var imageList_ReplaceIcon = comctl32.MustFindProc("ImageList_ReplaceIcon").Addr()
var _ = builtin3("ImageList_ReplaceIcon(himl, i, hicon)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(imageList_ReplaceIcon,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll bool Comctl32:ImageList_BeginDrag(
//		pointer himlTrack, long iTrack, long dxHotspot, long dyHotspot)
var imageList_BeginDrag = comctl32.MustFindProc("ImageList_BeginDrag").Addr()
var _ = builtin4("ImageList_BeginDrag(himlTrack, iTrack, dxHotspot, dyHotspot)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(imageList_BeginDrag,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return boolRet(rtn)
	})

// dll bool Comctl32:ImageList_DragEnter(pointer hwnd, long x, long y)
var imageList_DragEnter = comctl32.MustFindProc("ImageList_DragEnter").Addr()
var _ = builtin3("ImageList_DragEnter(hwnd, x, y)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(imageList_DragEnter,
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool Comctl32:ImageList_DragLeave(pointer hwnd)
var imageList_DragLeave = comctl32.MustFindProc("ImageList_DragLeave").Addr()
var _ = builtin1("ImageList_DragLeave(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(imageList_DragLeave,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool Comctl32:ImageList_DragMove(long x, long y)
var imageList_DragMove = comctl32.MustFindProc("ImageList_DragMove").Addr()
var _ = builtin2("ImageList_DragMove(x, y)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(imageList_DragMove,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll void Comctl32:ImageList_EndDrag()
var imageList_EndDrag = comctl32.MustFindProc("ImageList_EndDrag").Addr()
var _ = builtin0("ImageList_EndDrag()",
	func() Value {
		goc.Syscall0(imageList_EndDrag)
		return nil
	})

// dll pointer Comctl32:ImageList_Merge(pointer	himl1, long i1,
//		pointer himl2, long i2, long dx, long dy)
var imageList_Merge = comctl32.MustFindProc("ImageList_Merge").Addr()
var _ = builtin6("ImageList_Merge(himl1, i1, himl2, i2, dx, dy)",
	func(a, b, c, d, e, f Value) Value {
		rtn := goc.Syscall6(imageList_Merge,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f))
		return intRet(rtn)
	})

// dll long Comctl32:ImageList_Add(pointer imagelist, pointer image, pointer mask)
var imageList_Add = comctl32.MustFindProc("ImageList_Add").Addr()
var _ = builtin3("ImageList_Add(imagelist, image, mask)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(imageList_Add,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll long ComCtl32:ImageList_AddMasked(pointer himl, pointer hbmImage,
// 		long crMask)
var imageList_AddMasked = comctl32.MustFindProc("ImageList_AddMasked").Addr()
var _ = builtin3("ImageList_AddMasked(himl, hbmImage, crMask)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(imageList_AddMasked,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll void comctl32:DrawStatusText(pointer hdc, RECT* rect, [in] string text,
//		long flags)
var drawStatusText = comctl32.MustFindProc("DrawStatusTextA").Addr()
var _ = builtin4("DrawStatusText(himlTrack, iTrack, dxHotspot, dyHotspot)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		goc.Syscall4(drawStatusText,
			intArg(a),
			uintptr(rectArg(b, r)),
			uintptr(stringArg(c)),
			intArg(d))
		return nil
	})

// dll bool Comctl32:ImageList_GetImageInfo(pointer himl, long imageindex,
//		IMAGEINFO* pImageInfo)
var imageList_GetImageInfo = comctl32.MustFindProc("ImageList_GetImageInfo").Addr()
var _ = builtin3("ImageList_GetImageInfo(himl, imageindex, pImageInfo)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nIMAGEINFO)
		rtn := goc.Syscall3(imageList_GetImageInfo,
			intArg(a),
			intArg(b),
			uintptr(p))
		ii := *(*IMAGEINFO)(p)
		c.Put(nil, SuStr("hbmImage"), IntVal(int(ii.hbmImage)))
		c.Put(nil, SuStr("hbmMask"), IntVal(int(ii.hbmMask)))
		c.Put(nil, SuStr("rcImage"),
			rectToOb(&ii.rcImage, c.Get(nil, SuStr("rcImage"))))
		return boolRet(rtn)
	})

// dll bool Comctl32:ImageList_Draw(pointer himl, long imageindex,
//		pointer hdc, long x, long y, UINT fStyle)
var imageList_Draw = comctl32.MustFindProc("ImageList_Draw").Addr()
var _ = builtin6("ImageList_Draw(himl, imageindex, hdc, x, y, fStyle)",
	func(a, b, c, d, e, f Value) Value {
		rtn := goc.Syscall6(imageList_Draw,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f))
		return boolRet(rtn)
	})

type IMAGEINFO struct {
	hbmImage HANDLE
	hbmMask  HANDLE
	Unused1  int32
	Unused2  int32
	rcImage  RECT
}

const nIMAGEINFO = unsafe.Sizeof(IMAGEINFO{})
