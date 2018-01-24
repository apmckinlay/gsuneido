package base

import (
	"bytes"
	"fmt"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

const MaxArgs = 200

/*
SuFunc is a compiled function, method, or block.

Parameters at the start of names may be prefixed:
'@' for each, '_' for dynamic, '.' for member, or '=' for default.

There can only be one '@' parameter and if present it must be last.

Parameters with default values ('=' or '_')
must come after parameters without defaults.

'.' parameter names for public members are capitalized.
*/
type SuFunc struct {
	// Code is the actual byte code
	Code []byte

	// Nparams is the number of arguments required on the stack
	Nparams int

	// NDefaults is the number of default values for parameters
	// They are in the start of Values
	Ndefaults int

	// Flags specifies "types" of params
	Flags []Flag

	// Nlocals is the number of parameters plus local variables
	Nlocals int

	// Strings starts with the parameters, then the locals,
	// and then any argument or member names used in the code,
	// and any argument specs
	Strings []string

	// Values contains any literals in the function
	// starting with parameter defaults
	Values []Value
}

var _ Value = (*SuFunc)(nil) // confirm it implements Value

type Flag byte

const (
	AT_F Flag = 1 << iota
	DYN_F
	DOT_F
	PUB_F
)

func (f *SuFunc) ToInt() int32 {
	panic("cannot convert function to integer")
}

func (f *SuFunc) ToDnum() dnum.Dnum {
	panic("cannot convert function to number")
}

func (f *SuFunc) ToStr() string {
	panic("cannot convert function to string")
}

func (f *SuFunc) String() string {
	var buf bytes.Buffer
	buf.WriteString("function (")
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
	if flags == AT_F {
		p = "@" + p
	}
	if flags&PUB_F == PUB_F {
		p = str.Capitalize(p)
	}
	if flags&DYN_F == DYN_F {
		p = "_" + p
	}
	if flags&DOT_F == DOT_F {
		p = "." + p
	}
	return p
}

func (f *SuFunc) Get(key Value) Value {
	panic("function does not support get")
}

func (f *SuFunc) Put(key Value, val Value) {
	panic("function does not support put")
}

func (f *SuFunc) Hash() uint32 {
	panic("function hash not implemented")
}

func (f *SuFunc) hash2() uint32 {
	panic("function hash not implemented")
}

func (f *SuFunc) Equals(other interface{}) bool {
	if f2, ok := other.(*SuFunc); ok {
		return f == f2
	}
	return false
}

func (_ *SuFunc) TypeName() string {
	return "Function"
}

func (_ *SuFunc) Order() ord {
	return ordOther
}

func (f *SuFunc) Cmp(other Value) int {
	panic("function compare not implemented")
}
