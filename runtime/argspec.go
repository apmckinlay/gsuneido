package runtime

import "strings"

/*
ArgSpec specifies the arguments on the stack for a call.
The spec is normally embedded directly in the byte code
and sliced out of it without copying or processing.
*/
type ArgSpec struct {
	// Unnamed is the number of unnamed arguments, or the special values EACH or EACH1
	Unnamed byte
	// Spec has one entry per named argument, indexing into Names
	Spec []byte
	// Names is the locals names from the calling function
	Names []string
}

const (
	EACH  = 254
	EACH1 = 255
)

var ArgSpec0 = &ArgSpec{Unnamed: 0}
var ArgSpec1 = &ArgSpec{Unnamed: 1}
var ArgSpec2 = &ArgSpec{Unnamed: 2}
var ArgSpec3 = &ArgSpec{Unnamed: 3}
var ArgSpec4 = &ArgSpec{Unnamed: 4}

// Nargs returns the total number of arguments
func (as ArgSpec) Nargs() int {
	if as.Unnamed >= EACH {
		return 1
	}
	return int(as.Unnamed) + len(as.Spec)
}

// ArgName returns the name of the i'th argument
func (as ArgSpec) ArgName(i int) string {
	nu := int(as.Unnamed)
	if i < nu {
		return ""
	}
	ni := as.Spec[i-nu]
	if ni >= EACH1 {
		return ""
	}
	return as.Names[ni]
}

func (as ArgSpec) String() string {
	var buf strings.Builder
	sep := ""
	buf.WriteString("ArgSpec(")
	if as.Unnamed >= EACH {
		buf.WriteString("@")
		if as.Unnamed == EACH1 {
			buf.WriteString("+1")
		}
	} else {
		for i := byte(0); i < as.Unnamed; i++ {
			buf.WriteString(sep)
			buf.WriteString("?")
			sep = ", "
		}
		for _, i := range as.Spec {
			buf.WriteString(sep)
			buf.WriteString(as.Names[i])
			buf.WriteString(":")
			sep = ", "
		}
	}
	buf.WriteString(")")
	return buf.String()
}
