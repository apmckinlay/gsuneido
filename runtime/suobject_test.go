// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/pack"
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

	ob.Set(sv, iv)
	Assert(t).That(ob.String(), Equals("#(123, 'hello', hello: 123)"))
	ob.Set(iv, sv)
	Assert(t).That(ob.Size(), Equals(4))
}

func TestSuObjectString(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(k string, expected string) {
		ob := SuObject{}
		ob.Set(SuStr(k), SuInt(123))
		Assert(t).That(ob.String(), Equals(expected))
	}
	test("foo", "#(foo: 123)")
	test("123", "#('123': 123)")
	test("foo bar", "#('foo bar': 123)")
}

func TestSuObjectObjectAsKey(t *testing.T) {
	ob := SuObject{}
	ob.Set(&SuObject{}, SuInt(123))
	Assert(t).That(ob.Get(nil, &SuObject{}), Equals(SuInt(123)))
}

func TestSuObjectMigrate(t *testing.T) {
	ob := SuObject{}
	for i := 1; i < 5; i++ {
		ob.Set(SuInt(i), SuInt(i))
	}
	Assert(t).That(ob.NamedSize(), Equals(4))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Add(Zero)
	Assert(t).That(ob.NamedSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(5))
}

func TestSuObjectPut(t *testing.T) {
	ob := SuObject{}
	ob.Set(One, One) // put
	Assert(t).That(ob.NamedSize(), Equals(1))
	Assert(t).That(ob.ListSize(), Equals(0))
	ob.Set(Zero, Zero) // add + migrate
	Assert(t).That(ob.NamedSize(), Equals(0))
	Assert(t).That(ob.ListSize(), Equals(2))
	ob.Set(Zero, SuInt(10)) // set
	ob.Set(One, SuInt(11))  // set
	Assert(t).That(ob.Get(nil, Zero), Equals(SuInt(10)))
	Assert(t).That(ob.Get(nil, One), Equals(SuInt(11)))
}

func TestSuObjectDelete(t *testing.T) {
	ob := SuObject{}
	ob.Delete(nil, Zero)
	ob.Delete(nil, SuStr("baz"))
	for i := 0; i < 5; i++ {
		ob.Add(SuInt(i))
	}
	ob.Set(SuStr("foo"), SuInt(8))
	ob.Set(SuStr("bar"), SuInt(9))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9, foo: 8)"))
	ob.Delete(nil, SuStr("foo"))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9)"))
	ob.Delete(nil, SuInt(2))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 3, 4, bar: 9)"))
	ob.Delete(nil, Zero)
	Assert(t).That(ob.Show(), Equals("#(1, 3, 4, bar: 9)"))
	ob.Delete(nil, SuInt(2))
	Assert(t).That(ob.Show(), Equals("#(1, 3, bar: 9)"))

	ob.DeleteAll()
	Assert(t).That(ob.Show(), Equals("#()"))
	Assert(t).That(ob.Size(), Equals(0))
}

func TestSuObjectErase(t *testing.T) {
	ob := SuObject{}
	ob.Erase(nil, Zero)
	ob.Erase(nil, SuStr("baz"))
	for i := 0; i < 5; i++ {
		ob.Add(SuInt(i))
	}
	ob.Set(SuStr("foo"), SuInt(8))
	ob.Set(SuStr("bar"), SuInt(9))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9, foo: 8)"))
	ob.Erase(nil, SuStr("foo"))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 2, 3, 4, bar: 9)"))
	ob.Erase(nil, SuInt(2))
	Assert(t).That(ob.Show(), Equals("#(0, 1, 3: 3, 4: 4, bar: 9)"))
	ob.Erase(nil, One)
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
	x.Set(SuInt(4), SuInt(6))
	neq(t, x, y)
	y.Set(SuInt(4), SuInt(7))
	neq(t, x, y)
	y.Set(SuInt(4), SuInt(6))
	eq(t, x, y)
	x.Set(SuInt(9), x) // recursive
	neq(t, x, y)
	y.Set(SuInt(9), y)
	eq(t, x, y)

	a := &SuObject{}
	a.Set(SuStr("a"), SuStr("aa"))
	b := &SuObject{}
	b.Set(SuStr("x"), SuStr("aa"))
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
	ob.Set(SuStr("a"), SuInt(123))
	Assert(t).That(ob.String(), Equals("#(12, 34, 56, a: 123)"))
	ob2 := ob.Slice(0)
	Assert(t).True(ob.Equal(ob2))
	ob2 = ob.Slice(1)
	Assert(t).That(ob2.String(), Equals("#(34, 56, a: 123)"))
	ob2 = ob.Slice(10)
	Assert(t).That(ob2.String(), Equals("#(a: 123)"))
}

func TestSuObjectPackValue(t *testing.T) {
	test := func(v1 Value) {
		enc := pack.NewEncoder(50)
		packValue(v1, 0, enc)
		s := enc.String()
		dec := pack.NewDecoder(s)
		v2 := unpackValue(dec)
		Assert(t).That(v2, Equals(v1))
	}
	test(SuInt(123))
	test(SuStr("hello"))
}

func TestSuObjectPack(t *testing.T) {
	ob := &SuObject{}
	check := func() {
		t.Helper()
		s := Pack(ob)
		Assert(t).That(Unpack(s), Equals(ob))
	}
	check()
	ob.Add(SuStr(1))
	check()
	ob.Add(SuInt(2))
	check()
	ob.Set(SuStr("a"), SuInt(3))
	check()
	ob.Set(SuStr("b"), SuInt(4))
	check()
	ob.Add(SuStr(strings.Repeat("helloworld", 100)))
}

func TestSuObjectPack2(t *testing.T) {
	ob := &SuObject{}
	ob.Add(One)
	ob.Set(SuStr("a"), SuInt(2))
	buf := Pack(ob)
	expected := []byte{6, 1, 3, 3, 129, 10, 1, 2, 4, 97, 3, 3, 129, 20}
	Assert(t).That([]byte(buf), Equals(expected))
}

func TestSuObjectCompare(t *testing.T) {
	x := &SuObject{}
	x.Add(Zero)
	y := &SuObject{}
	y.Add(One)
	for i := 0; i < 2; i++ {
		Assert(t).That(x.Compare(y), Equals(-1))
		Assert(t).That(y.Compare(x), Equals(1))
		x.SetConcurrent()
		y.SetConcurrent()
	}
}
