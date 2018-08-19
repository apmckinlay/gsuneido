package base

// MaxArgs is the maximum number of arguments allowed
const MaxArgs = 200

// SuFunc is a compiled Suneido function, method, or block
type SuFunc struct {
	Func

	// Nlocals is the number of parameters and local variables
	Nlocals int

	// Code is the actual byte code
	Code []byte
}

var _ Callable = (*SuFunc)(nil) // verify SuFunc satisfies Callable

// Call invokes the SuFunc
func (f *SuFunc) Call(ct Context, self Value, _ ...Value) Value {
	// TODO push args ???
	return ct.Call(f, self)
}

var _ Value = (*SuFunc)(nil) // verify *SuFunc satisfies Value

// TypeName returns the Suneido name for the type (Value interface)
func (*SuFunc) TypeName() string {
	return "Function"
}
