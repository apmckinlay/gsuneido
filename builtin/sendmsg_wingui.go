// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"runtime"
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

// dll User32:SendMessage(pointer hwnd, long msg, pointer wParam,
// pointer lParam) pointer
var sendMessage = user32.MustFindProc("SendMessageA").Addr()
var _ = builtin(SendMessage, "(hwnd, msg, wParam, lParam)")

func SendMessage(a, b, c, d Value) Value {
	rtn, _, _ := syscall.SyscallN(sendMessage,
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
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(zstrArg(d))))
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
	n := uintptr(ToInt(d) + 1)
	buf := make([]byte, n)
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&buf[0])))
	ob := &SuObject{}
	ob.Put(nil, SuStr("text"), bufZstr(buf))
	ob.Put(nil, SuStr("result"), intRet(rtn))
	return ob
}

var _ = builtin(SendMessagePoint, "(hwnd, msg, wParam, point)")

func SendMessagePoint(a, b, c, d Value) Value {
	pt := toPoint(d)
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(pt)))
	fromPoint(pt, d)
	return intRet(rtn)
}

var _ = builtin(SendMessageRect, "(hwnd, msg, wParam, rect)")

func SendMessageRect(a, b, c, d Value) Value {
	r := toRect(d)
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(r)))
	fromRect(r, d)
	return intRet(rtn)
}

var _ = builtin(SendMessageTcitem, "(hwnd, msg, wParam, tcitem)")

func SendMessageTcitem(a, b, c, d Value) Value {
	t := getZstr(d, "pszText")
	n := getInt32(d, "cchTextMax")
	var buf []byte
	if n > 0 {
		buf = make([]byte, n)
		t = &buf[0]
	}
	tci := stTCItem{
		mask:        getUint32(d, "mask"),
		dwState:     getUint32(d, "dwState"),
		dwStateMask: getUint32(d, "dwStateMask"),
		pszText:     t,
		cchTextMax:  n,
		iImage:      getInt32(d, "iImage"),
		lParam:      getInt32(d, "lParam"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&tci)))
	if n > 0 {
		d.Put(nil, SuStr("pszText"), bufZstr(buf))
	}
	d.Put(nil, SuStr("iImage"), IntVal(int(tci.iImage)))
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

var _ = builtin(SendMessageTextRange, "(hwnd, msg, cpMin, cpMax, each = 1)")

func SendMessageTextRange(a, b, c, d, e Value) Value {
	cpMin := ToInt(c)
	cpMax := ToInt(d)
	if cpMax <= cpMin {
		return EmptyStr
	}
	each := ToInt(e)
	n := (cpMax - cpMin) * each
	buf := make([]byte, n+each)
	tr := stTextRange{
		chrg:      stCharRange{cpMin: int32(cpMin), cpMax: int32(cpMax)},
		lpstrText: &buf[0],
	}
	syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		0,
		uintptr(unsafe.Pointer(&tr)))
	return SuStr(hacks.BStoS(buf[:n]))
}

var _ = builtin(SendMessageTOOLINFO, "(hwnd, msg, wParam, lParam)")

func SendMessageTOOLINFO(a, b, c, d Value) Value {
	ti := stToolInfo{
		cbSize:   uint32(nToolInfo),
		uFlags:   getUint32(d, "uFlags"),
		hwnd:     getUintptr(d, "hwnd"),
		uId:      getUintptr(d, "uId"),
		hinst:    getUintptr(d, "hinst"),
		lpszText: getZstr(d, "lpszText"),
		lParam:   getUintptr(d, "lParam"),
		rect:     getRect(d, "rect"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&ti)))
	return intRet(rtn)
}

var _ = builtin(SendMessageTOOLINFO2, "(hwnd, msg, wParam, lParam)")

func SendMessageTOOLINFO2(a, b, c, d Value) Value {
	ti2 := stToolInfo2{
		cbSize:   uint32(nToolInfo),
		uFlags:   getUint32(d, "uFlags"),
		hwnd:     getUintptr(d, "hwnd"),
		uId:      getUintptr(d, "uId"),
		hinst:    getUintptr(d, "hinst"),
		lpszText: getUintptr(d, "lpszText"), // for LPSTR_TEXTCALLBACK
		lParam:   getUintptr(d, "lParam"),
		rect:     getRect(d, "rect"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&ti2)))
	return intRet(rtn)
}

