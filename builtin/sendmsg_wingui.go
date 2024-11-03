// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
)

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
// pointer lParam) pointer
var sendMessage = user32.MustFindProc("SendMessageA").Addr()
var _ = builtin(SendMessage, "(hwnd, msg, wParam, lParam)")

func SendMessage(a, b, c, d Value) Value {
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d))
	return intRet(rtn)
}

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
// string text) pointer
var _ = builtin(SendMessageText, "(hwnd, msg, wParam, text)")

func SendMessageText(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(stringArg(d)))
	return intRet(rtn)
}

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
// [in] string text) pointer
var _ = builtin(SendMessageTextIn, "(hwnd, msg, wParam, text)")

func SendMessageTextIn(a, b, c, d Value) Value {
	// can't pass direct pointer so same as SendMessageText (i.e. still copies)
	return SendMessageText(a, b, c, d)
}

var _ = builtin(SendMessageTextOut, "(hwnd, msg, wParam = 0, bufsize = 1024)")

func SendMessageTextOut(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	n := uintptr(ToInt(d) + 1)
	buf := heap.Alloc(n)
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(buf))
	ob := &SuObject{}
	ob.Put(nil, SuStr("text"), SuStr(heap.GetStrZ(buf, int(n))))
	ob.Put(nil, SuStr("result"), intRet(rtn))
	return ob
}

var _ = builtin(SendMessagePoint, "(hwnd, msg, wParam, point)")

func SendMessagePoint(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	pt := heap.Alloc(nPoint)
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(pointArg(d, pt)))
	upointToOb(pt, d)
	return intRet(rtn)
}

var _ = builtin(SendMessageRect, "(hwnd, msg, wParam, rect)")

func SendMessageRect(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(rectArg(d, r)))
	urectToOb(r, d)
	return intRet(rtn)
}

var _ = builtin(SendMessageTcitem, "(hwnd, msg, wParam, tcitem)")

func SendMessageTcitem(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	t := getStr(d, "pszText")
	n := getInt32(d, "cchTextMax")
	if n > 0 {
		t = (*byte)(heap.Alloc(uintptr(n)))
	}
	p := heap.Alloc(nTCItem)
	*(*stTCItem)(p) = stTCItem{
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
		d.Put(nil, SuStr("pszText"),
			SuStr(heap.GetStrZ(unsafe.Pointer(t), int(n))))
	}
	d.Put(nil, SuStr("iImage"), IntVal(int((*stTCItem)(p).iImage)))
	return intRet(rtn)
}

type stTCItem struct {
	mask        uint32
	dwState     uint32
	dwStateMask uint32
	pszText     *byte
	cchTextMax  int32
	iImage      int32
	lParam      int32
	_           [4]byte // padding
}

const nTCItem = unsafe.Sizeof(stTCItem{})

var _ = builtin(SendMessageTextRange, "(hwnd, msg, cpMin, cpMax, each = 1)")

func SendMessageTextRange(a, b, c, d, e Value) Value {
	defer heap.FreeTo(heap.CurSize())
	cpMin := ToInt(c)
	cpMax := ToInt(d)
	if cpMax <= cpMin {
		return EmptyStr
	}
	each := ToInt(e)
	n := (cpMax - cpMin) * each
	buf := heap.Alloc(uintptr(n + each))
	p := heap.Alloc(nTextRange)
	*(*stTextRange)(p) = stTextRange{
		chrg:      stCharRange{cpMin: int32(cpMin), cpMax: int32(cpMax)},
		lpstrText: (*byte)(buf),
	}
	goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		0,
		uintptr(p))
	return SuStr(heap.GetStrN(buf, n))
}

var _ = builtin(SendMessageTOOLINFO, "(hwnd, msg, wParam, lParam)")

func SendMessageTOOLINFO(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nToolInfo)
	*(*stToolInfo)(p) = stToolInfo{
		cbSize:   uint32(nToolInfo),
		uFlags:   getUint32(d, "uFlags"),
		hwnd:     getUintptr(d, "hwnd"),
		uId:      getUintptr(d, "uId"),
		hinst:    getUintptr(d, "hinst"),
		lpszText: getStr(d, "lpszText"),
		lParam:   getUintptr(d, "lParam"),
		rect:     getRect(d, "rect"),
	}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	return intRet(rtn)
}

var _ = builtin(SendMessageTOOLINFO2, "(hwnd, msg, wParam, lParam)")

