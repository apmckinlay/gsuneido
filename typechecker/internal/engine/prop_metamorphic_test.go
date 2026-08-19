// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// ---- alpha-rename (length-preserving, so Pos is stable) ----

func renamePrefix(c byte) byte {
	switch c {
	case 'v':
		return 'u'
	case 'M':
		return 'N'
	case 'f':
		return 'g'
	case 'p':
		return 'q'
	}
	return c
}

func renIdent(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	b[0] = renamePrefix(b[0])
	return string(b)
}

func renameExpr(e synth.Expr) synth.Expr {
	switch e.Kind {
	case synth.ExVar, synth.ExMember, synth.ExCall:
		e.Name = renIdent(e.Name)
	}
	na := make([]synth.Expr, len(e.Args))
	for i := range e.Args {
		na[i] = renameExpr(e.Args[i])
	}
	e.Args = na
	return e
}

func renameStmts(ss []synth.Stmt) []synth.Stmt {
	out := make([]synth.Stmt, len(ss))
	for i := range ss {
		s := ss[i]
		switch s.Kind {
		case synth.StAssign, synth.StMemberAssign:
			s.Name = renIdent(s.Name)
		}
		s.Expr = renameExpr(s.Expr)
		s.Body = renameStmts(s.Body)
		s.Else = renameStmts(s.Else)
		out[i] = s
	}
	return out
}

func renameProgram(p synth.Program) synth.Program {
	nm := make([]synth.Member, len(p.Members))
	for i, mem := range p.Members {
		mem.Name = renIdent(mem.Name)
		nm[i] = mem
	}
	nmeth := make([]synth.Method, len(p.Methods))
	for i, m := range p.Methods {
		m.Name = renIdent(m.Name)
		np := make([]synth.Param, len(m.Params))
		for j, pr := range m.Params {
			pr.Name = renIdent(pr.Name)
			np[j] = pr
		}
		m.Params = np
		m.Body = renameStmts(m.Body)
		nmeth[i] = m
	}
	p.Members, p.Methods = nm, nmeth
	return p
}

var identRe = regexp.MustCompile(`\b[vMfp][0-9]+\b`)

func renameMsg(s string) string {
	return identRe.ReplaceAllStringFunc(s, func(m string) string {
		return string(renamePrefix(m[0])) + m[1:]
	})
}

func renamedDiagKey(d Diagnostic) string {
	return fmt.Sprintf("%d\x00%s\x00%d\x00%s", d.Severity, renameMsg(d.Method), d.Pos, renameMsg(d.Msg))
}

func renameKeys(m map[string]DynType) map[string]DynType {
	out := make(map[string]DynType, len(m))
	for k, v := range m {
		out[renIdent(k)] = v
	}
	return out
}

func TestPropMetaAlphaRename(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		rp := renameProgram(p)
		rsrc := synth.Render(&rp)
		clsB, okB := safeParse(rsrc, "T")
		if !okB {
			rt.Fatalf("renamed did not parse but original did:\n%s", rsrc)
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (renamed): %s", pb)
		}
		for _, m := range []struct {
			Name string
			a, b map[string]DynType
		}{
			{"Members", renameKeys(envA.Members), envB.Members},
			{"Returns", renameKeys(envA.Returns), envB.Returns},
			{"PostCtorMembers", renameKeys(envA.PostCtorMembers), envB.PostCtorMembers},
			{"PreCtorReturns", renameKeys(envA.PreCtorReturns), envB.PreCtorReturns},
		} {
			if !typeMapsSemEq(m.a, m.b) {
				rt.Fatalf("%s not isomorphic under rename\n a=%v\n b=%v\nsrc:\n%s", m.Name, m.a, m.b, src)
			}
		}
		ka := sortedDiagKeys(diagList(envA), renamedDiagKey)
		kb := sortedDiagKeys(diagList(envB), diagKey)
		if !strSlicesEqual(ka, kb) {
			rt.Fatalf("diagnostics not isomorphic under rename\n a=%v\n b=%v\nsrc:\n%s", ka, kb, src)
		}
	})
}

// ---- method/member reorder ----

func genPerm(t *rapid.T, n int) []int {
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := rapid.IntRange(0, i).Draw(t, "perm")
		idx[i], idx[j] = idx[j], idx[i]
	}
	return idx
}

