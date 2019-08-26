package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var shell32 = windows.NewLazyDLL("shell32.dll")

// dll void Shell32:DragAcceptFiles(pointer hWnd, bool fAccept)
var dragAcceptFiles = shell32.NewProc("DragAcceptFiles")
var _ = builtin2("DragAcceptFiles(hWnd, fAccept)",
	func(a, b Value) Value {
		dragAcceptFiles.Call(
			intArg(a),
			boolArg(b))
		return nil
	})

// dll bool Shell32:SHGetPathFromIDList(pointer pidl, string path)
var shGetPathFromIDList = shell32.NewProc("SHGetPathFromIDList")
var _ = builtin2("SHGetPathFromIDList(pidl, path)",
	func(a, b Value) Value {
		rtn, _, _ := shGetPathFromIDList.Call(
			intArg(a),
			uintptr(stringArg(b)))
		return boolRet(rtn)
	})
