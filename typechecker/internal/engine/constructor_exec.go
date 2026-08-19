// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// ```suneido
//
//	class {
//		x: false
//		New() { .x = 5 }
//		Get() { return .x }
//		                ^^ post-New (instance) view: Number - the False seed arm
//	}                    is demoted because New provably always assigns non-false
//
// ```
func ConstructorExecPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	newFn, ok := cls.Methods["New"]
	if !ok {
		return false
	}

	seedLits := seedLiterals(cls)

	exec := &ctorExec{env: env, cls: cls, stack: map[string]bool{"New": true}}
	st := &ctorState{members: cloneTypeMap(seedLits), locals: map[string]DynType{}}
	exec.bindParams(newFn, nil, st, st)
	fell, exits := exec.runBody(newFn.Body, st)
	if exec.bail {
		return false
	}
	states := exits
	if fell {
		states = append(states, st.members)
	}
	if len(states) == 0 {
		return false
	}
	postNew := joinMembers(states)

	falseAssigned := membersAssignedFalse(cls, env)

	demotions := computeDemotions(seedLits, postNew, falseAssigned, env.Members)
	if len(demotions) == 0 {
		return false
	}
	for name, t := range demotions {
		env.PostCtorMembers[name] = t
	}
	restampMemberReads(cls, env, demotions)
	return true
}

func seedLiterals(cls *ClassObject) map[string]DynType {
	seedLits := map[string]DynType{}
	for name, val := range cls.Members {
		t := DynTypeOfSuValue(val)
		if isPublicMember(name) {
			t = widenBool(t)
		}
		seedLits[name] = t
	}
	return seedLits
}

func computeDemotions(seedLits, postNew map[string]DynType,
	falseAssigned map[string]bool, members map[string]DynType) map[string]DynType {
	demotions := map[string]DynType{}
	for name, lit := range seedLits {
		if isPublicMember(name) {
			continue // externally writable - New-init isn't a class-wide guarantee
		}
		if !containsFalse(lit) {
			continue // only `: false` sentinels are candidates
		}
		seed, ok := members[name]
		if !ok || !containsFalse(seed) {
			continue
		}
		pn, ok := postNew[name]
		if !ok || pn == TUnknown || containsFalse(pn) {
			continue // C1: New must definitely overwrite the seed non-false
		}
		if falseAssigned[name] {
			continue // C2: some method puts false back
		}
		demoted := removeFalse(seed)
		if demoted == TUnknown {
			continue
		}
		demotions[name] = demoted
	}
	return demotions
}

type ctorState struct {
	members map[string]DynType
	locals  map[string]DynType
}

func (st *ctorState) clone() *ctorState {
	return &ctorState{members: cloneTypeMap(st.members), locals: cloneTypeMap(st.locals)}
}

type ctorExec struct {
	env   TypeEnv
	cls   *ClassObject
	stack map[string]bool // methods on the current call stack (recursion guard)
	bail  bool            // dynamic member write seen - abort demotion entirely
}

func (e *ctorExec) runBody(stmts []ast.Statement, st *ctorState) (bool, []map[string]DynType) {
	var exits []map[string]DynType
	for _, s := range stmts {
		if e.bail {
			return false, exits
		}
		fell, ex := e.runStmt(s, st)
		exits = append(exits, ex...)
		if !fell {
			return false, exits // unconditional exit - rest is dead
		}
	}
	return true, exits
}

func (e *ctorExec) runStmt(s ast.Statement, st *ctorState) (bool, []map[string]DynType) {
	switch n := s.(type) {
	case *ast.Compound:
		return e.runBody(n.Body, st)
	case *ast.Return:
		if n.ReturnThrow {
			return false, nil // divergence
		}
		return false, []map[string]DynType{cloneTypeMap(st.members)}
	case *ast.Throw:
		return false, nil
	case *ast.ExprStmt:
		e.runExpr(n.E, st)
		return true, nil
	case *ast.If:
		return e.runIf(n, st)
	case *ast.Switch:
		return e.runSwitch(n, st)
	case *ast.ForIn:
		return e.runLoop(n.Body, st)
	case *ast.For:
		return e.runLoop(n.Body, st)
	case *ast.While:
		return e.runLoop(n.Body, st)
	case *ast.DoWhile:
		return e.runLoop(n.Body, st)
	case *ast.Forever:
		return e.runLoop(n.Body, st)
	case *ast.TryCatch:
		return e.runTryCatch(n, st)
	default:
		return true, nil
	}
}

