package language

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/interp"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuFuncCall(t *testing.T) {
	fn := compile.Constant("function () { 123 }")
	th := NewThread()
	result := fn.Call(th, SimpleArgSpecs[0])
	Assert(t).That(result, Equals(SuInt(123)))
	AddG("F", fn)
	fn = compile.Constant("function () { F() }")
	result = fn.Call(th, SimpleArgSpecs[0])
	Assert(t).That(result, Equals(SuInt(123)))
}
