package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("BuiltinNames()", func() Value {
	return NewSuObject(BuiltinNames()...)
})