func (e *ctorExec) runIf(n *ast.If, st *ctorState) (bool, []map[string]DynType) {
	feas := e.feasibility(n.Cond, st)
	var exits []map[string]DynType
	var fellStates []map[string]DynType

	if feas >= 0 { // then reachable (unknown or provably true)
		thenSt := st.clone()
		fell, ex := e.runBody(stmtList(n.Then), thenSt)
		exits = append(exits, ex...)
		if fell {
			fellStates = append(fellStates, thenSt.members)
		}
	}
	if feas <= 0 { // else reachable (unknown or provably false)
		elseSt := st.clone()
		if n.Else != nil {
			fell, ex := e.runBody(stmtList(n.Else), elseSt)
			exits = append(exits, ex...)
			if fell {
				fellStates = append(fellStates, elseSt.members)
			}
		} else {
			fellStates = append(fellStates, elseSt.members) // implicit empty else
		}
	}
	if len(fellStates) == 0 {
		return false, exits // both arms exit
	}
	st.members = joinMembers(fellStates)
	return true, exits
}

func (e *ctorExec) runSwitch(n *ast.Switch, st *ctorState) (bool, []map[string]DynType) {
	var exits []map[string]DynType
	var fellStates []map[string]DynType
	for i := range n.Cases {
		caseSt := st.clone()
		fell, ex := e.runBody(n.Cases[i].Body, caseSt)
		exits = append(exits, ex...)
		if fell {
			fellStates = append(fellStates, caseSt.members)
		}
	}
	if n.Default != nil {
		defSt := st.clone()
		fell, ex := e.runBody(n.Default, defSt)
		exits = append(exits, ex...)
		if fell {
			fellStates = append(fellStates, defSt.members)
		}
	} else {
		fellStates = append(fellStates, cloneTypeMap(st.members)) // no-match path
	}
	if len(fellStates) == 0 {
		return false, exits
	}
	st.members = joinMembers(fellStates)
	return true, exits
}

func (e *ctorExec) runLoop(body ast.Statement, st *ctorState) (bool, []map[string]DynType) {
	bodySt := st.clone()
	_, exits := e.runBody(stmtList(body), bodySt)
	st.members = joinMembers([]map[string]DynType{st.members, bodySt.members})
	return true, exits
}

func (e *ctorExec) runTryCatch(n *ast.TryCatch, st *ctorState) (bool, []map[string]DynType) {
	var exits []map[string]DynType
	var fellStates []map[string]DynType
	trySt := st.clone()
	tfell, tex := e.runBody(stmtList(n.Try), trySt)
	exits = append(exits, tex...)
	if tfell {
		fellStates = append(fellStates, trySt.members)
	}
	if n.Catch != nil {
		cSt := st.clone()
		cfell, cex := e.runBody(stmtList(n.Catch), cSt)
		exits = append(exits, cex...)
		if cfell {
			fellStates = append(fellStates, cSt.members)
		}
	} else {
		fellStates = append(fellStates, cloneTypeMap(st.members))
	}
	if len(fellStates) == 0 {
		return false, exits
	}
	st.members = joinMembers(fellStates)
	return true, exits
}

