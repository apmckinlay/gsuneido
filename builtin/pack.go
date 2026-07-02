// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(packSize, "(value) :number")

func packSize(arg Value) Value {
	return IntVal(PackSize(arg))
}

var _ = builtin(pack, "(value) :string")

func pack(arg Value) Value {
	return SuStr(PackValue(arg))
}

var _ = builtin(unpack, "(string) :unknown")

func unpack(arg Value) Value {
	defer func() {
		if e := recover(); e != nil {
			panic("Unpack: not a valid packed value")
		}
	}()
	return Unpack(ToStr(arg))
}
