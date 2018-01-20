package builtin

import (
	. "github.com/apmckinlay/gsuneido/interp"
)

// Type is a builtin function that returns a value's type as a string
func Type(x Value) Value {
	return SuStr(x.TypeName())
}
