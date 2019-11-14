package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		pointer lParam) pointer
var sendMessage = user32.MustFindProc("SendMessageA").Addr()
var _ = builtin4("SendMessage(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return intRet(rtn)
	})

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		string text) pointer
var _ = builtin4("SendMessageText(hwnd, msg, wParam, text)",
	sendMessageText)

func sendMessageText(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(stringArg(d)))
	return intRet(rtn)
}

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
//		[in] string text) pointer
var _ = builtin4("SendMessageTextIn(hwnd, msg, wParam, text)",
	// can't pass direct pointer so same as SendMessageText (i.e. still copies)
	sendMessageText)

var _ = builtin4("SendMessageTextOut(hwnd, msg, wParam = 0, bufsize = 1024)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		n := uintptr(ToInt(d) + 1)
		p := heap.Alloc(n)
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		ob := NewSuObject()
		ob.Put(nil, SuStr("text"), bufToStr(p, n))
		ob.Put(nil, SuStr("result"), intRet(rtn))
		return ob
	})

var _ = builtin4("SendMessagePoint(hwnd, msg, wParam, point)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		pt := heap.Alloc(nPOINT)
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(pointArg(d, pt)))
		upointToOb(pt, d)
		return intRet(rtn)
	})

var _ = builtin4("SendMessageRect(hwnd, msg, wParam, rect)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(rectArg(d, r)))
		urectToOb(r, d)
		return intRet(rtn)
	})

var _ = builtin4("SendMessageTcitem(hwnd, msg, wParam, tcitem)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		t := getStr(d, "pszText")
		n := getInt32(d, "cchTextMax")
		if n > 0 {
			t = (*byte)(heap.Alloc(uintptr(n)))
		}
		p := heap.Alloc(nTCITEM)
		*(*TCITEM)(p) = TCITEM{
			mask:        getUint32(d, "mask"),
			dwState:     getUint32(d, "dwState"),
			dwStateMask: getUint32(d, "dwStateMask"),
			pszText:     t,
			cchTextMax:  n,
			iImage:      getInt32(d, "iImage"),
			lParam:      getInt32(d, "lParam"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		if n > 0 {
			d.Put(nil, SuStr("pszText"), bufToStr(unsafe.Pointer(t), uintptr(n)))
		}
		d.Put(nil, SuStr("iImage"), IntVal(int((*TCITEM)(p).iImage)))
		return intRet(rtn)
	})

type TCITEM struct {
	mask        uint32
	dwState     uint32
	dwStateMask uint32
	pszText     *byte
	cchTextMax  int32
	iImage      int32
	lParam      int32
	_           [4]byte // padding
}

const nTCITEM = unsafe.Sizeof(TCITEM{})

var _ = builtin5("SendMessageTextRange(hwnd, msg, cpMin, cpMax, each = 1)",
	func(a, b, c, d, e Value) Value {
		defer heap.FreeTo(heap.CurSize())
		cpMin := ToInt(c)
		cpMax := ToInt(d)
		if cpMax <= cpMin {
			return EmptyStr
		}
		each := uintptr(ToInt(e))
		n := uintptr(cpMax-cpMin) * each
		buf := heap.Alloc(n + each)
		p := heap.Alloc(nTEXTRANGE)
		*(*TEXTRANGE)(p) = TEXTRANGE{
			chrg:      CHARRANGE{cpMin: int32(cpMin), cpMax: int32(cpMax)},
			lpstrText: (*byte)(buf),
		}
		goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			0,
			uintptr(p))
		return bufRet(buf, n)
	})

