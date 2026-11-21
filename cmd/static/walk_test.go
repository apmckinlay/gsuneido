// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func globalsOf(src string) []string {
	var globals []string
	walkRefs(parse("", src), func(typ, value string) {
		if typ == "global" {
			globals = append(globals, value)
		}
	})
	return globals
}

func TestWalkRefs(t *testing.T) {
	src := `function () {
		Foo(bar)
		Baz.Qux
		x = function () { Deep() }
	}`
	assert.T(t).This(globalsOf(src)).Is([]string{"Foo", "Baz", "Deep"})
}

func TestWalkRefsExprPos(t *testing.T) {
	assert.T(t).This(globalsOf(`function () {
		for x in Foo
			{ }
	}`)).Is([]string{"Foo"})
}

func TestWalkRefsClass(t *testing.T) {
	src := `Foo {
		Bar() { Baz }
		Qux() { x = function () { Deep() } }
		Member: 123
	}`
	assert.T(t).This(globalsOf(src)).Is([]string{"Baz", "Deep"})
}

func TestWalkRefsIgnoresNonFunction(t *testing.T) {
	assert.T(t).This(globalsOf(`123`)).Is(nil)
	assert.T(t).This(globalsOf(`#(Foo: Bar)`)).Is(nil)
}

func TestWalkRefsMembers(t *testing.T) {
	assert.T(t).This(globalsOf(`function () { a.b.c }`)).
		Is([]string{"Rule_b", "Rule_c"})
	assert.T(t).This(globalsOf(`function () { a.B }`)).Is(nil)
	assert.T(t).This(globalsOf(`function () { .x }`)).Is(nil)
	assert.T(t).This(globalsOf(`function () { this.x }`)).Is(nil)
	assert.T(t).This(globalsOf(`function () { super.x }`)).Is(nil)
}

func TestWalkStringConstants(t *testing.T) {
	var strings []string
	walkRefs(parse("", `function () {
		Foo("bar")
		x = function () { "nested" }
	}`), func(typ, value string) {
		if typ == "string" {
			strings = append(strings, value)
		}
	})
	assert.T(t).This(strings).Is([]string{"bar", "nested"})
}
