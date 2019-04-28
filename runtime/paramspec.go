package runtime

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/str"
)

// ParamSpec describes the parameters of a function
// See also ArgSpec which describes the arguments of a function call
// It also serves as the basis for callables like SuFunc
type ParamSpec struct {
	// Nparams is the number of arguments required on the stack
	Nparams uint8

	// NDefaults is the number of default values for parameters
	// They are in the start of Values
	Ndefaults uint8

	// Offset is the location of the parameter names within Names
	// It is used for closure blocks, for normal functions it is 0
	// Offset is only required for parameters (or arguments),
	// it is not needed for local variables
	Offset uint8

	// Signature is used for fast matching of simple Argspec to ParamSpec
	Signature byte

	// Flags specifies "types" of params
	Flags []Flag

	// Names normally starts with the parameters, followed by local variables
	Names []string

	// Values contains any literals in the function
	// starting with parameter defaults
	Values []Value

	// Name for library records will be the record name;
	// for class methods, the method name;
	// if assigned to a local variable, the variable name.
	// It is primarily for debugging, i.e. call stack traces.
	// See also: Name(value) builtin function
	Name string

	CantConvert
}

// Flag is a bit set of parameter options
type Flag byte

const (
	AtParam Flag = 1 << iota
	DynParam
	DotParam
	PubParam
)

var ParamSpec0 = ParamSpec{Nparams: 0, Signature: ^Sig0}
var ParamSpec1 = ParamSpec{Nparams: 1, Signature: ^Sig1}
var ParamSpec2 = ParamSpec{Nparams: 2, Signature: ^Sig2}

var ParamSpecOptionalBlock = ParamSpec{Nparams: 1, Ndefaults: 1,
	Flags: []Flag{0}, Names: []string{"block"}, Values: []Value{False}}

// Value interface (except TypeName) --------------------------------

func (f *ParamSpec) String() string {
	var buf strings.Builder
	// easier to add "function" here and strip it in Params
	// than to implement String in a bunch of places to add it
	buf.WriteString("function(")
	sep := ""
	v := 0 // index into Values
	for i := 0; i < int(f.Nparams); i++ {
		buf.WriteString(sep)
		buf.WriteString(flagsToName(f.Names[i+int(f.Offset)], f.Flags[i]))
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

func (*ParamSpec) Get(*Thread, Value) Value {
	panic("function does not support get")
}

func (*ParamSpec) Put(*Thread, Value, Value) {
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
	// interface check and double dispatch
	// to work with anything that embeds ParamSpec
	if ep, ok := other.(eqps); ok {
		return ep.equalParamSpec(f)
	}
	return false
}

func (f *ParamSpec) equalParamSpec(ps *ParamSpec) bool {
	return f == ps
}

type eqps interface{ equalParamSpec(*ParamSpec) bool }

func (f *ParamSpec) Compare(other Value) int {
	if cmp := ints.Compare(OrdOther, Order(other)); cmp != 0 {
		return cmp
	}
	return 0 // ???
}

// ParamsMethods is initialized by the builtin package
var ParamsMethods Methods

func (*ParamSpec) Lookup(method string) Callable {
	return ParamsMethods[method]
}

func (f *ParamSpec) Params() string {
	return f.String()[8:] // skip "function"
}

func (f *ParamSpec) Show() string {
	return f.String()
}

// Named interface --------------------------------------------------

var _ Named = (*ParamSpec)(nil)

func (f *ParamSpec) GetName() string {
	return f.Name
}
