package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/runtime/op"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestInterp(t *testing.T) {
	test := func(expected Value, code ...byte) {
		fn := &SuFunc{Code: code}
		th := NewThread()
		result := th.Call(fn)
		Assert(t).That(result, Equals(SuInt(8)))
	}
	test(SuInt(8), op.INT, 0, 3, op.INT, 0, 5, op.ADD, op.RETURN)
}

// compare to BenchmarkInterp in execute_test.go
func BenchmarkJit(b *testing.B) {
	th := &Thread{}
	for n := 0; n < b.N; n++ {
		th.Reset()
		result := jitfn(th)
		if !result.Equal(SuInt(4950)) {
			panic("wrong result")
		}
	}
}

var hundred = SuInt(100)

func jitfn(th *Thread) Value {
	th.sp += 2
	th.stack[0] = Zero // sum
	th.stack[1] = Zero // i
	for {
		th.stack[0] = Add(th.stack[0], th.stack[1]) // sum += i
		th.stack[1] = Add(th.stack[1], One)         // ++i
		if Lt(th.stack[1], hundred) != True {
			break
		}
	}
	return th.stack[0] // return sum
}