func SendMessageTOOLINFO2(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nToolInfo)
	*(*stToolInfo2)(p) = stToolInfo2{
		cbSize:   uint32(nToolInfo),
		uFlags:   getUint32(d, "uFlags"),
		hwnd:     getUintptr(d, "hwnd"),
		uId:      getUintptr(d, "uId"),
		hinst:    getUintptr(d, "hinst"),
		lpszText: getUintptr(d, "lpszText"), // for LPSTR_TEXTCALLBACK
		lParam:   getUintptr(d, "lParam"),
		rect:     getRect(d, "rect"),
	}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	return intRet(rtn)
}

var _ = builtin(SendMessageTreeItem, "(hwnd, msg, wParam, tvitem)")

func SendMessageTreeItem(a, b, c, d Value) Value {
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
	p := heap.Alloc(nTVItemEx)
	tvi := (*stTVItemEx)(p)
	*tvi = stTVItemEx{
		stTVItem: stTVItem{
			mask:           getUint32(d, "mask"),
			hItem:          getUintptr(d, "hItem"),
			state:          getUint32(d, "state"),
			stateMask:      getUint32(d, "stateMask"),
			pszText:        pszText,
			cchTextMax:     int32(n),
			iImage:         getInt32(d, "iImage"),
			iSelectedImage: getInt32(d, "iSelectedImage"),
			cChildren:      getInt32(d, "cChildren"),
			lParam:         getUintptr(d, "lParam"),
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
		d.Put(nil, SuStr("pszText"), SuStr(heap.GetStrZ(buf, int(n))))
	}
	d.Put(nil, SuStr("cchTextMax"), IntVal(int(tvi.cchTextMax)))
	d.Put(nil, SuStr("iImage"), IntVal(int(tvi.iImage)))
	d.Put(nil, SuStr("iSelectedImage"), IntVal(int(tvi.iSelectedImage)))
	d.Put(nil, SuStr("cChildren"), IntVal(int(tvi.cChildren)))
	d.Put(nil, SuStr("lParam"), IntVal(int(tvi.lParam)))
	return intRet(rtn)
}

var _ = builtin(SendMessageTreeSort, "(hwnd, msg, wParam, tvitem)")

func SendMessageTreeSort(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nTVSortCB)
	*(*stTVSortCB)(p) = stTVSortCB{
		hParent:     getUintptr(d, "hParent"),
		lpfnCompare: getCallback(d, "lpfnCompare", 3),
		lParam:      getUintptr(d, "lParam"),
	}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	return boolRet(rtn)
}

var _ = builtin(SendMessageTreeInsert, "(hwnd, msg, wParam, tvins)")

func SendMessageTreeInsert(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nTVInsertStruct)
	tvins := (*stTVInsertStruct)(p)
	*tvins = stTVInsertStruct{
		hParent:      getUintptr(d, "hParent"),
		hInsertAfter: getUintptr(d, "hInsertAfter"),
	}
	item := d.Get(nil, SuStr("item"))
	tvi := &tvins.item
	*tvi = stTVItemEx{
		stTVItem: stTVItem{
			mask:           getUint32(item, "mask"),
			hItem:          getUintptr(item, "hItem"),
			state:          getUint32(item, "state"),
			stateMask:      getUint32(item, "stateMask"),
			pszText:        getStr(item, "pszText"),
			cchTextMax:     getInt32(item, "cchTextMax"),
			iImage:         getInt32(item, "iImage"),
			iSelectedImage: getInt32(item, "iSelectedImage"),
			cChildren:      getInt32(item, "cChildren"),
			lParam:         getUintptr(item, "lParam"),
		}}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	return intRet(rtn)
}

var _ = builtin(SendMessageSBPART, "(hwnd, msg, wParam, sbpart)")

func SendMessageSBPART(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	np := ToInt(c)
	n := uintptr(np) * int32Size
	p := heap.Alloc(n)
	var ob Container
	if parts := d.Get(nil, SuStr("parts")); parts != nil {
		ob = ToContainer(parts)
		for i := range min(np, ob.ListSize()) {
			*(*int32)(unsafe.Pointer(uintptr(p) + int32Size*uintptr(i))) =
				int32(ToInt(ob.ListGet(i)))
		}
	} else {
		ob = &SuObject{}
		d.Put(nil, SuStr("parts"), ob)
	}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	for i := range np {
		x := *(*int32)(unsafe.Pointer(uintptr(p) + int32Size*uintptr(i)))
		ob.Put(nil, SuInt(i), IntVal(int(x)))
	}
	return intRet(rtn)
}

