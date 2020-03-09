// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var shlwapi = MustLoadDLL("shlwapi.dll")

// dll long SHCreateStreamOnFileA(string pszFile, int32 grfMode, POINTER* ppstm)
var shCreateStreamOnFile = shlwapi.MustFindProc("SHCreateStreamOnFileA").Addr()
var _ = builtin3("SHCreateStreamOnFile(pszFile, grfMode, ppstm)",
	func(a, b, c Value) Value {
		defer heap.FreeTo(heap.CurSize())
		p := heap.Alloc(uintptrSize)
		rtn := goc.Syscall3(shCreateStreamOnFile,
			uintptr(stringArg(a)),
			intArg(b),
			uintptr(p))
		pstm := *(*uintptr)(p)
		c.Put(nil, SuStr("x"), IntVal(int(pstm)))
		return intRet(rtn)
	})
