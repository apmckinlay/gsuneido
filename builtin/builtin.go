package builtin

import (
	. "github.com/apmckinlay/gsuneido/base"
)

// Builtin wraps a built-in function
// so we can dynamically create a Value from a bare function
type Builtin struct {
	Func
	fn func(ct Context, self Value, args ...Value) Value
}

var _ Callable = (*Builtin)(nil)

func (f *Builtin) Call(ct Context, self Value, args ...Value) Value {
	return f.fn(ct, self, args...)
}

func (*Builtin) TypeName() string {
	return "BuiltinFunction" //TODO
}

var _ Value = (*Builtin)(nil) // verify Func satisfies Value
