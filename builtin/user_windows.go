// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"log"
	"time"
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/verify"
	"golang.org/x/sys/windows"
)

type RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

const nRECT = unsafe.Sizeof(RECT{})

type PAINTSTRUCT struct {
	hdc         HANDLE
	fErase      BOOL
	rcPaint     RECT
	fRestore    BOOL
	fIncUpdate  BOOL
	rgbReserved [32]byte
	_           [4]byte // padding
}

const nPAINTSTRUCT = unsafe.Sizeof(PAINTSTRUCT{})

type MONITORINFO struct {
	cbSize    uint32
	rcMonitor RECT
	rcWork    RECT
	dwFlags   uint32
}

const nMONITORINFO = unsafe.Sizeof(MONITORINFO{})

type SCROLLINFO struct {
	cbSize    uint32
	fMask     uint32
	nMin      int32
	nMax      int32
	nPage     uint32
	nPos      int32
	nTrackPos int32
}

const nSCROLLINFO = unsafe.Sizeof(SCROLLINFO{})

type POINT struct {
	x int32
	y int32
}

const nPOINT = unsafe.Sizeof(POINT{})

type WINDOWPLACEMENT struct {
	length           uint32
	flags            uint32
	showCmd          uint32
	ptMinPosition    POINT
	ptMaxPosition    POINT
	rcNormalPosition RECT
}

const nWINDOWPLACEMENT = unsafe.Sizeof(WINDOWPLACEMENT{})

type MENUITEMINFO struct {
	cbSize        uint32
	fMask         uint32
	fType         uint32
	fState        uint32
	wID           uint32
	hSubMenu      HANDLE
	hbmpChecked   HANDLE
	hbmpUnchecked HANDLE
	dwItemData    uintptr
	dwTypeData    *byte
	cch           uint32
	hbmpItem      HANDLE
}

const nMENUITEMINFO = unsafe.Sizeof(MENUITEMINFO{})

type WNDCLASS struct {
	style      uint32
	wndProc    uintptr
	clsExtra   int32
	wndExtra   int32
	instance   HANDLE
	icon       HANDLE
	cursor     HANDLE
	background HANDLE
	menuName   *byte
	className  *byte
}

const nWNDCLASS = unsafe.Sizeof(WNDCLASS{})

type CHARRANGE struct {
	cpMin int32
	cpMax int32
}

type TEXTRANGE struct {
	chrg      CHARRANGE
	lpstrText *byte
}

const nTEXTRANGE = unsafe.Sizeof(TEXTRANGE{})

type MSG struct {
	hwnd    HANDLE
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      POINT
	_       [4]byte // padding
}

const nMSG = unsafe.Sizeof(MSG{})

type TOOLINFO struct {
	cbSize     uint32
	uFlags     uint32
	hwnd       HANDLE
	uId        uint32
	rect       RECT
	hinst      HANDLE
	lpszText   *byte
	lParam     int32
	lpReserved uintptr
}

const nTOOLINFO = unsafe.Sizeof(TOOLINFO{})

type TOOLINFO2 struct {
	cbSize     uint32
	uFlags     uint32
	hwnd       HANDLE
	uId        uint32
	rect       RECT
	hinst      HANDLE
	lpszText   uintptr
	lParam     int32
	lpReserved uintptr
}

const nTOOLINFO2 = unsafe.Sizeof(TOOLINFO{})

type TVITEM struct {
	mask           uint32
	hItem          HANDLE
	state          uint32
	stateMask      uint32
	pszText        *byte
	cchTextMax     int32
	iImage         int32
	iSelectedImage int32
	cChildren      int32
	lParam         HANDLE
}

type TVITEMEX struct {
	TVITEM
	iIntegral      int32
	uStateEx       int32
	hwnd           HANDLE
	iExpandedImage int32
	iReserved      int32
}

const nTVITEMEX = unsafe.Sizeof(TVITEMEX{})

type TVINSERTSTRUCT struct {
	hParent      HANDLE
	hInsertAfter HANDLE
	item         TVITEMEX
}

const nTVINSERTSTRUCT = unsafe.Sizeof(TVINSERTSTRUCT{})

type TVSORTCB struct {
	hParent     HANDLE
	lpfnCompare uintptr
	lParam      HANDLE
}

const nTVSORTCB = unsafe.Sizeof(TVSORTCB{})

//-------------------------------------------------------------------

// dll User32:GetDesktopWindow() hwnd
var getDesktopWindow = user32.MustFindProc("GetDesktopWindow").Addr()
var _ = builtin0("GetDesktopWindow()",
	func() Value {
		rtn := goc.Syscall0(getDesktopWindow)
		return intRet(rtn)
	})

// dll User32:GetSysColor(long nIndex) long
var getSysColor = user32.MustFindProc("GetSysColor").Addr()
var _ = builtin1("GetSysColor(index)",
	func(a Value) Value {
		rtn := goc.Syscall1(getSysColor,
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetWindowRect(pointer hwnd, RECT* rect) bool
var getWindowRect = user32.MustFindProc("GetWindowRect").Addr()
var _ = builtin2("GetWindowRectApi(hwnd, rect)",
	func(a Value, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall2(getWindowRect,
			intArg(a),
			uintptr(r))
		urectToOb(r, b)
		return boolRet(rtn)
	})

// dll long User32:MessageBox(pointer window, [in] string text,
//		[in] string caption, long flags)
var messageBox = user32.MustFindProc("MessageBoxA").Addr()
var _ = builtin4("MessageBox(hwnd, text, caption, flags)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall4(messageBox,
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(stringArg(c)),
			intArg(d))
		return intRet(rtn)
	})

// dll User32:AdjustWindowRectEx(RECT* rect, long style, bool menu,
// 		long exStyle) bool
var adjustWindowRectEx = user32.MustFindProc("AdjustWindowRectEx").Addr()
var _ = builtin4("AdjustWindowRectEx(lpRect, dwStyle, bMenu, dwExStyle)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall4(adjustWindowRectEx,
			uintptr(rectArg(a, r)),
			intArg(b),
			boolArg(c),
			intArg(d))
		urectToOb(r, a)
		return boolRet(rtn)
	})

// dll User32:CreateMenu() pointer
var createMenu = user32.MustFindProc("CreateMenu").Addr()
var _ = builtin0("CreateMenu()",
	func() Value {
		rtn := goc.Syscall0(createMenu)
		return intRet(rtn)
	})

// dll User32:CreatePopupMenu() pointer
var createPopupMenu = user32.MustFindProc("CreatePopupMenu").Addr()
var _ = builtin0("CreatePopupMenu()",
	func() Value {
		rtn := goc.Syscall0(createPopupMenu)
		return intRet(rtn)
	})

// dll User32:AppendMenu(pointer hmenu, long flags, pointer item,
//		[in] string name) bool
var appendMenu = user32.MustFindProc("AppendMenuA").Addr()
var _ = builtin4("AppendMenu(hmenu, flags, item, name)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall4(appendMenu,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(stringArg(d)))
		return boolRet(rtn)
	})

