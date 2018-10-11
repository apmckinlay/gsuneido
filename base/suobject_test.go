package base

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/dnum"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuObject(t *testing.T) {
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

func TestSuObjectString(t *testing.T) {
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

func TestSuObjectObjectAsKey(t *testing.T) {
	ob := SuObject{}
	ob.Put(&SuObject{}, SuInt(123))
	Assert(t).That(ob.Get(&SuObject{}), Equals(SuInt(123)))
}

func TestSuObjectMigrate(t *testing.T) {
	ob := SuObject{}
	for i := 1; i < 5; i++ {
		ob.Put(SuInt(i), SuInt(i))
	}
	Assert(t).That(ob.NamedSize(), Equals(4))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Add(SuInt(0))
	Assert(t).That(ob.NamedSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(5))
}

func TestSuObjectPut(t *testing.T) {
	ob := SuObject{}
	ob.Put(SuInt(1), SuInt(1)) // put
	Assert(t).That(ob.NamedSize(), Equals(1))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Put(SuInt(0), SuInt(0)) // add + migrate
	Assert(t).That(ob.NamedSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(2))
	ob.Put(SuInt(0), SuInt(10)) // set
	ob.Put(SuInt(1), SuInt(11)) // set
	Assert(t).That(ob.Get(SuInt(0)), Equals(SuInt(10)))
	Assert(t).That(ob.Get(SuInt(1)), Equals(SuInt(11)))
}

func TestSuObjectEquals(t *testing.T) {
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

	a := &SuObject{}
	a.Put(SuStr("a"), SuStr("aa"))
	b := &SuObject{}
	b.Put(SuStr("x"), SuStr("aa"))
	neq(t, a, b)
}

func eq(t *testing.T, x *SuObject, y *SuObject) {
	Assert(t).True(x.Equal(y))
	Assert(t).True(y.Equal(x))
}

func neq(t *testing.T, x *SuObject, y *SuObject) {
	Assert(t).False(x.Equal(y))
	Assert(t).False(y.Equal(x))
}

func TestSuObjectSlice(t *testing.T) {
	ob := SuObject{}
	ob.Add(SuInt(12))
	ob.Add(SuInt(34))
	ob.Add(SuInt(56))
	ob.Put(SuStr("a"), SuInt(123))
	Assert(t).That(ob.String(), Equals("#(12, 34, 56, a: 123)"))
	ob2 := ob.Slice(0)
	Assert(t).True(ob.Equal(ob2))
	ob2 = ob.Slice(1)
	Assert(t).That(ob2.String(), Equals("#(34, 56, a: 123)"))
	ob2 = ob.Slice(10)
	Assert(t).That(ob2.String(), Equals("#(a: 123)"))
}

func TestSuObjectIndex(t *testing.T) {
	Assert(t).That(index(SuInt(123)), Equals(123))
	Assert(t).That(index(SuDnum{dnum.FromInt(123)}), Equals(123))
}
