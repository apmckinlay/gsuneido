package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var comctl32 = windows.NewLazyDLL("comctl32.dll")

type INITCOMMONCONTROLSEX struct {
	dwSize int32
	dwICC  int32
}

var initCommonControlsEx = comctl32.NewProc("InitCommonControlsEx")
var _ = builtin1("InitCommonControlsEx(picce)", func(a Value) Value {
	a1 := INITCOMMONCONTROLSEX{
		dwSize: int32(getInt(a, "dwSize")),
		dwICC:  int32(getInt(a, "dwICC")),
	}
	rtn, _, _ := initCommonControlsEx.Call(uintptr(unsafe.Pointer(&a1)))
	if rtn == 0 {
		return False
	}
	return True
})

// dll Comctl32:ImageList_Create(long x, long y, long flags, long initial, long grow) pointer
var imageList_Create = user32.NewProc("ImageList_Create")
var _ = builtin5("ImageList_Create(cx, cy, flags, cInitial, cGrow)", func(a, b, c, d, e Value) Value {
	rtn, _, _ := imageList_Create.Call(
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e))
	return IntVal(int(rtn))
})

// dll Comctl32:ImageList_Destroy(pointer himl) bool
var imageList_Destroy = comctl32.NewProc("ImageList_Destroy")
var _ = builtin1("ImageList_Destroy(himl)", func(a Value) Value {
	rtn, _, _ := initCommonControlsEx.Call(intArg(a))
	if rtn == 0 {
		return False
	}
	return True
})

// dll Comctl32:ImageList_ReplaceIcon(pointer imagelist, long i, pointer hicon) long
// Returns the index of the image if successful, or -1 otherwise
var imageList_ReplaceIcon = user32.NewProc("ImageList_ReplaceIcon")
var _ = builtin3("ImageList_ReplaceIcon(himl, i, hicon)", func(a, b, c Value) Value {
	rtn, _, _ := imageList_ReplaceIcon.Call(intArg(a), intArg(b), intArg(c))
	return IntVal(int(rtn))
})