func (e *ctorExec) runExpr(expr ast.Expr, st *ctorState) {
	expr = ctorPeel(expr)
	switch x := expr.(type) {
	case *ast.Binary:
		if x.Tok == tok.Eq {
			e.runAssign(x.Lhs, x.Rhs, st)
			return
		}
		if x.Tok.IsAssign() { // compound: result type isn't trivially the RHS
			e.clearMember(x.Lhs, st)
			return
		}
	case *ast.Unary:
		switch x.Tok {
		case tok.Inc, tok.Dec, tok.PostInc, tok.PostDec:
			e.clearMember(x.E, st)
			return
		}
	case *ast.Call:
		e.runCall(x, st)
	}
}

func (e *ctorExec) runAssign(lhs, rhs ast.Expr, st *ctorState) {
	if name, _, ok := unwrapThisMember(lhs); ok {
		st.members[name] = e.evalType(rhs, st)
		return
	}
	peeled := ctorPeel(lhs)
	if mem, ok := peeled.(*ast.Mem); ok {
		if id, ok := mem.E.(*ast.Ident); ok && id.Name == "this" {
			// this[expr] = ... : dynamic member name, can't track which -> bail.
			e.bail = true
			return
		}
		return // .obj.y = ... : mutates obj, doesn't rebind a this-member
	}
	if id, ok := peeled.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		st.locals[id.Name] = e.evalType(rhs, st)
	}
}

func (e *ctorExec) clearMember(lhs ast.Expr, st *ctorState) {
	if name, _, ok := unwrapThisMember(lhs); ok {
		st.members[name] = TUnknown
	}
}

func (e *ctorExec) runCall(call *ast.Call, st *ctorState) {
	fn := ctorPeel(call.Fn)
	mem, ok := fn.(*ast.Mem)
	if !ok {
		return // free call - can't rebind this-members (closed world)
	}
	id, ok := mem.E.(*ast.Ident)
	if !ok || id.Name != "this" {
		return // foreign receiver / super - not followed (leaf-local)
	}
	name, ok := memberName(mem.M)
	if !ok {
		return
	}
	e.stepIn(name, call.Args, st)
}

func (e *ctorExec) stepIn(name string, args []ast.Arg, st *ctorState) {
	if e.stack[name] {
		return // recursion - cut conservatively (no proven effect)
	}
	fn, ok := e.cls.Methods[name]
	if !ok {
		return // builtin / inherited - guessed not to rebind privates
	}
	e.stack[name] = true
	defer delete(e.stack, name)

	frame := &ctorState{members: cloneTypeMap(st.members), locals: map[string]DynType{}}
	e.bindParams(fn, args, st, frame)
	fell, exits := e.runBody(fn.Body, frame)
	if e.bail {
		return
	}
	states := exits
	if fell {
		states = append(states, frame.members)
	}
	if len(states) == 0 {
		return // helper always diverges - leave caller state unchanged
	}
	st.members = joinMembers(states)
}

func (e *ctorExec) bindParams(fn *ast.Function, args []ast.Arg, callerSt, frame *ctorState) {
	for i := range fn.Params {
		p := &fn.Params[i]
		pname := p.Name.ParamName()
		var val DynType = TUnknown
		if args != nil && i < len(args) && args[i].Name == nil {
			val = e.evalType(args[i].E, callerSt)
		} else if p.DefVal != nil {
			val = DynTypeOfSuValue(p.DefVal)
		}
		if len(p.Name.Name) > 0 && p.Name.Name[0] == '.' {
			val = widenBool(val) // mirror LocalInference for `.x` params
			frame.members[pname] = val
		}
		frame.locals[pname] = val
	}
}

func (e *ctorExec) evalType(expr ast.Expr, st *ctorState) DynType {
	if name, _, ok := unwrapThisMember(expr); ok {
		if t, ok := st.members[name]; ok {
			return t
		}
		if t, ok := e.env.LookupMember(name); ok {
			return t
		}
		return TUnknown
	}
	peeled := ctorPeel(expr)
	if id, ok := peeled.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		if t, ok := st.locals[id.Name]; ok {
			return t
		}
	}
	if tri, ok := peeled.(*ast.Trinary); ok {
		return e.evalTrinary(tri, st)
	}
	if t := e.env.GetType(expr); t != TUnknown {
		return t
	}
	if t := e.env.GetType(peeled); t != TUnknown {
		return t
	}
	return TUnknown
}