// dll User32:DestroyMenu(pointer hmenu) bool
var destroyMenu = user32.MustFindProc("DestroyMenu").Addr()
var _ = builtin1("DestroyMenu(hmenu)",
	func(a Value) Value {
		rtn := goc.Syscall1(destroyMenu,
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:CreateWindowEx(long exStyle, resource classname, [in] string name,
//		long style, long x, long y, long w, long h, pointer parent, pointer menu,
//		pointer instance, pointer param) pointer
var createWindowEx = user32.MustFindProc("CreateWindowExA").Addr()
var _ = builtin("CreateWindowEx(exStyle, classname, name, style, x, y, w, h,"+
	" parent, menu, instance, param)",
	func(_ *Thread, a []Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall12(createWindowEx,
			intArg(a[0]),
			uintptr(stringArg(a[1])),
			uintptr(stringArg(a[2])),
			intArg(a[3]),
			intArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			intArg(a[8]),
			intArg(a[9]),
			intArg(a[10]),
			intArg(a[11]))
		return intRet(rtn)
	})

// dll User32:GetSystemMenu(pointer hWnd, bool bRevert) pointer
var getSystemMenu = user32.MustFindProc("GetSystemMenu").Addr()
var _ = builtin2("GetSystemMenu(hwnd, bRevert)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getSystemMenu,
			intArg(a),
			boolArg(b))
		return intRet(rtn)
	})

// dll User32:SetMenu(pointer hwnd, pointer hmenu) bool
var setMenu = user32.MustFindProc("SetMenu").Addr()
var _ = builtin2("SetMenu(hwnd, hmenu)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setMenu,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll User32:BeginPaint(pointer hwnd, PAINTSTRUCT* ps) pointer
var beginPaint = user32.MustFindProc("BeginPaint").Addr()
var _ = builtin2("BeginPaint(hwnd, ps)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nPAINTSTRUCT)
		rtn := goc.Syscall2(beginPaint,
			intArg(a),
			uintptr(psArg(b, p)))
		ps := (*PAINTSTRUCT)(p)
		b.Put(nil, SuStr("hdc"), IntVal(int(ps.hdc)))
		b.Put(nil, SuStr("fErase"), SuBool(ps.fErase != 0))
		b.Put(nil, SuStr("rcPaint"),
			rectToOb(&ps.rcPaint, b.Get(nil, SuStr("rcPaint"))))
		b.Put(nil, SuStr("fRestore"), SuBool(ps.fRestore != 0))
		b.Put(nil, SuStr("fIncUpdate"), SuBool(ps.fIncUpdate != 0))
		return intRet(rtn)
	})

// dll User32:EndPaint(pointer hwnd, PAINTSTRUCT* ps) bool
var endPaint = user32.MustFindProc("EndPaint").Addr()
var _ = builtin2("EndPaint(hwnd, ps)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nPAINTSTRUCT)
		rtn := goc.Syscall2(endPaint,
			intArg(a),
			uintptr(psArg(b, p)))
		return boolRet(rtn)
	})

func psArg(ob Value, p unsafe.Pointer) unsafe.Pointer {
	ps := (*PAINTSTRUCT)(p)
	ps.hdc = getHandle(ob, "hdc")
	ps.fErase = getBool(ob, "fErase")
	ps.rcPaint = getRect(ob, "rcPaint")
	ps.fRestore = getBool(ob, "fRestore")
	ps.fIncUpdate = getBool(ob, "fIncUpdate")
	return p
}

// dll User32:CallWindowProc(pointer wndprcPrev, pointer hwnd, long msg,
//		pointer wParam, pointer lParam) pointer
var callWindowProc = user32.MustFindProc("CallWindowProcA").Addr()
var _ = builtin5("CallWindowProc(wndprcPrev, hwnd, msg, wParam, lParam)",
	func(a, b, c, d, e Value) Value {
		rtn := goc.Syscall5(callWindowProc,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return intRet(rtn)
	})

// dll User32:CreateAcceleratorTable([in] string lpaccel, long cEntries) pointer
var createAcceleratorTable = user32.MustFindProc("CreateAcceleratorTableA").Addr()
var _ = builtin2("CreateAcceleratorTable(lpaccel, cEntries)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall2(createAcceleratorTable,
			uintptr(stringArg(a)),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:DestroyAcceleratorTable(pointer hAccel) bool
var destroyAcceleratorTable = user32.MustFindProc("DestroyAcceleratorTable").Addr()
var _ = builtin1("DestroyAcceleratorTable(hAccel)",
	func(a Value) Value {
		rtn := goc.Syscall1(destroyAcceleratorTable,
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:DestroyWindow(pointer hwnd) bool
var destroyWindow = user32.MustFindProc("DestroyWindow").Addr()
var _ = builtin1("DestroyWindow(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(destroyWindow,
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:DrawFrameControl(pointer hdc, RECT* lprc, long uType,
//		long uState) bool
var drawFrameControl = user32.MustFindProc("DrawFrameControl").Addr()
var _ = builtin4("DrawFrameControl(hdc, lprc, uType, uState)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall4(drawFrameControl,
			intArg(a),
			uintptr(rectArg(b, r)),
			intArg(c),
			intArg(d))
		return boolRet(rtn)
	})

// dll User32:DrawText(pointer hdc, [in] string lpsz, long cb, RECT* lprc,
//		long uFormat) long
var drawText = user32.MustFindProc("DrawTextA").Addr()
var _ = builtin5("DrawText(hdc, lpsz, cb, lprc, uFormat)",
	func(a, b, c, d, e Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall5(drawText,
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c),
			uintptr(rectArg(d, r)),
			intArg(e))
		urectToOb(r, d) // for CALCRECT
		return intRet(rtn)
	})

// dll User32:FillRect(pointer hdc, RECT* lpRect, pointer hBrush) long
var fillRect = user32.MustFindProc("FillRect").Addr()
var _ = builtin3("FillRect(hdc, lpRect, hBrush)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		verify.That(b != Zero)
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall3(fillRect,
			intArg(a),
			uintptr(rectArg(b, r)),
			intArg(c))
		return intRet(rtn)
	})

// dll User32:GetActiveWindow() pointer
var getActiveWindow = user32.MustFindProc("GetActiveWindow").Addr()
var _ = builtin0("GetActiveWindow()",
	func() Value {
		rtn := goc.Syscall0(getActiveWindow)
		return intRet(rtn)
	})

// dll User32:GetFocus() pointer
var getFocus = user32.MustFindProc("GetFocus").Addr()
var _ = builtin0("GetFocus()",
	func() Value {
		rtn := goc.Syscall0(getFocus)
		return intRet(rtn)
	})

// dll User32:GetClientRect(pointer hwnd, RECT* rect) bool
var getClientRect = user32.MustFindProc("GetClientRect").Addr()
var _ = builtin2("GetClientRect(hwnd, rect)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall2(getClientRect,
			intArg(a),
			uintptr(r))
		urectToOb(r, b)
		return boolRet(rtn)
	})

// dll User32:GetDC(pointer hwnd) pointer
var getDC = user32.MustFindProc("GetDC").Addr()
var _ = builtin1("GetDC(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(getDC,
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetMonitorInfo(pointer hMonitor, MONITORINFO* lpmi) bool
var getMonitorInfo = user32.MustFindProc("GetMonitorInfoA").Addr()
var _ = builtin2("GetMonitorInfoApi(hwnd, mInfo)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nMONITORINFO)
		mi := (*MONITORINFO)(p)
		mi.cbSize = uint32(nMONITORINFO)
		rtn := goc.Syscall2(getMonitorInfo,
			intArg(a),
			uintptr(p))
		b.Put(nil, SuStr("rcMonitor"), rectToOb(&mi.rcMonitor, nil))
		b.Put(nil, SuStr("rcWork"), rectToOb(&mi.rcWork, nil))
		b.Put(nil, SuStr("dwFlags"), IntVal(int(mi.dwFlags)))
		return boolRet(rtn)
	})

// dll User32:GetScrollInfo(pointer hwnd, long fnBar, SCROLLINFO* lpsi) bool
var getScrollInfo = user32.MustFindProc("GetScrollInfo").Addr()
var _ = builtin3("GetScrollInfo(hwnd, fnBar, lpsi)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nSCROLLINFO)
		si := (*SCROLLINFO)(p)
		*si = SCROLLINFO{
			cbSize:    uint32(nSCROLLINFO),
			fMask:     getUint32(c, "fMask"),
			nMin:      getInt32(c, "nMin"),
			nMax:      getInt32(c, "nMax"),
			nPage:     getUint32(c, "nPage"),
			nPos:      getInt32(c, "nPos"),
			nTrackPos: getInt32(c, "nTrackPos"),
		}
		rtn := goc.Syscall3(getScrollInfo,
			intArg(a),
			intArg(b),
			uintptr(p))
		c.Put(nil, SuStr("nMin"), IntVal(int(si.nMin)))
		c.Put(nil, SuStr("nMax"), IntVal(int(si.nMax)))
		c.Put(nil, SuStr("nPage"), IntVal(int(si.nPage)))
		c.Put(nil, SuStr("nPos"), IntVal(int(si.nPos)))
		c.Put(nil, SuStr("nTrackPos"), IntVal(int(si.nTrackPos)))
		return boolRet(rtn)
	})

// dll User32:GetScrollPos(pointer hwnd, int nBar) int
var getScrollPos = user32.MustFindProc("GetScrollPos").Addr()
var _ = builtin2("GetScrollPos(hwnd, nBar)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getScrollPos,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetSysColorBrush(long nIndex) pointer
var getSysColorBrush = user32.MustFindProc("GetSysColorBrush").Addr()
var _ = builtin1("GetSysColorBrush(nIndex)",
	func(a Value) Value {
		rtn := goc.Syscall1(getSysColorBrush,
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetSystemMetrics(long nIndex) long
var getSystemMetrics = user32.MustFindProc("GetSystemMetrics").Addr()
var _ = builtin1("GetSystemMetrics(nIndex)",
	func(a Value) Value {
		rtn := goc.Syscall1(getSystemMetrics,
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetWindowLong(pointer hwnd, long offset) long
var getWindowLong = user32.MustFindProc("GetWindowLongA").Addr()
var _ = builtin2("GetWindowLong(hwnd, offset)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getWindowLong,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetWindowLong(pointer hwnd, long offset) long
var getWindowLongPtr = user32.MustFindProc("GetWindowLongPtrA").Addr()
var _ = builtin2("GetWindowLongPtr(hwnd, offset)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getWindowLongPtr,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetWindowPlacement(pointer hwnd, WINDOWPLACEMENT* lpwndpl) bool
var getWindowPlacement = user32.MustFindProc("GetWindowPlacement").Addr()
var _ = builtin2("GetWindowPlacement(hwnd, ps)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nWINDOWPLACEMENT)
		wp := (*WINDOWPLACEMENT)(p)
		wp.length = uint32(nWINDOWPLACEMENT)
		rtn := goc.Syscall2(getWindowPlacement,
			intArg(a),
			uintptr(p))
		b.Put(nil, SuStr("flags"), IntVal(int(wp.flags)))
		b.Put(nil, SuStr("showCmd"), IntVal(int(wp.showCmd)))
		b.Put(nil, SuStr("ptMinPosition"), pointToOb(&wp.ptMinPosition, nil))
		b.Put(nil, SuStr("ptMaxPosition"), pointToOb(&wp.ptMaxPosition, nil))
		b.Put(nil, SuStr("rcNormalPosition"),
			rectToOb(&wp.rcNormalPosition, nil))
		return boolRet(rtn)
	})

// dll User32:GetWindowText(pointer hwnd, string buf, long len) long
var getWindowText = user32.MustFindProc("GetWindowTextA").Addr()
var getWindowTextLength = user32.MustFindProc("GetWindowTextLengthA").Addr()
var _ = builtin1("GetWindowText(hwnd)",
	func(hwnd Value) Value {
		defer heap.FreeTo(heap.CurSize())
		n := goc.Syscall1(getWindowTextLength,
			intArg(hwnd))
		buf := heap.Alloc(n + 1)
		n = goc.Syscall3(getWindowText,
			intArg(hwnd),
			uintptr(buf),
			n+1)
		return bufRet(buf, n)
	})

// dll User32:InflateRect(RECT* rect, long dx, long dy) bool
var inflateRect = user32.MustFindProc("InflateRect").Addr()
var _ = builtin3("InflateRect(rect, dx, dy)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall3(inflateRect,
			uintptr(rectArg(a, r)),
			intArg(b),
			intArg(c))
		urectToOb(r, a)
		return boolRet(rtn)
	})

// dll User32:InsertMenuItem(pointer hMenu, long uItem, bool fByPosition,
//		MENUITEMINFO* lpmii) bool
var insertMenuItem = user32.MustFindProc("InsertMenuItemA").Addr()
var _ = builtin4("InsertMenuItem(hMenu, uItem, fByPosition, lpmii)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nMENUITEMINFO)
		*(*MENUITEMINFO)(p) = MENUITEMINFO{
			cbSize:        uint32(nMENUITEMINFO),
			fMask:         getUint32(d, "fMask"),
			fType:         getUint32(d, "fType"),
			fState:        getUint32(d, "fState"),
			wID:           getUint32(d, "wID"),
			hSubMenu:      getHandle(d, "hSubMenu"),
			hbmpChecked:   getHandle(d, "hbmpChecked"),
			hbmpUnchecked: getHandle(d, "hbmpUnchecked"),
			dwItemData:    uintptr(getInt(d, "dwItemData")),
			dwTypeData:    getStr(d, "dwTypeData"),
			cch:           getUint32(d, "cch"),
			hbmpItem:      getHandle(d, "hbmpItem"),
		}
		rtn := goc.Syscall4(insertMenuItem,
			intArg(a),
			intArg(b),
			boolArg(c),
			uintptr(p))
		return boolRet(rtn)
	})

// dll long User32:GetMenuItemCount(pointer hMenu)
var getMenuItemCount = user32.MustFindProc("GetMenuItemCount").Addr()
var _ = builtin1("GetMenuItemCount(hMenu)",
	func(a Value) Value {
		rtn := goc.Syscall1(getMenuItemCount,
			intArg(a))
		return intRet(rtn)
	})

// dll long User32:GetMenuItemID(pointer hMenu, long nPos)
var getMenuItemID = user32.MustFindProc("GetMenuItemID").Addr()
var _ = builtin2("GetMenuItemID(hMenu, nPos)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getMenuItemID,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetMenuItemInfo(pointer hMenu, long uItem, bool fByPosition,
//		MENUITEMINFO* lpmii) bool
var getMenuItemInfo = user32.MustFindProc("GetMenuItemInfoA").Addr()
var _ = builtin2("GetMenuItemInfoText(hMenu, uItem)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		const MMIM_TYPE = 0x10
		const MFT_STRING = 0
		p := heap.Alloc(nMENUITEMINFO)
		mii := (*MENUITEMINFO)(p)
		*mii = MENUITEMINFO{
			cbSize:     uint32(nMENUITEMINFO),
			fMask:      MMIM_TYPE,
			fType:      MFT_STRING,
			dwTypeData: nil,
		}
		// get the length
		rtn := goc.Syscall4(getMenuItemInfo,
			intArg(a),
			intArg(b),
			0,
			uintptr(p))
		if rtn == 0 {
			return False
		}
		mii.cch++
		n := uintptr(mii.cch)
		buf := heap.Alloc(n)
		mii.dwTypeData = (*byte)(buf)
		rtn = goc.Syscall4(getMenuItemInfo,
			intArg(a),
			intArg(b),
			0,
			uintptr(p))
		return bufRet(buf, n-1) // -1 to omit nul terminator
	})

// dll User32:SetMenuItemInfo(pointer hMenu, long uItem, long fByPosition,
//		MENUITEMINFO* lpmii) bool
var setMenuItemInfo = user32.MustFindProc("SetMenuItemInfoA").Addr()
var _ = builtin4("SetMenuItemInfo(hMenu, uItem, fByPosition, lpmii)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nMENUITEMINFO)
		*(*MENUITEMINFO)(p) = MENUITEMINFO{
			cbSize:        uint32(nMENUITEMINFO),
			fMask:         getUint32(d, "fMask"),
			fType:         getUint32(d, "fType"),
			fState:        getUint32(d, "fState"),
			wID:           getUint32(d, "wID"),
			hSubMenu:      getHandle(d, "hSubMenu"),
			hbmpChecked:   getHandle(d, "hbmpChecked"),
			hbmpUnchecked: getHandle(d, "hbmpUnchecked"),
			dwItemData:    uintptr(getInt(d, "dwItemData")),
			//dwTypeData:    getStr(d, "dwTypeData"),
			cch:      getUint32(d, "cch"),
			hbmpItem: getHandle(d, "hbmpItem"),
		}
		rtn := goc.Syscall4(setMenuItemInfo,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		return boolRet(rtn)
	})

// dll User32:InvalidateRect(pointer hwnd, RECT* rect, bool erase) bool
var invalidateRect = user32.MustFindProc("InvalidateRect").Addr()
var _ = builtin3("InvalidateRect(hwnd, rect, erase)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall3(invalidateRect,
			intArg(a),
			uintptr(rectArg(b, r)),
			boolArg(c))
		return boolRet(rtn)
	})

// dll User32:IsWindowEnabled(pointer hwnd) bool
var isWindowEnabled = user32.MustFindProc("IsWindowEnabled").Addr()
var _ = builtin1("IsWindowEnabled(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(isWindowEnabled,
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:LoadCursor(pointer hinst, resource pszCursor) pointer
var loadCursor = user32.MustFindProc("LoadCursorA").Addr()
var _ = builtin2("LoadCursor(hinst, pszCursor)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(loadCursor,
			intArg(a),
			intArg(b)) // could be a string but we never use that
		return intRet(rtn)
	})

// dll User32:LoadIcon(pointer hInstance, resource lpIconName) pointer
var loadIcon = user32.MustFindProc("LoadIconA").Addr()
var _ = builtin2("LoadIcon(hinst, lpIconName)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(loadIcon,
			intArg(a),
			intArg(b)) // could be a string but we never use that
		return intRet(rtn)
	})

// dll User32:MonitorFromRect(RECT* lprc, long dwFlags) pointer
var monitorFromRect = user32.MustFindProc("MonitorFromRect").Addr()
var _ = builtin2("MonitorFromRect(lprc, dwFlags)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall2(monitorFromRect,
			uintptr(rectArg(a, r)),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:MoveWindow(pointer hwnd, long left, long top, long width,
//		long height, bool repaint) bool
var moveWindow = user32.MustFindProc("MoveWindow").Addr()
var _ = builtin6("MoveWindow(hwnd, left, top, width, height, repaint)",
	func(a, b, c, d, e, f Value) Value {
		rtn := goc.Syscall6(moveWindow,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			boolArg(f))
		return boolRet(rtn)
	})

// dll User32:RegisterClass(WNDCLASS* wc) short
var registerClass = user32.MustFindProc("RegisterClassA").Addr()
var _ = builtin1("RegisterClass(wc)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nWNDCLASS)
		*(*WNDCLASS)(p) = WNDCLASS{
			style:      getUint32(a, "style"),
			wndProc:    getHandle(a, "wndProc"),
			clsExtra:   getInt32(a, "clsExtra"),
			wndExtra:   getInt32(a, "wndExtra"),
			instance:   getHandle(a, "instance"),
			icon:       getHandle(a, "icon"),
			cursor:     getHandle(a, "cursor"),
			background: getHandle(a, "background"),
			menuName:   getStr(a, "menuName"),
			className:  getStr(a, "className"),
		}
		rtn := goc.Syscall1(registerClass,
			uintptr(p))
		return intRet(rtn)
	})

// dll User32:RegisterClipboardFormat([in] string lpszFormat) long
var registerClipboardFormat = user32.MustFindProc("RegisterClipboardFormatA").Addr()
var _ = builtin1("RegisterClipboardFormat(lpszFormat)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall1(registerClipboardFormat,
			uintptr(stringArg(a)))
		return intRet(rtn)
	})