var _ = builtin(SendMessageMSG, "(hwnd, msg, wParam, lParam)")

func SendMessageMSG(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nMsg)
	*(*stMsg)(p) = stMsg{
		hwnd:    getUintptr(d, "hwnd"),
		message: getUint32(d, "message"),
		wParam:  getUintptr(d, "wParam"),
		lParam:  getUintptr(d, "lParam"),
		time:    getUint32(d, "time"),
		pt:      getPoint(d, "pt"),
	}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	return intRet(rtn)
}

var _ = builtin(SendMessageHditem, "(hwnd, msg, wParam, lParam)")

func SendMessageHditem(a, b, c, d Value) Value {
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
	p := heap.Alloc(nHdItem)
	hdi := (*stHdItem)(p)
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
		d.Put(nil, SuStr("pszText"), SuStr(heap.GetStrZ(buf, int(n))))
	}
	return intRet(rtn)
}

func obToHDITEM(ob Value, hdi *stHdItem) {
	*hdi = stHdItem{
		mask:       getInt32(ob, "mask"),
		cxy:        getInt32(ob, "cxy"),
		pszText:    hdi.pszText, // not handled
		hbm:        getUintptr(ob, "hbm"),
		cchTextMax: getInt32(ob, "cchTextMax"),
		fmt:        getInt32(ob, "fmt"),
		lParam:     getUintptr(ob, "lParam"),
		iImage:     getInt32(ob, "iImage"),
		iOrder:     getInt32(ob, "iOrder"),
		typ:        getUint32(ob, "type"),
		pvFilter:   getUintptr(ob, "pvFilter"),
		state:      getUint32(ob, "state"),
	}
}

