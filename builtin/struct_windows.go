package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
)

type SuStructGlobal struct {
	SuBuiltin
	size int
}

func init() {
	Global.Builtin("INITCOMMONCONTROLSEX",
		&SuStructGlobal{size: int(unsafe.Sizeof(INITCOMMONCONTROLSEX{}))})
	Global.Builtin("MENUITEMINFO",
		&SuStructGlobal{size: int(unsafe.Sizeof(MENUITEMINFO{}))})
	Global.Builtin("MONITORINFO",
		&SuStructGlobal{size: int(unsafe.Sizeof(MONITORINFO{}))})
	Global.Builtin("SCROLLINFO",
		&SuStructGlobal{size: int(unsafe.Sizeof(SCROLLINFO{}))})
	Global.Builtin("TRACKMOUSEEVENT",
		&SuStructGlobal{size: int(unsafe.Sizeof(TRACKMOUSEEVENT{}))})
	Global.Builtin("OSVERSIONINFOEX",
		&SuStructGlobal{size: int(unsafe.Sizeof(OSVERSIONINFOEX{}))})
	Global.Builtin("SHELLEXECUTEINFO",
		&SuStructGlobal{size: int(unsafe.Sizeof(SHELLEXECUTEINFO{}))})
	Global.Builtin("TOOLINFO",
		&SuStructGlobal{size: int(unsafe.Sizeof(TOOLINFO{}))})
	Global.Builtin("WINDOWPLACEMENT",
		&SuStructGlobal{size: int(unsafe.Sizeof(WINDOWPLACEMENT{}))})
	Global.Builtin("OPENFILENAME",
		&SuStructGlobal{size: int(nOPENFILENAME)})
	Global.Builtin("TPMPARAMS",
		&SuStructGlobal{size: int(nTPMPARAMS)})
	Global.Builtin("DRAWTEXTPARAMS",
		&SuStructGlobal{size: int(nDRAWTEXTPARAMS)})
}

func (*SuStructGlobal) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call struct")
}

type callableStruct struct {
	SuStructGlobal
}

func (cs *callableStruct) Call(t *Thread, this Value, as *ArgSpec) Value {
	return cs.SuBuiltin.Call(t, this, as)
}

func (sg *SuStructGlobal) Size() int {
	return sg.size
}

type sizeable interface{ Size() int }

var structSize = method0(func(this Value) Value {
	return IntVal(this.(sizeable).Size())
})

func (*SuStructGlobal) Lookup(_ *Thread, method string) Callable {
	if method == "Size" {
		return structSize
	}
	return nil
}

func (*SuStructGlobal) String() string {
	return "/* builtin struct */"
}

//-------------------------------------------------------------------

type Struct interface {
	structToOb(p unsafe.Pointer) Value
	updateStruct(ob Value, p unsafe.Pointer)
}

// StructModify allows modifying a struct at a given address
// WARNING: address must be valid for type
var _ = builtin("StructModify(type, address, block)",
	func(t *Thread, args []Value) Value {
		typ, ok := args[0].(Struct)
		if !ok {
			panic("StructModify invalid type " + ErrType(args[0]))
		}
		p := unsafe.Pointer(uintptr(ToInt(args[1])))
		if p == nil {
			panic("StructModify: address cannot be zero")
		}
		ob := typ.structToOb(p)
		t.Call(args[2], ob) // call the block, which modifies ob
		typ.updateStruct(ob, p)
		return nil
	})

//-------------------------------------------------------------------

var _ = Global.Builtin("RECT",
	&SuRect{callableStruct{SuStructGlobal{size: int(unsafe.Sizeof(RECT{})),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            rect}}}})

type SuRect struct {
	callableStruct
}

func (*SuRect) structToOb(p unsafe.Pointer) Value {
	return urectToOb(p, nil)
}

func (*SuRect) updateStruct(ob Value, p unsafe.Pointer) {
	*(*RECT)(p) = obToRect(ob)
}

func rect(_ *Thread, args []Value) Value {
	r := obToRect(args[0])
	return bufRet(unsafe.Pointer(&r), nRECT)
}

//-------------------------------------------------------------------

var _ = Global.Builtin("MINMAXINFO",
	&SuMinMaxInfo{SuStructGlobal{size: int(unsafe.Sizeof(MINMAXINFO{}))}})

