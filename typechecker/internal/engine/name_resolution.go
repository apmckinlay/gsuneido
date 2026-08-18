// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

func memberName(n ast.Node) (string, bool) {
	switch x := n.(type) {
	case *ast.Constant:
		return x.Val.ToStr()
	case *ast.Symbol:
		return x.Val.ToStr()
	}
	return "", false
}

type scope map[string]DynType

func NameResolutionPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for _, fn := range cls.SortedMethods {
		sc := make(scope, len(fn.Params)+8)
		for i := range fn.Params {
			p := &fn.Params[i]
			paramName := p.Name.ParamName()
			if t, ok := env.Params[p]; ok {
				sc[paramName] = t
			} else if len(p.Name.Name) > 0 && p.Name.Name[0] == '.' {
				if memberType, ok := env.LookupMember(paramName); ok {
					sc[paramName] = memberType
				} else {
					sc[paramName] = TUnknown
				}
			} else {
				sc[paramName] = TUnknown
			}
		}
		walkStmtList(fn.Body, env, sc)
	}
	return false
}

func walkStmtList(stmts []ast.Statement, env TypeEnv, sc scope) {
	for _, stmt := range stmts {
		if stmt == nil {
			continue
		}
		walkNode(stmt, env, sc)
		if iff, ok := stmt.(*ast.If); ok {
			narrowSeedAfterDivergingIf(iff, sc)
		}
	}
}

func narrowSeedAfterDivergingIf(iff *ast.If, sc scope) {
	if iff.Else != nil || !branchAlwaysExits(iff.Then) {
		return
	}
	cond := iff.Cond
	if ep, ok := cond.(*ast.ExprPos); ok && ep.Expr != nil {
		cond = ep.Expr
	}
	b, ok := cond.(*ast.Binary)
	if !ok || b.Tok != tok.Is {
		return
	}
	tgt, lit, ok := targetAndLiteral(b)
	if !ok || tgt.isMember {
		return
	}
	t := DynTypeOfSuValue(lit.Val)
	if t != TFalse && t != TTrue {
		return
	}
	if cur, ok := sc[tgt.name]; ok {
		sc[tgt.name] = narrowAwaySet(cur, []DynType{t})
	}
}

func walkNode(n ast.Node, env TypeEnv, sc scope) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *ast.ExprPos:
		walkExprPos(x, env, sc)
		return
	case *ast.Compound:
		walkStmtList(x.Body, env, sc)
		return
	case *ast.If:
		walkIf(x, env, sc)
		return
	case *ast.Switch:
		walkSwitch(x, env, sc)
		return
	case *ast.TryCatch:
		walkTryCatch(x, env, sc)
		return
	case *ast.While:
		walkLoopFixpoint(sc, func(s scope) {
			walkNode(x.Cond, env, s)
			walkNode(x.Body, env, s)
		})
		return
	case *ast.DoWhile:
		walkLoopFixpoint(sc, func(s scope) {
			walkNode(x.Body, env, s)
			walkNode(x.Cond, env, s)
		})
		return
	case *ast.Forever:
		walkLoopFixpoint(sc, func(s scope) {
			walkNode(x.Body, env, s)
		})
		return
	case *ast.For:
		walkFor(x, env, sc)
		return
	case *ast.ForIn:
		walkForIn(x, env, sc)
		return
	case *ast.Binary:
		if x.Tok == tok.Eq {
			walkAssign(x, env, sc)
			return
		}
		if x.Tok.IsAssign() {
			walkCompoundAssign(x, env, sc)
			return
		}
	case *ast.Unary:
		if isIncDec(x.Tok) && walkIncDec(x, env, sc) {
			return
		}
	}

	n.Children(func(c ast.Node) ast.Node {
		walkNode(c, env, sc)
		return c
	})
	stampNodeType(n, env, sc)
}

func walkExprPos(ep *ast.ExprPos, env TypeEnv, sc scope) {
	if ep.Expr == nil {
		return
	}
	walkNode(ep.Expr, env, sc)
	if t := env.GetType(ep.Expr); t != TUnknown {
		env.SetType(ep, t)
	}
}

func walkIf(x *ast.If, env TypeEnv, sc scope) {
	walkNode(x.Cond, env, sc)
	thenSc := cloneScope(sc)
	walkNode(x.Then, env, thenSc)
	elseSc := sc
	if x.Else != nil {
		elseSc = cloneScope(sc)
		walkNode(x.Else, env, elseSc)
	}
	mergeScopesN(sc, []scope{thenSc, elseSc})
}

