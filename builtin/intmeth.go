package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

// These methods are specific for int values

func init() {
	IntMethods = Methods{
		"Int": method0(func(self Value) Value {
			return self
		}),
		"Frac": method0(func(self Value) Value {
			return Zero
		}),
	}
}
