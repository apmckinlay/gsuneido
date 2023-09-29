// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
)

// DefDef must be called to make Def available
func DefDef() {
	builtin(Def, "(name, definition)")
}

func Def(nameVal, val Value) Value {
	name := string(nameVal.(SuStr))
	if ss, ok := val.(SuStr); ok {
		val = compile.NamedConstant("Def", name, string(ss), nil)
	}
	Global.TestDef(name, val)
	return val
}
