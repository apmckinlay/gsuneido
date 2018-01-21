package base

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func ExampleIntConvert() {
	v := SuInt(123)
	fmt.Printf("%d %s\n", v.ToInt(), v.ToStr())
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
