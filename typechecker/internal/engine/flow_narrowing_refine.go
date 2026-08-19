// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

func isNarrower(a, b DynType) bool {
	if a == nil || b == nil {
		return false
	}
	aP, aIsP := a.(Primitive)
	bP, bIsP := b.(Primitive)
	if aIsP && bIsP && aP == bP {
		return false
	}
	if bIsP && bP == TUnknown {
		return !aIsP || aP != TUnknown
	}
	if aIsP && aP == TUnknown {
		return false
	}
	// any clean type is strictly narrower than any dirty one: dirty admits
	// unknown values, without this an assignment whose RHS
	// narrowing cleaned (e.g. a ternary arm reading a guarded param) cannot
	// rescue the local from its stale pre-narrowing dirty flow stamp
	if _, bDirty := decomposeForCheck(b); bDirty {
		if _, aDirty := decomposeForCheck(a); !aDirty {
			return true
		}
	}
	aU, aIsU := a.(Union)
	bU, bIsU := b.(Union)
	switch {
	case !aIsU && !bIsU:
		return subtypeOf(a, b) && a != b
	case !aIsU && bIsU:
		for _, t := range bU.Types {
			if a == t || subtypeOf(a, t) {
				return true
			}
		}
		return false
	case aIsU && !bIsU:
		return false
	default: // both Unions
		for _, ta := range aU.Types {
			found := false
			for _, tb := range bU.Types {
				if ta == tb || subtypeOf(ta, tb) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		switch {
		case len(aU.Types) < len(bU.Types):
			return true
		case len(aU.Types) == len(bU.Types) && !aU.IsDirty && bU.IsDirty:
			return true
		}
		return false
	}
}

// TFalse/TTrue <: TBoolean, TSequence <: TObject.
func subtypeOf(sub, sup DynType) bool {
	subP, subOk := sub.(Primitive)
	supP, supOk := sup.(Primitive)
	if !subOk || !supOk {
		return false
	}
	if subP == supP {
		return true
	}
	if supP == TBoolean {
		return subP == TTrue || subP == TFalse
	}
	if supP == TObject {
		return subP == TSequence || subP == TClass
	}
	return false
}

func refineCond(cond ast.Expr, sc narrowScope, polarity bool, env TypeEnv, allowMembers bool) narrowScope {
	out := sc.clone()
	applyRefinement(cond, out, polarity, env, allowMembers)
	return out
}

func applyRefinement(cond ast.Expr, sc narrowScope, polarity bool, env TypeEnv, allowMembers bool) {
	if ep, ok := cond.(*ast.ExprPos); ok && ep.Expr != nil {
		applyRefinement(ep.Expr, sc, polarity, env, allowMembers)
		return
	}
	switch n := cond.(type) {
	case *ast.Unary:
		switch n.Tok {
		case tok.Not:
			applyRefinement(n.E, sc, !polarity, env, allowMembers)
		case tok.LParen:
			applyRefinement(n.E, sc, polarity, env, allowMembers)
		}
	case *ast.Nary:
		switch {
		case (n.Tok == tok.And && polarity) || (n.Tok == tok.Or && !polarity):
			for _, e := range n.Exprs {
				applyRefinement(e, sc, polarity, env, allowMembers)
			}
		case (n.Tok == tok.Or && polarity) || (n.Tok == tok.And && !polarity):
			mergeForkedRefinements(n.Exprs, sc, polarity, env, allowMembers)
		}
	case *ast.Binary:
		refineBinary(n, sc, polarity, env, allowMembers)
	case *ast.Call:
		refinePredicateCall(n, sc, polarity, env, allowMembers)
		refineHelperPostcondition(n, sc, polarity, allowMembers)
	}
}

func mergeForkedRefinements(exprs []ast.Expr, sc narrowScope, polarity bool, env TypeEnv, allowMembers bool) {
	if len(exprs) == 0 {
		return
	}
	perOp := make([]narrowScope, len(exprs))
	for i, e := range exprs {
		saved := overrideEnvWithScopeBaseline(e, sc, env)
		perOp[i] = sc.clone()
		applyRefinement(e, perOp[i], polarity, env, allowMembers)
		restoreEnvStamps(saved, env)
	}
	mergeForkedKind(sc.Locals, perOp, false)
	if allowMembers {
		mergeForkedKind(sc.Members, perOp, true)
	}
}

type savedNode struct {
	node    ast.Node
	present bool
	ty      DynType
}

// the chain walk stamps speculative types into env.Nodes; restoreEnvStamps
// must undo them or narrowed types leak outside the guard.
func overrideEnvWithScopeBaseline(e ast.Expr, sc narrowScope, env TypeEnv) []savedNode {
	var saved []savedNode
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		switch x := n.(type) {
		case *ast.Ident:
			if !isGlobalIdent(x.Name) {
				if t, ok := sc.Locals.typ(x.Name); ok {
					prev, present := env.Nodes[n]
					saved = append(saved, savedNode{node: n, present: present, ty: prev})
					env.Nodes[n] = t
				}
			}
		case *ast.Mem:
			if name, mem, ok := unwrapThisMember(x); ok {
				if t, ok2 := sc.Members.typ(name); ok2 {
					prev, present := env.Nodes[mem]
					saved = append(saved, savedNode{node: mem, present: present, ty: prev})
					env.Nodes[mem] = t
				}
			}
		}
		n.Children(func(c ast.Node) ast.Node {
			walk(c)
			return c
		})
	}
	walk(e)
	return saved
}

// undoes overrideEnvWithScopeBaseline. skipping this looks safe and is not.
func restoreEnvStamps(saved []savedNode, env TypeEnv) {
	for _, s := range saved {
		if s.present {
			env.Nodes[s.node] = s.ty
		} else {
			delete(env.Nodes, s.node)
		}
	}
}

func mergeForkedKind(r refinements, perOp []narrowScope, members bool) {
	cand := map[string]bool{}
	for _, op := range perOp {
		for name, f := range op.kind(members) {
			if f.InGuard {
				cand[name] = true
			}
		}
	}
	for name := range cand {
		merged, every := mergeForkedName(name, r, perOp, members)
		if !every || merged == nil {
			continue
		}
		r.prove(name, merged)
	}
}

func mergeForkedName(name string, r refinements,
	perOp []narrowScope, members bool) (DynType, bool) {
	var merged DynType
	cur := r[name]
	for _, op := range perOp {
		opF, ok := op.kind(members)[name]
		if !ok || !opF.InGuard {
			return nil, false
		}
		if cur.InGuard && dynEqual(opF.Typ, cur.Typ) {
			return nil, false
		}
		if merged == nil {
			merged = opF.Typ
		} else {
			merged = U(merged, opF.Typ)
		}
	}
	return merged, true
}

func existingType(sc narrowScope, tgt narrowTarget, env TypeEnv) DynType {
	if tgt.isMember {
		if t, ok := sc.Members.guarded(tgt.name); ok {
			return t
		}
		if t := env.GetType(tgt.node); t != TUnknown {
			return t
		}
		if t, ok := env.LookupMember(tgt.name); ok {
			return t
		}
		return TUnknown
	}
	if t, ok := sc.Locals.guarded(tgt.name); ok {
		return t
	}
	return env.GetType(tgt.node)
}

// ```suneido
// if x is 5    { return x }    // x narrows toward TNumber inside { ... }
//
//	^^^^
//
// if x isnt false { return x } // TFalse dropped from x's type
//
//	^^^^^^^^^
//
// ```
func refineBinary(b *ast.Binary, sc narrowScope, polarity bool, env TypeEnv, allowMembers bool) {
	switch b.Tok {
	case tok.Is:
		applyEqRefinement(b, sc, polarity, env, allowMembers)
	case tok.Isnt:
		applyEqRefinement(b, sc, !polarity, env, allowMembers)
	}
}

func applyEqRefinement(b *ast.Binary, sc narrowScope, equal bool, env TypeEnv, allowMembers bool) {
	if tgt, lit, ok := targetAndLiteral(b); ok {
		t := DynTypeOfSuValue(lit.Val)
		existing := existingType(sc, tgt, env)
		if equal {
			storeRefinement(sc, tgt, narrowTowardSet(existing, []DynType{t}), allowMembers)
		} else if t == TFalse || t == TTrue {
			storeRefinement(sc, tgt, narrowAwaySet(existing, []DynType{t}), allowMembers)
		}
		return
	}
	refineTypeStringEq(b, sc, equal, env, allowMembers)
}

func storeRefinement(sc narrowScope, tgt narrowTarget, t DynType, allowMembers bool) {
	if t == nil {
		return
	}
	if tgt.isMember {
		if !allowMembers {
			return
		}
		sc.Members.prove(tgt.name, t)
		return
	}
	sc.Locals.prove(tgt.name, t)
}

func targetAndLiteral(b *ast.Binary) (narrowTarget, *ast.Constant, bool) {
	if tgt, ok := unwrapTarget(b.Lhs); ok {
		if c, ok := unwrapConstant(b.Rhs); ok {
			return tgt, c, true
		}
	}
	if tgt, ok := unwrapTarget(b.Rhs); ok {
		if c, ok := unwrapConstant(b.Lhs); ok {
			return tgt, c, true
		}
	}
	return narrowTarget{}, nil, false
}

// peels ExprPos, parens, and inline assignments to the underlying name.
//
// ```suneido
// if false is x = .Maybe()    return  // RHS is Binary(Eq), not Ident
//
//	^^^^^^^^^^^^^^^         // we still need to refine on x
//
// if ((x = .Maybe()) is false) ...    // LHS wrapped in Unary(LParen)
//
//	^^^^^^^^^^^^^^^                    // peel paren AND assign to find x
//
// ```
func unwrapIdent(e ast.Expr) (*ast.Ident, bool) {
	for {
		switch x := e.(type) {
		case *ast.ExprPos:
			if x.Expr == nil {
				return nil, false
			}
			e = x.Expr
		case *ast.Unary:
			if x.Tok != tok.LParen {
				return nil, false
			}
			e = x.E
		case *ast.Binary:
			if x.Tok != tok.Eq {
				return nil, false
			}
			e = x.Lhs
		case *ast.Ident:
			return x, true
		default:
			return nil, false
		}
	}
}

func unwrapConstant(e ast.Expr) (*ast.Constant, bool) {
	if ep, ok := e.(*ast.ExprPos); ok && ep.Expr != nil {
		e = ep.Expr
	}
	c, ok := e.(*ast.Constant)
	return c, ok
}

func unwrapCall(e ast.Expr) (*ast.Call, bool) {
	if ep, ok := e.(*ast.ExprPos); ok && ep.Expr != nil {
		e = ep.Expr
	}
	c, ok := e.(*ast.Call)
	return c, ok
}

// ```suneido
// if Type(x) is "Number"   { ... } // x narrows toward TNumber
//
//	^      ^^^^^^^^
//
// if Type(.foo) isnt "String" { ... } // TString removed from .foo
//
//	^^^^       ^^^^^^^^
//
// ```
// x may be a local Ident or a `this`-member Mem.
func refineTypeStringEq(b *ast.Binary, sc narrowScope, equal bool, env TypeEnv, allowMembers bool) {
	tryPair := func(callExpr, strExpr ast.Expr) bool {
		call, ok := unwrapCall(callExpr)
		if !ok {
			return false
		}
		fnId, ok := unwrapIdent(call.Fn)
		if !ok || fnId.Name != "Type" || len(call.Args) != 1 {
			return false
		}
		tgt, ok := unwrapTarget(call.Args[0].E)
		if !ok {
			return false
		}
		c, ok := unwrapConstant(strExpr)
		if !ok {
			return false
		}
		s, ok := c.Val.ToStr()
		if !ok {
			return false
		}
		targets, ok := typeStringTargets[s]
		if !ok {
			return false
		}
		existing := existingType(sc, tgt, env)
		if equal {
			storeRefinement(sc, tgt, narrowTowardSet(existing, targets), allowMembers)
		} else {
			storeRefinement(sc, tgt, narrowAwaySet(existing, targets), allowMembers)
		}
		return true
	}
	if tryPair(b.Lhs, b.Rhs) {
		return
	}
	tryPair(b.Rhs, b.Lhs)
}

// ```suneido
// if Number?(x)   { return x } // x narrows to TNumber inside the body
//
//	^^^^^^^^^^                  // polarity=true  -> narrowTowardSet
//
// if not String?(x) return     // post-this-stmt x narrows to TString
//
//	^^^^^^^^^^^^^^              // polarity=false (from outer `not`)
//
// ```
func refinePredicateCall(call *ast.Call, sc narrowScope, polarity bool, env TypeEnv, allowMembers bool) {
	fn, ok := unwrapIdent(call.Fn)
	if !ok {
		return
	}
	targets, ok := predicateTargets[fn.Name]
	if !ok || len(call.Args) != 1 {
		return
	}
	tgt, ok := unwrapTarget(call.Args[0].E)
	if !ok {
		return
	}
	existing := existingType(sc, tgt, env)
	if polarity {
		storeRefinement(sc, tgt, narrowTowardSet(existing, targets), allowMembers)
	} else {
		storeRefinement(sc, tgt, narrowAwaySet(existing, targets), allowMembers)
	}
}

func boolIntersect(targets []DynType) DynType {
	keepF, keepT := false, false
	for _, tgt := range targets {
		switch tgt {
		case TFalse:
			keepF = true
		case TTrue:
			keepT = true
		case TBoolean:
			keepF, keepT = true, true
		}
	}
	switch {
	case keepF && keepT:
		return TBoolean
	case keepF:
		return TFalse
	case keepT:
		return TTrue
	}
	return TVoid
}

// dual of boolIntersect removing [TFalse] leaves TTrue, etc.
func boolComplement(targets []DynType) DynType {
	removeF, removeT := false, false
	for _, tgt := range targets {
		switch tgt {
		case TFalse:
			removeF = true
		case TTrue:
			removeT = true
		case TBoolean:
			removeF, removeT = true, true
		}
	}
	switch {
	case removeF && removeT:
		return TVoid
	case removeF:
		return TTrue
	case removeT:
		return TFalse
	}
	return TBoolean
}

func narrowTowardSet(existing DynType, targets []DynType) DynType {
	if existing == nil || existing == TUnknown {
		return foldTargetUnion(targets)
	}
	if u, ok := existing.(Union); ok {
		return narrowUnionToward(u, targets)
	}
	if existing == TBoolean {
		if bi := boolIntersect(targets); bi != TVoid {
			return bi
		}
		return foldTargetUnion(targets)
	}
	for _, target := range targets {
		if subtypeOf(existing, target) {
			return existing
		}
	}
	return foldTargetUnion(targets)
}

func narrowUnionToward(u Union, targets []DynType) DynType {
	var kept []DynType
	for _, t := range u.Types {
		if t == TBoolean {
			if bi := boolIntersect(targets); bi != TVoid {
				kept = append(kept, bi)
			}
			continue
		}
		for _, target := range targets {
			if subtypeOf(t, target) {
				kept = append(kept, t)
				break
			}
		}
	}
	if u.IsDirty {
		kept = append(kept, targets...)
	}
	if len(kept) == 0 {
		return foldTargetUnion(targets)
	}
	return Union{Types: kept}.Fold()
}

func isAtomicType(t DynType) bool {
	p, ok := t.(Primitive)
	if !ok {
		return false
	}
	return p == TFalse || p == TTrue || p == TBoolean
}

func narrowAwaySet(existing DynType, targets []DynType) DynType {
	if existing == nil || existing == TUnknown {
		return existing
	}
	matches := func(t DynType) bool {
		for _, target := range targets {
			if subtypeOf(t, target) {
				return true
			}
		}
		return false
	}
	if u, ok := existing.(Union); ok {
		return narrowUnionAway(u, targets, matches)
	}
	if existing == TBoolean {
		if c := boolComplement(targets); c != TVoid {
			return c
		}
		return TUnknown
	}
	if matches(existing) {
		if isAtomicType(existing) {
			return TUnknown
		}
		return existing
	}
	return existing
}

func narrowUnionAway(u Union, targets []DynType, matches func(DynType) bool) DynType {
	var kept []DynType
	allAtomicRemoved := true
	for _, t := range u.Types {
		if t == TBoolean {
			if c := boolComplement(targets); c != TVoid {
				kept = append(kept, c)
			}
			continue
		}
		if !matches(t) {
			kept = append(kept, t)
			continue
		}
		if !isAtomicType(t) {
			allAtomicRemoved = false
		}
	}
	if len(kept) == 0 {
		if allAtomicRemoved {
			return TUnknown
		}
		return u
	}
	return Union{Types: kept, IsDirty: u.IsDirty}.Fold()
}

func foldTargetUnion(targets []DynType) DynType {
	if len(targets) == 0 {
		return TUnknown
	}
	if len(targets) == 1 {
		return targets[0]
	}
	return Union{Types: targets}.Fold()
}
