package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Use(library)",
	func(t *Thread, args ...Value) Value {
		return SuBool(t.Dbms().Use(IfStr(args[0])))
	})
