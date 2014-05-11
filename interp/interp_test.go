package interp

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/value"
)

func TestInterp(t *testing.T) {
	code := []byte{INT, 3 << 1, INT, 5 << 1, ADD, RETURN}
	fn := &value.SuFunc{Code: code}
	th := Thread{}
	result := th.Call(fn, SimpleArgSpecs[0])
	Assert(t).That(result, Equals(value.SuInt(8)))
}
