// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/util/str"
)

type sufunction struct{}

// ParamSpec describes the parameters of a function
// See also ArgSpec which describes the arguments of a function call
// It also serves as the basis for callables like SuFunc
type ParamSpec struct {
	ValueBase[sufunction]

	// Lib is the library that the function came from
	Lib string

	// Name for library records will be the record name;
	// for class methods, the method name;
	// if assigned to a local variable, the variable name.
	// It is primarily for debugging, i.e. call stack traces.
	// See also: Name(value) builtin function
	Name string

	// Values contains any literals in the function
	// starting with parameter defaults
	Values []Value

	// Flags specifies "types" of params
	Flags []Flag

	// Names normally starts with the parameters, followed by local variables
	Names []string

	// Nparams is the number of arguments required on the stack
	Nparams uint8

	// NDefaults is the number of default values for parameters
	// They are in the start of Values
	Ndefaults uint8

	// Offset is the location of the parameter names within Names
	// and the arguments within Locals
	// It is used for closure blocks, for normal functions it is 0
	// Offset is only required for parameters (or arguments),
	// it is not needed for local variables
	Offset uint8

	// Signature is used for fast matching of simple Argspec to ParamSpec
	Signature byte
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
var ParamSpec1 = ParamSpec{Nparams: 1, Signature: ^Sig1, Flags: []Flag{0}}

var ParamSpecAt = ParamSpec{Nparams: 1, Flags: []Flag{AtParam}}

var ParamSpecOptionalBlock = ParamSpec{Nparams: 1, Ndefaults: 1,
	Flags: []Flag{0}, Names: []string{"block"}, Values: []Value{False}}

// Value interface (except TypeName) --------------------------------

func (ps *ParamSpec) String() string {
	var buf strings.Builder
	// easier to add "function" here and strip it in Params
	// than to implement String in a bunch of places to add it
	buf.WriteString("function(")
	sep := ""
	v := 0 // index into Values
	for i := range int(ps.Nparams) {
		buf.WriteString(sep)
		buf.WriteString(flagsToName(ps.ParamName(i), ps.Flags[i]))
		if i >= int(ps.Nparams-ps.Ndefaults) {
			buf.WriteString("=")
			buf.WriteString(fmt.Sprint(ps.Values[v]))
			v++
		}
		sep = ","
	}
	buf.WriteString(")")
	return buf.String()
}

func (ps *ParamSpec) ParamName(i int) string {
	return ps.Names[i+int(ps.Offset)]
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

func (ps *ParamSpec) Equal(other any) bool {
	// interface check and double dispatch
	// to work with anything that embeds ParamSpec
	ep, ok := other.(eqps)
	return ok && ep.equalParamSpec(ps)
}

func (ps *ParamSpec) equalParamSpec(ps2 *ParamSpec) bool {
	return ps == ps2
}

type eqps interface{ equalParamSpec(*ParamSpec) bool }

func (ps *ParamSpec) Compare(other Value) int {
	if cmp := cmp.Compare(OrdOther, Order(other)); cmp != 0 {
		return cmp
	}
	return 0 // ???
}

func (ParamSpec) SetConcurrent() {
	// immutable so ok
}

// ParamsMethods is initialized by the builtin package
var ParamsMethods Methods

func (*ParamSpec) Lookup(_ *Thread, method string) Value {
	return ParamsMethods[method]
}

func (ps *ParamSpec) Params() string {
	return ps.String()[8:] // skip "function"
}

func (ps *ParamSpec) Show() string {
	return ps.String()
}

// Named interface --------------------------------------------------

var _ Named = (*ParamSpec)(nil)

func (ps *ParamSpec) GetName() string {
	return ps.Name
}
