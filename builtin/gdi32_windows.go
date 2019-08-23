package builtin

import (
	"unsafe"

	. "github.com/apmckinlay/gsuneido/runtime"
	"golang.org/x/sys/windows"
)

var gdi32 = windows.NewLazyDLL("gdi32.dll")

const LF_FACESIZE = 32

type LOGFONTA struct {
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

type TEXTMETRIC struct {
	Height           int32
	Ascent           int32
	Descent          int32
	InternalLeading  int32
	ExternalLeading  int32
	AveCharWidth     int32
	MaxCharWidth     int32
	Weight           int32
	Italic           byte
	Underlined       byte
	StruckOut        byte
	FirstChar        byte
	LastChar         byte
	DefaultChar      byte
	BreakChar        byte
	PitchAndFamily   byte
	CharSet          byte
	Overhang         int32
	DigitizedAspectX int32
	DigitizedAspectY int32
}

// dll Gdi32:CreateFontIndirect(LOGFONT* lf) gdiobj
// e.g. CreateFontIndirect(#(lfFaceName: "Segoe UI", lfHeight: -27))
var createFontIndirect = gdi32.NewProc("CreateFontIndirectA")
var _ = builtin1("CreateFontIndirect(logfont)", func(a Value) Value {
	f := LOGFONTA{
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
	copy(f.lfFaceName[:], ToStr(a.Get(nil, SuStr("lfFaceName"))))
	rtn, _, _ := createFontIndirect.Call(uintptr(unsafe.Pointer(&f)))
	return intRet(rtn)
})

// dll bool Gdi32:GetTextMetrics(pointer hdc, TEXTMETRIC* tm)
var getTextMetrics = gdi32.NewProc("GetTextMetricsA")
var _ = builtin2("GetTextMetrics(hdc, tm)",
	func(a, b Value) Value {
		var tm TEXTMETRIC
		rtn, _, _ := getTextMetrics.Call(
			intArg(a),
			uintptr(unsafe.Pointer(&tm)))
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

var getStockObject = gdi32.NewProc("GetStockObject")
var _ = builtin1("GetStockObject(i)", func(a Value) Value {
	rtn, _, _ := getStockObject.Call(intArg(a))
	return intRet(rtn)
})

var getDeviceCaps = gdi32.NewProc("GetDeviceCaps")
var _ = builtin2("GetDeviceCaps(hdc, nIndex)", func(a, b Value) Value {
	rtn, _, _ := getDeviceCaps.Call(
		intArg(a),
		intArg(b))
	return intRet(rtn)
})

// dll GDI32:BitBlt(pointer hdcDest, long nXDest, long nYDest, long nWidth, long nHeight, pointer hdcSrc, long nXSrc, long nYSrc, long dwRop) bool
var bitBlt = gdi32.NewProc("BitBlt")
var _ = builtin("BitBlt(hdc, x, y, cx, cy, hdcSrc, x1, y1, rop)",
	func(_ *Thread, a []Value) Value {
		rtn, _, _ := bitBlt.Call(
			intArg(a[0]),
			intArg(a[1]),
			intArg(a[2]),
			intArg(a[3]),
			intArg(a[4]),
			intArg(a[5]),
			intArg(a[6]),
			intArg(a[7]),
			intArg(a[8]))
		return intRet(rtn)
	})

// dll Gdi32:CreateCompatibleBitmap(pointer hdc, long nWidth, long nHeight) gdiobj
var createCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
var _ = builtin3("CreateCompatibleBitmap(hdc, cx, cy)", func(a, b, c Value) Value {
	rtn, _, _ := createCompatibleBitmap.Call(
		intArg(a),
		intArg(b),
		intArg(c))
	return intRet(rtn)
})

// dll Gdi32:CreateCompatibleDC(pointer hdc) pointer
var createCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
var _ = builtin1("CreateCompatibleDC(hdc)", func(a Value) Value {
	rtn, _, _ := createCompatibleDC.Call(intArg(a))
	return intRet(rtn)
})

// dll Gdi32:CreateSolidBrush(long rgb) gdiobj
var createSolidBrush = gdi32.NewProc("CreateSolidBrush")
var _ = builtin1("CreateSolidBrush(i)", func(a Value) Value {
	rtn, _, _ := createSolidBrush.Call(intArg(a))
	return intRet(rtn)
})

// dll pointer Gdi32:SelectObject(pointer hdc, pointer obj)
var selectObject = gdi32.NewProc("SelectObject")
var _ = builtin2("SelectObject(hdc, obj)",
	func(a, b Value) Value {
		rtn, _, _ := selectObject.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll bool Gdi32:GetTextExtentPoint32(pointer hdc, [in] string text, long len, POINT* p)
var getTextExtentPoint32 = gdi32.NewProc("GetTextExtentPoint32A")
var _ = builtin4("GetTextExtentPoint32(hdc, text, len, p)",
	func(a, b, c, d Value) Value {
		var pt POINT
		rtn, _, _ := getTextExtentPoint32.Call(
			intArg(a),
			stringArg(b),
			uintptr(len(ToStr(b))),
			uintptr(unsafe.Pointer(&pt)))
		pointToOb(&pt, d)
		return boolRet(rtn)
	})

var getTextFace = gdi32.NewProc("GetTextFaceA")
var _ = builtin1("GetTextFace(hdc)",
	func(a Value) Value {
		const bufsize = 512
		var buf [bufsize]byte
		n, _, _ := getTextFace.Call(
			intArg(a),
			bufsize,
			uintptr(unsafe.Pointer(&buf[0])))
		return SuStr(string(buf[:n]))
	})

// dll long Gdi32:SetBkMode(pointer hdc, long mode)
var setBkMode = gdi32.NewProc("SetBkMode")
var _ = builtin2("SetBkMode(hdc, color)",
	func(a, b Value) Value {
		rtn, _, _ := setBkMode.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll long Gdi32:SetBkColor(pointer hdc, long color)
var setBkColor = gdi32.NewProc("SetBkColor")
var _ = builtin2("SetBkColor(hdc, color)",
	func(a, b Value) Value {
		rtn, _, _ := setBkColor.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Gdi32:DeleteDC(pointer hdc) bool
var deleteDC = gdi32.NewProc("DeleteDC")
var _ = builtin1("DeleteDC(hdc)",
	func(a Value) Value {
		rtn, _, _ := deleteDC.Call(intArg(a))
		return intRet(rtn)
	})

// dll Gdi32:DeleteObject(pointer hgdiobj) bool
var deleteObject = gdi32.NewProc("DeleteObject")
var _ = builtin1("DeleteObject(hgdiobj)", func(a Value) Value {
	rtn, _, _ := deleteObject.Call(intArg(a))
	return boolRet(rtn)
})

// dll Gdi32:GetClipBox(pointer hdc, RECT* rect) long
var getClipBox = gdi32.NewProc("GetClipBox")
var _ = builtin2("GetClipBox(hdc, rect)",
	func(a, b Value) Value {
		var r RECT
		rtn, _, _ := getClipBox.Call(
			intArg(a),
			rectArg(b, &r))
		return intRet(rtn)
	})

// dll Gdi32:GetPixel(pointer hdc, long nXPos, long nYPos) long
var getPixel = gdi32.NewProc("GetPixel")
var _ = builtin3("GetPixel(hdc, nXPos, nYPos)",
	func(a, b, c Value) Value {
		rtn, _, _ := getPixel.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		return intRet(rtn)
	})

// dll Gdi32:PtVisible(pointer hdc, long nXPos, long nYPos) bool
var ptVisible = gdi32.NewProc("PtVisible")
var _ = builtin3("PtVisible(hdc, nXPos, nYPos)",
	func(a, b, c Value) Value {
		rtn, _, _ := ptVisible.Call(
			intArg(a),
			intArg(b),
			intArg(c))
		if rtn == 0xffffffff { // 0xffffffff: 32 bit -1
			return IntVal(-1)
		}
		return boolRet(rtn)
	})

// dll Gdi32:Rectangle(pointer hdc, long left, long top, long right, long bottom) bool
var rectangle = gdi32.NewProc("Rectangle")
var _ = builtin5("Rectangle(hdc, left, top, right, bottom)",
	func(a, b, c, d, e Value) Value {
		rtn, _, _ := rectangle.Call(
			intArg(a),
			intArg(b),
			intArg(c),
			intArg(d),
			intArg(e))
		return intRet(rtn)
	})

// dll Gdi32:SetStretchBltMode(pointer hdc, long iStretchMode) long
var setStretchBltMode = gdi32.NewProc("SetStretchBltMode")
var _ = builtin2("SetStretchBltMode(hdc, iStretchMode)",
	func(a, b Value) Value {
		rtn, _, _ := setStretchBltMode.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Gdi32:SetTextColor(pointer hdc, long color) long
var setTextColor = gdi32.NewProc("SetTextColor")
var _ = builtin2("SetTextColor(hdc, color)",
	func(a, b Value) Value {
		rtn, _, _ := setTextColor.Call(
			intArg(a),
			intArg(b))
		return intRet(rtn)
	})

// dll Gdi32:StretchBlt(pointer hdcDest, long nXOriginDest, long nYOriginDest, long nWidthDest, long nHeightDest, pointer hdcSrc,
//	long nXOriginSrc, long nYOriginSrc, long nWidthSrc, long nHeightSrc, long dwRop) bool
var stretchBlt = gdi32.NewProc("StretchBlt")
var _ = builtin("StretchBlt(hdcDest, nXOriginDest, nYOriginDest, nWidthDest, nHeightDest, hdcSrc, nXOriginSrc, nYOriginSrc, nWidthSrc, nHeightSrc, dwRop)",
	func(_ *Thread, a []Value) Value {
		rtn, _, _ := stretchBlt.Call(
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
