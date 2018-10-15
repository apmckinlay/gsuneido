package runtime

// ParamSpec describes the parameters of a function
type ParamSpec struct {
	// Nparams is the number of arguments required on the stack
	Nparams int

	// NDefaults is the number of default values for parameters
	// They are in the start of Values
	Ndefaults int

	// Flags specifies "types" of params
	Flags []Flag

	// Strings starts with the parameter names, then the local names,
	// and then any argument or member names used in the code,
	// and any argument specs
	Strings []string

	// Values contains any literals in the function
	// starting with parameter defaults
	Values []Value
}

// Flag is a bit set of parameter options
type Flag byte

const (
	AtParam Flag = 1 << iota
	DynParam
	DotParam
	PubParam
)

func (f *ParamSpec) Params() *ParamSpec {
	return f
}
