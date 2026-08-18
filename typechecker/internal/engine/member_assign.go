// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// backfills env.Members from `.foo = ...` assignments. must run after
// NameResolutionPass so RHS nodes are typed. multi-site assignments are
// unioned - up to the programmer to avoid wide unions.
//
// ```suneido
//
//	class {
//	    Init() { .count = 5 }
//	                      ^ RHS type (TNumber) backfilled onto .count
//	    Reset() { .count = "no" }
//	                       ^^^^ unioned in -> .count becomes Number | String
//	}
//
// ```
func MemberAssignmentPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				collectThisAssignments(stmt, env)
			}
		}
	}
	return false
}

// matches `this.name = rhs` with a literal member name
func unwrapThisAssign(n ast.Node) (string, ast.Expr, bool) {
	b, ok := n.(*ast.Binary)
	if !ok || b.Tok != tok.Eq {
		return "", nil, false
	}
	mem, ok := b.Lhs.(*ast.Mem)
	if !ok {
		return "", nil, false
	}
	id, ok := mem.E.(*ast.Ident)
	if !ok || id.Name != "this" {
		return "", nil, false
	}
	name, ok := memberName(mem.M)
	if !ok {
		return "", nil, false
	}
	return name, b.Rhs, true
}

func collectThisAssignments(n ast.Node, env TypeEnv) {
	if n == nil {
		return
	}
	if name, rhs, ok := unwrapThisAssign(n); ok {
		arms, _ := decomposeForCheck(env.GetType(rhs))
		for _, m := range arms {
			env.UnionMember(name, m)
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		collectThisAssignments(c, env)
		return c
	})
}

// runs once after the main loop. by now all resolvable callsites have
// concrete types, so a still-TUnknown RHS is unknown
// mark the member dirty.
//
// ```suneido
//
//	class {
//	    Set(v) { .x = _parent.Custom }
//	                  ^^^^^^^^^^^^^^^ unresolvable -> .x picks up TUnknown,
//	                                  so .x becomes dirty (e.g. Number | ?)
//	}
//
// ```
func MemberDirtyPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				dirtyThisAssignments(stmt, env)
			}
		}
	}
	return false
}

func dirtyThisAssignments(n ast.Node, env TypeEnv) {
	if n == nil {
		return
	}
	if name, rhs, ok := unwrapThisAssign(n); ok {
		dirtyMemberAssign(name, rhs, env)
	}
	n.Children(func(c ast.Node) ast.Node {
		dirtyThisAssignments(c, env)
		return c
	})
}

func dirtyMemberAssign(name string, rhs ast.Expr, env TypeEnv) {
	// peel wrappers: only the inner expression carries a type
	// this early, so `.x = (cond ? a : b)` would otherwise read
	if _, dirty := decomposeForCheck(env.GetType(peelParens(rhs))); !dirty {
		return
	}
	if existing, ok := env.LookupMember(name); ok && existing != TUnknown {
		env.UnionMember(name, TUnknown)
	}
}
