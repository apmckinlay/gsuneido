// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"maps"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// a narrowed type and whether a guard proved it. Only in-guard entries are
// visible to reads, so the two travel together - see refinements.guarded.
type refinement struct {
	Typ     DynType
	InGuard bool
}

type refinements map[string]refinement

// typ answers whatever type is recorded, guard-proven or not
func (r refinements) typ(name string) (DynType, bool) {
	f, ok := r[name]
	if !ok || f.Typ == nil {
		return nil, false
	}
	return f.Typ, true
}

// guarded answers only when a guard proved the type
func (r refinements) guarded(name string) (DynType, bool) {
	f, ok := r[name]
	if !ok || !f.InGuard || f.Typ == nil {
		return nil, false
	}
	return f.Typ, true
}

func (r refinements) prove(name string, t DynType) {
	r[name] = refinement{Typ: t, InGuard: true}
}

// setType types a name without claiming a guard proved it
func (r refinements) setType(name string, t DynType) {
	f := r[name]
	f.Typ = t
	r[name] = f
}

type narrowScope struct {
	Locals              refinements
	Members             refinements
	memberAssignedFalse map[string]bool // computed once per pass, shared read-only across forks; never mutate
	writes              *classMemberWrites
	postconds           boolPostconds
	postHook            func(*ast.Return, narrowScope)
}

func (s narrowScope) kind(members bool) refinements {
	if members {
		return s.Members
	}
	return s.Locals
}

func newNarrowScope(size int) narrowScope {
	return narrowScope{
		Locals:  make(refinements, size),
		Members: make(refinements, size),
	}
}

func (s narrowScope) clone() narrowScope {
	c := newNarrowScope(len(s.Locals))
	maps.Copy(c.Locals, s.Locals)
	maps.Copy(c.Members, s.Members)
	c.memberAssignedFalse = s.memberAssignedFalse // shared, read-only
	c.writes = s.writes                           // shared, read-only
	c.postconds = s.postconds                     // shared, read-only
	c.postHook = s.postHook                       // shared, read-only
	return c
}

type narrowTarget struct {
	name     string // bare name, no leading dot for members
	isMember bool
	node     ast.Node // *ast.Ident or *ast.Mem
}

func unwrapTarget(e ast.Expr) (narrowTarget, bool) {
	if id, ok := unwrapIdent(e); ok && !isGlobalIdent(id.Name) {
		return narrowTarget{name: id.Name, node: id}, true
	}
	if name, mem, ok := unwrapThisMember(e); ok {
		return narrowTarget{name: name, isMember: true, node: mem}, true
	}
	return narrowTarget{}, false
}

// peelParens strips position and parenthesis wrappers to reach the expression
// that actually carries a type stam
func peelParens(e ast.Expr) ast.Expr {
	for {
		switch x := e.(type) {
		case *ast.ExprPos:
			if x.Expr == nil {
				return e
			}
			e = x.Expr
		case *ast.Unary:
			if x.Tok != tok.LParen || x.E == nil {
				return e
			}
			e = x.E
		default:
			return e
		}
	}
}

func unwrapThisMember(e ast.Expr) (string, *ast.Mem, bool) {
	for {
		switch x := e.(type) {
		case *ast.ExprPos:
			if x.Expr == nil {
				return "", nil, false
			}
			e = x.Expr
		case *ast.Unary:
			if x.Tok != tok.LParen {
				return "", nil, false
			}
			e = x.E
		default:
			mem, ok := e.(*ast.Mem)
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
			return name, mem, true
		}
	}
}

var predicateTargets = map[string][]DynType{
	"Boolean?":  {TBoolean},
	"Number?":   {TNumber},
	"String?":   {TString},
	"Date?":     {TDate},
	"Object?":   {TObject},
	"Record?":   {TObject},
	"Class?":    {TClass},
	"Function?": {TFunction, TBlock},
}

var typeStringTargets = map[string][]DynType{
	"Boolean":  {TBoolean},
	"Number":   {TNumber},
	"String":   {TString},
	"Date":     {TDate},
	"Object":   {TObject},
	"Record":   {TObject},
	"Class":    {TClass},
	"Function": {TFunction},
	"Block":    {TBlock},
}

func FlowNarrowingPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	assignedFalse := membersAssignedFalse(cls, env)
	writes := ComputeMemberWrites(cls)
	postconds := ComputeBoolPostconditions(cls, env, assignedFalse, writes)
	for _, fn := range cls.SortedMethods {
		sc := initialNarrowScope(fn, env)
		sc.memberAssignedFalse = assignedFalse
		sc.writes = writes
		sc.postconds = postconds
		walkBlock(fn.Body, env, sc)
	}
	return false
}

