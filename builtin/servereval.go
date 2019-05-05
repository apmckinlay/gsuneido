package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("ServerEval(@args)", func(t *Thread, args ...Value) Value {
	return nilToEmptyStr(t.Dbms().Exec(t, args[0]))
})

func nilToEmptyStr(v Value) Value {
	if v == nil {
		return EmptyStr
	}
	return v
}
