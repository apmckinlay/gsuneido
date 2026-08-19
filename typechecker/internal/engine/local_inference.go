// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"

	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

// collects and types simple class literals for a given class
func LocalInference(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	// these are simple members for which we can 100% know the types like
	// ```suneido
	// class {
	// 		foo: -1
	// 		bar: false
	// 		baz: #()
	//     ^^^^^^^^^^^^ we are inferring these
	// }
	// ```
	for name, val := range cls.Members {
		t := DynTypeOfSuValue(val)
		if isPublicMember(name) {
			t = widenBool(t)
		}
		env.SeedMember(name, t)
	}

	for name, fn := range cls.SortedMethods {
		if fn.ReturnAnnotation != "" {
			if t, err := ParseTypeAnnotation(fn.ReturnAnnotation); err != nil {
				env.Report(&Diagnostic{
					Severity: SeverityWarning,
					Method:   name,
					Pos:      int(fn.Pos1),
					Msg:      fmt.Sprintf("return %v", err),
				})
			} else if t != TUnknown {
				env.AnnotatedReturns[name] = t
			}
		}

		// (locally) infer each param by itself
		//
		// these are default parameters with a value, so we can infer their type
		// ```suneido
		// Foo(x = 0, a = "name", o = #(), b = false) {}
		//     ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ we are inferring these
		// ```
		for i := range fn.Params {
			p := &fn.Params[i]
			annotated, annotatedOk := paramAnnotationType(p, name, env)
			if annotatedOk {
				env.SetParam(p, annotated)
				if len(p.Name.Name) > 0 && p.Name.Name[0] == '.' {
					env.UnionMember(p.Name.ParamName(), annotated)
				}
				if p.DefVal != nil {
					defType := DynTypeOfSuValue(p.DefVal)
					if !typeFits(defType, annotated) {
						env.Report(&Diagnostic{
							Severity: SeverityError,
							Method:   name,
							Pos:      int(p.Name.Pos),
							Msg: fmt.Sprintf(
								"param %q default has type %v, annotation says %v",
								p.Name.ParamName(), defType, annotated),
						})
					}
				}
				continue
			}
			if p.DefVal != nil {
				t := DynTypeOfSuValue(p.DefVal)
				if len(p.Name.Name) > 0 && p.Name.Name[0] == '.' {
					t = widenBool(t)
					env.UnionMember(p.Name.ParamName(), t)
				}
				env.SetParam(p, t)
			} else if len(p.Name.Name) > 0 && p.Name.Name[0] == '@' {
				env.SetParam(p, TObject)
			}
		}

		for _, stmt := range fn.Body {
			if stmt != nil {
				inferStmtOrExpr(stmt, env)
			}
		}
	}
	return false
}

func paramAnnotationType(p *ast.Param, fnName string, env TypeEnv) (DynType, bool) {
	if p.Annotations == "" {
		return nil, false
	}
	t, err := ParseTypeAnnotation(p.Annotations)
	if err != nil {
		env.Report(&Diagnostic{
			Severity: SeverityWarning,
			Method:   fnName,
			Pos:      int(p.Name.Pos),
			Msg:      fmt.Sprintf("param %q %v", p.Name.ParamName(), err),
		})
	}
	if t == TUnknown || t == nil {
		return nil, false
	}
	return t, true
}

func typeFits(sub, sup DynType) bool {
	if sub == nil || sup == nil || sub == TUnknown || sup == TUnknown {
		return true
	}
	if dynEqual(sub, sup) {
		return true
	}
	if subtypeOf(sub, sup) {
		return true
	}
	if u, ok := sup.(Union); ok {
		for _, alt := range u.Types {
			if typeFits(sub, alt) {
				return true
			}
		}
	}
	return false
}

// works regardless of operand types - just infers the result type assuming
// operands are well-formed.
//
// if x and y are well typed then by the nature of + z must be a TNumber
// ```suneido
// z = x + y
// ^ we are inferring these
// ```
// a `: false` / `= false` seed is a sentinel - "unset, a real value comes
// later" - not a boolean. model it as False | ? so `isnt false` narrows to
// unknown, not the bogus True that Boolean would give. `: true` stays Boolean:
// a true default is usually a genuine flag, set true/false externally.
// ```suneido
//
//	class {
//		x: false
//		^^^^^^^^ "unset for now" sentinel -> False | ? so `isnt false` narrows to unknown, not True
//		y: true
//		^^^^^^^ real flag -> Boolean
//	}
//
// ```
func widenBool(t DynType) DynType {
	if t == TFalse {
		return markDirty(TFalse)
	}
	if t == TTrue {
		return TBoolean
	}
	return t
}

// uppercase = publicly writable by Suneido convention, so the seed is widened; lowercase is closed-world
func isPublicMember(name string) bool { return startsUpper(name) }

func inferResultTypeOfOperator(t tok.Token) DynType {
	return typealgebra.ResultTypeOfOp(t.String())
}

