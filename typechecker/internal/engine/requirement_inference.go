// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// soundness: a demand is only taken from a clean primitive on a non-guessed
// signature, only for a param the method never reassigns, and only where the
// arg's own type is still unknown (a narrowing-guarded use proves nothing
// about the caller). two disagreeing demands conflict: the slot drops to
// unconstrained, permanently (env.ReqConflicts)
//
// returns whether any signature changed, so callers can drive a fixpoint
//
// confluence: a demand read from a slot that LATER conflicts is unjustified -
// the source turned out polymorphic. so any sweep that grows the conflict set
// clears every inferred slot and the iteration re-derives from scratch under
// the larger (sticky, grow-only) conflict set. once conflicts stabilize the
// system is monotone, so the result is order-independent
func RequirementPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	changed := false
	for range maxFixpointPasses {
		before := len(env.ReqConflicts)
		swept := requirementSweep(cls, env)
		if len(env.ReqConflicts) > before {
			clearInferred(env.MethodSigs)
		}
		if !swept {
			break
		}
		changed = true
	}
	return changed
}

func clearInferred(sigs map[string]*Signature) {
	for _, sig := range sigs {
		if sig == nil {
			continue
		}
		for i := range sig.Params {
			if sig.Params[i].Inferred {
				sig.Params[i].Typ = TUnknown
				sig.Params[i].Inferred = false
				sig.Params[i].Why = ""
			}
		}
	}
}

func requirementSweep(cls *ClassObject, env TypeEnv) bool {
	changed := false
	for name, fn := range cls.SortedMethods {
		sig := env.MethodSigs[name]
		if fn == nil || sig == nil || sig.AtParam {
			continue
		}
		slots := fillableSlots(fn, sig)
		if len(slots) == 0 {
			continue
		}
		params := paramNameSet(fn)
		for p := range falseTestedParams(fn, params) {
			delete(slots, p)
		}
		if len(slots) == 0 {
			continue
		}
		d := demandCtx{
			env:      env,
			method:   name,
			sig:      sig,
			slots:    slots,
			assigned: assignedNames(fn),
		}
		for _, stmt := range fn.Body[:sentinelCutoff(fn, params)] {
			if stmt != nil && d.walk(stmt) {
				changed = true
			}
		}
	}
	return changed
}

func paramNameSet(fn *ast.Function) map[string]bool {
	out := make(map[string]bool, len(fn.Params))
	for i := range fn.Params {
		out[fn.Params[i].Name.ParamName()] = true
	}
	return out
}

func falseTestedParams(fn *ast.Function, params map[string]bool) map[string]bool {
	out := map[string]bool{}
	for _, stmt := range fn.Body {
		if stmt != nil {
			collectFalseTests(stmt, params, out)
		}
	}
	return out
}

func collectFalseTests(n ast.Node, params map[string]bool, out map[string]bool) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		collectFalseTests(ep.Expr, params, out)
		return
	}
	if b, ok := n.(*ast.Binary); ok {
		if id, _, ok := identFalseTest(b); ok && params[id.Name] {
			out[id.Name] = true
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		collectFalseTests(c, params, out)
		return c
	})
}

// index of the first top-level `if <false-test of any param> { ...always exits }`:
// past that guard the body never sees the sentinel, so demands taken there
// would wrongly bind callers that pass it
func sentinelCutoff(fn *ast.Function, params map[string]bool) int {
	for i, stmt := range fn.Body {
		iff, ok := stmt.(*ast.If)
		if !ok || !branchAlwaysExits(iff.Then) {
			continue
		}
		found := map[string]bool{}
		collectFalseTests(iff.Cond, params, found)
		if len(found) > 0 {
			return i
		}
	}
	return len(fn.Body)
}

// plain params whose signature slot an inferred requirement may fill:
// unannotated (or previously inferred), no default, and named by a bare
// identifier in the body (excludes .member-init, @args, _dynamic params).
func fillableSlots(fn *ast.Function, sig *Signature) map[string]int {
	slots := map[string]int{}
	for i := range sig.Params {
		if i >= len(fn.Params) {
			break
		}
		p := &sig.Params[i]
		if p.HasDefault || (p.Typ != nil && p.Typ != TUnknown && !p.Inferred) {
			continue
		}
		raw := fn.Params[i].Name.Name
		if raw == "" || raw[0] == '.' || raw[0] == '@' || raw[0] == '_' {
			continue
		}
		slots[p.Name] = i
	}
	return slots
}

