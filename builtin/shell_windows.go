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
var shGetPathFromIDList = shell32.NewProc("SHGetPathFromIDListA")
var _ = builtin1("SHGetPathFromIDList(pidl)",
	func(a Value) Value {
		var buf [MAX_PATH]byte
		rtn, _, _ := shGetPathFromIDList.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&buf)))
		if rtn == 0 {
			return EmptyStr
		}
		return strRet(buf[:])
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
			0,
			0)
		buf := make([]byte, n)
		dragQueryFile.Call(
			intArg(a),
			intArg(b),
			uintptr(unsafe.Pointer(&buf[0])),
			n)
		return SuStr(str.BeforeFirst(string(buf), "\x00"))
	})

// dll bool Shell32:Shell_NotifyIcon(long dwMessage, NOTIFYICONDATA* lpdata)
var shell_NotifyIcon = shell32.NewProc("Shell_NotifyIconA")
var _ = builtin2("Shell_NotifyIcon(dwMessage, lpdata)",
	func(a, b Value) Value {
		nid := NOTIFYICONDATA{
			cbSize:               uint32(unsafe.Sizeof(NOTIFYICONDATA{})),
			hWnd:                 getHandle(b, "hWnd"),
			uID:                  getInt32(b, "uID"),
			uFlags:               getInt32(b, "uFlags"),
			uCallbackMessage:     getInt32(b, "uCallbackMessage"),
			hIcon:                getHandle(b, "hIcon"),
			dwState:              getInt32(b, "dwState"),
			dwStateMask:          getInt32(b, "dwStateMask"),
			uTimeoutVersionUnion: getUint32(b, "uTimeoutVersionUnion"),
			dwInfoFlags:          getInt32(b, "dwInfoFlags"),
		}
		copyStr(nid.szTip[:], ToStr(b.Get(nil, SuStr("szTip"))))
		copyStr(nid.szInfo[:], ToStr(b.Get(nil, SuStr("szInfo"))))
		copyStr(nid.szInfoTitle[:], ToStr(b.Get(nil, SuStr("szInfoTitle"))))
		rtn, _, _ := shell_NotifyIcon.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&nid)))
		return boolRet(rtn)
	})

type NOTIFYICONDATA struct {
	cbSize               uint32
	hWnd                 HANDLE
	uID                  int32
	uFlags               int32
	uCallbackMessage     int32
	hIcon                HANDLE
	szTip                [64]byte
	dwState              int32
	dwStateMask          int32
	szInfo               [256]byte
	uTimeoutVersionUnion uint32
	szInfoTitle          [64]byte
	dwInfoFlags          int32
}

// copyStr copies the string into the byte slice and adds a nul terminator
func copyStr(dst []byte, src string) {
	copy(dst[:], src)
	dst[len(src)] = 0
}

// dll bool Shell32:ShellExecuteEx(SHELLEXECUTEINFO* lpExecInfo)
var shellExecuteEx = shell32.NewProc("ShellExecuteExA")
var _ = builtin1("ShellExecuteEx(lpExecInfo)",
	func(a Value) Value {
		sei := SHELLEXECUTEINFO{
			cbSize:       int32(unsafe.Sizeof(SHELLEXECUTEINFO{})),
			fMask:        getInt32(a, "fMask"),
			hwnd:         getHandle(a, "hwnd"),
			lpVerb:       getStr(a, "lpVerb"),
			lpFile:       getStr(a, "lpFile"),
			lpDirectory:  getStr(a, "lpDirectory"),
			lpParameters: getStr(a, "lpParameters"),
			nShow:        getInt32(a, "nShow"),
			hInstApp:     getHandle(a, "hInstApp"),
			lpIDList:     getHandle(a, "lpIDList"),
			lpClass:      getStr(a, "lpClass"),
			hkeyClass:    getHandle(a, "hkeyClass"),
			dwHotKey:     getInt32(a, "dwHotKey"),
			hIcon:        getHandle(a, "hIcon"),
			hProcess:     getHandle(a, "hProcess"),
		}
		rtn, _, _ := shellExecuteEx.Call(
			uintptr(unsafe.Pointer(&sei)))
		return boolRet(rtn)
	})

type SHELLEXECUTEINFO struct {
	cbSize       int32
	fMask        int32
	hwnd         HANDLE
	lpVerb       *byte
	lpFile       *byte
	lpParameters *byte
	lpDirectory  *byte
	nShow        int32
	hInstApp     HANDLE
	lpIDList     HANDLE
	lpClass      *byte
	hkeyClass    HANDLE
	dwHotKey     int32
	hIcon        HANDLE
	hProcess     HANDLE
}

const MAX_PATH = 260

// dll pointer Shell32:SHBrowseForFolder(BROWSEINFO* lpbi)
var sHBrowseForFolder = shell32.NewProc("SHBrowseForFolderA")
var _ = builtin1("SHBrowseForFolder(lpbi)",
	func(a Value) Value {
		bi := BROWSEINFO{
			hwndOwner:      getHandle(a, "hwndOwner"),
			pidlRoot:       getHandle(a, "pidlRoot"),
			pszDisplayName: nil,
			lpszTitle:      getStr(a, "lpszTitle"),
			ulFlags:        getInt32(a, "ulFlags"),
			lpfn:           getCallback(a, "lpfn", 4),
			lParam:         getHandle(a, "lParam"),
			iImage:         getInt32(a, "iImage"),
		}
		rtn, _, _ := sHBrowseForFolder.Call(
			uintptr(unsafe.Pointer(&bi)))
		return intRet(rtn)
	})

type BROWSEINFO struct {
	hwndOwner      HANDLE
	pidlRoot       HANDLE
	pszDisplayName *byte
	lpszTitle      *byte
	ulFlags        int32
	lpfn           uintptr
	lParam         HANDLE
	iImage         int32
}
