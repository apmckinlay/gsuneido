// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func crashDiags(env TypeEnv, method string) []Diagnostic {
	var out []Diagnostic
	for _, d := range methodDiags(env, method) {
		if strings.Contains(d.Msg, "may fail at runtime") {
			out = append(out, d)
		}
	}
	return out
}

func TestCatCrashMemberPredicate(t *testing.T) {
	a := assert.T(t)
	for _, m := range []DynType{TObject, TSequence, TDate, TFunction, TBlock,
		TClass, Instance{Class: "Foo"}} {
		a.That(catCrashMember(m))
	}
	for _, m := range []DynType{TString, TNumber, TBoolean, TTrue, TFalse,
		TVoid, TUnknown} {
		a.That(!catCrashMember(m))
	}
}

func TestCatCrashConcreteObject(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { ob = Object(); return "x" $ ob }
	}`, "T")
	ds := crashDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	d := ds[0]
	a.This(d.Severity).Is(SeverityError)
	a.This(d.Flag).Is(FlagNone)
	a.That(strings.Contains(d.Msg, `can't convert object to String`))
	a.This(ScoreConfidence(&d)).Is(0.90)
}

func TestCatCrashGuardNarrowedDate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) {
			if Date?(x)
				return "n" $ x
			return ""
		}
	}`, "T")
	ds := crashDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	d := ds[0]
	a.This(d.Severity).Is(SeverityError)
	a.This(d.Flag).Is(FlagNone)
	a.That(strings.Contains(d.Msg, `can't convert date to String`))
	a.This(ScoreConfidence(&d)).Is(0.90)
}

func TestCatCrashUnionArm(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(p) { x = p is 1 ? "s" : Object(); return x $ "y" }
	}`, "T")
	ds := crashDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	d := ds[0]
	a.That(strings.Contains(d.Msg, `can't convert object to String`))
	a.That(strings.Contains(d.Msg, "operand type"))
	a.This(ScoreConfidence(&d)).Is(0.70)
}

func TestCatCrashFunctionGuard(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(f) {
			if Function?(f)
				return "x" $ f
			return ""
		}
	}`, "T")
	ds := crashDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.That(strings.Contains(ds[0].Msg, "to String"))
}

func TestCatCrashDirtyStaysBelowGate(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(p) {
			x = Object()
			if p is 1
				x = .undefined()
			return "s" $ x
		}
	}`, "T")
	ds := crashDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	d := ds[0]
	a.This(d.Severity).Is(SeverityError)
	a.This(ScoreConfidence(&d)).Is(0.40)
}

func TestCatEqCrashLhs(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { ob = Object(); ob $= "x"; return ob }
	}`, "T")
	ds := crashDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.That(strings.Contains(ds[0].Msg, `operator "$="`))
}

func TestCatCrashNaryChain(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { ob = Object(); return "a" $ ob $ "c" }
	}`, "T")
	ds := crashDiags(env, "Foo")
	a.This(len(ds)).Is(1)
	a.That(strings.Contains(ds[0].Msg, `can't convert object to String`))
}

func TestCatNumberStaysBehindFlag(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo() { n = 5; return "x" $ n }
	}`, "T")
	a.This(len(crashDiags(env, "Foo"))).Is(0)
	var found bool
	for _, d := range methodDiags(env, "Foo") {
		if strings.Contains(d.Msg, `operator "$" expects string, got number`) {
			found = true
			a.This(d.Flag).Is(FlagStrictStringConcat)
			a.This(ScoreConfidence(&d)).Is(0.20)
		}
	}
	a.That(found)
}

func TestCatCrashGuessedReceiverNotPromoted(t *testing.T) {
	withSigs(t,
		Reg{Receiver: "string", Name: "Chunks", Sig: "() :object"},
		Reg{Receiver: "object", Name: "Chunks", Sig: "() :object"},
	)
	a := assert.T(t)
	_, env := runPasses(`class {
		Foo(x) { return "s" $ x.Chunks() }
	}`, "T")
	a.This(len(crashDiags(env, "Foo"))).Is(0)
	var warned bool
	for _, d := range methodDiags(env, "Foo") {
		if d.Severity == SeverityWarning && strings.Contains(d.Msg, "type guessed") {
			warned = true
			a.This(d.Flag).Is(FlagStrictStringConcat)
		}
	}
	a.That(warned)
}
