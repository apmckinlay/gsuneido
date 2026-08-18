// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import "github.com/apmckinlay/gsuneido/compile/ast"

func narrowSwitch(sw *ast.Switch, env TypeEnv, sc narrowScope) {
	narrowWalk(sw.E, env, sc)
	if isImplicitTrue(sw.E) {
		narrowSwitchCondMode(sw, env, sc)
	} else {
		narrowSwitchScrutineeMode(sw, env, sc)
	}
}

// matches the bare `switch` form's implicit true and the explicit
// `switch (true)` spelling, whose parens parse as Unary{LParen}
func isImplicitTrue(e ast.Expr) bool {
	c, ok := unwrapConstant(peelParens(e))
	return ok && DynTypeOfSuValue(c.Val) == TTrue
}

func narrowSwitchCondMode(sw *ast.Switch, env TypeEnv, sc narrowScope) {
	defScope := sc.clone()
	for i := range sw.Cases {
		c := &sw.Cases[i]
		for _, e := range c.Exprs {
			narrowWalk(e, env, sc)
		}
		caseScope := sc.clone()
		if len(c.Exprs) == 1 {
			caseScope = refineCond(c.Exprs[0], sc, true, env, false)
		}
		for _, e := range c.Exprs {
			applyRefinement(e, defScope, false, env, false)
		}
		for _, stmt := range c.Body {
			narrowWalk(stmt, env, caseScope)
		}
	}
	for _, stmt := range sw.Default {
		narrowWalk(stmt, env, defScope)
	}
}

func narrowSwitchScrutineeMode(sw *ast.Switch, env TypeEnv, sc narrowScope) {
	id, isIdent := unwrapIdent(sw.E)
	canNarrow := isIdent && !isGlobalIdent(id.Name)

	defScope := sc.clone()
	var allTargets []DynType

	for i := range sw.Cases {
		c := &sw.Cases[i]
		for _, e := range c.Exprs {
			narrowWalk(e, env, sc)
		}
		caseScope := sc.clone()

		if canNarrow {
			targets := caseLiteralTargets(c.Exprs)
			if len(targets) > 0 {
				tgt := narrowTarget{name: id.Name, node: id}
				existing := existingType(sc, tgt, env)
				// allowMembers moot - scrutinee mode is always Ident
				storeRefinement(caseScope, tgt, narrowTowardSet(existing, targets), false)
				allTargets = append(allTargets, targets...)
			}
		}

		for _, stmt := range c.Body {
			narrowWalk(stmt, env, caseScope)
		}
	}

	// `x isnt 1` says nothing about x not being a Number, so only a literal
	// whose type holds a single value can be subtracted here - the same
	// distinction applyEqRefinement makes on its isnt path
	if removable := atomicTargets(allTargets); canNarrow && len(removable) > 0 {
		tgt := narrowTarget{name: id.Name, node: id}
		existing := existingType(defScope, tgt, env)
		storeRefinement(defScope, tgt, narrowAwaySet(existing, removable), false)
	}

	for _, stmt := range sw.Default {
		narrowWalk(stmt, env, defScope)
	}
}

func caseLiteralTargets(exprs []ast.Expr) []DynType {
	var targets []DynType
	for _, e := range exprs {
		if lit, ok := unwrapConstant(e); ok {
			targets = append(targets, DynTypeOfSuValue(lit.Val))
		}
	}
	return targets
}

func atomicTargets(targets []DynType) []DynType {
	var out []DynType
	for _, t := range targets {
		if isAtomicType(t) {
			out = append(out, t)
		}
	}
	return out
}