func hditemToOb(hdi *stHdItem, ob Value) Value {
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

type stHdItem struct {
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

const nHdItem = unsafe.Sizeof(stHdItem{})

var _ = builtin(SendMessageHDHITTESTINFO, "(hwnd, msg, wParam, lParam)")

func SendMessageHDHITTESTINFO(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nHdHitTestInfo)
	ht := (*stHdHitTestInfo)(p)
	*ht = stHdHitTestInfo{
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
}

type stHdHitTestInfo struct {
	pt    stPoint
	flags int32
	iItem int32
}

const nHdHitTestInfo = unsafe.Sizeof(stHdHitTestInfo{})

var _ = builtin(SendMessageTreeHitTest, "(hwnd, msg, wParam, lParam)")

func SendMessageTreeHitTest(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nTVHitTestInfo)
	ht := (*stTVHitTestInfo)(p)
	*ht = stTVHitTestInfo{
		pt:    getPoint(d, "pt"),
		flags: getInt32(d, "flags"),
		iItem: getUintptr(d, "iItem"),
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
}

type stTVHitTestInfo struct {
	pt    stPoint
	flags int32
	iItem HANDLE
}

const nTVHitTestInfo = unsafe.Sizeof(stTVHitTestInfo{})

var _ = builtin(SendMessageTabHitTest, "(hwnd, msg, wParam, lParam)")

func SendMessageTabHitTest(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nTCHitTestInfo)
	ht := (*stTCHitTestInfo)(p)
	*ht = stTCHitTestInfo{
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
}

type stTCHitTestInfo struct {
	pt    stPoint
	flags int32
}

const nTCHitTestInfo = unsafe.Sizeof(stTCHitTestInfo{})

var _ = builtin(SendMessageListColumn, "(hwnd, msg, wParam, lParam)")

func SendMessageListColumn(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	var p unsafe.Pointer
	if !d.Equal(Zero) {
		p = heap.Alloc(nLVColumn)
		*(*stLVColumn)(p) = stLVColumn{
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
}

type stLVColumn struct {
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

const nLVColumn = unsafe.Sizeof(stLVColumn{})

var _ = builtin(SendMessageListItem, "(hwnd, msg, wParam, lParam)")

func SendMessageListItem(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p, li := obToLVITEM(d)
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	d.Put(nil, SuStr("lParam"), IntVal(int(li.lParam)))
	return intRet(rtn)
}

func obToLVITEM(ob Value) (unsafe.Pointer, *stLVItem) {
	var p unsafe.Pointer
	var li *stLVItem
	if !ob.Equal(Zero) {
		p = heap.Alloc(nLVItem)
		li = (*stLVItem)(p)
		*li = stLVItem{
			mask:       getInt32(ob, "mask"),
			iItem:      getInt32(ob, "iItem"),
			iSubItem:   getInt32(ob, "iSubItem"),
			state:      getInt32(ob, "state"),
			stateMask:  getInt32(ob, "stateMask"),
			pszText:    getStr(ob, "pszText"),
			cchTextMax: getInt32(ob, "cchTextMax"),
			iImage:     getInt32(ob, "iImage"),
			lParam:     getUintptr(ob, "lParam"),
			iIndent:    getInt32(ob, "iIndent"),
		}
	}
	return p, li
}

const LVIF_TEXT = 1
const LVM_ITEM = 4101

var _ = builtin(SendMessageListItemOut, "(hwnd, iItem, iSubItem)")

func SendMessageListItemOut(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nLVItem)
	const bufsize = 256
	buf := heap.Alloc(bufsize)
	li := (*stLVItem)(p)
	*li = stLVItem{
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
	return SuStr(heap.GetStrZ(buf, bufsize))
}

type stLVItem struct {
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

const nLVItem = unsafe.Sizeof(stLVItem{})

var _ = builtin(SendMessageListColumnOrder, "(hwnd, msg, wParam, lParam)")

func SendMessageListColumnOrder(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	n := ToInt(c)
	p := heap.Alloc(uintptr(n) * int32Size)
	colsob := d.Get(nil, SuStr("order"))
	if colsob == nil {
		colsob = &SuObject{}
		d.Put(nil, SuStr("order"), colsob)
	}
	for i := range n {
		if x := colsob.Get(nil, IntVal(i)); x != nil {
			*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(i)*int32Size)) =
				int32(ToInt(x))
		}
	}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	for i := range n {
		colsob.Put(nil, IntVal(i), IntVal(int(*(*int32)(
			unsafe.Pointer(uintptr(p) + uintptr(i)*int32Size)))))
	}
	return intRet(rtn)
}

var _ = builtin(SendMessageLVHITTESTINFO, "(hwnd, msg, wParam, lParam)")

func SendMessageLVHITTESTINFO(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nLVHitTestInfo)
	ht := (*stLVHitTestInfo)(p)
	*ht = stLVHitTestInfo{
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
}

type stLVHitTestInfo struct {
	pt       stPoint
	flags    int32
	iItem    int32
	iSubItem int32
	iGroup   int32
}

const nLVHitTestInfo = unsafe.Sizeof(stLVHitTestInfo{})

var _ = builtin(SendMessageSystemTime, "(hwnd, msg, wParam, lParam)")

func SendMessageSystemTime(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nSystemTime)
	st := (*stSystemTime)(p)
	*st = obToSYSTEMTIME(d)
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	SYSTEMTIMEtoOb(st, d)
	return intRet(rtn)
}

func obToSYSTEMTIME(ob Value) stSystemTime {
	return stSystemTime{
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

func SYSTEMTIMEtoOb(st *stSystemTime, ob Value) Value {
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

type stSystemTime struct {
	wYear         int16
	wMonth        int16
	wDayOfWeek    int16
	wDay          int16
	wHour         int16
	wMinute       int16
	wSecond       int16
	wMilliseconds int16
}

const nSystemTime = unsafe.Sizeof(stSystemTime{})

var _ = builtin(SendMessageSTRange, "(hwnd, msg, wParam, lParam)")

func SendMessageSTRange(a, b, c, d Value) Value {
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
}

type SystemTimeRange struct {
	min stSystemTime
	max stSystemTime
}

const nSystemTimeRange = unsafe.Sizeof(SystemTimeRange{})

var _ = builtin(SendMessageEDITBALLOONTIP, "(hwnd, msg, wParam, lParam)")

func SendMessageEDITBALLOONTIP(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nEditBalloonTip)
	*(*stEditBalloonTip)(p) = stEditBalloonTip{
		cbStruct: int32(nEditBalloonTip),
		pszTitle: getStr(d, "pszTitle"),
		pszText:  getStr(d, "pszText"),
		ttiIcon:  getInt32(d, "ttiIcon"),
	}
	rtn := goc.Syscall4(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(p))
	return intRet(rtn)
}

type stEditBalloonTip struct {
	cbStruct int32
	pszTitle *byte
	pszText  *byte
	ttiIcon  int32
	_        [4]byte // padding
}

const nEditBalloonTip = unsafe.Sizeof(stEditBalloonTip{})
