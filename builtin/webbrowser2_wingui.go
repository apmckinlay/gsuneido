// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
)

// must match webview2_ops in webview2.cpp
const (
	webview2_create = iota
	webview2_close
	webview2_resize
	webview2_navigate
	webview2_navigate_to_string
	webview2_execute_script
	webview2_get_source
	webview2_print
	webview2_set_focus
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

func WebBrowser2(a, b, s, cb Value) Value {
	defer heap.FreeTo(heap.CurSize())
	iunk := heap.Alloc(int64Size)
	rtn := goc.WebBrowser2(webview2_create, intArg(a), uintptr(iunk), uintptr(stringArg(b)), uintptr(stringArg(s)), NewCallback(cb, 4))
	if rtn != 0 {
		return intRet(rtn)
	}
	iOleObject := *(*uintptr)(iunk)
	swb := &suWebBrowser2{iOleObject: iOleObject}
	return swb
}

var suWebBrowser2Methods = methods()

var _ = method(web2_Release, "()")

func web2_Release(this Value) Value {
	wb := this.(*suWebBrowser2)
	rtn := goc.WebBrowser2(webview2_close, wb.iOleObject, 0, 0, 0, 0)
	return intRet(rtn)
}

var _ = method(web2_Resize, "(w, h)")

func web2_Resize(this Value, _w Value, _h Value) Value {
	wb := this.(*suWebBrowser2)
	w := ToInt(_w)
	h := ToInt(_h)
	rtn := goc.WebBrowser2(webview2_resize, wb.iOleObject, uintptr(w), uintptr(h), 0, 0)
	return intRet(rtn)
}

var _ = method(web2_Navigate, "(s)")

func web2_Navigate(this Value, s Value) Value {
	defer heap.FreeTo(heap.CurSize())
	wb := this.(*suWebBrowser2)
	rtn := goc.WebBrowser2(webview2_navigate, wb.iOleObject, uintptr(stringArg(s)), 0, 0, 0)
	return intRet(rtn)
}

var _ = method(web2_NavigateToString, "(s)")

func web2_NavigateToString(this Value, s Value) Value {
	defer heap.FreeTo(heap.CurSize())
	wb := this.(*suWebBrowser2)
	rtn := goc.WebBrowser2(webview2_navigate_to_string, wb.iOleObject, uintptr(stringArg(s)), 0, 0, 0)
	return intRet(rtn)
}

var _ = method(web2_ExecuteScript, "(script)")

func web2_ExecuteScript(this Value, script Value) Value {
	defer heap.FreeTo(heap.CurSize())
	wb := this.(*suWebBrowser2)
	rtn := goc.WebBrowser2(webview2_execute_script, wb.iOleObject, uintptr(stringArg(script)), 0, 0, 0)
	return intRet(rtn)
}

var _ = method(web2_GetSource, "()")

func web2_GetSource(this Value) Value {
	defer heap.FreeTo(heap.CurSize())
	wb := this.(*suWebBrowser2)
	buf := heap.Alloc(MAX_PATH)
	rtn := goc.WebBrowser2(webview2_get_source, wb.iOleObject, uintptr(buf), 0, 0, 0)
	if rtn != 0 {
		return EmptyStr
	}
	return SuStr(heap.GetStrZ(buf, MAX_PATH))
}

var _ = method(web2_Print, "()")

func web2_Print(this Value) Value {
	defer heap.FreeTo(heap.CurSize())
	wb := this.(*suWebBrowser2)
	rtn := goc.WebBrowser2(webview2_print, wb.iOleObject, 0, 0, 0, 0)
	return intRet(rtn)
}

var _ = method(web2_SetFocus, "()")

func web2_SetFocus(this Value) Value {
	defer heap.FreeTo(heap.CurSize())
	wb := this.(*suWebBrowser2)
	rtn := goc.WebBrowser2(webview2_set_focus, wb.iOleObject, 0, 0, 0, 0)
	return intRet(rtn)
}

func (swb *suWebBrowser2) Lookup(th *Thread, method string) Callable {
	return suWebBrowser2Methods[method];
}