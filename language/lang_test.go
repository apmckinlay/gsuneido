package language

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/interp/global"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuFuncCall(t *testing.T) {
	fn := compile.Constant("function (a, b) { a - b }")
	th := NewThread()
	th.Push(SuInt(100))
	th.Push(SuInt(1))
	result := th.Call(fn.(*SuFunc), ArgSpec{Unnamed: 2})
	Assert(t).That(result, Equals(SuInt(99)))
	global.Add("F", fn)

	fn = compile.Constant("function () { F(100, 1) }")
	result = th.Call(fn.(*SuFunc), ArgSpec{})
	Assert(t).That(result, Equals(SuInt(99)))

	// fn = compile.Constant("function () { F(b: 1, a: 100) }")
	// result = th.Call(fn.(*SuFunc), ArgSpec{})
	// Assert(t).That(result, Equals(SuInt(99)))
}