func reorderProgram(p synth.Program, methodPerm, memberPerm []int) synth.Program {
	nm := make([]synth.Member, len(p.Members))
	for i, j := range memberPerm {
		nm[i] = p.Members[j]
	}
	nmeth := make([]synth.Method, len(p.Methods))
	for i, j := range methodPerm {
		nmeth[i] = p.Methods[j]
	}
	p.Members, p.Methods = nm, nmeth
	return p
}

func TestPropMetaReorder(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		rp := reorderProgram(p, genPerm(rt, len(p.Methods)), genPerm(rt, len(p.Members)))
		clsB, okB := safeParse(synth.Render(&rp), "T")
		if !okB {
			rt.Fatalf("reordered did not parse:\n%s", synth.Render(&rp))
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (reordered): %s", pb)
		}
		for _, m := range []struct {
			Name string
			a, b map[string]DynType
		}{
			{"Members", envA.Members, envB.Members},
			{"Returns", envA.Returns, envB.Returns},
			{"PostCtorMembers", envA.PostCtorMembers, envB.PostCtorMembers},
			{"PreCtorReturns", envA.PreCtorReturns, envB.PreCtorReturns},
		} {
			if !typeMapsSemEq(m.a, m.b) {
				rt.Fatalf("%s changed under reorder\n a=%v\n b=%v\nsrc:\n%s", m.Name, m.a, m.b, src)
			}
		}
		// diagnostics as (sev, method, msg) multiset, ignoring Pos
		ka := sortedDiagKeys(diagList(envA), diagKeyNoPos)
		kb := sortedDiagKeys(diagList(envB), diagKeyNoPos)
		if !strSlicesEqual(ka, kb) {
			rt.Fatalf("diagnostics (ignoring Pos) changed under reorder\n a=%v\n b=%v\nsrc:\n%s", ka, kb, src)
		}
	})
}

// ---- full-line comment insertion ----

func insertComments(t *rapid.T, src string) string {
	lines := strings.Split(src, "\n")
	n := rapid.IntRange(1, 4).Draw(t, "ncomments")
	for k := range n {
		pos := rapid.IntRange(0, len(lines)).Draw(t, "cpos")
		c := []string{"// comment " + fmt.Sprint(k)}
		lines = append(lines[:pos], append(c, lines[pos:]...)...)
	}
	return strings.Join(lines, "\n")
}

func TestPropMetaComments(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		csrc := insertComments(rt, src)
		clsB, okB := safeParse(csrc, "T")
		if !okB {
			rt.Fatalf("commented source did not parse:\n%s", csrc)
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (commented): %s", pb)
		}
		if !typeMapsSemEq(envA.Members, envB.Members) || !typeMapsSemEq(envA.Returns, envB.Returns) {
			rt.Fatalf("types changed under comment insertion\nsrc:\n%s", csrc)
		}
		ka := sortedDiagKeys(diagList(envA), diagKeyNoPos)
		kb := sortedDiagKeys(diagList(envB), diagKeyNoPos)
		if !strSlicesEqual(ka, kb) {
			rt.Fatalf("diagnostics (ignoring Pos) changed under comments\n a=%v\n b=%v\nsrc:\n%s", ka, kb, csrc)
		}
	})
}

// ---- fresh method with k literal returns -> Returns == fold of literal types ----

func litDynType(l synth.Lit) DynType {
	switch l.Kind {
	case synth.LitNum:
		return TNumber
	case synth.LitStr:
		return TString
	case synth.LitDate:
		return TDate
	case synth.LitTrue:
		return TTrue
	case synth.LitFalse:
		return TFalse
	default:
		return TObject
	}
}

