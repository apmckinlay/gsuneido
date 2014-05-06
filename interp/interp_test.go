package interp

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/value"
)

func TestInterp(t *testing.T) {
	code := []byte{PUSHINT, 3 << 1, PUSHINT, 5 << 1, ADD, RETURN}
	fn := &Function{Code: code}
	th := Thread{}
	result := th.Call(fn, SimpleArgSpecs[0])
	Assert(t).That(result, Equals(value.SuInt(8)))
}
