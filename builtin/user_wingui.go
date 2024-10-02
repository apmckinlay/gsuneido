// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"log"
	"syscall"
	"time"
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"golang.org/x/sys/windows"
)

type stRect struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

const nRect = unsafe.Sizeof(stRect{})

type stPaintStruct struct {
	hdc         HANDLE
	fErase      BOOL
	rcPaint     stRect
	fRestore    BOOL
	fIncUpdate  BOOL
	rgbReserved [32]byte
	_           [4]byte // padding
}

const nPaintStruct = unsafe.Sizeof(stPaintStruct{})

type stMonitorInfo struct {
	cbSize    uint32
	rcMonitor stRect
	rcWork    stRect
	dwFlags   uint32
}

const nMonitorInfo = unsafe.Sizeof(stMonitorInfo{})

type stScrollInfo struct {
	cbSize    uint32
	fMask     uint32
	nMin      int32
	nMax      int32
	nPage     uint32
	nPos      int32
	nTrackPos int32
}

const nScrollInfo = unsafe.Sizeof(stScrollInfo{})

type stPoint struct {
	x int32
	y int32
}

const nPoint = unsafe.Sizeof(stPoint{})

type stWindowPlacement struct {
	length           uint32
	flags            uint32
	showCmd          uint32
	ptMinPosition    stPoint
	ptMaxPosition    stPoint
	rcNormalPosition stRect
}

const nWindowPlacement = unsafe.Sizeof(stWindowPlacement{})

type stMenuItemInfo struct {
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

const nMenuItemInfo = unsafe.Sizeof(stMenuItemInfo{})

type stWndClass struct {
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

const nWndClass = unsafe.Sizeof(stWndClass{})

type stCharRange struct {
	cpMin int32
	cpMax int32
}

type stTextRange struct {
	chrg      stCharRange
	lpstrText *byte
}

const nTextRange = unsafe.Sizeof(stTextRange{})

type stMsg struct {
	hwnd    HANDLE
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      stPoint
	_       [4]byte // padding
}

const nMsg = unsafe.Sizeof(stMsg{})

type stToolInfo struct {
	cbSize     uint32
	uFlags     uint32
	hwnd       HANDLE
	uId        uintptr
	rect       stRect
	hinst      HANDLE
	lpszText   *byte
	lParam     uintptr
	lpReserved uintptr
}

const nToolInfo = unsafe.Sizeof(stToolInfo{})

type stToolInfo2 struct {
	cbSize     uint32
	uFlags     uint32
	hwnd       HANDLE
	uId        uintptr
	rect       stRect
	hinst      HANDLE
	lpszText   uintptr // the difference from TOOLINFO
	lParam     uintptr
	lpReserved uintptr
}

type stTVItem struct {
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

type stTVItemEx struct {
	stTVItem
	iIntegral      int32
	uStateEx       int32
	hwnd           HANDLE
	iExpandedImage int32
	iReserved      int32
}

const nTVItemEx = unsafe.Sizeof(stTVItemEx{})

type stTVInsertStruct struct {
	hParent      HANDLE
	hInsertAfter HANDLE
	item         stTVItemEx
}

const nTVInsertStruct = unsafe.Sizeof(stTVInsertStruct{})

type stTVSortCB struct {
	hParent     HANDLE
	lpfnCompare uintptr
	lParam      HANDLE
}

const nTVSortCB = unsafe.Sizeof(stTVSortCB{})

//-------------------------------------------------------------------

// dll User32:GetDesktopWindow() hwnd
var getDesktopWindow = user32.MustFindProc("GetDesktopWindow").Addr()
var _ = builtin(GetDesktopWindow, "()")

func GetDesktopWindow() Value {
	rtn := goc.Syscall0(getDesktopWindow)
	return intRet(rtn)
}

// dll User32:GetSysColor(long nIndex) long
var getSysColor = user32.MustFindProc("GetSysColor").Addr()
var _ = builtin(GetSysColor, "(index)")

func GetSysColor(a Value) Value {
	rtn := goc.Syscall1(getSysColor,
		intArg(a))
	return intRet(rtn)
}

// dll User32:GetWindowRect(pointer hwnd, RECT* rect) bool
var getWindowRect = user32.MustFindProc("GetWindowRect").Addr()
var _ = builtin(GetWindowRectApi, "(hwnd, rect)")

func GetWindowRectApi(a Value, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall2(getWindowRect,
		intArg(a),
		uintptr(r))
	urectToOb(r, b)
	return boolRet(rtn)
}

// dll long User32:MessageBox(pointer window, [in] string text,
// [in] string caption, long flags)
var messageBox = user32.MustFindProc("MessageBoxA").Addr()
var _ = builtin(MessageBox, "(hwnd, text, caption, flags)")

func MessageBox(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall4(messageBox,
		intArg(a),
		uintptr(stringArg(b)),
		uintptr(stringArg(c)),
		intArg(d))
	return intRet(rtn)
}

// dll User32:AdjustWindowRectEx(RECT* rect, long style, bool menu,
// long exStyle) bool
var adjustWindowRectEx = user32.MustFindProc("AdjustWindowRectEx").Addr()
var _ = builtin(AdjustWindowRectEx, "(lpRect, dwStyle, bMenu, dwExStyle)")

func AdjustWindowRectEx(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall4(adjustWindowRectEx,
		uintptr(rectArg(a, r)),
		intArg(b),
		boolArg(c),
		intArg(d))
	urectToOb(r, a)
	return boolRet(rtn)
}

// dll User32:CreateMenu() pointer
var createMenu = user32.MustFindProc("CreateMenu").Addr()
var _ = builtin(CreateMenu, "()")

func CreateMenu() Value {
	rtn := goc.Syscall0(createMenu)
	return intRet(rtn)
}

// dll User32:CreatePopupMenu() pointer
var createPopupMenu = user32.MustFindProc("CreatePopupMenu").Addr()
var _ = builtin(CreatePopupMenu, "()")

func CreatePopupMenu() Value {
	rtn := goc.Syscall0(createPopupMenu)
	return intRet(rtn)
}

// dll User32:AppendMenu(pointer hmenu, long flags, pointer item,
// [in] string name) bool
var appendMenu = user32.MustFindProc("AppendMenuA").Addr()
var _ = builtin(AppendMenu, "(hmenu, flags, item, name)")

func AppendMenu(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall4(appendMenu,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(stringArg(d)))
	return boolRet(rtn)
}

// dll User32:DestroyMenu(pointer hmenu) bool
var destroyMenu = user32.MustFindProc("DestroyMenu").Addr()
var _ = builtin(DestroyMenu, "(hmenu)")

func DestroyMenu(a Value) Value {
	rtn := goc.Syscall1(destroyMenu,
		intArg(a))
	return boolRet(rtn)
}

// dll User32:CreateWindowEx(long exStyle, resource classname, [in] string name,
// long style, long x, long y, long w, long h, pointer parent, pointer menu,
// pointer instance, pointer param) pointer
var createWindowEx = user32.MustFindProc("CreateWindowExA").Addr()
var _ = builtin(CreateWindowEx, "(exStyle, classname, name, style, x, y, w, h,"+" parent, menu, instance, param)")

func CreateWindowEx(_ *Thread, a []Value) Value {
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
}

// dll User32:GetSystemMenu(pointer hWnd, bool bRevert) pointer
var getSystemMenu = user32.MustFindProc("GetSystemMenu").Addr()
var _ = builtin(GetSystemMenu, "(hwnd, bRevert)")

func GetSystemMenu(a, b Value) Value {
	rtn := goc.Syscall2(getSystemMenu,
		intArg(a),
		boolArg(b))
	return intRet(rtn)
}

// dll User32:SetMenu(pointer hwnd, pointer hmenu) bool
var setMenu = user32.MustFindProc("SetMenu").Addr()
var _ = builtin(SetMenu, "(hwnd, hmenu)")

func SetMenu(a, b Value) Value {
	rtn := goc.Syscall2(setMenu,
		intArg(a),
		intArg(b))
	return boolRet(rtn)
}

// dll User32:BeginPaint(pointer hwnd, PAINTSTRUCT* ps) pointer
var beginPaint = user32.MustFindProc("BeginPaint").Addr()
var _ = builtin(BeginPaint, "(hwnd, ps)")

func BeginPaint(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nPaintStruct)
	rtn := goc.Syscall2(beginPaint,
		intArg(a),
		uintptr(psArg(b, p)))
	ps := (*stPaintStruct)(p)
	b.Put(nil, SuStr("hdc"), IntVal(int(ps.hdc)))
	b.Put(nil, SuStr("fErase"), SuBool(ps.fErase != 0))
	b.Put(nil, SuStr("rcPaint"),
		rectToOb(&ps.rcPaint, b.Get(nil, SuStr("rcPaint"))))
	b.Put(nil, SuStr("fRestore"), SuBool(ps.fRestore != 0))
	b.Put(nil, SuStr("fIncUpdate"), SuBool(ps.fIncUpdate != 0))
	return intRet(rtn)
}

// dll User32:EndPaint(pointer hwnd, PAINTSTRUCT* ps) bool
var endPaint = user32.MustFindProc("EndPaint").Addr()
var _ = builtin(EndPaint, "(hwnd, ps)")

func EndPaint(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nPaintStruct)
	rtn := goc.Syscall2(endPaint,
		intArg(a),
		uintptr(psArg(b, p)))
	return boolRet(rtn)
}

func psArg(ob Value, p unsafe.Pointer) unsafe.Pointer {
	ps := (*stPaintStruct)(p)
	ps.hdc = getUintptr(ob, "hdc")
	ps.fErase = getBool(ob, "fErase")
	ps.rcPaint = getRect(ob, "rcPaint")
	ps.fRestore = getBool(ob, "fRestore")
	ps.fIncUpdate = getBool(ob, "fIncUpdate")
	return p
}

// dll User32:CallWindowProc(pointer wndprcPrev, pointer hwnd, long msg,
// pointer wParam, pointer lParam) pointer
var callWindowProc = user32.MustFindProc("CallWindowProcA").Addr()
var _ = builtin(CallWindowProc, "(wndprcPrev, hwnd, msg, wParam, lParam)")

func CallWindowProc(th *Thread, a []Value) Value {
	if a[0].Type() != types.Number {
		// presumably a previous callback returned by SetWindowProc
		return th.Call(a[0], a[1], a[2], a[3], a[4])
	}
	rtn := goc.Syscall5(callWindowProc,
		intArg(a[0]),
		intArg(a[1]),
		intArg(a[2]),
		intArg(a[3]),
		intArg(a[4]))
	return intRet(rtn)
}

// dll User32:CreateAcceleratorTable([in] string lpaccel, long cEntries) pointer
var createAcceleratorTable = user32.MustFindProc("CreateAcceleratorTableA").Addr()
var _ = builtin(CreateAcceleratorTable, "(lpaccel, cEntries)")

func CreateAcceleratorTable(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall2(createAcceleratorTable,
		uintptr(stringArg(a)),
		intArg(b))
	return intRet(rtn)
}

// dll User32:DestroyAcceleratorTable(pointer hAccel) bool
var destroyAcceleratorTable = user32.MustFindProc("DestroyAcceleratorTable").Addr()
var _ = builtin(DestroyAcceleratorTable, "(hAccel)")

func DestroyAcceleratorTable(a Value) Value {
	rtn := goc.Syscall1(destroyAcceleratorTable,
		intArg(a))
	return boolRet(rtn)
}

// dll User32:DestroyWindow(pointer hwnd) bool
var destroyWindow = user32.MustFindProc("DestroyWindow").Addr()
var _ = builtin(DestroyWindow, "(hwnd)")

func DestroyWindow(a Value) Value {
	rtn := goc.Syscall1(destroyWindow,
		intArg(a))
	return boolRet(rtn)
}

// dll User32:DrawFrameControl(pointer hdc, RECT* lprc, long uType,
// long uState) bool
var drawFrameControl = user32.MustFindProc("DrawFrameControl").Addr()
var _ = builtin(DrawFrameControl, "(hdc, lprc, uType, uState)")

func DrawFrameControl(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall4(drawFrameControl,
		intArg(a),
		uintptr(rectArg(b, r)),
		intArg(c),
		intArg(d))
	return boolRet(rtn)
}

// dll User32:DrawText(pointer hdc, [in] string lpsz, long cb, RECT* lprc,
// long uFormat) long
var drawText = user32.MustFindProc("DrawTextA").Addr()
var _ = builtin(DrawText, "(hdc, lpsz, cb, lprc, uFormat)")

func DrawText(a, b, c, d, e Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall5(drawText,
		intArg(a),
		uintptr(stringArg(b)),
		intArg(c),
		uintptr(rectArg(d, r)),
		intArg(e))
	urectToOb(r, d) // for CALCRECT
	return intRet(rtn)
}

// dll User32:FillRect(pointer hdc, RECT* lpRect, pointer hBrush) long
var fillRect = user32.MustFindProc("FillRect").Addr()
var _ = builtin(FillRect, "(hdc, lpRect, hBrush)")

func FillRect(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	assert.That(b != Zero)
	r := heap.Alloc(nRect)
	rtn := goc.Syscall3(fillRect,
		intArg(a),
		uintptr(rectArg(b, r)),
		intArg(c))
	return intRet(rtn)
}

// dll User32:GetActiveWindow() pointer
var getActiveWindow = user32.MustFindProc("GetActiveWindow").Addr()
var _ = builtin(GetActiveWindow, "()")

func GetActiveWindow() Value {
	rtn := goc.Syscall0(getActiveWindow)
	return intRet(rtn)
}

// dll User32:GetFocus() pointer
var getFocus = user32.MustFindProc("GetFocus").Addr()
var _ = builtin(GetFocus, "()")

func GetFocus() Value {
	rtn := goc.Syscall0(getFocus)
	return intRet(rtn)
}

// dll User32:GetClientRect(pointer hwnd, RECT* rect) bool
var getClientRect = user32.MustFindProc("GetClientRect").Addr()
var _ = builtin(GetClientRect, "(hwnd, rect)")

func GetClientRect(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall2(getClientRect,
		intArg(a),
		uintptr(r))
	urectToOb(r, b)
	return boolRet(rtn)
}

// dll User32:GetDC(pointer hwnd) pointer
var getDC = user32.MustFindProc("GetDC").Addr()
var _ = builtin(GetDC, "(hwnd)")

func GetDC(a Value) Value {
	rtn := goc.Syscall1(getDC,
		intArg(a))
	return intRet(rtn)
}

// dll User32:GetMonitorInfo(pointer hMonitor, MONITORINFO* lpmi) bool
var getMonitorInfo = user32.MustFindProc("GetMonitorInfoA").Addr()
var _ = builtin(GetMonitorInfoApi, "(hwnd, mInfo)")

func GetMonitorInfoApi(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nMonitorInfo)
	mi := (*stMonitorInfo)(p)
	mi.cbSize = uint32(nMonitorInfo)
	rtn := goc.Syscall2(getMonitorInfo,
		intArg(a),
		uintptr(p))
	b.Put(nil, SuStr("rcMonitor"), rectToOb(&mi.rcMonitor, nil))
	b.Put(nil, SuStr("rcWork"), rectToOb(&mi.rcWork, nil))
	b.Put(nil, SuStr("dwFlags"), IntVal(int(mi.dwFlags)))
	return boolRet(rtn)
}

