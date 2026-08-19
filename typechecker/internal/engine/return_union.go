// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
)

// `Foo() : T1 | T2` annotations are hard contracts. when present, the
// annotation is what gets published (so callers see the declared shape,
// not whatever the body actually returned), and every `return expr` is
// verified per-site against it - errors for proven mismatches, warnings
// for unprovable ones (TUnknown / dirty unions).
//
// ```suneido
//
//	class {
//	    Foo(b) {
//	        if b { return 1 }
//	                      ^ TNumber
//	        return "x"
//	               ^^^ TString
//	    }            // -> Foo returns Number | String
//	    Bar() : string {
//	        return 42        // proven mismatch - SeverityError
//	               ^^
//	    }                    // -> Bar publishes :string regardless
//	}
//
// ```
func ReturnUnionPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for name, fn := range cls.SortedMethods {
		if declared, ok := env.AnnotatedReturns[name]; ok {
			checkBodyReturns(name, fn, env, declared)
			env.PublishReturn(name, declared)
		} else {
			env.PublishReturn(name, inferBodyReturnReachable(fn, freeParamOverrides(fn), env))
		}
	}
	return false
}

func collectReturns(n ast.Node, env TypeEnv, acc *DynType) {
	if n == nil {
		return
	}
	if r, ok := n.(*ast.Return); ok && !r.ReturnThrow {
		var ty DynType
		if len(r.Exprs) == 0 {
			ty = TVoid
		} else {
			ty = env.GetType(r.Exprs[0])
			for _, expr := range r.Exprs[1:] {
				ty = U(ty, env.GetType(expr))
			}
		}
		*acc = mergeReturn(*acc, ty)
	}
	n.Children(func(c ast.Node) ast.Node {
		collectReturns(c, env, acc)
		return c
	})
}

func checkBodyReturns(name string, fn *ast.Function, env TypeEnv, declared DynType) {
	for _, stmt := range fn.Body {
		if stmt != nil {
			checkReturnsIn(stmt, name, env, declared)
		}
	}
}

func checkReturnsIn(n ast.Node, name string, env TypeEnv, declared DynType) {
	if n == nil {
		return
	}
	if r, ok := n.(*ast.Return); ok && !r.ReturnThrow {
		var ty DynType
		if len(r.Exprs) == 0 {
			ty = TVoid
		} else {
			ty = env.GetType(r.Exprs[0])
			for _, expr := range r.Exprs[1:] {
				ty = U(ty, env.GetType(expr))
			}
		}
		if d, ok := checkReturnAgainst(name, r.Position(), ty, declared); ok {
			env.Report(&d)
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		checkReturnsIn(c, name, env, declared)
		return c
	})
}

func checkReturnAgainst(method string, pos int, inferred, declared DynType) (Diagnostic, bool) {
	if inferred == TVoid {
		return Diagnostic{}, false
	}
	if declared == nil || declared == TUnknown {
		return Diagnostic{}, false
	}

	members, dirty := decomposeForCheck(inferred)
	var badMembers []DynType
	for _, m := range members {
		if !memberFits(m, declared) {
			badMembers = append(badMembers, m)
		}
	}

	switch {
	case len(badMembers) > 0:
		return Diagnostic{
			Severity: SeverityError,
			Method:   method,
			Pos:      pos,
			Msg: fmt.Sprintf("return type %v not assignable to declared %v",
				inferred, declared),
		}, true
	case dirty:
		return Diagnostic{
			Severity: SeverityWarning,
			Method:   method,
			Pos:      pos,
			Msg: fmt.Sprintf("return type %v contains unknown cannot prove it is assignable to declared %v",
				inferred, declared),
		}, true
	}
	return Diagnostic{}, false
}

func decomposeForCheck(inferred DynType) (members []DynType, dirty bool) {
	if inferred == TUnknown {
		return nil, true
	}
	if u, ok := inferred.(Union); ok {
		out := make([]DynType, 0, len(u.Types))
		for _, t := range u.Types {
			if t == TUnknown {
				dirty = true
				continue
			}
			out = append(out, t)
		}
		return out, dirty || u.IsDirty
	}
	return []DynType{inferred}, false
}

func memberFits(sub, declared DynType) bool {
	if sub == nil || declared == nil {
		return false
	}
	if sub == TVoid {
		return true
	}
	if dynEqual(sub, declared) {
		return true
	}
	if subtypeOf(sub, declared) {
		return true
	}
	if u, ok := declared.(Union); ok {
		for _, alt := range u.Types {
			if memberFits(sub, alt) {
				return true
			}
		}
	}
	return false
}
