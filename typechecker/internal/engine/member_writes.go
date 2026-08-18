// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

type classMemberWrites struct {
	perMethod map[string]*methodWrites
	methods   map[string]*ast.Function
}

type methodWrites struct {
	members map[string]bool
	opaque  bool
}

// ```suneido
// if .x isnt false
//
//	{
//	.Helper()
//	^^^^^^^^^ write-set unbounded (opaque call) - invalidates the .x narrowing
//	.x.Size()
//	}
//
// ```
func ComputeMemberWrites(cls *ClassObject) *classMemberWrites {
	w := &classMemberWrites{
		perMethod: make(map[string]*methodWrites, len(cls.Methods)),
		methods:   cls.Methods,
	}
	edges := make(map[string]map[string]bool, len(cls.Methods))
	for name, fn := range cls.SortedMethods {
		mw := &methodWrites{members: map[string]bool{}}
		callees := map[string]bool{}
		for _, stmt := range fn.Body {
			if stmt != nil {
				w.scanDirect(stmt, mw, callees)
			}
		}
		w.perMethod[name] = mw
		edges[name] = callees
	}
	w.propagate(edges)
	return w
}

func (w *classMemberWrites) propagate(edges map[string]map[string]bool) {
	for changed := true; changed; {
		changed = false
		for name, mw := range w.perMethod {
			if mw.opaque {
				continue
			}
			if w.absorbCallees(mw, edges[name]) {
				changed = true
			}
		}
	}
}

func (w *classMemberWrites) absorbCallees(mw *methodWrites, callees map[string]bool) bool {
	changed := false
	for callee := range callees {
		cw := w.perMethod[callee]
		if cw == nil {
			continue
		}
		if cw.opaque {
			mw.opaque = true
			return true
		}
		for m := range cw.members {
			if !mw.members[m] {
				mw.members[m] = true
				changed = true
			}
		}
	}
	return changed
}

func (w *classMemberWrites) scanDirect(n ast.Node, mw *methodWrites, callees map[string]bool) {
	if n == nil || mw.opaque {
		return
	}
	switch x := n.(type) {
	case *ast.Binary:
		if x.Tok.IsAssign() {
			if name, _, ok := unwrapThisMember(x.Lhs); ok {
				mw.members[name] = true
			}
		}
	case *ast.Unary:
		if x.Tok == tok.Inc || x.Tok == tok.Dec ||
			x.Tok == tok.PostInc || x.Tok == tok.PostDec {
			if name, _, ok := unwrapThisMember(x.E); ok {
				mw.members[name] = true
			}
		}
	case *ast.Call:
		w.classifyCall(x, mw, callees)
	}
	n.Children(func(c ast.Node) ast.Node {
		w.scanDirect(c, mw, callees)
		return c
	})
}

func (w *classMemberWrites) classifyCall(call *ast.Call, mw *methodWrites, callees map[string]bool) {
	if thisEscapes(call) {
		mw.opaque = true
		return
	}
	switch fn := call.Fn.(type) {
	case *ast.Ident:
		return
	case *ast.Mem:
		if id, ok := fn.E.(*ast.Ident); ok {
			switch id.Name {
			case "this":
				name, ok := memberName(fn.M)
				if !ok {
					mw.opaque = true
					return
				}
				if _, defined := w.methods[name]; defined {
					callees[name] = true
					return
				}
				mw.opaque = true
				return
			case "super":
				mw.opaque = true
				return
			default:
				return
			}
		}
		// receiver is a computed expression - other object, fine.
		return
	default:
		// computed callee: `(.fns[k])(...)` may be a bound method of this.
		mw.opaque = true
	}
}

// a call that passes `this` itself hands the callee write access.
func thisEscapes(call *ast.Call) bool {
	for i := range call.Args {
		if argMentionsBareThis(call.Args[i].E) {
			return true
		}
	}
	return false
}

func argMentionsBareThis(n ast.Node) bool {
	if n == nil {
		return false
	}
	if id, ok := n.(*ast.Ident); ok {
		return id.Name == "this"
	}
	if mem, ok := n.(*ast.Mem); ok {
		if id, ok := mem.E.(*ast.Ident); ok && id.Name == "this" {
			return argMentionsBareThis(mem.M)
		}
	}
	found := false
	n.Children(func(c ast.Node) ast.Node {
		if !found && argMentionsBareThis(c) {
			found = true
		}
		return c
	})
	return found
}

func (w *classMemberWrites) callEffect(call *ast.Call) (map[string]bool, bool) {
	if w == nil {
		return nil, true
	}
	probe := &methodWrites{members: map[string]bool{}}
	callees := map[string]bool{}
	w.classifyCall(call, probe, callees)
	if probe.opaque {
		return nil, true
	}
	eff := probe.members
	for callee := range callees {
		cw := w.perMethod[callee]
		if cw == nil || cw.opaque {
			return nil, true
		}
		for m := range cw.members {
			eff[m] = true
		}
	}
	// arg subtrees: blocks (and inline assigns) the callee may run.
	for i := range call.Args {
		collectMemberWritesIn(call.Args[i].E, eff)
	}
	return eff, false
}

func collectMemberWritesIn(n ast.Node, out map[string]bool) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *ast.Binary:
		if x.Tok.IsAssign() {
			if name, _, ok := unwrapThisMember(x.Lhs); ok {
				out[name] = true
			}
		}
	case *ast.Unary:
		if x.Tok == tok.Inc || x.Tok == tok.Dec ||
			x.Tok == tok.PostInc || x.Tok == tok.PostDec {
			if name, _, ok := unwrapThisMember(x.E); ok {
				out[name] = true
			}
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		collectMemberWritesIn(c, out)
		return c
	})
}
