// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// anti-faults: the soundness mirror of fault injection. each injected shape
// would only type as a clean wrong operand under an ILLEGAL narrowing - a
// refinement leaking past its guard join, surviving a kill, or applied at the
// wrong polarity. injecting one into an error-free wellTyped program must not
// introduce a SeverityError anywhere; a dirty-operand warning is the ceiling.

func fnPredOn(g synth.GType, target synth.Expr) synth.Expr {
	return synth.Expr{Kind: synth.ExFnPred, Pred: synth.PredName(g), Args: []synth.Expr{target}}
}

func exprMentionsVar(e synth.Expr, name string) bool {
	if e.Kind == synth.ExVar && e.Name == name {
		return true
	}
	for i := range e.Args {
		if exprMentionsVar(e.Args[i], name) {
			return true
		}
	}
	return false
}

func stmtsMentionVar(ss []synth.Stmt, name string) bool {
	for i := range ss {
		s := &ss[i]
		if s.Kind == synth.StAssign && s.Name == name {
			return true
		}
		if exprMentionsVar(s.Expr, name) {
			return true
		}
		if stmtsMentionVar(s.Body, name) || stmtsMentionVar(s.Else, name) {
			return true
		}
	}
	return false
}

// unreferencedParams: params the body never mentions. only these are provably
// Unknown at every top-level splice site - a mentioned param can be legally
// clean-narrowed there (e.g. a guard whose else branch always returns flows
// its refinement past the if), which would make the anti-fault a real fault.
func unreferencedParams(m synth.Method) []synth.Param {
	var out []synth.Param
	for _, p := range m.Params {
		if !stmtsMentionVar(m.Body, p.Name) {
			out = append(out, p)
		}
	}
	return out
}

// writesMemberIR reports whether the statement tree directly assigns member name.
func writesMemberIR(ss []synth.Stmt, name string) bool {
	for i := range ss {
		s := &ss[i]
		if s.Kind == synth.StMemberAssign && s.Name == name {
			return true
		}
		if writesMemberIR(s.Body, name) || writesMemberIR(s.Else, name) {
			return true
		}
	}
	return false
}

// stalePair is a member plus a method that writes it, for the kill-on-call kind.
type stalePair struct {
	member synth.Member
	writer synth.Method
}

// stalePairs finds members of base type Number or String with a writer method.
// faultyOpExpr against the guard type is then legal for the base type, so the
// post-kill read is silent-or-warning and only a stale refinement can error.
func stalePairs(p synth.Program) []stalePair {
	var out []stalePair
	for _, mem := range p.Members {
		b := synth.LitGType(mem.Seed)
		if b != synth.GtNum && b != synth.GtStr {
			continue
		}
		for _, m := range p.Methods {
			// never inject an explicit .New() call - construction semantics
			if m.Name == "New" {
				continue
			}
			if writesMemberIR(m.Body, mem.Name) {
				out = append(out, stalePair{member: mem, writer: m})
			}
		}
	}
	return out
}

// staleGuardType picks a guard type G != base whose faultyOpExpr op is legal
// for base: Number base fails nothing under "+", String base nothing under "$".
func staleGuardType(rt *rapid.T, base synth.GType) synth.GType {
	if base == synth.GtStr {
		return synth.GtNum // faultyOpExpr(_, synth.GtNum) is "$" - legal for String
	}
	return rapid.SampledFrom([]synth.GType{synth.GtStr, synth.GtDate, synth.GtObj, synth.GtBool}).Draw(rt, "afstaleg")
}

// genAntiFaultStmts builds consecutive statements to splice into method m.
// ok=false means the program offers no material for any kind.
func genAntiFaultStmts(rt *rapid.T, p synth.Program, m synth.Method) (stmts []synth.Stmt, label string, ok bool) {
	pairs := stalePairs(p)
	free := unreferencedParams(m)
	var kinds []string
	if len(free) > 0 {
		kinds = append(kinds, "leakAfterIf", "leakAfterWhile", "elseBranch", "notThenBranch")
	}
	if len(free) > 1 {
		kinds = append(kinds, "killByAssign")
	}
	if len(pairs) > 0 {
		kinds = append(kinds, "staleMemberKill")
	}
	if len(kinds) == 0 {
		return nil, "", false
	}
	freeVar := func() synth.Expr {
		return synth.Expr{Kind: synth.ExVar, Name: rapid.SampledFrom(free).Draw(rt, "afp").Name}
	}
	kind := rapid.SampledFrom(kinds).Draw(rt, "afkind")
	g := rapid.SampledFrom(synth.BaseTypes).Draw(rt, "afg")
	switch kind {
	case "leakAfterIf":
		pv := freeVar()
		return []synth.Stmt{
			{Kind: synth.StIf, Expr: fnPredOn(g, pv)},
			{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(pv, g)},
		}, kind, true
	case "leakAfterWhile":
		pv := freeVar()
		return []synth.Stmt{
			{Kind: synth.StWhile, Expr: fnPredOn(g, pv)},
			{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(pv, g)},
		}, kind, true
	case "elseBranch":
		pv := freeVar()
		return []synth.Stmt{
			{Kind: synth.StIf, Expr: fnPredOn(g, pv),
				Else: []synth.Stmt{{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(pv, g)}}},
		}, kind, true
	case "notThenBranch":
		pv := freeVar()
		return []synth.Stmt{
			{Kind: synth.StIf, Expr: synth.Expr{Kind: synth.ExUnary, Op: "not", Args: []synth.Expr{fnPredOn(g, pv)}},
				Body: []synth.Stmt{{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(pv, g)}}},
		}, kind, true
	case "killByAssign":
		// distinct free Params: guard p, overwrite from q, then use p at the guard type
		pi := rapid.IntRange(0, len(free)-1).Draw(rt, "afpi")
		qi := (pi + 1 + rapid.IntRange(0, len(free)-2).Draw(rt, "afqi")) % len(free)
		pn, qn := free[pi].Name, free[qi].Name
		pv := synth.Expr{Kind: synth.ExVar, Name: pn}
		return []synth.Stmt{
			{Kind: synth.StIf, Expr: fnPredOn(g, pv), Body: []synth.Stmt{
				{Kind: synth.StAssign, Name: pn, Expr: synth.Expr{Kind: synth.ExVar, Name: qn}},
				{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(pv, g)},
			}},
		}, kind, true
	default: // staleMemberKill
		pair := rapid.SampledFrom(pairs).Draw(rt, "afpair")
		gm := staleGuardType(rt, synth.LitGType(pair.member.Seed))
		mem := synth.Expr{Kind: synth.ExMember, Name: pair.member.Name}
		callArgs := make([]synth.Expr, len(pair.writer.Params))
		for i := range callArgs {
			callArgs[i] = synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitNum, N: 1}}
		}
		return []synth.Stmt{
			{Kind: synth.StIf, Expr: fnPredOn(gm, mem), Body: []synth.Stmt{
				{Kind: synth.StExpr, Expr: synth.Expr{Kind: synth.ExCall, Name: pair.writer.Name, Args: callArgs}},
				{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(mem, gm)},
			}},
		}, kind, true
	}
}

