// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// treats `Assert(cond)` whose condition constrains a class
// member as a class-wide ground truth
//
// ```suneido
//
//	class {
//	    New(.x, .y) { Assert(Number?(x) and Number?(y)) }
//	    Bar() { return .x $ .y }   // ERROR: .x is a Number ground truth
//	}
//
// ```

// ```suneido
// New() { Assert(Number?(.count)) }
//
//	^^^^^^^^^^^^^^^ ground truth: .count is pinned to Number
//	                class-wide, overriding the assignment union
//
// ```
func AssertMemberPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	gts := collectAssertGroundTruths(cls, env)
	checkLocalAssertContradictions(cls, env, gts)
	if len(gts) == 0 {
		return false
	}
	for name, gt := range gts {
		env.SeedMember(name, gt.ty)
		env.AssertedMembers[name] = AssertFact{Ty: gt.ty, Method: gt.method, Pos: gt.pos}
	}
	return true
}

type AssertFact struct {
	Ty     DynType
	Method string
	Pos    int
}

func AssertAssignmentCheckPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	if len(env.AssertedMembers) == 0 {
		return false
	}
	gts := map[string]assertGroundTruth{}
	for name, f := range env.AssertedMembers {
		gts[name] = assertGroundTruth{ty: f.Ty, method: f.Method, pos: f.Pos}
	}
	checkAssertAssignments(cls, env, gts)
	return false
}

func checkLocalAssertContradictions(cls *ClassObject, env TypeEnv, gts map[string]assertGroundTruth) {
	for method, fn := range cls.SortedMethods {
		if method == "New" {
			continue
		}
		inlineInit := privateInlineInitParams(fn)
		forEachAssert(fn, func(cond ast.Expr, pos int) {
			for name, ty := range assertedMembers(cond, env, inlineInit) {
				gt, pinned := gts[name]
				if !pinned {
					continue
				}
				if _, ok := combineAsserts(gt.ty, ty); ok {
					continue
				}
				env.Report(&Diagnostic{
					Severity: SeverityError,
					Method:   method,
					Pos:      pos,
					Msg: fmt.Sprintf("conflicting Assert ground truths for member .%s: "+
						"%v asserted in %s, %v asserted in %s",
						name, gt.ty, gt.method, ty, method),
				})
			}
		})
	}
}

type assertGroundTruth struct {
	ty     DynType
	method string
	pos    int
}

func collectAssertGroundTruths(cls *ClassObject, env TypeEnv) map[string]assertGroundTruth {
	gts := map[string]assertGroundTruth{}
	for method, fn := range cls.SortedMethods {
		if method != "New" {
			continue
		}
		inlineInit := privateInlineInitParams(fn)
		forEachAssert(fn, func(cond ast.Expr, pos int) {
			for name, ty := range assertedMembers(cond, env, inlineInit) {
				prev, seen := gts[name]
				if !seen {
					gts[name] = assertGroundTruth{ty: ty, method: method, pos: pos}
					continue
				}
				if combined, ok := combineAsserts(prev.ty, ty); ok {
					gts[name] = assertGroundTruth{ty: combined, method: prev.method, pos: prev.pos}
					continue
				}
				env.Report(&Diagnostic{
					Severity: SeverityError,
					Method:   method,
					Pos:      pos,
					Msg: fmt.Sprintf("conflicting Assert ground truths for member .%s: "+
						"%v asserted in %s, %v asserted in %s",
						name, prev.ty, prev.method, ty, method),
				})
			}
		})
	}
	return gts
}

func privateInlineInitParams(fn *ast.Function) map[string]bool {
	out := map[string]bool{}
	for i := range fn.Params {
		p := &fn.Params[i]
		raw := p.Name.Name
		if len(raw) > 1 && raw[0] == '.' && !isPublicMember(raw[1:]) {
			out[p.Name.ParamName()] = true
		}
	}
	return out
}

func assertedMembers(cond ast.Expr, env TypeEnv, inlineInit map[string]bool) map[string]DynType {
	if cond == nil {
		return nil
	}
	sc := newNarrowScope(4)
	applyRefinement(cond, sc, true, env, true)

	out := map[string]DynType{}
	for name, f := range sc.Members {
		// public members are externally mutable - not a sound ground truth
		if isPublicMember(name) {
			continue
		}
		if f.Typ != nil && f.Typ != TUnknown {
			out[name] = f.Typ
		}
	}
	for name, f := range sc.Locals {
		if inlineInit[name] && f.Typ != nil && f.Typ != TUnknown {
			out[name] = f.Typ
		}
	}
	return out
}

func combineAsserts(a, b DynType) (DynType, bool) {
	switch {
	case dynEqual(a, b):
		return a, true
	case isNarrower(a, b):
		return a, true
	case isNarrower(b, a):
		return b, true
	case typeFits(a, b):
		return a, true
	case typeFits(b, a):
		return b, true
	}
	return nil, false
}

func checkAssertAssignments(cls *ClassObject, env TypeEnv, gts map[string]assertGroundTruth) {
	for method, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				checkAssertAssignsIn(stmt, method, env, gts)
			}
		}
	}
}

func checkAssertAssignsIn(n ast.Node, method string, env TypeEnv, gts map[string]assertGroundTruth) {
	if n == nil {
		return
	}
	if b, name, ok := thisMemberAssign(n); ok {
		if gt, asserted := gts[name]; asserted {
			reportAssignViolation(b, name, method, env, gt)
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		checkAssertAssignsIn(c, method, env, gts)
		return c
	})
}

// unwraps a `this.member = rhs` assignment into the assign node and member name
func thisMemberAssign(n ast.Node) (*ast.Binary, string, bool) {
	b, ok := n.(*ast.Binary)
	if !ok || b.Tok != tok.Eq {
		return nil, "", false
	}
	mem, ok := b.Lhs.(*ast.Mem)
	if !ok {
		return nil, "", false
	}
	id, ok := mem.E.(*ast.Ident)
	if !ok || id.Name != "this" {
		return nil, "", false
	}
	name, ok := memberName(mem.M)
	if !ok {
		return nil, "", false
	}
	return b, name, true
}

func reportAssignViolation(b *ast.Binary, name, method string, env TypeEnv, gt assertGroundTruth) {
	rhsT := env.GetType(b.Rhs)
	if rhsT == nil || rhsT == TVoid {
		return
	}
	members, _ := decomposeForCheck(rhsT)
	for _, m := range members {
		if !memberFits(m, gt.ty) {
			env.Report(&Diagnostic{
				Severity: SeverityError,
				Method:   method,
				Pos:      nodePos(b.Lhs),
				Msg: fmt.Sprintf("member .%s is a %v ground truth (Assert in %s), "+
					"cannot assign %v", name, gt.ty, gt.method, rhsT),
			})
			return
		}
	}
}

func forEachAssert(fn *ast.Function, visit func(cond ast.Expr, pos int)) {
	for _, stmt := range fn.Body {
		if stmt != nil {
			findAsserts(stmt, visit)
		}
	}
}

func findAsserts(n ast.Node, visit func(cond ast.Expr, pos int)) {
	if n == nil {
		return
	}
	if call, ok := n.(*ast.Call); ok && isAssertCall(call) {
		visit(call.Args[0].E, nodePos(call.Args[0].E))
	}
	n.Children(func(c ast.Node) ast.Node {
		findAsserts(c, visit)
		return c
	})
}

func isAssertCall(call *ast.Call) bool {
	id, ok := call.Fn.(*ast.Ident)
	if !ok || id.Name != "Assert" {
		return false
	}
	return len(call.Args) == 1 && call.Args[0].Name == nil
}
