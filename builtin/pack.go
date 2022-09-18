// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(packSize, "(value)")

func packSize(arg Value) Value {
	return IntVal(PackSize(arg))
}

var _ = builtin(pack, "(value)")

func pack(arg Value) Value {
	return SuStr(PackValue(arg))
}

var _ = builtin(unpack, "(string)")

func unpack(arg Value) Value {
	defer func() {
		if e := recover(); e != nil {
			panic("Unpack: not a valid packed value")
		}
	}()
	return Unpack(ToStr(arg))
}
