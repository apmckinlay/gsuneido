package runtime

// MaxArgs is the maximum number of arguments allowed
const MaxArgs = 200

// SuFunc is a compiled Suneido function, method, or block
// Not worth specializing SuFunc0..4 for interpreted functions
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
func (f *SuFunc) Call0(t *Thread) Value {
	return f.Call(t, ArgSpec0)
}
func (f *SuFunc) Call1(t *Thread, _ Value) Value {
	return f.Call(t, ArgSpec1)
}
func (f *SuFunc) Call2(t *Thread, _, _ Value) Value {
	return f.Call(t, ArgSpec2)
}
func (f *SuFunc) Call3(t *Thread, _, _, _ Value) Value {
	return f.Call(t, ArgSpec3)
}
func (f *SuFunc) Call4(t *Thread, _, _, _, _ Value) Value {
	return f.Call(t, ArgSpec4)
}

var _ Value = (*SuFunc)(nil) // verify *SuFunc satisfies Value

// TypeName returns the Suneido name for the type (Value interface)
func (*SuFunc) TypeName() string {
	return "Function"
}

// SuFuncMethods is initialized by the builtin package
var SuFuncMethods Methods

func (*SuFunc) Lookup(method string) Callable {
	return SuFuncMethods[method]
}

func (*SuFunc) String() string {
	return "/* function */" // TODO name and library
}
