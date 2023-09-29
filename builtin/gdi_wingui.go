// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

var gdi32 = MustLoadDLL("gdi32.dll")

const LF_FACESIZE = 32

type stLogFont struct {
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

const nLogFont = unsafe.Sizeof(stLogFont{})

type stTextMetric struct {
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

const nTextMetric = unsafe.Sizeof(stTextMetric{})

// dll Gdi32:CreateFontIndirect(LOGFONT* lf) gdiobj
var createFontIndirect = gdi32.MustFindProc("CreateFontIndirectA").Addr()
var _ = builtin(CreateFontIndirect, "(logfont)")

func CreateFontIndirect(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nLogFont)
	f := (*stLogFont)(p)
	*f = stLogFont{
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
}

// dll pointer Gdi32:AddFontMemResourceEx([in] string pFileView, long cjSize,
// pointer pvResrved, LONG* pNumFonts)
var addFontMemResourceEx = gdi32.MustFindProc("AddFontMemResourceEx").Addr()
var _ = builtin(AddFontMemResourceEx, "(pFileView, cjSize, pvResrved, pNumFonts)")

func AddFontMemResourceEx(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	pFileView := ToStr(a)
	cjSize := ToInt(b)
	pNumFonts := heap.Alloc(int32Size)
	rtn := goc.Syscall4(addFontMemResourceEx,
		uintptr(heap.Copy(pFileView, cjSize)),
		uintptr(cjSize),
		0,
		uintptr(pNumFonts))
	d.Put(nil, SuStr("x"), IntVal(int(*(*int32)(pNumFonts))))
	return intRet(rtn)
}

// dll bool Gdi32:GetTextMetrics(pointer hdc, TEXTMETRIC* tm)
var getTextMetrics = gdi32.MustFindProc("GetTextMetricsA").Addr()
var _ = builtin(GetTextMetrics, "(hdc, tm)")

func GetTextMetrics(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nTextMetric)
	rtn := goc.Syscall2(getTextMetrics,
		intArg(a),
		uintptr(p))
	tm := (*stTextMetric)(p)
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
}

var getStockObject = gdi32.MustFindProc("GetStockObject").Addr()
var _ = builtin(GetStockObject, "(i)")

func GetStockObject(a Value) Value {
	rtn := goc.Syscall1(getStockObject,
		intArg(a))
	return intRet(rtn)
}

var getDeviceCaps = gdi32.MustFindProc("GetDeviceCaps").Addr()
var _ = builtin(GetDeviceCaps, "(hdc, nIndex)")

func GetDeviceCaps(a, b Value) Value {
	rtn := goc.Syscall2(getDeviceCaps,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll Gdi32:BitBlt(pointer hdcDest,
// long nXDest, long nYDest, long nWidth, long nHeight,
// pointer hdcSrc, long nXSrc, long nYSrc, long dwRop) bool
var bitBlt = gdi32.MustFindProc("BitBlt").Addr()
var _ = builtin(BitBlt, "(hdc, x, y, cx, cy, hdcSrc, x1, y1, rop)")

func BitBlt(_ *Thread, a []Value) Value {
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
}

// dll Gdi32:CreateCompatibleBitmap(pointer hdc, long nWidth, long nHeight) gdiobj
var createCompatibleBitmap = gdi32.MustFindProc("CreateCompatibleBitmap").Addr()
var _ = builtin(CreateCompatibleBitmap, "(hdc, cx, cy)")

func CreateCompatibleBitmap(a, b, c Value) Value {
	rtn := goc.Syscall3(createCompatibleBitmap,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll Gdi32:CreateCompatibleDC(pointer hdc) pointer
var createCompatibleDC = gdi32.MustFindProc("CreateCompatibleDC").Addr()
var _ = builtin(CreateCompatibleDC, "(hdc)")

func CreateCompatibleDC(a Value) Value {
	rtn := goc.Syscall1(createCompatibleDC,
		intArg(a))
	return intRet(rtn)
}

// dll Gdi32:CreateSolidBrush(long rgb) gdiobj
var createSolidBrush = gdi32.MustFindProc("CreateSolidBrush").Addr()
var _ = builtin(CreateSolidBrush, "(i)")

func CreateSolidBrush(a Value) Value {
	rtn := goc.Syscall1(createSolidBrush,
		intArg(a))
	return intRet(rtn)
}

// dll pointer Gdi32:SelectObject(pointer hdc, pointer obj)
var selectObject = gdi32.MustFindProc("SelectObject").Addr()
var _ = builtin(SelectObject, "(hdc, obj)")

func SelectObject(a, b Value) Value {
	rtn := goc.Syscall2(selectObject,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll bool Gdi32:GetTextExtentPoint32(pointer hdc, [in] string text, long len, POINT* p)
var getTextExtentPoint32 = gdi32.MustFindProc("GetTextExtentPoint32A").Addr()
var _ = builtin(GetTextExtentPoint32, "(hdc, text, len, p)")

func GetTextExtentPoint32(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	pt := heap.Alloc(nPoint)
	rtn := goc.Syscall4(getTextExtentPoint32,
		intArg(a),
		uintptr(stringArg(b)),
		uintptr(len(ToStr(b))),
		uintptr(pt))
	upointToOb(pt, d)
	return boolRet(rtn)
}

var getTextFace = gdi32.MustFindProc("GetTextFaceA").Addr()
var _ = builtin(GetTextFace, "(hdc)")

func GetTextFace(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	const bufsize = 512
	p := heap.Alloc(bufsize)
	n := goc.Syscall3(getTextFace,
		intArg(a),
		bufsize,
		uintptr(p))
	return SuStr(heap.GetStrN(p, int(n)))
}

// dll long Gdi32:SetBkMode(pointer hdc, long mode)
var setBkMode = gdi32.MustFindProc("SetBkMode").Addr()
var _ = builtin(SetBkMode, "(hdc, color)")

func SetBkMode(a, b Value) Value {
	rtn := goc.Syscall2(setBkMode,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll long Gdi32:SetBkColor(pointer hdc, long color)
var setBkColor = gdi32.MustFindProc("SetBkColor").Addr()
var _ = builtin(SetBkColor, "(hdc, color)")

func SetBkColor(a, b Value) Value {
	rtn := goc.Syscall2(setBkColor,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll Gdi32:DeleteDC(pointer hdc) bool

var deleteDC = gdi32.MustFindProc("DeleteDC").Addr()
var _ = builtin(DeleteDC, "(hdc)")

func DeleteDC(a Value) Value {
	rtn := goc.Syscall1(deleteDC,
		intArg(a))
	return boolRet(rtn)
}

// dll Gdi32:DeleteObject(pointer hgdiobj) bool
var deleteObject = gdi32.MustFindProc("DeleteObject").Addr()
var _ = builtin(DeleteObject, "(hgdiobj)")

func DeleteObject(a Value) Value {
	rtn := goc.Syscall1(deleteObject,
		intArg(a))
	return boolRet(rtn)
}

// dll Gdi32:GetClipBox(pointer hdc, RECT* rect) long
var getClipBox = gdi32.MustFindProc("GetClipBox").Addr()
var _ = builtin(GetClipBox, "(hdc, rect)")

func GetClipBox(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall3(getClipBox,
		intArg(a),
		uintptr(rectArg(b, r)),
		0)
	urectToOb(r, b)
	return intRet(rtn)
}

// dll Gdi32:GetPixel(pointer hdc, long nXPos, long nYPos) long
var getPixel = gdi32.MustFindProc("GetPixel").Addr()
var _ = builtin(GetPixel, "(hdc, nXPos, nYPos)")

func GetPixel(a, b, c Value) Value {
	rtn := goc.Syscall3(getPixel,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll Gdi32:PtVisible(pointer hdc, long nXPos, long nYPos) bool
var ptVisible = gdi32.MustFindProc("PtVisible").Addr()
var _ = builtin(PtVisible, "(hdc, nXPos, nYPos)")

func PtVisible(a, b, c Value) Value {
	rtn := goc.Syscall3(ptVisible,
		intArg(a),
		intArg(b),
		intArg(c))
	if rtn == 0xffffffff { // 0xffffffff: 32 bit -1
		return IntVal(-1)
	}
	return boolRet(rtn)
}

// dll Gdi32:Rectangle(pointer hdc, long left, long top, long right, long bottom) bool
var rectangle = gdi32.MustFindProc("Rectangle").Addr()
var _ = builtin(Rectangle, "(hdc, left, top, right, bottom)")

func Rectangle(a, b, c, d, e Value) Value {
	rtn := goc.Syscall5(rectangle,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e))
	return boolRet(rtn)
}

// dll Gdi32:SetStretchBltMode(pointer hdc, long iStretchMode) long
var setStretchBltMode = gdi32.MustFindProc("SetStretchBltMode").Addr()
var _ = builtin(SetStretchBltMode, "(hdc, iStretchMode)")

func SetStretchBltMode(a, b Value) Value {
	rtn := goc.Syscall2(setStretchBltMode,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll Gdi32:SetTextColor(pointer hdc, long color) long
var setTextColor = gdi32.MustFindProc("SetTextColor").Addr()
var _ = builtin(SetTextColor, "(hdc, color)")

func SetTextColor(a, b Value) Value {
	rtn := goc.Syscall2(setTextColor,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll Gdi32:StretchBlt(pointer hdcDest, long nXOriginDest, long nYOriginDest,
// long nWidthDest, long nHeightDest, pointer hdcSrc, long nXOriginSrc,
// long nYOriginSrc, long nWidthSrc, long nHeightSrc, long dwRop) bool
var stretchBlt = gdi32.MustFindProc("StretchBlt").Addr()
var _ = builtin(StretchBlt, "(hdcDest, nXOriginDest, nYOriginDest, nWidthDest, "+
	"nHeightDest, hdcSrc, nXOriginSrc, nYOriginSrc, nWidthSrc, nHeightSrc, dwRop)")

func StretchBlt(_ *Thread, a []Value) Value {
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
}

// dll pointer Gdi32:CloseEnhMetaFile(pointer dc)
var closeEnhMetaFile = gdi32.MustFindProc("CloseEnhMetaFile").Addr()
var _ = builtin(CloseEnhMetaFile, "(dc)")

func CloseEnhMetaFile(a Value) Value {
	rtn := goc.Syscall1(closeEnhMetaFile,
		intArg(a))
	return intRet(rtn)
}

// dll pointer Gdi32:DeleteEnhMetaFile(pointer emf)
var deleteEnhMetaFile = gdi32.MustFindProc("DeleteEnhMetaFile").Addr()
var _ = builtin(DeleteEnhMetaFile, "(dc)")

func DeleteEnhMetaFile(a Value) Value {
	rtn := goc.Syscall1(deleteEnhMetaFile,
		intArg(a))
	return intRet(rtn)
}

// dll pointer Gdi32:CreateDC(
// [in] string driver,
// [in] string device,
// [in] string output,
// pointer devmode)
var createDC = gdi32.MustFindProc("CreateDCA").Addr()
var _ = builtin(CreateDC, "(driver, device, output, devmode)")

func CreateDC(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall4(createDC,
		uintptr(stringArg(a)),
		uintptr(stringArg(b)),
		uintptr(stringArg(c)),
		intArg(d))
	return intRet(rtn)
}

// dll bool Gdi32:Ellipse(pointer hdc, long left, long top, long right, long bottom)
var ellipse = gdi32.MustFindProc("Ellipse").Addr()
var _ = builtin(Ellipse, "(hdc, left, top, right, bottom)")

func Ellipse(a, b, c, d, e Value) Value {
	rtn := goc.Syscall5(ellipse,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e))
	return boolRet(rtn)
}

// dll long Gdi32:EndDoc(pointer hdc)
var endDoc = gdi32.MustFindProc("EndDoc").Addr()
var _ = builtin(EndDoc, "(hdc)")

func EndDoc(a Value) Value {
	rtn := goc.Syscall1(endDoc,
		intArg(a))
	return intRet(rtn)
}

// dll long Gdi32:EndPage(pointer hdc)
var endPage = gdi32.MustFindProc("EndPage").Addr()
var _ = builtin(EndPage, "(hdc)")

func EndPage(a Value) Value {
	rtn := goc.Syscall1(endPage,
		intArg(a))
	return intRet(rtn)
}

// dll long Gdi32:ExcludeClipRect(pointer hdc, long l, long t, long r, long b)
var excludeClipRect = gdi32.MustFindProc("ExcludeClipRect").Addr()
var _ = builtin(ExcludeClipRect, "(hdc, l, t, r, b)")

func ExcludeClipRect(a, b, c, d, e Value) Value {
	rtn := goc.Syscall5(excludeClipRect,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e))
	return intRet(rtn)
}

// dll long Gdi32:GetClipRgn(pointer hdc, pointer hrgn)
var getClipRgn = gdi32.MustFindProc("GetClipRgn").Addr()
var _ = builtin(GetClipRgn, "(hdc, hrgn)")

func GetClipRgn(a, b Value) Value {
	rtn := goc.Syscall2(getClipRgn,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll pointer Gdi32:GetCurrentObject(pointer hdc, long uObjectType)
var getCurrentObject = gdi32.MustFindProc("GetCurrentObject").Addr()
var _ = builtin(GetCurrentObject, "(hdc, uObjectType)")

func GetCurrentObject(a, b Value) Value {
	rtn := goc.Syscall2(getCurrentObject,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll pointer Gdi32:GetEnhMetaFile(string filename)
var getEnhMetaFile = gdi32.MustFindProc("GetEnhMetaFileA").Addr()
var _ = builtin(GetEnhMetaFile, "(filename)")

func GetEnhMetaFile(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall1(getEnhMetaFile,
		uintptr(stringArg(a)))
	return intRet(rtn)
}

// dll bool Gdi32:LineTo(pointer hdc, long x, long y)
var lineTo = gdi32.MustFindProc("LineTo").Addr()
var _ = builtin(LineTo, "(hdc, x, y)")

func LineTo(a, b, c Value) Value {
	rtn := goc.Syscall3(lineTo,
		intArg(a),
		intArg(b),
		intArg(c))
	return boolRet(rtn)
}

// dll bool Gdi32:PatBlt(
// pointer hdc,  // Destination device context
// long	nXLeft,  // x-coordinate of upper left corner of destination rectangle
// long	nYLeft,  // y-coordinate of upper left corner of destination rectangle
// long	nWidth,  // width of destination rectangle
// long	nHeight, // height of destination rectangle
// long	dwRop)   // Raster operation
var patBlt = gdi32.MustFindProc("PatBlt").Addr()
var _ = builtin(PatBlt, "(hdc, nXLeft, nYLeft, nWidth, nHeight, dwRop)")

func PatBlt(a, b, c, d, e, f Value) Value {
	rtn := goc.Syscall6(patBlt,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e),
		intArg(f))
	return boolRet(rtn)
}

// dll bool Gdi32:PolygonApi(
// pointer hdc,		 // handle to device context
// [in] string lppt, // array of points
// long cCount)		 // count of points
var polygon = gdi32.MustFindProc("Polygon").Addr()
var _ = builtin(Polygon, "(hdc, points, npoints = false)")

func Polygon(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	ob := ToContainer(b)
	var n int
	if c == False {
		n = ob.ListSize()
	} else {
		n = ToInt(c)
	}
	p := heap.Alloc(uintptr(n) * nPoint)
	for i := 0; i < min(n, ob.ListSize()); i++ {
		*(*stPoint)(unsafe.Pointer(uintptr(p) + uintptr(i)*nPoint)) =
			obToPoint(ob.ListGet(i))
	}
	rtn := goc.Syscall3(polygon,
		intArg(a),
		uintptr(p),
		intArg(c))
	return boolRet(rtn)
}

// dll bool Gdi32:RestoreDC(
// pointer	hdc,  // handle to DC
// long	nSavedDC) // restore state returned by SaveDC
var restoreDC = gdi32.MustFindProc("RestoreDC").Addr()
var _ = builtin(RestoreDC, "(hdc, nSavedDC)")

func RestoreDC(a, b Value) Value {
	rtn := goc.Syscall2(restoreDC,
		intArg(a),
		intArg(b))
	return boolRet(rtn)
}

// dll bool Gdi32:RoundRect(pointer hdc,
// long left, long top, long right, long bottom,
// long ellipse_width, long ellipse_height)
var roundRect = gdi32.MustFindProc("RoundRect").Addr()
var _ = builtin(RoundRect,
	"(hdc, left, top, right, bottom, ellipse_width, ellipse_height)")

func RoundRect(a, b, c, d, e, f, g Value) Value {
	rtn := goc.Syscall7(roundRect,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d),
		intArg(e),
		intArg(f),
		intArg(g))
	return boolRet(rtn)
}

// dll long Gdi32:SaveDC(pointer hdc)
var saveDC = gdi32.MustFindProc("SaveDC").Addr()
var _ = builtin(SaveDC, "(hdc)")

func SaveDC(a Value) Value {
	rtn := goc.Syscall1(saveDC,
		intArg(a))
	return intRet(rtn)
}

// dll long Gdi32:SelectClipRgn(pointer hdc, pointer hrgn)
var selectClipRgn = gdi32.MustFindProc("SelectClipRgn").Addr()
var _ = builtin(SelectClipRgn, "(hdc, hrgn)")

func SelectClipRgn(a, b Value) Value {
	rtn := goc.Syscall2(selectClipRgn,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll pointer Gdi32:SetEnhMetaFileBits(long cbBuffer, [in] string lpData)

var setEnhMetaFileBits = gdi32.MustFindProc("SetEnhMetaFileBits").Addr()
var _ = builtin(SetEnhMetaFileBits, "(cbBuffer, lpData)")

func SetEnhMetaFileBits(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall2(setEnhMetaFileBits,
		intArg(a),
		uintptr(stringArg(b)))
	return intRet(rtn)
}

// dll long Gdi32:SetMapMode(pointer hdc, long mode)

var setMapMode = gdi32.MustFindProc("SetMapMode").Addr()
var _ = builtin(SetMapMode, "(hdc, mode)")

func SetMapMode(a, b Value) Value {
	rtn := goc.Syscall2(setMapMode,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll long Gdi32:SetROP2(pointer hdc, long fnDrawMode)
var setROP2 = gdi32.MustFindProc("SetROP2").Addr()
var _ = builtin(SetROP2, "(hdc, fnDrawMode)")

func SetROP2(a, b Value) Value {
	rtn := goc.Syscall2(setROP2,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll long Gdi32:SetTextAlign(pointer hdc, long mode)
var setTextAlign = gdi32.MustFindProc("SetTextAlign").Addr()
var _ = builtin(SetTextAlign, "(hdc, mode)")

func SetTextAlign(a, b Value) Value {
	rtn := goc.Syscall2(setTextAlign,
		intArg(a),
		intArg(b))
	return intRet(rtn)
}

// dll long Gdi32:StartPage(pointer hdc)
var startPage = gdi32.MustFindProc("StartPage").Addr()
var _ = builtin(StartPage, "(hdc)")

func StartPage(a Value) Value {
	rtn := goc.Syscall1(startPage,
		intArg(a))
	return intRet(rtn)
}

// dll bool Gdi32:TextOut(pointer hdc, long x, long y, [in] string text, long n)
var textOut = gdi32.MustFindProc("TextOutA").Addr()
var _ = builtin(TextOut, "(hdc, x, y, text, n)")

func TextOut(a, b, c, d, e Value) Value {
	defer heap.FreeTo(heap.CurSize())
	rtn := goc.Syscall5(textOut,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(stringArg(d)),
		intArg(e))
	return boolRet(rtn)
}

// dll bool Gdi32:Arc(pointer hdc, long nLeftRect, long nTopRect,
// long nRightRect, long nBottomRect, long nXStartArc, long nYStartArc,
// long nXEndArc, long nYEndArc)
var arc = gdi32.MustFindProc("Arc").Addr()
var _ = builtin(Arc, "(hdc, nLeftRect, nTopRect, nRightRect, nBottomRect, "+
	"nXStartArc, nYStartArc, nXEndArc, nYEndArc)")

func Arc(_ *Thread, a []Value) Value {
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
}

// dll pointer Gdi32:CreateEnhMetaFile(pointer hdcRef, [in] string filename,
// RECT* rect, [in] string desc)
var createEnhMetaFile = gdi32.MustFindProc("CreateEnhMetaFileA").Addr()
var _ = builtin(CreateEnhMetaFile, "(hdcRef, filename, rect, desc)")

func CreateEnhMetaFile(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall4(createEnhMetaFile,
		intArg(a),
		uintptr(stringArg(b)),
		uintptr(rectArg(c, r)),
		uintptr(stringArg(d)))
	return intRet(rtn)
}

// dll gdiobj Gdi32:CreateFont(long nHeight, long nWidth, long nEscapement,
// long nOrientation, long fnWeight, long fdwItalic, long fdwUnderline,
// long fdwStrikeOut, long fdwCharSet, long fdwOutputPrecision,
// long fdwClipPrecision, long fdwQuality, long fdwPitchAndFamily,
// [in] string lpszFace)
var createFont = gdi32.MustFindProc("CreateFontA").Addr()
var _ = builtin(CreateFont, "(hdc, x, y, cx, cy, hdcSrc, x1, y1, rop)")

func CreateFont(_ *Thread, a []Value) Value {
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
}

// dll bool Gdi32:ExtTextOut(pointer hdc, long X, long Y, long fuOptions,
// RECT* lprc, [in] string lpString, long cbCount, LONG* lpDx)
var extTextOut = gdi32.MustFindProc("ExtTextOutA").Addr()
var _ = builtin(ExtTextOut,
	"(hdc, x, y, fuOptions, lprc, lpString, cbCount, lpDx/*unused*/)")

func ExtTextOut(_ *Thread, a []Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	assert.That(a[7].Equal(Zero))
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
}

// dll gdiobj Gdi32:CreatePen(long fnPenStyle, long nWidth, long clrref)
var createPen = gdi32.MustFindProc("CreatePen").Addr()
var _ = builtin(CreatePen, "(fnPenStyle, nWidth, clrref)")

func CreatePen(a, b, c Value) Value {
	rtn := goc.Syscall3(createPen,
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
}

// dll gdiobj Gdi32:ExtCreatePen(long dwPenStyle, long dwWidth, LOGBRUSH* brush,
// long dwStyleCount, pointer lpStyle)
var extCreatePen = gdi32.MustFindProc("ExtCreatePen").Addr()
var _ = builtin(ExtCreatePen, "(dwPenStyle, dwWidth, brush, "+
	"dwStyleCount/*unused*/, lpStyle/*unused*/)")

func ExtCreatePen(a, b, c, d, e Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nLogBrush)
	*(*stLogBrush)(p) = stLogBrush{
		lbStyle: getInt32(c, "lbStyle"),
		lbColor: getInt32(c, "lbColor"),
		lbHatch: getUintptr(c, "lbHatch"),
	}
	rtn := goc.Syscall5(extCreatePen,
		intArg(a),
		intArg(b),
		uintptr(p),
		0,
		0)
	return intRet(rtn)
}

type stLogBrush struct {
	lbStyle int32
	lbColor int32
	lbHatch uintptr
}

const nLogBrush = unsafe.Sizeof(stLogBrush{})

// dll long Gdi32:GetGlyphOutline(pointer hdc, long uChar, long uFormat,
// GLYPHMETRICS*  lpgm, long cbBuffer, pointer lpvBuffer, MAT2* lpmat2)
var getGlyphOutline = gdi32.MustFindProc("GetGlyphOutlineA").Addr()
var _ = builtin(GetGlyphOutline, "(hdc, uChar, uFormat, lpgm, "+
	"cbBuffer/*unused*/, lpvBuffer/*unused*/, lpmat2)")

func GetGlyphOutline(a, b, c, d, e, f, g Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p1 := heap.Alloc(nGlyphMetrics)
	p2 := heap.Alloc(nMat2)
	*(*stMat2)(p2) = stMat2{
		eM11: stFixed{
			fract: getInt16(g.Get(nil, SuStr("eM11")), "fract"),
			value: getInt16(g.Get(nil, SuStr("eM11")), "value")},
		eM12: stFixed{
			fract: getInt16(g.Get(nil, SuStr("eM12")), "fract"),
			value: getInt16(g.Get(nil, SuStr("eM12")), "value")},
		eM21: stFixed{
			fract: getInt16(g.Get(nil, SuStr("eM21")), "fract"),
			value: getInt16(g.Get(nil, SuStr("eM21")), "value")},
		eM22: stFixed{
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
	gm := *(*stGlyphMetrics)(p1)
	d.Put(nil, SuStr("gmBlackBoxX"), IntVal(int(gm.gmBlackBoxX)))
	d.Put(nil, SuStr("gmBlackBoxY"), IntVal(int(gm.gmBlackBoxY)))
	d.Put(nil, SuStr("gmptGlyphOrigin"),
		pointToOb(&gm.gmptGlyphOrigin, d.Get(nil, SuStr("gmptGlyphOrigin"))))
	d.Put(nil, SuStr("gmCellIncX"), IntVal(int(gm.gmCellIncX)))
	d.Put(nil, SuStr("gmCellIncY"), IntVal(int(gm.gmCellIncY)))
	return intRet(rtn)
}

type stGlyphMetrics struct {
	gmBlackBoxX     int32
	gmBlackBoxY     int32
	gmptGlyphOrigin stPoint
	gmCellIncX      int16
	gmCellIncY      int16
}

const nGlyphMetrics = unsafe.Sizeof(stGlyphMetrics{})

type stFixed struct {
	fract int16
	value int16
}

type stMat2 struct {
	eM11 stFixed
	eM12 stFixed
	eM21 stFixed
	eM22 stFixed
}

const nMat2 = unsafe.Sizeof(stMat2{})

// dll long Gdi32:StartDoc(pointer hdc, DOCINFO* di)
var startDoc = gdi32.MustFindProc("StartDocA").Addr()
var _ = builtin(StartDoc, "(hdc, di)")

func StartDoc(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nDocInfo)
	*(*stDocInfo)(p) = stDocInfo{
		cbSize:       uint32(nDocInfo),
		lpszDocName:  getStr(b, "lpszDocName"),
		lpszOutput:   getStr(b, "lpszOutput"),
		lpszDatatype: getStr(b, "lpszDatatype"),
		fwType:       getInt32(b, "fwType"),
	}
	rtn := goc.Syscall2(startDoc,
		intArg(a),
		uintptr(p))
	return intRet(rtn)
}

type stDocInfo struct {
	cbSize       uint32
	lpszDocName  *byte
	lpszOutput   *byte
	lpszDatatype *byte
	fwType       int32
	_            [4]byte // padding
}

const nDocInfo = unsafe.Sizeof(stDocInfo{})

// dll bool Gdi32:SetWindowExtEx(pointer hdc, long x, long y, POINT* p)
var setWindowExtEx = gdi32.MustFindProc("SetWindowExtEx").Addr()
var _ = builtin(SetWindowExtEx, "(hdc, x, y, p)")

func SetWindowExtEx(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	pt := heap.Alloc(nPoint)
	rtn := goc.Syscall4(setWindowExtEx,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(pt))
	if !d.Equal(Zero) {
		upointToOb(pt, d)
	}
	return boolRet(rtn)
}

// dll bool Gdi32:SetViewportOrgEx(pointer hdc, long x, long y, POINT* p)
var setViewportOrgEx = gdi32.MustFindProc("SetViewportOrgEx").Addr()
var _ = builtin(SetViewportOrgEx, "(hdc, x, y, p)")

func SetViewportOrgEx(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	pt := heap.Alloc(nPoint)
	rtn := goc.Syscall4(setViewportOrgEx,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(pt))
	if !d.Equal(Zero) {
		upointToOb(pt, d)
	}
	return boolRet(rtn)
}

// dll bool Gdi32:SetViewportExtEx(pointer hdc, long x, long y, POINT* p)
var setViewportExtEx = gdi32.MustFindProc("SetViewportExtEx").Addr()
var _ = builtin(SetViewportExtEx, "(hdc, x, y, p)")

func SetViewportExtEx(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	pt := heap.Alloc(nPoint)
	rtn := goc.Syscall4(setViewportExtEx,
		intArg(a),
		intArg(b),
		intArg(c),
		uintptr(pt))
	if !d.Equal(Zero) {
		upointToOb(pt, d)
	}
	return boolRet(rtn)
}

// dll bool Gdi32:PlayEnhMetaFile(pointer hdc, pointer hemf, RECT* rect)
var playEnhMetaFile = gdi32.MustFindProc("PlayEnhMetaFile").Addr()
var _ = builtin(PlayEnhMetaFile, "(hdc, hemf, rect)")

func PlayEnhMetaFile(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	r := heap.Alloc(nRect)
	rtn := goc.Syscall3(playEnhMetaFile,
		intArg(a),
		intArg(b),
		uintptr(rectArg(c, r)))
	return boolRet(rtn)
}

// dll bool Gdi32:MoveToEx(pointer hdc, long x, long y, POINT* p)
var moveToEx = gdi32.MustFindProc("MoveToEx").Addr()
var _ = builtin(MoveToEx, "(hdc, x, y, p)")

func MoveToEx(a, b, c, d Value) Value {
	defer heap.FreeTo(heap.CurSize())
	var pt unsafe.Pointer
	if !d.Equal(Zero) {
		pt = heap.Alloc(nPoint)
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
}

// dll long Gdi32:GetObject(pointer hgdiobj, long bufsize, buffer buf)
var getObject = gdi32.MustFindProc("GetObjectA").Addr()

var _ = builtin(GetObjectBrush, "(h)")

func GetObjectBrush(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nLogBrush)
	lb := (*stLogBrush)(p)
	rtn := goc.Syscall3(getObject,
		intArg(a),
		nLogBrush,
		uintptr(p))
	if rtn != nLogBrush {
		return False
	}
	ob := &SuObject{}
	ob.Put(nil, SuStr("lbStyle"), IntVal(int(lb.lbStyle)))
	ob.Put(nil, SuStr("lbColor"), IntVal(int(lb.lbColor)))
	ob.Put(nil, SuStr("lbHatch"), IntVal(int(lb.lbHatch)))
	return ob
}

var _ = builtin(GetObjectBitmap, "(h)")

func GetObjectBitmap(a Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nBitMap)
	bm := (*stBitMap)(p)
	rtn := goc.Syscall3(getObject,
		intArg(a),
		nBitMap,
		uintptr(p))
	if rtn != nBitMap {
		return False
	}
	ob := &SuObject{}
	ob.Put(nil, SuStr("bmType"), IntVal(int(bm.bmType)))
	ob.Put(nil, SuStr("bmWidth"), IntVal(int(bm.bmWidth)))
	ob.Put(nil, SuStr("bmHeight"), IntVal(int(bm.bmHeight)))
	ob.Put(nil, SuStr("bmWidthBytes"), IntVal(int(bm.bmWidthBytes)))
	ob.Put(nil, SuStr("bmPlanes"), IntVal(int(bm.bmPlanes)))
	ob.Put(nil, SuStr("bmBitsPixel"), IntVal(int(bm.bmBitsPixel)))
	// bmBits not used
	return ob
}

type stBitMap struct {
	bmType       uint32
	bmWidth      uint32
	bmHeight     uint32
	bmWidthBytes uint32
	bmPlanes     int16
	bmBitsPixel  int16
	bmBits       uintptr
}

const nBitMap = unsafe.Sizeof(stBitMap{})

// dll gdiobj Gdi32:CreateRectRgn(long x1, long y1, long x2, long y2)
var createRectRegion = gdi32.MustFindProc("CreateRectRgn").Addr()
var _ = builtin(CreateRectRgn, "(x1, y1, x2, y2)")

func CreateRectRgn(a, b, c, d Value) Value {
	rtn := goc.Syscall4(createRectRegion,
		intArg(a),
		intArg(b),
		intArg(c),
		intArg(d))
	return intRet(rtn)
}

// dll long Gdi32:GetDIBits(pointer hdc, pointer hbmp, long uStartScan,
// long cScanLines, pointer lpvBits, BITMAPINFO* lpbi, long uUsage)
var getDIBits = gdi32.MustFindProc("GetDIBits").Addr()
var _ = builtin(GetDIBits,
	"(hdc, hbmp, uStartScan, cScanLines, lpvBits, lpbi, uUsage)")

func GetDIBits(a, b, c, d, e, f, g Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(nBitMapInfoHeader)
	hdr := f.Get(nil, SuStr("bmiHeader"))
	bmih := (*stBitMapInfoHeader)(p)
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
}

type stBitMapInfoHeader struct {
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

const nBitMapInfoHeader = unsafe.Sizeof(stBitMapInfoHeader{})
