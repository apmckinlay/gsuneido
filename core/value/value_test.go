package value

import "fmt"
import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func IntConvert() {
	v := IntVal(123)
	fmt.Printf("%d %s\n", v.ToInt(), v.ToStr())
	// Output: 123 123
}

func StrConvert() {
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
	VerifyPanic(t, func() { v.Get(v) }, "number does not support get")

	var ob Object
	VerifyPanic(t, func() { ob.ToInt() }, "cannot convert object to integer")
}

func VerifyPanic(t *testing.T, f func(), expected string) {
	defer func() {
		if e := recover(); e != nil {
			Assert(t).That(e, Equals(expected))
		}
	}()
	f()
	panic("expected exception")
}
