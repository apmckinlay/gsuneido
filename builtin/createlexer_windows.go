// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("CreateLexer(name)",
	func(a Value) Value {
		defer heap.FreeTo(heap.CurSize())
		rtn := goc.CreateLexer(uintptr(heap.CopyStr(ToStr(a))))
		return intRet(rtn)
	})
