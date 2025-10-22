// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && gui

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	. "github.com/apmckinlay/gsuneido/core"
)

type suWebBrowser struct {
	suCOMObject
	iOleObject uintptr
	ptr        uintptr
}

func (*suWebBrowser) String() string {
	return "WebBrowser"
}

var _ = builtin(WebBrowser, "(hwnd)")

func WebBrowser(a Value) Value {
	rtn, iunk, ptr := goc.EmbedBrowserObject(intArg(a))
	if rtn != 0 {
		return intRet(rtn)
	}
	swb := &suWebBrowser{iOleObject: iunk, ptr: ptr}
	idisp := goc.QueryIDispatch(iunk)
	swb.suCOMObject = suCOMObject{ptr: idisp, idisp: true}
	return swb
}

var suWebBrowserMethods = methods("web")

var _ = method(web_Release, "()")

func web_Release(this Value) Value {
	wb := this.(*suWebBrowser)
	goc.Release(wb.suCOMObject.ptr)
	goc.UnEmbedBrowserObject(wb.iOleObject, wb.ptr)
	return nil
}

var _ = method(web_GetIOleObject, "()")

func web_GetIOleObject(this Value) Value {
	return IntVal((int)(this.(*suWebBrowser).iOleObject))
}

func (swb *suWebBrowser) Lookup(th *Thread, method string) Value {
	if f, ok := suWebBrowserMethods[method]; ok {
		return f
	}
	return swb.suCOMObject.Lookup(th, method)
}
