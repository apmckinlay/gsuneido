package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
	"golang.org/x/sys/windows"
)

func init() {
	windows.MustLoadDLL("scilexer.dll")
}

var user32 = windows.NewLazyDLL("user32.dll")

type HANDLE = uintptr

type RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

type PAINTSTRUCT struct {
	hdc         HANDLE
	fErase      bool
	rcPaint     RECT
	fRestore    bool
	fIncUpdate  bool
	rgbReserved [32]byte
}

type MONITORINFO struct {
	cbSize    uint32
	rcMonitor RECT
	rcWork    RECT
	dwFlags   uint32
}

type SCROLLINFO struct {
	cbSize    uint32
	fMask     uint32
	nMin      int32
	nMax      int32
	nPage     uint32
	nPos      int32
	nTrackPos int32
}

type POINT struct {
	x int32
	y int32
}

type WINDOWPLACEMENT struct {
	length           uint32
	flags            uint32
	showCmd          uint32
	ptMinPosition    POINT
	ptMaxPosition    POINT
	rcNormalPosition RECT
	rcDevice         RECT // stdlib does not have this member
}

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

type TCITEM struct {
	mask        uint32
	dwState     uint32
	dwStateMask uint32
	pszText     *byte
	cchTextMax  int32
	iImage      int32
	lParam      int32
}

type CHARRANGE struct {
	cpMin int32
	cpMax int32
}

type TEXTRANGE struct {
	chrg      CHARRANGE
	lpstrText *byte
}

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
	lParam         int32
}

type TV_INSERTSTRUCT struct {
	hParent      HANDLE
	hInsertAfter HANDLE
	item         TVITEM
}

// helper functions

func boolArg(arg Value) uintptr {
	if ToBool(arg) {
		return 1
	}
	return 0
}

func boolRet(rtn uintptr) Value {
	if rtn == 0 {
		return False
	}
	return True
}

func intArg(arg Value) uintptr {
	if arg.Equal(True) {
		return 1
	}
	if arg.Equal(False) {
		return 0
	}
	return uintptr(ToInt(arg))
}

func intRet(rtn uintptr) Value {
	return IntVal(int(rtn))
}

var zero byte

// bytePtrFromString returns a pointer to a nul terminated copy of the string.
// The copy is required to get nul termination.
// uintptr(unsafe.Pointer( must be in final call to prevent garbage collection.
func bytePtrFromString(v Value) *byte {
	if v.Equal(Zero) {
		return nil
	}
	s := ToStr(v)
	if s == "" {
		return nil
	}
	a := make([]byte, len(s)+1) // +1 for nul terminator
	copy(a, s)
	return &a[0]
}

func stringArg(v Value) unsafe.Pointer {
	return unsafe.Pointer(bytePtrFromString(v))
}

func getBool(ob Value, mem string) bool {
	x := ob.Get(nil, SuStr(mem))
	if x == nil {
		return false
	}
	return ToBool(x)
}

func getInt(ob Value, mem string) int {
	x := ob.Get(nil, SuStr(mem))
	if x == nil || x.Equal(False) {
		return 0
	}
	if x.Equal(True) {
		return 1
	}
	return ToInt(x)
}

func getInt32(ob Value, mem string) int32 {
	return int32(getInt(ob, mem))
}

func getUint32(ob Value, mem string) uint32 {
	return uint32(getInt(ob, mem))
}

func getStr(ob Value, mem string) *byte {
	x := ob.Get(nil, SuStr(mem))
	if x == nil || x.Equal(Zero) || x.Equal(False) {
		return nil
	}
	return bytePtrFromString(x)
}

func rectArg(ob Value, r *RECT) unsafe.Pointer {
	if ob.Equal(Zero) {
		return nil
	}
	*r = obToRect(ob)
	return unsafe.Pointer(r)
}

func obToRect(ob Value) RECT {
	return RECT{
		left:   getInt32(ob, "left"),
		top:    getInt32(ob, "top"),
		right:  getInt32(ob, "right"),
		bottom: getInt32(ob, "bottom"),
	}
}

func getRect(ob Value, mem string) RECT {
	x := ob.Get(nil, SuStr(mem))
	if x == nil {
		return RECT{}
	}
	return obToRect(x)
}

func rectToOb(r *RECT, ob Value) Value {
	if ob == nil {
		ob = NewSuObject()
	} else if ob.Equal(Zero) {
		return ob
	}
	ob.Put(nil, SuStr("left"), IntVal(int(r.left)))
	ob.Put(nil, SuStr("top"), IntVal(int(r.top)))
	ob.Put(nil, SuStr("right"), IntVal(int(r.right)))
	ob.Put(nil, SuStr("bottom"), IntVal(int(r.bottom)))
	return ob
}

func obToPoint(ob Value) POINT {
	return POINT{
		x: getInt32(ob, "x"),
		y: getInt32(ob, "y"),
	}
}

func pointToOb(pt *POINT, ob Value) Value {
	if ob == nil {
		ob = NewSuObject()
	}
	ob.Put(nil, SuStr("x"), IntVal(int(pt.x)))
	ob.Put(nil, SuStr("y"), IntVal(int(pt.y)))
	return ob
}

func getPoint(ob Value, mem string) POINT {
	x := ob.Get(nil, SuStr(mem))
	if x == nil {
		return POINT{}
	}
	return obToPoint(x)
}

func pointArg(ob Value, pt *POINT) unsafe.Pointer {
	pt.x = getInt32(ob, "x")
	pt.y = getInt32(ob, "y")
	return unsafe.Pointer(&pt)
}

func getHandle(ob Value, mem string) HANDLE {
	return HANDLE(getInt(ob, mem))
}

//===================================================================

// dll User32:GetDesktopWindow() hwnd
var getDesktopWindow = user32.NewProc("GetDesktopWindow")
var _ = builtin0("GetDesktopWindow()",
	func() Value {
		n, _, _ := getDesktopWindow.Call()
		return IntVal(int(n))
	})

// dll User32:GetSysColor(long nIndex) long
var getSysColor = user32.NewProc("GetSysColor")
var _ = builtin1("GetSysColor(index)",
	func(a Value) Value {
		n, _, _ := getSysColor.Call(
			intArg(a))
		return IntVal(int(n))
	})

// dll User32:GetWindowRect(pointer hwnd, RECT* rect) bool
var getWindowRect = user32.NewProc("GetWindowRect")
var _ = builtin2("GetWindowRectApi(hwnd, rect)",
	func(a Value, b Value) Value {
		var r RECT
		rtn, _, _ := getWindowRect.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&r)))
		rectToOb(&r, b)
		return boolRet(rtn)
	})

// dll long User32:MessageBox(pointer window, [in] string text,
//		[in] string caption, long flags)
var messageBox = user32.NewProc("MessageBoxA")
var _ = builtin4("MessageBox(hwnd, text, caption, flags)",
	func(a, b, c, d Value) Value {
		n, _, _ := messageBox.Call(
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(stringArg(c)),
			intArg(d))
		return IntVal(int(n))
	})

// dll User32:AdjustWindowRectEx(RECT* rect, long style, bool menu,
// 		long exStyle) bool
var adjustWindowRectEx = user32.NewProc("AdjustWindowRectEx")
var _ = builtin4("AdjustWindowRectEx(lpRect, dwStyle, bMenu, dwExStyle)",
	func(a, b, c, d Value) Value {
		var r RECT
		rtn, _, _ := adjustWindowRectEx.Call(
			uintptr(rectArg(a, &r)),
			intArg(b),
			boolArg(c),
			intArg(d))
		rectToOb(&r, a)
		return boolRet(rtn)
	})

