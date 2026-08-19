// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"maps"
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

// dirt model: an independent prediction of where the dirty flag may legally
// appear in a wellTyped program's published types. dirt enters only through
// param reads (Unknown until narrowing) and flows only through ternary arms,
// local assignments, member assignments, and call returns - never through
// operators, whose result types are fixed (typealgebra.ResultTypeOfOp).
//
// two worlds per expression:
//   pre:  narrowing has not run; every param read is dirty. member types are
//         built in this world (MemberAssignmentPass runs before narrowing).
//   post: narrowing has run; param reads are clean (the wellTyped generator
//         only emits narrowed reads). published Returns live in this world.
//
// at a branch join the post bit degrades to the pre bit for locals assigned
// inside a branch: narrowing refinements die at the join, so later reads see
// the pre-world union. the model is conservative toward dirt, which is safe:
// the property asserts only the clean side - model-clean must be published
// clean. model-dirty carries no assertion (the engine may be more precise).

type dirtBits struct{ pre, post bool }

type dirtModel struct {
	member  map[string]bool
	preRet  map[string]bool
	postRet map[string]bool
}

type dirtWalker struct {
	dm      *dirtModel
	params  map[string]bool
	locals  map[string]dirtBits
	forced  map[string]bool // names is/isnt-compared by an enclosing condition
	retPre  bool
	retPost bool
	changed *bool
}

// eqOperandNames collects var and member names used as is/isnt operands.
// equality narrowing can empty such a name's type in the negative region,
// which narrowAwaySet deliberately reverts to TUnknown (see the emptied-union
// downgrade in flow_narrowing) - so reads inside the branch may be dirty.
func eqOperandNames(e synth.Expr, out map[string]bool) {
	if e.Kind == synth.ExBin && (e.Op == "is" || e.Op == "isnt") {
		for i := range e.Args {
			if a := &e.Args[i]; a.Kind == synth.ExVar || a.Kind == synth.ExMember {
				out[a.Name] = true
			}
		}
	}
	for i := range e.Args {
		eqOperandNames(e.Args[i], out)
	}
}

// withForced returns the forced set extended by cond's eq operands (nil-safe).
func (w *dirtWalker) withForced(cond synth.Expr) map[string]bool {
	names := map[string]bool{}
	eqOperandNames(cond, names)
	if len(names) == 0 {
		return w.forced
	}
	maps.Copy(names, w.forced)
	return names
}

func (w *dirtWalker) exprDirt(e synth.Expr) dirtBits {
	switch e.Kind {
	case synth.ExVar:
		if w.params[e.Name] {
			return dirtBits{pre: true, post: w.forced[e.Name]}
		}
		d := w.locals[e.Name]
		d.post = d.post || w.forced[e.Name]
		return d
	case synth.ExMember:
		d := w.dm.member[e.Name]
		return dirtBits{pre: d, post: d || w.forced[e.Name]}
	case synth.ExCall:
		return dirtBits{pre: w.dm.preRet[e.Name], post: w.dm.postRet[e.Name]}
	case synth.ExTern:
		saved := w.forced
		w.forced = w.withForced(e.Args[0])
		a, b := w.exprDirt(e.Args[1]), w.exprDirt(e.Args[2])
		w.forced = saved
		return dirtBits{pre: a.pre || b.pre, post: a.post || b.post}
	default:
		// literals; operators (fixed result type regardless of operand dirt)
		return dirtBits{}
	}
}

func cloneLocals(m map[string]dirtBits) map[string]dirtBits {
	return maps.Clone(m)
}

// joinLocals merges two branch exits relative to their shared entry. locals
// untouched by both branches keep their entry bits; locals assigned in either
// branch OR the bits and degrade post to pre (refinements die at the join).
func joinLocals(entry, a, b map[string]dirtBits) map[string]dirtBits {
	out := make(map[string]dirtBits, len(a)+len(b))
	names := map[string]bool{}
	for k := range a {
		names[k] = true
	}
	for k := range b {
		names[k] = true
	}
	for n := range names {
		av, aok := a[n]
		bv, bok := b[n]
		ev, eok := entry[n]
		if eok && (!aok || av == ev) && (!bok || bv == ev) {
			out[n] = ev
			continue
		}
		j := dirtBits{pre: av.pre || bv.pre, post: av.post || bv.post}
		j.post = j.post || j.pre
		out[n] = j
	}
	return out
}

func localsEqual(a, b map[string]dirtBits) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// loopWrittenLocals mirrors collectLoopWrites: locals assigned anywhere in a
// loop body (recursing into nested blocks) lose their loop-entry narrowing. a
// self-assign (`v = v`) can't change the value, so it doesn't count (see
// selfAssign in flow_narrowing).
func loopWrittenLocals(ss []synth.Stmt, out map[string]bool) {
	for i := range ss {
		s := &ss[i]
		switch s.Kind {
		case synth.StAssign:
			if s.Expr.Kind != synth.ExVar || s.Expr.Name != s.Name {
				out[s.Name] = true
			}
		case synth.StIf:
			loopWrittenLocals(s.Body, out)
			loopWrittenLocals(s.Else, out)
		case synth.StWhile:
			loopWrittenLocals(s.Body, out)
		}
	}
}