func TestPropAntiFaultNoFalseErrors(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{WellTyped: true})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		if _, bad := anySeverityError(envA); bad {
			rt.Skip("base wellTyped program already errors (covered by TestPropWellTypedNoError)")
		}

		p2 := p
		p2.Methods = append([]synth.Method{}, p.Methods...)
		mi := rapid.IntRange(0, len(p2.Methods)-1).Draw(rt, "afmeth")
		m := p2.Methods[mi]
		stmts, label, ok := genAntiFaultStmts(rt, p, m)
		if !ok {
			rt.Skip("no anti-fault material (method has no params, no writable member)")
		}
		si := rapid.IntRange(0, len(m.Body)).Draw(rt, "afsite")
		nb := make([]synth.Stmt, 0, len(m.Body)+len(stmts))
		nb = append(nb, m.Body[:si]...)
		nb = append(nb, stmts...)
		nb = append(nb, m.Body[si:]...)
		m.Body = nb
		p2.Methods[mi] = m

		fsrc := synth.Render(&p2)
		clsB, okB := safeParse(fsrc, "T")
		if !okB {
			rt.Fatalf("anti-fault source did not parse (%s):\n%s", label, fsrc)
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic after %s injection: %s\nsrc:\n%s", label, pb, fsrc)
		}
		for _, d := range diagList(envB) {
			if d.Severity == SeverityError {
				rt.Fatalf("anti-fault %q introduced a false error in %s: %s\nsrc:\n%s",
					label, d.Method, d.Msg, fsrc)
			}
		}
	})
}

// pins each anti-fault template: no error, and for the param kinds the use is
// reached as a dirty operand (the "cannot prove" warning), so a silently
// unchecked expression cannot fake a pass.
func TestAntiFaultTemplatesStayClean(t *testing.T) {
	cases := []struct {
		label, src   string
		wantDirtyUse bool
	}{
		{"leakAfterIf", `class { M0(p0)
{
if (Number?(p0))
{
}
v9 = (p0 $ "x")
return 0
} }`, true},
		{"leakAfterWhile", `class { M0(p0)
{
while (Number?(p0))
{
}
v9 = (p0 $ "x")
return 0
} }`, true},
		{"elseBranch", `class { M0(p0)
{
if (Number?(p0))
{
}
else
{
v9 = (p0 $ "x")
}
return 0
} }`, true},
		{"notThenBranch", `class { M0(p0)
{
if (not Number?(p0))
{
v9 = (p0 $ "x")
}
return 0
} }`, true},
		{"killByAssign", `class { M0(p0, p1)
{
if (Number?(p0))
{
p0 = p1
v9 = (p0 $ "x")
}
return 0
} }`, true},
		{"staleMemberKill", `class { f0: "a"
M0()
{
if (Number?(.f0))
{
.M1()
v9 = (.f0 $ "x")
}
return 0
}
M1()
{
.f0 = "b"
return 0
} }`, false},
	}
	for _, c := range cases {
		cls, ok := safeParse(c.src, "T")
		if !ok {
			t.Fatalf("%s: template did not parse:\n%s", c.label, c.src)
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			t.Fatalf("%s: pipeline panicked: %s", c.label, panicMsg)
		}
		dirtyUse := false
		for _, d := range diagList(env) {
			if d.Severity == SeverityError {
				t.Errorf("%s: false error: %s", c.label, d.Msg)
			}
			if d.Severity == SeverityWarning && strings.Contains(d.Msg, "cannot prove") {
				dirtyUse = true
			}
		}
		if c.wantDirtyUse && !dirtyUse {
			t.Errorf("%s: expected the dirty-operand warning (use not reached?), diagnostics: %v",
				c.label, diagList(env))
		}
	}
}