var _ = builtin4("SendMessageTOOLINFO(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nTOOLINFO)
		*(*TOOLINFO)(p) = TOOLINFO{
			cbSize:   uint32(nTOOLINFO),
			uFlags:   getUint32(d, "uFlags"),
			hwnd:     getHandle(d, "hwnd"),
			uId:      getUint32(d, "uId"),
			hinst:    getHandle(d, "hinst"),
			lpszText: getStr(d, "lpszText"),
			lParam:   getInt32(d, "lParam"),
			rect:     getRect(d, "rect"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		return intRet(rtn)
	})

var _ = builtin4("SendMessageTOOLINFO2(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nTOOLINFO2)
		*(*TOOLINFO2)(p) = TOOLINFO2{
			cbSize:   uint32(nTOOLINFO2),
			uFlags:   getUint32(d, "uFlags"),
			hwnd:     getHandle(d, "hwnd"),
			uId:      getUint32(d, "uId"),
			hinst:    getHandle(d, "hinst"),
			lpszText: uintptr(getInt(d, "lpszText")), // for LPSTR_TEXTCALLBACK
			lParam:   getInt32(d, "lParam"),
			rect:     getRect(d, "rect"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		return intRet(rtn)
	})

var _ = builtin4("SendMessageTreeItem(hwnd, msg, wParam, tvitem)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		n := uintptr(getInt32(d, "cchTextMax"))
		var pszText *byte
		var buf unsafe.Pointer
		if n == 0 {
			pszText = getStr(d, "pszText")
		} else {
			buf = heap.Alloc(n)
			pszText = (*byte)(buf)
		}
		p := heap.Alloc(nTVITEMEX)
		tvi := (*TVITEMEX)(p)
		*tvi = TVITEMEX{
			TVITEM: TVITEM{
				mask:           getUint32(d, "mask"),
				hItem:          getHandle(d, "hItem"),
				state:          getUint32(d, "state"),
				stateMask:      getUint32(d, "stateMask"),
				pszText:        pszText,
				cchTextMax:     int32(n),
				iImage:         getInt32(d, "iImage"),
				iSelectedImage: getInt32(d, "iSelectedImage"),
				cChildren:      getInt32(d, "cChildren"),
				lParam:         getHandle(d, "lParam"),
			}}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		d.Put(nil, SuStr("mask"), IntVal(int(tvi.mask)))
		d.Put(nil, SuStr("hItem"), IntVal(int(tvi.hItem)))
		d.Put(nil, SuStr("state"), IntVal(int(tvi.state)))
		d.Put(nil, SuStr("stateMask"), IntVal(int(tvi.stateMask)))
		if n != 0 {
			d.Put(nil, SuStr("pszText"), bufToStr(buf, n))
		}
		d.Put(nil, SuStr("cchTextMax"), IntVal(int(tvi.cchTextMax)))
		d.Put(nil, SuStr("iImage"), IntVal(int(tvi.iImage)))
		d.Put(nil, SuStr("iSelectedImage"), IntVal(int(tvi.iSelectedImage)))
		d.Put(nil, SuStr("cChildren"), IntVal(int(tvi.cChildren)))
		d.Put(nil, SuStr("lParam"), IntVal(int(tvi.lParam)))
		return intRet(rtn)
	})

var _ = builtin4("SendMessageTreeInsert(hwnd, msg, wParam, tvins)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nTVINSERTSTRUCT)
		tvins := (*TVINSERTSTRUCT)(p)
		*tvins = TVINSERTSTRUCT{
			hParent:      getHandle(d, "hParent"),
			hInsertAfter: getHandle(d, "hInsertAfter"),
		}
		item := d.Get(nil, SuStr("item"))
		tvi := &tvins.item
		*tvi = TVITEMEX{
			TVITEM: TVITEM{
				mask:           getUint32(item, "mask"),
				hItem:          getHandle(item, "hItem"),
				state:          getUint32(item, "state"),
				stateMask:      getUint32(item, "stateMask"),
				pszText:        getStr(item, "pszText"),
				cchTextMax:     getInt32(item, "cchTextMax"),
				iImage:         getInt32(item, "iImage"),
				iSelectedImage: getInt32(item, "iSelectedImage"),
				cChildren:      getInt32(item, "cChildren"),
				lParam:         getHandle(item, "lParam"),
			}}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		return intRet(rtn)
	})

var _ = builtin4("SendMessageSBPART(hwnd, msg, wParam, sbpart)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		np := ToInt(c)
		n := uintptr(np) * int32Size
		p := heap.Alloc(n)
		ob := ToContainer(d.Get(nil, SuStr("parts")))
		for i := 0; i < np && i < ob.ListSize(); i++ {
			*(*int32)(unsafe.Pointer(uintptr(p) + int32Size*uintptr(i))) =
				int32(ToInt(ob.ListGet(i)))
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		for i := 0; i < np; i++ {
			x := *(*int32)(unsafe.Pointer(uintptr(p) + int32Size*uintptr(i)))
			ob.Put(nil, SuInt(i), IntVal(int(x)))
		}
		return intRet(rtn)
	})

