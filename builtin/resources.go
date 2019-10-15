package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin0("ResourceCounts()", func() Value {
	return NewSuObject()
})
