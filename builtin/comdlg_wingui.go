// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"bytes"
	"syscall"
	"unsafe"

	. "github.com/apmckinlay/gsuneido/core"
)

var comdlg32 = MustLoadDLL("comdlg32.dll")

// dll long ComDlg32:CommDlgExtendedError()
var commDlgExtendedError = comdlg32.MustFindProc("CommDlgExtendedError").Addr()
var _ = builtin(CommDlgExtendedError, "()")

func CommDlgExtendedError() Value {
	rtn, _, _ := syscall.SyscallN(commDlgExtendedError)
	return intRet(rtn)
}

// dll bool ComDlg32:PrintDlg(PRINTDLG* printdlg)
var printDlg = comdlg32.MustFindProc("PrintDlgA").Addr()
var _ = builtin(PrintDlg, "(printdlg)")

func PrintDlg(a Value) Value {
	pd := stPrintDlg{
		lStructSize:         uint32(nPrintDlg),
		hwndOwner:           getUintptr(a, "hwndOwner"),
		hDevMode:            getUintptr(a, "hDevMode"),
		hDevNames:           getUintptr(a, "hDevNames"),
		hDC:                 getUintptr(a, "hDC"),
		Flags:               getInt32(a, "Flags"),
		nFromPage:           getInt16(a, "nFromPage"),
		nToPage:             getInt16(a, "nToPage"),
		nMinPage:            getInt16(a, "nMinPage"),
		nMaxPage:            getInt16(a, "nMaxPage"),
		nCopies:             getInt16(a, "nCopies"),
		hInstance:           getUintptr(a, "hInstance"),
		lCustData:           getUintptr(a, "lCustData"),
		lpPrintTemplateName: getZstr(a, "lpPrintTemplateName"),
		lpSetupTemplateName: getZstr(a, "lpSetupTemplateName"),
		hPrintTemplate:      getUintptr(a, "hPrintTemplate"),
		hSetupTemplate:      getUintptr(a, "hSetupTemplate"),
	}
	rtn, _, _ := syscall.SyscallN(printDlg,
		uintptr(unsafe.Pointer(&pd)))
	a.Put(nil, SuStr("hwndOwner"), IntVal(int(pd.hwndOwner)))
	a.Put(nil, SuStr("hDevMode"), IntVal(int(pd.hDevMode)))
	a.Put(nil, SuStr("hDevNames"), IntVal(int(pd.hDevNames)))
	a.Put(nil, SuStr("hDC"), IntVal(int(pd.hDC)))
	a.Put(nil, SuStr("Flags"), IntVal(int(pd.Flags)))
	a.Put(nil, SuStr("nFromPage"), IntVal(int(pd.nFromPage)))
	a.Put(nil, SuStr("nToPage"), IntVal(int(pd.nToPage)))
	a.Put(nil, SuStr("nMinPage"), IntVal(int(pd.nMinPage)))
	a.Put(nil, SuStr("nMaxPage"), IntVal(int(pd.nMaxPage)))
	a.Put(nil, SuStr("nCopies"), IntVal(int(pd.nCopies)))
	a.Put(nil, SuStr("hInstance"), IntVal(int(pd.hInstance)))
	a.Put(nil, SuStr("lCustData"), IntVal(int(pd.lCustData)))
	a.Put(nil, SuStr("hPrintTemplate"), IntVal(int(pd.hPrintTemplate)))
	a.Put(nil, SuStr("hSetupTemplate"), IntVal(int(pd.hSetupTemplate)))
	return boolRet(rtn)
}

type stPrintDlg struct {
	lStructSize         uint32
	hwndOwner           HANDLE
	hDevMode            HANDLE
	hDevNames           HANDLE
	hDC                 HANDLE
	Flags               int32
	nFromPage           int16
	nToPage             int16
	nMinPage            int16
	nMaxPage            int16
	nCopies             int16
	hInstance           HANDLE
	lCustData           HANDLE
	lpfnPrintHook       HANDLE
	lpfnSetupHook       HANDLE
	lpPrintTemplateName *byte
	lpSetupTemplateName *byte
	hPrintTemplate      HANDLE
	hSetupTemplate      HANDLE
}

const nPrintDlg = unsafe.Sizeof(stPrintDlg{})

