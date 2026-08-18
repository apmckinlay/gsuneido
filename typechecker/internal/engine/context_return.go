// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

type condVerdict int

const (
	condUnknown condVerdict = iota
	condTrue
	condFalse
)

func mustBeFalse(t DynType) bool { return t == TFalse }

func cannotBeFalse(t DynType) bool {
	switch x := t.(type) {
	case Primitive:
		switch x {
		case TUnknown, TVoid, TFalse, TBoolean:
			return false
		default:
			return true
		}
	case Instance:
		return true
	case Union:
		if x.IsDirty {
			return false
		}
		for _, m := range x.Types {
			if !cannotBeFalse(m) {
				return false
			}
		}
		return len(x.Types) > 0
	}
	return false
}

func unwrapBinary(e ast.Expr) (*ast.Binary, bool) {
	for {
		switch x := e.(type) {
		case *ast.ExprPos:
			e = x.Expr
		case *ast.Binary:
			return x, true
		default:
			return nil, false
		}
	}
}

// matches `ident is false` / `ident isnt false` (either operand order).
func identFalseTest(b *ast.Binary) (id *ast.Ident, negated, ok bool) {
	switch b.Tok {
	case tok.Is:
		negated = false
	case tok.Isnt:
		negated = true
	default:
		return nil, false, false
	}
	if id, ok := unwrapIdent(b.Lhs); ok {
		if c, ok := unwrapConstant(b.Rhs); ok && DynTypeOfSuValue(c.Val) == TFalse {
			return id, negated, true
		}
	}
	if id, ok := unwrapIdent(b.Rhs); ok {
		if c, ok := unwrapConstant(b.Lhs); ok && DynTypeOfSuValue(c.Val) == TFalse {
			return id, negated, true
		}
	}
	return nil, false, false
}

func decideCond(cond ast.Expr, overrides map[string]DynType, env TypeEnv) condVerdict {
	if c, ok := unwrapConstant(peelParens(cond)); ok {
		switch DynTypeOfSuValue(c.Val) {
		case TTrue:
			return condTrue
		case TFalse:
			return condFalse
		}
	}
	b, ok := unwrapBinary(cond)
	if !ok {
		return condUnknown
	}
	id, negated, ok := identFalseTest(b)
	if !ok {
		return condUnknown
	}
	t := env.GetType(id)
	if overrides != nil {
		if ot, ok := overrides[id.Name]; ok {
			t = ot
		}
	}
	var v condVerdict
	switch {
	case mustBeFalse(t):
		v = condTrue
	case cannotBeFalse(t):
		v = condFalse
	default:
		return condUnknown
	}
	if negated {
		if v == condTrue {
			return condFalse
		}
		return condTrue
	}
	return v
}

func addReturnType(r *ast.Return, env TypeEnv, acc *DynType) {
	if r.ReturnThrow {
		return
	}
	var ty DynType
	if len(r.Exprs) == 0 {
		ty = TVoid
	} else {
		ty = env.GetType(r.Exprs[0])
		for _, e := range r.Exprs[1:] {
			ty = U(ty, env.GetType(e))
		}
	}
	if *acc == nil {
		*acc = ty
	} else {
		*acc = U(*acc, ty)
	}
}

func reachStmt(n ast.Node, overrides map[string]DynType, env TypeEnv, acc *DynType) bool {
	if n == nil {
		return false
	}
	switch x := n.(type) {
	case *ast.Return:
		addReturnType(x, env, acc)
		return true
	case *ast.Compound:
		return reachStmtList(x.Body, overrides, env, acc)
	case *ast.If:
		switch decideCond(x.Cond, overrides, env) {
		case condTrue:
			return reachStmt(x.Then, overrides, env, acc)
		case condFalse:
			return reachStmt(x.Else, overrides, env, acc)
		default:
			tThen := reachStmt(x.Then, overrides, env, acc)
			tElse := reachStmt(x.Else, overrides, env, acc)
			return tThen && tElse
		}
	default:
		collectReturns(n, env, acc)
		return false
	}
}

func reachStmtList(stmts []ast.Statement, overrides map[string]DynType, env TypeEnv, acc *DynType) bool {
	for _, s := range stmts {
		if s != nil && reachStmt(s, overrides, env, acc) {
			return true
		}
	}
	return false
}

func inferBodyReturnReachable(fn *ast.Function, overrides map[string]DynType, env TypeEnv) DynType {
	var result DynType
	if !reachStmtList(fn.Body, overrides, env, &result) {
		if es, ok := lastExprStmt(fn.Body); ok {
			result = mergeReturn(result, env.GetType(es.E))
		}
	}
	if u, ok := result.(Union); ok {
		return u.Fold()
	}
	if result != nil {
		return result
	}
	return TVoid
}

