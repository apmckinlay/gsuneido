// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// fault injection: the completeness dual of TestPropWellTypedNoError.
// splice one provably-ill-typed statement into an error-free wellTyped
// program; the checker must report an error in that method and nowhere else.
// guarded faults only error if narrowing produced a clean type, so this also
// fails if flow narrowing silently stops working.

// faultyOpExpr wraps operand in an operator that provably rejects type g:
// Number fails `$` (wants String), everything else fails `+` (wants Number).
func faultyOpExpr(operand synth.Expr, g synth.GType) synth.Expr {
	if g == synth.GtNum {
		return synth.Expr{Kind: synth.ExBin, Op: "$", Args: []synth.Expr{
			operand, {Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitStr, S: "x"}}}}
	}
	return synth.Expr{Kind: synth.ExBin, Op: "+", Args: []synth.Expr{
		operand, {Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitNum, N: 1}}}}
}

func fnPredExpr(g synth.GType, param string) synth.Expr {
	return synth.Expr{Kind: synth.ExFnPred, Pred: synth.PredName(g), Args: []synth.Expr{{Kind: synth.ExVar, Name: param}}}
}

var nonBoolGtypes = []synth.GType{synth.GtNum, synth.GtStr, synth.GtDate, synth.GtObj}

// genFaultStmt builds one self-contained ill-typed statement for method m.
// v9 is outside the generator's v0-v4 pool and no fault contains a return,
// so Returns and every other method's typing are untouched - which is what
// lets the property demand fault locality.
// memberFaultG picks the type faultyOpExpr must reject for a member read, and
// whether `if (.mem)` is a provable condition fault. optional members do not
// read at their seed type: an always-assigned one is demoted to its payload,
// a maybe-assigned one reads False|Payload (the False arm rejects "+", the
// payload arm rejects Boolean conditions), and a never-assigned one reads
// False - a legal condition, so only the operator fault convicts.
func memberFaultG(mem synth.Member) (g synth.GType, condOK bool) {
	switch mem.Ctor {
	case synth.CtorAlways:
		return mem.Payload, true
	case synth.CtorMaybe:
		return synth.GtBool, true
	case synth.CtorNever:
		return synth.GtBool, false
	default:
		g = synth.LitGType(mem.Seed)
		return g, g != synth.GtBool
	}
}

func genFaultStmt(rt *rapid.T, members []synth.Member, m synth.Method) (synth.Stmt, string, string) {
	var nonBoolMembers []synth.Member
	for _, mem := range members {
		if _, condOK := memberFaultG(mem); condOK {
			nonBoolMembers = append(nonBoolMembers, mem)
		}
	}
	var kinds []string
	if len(members) > 0 {
		kinds = append(kinds, "memberOp")
	}
	if len(nonBoolMembers) > 0 {
		kinds = append(kinds, "memberCond")
	}
	if len(m.Params) > 0 {
		kinds = append(kinds, "guardOp", "guardCond", "notGuardOp")
	}
	kind := rapid.SampledFrom(kinds).Draw(rt, "fikind")
	switch kind {
	case "memberOp":
		mem := rapid.SampledFrom(members).Draw(rt, "fimem")
		g, _ := memberFaultG(mem)
		return synth.Stmt{Kind: synth.StAssign, Name: "v9",
			Expr: faultyOpExpr(synth.Expr{Kind: synth.ExMember, Name: mem.Name}, g)}, kind, mem.Name
	case "memberCond":
		mem := rapid.SampledFrom(nonBoolMembers).Draw(rt, "fimemc")
		return synth.Stmt{Kind: synth.StIf, Expr: synth.Expr{Kind: synth.ExMember, Name: mem.Name}}, kind, mem.Name
	case "guardCond":
		pn := rapid.SampledFrom(m.Params).Draw(rt, "fiparam").Name
		g := rapid.SampledFrom(nonBoolGtypes).Draw(rt, "figt")
		return synth.Stmt{Kind: synth.StIf, Expr: fnPredExpr(g, pn),
			Body: []synth.Stmt{{Kind: synth.StIf, Expr: synth.Expr{Kind: synth.ExVar, Name: pn}}}}, kind, ""
	case "notGuardOp":
		pn := rapid.SampledFrom(m.Params).Draw(rt, "fiparam").Name
		g := rapid.SampledFrom(synth.BaseTypes).Draw(rt, "figt")
		return synth.Stmt{Kind: synth.StIf,
			Expr: synth.Expr{Kind: synth.ExUnary, Op: "not", Args: []synth.Expr{fnPredExpr(g, pn)}},
			Else: []synth.Stmt{{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(synth.Expr{Kind: synth.ExVar, Name: pn}, g)}}}, kind, ""
	default: // guardOp
		pn := rapid.SampledFrom(m.Params).Draw(rt, "fiparam").Name
		g := rapid.SampledFrom(synth.BaseTypes).Draw(rt, "figt")
		return synth.Stmt{Kind: synth.StIf, Expr: fnPredExpr(g, pn),
			Body: []synth.Stmt{{Kind: synth.StAssign, Name: "v9", Expr: faultyOpExpr(synth.Expr{Kind: synth.ExVar, Name: pn}, g)}}}, kind, ""
	}
}

// index of the first top-level statement referencing member name. a member
// fault is only provable while the member's seed type still holds; past a
// statement that mentions the member, narrowing may have refined it - a
// diverging `if (true is .f0)` leaves the fall-through provably dead for f0,
// where the fault could never execute and the checker rightly stays silent.
func firstMemberRef(body []synth.Stmt, name string) int {
	for i := range body {
		if stmtRefsMember(&body[i], name) {
			return i
		}
	}
	return len(body)
}