// dll bool ComDlg32:PageSetupDlg(PAGESETUPDLG* pagesetupdlg)
var pageSetupDlg = comdlg32.MustFindProc("PageSetupDlgA").Addr()
var _ = builtin(PageSetupDlg, "(pagesetupdlg)")

func PageSetupDlg(a Value) Value {
	psd := stPageSetupDlg{
		lStructSize:             uint32(nPageSetupDlg),
		ptPaperSize:             getPoint(a, "ptPaperSize"),
		rtMinMargin:             getRect(a, "rtMinMargin"),
		rtMargin:                getRect(a, "rtMargin"),
		hwndOwner:               getUintptr(a, "hwndOwner"),
		hDevMode:                getUintptr(a, "hDevMode"),
		hDevNames:               getUintptr(a, "hDevNames"),
		Flags:                   getInt32(a, "Flags"),
		hInstance:               getUintptr(a, "hInstance"),
		lCustData:               getUintptr(a, "lCustData"),
		lpfnPageSetupHook:       0,
		lpfnPagePaintHook:       0,
		lpPageSetupTemplateName: getZstr(a, "lpPageSetupTemplateName"),
		hPageSetupTemplate:      getUintptr(a, "hPageSetupTemplate"),
	}
	rtn, _, _ := syscall.SyscallN(pageSetupDlg,
		uintptr(unsafe.Pointer(&psd)))
	a.Put(nil, SuStr("hwndOwner"), IntVal(int(psd.hwndOwner)))
	a.Put(nil, SuStr("hDevMode"), IntVal(int(psd.hDevMode)))
	a.Put(nil, SuStr("hDevNames"), IntVal(int(psd.hDevNames)))
	a.Put(nil, SuStr("Flags"), IntVal(int(psd.Flags)))
	a.Put(nil, SuStr("ptPaperSize"),
		fromPoint(&psd.ptPaperSize, a.Get(nil, SuStr("ptPaperSize"))))
	a.Put(nil, SuStr("rtMinMargin"),
		fromRect(&psd.rtMinMargin, a.Get(nil, SuStr("rtMinMargin"))))
	a.Put(nil, SuStr("rtMargin"),
		fromRect(&psd.rtMargin, a.Get(nil, SuStr("rtMargin"))))
	a.Put(nil, SuStr("hInstance"), IntVal(int(psd.hInstance)))
	a.Put(nil, SuStr("lCustData"), IntVal(int(psd.lCustData)))
	a.Put(nil, SuStr("lpfnPageSetupHook"), IntVal(int(psd.lpfnPageSetupHook)))
	a.Put(nil, SuStr("lpfnPagePaintHook"), IntVal(int(psd.lpfnPagePaintHook)))
	a.Put(nil, SuStr("hPageSetupTemplate"), IntVal(int(psd.hPageSetupTemplate)))
	return boolRet(rtn)
}

type stPageSetupDlg struct {
	lStructSize             uint32
	hwndOwner               HANDLE
	hDevMode                HANDLE
	hDevNames               HANDLE
	Flags                   int32
	ptPaperSize             stPoint
	rtMinMargin             stRect
	rtMargin                stRect
	hInstance               HANDLE
	lCustData               uintptr
	lpfnPageSetupHook       HANDLE
	lpfnPagePaintHook       HANDLE
	lpPageSetupTemplateName *byte
	hPageSetupTemplate      HANDLE
}

const nPageSetupDlg = unsafe.Sizeof(stPageSetupDlg{})

// dll bool ComDlg32:GetSaveFileName(OPENFILENAME* ofn)
var getSaveFileName = comdlg32.MustFindProc("GetSaveFileNameA").Addr()
var _ = builtin(GetSaveFileName, "(a)")

func GetSaveFileName(a Value) Value {
	ofn, buf := buildOPENFILENAME(a)
	rtn, _, _ := syscall.SyscallN(getSaveFileName,
		uintptr(unsafe.Pointer(ofn)))
	if rtn != 0 {
		a.Put(nil, SuStr("file"), bufZstr(buf))
	}
	return boolRet(rtn)
}

// dll bool ComDlg32:GetOpenFileName(OPENFILENAME* ofn)
var getOpenFileName = comdlg32.MustFindProc("GetOpenFileNameA").Addr()
var _ = builtin(GetOpenFileName, "(a)")