func TestPropMetaKLiteralReturns(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		k := rapid.IntRange(1, 4).Draw(rt, "k")
		lits := make([]synth.Lit, k)
		var want DynType
		for i := range k {
			lits[i] = synth.GTypeLit(rt, rapid.SampledFrom(synth.BaseTypes).Draw(rt, "krt"))
			if i == 0 {
				want = litDynType(lits[i])
			} else {
				want = U(want, litDynType(lits[i]))
			}
		}
		var body []synth.Stmt
		for i := 0; i < k-1; i++ {
			body = append(body, synth.Stmt{Kind: synth.StIf, Expr: synth.Expr{Kind: synth.ExVar, Name: "p0"},
				Body: []synth.Stmt{{Kind: synth.StReturn, Expr: synth.Expr{Kind: synth.ExLit, Lit: lits[i]}}}})
		}
		body = append(body, synth.Stmt{Kind: synth.StReturn, Expr: synth.Expr{Kind: synth.ExLit, Lit: lits[k-1]}})
		prog := synth.Program{Name: "T", Methods: []synth.Method{{Name: "Mret", Params: []synth.Param{{Name: "p0"}}, Body: body}}}
		cls, ok := safeParse(synth.Render(&prog), "T")
		if !ok {
			rt.Fatalf("did not parse:\n%s", synth.Render(&prog))
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", panicMsg, synth.Render(&prog))
		}
		if got := env.Returns["Mret"]; !semEq(got, want) {
			rt.Fatalf("Returns[Mret]=%v want fold=%v\nsrc:\n%s", got, want, synth.Render(&prog))
		}
	})
}

// ---- append an unrelated, self-contained method at the end ----

func TestPropMetaAppendMethod(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		// "Zed": no members, no calls, no params - self-contained.
		fresh := synth.Method{Name: "Zed", Body: []synth.Stmt{
			{Kind: synth.StAssign, Name: "v0", Expr: synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitNum, N: 7}}},
			{Kind: synth.StReturn, Expr: synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitNum, N: 9}}},
		}}
		p2 := p
		p2.Methods = append(append([]synth.Method{}, p.Methods...), fresh)
		clsB, okB := safeParse(synth.Render(&p2), "T")
		if !okB {
			rt.Fatalf("appended did not parse:\n%s", synth.Render(&p2))
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (appended): %s", pb)
		}
		// all existing entries unchanged
		if !typeMapsSemEq(envA.Members, envB.Members) {
			rt.Fatalf("Members changed by appending a method\nsrc:\n%s", src)
		}
		for k, v := range envA.Returns {
			if !semEq(v, envB.Returns[k]) {
				rt.Fatalf("Returns[%q] changed by appending a method\nsrc:\n%s", k, src)
			}
		}
		// existing-method diagnostics unchanged (positions stable: appended at end)
		ka := sortedDiagKeys(diagList(envA), diagKey)
		var existing []Diagnostic
		for _, d := range diagList(envB) {
			if d.Method != "Zed" {
				existing = append(existing, d)
			}
		}
		kb := sortedDiagKeys(existing, diagKey)
		if !strSlicesEqual(ka, kb) {
			rt.Fatalf("existing diagnostics changed by appending a method\n a=%v\n b=%v\nsrc:\n%s", ka, kb, src)
		}
	})
}

// ---- prepend `.freshMember = literal` to the last method's body ----
// prepended, not appended: a write at the tail would become the method's
// implicit return, and that change propagates to every transitive caller
// whose own tail calls it - no per-method exemption can express that

func TestPropMetaAppendMember(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		p := synth.GenProgram(rt, synth.Config{})
		src := synth.Render(&p)
		clsA, okA := safeParse(src, "T")
		if !okA {
			rt.Skip("unparseable")
		}
		envA, pa := safeRun(clsA)
		if pa != "" {
			rt.Fatalf("panic: %s\nsrc:\n%s", pa, src)
		}
		p2 := p
		p2.Methods = append([]synth.Method{}, p.Methods...)
		li := len(p2.Methods) - 1
		last := p2.Methods[li]
		if len(last.Body) == 0 {
			rt.Skip("empty body: the write would be the tail, hence the implicit return")
		}
		last.Body = append([]synth.Stmt{
			{Kind: synth.StMemberAssign, Name: "fz", Expr: synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitNum, N: 3}}}},
			last.Body...)
		p2.Methods[li] = last
		clsB, okB := safeParse(synth.Render(&p2), "T")
		if !okB {
			rt.Fatalf("appended member did not parse:\n%s", synth.Render(&p2))
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (appended member): %s", pb)
		}
		// Members gains exactly fz (Number); everything else identical
		if _, existed := envA.Members["fz"]; existed {
			rt.Skip("generator already used fz (impossible with current names)")
		}
		if got, ok := envB.Members["fz"]; !ok || !semEq(got, TNumber) {
			rt.Fatalf("Members did not gain fz:Number (got %v, ok=%v)\nsrc:\n%s", got, ok, synth.Render(&p2))
		}
		if len(envB.Members) != len(envA.Members)+1 {
			rt.Fatalf("Members gained more than fz: a=%v b=%v", envA.Members, envB.Members)
		}
		for k, v := range envA.Members {
			if !semEq(v, envB.Members[k]) {
				rt.Fatalf("Members[%q] changed by appending a member\nsrc:\n%s", k, src)
			}
		}
		if !typeMapsSemEq(envA.Returns, envB.Returns) {
			rt.Fatalf("Returns changed by prepending a member\nsrc:\n%s", src)
		}
	})
}

