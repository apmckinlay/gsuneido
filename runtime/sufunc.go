package runtime

import "github.com/apmckinlay/gsuneido/util/str"

// MaxArgs is the maximum number of arguments allowed
const MaxArgs = 200

// SuFunc is a compiled Suneido function, method, or block
type SuFunc struct {
	ParamSpec

	// Nlocals is the number of parameters and local variables
	Nlocals uint8

	// Code is the actual byte code
	Code []byte

	// ArgSpecs used by calls in the code
	ArgSpecs []*ArgSpec

	// ClassName is used to privatize dot params
	ClassName string
}

var _ Value = (*SuFunc)(nil) // verify *SuFunc satisfies Value

func (f *SuFunc) Call(t *Thread, as *ArgSpec) Value {
	args := t.Args(&f.ParamSpec, as)
	for i, flag := range f.Flags {
		if flag&DotParam == DotParam {
			name := f.Names[i]
			if flag&PubParam == PubParam {
				name = str.Capitalize(name)
			} else { // privatize
				name = f.ClassName + "_" + name
			}
			t.this.Put(SuStr(name), args[i])
		}
	}
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

func (f *SuFunc) String() string {
	s := ""
	if f.Name != "" {
		s = f.Name + " "
	}
	s += "/* "
	if f.ClassName != "" {
		s += "method"
	} else {
		s += "function"
	}
	s += " */"
	return s
}