var _ = builtin(SendMessageTreeItem, "(hwnd, msg, wParam, tvitem)")

func SendMessageTreeItem(a, b, c, d Value) Value {
	n := uintptr(getInt32(d, "cchTextMax"))
	var pszText *byte
	var buf []byte
	if n == 0 {
		pszText = getZstr(d, "pszText")
	} else {
		buf = make([]byte, n)
		pszText = &buf[0]
	}
	tvi := stTVItemEx{
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
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&tvi)))
	d.Put(nil, SuStr("mask"), IntVal(int(tvi.mask)))
	d.Put(nil, SuStr("hItem"), IntVal(int(tvi.hItem)))
	d.Put(nil, SuStr("state"), IntVal(int(tvi.state)))
	d.Put(nil, SuStr("stateMask"), IntVal(int(tvi.stateMask)))
	if n != 0 {
		d.Put(nil, SuStr("pszText"), bufZstr(buf))
	}
	d.Put(nil, SuStr("cchTextMax"), IntVal(int(tvi.cchTextMax)))
	d.Put(nil, SuStr("iImage"), IntVal(int(tvi.iImage)))
	d.Put(nil, SuStr("iSelectedImage"), IntVal(int(tvi.iSelectedImage)))
	d.Put(nil, SuStr("cChildren"), IntVal(int(tvi.cChildren)))
	d.Put(nil, SuStr("lParam"), IntVal(int(tvi.lParam)))
	return intRet(rtn)
}

var _ = builtin(SendMessageTreeSort, "(hwnd, msg, wParam, tvitem)")

func SendMessageTreeSort(th *Thread, args []Value) Value {
	tvcb := stTVSortCB{
		hParent:     getUintptr(args[3], "hParent"),
		lpfnCompare: getCallback(th, args[3], "lpfnCompare", 3),
		lParam:      getUintptr(args[3], "lParam"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(args[0]),
		intArg(args[1]),
		intArg(args[2]),
		uintptr(unsafe.Pointer(&tvcb)))
	return boolRet(rtn)
}

var _ = builtin(SendMessageTreeInsert, "(hwnd, msg, wParam, tvins)")

func SendMessageTreeInsert(a, b, c, d Value) Value {
	tvins := stTVInsertStruct{
		hParent:      getUintptr(d, "hParent"),
		hInsertAfter: getUintptr(d, "hInsertAfter"),
	}
	item := d.Get(nil, SuStr("item"))
	tvins.item.stTVItem = stTVItem{
		mask:           getUint32(item, "mask"),
		hItem:          getUintptr(item, "hItem"),
		state:          getUint32(item, "state"),
		stateMask:      getUint32(item, "stateMask"),
		pszText:        getZstr(item, "pszText"),
		cchTextMax:     getInt32(item, "cchTextMax"),
		iImage:         getInt32(item, "iImage"),
		iSelectedImage: getInt32(item, "iSelectedImage"),
		cChildren:      getInt32(item, "cChildren"),
		lParam:         getUintptr(item, "lParam"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&tvins)))
	return intRet(rtn)
}

var _ = builtin(SendMessageSBPART, "(hwnd, msg, wParam, sbpart)")

func SendMessageSBPART(a, b, c, d Value) Value {
	np := ToInt(c)
	p := make([]int32, np)
	var ob Container
	if parts := d.Get(nil, SuStr("parts")); parts != nil {
		ob = ToContainer(parts)
		for i := range min(np, ob.ListSize()) {
			p[i] = int32(ToInt(ob.ListGet(i)))
		}
	} else {
		ob = &SuObject{}
		d.Put(nil, SuStr("parts"), ob)
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(unsafe.SliceData(p))))
	for i := range np {
		ob.Put(nil, SuInt(i), IntVal(int(p[i])))
	}
	return intRet(rtn)
}

var _ = builtin(SendMessageMSG, "(hwnd, msg, wParam, lParam)")

func SendMessageMSG(a, b, c, d Value) Value {
	msg := stMsg{
		hwnd:    getUintptr(d, "hwnd"),
		message: getUint32(d, "message"),
		wParam:  getUintptr(d, "wParam"),
		lParam:  getUintptr(d, "lParam"),
		time:    getUint32(d, "time"),
		pt:      getPoint(d, "pt"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&msg)))
	return intRet(rtn)
}

