// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"encoding/binary"
	"strings"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type suStructGlobal struct {
	SuBuiltin
	size int
}

func init() {
	Global.Builtin("INITCOMMONCONTROLSEX",
		&suStructGlobal{size: int(unsafe.Sizeof(stInitCommonControlsEx{}))})
	Global.Builtin("MENUITEMINFO",
		&suStructGlobal{size: int(unsafe.Sizeof(stMenuItemInfo{}))})
	Global.Builtin("MONITORINFO",
		&suStructGlobal{size: int(unsafe.Sizeof(stMonitorInfo{}))})
	Global.Builtin("SCROLLINFO",
		&suStructGlobal{size: int(unsafe.Sizeof(stScrollInfo{}))})
	Global.Builtin("TRACKMOUSEEVENT",
		&suStructGlobal{size: int(unsafe.Sizeof(stTrackMouseEvent{}))})
	Global.Builtin("SHELLEXECUTEINFO",
		&suStructGlobal{size: int(unsafe.Sizeof(stShellExecuteInfo{}))})
	Global.Builtin("TOOLINFO",
		&suStructGlobal{size: int(unsafe.Sizeof(stToolInfo{}))})
	Global.Builtin("WINDOWPLACEMENT",
		&suStructGlobal{size: int(unsafe.Sizeof(stWindowPlacement{}))})
	Global.Builtin("OPENFILENAME",
		&suStructGlobal{size: int(nOpenFileName)})
	Global.Builtin("TPMPARAMS",
		&suStructGlobal{size: int(nTPMParams)})
	Global.Builtin("DRAWTEXTPARAMS",
		&suStructGlobal{size: int(nDrawTextParams)})
	Global.Builtin("SYSTEMTIME",
		&suStructGlobal{size: int(nSystemTime)})
	Global.Builtin("PRINTDLG",
		&suStructGlobal{size: int(nPrintDlg)})
	Global.Builtin("PAGESETUPDLG",
		&suStructGlobal{size: int(nPageSetupDlg)})
	Global.Builtin("DOCINFO",
		&suStructGlobal{size: int(nDocInfo)})
	Global.Builtin("PRINTDLGEX",
		&suStructGlobal{size: int(nPrintDlgEx)})
	Global.Builtin("CHOOSEFONT",
		&suStructGlobal{size: int(nChooseFont)})
	Global.Builtin("CHOOSECOLOR",
		&suStructGlobal{size: int(nChooseColor)})
	Global.Builtin("FLASHWINFO",
		&suStructGlobal{size: int(nFlashWInfo)})
	Global.Builtin("EDITBALLOONTIP",
		&suStructGlobal{size: int(nEditBalloonTip)})
}

func (*suStructGlobal) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call struct")
}

func (sg *suStructGlobal) Size() int {
	return sg.size
}

type sizeable interface{ Size() int }

var structMethods = methods("st")

var _ = method(st_Size, "()")

func st_Size(this Value) Value {
	return IntVal(this.(sizeable).Size())
}

func (*suStructGlobal) Lookup(_ *Thread, method string) Value {
	return structMethods[method]
}

func (*suStructGlobal) String() string {
	return "/* builtin struct */"
}

type callableStruct struct {
	suStructGlobal
}

func (cs *callableStruct) Call(th *Thread, this Value, as *ArgSpec) Value {
	return cs.SuBuiltin.Call(th, this, as)
}

//-------------------------------------------------------------------

type Struct interface {
	fromStruct(p unsafe.Pointer) Value
	updateStruct(ob Value, p unsafe.Pointer)
}

// StructModify allows modifying a struct at a given address
// WARNING: address must be valid for type
var _ = builtin(StructModify, "(type, address, block)")

