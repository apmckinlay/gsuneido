package runtime

import (
	"strings"

	"github.com/apmckinlay/gsuneido/lexer"
)

// ArgSpec describes the arguments on the stack for a call
// See also ParamSpec
type ArgSpec struct {
	// Nargs is the number of argument (unnamed + named)
	Nargs byte

	// Each is 1 for @args, 2 for @+1args, 0 otherwise
	Each byte

	// Signature is used for fast matching of simple Argspec to ParamSpec
	Signature byte

	// Spec has one entry per named argument, indexing into Names
	Spec []byte

	// Names is the argument names from the calling function
	Names []Value
}

// values for ArgSpec.Each
const (
	EACH  = 1
	EACH1 = 2
)

const (
	Sig0 byte = iota + 1
	Sig1
	Sig2
	Sig3
	Sig4
)

var ArgSpec0 = &ArgSpec{Nargs: 0, Signature: Sig0}
var ArgSpec1 = &ArgSpec{Nargs: 1, Signature: Sig1}
var ArgSpec2 = &ArgSpec{Nargs: 2, Signature: Sig2}
var ArgSpec3 = &ArgSpec{Nargs: 3, Signature: Sig3}
var ArgSpec4 = &ArgSpec{Nargs: 4, Signature: Sig4}
var ArgSpecEach = &ArgSpec{Nargs: 1, Each: EACH}
var ArgSpecEach1 = &ArgSpec{Nargs: 1, Each: EACH1}
var ArgSpecBlock = &ArgSpec{Nargs: 1,
	Spec: []byte{0}, Names: []Value{SuStr("block")}}

var StdArgSpecs = [...]*ArgSpec{
	ArgSpec0,
	ArgSpec1,
	ArgSpec2,
	ArgSpec3,
	ArgSpec4,
	ArgSpecEach,
	ArgSpecEach1,
	ArgSpecBlock,
}

const (
	AsEach = iota + 5
	AsEach1
	AsBlock
)

// Nargs returns the total number of arguments
func (as *ArgSpec) Unnamed() int {
	return int(as.Nargs) - len(as.Spec)
}

func (as *ArgSpec) Equal(a2 *ArgSpec) bool {
	if as.Nargs != a2.Nargs || as.Each != a2.Each || len(as.Spec) != len(a2.Spec) {
		return false
	}
	for i := range as.Spec {
		if !as.Names[as.Spec[i]].Equal(a2.Names[a2.Spec[i]]) {
			return false
		}
	}
	return true
}

func (as *ArgSpec) String() string {
	var buf strings.Builder
	sep := ""
	buf.WriteString("ArgSpec(")
	if as.Each >= EACH {
		buf.WriteString("@")
		if as.Each == EACH1 {
			buf.WriteString("+1")
		}
	} else {
		for i := 0; i < as.Unnamed(); i++ {
			buf.WriteString(sep)
			buf.WriteString("?")
			sep = ", "
		}
		for _, i := range as.Spec {
			buf.WriteString(sep)
			if s, ok := as.Names[i].(SuStr); ok && lexer.IsIdentifier(string(s)) {
				buf.WriteString(string(s))
			} else {
				buf.WriteString(as.Names[i].String())
			}
			buf.WriteString(":")
			sep = ", "
		}
	}
	buf.WriteString(")")
	return buf.String()
}
