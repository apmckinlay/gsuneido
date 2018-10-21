package language

import (
	"math"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuFuncCall(t *testing.T) {
	fn := compile.Constant("function (a, b) { a - b }").(*SuFunc)
	th := NewThread()
	th.Push(SuInt(100))
	th.Push(SuInt(1))
	result := th.Call(fn)
	Assert(t).That(result, Equals(SuInt(99)))
	AddGlobal("F", fn)

	fn = compile.Constant("function () { F(100, 1) }").(*SuFunc)
	result = th.Call(fn)
	Assert(t).That(result, Equals(SuInt(99)))

	fn = compile.Constant("function () { F(b: 1, a: 100) }").(*SuFunc)
	result = th.Call(fn)
	Assert(t).That(result, Equals(SuInt(99)))
}

func BenchmarkInt(b *testing.B) {
	fn := compile.Constant("function (n) { for (i = 0; i < n; ++i){} }")
	th := NewThread()
	m := 1
	n := b.N
	for n > math.MaxInt16 {
		n /= 2
		m *= 2
	}
	for i := 0; i < m; i++ {
		th.Push(SuInt(n))
		th.Call(fn.(*SuFunc))
	}
}
