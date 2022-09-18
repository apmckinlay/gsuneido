// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build windows && !portable

package builtin

import (
	"github.com/apmckinlay/gsuneido/builtin/goc"
	"github.com/apmckinlay/gsuneido/builtin/heap"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(Traccel, "(ob, msg)")

func Traccel(a, b Value) Value {
	defer heap.FreeTo(heap.CurSize())
	return IntVal(goc.Traccel(ToInt(a), obToMSG(b)))
}