// the union of a trinary's arms, minus arms the condition proves impossible:
// in `x is false ? A : x` the else arm never sees False, so it must not
// contribute one. same rule ConstructorExecPass applies (evalTrinary).
func trinaryType(x *ast.Trinary, env TypeEnv) DynType {
	tType, fType := env.GetType(x.T), env.GetType(x.F)
	tEmpty, fEmpty := false, false
	if tgt, ok := ctorTarget(x.T); ok && condForcesNonFalse(x.Cond, tgt, true) {
		tType, tEmpty = dropFalseClean(tType)
	}
	if tgt, ok := ctorTarget(x.F); ok && condForcesNonFalse(x.Cond, tgt, false) {
		fType, fEmpty = dropFalseClean(fType)
	}
	switch {
	case tEmpty && fEmpty:
		return TUnknown
	case tEmpty:
		return fType
	case fEmpty:
		return tType
	default:
		return U(tType, fType)
	}
}

func RefreshTrinaryTypes(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				refreshExprWalk(stmt, env)
			}
		}
	}
	return false
}

func refreshExprWalk(n ast.Node, env TypeEnv) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *ast.ExprPos:
		refreshInner(x, x.Expr, env)
		return
	case *ast.Trinary:
		refreshExprWalk(x.Cond, env)
		refreshExprWalk(x.T, env)
		refreshExprWalk(x.F, env)
		if t := trinaryType(x, env); t != TUnknown {
			env.SetType(x, t)
		}
		return
	case *ast.Unary:
		if x.Tok == tok.LParen && x.E != nil {
			refreshInner(x, x.E, env)
			return
		}
	case *ast.RangeTo, *ast.RangeLen:
		n.Children(func(c ast.Node) ast.Node {
			refreshExprWalk(c, env)
			return c
		})
		stampRangeType(n, env)
		return
	}
	n.Children(func(c ast.Node) ast.Node {
		refreshExprWalk(c, env)
		return c
	})
}

func stampRangeType(n ast.Node, env TypeEnv) {
	var base ast.Expr
	switch x := n.(type) {
	case *ast.RangeTo:
		base = x.E
	case *ast.RangeLen:
		base = x.E
	default:
		return
	}
	if t := rangeResultType(base, env); t != TUnknown {
		env.SetType(n, t)
	}
}

// rangeResultType projects a base type through E[i..j] / E[i::n].
// String arms stay string; object, record, and sequence arms yield object
// (SuSequence ranges by instantiating to an SuObject). Non-rangeable arms
// are dropped — the capability check reports those — and unknown arms keep
// the result dirty.
func rangeResultType(base ast.Expr, env TypeEnv) DynType {
	if base == nil {
		return TUnknown
	}
	members, dirty := decomposeForCheck(env.GetType(base))
	var out DynType
	for _, m := range members {
		var r DynType
		switch m {
		case TString:
			r = TString
		case TObject, TSequence:
			r = TObject
		default:
			continue
		}
		if out == nil {
			out = r
		} else {
			out = U(out, r)
		}
	}
	if out == nil {
		return TUnknown
	}
	if dirty {
		out = markDirty(out)
	}
	return out
}

// walks a wrapper's inner expr and lifts its refreshed type onto the wrapper
func refreshInner(outer ast.Node, inner ast.Expr, env TypeEnv) {
	if inner == nil {
		return
	}
	refreshExprWalk(inner, env)
	if t := env.GetType(inner); t != TUnknown {
		env.SetType(outer, t)
	}
}

func inferStmtOrExpr(child ast.Node, env TypeEnv) ast.Node {
	if child == nil {
		return nil
	}
	switch n := child.(type) {
	case ast.Expr:
		inferExpression(n, env)
		return n
	case ast.Statement:
		inferStatement(n, env)
		return n
	default:
		return child
	}
}

func inferStatement(stmt ast.Statement, env TypeEnv) {
	stmt.Children(func(n ast.Node) ast.Node { return inferStmtOrExpr(n, env) })
}

func inferExpression(expr ast.Expr, env TypeEnv) {
	if ep, ok := expr.(*ast.ExprPos); ok {
		if ep.Expr != nil {
			inferExpression(ep.Expr, env)
			if t := env.GetType(ep.Expr); t != TUnknown {
				env.SetType(ep, t)
			}
		}
		return
	}
	// we need to recurse here or else we miss some trivial cases for inference
	expr.Children(func(n ast.Node) ast.Node { return inferStmtOrExpr(n, env) })
	switch e := expr.(type) {
	case *ast.Symbol:
		env.SetType(e, TString)
	case *ast.Constant:
		ty := DynTypeOfSuValue(e.Val)
		env.SetType(e, ty)
	case *ast.Nary:
		if ty := inferResultTypeOfOperator(e.Tok); ty != TUnknown {
			env.SetType(e, ty)
		}
	case *ast.Trinary:
		if t := trinaryType(e, env); t != TUnknown {
			env.SetType(e, t)
		}
	case *ast.Binary:
		if ty := inferResultTypeOfOperator(e.Tok); ty != TUnknown {
			env.SetType(e, ty)
		}
	case *ast.Unary:
		if ty := inferResultTypeOfOperator(e.Tok); ty != TUnknown {
			env.SetType(e, ty)
		}
	case *ast.In, *ast.InRange:
		env.SetType(e, TBoolean)
	case *ast.RangeTo, *ast.RangeLen:
		stampRangeType(e, env)
	}
}
