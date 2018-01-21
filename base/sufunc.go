package base

import (
	"bytes"
	"fmt"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

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
	// nparams is the number of values required on the stack
	Nparams int
	// nlocals is the number of parameters and local variables
	Nlocals int
	// strings starts with the parameters, then the locals,
	// and then any argument or member names used in the code,
	// and any argument specs
	Strings []string
	Values  []Value
}

var _ Value = &SuFunc{} // confirm it implements Value

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
		p := f.Strings[i]
		if p[0] == '=' {
			buf.WriteString(p[1:])
			buf.WriteString("=")
			buf.WriteString(fmt.Sprint(f.Values[v]))
			v++
		} else {
			buf.WriteString(p)
		}
		sep = ", "
	}
	buf.WriteString(")")
	return buf.String()
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

func (f *SuFunc) Call(c CallContext) Value {
	return c.CallSuFunc(f)
}