func initialNarrowScope(fn *ast.Function, env TypeEnv) narrowScope {
	sc := newNarrowScope(len(fn.Params) + 8)
	for i := range fn.Params {
		p := &fn.Params[i]
		name := p.Name.ParamName()
		if t, ok := env.Params[p]; ok {
			sc.Locals.setType(name, t)
		} else if len(p.Name.Name) > 0 && p.Name.Name[0] == '.' {
			if t, ok := env.LookupMember(name); ok {
				sc.Locals.setType(name, t)
			}
		}
	}
	return sc
}

func narrowWalk(n ast.Node, env TypeEnv, sc narrowScope) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *ast.ExprPos:
		if x.Expr != nil {
			narrowWalk(x.Expr, env, sc)
			if t, ok := env.Nodes[x.Expr]; ok {
				env.SetType(x, t)
			}
		}
		return
	case *ast.Compound:
		walkBlock(x.Body, env, sc)
		return
	case *ast.If:
		narrowIf(x, env, sc)
		return
	case *ast.Trinary:
		narrowTrinary(x, env, sc)
		return
	case *ast.Switch:
		narrowSwitch(x, env, sc)
		return
	case *ast.While:
		narrowWhile(x, env, sc)
		return
	case *ast.For:
		narrowFor(x, env, sc)
		return
	case *ast.DoWhile:
		loopSc := loopEntryScope(x.Body, nil, sc, env)
		narrowWalk(x.Body, env, loopSc)
		narrowWalk(x.Cond, env, loopSc)
		return
	case *ast.ForIn, *ast.Forever:
		// no boolean guard to exploit; clone + strip loop-carried.
		narrowChildren(n, env, loopEntryScope(loopBodyOf(n), nil, sc, env))
		return
	case *ast.Nary:
		if x.Tok == tok.And || x.Tok == tok.Or {
			narrowAndOr(x, env, sc)
			return
		}
	case *ast.Binary:
		if x.Tok == tok.Eq {
			narrowEqAssign(x, env, sc)
			return
		}
		if x.Tok.IsAssign() {
			narrowChildren(n, env, sc)
			killRefinement(x.Lhs, sc)
			return
		}
	case *ast.Unary:
		if isIncDec(x.Tok) {
			narrowChildren(n, env, sc)
			killRefinement(x.E, sc)
			return
		}
	case *ast.Ident:
		if !isGlobalIdent(x.Name) {
			if t, ok := sc.Locals.guarded(x.Name); ok {
				env.SetType(x, t)
			}
		}
		return
	case *ast.Mem:
		if name, _, ok := unwrapThisMember(x); ok {
			if t, ok2 := sc.Members.guarded(name); ok2 {
				env.SetType(x, t)
			}
		}
		// no return: a Mem still walks its children below
	case *ast.Return:
		narrowChildren(n, env, sc)
		if sc.postHook != nil {
			sc.postHook(x, sc)
		}
		return
	case *ast.Call:
		narrowCall(x, env, sc)
		return
	}
	narrowChildren(n, env, sc)
}

func narrowChildren(n ast.Node, env TypeEnv, sc narrowScope) {
	n.Children(func(c ast.Node) ast.Node {
		narrowWalk(c, env, sc)
		return c
	})
}

// writing a target discards whatever a guard proved about it
func killRefinement(target ast.Expr, sc narrowScope) {
	if name, _, ok := unwrapThisMember(target); ok {
		delete(sc.Members, name)
		return
	}
	if id, ok := target.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		delete(sc.Locals, id.Name)
	}
}

func narrowTrinary(x *ast.Trinary, env TypeEnv, sc narrowScope) {
	narrowWalk(x.Cond, env, sc)
	tScope := refineCond(x.Cond, sc, true, env, true)
	narrowWalk(x.T, env, tScope)
	fScope := refineCond(x.Cond, sc, false, env, true)
	narrowWalk(x.F, env, fScope)
	if t := trinaryType(x, env); t != TUnknown {
		env.SetType(x, t)
	}
}

func narrowWhile(x *ast.While, env TypeEnv, sc narrowScope) {
	loopSc := loopEntryScope(x.Body, nil, sc, env)
	narrowWalk(x.Cond, env, loopSc)
	body := refineCond(x.Cond, loopSc, true, env, true)
	narrowWalk(x.Body, env, body)
}

func narrowFor(x *ast.For, env TypeEnv, sc narrowScope) {
	loopSc := loopEntryScope(x.Body, x.Inc, sc, env)
	for _, e := range x.Init {
		narrowWalk(e, env, loopSc)
	}
	body := loopSc
	if x.Cond != nil {
		narrowWalk(x.Cond, env, loopSc)
		body = refineCond(x.Cond, loopSc, true, env, true)
	}
	narrowWalk(x.Body, env, body)
	for _, e := range x.Inc {
		narrowWalk(e, env, body)
	}
}

