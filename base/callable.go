package base

import (
	"bytes"
	"fmt"

	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Context is the current Thread
// can't use actual Thread because that causes circular dependency
type Context interface {
	// Call interprets a compiled Suneido function
	Call(fn *SuFunc, self Value) Value
}

// Callable is how things are called
// Assumes args on stack, already massaged if necessary
type Callable interface {
	Params() *Func
	Call(ct Context, self Value, args ...Value) Value
}

type Func struct {
	// Nparams is the number of arguments required on the stack
	Nparams int

	// NDefaults is the number of default values for parameters
	// They are in the start of Values
	Ndefaults int

	// Flags specifies "types" of params
	Flags []Flag

	// Strings starts with the parameters, then the locals,
	// and then any argument or member names used in the code,
	// and any argument specs
	Strings []string

	// Values contains any literals in the function
	// starting with parameter defaults
	Values []Value
}

type Flag byte

const (
	AT_F Flag = 1 << iota
	DYN_F
	DOT_F
	PUB_F
)

func (f *Func) Params() *Func {
	return f
}

// Value stuff ------------------------------------------------------
// note: does not implement TypeName

func (f *Func) ToInt() int {
	panic("cannot convert function to integer")
}

func (f *Func) ToDnum() dnum.Dnum {
	panic("cannot convert function to number")
}

func (f *Func) ToStr() string {
	panic("cannot convert function to string")
}

func (f *Func) String() string {
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

func (f *Func) Get(key Value) Value {
	panic("function does not support get")
}

func (f *Func) Put(key Value, val Value) {
	panic("function does not support put")
}

func (f *Func) Hash() uint32 {
	panic("function hash not implemented")
}

func (f *Func) hash2() uint32 {
	panic("function hash not implemented")
}

func (f *Func) Equals(other interface{}) bool {
	if f2, ok := other.(*Func); ok {
		return f == f2
	}
	return false
}

func (_ *Func) Order() ord {
	return ordOther
}

func (f *Func) Cmp(other Value) int {
	panic("function compare not implemented")
}
