// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/verify"
)

var gdi32 = MustLoadDLL("gdi32.dll")

const LF_FACESIZE = 32

type LOGFONT struct {
	lfHeight         int32
	lfWidth          int32
	lfEscapement     int32
	lfOrientation    int32
	lfWeight         int32
	lfItalic         byte
	lfUnderline      byte
	lfStrikeOut      byte
	lfCharSet        byte
	lfOutPrecision   byte
	lfClipPrecision  byte
	lfQuality        byte
	lfPitchAndFamily byte
	lfFaceName       [LF_FACESIZE]byte
}

const nLOGFONT = unsafe.Sizeof(LOGFONT{})

type TEXTMETRIC struct {
	Height           int32
	Ascent           int32
	Descent          int32
	InternalLeading  int32
	ExternalLeading  int32
	AveCharWidth     int32
	MaxCharWidth     int32
	Weight           int32
	Overhang         int32
	DigitizedAspectX int32
	DigitizedAspectY int32
	FirstChar        byte
	LastChar         byte
	DefaultChar      byte
	BreakChar        byte
	Italic           byte
	Underlined       byte
	StruckOut        byte
	PitchAndFamily   byte
	CharSet          byte
	_                [3]byte // padding
}

const nTEXTMETRIC = unsafe.Sizeof(TEXTMETRIC{})

// dll Gdi32:CreateFontIndirect(LOGFONT* lf) gdiobj
var createFontIndirect = gdi32.MustFindProc("CreateFontIndirectA").Addr()
var _ = builtin1("CreateFontIndirect(logfont)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nLOGFONT)
		f := (*LOGFONT)(p)
		*f = LOGFONT{
			lfHeight:         getInt32(a, "lfHeight"),
			lfWidth:          getInt32(a, "lfWidth"),
			lfEscapement:     getInt32(a, "lfEscapement"),
			lfOrientation:    getInt32(a, "lfOrientation"),
			lfWeight:         getInt32(a, "lfWeight"),
			lfItalic:         byte(getInt(a, "lfItalic")),
			lfUnderline:      byte(getInt(a, "lfUnderline")),
			lfStrikeOut:      byte(getInt(a, "lfStrikeOut")),
			lfCharSet:        byte(getInt(a, "lfCharSet")),
			lfOutPrecision:   byte(getInt(a, "lfOutPrecision")),
			lfClipPrecision:  byte(getInt(a, "lfClipPrecision")),
			lfQuality:        byte(getInt(a, "lfQuality")),
			lfPitchAndFamily: byte(getInt(a, "lfPitchAndFamily")),
		}
		getStrZbs(a, "lfFaceName", f.lfFaceName[:])
		rtn := goc.Syscall1(createFontIndirect,
			uintptr(p))
		return intRet(rtn)
	})

// dll bool Gdi32:GetTextMetrics(pointer hdc, TEXTMETRIC* tm)
var getTextMetrics = gdi32.MustFindProc("GetTextMetricsA").Addr()
var _ = builtin2("GetTextMetrics(hdc, tm)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nTEXTMETRIC)
		rtn := goc.Syscall2(getTextMetrics,
			intArg(a),
			uintptr(p))
		tm := (*TEXTMETRIC)(p)
		b.Put(nil, SuStr("Height"), IntVal(int(tm.Height)))
		b.Put(nil, SuStr("Ascent"), IntVal(int(tm.Ascent)))
		b.Put(nil, SuStr("Descent"), IntVal(int(tm.Descent)))
		b.Put(nil, SuStr("InternalLeading"), IntVal(int(tm.InternalLeading)))
		b.Put(nil, SuStr("ExternalLeading"), IntVal(int(tm.ExternalLeading)))
		b.Put(nil, SuStr("AveCharWidth"), IntVal(int(tm.AveCharWidth)))
		b.Put(nil, SuStr("MaxCharWidth"), IntVal(int(tm.MaxCharWidth)))
		b.Put(nil, SuStr("Italic"), IntVal(int(tm.Italic)))
		b.Put(nil, SuStr("Underlined"), IntVal(int(tm.Underlined)))
		b.Put(nil, SuStr("StruckOut"), IntVal(int(tm.StruckOut)))
		b.Put(nil, SuStr("FirstChar"), IntVal(int(tm.FirstChar)))
		b.Put(nil, SuStr("LastChar"), IntVal(int(tm.LastChar)))
		b.Put(nil, SuStr("DefaultChar"), IntVal(int(tm.DefaultChar)))
		b.Put(nil, SuStr("BreakChar"), IntVal(int(tm.BreakChar)))
		b.Put(nil, SuStr("PitchAndFamily"), IntVal(int(tm.PitchAndFamily)))
		b.Put(nil, SuStr("CharSet"), IntVal(int(tm.CharSet)))
		b.Put(nil, SuStr("Overhang"), IntVal(int(tm.Overhang)))
		b.Put(nil, SuStr("DigitizedAspectX"), IntVal(int(tm.DigitizedAspectX)))
		b.Put(nil, SuStr("DigitizedAspectY"), IntVal(int(tm.DigitizedAspectY)))
		return boolRet(rtn)
	})

