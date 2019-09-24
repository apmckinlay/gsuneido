package builtin

import (
	"sync"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var builtinNamesOnce sync.Once
var builtinNames *SuObject

var _ = builtin0("BuiltinNames()", func() Value {
	builtinNamesOnce.Do(func() {
		builtinNames = NewSuObject(BuiltinNames()...)
		builtinNames.Sort(nil, False)
		builtinNames.SetReadOnly()
	})
	return builtinNames
})
