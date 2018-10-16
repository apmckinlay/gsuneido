package runtime

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

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

// Value interface (except TypeName)

func (f *ParamSpec) ToInt() int {
	panic("cannot convert function to integer")
}

func (f *ParamSpec) ToDnum() dnum.Dnum {
	panic("cannot convert function to number")
}

func (f *ParamSpec) ToStr() string {
	panic("cannot convert function to string")
}

func (f *ParamSpec) String() string {
	var buf strings.Builder
	buf.WriteString("function(")
	sep := ""
	v := 0 // index into Values
	for i := 0; i < f.Nparams; i++ {
		buf.WriteString(sep)
		buf.WriteString(flagsToName(f.Strings[i], f.Flags[i]))
		if i >= f.Nparams-f.Ndefaults {
			buf.WriteString("=")
			buf.WriteString(fmt.Sprint(f.Values[v]))
			v++
		}
		sep = ", "
	}
	buf.WriteString(")")
	return buf.String()
}

func flagsToName(p string, flags Flag) string {
	if flags == AtParam {
		p = "@" + p
	}
	if flags&PubParam == PubParam {
		p = str.Capitalize(p)
	}
	if flags&DynParam == DynParam {
		p = "_" + p
	}
	if flags&DotParam == DotParam {
		p = "." + p
	}
	return p
}

func (f *ParamSpec) Get(Value) Value {
	panic("function does not support get")
}

func (f *ParamSpec) Put(Value, Value) {
	panic("function does not support put")
}

func (ParamSpec) RangeTo(int, int) Value {
	panic("function does not support range")
}

func (ParamSpec) RangeLen(int, int) Value {
	panic("function does not support range")
}

func (f *ParamSpec) Hash() uint32 {
	panic("function hash not implemented")
}

func (f *ParamSpec) Hash2() uint32 {
	panic("function hash not implemented")
}

func (f *ParamSpec) Equal(other interface{}) bool {
	if f2, ok := other.(*ParamSpec); ok {
		return f == f2
	}
	return false
}

func (*ParamSpec) Order() Ord {
	return OrdOther
}

func (f *ParamSpec) Compare(Value) int {
	panic("function compare not implemented")
}
