package interp

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/value"
)

func TestInterp(t *testing.T) {
	test := func(expected value.Value, code ...byte) {
		fn := &value.SuFunc{Code: code}
		th := NewThread()
		result := th.Call(fn, SimpleArgSpecs[0])
		Assert(t).That(result, Equals(value.SuInt(8)))
	}
	test(value.SuInt(8), INT, 3<<1, INT, 5<<1, ADD, RETURN)
}
