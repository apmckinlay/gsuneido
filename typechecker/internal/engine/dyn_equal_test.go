// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// Union holds a slice, so `==` on two Union values panics. Two overloads of one
// name with the same union return - `Find :number|false` on both String and
// Object - reach allReturnsAgree with a Union on each side.
func TestUnionReturnsAgree(t *testing.T) {
	a := assert.T(t)
	numOrFalse, err := ParseTypeAnnotation("number|false")
	a.That(err == nil)
	dateOrFalse, err := ParseTypeAnnotation("date|false")
	a.That(err == nil)

	agree := scanByName([]Signature{
		{Receiver: TString, Returns: numOrFalse},
		{Receiver: TObject, Returns: numOrFalse},
	}, "Find")
	a.This(agree.Kind).Is(dispatchGuessAgree)
	a.That(dynEqual(agree.Returns, numOrFalse))

	disagree := scanByName([]Signature{
		{Receiver: TString, Returns: numOrFalse},
		{Receiver: TObject, Returns: dateOrFalse},
	}, "Find")
	a.This(disagree.Kind).Is(dispatchGuessDisagree)
}

// combineAsserts falls through to typeFits when dynEqual says the two facts
// differ, so two Asserts with different union shapes land there with a Union
// on each side too.
func TestAssertUnionShapes(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		New(.x) { Assert(Number?(.x) or String?(.x)) }
		Bar() { Assert(Number?(.x) or Date?(.x)); return .x }
	}`, "T")
	a.That(hasDiag(env, "Bar", SeverityError, "conflicting Assert ground truths"))
}
