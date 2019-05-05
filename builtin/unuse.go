package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Unuse(library)",
	func(t *Thread, args ...Value) Value {
		return SuBool(t.Dbms().Unuse(IfStr(args[0])))
	})
