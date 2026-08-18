// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import "github.com/apmckinlay/gsuneido/compile/ast"

func MethodVarTypes(cls *ClassObject, env TypeEnv) map[string]map[string]DynType {
	out := make(map[string]map[string]DynType, len(cls.Methods))
	for name, fn := range cls.SortedMethods {
		vars := map[string]DynType{}
		if t, ok := env.Returns[name]; ok && t != nil && t != TUnknown {
			vars["$return"] = t
		}
		isParam := map[string]bool{}
		for i := range fn.Params {
			p := &fn.Params[i]
			pn := p.Name.ParamName()
			isParam[pn] = true
			if t, ok := env.Params[p]; ok && t != nil && t != TUnknown {
				vars[pn] = t
			}
		}
		locals := map[string]DynType{}
		eachIdent(fn, func(id *ast.Ident) {
			if id.Implicit || isGlobalIdent(id.Name) || isParam[id.Name] {
				return
			}
			addLocal(locals, id.Name, env.GetType(id))
		})
		mergeLocals(vars, locals)
		out[name] = vars
	}
	return out
}

func addLocal(locals map[string]DynType, name string, t DynType) {
	if t == nil || t == TUnknown {
		return
	}
	if prev, ok := locals[name]; ok {
		locals[name] = U(prev, t)
	} else {
		locals[name] = t
	}
}

func mergeLocals(vars, locals map[string]DynType) {
	for n, t := range locals {
		if t != TUnknown {
			vars[n] = t
		}
	}
}

func MemberTypes(cls *ClassObject, env TypeEnv) map[string]DynType {
	out := make(map[string]DynType, len(cls.Members))
	for name := range cls.Members {
		if t, ok := env.Members[name]; ok && t != nil && t != TUnknown {
			out[name] = t
		}
	}
	return out
}

func eachIdent(fn *ast.Function, visit func(*ast.Ident)) {
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n == nil {
			return
		}
		if ep, ok := n.(*ast.ExprPos); ok {
			walk(ep.Expr)
			return
		}
		if id, ok := n.(*ast.Ident); ok {
			visit(id)
		}
		if fi, ok := n.(*ast.ForIn); ok {
			if fi.Var.Name != "" {
				visit(&fi.Var)
			}
			if fi.Var2.Name != "" {
				visit(&fi.Var2)
			}
		}
		switch x := n.(type) {
		case ast.Statement:
			x.Children(func(c ast.Node) ast.Node { walk(c); return c })
		case ast.Expr:
			x.Children(func(c ast.Node) ast.Node { walk(c); return c })
		}
	}
	for _, stmt := range fn.Body {
		walk(stmt)
	}
}