var _ = builtin4("SendMessageMSG(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nMSG)
		*(*MSG)(p) = MSG{
			hwnd:    getHandle(d, "hwnd"),
			message: getUint32(d, "message"),
			wParam:  getHandle(d, "message"),
			lParam:  getHandle(d, "message"),
			time:    getUint32(d, "message"),
			pt:      getPoint(d, "pt"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		return intRet(rtn)
	})

var _ = builtin4("SendMessageHditem(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		n := getInt32(d, "cchTextMax")
		var pszText *byte
		var buf unsafe.Pointer
		pszText = getStr(d, "pszText")
		if pszText != nil {
			s := ToStr(d.Get(nil, SuStr("pszText")))
			n = int32(len(s))
		} else if n > 0 {
			buf = heap.Alloc(uintptr(n))
			pszText = (*byte)(buf)
		}
		p := heap.Alloc(nHDITEM)
		hdi := (*HDITEM)(p)
		obToHDITEM(d, hdi)
		hdi.pszText = pszText
		hdi.cchTextMax = n
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		hditemToOb(hdi, d)
		if buf != nil {
			d.Put(nil, SuStr("pszText"), bufToStr(buf, uintptr(n)))
		}
		return intRet(rtn)
	})

func obToHDITEM(ob Value, hdi *HDITEM) {
	*hdi = HDITEM{
		mask:       getInt32(ob, "mask"),
		cxy:        getInt32(ob, "cxy"),
		pszText:    hdi.pszText, // not handled
		hbm:        getHandle(ob, "hbm"),
		cchTextMax: getInt32(ob, "cchTextMax"),
		fmt:        getInt32(ob, "fmt"),
		lParam:     getHandle(ob, "lParam"),
		iImage:     getInt32(ob, "iImage"),
		iOrder:     getInt32(ob, "iOrder"),
		typ:        getUint32(ob, "type"),
		pvFilter:   getHandle(ob, "pvFilter"),
		state:      getUint32(ob, "state"),
	}
}

func hditemToOb(hdi *HDITEM, ob Value) Value {
	ob.Put(nil, SuStr("mask"), IntVal(int(hdi.mask)))
	ob.Put(nil, SuStr("cxy"), IntVal(int(hdi.cxy)))
	// pszText must be handled by caller
	ob.Put(nil, SuStr("hbm"), IntVal(int(hdi.hbm)))
	ob.Put(nil, SuStr("cchTextMax"), IntVal(int(hdi.cchTextMax)))
	ob.Put(nil, SuStr("fmt"), IntVal(int(hdi.fmt)))
	ob.Put(nil, SuStr("lParam"), IntVal(int(hdi.lParam)))
	ob.Put(nil, SuStr("iImage"), IntVal(int(hdi.iImage)))
	ob.Put(nil, SuStr("iOrder"), IntVal(int(hdi.iOrder)))
	ob.Put(nil, SuStr("type"), IntVal(int(hdi.typ)))
	ob.Put(nil, SuStr("pvFilter"), IntVal(int(hdi.pvFilter)))
	ob.Put(nil, SuStr("state"), IntVal(int(hdi.state)))
	return ob
}

type HDITEM struct {
	mask       int32
	cxy        int32
	pszText    *byte
	hbm        HANDLE
	cchTextMax int32
	fmt        int32
	lParam     uintptr
	iImage     int32
	iOrder     int32
	typ        uint32
	pvFilter   uintptr
	state      uint32
	_          [4]byte // padding
}

const nHDITEM = unsafe.Sizeof(HDITEM{})

