package interp

/*
ArgSpec specifies the arguments on the stack

- first byte is the number of un-named arguments
- additional bytes (if any) are one of:
-- index into the names array for a named argument
-- EACH for @arg
-- NONAME for unnamed args following named or each

The spec is normally embedded directly in the byte code
and sliced out of it without copying or processing.
Shortcut op codes
that specify only the number of unnamed arguments
use the predefined SimpleArgSpecs
since the code does not contain an actual spec.
*/
type ArgSpec struct {
	spec  []byte
	names []string // the names from the calling Function
}

const (
	EACH   = 254
	NONAME = 255
)

// ArgSpecs are predefined ArgSpec
// for small numbers of unnamed arguments
// with no each or named
var SimpleArgSpecs = [...]ArgSpec{
	{[]byte{0}, []string{}},
	{[]byte{1}, []string{}},
	{[]byte{2}, []string{}},
	{[]byte{3}, []string{}},
	{[]byte{4}, []string{}},
	{[]byte{5}, []string{}},
	{[]byte{6}, []string{}},
	{[]byte{7}, []string{}},
}

func (as *ArgSpec) Nargs() int {
	return as.N_unnamed() + len(as.spec) - 1
}

func (as *ArgSpec) N_unnamed() int {
	return int(as.spec[0])
}

// ArgName returns the name of the i'th argument
func (as *ArgSpec) ArgName(i int) string {
	nu := as.N_unnamed()
	if i < nu {
		return ""
	}
	ni := as.spec[i-nu+1]
	if ni >= EACH {
		return ""
	}
	return as.names[ni]
}
