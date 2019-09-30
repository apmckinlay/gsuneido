package builtin

import (
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/sys/windows"
)

var shell32 = windows.MustLoadDLL("shell32.dll")

// dll void Shell32:DragAcceptFiles(pointer hWnd, bool fAccept)
var dragAcceptFiles = shell32.MustFindProc("DragAcceptFiles").Addr()
var _ = builtin2("DragAcceptFiles(hWnd, fAccept)",
	func(a, b Value) Value {
		syscall.Syscall(dragAcceptFiles, 2,
			intArg(a),
			boolArg(b),
			0)
		return nil
	})

// dll bool Shell32:SHGetPathFromIDList(pointer pidl, string path)
var shGetPathFromIDList = shell32.MustFindProc("SHGetPathFromIDListA").Addr()
var _ = builtin1("SHGetPathFromIDList(pidl)",
	func(a Value) Value {
		var buf [MAX_PATH]byte
		rtn, _, _ := syscall.Syscall(shGetPathFromIDList, 2,
			intArg(a),
			uintptr(unsafe.Pointer(&buf)),
			0)
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
var dragQueryFile = shell32.MustFindProc("DragQueryFile").Addr()
var _ = builtin2("DragQueryFile(hDrop, iFile)",
	func(a, b Value) Value {
		n, _, _ := syscall.Syscall6(dragQueryFile, 4,
			intArg(a),
			intArg(b),
			0,
			0,
			0, 0)
		buf := make([]byte, n)
		syscall.Syscall6(dragQueryFile, 4,
			intArg(a),
			intArg(b),
			uintptr(unsafe.Pointer(&buf[0])),
			n,
			0, 0)
		return SuStr(str.BeforeFirst(string(buf), "\x00"))
	})

// dll bool Shell32:Shell_NotifyIcon(long dwMessage, NOTIFYICONDATA* lpdata)
var shell_NotifyIcon = shell32.MustFindProc("Shell_NotifyIconA").Addr()
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
		rtn, _, _ := syscall.Syscall(shell_NotifyIcon, 2,
			intArg(a),
			uintptr(unsafe.Pointer(&nid)),
			0)
		return boolRet(rtn)
	})

type NOTIFYICONDATA struct {
	cbSize               uint32
	hWnd                 HANDLE
	uID                  int32
	uFlags               int32
	uCallbackMessage     int32
	hIcon                HANDLE
	szTip                [128]byte
	dwState              int32
	dwStateMask          int32
	szInfo               [256]byte
	uTimeoutVersionUnion uint32
	szInfoTitle          [64]byte
	dwInfoFlags          int32
	guidItem             GUID
	hBalloonIcon         HANDLE
}

// copyStr copies the string into the byte slice and adds a nul terminator
func copyStr(dst []byte, src string) {
	copy(dst[:], src)
	dst[len(src)] = 0
}

// dll bool Shell32:ShellExecuteEx(SHELLEXECUTEINFO* lpExecInfo)
var shellExecuteEx = shell32.MustFindProc("ShellExecuteExA").Addr()
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
		rtn, _, _ := syscall.Syscall(shellExecuteEx, 1,
			uintptr(unsafe.Pointer(&sei)),
			0, 0)
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
var sHBrowseForFolder = shell32.MustFindProc("SHBrowseForFolderA").Addr()
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
		rtn, _, _ := syscall.Syscall(sHBrowseForFolder, 1,
			uintptr(unsafe.Pointer(&bi)),
			0, 0)
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
	_              [4]byte // padding
}
