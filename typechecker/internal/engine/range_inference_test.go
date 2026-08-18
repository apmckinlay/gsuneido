// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// stdlib:SelectFields.ScanFields regression: once the base is a known string,
// a range of it is a string, so Extract resolves to String.Extract and Size
// on the false arm is a hard error - same as the unsliced ScanFormula twin.
func TestRangeToOfKnownStringExtractSizeErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Extract", Sig: "(member, x = #extract_no_default) :unknown"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Trim", Sig: "(chars = ' \t\r\n') :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(src) {
			src = src.Trim()
			w = src[0 ..].Extract('^x')
			return w.Size()
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "no overload for false"))
	a.That(!hasDiag(env, "Foo", SeverityWarning, "unknown arms not checked"))
}

func TestRangeLenOfKnownStringExtractSizeErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Extract", Sig: "(member, x = #extract_no_default) :unknown"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Trim", Sig: "(chars = ' \t\r\n') :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(src) {
			src = src.Trim()
			w = src[0 :: 5].Extract('^x')
			return w.Size()
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "no overload for false"))
	a.That(!hasDiag(env, "Foo", SeverityWarning, "unknown arms not checked"))
}

// a range of an object literal is an object, so calls on it resolve
// directly instead of dispatching by name across builtin overloads
func TestRangeOfObjectIsObject(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() {
			ob = #(1, 2, 3)
			return ob[1 ..].Size()
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityWarning, "receiver type unknown"))
	a.That(!hasDiag(env, "Foo", SeverityError, "not applicable"))
}

// non-rangeable arms are dropped from the projected type - the capability
// check already reports them - so a false arm does not cascade into a
// second finding on the range result
func TestRangeDropsNonRangeableArms(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Extract", Sig: "(member, x = #extract_no_default) :unknown"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Trim", Sig: "(chars = ' \t\r\n') :string"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(src) {
			src = src.Trim()
			w = src.Extract('^x')
			return w[1 ..].Size()
		}
	}`, "T")
	a.That(hasDiag(env, "Foo", SeverityError, "not rangeable"))
	a.That(!hasDiag(env, "Foo", SeverityError, "no overload for false"))
}

// an annotated param is enough to anchor the chain: the exact
// SelectFields.ScanFields shape becomes a hard error once src is declared
func TestRangeOfAnnotatedParamErrors(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Extract", Sig: "(member, x = #extract_no_default) :unknown"},
		Reg{Receiver: "object", Name: "Size", Sig: "(list :boolean =false,named :boolean =false) :number"},
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "class", Name: "Size", Sig: "() :number"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		ScanFields(src: string) {
			pos = 0
			while pos < src.Size() {
				whitespace = src[pos ..].Extract('^(\s*?)\S')
				pos += whitespace.Size()
			}
			return pos
		}
	}`, "T")
	a.That(hasDiag(env, "ScanFields", SeverityError, "no overload for false"))
	a.That(!hasDiag(env, "ScanFields", SeverityWarning, "unknown arms not checked"))
}

// guess provenance flows through a range: when the base's type rests on a
// dispatch-by-name guess, findings on the range result cap at warning
func TestRangeOfGuessTypedBaseCapsAtWarning(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(data) {
			data = data.Replace('\t', ' ')
			indent = data[.. 2]
			indent $= 'x'
			return indent
		}
	}`, "T")
	a.That(!hasDiag(env, "Foo", SeverityError, `operator "$="`))
}

// an untyped base leaves the range untyped - dispatch-by-name behavior
// (TestDirtyUnionMixedReceiverWarns) is unchanged
func TestRangeOfUnknownStaysUnknown(t *testing.T) {
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
	a.That(hasDiag(env, "Foo", SeverityWarning, "unknown arms not checked"))
	a.That(!hasDiag(env, "Foo", SeverityError, "no overload for false"))
}