func GetOpenFileName(a Value) Value {
	ofn, buf := buildOPENFILENAME(a)
	rtn, _, _ := syscall.SyscallN(getOpenFileName,
		uintptr(unsafe.Pointer(ofn)))
	if rtn != 0 {
		i := bytes.Index(buf, []byte{0, 0}) // double nul terminated
		if i == -1 {
			i = len(buf) - 2
		}
		a.Put(nil, SuStr("file"), SuStr(string(buf[:i+2])))
	}
	return boolRet(rtn)
}

func buildOPENFILENAME(a Value) (*stOpenFileName, []byte) {
	bufsize := getInt(a, "maxFile")
	buf := make([]byte, bufsize)
	copy(buf, ToStr(a.Get(nil, SuStr("file"))))
	ofn := &stOpenFileName{
		structSize: int32(nOpenFileName),
		hwndOwner:  getUintptr(a, "hwndOwner"),
		file:       &buf[0],
		maxFile:    int32(bufsize),
		filter:     getZstr(a, "filter"),
		flags:      getInt32(a, "flags"),
		defExt:     getZstr(a, "defExt"),
		initialDir: getZstr(a, "initialDir"),
		title:      getZstr(a, "title"),
	}
	return ofn, buf
}

type stOpenFileName struct {
	structSize     int32
	hwndOwner      HANDLE
	instance       HANDLE
	filter         *byte
	customFilter   *byte
	nMaxCustFilter int32
	nFilterIndex   int32
	file           *byte
	maxFile        int32
	fileTitle      *byte
	maxFileTitle   int32
	initialDir     *byte
	title          *byte
	flags          int32
	fileOffset     int16
	fileExtension  int16
	defExt         *byte
	custData       HANDLE
	hook           HANDLE
	templateName   *byte
	pvReserved     uintptr
	dwReserved     int32
	FlagsEx        int32
}

const nOpenFileName = unsafe.Sizeof(stOpenFileName{})

// dll bool ComDlg32:ChooseColor(CHOOSECOLOR* x)
var chooseColor = comdlg32.MustFindProc("ChooseColorA").Addr()
var _ = builtin(ChooseColor, "(x)")

func ChooseColor(a Value) Value {
	var custColors CustColors
	ccs := a.Get(nil, SuStr("custColors"))
	if ccs != nil {
		for i := range nCustColors {
			if x := ccs.Get(nil, SuInt(i)); x != nil {
				custColors[i] = int32(ToInt(x))
			}
		}
	}
	cc := stChooseColor{
		size:         int32(nChooseColor),
		owner:        getUintptr(a, "owner"),
		instance:     getUintptr(a, "instance"),
		rgbResult:    getInt32(a, "rgbResult"),
		custColors:   &custColors,
		flags:        getInt32(a, "flags"),
		custData:     getUintptr(a, "custData"),
		hook:         getUintptr(a, "hook"),
		templateName: getZstr(a, "templateName"),
	}
	rtn, _, _ := syscall.SyscallN(chooseColor,
		uintptr(unsafe.Pointer(&cc)))
	a.Put(nil, SuStr("rgbResult"), IntVal(int(cc.rgbResult)))
	a.Put(nil, SuStr("flags"), IntVal(int(cc.flags)))
	if ccs != nil {
		for i := range custColors {
			ccs.Put(nil, SuInt(i), IntVal(int(custColors[i])))
		}
	}
	return boolRet(rtn)
}

type stChooseColor struct {
	size         int32
	owner        HANDLE
	instance     HANDLE
	rgbResult    int32
	custColors   *CustColors
	flags        int32
	custData     HANDLE
	hook         HANDLE
	templateName *byte
}

const nChooseColor = unsafe.Sizeof(stChooseColor{})

const nCustColors = 16

type CustColors [nCustColors]int32

// dll bool ComDlg32:ChooseFont(CHOOSEFONT* cf)
var chooseFont = comdlg32.MustFindProc("ChooseFontA").Addr()
var _ = builtin(ChooseFont, "(cf)")