func (w *dirtWalker) walkStmts(ss []synth.Stmt) {
	for i := range ss {
		s := &ss[i]
		switch s.Kind {
		case synth.StAssign:
			w.locals[s.Name] = w.exprDirt(s.Expr)
		case synth.StMemberAssign:
			if w.exprDirt(s.Expr).pre && !w.dm.member[s.Name] {
				w.dm.member[s.Name] = true
				*w.changed = true
			}
		case synth.StReturn:
			if !s.Bare {
				d := w.exprDirt(s.Expr)
				w.retPre = w.retPre || d.pre
				w.retPost = w.retPost || d.post
			}
		case synth.StIf:
			savedForced := w.forced
			w.forced = w.withForced(s.Expr)
			entry := cloneLocals(w.locals)
			w.walkStmts(s.Body)
			thenExit := w.locals
			w.locals = cloneLocals(entry)
			w.walkStmts(s.Else)
			w.locals = joinLocals(entry, thenExit, w.locals)
			w.forced = savedForced
			// a branch that always returns keeps the guard's refinement alive
			// for the rest of the block - and with it the chance the negative
			// side emptied a name's type to TUnknown (which polarity survives
			// depends on is vs isnt, so force either way)
			if stmtsTerminate(s.Body) || stmtsTerminate(s.Else) {
				w.forced = w.withForced(s.Expr)
			}
		case synth.StWhile:
			savedForced := w.forced
			w.forced = w.withForced(s.Expr)
			// loop-written locals lose their narrowing for the whole body: the
			// engine's loopEntryScope deletes it, so a read before the in-loop
			// reassignment (e.g. `return v` above `v = p`) sees the un-narrowed
			// param type. degrade post to pre; a later reassignment overwrites it.
			written := map[string]bool{}
			loopWrittenLocals(s.Body, written)
			for name := range written {
				if d, ok := w.locals[name]; ok {
					d.post = d.post || d.pre
					w.locals[name] = d
				}
			}
			// dirt only grows, so iterate the body join to a fixpoint
			for {
				entry := cloneLocals(w.locals)
				w.walkStmts(s.Body)
				w.locals = joinLocals(entry, entry, w.locals)
				if localsEqual(entry, w.locals) {
					break
				}
			}
			w.forced = savedForced
		}
		// synth.StExpr: no local effect; a callee's member writes are handled by
		// walking the callee's own body in the global fixpoint
	}
}

// stmtsTerminate reports whether every path through the list ends in a return.
func stmtsTerminate(ss []synth.Stmt) bool {
	if len(ss) == 0 {
		return false
	}
	last := &ss[len(ss)-1]
	if last.Kind == synth.StReturn {
		return true
	}
	return last.Kind == synth.StIf &&
		stmtsTerminate(last.Body) && stmtsTerminate(last.Else)
}

func computeDirtModel(p synth.Program) dirtModel {
	dm := dirtModel{
		member:  map[string]bool{},
		preRet:  map[string]bool{},
		postRet: map[string]bool{},
	}
	for {
		changed := false
		for _, m := range p.Methods {
			params := make(map[string]bool, len(m.Params))
			for _, pr := range m.Params {
				params[pr.Name] = true
			}
			w := &dirtWalker{dm: &dm, params: params,
				locals: map[string]dirtBits{}, changed: &changed}
			w.walkStmts(m.Body)
			if w.retPre && !dm.preRet[m.Name] {
				dm.preRet[m.Name] = true
				changed = true
			}
			if w.retPost && !dm.postRet[m.Name] {
				dm.postRet[m.Name] = true
				changed = true
			}
		}
		if !changed {
			return dm
		}
	}
}

func isDirty(t DynType) bool {
	_, dirty := decomposeForCheck(t)
	return dirty
}

// model-clean types must be published clean; model-dirty is unasserted.
func TestPropDirtModel(t *testing.T) {
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
		dm := computeDirtModel(p)
		for _, mem := range p.Members {
			if dm.member[mem.Name] {
				continue
			}
			if got := env.Members[mem.Name]; isDirty(got) {
				rt.Fatalf("model-clean member %q published dirty: %v\nsrc:\n%s",
					mem.Name, got, src)
			}
		}
		for _, m := range p.Methods {
			if !m.RetKnown || dm.postRet[m.Name] {
				continue
			}
			if got := env.Returns[m.Name]; isDirty(got) {
				rt.Fatalf("model-clean Returns[%q] published dirty: %v\nsrc:\n%s",
					m.Name, got, src)
			}
		}
	})
}