func walkSwitch(sw *ast.Switch, env TypeEnv, sc scope) {
	if sw.E != nil {
		walkNode(sw.E, env, sc)
	}
	for i := range sw.Cases {
		for _, e := range sw.Cases[i].Exprs {
			walkNode(e, env, sc)
		}
	}
	armScopes := make([]scope, 0, len(sw.Cases)+1)
	for i := range sw.Cases {
		armSc := cloneScope(sc)
		for _, stmt := range sw.Cases[i].Body {
			walkNode(stmt, env, armSc)
		}
		armScopes = append(armScopes, armSc)
	}
	if sw.Default != nil {
		defSc := cloneScope(sc)
		for _, stmt := range sw.Default {
			walkNode(stmt, env, defSc)
		}
		armScopes = append(armScopes, defSc)
	} else {
		// no-default fall-through preserves entry-state types
		armScopes = append(armScopes, sc)
	}
	mergeScopesN(sc, armScopes)
}

func walkTryCatch(tc *ast.TryCatch, env TypeEnv, sc scope) {
	trySc := cloneScope(sc)
	walkNode(tc.Try, env, trySc)
	catchSc := cloneScope(sc)
	joinScopeInto(catchSc, trySc)
	if tc.CatchVar.Name != "" {
		catchSc[tc.CatchVar.Name] = TString
		env.SetType(&tc.CatchVar, TString)
	}
	if tc.Catch != nil {
		walkNode(tc.Catch, env, catchSc)
	}
	mergeScopesN(sc, []scope{trySc, catchSc})
}

func walkFor(x *ast.For, env TypeEnv, sc scope) {
	for _, e := range x.Init {
		walkNode(e, env, sc)
	}
	walkLoopFixpoint(sc, func(s scope) {
		if x.Cond != nil {
			walkNode(x.Cond, env, s)
		}
		walkNode(x.Body, env, s)
		for _, e := range x.Inc {
			walkNode(e, env, s)
		}
	})
}

// seed loop vars before walking the body so refs resolve.
//
// ```suneido
// for c in "abc" { ... }   // c -> TString (1-char string per element)
//
//	^                ^
//
// for x in ob    { ... }   // x -> TUnknown (no per-element tracking)
//
//	^
//
// for k, v in ob { ... }   // both stay TUnknown - semantics vary by iterable
//
//	^  ^
//
// ```
func walkForIn(x *ast.ForIn, env TypeEnv, sc scope) {
	if x.E != nil {
		walkNode(x.E, env, sc)
	}
	if x.E2 != nil {
		walkNode(x.E2, env, sc)
	}
	varType := DynType(TUnknown)
	if x.Var2.Name == "" && env.GetType(x.E) == TString {
		varType = TString
	}
	walkLoopFixpoint(sc, func(s scope) {
		if x.Var.Name != "" {
			s[x.Var.Name] = varType
			env.SetType(&x.Var, varType)
		}
		if x.Var2.Name != "" {
			s[x.Var2.Name] = TUnknown
			env.SetType(&x.Var2, TUnknown)
		}
		if x.Body != nil {
			walkNode(x.Body, env, s)
		}
	})
}

func walkAssign(b *ast.Binary, env TypeEnv, sc scope) {
	walkNode(b.Rhs, env, sc)
	rhsType := env.GetType(b.Rhs)
	if id, ok := b.Lhs.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		sc[id.Name] = rhsType
		env.SetType(id, rhsType)
	} else {
		walkNode(b.Lhs, env, sc)
	}
	env.SetType(b, rhsType)
}

func walkCompoundAssign(b *ast.Binary, env TypeEnv, sc scope) {
	walkNode(b.Rhs, env, sc)
	resultType := inferResultTypeOfOperator(b.Tok)
	if id, ok := b.Lhs.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		if pre, ok := sc[id.Name]; ok {
			env.SetType(id, pre)
		}
		sc[id.Name] = resultType
	} else {
		walkNode(b.Lhs, env, sc)
	}
	env.SetType(b, resultType)
}

func isIncDec(t tok.Token) bool {
	return t == tok.Inc || t == tok.Dec || t == tok.PostInc || t == tok.PostDec
}