func StructModify(th *Thread, args []Value) Value {
	typ, ok := args[0].(Struct)
	if !ok {
		panic("StructModify invalid type " + ErrType(args[0]))
	}
	p := unsafe.Pointer(uintptr(ToInt(args[1])))
	if p == nil {
		panic("StructModify: address can't be zero")
	}
	ob := typ.fromStruct(p)
	th.Call(args[2], ob) // call the block, which modifies ob
	typ.updateStruct(ob, p)
	return nil
}

//-------------------------------------------------------------------

var _ = Global.Builtin("RECT",
	&suRect{callableStruct{suStructGlobal{size: int(unsafe.Sizeof(stRect{})),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            rect}}}})

type suRect struct {
	callableStruct
}

func (*suRect) fromStruct(p unsafe.Pointer) Value {
	return fromRect((*stRect)(p), nil)
}

func (*suRect) updateStruct(ob Value, p unsafe.Pointer) {
	*(*stRect)(p) = *toRect(ob)
}

func rect(_ *Thread, args []Value) Value {
	r := toRect(args[0])
	return ptrNstr(unsafe.Pointer(r), nRect)
}

//-------------------------------------------------------------------

var _ = Global.Builtin("MINMAXINFO",
	&suMinMaxInfo{suStructGlobal{size: int(unsafe.Sizeof(stMinMaxInfo{}))}})

type suMinMaxInfo struct {
	suStructGlobal
}

type stMinMaxInfo struct {
	ptReserved   stPoint
	maxSize      stPoint
	maxPosition  stPoint
	minTrackSize stPoint
	maxTrackSize stPoint
}

func (typ *suMinMaxInfo) fromStruct(p unsafe.Pointer) Value {
	mmi := (*stMinMaxInfo)(p)
	ob := &SuObject{}
	ob.Put(nil, SuStr("maxSize"), fromPoint(&mmi.maxSize, nil))
	ob.Put(nil, SuStr("maxPosition"), fromPoint(&mmi.maxPosition, nil))
	ob.Put(nil, SuStr("minTrackSize"), fromPoint(&mmi.minTrackSize, nil))
	ob.Put(nil, SuStr("maxTrackSize"), fromPoint(&mmi.maxTrackSize, nil))
	return ob
}

func (typ *suMinMaxInfo) updateStruct(ob Value, p unsafe.Pointer) {
	mmi := (*stMinMaxInfo)(p)
	mmi.maxSize = getPoint(ob, "maxSize")
	mmi.maxPosition = getPoint(ob, "maxPosition")
	mmi.minTrackSize = getPoint(ob, "minTrackSize")
	mmi.maxTrackSize = getPoint(ob, "maxTrackSize")
}

//-------------------------------------------------------------------

type stNMHdr struct {
	hwndFrom uintptr
	idFrom   uintptr
	code     uint32
	_        [4]byte // padding
}

var _ = builtin(NMHDR, "(address)")

func NMHDR(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	nmh := (*stNMHdr)(unsafe.Pointer(uintptr(adr)))
	return fromNMHdr(nmh)
}

func fromNMHdr(nmh *stNMHdr) *SuObject {
	ob := &SuObject{}
	ob.Put(nil, SuStr("hwndFrom"), IntVal(int(nmh.hwndFrom)))
	ob.Put(nil, SuStr("idFrom"), IntVal(int(nmh.idFrom)))
	ob.Put(nil, SuStr("code"), IntVal(int(nmh.code)))
	return ob
}

//-------------------------------------------------------------------

type stNMTVDispInfo struct {
	nmhdr stNMHdr
	item  stTVItem
}

var _ = builtin(NMTVDISPINFO, "(address)")

func NMTVDISPINFO(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	di := (*stNMTVDispInfo)(unsafe.Pointer(uintptr(adr)))
	ob := fromNMHdr(&di.nmhdr)
	ob.Put(nil, SuStr("nmhdr"), fromNMHdr(&di.nmhdr))
	tvi := fromTVItem(&di.item)
	tvi.Put(nil, SuStr("pszText"),
		ptrZstr(unsafe.Pointer(di.item.pszText), 1024))
	ob.Put(nil, SuStr("item"), tvi)
	return ob
}