// dll User32:ReleaseDC(pointer hWnd, pointer hDC) long
var releaseDC = user32.MustFindProc("ReleaseDC").Addr()
var _ = builtin2("ReleaseDC(hwnd, hDC)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(releaseDC,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:SetFocus(pointer hwnd) pointer
var setFocus = user32.MustFindProc("SetFocus").Addr()
var _ = builtin1("SetFocus(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(setFocus,
			intArg(a))
		return intRet(rtn)
	})

var postThreadMessage = user32.MustFindProc("PostThreadMessageA").Addr()

const WM_USER = 0x400

// dll User32:SetTimer(pointer hwnd, long id, long ms, TIMERPROC f) long
var setTimer = user32.MustFindProc("SetTimer").Addr()
var _ = builtin4("SetTimer(hwnd, id, ms, f)",
	func(a, b, c, d Value) Value {
		if options.TimersDisabled {
			return Zero
		}
		if windows.GetCurrentThreadId() != uiThreadId {
			// WARNING: don't use heap from background thread
			d.SetConcurrent() // since callback will be from different thread
			ts := timerSpec{hwnd: a, id: b, ms: c, cb: d, ret: make(chan Value, 1)}
			setTimerChan <- ts
			notifyMessageLoop()
			first := true
			for {
				select {
				case id := <-ts.ret:
					return id
				case <-time.After(5 * time.Second):
					if first {
						first = false
						log.Println("SetTimer timeout")
					}
				}
			}
		}
		return gocSetTimer(a, b, c, d)
	})