var getStockObject = gdi32.MustFindProc("GetStockObject").Addr()
var _ = builtin1("GetStockObject(i)",
	func(a Value) Value {
		rtn := goc.Syscall1(getStockObject,
			intArg(a))
		return intRet(rtn)
	})

var getDeviceCaps = gdi32.MustFindProc("GetDeviceCaps").Addr()
var _ = builtin2("GetDeviceCaps(hdc, nIndex)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getDeviceCaps,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Gdi32:BitBlt(pointer hdcDest,
//		long nXDest, long nYDest, long nWidth, long nHeight,
//		pointer hdcSrc, long nXSrc, long nYSrc, long dwRop) bool
var bitBlt = gdi32.MustFindProc("BitBlt").Addr()
var _ = builtin("BitBlt(hdc, x, y, cx, cy, hdcSrc, x1, y1, rop)",
	func(_ *Thread, a []Value) Value {
		rtn := goc.Syscall9(bitBlt,
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			intArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			intArg(a[8]))
		return boolRet(rtn)
	})

// dll Gdi32:CreateCompatibleBitmap(pointer hdc, long nWidth, long nHeight) gdiobj
var createCompatibleBitmap = gdi32.MustFindProc("CreateCompatibleBitmap").Addr()
var _ = builtin3("CreateCompatibleBitmap(hdc, cx, cy)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(createCompatibleBitmap,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll Gdi32:CreateCompatibleDC(pointer hdc) pointer
var createCompatibleDC = gdi32.MustFindProc("CreateCompatibleDC").Addr()
var _ = builtin1("CreateCompatibleDC(hdc)",
	func(a Value) Value {
		rtn := goc.Syscall1(createCompatibleDC,
			intArg(a))
		return intRet(rtn)
	})

// dll Gdi32:CreateSolidBrush(long rgb) gdiobj
var createSolidBrush = gdi32.MustFindProc("CreateSolidBrush").Addr()
var _ = builtin1("CreateSolidBrush(i)",
	func(a Value) Value {
		rtn := goc.Syscall1(createSolidBrush,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer Gdi32:SelectObject(pointer hdc, pointer obj)
var selectObject = gdi32.MustFindProc("SelectObject").Addr()
var _ = builtin2("SelectObject(hdc, obj)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(selectObject,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll bool Gdi32:GetTextExtentPoint32(pointer hdc, [in] string text, long len, POINT* p)
var getTextExtentPoint32 = gdi32.MustFindProc("GetTextExtentPoint32A").Addr()
var _ = builtin4("GetTextExtentPoint32(hdc, text, len, p)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		pt := heap.Alloc(nPOINT)
		rtn := goc.Syscall4(getTextExtentPoint32,
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(len(ToStr(b))),
			uintptr(pt))
		upointToOb(pt, d)
		return boolRet(rtn)
	})

var getTextFace = gdi32.MustFindProc("GetTextFaceA").Addr()
var _ = builtin1("GetTextFace(hdc)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		const bufsize = 512
		p := heap.Alloc(bufsize)
		n := goc.Syscall3(getTextFace,
			intArg(a),
			bufsize,
			uintptr(p))
		return SuStr(heap.GetStrN(p, int(n)))
	})

// dll long Gdi32:SetBkMode(pointer hdc, long mode)
var setBkMode = gdi32.MustFindProc("SetBkMode").Addr()
var _ = builtin2("SetBkMode(hdc, color)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setBkMode,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll long Gdi32:SetBkColor(pointer hdc, long color)
var setBkColor = gdi32.MustFindProc("SetBkColor").Addr()
var _ = builtin2("SetBkColor(hdc, color)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setBkColor,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Gdi32:DeleteDC(pointer hdc) bool
var deleteDC = gdi32.MustFindProc("DeleteDC").Addr()
var _ = builtin1("DeleteDC(hdc)",
	func(a Value) Value {
		rtn := goc.Syscall1(deleteDC,
			intArg(a))
		return boolRet(rtn)
	})

// dll Gdi32:DeleteObject(pointer hgdiobj) bool
var deleteObject = gdi32.MustFindProc("DeleteObject").Addr()
var _ = builtin1("DeleteObject(hgdiobj)",
	func(a Value) Value {
		rtn := goc.Syscall1(deleteObject,
			intArg(a))
		return boolRet(rtn)
	})

// dll Gdi32:GetClipBox(pointer hdc, RECT* rect) long
var getClipBox = gdi32.MustFindProc("GetClipBox").Addr()
var _ = builtin2("GetClipBox(hdc, rect)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall3(getClipBox,
			intArg(a),
			uintptr(rectArg(b, r)),
			0)
		urectToOb(r, b)
		return intRet(rtn)
	})

// dll Gdi32:GetPixel(pointer hdc, long nXPos, long nYPos) long
var getPixel = gdi32.MustFindProc("GetPixel").Addr()
var _ = builtin3("GetPixel(hdc, nXPos, nYPos)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(getPixel,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll Gdi32:PtVisible(pointer hdc, long nXPos, long nYPos) bool
var ptVisible = gdi32.MustFindProc("PtVisible").Addr()
var _ = builtin3("PtVisible(hdc, nXPos, nYPos)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(ptVisible,
			intArg(a),
			intArg(b),
			intArg(c))
		if rtn == 0xffffffff { // 0xffffffff: 32 bit -1
			return IntVal(-1)
		}
		return boolRet(rtn)
	})

// dll Gdi32:Rectangle(pointer hdc, long left, long top, long right, long bottom) bool
var rectangle = gdi32.MustFindProc("Rectangle").Addr()
var _ = builtin5("Rectangle(hdc, left, top, right, bottom)",
	func(a, b, c, d, e Value) Value {
		rtn := goc.Syscall5(rectangle,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return boolRet(rtn)
	})

// dll Gdi32:SetStretchBltMode(pointer hdc, long iStretchMode) long
var setStretchBltMode = gdi32.MustFindProc("SetStretchBltMode").Addr()
var _ = builtin2("SetStretchBltMode(hdc, iStretchMode)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setStretchBltMode,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Gdi32:SetTextColor(pointer hdc, long color) long
var setTextColor = gdi32.MustFindProc("SetTextColor").Addr()
var _ = builtin2("SetTextColor(hdc, color)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setTextColor,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Gdi32:StretchBlt(pointer hdcDest, long nXOriginDest, long nYOriginDest, long nWidthDest, long nHeightDest, pointer hdcSrc,
//	long nXOriginSrc, long nYOriginSrc, long nWidthSrc, long nHeightSrc, long dwRop) bool
var stretchBlt = gdi32.MustFindProc("StretchBlt").Addr()
var _ = builtin("StretchBlt(hdcDest, nXOriginDest, nYOriginDest, nWidthDest, nHeightDest, hdcSrc, nXOriginSrc, nYOriginSrc, nWidthSrc, nHeightSrc, dwRop)",
	func(_ *Thread, a []Value) Value {
		rtn := goc.Syscall11(stretchBlt,
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			intArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			intArg(a[8]),
			intArg(a[9]),
			intArg(a[10]))
		return boolRet(rtn)
	})

// dll pointer Gdi32:CloseEnhMetaFile(pointer dc)
var closeEnhMetaFile = gdi32.MustFindProc("CloseEnhMetaFile").Addr()
var _ = builtin1("CloseEnhMetaFile(dc)",
	func(a Value) Value {
		rtn := goc.Syscall1(closeEnhMetaFile,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer Gdi32:DeleteEnhMetaFile(pointer emf)
var deleteEnhMetaFile = gdi32.MustFindProc("DeleteEnhMetaFile").Addr()
var _ = builtin1("DeleteEnhMetaFile(dc)",
	func(a Value) Value {
		rtn := goc.Syscall1(deleteEnhMetaFile,
			intArg(a))
		return intRet(rtn)
	})

// dll pointer Gdi32:CreateDC(
// 	[in] string driver,
// 	[in] string device,
// 	[in] string output,
// 	pointer devmode)
var createDC = gdi32.MustFindProc("CreateDCA").Addr()
var _ = builtin4("CreateDC(driver, device, output, devmode)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall4(createDC,
			uintptr(stringArg(a)),
			uintptr(stringArg(b)),
			uintptr(stringArg(c)),
			intArg(d))
		return intRet(rtn)
	})

// dll bool Gdi32:Ellipse(pointer hdc, long left, long top, long right, long bottom)
var ellipse = gdi32.MustFindProc("Ellipse").Addr()
var _ = builtin5("Ellipse(hdc, left, top, right, bottom)",
	func(a, b, c, d, e Value) Value {
		rtn := goc.Syscall5(ellipse,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return boolRet(rtn)
	})

// dll long Gdi32:EndDoc(pointer hdc)
var endDoc = gdi32.MustFindProc("EndDoc").Addr()
var _ = builtin1("EndDoc(hdc)",
	func(a Value) Value {
		rtn := goc.Syscall1(endDoc,
			intArg(a))
		return intRet(rtn)
	})

// dll long Gdi32:EndPage(pointer hdc)
var endPage = gdi32.MustFindProc("EndPage").Addr()
var _ = builtin1("EndPage(hdc)",
	func(a Value) Value {
		rtn := goc.Syscall1(endPage,
			intArg(a))
		return intRet(rtn)
	})

// dll long Gdi32:ExcludeClipRect(pointer hdc, long l, long t, long r, long b)
var excludeClipRect = gdi32.MustFindProc("ExcludeClipRect").Addr()
var _ = builtin5("ExcludeClipRect(hdc, l, t, r, b)",
	func(a, b, c, d, e Value) Value {
		rtn := goc.Syscall5(excludeClipRect,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return intRet(rtn)
	})

// dll long Gdi32:GetClipRgn(pointer hdc, pointer hrgn)
var getClipRgn = gdi32.MustFindProc("GetClipRgn").Addr()
var _ = builtin2("GetClipRgn(hdc, hrgn)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getClipRgn,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer Gdi32:GetCurrentObject(pointer hdc, long uObjectType)
var getCurrentObject = gdi32.MustFindProc("GetCurrentObject").Addr()
var _ = builtin2("GetCurrentObject(hdc, uObjectType)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(getCurrentObject,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer Gdi32:GetEnhMetaFile(string filename)
var getEnhMetaFile = gdi32.MustFindProc("GetEnhMetaFileA").Addr()
var _ = builtin1("GetEnhMetaFile(filename)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall1(getEnhMetaFile,
			uintptr(stringArg(a)))
		return intRet(rtn)
	})

// dll bool Gdi32:LineTo(pointer hdc, long x, long y)
var lineTo = gdi32.MustFindProc("LineTo").Addr()
var _ = builtin3("LineTo(hdc, x, y)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(lineTo,
			intArg(a),
			intArg(b),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool Gdi32:PatBlt(
// 	pointer	hdc,        // Destination device context
// 	long	nXLeft,     // x-coordinate of upper left corner of destination rectangle
// 	long	nYLeft,     // y-coordinate of upper left corner of destination rectangle
// 	long	nWidth,     // width of destination rectangle
// 	long	nHeight,    // height of destination rectangle
// 	long	dwRop       // Raster operation
// 	)
var patBlt = gdi32.MustFindProc("PatBlt").Addr()
var _ = builtin6("PatBlt(hdc, nXLeft, nYLeft, nWidth, nHeight, dwRop)",
	func(a, b, c, d, e, f Value) Value {
		rtn := goc.Syscall6(patBlt,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f))
		return boolRet(rtn)
	})

// dll bool Gdi32:Polygon(
// 	pointer hdc,		// handle to device context
// 	[in] string lppt,		// array of points
// 	long cCount		// count of points
// 	)
var polygon = gdi32.MustFindProc("Polygon").Addr()
var _ = builtin3("Polygon(hdc, lppt, cCount)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall3(polygon,
			intArg(a),
			uintptr(stringArg(b)),
			intArg(c))
		return boolRet(rtn)
	})

// dll bool Gdi32:RestoreDC(
// 	pointer	hdc,        // handle to DC
// 	long	nSavedDC    // restore state returned by SaveDC
// )
var restoreDC = gdi32.MustFindProc("RestoreDC").Addr()
var _ = builtin2("RestoreDC(hdc, nSavedDC)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(restoreDC,
			intArg(a),
			intArg(b))
		return boolRet(rtn)
	})

// dll bool Gdi32:RoundRect(pointer hdc, long left, long top, long right, long bottom,
// 	long ellipse_width, long ellipse_height)
var roundRect = gdi32.MustFindProc("RoundRect").Addr()
var _ = builtin7("RoundRect(hdc, left, top, right, bottom, ellipse_width, ellipse_height)",
	func(a, b, c, d, e, f, g Value) Value {
		rtn := goc.Syscall7(roundRect,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			intArg(f),
			intArg(g))
		return boolRet(rtn)
	})

// dll long Gdi32:SaveDC(pointer hdc)
var saveDC = gdi32.MustFindProc("SaveDC").Addr()
var _ = builtin1("SaveDC(hdc)",
	func(a Value) Value {
		rtn := goc.Syscall1(saveDC,
			intArg(a))
		return intRet(rtn)
	})

// dll long Gdi32:SelectClipRgn(pointer hdc, pointer hrgn)
var selectClipRgn = gdi32.MustFindProc("SelectClipRgn").Addr()
var _ = builtin2("SelectClipRgn(hdc, hrgn)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(selectClipRgn,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll pointer Gdi32:SetEnhMetaFileBits(long cbBuffer, [in] string lpData)
var setEnhMetaFileBits = gdi32.MustFindProc("SetEnhMetaFileBits").Addr()
var _ = builtin2("SetEnhMetaFileBits(cbBuffer, lpData)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall2(setEnhMetaFileBits,
			intArg(a),
			uintptr(stringArg(b)))
		return intRet(rtn)
	})

// dll long Gdi32:SetMapMode(pointer hdc, long mode)
var setMapMode = gdi32.MustFindProc("SetMapMode").Addr()
var _ = builtin2("SetMapMode(hdc, mode)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setMapMode,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll long Gdi32:SetROP2(pointer hdc, long fnDrawMode)
var setROP2 = gdi32.MustFindProc("SetROP2").Addr()
var _ = builtin2("SetROP2(hdc, fnDrawMode)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setROP2,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll long Gdi32:SetTextAlign(pointer hdc, long mode)
var setTextAlign = gdi32.MustFindProc("SetTextAlign").Addr()
var _ = builtin2("SetTextAlign(hdc, mode)",
	func(a, b Value) Value {
		rtn := goc.Syscall2(setTextAlign,
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll long Gdi32:StartPage(pointer hdc)
var startPage = gdi32.MustFindProc("StartPage").Addr()
var _ = builtin1("StartPage(hdc)",
	func(a Value) Value {
		rtn := goc.Syscall1(startPage,
			intArg(a))
		return intRet(rtn)
	})

// dll bool Gdi32:TextOut(pointer hdc, long x, long y, [in] string text, long n)
var textOut = gdi32.MustFindProc("TextOutA").Addr()
var _ = builtin5("TextOut(hdc, x, y, text, n)",
	func(a, b, c, d, e Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall5(textOut,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(stringArg(d)),
			intArg(e))
		return boolRet(rtn)
	})

// dll bool Gdi32:Arc(pointer hdc, long nLeftRect, long nTopRect,
//		long nRightRect, long nBottomRect, long nXStartArc, long nYStartArc,
//		long nXEndArc, long nYEndArc)
var arc = gdi32.MustFindProc("Arc").Addr()
var _ = builtin("Arc(hdc, nLeftRect, nTopRect, nRightRect, nBottomRect,"+
	"nXStartArc, nYStartArc, nXEndArc, nYEndArc)",
	func(_ *Thread, a []Value) Value {
		rtn := goc.Syscall9(arc,
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			intArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			intArg(a[8]))
		return boolRet(rtn)
	})

// dll pointer Gdi32:CreateEnhMetaFile(pointer hdcRef, [in] string filename,
//		RECT* rect, [in] string desc)
var createEnhMetaFile = gdi32.MustFindProc("CreateEnhMetaFileA").Addr()
var _ = builtin4("CreateEnhMetaFile(hdcRef, filename, rect, desc)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall4(createEnhMetaFile,
			intArg(a),
			uintptr(stringArg(b)),
			uintptr(rectArg(c, r)),
			uintptr(stringArg(d)))
		return intRet(rtn)
	})

// dll gdiobj Gdi32:CreateFont(long nHeight, long nWidth, long nEscapement,
//		long nOrientation, long fnWeight, long fdwItalic, long fdwUnderline,
//		long fdwStrikeOut, long fdwCharSet, long fdwOutputPrecision,
//		long fdwClipPrecision, long fdwQuality, long fdwPitchAndFamily,
//		[in] string lpszFace)
var createFont = gdi32.MustFindProc("CreateFontA").Addr()
var _ = builtin("CreateFont(hdc, x, y, cx, cy, hdcSrc, x1, y1, rop)",
	func(_ *Thread, a []Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.Syscall14(createFont,
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			intArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			intArg(a[8]),
			intArg(a[9]),
			intArg(a[10]),
			intArg(a[11]),
			intArg(a[12]),
			uintptr(stringArg(a[13])))
		return intRet(rtn)
	})

// dll bool Gdi32:ExtTextOut(pointer hdc, long X, long Y, long fuOptions,
//		RECT* lprc, [in] string lpString, long cbCount, LONG* lpDx)
var extTextOut = gdi32.MustFindProc("ExtTextOutA").Addr()
var _ = builtin("ExtTextOut(hdc, x, y, fuOptions, lprc, lpString, cbCount,"+
	" lpDx/*unused*/)",
	func(_ *Thread, a []Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		verify.That(a[7].Equal(Zero))
		rtn := goc.Syscall8(extTextOut,
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			uintptr(rectArg(a[4], r)),
			uintptr(stringArg(a[5])),
			intArg(a[6]),
			0)
		return boolRet(rtn)
	})

// dll gdiobj Gdi32:CreatePen(long fnPenStyle, long nWidth, long clrref)
var createPen = gdi32.MustFindProc("CreatePen").Addr()
var _ = builtin3("CreatePen(fnPenStyle, nWidth, clrref)",
	func(a, b, c Value) Value {
		rtn := goc.Syscall3(createPen,
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll gdiobj Gdi32:ExtCreatePen(long dwPenStyle, long dwWidth, LOGBRUSH* brush,
//		long dwStyleCount, pointer lpStyle)
var extCreatePen = gdi32.MustFindProc("ExtCreatePen").Addr()
var _ = builtin5("ExtCreatePen(dwPenStyle, dwWidth, brush, "+
	"dwStyleCount/*unused*/, lpStyle/*unused*/)",
	func(a, b, c, d, e Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nLOGBRUSH)
		*(*LOGBRUSH)(p) = LOGBRUSH{
			lbStyle: getInt32(c, "lbStyle"),
			lbColor: getInt32(c, "lbColor"),
			lbHatch: uintptr(getInt(c, "lbHatch")),
		}
		rtn := goc.Syscall5(extCreatePen,
			intArg(a),
			intArg(b),
			uintptr(p),
			0,
			0)
		return intRet(rtn)
	})

type LOGBRUSH struct {
	lbStyle int32
	lbColor int32
	lbHatch uintptr
}

const nLOGBRUSH = unsafe.Sizeof(LOGBRUSH{})

// dll long Gdi32:GetGlyphOutline(pointer hdc, long uChar, long uFormat,
//		GLYPHMETRICS*  lpgm, long cbBuffer, pointer lpvBuffer, MAT2* lpmat2)
var getGlyphOutline = gdi32.MustFindProc("GetGlyphOutlineA").Addr()
var _ = builtin7("GetGlyphOutline(hdc, uChar, uFormat, lpgm, "+
	"cbBuffer/*unused*/, lpvBuffer/*unused*/, lpmat2)",
	func(a, b, c, d, e, f, g Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p1 := heap.Alloc(nGLYPHMETRICS)
		p2 := heap.Alloc(nMAT2)
		*(*MAT2)(p2) = MAT2{
			eM11: FIXED{
				fract: getInt16(g.Get(nil, SuStr("eM11")), "fract"),
				value: getInt16(g.Get(nil, SuStr("eM11")), "value")},
			eM12: FIXED{
				fract: getInt16(g.Get(nil, SuStr("eM12")), "fract"),
				value: getInt16(g.Get(nil, SuStr("eM12")), "value")},
			eM21: FIXED{
				fract: getInt16(g.Get(nil, SuStr("eM21")), "fract"),
				value: getInt16(g.Get(nil, SuStr("eM21")), "value")},
			eM22: FIXED{
				fract: getInt16(g.Get(nil, SuStr("eM22")), "fract"),
				value: getInt16(g.Get(nil, SuStr("eM22")), "value")},
		}
		rtn := goc.Syscall7(getGlyphOutline,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(p1),
			0,
			0,
			uintptr(p2))
		gm := *(*GLYPHMETRICS)(p1)
		d.Put(nil, SuStr("gmBlackBoxX"), IntVal(int(gm.gmBlackBoxX)))
		d.Put(nil, SuStr("gmBlackBoxY"), IntVal(int(gm.gmBlackBoxY)))
		d.Put(nil, SuStr("gmptGlyphOrigin"),
			pointToOb(&gm.gmptGlyphOrigin, d.Get(nil, SuStr("gmptGlyphOrigin"))))
		d.Put(nil, SuStr("gmCellIncX"), IntVal(int(gm.gmCellIncX)))
		d.Put(nil, SuStr("gmCellIncY"), IntVal(int(gm.gmCellIncY)))
		return intRet(rtn)
	})

type GLYPHMETRICS struct {
	gmBlackBoxX     int32
	gmBlackBoxY     int32
	gmptGlyphOrigin POINT
	gmCellIncX      int16
	gmCellIncY      int16
}

const nGLYPHMETRICS = unsafe.Sizeof(GLYPHMETRICS{})

type FIXED struct {
	fract int16
	value int16
}

type MAT2 struct {
	eM11 FIXED
	eM12 FIXED
	eM21 FIXED
	eM22 FIXED
}

const nMAT2 = unsafe.Sizeof(MAT2{})

// dll long Gdi32:StartDoc(pointer hdc, DOCINFO* di)
var startDoc = gdi32.MustFindProc("StartDocA").Addr()
var _ = builtin2("StartDoc(hdc, di)",
	func(a, b Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nDOCINFO)
		*(*DOCINFO)(p) = DOCINFO{
			cbSize:       uint32(nDOCINFO),
			lpszDocName:  getStr(b, "lpszDocName"),
			lpszOutput:   getStr(b, "lpszOutput"),
			lpszDatatype: getStr(b, "lpszDatatype"),
			fwType:       getInt32(b, "fwType"),
		}
		rtn := goc.Syscall2(startDoc,
			intArg(a),
			uintptr(p))
		return intRet(rtn)
	})

type DOCINFO struct {
	cbSize       uint32
	lpszDocName  *byte
	lpszOutput   *byte
	lpszDatatype *byte
	fwType       int32
	_            [4]byte // padding
}

const nDOCINFO = unsafe.Sizeof(DOCINFO{})

// dll bool Gdi32:SetWindowExtEx(pointer hdc, long x, long y, POINT* p)
var setWindowExtEx = gdi32.MustFindProc("SetWindowExtEx").Addr()
var _ = builtin4("SetWindowExtEx(hdc, x, y, p)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		pt := heap.Alloc(nPOINT)
		rtn := goc.Syscall4(setWindowExtEx,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(pt))
		if !d.Equal(Zero) {
			upointToOb(pt, d)
		}
		return boolRet(rtn)
	})

// dll bool Gdi32:SetViewportOrgEx(pointer hdc, long x, long y, POINT* p)
var setViewportOrgEx = gdi32.MustFindProc("SetViewportOrgEx").Addr()
var _ = builtin4("SetViewportOrgEx(hdc, x, y, p)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		pt := heap.Alloc(nPOINT)
		rtn := goc.Syscall4(setViewportOrgEx,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(pt))
		if !d.Equal(Zero) {
			upointToOb(pt, d)
		}
		return boolRet(rtn)
	})

// dll bool Gdi32:SetViewportExtEx(pointer hdc, long x, long y, POINT* p)
var setViewportExtEx = gdi32.MustFindProc("SetViewportExtEx").Addr()
var _ = builtin4("SetViewportExtEx(hdc, x, y, p)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		pt := heap.Alloc(nPOINT)
		rtn := goc.Syscall4(setViewportExtEx,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(pt))
		if !d.Equal(Zero) {
			upointToOb(pt, d)
		}
		return boolRet(rtn)
	})

// dll bool Gdi32:PlayEnhMetaFile(pointer hdc, pointer hemf, RECT* rect)
var playEnhMetaFile = gdi32.MustFindProc("PlayEnhMetaFile").Addr()
var _ = builtin3("PlayEnhMetaFile(hdc, hemf, rect)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		r := heap.Alloc(nRECT)
		rtn := goc.Syscall3(playEnhMetaFile,
			intArg(a),
			intArg(b),
			uintptr(rectArg(c, r)))
		return boolRet(rtn)
	})

// dll bool Gdi32:MoveToEx(pointer hdc, long x, long y, POINT* p)
var moveToEx = gdi32.MustFindProc("MoveToEx").Addr()
var _ = builtin4("MoveToEx(hdc, x, y, p)",
	func(a, b, c, d Value) Value {
		defer heap.FreeTo(heap.CurSize())
		var pt unsafe.Pointer
		if !d.Equal(Zero) {
			pt = heap.Alloc(nPOINT)
		}
		rtn := goc.Syscall4(moveToEx,
			intArg(a),
			intArg(b),
			intArg(c),
			uintptr(pt))
		if pt != nil {
			upointToOb(pt, d)
		}
		return boolRet(rtn)
	})

// dll long Gdi32:GetObject(pointer hgdiobj, long bufsize, buffer buf)
var getObject = gdi32.MustFindProc("GetObjectA").Addr()
var _ = builtin1("GetObjectBitmap(h)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nBITMAP)
		bm := (*BITMAP)(p)
		rtn := goc.Syscall3(getObject,
			intArg(a),
			nBITMAP,
			uintptr(p))
		if rtn != nBITMAP {
			return False
		}
		ob := NewSuObject()
		ob.Put(nil, SuStr("bmType"), IntVal(int(bm.bmType)))
		ob.Put(nil, SuStr("bmWidth"), IntVal(int(bm.bmWidth)))
		ob.Put(nil, SuStr("bmHeight"), IntVal(int(bm.bmHeight)))
		ob.Put(nil, SuStr("bmWidthBytes"), IntVal(int(bm.bmWidthBytes)))
		ob.Put(nil, SuStr("bmPlanes"), IntVal(int(bm.bmPlanes)))
		ob.Put(nil, SuStr("bmBitsPixel"), IntVal(int(bm.bmBitsPixel)))
		// bmBits not used
		return ob
	})

type BITMAP struct {
	bmType       uint32
	bmWidth      uint32
	bmHeight     uint32
	bmWidthBytes uint32
	bmPlanes     int16
	bmBitsPixel  int16
	bmBits       uintptr
}

const nBITMAP = unsafe.Sizeof(BITMAP{})

// dll gdiobj Gdi32:CreateRectRgn(long x1, long y1, long x2, long y2)
var createRectRegion = gdi32.MustFindProc("CreateRectRgn").Addr()
var _ = builtin4("CreateRectRgn(x1, y1, x2, y2)",
	func(a, b, c, d Value) Value {
		rtn := goc.Syscall4(createRectRegion,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d))
		return intRet(rtn)
	})

// dll long Gdi32:GetDIBits(pointer hdc, pointer hbmp, long uStartScan,
//		long cScanLines, pointer lpvBits, BITMAPINFO* lpbi, long uUsage)
var getDIBits = gdi32.MustFindProc("GetDIBits").Addr()
var _ = builtin7("GetDIBits(hdc, hbmp, uStartScan, cScanLines, lpvBits,"+
	" lpbi, uUsage)",
	func(a, b, c, d, e, f, g Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(nBITMAPINFOHEADER)
		hdr := f.Get(nil, SuStr("bmiHeader"))
		bmih := (*BITMAPINFOHEADER)(p)
		*bmih = obToBMIH(hdr)
		rtn := goc.Syscall7(getDIBits,
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e),
			uintptr(p),
			intArg(g))
		hdr.Put(nil, SuStr("biSizeImage"), IntVal(int(bmih.biSizeImage)))
		hdr.Put(nil, SuStr("biXPelsPerMeter"), IntVal(int(bmih.biXPelsPerMeter)))
		hdr.Put(nil, SuStr("biYPelsPerMeter"), IntVal(int(bmih.biYPelsPerMeter)))
		hdr.Put(nil, SuStr("biClrUsed"), IntVal(int(bmih.biClrUsed)))
		hdr.Put(nil, SuStr("biClrImportant"), IntVal(int(bmih.biClrImportant)))
		return intRet(rtn)
	})

type BITMAPINFOHEADER struct {
	biSize          int32
	biWidth         int32
	biHeight        int32
	biPlanes        int16
	biBitCount      int16
	biCompression   int32
	biSizeImage     int32
	biXPelsPerMeter int32
	biYPelsPerMeter int32
	biClrUsed       int32
	biClrImportant  int32
}

const nBITMAPINFOHEADER = unsafe.Sizeof(BITMAPINFOHEADER{})
