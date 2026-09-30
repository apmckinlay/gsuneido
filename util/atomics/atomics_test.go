// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package atomics

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestString(t *testing.T) {
	var s String
	assert.T(t).This(s.Load()).Is("")
	s.Store("hello")
	assert.T(t).This(s.Load()).Is("hello")
	s.Store("world")
	assert.T(t).This(s.Load()).Is("world")
}

func TestValue(t *testing.T) {
	var v Value[int]
	assert.T(t).This(v.Load()).Is(0)
	v.Store(42)
	assert.T(t).This(v.Load()).Is(42)
	old := v.Swap(20)
	assert.T(t).This(old).Is(42)
	assert.T(t).This(v.Load()).Is(20)
}

type testIntfc interface {
	Foo() string
}

type testImpl struct {
	s string
}

func (t testImpl) Foo() string {
	return t.s
}

type testImpl2 struct {
	s string
}

func (t testImpl2) Foo() string {
	return t.s
}

func TestIntfcStoreLoad(t *testing.T) {
	var a Intfc[testIntfc]
	assert.T(t).This(a.Load()).Is(nil)
	var val testIntfc = testImpl{s: "hello"}
	a.Store(val)
	got := a.Load()
	assert.T(t).This(got).Is(val)
	assert.T(t).This(got.Foo()).Is("hello")

	val = testImpl2{s: "world"}
	a.Store(val)
	got = a.Load()
	assert.T(t).This(got).Is(val)
	assert.T(t).This(got.Foo()).Is("world")
}
