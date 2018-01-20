package interp

import "testing"
import . "github.com/apmckinlay/gsuneido/util/hamcrest"

func TestBasic(t *testing.T) {
	ob := SuObject{}
	Assert(t).That(ob.String(), Equals("#()"))
	Assert(t).That(ob.Size(), Equals(0))
	iv := SuInt(123)
	ob.Add(iv)
	Assert(t).That(ob.Size(), Equals(1))
	Assert(t).That(ob.String(), Equals("#(123)"))
	sv := SuStr("hello")
	ob.Add(sv)
	Assert(t).That(ob.Size(), Equals(2))
	Assert(t).That(ob.Get(SuInt(0)), Equals(iv))
	Assert(t).That(ob.Get(SuInt(1)), Equals(sv))

	ob.Put(sv, iv)
	Assert(t).That(ob.String(), Equals("#(123, 'hello', hello: 123)"))
	ob.Put(iv, sv)
	Assert(t).That(ob.Size(), Equals(4))
}

func TestString(t *testing.T) {
	test := func(k string, expected string) {
		ob := SuObject{}
		ob.Put(SuStr(k), SuInt(123))
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
	ob := SuObject{}
	ob.Put(&SuObject{}, SuInt(123))
	Assert(t).That(ob.Get(&SuObject{}), Equals(SuInt(123)))
}

func TestMigrate(t *testing.T) {
	ob := SuObject{}
	for i := 1; i < 5; i++ {
		ob.Put(SuInt(i), SuInt(i))
	}
	Assert(t).That(ob.HashSize(), Equals(4))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Add(SuInt(0))
	Assert(t).That(ob.HashSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(5))
}

func TestPut(t *testing.T) {
	ob := SuObject{}
	ob.Put(SuInt(1), SuInt(1)) // put
	Assert(t).That(ob.HashSize(), Equals(1))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Put(SuInt(0), SuInt(0)) // add + migrate
	Assert(t).That(ob.HashSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(2))
	ob.Put(SuInt(0), SuInt(10)) // set
	ob.Put(SuInt(1), SuInt(11)) // set
	Assert(t).That(ob.Get(SuInt(0)), Equals(SuInt(10)))
	Assert(t).That(ob.Get(SuInt(1)), Equals(SuInt(11)))
}

func TestEquals(t *testing.T) {
	x := &SuObject{}
	y := &SuObject{}
	eq(t, x, y)
	x.Add(SuInt(1))
	neq(t, x, y)
	y.Add(SuInt(1))
	eq(t, x, y)
	x.Put(SuInt(4), SuInt(6))
	neq(t, x, y)
	y.Put(SuInt(4), SuInt(7))
	neq(t, x, y)
	y.Put(SuInt(4), SuInt(6))
	eq(t, x, y)
	x.Put(SuInt(9), x) // recursive
	neq(t, x, y)
	y.Put(SuInt(9), y)
	eq(t, x, y)
}

func eq(t *testing.T, x Value, y Value) {
	Assert(t).True(x.Equals(y))
	Assert(t).True(y.Equals(x))
}

func neq(t *testing.T, x Value, y Value) {
	Assert(t).False(x.Equals(y))
	Assert(t).False(y.Equals(x))
}
