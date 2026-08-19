// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"slices"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

func narrowIf(x *ast.If, env TypeEnv, sc narrowScope) {
	narrowWalk(x.Cond, env, sc)
	// snapshot before the branches fork - joinBranchKills needs the pre-if
	// facts both as the key set to re-check and as the no-else fall-through
	preTypes := factsInGuard(sc.Locals)
	preMembers := factsInGuard(sc.Members)
	thenScope := refineCond(x.Cond, sc, true, env, true)
	narrowWalk(x.Then, env, thenScope)
	var elseScope narrowScope
	if x.Else != nil {
		elseScope = refineCond(x.Cond, sc, false, env, true)
		narrowWalk(x.Else, env, elseScope)
	}
	// kills first: the merge blocks below rebuild from sc, so dropping a dead
	// refinement here keeps them from resurrecting it
	joinBranchKills(x, sc, preTypes, preMembers, thenScope, elseScope)
	if x.Else == nil && !branchAlwaysExits(x.Then) {
		joinNoElseLocals(x, env, sc)
	}
	if x.Else == nil && !branchAlwaysExits(x.Then) {
		joinNoElseMembers(x, env, sc)
	}
	// early-return: if exactly one branch always exits, siblings after the
	// if inherit the other branch's refinement.
	//
	// ```suneido
	// Foo(x) {
	//     if not Number?(x)
	//         return false       // then-branch exits
	//     ^^^^^^^^^^^^^^^^^
	//     return x + 1           // <-- siblings see x narrowed to TNumber
	//            ^                   (Number? polarity flipped via the `not`)
	// }
	// ```
	// member refinements always installed - narrowWalk's Call/assign/Inc-Dec
	// handlers clear them position-aware as sibling statements walk.
	thenExits := branchAlwaysExits(x.Then)
	elseExits := x.Else != nil && branchAlwaysExits(x.Else)
	if thenExits && !elseExits {
		applyRefinement(x.Cond, sc, false, env, true)
	} else if elseExits && !thenExits {
		applyRefinement(x.Cond, sc, true, env, true)
	}
}

func joinNoElseLocals(x *ast.If, env TypeEnv, sc narrowScope) {
	assigned := topLevelAssignedLocals(x.Then)
	if len(assigned) == 0 {
		return
	}
	negScope := refineCond(x.Cond, sc, false, env, false)
	for name := range assigned {
		elseT, ok := negScope.Locals.guarded(name)
		if !ok {
			continue
		}
		thenT := lastTopLevelAssignType(x.Then, name, env)
		if thenT == nil {
			continue
		}
		var merged DynType
		if elseT == TUnknown {
			entryT := condEntryType(x.Cond, name, env)
			if entryT == nil || entryT == TUnknown {
				continue
			}
			merged = thenT
		} else {
			merged = U(thenT, elseT)
		}
		if merged == nil || merged == TUnknown {
			continue
		}
		sc.Locals.prove(name, merged)
	}
}

func joinNoElseMembers(x *ast.If, env TypeEnv, sc narrowScope) {
	assignedM := topLevelAssignedMembers(x.Then)
	if len(assignedM) == 0 {
		return
	}
	negScope := refineCond(x.Cond, sc, false, env, true)
	for name := range assignedM {
		elseT, ok := negScope.Members.guarded(name)
		if !ok || elseT == TUnknown {
			continue
		}
		thenT := lastTopLevelMemberAssignType(x.Then, name, env)
		if thenT == nil {
			continue
		}
		merged := U(thenT, elseT)
		if merged == TUnknown || typeHasBoolish(merged) {
			continue
		}
		sc.Members.prove(name, merged)
	}
}

// factsInGuard snapshots a scope's live refinements as a plain name->type map.
// Only in-guard entries count: being in a guard is what makes a refinement
// visible to reads (see narrowWalk's *ast.Ident and *ast.Mem cases).
func factsInGuard(r refinements) map[string]DynType {
	out := make(map[string]DynType, len(r))
	for name, f := range r {
		if f.InGuard && f.Typ != nil {
			out[name] = f.Typ
		}
	}
	return out
}

// joinBranchKills propagates branch *removals* back to the enclosing scope.
//
// refineCond hands each branch its own clone, so a refinement a branch
// destroys - a member-writing call, a reassignment nested below the top level
// - dies with that clone and sc walks on believing the stale fact. narrowIf's
// blocks merge back what a branch *adds*; this is the missing mirror for what
// it takes away.
//
// ```suneido
//
//	M1() { .f1 = "x" }
//	M0() {
//	    .f1 = #20200101      // refines .f1 -> Date
//	    if (true) { .M1() }  // the kill lands on the clone, not on sc
//	    return .f1           // ^^^ so this read answered Date, not String
//	}
//
// ```
func joinBranchKills(x *ast.If, sc narrowScope, preTypes, preMembers map[string]DynType,
	thenScope, elseScope narrowScope) {
	var types, members []map[string]DynType
	if !branchAlwaysExits(x.Then) {
		types = append(types, factsInGuard(thenScope.Locals))
		members = append(members, factsInGuard(thenScope.Members))
	}
	switch {
	case x.Else == nil:
		// skipping the if is itself a route out, and it carries the pre-if facts
		types = append(types, preTypes)
		members = append(members, preMembers)
	case !branchAlwaysExits(x.Else):
		types = append(types, factsInGuard(elseScope.Locals))
		members = append(members, factsInGuard(elseScope.Members))
	}
	if len(types) == 0 {
		return // every route exits - no siblings below the if to protect
	}
	joinReachingFacts(preTypes, types, sc.Locals)
	joinReachingFacts(preMembers, members, sc.Members)
}

