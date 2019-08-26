package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
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

// dll long Shell32:DragQueryFile(
//  pointer hDrop,
//  long iFile,
//  string lpszFile,
//  long cch)
var dragQueryFile = shell32.NewProc("DragQueryFile")
var _ = builtin2("DragQueryFile(hDrop, iFile)",
	func(a, b Value) Value {
		n, _, _ := dragQueryFile.Call(
			intArg(a),
			intArg(b),
			uintptr(0),
			uintptr(0))
		buf := make([]byte, n)
		dragQueryFile.Call(
			intArg(a),
			intArg(b),
			uintptr(unsafe.Pointer(&buf[0])),
			n)
		return SuStr(str.BeforeFirst(string(buf), "\x00"))
	})
