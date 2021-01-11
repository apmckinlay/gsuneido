// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
)

// PropFold traverses an AST and does constant propagation and folding,
// modifying the AST
func PropFold(fn *Function) *Function {
	// Final variables (set once, not modified) are determined during parse
	propfold(fn, fn.Final)
	return fn
}

// propfold - constant propagation and folding
func propfold(fn *Function, vars map[string]int) {
	f := fold{vars: map[string]Value{}, srcpos: -1}
	for id, lev := range vars {
		if lev != -1 {
			f.vars[id] = nil
		}
	}
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Sprintf("compile error @%d %s", f.srcpos, e))
		}
	}()
	fn.Children(f.visit)
}

type fold struct {
	Folder
	vars   map[string]Value
	srcpos int
}

func (f *fold) visit(node Node) Node {
	node.Children(f.visit) // recurse
	if stmt, ok := node.(Statement); ok {
		f.srcpos = stmt.Position() // for error reporting
	}
	node = f.fold(node)
	if node != nil {
		node = f.findConst(node)
	}
	return node
}

func (f *fold) fold(node Node) Node {
	switch node := node.(type) {
	case *Ident:
		return f.ident(node)
	case *Unary:
		return f.unary(node)
	case *Binary:
		return f.binary(node)
	case *Trinary:
		return f.trinary(node)
	case *In:
		return f.in(node)
	case *Nary:
		return f.nary(node)
	case *If:
		return f.ifStmt(node)
	}
	// TODO switch
	return node
}

func (f *fold) ident(node *Ident) Node {
	if val := f.vars[node.Name]; val != nil {
		return &Constant{Val: val}
	}
	return node
}

func (f *fold) unary(u *Unary) Node {
	return f.foldUnary(u)
}

func (f *fold) binary(b *Binary) Node {
	return f.foldBinary(b)
}

func (f *fold) in(in *In) Node {
	return f.foldIn(in)
}

func (f *fold) trinary(t *Trinary) Node {
	return f.foldTrinary(t)
}

func (f *fold) ifStmt(t *If) Node {
	c, ok := t.Cond.(*Constant)
	if !ok {
		return t
	}
	if c.Val == True {
		return t.Then
	}
	if c.Val == False {
		return t.Else
	}
	panic("if requires boolean")
}

func (f *fold) nary(n *Nary) Node {
	return f.foldNary(n)
}

func (f *fold) findConst(node Node) Node {
	if node, ok := node.(*Binary); ok {
		if node.Tok == tok.Eq {
			if id, ok := node.Lhs.(*Ident); ok {
				if _, ok := f.vars[id.Name]; ok {
					if val, ok := node.Rhs.(*Constant); ok {
						f.vars[id.Name] = val.Val
						return node.Rhs // remove assignment
					}
				}
			}
		}
	}
	return node
}
