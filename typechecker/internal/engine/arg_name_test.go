// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// `f(1: x)` parses - the parser takes the name from any constant - so every
// place that reads an argument name has to cope with a name that is not a
// string. Same for a matcher or Type() argument that is not a type name.
func TestNumericArgumentName(t *testing.T) {
	a := assert.T(t)
	for _, src := range []string{
		`class { Go() { return this.Foo(1: 2) } Foo(a) { return a } }`,
		`class { Go(ob: object) { return ob.Map(1: 2) } }`,
		`class { Go(x) { Assert(x, 1: 2) return x } }`,
		`class { Go(x) { Assert(x, isType: 5) return x } }`,
		`class { Go(x) { if Type(x) is 5 { return 1 } return 0 } }`,
	} {
		cls, ok := safeParse(src, "T")
		a.Msg(src).That(ok)
		_, panicMsg := safeRun(cls)
		a.Msg(src).This(panicMsg).Is("")
	}
}

// a numeric member key names no method and no member
func TestNumericMemberKey(t *testing.T) {
	a := assert.T(t)
	cls := NewClassObject("T", ParseClass(`#(1: 2)`))
	a.This(len(cls.Methods)).Is(0)
	a.This(len(cls.Members)).Is(0)
}