// joinReachingFacts intersects the pre-if refinements against every route that
// reaches past the if, writing the survivors back into vals/guard. A name keeps
// its refinement only if all routes still refine it, at the union of the types
// they give it; whatever any route dropped is dropped here too!
func joinReachingFacts(pre map[string]DynType, routes []map[string]DynType,
	r refinements) {
	for name := range pre {
		var merged DynType
		for _, r := range routes {
			t, ok := r[name]
			if !ok || t == nil || t == TUnknown {
				merged = nil
				break
			}
			if merged == nil {
				merged = t
			} else {
				merged = U(merged, t)
			}
		}
		if merged == nil || merged == TUnknown {
			delete(r, name)
			continue
		}
		r.prove(name, merged)
	}
}

func branchAlwaysExits(s ast.Statement) bool {
	if s == nil {
		return false
	}
	switch x := s.(type) {
	case *ast.Return, *ast.Throw, *ast.Break, *ast.Continue:
		return true
	case *ast.Compound:
		return slices.ContainsFunc(x.Body, branchAlwaysExits)
	case *ast.If:
		if x.Else == nil {
			return false
		}
		return branchAlwaysExits(x.Then) && branchAlwaysExits(x.Else)
	}
	return false
}

func topLevelAssignedLocals(s ast.Statement) map[string]bool {
	out := map[string]bool{}
	for _, st := range stmtList(s) {
		if id, ok := unwrapStmtAssign(st); ok {
			out[id.Name] = true
		}
	}
	return out
}

func lastTopLevelAssignType(s ast.Statement, name string, env TypeEnv) DynType {
	var lastT DynType
	for _, st := range stmtList(s) {
		id, ok := unwrapStmtAssign(st)
		if !ok || id.Name != name {
			continue
		}
		if t := env.GetType(id); t != TUnknown {
			lastT = t
		}
	}
	return lastT
}

func topLevelAssignedMembers(s ast.Statement) map[string]bool {
	out := map[string]bool{}
	for _, st := range stmtList(s) {
		if name, _, ok := stmtMemberAssign(st); ok {
			out[name] = true
		}
	}
	return out
}

func lastTopLevelMemberAssignType(s ast.Statement, name string, env TypeEnv) DynType {
	var lastT DynType
	for _, st := range stmtList(s) {
		mname, b, ok := stmtMemberAssign(st)
		if !ok || mname != name {
			continue
		}
		if t := env.GetType(b.Rhs); t != TUnknown {
			lastT = t
		}
	}
	return lastT
}

func stmtMemberAssign(s ast.Statement) (string, *ast.Binary, bool) {
	es, ok := s.(*ast.ExprStmt)
	if !ok || es.E == nil {
		return "", nil, false
	}
	expr := es.E
	if ep, ok := expr.(*ast.ExprPos); ok && ep.Expr != nil {
		expr = ep.Expr
	}
	b, ok := expr.(*ast.Binary)
	if !ok || b.Tok != tok.Eq {
		return "", nil, false
	}
	name, _, ok := unwrapThisMember(b.Lhs)
	if !ok {
		return "", nil, false
	}
	return name, b, true
}

func unwrapStmtAssign(s ast.Statement) (*ast.Ident, bool) {
	es, ok := s.(*ast.ExprStmt)
	if !ok || es.E == nil {
		return nil, false
	}
	expr := es.E
	if ep, ok := expr.(*ast.ExprPos); ok && ep.Expr != nil {
		expr = ep.Expr
	}
	b, ok := expr.(*ast.Binary)
	if !ok || b.Tok != tok.Eq {
		return nil, false
	}
	id, ok := b.Lhs.(*ast.Ident)
	if !ok || isGlobalIdent(id.Name) {
		return nil, false
	}
	return id, true
}

func condEntryType(cond ast.Node, name string, env TypeEnv) DynType {
	var found *ast.Ident
	var walk func(ast.Node)
	walk = func(n ast.Node) {
		if found != nil || n == nil {
			return
		}
		if id, ok := n.(*ast.Ident); ok && id.Name == name {
			found = id
			return
		}
		n.Children(func(c ast.Node) ast.Node {
			walk(c)
			return c
		})
	}
	walk(cond)
	if found == nil {
		return nil
	}
	return env.GetType(found)
}
