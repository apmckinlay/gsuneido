// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// ```suneido
// Query1(.q).name
//
//	^^^^ throws on every path where Query1 returned false
//
// ```
func CapabilityCheckPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for name, fn := range cls.SortedMethods {
		ctx := capCheckCtx{
			env:    env,
			method: name,
			calls:  map[ast.Node]bool{},
			writes: map[ast.Node]bool{},
		}
		for _, stmt := range fn.Body {
			collectCallTargets(stmt, ctx.calls)
			collectWriteTargets(stmt, ctx.writes)
		}
		for _, stmt := range fn.Body {
			ctx.walk(stmt, 0)
		}
	}
	return false
}

type capability int

const (
	capRead capability = iota
	capWrite
	capRange
	capCall
	capIterate
)

//	read     everything but a boolean - SuBool.Get panics before the method
//	         lookup can bind a name, every other type falls through to it
//	write    Instance, Object, Record, Sequence (no lookup fallback for Put)
//	range    String, Object, Record, Sequence
//	iterate  String, Object, Record, Sequence - OpIter tests a Go interface,
//	         so a Suneido class cannot opt in by defining Iter
//	call     Function, Class, Block, String, Object, and an Instance whose
//	         class defines Call
//
// Record folds into TObject in this lattice, so anything Record supports the
// rule must treat TObject as supporting.
func supports(t DynType, c capability) bool {
	if t == TVoid {
		return true // bottom: nothing inhabits it, so nothing to prove
	}
	if _, ok := t.(Instance); ok {
		return c != capRange && c != capIterate
	}
	p, ok := t.(Primitive)
	if !ok {
		return true
	}
	switch c {
	case capRead:
		return !isBoolean(p)
	case capWrite:
		return p == TObject || p == TSequence
	case capRange, capIterate:
		return p == TString || p == TObject || p == TSequence
	case capCall:
		return p == TString || p == TObject ||
			p == TFunction || p == TBlock || p == TClass
	}
	return true
}

func (c capability) verb() string {
	switch c {
	case capRead:
		return "readable"
	case capWrite:
		return "writable"
	}
	return "applicable"
}

func (c capability) note() string {
	switch c {
	case capRead:
		return "no members on %s"
	case capWrite:
		return "%s has no writable members"
	case capRange:
		return "%s is not rangeable"
	case capCall:
		return "%s is not callable"
	}
	return "%s is not iterable"
}

type capCheckCtx struct {
	env    TypeEnv
	method string
	calls  map[ast.Node]bool
	writes map[ast.Node]bool
}

func (c *capCheckCtx) walk(n ast.Node, pos int) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		if ep.Expr != nil {
			c.walk(ep.Expr, int(ep.Pos))
		}
		return
	}
	if s, ok := n.(ast.Statement); ok {
		pos = s.Position()
	}
	n.Children(func(child ast.Node) ast.Node {
		c.walk(child, pos)
		return child
	})
	switch x := n.(type) {
	case *ast.Mem:
		if !c.calls[x] {
			c.checkMem(x, pos)
		}
	case *ast.RangeTo:
		c.check(x.E, pos, "range", capRange)
	case *ast.RangeLen:
		c.check(x.E, pos, "range", capRange)
	case *ast.ForIn:
		if x.E2 == nil { // E2 set means for-range (`for i in 1..10`), not iteration
			c.check(x.E, pos, "iteration", capIterate)
		}
	case *ast.Call:
		c.checkCallee(x, pos)
	}
}

func (c *capCheckCtx) checkMem(mem *ast.Mem, pos int) {
	if !ownedHere(mem.E) {
		return
	}
	if p := int(mem.DotPos); p > 0 {
		pos = p
	}
	cap := capRead
	if c.writes[mem] {
		cap = capWrite
	}
	// a non-identifier member is a subscript: `ob[0]`, not `ob.name`
	subject := "subscript"
	if name, ok := memberName(mem.M); ok {
		subject = fmt.Sprintf("member %q", name)
	}
	c.check(mem.E, pos, subject, cap)
}

func (c *capCheckCtx) checkCallee(call *ast.Call, pos int) {
	if _, ok := call.Fn.(*ast.Mem); ok {
		return
	}
	if !ownedHere(call.Fn) {
		return
	}
	c.check(call.Fn, pos, "call", capCall)
}

func (c *capCheckCtx) check(e ast.Expr, pos int, subject string, cap capability) {
	target := peelParens(e)
	receiver := c.env.GetType(target)
	arms, dirty := decomposeForCheck(receiver)
	if dirty {
		return
	}
	var bad []DynType
	for _, arm := range arms {
		if !supports(arm, cap) {
			bad = append(bad, arm)
		}
	}
	if len(bad) == 0 {
		return
	}

	var msg string
	if _, isUnion := receiver.(Union); isUnion {
		msg = fmt.Sprintf("%s%s not %s on at least one path of union receiver %v; %s",
			subject, onExpr(target), cap.verb(), receiver,
			fmt.Sprintf(cap.note(), joinTypes(bad)))
	} else {
		msg = fmt.Sprintf("%s%s not %s on receiver of type %v",
			subject, onExpr(target), cap.verb(), receiver)
	}

	sev := SeverityError
	if exprGuessed(target, c.method, c.env) {
		sev = SeverityWarning
		msg += " (receiver type guessed from builtin overloads)"
	}
	c.env.Report(&Diagnostic{Severity: sev, Method: c.method, Pos: pos, Msg: msg})
}

func ownedHere(e ast.Expr) bool {
	id, ok := peelParens(e).(*ast.Ident)
	if !ok {
		return true
	}
	return id.Name != "this" && id.Name != "super" && !isGlobalIdent(id.Name)
}

func isBoolean(t DynType) bool {
	return t == TFalse || t == TTrue || t == TBoolean
}

func collectWriteTargets(n ast.Node, targets map[ast.Node]bool) {
	if n == nil {
		return
	}
	if b, ok := n.(*ast.Binary); ok && b.Tok == tok.Eq {
		if mem, ok := b.Lhs.(*ast.Mem); ok {
			targets[mem] = true
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		collectWriteTargets(c, targets)
		return c
	})
}
