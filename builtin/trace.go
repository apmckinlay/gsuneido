package builtin

import (
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("Trace(flags, block)",
	func(t *Thread, args []Value) Value {
		defer func(oldflags int) {
			options.Trace = oldflags
		}(options.Trace)
		options.Trace = ToInt(args[0])
		t.Call(args[1])
		return nil
	})
