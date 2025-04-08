// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/pack"
)

func TestSuObject(t *testing.T) {
	assert := assert.T(t).This
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	ob := SuObject{}
	assert(ob.String()).Is("#()")
	assert(ob.Size()).Is(0)
	iv := SuInt(123)
	ob.Add(iv)
	assert(ob.Size()).Is(1)
	assert(ob.String()).Is("#(123)")
	sv := SuStr("hello")
	ob.Add(sv)
	assert(ob.Size()).Is(2)
	assert(ob.Get(nil, Zero)).Is(iv)
	assert(ob.Get(nil, One)).Is(sv)

	ob.Set(sv, iv)
	assert(ob.String()).Is("#(123, 'hello', hello: 123)")
	ob.Set(iv, sv)
	assert(ob.Size()).Is(4)
}

func TestSuObjectString(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()

	ob := SuObject{}
	assert.T(t).This(ob.String()).Is("#()")
	ob.Add(Zero)
	assert.T(t).This(ob.String()).Is("#(0)")
	ob.Add(One)
	assert.T(t).This(ob.String()).Is("#(0, 1)")

	ob = SuObject{}
	ob.Set(SuInt(123), Zero)
	assert.T(t).This(ob.String()).Is("#(123: 0)")
	ob.Set(SuInt(456), SuStr("abc"))
	assert.T(t).This(ob.Show()).Is("#(123: 0, 456: 'abc')")

	ob = SuObject{}
	ob.Set(EmptyStr, False)
	assert.T(t).This(ob.String()).Is("#('': false)")
	ob.Set(SuStr("a"), True)
	assert.T(t).This(ob.Show()).Is("#('': false, a:)")
	ob.Add(True)
	assert.T(t).This(ob.Show()).Is("#(true, '': false, a:)")

	test := func(k string, expected string) {
		t.Helper()
		ob := SuObject{}
		ob.Set(SuStr(k), SuInt(123))
		assert.T(t).This(ob.String()).Is(expected)
	}
	test("foo", "#(foo: 123)")
	test("123", "#('123': 123)")
	test("foo bar", "#('foo bar': 123)")
}

func TestSuObjectObjectAsKey(t *testing.T) {
	ob := SuObject{}
	ob.Set(&SuObject{}, SuInt(123))
	assert.T(t).This(ob.Get(nil, &SuObject{})).Is(SuInt(123))
}

func TestSuObjectMigrate(t *testing.T) {
	assert := assert.T(t).This
	ob := SuObject{}
	for i := 1; i < 5; i++ {
		ob.Set(SuInt(i), SuInt(i))
	}
	assert(ob.NamedSize()).Is(4)
	assert(ob.ListSize()).Is(0)
	ob.Add(Zero)
	assert(ob.NamedSize()).Is(0)
	assert(ob.ListSize()).Is(5)
}

func TestSuObjectPut(t *testing.T) {
	assert := assert.T(t).This
	ob := SuObject{}
	ob.Set(One, One) // put
	assert(ob.NamedSize()).Is(1)
	assert(ob.ListSize()).Is(0)
	ob.Set(Zero, Zero) // add + migrate
	assert(ob.NamedSize()).Is(0)
	assert(ob.ListSize()).Is(2)
	ob.Set(Zero, SuInt(10)) // set
	ob.Set(One, SuInt(11))  // set
	assert(ob.Get(nil, Zero)).Is(SuInt(10))
	assert(ob.Get(nil, One)).Is(SuInt(11))
}

func TestSuObjectDelete(t *testing.T) {
	assert := assert.T(t).This
	ob := SuObject{}
	ob.Delete(nil, Zero)
	ob.Delete(nil, SuStr("baz"))
	for i := range 5 {
		ob.Add(SuInt(i))
	}
	ob.Set(SuStr("foo"), SuInt(8))
	ob.Set(SuStr("bar"), SuInt(9))
	assert(ob.Show()).Is("#(0, 1, 2, 3, 4, bar: 9, foo: 8)")
	ob.Delete(nil, SuStr("foo"))
	assert(ob.Show()).Is("#(0, 1, 2, 3, 4, bar: 9)")
	ob.Delete(nil, SuInt(2))
	assert(ob.Show()).Is("#(0, 1, 3, 4, bar: 9)")
	ob.Delete(nil, Zero)
	assert(ob.Show()).Is("#(1, 3, 4, bar: 9)")
	ob.Delete(nil, SuInt(2))
	assert(ob.Show()).Is("#(1, 3, bar: 9)")

	ob.DeleteAll()
	assert(ob.Show()).Is("#()")
	assert(ob.Size()).Is(0)
}

