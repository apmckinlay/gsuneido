package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin0("Timestamp()", func() Value {
	//TODO client/server, guarantee unique value
	return Now()
})
