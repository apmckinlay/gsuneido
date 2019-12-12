package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("ServerEval(@args)", func(t *Thread, args []Value) Value {
	return t.Dbms().Exec(t, args[0])
})
