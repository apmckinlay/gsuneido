package value

import "fmt"
import "testing"

func IntConvert() {
	v := IntValue{i: 123}
	fmt.Printf("%d %s\n", v.ToInt(), v.ToStr())
	// Output: 123 123
}

func StrConvert() {
	v := StrValue("123")
	fmt.Printf("%d %s\n", v.ToInt(), v.ToStr())
	// Output: 123 123
}

func TestStringGet(t *testing.T) {
	var v Value = StrValue("hello")
	v = v.Get(IntValue{i: 1})
	verify(t, v == Value(StrValue("e")), fmt.Sprintf("expected 'e' but got %v", v))
}

func TestPanics(t *testing.T) {
	v := IntValue{i: 123}
	verifyPanic(t, func() { v.Get(v) }, "doesn't support get")

	var ob Object
	verifyPanic(t, func() { ob.ToInt() }, "can't convert to integer")
}

func verify(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

func verifyPanic(t *testing.T, f func(), expected string) {
	defer func() {
		if e := recover(); e != nil {
			verify(t, e == expected,
				"expected panic of: "+expected+" but got: "+e.(string))
		}
	}()
	f()
	panic("expected exception")
}
