// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(UninitQ, "(name :string) :boolean")

func UninitQ(th *Thread, args []Value) (result Value) {
	name := ToStr(args[0])
	if !compile.IsGlobal(name) {
		panic("Uninit? requires a global")
	}
	defer func() {
		if e := recover(); e != nil {
			OpCatch(th, e, "error loading")
			result = True
		}
	}()
	return SuBool(Global.FindName(th, name) == nil)
}