// ---- wrap a random statement in `if (true) { ... }` on a wellTyped program ----

func TestPropMetaIfTrueWrap(t *testing.T) {
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
			rt.Skip("base wellTyped program already has an error (covered by Property D)")
		}
		// pick a method with >=1 top-level statement, wrap one statement in if(true){}
		mi := rapid.IntRange(0, len(p.Methods)-1).Draw(rt, "wmeth")
		if len(p.Methods[mi].Body) == 0 {
			rt.Skip("empty method")
		}
		// never wrap the final statement: an if(true) around the tail kills
		// its implicit-return role, which is a real semantic change
		if len(p.Methods[mi].Body) < 2 {
			rt.Skip("only a final statement to wrap")
		}
		si := rapid.IntRange(0, len(p.Methods[mi].Body)-2).Draw(rt, "wstmt")
		p2 := p
		p2.Methods = append([]synth.Method{}, p.Methods...)
		m := p2.Methods[mi]
		nb := append([]synth.Stmt{}, m.Body...)
		nb[si] = synth.Stmt{Kind: synth.StIf, Expr: synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitTrue}},
			Body: []synth.Stmt{m.Body[si]}}
		m.Body = nb
		p2.Methods[mi] = m
		clsB, okB := safeParse(synth.Render(&p2), "T")
		if !okB {
			rt.Fatalf("wrapped did not parse:\n%s", synth.Render(&p2))
		}
		envB, pb := safeRun(clsB)
		if pb != "" {
			rt.Fatalf("panic (wrapped): %s", pb)
		}
		if d, bad := anySeverityError(envB); bad {
			rt.Fatalf("if(true) wrap introduced an error [%s]: %s\nsrc:\n%s", d.Method, d.Msg, synth.Render(&p2))
		}
		// old types fit new (widening allowed, dropping arms not)
		for k, v := range envA.Members {
			if !fits(v, envB.Members[k]) {
				rt.Fatalf("Members[%q]: old %v does not fit new %v\nsrc:\n%s", k, v, envB.Members[k], synth.Render(&p2))
			}
		}
		for k, v := range envA.Returns {
			if !fits(v, envB.Returns[k]) {
				rt.Fatalf("Returns[%q]: old %v does not fit new %v\nsrc:\n%s", k, v, envB.Returns[k], synth.Render(&p2))
			}
		}
	})
}

// ---- add a correct return annotation to a known-type fresh method ----

func TestPropMetaCorrectAnnotation(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		g := rapid.SampledFrom(synth.BaseTypes).Draw(rt, "anng")
		lit := synth.GTypeLit(rt, g)
		mk := func(ann string) TypeEnv {
			prog := synth.Program{Name: "T", Methods: []synth.Method{
				{Name: "Mann", RetAnn: ann, Body: []synth.Stmt{{Kind: synth.StReturn, Expr: synth.Expr{Kind: synth.ExLit, Lit: lit}}}}}}
			cls, ok := safeParse(synth.Render(&prog), "T")
			if !ok {
				rt.Fatalf("did not parse:\n%s", synth.Render(&prog))
			}
			env, panicMsg := safeRun(cls)
			if panicMsg != "" {
				rt.Fatalf("panic: %s", panicMsg)
			}
			return env
		}
		envA := mk("")
		envB := mk(synth.AnnName(g))
		// no new diagnostics from adding the correct annotation (ignore Pos)
		ka := sortedDiagKeys(diagList(envA), diagKeyNoPos)
		kb := sortedDiagKeys(diagList(envB), diagKeyNoPos)
		if !strSlicesEqual(ka, kb) {
			rt.Fatalf("correct annotation :%s changed diagnostics\n a=%v\n b=%v", synth.AnnName(g), ka, kb)
		}
	})
}
