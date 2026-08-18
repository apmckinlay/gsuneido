// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"sort"
	"strings"
)

// PassCtx carries what individual passes need beyond the class and the env.
// Most passes use none of it.
type PassCtx struct {
	Annotations   AnnotationSet
	ParentReturns map[string]DynType
}

// NewPassCtx builds the context for one run. Annotations() copies the builtin
// signature table, so hoist this out of any loop over classes.
func NewPassCtx() *PassCtx {
	return &PassCtx{Annotations: Annotations()}
}

// Pass reports whether it changed anything that warrants re-running other
// passes. Only AssertMemberPass, ConstructorExecPass and RequirementPass can
// say yes; the rest always return false.
type Pass func(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool

// every pipeline pass, checked against Pass at compile time
var _ = [...]Pass{
	ArityCheckPass,
	AssertAssignmentCheckPass,
	AssertMemberPass,
	BooleanConditionCheckPass,
	CallsiteCheckPass,
	CallsiteResolutionPass,
	CapabilityCheckPass,
	ComputeSummaries,
	ConstructorExecPass,
	DateNarrowingPass,
	FlowNarrowingPass,
	GuessTaintPass,
	LocalInference,
	MemberAssignmentPass,
	MemberDirtyPass,
	NameResolutionPass,
	RefreshTrinaryTypes,
	RequirementPass,
	ReturnUnionPass,
	StaticMemberCheckPass,
	SuperCallsiteResolutionPass,
	TypeCheckPass,
}

const maxFixpointPasses = 64

func envTypeSignature(env TypeEnv) string {
	sig := func(t DynType) string {
		if t == nil {
			return "-"
		}
		return t.String()
	}
	parts := make([]string, 0, len(env.Nodes)+len(env.Params)+len(env.Returns)+len(env.Members))
	for n, ty := range env.Nodes {
		parts = append(parts, fmt.Sprintf("n%p=%s", n, sig(ty)))
	}
	for p, ty := range env.Params {
		parts = append(parts, fmt.Sprintf("a%p=%s", p, sig(ty)))
	}
	for k, v := range env.Returns {
		parts = append(parts, "r:"+k+"="+sig(v))
	}
	for k, v := range env.Members {
		parts = append(parts, "m:"+k+"="+sig(v))
	}
	for k, v := range env.PostCtorMembers {
		parts = append(parts, "o:"+k+"="+sig(v))
	}
	for k, v := range env.PreCtorReturns {
		parts = append(parts, "q:"+k+"="+sig(v))
	}
	sort.Strings(parts)
	return strings.Join(parts, "|")
}

func iterateToFixpoint(env TypeEnv, body func()) {
	prev := envTypeSignature(env)
	for range maxFixpointPasses {
		body()
		cur := envTypeSignature(env)
		if cur == prev {
			return
		}
		prev = cur
	}
}

func RunPipeline(cls *ClassObject, env TypeEnv, pctx *PassCtx,
	parentReturns map[string]DynType) {
	if parentReturns == nil {
		parentReturns = map[string]DynType{}
	}

	// sigs must be bound before the first CallsiteResolutionPass
	env = env.WithClass(cls, buildMethodSigs(cls))

	pctx.ParentReturns = parentReturns

	DateNarrowingPass(cls, env, pctx)

	LocalInference(cls, env, pctx)
	NameResolutionPass(cls, env, pctx)
	MemberAssignmentPass(cls, env, pctx)
	ReturnUnionPass(cls, env, pctx)

	iterateToFixpoint(env, func() {
		CallsiteResolutionPass(cls, env, pctx)
		SuperCallsiteResolutionPass(cls, env, pctx)
		NameResolutionPass(cls, env, pctx)
		MemberAssignmentPass(cls, env, pctx)
		ReturnUnionPass(cls, env, pctx)
	})

	MemberDirtyPass(cls, env, pctx)
	// MemberAssignmentPass is deliberately NOT re-run after MemberDirtyPass
	NameResolutionPass(cls, env, pctx)
	ReturnUnionPass(cls, env, pctx)

	// after member unions have settled, before the narrowing phase; re-stamp only if it pinned something
	if AssertMemberPass(cls, env, pctx) {
		NameResolutionPass(cls, env, pctx)
		ReturnUnionPass(cls, env, pctx)
	}

	if ConstructorExecPass(cls, env, pctx) {
		env.CaptureSeedReturns()
		CallsiteResolutionPass(cls, env, pctx)
		ReturnUnionPass(cls, env, pctx)
	}

	// NameResolutionPass is intentionally omitted from this loop (it would undo the narrowing)
	iterateToFixpoint(env, func() {
		FlowNarrowingPass(cls, env, pctx)
		RefreshTrinaryTypes(cls, env, pctx)
		CallsiteResolutionPass(cls, env, pctx)
		SuperCallsiteResolutionPass(cls, env, pctx)
		clearReturnDiagnostics(env) // scrub warnings from earlier iterations that later ones invalidated
		ReturnUnionPass(cls, env, pctx)
	})

	// check passes from here down: order-independent among themselves, strictly after the narrowing loop

	// taint locals fed by guesses so the checks below can downgrade
	GuessTaintPass(cls, env, pctx)

	// after taint (so guessed sigs are excluded), before the checks that
	// consume the inferred param requirements
	// its "changed" result matters only to the registry fixpoint, not here
	RequirementPass(cls, env, pctx)

	// after the narrowing loop so RHS stamps inside guards are narrowed.
	AssertAssignmentCheckPass(cls, env, pctx)

	TypeCheckPass(cls, env, pctx)

	// after FlowNarrowing (above) so narrowed condition types are visible.
	BooleanConditionCheckPass(cls, env, pctx)

	CallsiteCheckPass(cls, env, pctx)

	ArityCheckPass(cls, env, pctx)

	StaticMemberCheckPass(cls, env, pctx)

	CapabilityCheckPass(cls, env, pctx)

	ComputeSummaries(cls, env, pctx)
}

// keep in sync with the "return type ..." diagnostics ReturnUnionPass emits
func clearReturnDiagnostics(env TypeEnv) {
	if env.Diagnostics == nil {
		return
	}
	kept := (*env.Diagnostics)[:0]
	for _, d := range *env.Diagnostics {
		if !strings.HasPrefix(d.Msg, "return type ") {
			kept = append(kept, d)
		}
	}
	*env.Diagnostics = kept
}

func extensionThisType(className string) DynType {
	switch className {
	case "Strings":
		return TString
	case "Numbers":
		return TNumber
	case "Dates":
		return TDate
	case "Objects", "Records":
		return TObject
	}
	return nil
}