var _ = builtin(SendMessageHditem, "(hwnd, msg, wParam, lParam)")

func SendMessageHditem(a, b, c, d Value) Value {
	pszTextGiven := d.Get(nil, SuStr("pszText")) != nil
	cchTextMaxGiven := d.Get(nil, SuStr("cchTextMax")) != nil
	if pszTextGiven && cchTextMaxGiven {
		panic("SendMessageHditem: pszText and cchTextMax cannot be used together")
	}
	hdi := toHdItem(d)
	cchTextMax := getInt32(d, "cchTextMax")
	if cchTextMaxGiven {
		buf := make([]byte, cchTextMax)
		hdi.pszText = &buf[0]
	}
	if pszTextGiven {
		hdi.pszText = getZstr(d, "pszText")
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&hdi)))
	fromHdItem(&hdi, d)
	if cchTextMaxGiven {
		d.Put(nil, SuStr("pszText"),
			ptrZstr(unsafe.Pointer(hdi.pszText), int(cchTextMax)))
		ToContainer(d).Delete(nil, SuStr("cchTextMax"))
		runtime.KeepAlive(hdi.pszText)
	}
	return intRet(rtn)
}

func toHdItem(ob Value) stHdItem {
	return stHdItem{
		mask: getInt32(ob, "mask"),
		cxy:  getInt32(ob, "cxy"),
		// pszText not handled
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

func fromHdItem(hdi *stHdItem, ob Value) Value {
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

var _ = builtin(SendMessageHDHITTESTINFO, "(hwnd, msg, wParam, lParam)")

func SendMessageHDHITTESTINFO(a, b, c, d Value) Value {
	ht := stHdHitTestInfo{
		pt:    getPoint(d, "pt"),
		flags: getInt32(d, "flags"),
		iItem: getInt32(d, "iItem"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&ht)))
	d.Put(nil, SuStr("pt"), fromPoint(&ht.pt, nil))
	d.Put(nil, SuStr("flags"), IntVal(int(ht.flags)))
	d.Put(nil, SuStr("iItem"), IntVal(int(ht.iItem)))
	return intRet(rtn)
}

type stHdHitTestInfo struct {
	pt    stPoint
	flags int32
	iItem int32
}

var _ = builtin(SendMessageTreeHitTest, "(hwnd, msg, wParam, lParam)")

func SendMessageTreeHitTest(a, b, c, d Value) Value {
	ht := stTVHitTestInfo{
		pt:    getPoint(d, "pt"),
		flags: getInt32(d, "flags"),
		iItem: getUintptr(d, "iItem"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&ht)))
	d.Put(nil, SuStr("pt"), fromPoint(&ht.pt, nil))
	d.Put(nil, SuStr("flags"), IntVal(int(ht.flags)))
	d.Put(nil, SuStr("iItem"), IntVal(int(ht.iItem)))
	return intRet(rtn)
}

type stTVHitTestInfo struct {
	pt    stPoint
	flags int32
	iItem HANDLE
}

var _ = builtin(SendMessageTabHitTest, "(hwnd, msg, wParam, lParam)")

func SendMessageTabHitTest(a, b, c, d Value) Value {
	ht := stTCHitTestInfo{
		pt:    getPoint(d, "pt"),
		flags: getInt32(d, "flags"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&ht)))
	d.Put(nil, SuStr("pt"), fromPoint(&ht.pt, nil))
	d.Put(nil, SuStr("flags"), IntVal(int(ht.flags)))
	return intRet(rtn)
}

type stTCHitTestInfo struct {
	pt    stPoint
	flags int32
}

var _ = builtin(SendMessageListColumn, "(hwnd, msg, wParam, lParam)")

func SendMessageListColumn(a, b, c, d Value) Value {
	var lvc *stLVColumn
	if !d.Equal(Zero) {
		lvc = &stLVColumn{
			mask:       getInt32(d, "mask"),
			fmt:        getInt32(d, "fmt"),
			cx:         getInt32(d, "cx"),
			pszText:    getZstr(d, "pszText"),
			cchTextMax: getInt32(d, "cchTextMax"),
			iSubItem:   getInt32(d, "iSubItem"),
			iImage:     getInt32(d, "iImage"),
			iOrder:     getInt32(d, "iOrder"),
		}
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(lvc)))
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

