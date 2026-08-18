// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package annotations

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func findSig(set Set, name string, recv typealgebra.DynType) (Signature, bool) {
	for _, s := range set[name] {
		if receiverKey(s.Receiver) == receiverKey(recv) {
			return s, true
		}
	}
	return Signature{}, false
}

func TestRegistrationKinds(t *testing.T) {
	a := assert.T(t)
	set, err := LoadRegistrations([]Registration{
		{Receiver: "record", Name: "Zed", Sig: "() :boolean"},
		{Kind: "free", Name: "ZedFn", Sig: "(x) :string"},
		{Kind: "static", Class: "ZedClass", Name: "Go", Sig: "() :number"},
	})
	a.That(err == nil)
	s, ok := findSig(set, "Zed", typealgebra.TObject) // record folds to object
	a.That(ok)
	a.This(s.Returns).Is(typealgebra.TBoolean)
	a.That(set["ZedFn"] != nil)
	a.That(set["ZedClass.Go"] != nil)
}

func TestRegistrationDefaultsUnknownReturn(t *testing.T) {
	a := assert.T(t)
	set, err := LoadRegistrations([]Registration{
		{Receiver: "string", Name: "Zed", Sig: "(x)"},
	})
	a.That(err == nil)
	s, ok := findSig(set, "Zed", typealgebra.TString)
	a.That(ok)
	a.This(s.Returns).Is(typealgebra.TUnknown)
}

func TestRegistrationErrors(t *testing.T) {
	a := assert.T(t)
	bad := []struct {
		reg  Registration
		want string
	}{
		{Registration{Receiver: "stirng", Name: "X", Sig: "() :string"}, "unknown receiver"},
		{Registration{Receiver: "string", Sig: "() :string"}, "missing name"},
		{Registration{Receiver: "string", Name: "X", Sig: "(((("}, "bad sig"},
		{Registration{Kind: "static", Name: "X", Sig: "() :string"}, "missing class"},
		{Registration{Kind: "free", Receiver: "string", Name: "X", Sig: "() :string"},
			"must not have receiver"},
		{Registration{Kind: "wat", Name: "X", Sig: "() :string"}, "unknown kind"},
		{Registration{Receiver: "string", Name: "X", Class: "C", Sig: "() :string"},
			"must not have class"},
	}
	for _, tc := range bad {
		_, err := LoadRegistrations([]Registration{tc.reg})
		a.That(err != nil)
		a.That(strings.Contains(err.Error(), tc.want))
	}
}

func TestRegistrationDuplicateErrors(t *testing.T) {
	a := assert.T(t)
	_, err := LoadRegistrations([]Registration{
		{Receiver: "string", Name: "Zed", Sig: "() :string"},
		{Receiver: "string", Name: "Zed", Sig: "() :number"},
	})
	a.That(err != nil)
	a.That(strings.Contains(err.Error(), "duplicate registration"))
	// record folds to object, so an object/record pair is a duplicate too
	_, err = LoadRegistrations([]Registration{
		{Receiver: "object", Name: "Zed", Sig: "() :string"},
		{Receiver: "record", Name: "Zed", Sig: "() :number"},
	})
	a.That(err != nil)
	// same name on a genuinely different receiver is fine
	_, err = LoadRegistrations([]Registration{
		{Receiver: "string", Name: "Zed", Sig: "() :string"},
		{Receiver: "object", Name: "Zed", Sig: "() :number"},
	})
	a.That(err == nil)
}
