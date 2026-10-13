// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"encoding/binary"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
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

var _ = method(st_Size, "() :number")

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
var _ = builtin(StructModify, "(type, address, block) :void")

func StructModify(th *Thread, args []Value) Value {
	typ, ok := args[0].(Struct)
	if !ok {
		panic("StructModify invalid type " + ErrType(args[0]))
	}
	p := toptr(ToInt(args[1]))
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
	buf := unsafe.Slice((*byte)(unsafe.Pointer(r)), nRect)
	return StrVal(hacks.BStoS(buf))
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
	ob.Put(nil, StrVal("maxSize"), fromPoint(&mmi.maxSize, nil))
	ob.Put(nil, StrVal("maxPosition"), fromPoint(&mmi.maxPosition, nil))
	ob.Put(nil, StrVal("minTrackSize"), fromPoint(&mmi.minTrackSize, nil))
	ob.Put(nil, StrVal("maxTrackSize"), fromPoint(&mmi.maxTrackSize, nil))
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

var _ = builtin(NMHDR, "(address) :false|object")

func NMHDR(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	nmh := (*stNMHdr)(toptr(adr))
	return fromNMHdr(nmh)
}

func fromNMHdr(nmh *stNMHdr) *SuObject {
	ob := &SuObject{}
	ob.Put(nil, StrVal("hwndFrom"), IntVal(nmh.hwndFrom))
	ob.Put(nil, StrVal("idFrom"), IntVal(nmh.idFrom))
	ob.Put(nil, StrVal("code"), IntVal(nmh.code))
	return ob
}

//-------------------------------------------------------------------

type stNMTVDispInfo struct {
	nmhdr stNMHdr
	item  stTVItem
}

var _ = builtin(NMTVDISPINFO, "(address) :false|object")

func NMTVDISPINFO(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	di := (*stNMTVDispInfo)(toptr(adr))
	ob := fromNMHdr(&di.nmhdr)
	ob.Put(nil, StrVal("nmhdr"), fromNMHdr(&di.nmhdr))
	tvi := fromTVItem(&di.item)
	tvi.Put(nil, StrVal("pszText"), ptrZstr(di.item.pszText, 1024))
	ob.Put(nil, StrVal("item"), tvi)
	return ob
}

func fromTVItem(tvi *stTVItem) *SuObject {
	ob := &SuObject{}
	ob.Put(nil, StrVal("mask"), IntVal(tvi.mask))
	ob.Put(nil, StrVal("hItem"), IntVal(tvi.hItem))
	ob.Put(nil, StrVal("state"), IntVal(tvi.state))
	ob.Put(nil, StrVal("stateMask"), IntVal(tvi.stateMask))
	// pszText must be handled by caller
	ob.Put(nil, StrVal("cchTextMax"), IntVal(tvi.cchTextMax))
	ob.Put(nil, StrVal("iImage"), IntVal(tvi.iImage))
	ob.Put(nil, StrVal("iSelectedImage"), IntVal(tvi.iSelectedImage))
	ob.Put(nil, StrVal("cChildren"), IntVal(tvi.cChildren))
	ob.Put(nil, StrVal("lParam"), IntVal(tvi.lParam))
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
	ob.Put(nil, StrVal("hdr"), fromNMHdr(&x.nmhdr))
	ob.Put(nil, StrVal("item"), fromTVItem(&x.item))
	return ob
}

func (*suNMTVDISPINFO) updateStruct(ob Value, p unsafe.Pointer) {
	tvi := ob.Get(nil, StrVal("item"))
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
	return x.fromStruct(toptr(adr))
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
	ob.Put(nil, StrVal("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, StrVal("szText"), bufZstr(x.szText[:]))
	ob.Put(nil, StrVal("lpszText"), IntVal(x.lpszText))
	ob.Put(nil, StrVal("hinst"), IntVal(x.hinst))
	ob.Put(nil, StrVal("uFlags"), IntVal(x.uFlags))
	ob.Put(nil, StrVal("lParam"), IntVal(x.lParam))
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

var _ = builtin(NMHEADER, "(address) :false|object")

func NMHEADER(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMHeader)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, StrVal("iItem"), IntVal(x.iItem))
	ob.Put(nil, StrVal("iButton"), IntVal(x.iButton))
	if x.pitem != nil {
		hdi := fromHdItem(x.pitem, &SuObject{})
		hdi.Put(nil, StrVal("pszText"),
			IntVal(uintptr(unsafe.Pointer(x.pitem.pszText))))
		ob.Put(nil, StrVal("pitem"), hdi)
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

var _ = builtin(NMTREEVIEW, "(address) :false|object")

func NMTREEVIEW(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMTreeView)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, StrVal("action"), IntVal(x.action))
	ob.Put(nil, StrVal("itemOld"), fromTVItem(&x.itemOld))
	ob.Put(nil, StrVal("itemNew"), fromTVItem(&x.itemNew))
	ob.Put(nil, StrVal("ptDrag"), fromPoint(&x.ptDrag, nil))
	return ob
}

//-------------------------------------------------------------------

type stNMTVKeyDown struct {
	hdr   stNMHdr
	wVKey int16
	flags int32
}

var _ = builtin(NMTVKEYDOWN, "(address) :false|object")

func NMTVKEYDOWN(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMTVKeyDown)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, StrVal("wVKey"), IntVal(x.wVKey))
	ob.Put(nil, StrVal("flags"), IntVal(x.flags))
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
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(n))
	return StrVal(string(buf[:]))
}

