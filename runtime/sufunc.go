package runtime

// MaxArgs is the maximum number of arguments allowed
const MaxArgs = 200

// SuFunc is a compiled Suneido function, method, or block
type SuFunc struct {
	ParamSpec

	// Nlocals is the number of parameters and local variables
	Nlocals int

	// Code is the actual byte code
	Code []byte
}

// Call invokes the SuFunc
func (f *SuFunc) Call(t *Thread, as *ArgSpec) Value {
	t.Args(&f.ParamSpec, as)
	return t.Call(f)
}

var _ Value = (*SuFunc)(nil) // verify *SuFunc satisfies Value

// TypeName returns the Suneido name for the type (Value interface)
func (*SuFunc) TypeName() string {
	return "Function"
}

func (*SuFunc) Lookup(string) Callable {
	return nil // TODO
}
