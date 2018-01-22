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
	fn := compile.Constant("function (n) { n + n }")
	th := NewThread()
	th.Push(SuInt(123))
	result := th.Call(fn.(*SuFunc), ArgSpec{1, nil, nil})
	Assert(t).That(result, Equals(SuInt(246)))
	global.Add("F", fn)

	fn = compile.Constant("function () { F(321) }")
	result = th.Call(fn.(*SuFunc), ArgSpec{})
	Assert(t).That(result, Equals(SuInt(642)))
}
