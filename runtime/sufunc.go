package runtime

// MaxArgs is the maximum number of arguments allowed
const MaxArgs = 200

// SuFunc is a compiled Suneido function, method, or block
type SuFunc struct {
	ParamSpec

	// Nlocals is the number of parameters and local variables
	Nlocals uint8

	// Code is the actual byte code
	Code []byte

	// the ArgSpec's used by calls in the code
	ArgSpecs []*ArgSpec
}

var _ Value = (*SuFunc)(nil) // verify *SuFunc satisfies Value

func (f *SuFunc) Call(t *Thread, as *ArgSpec) Value {
	t.Args(&f.ParamSpec, as)
	return t.Call(f)
}

// TypeName returns the Suneido name for the type (Value interface)
func (*SuFunc) TypeName() string {
	return "Function"
}

// SuFuncMethods is initialized by the builtin package
var SuFuncMethods Methods

func (*SuFunc) Lookup(method string) Value {
	return SuFuncMethods[method]
}

func (*SuFunc) String() string {
	return "/* function */" // TODO name and library
}
