// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
)

// ```suneido
// Foo(a, b = 1) { }
// .Foo()
//
//	^^^^^ missing argument: a
//
// .Foo(1, 2, 3)
//
//	^^^^^^^^^^^^ too many arguments
//
// ```
func ArityCheckPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for name, fn := range cls.SortedMethods {
		if fn == nil {
			continue
		}
		ctx := arityCtx{env: env, method: name}
		for _, stmt := range fn.Body {
			if stmt != nil {
				ctx.walk(stmt, 0)
			}
		}
	}
	return false
}

type arityCtx struct {
	env    TypeEnv
	method string
}

// Calls are checked on the way up
func (c *arityCtx) walk(n ast.Node, pos int) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		newPos := int(ep.Pos)
		if ep.Expr != nil {
			c.walk(ep.Expr, newPos)
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
	if call, ok := n.(*ast.Call); ok {
		c.checkArity(call, pos)
	}
}

func (c *arityCtx) checkArity(call *ast.Call, pos int) {
	sig := c.env.CallSig(call)
	if sig == nil || sig.AtParam {
		return
	}

	shape, ok := argShape(call)
	if !ok {
		return // an argument is an @ / @+1 spread; positional count is unknown
	}

	sev := SeverityError
	if c.env.GuessedCall(call) || callReceiverGuessed(call, c.method, c.env) {
		sev = SeverityWarning
	}
	name := calleeName(call)

	if shape.positional > len(sig.Params) {
		c.env.Report(&Diagnostic{
			Severity: sev,
			Method:   c.method,
			Pos:      pos,
			Msg: fmt.Sprintf("too many arguments to %s: %d given, takes at most %d",
				name, shape.positional, len(sig.Params)),
		})
		return
	}

	for i := range sig.Params {
		p := &sig.Params[i]
		if i < shape.positional && shape.named[p.Name] {
			c.env.Report(&Diagnostic{
				Severity:   SeverityWarning,
				Method:     c.method,
				Pos:        pos,
				Confidence: 0.70,
				Msg: fmt.Sprintf("argument %q to %s is passed both positionally and by name; the named value overrides the positional one (this may be confusing)",
					p.Name, name),
			})
		}
	}

	var missing []string
	for i := range sig.Params {
		p := &sig.Params[i]
		if p.HasDefault || i < shape.positional || shape.named[p.Name] {
			continue
		}
		missing = append(missing, p.Name)
	}
	if len(missing) > 0 {
		plural := ""
		if len(missing) > 1 {
			plural = "s"
		}
		c.env.Report(&Diagnostic{
			Severity: sev,
			Method:   c.method,
			Pos:      pos,
			Msg: fmt.Sprintf("missing argument%s to %s: %s",
				plural, name, strings.Join(missing, ", ")),
		})
	}
}

type callArgShape struct {
	positional int
	named      map[string]bool
}

func argShape(call *ast.Call) (callArgShape, bool) {
	s := callArgShape{named: map[string]bool{}}
	for i := range call.Args {
		arg := &call.Args[i]
		if isAtArg(arg) {
			return s, false
		}
		if arg.Name != nil {
			// a non-string name binds to no parameter, and is not positional
			if name, ok := argName(arg); ok {
				s.named[name] = true
			}
			continue
		}
		s.positional++
	}
	return s, true
}