func fromTVItem(tvi *stTVItem) *SuObject {
	ob := &SuObject{}
	ob.Put(nil, SuStr("mask"), IntVal(int(tvi.mask)))
	ob.Put(nil, SuStr("hItem"), IntVal(int(tvi.hItem)))
	ob.Put(nil, SuStr("state"), IntVal(int(tvi.state)))
	ob.Put(nil, SuStr("stateMask"), IntVal(int(tvi.stateMask)))
	// pszText must be handled by caller
	ob.Put(nil, SuStr("cchTextMax"), IntVal(int(tvi.cchTextMax)))
	ob.Put(nil, SuStr("iImage"), IntVal(int(tvi.iImage)))
	ob.Put(nil, SuStr("iSelectedImage"), IntVal(int(tvi.iSelectedImage)))
	ob.Put(nil, SuStr("cChildren"), IntVal(int(tvi.cChildren)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(tvi.lParam)))
	return ob
}

//-------------------------------------------------------------------

var _ = Global.Builtin("NMTVDISPINFO2",
	&suNMTVDISPINFO{callableStruct{
		suStructGlobal{size: int(unsafe.Sizeof(stNMTVDispInfo{})),
			SuBuiltin: SuBuiltin{
				BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
				Fn:            fromNMTVDispInfo}}}})

type suNMTVDISPINFO struct {
	callableStruct
}

func (*suNMTVDISPINFO) fromStruct(p unsafe.Pointer) Value {
	x := (*stNMTVDispInfo)(p)
	ob := &SuObject{}
	ob.Put(nil, SuStr("hdr"), fromNMHdr(&x.nmhdr))
	ob.Put(nil, SuStr("item"), fromTVItem(&x.item))
	return ob
}

func (*suNMTVDISPINFO) updateStruct(ob Value, p unsafe.Pointer) {
	tvi := ob.Get(nil, SuStr("item"))
	x := (*stNMTVDispInfo)(p)
	x.item.mask = getUint32(tvi, "mask")
	x.item.hItem = getUintptr(tvi, "hItem")
	x.item.state = getUint32(tvi, "state")
	x.item.stateMask = getUint32(tvi, "stateMask")
	x.item.cchTextMax = getInt32(tvi, "cchTextMax")
	x.item.iImage = getInt32(tvi, "iImage")
	x.item.iSelectedImage = getInt32(tvi, "iSelectedImage")
	x.item.cChildren = getInt32(tvi, "cChildren")
	x.item.lParam = getUintptr(tvi, "lParam")
}

func fromNMTVDispInfo(_ *Thread, args []Value) Value {
	adr := ToInt(args[0])
	if adr == 0 {
		return False
	}
	var x *suNMTVDISPINFO
	return x.fromStruct(unsafe.Pointer(uintptr(adr)))
}

//-------------------------------------------------------------------

var _ = Global.Builtin("NMTTDISPINFO2",
	&suNMTTDISPINFO{suStructGlobal{size: int(unsafe.Sizeof(stNMTTDispInfo{}))}})

type suNMTTDISPINFO struct {
	suStructGlobal
}

