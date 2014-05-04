package value

import (
	"fmt"
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func ExampleIntConvert() {
	v := IntVal(123)
	fmt.Printf("%d %s\n", v.ToInt(), v.ToStr())
	// Output: 123 123
}

func ExampleStrConvert() {
	v := StrVal("123")
	fmt.Printf("%d %s\n", v.ToInt(), v.ToStr())
	// Output: 123 123
}

func TestStringGet(t *testing.T) {
	var v Value = StrVal("hello")
	v = v.Get(IntVal(1))
	Assert(t).That(v, Equals(Value(StrVal("e"))))
}

func TestPanics(t *testing.T) {
	v := IntVal(123)
	Assert(t).That(func() { v.Get(v) }, Panics("number does not support get"))

	var ob Object
	Assert(t).That(func() { ob.ToInt() }, Panics("cannot convert object to integer"))
}
