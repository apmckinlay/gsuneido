package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

// These methods are specific for int values

func init() {
	IntMethods = Methods{
		"Int": method0(func(this Value) Value {
			return this
		}),
		"Frac": method0(func(Value) Value {
			return Zero
		}),
	}
}
