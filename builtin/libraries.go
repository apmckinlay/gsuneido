package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Libraries()", func(t *Thread, args ...Value) Value {
	return t.Dbms().Libraries()
})