type SuMinMaxInfo struct {
	SuStructGlobal
}

type MINMAXINFO struct {
	ptReserved   POINT
	maxSize      POINT
	maxPosition  POINT
	minTrackSize POINT
	maxTrackSize POINT
}

func (typ *SuMinMaxInfo) structToOb(p unsafe.Pointer) Value {
	mmi := (*MINMAXINFO)(p)
	ob := NewSuObject()
	ob.Put(nil, SuStr("maxSize"), pointToOb(&mmi.maxSize, nil))
	ob.Put(nil, SuStr("maxPosition"), pointToOb(&mmi.maxPosition, nil))
	ob.Put(nil, SuStr("minTrackSize"), pointToOb(&mmi.minTrackSize, nil))
	ob.Put(nil, SuStr("maxTrackSize"), pointToOb(&mmi.maxTrackSize, nil))
	return ob
}

func (typ *SuMinMaxInfo) updateStruct(ob Value, p unsafe.Pointer) {
	mmi := (*MINMAXINFO)(p)
	mmi.maxSize = getPoint(ob, "maxSize")
	mmi.maxPosition = getPoint(ob, "maxPosition")
	mmi.minTrackSize = getPoint(ob, "minTrackSize")
	mmi.maxTrackSize = getPoint(ob, "maxTrackSize")
}

//-------------------------------------------------------------------

type NMHDR struct {
	hwndFrom uintptr
	idFrom   uintptr
	code     int32
	_        [4]byte // padding
}

var _ = builtin1("NMHDR(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		nmh := (*NMHDR)(unsafe.Pointer(uintptr(adr)))
		return nmhdrToOb(nmh)
	})

func nmhdrToOb(nmh *NMHDR) *SuObject {
	ob := NewSuObject()
	ob.Put(nil, SuStr("hwndFrom"), IntVal(int(nmh.hwndFrom)))
	ob.Put(nil, SuStr("idFrom"), IntVal(int(nmh.idFrom)))
	ob.Put(nil, SuStr("code"), IntVal(int(nmh.code)))
	return ob
}

//-------------------------------------------------------------------

type NMTVDISPINFO struct {
	nmhdr NMHDR
	item  TVITEM
}

var _ = builtin1("NMTVDISPINFO(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		di := (*NMTVDISPINFO)(unsafe.Pointer(uintptr(adr)))
		ob := nmhdrToOb(&di.nmhdr)
		ob.Put(nil, SuStr("nmhdr"), nmhdrToOb(&di.nmhdr))
		tvi := tvitemToOb(&di.item)
		tvi.Put(nil, SuStr("pszText"),
			bufToStr(unsafe.Pointer(di.item.pszText), 1024))
		ob.Put(nil, SuStr("item"), tvi)
		return ob
	})

func tvitemToOb(tvi *TVITEM) *SuObject {
	ob := NewSuObject()
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
	&SuNMTVDISPINFO{callableStruct{SuStructGlobal{size: int(unsafe.Sizeof(NMTVDISPINFO{})),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            nmtvdispinfo}}}})

type SuNMTVDISPINFO struct {
	callableStruct
}

func (*SuNMTVDISPINFO) structToOb(p unsafe.Pointer) Value {
	x := (*NMTVDISPINFO)(p)
	ob := NewSuObject()
	ob.Put(nil, SuStr("hdr"), nmhdrToOb(&x.nmhdr))
	ob.Put(nil, SuStr("item"), tvitemToOb(&x.item))
	return ob
}

func (*SuNMTVDISPINFO) updateStruct(ob Value, p unsafe.Pointer) {
	tvi := ob.Get(nil, SuStr("item"))
	x := (*NMTVDISPINFO)(p)
	x.item.mask = getUint32(tvi, "mask")
	x.item.hItem = getHandle(tvi, "hItem")
	x.item.state = getUint32(tvi, "state")
	x.item.stateMask = getUint32(tvi, "stateMask")
	x.item.cchTextMax = getInt32(tvi, "cchTextMax")
	x.item.iImage = getInt32(tvi, "iImage")
	x.item.iSelectedImage = getInt32(tvi, "iSelectedImage")
	x.item.cChildren = getInt32(tvi, "cChildren")
	x.item.lParam = getHandle(tvi, "lParam")
}

