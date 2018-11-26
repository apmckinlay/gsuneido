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
	Nparams uint8

	// NDefaults is the number of default values for parameters
	// They are in the start of Values
	Ndefaults uint8

	// Flags specifies "types" of params
	Flags []Flag

	// Names starts with the parameter names, then the local names
	Names []string

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

// Value interface (except TypeName)

func (*ParamSpec) ToInt() int {
	panic("cannot convert function to integer")
}

func (*ParamSpec) ToDnum() dnum.Dnum {
	panic("cannot convert function to number")
}

func (*ParamSpec) ToStr() string {
	panic("cannot convert function to string")
}

func (f *ParamSpec) String() string {
	var buf strings.Builder
	// easier to add "function" here and strip it in Params
	// than to implement String in a bunch of places to add it
	buf.WriteString("function(")
	sep := ""
	v := 0 // index into Values
	for i := 0; i < int(f.Nparams); i++ {
		if i == 0 && f.Names[0] == "this" {
			continue
		}
		buf.WriteString(sep)
		buf.WriteString(flagsToName(f.Names[i], f.Flags[i]))
		if i >= int(f.Nparams-f.Ndefaults) {
			buf.WriteString("=")
			buf.WriteString(fmt.Sprint(f.Values[v]))
			v++
		}
		sep = ","
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
	// don't include DotParam - not relevant to caller
	return p
}

func (*ParamSpec) Get(Value) Value {
	panic("function does not support get")
}

func (*ParamSpec) Put(Value, Value) {
	panic("function does not support put")
}

func (ParamSpec) RangeTo(int, int) Value {
	panic("function does not support range")
}

func (ParamSpec) RangeLen(int, int) Value {
	panic("function does not support range")
}

func (*ParamSpec) Hash() uint32 {
	panic("function hash not implemented")
}

func (*ParamSpec) Hash2() uint32 {
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

func (*ParamSpec) Compare(Value) int {
	panic("function compare not implemented")
}

// Params is set in the builtin package
var Params Callable

func (*ParamSpec) Lookup(method string) Callable {
	if method == "Params" {
		return Params
	}
	return nil
}

func (f *ParamSpec) Params() string {
	return f.String()[8:] // skip "function"
}

// defaults for Builtin# and Method#
// not enough arguments is handled elsewhere

func (f *ParamSpec) Call0(*Thread) Value {
	panic("too many arguments")
}
func (f *ParamSpec) Call1(_ *Thread, _ Value) Value {
	panic("too many arguments")
}
func (f *ParamSpec) Call2(_ *Thread, _, _ Value) Value {
	panic("too many arguments")
}
func (f *ParamSpec) Call3(_ *Thread, _, _, _ Value) Value {
	panic("too many arguments")
}
func (f *ParamSpec) Call4(_ *Thread, _, _, _, _ Value) Value {
	panic("too many arguments")
}

func (f *ParamSpec) fillin(t *Thread, i int) Value {
	if f.Flags[i]&DynParam != 0 {
		if x := t.dyn("_" + f.Names[i]); x != nil {
			return x
		}
	}
	if i < int(f.Nparams-f.Ndefaults) {
		panic("missing argument(s)")
	}
	return f.Values[i-int(f.Nparams-f.Ndefaults)]
}

func (f *ParamSpec) Show() string {
	return f.String()
}