//-------------------------------------------------------------------

var _ = Global.Builtin("ACCEL",
	&callableStruct{suStructGlobal{size: int(unsafe.Sizeof(stAccel{})),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            accel}}})

type stAccel struct {
	fVirt byte
	_     byte
	key   int16
	cmd   int16
}

const nAccel = unsafe.Sizeof(stAccel{})

func accel(_ *Thread, args []Value) Value {
	arg := args[0]
	if a, ok := arg.ToInt(); ok {
		// address => ob
		if a == 0 {
			return False
		}
		ac := (*stAccel)(toptr(a))
		ob := &SuObject{}
		ob.Put(nil, StrVal("fVirt"), IntVal(ac.fVirt))
		ob.Put(nil, StrVal("key"), IntVal(ac.key))
		ob.Put(nil, StrVal("cmd"), IntVal(ac.cmd))
		return ob
	}
	// else ob => string
	ac := stAccel{
		fVirt: byte(getInt(arg, "fVirt")),
		key:   int16(getInt(arg, "key")),
		cmd:   int16(getInt(arg, "cmd")),
	}
	buf := make([]byte, nAccel)
	copy(buf, (*(*[nAccel]byte)(unsafe.Pointer(&ac)))[:])
	return StrVal(hacks.BStoS(buf))
}

//-------------------------------------------------------------------

type stSCNotification struct {
	nmhdr                stNMHdr
	position             int
	ch                   int32
	modifiers            int32
	modificationType     int32
	text                 *byte
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

var _ = builtin(SCNotification, "(address) :false|object")

func SCNotification(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	scn := (*stSCNotification)(toptr(adr))
	return fromSCNotification(scn)
}

var _ = builtin(SCNotificationText, "(address) :false|object")

func SCNotificationText(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	scn := (*stSCNotification)(toptr(adr))
	ob := fromSCNotification(scn)
	ob.Put(nil, StrVal("text"), ptrZstr(scn.text, 1024))
	return ob
}

func fromSCNotification(scn *stSCNotification) *SuObject {
	ob := fromNMHdr(&scn.nmhdr)
	ob.Put(nil, StrVal("position"), IntVal(scn.position))
	ob.Put(nil, StrVal("ch"), IntVal(scn.ch))
	ob.Put(nil, StrVal("modifiers"), IntVal(scn.modifiers))
	ob.Put(nil, StrVal("modificationType"), IntVal(scn.modificationType))
	// NOT scn.text
	ob.Put(nil, StrVal("length"), IntVal(scn.length))
	ob.Put(nil, StrVal("linesAdded"), IntVal(scn.linesAdded))
	ob.Put(nil, StrVal("message"), IntVal(scn.message))
	ob.Put(nil, StrVal("wParam"), IntVal(scn.wParam))
	ob.Put(nil, StrVal("lParam"), IntVal(scn.lParam))
	ob.Put(nil, StrVal("line"), IntVal(scn.line))
	ob.Put(nil, StrVal("foldLevelNow"), IntVal(scn.foldLevelNow))
	ob.Put(nil, StrVal("foldLevelPrev"), IntVal(scn.foldLevelPrev))
	ob.Put(nil, StrVal("margin"), IntVal(scn.margin))
	ob.Put(nil, StrVal("listType"), IntVal(scn.listType))
	ob.Put(nil, StrVal("x"), IntVal(scn.x))
	ob.Put(nil, StrVal("y"), IntVal(scn.y))
	ob.Put(nil, StrVal("token"), IntVal(scn.token))
	ob.Put(nil, StrVal("updated"), IntVal(scn.updated))
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

var _ = builtin(DRAWITEMSTRUCT, "(address) :false|object")

func DRAWITEMSTRUCT(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	dis := (*stDrawItemStruct)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("CtlType"), IntVal(dis.CtlType))
	ob.Put(nil, StrVal("CtlID"), IntVal(dis.CtlID))
	ob.Put(nil, StrVal("itemID"), IntVal(dis.itemID))
	ob.Put(nil, StrVal("itemAction"), IntVal(dis.itemAction))
	ob.Put(nil, StrVal("itemState"), IntVal(dis.itemState))
	ob.Put(nil, StrVal("hwndItem"), IntVal(dis.hwndItem))
	ob.Put(nil, StrVal("hDC"), IntVal(dis.hDC))
	ob.Put(nil, StrVal("rcItem"), fromRect(&dis.rcItem, nil))
	ob.Put(nil, StrVal("itemData"), IntVal(dis.itemData))
	return ob
}

//-------------------------------------------------------------------

var _ = builtin(MSG, "(address) :false|object")

func MSG(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	msg := (*stMsg)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("hwnd"), IntVal(msg.hwnd))
	ob.Put(nil, StrVal("message"), IntVal(msg.message))
	ob.Put(nil, StrVal("wParam"), IntVal(msg.wParam))
	ob.Put(nil, StrVal("lParam"), IntVal(msg.lParam))
	ob.Put(nil, StrVal("time"), IntVal(msg.time))
	ob.Put(nil, StrVal("pt"), fromPoint(&msg.pt, nil))
	return ob
}

