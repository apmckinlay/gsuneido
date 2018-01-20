package interp

import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func TestSuFunc(t *testing.T) {
	sf := SuFunc{}
	Assert(t).That(sf.String(), Equals("function ()"))
	sf.Nparams = 3
	sf.Strings = []string{"a", "b", "c"}
	Assert(t).That(sf.String(), Equals("function (a, b, c)"))
	sf.Strings = []string{"a", "b", "=c"}
	sf.Values = []Value{SuInt(123)}
	Assert(t).That(sf.String(), Equals("function (a, b, c=123)"))
}
