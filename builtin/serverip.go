package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("ServerIP()", func() Value {
	return SuStr("127.0.0.1") // TODO
})
