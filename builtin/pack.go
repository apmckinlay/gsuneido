// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

var _ = builtin(packSize, "(value)")

func packSize(arg Value) Value {
	return IntVal(PackSize(arg))
}

var _ = builtin(pack, "(value)")

func pack(arg Value) Value {
	p, ok := arg.(Packable)
	if !ok {
		panic("can't pack " + ErrType(arg))
	}
	enc := Pack2(p)
	buf := enc.Buffer()
	if PackString != 0 {
		// convert to old pack format for interoperability
		RevertValue(buf, enc.String())
	}
	return SuStr(hacks.BStoS(buf))
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
