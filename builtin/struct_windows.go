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
		
	Global.Builtin("MINMAXINFO",
		&SuMinMaxInfo{SuStructGlobal{size: int(unsafe.Sizeof(MINMAXINFO{}))}})
	Global.Builtin("RECT",
		&SuRect{SuStructGlobal{size: int(unsafe.Sizeof(RECT{}))}})
	Global.Builtin("ACCEL",
		&SuAccel{SuStructGlobal{size: int(unsafe.Sizeof(ACCEL{})),
			SuBuiltin: SuBuiltin{
				BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1},
				Fn: func(_ *Thread, args []Value) Value {
					return accel(args[0])
				}}}})
}

func (*SuStructGlobal) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call struct")
}

func (sg *SuStructGlobal) Size() int {
	return sg.size
}

type sizeable interface { Size() int }

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
	SuStructGlobal
}

func (typ *SuRect) structToOb(p unsafe.Pointer) Value {
	return rectToOb((*RECT)(p), nil)
}

func (typ *SuRect) updateStruct(ob Value, p unsafe.Pointer) {
	*(*RECT)(p) = obToRect(ob)
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

type SuAccel struct {
	SuStructGlobal
}

type ACCEL struct {
	fVirt byte
	pad   byte
	key   int16
	cmd   int16
}

func accel(arg Value) Value {
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
		pad: byte(getInt(arg, "pad")),
		key: int16(getInt(arg, "key")),
		cmd: int16(getInt(arg, "cmd")),
	}
	return SuStr(memToStr(unsafe.Pointer(&ac), unsafe.Sizeof(ac)))
}

func memToStr(p unsafe.Pointer, n uintptr) string {
	var sb strings.Builder
	for ; n > 0; n--{
		sb.WriteByte(*(*byte)(unsafe.Pointer(uintptr(p) + n)))
	}
	return sb.String()
}

//-------------------------------------------------------------------

type SCNotification struct {
	nmhdr            NMHDR
	position         int32
	ch               int32
	modifiers        int32
	modificationType int32
	text             uintptr
	length           int32
	linesAdded       int32
	message          int32
	wParam           uintptr
	lParam           uintptr
	line             int32
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
		ob.Put(nil, SuStr("text"), SuStr(strFromAddr(scn.text)))
		return ob
	})

func scnToOb(scn *SCNotification) *SuObject {
	ob := nmhdrToOb(&scn.nmhdr)
	ob.Put(nil, SuStr("position"), IntVal(int(scn.position)))
	ob.Put(nil, SuStr("ch"), IntVal(int(scn.ch)))
	ob.Put(nil, SuStr("modifiers"), IntVal(int(scn.modifiers)))
	ob.Put(nil, SuStr("modificationType"), IntVal(int(scn.modificationType)))
	// NOT text
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

func strFromAddr(a uintptr) string {
	var sb strings.Builder
	for ; ; a++ {
		c := *(*byte)(unsafe.Pointer(a))
		if c == 0 {
			break
		}
		sb.WriteByte(c)
	}
	return sb.String()
}
