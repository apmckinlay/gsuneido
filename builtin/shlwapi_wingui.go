// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"syscall"

	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/core"
)

var shlwapi = MustLoadDLL("shlwapi.dll")

// dll long SHCreateStreamOnFileA(string pszFile, int32 grfMode, POINTER* ppstm)
var shCreateStreamOnFile = shlwapi.MustFindProc("SHCreateStreamOnFileA").Addr()
var _ = builtin(SHCreateStreamOnFile, "(pszFile, grfMode, ppstm)")

func SHCreateStreamOnFile(a, b, c Value) Value {
	defer heap.FreeTo(heap.CurSize())
	p := heap.Alloc(uintptrSize)
	rtn, _, _ := syscall.SyscallN(shCreateStreamOnFile,
		uintptr(stringArg(a)),
		intArg(b),
		uintptr(p))
	pstm := *(*uintptr)(p)
	c.Put(nil, SuStr("x"), IntVal(int(pstm)))
	return intRet(rtn)
}
