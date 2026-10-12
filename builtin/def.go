// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
)

// DefDef must be called to make Def available
func DefDef() {
	builtin(Def, "(name :string, definition)")
}

func Def(nameVal, val Value) Value {
	name := ToStr(nameVal)
	if s, ok := val.ToStr(); ok {
		val = compile.NamedConstant("Def", name, s, nil)
	}
	Global.TestDef(name, val)
	return val
}
