package builtin

import (
	v "github.com/apmckinlay/gsuneido/value"
)

// Type is a builtin function that returns a value's type as a string
func Type(x v.Value) v.Value {
	return v.SuStr(x.TypeName())
}
