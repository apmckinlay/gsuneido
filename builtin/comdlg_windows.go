package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var comdlg32 = windows.NewLazyDLL("comdlg32.dll")

// dll long ComDlg32:CommDlgExtendedError()
var commDlgExtendedError = comdlg32.NewProc("CommDlgExtendedError")
var _ = builtin0("CommDlgExtendedError()",
	func() Value {
		rtn, _, _ := commDlgExtendedError.Call()
		return intRet(rtn)
	})

// dll bool ComDlg32:PrintDlg(PRINTDLG* printdlg)
var printDlg = comdlg32.NewProc("PrintDlgA")
var _ = builtin1("PrintDlg(printdlg)",
	func(a Value) Value {
		pd := PRINTDLG{
			lStructSize:         uint32(unsafe.Sizeof(PRINTDLG{})),
			hwndOwner:           getHandle(a, "hwndOwner"),
			hDevMode:            getHandle(a, "hDevMode"),
			hDevNames:           getHandle(a, "hDevNames"),
			hDC:                 getHandle(a, "hDC"),
			Flags:               getInt32(a, "Flags"),
			nFromPage:           getInt16(a, "nFromPage"),
			nToPage:             getInt16(a, "nToPage"),
			nMinPage:            getInt16(a, "nMinPage"),
			nMaxPage:            getInt16(a, "nMaxPage"),
			nCopies:             getInt16(a, "nCopies"),
			hInstance:           getHandle(a, "hInstance"),
			lCustData:           getHandle(a, "lCustData"),
			lpPrintTemplateName: getStr(a, "lpPrintTemplateName"),
			lpSetupTemplateName: getStr(a, "lpSetupTemplateName"),
			hPrintTemplate:      getHandle(a, "hPrintTemplate"),
			hSetupTemplate:      getHandle(a, "hSetupTemplate"),
		}
		rtn, _, _ := printDlg.Call(
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
	})

type PRINTDLG struct {
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

// dll bool ComDlg32:PageSetupDlg(PAGESETUPDLG* pagesetupdlg)
var pageSetupDlg = comdlg32.NewProc("PageSetupDlgA")
var _ = builtin1("PageSetupDlg(pagesetupdlg)",
	func(a Value) Value {
		psob := a.Get(nil, SuStr("ptPaperSize"))
		mmob := a.Get(nil, SuStr("rtMinMargin"))
		rmob := a.Get(nil, SuStr("rtMargin"))
		psd := PAGESETUPDLG{
			lStructSize:             uint32(unsafe.Sizeof(PAGESETUPDLG{})),
			ptPaperSize:             getPoint(psob, "ptPaperSize"),
			rtMinMargin:             getRect(mmob, "rtMinMargin"),
			rtMargin:                getRect(rmob, "rtMargin"),
			hwndOwner:               getHandle(a, "hwndOwner"),
			hDevMode:                getHandle(a, "hDevMode"),
			hDevNames:               getHandle(a, "hDevNames"),
			Flags:                   getInt32(a, "Flags"),
			hInstance:               getHandle(a, "hInstance"),
			lCustData:               getHandle(a, "lCustData"),
			lpfnPageSetupHook:       0,
			lpfnPagePaintHook:       0,
			lpPageSetupTemplateName: getStr(a, "lpPageSetupTemplateName"),
			hPageSetupTemplate:      getHandle(a, "hPageSetupTemplate"),
		}
		rtn, _, _ := pageSetupDlg.Call(
			uintptr(unsafe.Pointer(&psd)))
		a.Put(nil, SuStr("hwndOwner"), IntVal(int(psd.hwndOwner)))
		a.Put(nil, SuStr("hDevMode"), IntVal(int(psd.hDevMode)))
		a.Put(nil, SuStr("hDevNames"), IntVal(int(psd.hDevNames)))
		a.Put(nil, SuStr("Flags"), IntVal(int(psd.Flags)))
		a.Put(nil, SuStr("ptPaperSize"), pointToOb(&psd.ptPaperSize, psob))
		a.Put(nil, SuStr("rtMinMargin"), rectToOb(&psd.rtMinMargin, mmob))
		a.Put(nil, SuStr("rtMargin"), rectToOb(&psd.rtMargin, rmob))
		a.Put(nil, SuStr("hInstance"), IntVal(int(psd.hInstance)))
		a.Put(nil, SuStr("lCustData"), IntVal(int(psd.lCustData)))
		a.Put(nil, SuStr("lpfnPageSetupHook"), IntVal(int(psd.lpfnPageSetupHook)))
		a.Put(nil, SuStr("lpfnPagePaintHook"), IntVal(int(psd.lpfnPagePaintHook)))
		a.Put(nil, SuStr("hPageSetupTemplate"), IntVal(int(psd.hPageSetupTemplate)))
		return boolRet(rtn)
	})

type PAGESETUPDLG struct {
	lStructSize             uint32
	hwndOwner               HANDLE
	hDevMode                HANDLE
	hDevNames               HANDLE
	Flags                   int32
	ptPaperSize             POINT
	rtMinMargin             RECT
	rtMargin                RECT
	hInstance               HANDLE
	lCustData               uintptr
	lpfnPageSetupHook       HANDLE
	lpfnPagePaintHook       HANDLE
	lpPageSetupTemplateName *byte
	hPageSetupTemplate      HANDLE
}

// dll bool ComDlg32:GetSaveFileName(OPENFILENAME* ofn)
var getSaveFileName = comdlg32.NewProc("GetSaveFileNameA")
var _ = builtin1("GetSaveFileName(a)",
	func(a Value) Value {
		const bufsize = 8192
		var buf [bufsize + 1]byte
		file := ToStr(a.Get(nil, SuStr("file")))
		copyStr(buf[:], file)
		buf[len(file)+1] = 0 // need double nuls
		ofn := OPENFILENAME{
			structSize: int32(unsafe.Sizeof(OPENFILENAME{})),
			file:       &buf[0],
			maxFile:    bufsize,
			filter:     getStr(a, "filter"),
			flags:      getInt32(a, "flags"),
			defExt:     getStr(a, "defExt"),
			initialDir: getStr(a, "initialDir"),
		}
		rtn, _, _ := getSaveFileName.Call(uintptr(unsafe.Pointer(&ofn)))
		if rtn == 0 {
			return EmptyStr
		}
		return strRet(buf[:])
	})

type OPENFILENAME struct {
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

// dll bool ComDlg32:ChooseColor(CHOOSECOLOR* x)
var chooseColor = comdlg32.NewProc("ChooseColorA")
var _ = builtin1("ChooseColor(x)",
	func(a Value) Value {
		var custColors CustColors
		ccs := a.Get(nil, SuStr("custColors"))
		for i := 0; i < nCustColors; i++ {
			custColors[i] = int32(ToInt(ccs.Get(nil, SuInt(i))))
		}
		cc := CHOOSECOLOR{
			size:       getInt32(a, "size"),
			owner:      getHandle(a, "owner"),
			flags:      getInt32(a, "flags"),
			resource:   getStr(a, "resource"),
			custColors: &custColors,
		}
		rtn, _, _ := chooseColor.Call(
			uintptr(unsafe.Pointer(&cc)))
		a.Put(nil, SuStr("rgbResult"), IntVal(int(cc.rgbResult)))
		a.Put(nil, SuStr("flags"), IntVal(int(cc.flags)))
		for i := 0; i < nCustColors; i++ {
			ccs.Put(nil, SuInt(i), IntVal(int(custColors[i])))
		}
		return boolRet(rtn)
	})

type CHOOSECOLOR struct {
	size       int32
	owner      HANDLE
	instance   HANDLE
	rgbResult  int32
	custColors *CustColors
	flags      int32
	custData   HANDLE
	hook       HANDLE
	resource   *byte
}

const nCustColors = 16

type CustColors [nCustColors]int32

// dll bool ComDlg32:ChooseFont(CHOOSEFONT* cf)
var chooseFont = comdlg32.NewProc("ChooseFontA")
var _ = builtin1("ChooseFont(cf)",
	func(a Value) Value {
		lfob := a.Get(nil, SuStr("lpLogFont"))
		lf := LOGFONT{
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
		copyStr(lf.lfFaceName[:], ToStr(lfob.Get(nil, SuStr("lfFaceName"))))
		cf := CHOOSEFONT{
			lStructSize:    uint32(unsafe.Sizeof(CHOOSEFONT{})),
			hwndOwner:      getHandle(a, "hwndOwner"),
			hDC:            getHandle(a, "hDC"),
			lpLogFont:      &lf,
			iPointSize:     getInt32(a, "iPointSize"),
			Flags:          getInt32(a, "Flags"),
			rgbColors:      getInt32(a, "rgbColors"),
			lCustData:      getHandle(a, "lCustData"),
			lpfnHook:       getHandle(a, "lpfnHook"),
			lpTemplateName: getStr(a, "lpTemplateName"),
			hInstance:      getHandle(a, "hInstance"),
			lpszStyle:      getStr(a, "lpszStyle"),
			nFontType:      getInt16(a, "nFontType"),
			nSizeMin:       getInt32(a, "nSizeMin"),
			nSizeMax:       getInt32(a, "nSizeMax"),
		}
		rtn, _, _ := chooseFont.Call(uintptr(unsafe.Pointer(&cf)))
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
		return boolRet(rtn)
	})

type CHOOSEFONT struct {
	lStructSize    uint32
	hwndOwner      HANDLE
	hDC            HANDLE
	lpLogFont      *LOGFONT
	iPointSize     int32
	Flags          int32
	rgbColors      int32
	lCustData      uintptr
	lpfnHook       HANDLE
	lpTemplateName *byte
	hInstance      HANDLE
	lpszStyle      *byte
	nFontType      int16
	_              [4]byte // padding
	nSizeMin       int32
	nSizeMax       int32
}
