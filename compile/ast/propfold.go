// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/assert"
)

var trueConst = &Constant{Val: True}
var falseConst = &Constant{Val: False}

// PropFold traverses an AST and does constant propagation and folding,
// modifying the AST.
// Propagation is done on the way down the tree (top down).
// Folding is (mostly) done on the way up the tree (bottom up).
// Modifications are done by the update function of node.Children.
// https://thesoftwarelife.blogspot.com/2020/04/constant-propogation-and-folding.html
func PropFold(fn *Function) *Function {
	// Final variables (set once, not modified) are determined during parse
	propfold(fn, fn.Final)
	return fn
}

// propfold - constant propagation and folding
func propfold(fn *Function, final map[string]uint8) {
	f := fold{final: final, srcpos: -1}
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Sprintf("compile error @%d %s", f.srcpos, e))
		}
	}()
	fn.Children(f.visit)
}

type fold struct {
	Folder
	final  map[string]byte
	lvalue string
	values []pair
	srcpos int
}

func (f *fold) visit(node Node) Node {
	if stmt, ok := node.(Statement); ok {
		f.srcpos = stmt.Position() // for error reporting
	}
	f.children(node) // RECURSE
	node = f.fold(node)
	if node != nil {
		node = f.findConst(node)
	}
	return node
}

func (f *fold) children(node Node) {
	save := f.values
	switch node := node.(type) {
	case *If:
		f.childExpr(&node.Cond)
		save = f.values
		f.childStmt(&node.Then)
		f.values = save
		f.childStmt(&node.Else)
		f.values = save
		return
	case *Switch:
		f.childExpr(&node.E)
		save = f.values
		for i := range node.Cases {
			c := &node.Cases[i]
			for j := range c.Exprs {
				f.childExpr(&c.Exprs[j])
			}
			for j := range c.Body {
				f.childStmt(&c.Body[j])
			}
			f.values = save
		}
		for i := range node.Default {
			f.childStmt(&node.Default[i])
		}
		f.values = save
		return
	case *Return:
		for i := range node.Exprs {
			f.childExpr(&node.Exprs[i])
		}
		f.values = save
		return
	case *MultiAssign:
		for i := range node.Lhs {
			id := node.Lhs[i].(*Ident)
			if id.Name != "unused" {
				f.lvalue = id.Name
			}
			f.childExpr(&node.Lhs[i])
		}
		f.childExpr(&node.Rhs)
		f.values = save
		return
	case *Throw:
		f.childExpr(&node.E)
		f.values = save
		return
	case *TryCatch:
		f.childStmt(&node.Try)
		f.values = save
		f.childStmt(&node.Catch)
		f.values = save
		return
	case *Forever:
		f.childStmt(&node.Body)
		f.values = save
		return
	case *While:
		f.childExpr(&node.Cond)
		save = f.values
		f.childStmt(&node.Body)
		f.values = save
		return
	case *ForIn:
		f.childExpr(&node.E)
		f.childExpr(&node.E2)
		save = f.values
		f.childStmt(&node.Body)
		f.values = save
		return
	case *For:
		for i := range node.Init {
			f.childExpr(&node.Init[i])
		}
		save = f.values
		f.childExpr(&node.Cond)
		f.childStmt(&node.Body)
		for i := range node.Inc {
			f.childExpr(&node.Inc[i])
		}
		f.values = save
		return
	case *Binary:
		// prevent lvalues from being flagged as uninitialized
		if node.Tok == tok.Eq {
			if id, ok := node.Lhs.(*Ident); ok {
				f.lvalue = id.Name
			}
		}
	case *Nary:
		// We short circuit here rather than in folder to do it top-down
		// because bottom up evaluations could cause StrictCompare errors.
		// We discard assigned f.values from right hand side of And/Or
		// since they are only conditionally executed.
		switch node.Tok {
		case tok.And:
			f.childExpr(&node.Exprs[0])
			save = f.values
			isFalse := false
			for i := 1; i < len(node.Exprs); i++ {
				if c, ok := node.Exprs[i-1].(*Constant); ok && c.Val == False {
					isFalse = true
				}
				if isFalse {
					node.Exprs[i] = falseConst
				} else {
					f.childExpr(&node.Exprs[i])
				}
			}
			f.values = save
			return
		case tok.Or:
			f.childExpr(&node.Exprs[0])
			save = f.values
			isTrue := false
			for i := 1; i < len(node.Exprs); i++ {
				if c, ok := node.Exprs[i-1].(*Constant); ok && c.Val == True {
					isTrue = true
				}
				if isTrue {
					node.Exprs[i] = trueConst
				} else {
					f.childExpr(&node.Exprs[i])
				}
			}
			f.values = save
			return
		}
	case *Trinary:
		// discard assigned values from right hand side of Trinary
		// since they are only conditionally executed
		f.childExpr(&node.Cond)
		save = f.values
		f.childExpr(&node.T)
		f.values = save
		f.childExpr(&node.F)
		f.values = save
		return
	}
	node.Children(f.visit) // recurse
}

