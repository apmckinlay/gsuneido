// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"maps"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// rhsOfMemberAssign returns the right-hand side of the first `.name = e` in fn
func rhsOfMemberAssign(fn *ast.Function, name string) ast.Expr {
	var found ast.Expr
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n == nil || found != nil {
			return
		}
		if m, rhs, ok := unwrapThisAssign(n); ok && m == name {
			found = rhs
			return
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

// a parameter default only says what the caller may omit; inside a
// predicate guard the assigned value is the predicate's type
func TestMemberAssignUnderPredicateGuardIgnoresParamDefault(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		dirty: true
		Dirty?(dirty = "") {
			if Boolean?(dirty)
				.dirty = dirty
			return .dirty
		}
	}`, "T")
	a.This(env.Members["dirty"]).Is(TBoolean)
	a.This(env.Returns["Dirty?"]).Is(TBoolean)

	_, env = runPasses(`class {
		dirty: true
		Dirty?(dirty = 0) {
			if Boolean?(dirty)
				.dirty = dirty
		}
	}`, "T")
	a.This(env.Members["dirty"]).Is(TBoolean)
}

// an unannotated parameter is unknown, and a predicate on unknown proves the
// predicate type, so the member is clean
func TestMemberAssignUnderPredicateGuardUnannotatedParam(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		dirty: true
		Dirty?(dirty) {
			if Boolean?(dirty)
				.dirty = dirty
		}
	}`, "T")
	a.This(env.Members["dirty"]).Is(TBoolean)
}

// the local's type only arrives once .Get() resolves inside the first
// fixpoint, so the guard has to be re-applied on every round or the
// pre-resolution false arm ratchets into .m
func TestMemberAssignUnderGuardOnLateResolvedLocal(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		m: 0
		flag: false
		Compute() {
			if .flag is true { return "a" }
			return false
		}
		Get() { return .Compute() }
		Set() {
			v = .Get()
			if v isnt false
				.m = v
		}
	}`, "T")
	u, ok := env.Members["m"].(Union)
	a.That(ok)
	a.That(u.Contains(TNumber) && u.Contains(TString))
	a.That(!u.Contains(TFalse))
	a.That(!u.IsDirty)
}

// the guarded assignment removes the only path that could put false back,
// so New's overwrite demotes .x; before the fix .f ratcheted to ?|string
// and blocked the demotion
func TestMemberAssignUnderGuardAllowsDemotion(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: false
		f: "a"
		New() { .x = .f }
		M0(p0) {
			if String?(p0)
				.f = p0
		}
	}`, "T")
	a.This(env.Members["f"]).Is(TString)
	a.This(env.PostCtorMembers["x"]).Is(TString)
}

// narrowing a plain false away from false proves nothing about the value,
// so the member takes the unknown rather than a false arm the guard excludes
func TestMemberAssignUnderIsntFalseOnFalseParamIsDirty(t *testing.T) {
	a := assert.T(t)
	_, env := runPasses(`class {
		x: false
		Set(v = false) {
			if v isnt false
				.x = v
		}
	}`, "T")
	got := env.Members["x"]
	a.That(isDirty(got))
	a.That(containsFalse(got))
}

// the pass refines locals and nothing else: a member guard in the same
// method leaves both the member's reads and env.Members untouched
func TestLocalNarrowingPassLeavesMembersAlone(t *testing.T) {
	a := assert.T(t)
	cls := NewClassObject("T", ParseClass(`class {
		x: false
		y: 0
		z: 0
		Foo(p, .w) {
			if Number?(p)
				.y = p
			if .x isnt false
				.z = .x
			if Number?(.w)
				.y = .w
		}
	}`))
	env := NewTypeEnv().WithClass(cls, buildMethodSigs(cls))
	pctx := NewPassCtx()
	LocalInference(cls, env, pctx)
	NameResolutionPass(cls, env, pctx)
	before := map[string]DynType{}
	maps.Copy(before, env.Members)
	LocalNarrowingPass(cls, env, pctx)

	fn := cls.Methods["Foo"]
	a.This(env.GetType(rhsOfMemberAssign(fn, "y"))).Is(TNumber)
	a.This(env.GetType(rhsOfMemberAssign(fn, "z"))).Is(TFalse)
	a.This(len(env.Members)).Is(len(before))
	for k, v := range before {
		a.This(env.Members[k]).Is(v)
	}
}