func TestSuObjectErase(t *testing.T) {
	assert := assert.T(t).This
	ob := SuObject{}
	ob.Erase(nil, Zero)
	ob.Erase(nil, SuStr("baz"))
	for i := range 5 {
		ob.Add(SuInt(i))
	}
	ob.Set(SuInt(88), SuInt(8))
	ob.Set(SuInt(99), SuInt(9))
	assert(ob.Show()).Is("#(0, 1, 2, 3, 4, 88: 8, 99: 9)")
	ob.Erase(nil, SuInt(88))
	assert(ob.Show()).Is("#(0, 1, 2, 3, 4, 99: 9)")
	ob.Erase(nil, SuInt(2))
	assert(ob.Show()).Is("#(0, 1, 3: 3, 4: 4, 99: 9)")
	ob.Erase(nil, One)
	assert(ob.Show()).Is("#(0, 3: 3, 4: 4, 99: 9)")
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
	assert.T(t).True(x.Equal(y))
	assert.T(t).True(y.Equal(x))
}

func neq(t *testing.T, x *SuObject, y *SuObject) {
	assert.T(t).False(x.Equal(y))
	assert.T(t).False(y.Equal(x))
}

func TestSuObjectSlice(t *testing.T) {
	assert := assert.T(t)
	ob := SuObject{}
	ob.Add(SuInt(12))
	ob.Add(SuInt(34))
	ob.Add(SuInt(56))
	ob.Set(SuStr("a"), SuInt(123))
	assert.This(ob.String()).Is("#(12, 34, 56, a: 123)")
	ob2 := ob.Slice(0)
	assert.True(ob.Equal(ob2))
	ob2 = ob.Slice(1)
	assert.This(ob2.String()).Is("#(34, 56, a: 123)")
	ob2 = ob.Slice(10)
	assert.This(ob2.String()).Is("#(a: 123)")
}

func TestSuObjectPackValue(t *testing.T) {
	test := func(v1 Value) {
		t.Helper()
		pk := newPacking(50)
		packValue(v1, pk)
		s := pk.String()
		dec := pack.MakeDecoder(s)
		v2 := unpackValue(&dec, false)
		assert.T(t).This(v2).Is(v1)
	}
	test(SuInt(123))
	test(SuStr("hello"))
}

func TestSuObjectPack(t *testing.T) {
	ob := &SuObject{}
	check := func() {
		t.Helper()
		s := Pack(ob)
		assert.T(t).This(Unpack(s)).Is(ob)
	}
	check()
	ob.Add(SuInt(1))
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
	expected := []byte{6, 1, 3, PackPlus, 129, 10, 1, 2, PackString, 97, 3,
		PackPlus, 129, 20}
	assert.T(t).This([]byte(buf)).Is(expected)
}

func TestSuObjectPack3(t *testing.T) {
	ob := &SuObject{}
	ob.Add(One)
	pk1 := newPacking(0)
	size := ob.PackSize(pk1)
	ob.DeleteAll()
	pk2 := newPacking(size)
	ob.Pack(pk2)
	assert.T(t).This(pk1.hash).Isnt(pk2.hash)
	assert.T(t).This(size).Isnt(pk2.Len())

	size = ob.PackSize(pk1)
	ob.Add(One)
	pk2 = newPacking(size)
	assert.T(t).This(func() { ob.Pack(pk2) }).Panics("out of range")
}

func TestSuObjectCompare(t *testing.T) {
	x := &SuObject{}
	x.Add(Zero)
	y := &SuObject{}
	y.Add(One)
	for range 2 {
		assert.T(t).This(x.Compare(y)).Is(-1)
		assert.T(t).This(y.Compare(x)).Is(1)
		x.SetConcurrent()
		y.SetConcurrent()
	}
}

func TestSuObjectCopyOnWrite(t *testing.T) {
	x := &SuObject{}
	x.Add(IntVal(123))
	x.Set(SuStr("abc"), IntVal(456))
	assert.This(x.String()).Is("#(123, abc: 456)")
	y := x.Copy().(*SuObject)
	assert.This(y.String()).Is("#(123, abc: 456)")
	assert.True(slc.Same(x.list, y.list))

	x.Add(Zero)
	assert.This(x.String()).Is("#(123, 0, abc: 456)")
	assert.This(y.String()).Is("#(123, abc: 456)")
	assert.False(slc.Same(x.list, y.list))
}

func BenchmarkGetPut(b *testing.B) {
	x := &SuObject{}
	for _, m := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		x.Set(SuStr(m), SuInt(123))
	}
	for b.Loop() {
		x.GetPut(nil, SuStr("a"), One, F, false)
	}
}

var F = OpAdd
