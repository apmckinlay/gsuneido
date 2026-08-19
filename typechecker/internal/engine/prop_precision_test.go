// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

func anySeverityError(env TypeEnv) (Diagnostic, bool) {
	for _, d := range diagList(env) {
		if d.Severity == SeverityError {
			return d, true
		}
	}
	return Diagnostic{}, false
}

func TestPropWellTypedNoError(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		src := genSource(rt, synth.Config{WellTyped: true})
		cls, ok := safeParse(src, "T")
		if !ok {
			rt.Skip("unparseable render")
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			rt.Fatalf("pipeline panicked: %s\nsrc:\n%s", panicMsg, src)
		}
		if d, ok := anySeverityError(env); ok {
			rt.Fatalf("wellTyped program produced SeverityError [%s]: %s\nsrc:\n%s", d.Method, d.Msg, src)
		}
	})
}

// --- precision: published types must match the generator's ground truth ---

func gtypeDynType(g synth.GType) DynType {
	switch g {
	case synth.GtNum:
		return TNumber
	case synth.GtStr:
		return TString
	case synth.GtDate:
		return TDate
	case synth.GtObj:
		return TObject
	default:
		return TBoolean
	}
}

// armsFit: every clean arm of got fits intended (dirt tolerated - see below).
func armsFit(got, intended DynType) bool {
	arms, _ := decomposeForCheck(got)
	for _, a := range arms {
		if !memberFits(a, intended) {
			return false
		}
	}
	return true
}

// eqNarrows reports whether any method body contains an `is`/`isnt`
// comparison. equality narrowing can empty a name's type in the negative
// region, which narrowAwaySet downgrades to Unknown (see eqOperandNames in the
// dirt model) - a dirt source that does not require a param. a program with
// such a comparison is therefore not a clean oracle for the pure-dirt check.
func eqNarrows(p synth.Program) bool {
	var exprHas func(e synth.Expr) bool
	exprHas = func(e synth.Expr) bool {
		if e.Kind == synth.ExBin && (e.Op == "is" || e.Op == "isnt") {
			return true
		}
		for i := range e.Args {
			if exprHas(e.Args[i]) {
				return true
			}
		}
		return false
	}
	var stmtsHave func(ss []synth.Stmt) bool
	stmtsHave = func(ss []synth.Stmt) bool {
		for i := range ss {
			s := &ss[i]
			if exprHas(s.Expr) || stmtsHave(s.Body) || stmtsHave(s.Else) {
				return true
			}
		}
		return false
	}
	for _, m := range p.Methods {
		if stmtsHave(m.Body) {
			return true
		}
	}
	return false
}

// the wellTyped generator constructs every return and member assignment at a
// known type, so the published env is checkable against that ground truth:
// no arm may fall outside the intended type (over-widening), and in a
// param-free program with no equality narrowing - where no Unknown ever
// enters - no published type may be dirty. with params, dirt is doctrine:
// MemberDirtyPass runs before the narrowing loop, so a member assigned from a
// param-derived expression is deliberately dirtied, and returns reading such
// members inherit it. `is`/`isnt` is a separate, param-free dirt source: it
// can empty a union in the negative region (downgraded to Unknown).
func TestPropWellTypedPublishedTypes(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{WellTyped: true})
		src := synth.Render(&p)
		cls, ok := safeParse(src, "T")
		if !ok {
			rt.Skip("unparseable render")
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			rt.Fatalf("pipeline panicked: %s\nsrc:\n%s", panicMsg, src)
		}
		if d, bad := anySeverityError(env); bad {
			rt.Skip("base wellTyped program errors (covered by TestPropWellTypedNoError): " + d.Msg)
		}
		pure := !eqNarrows(p)
		for _, m := range p.Methods {
			if len(m.Params) > 0 {
				pure = false
			}
		}
		for _, m := range p.Methods {
			if !m.RetKnown {
				continue
			}
			intended := gtypeDynType(m.Ret)
			got, ok := env.Returns[m.Name]
			if !ok {
				rt.Fatalf("Returns[%q] missing\nsrc:\n%s", m.Name, src)
			}
			if !armsFit(got, intended) {
				rt.Fatalf("Returns[%q]=%v has an arm outside intended %v\nsrc:\n%s",
					m.Name, got, intended, src)
			}
			if _, dirty := decomposeForCheck(got); pure && dirty {
				rt.Fatalf("param-free program published dirty Returns[%q]=%v\nsrc:\n%s",
					m.Name, got, src)
			}
		}
		for _, mem := range p.Members {
			intended := gtypeDynType(synth.LitGType(mem.Seed))
			if mem.Ctor != synth.CtorNone {
				// optional member: False seed plus the payload New assigns
				intended = U(TFalse, gtypeDynType(mem.Payload))
			}
			got, ok := env.Members[mem.Name]
			if !ok {
				rt.Fatalf("Members[%q] missing\nsrc:\n%s", mem.Name, src)
			}
			if !armsFit(got, intended) {
				rt.Fatalf("Members[%q]=%v has an arm outside intended %v\nsrc:\n%s",
					mem.Name, got, intended, src)
			}
			if _, dirty := decomposeForCheck(got); pure && dirty {
				rt.Fatalf("param-free program published dirty Members[%q]=%v\nsrc:\n%s",
					mem.Name, got, src)
			}
		}
	})
}

// --- completeness: one injected provable fault must surface an error ---

func TestCompletenessOperatorFault(t *testing.T) {
	cases := []struct{ src, want string }{
		{`class { f0: true  M0() { return (.f0 + 1) } }`, "expects number"},
		{`class { f0: 42  M0() { return (.f0 $ "x") } }`, "expects string"},
	}
	for _, c := range cases {
		_, env := runPasses(c.src, "T")
		d, ok := anySeverityError(env)
		if !ok {
			t.Fatalf("expected a SeverityError for %s; got none", c.src)
		}
		if !strings.Contains(d.Msg, c.want) {
			t.Fatalf("for %s want %q, got: %s", c.src, c.want, d.Msg)
		}
	}
}

func TestCompletenessConditionFault(t *testing.T) {
	// a number-typed condition is provably non-boolean
	_, env := runPasses(`class { f0: 42  M0() { if (.f0) { return 1 } return 2 } }`, "T")
	if d, ok := anySeverityError(env); !ok {
		t.Fatalf("expected a SeverityError for a non-boolean condition; got none")
	} else if !strings.Contains(d.Msg, "expects boolean") {
		t.Fatalf("expected a boolean-condition error, got: %s", d.Msg)
	}
}

func TestCompletenessReturnAnnotationFault(t *testing.T) {
	// declared :string but returns a Number literal -> provable mismatch
	_, env := runPasses(`class { M0() : string { return 42 } }`, "T")
	if d, ok := anySeverityError(env); !ok {
		t.Fatalf("expected a SeverityError for a return-annotation mismatch; got none")
	} else if !strings.Contains(d.Msg, "not assignable to declared") {
		t.Fatalf("expected a return-annotation error, got: %s", d.Msg)
	}
}
