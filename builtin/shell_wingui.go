// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"bytes"
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

var shell32 = MustLoadDLL("shell32.dll")

// dll void Shell32:DragAcceptFiles(pointer hWnd, bool fAccept)
var dragAcceptFiles = shell32.MustFindProc("DragAcceptFiles").Addr()
var _ = builtin(DragAcceptFiles, "(hWnd, fAccept)")

func DragAcceptFiles(a, b Value) Value {
	syscall.SyscallN(dragAcceptFiles,
		intArg(a),
		boolArg(b))
	return nil
}

// dll bool Shell32:SHGetPathFromIDList(pointer pidl, string path)
var shGetPathFromIDList = shell32.MustFindProc("SHGetPathFromIDListA").Addr()
var _ = builtin(SHGetPathFromIDList, "(pidl)")

func SHGetPathFromIDList(a Value) Value {
	var buf [MAX_PATH]byte
	rtn, _, _ := syscall.SyscallN(shGetPathFromIDList,
		intArg(a),
		uintptr(unsafe.Pointer(&buf)))
	if rtn == 0 {
		return EmptyStr
	}
	return bufZstr(buf[:])
}

// dll long Shell32:DragQueryFile(
// pointer hDrop, long iFile, string lpszFile, long cch)
var dragQueryFile = shell32.MustFindProc("DragQueryFile").Addr()

var _ = builtin(DragQueryFile, "(hDrop, iFile)")

func DragQueryFile(a, b Value) Value {
	n, _, _ := syscall.SyscallN(dragQueryFile,
		intArg(a),
		intArg(b),
		0,
		0)
	buf := make([]byte, n+1)
	syscall.SyscallN(dragQueryFile,
		intArg(a),
		intArg(b),
		uintptr(unsafe.Pointer(&buf[0])),
		n+1)
	return SuStr(hacks.BStoS(buf[:n]))
}

var _ = builtin(DragQueryFileCount, "(hDrop)")

func DragQueryFileCount(a Value) Value {
	rtn, _, _ := syscall.SyscallN(dragQueryFile,
		intArg(a),
		0xffffffff,
		0,
		0)
	return intRet(rtn)
}

// dll bool Shell32:Shell_NotifyIcon(long dwMessage, NOTIFYICONDATA* lpdata)
var shell_NotifyIcon = shell32.MustFindProc("Shell_NotifyIconA").Addr()
var _ = builtin(Shell_NotifyIcon, "(dwMessage, lpdata)")

func Shell_NotifyIcon(a, b Value) Value {
	nid := stNotifyIconData{
		cbSize:               uint32(nNotifyIconData),
		hWnd:                 getUintptr(b, "hWnd"),
		uID:                  getInt32(b, "uID"),
		uFlags:               getInt32(b, "uFlags"),
		uCallbackMessage:     getInt32(b, "uCallbackMessage"),
		hIcon:                getUintptr(b, "hIcon"),
		dwState:              getInt32(b, "dwState"),
		dwStateMask:          getInt32(b, "dwStateMask"),
		uTimeoutVersionUnion: getUint32(b, "uTimeoutVersionUnion"),
		dwInfoFlags:          getInt32(b, "dwInfoFlags"),
	}
	getZstrBs(b, "szTip", nid.szTip[:])
	getZstrBs(b, "szInfo", nid.szInfo[:])
	getZstrBs(b, "szInfoTitle", nid.szInfoTitle[:])
	rtn, _, _ := syscall.SyscallN(shell_NotifyIcon,
		intArg(a),
		uintptr(unsafe.Pointer(&nid)))
	return boolRet(rtn)
}

type stNotifyIconData struct {
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

const nNotifyIconData = unsafe.Sizeof(stNotifyIconData{})

// dll bool Shell32:ShellExecuteEx(SHELLEXECUTEINFO* lpExecInfo)
var shellExecuteEx = shell32.MustFindProc("ShellExecuteExA").Addr()
var _ = builtin(ShellExecuteEx, "(lpExecInfo)")

func ShellExecuteEx(a Value) Value {
	sei := stShellExecuteInfo{
		cbSize:       int32(nShellExecuteInfo),
		fMask:        getInt32(a, "fMask"),
		hwnd:         getUintptr(a, "hwnd"),
		lpVerb:       getZstr(a, "lpVerb"),
		lpFile:       getZstr(a, "lpFile"),
		lpDirectory:  getZstr(a, "lpDirectory"),
		lpParameters: getZstr(a, "lpParameters"),
		nShow:        getInt32(a, "nShow"),
		hInstApp:     getUintptr(a, "hInstApp"),
		lpIDList:     getUintptr(a, "lpIDList"),
		lpClass:      getZstr(a, "lpClass"),
		hkeyClass:    getUintptr(a, "hkeyClass"),
		dwHotKey:     getInt32(a, "dwHotKey"),
		hIcon:        getUintptr(a, "hIcon"),
		hProcess:     getUintptr(a, "hProcess"),
	}
	rtn, _, _ := syscall.SyscallN(shellExecuteEx,
		uintptr(unsafe.Pointer(&sei)))
	return boolRet(rtn)
}

type stShellExecuteInfo struct {
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

const nShellExecuteInfo = unsafe.Sizeof(stShellExecuteInfo{})

const MAX_PATH = 260

// dll pointer Shell32:SHBrowseForFolder(BROWSEINFO* lpbi)
var sHBrowseForFolder = shell32.MustFindProc("SHBrowseForFolderA").Addr()
var _ = builtin(SHBrowseForFolder, "(lpbi)")

func SHBrowseForFolder(th *Thread, args []Value) Value {
	b := args[0]
	bi := stBrowseInfo{
		hwndOwner:      getUintptr(b, "hwndOwner"),
		pidlRoot:       getUintptr(b, "pidlRoot"),
		pszDisplayName: nil,
		lpszTitle:      getZstr(b, "lpszTitle"),
		ulFlags:        getInt32(b, "ulFlags"),
		lpfn:           getCallback(th, b, "lpfn", 4),
		lParam:         getUintptr(b, "lParam"),
		iImage:         getInt32(b, "iImage"),
	}
	rtn, _, _ := syscall.SyscallN(sHBrowseForFolder,
		uintptr(unsafe.Pointer(&bi)))
	return intRet(rtn)
}

type stBrowseInfo struct {
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

const nBrowseInfo = unsafe.Sizeof(stBrowseInfo{})

var shGetFolderPath = shell32.MustFindProc("SHGetFolderPathA").Addr()

func ErrlogDir() string {
	const CSIDL_APPDATA = 0x001a
	const CSIDL_FLAG_CREATE = 0x8000
	var buf [MAX_PATH]byte
	rtn, _, _ := syscall.SyscallN(shGetFolderPath,
		0,
		CSIDL_APPDATA|CSIDL_FLAG_CREATE,
		0,
		0,
		uintptr(unsafe.Pointer(&buf)))
	i := bytes.IndexByte(buf[:], 0)
	if rtn != 0 || i == -1 {
		return "" // failed
	}
	return string(buf[:i]) + `\`
}