func (*suNMTTDISPINFO) fromStruct(p unsafe.Pointer) Value {
	x := (*stNMTTDispInfo)(p)
	ob := &SuObject{}
	ob.Put(nil, SuStr("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, SuStr("szText"), bufZstr(x.szText[:]))
	ob.Put(nil, SuStr("lpszText"), IntVal(int(x.lpszText)))
	ob.Put(nil, SuStr("hinst"), IntVal(int(x.hinst)))
	ob.Put(nil, SuStr("uFlags"), IntVal(int(x.uFlags)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(x.lParam)))
	return ob
}

func (*suNMTTDISPINFO) updateStruct(ob Value, p unsafe.Pointer) {
	x := (*stNMTTDispInfo)(p)
	x.lpszText = getUintptr(ob, "lpszText")
	getZstrBs(ob, "szText", x.szText[:])
	x.hinst = getUintptr(ob, "hinst")
	x.uFlags = getInt32(ob, "uFlags")
	x.lParam = getUintptr(ob, "lParam")
}

type stNMTTDispInfo struct {
	hdr      stNMHdr
	lpszText uintptr
	szText   [80]byte
	hinst    uintptr
	uFlags   int32
	lParam   uintptr
}

//-------------------------------------------------------------------

type stNMHeader struct {
	hdr     stNMHdr
	iItem   int32
	iButton int32
	pitem   *stHdItem
}

var _ = builtin(NMHEADER, "(address)")

func NMHEADER(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMHeader)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, SuStr("iItem"), IntVal(int(x.iItem)))
	ob.Put(nil, SuStr("iButton"), IntVal(int(x.iButton)))
	if x.pitem != nil {
		hdi := fromHdItem(x.pitem, &SuObject{})
		hdi.Put(nil, SuStr("pszText"),
			IntVal(int(uintptr(unsafe.Pointer(x.pitem.pszText)))))
		ob.Put(nil, SuStr("pitem"), hdi)
	}
	return ob
}

//-------------------------------------------------------------------

type stNMTreeView struct {
	hdr     stNMHdr
	action  int32
	itemOld stTVItem
	itemNew stTVItem
	ptDrag  stPoint
}

var _ = builtin(NMTREEVIEW, "(address)")

func NMTREEVIEW(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMTreeView)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, SuStr("action"), IntVal(int(x.action)))
	ob.Put(nil, SuStr("itemOld"), fromTVItem(&x.itemOld))
	ob.Put(nil, SuStr("itemNew"), fromTVItem(&x.itemNew))
	ob.Put(nil, SuStr("ptDrag"), fromPoint(&x.ptDrag, nil))
	return ob
}

//-------------------------------------------------------------------

type stNMTVKeyDown struct {
	hdr   stNMHdr
	wVKey int16
	flags int32
}

var _ = builtin(NMTVKEYDOWN, "(address)")

func NMTVKEYDOWN(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMTVKeyDown)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, SuStr("wVKey"), IntVal(int(x.wVKey)))
	ob.Put(nil, SuStr("flags"), IntVal(int(x.flags)))
	return ob
}

//-------------------------------------------------------------------

var _ = Global.Builtin("LONG",
	&callableStruct{suStructGlobal{size: int(int32Size),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            long}}})

func long(_ *Thread, args []Value) Value {
	n := getInt32(args[0], "x")
	var buf strings.Builder
	binary.Write(&buf, binary.LittleEndian, n)
	return SuStr(buf.String())
}

//-------------------------------------------------------------------

var _ = Global.Builtin("ACCEL",
	&callableStruct{suStructGlobal{size: int(unsafe.Sizeof(stAccel{})),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            accel}}})

type stAccel struct {
	fVirt byte
	pad   byte
	key   int16
	cmd   int16
}

func accel(_ *Thread, args []Value) Value {
	arg := args[0]
	if a, ok := arg.ToInt(); ok {
		// address => ob
		if a == 0 {
			return False
		}
		ac := (*stAccel)(unsafe.Pointer(uintptr(a)))
		ob := &SuObject{}
		ob.Put(nil, SuStr("fVirt"), IntVal(int(ac.fVirt)))
		ob.Put(nil, SuStr("pad"), IntVal(int(ac.pad)))
		ob.Put(nil, SuStr("key"), IntVal(int(ac.key)))
		ob.Put(nil, SuStr("cmd"), IntVal(int(ac.cmd)))
		return ob
	}
	// else ob => string
	ac := stAccel{
		fVirt: byte(getInt(arg, "fVirt")),
		pad:   byte(getInt(arg, "pad")),
		key:   int16(getInt(arg, "key")),
		cmd:   int16(getInt(arg, "cmd")),
	}
	return ptrNstr(unsafe.Pointer(&ac), unsafe.Sizeof(ac))
}

