package runtime

import "github.com/apmckinlay/gsuneido/util/str"

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

	// Id is a unique identifier for a function defining blocks
	Id uint32
	// OuterId is the Id of the outer SuFunc
	// It is used by interp to handle block return
	OuterId uint32
}

// Value interface (mostly handled by ParamSpec) --------------------

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

func (f *SuFunc) TypeName() string {
	if f.OuterId != 0 {
		return "Block"
	}
	return "Function"
}

// SuFuncMethods is initialized by the builtin package
var SuFuncMethods Methods

func (*SuFunc) Lookup(method string) Value {
	if m,ok := ParamsMethods[method]; ok {
		return m
	}
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
