package interp

import (
	"testing"

	"github.com/apmckinlay/gsuneido/core/value"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestInterp(t *testing.T) {
	code := []byte{PUSHINT, 3, PUSHINT, 5, ADD, RETURN}
	fn := Function{code: code}
	th := Thread{}
	result := th.Call(fn, SimpleArgSpecs[0])
	Assert(t).That(result, Equals(value.IntVal(8)))
}