// gocSetTimer is called by SetTimer directly if on main UI thread
// and via updateUI2 if from background thread
func gocSetTimer(hwnd, id, ms, cb Value) Value {
	rtn := goc.Syscall4(setTimer,
		intArg(hwnd),
		intArg(id),
		intArg(ms),
		NewCallback(cb, 4))
	return intRet(rtn)
}

// dll User32:KillTimer(pointer hwnd, long id) bool
var killTimer = user32.MustFindProc("KillTimer").Addr()
var _ = builtin2("KillTimer(hwnd, id)",
	func(a, b Value) Value {
		if options.TimersDisabled {
			return False
		}
		rtn := goc.Syscall2(killTimer,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll User32:SetWindowLong(pointer hwnd, int offset, long value) long
var setWindowLong = user32.MustFindProc("SetWindowLongA").Addr()
var _ = builtin3("SetWindowLong(hwnd, offset, value)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(setWindowLong,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll User32:SetWindowLong(pointer hwnd, long offset, long value) long
var setWindowLongPtr = user32.MustFindProc("SetWindowLongPtrA").Addr()
var _ = builtin3("SetWindowLongPtr(hwnd, offset, value)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(setWindowLongPtr,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll User32:SetWindowProc(pointer hwnd, long offset, WNDPROC proc) pointer
var _ = builtin3("SetWindowProc(hwnd, offset, proc)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(setWindowLongPtr,
			intArg(a),
			intArg(b),
			NewCallback(c, 4))
		return intRet(rtn)
	})

// dll User32:SetWindowPlacement(pointer hwnd, WINDOWPLACEMENT* lpwndpl) bool
var setWindowPlacement = user32.MustFindProc("SetWindowPlacement").Addr()
var _ = builtin2("SetWindowPlacement(hwnd, lpwndpl)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nWINDOWPLACEMENT)
		*(*WINDOWPLACEMENT)(p) = WINDOWPLACEMENT{
			length:           getUint32(b, "length"),
			flags:            getUint32(b, "flags"),
			showCmd:          getUint32(b, "showCmd"),
			ptMinPosition:    getPoint(b, "ptMinPosition"),
			ptMaxPosition:    getPoint(b, "ptMaxPosition"),
			rcNormalPosition: getRect(b, "rcNormalPosition"),
		}
		rtn := goc.Syscall2(setWindowPlacement,
			intArg(a),
			uintptr(p))
		return boolRet(rtn)
	})

// dll User32:SetWindowPos(pointer hWnd, pointer hWndInsertAfter,
//		long X, long Y, long cx, long cy, long uFlags) bool
var setWindowPos = user32.MustFindProc("SetWindowPos").Addr()
var _ = builtin7("SetWindowPos(hWnd, hWndInsertAfter, X, Y, cx, cy, uFlags)",
	func(a, b, c, d, e, f, g Value) Value {
		rtn := goc.Syscall7(setWindowPos,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f),
			intArg(g))
		return boolRet(rtn)
	})

// dll User32:SetWindowText(pointer hwnd, [in] string text) bool
var setWindowText = user32.MustFindProc("SetWindowTextA").Addr()
var _ = builtin2("SetWindowText(hwnd, lpwndpl)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall2(setWindowText,
			intArg(a),
			uintptr(stringArg(b)))
		return boolRet(rtn)
	})

// dll User32:ShowWindow(pointer hwnd, long ncmd) bool
var showWindow = user32.MustFindProc("ShowWindow").Addr()
var _ = builtin2("ShowWindow(hwnd, ncmd)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(showWindow,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll User32:SystemParametersInfo(long uiAction, long uiParam, ? pvParam,
//		long fWinIni) bool
var systemParametersInfo = user32.MustFindProc("SystemParametersInfoA").Addr()

var _ = builtin0("SPI_GetFocusBorderHeight()",
	func() Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(4)
		const SPI_GETFOCUSBORDERHEIGHT = 0x2010
		goc.Syscall4(systemParametersInfo,
			SPI_GETFOCUSBORDERHEIGHT,
			0,
			uintptr(p),
			0)
		return IntVal(int((*(*int32)(p))))
	})