func ChooseFont(a Value) Value {
	lfob := a.Get(nil, SuStr("lpLogFont"))
	lf := stLogFont{
		lfHeight:         getInt32(lfob, "lfHeight"),
		lfWidth:          getInt32(lfob, "lfWidth"),
		lfEscapement:     getInt32(lfob, "lfEscapement"),
		lfOrientation:    getInt32(lfob, "lfOrientation"),
		lfWeight:         getInt32(lfob, "lfWeight"),
		lfItalic:         byte(getInt(lfob, "lfItalic")),
		lfUnderline:      byte(getInt(lfob, "lfUnderline")),
		lfStrikeOut:      byte(getInt(lfob, "lfStrikeOut")),
		lfCharSet:        byte(getInt(lfob, "lfCharSet")),
		lfOutPrecision:   byte(getInt(lfob, "lfOutPrecision")),
		lfClipPrecision:  byte(getInt(lfob, "lfClipPrecision")),
		lfQuality:        byte(getInt(lfob, "lfQuality")),
		lfPitchAndFamily: byte(getInt(lfob, "lfPitchAndFamily")),
	}
	getZstrBs(lfob, "lfFaceName", lf.lfFaceName[:])
	cf := stChooseFont{
		lStructSize:    uint32(nChooseFont),
		hwndOwner:      getUintptr(a, "hwndOwner"),
		hDC:            getUintptr(a, "hDC"),
		lpLogFont:      &lf,
		iPointSize:     getInt32(a, "iPointSize"),
		Flags:          getInt32(a, "Flags"),
		rgbColors:      getInt32(a, "rgbColors"),
		lCustData:      getUintptr(a, "lCustData"),
		lpfnHook:       getUintptr(a, "lpfnHook"),
		lpTemplateName: getZstr(a, "lpTemplateName"),
		hInstance:      getUintptr(a, "hInstance"),
		lpszStyle:      getZstr(a, "lpszStyle"),
		nFontType:      getInt16(a, "nFontType"),
		nSizeMin:       getInt32(a, "nSizeMin"),
		nSizeMax:       getInt32(a, "nSizeMax"),
	}
	rtn, _, _ := syscall.SyscallN(chooseFont,
		uintptr(unsafe.Pointer(&cf)))
	lfob.Put(nil, SuStr("lfHeight"), IntVal(int(lf.lfHeight)))
	lfob.Put(nil, SuStr("lfWidth"), IntVal(int(lf.lfWidth)))
	lfob.Put(nil, SuStr("lfEscapement"), IntVal(int(lf.lfEscapement)))
	lfob.Put(nil, SuStr("lfOrientation"), IntVal(int(lf.lfOrientation)))
	lfob.Put(nil, SuStr("lfWeight"), IntVal(int(lf.lfWeight)))
	lfob.Put(nil, SuStr("lfItalic"), IntVal(int(lf.lfItalic)))
	lfob.Put(nil, SuStr("lfUnderline"), IntVal(int(lf.lfUnderline)))
	lfob.Put(nil, SuStr("lfStrikeOut"), IntVal(int(lf.lfStrikeOut)))
	lfob.Put(nil, SuStr("lfCharSet"), IntVal(int(lf.lfCharSet)))
	lfob.Put(nil, SuStr("lfOutPrecision"), IntVal(int(lf.lfOutPrecision)))
	lfob.Put(nil, SuStr("lfClipPrecision"), IntVal(int(lf.lfClipPrecision)))
	lfob.Put(nil, SuStr("lfQuality"), IntVal(int(lf.lfQuality)))
	lfob.Put(nil, SuStr("lfPitchAndFamily"), IntVal(int(lf.lfPitchAndFamily)))
	lfob.Put(nil, SuStr("lfPitchAndFamily"), IntVal(int(lf.lfPitchAndFamily)))
	lfob.Put(nil, SuStr("lfFaceName"), bufZstr(lf.lfFaceName[:]))
	return boolRet(rtn)
}

type stChooseFont struct {
	lStructSize    uint32
	hwndOwner      HANDLE
	hDC            HANDLE
	lpLogFont      *stLogFont
	iPointSize     int32
	Flags          int32
	rgbColors      int32
	lCustData      uintptr
	lpfnHook       HANDLE
	lpTemplateName *byte
	hInstance      HANDLE
	lpszStyle      *byte
	nFontType      int16
	_              int16 // padding
	nSizeMin       int32
	nSizeMax       int32
	_              int32 // padding
}

const nChooseFont = unsafe.Sizeof(stChooseFont{})

// dll HRESULT ComDlg32:PrintDlgEx(PRINTDLGEX* printdlgex)
var printDlgEx = comdlg32.MustFindProc("PrintDlgExA").Addr()
var _ = builtin(PrintDlgEx, "(printdlgex)")

