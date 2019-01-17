package ast

import "github.com/apmckinlay/gsuneido/util/ascii"

// Vars returns a list of variable names used in an AST
// This includes function/block parameters
func Vars(ast Node) []string {
	vv := varVisitor{vars: map[string]bool{}}
	Traverse(ast, &vv)
	return vv.List()
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

func (v *varVisitor) List() []string {
	keys := make([]string, len(v.vars))
	i := 0
	for k := range v.vars {
		keys[i] = k
		i++
	}
	return keys
}
