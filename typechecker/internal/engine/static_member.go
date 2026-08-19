// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/compile/ast"
)

func ClassStaticMemberTypes(cls *ClassObject) map[string]DynType {
	out := make(map[string]DynType, len(cls.Members))
	for name, val := range cls.Members {
		out[name] = DynTypeOfSuValue(val)
	}
	return out
}

// ```suneido
// b = Config.Timeout
//
//	^^^^^^^^^^^^^^ resolved against Config's class-level members
//	               (SuClass get1: own members, then the base chain)
//
// ```
func StaticMemberCheckPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	for name, fn := range cls.SortedMethods {
		targets := map[ast.Node]bool{}
		for _, stmt := range fn.Body {
			collectCallTargets(stmt, targets)
		}
		for _, stmt := range fn.Body {
			checkStaticReads(stmt, name, env, targets)
		}
	}
	return false
}

func collectCallTargets(n ast.Node, targets map[ast.Node]bool) {
	if n == nil {
		return
	}
	if call, ok := n.(*ast.Call); ok {
		if mem, ok := call.Fn.(*ast.Mem); ok {
			targets[mem] = true
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		collectCallTargets(c, targets)
		return c
	})
}

func checkStaticReads(n ast.Node, method string, env TypeEnv, targets map[ast.Node]bool) {
	if n == nil {
		return
	}
	if mem, ok := n.(*ast.Mem); ok && !targets[mem] {
		if cls, name, ok := staticReadTarget(mem); ok {
			if acc, resolved := env.ClassStaticAccessible(cls, name); resolved && !acc {
				env.Report(&Diagnostic{
					Severity: SeverityError,
					Method:   method,
					Pos:      int(mem.DotPos),
					Msg:      fmt.Sprintf("no member %q on class %s", name, cls),
				})
			}
		}
	}
	n.Children(func(c ast.Node) ast.Node {
		checkStaticReads(c, method, env, targets)
		return c
	})
}

func staticReadTarget(mem *ast.Mem) (string, string, bool) {
	id, ok := mem.E.(*ast.Ident)
	if !ok || !isGlobalIdent(id.Name) {
		return "", "", false
	}
	name, ok := memberName(mem.M)
	if !ok || looksPrivatized(name) {
		return "", "", false
	}
	return id.Name, name, true
}
