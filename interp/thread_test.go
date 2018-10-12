package interp

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestThread_stack(t *testing.T) {
	th := &Thread{}
	setStack := func(nums ...int) {
		th.sp = 0
		for _, n := range nums {
			th.Push(SuInt(n))
		}
	}
	ckStack := func(vals ...int) {
		t.Helper()
		Assert(t).That(fmt.Sprint(th.stack[:th.sp]), Equals(fmt.Sprint(vals)))
	}

	setStack(11, 22)
	ckStack(11, 22)
	th.Dup2()
	ckStack(11, 22, 11, 22)

	setStack(11, 22, 33, 44)
	th.Dupx2()
	ckStack(11, 44, 22, 33, 44)
}