func (e *ctorExec) evalTrinary(tri *ast.Trinary, st *ctorState) DynType {
	tType := e.evalType(tri.T, st)
	fType := e.evalType(tri.F, st)
	tEmpty, fEmpty := false, false
	if tgt, ok := ctorTarget(tri.T); ok && condForcesNonFalse(tri.Cond, tgt, true) {
		tType, tEmpty = dropFalseClean(tType)
	}
	if tgt, ok := ctorTarget(tri.F); ok && condForcesNonFalse(tri.Cond, tgt, false) {
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

func ctorTarget(e ast.Expr) (string, bool) {
	if name, _, ok := unwrapThisMember(e); ok {
		return "." + name, true
	}
	p := ctorPeel(e)
	if id, ok := p.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		return id.Name, true
	}
	return "", false
}

func condForcesNonFalse(cond ast.Expr, target string, polarity bool) bool {
	c := ctorPeel(cond)
	switch n := c.(type) {
	case *ast.Unary:
		if n.Tok == tok.Not {
			return condForcesNonFalse(n.E, target, !polarity)
		}
	case *ast.Binary:
		if (n.Tok == tok.Isnt && polarity) || (n.Tok == tok.Is && !polarity) {
			if t, ok := falseCmpTarget(n); ok {
				return t == target
			}
		}
	default:
		if polarity {
			if t, ok := ctorTarget(c); ok {
				return t == target
			}
		}
	}
	return false
}

func falseCmpTarget(b *ast.Binary) (string, bool) {
	if isFalseLiteral(b.Rhs) {
		return ctorTarget(b.Lhs)
	}
	if isFalseLiteral(b.Lhs) {
		return ctorTarget(b.Rhs)
	}
	return "", false
}

func isFalseLiteral(e ast.Expr) bool {
	c, ok := unwrapConstant(e)
	return ok && DynTypeOfSuValue(c.Val) == TFalse
}

func dropFalseClean(t DynType) (DynType, bool) {
	switch x := t.(type) {
	case Primitive:
		switch x {
		case TFalse:
			return TUnknown, true
		case TBoolean:
			return TTrue, false
		}
		return t, false
	case Union:
		var kept []DynType
		for _, m := range x.Types {
			switch m {
			case TFalse:
				// drop
			case TBoolean:
				kept = append(kept, TTrue)
			default:
				kept = append(kept, m)
			}
		}
		if len(kept) == 0 {
			return TUnknown, true
		}
		return Union{Types: kept, IsDirty: x.IsDirty}.Fold(), false
	}
	return t, false
}

func (e *ctorExec) feasibility(cond ast.Expr, st *ctorState) int {
	c := ctorPeel(cond)
	switch n := c.(type) {
	case *ast.Binary:
		return e.feasCompare(n, st)
	case *ast.Unary:
		if n.Tok == tok.Not {
			return -e.feasibility(n.E, st)
		}
	case *ast.Nary:
		switch n.Tok {
		case tok.And:
			return e.feasAnd(n.Exprs, st)
		case tok.Or:
			return e.feasOr(n.Exprs, st)
		}
	}
	return 0
}

func (e *ctorExec) feasCompare(n *ast.Binary, st *ctorState) int {
	if n.Tok != tok.Is && n.Tok != tok.Isnt {
		return 0
	}
	if !disjoint(e.evalType(n.Lhs, st), e.evalType(n.Rhs, st)) {
		return 0
	}
	if n.Tok == tok.Is {
		return -1 // disjoint values can never be equal
	}
	return +1 // disjoint values are never equal, so isnt is always true
}

func (e *ctorExec) feasAnd(exprs []ast.Expr, st *ctorState) int {
	allTrue := true
	for _, sub := range exprs {
		switch e.feasibility(sub, st) {
		case -1:
			return -1
		case 0:
			allTrue = false
		}
	}
	if allTrue {
		return +1
	}
	return 0
}

func (e *ctorExec) feasOr(exprs []ast.Expr, st *ctorState) int {
	allFalse := true
	for _, sub := range exprs {
		switch e.feasibility(sub, st) {
		case +1:
			return +1
		case 0:
			allFalse = false
		}
	}
	if allFalse {
		return -1
	}
	return 0
}

func disjoint(a, b DynType) bool {
	pa, aok := a.(Primitive)
	pb, bok := b.(Primitive)
	if !aok || !bok {
		return false
	}
	if pa == TUnknown || pb == TUnknown || pa == TVoid || pb == TVoid {
		return false
	}
	if pa == pb || subtypeOf(pa, pb) || subtypeOf(pb, pa) {
		return false
	}
	return true
}

func membersAssignedFalse(cls *ClassObject, env TypeEnv) map[string]bool {
	out := map[string]bool{}
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		if b, ok := n.(*ast.Binary); ok && b.Tok == tok.Eq {
			if name, _, ok := unwrapThisMember(b.Lhs); ok {
				if containsFalse(env.GetType(b.Rhs)) {
					out[name] = true
				}
			}
		}
		n.Children(func(c ast.Node) ast.Node {
			walk(c)
			return c
		})
	}
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			walk(stmt)
		}
	}
	return out
}