// pins the model itself and the engine's agreement on hand-checkable cases.
func TestDirtModelKnownCases(t *testing.T) {
	litNumExpr := func(n int) synth.Expr { return synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitNum, N: n}} }
	numMember := synth.Member{Name: "f0", Seed: synth.Lit{Kind: synth.LitNum, N: 1}}
	ret := func(e synth.Expr) synth.Stmt { return synth.Stmt{Kind: synth.StReturn, Expr: e} }

	cases := []struct {
		label        string
		prog         synth.Program
		wantMemDirty map[string]bool
		wantRetPost  map[string]bool
	}{
		{
			label: "param-assigned member dirties member and its readers",
			prog: synth.Program{Name: "T", Members: []synth.Member{numMember}, Methods: []synth.Method{
				{Name: "M0", Params: []synth.Param{{Name: "p0"}}, RetKnown: true, Ret: synth.GtNum, Body: []synth.Stmt{
					{Kind: synth.StIf, Expr: fnPredOn(synth.GtNum, synth.Expr{Kind: synth.ExVar, Name: "p0"}), Body: []synth.Stmt{
						{Kind: synth.StMemberAssign, Name: "f0", Expr: synth.Expr{Kind: synth.ExVar, Name: "p0"}},
					}},
					ret(litNumExpr(1)),
				}},
				{Name: "M1", RetKnown: true, Ret: synth.GtNum, Body: []synth.Stmt{
					ret(synth.Expr{Kind: synth.ExMember, Name: "f0"}),
				}},
			}},
			wantMemDirty: map[string]bool{"f0": true},
			wantRetPost:  map[string]bool{"M0": false, "M1": true},
		},
		{
			label: "operator laundering: param inside a binop stays clean",
			prog: synth.Program{Name: "T", Members: []synth.Member{numMember}, Methods: []synth.Method{
				{Name: "M0", Params: []synth.Param{{Name: "p0"}}, RetKnown: true, Ret: synth.GtNum, Body: []synth.Stmt{
					{Kind: synth.StIf, Expr: fnPredOn(synth.GtNum, synth.Expr{Kind: synth.ExVar, Name: "p0"}), Body: []synth.Stmt{
						{Kind: synth.StMemberAssign, Name: "f0",
							Expr: synth.Expr{Kind: synth.ExBin, Op: "+", Args: []synth.Expr{
								synth.Expr{Kind: synth.ExVar, Name: "p0"}, litNumExpr(1)}}},
					}},
					ret(litNumExpr(1)),
				}},
			}},
			wantMemDirty: map[string]bool{"f0": false},
			wantRetPost:  map[string]bool{"M0": false},
		},
		{
			label: "unreachable is-arm reverts to Unknown and dirties the ternary",
			prog: synth.Program{Name: "T", Methods: []synth.Method{
				{Name: "M0", RetKnown: true, Ret: synth.GtBool, Body: []synth.Stmt{
					{Kind: synth.StAssign, Name: "v0", Expr: synth.Expr{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitTrue}}},
					ret(synth.Expr{Kind: synth.ExTern, Args: []synth.Expr{
						{Kind: synth.ExBin, Op: "is", Args: []synth.Expr{
							{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitTrue}}, {Kind: synth.ExVar, Name: "v0"}}},
						{Kind: synth.ExLit, Lit: synth.Lit{Kind: synth.LitFalse}},
						{Kind: synth.ExVar, Name: "v0"}}}),
				}},
			}},
			wantRetPost: map[string]bool{"M0": true},
		},
		{
			label: "guard-local param carried out of the guard dirties the return",
			prog: synth.Program{Name: "T", Methods: []synth.Method{
				{Name: "M0", Params: []synth.Param{{Name: "p0"}}, RetKnown: true, Ret: synth.GtNum, Body: []synth.Stmt{
					{Kind: synth.StAssign, Name: "v0", Expr: litNumExpr(1)},
					{Kind: synth.StIf, Expr: fnPredOn(synth.GtNum, synth.Expr{Kind: synth.ExVar, Name: "p0"}), Body: []synth.Stmt{
						{Kind: synth.StAssign, Name: "v0", Expr: synth.Expr{Kind: synth.ExVar, Name: "p0"}},
					}},
					ret(synth.Expr{Kind: synth.ExVar, Name: "v0"}),
				}},
			}},
			wantRetPost: map[string]bool{"M0": true},
		},
	}
	for _, c := range cases {
		dm := computeDirtModel(c.prog)
		for name, want := range c.wantMemDirty {
			if dm.member[name] != want {
				t.Errorf("%s: model member[%s]=%v want %v", c.label, name, dm.member[name], want)
			}
		}
		for name, want := range c.wantRetPost {
			if dm.postRet[name] != want {
				t.Errorf("%s: model postRet[%s]=%v want %v", c.label, name, dm.postRet[name], want)
			}
		}
		src := synth.Render(&c.prog)
		cls, ok := safeParse(src, "T")
		if !ok {
			t.Fatalf("%s: did not parse:\n%s", c.label, src)
		}
		env, panicMsg := safeRun(cls)
		if panicMsg != "" {
			t.Fatalf("%s: pipeline panicked: %s", c.label, panicMsg)
		}
		// engine agreement on the clean side
		for name, want := range c.wantMemDirty {
			if !want && isDirty(env.Members[name]) {
				t.Errorf("%s: engine published dirty member %s: %v", c.label, name, env.Members[name])
			}
		}
		for name, want := range c.wantRetPost {
			if !want && isDirty(env.Returns[name]) {
				t.Errorf("%s: engine published dirty Returns[%s]: %v", c.label, name, env.Returns[name])
			}
		}
	}
}