func nmtvdispinfo(_ *Thread, args []Value) Value {
	adr := ToInt(args[0])
	if adr == 0 {
		return False
	}
	var x *SuNMTVDISPINFO
	return x.structToOb(unsafe.Pointer(uintptr(adr)))
}

//-------------------------------------------------------------------

var _ = Global.Builtin("NMTTDISPINFO2",
	&SuNMTTDISPINFO{SuStructGlobal{size: int(unsafe.Sizeof(NMTTDISPINFO{}))}})

type SuNMTTDISPINFO struct {
	SuStructGlobal
}

func (*SuNMTTDISPINFO) structToOb(p unsafe.Pointer) Value {
	x := (*NMTTDISPINFO)(p)
	ob := NewSuObject()
	ob.Put(nil, SuStr("hdr"), nmhdrToOb(&x.hdr))
	ob.Put(nil, SuStr("szText"), strRet(x.szText[:]))
	ob.Put(nil, SuStr("lpszText"), IntVal(int(x.lpszText)))
	ob.Put(nil, SuStr("hinst"), IntVal(int(x.hinst)))
	ob.Put(nil, SuStr("uFlags"), IntVal(int(x.uFlags)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(x.lParam)))
	return ob
}

func (*SuNMTTDISPINFO) updateStruct(ob Value, p unsafe.Pointer) {
	x := (*NMTTDISPINFO)(p)
	x.lpszText = getHandle(ob, "lpszText")
	copyStr(x.szText[:], ob, "szText")
	x.hinst = getHandle(ob, "hinst")
	x.uFlags = getInt32(ob, "uFlags")
	x.lParam = getHandle(ob, "lParam")
}

type NMTTDISPINFO struct {
	hdr      NMHDR
	lpszText uintptr
	szText   [80]byte
	hinst    uintptr
	uFlags   int32
	lParam   uintptr
}

//-------------------------------------------------------------------

type NMHEADER struct {
	hdr     NMHDR
	iItem   int32
	iButton int32
	pitem   *HDITEM
}

var _ = builtin1("NMHEADER(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		x := (*NMHEADER)(unsafe.Pointer(uintptr(adr)))
		ob := NewSuObject()
		ob.Put(nil, SuStr("hdr"), nmhdrToOb(&x.hdr))
		ob.Put(nil, SuStr("iItem"), IntVal(int(x.iItem)))
		ob.Put(nil, SuStr("iButton"), IntVal(int(x.iButton)))
		if x.pitem != nil {
			hdi := hditemToOb(x.pitem, NewSuObject())
			hdi.Put(nil, SuStr("pszText"),
				IntVal(int(uintptr(unsafe.Pointer(x.pitem.pszText)))))
			ob.Put(nil, SuStr("pitem"), hdi)
		}
		return ob
	})

//-------------------------------------------------------------------

type NMTREEVIEW struct {
	hdr     NMHDR
	action  int32
	itemOld TVITEM
	itemNew TVITEM
	ptDrag  POINT
}

var _ = builtin1("NM_TREEVIEW(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		x := (*NMTREEVIEW)(unsafe.Pointer(uintptr(adr)))
		ob := NewSuObject()
		ob.Put(nil, SuStr("hdr"), nmhdrToOb(&x.hdr))
		ob.Put(nil, SuStr("action"), IntVal(int(x.action)))
		ob.Put(nil, SuStr("itemOld"), tvitemToOb(&x.itemOld))
		ob.Put(nil, SuStr("itemNew"), tvitemToOb(&x.itemNew))
		ob.Put(nil, SuStr("ptDrag"), pointToOb(&x.ptDrag, nil))
		return ob
	})

//-------------------------------------------------------------------

type NMTVKEYDOWN struct {
	hdr   NMHDR
	wVKey int16
	flags int32
}

var _ = builtin1("NMTVKEYDOWN(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		x := (*NMTVKEYDOWN)(unsafe.Pointer(uintptr(adr)))
		ob := NewSuObject()
		ob.Put(nil, SuStr("hdr"), nmhdrToOb(&x.hdr))
		ob.Put(nil, SuStr("wVKey"), IntVal(int(x.wVKey)))
		ob.Put(nil, SuStr("flags"), IntVal(int(x.flags)))
		return ob
	})

//-------------------------------------------------------------------

