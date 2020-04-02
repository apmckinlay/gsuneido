// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var shell32 = MustLoadDLL("shell32.dll")

// dll void Shell32:DragAcceptFiles(pointer hWnd, bool fAccept)
var dragAcceptFiles = shell32.MustFindProc("DragAcceptFiles").Addr()
var _ = builtin2("DragAcceptFiles(hWnd, fAccept)",
	func(a, b Value) Value {
		goc.Syscall2(dragAcceptFiles,
			intArg(a),
			boolArg(b))
		return nil
	})

// dll bool Shell32:SHGetPathFromIDList(pointer pidl, string path)
var shGetPathFromIDList = shell32.MustFindProc("SHGetPathFromIDListA").Addr()
var _ = builtin1("SHGetPathFromIDList(pidl)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		buf := heap.Alloc(MAX_PATH)
		rtn := goc.Syscall2(shGetPathFromIDList,
			intArg(a),
			uintptr(buf))
		if rtn == 0 {
			return EmptyStr
		}
		return SuStr(heap.GetStrZ(buf, MAX_PATH))
	})

// dll long Shell32:DragQueryFile(
//		pointer hDrop, long iFile, string lpszFile, long cch)
var dragQueryFile = shell32.MustFindProc("DragQueryFile").Addr()

var _ = builtin2("DragQueryFile(hDrop, iFile)",
	func(a, b Value) Value {
		n := goc.Syscall4(dragQueryFile,
			intArg(a),
			intArg(b),
			0,
			0)
		defer heap.FreeTo(heap.CurSize())
		buf := heap.Alloc(n + 1)
		goc.Syscall4(dragQueryFile,
			intArg(a),
			intArg(b),
			uintptr(buf),
			n)
		return SuStr(heap.GetStrN(buf, int(n)))
	})

var _ = builtin1("DragQueryFileCount(hDrop)",
	func(a Value) Value {
		rtn := goc.Syscall4(dragQueryFile,
			intArg(a),
			0xffffffff,
			0,
			0)
		return intRet(rtn)
	})

// dll bool Shell32:Shell_NotifyIcon(long dwMessage, NOTIFYICONDATA* lpdata)
var shell_NotifyIcon = shell32.MustFindProc("Shell_NotifyIconA").Addr()
var _ = builtin2("Shell_NotifyIcon(dwMessage, lpdata)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nNOTIFYICONDATA)
		nid := (*NOTIFYICONDATA)(p)
		*nid = NOTIFYICONDATA{
			cbSize:               uint32(nNOTIFYICONDATA),
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
		copyStr(nid.szTip[:], b, "szTip")
		copyStr(nid.szInfo[:], b, "szInfo")
		copyStr(nid.szInfoTitle[:], b, "szInfoTitle")
		rtn := goc.Syscall2(shell_NotifyIcon,
			intArg(a),
			uintptr(p))
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

const nNOTIFYICONDATA = unsafe.Sizeof(NOTIFYICONDATA{})

// dll bool Shell32:ShellExecuteEx(SHELLEXECUTEINFO* lpExecInfo)
var shellExecuteEx = shell32.MustFindProc("ShellExecuteExA").Addr()
var _ = builtin1("ShellExecuteEx(lpExecInfo)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nSHELLEXECUTEINFO)
		*(*SHELLEXECUTEINFO)(p) = SHELLEXECUTEINFO{
			cbSize:       int32(nSHELLEXECUTEINFO),
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
		rtn := goc.Syscall1(shellExecuteEx,
			uintptr(p))
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

const nSHELLEXECUTEINFO = unsafe.Sizeof(SHELLEXECUTEINFO{})

const MAX_PATH = 260

// dll pointer Shell32:SHBrowseForFolder(BROWSEINFO* lpbi)
var sHBrowseForFolder = shell32.MustFindProc("SHBrowseForFolderA").Addr()
var _ = builtin1("SHBrowseForFolder(lpbi)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nBROWSEINFO)
		*(*BROWSEINFO)(p) = BROWSEINFO{
			hwndOwner:      getHandle(a, "hwndOwner"),
			pidlRoot:       getHandle(a, "pidlRoot"),
			pszDisplayName: nil,
			lpszTitle:      getStr(a, "lpszTitle"),
			ulFlags:        getInt32(a, "ulFlags"),
			lpfn:           getCallback(a, "lpfn", 4),
			lParam:         getHandle(a, "lParam"),
			iImage:         getInt32(a, "iImage"),
		}
		rtn := goc.Syscall1(sHBrowseForFolder,
			uintptr(p))
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

const nBROWSEINFO = unsafe.Sizeof(BROWSEINFO{})
