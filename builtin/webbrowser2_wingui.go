// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"unsafe"

	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/core"
)

type suWebBrowser2 struct {
	ValueBase[*suWebBrowser2]
	iOleObject uintptr
}

func (*suWebBrowser2) String() string {
	return "WebBrowser2"
}

func (swb *suWebBrowser2) Equal(other any) bool {
	return swb == other
}

func (*suWebBrowser2) SetConcurrent() {
	// ok since immutable (assuming the COM object is thread safe)
}

var _ = builtin(WebBrowser2, "(hwnd, dllPath, userDataFolder, cb)")

func WebBrowser2(th *Thread, args []Value) Value {
	var iunk uintptr
	rtn := goc.WebView2_Create(
		intArg(args[0]),
		unsafe.Pointer(&iunk),
		ToStr(args[1]),
		ToStr(args[2]),
		NewCallback(th, args[3], 4))
	if rtn != 0 {
		return intRet(rtn)
	}
	return &suWebBrowser2{iOleObject: iunk}
}

var suWebBrowser2Methods = methods("web2")

var _ = method(web2_Release, "()")

func web2_Release(this Value) Value {
	wb := this.(*suWebBrowser2)
	rtn := goc.WebView2_Close(wb.iOleObject)
	return intRet(rtn)
}

var _ = method(web2_Resize, "(w, h)")

func web2_Resize(this Value, _w Value, _h Value) Value {
	wb := this.(*suWebBrowser2)
	w := ToInt(_w)
	h := ToInt(_h)
	rtn := goc.WebView2_Resize(wb.iOleObject, uintptr(w), uintptr(h))
	return intRet(rtn)
}

var _ = method(web2_Navigate, "(s)")

func web2_Navigate(this Value, s Value) Value {
	wb := this.(*suWebBrowser2)
	rtn := goc.WebView2_Navigate(wb.iOleObject, ToStr(s))
	return intRet(rtn)
}

var _ = method(web2_NavigateToString, "(s)")

func web2_NavigateToString(this Value, s Value) Value {
	wb := this.(*suWebBrowser2)
	rtn := goc.WebView2_NavigateToString(wb.iOleObject, ToStr(s))
	return intRet(rtn)
}

var _ = method(web2_ExecuteScript, "(script)")

func web2_ExecuteScript(this Value, script Value) Value {
	wb := this.(*suWebBrowser2)
	rtn := goc.WebView2_ExecuteScript(wb.iOleObject, ToStr(script))
	return intRet(rtn)
}

var _ = method(web2_GetSource, "()")

func web2_GetSource(this Value) Value {
	wb := this.(*suWebBrowser2)
	buf := make([]byte, MAX_PATH)
	rtn := goc.WebView2_GetSource(wb.iOleObject, &buf[0])
	if rtn != 0 {
		return EmptyStr
	}
	return bufZstr(buf)
}

var _ = method(web2_Print, "()")

func web2_Print(this Value) Value {
	wb := this.(*suWebBrowser2)
	rtn := goc.WebView2_Print(wb.iOleObject)
	return intRet(rtn)
}

var _ = method(web2_SetFocus, "()")

func web2_SetFocus(this Value) Value {
	wb := this.(*suWebBrowser2)
	rtn := goc.WebView2_SetFocus(wb.iOleObject)
	return intRet(rtn)
}

func (swb *suWebBrowser2) Lookup(th *Thread, method string) Value {
	return suWebBrowser2Methods[method]
}
