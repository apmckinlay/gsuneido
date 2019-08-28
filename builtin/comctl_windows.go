package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var comctl32 = windows.NewLazyDLL("comctl32.dll")

var initCommonControlsEx = comctl32.NewProc("InitCommonControlsEx")
var _ = builtin1("InitCommonControlsEx(picce)",
	func(a Value) Value {
		a1 := INITCOMMONCONTROLSEX{
			dwSize: int32(unsafe.Sizeof(INITCOMMONCONTROLSEX{})),
			dwICC:  int32(getInt(a, "dwICC")),
		}
		rtn, _, _ := initCommonControlsEx.Call(uintptr(unsafe.Pointer(&a1)))
		if rtn == 0 {
			return False
		}
		return True
	})

type INITCOMMONCONTROLSEX struct {
	dwSize int32
	dwICC  int32
}

// dll Comctl32:ImageList_Create(
//		long x, long y, long flags, long initial, long grow) pointer
var imageList_Create = user32.NewProc("ImageList_Create")
var _ = builtin5("ImageList_Create(cx, cy, flags, cInitial, cGrow)",
	func(a, b, c, d, e Value) Value {
		rtn, _, _ := imageList_Create.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return intRet(rtn)
	})

// dll Comctl32:ImageList_Destroy(pointer himl) bool
var imageList_Destroy = comctl32.NewProc("ImageList_Destroy")
var _ = builtin1("ImageList_Destroy(himl)",
	func(a Value) Value {
		rtn, _, _ := initCommonControlsEx.Call(intArg(a))
		if rtn == 0 {
			return False
		}
		return True
	})

// dll Comctl32:ImageList_ReplaceIcon(pointer imagelist, long i, pointer hicon) long
var imageList_ReplaceIcon = user32.NewProc("ImageList_ReplaceIcon")
var _ = builtin3("ImageList_ReplaceIcon(himl, i, hicon)",
	func(a, b, c Value) Value {
		rtn, _, _ := imageList_ReplaceIcon.Call(intArg(a), intArg(b), intArg(c))
		return intRet(rtn)
	})

// dll bool Comctl32:ImageList_BeginDrag(
//		pointer himlTrack, long iTrack, long dxHotspot, long dyHotspot)
var imageList_BeginDrag = comctl32.NewProc("ImageList_BeginDrag")
var _ = builtin4("ImageList_BeginDrag(himlTrack, iTrack, dxHotspot, dyHotspot)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := imageList_BeginDrag.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return boolRet(rtn)
	})

// dll bool Comctl32:ImageList_DragEnter(pointer hwnd, long x, long y)
var imageList_DragEnter = comctl32.NewProc("ImageList_DragEnter")
var _ = builtin3("ImageList_DragEnter(hwnd, x, y)",
	func(a, b, c Value) Value {
		rtn, _, _ := imageList_DragEnter.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool Comctl32:ImageList_DragLeave(pointer hwnd)
var imageList_DragLeave = comctl32.NewProc("ImageList_DragLeave")
var _ = builtin1("ImageList_DragLeave(hwnd)",
	func(a Value) Value {
		rtn, _, _ := imageList_DragLeave.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool Comctl32:ImageList_DragMove(long x, long y)
var imageList_DragMove = comctl32.NewProc("ImageList_DragMove")
var _ = builtin2("ImageList_DragMove(x, y)",
	func(a, b Value) Value {
		rtn, _, _ := imageList_DragMove.Call(
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll void Comctl32:ImageList_EndDrag()
var imageList_EndDrag = comctl32.NewProc("ImageList_EndDrag")
var _ = builtin0("ImageList_EndDrag()",
	func() Value {
		imageList_EndDrag.Call()
		return nil
	})

// dll pointer Comctl32:ImageList_Merge(pointer	himl1, long i1,
//		pointer himl2, long i2, long dx, long dy)
var imageList_Merge = comctl32.NewProc("ImageList_Merge")
var _ = builtin6("ImageList_Merge(himl1, i1, himl2, i2, dx, dy)",
	func(a, b, c, d, e, f Value) Value {
		rtn, _, _ := imageList_Merge.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f))
		return intRet(rtn)
	})

// dll void comctl32:DrawStatusText(pointer hdc, RECT* rect, [in] string text,
//		long flags)
var drawStatusText = comctl32.NewProc("DrawStatusTextA")
var _ = builtin4("DrawStatusText(himlTrack, iTrack, dxHotspot, dyHotspot)",
	func(a, b, c, d Value) Value {
		var r RECT
		drawStatusText.Call(
			intArg(a),
			uintptr(rectArg(b, &r)),
			uintptr(stringArg(c)),
			intArg(d))
		return nil
	})