func demotedMemberRead(n ast.Node, demotions map[string]DynType) (*ast.Mem, DynType, bool) {
	mem, ok := n.(*ast.Mem)
	if !ok {
		return nil, nil, false
	}
	id, ok := mem.E.(*ast.Ident)
	if !ok || id.Name != "this" {
		return nil, nil, false
	}
	name, ok := memberName(mem.M)
	if !ok {
		return nil, nil, false
	}
	t, ok := demotions[name]
	if !ok {
		return nil, nil, false
	}
	return mem, t, true
}

func restampMemberReads(cls *ClassObject, env TypeEnv, demotions map[string]DynType) {
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		if mem, t, ok := demotedMemberRead(n, demotions); ok {
			env.SetType(mem, t)
		}
		n.Children(func(c ast.Node) ast.Node {
			walk(c)
			return c
		})
	}
	for name, fn := range cls.SortedMethods {
		if name == "New" {
			continue
		}
		for _, stmt := range fn.Body {
			walk(stmt)
		}
	}
}

func cloneTypeMap(m map[string]DynType) map[string]DynType {
	out := make(map[string]DynType, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func joinMembers(states []map[string]DynType) map[string]DynType {
	if len(states) == 1 {
		return states[0]
	}
	keys := map[string]bool{}
	for _, s := range states {
		for k := range s {
			keys[k] = true
		}
	}
	out := make(map[string]DynType, len(keys))
	for k := range keys {
		var t DynType
		for _, s := range states {
			v, ok := s[k]
			if !ok {
				v = TUnknown
			}
			if t == nil {
				t = v
			} else {
				t = U(t, v)
			}
		}
		out[k] = t
	}
	return out
}

func containsFalse(t DynType) bool {
	switch x := t.(type) {
	case Primitive:
		return x == TFalse || x == TBoolean || x == TUnknown
	case Union:
		if x.IsDirty {
			return true
		}
		for _, m := range x.Types {
			if containsFalse(m) {
				return true
			}
		}
	}
	return false
}

func removeFalse(t DynType) DynType {
	r, _ := dropFalseClean(t)
	return r
}

func ctorPeel(e ast.Expr) ast.Expr {
	for {
		switch x := e.(type) {
		case *ast.ExprPos:
			if x.Expr == nil {
				return e
			}
			e = x.Expr
		case *ast.Unary:
			if x.Tok != tok.LParen {
				return e
			}
			e = x.E
		default:
			return e
		}
	}
}

func stmtList(s ast.Statement) []ast.Statement {
	if s == nil {
		return nil
	}
	if c, ok := s.(*ast.Compound); ok {
		return c.Body
	}
	return []ast.Statement{s}
}