func lastExprStmt(stmts []ast.Statement) (*ast.ExprStmt, bool) {
	if n := len(stmts); n > 0 {
		if es, ok := stmts[n-1].(*ast.ExprStmt); ok && es.E != nil {
			return es, true
		}
	}
	return nil, false
}

func mergeReturn(acc, ty DynType) DynType {
	if acc == nil {
		return ty
	}
	return U(acc, ty)
}

func freeParamOverrides(fn *ast.Function) map[string]DynType {
	ov := make(map[string]DynType, len(fn.Params))
	for i := range fn.Params {
		ov[fn.Params[i].Name.ParamName()] = TUnknown
	}
	return ov
}

func bindArgs(fn *ast.Function, args []ast.Arg, env TypeEnv) map[string]DynType {
	ov := make(map[string]DynType, len(fn.Params))
	for i := range fn.Params {
		p := &fn.Params[i]
		if p.DefVal != nil {
			ov[p.Name.ParamName()] = DynTypeOfSuValue(p.DefVal)
		} else {
			ov[p.Name.ParamName()] = TUnknown
		}
	}
	pos := 0
	for i := range args {
		arg := &args[i]
		if isAtArg(arg) {
			return nil
		}
		if arg.Name != nil {
			if name, ok := argName(arg); ok {
				ov[name] = env.GetType(arg.E)
			}
			continue
		}
		if pos < len(fn.Params) {
			ov[fn.Params[pos].Name.ParamName()] = env.GetType(arg.E)
		}
		pos++
	}
	return ov
}

func contextReturn(fn *ast.Function, args []ast.Arg, env TypeEnv) (DynType, bool) {
	ov := bindArgs(fn, args, env)
	if ov == nil {
		return nil, false
	}
	r := inferBodyReturnReachable(fn, ov, env)
	if r == nil || r == TVoid || r == TUnknown {
		return nil, false
	}
	return r, true
}

type ReturnSummary struct {
	Param      string  // discriminating param name ("" = none)
	ParamIndex int     // its positional index, for positional-arg binding
	DefaultArg DynType // its type when the arg is omitted (the param's default)
	IfFalse    DynType // return when the arg is provably false
	IfTrue     DynType // return when the arg is provably not false
}

func computeSummary(fn *ast.Function, env TypeEnv) ReturnSummary {
	free := freeParamOverrides(fn)
	for i := range fn.Params {
		p := &fn.Params[i]
		name := p.Name.ParamName()
		ovF := cloneTypeMap(free)
		ovF[name] = TFalse
		ovT := cloneTypeMap(free)
		ovT[name] = TNumber // any provably-not-false type drives the not-false arm
		rF := inferBodyReturnReachable(fn, ovF, env)
		rT := inferBodyReturnReachable(fn, ovT, env)
		if dynEqual(rF, rT) {
			continue
		}
		var def DynType = TUnknown
		if p.DefVal != nil {
			def = DynTypeOfSuValue(p.DefVal)
		}
		return ReturnSummary{Param: name, ParamIndex: i, DefaultArg: def, IfFalse: rF, IfTrue: rT}
	}
	return ReturnSummary{}
}

// ```suneido
// Get(x = false) { return x is false ? "default" : x }
// .Get()
//
//	^^^^^ discriminator x known false -> String
//
// .Get(5)
//
//	^^^^^^^ non-false arm -> the argument's type
//
// ```
func ComputeSummaries(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	if env.Summaries == nil {
		return false
	}
	for name, fn := range cls.SortedMethods {
		if s := computeSummary(fn, env); s.Param != "" {
			env.Summaries[name] = s
		}
	}
	return false
}

func resolveDiscriminatorType(s ReturnSummary, args []ast.Arg, env TypeEnv) (DynType, bool) {
	positional := 0
	for i := range args {
		arg := &args[i]
		if isAtArg(arg) {
			return nil, false
		}
		if arg.Name != nil {
			if name, ok := argName(arg); ok && name == s.Param {
				return env.GetType(arg.E), true
			}
			continue
		}
		if positional == s.ParamIndex {
			return env.GetType(arg.E), true
		}
		positional++
	}
	return s.DefaultArg, true
}

func summaryReturn(env TypeEnv, class, method string, args []ast.Arg) (DynType, bool) {
	if env.ClassSummaries == nil {
		return nil, false
	}
	m, ok := env.ClassSummaries[class]
	if !ok {
		return nil, false
	}
	s, ok := m[method]
	if !ok || s.Param == "" {
		return nil, false
	}
	dt, ok := resolveDiscriminatorType(s, args, env)
	if !ok {
		return nil, false
	}
	switch {
	case mustBeFalse(dt):
		return s.IfFalse, true
	case cannotBeFalse(dt):
		return s.IfTrue, true
	}
	return nil, false
}
