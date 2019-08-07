package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin0("BuiltinNames()", func() Value {
	list := NewSuObject(BuiltinNames()...)
	list.Sort(nil, False)
	return list
})
