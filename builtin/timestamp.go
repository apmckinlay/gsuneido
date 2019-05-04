package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var prevTimestamp SuDate

var _ = builtin("Timestamp()", func(t *Thread, args ...Value) Value {
	return t.Dbms().Timestamp()
})