var _ = builtin4("SendMessageHDHITTESTINFO(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nHDHITTESTINFO)
		ht := (*HDHITTESTINFO)(p)
		*ht = HDHITTESTINFO{
			pt:    getPoint(d, "pt"),
			flags: getInt32(d, "flags"),
			iItem: getInt32(d, "iItem"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		d.Put(nil, SuStr("pt"), pointToOb(&ht.pt, nil))
		d.Put(nil, SuStr("flags"), IntVal(int(ht.flags)))
		d.Put(nil, SuStr("iItem"), IntVal(int(ht.iItem)))
		return intRet(rtn)
	})

type HDHITTESTINFO struct {
	pt    POINT
	flags int32
	iItem int32
}

const nHDHITTESTINFO = unsafe.Sizeof(HDHITTESTINFO{})

var _ = builtin4("SendMessageTreeHitTest(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nTVHITTESTINFO)
		ht := (*TVHITTESTINFO)(p)
		*ht = TVHITTESTINFO{
			pt:    getPoint(d, "pt"),
			flags: getInt32(d, "flags"),
			iItem: getHandle(d, "iItem"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		d.Put(nil, SuStr("pt"), pointToOb(&ht.pt, nil))
		d.Put(nil, SuStr("flags"), IntVal(int(ht.flags)))
		d.Put(nil, SuStr("iItem"), IntVal(int(ht.iItem)))
		return intRet(rtn)
	})

type TVHITTESTINFO struct {
	pt    POINT
	flags int32
	iItem HANDLE
}

const nTVHITTESTINFO = unsafe.Sizeof(TVHITTESTINFO{})

var _ = builtin4("SendMessageTabHitTest(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nTCHITTESTINFO)
		ht := (*TCHITTESTINFO)(p)
		*ht = TCHITTESTINFO{
			pt:    getPoint(d, "pt"),
			flags: getInt32(d, "flags"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		d.Put(nil, SuStr("pt"), pointToOb(&ht.pt, nil))
		d.Put(nil, SuStr("flags"), IntVal(int(ht.flags)))
		return intRet(rtn)
	})

type TCHITTESTINFO struct {
	pt    POINT
	flags int32
}

const nTCHITTESTINFO = unsafe.Sizeof(TCHITTESTINFO{})

var _ = builtin4("SendMessageListColumn(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		var p unsafe.Pointer
		if !d.Equal(Zero) {
			p = heap.Alloc(nLVCOLUMN)
			*(*LVCOLUMN)(p) = LVCOLUMN{
				mask:       getInt32(d, "mask"),
				fmt:        getInt32(d, "fmt"),
				cx:         getInt32(d, "cx"),
				pszText:    getStr(d, "pszText"),
				cchTextMax: getInt32(d, "cchTextMax"),
				iSubItem:   getInt32(d, "iSubItem"),
				iImage:     getInt32(d, "iImage"),
				iOrder:     getInt32(d, "iOrder"),
			}
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		return intRet(rtn)
	})

type LVCOLUMN struct {
	mask       int32
	fmt        int32
	cx         int32
	pszText    *byte
	cchTextMax int32
	iSubItem   int32
	iImage     int32
	iOrder     int32
	cxMin      int32
	cxDefault  int32
	cxIdeal    int32
	_          [4]byte // padding
}

const nLVCOLUMN = unsafe.Sizeof(LVCOLUMN{})

var _ = builtin4("SendMessageListItem(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p, li := obToLVITEM(d)
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		d.Put(nil, SuStr("lParam"), IntVal(int(li.lParam)))
		return intRet(rtn)
	})

	func obToLVITEM(ob Value) (unsafe.Pointer, *LVITEM) {
		var p unsafe.Pointer
		var li *LVITEM
		if !ob.Equal(Zero) {
			p = heap.Alloc(nLVITEM)
			li = (*LVITEM)(p)
			*li = LVITEM{
				mask:       getInt32(ob, "mask"),
				iItem:      getInt32(ob, "iItem"),
				iSubItem:   getInt32(ob, "iSubItem"),
				state:      getInt32(ob, "state"),
				stateMask:  getInt32(ob, "stateMask"),
				pszText:    getStr(ob, "pszText"),
				cchTextMax: getInt32(ob, "cchTextMax"),
				iImage:     getInt32(ob, "iImage"),
				lParam:     getHandle(ob, "lParam"),
				iIndent:    getInt32(ob, "iIndent"),
			}
		}
		return p, li
	}
	
const LVIF_TEXT = 1
const LVM_ITEM = 4101

var _ = builtin3("SendMessageListItemOut(hwnd, iItem, iSubItem)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nLVITEM)
		const bufsize = 256
		buf := heap.Alloc(bufsize)
		li := (*LVITEM)(p)
		*li = LVITEM{
			mask:       LVIF_TEXT,
			iItem:      int32(ToInt(b)),
			iSubItem:   int32(ToInt(c)),
			pszText:    (*byte)(buf),
			cchTextMax: bufsize,
		}
		goc.Syscall4(sendMessage,
			intArg(a),
			uintptr(LVM_ITEM),
			0,
			uintptr(p))
		return bufRet(buf, bufsize)
	})

type LVITEM struct {
	mask       int32
	iItem      int32
	iSubItem   int32
	state      int32
	stateMask  int32
	pszText    *byte
	cchTextMax int32
	iImage     int32
	lParam     uintptr
	iIndent    int32
	iGroupId   int32
	cColumns   uint32
	puColumns  uintptr
	piColFmt   uintptr
	iGroup     int32
	_          [4]byte // padding
}

const nLVITEM = unsafe.Sizeof(LVITEM{})

var _ = builtin4("SendMessageListColumnOrder(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		msg := ToInt(b)
		n := ToInt(c)
		p := heap.Alloc(uintptr(n) * int32Size)
		colsob := d.Get(nil, SuStr("order"))
		const LVM_SETCOLUMNORDERARRAY = 4154
		if msg == LVM_SETCOLUMNORDERARRAY {
			for i := 0; i < n; i++ {
				*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(i)*int32Size)) =
					int32(ToInt(colsob.Get(nil, IntVal(i))))
			}
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		if msg != LVM_SETCOLUMNORDERARRAY {
			for i := 0; i < n; i++ {
				colsob.Put(nil, IntVal(i), IntVal(int(*(*int32)(
					unsafe.Pointer(uintptr(p) + uintptr(i)*int32Size)))))
			}
		}
		return intRet(rtn)
	})