var _ = Global.Builtin("ACCEL",
	&SuAccel{callableStruct{SuStructGlobal{size: int(unsafe.Sizeof(ACCEL{})),
		SuBuiltin: SuBuiltin{
			BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
			Fn:            accel}}}})

type SuAccel struct {
	callableStruct
}

type ACCEL struct {
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
		ac := (*ACCEL)(unsafe.Pointer(uintptr(a)))
		ob := NewSuObject()
		ob.Put(nil, SuStr("fVirt"), IntVal(int(ac.fVirt)))
		ob.Put(nil, SuStr("pad"), IntVal(int(ac.pad)))
		ob.Put(nil, SuStr("key"), IntVal(int(ac.key)))
		ob.Put(nil, SuStr("cmd"), IntVal(int(ac.cmd)))
		return ob
	}
	// else ob => string
	ac := ACCEL{
		fVirt: byte(getInt(arg, "fVirt")),
		pad:   byte(getInt(arg, "pad")),
		key:   int16(getInt(arg, "key")),
		cmd:   int16(getInt(arg, "cmd")),
	}
	return bufRet(unsafe.Pointer(&ac), unsafe.Sizeof(ac))
}

//-------------------------------------------------------------------

type SCNotification struct {
	nmhdr                NMHDR
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

var _ = builtin1("SCNotification(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		scn := (*SCNotification)(unsafe.Pointer(uintptr(adr)))
		return scnToOb(scn)
	})

var _ = builtin1("SCNotificationText(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		scn := (*SCNotification)(unsafe.Pointer(uintptr(adr)))
		ob := scnToOb(scn)
		ob.Put(nil, SuStr("text"), bufRet(unsafe.Pointer(scn.text), 1024))
		return ob
	})

func scnToOb(scn *SCNotification) *SuObject {
	ob := nmhdrToOb(&scn.nmhdr)
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
	return ob
}

//-------------------------------------------------------------------

type DRAWITEMSTRUCT struct {
	CtlType    int32
	CtlID      int32
	itemID     int32
	itemAction int32
	itemState  int32
	hwndItem   uintptr
	hDC        uintptr
	rcItem     RECT
	itemData   uintptr
}

var _ = builtin1("DRAWITEMSTRUCT(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		dis := (*DRAWITEMSTRUCT)(unsafe.Pointer(uintptr(adr)))
		ob := NewSuObject()
		ob.Put(nil, SuStr("CtlType"), IntVal(int(dis.CtlType)))
		ob.Put(nil, SuStr("CtlID"), IntVal(int(dis.CtlID)))
		ob.Put(nil, SuStr("itemID"), IntVal(int(dis.itemID)))
		ob.Put(nil, SuStr("itemAction"), IntVal(int(dis.itemAction)))
		ob.Put(nil, SuStr("itemState"), IntVal(int(dis.itemState)))
		ob.Put(nil, SuStr("hwndItem"), IntVal(int(dis.hwndItem)))
		ob.Put(nil, SuStr("hDC"), IntVal(int(dis.hDC)))
		ob.Put(nil, SuStr("rcItem"), rectToOb(&dis.rcItem, nil))
		ob.Put(nil, SuStr("itemData"), IntVal(int(dis.itemData)))
		return ob
	})

//-------------------------------------------------------------------

var _ = builtin1("MSG(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		msg := (*MSG)(unsafe.Pointer(uintptr(adr)))
		ob := NewSuObject()
		ob.Put(nil, SuStr("hwnd"), IntVal(int(msg.hwnd)))
		ob.Put(nil, SuStr("message"), IntVal(int(msg.message)))
		ob.Put(nil, SuStr("wParam"), IntVal(int(msg.wParam)))
		ob.Put(nil, SuStr("lParam"), IntVal(int(msg.lParam)))
		ob.Put(nil, SuStr("time"), IntVal(int(msg.time)))
		ob.Put(nil, SuStr("pt"), pointToOb(&msg.pt, nil))
		return ob
	})

//-------------------------------------------------------------------