func narrowAndOr(x *ast.Nary, env TypeEnv, sc narrowScope) {
	siblingPolarity := x.Tok == tok.And
	sub := sc.clone()
	for _, e := range x.Exprs {
		narrowWalk(e, env, sub)
		applyRefinement(e, sub, siblingPolarity, env, true)
	}
}

func narrowCall(x *ast.Call, env TypeEnv, sc narrowScope) {
	narrowChildren(x, env, sc)
	eff, opaque := sc.writes.callEffect(x)
	for k := range sc.Members {
		if !opaque && !eff[k] {
			continue // refinement survives: call can't touch k
		}
		if t, keep := memberTypeAcrossCall(k, sc, env); keep {
			sc.Members.setType(k, t)
			continue
		}
		delete(sc.Members, k)
	}
}

func narrowEqAssign(x *ast.Binary, env TypeEnv, sc narrowScope) {
	narrowWalk(x.Rhs, env, sc)
	if id, ok := x.Lhs.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		// new value invalidates any prior guard refinement
		delete(sc.Locals, id.Name)
		rhsT := env.GetType(x.Rhs)
		lhsT := env.GetType(id)
		if isNarrower(rhsT, lhsT) {
			sc.Locals.prove(id.Name, rhsT)
			env.SetType(id, rhsT)
		}
	} else if name, mem, ok := unwrapThisMember(x.Lhs); ok {
		delete(sc.Members, name)
		rhsT := env.GetType(x.Rhs)
		lhsT := env.GetType(mem)
		if isNarrower(rhsT, lhsT) {
			sc.Members.prove(name, rhsT)
			env.SetType(mem, rhsT)
		}
	} else {
		narrowWalk(x.Lhs, env, sc)
	}
}

func loopEntryScope(body ast.Statement, inc []ast.Expr, sc narrowScope, env TypeEnv) narrowScope {
	out := sc.clone()
	w := collectLoopWrites(body, sc.writes)
	for _, e := range inc {
		w.scan(e)
	}
	for name := range w.locals {
		delete(out.Locals, name)
	}
	for name := range out.Members {
		if !w.opaque && !w.callMembers[name] && !w.members[name] {
			continue
		}
		if !w.members[name] {
			if t, keep := memberTypeAcrossCall(name, sc, env); keep {
				out.Members.setType(name, t)
				continue
			}
		}
		delete(out.Members, name)
	}
	return out
}

type loopWrites struct {
	locals      map[string]bool
	members     map[string]bool
	callMembers map[string]bool
	opaque      bool
	writes      *classMemberWrites
}

// `v = v` cannot change v's value on any iteration, so it does not
// invalidate a loop-entry fact. only plain `=` qualifies: `v += v` widens
func selfAssign(b *ast.Binary) bool {
	if b.Tok != tok.Eq {
		return false
	}
	lhs, ok := b.Lhs.(*ast.Ident)
	if !ok {
		return false
	}
	rhs, ok := unwrapIdent(b.Rhs)
	return ok && rhs.Name == lhs.Name
}

func collectLoopWrites(body ast.Statement, writes *classMemberWrites) *loopWrites {
	w := &loopWrites{
		locals:      map[string]bool{},
		members:     map[string]bool{},
		callMembers: map[string]bool{},
		writes:      writes,
	}
	if body != nil {
		w.scan(body)
	}
	return w
}

func (w *loopWrites) scan(n ast.Node) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *ast.Call:
		eff, opaque := w.writes.callEffect(x)
		if opaque {
			w.opaque = true
		}
		for m := range eff {
			w.callMembers[m] = true
		}
	case *ast.Binary:
		if x.Tok.IsAssign() && !selfAssign(x) {
			w.mark(x.Lhs)
		}
	case *ast.Unary:
		if x.Tok == tok.Inc || x.Tok == tok.Dec ||
			x.Tok == tok.PostInc || x.Tok == tok.PostDec {
			w.mark(x.E)
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		w.scan(c)
		return c
	})
}

func (w *loopWrites) mark(e ast.Expr) {
	if id, ok := e.(*ast.Ident); ok && !isGlobalIdent(id.Name) {
		w.locals[id.Name] = true
	} else if name, _, ok := unwrapThisMember(e); ok {
		w.members[name] = true
	}
}

func loopBodyOf(n ast.Node) ast.Statement {
	switch x := n.(type) {
	case *ast.ForIn:
		return x.Body
	case *ast.Forever:
		return x.Body
	}
	return nil
}

func memberTypeAcrossCall(name string, sc narrowScope, env TypeEnv) (DynType, bool) {
	if sc.memberAssignedFalse == nil || sc.memberAssignedFalse[name] {
		return nil, false // can't prove it stayed non-false
	}
	cw, ok := env.LookupMember(name)
	if !ok {
		return nil, false
	}
	rf := removeFalse(cw)
	if rf == TUnknown || typeHasBoolish(rf) {
		return nil, false
	}
	return rf, true
}

