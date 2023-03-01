// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
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
	defer heap.FreeTo(heap.CurSize())
	iunk := heap.Alloc(int64Size)
	pPtr := heap.Alloc(int64Size)
	rtn := goc.EmbedBrowserObject(
		intArg(a),
		uintptr(iunk),
		uintptr(pPtr))
	if rtn != 0 {
		return intRet(rtn)
	}
	iOleObject := *(*uintptr)(iunk)
	swb := &suWebBrowser{iOleObject: iOleObject, ptr: *(*uintptr)(pPtr)}
	idisp := goc.QueryIDispatch(iOleObject)
	swb.suCOMObject = suCOMObject{ptr: idisp, idisp: true}
	return swb
}

var suWebBrowserMethods = methods()

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

func (swb *suWebBrowser) Lookup(th *Thread, method string) Callable {
	if f, ok := suWebBrowserMethods[method]; ok {
		return f
	}
	return swb.suCOMObject.Lookup(th, method)
}
