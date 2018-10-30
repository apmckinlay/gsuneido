package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuFuncString(t *testing.T) {
	sf := SuFunc{}
	sf.Flags = make([]Flag, 8)
	Assert(t).That(sf.Params(), Equals("()"))
	sf.Nparams = 3
	sf.Names = []string{"a", "b", "c"}
	Assert(t).That(sf.Params(), Equals("(a,b,c)"))
	sf.Names = []string{"a", "b", "c"}
	sf.Ndefaults = 1
	sf.Values = []Value{SuInt(123)}
	Assert(t).That(sf.Params(), Equals("(a,b,c=123)"))
}