func typeHasBoolish(t DynType) bool {
	switch x := t.(type) {
	case Primitive:
		return x == TTrue || x == TBoolean
	case Union:
		for _, m := range x.Types {
			if m == TTrue || m == TBoolean {
				return true
			}
		}
	}
	return false
}

func walkBlock(stmts []ast.Statement, env TypeEnv, sc narrowScope) {
	for _, stmt := range stmts {
		if ifStmt, ok := stmt.(*ast.If); ok {
			narrowIf(ifStmt, env, sc)
			continue
		}
		narrowWalk(stmt, env, sc)
		if call, ok := assertStmtCall(stmt); ok {
			applyAssertStmt(call, sc, env)
		}
	}
}

// the Assert(...) call of an expression-statement, if stmt is one
func assertStmtCall(stmt ast.Statement) (*ast.Call, bool) {
	es, ok := stmt.(*ast.ExprStmt)
	if !ok || es.E == nil {
		return nil, false
	}
	e := es.E
	if ep, ok := e.(*ast.ExprPos); ok && ep.Expr != nil {
		e = ep.Expr
	}
	call, ok := e.(*ast.Call)
	if !ok {
		return nil, false
	}
	if id, ok := call.Fn.(*ast.Ident); !ok || id.Name != "Assert" {
		return nil, false
	}
	if len(call.Args) == 0 || call.Args[0].Name != nil {
		return nil, false
	}
	return call, true
}

// executing past an Assert proves its claim (a failure throws), so the fact
// holds for the rest of the block. two shapes (see stdlib Assert):
//
// ```suneido
// Assert(String?(x))          // old style: a condition, msg optional
// Assert(x isString:)         // matcher style: Matcher_<name>.Match(x, args)
//
//	^ ^^^^^^^^^ the sole non-msg named arg names the matcher
//
// ```
func applyAssertStmt(call *ast.Call, sc narrowScope, env TypeEnv) {
	name, matcherArg, ok := assertMatcher(call)
	if !ok {
		applyRefinement(call.Args[0].E, sc, true, env, true)
		return
	}
	tgt, tok := unwrapTarget(call.Args[0].E)
	if !tok {
		return
	}
	existing := existingType(sc, tgt, env)
	if targets, ok := assertMatcherTargets[name]; ok {
		storeRefinement(sc, tgt, narrowTowardSet(existing, targets), true)
		return
	}
	switch name {
	case "is": // Matcher_is: equality with the expected literal
		if c, ok := unwrapConstant(matcherArg); ok {
			storeRefinement(sc, tgt,
				narrowTowardSet(existing, []DynType{DynTypeOfSuValue(c.Val)}), true)
		}
	case "isnt":
		if c, ok := unwrapConstant(matcherArg); ok {
			if t := DynTypeOfSuValue(c.Val); t == TFalse || t == TTrue {
				storeRefinement(sc, tgt, narrowAwaySet(existing, []DynType{t}), true)
			}
		}
	case "isType": // Matcher_isType: Type(value) is args
		if c, ok := unwrapConstant(matcherArg); ok {
			if s, ok := c.Val.ToStr(); ok {
				if targets, ok := typeStringTargets[s]; ok {
					storeRefinement(sc, tgt, narrowTowardSet(existing, targets), true)
				}
			}
		}
	}
}

// the sole non-msg named arg of a matcher-style Assert; ok=false means old
// style (condition form). unrecognized matchers still return ok=true so the
// condition path is not misapplied to a matcher subject
func assertMatcher(call *ast.Call) (name string, arg ast.Expr, ok bool) {
	for i := range call.Args {
		a := &call.Args[i]
		if a.Name == nil || isAtArg(a) {
			continue
		}
		n, isStr := argName(a)
		if !isStr || n == "msg" {
			continue
		}
		return n, a.E, true
	}
	return "", nil, false
}

// matchers whose Match is exactly a type predicate (see stdlib Matcher_is*);
// the any* forms subclass the is* forms unchanged. mirror predicateTargets.
var assertMatcherTargets = map[string][]DynType{
	"isString":         {TString},
	"anyString":        {TString},
	"isNumber":         {TNumber},
	"anyNumber":        {TNumber},
	"isInt":            {TNumber},
	"isIntNonNegative": {TNumber},
	"isIntInRange":     {TNumber},
	"isObject":         {TObject},
	"anyObject":        {TObject},
	"isBoolean":        {TBoolean},
	"isDate":           {TDate},
	"isCallable":       {TFunction, TBlock}, // Matcher_isCallable is Function?
}
