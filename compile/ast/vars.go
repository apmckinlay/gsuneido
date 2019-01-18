package ast

import "github.com/apmckinlay/gsuneido/util/ascii"

// VarSet returns a set (map to bool) of variable names used in an AST
// This includes function/block parameters
func VarSet(ast Node) map[string]bool {
	vv := varVisitor{vars: map[string]bool{}}
	Traverse(ast, &vv)
	return vv.vars
}

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
			v.vars[p.Name] = true
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