var _ = builtin0("SPI_GetWheelScrollLines()",
	func() Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(4)
		const SPI_GETWHEELSCROLLLINES = 104
		goc.Syscall4(systemParametersInfo,
			SPI_GETWHEELSCROLLLINES,
			0,
			uintptr(p),
			0)
		return IntVal(int((*(*int32)(p))))
	})

var _ = builtin0("SPI_GetWorkArea()",
	func() Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		const SPI_GETWORKAREA = 48
		goc.Syscall4(systemParametersInfo,
			SPI_GETWORKAREA,
			0,
			uintptr(r),
			0)
		return urectToOb(r, nil)
	})

// dll User32:PostQuitMessage(long exitcode) void
var postQuitMessage = user32.MustFindProc("PostQuitMessage").Addr()

func PostQuitMessage(arg uintptr) {
	goc.Syscall1(postQuitMessage, arg)
}

var _ = builtin1("PostQuitMessage(exitcode)",
	func(a Value) Value {
		PostQuitMessage(intArg(a))
		return nil
	})

// dll User32:GetNextDlgTabItem(pointer hDlg, pointer hCtl, bool prev) pointer
var getNextDlgTabItem = user32.MustFindProc("GetNextDlgTabItem").Addr()
var _ = builtin3("GetNextDlgTabItem(hDlg, hCtl, prev)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(getNextDlgTabItem,
			intArg(a),
			intArg(b),
			boolArg(c))
		return intRet(rtn)
	})