var _ = builtin(SendMessageListItem, "(hwnd, msg, wParam, lParam)")

func SendMessageListItem(a, b, c, d Value) Value {
	li := toLVItem(d)
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(li)))
	d.Put(nil, SuStr("lParam"), IntVal(int(li.lParam)))
	return intRet(rtn)
}

func toLVItem(ob Value) *stLVItem {
	if ob.Equal(Zero) {
		return nil
	}
	return &stLVItem{
		mask:       getInt32(ob, "mask"),
		iItem:      getInt32(ob, "iItem"),
		iSubItem:   getInt32(ob, "iSubItem"),
		state:      getInt32(ob, "state"),
		stateMask:  getInt32(ob, "stateMask"),
		pszText:    getZstr(ob, "pszText"),
		cchTextMax: getInt32(ob, "cchTextMax"),
		iImage:     getInt32(ob, "iImage"),
		lParam:     getUintptr(ob, "lParam"),
		iIndent:    getInt32(ob, "iIndent"),
	}
}

const LVIF_TEXT = 1
const LVM_ITEM = 4101

var _ = builtin(SendMessageListItemOut, "(hwnd, iItem, iSubItem)")

func SendMessageListItemOut(a, b, c Value) Value {
	const bufsize = 256
	buf := make([]byte, bufsize)
	li := stLVItem{
		mask:       LVIF_TEXT,
		iItem:      int32(ToInt(b)),
		iSubItem:   int32(ToInt(c)),
		pszText:    &buf[0],
		cchTextMax: bufsize,
	}
	syscall.SyscallN(sendMessage,
		intArg(a),
		uintptr(LVM_ITEM),
		0,
		uintptr(unsafe.Pointer(&li)))
	return bufZstr(buf)
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

var _ = builtin(SendMessageListColumnOrder, "(hwnd, msg, wParam, lParam)")

func SendMessageListColumnOrder(a, b, c, d Value) Value {
	n := ToInt(c)
	p := make([]int32, n)
	colsob := d.Get(nil, SuStr("order"))
	if colsob == nil {
		colsob = &SuObject{}
		d.Put(nil, SuStr("order"), colsob)
	}
	for i := range n {
		if x := colsob.Get(nil, IntVal(i)); x != nil {
			p[i] = int32(ToInt(x))
		}
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(unsafe.SliceData(p))))
	for i := range n {
		colsob.Put(nil, IntVal(i), IntVal(int(p[i])))
	}
	return intRet(rtn)
}

var _ = builtin(SendMessageLVHITTESTINFO, "(hwnd, msg, wParam, lParam)")

func SendMessageLVHITTESTINFO(a, b, c, d Value) Value {
	ht := stLVHitTestInfo{
		pt:       getPoint(d, "pt"),
		flags:    getInt32(d, "flags"),
		iItem:    getInt32(d, "iItem"),
		iSubItem: getInt32(d, "iSubItem"),
		iGroup:   getInt32(d, "iGroup"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&ht)))
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

var _ = builtin(SendMessageSystemTime, "(hwnd, msg, wParam, lParam)")

func SendMessageSystemTime(a, b, c, d Value) Value {
	st := toSystemTime(d)
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&st)))
	SYSTEMTIMEtoOb(&st, d)
	return intRet(rtn)
}

func toSystemTime(ob Value) stSystemTime {
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
	str := SystemTimeRange{
		min: toSystemTime(d.Get(nil, SuStr("min"))),
		max: toSystemTime(d.Get(nil, SuStr("max"))),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&str)))
	return intRet(rtn)
}

type SystemTimeRange struct {
	min stSystemTime
	max stSystemTime
}

var _ = builtin(SendMessageEDITBALLOONTIP, "(hwnd, msg, wParam, lParam)")

func SendMessageEDITBALLOONTIP(a, b, c, d Value) Value {
	ebt := stEditBalloonTip{
		cbStruct: int32(nEditBalloonTip),
		pszTitle: getZstr(d, "pszTitle"),
		pszText:  getZstr(d, "pszText"),
		ttiIcon:  getInt32(d, "ttiIcon"),
	}
	rtn, _, _ := syscall.SyscallN(sendMessage,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(unsafe.Pointer(&ebt)))
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
