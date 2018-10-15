package runtime

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func ExampleSuInt() {
	v := SuInt(123)
	fmt.Printf("%d %s\n", v.ToInt(), v.String())
	// Output: 123 123
}

func TestStrConvert(t *testing.T) {
	Assert(t).That(SuStr("123").ToStr(), Equals("123"))
}

func TestStringGet(t *testing.T) {
	var v Value = SuStr("hello")
	v = v.Get(SuInt(1))
	Assert(t).That(v, Equals(Value(SuStr("e"))))
}

func TestPanics(t *testing.T) {
	v := SuInt(123)
	Assert(t).That(func() { v.Get(v) }, Panics("number does not support get"))

	var ob SuObject
	Assert(t).That(func() { ob.ToInt() }, Panics("cannot convert object to integer"))
}

func TestCompare(t *testing.T) {
	vals := []Value{False, True, SuDnum{dnum.NegInf},
		SuInt(-1), SuInt(0), SuInt(+1), SuDnum{dnum.Inf},
		SuStr(""), SuStr("abc"), NewSuConcat().Add("foo"), SuStr("world")}
	for i := 1; i < len(vals); i++ {
		Assert(t).That(vals[i].Compare(vals[i]), Equals(0))
		Assert(t).That(vals[i-1].Compare(vals[i]), Equals(-1).Comment(vals[i-1], vals[i]))
		Assert(t).That(vals[i].Compare(vals[i-1]), Equals(+1))
	}
}