//-------------------------------------------------------------------

var _ = builtin(CWPRETSTRUCT, "(address) :false|object")

func CWPRETSTRUCT(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stCWPRetStruct)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("lResult"), IntVal(x.lResult))
	ob.Put(nil, StrVal("lParam"), IntVal(x.lParam))
	ob.Put(nil, StrVal("wParam"), IntVal(x.wParam))
	ob.Put(nil, StrVal("message"), IntVal(x.message))
	ob.Put(nil, StrVal("hwnd"), IntVal(x.hwnd))
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

var _ = builtin(NMLVDISPINFO, "(address) :false|object")

func NMLVDISPINFO(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMLVDispInfo)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("hdr"), fromNMHdr(&x.hdr))
	item := &SuObject{}
	item.Put(nil, StrVal("mask"), IntVal(x.item.mask))
	item.Put(nil, StrVal("iItem"), IntVal(x.item.iItem))
	item.Put(nil, StrVal("iSubItem"), IntVal(x.item.iSubItem))
	item.Put(nil, StrVal("state"), IntVal(x.item.state))
	item.Put(nil, StrVal("stateMask"), IntVal(x.item.stateMask))
	item.Put(nil, StrVal("pszText"),
		IntVal(uintptr(unsafe.Pointer(x.item.pszText))))
	item.Put(nil, StrVal("cchTextMax"), IntVal(x.item.cchTextMax))
	item.Put(nil, StrVal("iImage"), IntVal(x.item.iImage))
	item.Put(nil, StrVal("lParam"), IntVal(x.item.lParam))
	item.Put(nil, StrVal("iIndent"), IntVal(x.item.iIndent))
	ob.Put(nil, StrVal("item"), item)
	return ob
}

type stNMLVDispInfo struct {
	hdr  stNMHdr
	item stLVItem
}

//-------------------------------------------------------------------

var _ = builtin(NMLISTVIEW, "(address) :false|object")

func NMLISTVIEW(a Value) Value {
	adr := ToInt(a)
	if adr == 0 {
		return False
	}
	x := (*stNMListView)(toptr(adr))
	ob := &SuObject{}
	ob.Put(nil, StrVal("hdr"), fromNMHdr(&x.hdr))
	ob.Put(nil, StrVal("iItem"), IntVal(x.iItem))
	ob.Put(nil, StrVal("iSubItem"), IntVal(x.iSubItem))
	ob.Put(nil, StrVal("uNewState"), IntVal(x.uNewState))
	ob.Put(nil, StrVal("uOldState"), IntVal(x.uOldState))
	ob.Put(nil, StrVal("uChanged"), IntVal(x.uChanged))
	ob.Put(nil, StrVal("ptAction"), fromPoint(&x.ptAction, nil))
	ob.Put(nil, StrVal("lParam"), IntVal(x.lParam))
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
	ob.Put(nil, StrVal("stStart"), SYSTEMTIMEtoOb(&x.stStart, &SuObject{}))
	ob.Put(nil, StrVal("cDayState"), IntVal(x.cDayState))
	return ob
}

func (*suNMDAYSTATE) updateStruct(ob Value, p unsafe.Pointer) {
	x := (*stNMDayState)(p)
	x.stStart = toSystemTime(ob.Get(nil, StrVal("stStart")))
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
	var buf [nBitMapFileHeader]byte
	binary.LittleEndian.PutUint16(buf[0:], uint16(getInt16(args[0], "bfType")))
	binary.LittleEndian.PutUint32(buf[2:], uint32(getInt32(args[0], "bfSize")))
	// buf[6:10] reserved = 0
	binary.LittleEndian.PutUint32(buf[10:], uint32(getInt32(args[0], "bfOffBits")))
	return StrVal(string(buf[:]))
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
	buf := unsafe.Slice((*byte)(unsafe.Pointer(bmih)), nBitMapInfoHeader)
	return StrVal(hacks.BStoS(buf))
}

func toBitMapInfoHeader(hdr Value) *stBitMapInfoHeader {
	return &stBitMapInfoHeader{
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