//-------------------------------------------------------------------

type stSCNotification struct {
	nmhdr                stNMHdr
	position             int
	ch                   int32
	modifiers            int32
	modificationType     int32
	text                 uintptr
	length               int
	linesAdded           int
	message              int32
	wParam               uintptr
	lParam               uintptr
	line                 int
	foldLevelNow         int32
	foldLevelPrev        int32
	margin               int32
	listType             int32
	x                    int32
	y                    int32
	token                int32
	annotationLinesAdded int
	updated              int32
	listCompletionMethod int32
}

var _ = builtin(SCNotification, "(address)")

func SCNotification(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	scn := (*stSCNotification)(unsafe.Pointer(uintptr(adr)))
	return fromSCNotification(scn)
}

var _ = builtin(SCNotificationText, "(address)")

func SCNotificationText(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	scn := (*stSCNotification)(unsafe.Pointer(uintptr(adr)))
	ob := fromSCNotification(scn)
	ob.Put(nil, SuStr("text"), ptrZstr(unsafe.Pointer(scn.text), 1024))
	return ob
}

func fromSCNotification(scn *stSCNotification) *SuObject {
	ob := fromNMHdr(&scn.nmhdr)
	ob.Put(nil, SuStr("position"), IntVal(int(scn.position)))
	ob.Put(nil, SuStr("ch"), IntVal(int(scn.ch)))
	ob.Put(nil, SuStr("modifiers"), IntVal(int(scn.modifiers)))
	ob.Put(nil, SuStr("modificationType"), IntVal(int(scn.modificationType)))
	// NOT scn.text
	ob.Put(nil, SuStr("length"), IntVal(int(scn.length)))
	ob.Put(nil, SuStr("linesAdded"), IntVal(int(scn.linesAdded)))
	ob.Put(nil, SuStr("message"), IntVal(int(scn.message)))
	ob.Put(nil, SuStr("wParam"), IntVal(int(scn.wParam)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(scn.lParam)))
	ob.Put(nil, SuStr("line"), IntVal(int(scn.line)))
	ob.Put(nil, SuStr("foldLevelNow"), IntVal(int(scn.foldLevelNow)))
	ob.Put(nil, SuStr("foldLevelPrev"), IntVal(int(scn.foldLevelPrev)))
	ob.Put(nil, SuStr("margin"), IntVal(int(scn.margin)))
	ob.Put(nil, SuStr("listType"), IntVal(int(scn.listType)))
	ob.Put(nil, SuStr("x"), IntVal(int(scn.x)))
	ob.Put(nil, SuStr("y"), IntVal(int(scn.y)))
	ob.Put(nil, SuStr("token"), IntVal(int(scn.token)))
	ob.Put(nil, SuStr("updated"), IntVal(int(scn.updated)))
	return ob
}

//-------------------------------------------------------------------

type stDrawItemStruct struct {
	CtlType    int32
	CtlID      int32
	itemID     int32
	itemAction int32
	itemState  int32
	hwndItem   uintptr
	hDC        uintptr
	rcItem     stRect
	itemData   uintptr
}

var _ = builtin(DRAWITEMSTRUCT, "(address)")

func DRAWITEMSTRUCT(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	dis := (*stDrawItemStruct)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("CtlType"), IntVal(int(dis.CtlType)))
	ob.Put(nil, SuStr("CtlID"), IntVal(int(dis.CtlID)))
	ob.Put(nil, SuStr("itemID"), IntVal(int(dis.itemID)))
	ob.Put(nil, SuStr("itemAction"), IntVal(int(dis.itemAction)))
	ob.Put(nil, SuStr("itemState"), IntVal(int(dis.itemState)))
	ob.Put(nil, SuStr("hwndItem"), IntVal(int(dis.hwndItem)))
	ob.Put(nil, SuStr("hDC"), IntVal(int(dis.hDC)))
	ob.Put(nil, SuStr("rcItem"), fromRect(&dis.rcItem, nil))
	ob.Put(nil, SuStr("itemData"), IntVal(int(dis.itemData)))
	return ob
}