// dll User32:CreateMenu() pointer
var createMenu = user32.NewProc("CreateMenu")
var _ = builtin0("CreateMenu()",
	func() Value {
		rtn, _, _ := createMenu.Call()
		return intRet(rtn)
	})

// dll User32:CreatePopupMenu() pointer
var createPopupMenu = user32.NewProc("CreatePopupMenu")
var _ = builtin0("CreatePopupMenu()",
	func() Value {
		rtn, _, _ := createPopupMenu.Call()
		return intRet(rtn)
	})

// dll User32:AppendMenu(pointer hmenu, long flags, pointer item,
//		[in] string name) bool
var appendMenu = user32.NewProc("AppendMenuA")
var _ = builtin4("AppendMenu(hmenu, flags, item, name)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := appendMenu.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(stringArg(d)))
		return boolRet(rtn)
	})

// dll User32:DestroyMenu(pointer hmenu) bool
var destroyMenu = user32.NewProc("DestroyMenu")
var _ = builtin1("DestroyMenu(hmenu)",
	func(a Value) Value {
		rtn, _, _ := destroyMenu.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:CreateWindowEx(long exStyle, resource classname, [in] string name,
//		long style, long x, long y, long w, long h, pointer parent, pointer menu,
//		pointer instance, pointer param) pointer
var createWindowEx = user32.NewProc("CreateWindowExA")
var _ = builtin("CreateWindowEx(exStyle, classname, name, style, x, y, w, h,"+
	" parent, menu, instance, param)",
	func(_ *Thread, a []Value) Value {
		rtn, _, _ := createWindowEx.Call(
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
var getSystemMenu = user32.NewProc("GetSystemMenu")
var _ = builtin2("GetSystemMenu(hwnd, bRevert)",
	func(a, b Value) Value {
		rtn, _, _ := getSystemMenu.Call(
			intArg(a),
			boolArg(b))
		return intRet(rtn)
	})

// dll User32:SetMenu(pointer hwnd, pointer hmenu) bool
var setMenu = user32.NewProc("SetMenu")
var _ = builtin2("SetMenu(hwnd, hmenu)",
	func(a, b Value) Value {
		rtn, _, _ := setMenu.Call(
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll User32:BeginPaint(pointer hwnd, PAINTSTRUCT* ps) pointer
var beginPaint = user32.NewProc("BeginPaint")
var _ = builtin2("BeginPaint(hwnd, ps)",
	func(a, b Value) Value {
		var ps PAINTSTRUCT
		rtn, _, _ := beginPaint.Call(
			intArg(a),
			uintptr(psArg(b, &ps)))
		b.Put(nil, SuStr("hdc"), IntVal(int(ps.hdc)))
		b.Put(nil, SuStr("fErase"), SuBool(ps.fErase))
		b.Put(nil, SuStr("rcPaint"),
			rectToOb(&ps.rcPaint, b.Get(nil, SuStr("rcPaint"))))
		b.Put(nil, SuStr("fRestore"), SuBool(ps.fRestore))
		b.Put(nil, SuStr("fIncUpdate"), SuBool(ps.fIncUpdate))
		return intRet(rtn)
	})

// dll User32:EndPaint(pointer hwnd, PAINTSTRUCT* ps) bool
var endPaint = user32.NewProc("EndPaint")
var _ = builtin2("EndPaint(hwnd, ps)",
	func(a, b Value) Value {
		var ps PAINTSTRUCT
		rtn, _, _ := endPaint.Call(
			intArg(a),
			uintptr(psArg(b, &ps)))
		return boolRet(rtn)
	})

func psArg(ob Value, ps *PAINTSTRUCT) unsafe.Pointer {
	ps.hdc = getHandle(ob, "hdc")
	ps.fErase = getBool(ob, "fErase")
	ps.rcPaint = getRect(ob, "rcPaint")
	ps.fRestore = getBool(ob, "fRestore")
	ps.fIncUpdate = getBool(ob, "fIncUpdate")
	return unsafe.Pointer(&ps)
}

// dll User32:CallWindowProc(pointer wndprcPrev, pointer hwnd, long msg,
//		pointer wParam, pointer lParam) pointer
var callWindowProc = user32.NewProc("CallWindowProcA")
var _ = builtin5("CallWindowProc(wndprcPrev, hwnd, msg, wParam, lParam)",
	func(a, b, c, d, e Value) Value {
		rtn, _, _ := callWindowProc.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return intRet(rtn)
	})

// dll User32:CreateAcceleratorTable([in] string lpaccel, long cEntries) pointer
var createAcceleratorTable = user32.NewProc("CreateAcceleratorTable")
var _ = builtin2("CreateAcceleratorTable(lpaccel, cEntries)",
	func(a, b Value) Value {
		rtn, _, _ := createAcceleratorTable.Call(
			uintptr(stringArg(a)),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:DestroyAcceleratorTable(pointer hAccel) bool
var destroyAcceleratorTable = user32.NewProc("DestroyAcceleratorTable")
var _ = builtin1("DestroyAcceleratorTable(hAccel)",
	func(a Value) Value {
		rtn, _, _ := destroyAcceleratorTable.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:DestroyWindow(pointer hwnd) bool
var destroyWindow = user32.NewProc("DestroyWindow")
var _ = builtin1("DestroyWindow(hwnd)",
	func(a Value) Value {
		rtn, _, _ := destroyWindow.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:DrawFrameControl(pointer hdc, RECT* lprc, long uType,
//		long uState) bool
var drawFrameControl = user32.NewProc("DrawFrameControl")
var _ = builtin4("DrawFrameControl(hdc, lprc, uType, uState)",
	func(a, b, c, d Value) Value {
		var r RECT
		rtn, _, _ := drawFrameControl.Call(
			intArg(a),
			uintptr(rectArg(b, &r)),
			intArg(c),
			intArg(d))
		return boolRet(rtn)
	})

// dll User32:DrawText(pointer hdc, [in] string lpsz, long cb, RECT* lprc,
//		long uFormat) long
var drawText = user32.NewProc("DrawText")
var _ = builtin5("DrawText(hdc, lpsz, cb, lprc, uFormat)",
	func(a, b, c, d, e Value) Value {
		var r RECT
		rtn, _, _ := drawText.Call(
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c),
			uintptr(rectArg(d, &r)),
			intArg(e))
		rectToOb(&r, d) // for CALCRECT
		return intRet(rtn)
	})

// dll User32:FillRect(pointer hdc, RECT* lpRect, pointer hBrush) long
var fillRect = user32.NewProc("FillRect")
var _ = builtin3("FillRect(hdc, lpRect, hBrush)",
	func(a, b, c Value) Value {
		var r RECT
		rtn, _, _ := fillRect.Call(
			intArg(a),
			uintptr(rectArg(b, &r)),
			intArg(c))
		return intRet(rtn)
	})

// dll User32:GetActiveWindow() pointer
var getActiveWindow = user32.NewProc("GetActiveWindow")
var _ = builtin0("GetActiveWindow()",
	func() Value {
		rtn, _, _ := getActiveWindow.Call()
		return intRet(rtn)
	})

// dll User32:GetFocus() pointer
var getFocus = user32.NewProc("GetFocus")
var _ = builtin0("GetFocus()",
	func() Value {
		rtn, _, _ := getFocus.Call()
		return intRet(rtn)
	})

// dll User32:GetClientRect(pointer hwnd, RECT* rect) bool
var getClientRect = user32.NewProc("GetClientRect")
var _ = builtin2("GetClientRect(hwnd, rect)",
	func(a, b Value) Value {
		var r RECT
		rtn, _, _ := getClientRect.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&r)))
		rectToOb(&r, b)
		return boolRet(rtn)
	})

// dll User32:GetDC(pointer hwnd) pointer
var getDC = user32.NewProc("GetDC")
var _ = builtin1("GetDC(hwnd)",
	func(a Value) Value {
		rtn, _, _ := getDC.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetMonitorInfo(pointer hMonitor, MONITORINFO* lpmi) bool
var getMonitorInfo = user32.NewProc("GetMonitorInfoA")
var _ = builtin2("GetMonitorInfoApi(hwnd, mInfo)",
	func(a, b Value) Value {
		var mi MONITORINFO
		mi.cbSize = uint32(unsafe.Sizeof(mi))
		rtn, _, _ := getMonitorInfo.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&mi)))
		b.Put(nil, SuStr("rcMonitor"), rectToOb(&mi.rcMonitor, nil))
		b.Put(nil, SuStr("rcWork"), rectToOb(&mi.rcWork, nil))
		b.Put(nil, SuStr("dwFlags"), IntVal(int(mi.dwFlags)))
		return boolRet(rtn)
	})

// dll User32:GetScrollInfo(pointer hwnd, long fnBar, SCROLLINFO* lpsi) bool
var getScrollInfo = user32.NewProc("GetScrollInfo")
var _ = builtin3("GetScrollInfo(hwnd, fnBar, lpsi)",
	func(a, b, c Value) Value {
		si := SCROLLINFO{
			cbSize:    uint32(unsafe.Sizeof(SCROLLINFO{})),
			fMask:     getUint32(c, "fMask"),
			nMin:      getInt32(c, "nMin"),
			nMax:      getInt32(c, "nMax"),
			nPage:     getUint32(c, "nPage"),
			nPos:      getInt32(c, "nPos"),
			nTrackPos: getInt32(c, "nTrackPos"),
		}
		rtn, _, _ := getScrollInfo.Call(
			intArg(a),
			intArg(b),
			uintptr(unsafe.Pointer(&si)))
		c.Put(nil, SuStr("nMin"), IntVal(int(si.nMin)))
		c.Put(nil, SuStr("nMax"), IntVal(int(si.nMax)))
		c.Put(nil, SuStr("nPage"), IntVal(int(si.nPage)))
		c.Put(nil, SuStr("nPos"), IntVal(int(si.nPos)))
		c.Put(nil, SuStr("nTrackPos"), IntVal(int(si.nTrackPos)))
		return boolRet(rtn)
	})

// dll User32:GetScrollPos(pointer hwnd, int nBar) int
var getScrollPos = user32.NewProc("GetScrollPos")
var _ = builtin2("GetScrollPos(hwnd, nBar)",
	func(a, b Value) Value {
		rtn, _, _ := getScrollPos.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetSysColorBrush(long nIndex) pointer
var getSysColorBrush = user32.NewProc("GetSysColorBrush")
var _ = builtin1("GetSysColorBrush(nIndex)",
	func(a Value) Value {
		rtn, _, _ := getSysColorBrush.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetSystemMetrics(long nIndex) long
var getSystemMetrics = user32.NewProc("GetSystemMetrics")
var _ = builtin1("GetSystemMetrics(nIndex)",
	func(a Value) Value {
		rtn, _, _ := getSystemMetrics.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetWindowLong(pointer hwnd, long offset) long
var getWindowLong = user32.NewProc("GetWindowLongA")
var _ = builtin2("GetWindowLong(hwnd, offset)",
	func(a, b Value) Value {
		rtn, _, _ := getWindowLong.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetWindowLong(pointer hwnd, long offset) long
var getWindowLongPtr = user32.NewProc("GetWindowLongPtrA")
var _ = builtin2("GetWindowLongPtr(hwnd, offset)",
	func(a, b Value) Value {
		rtn, _, _ := getWindowLongPtr.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetWindowPlacement(pointer hwnd, WINDOWPLACEMENT* lpwndpl) bool
var getWindowPlacement = user32.NewProc("GetWindowPlacement")
var _ = builtin2("GetWindowPlacement(hwnd, ps)",
	func(a, b Value) Value {
		var wp WINDOWPLACEMENT
		wp.length = uint32(unsafe.Sizeof(wp))
		rtn, _, _ := getWindowPlacement.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&wp)))
		b.Put(nil, SuStr("flags"), IntVal(int(wp.flags)))
		b.Put(nil, SuStr("showCmd"), IntVal(int(wp.showCmd)))
		b.Put(nil, SuStr("ptMinPosition"), pointToOb(&wp.ptMinPosition, nil))
		b.Put(nil, SuStr("ptMaxPosition"), pointToOb(&wp.ptMaxPosition, nil))
		b.Put(nil, SuStr("rcNormalPosition"),
			rectToOb(&wp.rcNormalPosition, nil))
		b.Put(nil, SuStr("rcDevice"), rectToOb(&wp.rcDevice, nil))
		return boolRet(rtn)
	})

// dll User32:GetWindowText(pointer hwnd, string buf, long len) long
var getWindowText = user32.NewProc("GetWindowTextA")
var getWindowTextLength = user32.NewProc("GetWindowTextLengthA")
var _ = builtin1("GetWindowText(hwnd)",
	func(hwnd Value) Value {
		n, _, _ := getWindowTextLength.Call(intArg(hwnd))
		buf := make([]byte, n+1)
		n, _, _ = getWindowText.Call(
			intArg(hwnd),
			uintptr(unsafe.Pointer(&buf[0])),
			n)
		return SuStr(string(buf[:n]))
	})

// dll User32:InflateRect(RECT* rect, long dx, long dy) bool
var inflateRect = user32.NewProc("InflateRect")
var _ = builtin3("InflateRect(rect, dx, dy)",
	func(a, b, c Value) Value {
		var r RECT
		rtn, _, _ := inflateRect.Call(
			uintptr(rectArg(a, &r)),
			intArg(b),
			intArg(c))
		rectToOb(&r, a)
		return boolRet(rtn)
	})

// dll User32:InsertMenuItem(pointer hMenu, long uItem, bool fByPosition,
//		MENUITEMINFO* lpmii) bool
var insertMenuItem = user32.NewProc("InsertMenuItemA")
var _ = builtin4("InsertMenuItem(hMenu, uItem, fByPosition, lpmii)",
	func(a, b, c, d Value) Value {
		m := MENUITEMINFO{
			cbSize:        uint32(unsafe.Sizeof(MENUITEMINFO{})),
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
		rtn, _, _ := insertMenuItem.Call(
			intArg(a),
			intArg(b),
			boolArg(c),
			uintptr(unsafe.Pointer(&m)))
		return boolRet(rtn)
	})

// dll long User32:GetMenuItemCount(pointer hMenu)
var getMenuItemCount = user32.NewProc("GetMenuItemCount")
var _ = builtin1("GetMenuItemCount(hMenu)",
	func(a Value) Value {
		rtn, _, _ := getMenuItemCount.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll long User32:GetMenuItemID(pointer hMenu, long nPos)
var getMenuItemID = user32.NewProc("GetMenuItemID")
var _ = builtin2("GetMenuItemID(hMenu, nPos)",
	func(a, b Value) Value {
		rtn, _, _ := getMenuItemID.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:GetMenuItemInfo(pointer hMenu, long uItem, bool fByPosition,
//		MENUITEMINFO* lpmii) bool
var getMenuItemInfo = user32.NewProc("GetMenuItemInfoA")
var _ = builtin4("GetMenuItemInfo(hMenu, uItem, fByPosition, lpmii)",
	func(a, b, c, d Value) Value {
		mii := MENUITEMINFO{
			cbSize: uint32(unsafe.Sizeof(MENUITEMINFO{})),
			fMask:  getUint32(d, "fMask"),
		}
		rtn, _, _ := getMenuItemInfo.Call(
			intArg(a),
			intArg(b),
			boolArg(c),
			uintptr(unsafe.Pointer(&mii)))
		b.Put(nil, SuStr("fMask"), IntVal(int(mii.fMask)))
		b.Put(nil, SuStr("fType"), IntVal(int(mii.fType)))
		b.Put(nil, SuStr("fState"), IntVal(int(mii.fState)))
		b.Put(nil, SuStr("wID"), IntVal(int(mii.wID)))
		b.Put(nil, SuStr("hSubMenu"), IntVal(int(mii.hSubMenu)))
		b.Put(nil, SuStr("hbmpChecked"), IntVal(int(mii.hbmpChecked)))
		b.Put(nil, SuStr("hbmpUnchecked"), IntVal(int(mii.hbmpUnchecked)))
		b.Put(nil, SuStr("dwItemData"), IntVal(int(mii.dwItemData)))
		//b.Put(nil, SuStr("dwTypeData"), IntVal(int(mii.dwTypeData)))
		b.Put(nil, SuStr("cch"), IntVal(int(mii.cch)))
		b.Put(nil, SuStr("hbmpItem"), IntVal(int(mii.hbmpItem)))
		return boolRet(rtn)
	})

var _ = builtin2("GetMenuItemInfoText(hMenu, uItem)",
	func(a, b Value) Value {
		const MMIM_TYPE = 0x10
		const MFT_STRING = 0
		mii := MENUITEMINFO{
			cbSize:     uint32(unsafe.Sizeof(MENUITEMINFO{})),
			fMask:      MMIM_TYPE,
			fType:      MFT_STRING,
			dwTypeData: nil,
		}
		rtn, _, _ := getMenuItemInfo.Call(
			intArg(a),
			intArg(b),
			0,
			uintptr(unsafe.Pointer(&mii)))
		if rtn == 0 {
			return False
		}
		mii.cch++
		buf := make([]byte, mii.cch)
		mii.dwTypeData = (*byte)(unsafe.Pointer(&buf[0]))
		rtn, _, _ = getMenuItemInfo.Call(
			intArg(a),
			intArg(b),
			0,
			uintptr(unsafe.Pointer(&mii)))
		return SuStr(string(buf[:]))
	})

// dll User32:SetMenuItemInfo(pointer hMenu, long uItem, long fByPosition,
//		MENUITEMINFO* lpmii) bool
var setMenuItemInfo = user32.NewProc("SetMenuItemInfoA")
var _ = builtin4("SetMenuItemInfo(hMenu, uItem, fByPosition, lpmii)",
	func(a, b, c, d Value) Value {
		m := MENUITEMINFO{
			cbSize:        uint32(unsafe.Sizeof(MENUITEMINFO{})),
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
		rtn, _, _ := setMenuItemInfo.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&m)))
		return boolRet(rtn)
	})

// dll User32:InvalidateRect(pointer hwnd, RECT* rect, bool erase) bool
var invalidateRect = user32.NewProc("InvalidateRect")
var _ = builtin3("InvalidateRect(hwnd, rect, erase)",
	func(a, b, c Value) Value {
		var r RECT
		rtn, _, _ := invalidateRect.Call(
			intArg(a),
			uintptr(rectArg(b, &r)),
			boolArg(c))
		return boolRet(rtn)
	})

// dll User32:IsWindowEnabled(pointer hwnd) bool
var isWindowEnabled = user32.NewProc("IsWindowEnabled")
var _ = builtin1("IsWindowEnabled(hwnd)",
	func(a Value) Value {
		rtn, _, _ := isWindowEnabled.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll User32:KillTimer(pointer hwnd, long id) bool
var killTimer = user32.NewProc("KillTimer")
var _ = builtin2("KillTimer(hwnd, id)",
	func(a, b Value) Value {
		rtn, _, _ := killTimer.Call(
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll User32:LoadCursor(pointer hinst, resource pszCursor) pointer
var loadCursor = user32.NewProc("LoadCursorA")
var _ = builtin2("LoadCursor(hinst, pszCursor)",
	func(a, b Value) Value {
		rtn, _, _ := loadCursor.Call(
			intArg(a),
			intArg(b)) // could be a string but we never use that
		return intRet(rtn)
	})

// dll User32:LoadIcon(pointer hInstance, resource lpIconName) pointer
var loadIcon = user32.NewProc("LoadIconA")
var _ = builtin2("LoadIcon(hinst, lpIconName)",
	func(a, b Value) Value {
		rtn, _, _ := loadIcon.Call(
			intArg(a),
			intArg(b)) // could be a string but we never use that
		return intRet(rtn)
	})

// dll User32:MonitorFromRect(RECT* lprc, long dwFlags) pointer
var monitorFromRect = user32.NewProc("MonitorFromRect")
var _ = builtin2("MonitorFromRect(lprc, dwFlags)",
	func(a, b Value) Value {
		var r RECT
		rtn, _, _ := monitorFromRect.Call(
			uintptr(rectArg(a, &r)),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:MoveWindow(pointer hwnd, long left, long top, long width,
//		long height, bool repaint) bool
var moveWindow = user32.NewProc("MoveWindow")
var _ = builtin6("MoveWindow(hwnd, left, top, width, height, repaint)",
	func(a, b, c, d, e, f Value) Value {
		rtn, _, _ := moveWindow.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			boolArg(f))
		return boolRet(rtn)
	})

// dll User32:RegisterClass(WNDCLASS* wc) short
var registerClass = user32.NewProc("RegisterClassA")
var _ = builtin1("RegisterClass(wc)",
	func(a Value) Value {
		w := WNDCLASS{
			style:      getUint32(a, "style"),
			wndProc:    uintptr(getInt(a, "wndProc")),
			clsExtra:   getInt32(a, "clsExtra"),
			wndExtra:   getInt32(a, "wndExtra"),
			instance:   getHandle(a, "instance"),
			icon:       getHandle(a, "icon"),
			cursor:     getHandle(a, "cursor"),
			background: getHandle(a, "background"),
			menuName:   getStr(a, "menuName"),
			className:  getStr(a, "className"),
		}
		rtn, _, _ := registerClass.Call(
			uintptr(unsafe.Pointer(&w)))
		return intRet(rtn)
	})

// dll User32:RegisterClipboardFormat([in] string lpszFormat) long
var registerClipboardFormat = user32.NewProc("RegisterClipboardFormat")
var _ = builtin1("RegisterClipboardFormat(lpszFormat)",
	func(a Value) Value {
		rtn, _, _ := registerClipboardFormat.Call(
			uintptr(stringArg(a)))
		return intRet(rtn)
	})

// dll User32:ReleaseDC(pointer hWnd, pointer hDC) long
var releaseDC = user32.NewProc("ReleaseDC")
var _ = builtin2("ReleaseDC(hwnd, hDC)",
	func(a, b Value) Value {
		rtn, _, _ := releaseDC.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		pointer lParam) pointer
var sendMessage = user32.NewProc("SendMessageA")
var _ = builtin4("SendMessage(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return intRet(rtn)
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		string text) pointer
var _ = builtin4("SendMessageText(hwnd, msg, wParam, text)",
	func(a, b, c, d Value) Value {
		// Must pass a defensive mutable copy of the string
		// (even though we discard it)
		// since the function may modify it.
		// Use SendMessageTextIn if the function doesn't modify it.
		// Use SendMessageTextOut if the modified text is needed.
		buf := ([]byte)(ToStr(d))
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&buf[0])))
		return intRet(rtn)
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		[in] string text) pointer
var _ = builtin4("SendMessageTextIn(hwnd, msg, wParam, text)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(stringArg(d)))
		return intRet(rtn)
	})

var _ = builtin4("SendMessageTextOut(hwnd, msg, wParam = 0, bufsize = 1024)",
	func(a, b, c, d Value) Value {
		n := ToInt(d)
		buf := make([]byte, n)
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&buf[0])))
		ob := NewSuObject()
		text := str.BeforeFirst(string(buf), "\x00")
		ob.Put(nil, SuStr("text"), SuStr(text))
		ob.Put(nil, SuStr("result"), IntVal(int(rtn)))
		return ob
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		TCITEM* tcitem) pointer
var _ = builtin4("SendMessageTcitem(hwnd, msg, wParam, tcitem)",
	func(a, b, c, d Value) Value {
		verify.That(getInt32(d, "cchTextMax") == 0)
		t := TCITEM{
			mask:        getUint32(d, "mask"),
			dwState:     getUint32(d, "dwState"),
			dwStateMask: getUint32(d, "dwStateMask"),
			pszText:     getStr(d, "pszText"),
			iImage:      getInt32(d, "iImage"),
			lParam:      getInt32(d, "lParam"),
		}
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&t)))
		return intRet(rtn)
	})

var _ = builtin5("SendMessageTextRange(hwnd, msg, cpMin, cpMax, each = 1)",
	func(a, b, c, d, e Value) Value {
		cpMin := ToInt(c)
		cpMax := ToInt(d)
		if cpMax <= cpMin {
			return EmptyStr
		}
		each := ToInt(e)
		n := (cpMax - cpMin) * each
		buf := make([]byte, n+each)
		tr := TEXTRANGE{
			chrg:      CHARRANGE{cpMin: int32(cpMin), cpMax: int32(cpMax)},
			lpstrText: &buf[0],
		}
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			0,
			uintptr(unsafe.Pointer(&tr)))
		return SuStr(string(buf[:rtn]))
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		TOOLINFO* lParam) pointer
var _ = builtin4("SendMessageTOOLINFO(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		t := TOOLINFO{
			cbSize:   uint32(unsafe.Sizeof(TOOLINFO{})),
			uFlags:   getUint32(d, "uFlags"),
			hwnd:     getHandle(d, "hwnd"),
			uId:      getUint32(d, "uId"),
			hinst:    getHandle(d, "hinst"),
			lpszText: getStr(d, "lpszText"),
			lParam:   getInt32(d, "lParam"),
			rect:     getRect(d, "rect"),
		}
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&t)))
		return intRet(rtn)
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		TOOLINFO2* lParam) pointer
var _ = builtin4("SendMessageTOOLINFO2(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		t := TOOLINFO2{
			cbSize:   uint32(unsafe.Sizeof(TOOLINFO{})),
			uFlags:   getUint32(d, "uFlags"),
			hwnd:     getHandle(d, "hwnd"),
			uId:      getUint32(d, "uId"),
			hinst:    getHandle(d, "hinst"),
			lpszText: uintptr(getInt(d, "lpszText")), // for LPSTR_TEXTCALLBACK
			lParam:   getInt32(d, "lParam"),
			rect:     getRect(d, "rect"),
		}
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&t)))
		return intRet(rtn)
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		TV_ITEM* tvitem) pointer
var _ = builtin4("SendMessageTreeItem(hwnd, msg, wParam, tvitem)",
	func(a, b, c, d Value) Value {
		cchTextMax := getInt32(d, "cchTextMax")
		var pszText *byte
		var buf []byte
		if cchTextMax == 0 {
			pszText = getStr(d, "pszText")
		} else {
			buf = make([]byte, cchTextMax)
			pszText = &buf[0]
		}
		tvi := TVITEM{
			mask:           getUint32(d, "mask"),
			hItem:          getHandle(d, "hItem"),
			state:          getUint32(d, "state"),
			stateMask:      getUint32(d, "stateMask"),
			pszText:        pszText,
			cchTextMax:     cchTextMax,
			iImage:         getInt32(d, "iImage"),
			iSelectedImage: getInt32(d, "iSelectedImage"),
			cChildren:      getInt32(d, "cChildren"),
			lParam:         getInt32(d, "lParam"),
		}
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&tvi)))
		d.Put(nil, SuStr("mask"), IntVal(int(tvi.mask)))
		d.Put(nil, SuStr("hItem"), IntVal(int(tvi.hItem)))
		d.Put(nil, SuStr("state"), IntVal(int(tvi.state)))
		d.Put(nil, SuStr("stateMask"), IntVal(int(tvi.stateMask)))
		if cchTextMax != 0 {
			d.Put(nil, SuStr("pszText"), SuStr(buf))
		}
		d.Put(nil, SuStr("cchTextMax"), IntVal(int(tvi.cchTextMax)))
		d.Put(nil, SuStr("iImage"), IntVal(int(tvi.iImage)))
		d.Put(nil, SuStr("iSelectedImage"), IntVal(int(tvi.iSelectedImage)))
		d.Put(nil, SuStr("cChildren"), IntVal(int(tvi.cChildren)))
		d.Put(nil, SuStr("lParam"), IntVal(int(tvi.lParam)))
		return intRet(rtn)
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		TV_INSERTSTRUCT* tvins) pointer
var _ = builtin4("SendMessageTVINS(hwnd, msg, wParam, tvins)",
	func(a, b, c, d Value) Value {
		item := d.Get(nil, SuStr("item"))
		cchTextMax := getInt32(item, "cchTextMax")
		var pszText *byte
		var buf []byte
		if cchTextMax == 0 {
			pszText = getStr(item, "pszText")
		} else {
			buf = make([]byte, cchTextMax)
			pszText = &buf[0]
		}
		tvi := TVITEM{
			mask:           getUint32(d, "mask"),
			hItem:          getHandle(d, "hItem"),
			state:          getUint32(d, "state"),
			stateMask:      getUint32(d, "stateMask"),
			pszText:        pszText,
			cchTextMax:     cchTextMax,
			iImage:         getInt32(d, "iImage"),
			iSelectedImage: getInt32(d, "iSelectedImage"),
			cChildren:      getInt32(d, "cChildren"),
			lParam:         getInt32(d, "lParam"),
		}
		tvins := TV_INSERTSTRUCT{
			hParent:      getHandle(d, "hParent"),
			hInsertAfter: getHandle(d, "hInsertAfter"),
			item:         tvi,
		}
		rtn, _, _ := sendMessage.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(unsafe.Pointer(&tvins)))
		d.Put(nil, SuStr("hParent"), IntVal(int(tvins.hParent)))
		d.Put(nil, SuStr("hInsertAfter"), IntVal(int(tvins.hInsertAfter)))
		item.Put(nil, SuStr("mask"), IntVal(int(tvi.mask)))
		item.Put(nil, SuStr("hItem"), IntVal(int(tvi.hItem)))
		item.Put(nil, SuStr("state"), IntVal(int(tvi.state)))
		item.Put(nil, SuStr("stateMask"), IntVal(int(tvi.stateMask)))
		item.Put(nil, SuStr("cchTextMax"), IntVal(int(tvi.cchTextMax)))
		item.Put(nil, SuStr("iImage"), IntVal(int(tvi.iImage)))
		item.Put(nil, SuStr("iSelectedImage"), IntVal(int(tvi.iSelectedImage)))
		item.Put(nil, SuStr("cChildren"), IntVal(int(tvi.cChildren)))
		item.Put(nil, SuStr("lParam"), IntVal(int(tvi.lParam)))
		if cchTextMax != 0 {
			item.Put(nil, SuStr("pszText"), SuStr(buf))
		}
		return intRet(rtn)
	})

// dll User32:SetFocus(pointer hwnd) pointer
var setFocus = user32.NewProc("SetFocus")
var _ = builtin1("SetFocus(hwnd)",
	func(a Value) Value {
		rtn, _, _ := setFocus.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll User32:SetTimer(pointer hwnd, long id, long ms, TIMERPROC f) long
var setTimer = user32.NewProc("SetTimer")
var _ = builtin4("SetTimer(hwnd, id, ms, f)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := setTimer.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			NewCallback(d, 4))
		return intRet(rtn)
	})

// dll User32:SetWindowLong(pointer hwnd, int offset, long value) long
var setWindowLong = user32.NewProc("SetWindowLongA")
var _ = builtin3("SetWindowLong(hwnd, offset, value)",
	func(a, b, c Value) Value {
		rtn, _, _ := setWindowLong.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll User32:SetWindowLong(pointer hwnd, long offset, long value) long
var setWindowLongPtr = user32.NewProc("SetWindowLongPtrA")
var _ = builtin3("SetWindowLongPtr(hwnd, offset, value)",
	func(a, b, c Value) Value {
		rtn, _, _ := setWindowLongPtr.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll User32:SetWindowProc(pointer hwnd, long offset, WNDPROC proc) pointer
var _ = builtin3("SetWindowProc(hwnd, offset, value)",
	func(a, b, c Value) Value {
		rtn, _, _ := setWindowLongPtr.Call(
			intArg(a),
			intArg(b),
			NewCallback(c, 4))
		return intRet(rtn)
	})

// dll User32:SetWindowPlacement(pointer hwnd, WINDOWPLACEMENT* lpwndpl) bool
var setWindowPlacement = user32.NewProc("SetWindowPlacement")
var _ = builtin2("SetWindowPlacement(hwnd, lpwndpl)",
	func(a, b Value) Value {
		w := WINDOWPLACEMENT{
			length:           getUint32(b, "length"),
			flags:            getUint32(b, "flags"),
			showCmd:          getUint32(b, "showCmd"),
			ptMinPosition:    getPoint(b, "ptMinPosition"),
			ptMaxPosition:    getPoint(b, "ptMaxPosition"),
			rcNormalPosition: getRect(b, "rcNormalPosition"),
		}
		rtn, _, _ := setWindowPlacement.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&w)))
		return boolRet(rtn)
	})

// dll User32:SetWindowPos(pointer hWnd, pointer hWndInsertAfter,
//		long X, long Y, long cx, long cy, long uFlags) bool
var setWindowPos = user32.NewProc("SetWindowPos")
var _ = builtin7("SetWindowPos(hWnd, hWndInsertAfter, X, Y, cx, cy, uFlags)",
	func(a, b, c, d, e, f, g Value) Value {
		rtn, _, _ := setWindowPos.Call(
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
var setWindowText = user32.NewProc("SetWindowTextA")
var _ = builtin2("SetWindowText(hwnd, lpwndpl)",
	func(a, b Value) Value {
		rtn, _, _ := setWindowText.Call(
			intArg(a),
			uintptr(stringArg(b)))
		return boolRet(rtn)
	})

// dll User32:ShowWindow(pointer hwnd, long ncmd) bool
var _ = builtin2("ShowWindow(hwnd, ncmd)",
	func(a, b Value) Value {
		var showWindow = user32.NewProc("ShowWindow")
		rtn, _, _ := showWindow.Call(
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll User32:SystemParametersInfo(long uiAction, long uiParam, ? pvParam,
//		long fWinIni) bool
var systemParametersInfo = user32.NewProc("SystemParametersInfoA")

var _ = builtin0("SPI_GetFocusBorderHeight()",
	func() Value {
		var x int32
		const SPI_GETFOCUSBORDERHEIGHT = 0x2010
		systemParametersInfo.Call(
			SPI_GETFOCUSBORDERHEIGHT,
			0,
			uintptr(unsafe.Pointer(&x)),
			0)
		return IntVal(int(x))
	})
var _ = builtin0("SPI_GetWheelScrollLines()",
	func() Value {
		var x int32
		const SPI_GETWHEELSCROLLLINES = 104
		systemParametersInfo.Call(
			SPI_GETWHEELSCROLLLINES,
			0,
			uintptr(unsafe.Pointer(&x)),
			0)
		return IntVal(int(x))
	})
var _ = builtin0("SPI_GetWorkArea()",
	func() Value {
		var r RECT
		const SPI_GETWORKAREA = 48
		systemParametersInfo.Call(
			SPI_GETWORKAREA,
			0,
			uintptr(unsafe.Pointer(&r)),
			0)
		return rectToOb(&r, nil)
	})

// dll User32:PostQuitMessage(long exitcode) void
var postQuitMessage = user32.NewProc("PostQuitMessage")
var _ = builtin1("PostQuitMessage(exitcode)",
	func(a Value) Value {
		rtn, _, _ := postQuitMessage.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll User32:GetNextDlgTabItem(pointer hDlg, pointer hCtl, bool prev) pointer
var getNextDlgTabItem = user32.NewProc("GetNextDlgTabItem")
var _ = builtin3("GetNextDlgTabItem(hDlg, hCtl, prev)",
	func(a, b, c Value) Value {
		rtn, _, _ := getNextDlgTabItem.Call(
			intArg(a),
			intArg(b),
			boolArg(c))
		return intRet(rtn)
	})

// dll User32:UpdateWindow(pointer hwnd) bool
var updateWindow = user32.NewProc("UpdateWindow")
var _ = builtin1("UpdateWindow(hwnd)",
	func(a Value) Value {
		rtn, _, _ := updateWindow.Call(intArg(a))
		return boolRet(rtn)
	})

// dll User32:DefWindowProc(pointer hwnd, long msg, pointer wParam,
//		pointer lParam) pointer
var defWindowProc = user32.NewProc("DefWindowProcA")
var _ = builtin4("DefWindowProc(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := defWindowProc.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return intRet(rtn)
	})

var _ = builtin0("GetDefWindowProc()",
	func() Value {
		return IntVal(int(defWindowProc.Addr()))
	})

// dll User32:GetKeyState(long key) short
var getKeyState = user32.NewProc("GetKeyState")
var _ = builtin1("GetKeyState(nVirtKey)",
	func(a Value) Value {
		rtn, _, _ := getKeyState.Call(intArg(a))
		return intRet(rtn)
	})

type TPMPARAMS struct {
	cbSize    int32
	rcExclude RECT
}

// dll long User32:TrackPopupMenuEx(pointer hmenu, long fuFlags, long x, long y,
//		pointer hwnd, TPMPARAMS* lptpm)
var trackPopupMenuEx = user32.NewProc("TrackPopupMenuEx")
var _ = builtin6("TrackPopupMenuEx(hmenu, fuFlags, x, y, hwnd, lptpm)",
	func(a, b, c, d, e, f Value) Value {
		var lptpm *TPMPARAMS
		if f.Equal(Zero) {
			lptpm = nil
		} else {
			tpm := TPMPARAMS{
				cbSize:    int32(unsafe.Sizeof(TPMPARAMS{})),
				rcExclude: getRect(f, "rcExclude"),
			}
			lptpm = &tpm
		}
		rtn, _, _ := trackPopupMenuEx.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			uintptr(unsafe.Pointer(lptpm)))
		return intRet(rtn)
	})

// dll bool User32:OpenClipboard(pointer hwnd)
var openClipboard = user32.NewProc("OpenClipboard")
var _ = builtin1("OpenClipboard(hwnd)",
	func(a Value) Value {
		rtn, _, _ := openClipboard.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:EmptyClipboard()
var emptyClipboard = user32.NewProc("EmptyClipboard")
var _ = builtin0("EmptyClipboard()",
	func() Value {
		rtn, _, _ := emptyClipboard.Call()
		return boolRet(rtn)
	})

// dll pointer User32:GetClipboardData(long format)
var getClipboardData = user32.NewProc("GetClipboardData")
var _ = builtin1("GetClipboardData(format)",
	func(a Value) Value {
		rtn, _, _ := getClipboardData.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:SetClipboardData(long uFormat, pointer hMem)
var setClipboardData = user32.NewProc("SetClipboardData")
var _ = builtin2("SetClipboardData(uFormat, hMem)",
	func(a Value, b Value) Value {
		rtn, _, _ := setClipboardData.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll bool User32:CloseClipboard()
var closeClipboard = user32.NewProc("CloseClipboard")
var _ = builtin0("CloseClipboard()",
	func() Value {
		rtn, _, _ := closeClipboard.Call()
		return boolRet(rtn)
	})

// dll pointer User32:BeginDeferWindowPos(
// 	long nNumWindows		// initial number of windows to allocate space for
// 	)
var beginDeferWindowPos = user32.NewProc("BeginDeferWindowPos")
var _ = builtin1("BeginDeferWindowPos(nNumWindows)",
	func(a Value) Value {
		rtn, _, _ := beginDeferWindowPos.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:CallNextHookEx(		// returns an LRESULT
// 	pointer	hhk,	// handle to current hook [HHOOK]
// 	long	nCode,	// hook code passed to hook procedure [int]
// 	pointer	wParam,	// value passed to hook procedure [WPARAM]
// 	pointer	lParam	// value passed to hook procedure [LPARAM]
// )
var callNextHookEx = user32.NewProc("CallNextHookEx")
var _ = builtin4("CallNextHookEx(hhk, nCode, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := callNextHookEx.Call(
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
var checkMenuItem = user32.NewProc("CheckMenuItem")
var _ = builtin3("CheckMenuItem(hmenu, uIDCheckItem, uCheck)",
	func(a, b, c Value) Value {
		rtn, _, _ := checkMenuItem.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll bool user32:DeleteMenu(pointer hMenu, long uPosition, long uFlags)
var deleteMenu = user32.NewProc("DeleteMenu")
var _ = builtin3("DeleteMenu(hMenu, uPosition, uFlags)",
	func(a, b, c Value) Value {
		rtn, _, _ := deleteMenu.Call(
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
var enableMenuItem = user32.NewProc("EnableMenuItem")
var _ = builtin3("EnableMenuItem(hMenu, uIDEnableItem, uEnable)",
	func(a, b, c Value) Value {
		rtn, _, _ := enableMenuItem.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll bool User32:EnableWindow(pointer hWnd, bool bEnable)
var enableWindow = user32.NewProc("EnableWindow")
var _ = builtin2("EnableWindow(hWnd, bEnable)",
	func(a, b Value) Value {
		rtn, _, _ := enableWindow.Call(
			intArg(a),
			boolArg(b))
		return boolRet(rtn)
	})

// dll bool User32:EndDeferWindowPos(
// 	pointer hWinPosInfo		// handle to internal structure
// 	)
var endDeferWindowPos = user32.NewProc("EndDeferWindowPos")
var _ = builtin1("EndDeferWindowPos(hWinPosInfo)",
	func(a Value) Value {
		rtn, _, _ := endDeferWindowPos.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:EndDialog(pointer hwndDlg, long nResult)
var endDialog = user32.NewProc("EndDialog")
var _ = builtin2("EndDialog(hwndDlg, nResult)",
	func(a, b Value) Value {
		rtn, _, _ := endDialog.Call(
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll long User32:EnumClipboardFormats(long format)
var enumClipboardFormats = user32.NewProc("EnumClipboardFormats")
var _ = builtin1("EnumClipboardFormats(format)",
	func(a Value) Value {
		rtn, _, _ := enumClipboardFormats.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:FindWindow([in] string c, [in] string n)
var findWindow = user32.NewProc("FindWindow")
var _ = builtin2("FindWindow(c, n)",
	func(a, b Value) Value {
		rtn, _, _ := findWindow.Call(
			uintptr(stringArg(a)),
			uintptr(stringArg(b)))
		return intRet(rtn)
	})

// dll pointer User32:GetAncestor(pointer hwnd, long gaFlags)
var getAncestor = user32.NewProc("GetAncestor")
var _ = builtin2("GetAncestor(hwnd, gaFlags)",
	func(a, b Value) Value {
		rtn, _, _ := getAncestor.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll long User32:GetClipboardFormatName(
// 	long		format,				// clipboard format to retrieve
// 	string		lpszFormatName,		// buffer to receive format name
// 	long		cchMaxCount			// maximum length of string to copy into buffer
// 	)
var getClipboardFormatName = user32.NewProc("GetClipboardFormatName")
var _ = builtin3("GetClipboardFormatName(format, lpszFormatName, cchMaxCount)",
	func(a, b, c Value) Value {
		rtn, _, _ := getClipboardFormatName.Call(
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c))
		return intRet(rtn)
	})

// dll pointer User32:GetCursor()
var getCursor = user32.NewProc("GetCursor")
var _ = builtin0("GetCursor()",
	func() Value {
		rtn, _, _ := getCursor.Call()
		return intRet(rtn)
	})

// dll long User32:GetDoubleClickTime()
var getDoubleClickTime = user32.NewProc("GetDoubleClickTime")
var _ = builtin0("GetDoubleClickTime()",
	func() Value {
		rtn, _, _ := getDoubleClickTime.Call()
		return intRet(rtn)
	})

// dll long User32:GetMenuState(
// 	pointer hMenu, 	// handle to menu
// 	long uId, 		// menu item to query
// 	long uFlags		// options
// 	)
var getMenuState = user32.NewProc("GetMenuState")
var _ = builtin3("GetMenuState(hMenu, uId, uFlags)",
	func(a, b, c Value) Value {
		rtn, _, _ := getMenuState.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll long User32:GetMessagePos()
var getMessagePos = user32.NewProc("GetMessagePos")
var _ = builtin0("GetMessagePos()",
	func() Value {
		rtn, _, _ := getMessagePos.Call()
		return intRet(rtn)
	})

// dll pointer User32:GetNextDlgGroupItem(pointer hDlg, pointer hCtl, bool prev)
var getNextDlgGroupItem = user32.NewProc("GetNextDlgGroupItem")
var _ = builtin3("GetNextDlgGroupItem(hDlg, hCtl, prev)",
	func(a, b, c Value) Value {
		rtn, _, _ := getNextDlgGroupItem.Call(
			intArg(a),
			intArg(b),
			boolArg(c))
		return intRet(rtn)
	})

// dll pointer User32:GetParent(pointer hwnd)
var getParent = user32.NewProc("GetParent")
var _ = builtin1("GetParent(hwnd)",
	func(a Value) Value {
		rtn, _, _ := getParent.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:GetSubMenu(
// 	pointer hmenu,	//menu handle
// 	long position	//position
// 	)
var getSubMenu = user32.NewProc("GetSubMenu")
var _ = builtin2("GetSubMenu(hmenu, position)",
	func(a, b Value) Value {
		rtn, _, _ := getSubMenu.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer User32:GetTopWindow(pointer hWnd)
var getTopWindow = user32.NewProc("GetTopWindow")
var _ = builtin1("GetTopWindow(hWnd)",
	func(a Value) Value {
		rtn, _, _ := getTopWindow.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll long User32:GetUpdateRgn(pointer hwnd, pointer hRgn, bool bErase)
var getUpdateRgn = user32.NewProc("GetUpdateRgn")
var _ = builtin3("GetUpdateRgn(hwnd, hRgn, bErase)",
	func(a, b, c Value) Value {
		rtn, _, _ := getUpdateRgn.Call(
			intArg(a),
			intArg(b),
			boolArg(c))
		return intRet(rtn)
	})

// dll pointer User32:GetWindow(pointer hWnd, long uCmd)
var getWindow = user32.NewProc("GetWindow")
var _ = builtin2("GetWindow(hWnd, uCmd)",
	func(a, b Value) Value {
		rtn, _, _ := getWindow.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer User32:GetWindowDC(pointer hwnd)
var getWindowDC = user32.NewProc("GetWindowDC")
var _ = builtin1("GetWindowDC(hwnd)",
	func(a Value) Value {
		rtn, _, _ := getWindowDC.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll bool User32:IsClipboardFormatAvailable(long format)
var isClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
var _ = builtin1("IsClipboardFormatAvailable(format)",
	func(a Value) Value {
		rtn, _, _ := isClipboardFormatAvailable.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:IsWindow(pointer hwnd)
var isWindow = user32.NewProc("IsWindow")
var _ = builtin1("IsWindow(hwnd)",
	func(a Value) Value {
		rtn, _, _ := isWindow.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:IsWindowVisible(pointer hwnd)
var isWindowVisible = user32.NewProc("IsWindowVisible")
var _ = builtin1("IsWindowVisible(hwnd)",
	func(a Value) Value {
		rtn, _, _ := isWindowVisible.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:MessageBeep(long type)
var messageBeep = user32.NewProc("MessageBeep")
var _ = builtin1("MessageBeep(type)",
	func(a Value) Value {
		rtn, _, _ := messageBeep.Call(
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
var mouse_event = user32.NewProc("mouse_event")
var _ = builtin5("mouse_event(dwFlags, dx, dy, dwData, dwExtraInfo)",
	func(a, b, c, d, e Value) Value {
		mouse_event.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return nil
	})

// dll bool User32:PostMessage(pointer hwnd, long msg, pointer wParam, pointer lParam)
var postMessage = user32.NewProc("PostMessage")
var _ = builtin4("PostMessage(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := postMessage.Call(
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
var registerHotKey = user32.NewProc("RegisterHotKey")
var _ = builtin4("RegisterHotKey(hWnd, id, fsModifiers, vk)",
	func(a, b, c, d Value) Value {
		rtn, _, _ := registerHotKey.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return boolRet(rtn)
	})

// dll bool User32:ReleaseCapture()
var releaseCapture = user32.NewProc("ReleaseCapture")
var _ = builtin0("ReleaseCapture()",
	func() Value {
		rtn, _, _ := releaseCapture.Call()
		return boolRet(rtn)
	})

// dll pointer User32:SetActiveWindow(pointer hWnd)
var setActiveWindow = user32.NewProc("SetActiveWindow")
var _ = builtin1("SetActiveWindow(hWnd)",
	func(a Value) Value {
		rtn, _, _ := setActiveWindow.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:SetCapture(pointer hwnd)
var setCapture = user32.NewProc("SetCapture")
var _ = builtin1("SetCapture(hwnd)",
	func(a Value) Value {
		rtn, _, _ := setCapture.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll pointer User32:SetCursor(pointer hcursor)
var setCursor = user32.NewProc("SetCursor")
var _ = builtin1("SetCursor(hcursor)",
	func(a Value) Value {
		rtn, _, _ := setCursor.Call(
			intArg(a))
		return intRet(rtn)
	})

// dll bool User32:SetForegroundWindow(pointer hwnd)
// // puts creator thread into foreground and activates window
var setForegroundWindow = user32.NewProc("SetForegroundWindow")
var _ = builtin1("SetForegroundWindow(hwnd)",
	func(a Value) Value {
		rtn, _, _ := setForegroundWindow.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:SetMenuDefaultItem(
// 	pointer hMenu,
// 	long uItem,
// 	long fByPosition
// 	)
var setMenuDefaultItem = user32.NewProc("SetMenuDefaultItem")
var _ = builtin3("SetMenuDefaultItem(hMenu, uItem, fByPosition)",
	func(a, b, c Value) Value {
		rtn, _, _ := setMenuDefaultItem.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll pointer User32:SetParent(pointer hwndNewChild, pointer hwndNewParent)
var setParent = user32.NewProc("SetParent")
var _ = builtin2("SetParent(hwndNewChild, hwndNewParent)",
	func(a, b Value) Value {
		rtn, _, _ := setParent.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll bool User32:SetProp(pointer hwnd, [in] string name, pointer value)
var setProp = user32.NewProc("SetProp")
var _ = builtin3("SetProp(hwnd, name, value)",
	func(a, b, c Value) Value {
		rtn, _, _ := setProp.Call(
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool User32:UnhookWindowsHookEx(
// 	pointer hhk	// handle to hook procedure [HHOOK]
// )
var unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
var _ = builtin1("UnhookWindowsHookEx(hhk)",
	func(a Value) Value {
		rtn, _, _ := unhookWindowsHookEx.Call(
			intArg(a))
		return boolRet(rtn)
	})

// dll bool User32:UnregisterHotKey(pointer hWnd /*optional*/, long id)
var unregisterHotKey = user32.NewProc("UnregisterHotKey")
var _ = builtin2("UnregisterHotKey(hWnd, id)",
	func(a, b Value) Value {
		rtn, _, _ := unregisterHotKey.Call(
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll bool User32:ClientToScreen(pointer hwnd, POINT* point)
var clientToScreen = user32.NewProc("ClientToScreen")
var _ = builtin2("ClientToScreen(hWnd, point)",
	func(a, b Value) Value {
		var pt POINT
		rtn, _, _ := clientToScreen.Call(
			intArg(a),
			uintptr(pointArg(b, &pt)))
		return intRet(rtn)
	})

// dll bool User32:ClipCursor(RECT* rect)
var clipCursor = user32.NewProc("ClipCursor")
var _ = builtin1("ClipCursor(rect)",
	func(a Value) Value {
		var r RECT
		rtn, _, _ := clipCursor.Call(
			uintptr(rectArg(a, &r)))
		return boolRet(rtn)
	})

// dll pointer User32:DeferWindowPos(pointer hWinPosInfo, pointer hWnd,
//		pointer hWndInsertAfter, long x, long y, long cx, long cy, long flags)
var deferWindowPos = user32.NewProc("DeferWindowPos")
var _ = builtin("DeferWindowPos(hWinPosInfo, hWnd, hWndInsertAfter, "+
	"x, y, cx, cy, flags)",
	func(_ *Thread, a []Value) Value {
		rtn, _, _ := deferWindowPos.Call(
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
var drawFocusRect = user32.NewProc("DrawFocusRect")
var _ = builtin2("DrawFocusRect(hwnd, rect)",
	func(a, b Value) Value {
		var r RECT
		rtn, _, _ := drawFocusRect.Call(
			intArg(a),
			uintptr(rectArg(b, &r)))
		return intRet(rtn)
	})

// dll long User32:DrawTextEx(pointer hdc, [in] string lpsz, long cb,
// RECT* lprc, long uFormat, DRAWTEXTPARAMS* params)
var drawTextEx = user32.NewProc("DrawTextEx")
var _ = builtin6("DrawTextEx(hdc, lpsz, cb, lprc, uFormat, params)",
	func(a, b, c, d, e, f Value) Value {
		var r RECT
		dtp := DRAWTEXTPARAMS{
			cbSize:        int32(unsafe.Sizeof(DRAWTEXTPARAMS{})),
			iTabLength:    getInt32(f, "iTabLength"),
			iLeftMargin:   getInt32(f, "iLeftMargin"),
			iRightMargin:  getInt32(f, "iRightMargin"),
			uiLengthDrawn: getInt32(f, "uiLengthDrawn"),
		}
		rtn, _, _ := drawTextEx.Call(
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c),
			uintptr(rectArg(d, &r)),
			intArg(e),
			uintptr(unsafe.Pointer(&dtp)))
		return intRet(rtn)
	})

type DRAWTEXTPARAMS struct {
	cbSize        int32
	iTabLength    int32
	iLeftMargin   int32
	iRightMargin  int32
	uiLengthDrawn int32
}