func PrintDlgEx(a Value) Value {
	pd := stPrintDlgEx{
		lStructSize:         int32(nPrintDlgEx),
		hwndOwner:           getUintptr(a, "hwndOwner"),
		hDevMode:            getUintptr(a, "hDevMode"),
		hDevNames:           getUintptr(a, "hDevNames"),
		hDC:                 getUintptr(a, "hDC"),
		Flags:               getInt32(a, "Flags"),
		Flags2:              getInt32(a, "Flags2"),
		ExclusionFlags:      getInt32(a, "ExclusionFlags"),
		nMinPage:            getInt32(a, "nMinPage"),
		nMaxPage:            getInt32(a, "nMaxPage"),
		nCopies:             getInt32(a, "nCopies"),
		hInstance:           getUintptr(a, "hInstance"),
		lpPrintTemplateName: getZstr(a, "lpPrintTemplateName"),
		nStartPage:          getInt32(a, "nStartPage"),
		dwResultAction:      getInt32(a, "dwResultAction"),
	}
	prob := a.Get(nil, SuStr("lpPageRanges"))
	var pr *stPrintPageRange
	if prob != nil {
		pr = &stPrintPageRange{
			nFromPage: getInt32(prob, "nFromPage"),
			nToPage:   getInt32(prob, "nToPage"),
		}
		pd.lpPageRanges = pr
		pd.nPageRanges = 1
		pd.nMaxPageRanges = 1
	}
	rtn, _, _ := syscall.SyscallN(printDlgEx,
		uintptr(unsafe.Pointer(&pd)))
	a.Put(nil, SuStr("hwndOwner"), IntVal(int(pd.hwndOwner)))
	a.Put(nil, SuStr("hDevMode"), IntVal(int(pd.hDevMode)))
	a.Put(nil, SuStr("hDevNames"), IntVal(int(pd.hDevNames)))
	a.Put(nil, SuStr("hDC"), IntVal(int(pd.hDC)))
	a.Put(nil, SuStr("Flags"), IntVal(int(pd.Flags)))
	a.Put(nil, SuStr("Flags2"), IntVal(int(pd.Flags2)))
	a.Put(nil, SuStr("ExclusionFlags"), IntVal(int(pd.ExclusionFlags)))
	a.Put(nil, SuStr("nPageRanges"), IntVal(int(pd.nPageRanges)))
	a.Put(nil, SuStr("nMaxPageRanges"), IntVal(int(pd.nMaxPageRanges)))
	if prob != nil {
		prob.Put(nil, SuStr("nFromPage"), IntVal(int(pr.nFromPage)))
		prob.Put(nil, SuStr("nToPage"), IntVal(int(pr.nToPage)))
	}
	a.Put(nil, SuStr("nMinPage"), IntVal(int(pd.nMinPage)))
	a.Put(nil, SuStr("nMaxPage"), IntVal(int(pd.nMaxPage)))
	a.Put(nil, SuStr("nCopies"), IntVal(int(pd.nCopies)))
	a.Put(nil, SuStr("hInstance"), IntVal(int(pd.hInstance)))
	a.Put(nil, SuStr("nStartPage"), IntVal(int(pd.nStartPage)))
	a.Put(nil, SuStr("dwResultAction"), IntVal(int(pd.dwResultAction)))
	return intRet(rtn)
}

type stPrintDlgEx struct {
	lStructSize         int32
	hwndOwner           HANDLE
	hDevMode            HANDLE
	hDevNames           HANDLE
	hDC                 HANDLE
	Flags               int32
	Flags2              int32
	ExclusionFlags      int32
	nPageRanges         int32
	nMaxPageRanges      int32
	lpPageRanges        *stPrintPageRange
	nMinPage            int32
	nMaxPage            int32
	nCopies             int32
	hInstance           HANDLE
	lpPrintTemplateName *byte
	lpCallback          uintptr
	nPropertyPages      int32
	lphPropertyPages    uintptr
	nStartPage          int32
	dwResultAction      int32
}

const nPrintDlgEx = unsafe.Sizeof(stPrintDlgEx{})

type stPrintPageRange struct {
	nFromPage int32
	nToPage   int32
}

const nPrintPageRange = unsafe.Sizeof(stPrintPageRange{})
