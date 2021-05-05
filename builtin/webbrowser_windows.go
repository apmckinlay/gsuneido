// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
)

type suWebBrowser struct {
	suComObject
	iOleObject uintptr
	ptr uintptr
}

func (*suWebBrowser) String() string {
	return "WebBrowser"
}

var _ = builtin1("WebBrowser(hwnd)",
	func(a Value) Value {
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
		iOleObject :=*(*uintptr)(iunk)
		swb := &suWebBrowser{iOleObject: iOleObject, ptr: *(*uintptr)(pPtr)}
		idisp := goc.QueryIDispatch(iOleObject)
		swb.suComObject = suComObject{ptr: idisp, idisp: true}
		return swb
	})

var suWebBrowserMethods = Methods{
	"Release": method0(func(this Value) Value {
		wb := this.(*suWebBrowser)
		goc.Release(wb.suComObject.ptr)
		goc.UnEmbedBrowserObject(wb.iOleObject, wb.ptr)
		return nil
	}),
	"GetIOleObject": method0(func(this Value) Value {
		return IntVal((int)(this.(*suWebBrowser).iOleObject))
	}),
}

func (swb *suWebBrowser) Lookup(t *Thread, method string) Callable {
	if f, ok := suWebBrowserMethods[method]; ok {
		return f
	}
	return swb.suComObject.Lookup(t, method)
}
