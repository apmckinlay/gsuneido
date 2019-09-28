package builtin

import (
	"strings"
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

	Global.Builtin("MINMAXINFO",
		&SuMinMaxInfo{SuStructGlobal{size: int(unsafe.Sizeof(MINMAXINFO{}))}})
	Global.Builtin("RECT",
		&SuRect{callableStruct{SuStructGlobal{size: int(unsafe.Sizeof(RECT{})),
			SuBuiltin: SuBuiltin{
				BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
				Fn:            rect}}}})
	Global.Builtin("ACCEL",
		&SuAccel{callableStruct{SuStructGlobal{size: int(unsafe.Sizeof(ACCEL{})),
			SuBuiltin: SuBuiltin{
				BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
				Fn:            accel}}}})
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

//===================================================================

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
		ob := typ.structToOb(p)
		t.Call(args[2], ob) // call the block, which modifies ob
		typ.updateStruct(ob, p)
		return nil
	})

//-------------------------------------------------------------------

type SuRect struct {
	callableStruct
}

func (*SuRect) structToOb(p unsafe.Pointer) Value {
	return rectToOb((*RECT)(p), nil)
}

func (*SuRect) updateStruct(ob Value, p unsafe.Pointer) {
	*(*RECT)(p) = obToRect(ob)
}

func rect(_ *Thread, args []Value) Value {
	r := obToRect(args[0])
	return SuStr(memToStr(uintptr(unsafe.Pointer(&r)), unsafe.Sizeof(r)))
}

//-------------------------------------------------------------------

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
		nmh := (*NMHDR)(unsafe.Pointer(uintptr(ToInt(a))))
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

type TV_DISPINFO struct {
	nmhdr NMHDR
	item  TVITEM
}

var _ = builtin1("TV_DISPINFO(address)",
	func(a Value) Value {
		di := (*TV_DISPINFO)(unsafe.Pointer(uintptr(ToInt(a))))
		ob := nmhdrToOb(&di.nmhdr)
		ob.Put(nil, SuStr("nmhdr"), nmhdrToOb(&di.nmhdr))
		ob.Put(nil, SuStr("item"), tvitemToOb(&di.item))
		return ob
	})

func tvitemToOb(tvi *TVITEM) *SuObject {
	ob := NewSuObject()
	ob.Put(nil, SuStr("mask"), IntVal(int(tvi.mask)))
	ob.Put(nil, SuStr("hItem"), IntVal(int(tvi.hItem)))
	ob.Put(nil, SuStr("state"), IntVal(int(tvi.state)))
	ob.Put(nil, SuStr("stateMask"), IntVal(int(tvi.stateMask)))
	ob.Put(nil, SuStr("pszText"), strFromAddr(uintptr(unsafe.Pointer(tvi.pszText))))
	ob.Put(nil, SuStr("cchTextMax"), IntVal(int(tvi.cchTextMax)))
	ob.Put(nil, SuStr("iImage"), IntVal(int(tvi.iImage)))
	ob.Put(nil, SuStr("iSelectedImage"), IntVal(int(tvi.iSelectedImage)))
	ob.Put(nil, SuStr("cChildren"), IntVal(int(tvi.cChildren)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(tvi.lParam)))
	return ob
}

//-------------------------------------------------------------------

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
	return SuStr(memToStr(uintptr(unsafe.Pointer(&ac)), unsafe.Sizeof(ac)))
}

func memToStr(p uintptr, n uintptr) string {
	var sb strings.Builder
	for i := uintptr(0); i < n; i++ {
		sb.WriteByte(*(*byte)(unsafe.Pointer(p + i)))
	}
	return sb.String()
}

//-------------------------------------------------------------------

type SCNotification struct {
	nmhdr            NMHDR
	position         int
	ch               int32
	modifiers        int32
	modificationType int32
	text             uintptr
	length           int
	linesAdded       int
	message          int32
	wParam           uintptr
	lParam           uintptr
	line             int
	foldLevelNow     int32
	foldLevelPrev    int32
	margin           int32
	listType         int32
	x                int32
	y                int32
	token            int32
}

var _ = builtin1("SCNotification(address)",
	func(a Value) Value {
		scn := (*SCNotification)(unsafe.Pointer(uintptr(ToInt(a))))
		return scnToOb(scn)
	})

var _ = builtin1("SCNotificationText(address)",
	func(a Value) Value {
		scn := (*SCNotification)(unsafe.Pointer(uintptr(ToInt(a))))
		ob := scnToOb(scn)
		ob.Put(nil, SuStr("text"), strFromAddr(scn.text))
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

func strFromAddr(a uintptr) Value {
	if a == 0 {
		return False
	}
	var sb strings.Builder
	for ; ; a++ {
		c := *(*byte)(unsafe.Pointer(a))
		if c == 0 {
			break
		}
		sb.WriteByte(c)
	}
	return SuStr(sb.String())
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
		dis := (*DRAWITEMSTRUCT)(unsafe.Pointer(uintptr(ToInt(a))))
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
		msg := (*MSG)(unsafe.Pointer(uintptr(ToInt(a))))
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
	_ = x[unsafe.Sizeof(TV_DISPINFO{})-80]
	_ = x[unsafe.Sizeof(TV_INSERTSTRUCT{})-96]
	_ = x[unsafe.Sizeof(TPMPARAMS{})-20]
	_ = x[unsafe.Sizeof(DRAWTEXTPARAMS{})-20]
	_ = x[unsafe.Sizeof(TRACKMOUSEEVENT{})-24]
	_ = x[unsafe.Sizeof(FLASHWINFO{})-32]
	_ = x[unsafe.Sizeof(BITMAPINFOHEADER{})-40]
	_ = x[unsafe.Sizeof(INITCOMMONCONTROLSEX{})-8]
	_ = x[unsafe.Sizeof(IMAGEINFO{})-40]
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
}