// dll User32:UpdateWindow(pointer hwnd) bool
var updateWindow = user32.MustFindProc("UpdateWindow").Addr()
var _ = builtin1("UpdateWindow(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(updateWindow,
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:DefWindowProc(pointer hwnd, long msg, pointer wParam,
//		pointer lParam) pointer
var defWindowProc = user32.MustFindProc("DefWindowProcA").Addr()
var _ = builtin4("DefWindowProc(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(defWindowProc,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return intRet(rtn)
	})

var _ = builtin0("GetDefWindowProc()",
	func() Value {
		return IntVal(int(defWindowProc))
	})

// dll User32:GetKeyState(long key) short
var getKeyState = user32.MustFindProc("GetKeyState").Addr()
var _ = builtin1("GetKeyState(nVirtKey)",
	func(a Value) Value {
		rtn := goc.Syscall1(getKeyState,
			intArg(a))
		return intRet(rtn)
	})

type TPMPARAMS struct {
	cbSize    int32
	rcExclude RECT
}

const nTPMPARAMS = unsafe.Sizeof(TPMPARAMS{})

// dll long User32:TrackPopupMenuEx(pointer hmenu, long fuFlags, long x, long y,
//		pointer hwnd, TPMPARAMS* lptpm)
var trackPopupMenuEx = user32.MustFindProc("TrackPopupMenuEx").Addr()
var _ = builtin6("TrackPopupMenuEx(hmenu, fuFlags, x, y, hwnd, lptpm)",
	func(a, b, c, d, e, f Value) Value {
		var p unsafe.Pointer
		if !f.Equal(Zero) {
			defer heap.FreeTo(heap.CurSize())
			p := heap.Alloc(nTPMPARAMS)
			*(*TPMPARAMS)(p) = TPMPARAMS{
				cbSize:    int32(nTPMPARAMS),
				rcExclude: getRect(f, "rcExclude"),
			}
		}
		rtn := goc.Syscall6(trackPopupMenuEx,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			uintptr(p))
		return intRet(rtn)
	})

// dll bool User32:OpenClipboard(pointer hwnd)
var openClipboard = user32.MustFindProc("OpenClipboard").Addr()
var _ = builtin1("OpenClipboard(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(openClipboard,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:EmptyClipboard()
var emptyClipboard = user32.MustFindProc("EmptyClipboard").Addr()
var _ = builtin0("EmptyClipboard()",
	func() Value {
		rtn := goc.Syscall0(emptyClipboard)
		return boolRet(rtn)
	})

// dll pointer User32:GetClipboardData(long format)
var getClipboardData = user32.MustFindProc("GetClipboardData").Addr()
var _ = builtin1("GetClipboardData(format)",
	func(a Value) Value {
		rtn := goc.Syscall1(getClipboardData,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:SetClipboardData(long uFormat, pointer hMem)
var setClipboardData = user32.MustFindProc("SetClipboardData").Addr()
var _ = builtin2("SetClipboardData(uFormat, hMem)",
	func(a Value, b Value) Value {
		rtn := goc.Syscall2(setClipboardData,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll bool User32:CloseClipboard()
var closeClipboard = user32.MustFindProc("CloseClipboard").Addr()
var _ = builtin0("CloseClipboard()",
	func() Value {
		rtn := goc.Syscall0(closeClipboard)
		return boolRet(rtn)
	})

// dll pointer User32:BeginDeferWindowPos(
// 	long nNumWindows		// initial number of windows to allocate space for
// 	)
var beginDeferWindowPos = user32.MustFindProc("BeginDeferWindowPos").Addr()
var _ = builtin1("BeginDeferWindowPos(nNumWindows)",
	func(a Value) Value {
		rtn := goc.Syscall1(beginDeferWindowPos,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:CallNextHookEx(		// returns an LRESULT
// 	pointer	hhk,	// handle to current hook [HHOOK]
// 	long	nCode,	// hook code passed to hook procedure [int]
// 	pointer	wParam,	// value passed to hook procedure [WPARAM]
// 	pointer	lParam	// value passed to hook procedure [LPARAM]
// )
var callNextHookEx = user32.MustFindProc("CallNextHookEx").Addr()
var _ = builtin4("CallNextHookEx(hhk, nCode, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(callNextHookEx,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return intRet(rtn)
	})

// dll long User32:CheckMenuItem(
// 	pointer hmenu, 		// handle to menu
// 	long uIDCheckItem, 	// menu item to check or uncheck
// 	long uCheck 		// menu item options
// 	)
var checkMenuItem = user32.MustFindProc("CheckMenuItem").Addr()
var _ = builtin3("CheckMenuItem(hmenu, uIDCheckItem, uCheck)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(checkMenuItem,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll bool user32:DeleteMenu(pointer hMenu, long uPosition, long uFlags)
var deleteMenu = user32.MustFindProc("DeleteMenu").Addr()
var _ = builtin3("DeleteMenu(hMenu, uPosition, uFlags)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(deleteMenu,
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll long User32:EnableMenuItem(
// 	pointer hMenu,
// 	long uIDEnableItem,
// 	long uEnable
// 	)
var enableMenuItem = user32.MustFindProc("EnableMenuItem").Addr()
var _ = builtin3("EnableMenuItem(hMenu, uIDEnableItem, uEnable)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(enableMenuItem,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll bool User32:EnableWindow(pointer hWnd, bool bEnable)
var enableWindow = user32.MustFindProc("EnableWindow").Addr()
var _ = builtin2("EnableWindow(hWnd, bEnable)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(enableWindow,
			intArg(a),
			boolArg(b))
		return boolRet(rtn)
	})

// dll bool User32:EndDeferWindowPos(
// 	pointer hWinPosInfo		// handle to internal structure
// 	)
var endDeferWindowPos = user32.MustFindProc("EndDeferWindowPos").Addr()
var _ = builtin1("EndDeferWindowPos(hWinPosInfo)",
	func(a Value) Value {
		rtn := goc.Syscall1(endDeferWindowPos,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:EndDialog(pointer hwndDlg, long nResult)
var endDialog = user32.MustFindProc("EndDialog").Addr()
var _ = builtin2("EndDialog(hwndDlg, nResult)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(endDialog,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll long User32:EnumClipboardFormats(long format)
var enumClipboardFormats = user32.MustFindProc("EnumClipboardFormats").Addr()
var _ = builtin1("EnumClipboardFormats(format)",
	func(a Value) Value {
		rtn := goc.Syscall1(enumClipboardFormats,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:FindWindow([in] string c, [in] string n)
var findWindow = user32.MustFindProc("FindWindowA").Addr()
var _ = builtin2("FindWindow(c, n)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall2(findWindow,
			uintptr(stringArg(a)),
			uintptr(stringArg(b)))
		return intRet(rtn)
	})

// dll pointer User32:GetAncestor(pointer hwnd, long gaFlags)
var getAncestor = user32.MustFindProc("GetAncestor").Addr()
var _ = builtin2("GetAncestor(hwnd, gaFlags)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getAncestor,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll long User32:GetClipboardFormatName(
// 	long		format,				// clipboard format to retrieve
// 	string		lpszFormatName,		// buffer to receive format name
// 	long		cchMaxCount			// maximum length of string to copy into buffer
// 	)
var getClipboardFormatName = user32.MustFindProc("GetClipboardFormatNameA").Addr()
var _ = builtin3("GetClipboardFormatName(format, lpszFormatName, cchMaxCount)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall3(getClipboardFormatName,
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c))
		return intRet(rtn)
	})

// dll pointer User32:GetCursor()
var getCursor = user32.MustFindProc("GetCursor").Addr()
var _ = builtin0("GetCursor()",
	func() Value {
		rtn := goc.Syscall0(getCursor)
		return intRet(rtn)
	})

// dll long User32:GetDoubleClickTime()
var getDoubleClickTime = user32.MustFindProc("GetDoubleClickTime").Addr()
var _ = builtin0("GetDoubleClickTime()",
	func() Value {
		rtn := goc.Syscall0(getDoubleClickTime)
		return intRet(rtn)
	})

// dll long User32:GetMenuState(
// 	pointer hMenu, 	// handle to menu
// 	long uId, 		// menu item to query
// 	long uFlags		// options
// 	)
var getMenuState = user32.MustFindProc("GetMenuState").Addr()
var _ = builtin3("GetMenuState(hMenu, uId, uFlags)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(getMenuState,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll long User32:GetMessagePos()
var getMessagePos = user32.MustFindProc("GetMessagePos").Addr()
var _ = builtin0("GetMessagePos()",
	func() Value {
		rtn := goc.Syscall0(getMessagePos)
		return intRet(rtn)
	})

// dll pointer User32:GetNextDlgGroupItem(pointer hDlg, pointer hCtl, bool prev)
var getNextDlgGroupItem = user32.MustFindProc("GetNextDlgGroupItem").Addr()
var _ = builtin3("GetNextDlgGroupItem(hDlg, hCtl, prev)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(getNextDlgGroupItem,
			intArg(a),
			intArg(b),
			boolArg(c))
		return intRet(rtn)
	})

// dll pointer User32:GetParent(pointer hwnd)
var getParent = user32.MustFindProc("GetParent").Addr()
var _ = builtin1("GetParent(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(getParent,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:GetSubMenu(
// 	pointer hmenu,	//menu handle
// 	long position	//position
// 	)
var getSubMenu = user32.MustFindProc("GetSubMenu").Addr()
var _ = builtin2("GetSubMenu(hmenu, position)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getSubMenu,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer User32:GetTopWindow(pointer hWnd)
var getTopWindow = user32.MustFindProc("GetTopWindow").Addr()
var _ = builtin1("GetTopWindow(hWnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(getTopWindow,
			intArg(a))
		return intRet(rtn)
	})

// dll long User32:GetUpdateRgn(pointer hwnd, pointer hRgn, bool bErase)
var getUpdateRgn = user32.MustFindProc("GetUpdateRgn").Addr()
var _ = builtin3("GetUpdateRgn(hwnd, hRgn, bErase)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(getUpdateRgn,
			intArg(a),
			intArg(b),
			boolArg(c))
		return intRet(rtn)
	})

// dll pointer User32:GetWindow(pointer hWnd, long uCmd)
var getWindow = user32.MustFindProc("GetWindow").Addr()
var _ = builtin2("GetWindow(hWnd, uCmd)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getWindow,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer User32:GetWindowDC(pointer hwnd)
var getWindowDC = user32.MustFindProc("GetWindowDC").Addr()
var _ = builtin1("GetWindowDC(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(getWindowDC,
			intArg(a))
		return intRet(rtn)
	})

// dll bool User32:IsClipboardFormatAvailable(long format)
var isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable").Addr()
var _ = builtin1("IsClipboardFormatAvailable(format)",
	func(a Value) Value {
		rtn := goc.Syscall1(isClipboardFormatAvailable,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:IsWindow(pointer hwnd)
var isWindow = user32.MustFindProc("IsWindow").Addr()
var _ = builtin1("IsWindow(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(isWindow,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:IsChild(pointer hwndParent, pointer hwnd)
var isChild = user32.MustFindProc("IsChild").Addr()
var _ = builtin2("IsChild(hwndParent, hwnd)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(isChild,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll bool User32:IsWindowVisible(pointer hwnd)
var isWindowVisible = user32.MustFindProc("IsWindowVisible").Addr()
var _ = builtin1("IsWindowVisible(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(isWindowVisible,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:MessageBeep(long type)
var messageBeep = user32.MustFindProc("MessageBeep").Addr()
var _ = builtin1("MessageBeep(type)",
	func(a Value) Value {
		rtn := goc.Syscall1(messageBeep,
			intArg(a))
		return boolRet(rtn)
	})

// dll void User32:mouse_event(
// 	long	dwFlags,		// motion and click options
// 	long	dx,			// horizontal position or change
// 	long	dy,			// vertical position or change
// 	long	dwData,		// wheel movement
// 	pointer	dwExtraInfo	// (ULONG_PTR) application-defined information
// )
var mouse_event = user32.MustFindProc("mouse_event").Addr()
var _ = builtin5("Mouse_event(dwFlags, dx, dy, dwData, dwExtraInfo)",
	func(a, b, c, d, e Value) Value {
		goc.Syscall5(mouse_event,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return nil
	})

// dll bool User32:PostMessage(pointer hwnd, long msg, pointer wParam, pointer lParam)
var postMessage = user32.MustFindProc("PostMessageA").Addr()
var _ = builtin4("PostMessage(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(postMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return boolRet(rtn)
	})

// dll bool User32:RegisterHotKey(
// 	pointer hWnd /*optional*/,
// 	long id,
// 	long fsModifiers,
// 	long vk
// )
var registerHotKey = user32.MustFindProc("RegisterHotKey").Addr()
var _ = builtin4("RegisterHotKey(hWnd, id, fsModifiers, vk)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(registerHotKey,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return boolRet(rtn)
	})

// dll bool User32:ReleaseCapture()
var releaseCapture = user32.MustFindProc("ReleaseCapture").Addr()
var _ = builtin0("ReleaseCapture()",
	func() Value {
		rtn := goc.Syscall0(releaseCapture)
		return boolRet(rtn)
	})

// dll pointer User32:SetActiveWindow(pointer hWnd)
var setActiveWindow = user32.MustFindProc("SetActiveWindow").Addr()
var _ = builtin1("SetActiveWindow(hWnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(setActiveWindow,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:SetCapture(pointer hwnd)
var setCapture = user32.MustFindProc("SetCapture").Addr()
var _ = builtin1("SetCapture(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(setCapture,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:SetCursor(pointer hcursor)
var setCursor = user32.MustFindProc("SetCursor").Addr()
var _ = builtin1("SetCursor(hcursor)",
	func(a Value) Value {
		rtn := goc.Syscall1(setCursor,
			intArg(a))
		return intRet(rtn)
	})

// dll bool User32:SetForegroundWindow(pointer hwnd)
// // puts creator thread into foreground and activates window
var setForegroundWindow = user32.MustFindProc("SetForegroundWindow").Addr()
var _ = builtin1("SetForegroundWindow(hwnd)",
	func(a Value) Value {
		rtn := goc.Syscall1(setForegroundWindow,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:SetMenuDefaultItem(
// 	pointer hMenu,
// 	long uItem,
// 	long fByPosition
// 	)
var setMenuDefaultItem = user32.MustFindProc("SetMenuDefaultItem").Addr()
var _ = builtin3("SetMenuDefaultItem(hMenu, uItem, fByPosition)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(setMenuDefaultItem,
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll pointer User32:SetParent(pointer hwndNewChild, pointer hwndNewParent)
var setParent = user32.MustFindProc("SetParent").Addr()
var _ = builtin2("SetParent(hwndNewChild, hwndNewParent)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setParent,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll bool User32:SetProp(pointer hwnd, [in] string name, pointer value)
var setProp = user32.MustFindProc("SetPropA").Addr()
var _ = builtin3("SetProp(hwnd, name, value)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall3(setProp,
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool User32:UnhookWindowsHookEx(
// 	pointer hhk	// handle to hook procedure [HHOOK]
// )
var unhookWindowsHookEx = user32.MustFindProc("UnhookWindowsHookEx").Addr()
var _ = builtin1("UnhookWindowsHookEx(hhk)",
	func(a Value) Value {
		rtn := goc.Syscall1(unhookWindowsHookEx,
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:UnregisterHotKey(pointer hWnd /*optional*/, long id)
var unregisterHotKey = user32.MustFindProc("UnregisterHotKey").Addr()
var _ = builtin2("UnregisterHotKey(hWnd, id)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(unregisterHotKey,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll bool User32:ClientToScreen(pointer hwnd, POINT* point)
var clientToScreen = user32.MustFindProc("ClientToScreen").Addr()
var _ = builtin2("ClientToScreen(hWnd, point)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nPOINT)
		rtn := goc.Syscall2(clientToScreen,
			intArg(a),
			uintptr(pointArg(b, p)))
		pt := (*POINT)(p)
		b.Put(nil, SuStr("x"), IntVal(int(pt.x)))
		b.Put(nil, SuStr("y"), IntVal(int(pt.y)))
		return boolRet(rtn)
	})

// dll bool User32:ClipCursor(RECT* rect)
var clipCursor = user32.MustFindProc("ClipCursor").Addr()
var _ = builtin1("ClipCursor(rect)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall1(clipCursor,
			uintptr(rectArg(a, r)))
		return boolRet(rtn)
	})

// dll pointer User32:DeferWindowPos(pointer hWinPosInfo, pointer hWnd,
//		pointer hWndInsertAfter, long x, long y, long cx, long cy, long flags)
var deferWindowPos = user32.MustFindProc("DeferWindowPos").Addr()
var _ = builtin("DeferWindowPos(hWinPosInfo, hWnd, hWndInsertAfter, "+
	"x, y, cx, cy, flags)",
	func(_ *Thread, a []Value) Value {
		rtn := goc.Syscall8(deferWindowPos,
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			intArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]))
		return intRet(rtn)
	})

// dll bool User32:DrawFocusRect(pointer hdc, RECT* lprc)
var drawFocusRect = user32.MustFindProc("DrawFocusRect").Addr()
var _ = builtin2("DrawFocusRect(hwnd, rect)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall2(drawFocusRect,
			intArg(a),
			uintptr(rectArg(b, r)))
		return boolRet(rtn)
	})

// dll long User32:DrawTextEx(pointer hdc, [in] string lpsz, long cb,
// RECT* lprc, long uFormat, DRAWTEXTPARAMS* params)
var drawTextEx = user32.MustFindProc("DrawTextExA").Addr()

var _ = builtin6("DrawTextEx(hdc, lpsz, cb, lprc, uFormat, params)",
	func(a, b, c, d, e, f Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall6(drawTextEx,
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c),
			uintptr(rectArg(d, r)),
			intArg(e),
			uintptr(drawTextParams(f)))
		urectToOb(r, d)
		return intRet(rtn)
	})

var _ = builtin5("DrawTextExOut(hdc, text, rect, flags, params)",
	func(a, b, c, d, e Value) Value {
		defer heap.FreeTo(heap.CurSize())
		text := ToStr(b)
		bufsize := len(text) + 8
		buf := strToBuf(text, bufsize)
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall6(drawTextEx,
			intArg(a),
			uintptr(buf),
			uintptrMinusOne,
			uintptr(rectArg(c, r)),
			intArg(d),
			uintptr(drawTextParams(e)))
		urectToOb(r, c)
		ob := NewSuObject()
		ob.Put(nil, SuStr("text"), bufToStr(buf, uintptr(bufsize)))
		ob.Put(nil, SuStr("result"), intRet(rtn))
		return ob
	})

func drawTextParams(x Value) unsafe.Pointer {
	p := unsafe.Pointer(nil)
	if !x.Equal(Zero) {
		p = heap.Alloc(nDRAWTEXTPARAMS)
		*(*DRAWTEXTPARAMS)(p) = DRAWTEXTPARAMS{
			cbSize:        uint32(nDRAWTEXTPARAMS),
			iTabLength:    getInt32(x, "iTabLength"),
			iLeftMargin:   getInt32(x, "iLeftMargin"),
			iRightMargin:  getInt32(x, "iRightMargin"),
			uiLengthDrawn: getInt32(x, "uiLengthDrawn"),
		}
	}
	return p
}

type DRAWTEXTPARAMS struct {
	cbSize        uint32
	iTabLength    int32
	iLeftMargin   int32
	iRightMargin  int32
	uiLengthDrawn int32
}

const nDRAWTEXTPARAMS = unsafe.Sizeof(DRAWTEXTPARAMS{})

// dll bool User32:TrackMouseEvent(TRACKMOUSEEVENT* lpEventTrack)
var trackMouseEvent = user32.MustFindProc("TrackMouseEvent").Addr()
var _ = builtin1("TrackMouseEvent(lpEventTrack)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nTRACKMOUSEEVENT)
		*(*TRACKMOUSEEVENT)(p) = TRACKMOUSEEVENT{
			cbSize:      uint32(nTRACKMOUSEEVENT),
			dwFlags:     getInt32(a, "dwFlags"),
			hwndTrack:   getHandle(a, "hwndTrack"),
			dwHoverTime: getInt32(a, "dwHoverTime"),
		}
		rtn := goc.Syscall1(trackMouseEvent,
			uintptr(p))
		return boolRet(rtn)
	})

type TRACKMOUSEEVENT struct {
	cbSize      uint32
	dwFlags     int32
	hwndTrack   uintptr
	dwHoverTime int32
	_           [4]byte // padding
}

const nTRACKMOUSEEVENT = unsafe.Sizeof(TRACKMOUSEEVENT{})

// dll bool User32:FlashWindowEx(FLASHWINFO* fi)
var flashWindowEx = user32.MustFindProc("FlashWindowEx").Addr()
var _ = builtin1("FlashWindowEx(fi)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nFLASHWINFO)
		*(*FLASHWINFO)(p) = FLASHWINFO{
			cbSize:    uint32(nFLASHWINFO),
			hwnd:      getHandle(a, "hwnd"),
			dwFlags:   getInt32(a, "dwFlags"),
			uCount:    getInt32(a, "uCount"),
			dwTimeout: getInt32(a, "dwTimeout"),
		}
		rtn := goc.Syscall1(flashWindowEx,
			uintptr(p))
		return boolRet(rtn)
	})

type FLASHWINFO struct {
	cbSize    uint32
	hwnd      HANDLE
	dwFlags   int32
	uCount    int32
	dwTimeout int32
	_         [4]byte // padding
}

const nFLASHWINFO = unsafe.Sizeof(FLASHWINFO{})

// dll long User32:FrameRect(pointer hdc, RECT* rect, pointer brush)
var frameRect = user32.MustFindProc("FrameRect").Addr()
var _ = builtin3("FrameRect(hdc, rect, brush)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall3(frameRect,
			intArg(a),
			uintptr(rectArg(b, r)),
			intArg(c))
		return intRet(rtn)
	})

// dll bool User32:GetClipCursor(RECT* rect)
var getClipCursor = user32.MustFindProc("GetClipCursor").Addr()
var _ = builtin1("GetClipCursor(rect)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall1(getClipCursor,
			uintptr(r))
		urectToOb(r, a)
		return boolRet(rtn)
	})

// dll bool User32:GetCursorPos(POINT* p)
var getCursorPos = user32.MustFindProc("GetCursorPos").Addr()
var _ = builtin1("GetCursorPos(rect)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nPOINT)
		rtn := goc.Syscall1(getCursorPos,
			uintptr(p))
		upointToOb(p, a)
		return boolRet(rtn)
	})

// dll bool User32:EnumThreadWindows(long dwThreadId, WNDENUMPROC lpfn,
//		pointer lParam)
var enumThreadWindows = user32.MustFindProc("EnumThreadWindows").Addr()
var _ = builtin3("EnumThreadWindows(dwThreadId, lpfn, lParam)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(enumThreadWindows,
			intArg(a),
			NewCallback(b, 2),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool User32:EnumChildWindows(pointer hwnd, WNDENUMPROC lpEnumProc,
//		pointer lParam)
var enumChildWindows = user32.MustFindProc("EnumChildWindows").Addr()
var _ = builtin3("EnumChildWindowsApi(hwnd, lpEnumProc, lParam)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(enumChildWindows,
			intArg(a),
			NewCallback(b, 2),
			intArg(c))
		return boolRet(rtn)
	})

// dll pointer User32:WindowFromPoint(POINT pt)
var windowFromPoint = user32.MustFindProc("WindowFromPoint").Addr()
var _ = builtin1("WindowFromPoint(pt)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nPOINT)
		rtn := goc.Syscall1(windowFromPoint,
			uintptr(pointArg(a, p)))
		return intRet(rtn)
	})

// dll long User32:GetWindowThreadProcessId(pointer hwnd, LONG* lpdwProcessId)
var getWindowThreadProcessId = user32.MustFindProc("GetWindowThreadProcessId").Addr()
var _ = builtin2("GetWindowThreadProcessId(hwnd, lpdwProcessId)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(4)
		rtn := goc.Syscall2(getWindowThreadProcessId,
			intArg(a),
			uintptr(p))
		b.Put(nil, SuStr("x"), IntVal(int(*(*int32)(p))))
		return boolRet(rtn)
	})

// dll long User32:TrackPopupMenu(pointer hMenu, long uFlags, long x, long y,
//		long nReserved, pointer hWnd, RECT* prcRect)
var trackPopupMenu = user32.MustFindProc("TrackPopupMenu").Addr()
var _ = builtin7("TrackPopupMenu(hMenu, uFlags, x, y, nReserved, hWnd, prcRect)",
	func(a, b, c, d, e, f, g Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall7(trackPopupMenu,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f),
			uintptr(rectArg(g, r)))
		return intRet(rtn)
	})

// dll pointer User32:SetWindowsHookEx(long idHook, HOOKPROC lpfn, pointer hMod,
//		long dwThreadId)
var setWindowsHookEx = user32.MustFindProc("SetWindowsHookExA").Addr()
var _ = builtin4("SetWindowsHookEx(idHook, lpfn, hMod, dwThreadId)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(setWindowsHookEx,
			intArg(a),
			NewCallback(b, 3),
			intArg(c),
			intArg(d))
		return intRet(rtn)
	})

// dll long User32:SetScrollInfo(pointer hwnd, long bar, SCROLLINFO* si,
//		bool redraw)
var setScrollInfo = user32.MustFindProc("SetScrollInfo").Addr()
var _ = builtin4("SetScrollInfo(hwnd, bar, si, redraw)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nSCROLLINFO)
		*(*SCROLLINFO)(p) = SCROLLINFO{
			cbSize:    uint32(nSCROLLINFO),
			fMask:     getUint32(c, "fMask"),
			nMin:      getInt32(c, "nMin"),
			nMax:      getInt32(c, "nMax"),
			nPage:     getUint32(c, "nPage"),
			nPos:      getInt32(c, "nPos"),
			nTrackPos: getInt32(c, "nTrackPos"),
		}
		rtn := goc.Syscall4(setScrollInfo,
			intArg(a),
			intArg(b),
			uintptr(p),
			boolArg(d))
		return intRet(rtn)
	})

// dll long User32:ScrollWindowEx(pointer hwnd, long dx, long dy, RECT* scroll,
//		RECT* clip, pointer rgnUpdate, RECT* rcUpdate, long flags)
var scrollWindowEx = user32.MustFindProc("ScrollWindowEx").Addr()
var _ = builtin("ScrollWindowEx(hwnd, dx, dy, scroll, clip, rgnUpdate,"+
	"rcUpdate, flags)",
	func(_ *Thread, a []Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r1 := heap.Alloc(nRECT)
		r2 := heap.Alloc(nRECT)
		r3 := heap.Alloc(nRECT)
		rtn := goc.Syscall8(scrollWindowEx,
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			uintptr(rectArg(a[3], r1)),
			uintptr(rectArg(a[4], r2)),
			intArg(a[5]),
			uintptr(rectArg(a[6], r3)),
			intArg(a[7]))
		urectToOb(r3, a[6])
		return intRet(rtn)
	})

// dll bool User32:ScreenToClient(pointer hwnd, POINT* p)
var screenToClient = user32.MustFindProc("ScreenToClient").Addr()
var _ = builtin2("ScreenToClient(hWnd, p)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nPOINT)
		rtn := goc.Syscall2(screenToClient,
			intArg(a),
			uintptr(pointArg(b, p)))
		upointToOb(p, b)
		return boolRet(rtn)
	})

// dll pointer User32:LoadImage(pointer hInstance, resource lpszName,
//		long uType, long cxDesired, long cyDesired, long fuLoad)
var loadImage = user32.MustFindProc("LoadImageA").Addr()
var _ = builtin6("LoadImage(hInstance, lpszName, uType, cxDesired, cyDesired,"+
	" fuLoad)",
	func(a, b, c, d, e, f Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall6(loadImage,
			intArg(a),
			uintptr(stringArg(b)), // doesn't handle resource id
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f))
		return intRet(rtn)
	})

// dll long User32:GetMessage(
//		MSG* msg, pointer hwnd, long minfilter, long maxfilter)
var getMessage = user32.MustFindProc("GetMessageA").Addr()
var _ = builtin4("GetMessage(msg, hwnd, minfilter, maxfilter)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nMSG)
		rtn := goc.Syscall4(getMessage,
			uintptr(p),
			intArg(b),
			intArg(c),
			intArg(d))
		msg := (*MSG)(p)
		a.Put(nil, SuStr("hwnd"), IntVal(int(msg.hwnd)))
		a.Put(nil, SuStr("message"), IntVal(int(msg.message)))
		a.Put(nil, SuStr("wParam"), IntVal(int(msg.wParam)))
		a.Put(nil, SuStr("lParam"), IntVal(int(msg.lParam)))
		a.Put(nil, SuStr("time"), IntVal(int(msg.time)))
		a.Put(nil, SuStr("pt"), pointToOb(&msg.pt, nil))
		return intRet(rtn)
	})

// dll bool User32:IsDialogMessage(pointer hDlg, MSG* lpMsg)
var isDialogMessage = user32.MustFindProc("IsDialogMessageA").Addr()
var _ = builtin2("IsDialogMessage(hDlg, lpMsg)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := obToMSG(b)
		rtn := goc.Syscall2(isDialogMessage,
			intArg(a),
			uintptr(p))
		return boolRet(rtn)
	})

// dll bool User32:TranslateMessage(MSG* msg)
var translateMessage = user32.MustFindProc("TranslateMessage").Addr()
var _ = builtin1("TranslateMessage(msg)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := obToMSG(a)
		rtn := goc.Syscall1(translateMessage,
			uintptr(p))
		return boolRet(rtn)
	})

// dll long User32:DispatchMessage(MSG* msg)
var dispatchMessage = user32.MustFindProc("DispatchMessageA").Addr()
var _ = builtin1("DispatchMessage(msg)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := obToMSG(a)
		rtn := goc.Syscall1(dispatchMessage,
			uintptr(p))
		return intRet(rtn)
	})

// obToMSG callers must defer heap.FreeTo
func obToMSG(ob Value) unsafe.Pointer {
	p := heap.Alloc(nMSG)
	*(*MSG)(p) = MSG{
		hwnd:    getHandle(ob, "hwnd"),
		message: uint32(getInt(ob, "message")),
		wParam:  getHandle(ob, "wParam"),
		lParam:  getHandle(ob, "lParam"),
		time:    uint32(getInt(ob, "time")),
		pt:      getPoint(ob, "pt"),
	}
	return p
}

// dll long User32:MapWindowPoints(pointer hwndfrom, pointer hwndto, RECT* p, long n)
var mapWindowPoints = user32.MustFindProc("MapWindowPoints").Addr()
var _ = builtin3("MapWindowRect(hwndfrom, hwndto, r)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall4(mapWindowPoints,
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, r)),
			2)
		urectToOb(r, c)
		return intRet(rtn)
	})