var _ = builtin1("CWPRETSTRUCT(address)",
	func(a Value) Value {
		adr := ToInt(a)
		if adr == 0 {
			return False
		}
		x := (*CWPRETSTRUCT)(unsafe.Pointer(uintptr(adr)))
		ob := NewSuObject()
		ob.Put(nil, SuStr("lResult"), IntVal(int(x.lResult)))
		ob.Put(nil, SuStr("lParam"), IntVal(int(x.lParam)))
		ob.Put(nil, SuStr("wParam"), IntVal(int(x.wParam)))
		ob.Put(nil, SuStr("message"), IntVal(int(x.message)))
		ob.Put(nil, SuStr("hwnd"), IntVal(int(x.hwnd)))
		return ob
	})

type CWPRETSTRUCT struct {
	lResult uintptr
	lParam  uintptr
	wParam  uintptr
	message int32
	hwnd    HANDLE
}

//-------------------------------------------------------------------

func _() {
	var x [1]struct{}
	// generated by Visual C++ x64 (ascii A versions, not wide W versions)
	_ = x[unsafe.Sizeof(RECT{})-16]
	_ = x[unsafe.Sizeof(PAINTSTRUCT{})-72]
	_ = x[unsafe.Sizeof(SCROLLINFO{})-28]
	_ = x[unsafe.Sizeof(WINDOWPLACEMENT{})-44]
	_ = x[unsafe.Sizeof(MENUITEMINFO{})-80]
	_ = x[unsafe.Sizeof(WNDCLASS{})-72]
	_ = x[unsafe.Sizeof(TCITEM{})-40]
	_ = x[unsafe.Sizeof(CHARRANGE{})-8]
	_ = x[unsafe.Sizeof(TEXTRANGE{})-16]
	_ = x[unsafe.Sizeof(TOOLINFO{})-72]
	_ = x[unsafe.Sizeof(TVITEM{})-56]
	_ = x[unsafe.Sizeof(NMTVDISPINFO{})-80]
	_ = x[unsafe.Sizeof(TV_INSERTSTRUCT{})-96]
	_ = x[unsafe.Sizeof(TPMPARAMS{})-20]
	_ = x[unsafe.Sizeof(DRAWTEXTPARAMS{})-20]
	_ = x[unsafe.Sizeof(TRACKMOUSEEVENT{})-24]
	_ = x[unsafe.Sizeof(FLASHWINFO{})-32]
	_ = x[unsafe.Sizeof(BITMAPINFOHEADER{})-40]
	_ = x[unsafe.Sizeof(INITCOMMONCONTROLSEX{})-8]
	_ = x[unsafe.Sizeof(IMAGEINFO{})-40]
	_ = x[unsafe.Sizeof(HDITEM{})-72]
	_ = x[unsafe.Sizeof(HDHITTESTINFO{})-16]
	_ = x[unsafe.Sizeof(PRINTDLG{})-120]
	_ = x[unsafe.Sizeof(PAGESETUPDLG{})-128]
	_ = x[unsafe.Sizeof(OPENFILENAME{})-152]
	_ = x[unsafe.Sizeof(CHOOSECOLOR{})-72]
	_ = x[unsafe.Sizeof(CHOOSEFONT{})-104]
	_ = x[unsafe.Sizeof(LOGFONT{})-60]
	_ = x[unsafe.Sizeof(TEXTMETRIC{})-56]
	_ = x[unsafe.Sizeof(LOGBRUSH{})-16]
	_ = x[unsafe.Sizeof(GLYPHMETRICS{})-20]
	_ = x[unsafe.Sizeof(FIXED{})-4]
	_ = x[unsafe.Sizeof(DOCINFO{})-40]
	_ = x[unsafe.Sizeof(BITMAP{})-32]
	_ = x[unsafe.Sizeof(OSVERSIONINFOEX{})-156]
	_ = x[unsafe.Sizeof(MSG{})-48]
	_ = x[unsafe.Sizeof(GUID{})-16]
	_ = x[unsafe.Sizeof(NOTIFYICONDATA{})-528]
	_ = x[unsafe.Sizeof(SHELLEXECUTEINFO{})-112]
	_ = x[unsafe.Sizeof(BROWSEINFO{})-64]
	_ = x[unsafe.Sizeof(MINMAXINFO{})-40]
	_ = x[unsafe.Sizeof(NMHDR{})-24]
	_ = x[unsafe.Sizeof(ACCEL{})-6]
	_ = x[unsafe.Sizeof(DRAWITEMSTRUCT{})-64]
	_ = x[unsafe.Sizeof(SCNotification{})-152]
}
