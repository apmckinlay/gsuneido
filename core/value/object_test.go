package value

import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func TestBasic(t *testing.T) {
	ob := Object{}
	Assert(t).That(ob.String(), Equals("#()"))
	Assert(t).That(ob.Size(), Equals(0))
	iv := IntVal(123)
	ob.Add(iv)
	Assert(t).That(ob.Size(), Equals(1))
	Assert(t).That(ob.String(), Equals("#(123)"))
	sv := StrVal("hello")
	ob.Add(sv)
	Assert(t).That(ob.Size(), Equals(2))
	Assert(t).That(ob.Get(IntVal(0)), Equals(iv))
	Assert(t).That(ob.Get(IntVal(1)), Equals(sv))

	ob.Put(sv, iv)
	Assert(t).That(ob.String(), Equals("#(123, 'hello', hello: 123)"))
	ob.Put(iv, sv)
	Assert(t).That(ob.Size(), Equals(4))
}

func TestString(t *testing.T) {
	test := func(k string, expected string) {
		ob := Object{}
		ob.Put(StrVal(k), IntVal(123))
		Assert(t).That(ob.String(), Equals(expected))
	}
	test("foo", "#(foo: 123)")
	test("123", "#('123': 123)")
	test("foo bar", "#('foo bar': 123)")
}

func Test_isIdentifier(t *testing.T) {
	Assert(t).That(isIdentifier(""), Equals(false))
	Assert(t).That(isIdentifier("123"), Equals(false))
	Assert(t).That(isIdentifier("123bar"), Equals(false))
	Assert(t).That(isIdentifier("foo123"), Equals(true))
	Assert(t).That(isIdentifier("foo 123"), Equals(false))
	Assert(t).That(isIdentifier("_foo"), Equals(true))
	Assert(t).That(isIdentifier("Bar!"), Equals(true))
	Assert(t).That(isIdentifier("Bar?"), Equals(true))
	Assert(t).That(isIdentifier("Bar?x"), Equals(false))
}

func TestObjectAsKey(t *testing.T) {
	ob := Object{}
	ob.Put(&Object{}, IntVal(123))
	Assert(t).That(ob.Get(&Object{}), Equals(IntVal(123)))
}