//-------------------------------------------------------------------

var _ = builtin(MSG, "(address)")

func MSG(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	msg := (*stMsg)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("hwnd"), IntVal(int(msg.hwnd)))
	ob.Put(nil, SuStr("message"), IntVal(int(msg.message)))
	ob.Put(nil, SuStr("wParam"), IntVal(int(msg.wParam)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(msg.lParam)))
	ob.Put(nil, SuStr("time"), IntVal(int(msg.time)))
	ob.Put(nil, SuStr("pt"), fromPoint(&msg.pt, nil))
	return ob
}

//-------------------------------------------------------------------

var _ = builtin(CWPRETSTRUCT, "(address)")

func CWPRETSTRUCT(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stCWPRetStruct)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("lResult"), IntVal(int(x.lResult)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(x.lParam)))
	ob.Put(nil, SuStr("wParam"), IntVal(int(x.wParam)))
	ob.Put(nil, SuStr("message"), IntVal(int(x.message)))
	ob.Put(nil, SuStr("hwnd"), IntVal(int(x.hwnd)))
	return ob
}

type stCWPRetStruct struct {
	lResult uintptr
	lParam  uintptr
	wParam  uintptr
	message int32
	hwnd    HANDLE
}

//-------------------------------------------------------------------

var _ = builtin(NMLVDISPINFO, "(address)")

func NMLVDISPINFO(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMLVDispInfo)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("hdr"), fromNMHdr(&x.hdr))
	item := &SuObject{}
	item.Put(nil, SuStr("mask"), IntVal(int(x.item.mask)))
	item.Put(nil, SuStr("iItem"), IntVal(int(x.item.iItem)))
	item.Put(nil, SuStr("iSubItem"), IntVal(int(x.item.iSubItem)))
	item.Put(nil, SuStr("state"), IntVal(int(x.item.state)))
	item.Put(nil, SuStr("stateMask"), IntVal(int(x.item.stateMask)))
	item.Put(nil, SuStr("pszText"),
		IntVal(int(uintptr(unsafe.Pointer(x.item.pszText)))))
	item.Put(nil, SuStr("cchTextMax"), IntVal(int(x.item.cchTextMax)))
	item.Put(nil, SuStr("iImage"), IntVal(int(x.item.iImage)))
	item.Put(nil, SuStr("lParam"), IntVal(int(x.item.lParam)))
	item.Put(nil, SuStr("iIndent"), IntVal(int(x.item.iIndent)))
	ob.Put(nil, SuStr("item"), item)
	return ob
}

type stNMLVDispInfo struct {
	hdr  stNMHdr
	item stLVItem
}

//-------------------------------------------------------------------

var _ = builtin(NMLISTVIEW, "(address)")