func walkIncDec(u *ast.Unary, env TypeEnv, sc scope) bool {
	id, ok := u.E.(*ast.Ident)
	if !ok || isGlobalIdent(id.Name) {
		return false
	}
	if pre, ok := sc[id.Name]; ok {
		env.SetType(id, pre)
	}
	sc[id.Name] = TNumber
	env.SetType(u, TNumber)
	return true
}

func stampNodeType(n ast.Node, env TypeEnv, sc scope) {
	switch x := n.(type) {
	case *ast.Symbol:
		env.SetType(x, TString)
	case *ast.Constant:
		env.SetType(x, DynTypeOfSuValue(x.Val))
	case *ast.Ident:
		// Foo() {
		// 		x = 123
		// 		y = "h"
		// 		^^^^^^^^ we are inferring based on `=` (tok.Eq)
		// }
		if !isGlobalIdent(x.Name) {
			if t, ok := sc[x.Name]; ok {
				env.SetType(x, t)
			}
		}
	case *ast.Mem:
		// Foo() {
		// 		x = .a
		// 		y = .b
		// 		^^^^^^^^ we are inferring based on `=` and resolving the this.a and this.b call(s)
		// }
		inferMemRead(x, env)
	case *ast.Binary:
		if t := inferResultTypeOfOperator(x.Tok); t != TUnknown {
			env.SetType(x, t)
		}
	case *ast.Nary:
		if t := inferResultTypeOfOperator(x.Tok); t != TUnknown {
			env.SetType(x, t)
		}
	case *ast.Unary:
		if x.Tok == tok.LParen && x.E != nil {
			if t := env.GetType(x.E); t != TUnknown {
				env.SetType(x, t)
			}
		} else if t := inferResultTypeOfOperator(x.Tok); t != TUnknown {
			env.SetType(x, t)
		}
	case *ast.Trinary:
		if t := trinaryType(x, env); t != TUnknown {
			env.SetType(x, t)
		}
	case *ast.In:
		env.SetType(x, TBoolean)
	case *ast.InRange:
		env.SetType(x, TBoolean)
	case *ast.Block:
		env.SetType(x, TBlock)
	}
}

func inferMemRead(x *ast.Mem, env TypeEnv) {
	id, ok := x.E.(*ast.Ident)
	if !ok {
		return
	}
	name, ok := memberName(x.M)
	if !ok {
		return
	}
	switch {
	case id.Name == "this":
		if t, ok := env.LookupMember(name); ok {
			env.SetType(x, t)
		}
	case isGlobalIdent(id.Name):
		if t := env.ClassStaticType(id.Name, name); t != TUnknown {
			env.SetType(x, t)
		}
	}
}

const maxLoopIterations = 4

func walkLoopFixpoint(sc scope, walkBody func(scope)) {
	cur := cloneScope(sc)
	for range maxLoopIterations {
		body := cloneScope(cur)
		walkBody(body)
		next := cloneScope(sc)
		joinScopeInto(next, body)
		if scopesEqual(next, cur) {
			cur = next
			break
		}
		cur = next
	}
	for name, t := range cur {
		sc[name] = t
	}
}

func joinScopeInto(dst, src scope) {
	for name, t := range src {
		if dt, ok := dst[name]; ok {
			dst[name] = U(dt, t)
		} else {
			dst[name] = t
		}
	}
}

func scopesEqual(a, b scope) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok || !dynEqual(v, bv) {
			return false
		}
	}
	return true
}

func cloneScope(sc scope) scope {
	out := make(scope, len(sc))
	for k, v := range sc {
		out[k] = v
	}
	return out
}

func mergeScopesN(dst scope, branches []scope) {
	if len(branches) == 0 {
		return
	}
	seen := map[string]bool{}
	for _, b := range branches {
		for name := range b {
			seen[name] = true
		}
	}
	for name := range seen {
		var merged DynType
		for _, b := range branches {
			t, ok := b[name]
			if !ok {
				continue
			}
			if merged == nil {
				merged = t
			} else {
				merged = U(merged, t)
			}
		}
		if merged != nil {
			dst[name] = merged
		}
	}
}

func startsUpper(name string) bool {
	if name == "" {
		return false
	}
	c := name[0]
	return c >= 'A' && c <= 'Z'
}

func isGlobalIdent(name string) bool { return startsUpper(name) }
