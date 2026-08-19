// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
)

type boolPost struct {
	onTrue  map[string]DynType
	onFalse map[string]DynType
}

type boolPostconds map[string]*boolPost

// ```suneido
// hasCache?() { return .cache isnt false }
// Get() { if .hasCache?() { return .cache.Size() } }
//
//	^^^^^^^^^^^ the true branch inherits the helper's
//	            postcondition: .cache is non-false in the guard
//
// ```
func ComputeBoolPostconditions(cls *ClassObject, env TypeEnv,
	assignedFalse map[string]bool, writes *classMemberWrites) boolPostconds {

	out := boolPostconds{}
	for name, fn := range cls.SortedMethods {
		if isPublicMember(name) || name == "New" {
			continue
		}
		if !hasLiteralBoolReturn(fn) {
			continue
		}
		c := &postCollector{env: env}
		sc := initialNarrowScope(fn, env)
		sc.memberAssignedFalse = assignedFalse
		sc.writes = writes
		sc.postHook = c.collect
		walkBlock(fn.Body, env, sc)
		if !stmtsAlwaysExit(fn.Body) {
			c.join(&c.trueExits, sc)
			c.join(&c.falseExits, sc)
		}
		if post := c.finish(env); post != nil {
			out[name] = post
		}
	}
	return out
}

func hasLiteralBoolReturn(fn *ast.Function) bool {
	found := false
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n == nil || found {
			return
		}
		if r, ok := n.(*ast.Return); ok && !r.ReturnThrow && len(r.Exprs) == 1 {
			if _, ok := returnLiteralPolarity(r); ok {
				found = true
				return
			}
		}
		n.Children(func(c ast.Node) ast.Node {
			walk(c)
			return c
		})
	}
	for _, stmt := range fn.Body {
		walk(stmt)
	}
	return found
}

func returnLiteralPolarity(r *ast.Return) (bool, bool) {
	if r.ReturnThrow || len(r.Exprs) != 1 {
		return false, false
	}
	c, ok := unwrapConstant(r.Exprs[0])
	if !ok {
		return false, false
	}
	switch DynTypeOfSuValue(c.Val) {
	case TTrue:
		return true, true
	case TFalse:
		return false, true
	}
	return false, false
}

type postCollector struct {
	env        TypeEnv
	trueExits  []map[string]DynType
	falseExits []map[string]DynType
}

func (c *postCollector) collect(r *ast.Return, sc narrowScope) {
	if polarity, ok := returnLiteralPolarity(r); ok {
		if polarity {
			c.join(&c.trueExits, sc)
		} else {
			c.join(&c.falseExits, sc)
		}
		return
	}
	c.join(&c.trueExits, sc)
	c.join(&c.falseExits, sc)
}

func (c *postCollector) join(dst *[]map[string]DynType, sc narrowScope) {
	snap := make(map[string]DynType, len(sc.Members))
	for k, f := range sc.Members {
		snap[k] = f.Typ
	}
	*dst = append(*dst, snap)
}

func (c *postCollector) finish(env TypeEnv) *boolPost {
	onTrue := foldExitSnapshots(c.trueExits, env)
	onFalse := foldExitSnapshots(c.falseExits, env)
	if len(onTrue) == 0 && len(onFalse) == 0 {
		return nil
	}
	return &boolPost{onTrue: onTrue, onFalse: onFalse}
}

func foldExitSnapshots(snaps []map[string]DynType, env TypeEnv) map[string]DynType {
	if len(snaps) == 0 {
		return nil
	}
	names := map[string]bool{}
	for _, s := range snaps {
		for k := range s {
			names[k] = true
		}
	}
	out := map[string]DynType{}
	for name := range names {
		classWide, ok := env.LookupMember(name)
		if !ok {
			continue
		}
		merged, sound := mergeSnapshots(snaps, name, classWide)
		if !sound || merged == nil {
			continue
		}
		if !isNarrower(merged, classWide) || typeHasBoolish(merged) {
			continue
		}
		out[name] = merged
	}
	return out
}

func mergeSnapshots(snaps []map[string]DynType, name string, classWide DynType) (DynType, bool) {
	var merged DynType
	for _, s := range snaps {
		t, refined := s[name]
		if !refined {
			t = classWide
		}
		if t == nil || t == TUnknown {
			return nil, false
		}
		if merged == nil {
			merged = t
		} else {
			merged = U(merged, t)
		}
	}
	return merged, true
}

func stmtsAlwaysExit(stmts []ast.Statement) bool {
	for _, s := range stmts {
		if branchAlwaysExits(s) {
			return true
		}
	}
	return false
}

func refineHelperPostcondition(call *ast.Call, sc narrowScope, polarity bool, allowMembers bool) {
	if sc.postconds == nil || !allowMembers {
		return
	}
	mem, ok := call.Fn.(*ast.Mem)
	if !ok {
		return
	}
	id, ok := mem.E.(*ast.Ident)
	if !ok || id.Name != "this" {
		return
	}
	name, ok := memberName(mem.M)
	if !ok {
		return
	}
	post := sc.postconds[name]
	if post == nil {
		return
	}
	members := post.onTrue
	if !polarity {
		members = post.onFalse
	}
	for m, t := range members {
		storeRefinement(sc, narrowTarget{name: m, isMember: true}, t, true)
	}
}
