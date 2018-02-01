package interp

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/interp/op"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestInterp(t *testing.T) {
	test := func(expected Value, code ...byte) {
		fn := &SuFunc{Code: code}
		th := NewThread()
		result := th.Call(fn, nil)
		Assert(t).That(result, Equals(SuInt(8)))
	}
	test(SuInt(8), op.INT, 3<<1, op.INT, 5<<1, op.ADD, op.RETURN)
}
