package base

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

func (f *SuFunc) Call(ct Context, self Value, args ...Value) Value {
	return ct.Call(f, self)
}

var _ Value = (*SuFunc)(nil) // verify SuFunc satisfies Value

func (_ *SuFunc) TypeName() string {
	return "Function"
}