func NMLISTVIEW(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMListView)(unsafe.Pointer(uintptr(adr)))
	ob := &SuObject{}
	ob.Put(nil, SuStr("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, SuStr("iItem"), IntVal(int(x.iItem)))
	ob.Put(nil, SuStr("iSubItem"), IntVal(int(x.iSubItem)))
	ob.Put(nil, SuStr("uNewState"), IntVal(int(x.uNewState)))
	ob.Put(nil, SuStr("uOldState"), IntVal(int(x.uOldState)))
	ob.Put(nil, SuStr("uChanged"), IntVal(int(x.uChanged)))
	ob.Put(nil, SuStr("ptAction"), fromPoint(&x.ptAction, nil))
	ob.Put(nil, SuStr("lParam"), IntVal(int(x.lParam)))
	return ob
}

type stNMListView struct {
	hdr       stNMHdr
	iItem     int32
	iSubItem  int32
	uNewState int32
	uOldState int32
	uChanged  int32
	ptAction  stPoint
	lParam    uintptr
}

//-------------------------------------------------------------------

var _ = Global.Builtin("NMDAYSTATE",
	&suNMDAYSTATE{suStructGlobal{size: int(unsafe.Sizeof(stNMDayState{}))}})

type suNMDAYSTATE struct {
	suStructGlobal
}

func (*suNMDAYSTATE) fromStruct(p unsafe.Pointer) Value {
	x := (*stNMDayState)(p)
	ob := &SuObject{}
	ob.Put(nil, SuStr("stStart"), SYSTEMTIMEtoOb(&x.stStart, &SuObject{}))
	ob.Put(nil, SuStr("cDayState"), IntVal(int(x.cDayState)))
	return ob
}

func (*suNMDAYSTATE) updateStruct(ob Value, p unsafe.Pointer) {
	x := (*stNMDayState)(p)
	x.stStart = toSystemTime(ob.Get(nil, SuStr("stStart")))
	x.cDayState = getInt32(ob, "cDayState")
}

type stNMDayState struct {
	nmhdr       stNMHdr
	stStart     stSystemTime
	cDayState   int32
	prgDayState uintptr
}

//-------------------------------------------------------------------

var _ = Global.Builtin("BITMAPFILEHEADER",
	&callableStruct{suStructGlobal{size: int(nBitMapFileHeader),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            bmfh}}})

func bmfh(_ *Thread, args []Value) Value {
	var buf strings.Builder
	binary.Write(&buf, binary.LittleEndian, getInt16(args[0], "bfType"))
	binary.Write(&buf, binary.LittleEndian, getInt32(args[0], "bfSize"))
	binary.Write(&buf, binary.LittleEndian, int32(0)) // reserved
	binary.Write(&buf, binary.LittleEndian, getInt32(args[0], "bfOffBits"))
	assert.That(buf.Len() == nBitMapFileHeader)
	return SuStr(buf.String())
}

// NOTE: 2 byte alignment (no padding) so not compatible as Go struct
// type BITMAPFILEHEADER struct {
// 	bfType      int16
// 	bfSize      int32
// 	bfReserved1 int16
// 	bfReserved2 int16
// 	bfOffBits   int32
// }

const nBitMapFileHeader = 14

//-------------------------------------------------------------------

var _ = Global.Builtin("BITMAPINFOHEADER",
	&callableStruct{suStructGlobal{size: int(nBitMapInfoHeader),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            bmih}}})

func bmih(_ *Thread, args []Value) Value {
	bmih := toBitMapInfoHeader(args[0])
	return ptrNstr(unsafe.Pointer(&bmih), nBitMapInfoHeader)
}

func toBitMapInfoHeader(hdr Value) stBitMapInfoHeader {
	return stBitMapInfoHeader{
		biSize:          int32(nBitMapInfoHeader),
		biWidth:         getInt32(hdr, "biWidth"),
		biHeight:        getInt32(hdr, "biHeight"),
		biPlanes:        getInt16(hdr, "biPlanes"),
		biBitCount:      getInt16(hdr, "biBitCount"),
		biCompression:   getInt32(hdr, "biCompression"),
		biSizeImage:     getInt32(hdr, "biSizeImage"),
		biXPelsPerMeter: getInt32(hdr, "biXPelsPerMeter"),
		biYPelsPerMeter: getInt32(hdr, "biYPelsPerMeter"),
		biClrUsed:       getInt32(hdr, "biClrUsed"),
		biClrImportant:  getInt32(hdr, "biClrImportant"),
	}
}

//-------------------------------------------------------------------

