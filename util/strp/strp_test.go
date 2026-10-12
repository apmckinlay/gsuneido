// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package strp

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSuStrP(t *testing.T) {
	test := func(s string) {
		t.Helper()
		sp := From(s)
		assert.This(Len(sp)).Is(len(s))
		assert.This(String(sp)).Is(s)
		// bs := []byte(s)
		// sp = From(bs)
		// assert.This(Len(sp)).Is(len(s))
		// assert.This(String(sp)).Is(s)
	}
	test("")
	test("x")
	test("foo")
	test(strings.Repeat("x", 255))

	assert.This(func() { From(strings.Repeat("x", 256)) }).Panics("too large")
}

func TestFromShort(t *testing.T) {
	assert := assert.T(t)
	assert.That(From("") == &empty)
	assert.That(From("") == From(""))
	assert.This(Len(From(""))).Is(0)
	assert.That(From("x") == From("x"))
	// assert.That(From([]byte("x")) == From("x"))
	// assert.This(String(From("x"))).Is("x")
}

func TestConcat(t *testing.T) {
	a := assert.T(t)
	assert := a.This
	test := func(ss ...string) {
		t.Helper()
		total := 0
		for _, s := range ss {
			total += len(s)
		}
		assert(Len(Concat(ss...))).Is(total)
		assert(String(Concat(ss...))).Is(strings.Join(ss, ""))
	}
	test()
	test("")
	test("x")
	test("", "")
	test("", "x")
	test("x", "")
	test("x", "y")
	test("foo", "bar")
	test("a", "b", "c")
	test("", "a", "", "bcd", "")
	test("foo", "bar", "baz", "qux", "quux")
	test(strings.Repeat("x", 100), strings.Repeat("y", 55), strings.Repeat("z", 100))

	a.This(func() { Concat(strings.Repeat("x", 256), "") }).Panics("too large")
	a.This(func() { Concat("a", strings.Repeat("x", 256), "b") }).Panics("too large")
}
