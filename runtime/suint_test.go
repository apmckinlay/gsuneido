package runtime

import (
	"math"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuInt(t *testing.T) {
	test := func(n int) {
		Assert(t).That(SuInt(n).ToInt(), Equals(int(n)))
	}
	test(0)
	test(12345)
	test(-12345)
	test(math.MinInt16)
	test(math.MaxInt16)

	xtest := func(n int) {
		Assert(t).That(func() { SuInt(n) }, Panics("index out of range"))
	}
	xtest(123456)
	xtest(-123456)
}