func _() {
	var x [1]struct{}
	// generated by Visual C++ x64 (ascii A versions, not wide W versions)
	_ = x[unsafe.Sizeof(stAccel{})-6]
	_ = x[unsafe.Sizeof(stBitMap{})-32]
	_ = x[unsafe.Sizeof(stBitMapInfoHeader{})-40]
	_ = x[unsafe.Sizeof(stBrowseInfo{})-64]
	_ = x[unsafe.Sizeof(stCharRange{})-8]
	_ = x[unsafe.Sizeof(stChooseColor{})-72]
	_ = x[unsafe.Sizeof(stChooseFont{})-104]
	_ = x[unsafe.Sizeof(stCWPRetStruct{})-40]
	_ = x[unsafe.Sizeof(stDocInfo{})-40]
	_ = x[unsafe.Sizeof(stDrawItemStruct{})-64]
	_ = x[unsafe.Sizeof(stDrawTextParams{})-20]
	_ = x[unsafe.Sizeof(stEditBalloonTip{})-32]
	_ = x[unsafe.Sizeof(stFixed{})-4]
	_ = x[unsafe.Sizeof(stFlashWInfo{})-32]
	_ = x[unsafe.Sizeof(stGlyphMetrics{})-20]
	_ = x[unsafe.Sizeof(GUID{})-16]
	_ = x[unsafe.Sizeof(stHdHitTestInfo{})-16]
	_ = x[unsafe.Sizeof(stHdItem{})-72]
	_ = x[unsafe.Sizeof(stImageInfo{})-40]
	_ = x[unsafe.Sizeof(stInitCommonControlsEx{})-8]
	_ = x[unsafe.Sizeof(stLogBrush{})-16]
	_ = x[unsafe.Sizeof(stLogFont{})-60]
	_ = x[unsafe.Sizeof(stLVColumn{})-56]
	_ = x[unsafe.Sizeof(stLVHitTestInfo{})-24]
	_ = x[unsafe.Sizeof(stLVItem{})-88]
	_ = x[unsafe.Sizeof(stMenuItemInfo{})-80]
	_ = x[unsafe.Sizeof(stMinMaxInfo{})-40]
	_ = x[unsafe.Sizeof(stMsg{})-48]
	_ = x[unsafe.Sizeof(stNMDayState{})-56]
	_ = x[unsafe.Sizeof(stNMHdr{})-24]
	_ = x[unsafe.Sizeof(stNMListView{})-64]
	_ = x[unsafe.Sizeof(stNMLVDispInfo{})-112]
	_ = x[unsafe.Sizeof(stNMTTDispInfo{})-136]
	_ = x[unsafe.Sizeof(stNMTVDispInfo{})-80]
	_ = x[unsafe.Sizeof(stNotifyIconData{})-528]
	_ = x[unsafe.Sizeof(stOpenFileName{})-152]
	_ = x[unsafe.Sizeof(stPageSetupDlg{})-128]
	_ = x[unsafe.Sizeof(stPaintStruct{})-72]
	_ = x[unsafe.Sizeof(stPoint{})-8]
	_ = x[unsafe.Sizeof(stPrintDlg{})-120]
	_ = x[unsafe.Sizeof(stRect{})-16]
	_ = x[unsafe.Sizeof(stSCNotification{})-152]
	_ = x[unsafe.Sizeof(stScrollInfo{})-28]
	_ = x[unsafe.Sizeof(stShellExecuteInfo{})-112]
	_ = x[unsafe.Sizeof(stSystemTime{})-16]
	_ = x[unsafe.Sizeof(stTCItem{})-40]
	_ = x[unsafe.Sizeof(stTextMetric{})-56]
	_ = x[unsafe.Sizeof(stTextRange{})-16]
	_ = x[unsafe.Sizeof(stToolInfo{})-72]
	_ = x[unsafe.Sizeof(stTPMParams{})-20]
	_ = x[unsafe.Sizeof(stTrackMouseEvent{})-24]
	_ = x[unsafe.Sizeof(stTVInsertStruct{})-96]
	_ = x[unsafe.Sizeof(stTVItem{})-56]
	_ = x[unsafe.Sizeof(stWindowPlacement{})-44]
	_ = x[unsafe.Sizeof(stWndClass{})-72]
}
