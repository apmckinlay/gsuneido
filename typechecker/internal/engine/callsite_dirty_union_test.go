// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// stdlib:SelectFields.ScanFields regression: src is untyped, the slice is
// untyped, so Extract dispatches by name across String.Extract (false|string)
// and Object.Extract (unknown), folding to ?|false|string. The false arm is
// concretely inferred - String.Extract really returns false on no match - so
// Size on it deserves a finding even though the union is dirty. It caps at
// warning because the ? arm is unproven.
func TestDirtyUnionMixedReceiverWarns(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Extract", Sig: "(member, x = #extract_no_default) :unknown"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(src) {
			w = src[0 ..].Extract('^x')
			return w.Size()
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityWarning, "not applicable on at least one path"))
	a.That(hasDiag(env, "Foo", SeverityWarning, "unknown arms not checked"))
	a.That(!hasDiag(env, "Foo", SeverityError, "not applicable on at least one path"))
}

// when every concrete arm is bad the union usually reflects broken inference
// upstream, not user code - stay silent as before
func TestDirtyUnionAllBadArmsStaysSilent(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			s = false
			if x is 1
				s = x.NoSuchMethodZq()
			return s.Size()
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityWarning, "not applicable on at least one path"))
	a.That(!hasDiag(env, "Foo", SeverityError, "not applicable on at least one path"))
}
