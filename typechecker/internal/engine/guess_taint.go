// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// guessed calls carry a clean type, so the guess would otherwise leak through
// locals. this pass taints locals assigned from an guessed call (or an alias
// of one) so check passes can cap findings on them at warning.
//
// ```suneido
// v = p.Size()      // Number is a guess -> v tainted
// return v $ "x"
//
//	^ checks on v warn, never error
//
// ```
func GuessTaintPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for name, fn := range cls.SortedMethods {
		for changed := true; changed; {
			changed = false
			for _, stmt := range fn.Body {
				if stmt != nil && taintWalk(stmt, name, env) {
					changed = true
				}
			}
		}
	}
	return false
}

func taintWalk(n ast.Node, method string, env TypeEnv) bool {
	if n == nil {
		return false
	}
	changed := false
	if b, ok := n.(*ast.Binary); ok && b.Tok == tok.Eq {
		if id, ok := b.Lhs.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
			if exprGuessed(b.Rhs, method, env) && !env.GuessedVar(method, id.Name) {
				env.SetGuessedVar(method, id.Name)
				changed = true
			}
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		if taintWalk(c, method, env) {
			changed = true
		}
		return c
	})
	return changed
}

func callReceiverGuessed(call *ast.Call, method string, env TypeEnv) bool {
	if mem, ok := call.Fn.(*ast.Mem); ok {
		return exprGuessed(mem.E, method, env)
	}
	return false
}

// direct provenance only: an guessed call, an alias of a tainted local, or
// either behind parens or a range. compound expressions are not chased -
// their result type comes from the operator, not the guess. ranges are
// pass-throughs: their type comes from the base, so its provenance carries.
func exprGuessed(e ast.Expr, method string, env TypeEnv) bool {
	for {
		switch x := e.(type) {
		case *ast.ExprPos:
			if x.Expr == nil {
				return false
			}
			e = x.Expr
		case *ast.Unary:
			if x.Tok != tok.LParen {
				return false
			}
			e = x.E
		case *ast.RangeTo:
			e = x.E
		case *ast.RangeLen:
			e = x.E
		case *ast.Call:
			return env.GuessedCall(x)
		case *ast.Ident:
			return env.GuessedVar(method, x.Name)
		default:
			return false
		}
	}
}