var _ = builtin4("SendMessageLVHITTESTINFO(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nLVHITTESTINFO)
		ht := (*LVHITTESTINFO)(p)
		*ht = LVHITTESTINFO{
			pt:       getPoint(d, "pt"),
			flags:    getInt32(d, "flags"),
			iItem:    getInt32(d, "iItem"),
			iSubItem: getInt32(d, "iSubItem"),
			iGroup:   getInt32(d, "iGroup"),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		d.Put(nil, SuStr("iItem"), IntVal(int(ht.iItem)))
		d.Put(nil, SuStr("iSubItem"), IntVal(int(ht.iSubItem)))
		return intRet(rtn)
	})

type LVHITTESTINFO struct {
	pt       POINT
	flags    int32
	iItem    int32
	iSubItem int32
	iGroup   int32
}

const nLVHITTESTINFO = unsafe.Sizeof(LVHITTESTINFO{})

var _ = builtin4("SendMessageSystemTime(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nSYSTEMTIME)
		st := (*SYSTEMTIME)(p)
		*st = obToSYSTEMTIME(d)
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		SYSTEMTIMEtoOb(st, d)
		return intRet(rtn)
	})

func obToSYSTEMTIME(ob Value) SYSTEMTIME {
	return SYSTEMTIME{
		wYear:         getInt16(ob, "wYear"),
		wMonth:        getInt16(ob, "wMonth"),
		wDayOfWeek:    getInt16(ob, "wDayOfWeek"),
		wDay:          getInt16(ob, "wDay"),
		wHour:         getInt16(ob, "wHour"),
		wMinute:       getInt16(ob, "wMinute"),
		wSecond:       getInt16(ob, "wSecond"),
		wMilliseconds: getInt16(ob, "wMilliseconds"),
	}
}

func SYSTEMTIMEtoOb(st *SYSTEMTIME, ob Value) Value {
	ob.Put(nil, SuStr("wYear"), IntVal(int(st.wYear)))
	ob.Put(nil, SuStr("wMonth"), IntVal(int(st.wMonth)))
	ob.Put(nil, SuStr("wDayOfWeek"), IntVal(int(st.wDayOfWeek)))
	ob.Put(nil, SuStr("wDay"), IntVal(int(st.wDay)))
	ob.Put(nil, SuStr("wHour"), IntVal(int(st.wHour)))
	ob.Put(nil, SuStr("wMinute"), IntVal(int(st.wMinute)))
	ob.Put(nil, SuStr("wSecond"), IntVal(int(st.wSecond)))
	ob.Put(nil, SuStr("wMilliseconds"), IntVal(int(st.wMilliseconds)))
	return ob
}

type SYSTEMTIME struct {
	wYear         int16
	wMonth        int16
	wDayOfWeek    int16
	wDay          int16
	wHour         int16
	wMinute       int16
	wSecond       int16
	wMilliseconds int16
}

const nSYSTEMTIME = unsafe.Sizeof(SYSTEMTIME{})

var _ = builtin4("SendMessageSTRange(hwnd, msg, wParam, lParam)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nSystemTimeRange)
		*(*SystemTimeRange)(p) = SystemTimeRange{
			min: obToSYSTEMTIME(d.Get(nil, SuStr("min"))),
			max: obToSYSTEMTIME(d.Get(nil, SuStr("max"))),
		}
		rtn := goc.Syscall4(sendMessage,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p))
		return intRet(rtn)
	})

type SystemTimeRange struct {
	min SYSTEMTIME
	max SYSTEMTIME
}

const nSystemTimeRange = unsafe.Sizeof(SystemTimeRange{})