type demandCtx struct {
	env      TypeEnv
	method   string
	sig      *Signature
	slots    map[string]int
	assigned map[string]bool
}

func (d *demandCtx) walk(n ast.Node) bool {
	if n == nil {
		return false
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		return d.walk(ep.Expr)
	}
	changed := false
	n.Children(func(c ast.Node) ast.Node {
		if d.walk(c) {
			changed = true
		}
		return c
	})
	if call, ok := n.(*ast.Call); ok {
		if d.callDemands(call) {
			changed = true
		}
	}
	return changed
}

func (d *demandCtx) callDemands(call *ast.Call) bool {
	sig := d.env.CallSig(call)
	if sig == nil || sig.AtParam || d.env.GuessedCall(call) {
		return false
	}
	changed := false
	for i := range call.Args {
		arg := &call.Args[i]
		if isAtArg(arg) {
			break
		}
		id := bareIdent(arg.E)
		if id == nil {
			continue
		}
		k, ok := d.slots[id.Name]
		if !ok || d.assigned[id.Name] {
			continue
		}
		if d.env.GetType(arg.E) != TUnknown || d.env.GetType(id) != TUnknown {
			continue
		}
		callee := matchParam(sig, arg, i)
		if callee == nil || !cleanRequirement(callee.Typ) {
			continue
		}
		if d.apply(k, id.Name, callee.Typ, demandWhy(call, callee)) {
			changed = true
		}
	}
	return changed
}

// a demand chained through an inferred slot keeps that slot's root
func demandWhy(call *ast.Call, callee *Param) string {
	if callee.Inferred && callee.Why != "" {
		return callee.Why
	}
	return calleeDisplay(call)
}

func (d *demandCtx) apply(k int, param string, req DynType, why string) bool {
	key := d.method + "\x00" + param
	if d.env.ReqConflicts[key] {
		return false
	}
	slot := &d.sig.Params[k]
	switch {
	case slot.Typ == nil || slot.Typ == TUnknown:
		slot.Typ = req
		slot.Inferred = true
		slot.Why = why
		return true
	case slot.Inferred && slot.Typ != req:
		slot.Typ = TUnknown
		slot.Inferred = false
		slot.Why = ""
		d.env.ReqConflicts[key] = true
		return true
	}
	return false
}

// calleeName plus the class qualifier for Foo.Bar(...) static calls.
func calleeDisplay(call *ast.Call) string {
	if mem, ok := call.Fn.(*ast.Mem); ok {
		if name, ok := memberName(mem.M); ok && name != "*new*" {
			if id, ok := mem.E.(*ast.Ident); ok && isGlobalIdent(id.Name) {
				return id.Name + "." + name
			}
		}
	}
	return calleeName(call)
}

func cleanRequirement(t DynType) bool {
	p, ok := t.(Primitive)
	return ok && p != TUnknown && p != TVoid
}

func bareIdent(e ast.Expr) *ast.Ident {
	if ep, ok := e.(*ast.ExprPos); ok {
		e = ep.Expr
	}
	id, ok := e.(*ast.Ident)
	if !ok || isGlobalIdent(id.Name) {
		return nil
	}
	return id
}

func assignedNames(fn *ast.Function) map[string]bool {
	names := map[string]bool{}
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		markAssigned(n, names)
		n.Children(func(c ast.Node) ast.Node {
			walk(c)
			return c
		})
	}
	for _, stmt := range fn.Body {
		if stmt != nil {
			walk(stmt)
		}
	}
	return names
}

func markAssigned(n ast.Node, names map[string]bool) {
	switch x := n.(type) {
	case *ast.Binary:
		if x.Tok.IsAssign() {
			if id := bareIdent(x.Lhs); id != nil {
				names[id.Name] = true
			}
		}
	case *ast.Unary:
		if x.Tok == tok.Inc || x.Tok == tok.Dec ||
			x.Tok == tok.PostInc || x.Tok == tok.PostDec {
			if id := bareIdent(x.E); id != nil {
				names[id.Name] = true
			}
		}
	case *ast.ForIn:
		if x.Var.Name != "" {
			names[x.Var.Name] = true
		}
		if x.Var2.Name != "" {
			names[x.Var2.Name] = true
		}
	}
}
