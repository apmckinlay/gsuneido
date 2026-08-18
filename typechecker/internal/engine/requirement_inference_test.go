// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

func requirementDiags(env TypeEnv) (errs, warns []Diagnostic) {
	for _, d := range diagList(env) {
		if !strings.Contains(d.Msg, "is passed on to") {
			continue
		}
		if d.Severity == SeverityError {
			errs = append(errs, d)
		} else {
			warns = append(warns, d)
		}
	}
	return
}

// Bar hands x straight to Scanner (builtin, param :string), so x requires
// String and the Number literal at the intra-class call site must error.
func TestRequirement_OneHop(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
	)
	src := `class {
	Bar(x) { return Scanner(x) }
	Go() { return .Bar(0) }
	}`
	_, env := runPasses(src, "T")
	errs, _ := requirementDiags(env)
	if len(errs) != 1 {
		t.Fatalf("want 1 requirement error, got %d: %+v", len(errs), diagList(env))
	}
	if !strings.Contains(errs[0].Msg, "argument `0` to Bar has type number, but \"x\" is passed on to Scanner which requires string") {
		t.Errorf("wrong message: %s", errs[0].Msg)
	}
	if got := ScoreConfidence(&errs[0]); got != 0.75 {
		t.Errorf("confidence = %v, want 0.75", got)
	}
}

// the requirement must climb Inner -> Outer before Go's literal is checked.
func TestRequirement_TwoHop(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
	)
	src := `class {
	Inner(a) { return Scanner(a) }
	Outer(b) { return .Inner(b) }
	Go() { return .Outer(0) }
	}`
	_, env := runPasses(src, "T")
	errs, _ := requirementDiags(env)
	if len(errs) != 1 {
		t.Fatalf("want 1 requirement error, got %d: %+v", len(errs), diagList(env))
	}
	if !strings.Contains(errs[0].Msg, `"b" is passed on to Scanner`) {
		t.Errorf("error should anchor at Outer's call site: %s", errs[0].Msg)
	}
}

// a passthrough param demands nothing; the Number caller stays legal.
func TestRequirement_PassthroughSilent(t *testing.T) {
	src := `class {
	Either(x) { return x }
	Go() { return .Either(0) }
	}`
	_, env := runPasses(src, "T")
	if errs, warns := requirementDiags(env); len(errs)+len(warns) != 0 {
		t.Fatalf("passthrough produced requirement diagnostics: %+v %+v", errs, warns)
	}
}

// String and Number demands disagree: the slot drops to unconstrained and
// neither caller is flagged.
func TestRequirement_ConflictSilent(t *testing.T) {
	src := `class {
	F1(a: string) { return a }
	F2(b: number) { return b }
	Both(x) { .F1(x)
		return .F2(x)
		}
	Go() { return .Both(0) }
	Go2() { return .Both("s") }
	}`
	_, env := runPasses(src, "T")
	if errs, _ := requirementDiags(env); len(errs) != 0 {
		t.Fatalf("conflicting demands still errored: %+v", errs)
	}
}

// inside a String? guard the use is narrowing-typed and proves nothing about
// the caller, so no requirement forms.
func TestRequirement_GuardedDemandSilent(t *testing.T) {
	src := `class {
	Maybe(x) {
		if String?(x)
			return Scanner(x)
		return false
		}
	Go() { return .Maybe(0) }
	}`
	_, env := runPasses(src, "T")
	if errs, _ := requirementDiags(env); len(errs) != 0 {
		t.Fatalf("guarded demand errored: %+v", errs)
	}
}

// a reassigned param's uses say nothing about its entry value.
func TestRequirement_ReassignedParamSilent(t *testing.T) {
	src := `class {
	Re(x) { x = Display(x)
		return Scanner(x)
		}
	Go() { return .Re(0) }
	}`
	_, env := runPasses(src, "T")
	if errs, _ := requirementDiags(env); len(errs) != 0 {
		t.Fatalf("reassigned param errored: %+v", errs)
	}
}

// a human annotation is a contract, not a slot to fill: it keeps the declared
// wording and full confidence.
func TestRequirement_AnnotationUntouched(t *testing.T) {
	src := `class {
	An(x: number) { return x }
	Go() { return .An("s") }
	}`
	_, env := runPasses(src, "T")
	if errs, warns := requirementDiags(env); len(errs)+len(warns) != 0 {
		t.Fatalf("annotated param rendered as inferred: %+v %+v", errs, warns)
	}
	found := false
	for _, d := range diagList(env) {
		if d.Severity == SeverityError &&
			strings.Contains(d.Msg, "not assignable to declared number") {
			found = true
		}
	}
	if !found {
		t.Fatalf("annotation contract stopped firing: %+v", diagList(env))
	}
}

// requirement + dirty arg is guess-on-guess: no dirty-only warning.
func TestRequirement_DirtyArgSilent(t *testing.T) {
	src := `class {
	Bar(x) { return Scanner(x) }
	Go(q = false) { return .Bar(q) }
	}`
	_, env := runPasses(src, "T")
	for _, d := range diagList(env) {
		if strings.Contains(d.Msg, "is passed on to") &&
			strings.Contains(d.Msg, "cannot prove") {
			t.Fatalf("dirty-only warning on inferred requirement: %s", d.Msg)
		}
	}
}

