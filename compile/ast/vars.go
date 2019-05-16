package ast

import (
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
)

// VarSet returns a set (map to bool) of variable names used in an AST
// This includes function/block parameters
func VarSet(ast Node) map[string]bool {
	vv := varVisitor{vars: map[string]bool{}}
	Traverse(ast, &vv)
	return vv.vars
}

// VarList returns VarSet converted to a list
func VarList(ast Node) []string {
	return mapToList(VarSet(ast))
}

type varVisitor struct {
	vars map[string]bool
}

func (v *varVisitor) Before(node Node) bool {
	switch node := node.(type) {
	case *Function: // only top level, outermost
		for _, p := range node.Params {
			v.vars[paramToName(p.Name)] = true
		}
	case *Ident:
		if ascii.IsLower(node.Name[0]) {
			v.vars[node.Name] = true
		}
	case *ForIn:
		v.vars[node.Var] = true
	case *TryCatch:
		if node.CatchVar != "" {
			v.vars[node.CatchVar] = true
		}
	case *Block:
		return false // don't look inside blocks
	}
	return true // process children
}

func paramToName(name string) string {
	if name[0] == '@' {
		return name[1:]
	}
	if name[0] == '.' {
		name = name[1:]
	}
	if name[0] == '_' {
		name = name[1:]
	}
	return str.UnCapitalize(name)
}

func (*varVisitor) After(Node) {
}

func mapToList(vars map[string]bool) []string {
	keys := make([]string, len(vars))
	i := 0
	for k := range vars {
		keys[i] = k
		i++
	}
	return keys
}