func (f *fold) childExpr(pexpr *Expr) {
	childExpr(f.visit, pexpr)
}

func (f *fold) childStmt(pstmt *Statement) {
	childStmt(f.visit, pstmt)
}

func (f *fold) fold(node Node) Node { // NOT recursive
	switch node := node.(type) {
	case *Ident:
		return f.ident(node)
	case *Unary:
		return f.foldUnary(node)
	case *Binary:
		return f.foldBinary(node)
	case *Trinary:
		return f.foldTrinary(node)
	case *In:
		return f.foldIn(node)
	case *Nary:
		return f.foldNary(node)
	case *If:
		return f.ifStmt(node)
	case *Call:
		return f.foldCall(node)
	}
	// TODO switch
	return node
}

func (f *fold) foldCall(node *Call) Node {
	// if calling Number?, String?, or Date?
	// and the single argument is a constant
	// evaluate it
	if id, ok := node.Fn.(*Ident); ok && len(node.Args) == 1 {
		if arg, ok := node.Args[0].E.(*Constant); ok && node.Args[0].Name == nil {
			var result Value
			switch id.Name {
			case "Number?":
				result = SuBool(arg.Val.Type() == types.Number)
			case "String?":
				t := arg.Val.Type()
				result = SuBool(t == types.String || t == types.Except)
			case "Date?":
				result = SuBool(arg.Val.Type() == types.Date)
			}
			if result != nil {
				return &Constant{Val: result}
			}
		}
	}
	return node
}

func (f *fold) ident(node *Ident) Node {
	if _, ok := f.final[node.Name]; ok {
		if val := f.value(node.Name); val != nil {
			return &Constant{Val: val}
		}
		if f.lvalue == "" {
			panic("possibly uninitialized variable: " + node.Name)
		} else {
			assert.That(node.Name == f.lvalue)
		}
	}
	f.lvalue = ""
	return node
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

func (f *fold) findConst(node Node) Node {
	if node, ok := node.(*Binary); ok {
		if node.Tok == tok.Eq {
			if id, ok := node.Lhs.(*Ident); ok {
				if _, ok := f.final[id.Name]; ok {
					if val, ok := node.Rhs.(*Constant); ok {
						f.push(id.Name, val.Val)
						return node.Rhs // remove assignment
					}
					delete(f.final, id.Name) // so not uninitialized
				}
			}
		}
	}
	return node
}

type pair struct {
	val Value
	id  string
}

func (f *fold) push(id string, val Value) {
	f.values = append(f.values, pair{id: id, val: val})
}

func (f *fold) value(id string) Value {
	for i := len(f.values) - 1; i >= 0; i-- {
		if f.values[i].id == id {
			return f.values[i].val
		}
	}
	return nil
}