// the requirement crosses the reference registry: Lib.Wrap's param picks up
// String while Lib is processed as a reference, and the checked class's
// Number literal errors against it.
func TestRequirement_CrossReference(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	regs := BuildReferenceRegistry([]RefSource{
		{Name: "Lib", Src: `class {
	Wrap(s) { return Scanner(s) }
	}`},
	})
	src := `class {
	Use() { return Lib.Wrap(0) }
	}`
	_, env := runPassesWith(src, "T", func(e *TypeEnv) { regs.Seed(e) })
	errs, _ := requirementDiags(env)
	if len(errs) != 1 {
		t.Fatalf("want 1 cross-ref requirement error, got %d: %+v", len(errs), diagList(env))
	}
	if !strings.Contains(errs[0].Msg, "argument `0` to Lib.Wrap has type number, but \"s\" is passed on to Scanner which requires string") {
		t.Errorf("wrong message: %s", errs[0].Msg)
	}
}

// registry refs arrive in inheritance order, not call order: Outer is
// registered before the Inner it calls, so only the post-registry fixpoint
// can carry Inner's requirement into Outer.
func TestRequirement_ReferenceOrderIndependent(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
		Reg{Kind: "free", Name: "Use", Sig: "(library :string) :boolean"},
	)
	regs := BuildReferenceRegistry([]RefSource{
		{Name: "Outer", Src: `class {
	CallClass(q) { return Inner(q) }
	}`},
		{Name: "Inner", Src: `class {
	CallClass(s) { return Scanner(s) }
	}`},
	})
	src := `class {
	Use() { return Outer(0) }
	}`
	_, env := runPassesWith(src, "T", func(e *TypeEnv) { regs.Seed(e) })
	errs, _ := requirementDiags(env)
	if len(errs) != 1 {
		t.Fatalf("want 1 error through ref-order fixpoint, got %d: %+v", len(errs), diagList(env))
	}
	if !strings.Contains(errs[0].Msg, `"q" is passed on to Scanner`) {
		t.Errorf("error should anchor at Outer's param: %s", errs[0].Msg)
	}
}

// `if x is false return ...` proves false is a handled caller value: no
// requirement may form on x despite the Scanner use below the guard.
func TestRequirement_SentinelGuardSilent(t *testing.T) {
	src := `class {
	F(x) {
		if x is false
			return false
		return Scanner(x)
		}
	Go() { return .F(false) }
	Go2() { return .F(0) }
	}`
	_, env := runPasses(src, "T")
	if errs, _ := requirementDiags(env); len(errs) != 0 {
		t.Fatalf("sentinel-guarded param errored: %+v", errs)
	}
}

// disjunctive sentinel guard (`a is false or b is false`) covers both params.
func TestRequirement_DisjunctiveSentinelSilent(t *testing.T) {
	src := `class {
	F(a, b) {
		if a is false or b is false
			return
		Scanner(a)
		return Scanner(b)
		}
	Go() { return .F(false, 0) }
	}`
	_, env := runPasses(src, "T")
	if errs, _ := requirementDiags(env); len(errs) != 0 {
		t.Fatalf("disjunctive sentinel errored: %+v", errs)
	}
}

// PdfMergerEncrypt.Setup shape: the guard tests u only, but its early return
// also protects o's use below - neither param may pick up a requirement.
func TestRequirement_CrossParamDominanceSilent(t *testing.T) {
	src := `class {
	Setup(u, o) {
		if u is false
			return false
		Scanner(u)
		return Scanner(o)
		}
	Go() { return .Setup(false, false) }
	}`
	_, env := runPasses(src, "T")
	if errs, _ := requirementDiags(env); len(errs) != 0 {
		t.Fatalf("cross-param dominated use errored: %+v", errs)
	}
}

// an `is ”` test is not a sentinel guard - IsCustomized?'s shape must keep
// its demand and the flagship finding.
func TestRequirement_EmptyStringTestStillDemands(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
	)
	src := `class {
	F(q) {
		if q is ""
			return false
		return Scanner(q)
		}
	Go() { return .F(0) }
	}`
	_, env := runPasses(src, "T")
	errs, _ := requirementDiags(env)
	if len(errs) != 1 {
		t.Fatalf("empty-string test wrongly suppressed the demand: %+v", diagList(env))
	}
}

func TestRequirement_DemandBeforeGuardKept(t *testing.T) {
	withSigs(t,
		Reg{Kind: "free", Name: "Scanner", Sig: "(string :string) :unknown"},
	)
	src := `class {
	F(x, y) {
		Scanner(x)
		if y is false
			return
		return x
		}
	Go() { return .F(0, 1) }
	}`
	_, env := runPasses(src, "T")
	errs, _ := requirementDiags(env)
	if len(errs) != 1 {
		t.Fatalf("pre-guard demand lost: %+v", diagList(env))
	}
}