// dll User32:GetScrollInfo(pointer hwnd, long fnBar, SCROLLINFO* lpsi) bool
var getScrollInfo = user32.MustFindProc("GetScrollInfo").Addr()
var _ = builtin(GetScrollInfo, "(hwnd, fnBar, lpsi)")

func GetScrollInfo(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nScrollInfo)
	si := (*stScrollInfo)(p)
	*si = stScrollInfo{
		cbSize:    uint32(nScrollInfo),
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
}

// dll User32:GetScrollPos(pointer hwnd, int nBar) int
var getScrollPos = user32.MustFindProc("GetScrollPos").Addr()
var _ = builtin(GetScrollPos, "(hwnd, nBar)")

func GetScrollPos(a, b Value) Value {
	rtn := goc.Syscall2(getScrollPos,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll User32:GetSysColorBrush(long nIndex) pointer
var getSysColorBrush = user32.MustFindProc("GetSysColorBrush").Addr()
var _ = builtin(GetSysColorBrush, "(nIndex)")

func GetSysColorBrush(a Value) Value {
	rtn := goc.Syscall1(getSysColorBrush,
		intArg(a))
	return intRet(rtn)
}

// dll User32:GetSystemMetrics(long nIndex) long
var getSystemMetrics = user32.MustFindProc("GetSystemMetrics").Addr()
var _ = builtin(GetSystemMetrics, "(nIndex)")

func GetSystemMetrics(a Value) Value {
	rtn := goc.Syscall1(getSystemMetrics,
		intArg(a))
	return intRet(rtn)
}

// dll User32:GetWindowLong(pointer hwnd, long offset) long
var getWindowLong = user32.MustFindProc("GetWindowLongA").Addr()
var _ = builtin(GetWindowLong, "(hwnd, offset)")

func GetWindowLong(a, b Value) Value {
	rtn := goc.Syscall2(getWindowLong,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll User32:GetWindowLong(pointer hwnd, long offset) long
var getWindowLongPtr = user32.MustFindProc("GetWindowLongPtrA").Addr()
var _ = builtin(GetWindowLongPtr, "(hwnd, offset)")

func GetWindowLongPtr(a, b Value) Value {
	rtn := goc.Syscall2(getWindowLongPtr,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll User32:GetWindowPlacement(pointer hwnd, WINDOWPLACEMENT* lpwndpl) bool
var getWindowPlacement = user32.MustFindProc("GetWindowPlacement").Addr()
var _ = builtin(GetWindowPlacement, "(hwnd, ps)")

func GetWindowPlacement(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nWindowPlacement)
	wp := (*stWindowPlacement)(p)
	wp.length = uint32(nWindowPlacement)
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
}

// dll User32:GetWindowText(pointer hwnd, string buf, long len) long
var getWindowText = user32.MustFindProc("GetWindowTextA").Addr()
var getWindowTextLength = user32.MustFindProc("GetWindowTextLengthA").Addr()
var _ = builtin(GetWindowText, "(hwnd)")

func GetWindowText(hwnd Value) Value {
	defer heap.FreeTo(heap.CurSize())
	n := goc.Syscall1(getWindowTextLength,
		intArg(hwnd))
	buf := heap.Alloc(n + 1)
	n = goc.Syscall3(getWindowText,
		intArg(hwnd),
		uintptr(buf),
		n+1)
	return SuStr(heap.GetStrN(buf, int(n)))
}

// dll bool User32:HideCaret(pointer hWnd)
var hideCaret = user32.MustFindProc("HideCaret").Addr()
var _ = builtin(HideCaret, "(hwnd)")

func HideCaret(hwnd Value) Value {
	rtn := goc.Syscall1(hideCaret,
		intArg(hwnd))
	return boolRet(rtn)
}

// dll User32:InflateRect(RECT* rect, long dx, long dy) bool
var inflateRect = user32.MustFindProc("InflateRect").Addr()
var _ = builtin(InflateRect, "(rect, dx, dy)")

func InflateRect(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall3(inflateRect,
		uintptr(rectArg(a, r)),
		intArg(b),
		intArg(c))
	urectToOb(r, a)
	return boolRet(rtn)
}

// dll User32:InsertMenuItem(pointer hMenu, long uItem, bool fByPosition,
// MENUITEMINFO* lpmii) bool
var insertMenuItem = user32.MustFindProc("InsertMenuItemA").Addr()
var _ = builtin(InsertMenuItem, "(hMenu, uItem, fByPosition, lpmii)")

func InsertMenuItem(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nMenuItemInfo)
	*(*stMenuItemInfo)(p) = stMenuItemInfo{
		cbSize:        uint32(nMenuItemInfo),
		fMask:         getUint32(d, "fMask"),
		fType:         getUint32(d, "fType"),
		fState:        getUint32(d, "fState"),
		wID:           getUint32(d, "wID"),
		hSubMenu:      getUintptr(d, "hSubMenu"),
		hbmpChecked:   getUintptr(d, "hbmpChecked"),
		hbmpUnchecked: getUintptr(d, "hbmpUnchecked"),
		dwItemData:    getUintptr(d, "dwItemData"),
		dwTypeData:    getStr(d, "dwTypeData"),
		cch:           getUint32(d, "cch"),
		hbmpItem:      getUintptr(d, "hbmpItem"),
	}
	rtn := goc.Syscall4(insertMenuItem,
		intArg(a),
		intArg(b),
		boolArg(c),
		uintptr(p))
	return boolRet(rtn)
}

// dll long User32:GetMenuItemCount(pointer hMenu)
var getMenuItemCount = user32.MustFindProc("GetMenuItemCount").Addr()
var _ = builtin(GetMenuItemCount, "(hMenu)")

func GetMenuItemCount(a Value) Value {
	rtn := goc.Syscall1(getMenuItemCount,
		intArg(a))
	return intRet(rtn)
}

// dll long User32:GetMenuItemID(pointer hMenu, long nPos)
var getMenuItemID = user32.MustFindProc("GetMenuItemID").Addr()
var _ = builtin(GetMenuItemID, "(hMenu, nPos)")

func GetMenuItemID(a, b Value) Value {
	rtn := goc.Syscall2(getMenuItemID,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll User32:GetMenuItemInfo(pointer hMenu, long uItem, bool fByPosition,
//
//	MENUITEMINFO* lpmii) bool
var getMenuItemInfo = user32.MustFindProc("GetMenuItemInfoA").Addr()
var _ = builtin(GetMenuItemInfoText, "(hMenu, uItem)")

func GetMenuItemInfoText(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	const MMIM_TYPE = 0x10
	const MFT_STRING = 0
	p := heap.Alloc(nMenuItemInfo)
	mii := (*stMenuItemInfo)(p)
	*mii = stMenuItemInfo{
		cbSize:     uint32(nMenuItemInfo),
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
	if rtn == 0 {
		return False
	}
	return SuStr(heap.GetStrN(buf, int(n-1))) // -1 to omit nul terminator
}

// dll User32:SetMenuItemInfo(pointer hMenu, long uItem, long fByPosition,
// MENUITEMINFO* lpmii) bool
var setMenuItemInfo = user32.MustFindProc("SetMenuItemInfoA").Addr()
var _ = builtin(SetMenuItemInfo, "(hMenu, uItem, fByPosition, lpmii)")

func SetMenuItemInfo(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nMenuItemInfo)
	*(*stMenuItemInfo)(p) = stMenuItemInfo{
		cbSize:        uint32(nMenuItemInfo),
		fMask:         getUint32(d, "fMask"),
		fType:         getUint32(d, "fType"),
		fState:        getUint32(d, "fState"),
		wID:           getUint32(d, "wID"),
		hSubMenu:      getUintptr(d, "hSubMenu"),
		hbmpChecked:   getUintptr(d, "hbmpChecked"),
		hbmpUnchecked: getUintptr(d, "hbmpUnchecked"),
		dwItemData:    getUintptr(d, "dwItemData"),
		//dwTypeData:    getStr(d, "dwTypeData"),
		cch:      getUint32(d, "cch"),
		hbmpItem: getUintptr(d, "hbmpItem"),
	}
	rtn := goc.Syscall4(setMenuItemInfo,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	return boolRet(rtn)
}

// dll User32:InvalidateRect(pointer hwnd, RECT* rect, bool erase) bool
var invalidateRect = user32.MustFindProc("InvalidateRect").Addr()
var _ = builtin(InvalidateRect, "(hwnd, rect, erase)")

func InvalidateRect(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall3(invalidateRect,
		intArg(a),
		uintptr(rectArg(b, r)),
		boolArg(c))
	return boolRet(rtn)
}

// dll User32:IsWindowEnabled(pointer hwnd) bool
var isWindowEnabled = user32.MustFindProc("IsWindowEnabled").Addr()
var _ = builtin(IsWindowEnabled, "(hwnd)")

func IsWindowEnabled(a Value) Value {
	rtn := goc.Syscall1(isWindowEnabled,
		intArg(a))
	return boolRet(rtn)
}

// dll User32:LoadCursor(pointer hinst, resource pszCursor) pointer
var loadCursor = user32.MustFindProc("LoadCursorA").Addr()
var _ = builtin(LoadCursor, "(hinst, pszCursor)")

func LoadCursor(a, b Value) Value {
	rtn := goc.Syscall2(loadCursor,
		intArg(a),
		intArg(b)) // could be a string but we never use that
	return intRet(rtn)
}

// dll User32:LoadIcon(pointer hInstance, resource lpIconName) pointer
var loadIcon = user32.MustFindProc("LoadIconA").Addr()
var _ = builtin(LoadIcon, "(hinst, lpIconName)")

func LoadIcon(a, b Value) Value {
	rtn := goc.Syscall2(loadIcon,
		intArg(a),
		intArg(b)) // could be a string but we never use that
	return intRet(rtn)
}

// dll User32:MonitorFromRect(RECT* lprc, long dwFlags) pointer
var monitorFromRect = user32.MustFindProc("MonitorFromRect").Addr()
var _ = builtin(MonitorFromRect, "(lprc, dwFlags)")

func MonitorFromRect(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall2(monitorFromRect,
		uintptr(rectArg(a, r)),
		intArg(b))
	return intRet(rtn)
}

// dll User32:MoveWindow(pointer hwnd, long left, long top, long width,
// long height, bool repaint) bool
var moveWindow = user32.MustFindProc("MoveWindow").Addr()
var _ = builtin(MoveWindow, "(hwnd, left, top, width, height, repaint)")

func MoveWindow(a, b, c, d, e, f Value) Value {
	rtn := goc.Syscall6(moveWindow,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e),
		boolArg(f))
	return boolRet(rtn)
}

// dll User32:RegisterClass(WNDCLASS* wc) short
var registerClass = user32.MustFindProc("RegisterClassA").Addr()
var _ = builtin(RegisterClass, "(wc)")

func RegisterClass(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nWndClass)
	*(*stWndClass)(p) = stWndClass{
		style:      getUint32(a, "style"),
		wndProc:    getUintptr(a, "wndProc"),
		clsExtra:   getInt32(a, "clsExtra"),
		wndExtra:   getInt32(a, "wndExtra"),
		instance:   getUintptr(a, "instance"),
		icon:       getUintptr(a, "icon"),
		cursor:     getUintptr(a, "cursor"),
		background: getUintptr(a, "background"),
		menuName:   getStr(a, "menuName"),
		className:  getStr(a, "className"),
	}
	rtn := goc.Syscall1(registerClass,
		uintptr(p))
	return intRet(rtn)
}

// dll User32:RegisterClipboardFormat([in] string lpszFormat) long
var registerClipboardFormat = user32.MustFindProc("RegisterClipboardFormatA").Addr()
var _ = builtin(RegisterClipboardFormat, "(lpszFormat)")

func RegisterClipboardFormat(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall1(registerClipboardFormat,
		uintptr(stringArg(a)))
	return intRet(rtn)
}

// dll User32:ReleaseDC(pointer hWnd, pointer hDC) long
var releaseDC = user32.MustFindProc("ReleaseDC").Addr()
var _ = builtin(ReleaseDC, "(hwnd, hDC)")

func ReleaseDC(a, b Value) Value {
	rtn := goc.Syscall2(releaseDC,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll User32:SetFocus(pointer hwnd) pointer
var setFocus = user32.MustFindProc("SetFocus").Addr()
var _ = builtin(SetFocus, "(hwnd)")

func SetFocus(a Value) Value {
	rtn := goc.Syscall1(setFocus,
		intArg(a))
	return intRet(rtn)
}

const WM_USER = 0x400

// dll User32:SetTimer(pointer hwnd, long id, long ms, TIMERPROC f) long
var setTimer = user32.MustFindProc("SetTimer").Addr()
var _ = builtin(SetTimer, "(hwnd, id, ms, f)")

func SetTimer(a, b, c, d Value) Value {
	if options.TimersDisabled {
		return Zero
	}
	if windows.GetCurrentThreadId() == uiThreadId {
		return gocSetTimer(a, b, c, d)
	}
	// WARNING: don't use heap from background thread
	d.SetConcurrent() // since callback will be from different thread
	ret := make(chan Value, 1)
	rogsChan <- func() {
		ret <- gocSetTimer(a, b, c, d)
	}
	notifyCside()
	first := true
	for {
		select {
		case id := <-ret:
			return id
		case <-time.After(5 * time.Second):
			if first {
				first = false
				log.Println("SetTimer timeout")
			}
		}
	}
}

var nTimer = 0

const warnTimers = 32
const maxTimers = 64

var _ = AddInfo("windows.nTimer", &nTimer)

// gocSetTimer is called by SetTimer directly if on main UI thread
// and via runOnGoSide if from another thread
func gocSetTimer(hwnd, id, ms, cb Value) Value {
	if nTimer > warnTimers {
		if nTimer > maxTimers {
			log.Panicln("ERROR: SetTimer: over", maxTimers)
		}
		log.Println("WARNING: SetTimer: over", warnTimers)
	}
	rtn := goc.Syscall4(setTimer,
		intArg(hwnd),
		intArg(id),
		intArg(ms),
		NewCallback(cb, 4))
	if rtn != 0 {
		nTimer++
	}
	return intRet(rtn)
}

// dll User32:KillTimer(pointer hwnd, long id) bool
var killTimer = user32.MustFindProc("KillTimer").Addr()
var _ = builtin(KillTimer, "(hwnd, id)")

func KillTimer(a, b Value) Value {
	if options.TimersDisabled {
		return False
	}
	if windows.GetCurrentThreadId() == uiThreadId {
		return gocKillTimer(a, b)
	}
	ret := make(chan Value, 1)
	rogsChan <- func() {
		ret <- gocKillTimer(a, b)
	}
	notifyCside()
	first := true
	for {
		select {
		case id := <-ret:
			return id
		case <-time.After(5 * time.Second):
			if first {
				first = false
				log.Println("KillTimer timeout")
			}
		}
	}
}

// gocKillTimer is called by KillTimer directly if on main UI thread
// and via runOnGoSide if from another thread
func gocKillTimer(hwnd, id Value) Value {
	rtn := goc.Syscall2(killTimer,
		intArg(hwnd),
		intArg(id))
	if rtn != 0 {
		nTimer--
	}
	return boolRet(rtn)
}

const notifyMsg = WM_USER

// notifyCside is used by SetTimer and KillTimer
// It uses PostMessage (high priority) to C side
// to handle when we're running in the message loop.
func notifyCside() {
	// NOTE: this has to be the Go Syscall, not goc.Syscall
	r, _, _ := syscall.SyscallN(postMessage,
		goc.CHelperHwnd(), notifyMsg, 0, 0)
	if r == 0 {
		log.Panicln("notifyCside PostMessage failed")
	}
}

// dll User32:SetWindowLong(pointer hwnd, int offset, long value) long
var setWindowLong = user32.MustFindProc("SetWindowLongA").Addr()
var _ = builtin(SetWindowLong, "(hwnd, offset, value)")

func SetWindowLong(a, b, c Value) Value {
	rtn := goc.Syscall3(setWindowLong,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll User32:SetWindowLong(pointer hwnd, long offset, long value) long
var setWindowLongPtr = user32.MustFindProc("SetWindowLongPtrA").Addr()
var _ = builtin(SetWindowLongPtr, "(hwnd, offset, value)")

func SetWindowLongPtr(a, b, c Value) Value {
	rtn := goc.Syscall3(setWindowLongPtr,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll User32:SetWindowProc(pointer hwnd, long offset, WNDPROC proc) pointer
var _ = builtin(SetWindowProc, "(hwnd, offset, proc)")

func SetWindowProc(a, b, c Value) Value {
	hwnd := intArg(a)
	var cb uintptr
	var fn Value
	if c.Type() == types.Number {
		cb = uintptr(ToInt(c))
	} else {
		fn = hwndToCb[hwnd] // save the old one in case we're overwriting
		cb = WndProcCallback(hwnd, c)
	}
	rtn := goc.Syscall3(setWindowLongPtr,
		hwnd,
		intArg(b),
		cb)
	if rtn == wndProcCb && fn != nil { // if overwriting
		return fn // return the actual previous Suneido callback
	}
	return intRet(rtn)
}

// dll User32:SetWindowPlacement(pointer hwnd, WINDOWPLACEMENT* lpwndpl) bool
var setWindowPlacement = user32.MustFindProc("SetWindowPlacement").Addr()
var _ = builtin(SetWindowPlacement, "(hwnd, lpwndpl)")

func SetWindowPlacement(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nWindowPlacement)
	*(*stWindowPlacement)(p) = stWindowPlacement{
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
}

// dll User32:SetWindowPos(pointer hWnd, pointer hWndInsertAfter,
// long X, long Y, long cx, long cy, long uFlags) bool
var setWindowPos = user32.MustFindProc("SetWindowPos").Addr()
var _ = builtin(SetWindowPos, "(hWnd, hWndInsertAfter, X, Y, cx, cy, uFlags)")

func SetWindowPos(a, b, c, d, e, f, g Value) Value {
	rtn := goc.Syscall7(setWindowPos,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e),
		intArg(f),
		intArg(g))
	return boolRet(rtn)
}

// dll User32:SetWindowText(pointer hwnd, [in] string text) bool
var setWindowText = user32.MustFindProc("SetWindowTextA").Addr()
var _ = builtin(SetWindowText, "(hwnd, lpwndpl)")

func SetWindowText(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall2(setWindowText,
		intArg(a),
		uintptr(stringArg(b)))
	return boolRet(rtn)
}

// dll User32:ShowWindow(pointer hwnd, long ncmd) bool
var showWindow = user32.MustFindProc("ShowWindow").Addr()
var _ = builtin(ShowWindow, "(hwnd, ncmd)")

func ShowWindow(a, b Value) Value {
	rtn := goc.Syscall2(showWindow,
		intArg(a),
		intArg(b))
	return boolRet(rtn)
}

// dll User32:SystemParametersInfo(long uiAction, long uiParam, ? pvParam,
// long fWinIni) bool
var systemParametersInfo = user32.MustFindProc("SystemParametersInfoA").Addr()

var _ = builtin(SPI_GetFocusBorderHeight, "()")

func SPI_GetFocusBorderHeight() Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(4)
	const SPI_GETFOCUSBORDERHEIGHT = 0x2010
	goc.Syscall4(systemParametersInfo,
		SPI_GETFOCUSBORDERHEIGHT,
		0,
		uintptr(p),
		0)
	return IntVal(int((*(*int32)(p))))
}

var _ = builtin(SPI_GetWheelScrollLines, "()")

func SPI_GetWheelScrollLines() Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(4)
	const SPI_GETWHEELSCROLLLINES = 104
	goc.Syscall4(systemParametersInfo,
		SPI_GETWHEELSCROLLLINES,
		0,
		uintptr(p),
		0)
	return IntVal(int((*(*int32)(p))))
}

var _ = builtin(SPI_GetWorkArea, "()")

func SPI_GetWorkArea() Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	const SPI_GETWORKAREA = 48
	goc.Syscall4(systemParametersInfo,
		SPI_GETWORKAREA,
		0,
		uintptr(r),
		0)
	return urectToOb(r, nil)
}

// dll User32:PostQuitMessage(long exitcode) void
var postQuitMessage = user32.MustFindProc("PostQuitMessage").Addr()

func postQuit(arg uintptr) {
	goc.Syscall1(postQuitMessage, arg)
}

var _ = builtin(PostQuitMessage, "(exitcode)")

func PostQuitMessage(a Value) Value {
	postQuit(intArg(a))
	return nil
}

// dll User32:GetNextDlgTabItem(pointer hDlg, pointer hCtl, bool prev) pointer
var getNextDlgTabItem = user32.MustFindProc("GetNextDlgTabItem").Addr()
var _ = builtin(GetNextDlgTabItem, "(hDlg, hCtl, prev)")

func GetNextDlgTabItem(a, b, c Value) Value {
	rtn := goc.Syscall3(getNextDlgTabItem,
		intArg(a),
		intArg(b),
		boolArg(c))
	return intRet(rtn)
}

// dll User32:UpdateWindow(pointer hwnd) bool
var updateWindow = user32.MustFindProc("UpdateWindow").Addr()
var _ = builtin(UpdateWindow, "(hwnd)")

func UpdateWindow(a Value) Value {
	rtn := goc.Syscall1(updateWindow,
		intArg(a))
	return boolRet(rtn)
}

// dll User32:DefWindowProc(pointer hwnd, long msg, pointer wParam,
// pointer lParam) pointer
var defWindowProc = user32.MustFindProc("DefWindowProcA").Addr()
var _ = builtin(DefWindowProc, "(hwnd, msg, wParam, lParam)")

func DefWindowProc(a, b, c, d Value) Value {
	rtn := goc.Syscall4(defWindowProc,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d))
	return intRet(rtn)
}

var _ = builtin(GetDefWindowProc, "()")

func GetDefWindowProc() Value {
	return IntVal(int(defWindowProc))
}

// dll User32:GetKeyState(long key) short
var getKeyState = user32.MustFindProc("GetKeyState").Addr()
var _ = builtin(GetKeyState, "(nVirtKey)")

func GetKeyState(a Value) Value {
	rtn := goc.Syscall1(getKeyState,
		intArg(a))
	return intRet(rtn)
}

type stTPMParams struct {
	cbSize    int32
	rcExclude stRect
}

const nTPMParams = unsafe.Sizeof(stTPMParams{})

// dll long User32:TrackPopupMenuEx(pointer hmenu, long fuFlags, long x, long y,
// pointer hwnd, TPMPARAMS* lptpm)
var trackPopupMenuEx = user32.MustFindProc("TrackPopupMenuEx").Addr()
var _ = builtin(TrackPopupMenuEx, "(hmenu, fuFlags, x, y, hwnd, lptpm)")

func TrackPopupMenuEx(a, b, c, d, e, f Value) Value {
	var p unsafe.Pointer
	if !f.Equal(Zero) {
		defer heap.FreeTo(heap.CurSize())
		p = heap.Alloc(nTPMParams)
		*(*stTPMParams)(p) = stTPMParams{
			cbSize:    int32(nTPMParams),
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
}

// dll bool User32:OpenClipboard(pointer hwnd)
var openClipboard = user32.MustFindProc("OpenClipboard").Addr()
var _ = builtin(OpenClipboard, "(hwnd)")

func OpenClipboard(a Value) Value {
	rtn := goc.Syscall1(openClipboard,
		intArg(a))
	return boolRet(rtn)
}

// dll bool User32:EmptyClipboard()
var emptyClipboard = user32.MustFindProc("EmptyClipboard").Addr()
var _ = builtin(EmptyClipboard, "()")

func EmptyClipboard() Value {
	rtn := goc.Syscall0(emptyClipboard)
	return boolRet(rtn)
}

// dll pointer User32:GetClipboardData(long format)
var getClipboardData = user32.MustFindProc("GetClipboardData").Addr()
var _ = builtin(GetClipboardData, "(format)")

func GetClipboardData(a Value) Value {
	rtn := goc.Syscall1(getClipboardData,
		intArg(a))
	return intRet(rtn)
}

// dll pointer User32:SetClipboardData(long uFormat, pointer hMem)
var setClipboardData = user32.MustFindProc("SetClipboardData").Addr()
var _ = builtin(SetClipboardData, "(uFormat, hMem)")

func SetClipboardData(a Value, b Value) Value {
	rtn := goc.Syscall2(setClipboardData,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll bool User32:CloseClipboard()
var closeClipboard = user32.MustFindProc("CloseClipboard").Addr()
var _ = builtin(CloseClipboard, "()")

func CloseClipboard() Value {
	rtn := goc.Syscall0(closeClipboard)
	return boolRet(rtn)
}

// dll pointer User32:BeginDeferWindowPos(
// long nNumWindows) // initial number of windows to allocate space for
var beginDeferWindowPos = user32.MustFindProc("BeginDeferWindowPos").Addr()
var _ = builtin(BeginDeferWindowPos, "(nNumWindows)")

func BeginDeferWindowPos(a Value) Value {
	rtn := goc.Syscall1(beginDeferWindowPos,
		intArg(a))
	return intRet(rtn)
}

// dll pointer User32:CallNextHookEx(		// returns an LRESULT
// pointer	hhk,	// handle to current hook [HHOOK]
// long	nCode,	// hook code passed to hook procedure [int]
// pointer	wParam,	// value passed to hook procedure [WPARAM]
// pointer	lParam)	// value passed to hook procedure [LPARAM]
var callNextHookEx = user32.MustFindProc("CallNextHookEx").Addr()
var _ = builtin(CallNextHookEx, "(hhk, nCode, wParam, lParam)")

func CallNextHookEx(a, b, c, d Value) Value {
	rtn := goc.Syscall4(callNextHookEx,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d))
	return intRet(rtn)
}

// dll long User32:CheckMenuItem(
// pointer hmenu, 		// handle to menu
// long uIDCheckItem, 	// menu item to check or uncheck
// long uCheck)		// menu item options
var checkMenuItem = user32.MustFindProc("CheckMenuItem").Addr()
var _ = builtin(CheckMenuItem, "(hmenu, uIDCheckItem, uCheck)")

func CheckMenuItem(a, b, c Value) Value {
	rtn := goc.Syscall3(checkMenuItem,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll bool user32:DeleteMenu(pointer hMenu, long uPosition, long uFlags)
var deleteMenu = user32.MustFindProc("DeleteMenu").Addr()
var _ = builtin(DeleteMenu, "(hMenu, uPosition, uFlags)")

func DeleteMenu(a, b, c Value) Value {
	rtn := goc.Syscall3(deleteMenu,
		intArg(a),
		intArg(b),
		intArg(c))
	return boolRet(rtn)
}

// dll long User32:EnableMenuItem(
// pointer hMenu,
// long uIDEnableItem,
// long uEnable)
var enableMenuItem = user32.MustFindProc("EnableMenuItem").Addr()
var _ = builtin(EnableMenuItem, "(hMenu, uIDEnableItem, uEnable)")

func EnableMenuItem(a, b, c Value) Value {
	rtn := goc.Syscall3(enableMenuItem,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll bool User32:EnableWindow(pointer hWnd, bool bEnable)
var enableWindow = user32.MustFindProc("EnableWindow").Addr()
var _ = builtin(EnableWindow, "(hWnd, bEnable)")

func EnableWindow(a, b Value) Value {
	rtn := goc.Syscall2(enableWindow,
		intArg(a),
		boolArg(b))
	return boolRet(rtn)
}

// dll bool User32:EndDeferWindowPos(pointer hWinPosInfo)
var endDeferWindowPos = user32.MustFindProc("EndDeferWindowPos").Addr()
var _ = builtin(EndDeferWindowPos, "(hWinPosInfo)")

func EndDeferWindowPos(a Value) Value {
	rtn := goc.Syscall1(endDeferWindowPos,
		intArg(a))
	return boolRet(rtn)
}

// dll bool User32:EndDialog(pointer hwndDlg, long nResult)
var endDialog = user32.MustFindProc("EndDialog").Addr()
var _ = builtin(EndDialog, "(hwndDlg, nResult)")

func EndDialog(a, b Value) Value {
	rtn := goc.Syscall2(endDialog,
		intArg(a),
		intArg(b))
	return boolRet(rtn)
}

// dll long User32:EnumClipboardFormats(long format)
var enumClipboardFormats = user32.MustFindProc("EnumClipboardFormats").Addr()
var _ = builtin(EnumClipboardFormats, "(format)")

func EnumClipboardFormats(a Value) Value {
	rtn := goc.Syscall1(enumClipboardFormats,
		intArg(a))
	return intRet(rtn)
}

// dll pointer User32:FindWindow([in] string c, [in] string n)
var findWindow = user32.MustFindProc("FindWindowA").Addr()
var _ = builtin(FindWindow, "(c, n)")

func FindWindow(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall2(findWindow,
		uintptr(stringArg(a)),
		uintptr(stringArg(b)))
	return intRet(rtn)
}

// dll pointer User32:GetAncestor(pointer hwnd, long gaFlags)
var getAncestor = user32.MustFindProc("GetAncestor").Addr()
var _ = builtin(GetAncestor, "(hwnd, gaFlags)")

func GetAncestor(a, b Value) Value {
	rtn := goc.Syscall2(getAncestor,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll long User32:GetClipboardFormatName(
// long   format,			// clipboard format to retrieve
// string lpszFormatName,	// buffer to receive format name
// long   cchMaxCount)		// maximum length of string to copy into buffer
var getClipboardFormatName = user32.MustFindProc("GetClipboardFormatNameA").Addr()
var _ = builtin(GetClipboardFormatName, "(format, lpszFormatName, cchMaxCount)")

func GetClipboardFormatName(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall3(getClipboardFormatName,
		intArg(a),
		uintptr(stringArg(b)),
		intArg(c))
	return intRet(rtn)
}

// dll pointer User32:GetCursor()
var getCursor = user32.MustFindProc("GetCursor").Addr()
var _ = builtin(GetCursor, "()")

func GetCursor() Value {
	rtn := goc.Syscall0(getCursor)
	return intRet(rtn)
}

// dll long User32:GetDoubleClickTime()

var getDoubleClickTime = user32.MustFindProc("GetDoubleClickTime").Addr()
var _ = builtin(GetDoubleClickTime, "()")

func GetDoubleClickTime() Value {
	rtn := goc.Syscall0(getDoubleClickTime)
	return intRet(rtn)
}

// dll long User32:GetMenuState(
// pointer hMenu, 	// handle to menu
// long uId, 		// menu item to query
// long uFlags)	// options
var getMenuState = user32.MustFindProc("GetMenuState").Addr()
var _ = builtin(GetMenuState, "(hMenu, uId, uFlags)")

func GetMenuState(a, b, c Value) Value {
	rtn := goc.Syscall3(getMenuState,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll long User32:GetMessagePos()
var getMessagePos = user32.MustFindProc("GetMessagePos").Addr()
var _ = builtin(GetMessagePos, "()")

func GetMessagePos() Value {
	rtn := goc.Syscall0(getMessagePos)
	return intRet(rtn)
}

// dll pointer User32:GetNextDlgGroupItem(pointer hDlg, pointer hCtl, bool prev)
var getNextDlgGroupItem = user32.MustFindProc("GetNextDlgGroupItem").Addr()
var _ = builtin(GetNextDlgGroupItem, "(hDlg, hCtl, prev)")

func GetNextDlgGroupItem(a, b, c Value) Value {
	rtn := goc.Syscall3(getNextDlgGroupItem,
		intArg(a),
		intArg(b),
		boolArg(c))
	return intRet(rtn)
}

// dll pointer User32:GetParent(pointer hwnd)
var getParent = user32.MustFindProc("GetParent").Addr()
var _ = builtin(GetParent, "(hwnd)")

func GetParent(a Value) Value {
	rtn := goc.Syscall1(getParent,
		intArg(a))
	return intRet(rtn)
}

// dll pointer User32:GetSubMenu(
// pointer hmenu,	//menu handle
// long position)	//position
var getSubMenu = user32.MustFindProc("GetSubMenu").Addr()
var _ = builtin(GetSubMenu, "(hmenu, position)")

func GetSubMenu(a, b Value) Value {
	rtn := goc.Syscall2(getSubMenu,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll pointer User32:GetTopWindow(pointer hWnd)
var getTopWindow = user32.MustFindProc("GetTopWindow").Addr()
var _ = builtin(GetTopWindow, "(hWnd)")

func GetTopWindow(a Value) Value {
	rtn := goc.Syscall1(getTopWindow,
		intArg(a))
	return intRet(rtn)
}

// dll long User32:GetUpdateRgn(pointer hwnd, pointer hRgn, bool bErase)
var getUpdateRgn = user32.MustFindProc("GetUpdateRgn").Addr()
var _ = builtin(GetUpdateRgn, "(hwnd, hRgn, bErase)")

func GetUpdateRgn(a, b, c Value) Value {
	rtn := goc.Syscall3(getUpdateRgn,
		intArg(a),
		intArg(b),
		boolArg(c))
	return intRet(rtn)
}

// dll pointer User32:GetWindow(pointer hWnd, long uCmd)
var getWindow = user32.MustFindProc("GetWindow").Addr()
var _ = builtin(GetWindow, "(hWnd, uCmd)")

func GetWindow(a, b Value) Value {
	rtn := goc.Syscall2(getWindow,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll pointer User32:GetWindowDC(pointer hwnd)
var getWindowDC = user32.MustFindProc("GetWindowDC").Addr()
var _ = builtin(GetWindowDC, "(hwnd)")

func GetWindowDC(a Value) Value {
	rtn := goc.Syscall1(getWindowDC,
		intArg(a))
	return intRet(rtn)
}

// dll bool User32:IsClipboardFormatAvailable(long format)
var isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable").Addr()
var _ = builtin(IsClipboardFormatAvailable, "(format)")

func IsClipboardFormatAvailable(a Value) Value {
	rtn := goc.Syscall1(isClipboardFormatAvailable,
		intArg(a))
	return boolRet(rtn)
}

// dll bool User32:IsWindow(pointer hwnd)
var isWindow = user32.MustFindProc("IsWindow").Addr()
var _ = builtin(IsWindow, "(hwnd)")

func IsWindow(a Value) Value {
	rtn := goc.Syscall1(isWindow,
		intArg(a))
	return boolRet(rtn)
}

// dll bool User32:IsChild(pointer hwndParent, pointer hwnd)
var isChild = user32.MustFindProc("IsChild").Addr()
var _ = builtin(IsChild, "(hwndParent, hwnd)")

func IsChild(a, b Value) Value {
	rtn := goc.Syscall2(isChild,
		intArg(a),
		intArg(b))
	return boolRet(rtn)
}

// dll bool User32:IsWindowVisible(pointer hwnd)
var isWindowVisible = user32.MustFindProc("IsWindowVisible").Addr()
var _ = builtin(IsWindowVisible, "(hwnd)")

func IsWindowVisible(a Value) Value {
	rtn := goc.Syscall1(isWindowVisible,
		intArg(a))
	return boolRet(rtn)
}

// dll bool User32:MessageBeep(long type)
var messageBeep = user32.MustFindProc("MessageBeep").Addr()
var _ = builtin(MessageBeep, "(type)")

func MessageBeep(a Value) Value {
	rtn := goc.Syscall1(messageBeep,
		intArg(a))
	return boolRet(rtn)
}

// dll void User32:mouse_event(
// long	dwFlags,		// motion and click options
// long	dx,				// horizontal position or change
// long	dy,				// vertical position or change
// long	dwData,			// wheel movement
// pointer dwExtraInfo)	// (ULONG_PTR) application-defined information
var mouse_event = user32.MustFindProc("mouse_event").Addr()
var _ = builtin(Mouse_event, "(dwFlags, dx, dy, dwData, dwExtraInfo)")

func Mouse_event(a, b, c, d, e Value) Value {
	goc.Syscall5(mouse_event,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e))
	return nil
}

// dll bool User32:PostMessage(pointer hwnd, long msg, pointer wParam,
// pointer lParam)
var postMessage = user32.MustFindProc("PostMessageA").Addr()
var _ = builtin(PostMessage, "(hwnd, msg, wParam, lParam)")

func PostMessage(a, b, c, d Value) Value {
	rtn := goc.Syscall4(postMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d))
	return boolRet(rtn)
}

// dll bool User32:RegisterHotKey(
// pointer hWnd /*optional*/,
// long id,
// long fsModifiers,
// long vk)
var registerHotKey = user32.MustFindProc("RegisterHotKey").Addr()
var _ = builtin(RegisterHotKey, "(hWnd, id, fsModifiers, vk)")

func RegisterHotKey(a, b, c, d Value) Value {
	rtn := goc.Syscall4(registerHotKey,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d))
	return boolRet(rtn)
}

// dll bool User32:ReleaseCapture()
var releaseCapture = user32.MustFindProc("ReleaseCapture").Addr()
var _ = builtin(ReleaseCapture, "()")

func ReleaseCapture() Value {
	rtn := goc.Syscall0(releaseCapture)
	return boolRet(rtn)
}

// dll pointer User32:SetActiveWindow(pointer hWnd)
var setActiveWindow = user32.MustFindProc("SetActiveWindow").Addr()
var _ = builtin(SetActiveWindow, "(hWnd)")

func SetActiveWindow(a Value) Value {
	rtn := goc.Syscall1(setActiveWindow,
		intArg(a))
	return intRet(rtn)
}

// dll pointer User32:SetCapture(pointer hwnd)
var setCapture = user32.MustFindProc("SetCapture").Addr()
var _ = builtin(SetCapture, "(hwnd)")

func SetCapture(a Value) Value {
	rtn := goc.Syscall1(setCapture,
		intArg(a))
	return intRet(rtn)
}

// dll pointer User32:SetCursor(pointer hcursor)
var setCursor = user32.MustFindProc("SetCursor").Addr()
var _ = builtin(SetCursor, "(hcursor)")

func SetCursor(a Value) Value {
	rtn := goc.Syscall1(setCursor,
		intArg(a))
	return intRet(rtn)
}

// dll bool User32:SetForegroundWindow(pointer hwnd)
var setForegroundWindow = user32.MustFindProc("SetForegroundWindow").Addr()
var _ = builtin(SetForegroundWindow, "(hwnd)")

func SetForegroundWindow(a Value) Value {
	rtn := goc.Syscall1(setForegroundWindow,
		intArg(a))
	return boolRet(rtn)
}

// dll bool User32:SetMenuDefaultItem(
// pointer hMenu,
// long uItem,
// long fByPosition)
var setMenuDefaultItem = user32.MustFindProc("SetMenuDefaultItem").Addr()
var _ = builtin(SetMenuDefaultItem, "(hMenu, uItem, fByPosition)")

func SetMenuDefaultItem(a, b, c Value) Value {
	rtn := goc.Syscall3(setMenuDefaultItem,
		intArg(a),
		intArg(b),
		intArg(c))
	return boolRet(rtn)
}

// dll pointer User32:SetParent(pointer hwndNewChild, pointer hwndNewParent)
var setParent = user32.MustFindProc("SetParent").Addr()
var _ = builtin(SetParent, "(hwndNewChild, hwndNewParent)")

func SetParent(a, b Value) Value {
	rtn := goc.Syscall2(setParent,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll bool User32:SetProp(pointer hwnd, [in] string name, pointer value)
var setProp = user32.MustFindProc("SetPropA").Addr()
var _ = builtin(SetProp, "(hwnd, name, value)")

func SetProp(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall3(setProp,
		intArg(a),
		uintptr(stringArg(b)),
		intArg(c))
	return boolRet(rtn)
}

// dll bool User32:UnhookWindowsHookEx(pointer hhk)
var unhookWindowsHookEx = user32.MustFindProc("UnhookWindowsHookEx").Addr()
var _ = builtin(UnhookWindowsHookEx, "(hhk)")

func UnhookWindowsHookEx(a Value) Value {
	rtn := goc.Syscall1(unhookWindowsHookEx,
		intArg(a))
	return boolRet(rtn)
}

// dll bool User32:UnregisterHotKey(pointer hWnd /*optional*/, long id)
var unregisterHotKey = user32.MustFindProc("UnregisterHotKey").Addr()
var _ = builtin(UnregisterHotKey, "(hWnd, id)")

func UnregisterHotKey(a, b Value) Value {
	rtn := goc.Syscall2(unregisterHotKey,
		intArg(a),
		intArg(b))
	return boolRet(rtn)
}

// dll bool User32:ClientToScreen(pointer hwnd, POINT* point)
var clientToScreen = user32.MustFindProc("ClientToScreen").Addr()
var _ = builtin(ClientToScreen, "(hWnd, point)")

func ClientToScreen(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nPoint)
	rtn := goc.Syscall2(clientToScreen,
		intArg(a),
		uintptr(pointArg(b, p)))
	pt := (*stPoint)(p)
	b.Put(nil, SuStr("x"), IntVal(int(pt.x)))
	b.Put(nil, SuStr("y"), IntVal(int(pt.y)))
	return boolRet(rtn)
}

// dll bool User32:ClipCursor(RECT* rect)
var clipCursor = user32.MustFindProc("ClipCursor").Addr()
var _ = builtin(ClipCursor, "(rect)")

func ClipCursor(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall1(clipCursor,
		uintptr(rectArg(a, r)))
	return boolRet(rtn)
}

// dll pointer User32:DeferWindowPos(pointer hWinPosInfo, pointer hWnd,
// pointer hWndInsertAfter, long x, long y, long cx, long cy, long flags)
var deferWindowPos = user32.MustFindProc("DeferWindowPos").Addr()
var _ = builtin(DeferWindowPos,
	"(hWinPosInfo, hWnd, hWndInsertAfter, x, y, cx, cy, flags)")

func DeferWindowPos(_ *Thread, a []Value) Value {
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
}

// dll bool User32:DrawFocusRect(pointer hdc, RECT* lprc)
var drawFocusRect = user32.MustFindProc("DrawFocusRect").Addr()
var _ = builtin(DrawFocusRect, "(hwnd, rect)")

func DrawFocusRect(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall2(drawFocusRect,
		intArg(a),
		uintptr(rectArg(b, r)))
	return boolRet(rtn)
}

// dll long User32:DrawTextEx(pointer hdc, [in] string lpsz, long cb,
// RECT* lprc, long uFormat, DRAWTEXTPARAMS* params)
var drawTextEx = user32.MustFindProc("DrawTextExA").Addr()
var _ = builtin(DrawTextEx, "(hdc, lpsz, cb, lprc, uFormat, params)")

func DrawTextEx(a, b, c, d, e, f Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall6(drawTextEx,
		intArg(a),
		uintptr(stringArg(b)),
		intArg(c),
		uintptr(rectArg(d, r)),
		intArg(e),
		uintptr(drawTextParams(f)))
	urectToOb(r, d)
	return intRet(rtn)
}

var _ = builtin(DrawTextExOut, "(hdc, text, rect, flags, params)")

func DrawTextExOut(a, b, c, d, e Value) Value {
	defer heap.FreeTo(heap.CurSize())
	text := ToStr(b)
	bufsize := len(text) + 8
	buf := heap.Copy(text, bufsize)
	r := heap.Alloc(nRect)
	rtn := goc.Syscall6(drawTextEx,
		intArg(a),
		uintptr(buf),
		uintptrMinusOne,
		uintptr(rectArg(c, r)),
		intArg(d),
		uintptr(drawTextParams(e)))
	urectToOb(r, c)
	ob := &SuObject{}
	ob.Put(nil, SuStr("text"), SuStr(heap.GetStrZ(buf, bufsize)))
	ob.Put(nil, SuStr("result"), intRet(rtn))
	return ob
}

func drawTextParams(x Value) unsafe.Pointer {
	p := unsafe.Pointer(nil)
	if !x.Equal(Zero) {
		p = heap.Alloc(nDrawTextParams)
		*(*stDrawTextParams)(p) = stDrawTextParams{
			cbSize:        uint32(nDrawTextParams),
			iTabLength:    getInt32(x, "iTabLength"),
			iLeftMargin:   getInt32(x, "iLeftMargin"),
			iRightMargin:  getInt32(x, "iRightMargin"),
			uiLengthDrawn: getInt32(x, "uiLengthDrawn"),
		}
	}
	return p
}

type stDrawTextParams struct {
	cbSize        uint32
	iTabLength    int32
	iLeftMargin   int32
	iRightMargin  int32
	uiLengthDrawn int32
}

const nDrawTextParams = unsafe.Sizeof(stDrawTextParams{})

// dll bool User32:TrackMouseEvent(TRACKMOUSEEVENT* lpEventTrack)
var trackMouseEvent = user32.MustFindProc("TrackMouseEvent").Addr()
var _ = builtin(TrackMouseEvent, "(lpEventTrack)")

func TrackMouseEvent(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nTrackMouseEvent)
	*(*stTrackMouseEvent)(p) = stTrackMouseEvent{
		cbSize:      uint32(nTrackMouseEvent),
		dwFlags:     getInt32(a, "dwFlags"),
		hwndTrack:   getUintptr(a, "hwndTrack"),
		dwHoverTime: getInt32(a, "dwHoverTime"),
	}
	rtn := goc.Syscall1(trackMouseEvent,
		uintptr(p))
	return boolRet(rtn)
}

type stTrackMouseEvent struct {
	cbSize      uint32
	dwFlags     int32
	hwndTrack   uintptr
	dwHoverTime int32
	_           [4]byte // padding
}

const nTrackMouseEvent = unsafe.Sizeof(stTrackMouseEvent{})

// dll bool User32:FlashWindowEx(FLASHWINFO* fi)
var flashWindowEx = user32.MustFindProc("FlashWindowEx").Addr()
var _ = builtin(FlashWindowEx, "(fi)")

func FlashWindowEx(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nFlashWInfo)
	*(*stFlashWInfo)(p) = stFlashWInfo{
		cbSize:    uint32(nFlashWInfo),
		hwnd:      getUintptr(a, "hwnd"),
		dwFlags:   getInt32(a, "dwFlags"),
		uCount:    getInt32(a, "uCount"),
		dwTimeout: getInt32(a, "dwTimeout"),
	}
	rtn := goc.Syscall1(flashWindowEx,
		uintptr(p))
	return boolRet(rtn)
}

type stFlashWInfo struct {
	cbSize    uint32
	hwnd      HANDLE
	dwFlags   int32
	uCount    int32
	dwTimeout int32
	_         [4]byte // padding
}

const nFlashWInfo = unsafe.Sizeof(stFlashWInfo{})

// dll long User32:FrameRect(pointer hdc, RECT* rect, pointer brush)
var frameRect = user32.MustFindProc("FrameRect").Addr()
var _ = builtin(FrameRect, "(hdc, rect, brush)")

func FrameRect(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall3(frameRect,
		intArg(a),
		uintptr(rectArg(b, r)),
		intArg(c))
	return intRet(rtn)
}

// dll bool User32:GetClipCursor(RECT* rect)
var getClipCursor = user32.MustFindProc("GetClipCursor").Addr()
var _ = builtin(GetClipCursor, "(rect)")

func GetClipCursor(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall1(getClipCursor,
		uintptr(r))
	urectToOb(r, a)
	return boolRet(rtn)
}

// dll bool User32:GetCursorPos(POINT* p)
var getCursorPos = user32.MustFindProc("GetCursorPos").Addr()
var _ = builtin(GetCursorPos, "(rect)")

func GetCursorPos(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nPoint)
	rtn := goc.Syscall1(getCursorPos,
		uintptr(p))
	upointToOb(p, a)
	return boolRet(rtn)
}

// dll bool User32:EnumThreadWindows(long dwThreadId, WNDENUMPROC lpfn,
// pointer lParam)
var enumThreadWindows = user32.MustFindProc("EnumThreadWindows").Addr()
var _ = builtin(EnumThreadWindows, "(dwThreadId, lpfn, lParam)")

func EnumThreadWindows(a, b, c Value) Value {
	rtn := goc.Syscall3(enumThreadWindows,
		intArg(a),
		NewCallback(b, 2),
		intArg(c))
	return boolRet(rtn)
}

// dll bool User32:EnumChildWindows(pointer hwnd, WNDENUMPROC lpEnumProc,
// pointer lParam)
var enumChildWindows = user32.MustFindProc("EnumChildWindows").Addr()
var _ = builtin(EnumChildWindowsApi, "(hwnd, lpEnumProc, lParam)")

func EnumChildWindowsApi(a, b, c Value) Value {
	rtn := goc.Syscall3(enumChildWindows,
		intArg(a),
		NewCallback(b, 2),
		intArg(c))
	return boolRet(rtn)
}

// dll pointer User32:WindowFromPoint(POINT pt)
var windowFromPoint = user32.MustFindProc("WindowFromPoint").Addr()
var _ = builtin(WindowFromPoint, "(pt)")

func WindowFromPoint(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nPoint)
	rtn := goc.Syscall1(windowFromPoint,
		uintptr(pointArg(a, p)))
	return intRet(rtn)
}

// dll long User32:GetWindowThreadProcessId(pointer hwnd, LONG* lpdwProcessId)
var getWindowThreadProcessId = user32.MustFindProc("GetWindowThreadProcessId").Addr()
var _ = builtin(GetWindowThreadProcessId, "(hwnd, lpdwProcessId)")

func GetWindowThreadProcessId(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(4)
	rtn := goc.Syscall2(getWindowThreadProcessId,
		intArg(a),
		uintptr(p))
	b.Put(nil, SuStr("x"), IntVal(int(*(*int32)(p))))
	return boolRet(rtn)
}

// dll long User32:TrackPopupMenu(pointer hMenu, long uFlags, long x, long y,
// long nReserved, pointer hWnd, RECT* prcRect)
var trackPopupMenu = user32.MustFindProc("TrackPopupMenu").Addr()
var _ = builtin(TrackPopupMenu, "(hMenu, uFlags, x, y, nReserved, hWnd, prcRect)")

func TrackPopupMenu(a, b, c, d, e, f, g Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall7(trackPopupMenu,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e),
		intArg(f),
		uintptr(rectArg(g, r)))
	return intRet(rtn)
}

// dll pointer User32:SetWindowsHookEx(long idHook, HOOKPROC lpfn, pointer hMod,
// long dwThreadId)
var setWindowsHookEx = user32.MustFindProc("SetWindowsHookExA").Addr()
var _ = builtin(SetWindowsHookEx, "(idHook, lpfn, hMod, dwThreadId)")

func SetWindowsHookEx(a, b, c, d Value) Value {
	rtn := goc.Syscall4(setWindowsHookEx,
		intArg(a),
		NewCallback(b, 3),
		intArg(c),
		intArg(d))
	return intRet(rtn)
}

// dll long User32:SetScrollInfo(pointer hwnd, long bar, SCROLLINFO* si,
// bool redraw)
var setScrollInfo = user32.MustFindProc("SetScrollInfo").Addr()
var _ = builtin(SetScrollInfo, "(hwnd, bar, si, redraw)")

func SetScrollInfo(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nScrollInfo)
	*(*stScrollInfo)(p) = stScrollInfo{
		cbSize:    uint32(nScrollInfo),
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
}

// dll long User32:ScrollWindowEx(pointer hwnd, long dx, long dy, RECT* scroll,
// RECT* clip, pointer rgnUpdate, RECT* rcUpdate, long flags)
var scrollWindowEx = user32.MustFindProc("ScrollWindowEx").Addr()
var _ = builtin(ScrollWindowEx,
	"(hwnd, dx, dy, scroll, clip, rgnUpdate, rcUpdate, flags)")

func ScrollWindowEx(_ *Thread, a []Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r1 := heap.Alloc(nRect)
	r2 := heap.Alloc(nRect)
	r3 := heap.Alloc(nRect)
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
}

// dll bool User32:ScreenToClient(pointer hwnd, POINT* p)
var screenToClient = user32.MustFindProc("ScreenToClient").Addr()
var _ = builtin(ScreenToClient, "(hWnd, p)")

func ScreenToClient(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nPoint)
	rtn := goc.Syscall2(screenToClient,
		intArg(a),
		uintptr(pointArg(b, p)))
	upointToOb(p, b)
	return boolRet(rtn)
}

// dll pointer User32:LoadImage(pointer hInstance, resource lpszName,
// long uType, long cxDesired, long cyDesired, long fuLoad)
var loadImage = user32.MustFindProc("LoadImageA").Addr()
var _ = builtin(LoadImage,
	"(hInstance, lpszName, uType, cxDesired, cyDesired, fuLoad)")

func LoadImage(a, b, c, d, e, f Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall6(loadImage,
		intArg(a),
		uintptr(stringArg(b)), // doesn't handle resource id
		intArg(c),
		intArg(d),
		intArg(e),
		intArg(f))
	return intRet(rtn)
}

// dll long User32:MapWindowPoints(pointer hwndfrom, pointer hwndto, RECT* p,
// long n)
var mapWindowPoints = user32.MustFindProc("MapWindowPoints").Addr()
var _ = builtin(MapWindowRect, "(hwndfrom, hwndto, r)")

func MapWindowRect(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall4(mapWindowPoints,
		intArg(a),
		intArg(b),
		uintptr(rectArg(c, r)),
		2)
	urectToOb(r, c)
	return intRet(rtn)
}

func init() {
	SuneidoObjectMethods["CNullPointer"] =
		builtinVal("CNullPointer", CNullPointer, "()")
}

func CNullPointer() Value {
	goc.Syscall1(0, 0)
	return nil
}

var getGuiResources = user32.MustFindProc("GetGuiResources").Addr()

const GR_GDIOBJECTS = 0
const GR_USEROBJECTS = 1

func GetGuiResources() (int, int) {
	hProcess, _ := syscall.GetCurrentProcess()
	gdi, _, _ := syscall.SyscallN(getGuiResources,
		uintptr(hProcess),
		uintptr(GR_GDIOBJECTS))
	user, _, _ := syscall.SyscallN(getGuiResources,
		uintptr(hProcess),
		uintptr(GR_USEROBJECTS))
	return int(gdi), int(user)
}

var _ = AddInfo("windows.nGdiObject", func() int {
	hProcess, _ := syscall.GetCurrentProcess()
	n, _, _ := syscall.SyscallN(getGuiResources,
		uintptr(hProcess),
		uintptr(GR_GDIOBJECTS))
	return int(n)
})

var _ = AddInfo("windows.nUserObject", func() int {
	hProcess, _ := syscall.GetCurrentProcess()
	n, _, _ := syscall.SyscallN(getGuiResources,
		uintptr(hProcess),
		uintptr(GR_USEROBJECTS))
	return int(n)
})
