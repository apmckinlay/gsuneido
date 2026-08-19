// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// Runs after FlowNarrowingPass so narrowed condition types are visible.
// ```suneido
// if .count
//
//	^^^^^^ Number - not a valid boolean condition
//
// ```
func BooleanConditionCheckPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for name, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				checkConditionsIn(stmt, name, env)
			}
		}
	}
	return false
}

func checkConditionsIn(n ast.Node, method string, env TypeEnv) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		if ep.Expr != nil {
			checkConditionsIn(ep.Expr, method, env)
		}
		return
	}
	switch e := n.(type) {
	case *ast.If:
		checkBooleanContext(e.Cond, "if condition", method, env)
	case *ast.While:
		checkBooleanContext(e.Cond, "while condition", method, env)
	case *ast.DoWhile:
		checkBooleanContext(e.Cond, "do-while condition", method, env)
	case *ast.For:
		checkBooleanContext(e.Cond, "for condition", method, env)
	case *ast.Trinary:
		checkBooleanContext(e.Cond, "?: condition", method, env)
	case *ast.Nary:
		if e.Tok == tok.And || e.Tok == tok.Or {
			what := fmt.Sprintf("operator %q operand", e.Tok.String())
			for _, operand := range e.Exprs {
				checkBooleanContext(operand, what, method, env)
			}
		}
	case *ast.Unary:
		if e.Tok == tok.Not {
			what := fmt.Sprintf("operator %q operand", e.Tok.String())
			checkBooleanContext(e.E, what, method, env)
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		checkConditionsIn(c, method, env)
		return c
	})
}

func checkBooleanContext(cond ast.Expr, what, method string, env TypeEnv) {
	if cond == nil {
		return
	}
	ty := conditionType(env, cond)
	if ty == nil || ty == TVoid {
		return
	}
	// dirty deliberately ignored: an unprovable condition is guessed boolean.
	members, _ := decomposeForCheck(ty)
	for _, m := range members {
		if !memberFits(m, TBoolean) {
			if exprGuessed(cond, method, env) {
				env.Report(&Diagnostic{
					Severity: SeverityWarning,
					Method:   method,
					Pos:      condPos(cond),
					Msg: fmt.Sprintf("%s expects boolean, got %v "+
						"(type guessed from builtin overloads; receiver unknown)",
						what, ty),
				})
				return
			}
			env.Report(&Diagnostic{
				Severity: SeverityError,
				Method:   method,
				Pos:      condPos(cond),
				Msg:      fmt.Sprintf("%s expects boolean, got %v", what, ty),
			})
			return
		}
	}
}

func condPos(cond ast.Node) int {
	if ep, ok := cond.(*ast.ExprPos); ok {
		if ep.Pos > 0 {
			return int(ep.Pos)
		}
		if ep.Expr != nil {
			return nodePos(ep.Expr)
		}
	}
	return nodePos(cond)
}

func conditionType(env TypeEnv, n ast.Node) DynType {
	if ep, ok := n.(*ast.ExprPos); ok {
		if t := env.GetType(ep); t != TUnknown {
			return t
		}
		if ep.Expr != nil {
			return env.GetType(ep.Expr)
		}
	}
	return env.GetType(n)
}