// per-callee-shape dispatch, one arm per call form
func arityCalleeSig(call *ast.Call, d *callDispatch, env TypeEnv) *Signature {
	if d.Sig != nil {
		return d.Sig
	}
	switch fn := call.Fn.(type) {
	case *ast.Ident:
		if fn.Name == "super" {
			return inheritedMethodSig(env, env.ClassBase, "New") // super(...) -> base New
		}

		// bare X(...) invokes CallClass.
		if fn.Name == env.ClassName {
			return classCallSig(env)
		}

		if isGlobalIdent(fn.Name) {
			return refClassCallSig(env, fn.Name) // a reference class
		}
	case *ast.Mem:
		name, ok := memberName(fn.M)
		if !ok {
			return nil
		}
		if name == "*new*" { // `new X(...)` bypasses CallClass -> New
			id, ok := fn.E.(*ast.Ident)
			if !ok {
				return nil
			}
			if id.Name == "this" || id.Name == env.ClassName {
				return newExprSig(env)
			}
			if isGlobalIdent(id.Name) {
				return refMethodSig(env, id.Name, "New")
			}
			return nil
		}
		if id, ok := fn.E.(*ast.Ident); ok {
			if id.Name == "this" {
				if s := env.MethodSigs[name]; s != nil {
					return s // own method (overrides any inherited one)
				}
				return inheritedMethodSig(env, env.ClassBase, name)
			}
			if id.Name == "super" {
				return inheritedMethodSig(env, env.ClassBase, name) // super.M(...)
			}
			if isGlobalIdent(id.Name) {
				return refMethodSig(env, id.Name, name) // Foo.Bar(...) static
			}
		}
		// inst.Method(...) on a constructed instance
		if inst, ok := env.GetType(fn.E).(Instance); ok {
			if inst.Class == env.ClassName {
				return env.MethodSigs[name]
			}
			return refMethodSig(env, inst.Class, name)
		}
	}
	return nil
}

func refMethodSig(env TypeEnv, class, method string) *Signature {
	if m := env.ClassMethodSigs[class]; m != nil {
		return m[method]
	}
	return nil
}

// bare Foo(...) binds the way SuClass.Call does: CallClass anywhere in the
// chain wins, and only if there is none does it construct, binding New - also
// inherited. inheritedMethodSig bails to nil at the first class we hold no
// signatures for, so an unloaded base skips the check rather than guessing.
func refClassCallSig(env TypeEnv, class string) *Signature {
	m := env.ClassMethodSigs[class]
	if m == nil {
		return nil
	}
	if s := inheritedMethodSig(env, class, "CallClass"); s != nil {
		return s
	}
	if s := inheritedMethodSig(env, class, "New"); s != nil {
		return s
	}
	return m[class]
}

func inheritedMethodSig(env TypeEnv, base, method string) *Signature {
	for c := base; c != ""; c = env.ClassBases[c] {
		m := env.ClassMethodSigs[c]
		if m == nil {
			return nil
		}
		if s := m[method]; s != nil {
			return s
		}
	}
	return nil
}

var emptySig = Signature{}

func classCallSig(env TypeEnv) *Signature {
	if s := env.MethodSigs["CallClass"]; s != nil {
		return s
	}
	return newExprSig(env)
}

func newExprSig(env TypeEnv) *Signature {
	if s := env.MethodSigs["New"]; s != nil {
		return s
	}
	if env.ClassBase == "" {
		return &emptySig
	}
	return nil
}

func calleeName(call *ast.Call) string {
	switch fn := call.Fn.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.Mem:
		if name, ok := memberName(fn.M); ok {
			if name == "*new*" {
				if id, ok := fn.E.(*ast.Ident); ok {
					return "new " + id.Name
				}
				return "constructor"
			}
			return name
		}
	}
	return "function"
}

func ClassMethodSignatures(cls *ClassObject) map[string]*Signature {
	return buildMethodSigs(cls)
}

func buildMethodSigs(cls *ClassObject) map[string]*Signature {
	m := make(map[string]*Signature, len(cls.Methods))
	for name, fn := range cls.SortedMethods {
		if fn == nil || isAbstractStub(fn) {
			continue
		}
		sig, err := signatureFromAst(fn)
		if err != nil {
			continue
		}
		s := sig
		m[name] = &s
	}
	return m
}

// isAbstractStub reports a method whose entire body is a single `throw`
func isAbstractStub(fn *ast.Function) bool {
	if len(fn.Body) != 1 {
		return false
	}
	_, ok := fn.Body[0].(*ast.Throw)
	return ok
}
