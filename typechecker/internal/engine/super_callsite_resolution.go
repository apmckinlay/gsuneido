// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
)

// ```suneido
// New() { x = super.Setup() }
//
//	^^^^^^^^^^^^^ resolved against the base class's returns,
//	              not this class's overrides
//
// ```
func SuperCallsiteResolutionPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	parentReturns := pctx.ParentReturns
	if len(parentReturns) == 0 {
		return false
	}
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				annotateSuperCalls(stmt, env, parentReturns)
			}
		}
	}
	return false
}

func annotateSuperCalls(n ast.Node, env TypeEnv, parentReturns map[string]DynType) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		if ep.Expr != nil {
			annotateSuperCalls(ep.Expr, env, parentReturns)
			if t := env.GetType(ep.Expr); t != TUnknown {
				env.SetType(ep, t)
			}
		}
		return
	}
	n.Children(func(c ast.Node) ast.Node {
		annotateSuperCalls(c, env, parentReturns)
		return c
	})
	call, ok := n.(*ast.Call)
	if !ok {
		return
	}
	if env.GetType(call) != TUnknown {
		return
	}
	if t := resolveSuperCall(call, parentReturns); t != TUnknown {
		env.SetType(call, t)
	}
}

func resolveSuperCall(call *ast.Call, parentReturns map[string]DynType) DynType {
	mem, ok := call.Fn.(*ast.Mem)
	if !ok {
		return TUnknown
	}
	id, ok := mem.E.(*ast.Ident)
	if !ok || id.Name != "super" {
		return TUnknown
	}
	name, ok := memberName(mem.M)
	if !ok {
		return TUnknown
	}
	if t, ok := parentReturns[name]; ok && t != TUnknown && t != TVoid {
		return t
	}
	return TUnknown
}
