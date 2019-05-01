package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var prevTimestamp SuDate

var _ = builtin0("Timestamp()", func() Value {
	//TODO client/server, concurrency
	t := Now()
	if t.Equal(prevTimestamp) {
		t = t.Plus(0, 0, 0, 0, 0, 0, 1)
	}
	prevTimestamp = t
	return t
})
