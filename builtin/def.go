package builtin

import (
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// Def must be called to be available
func Def() {
	builtin2("Def(name, definition)", func(nameVal, val Value) Value {
		name := string(nameVal.(SuStr))
		if ss, ok := val.(SuStr); ok {
			val = compile.NamedConstant("Def", name, string(ss))
		}
		Global.TestDef(name, val)
		return val
	})
}
