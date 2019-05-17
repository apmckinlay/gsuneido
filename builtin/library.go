package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Libraries()", func(t *Thread, args ...Value) Value {
	return t.Dbms().Libraries()
})

var _ = builtin("Use(library)",
	func(t *Thread, args ...Value) Value {
		return SuBool(t.Dbms().Use(IfStr(args[0])))
	})

var _ = builtin("Unuse(library)",
	func(t *Thread, args ...Value) Value {
		return SuBool(t.Dbms().Unuse(IfStr(args[0])))
	})