func stmtRefsMember(s *synth.Stmt, name string) bool {
	if s.Kind == synth.StMemberAssign && s.Name == name {
		return true
	}
	if exprRefsMember(&s.Expr, name) {
		return true
	}
	for i := range s.Body {
		if stmtRefsMember(&s.Body[i], name) {
			return true
		}
	}
	for i := range s.Else {
		if stmtRefsMember(&s.Else[i], name) {
			return true
		}
	}
	return false
}

func exprRefsMember(e *synth.Expr, name string) bool {
	if e.Kind == synth.ExMember && e.Name == name {
		return true
	}
	for i := range e.Args {
		if exprRefsMember(&e.Args[i], name) {
			return true
		}
	}
	return false
}

func TestPropFaultInjectionCompleteness(t *testing.T) {
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
		mi := rapid.IntRange(0, len(p2.Methods)-1).Draw(rt, "fimeth")
		m := p2.Methods[mi]

		var target, label string
		injectable := len(p.Members) > 0 || len(m.Params) > 0
		if !injectable || rapid.IntRange(0, 5).Draw(rt, "fiann") == 0 {
			// annotation fault: fresh method declared :string returning a Number
			p2.Methods = append(p2.Methods, synth.Method{Name: "Mz", RetAnn: "string",
				Body: []synth.Stmt{{Kind: synth.StReturn, Expr: synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitNum, N: 42}}}}})
			target, label = "Mz", "annotation"
		} else {
			fault, l, faultMem := genFaultStmt(rt, p.Members, m)
			hi := len(m.Body)
			if faultMem != "" {
				hi = firstMemberRef(m.Body, faultMem)
			}
			si := rapid.IntRange(0, hi).Draw(rt, "fisite")
			nb := make([]synth.Stmt, 0, len(m.Body)+1)
			nb = append(nb, m.Body[:si]...)
			nb = append(nb, fault)
			nb = append(nb, m.Body[si:]...)
			m.Body = nb
			p2.Methods[mi] = m
			target, label = m.Name, l
		}

		fsrc := synth.Render(&p2)
		clsB, okB := safeParse(fsrc, "T")
		if !okB {
			rt.Fatalf("injected source did not parse (%s):\n%s", label, fsrc)
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic after %s injection: %s\nsrc:\n%s", label, pb, fsrc)
		}
		hit := 0
		for _, d := range diagList(envB) {
			if d.Severity != SeverityError {
				continue
			}
			if d.Method != target {
				rt.Fatalf("%s fault in %s cascaded an error into %s: %s\nsrc:\n%s",
					label, target, d.Method, d.Msg, fsrc)
			}
			hit++
		}
		if hit == 0 {
			rt.Fatalf("injected %s fault in %s went undetected\nsrc:\n%s", label, target, fsrc)
		}
	})
}

// the dual of the site restriction above: past a diverging guard that
// consumed the member's only value, the fault site is dead - the member has
// no clean type left and the checker must stay at a warning, never convict.
func TestMemberFaultPastGuard(t *testing.T) {
	src := `class { f0: true
M0()
{
if ((true is .f0))
{
return 0
}
v9 = (.f0 + 1)
return 0
} }`
	cls, ok := safeParse(src, "T")
	if !ok {
		t.Fatal("template did not parse")
	}
	env, panicMsg := safeRun(cls)
	if panicMsg != "" {
		t.Fatalf("pipeline panicked: %s", panicMsg)
	}
	for _, d := range diagList(env) {
		if d.Severity == SeverityError {
			t.Errorf("dead fall-through convicted: %s: %s", d.Method, d.Msg)
		}
	}
}

// pins each fault template individually so a template that stops parsing or
// stops erroring is caught here, not as a shrunk property counterexample.
func TestFaultTemplatesDetected(t *testing.T) {
	cases := []struct{ label, src, want string }{
		{"memberOp number", `class { f0: 42
M0()
{
v9 = (.f0 $ "x")
return 0
} }`, `operator "$" expects string`},
		{"memberOp string", `class { f0: "a"
M0()
{
v9 = (.f0 + 1)
return 0
} }`, `operator "+" expects number`},
		{"memberCond", `class { f0: 42
M0()
{
if (.f0)
{
}
return 0
} }`, "if condition expects boolean"},
		{"guardOp", `class { M0(p0)
{
if (Number?(p0))
{
v9 = (p0 $ "x")
}
return 0
} }`, `operator "$" expects string`},
		{"guardOp boolean", `class { M0(p0)
{
if (Boolean?(p0))
{
v9 = (p0 + 1)
}
return 0
} }`, `operator "+" expects number`},
		{"guardCond", `class { M0(p0)
{
if (String?(p0))
{
if (p0)
{
}
}
return 0
} }`, "if condition expects boolean"},
		{"notGuardOp", `class { M0(p0)
{
if (not Number?(p0))
{
}
else
{
v9 = (p0 $ "x")
}
return 0
} }`, `operator "$" expects string`},
		{"annotation", `class { Mz() : string
{
return 42
} }`, "not assignable to declared"},
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
		found := false
		for _, d := range diagList(env) {
			if d.Severity == SeverityError && strings.Contains(d.Msg, c.want) {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: expected a SeverityError containing %q, diagnostics: %v",
				c.label, c.want, diagList(env))
		}
	}
}
