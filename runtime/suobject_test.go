package runtime

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestSuObject(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
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
	Assert(t).That(ob.Get(nil, Zero), Equals(iv))
	Assert(t).That(ob.Get(nil, One), Equals(sv))

	ob.Put(sv, iv)
	Assert(t).That(ob.String(), Equals("#(123, 'hello', hello: 123)"))
	ob.Put(iv, sv)
	Assert(t).That(ob.Size(), Equals(4))
}

func TestSuObjectString(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(k string, expected string) {
		ob := SuObject{}
		ob.Put(SuStr(k), SuInt(123))
		Assert(t).That(ob.String(), Equals(expected))
	}
	test("foo", "#(foo: 123)")
	test("123", "#('123': 123)")
	test("foo bar", "#('foo bar': 123)")
}

func TestSuObjectObjectAsKey(t *testing.T) {
	ob := SuObject{}
	ob.Put(&SuObject{}, SuInt(123))
	Assert(t).That(ob.Get(nil, &SuObject{}), Equals(SuInt(123)))
}

func TestSuObjectMigrate(t *testing.T) {
	ob := SuObject{}
	for i := 1; i < 5; i++ {
		ob.Put(SuInt(i), SuInt(i))
	}
	Assert(t).That(ob.NamedSize(), Equals(4))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Add(Zero)
	Assert(t).That(ob.NamedSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(5))
}

func TestSuObjectPut(t *testing.T) {
	ob := SuObject{}
	ob.Put(One, One) // put
	Assert(t).That(ob.NamedSize(), Equals(1))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Put(Zero, Zero) // add + migrate
	Assert(t).That(ob.NamedSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(2))
	ob.Put(Zero, SuInt(10)) // set
	ob.Put(One, SuInt(11)) // set
	Assert(t).That(ob.Get(nil, Zero), Equals(SuInt(10)))
	Assert(t).That(ob.Get(nil, One), Equals(SuInt(11)))
}

func TestSuObjectDelete(t *testing.T) {
	ob := SuObject{}
	ob.Delete(Zero)
	ob.Delete(SuStr("baz"))
	for i := 0; i < 5; i++ {
		ob.Add(SuInt(i))
	}
	ob.Put(SuStr("foo"), SuInt(8))
	ob.Put(SuStr("bar"), SuInt(9))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9, foo: 8)"))
	ob.Delete(SuStr("foo"))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9)"))
	ob.Delete(SuInt(2))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 3, 4, bar: 9)"))
	ob.Delete(Zero)
	Assert(t).That(ob.Show(), Equals("#(1, 3, 4, bar: 9)"))
	ob.Delete(SuInt(2))
	Assert(t).That(ob.Show(), Equals("#(1, 3, bar: 9)"))

	ob.Clear()
	Assert(t).That(ob.Show(), Equals("#()"))
	Assert(t).That(ob.Size(), Equals(0))
}

func TestSuObjectErase(t *testing.T) {
	ob := SuObject{}
	ob.Erase(Zero)
	ob.Erase(SuStr("baz"))
	for i := 0; i < 5; i++ {
		ob.Add(SuInt(i))
	}
	ob.Put(SuStr("foo"), SuInt(8))
	ob.Put(SuStr("bar"), SuInt(9))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9, foo: 8)"))
	ob.Erase(SuStr("foo"))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9)"))
	ob.Erase(SuInt(2))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 3: 3, 4: 4, bar: 9)"))
	ob.Erase(One)
	Assert(t).That(ob.Show(), Equals("#(0, 3: 3, 4: 4, bar: 9)"))
}

func TestSuObjectEquals(t *testing.T) {
	x := &SuObject{}
	y := &SuObject{}
	eq(t, x, y)
	x.Add(One)
	neq(t, x, y)
	y.Add(One)
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

func TestSuObjectPack(t *testing.T) {
	ob := &SuObject{}
	check := func() {
		t.Helper()
		Assert(t).That(Unpack(Pack(ob)), Equals(ob))
	}
	check()
	ob.Add(SuStr(1))
	check()
	ob.Add(SuInt(2))
	check()
	ob.Put(SuStr("a"), SuInt(3))
	check()
	ob.Put(SuStr("b"), SuInt(4))
	check()
}

func TestSuObjectPack2(t *testing.T) {
	ob := &SuObject{}
	ob.Add(One)
	ob.Put(SuStr("a"), SuInt(2))
	buf := Pack(ob)
	expected := []byte{6, 128, 0, 0, 1, 128, 0, 0, 4, 3, 129, 0, 1, 128,
		0, 0, 1, 128, 0, 0, 2, 4, 97, 128, 0, 0, 4, 3, 129, 0, 2}
	Assert(t).That(buf, Equals(string(expected)))
}
