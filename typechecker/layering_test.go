// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/typechecker/internal/engine"
	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSignatureLayering(t *testing.T) {
	a := assert.T(t)
	prev := engine.Annotations()
	t.Cleanup(func() {
		SetBuiltinSignatures(nil)
		RegisterSignatures(nil)
		engine.SetAnnotations(prev)
	})
	a.That(nil == SetBuiltinSignatures([]Reg{
		{Kind: "free", Name: "ZedF", Sig: "() :number"},
		{Receiver: "record", Name: "ZedM", Sig: "() :number"},
	}))

	// a registration overrides the base entry with the same name and
	// receiver - record folds to object, so an object entry replaces the
	// record one rather than becoming a second signature
	_, err := RegisterSignatures([]Reg{
		{Kind: "free", Name: "ZedF", Sig: "() :string"},
		{Receiver: "object", Name: "ZedM", Sig: "() :string"},
	})
	a.That(err == nil)
	set := engine.Annotations()
	a.This(len(set["ZedF"])).Is(1)
	a.That(typealgebra.Equal(set["ZedF"][0].Returns, typealgebra.TString))
	a.This(len(set["ZedM"])).Is(1)
	a.That(typealgebra.Equal(set["ZedM"][0].Returns, typealgebra.TString))

	// replace-not-append: clearing the Suneido layer brings the base back
	_, err = RegisterSignatures(nil)
	a.That(err == nil)
	set = engine.Annotations()
	a.That(typealgebra.Equal(set["ZedF"][0].Returns, typealgebra.TNumber))

	// a bad base entry reports which layer it came from and changes nothing
	err = SetBuiltinSignatures([]Reg{{Kind: "free", Name: "Bad", Sig: "("}})
	a.That(err != nil && strings.Contains(err.Error(), "builtin[0]"))
	a.That(typealgebra.Equal(
		engine.Annotations()["ZedF"][0].Returns, typealgebra.TNumber))
}
