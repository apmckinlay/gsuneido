// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
)

// this implements the Value interface for AST nodes
// to expose them to the Suneido language
// See also: astContainer

// expressions ------------------------------------------------------

func (a *Ident) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Ident")
	case SuStr("name"):
		return SuStr(a.Name)
	case SuStr("pos"):
		return IntVal(a.GetPos())
	case SuStr("end"):
		return IntVal(a.GetEnd())
	}
	return noChildren(m)
}

func noChildren(m Value) Value {
	if m == SuStr("children") {
		return children{}
	}
	return nil
}

func (a *Symbol) Get(th *Thread, m Value) Value {
	if m == SuStr("symbol") {
		return True
	}
	return a.Constant.Get(th, m)
}

func (a *Constant) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Constant")
	case SuStr("symbol"):
		return False
	case SuStr("value"):
		return a.Val
	case SuStr("pos"):
		return IntVal(a.GetPos())
	case SuStr("end"):
		return IntVal(a.GetEnd())
	case SuStr("children"):
		return children{list: []Value{a.Val}}
	}
	return nil
}

func falsePos(a Node, m Value) Value {
	switch m {
	case SuStr("pos"):
		return False
	case SuStr("end"):
		return False
	case SuStr("children"):
		return newChildren(a)
	}
	return nil
}

func (a *Unary) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Unary")
	case SuStr("op"):
		return SuStr(a.Tok.String())
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("children"):
		return newChildren(a)
	case SuStr("pos"):
		return IntVal(a.GetPos())
	case SuStr("end"):
		return IntVal(a.GetEnd())
	}
	return nil
}

func (a *Binary) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Binary")
	case SuStr("op"):
		return SuStr(a.Tok.String())
	case SuStr("lhs"):
		return a.Lhs.(Value)
	case SuStr("rhs"):
		return a.Rhs.(Value)
	}
	return falsePos(a, m)
}

func (a *Trinary) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Trinary")
	case SuStr("cond"):
		return a.Cond.(Value)
	case SuStr("t"):
		return a.T.(Value)
	case SuStr("f"):
		return a.F.(Value)
	}
	return falsePos(a, m)
}

func (a *Nary) Get(_ *Thread, m Value) Value {
	if r := get(a, m); r != nil {
		return r
	}
	switch m {
	case SuStr("type"):
		return SuStr("Nary")
	case SuStr("op"):
		return SuStr(a.Tok.String())
	}
	return falsePos(a, m)
}

func (a *RangeTo) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("RangeTo")
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("from"):
		return nilToFalse(a.From)
	case SuStr("to"):
		return nilToFalse(a.To)
	}
	return falsePos(a, m)
}

func (a *RangeLen) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("RangeLen")
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("from"):
		return nilToFalse(a.From)
	case SuStr("len"):
		return nilToFalse(a.Len)
	}
	return falsePos(a, m)
}

func nilToFalse(node Node) Value {
	if node == nil {
		return False
	}
	return node.(Value)
}

func (a *Mem) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Mem")
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("mem"):
		return a.M.(Value)
	case SuStr("dotpos"):
		return IntVal(int(a.DotPos))
	}
	return falsePos(a, m)
}

func (a *In) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("In")
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("size"):
		return IntVal(len(a.Exprs))
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(a.Exprs) <= i {
			return False
		}
		return a.Exprs[i].(Value)
	}
	return falsePos(a, m)
}

func (a *Call) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Call")
	case SuStr("func"):
		return a.Fn.(Value)
	case SuStr("size"):
		return IntVal(len(a.Args))
	case SuStr("end"):
		return zeroToFalse(int(a.End))
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(a.Args) <= i {
			return False
		}
		return &a.Args[i]
	}
	return falsePos(a, m)
}

func (a *Arg) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Argument")
	case SuStr("name"):
		if a.Name == nil {
			return False
		}
		return a.Name
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("pos"):
		return zeroToFalse(a.GetPos())
	case SuStr("end"):
		return zeroToFalse(a.GetEnd())
	}
	return nil
}

func zeroToFalse(n int) Value {
	if n == 0 {
		return False
	}
	return IntVal(n)
}

func (a *Block) Get(th *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Block")
	}
	return a.Function.Get(th, m)
}

func (a *Function) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Function")
	case SuStr("params"):
		return Params{Params: a.Params}
	case SuStr("pos"):
		return IntVal(a.GetPos())
	case SuStr("pos1"):
		return IntVal(int(a.Pos1))
	case SuStr("pos2"):
		return IntVal(int(a.Pos2))
	case SuStr("end"):
		return IntVal(a.GetEnd())
	case SuStr("children"):
		c := newChildren(a)
		for i := range a.Params {
			if a.Params[i].DefVal != nil {
				c.list = append(c.list, a.Params[i].DefVal)
			}
		}
		return c
	}
	return get(a, m)
}

func (a Params) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Params")
	case SuStr("size"):
		return IntVal(len(a.Params))
	case SuStr("pos"):
		return False
	case SuStr("end"):
		return False
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(a.Params) <= i {
			return False
		}
		return &a.Params[i]
	}
	return nil
}

func (a *Param) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Param")
	case SuStr("pos"):
		return IntVal(a.GetPos())
	case SuStr("end"):
		return IntVal(a.GetEnd())
	case SuStr("name"):
		return SuStr(a.Name.Name)
	case SuStr("unused"):
		return SuBool(a.Unused)
	case SuStr("hasdef"):
		return SuBool(a.DefVal != nil)
	case SuStr("defval"):
		return a.DefVal
	}
	return nil
}

// statements -------------------------------------------------------

func stmtGet(a Statement, m Value) Value {
	switch m {
	case SuStr("pos"):
		return IntVal(a.GetPos())
	case SuStr("end"):
		return IntVal(a.GetEnd())
	}
	return get(a.(nodeVal), m)
}

func (a *ExprStmt) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("ExprStmt")
	case SuStr("expr"):
		return a.E.(Value)
	}
	return stmtGet(a, m)
}

func (a *Compound) Get(_ *Thread, m Value) Value {
	if r := get(a, m); r != nil {
		return r
	}
	switch m {
	case SuStr("type"):
		return SuStr("Compound")
	}
	return stmtGet(a, m)
}

func (a *If) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("If")
	case SuStr("cond"):
		return a.Cond.(Value)
	case SuStr("t"):
		return a.Then.(Value)
	case SuStr("f"):
		if a.Else == nil {
			return False
		}
		return a.Else.(Value)
	case SuStr("elseend"):
		return IntVal(int(a.ElseEnd))
	}
	return stmtGet(a, m)
}

func (a *Switch) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Switch")
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("size"):
		return IntVal(len(a.Cases))
	case SuStr("pos1"):
		return IntVal(int(a.Pos1))
	case SuStr("pos2"):
		return IntVal(int(a.Pos2))
	case SuStr("posdef"):
		return IntVal(int(a.PosDef))
	case SuStr("def"):
		if a.Default == nil {
			return False
		}
		return &Compound{Body: a.Default}
	case SuStr("children"):
		return newChildren(a)
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(a.Cases) <= i {
			return False
		}
		return &a.Cases[i]
	}
	return stmtGet(a, m)
}

func (a *Case) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Case")
	case SuStr("size"):
		return IntVal(len(a.Exprs))
	case SuStr("body"):
		return &Compound{Body: a.Body}
	case SuStr("pos"):
		return IntVal(a.GetPos())
	case SuStr("end"):
		return IntVal(a.GetEnd())
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(a.Exprs) <= i {
			return False
		}
		return a.Exprs[i].(Value)
	}
	return nil
}

func (a *Return) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Return")
	case SuStr("expr"):
		return nilToFalse(a.E)
	}
	return stmtGet(a, m)
}

func (a *Throw) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Throw")
	case SuStr("expr"):
		return a.E.(Value)
	}
	return stmtGet(a, m)
}

func (a *TryCatch) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("TryCatch")
	case SuStr("try"):
		return a.Try.(Value)
	case SuStr("catchend"):
		return IntVal(int(a.CatchEnd))
	case SuStr("catch"):
		return nilToFalse(a.Catch)
	case SuStr("catchvar"):
		if a.Catch == nil || a.CatchVar.Name == "" {
			return False
		}
		return &a.CatchVar
	case SuStr("catchpat"):
		if a.Catch == nil || a.CatchFilter == "" {
			return False
		}
		return SuStr(a.CatchFilter)
	}
	return stmtGet(a, m)
}

func (a *Forever) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Forever")
	case SuStr("body"):
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *ForIn) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("ForIn")
	case SuStr("var"):
		return SuStr(a.Var.Name)
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("body"):
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *For) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("For")
	case SuStr("init"):
		ob := &SuObject{}
		for _, e := range a.Init {
			ob.Add(e.(Value))
		}
		return ob
	case SuStr("cond"):
		return nilToFalse(a.Cond)
	case SuStr("inc"):
		ob := &SuObject{}
		for _, e := range a.Inc {
			ob.Add(e.(Value))
		}
		return ob
	case SuStr("body"):
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *While) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("While")
	case SuStr("cond"):
		return a.Cond.(Value)
	case SuStr("body"):
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *DoWhile) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("DoWhile")
	case SuStr("cond"):
		return a.Cond.(Value)
	case SuStr("body"):
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *Break) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Break")
	}
	return stmtGet(a, m)
}

func (a *Continue) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Continue")
	}
	return stmtGet(a, m)
}

//-------------------------------------------------------------------

type nodeVal interface {
	Node
	Value
}

func get[T nodeVal](node T, m Value) Value {
	if m == SuStr("size") {
		n := 0
		node.Children(func(nd Node) Node {
			n++
			return nd
		})
		return IntVal(n)
	} else if m == SuStr("children") {
		return node
	}
	if i, ok := m.ToInt(); ok {
		child := False
		node.Children(func(nd Node) Node {
			if i == 0 {
				child = nd.(Value)
			}
			i--
			return nd
		})
		return child
	}
	return nil // not handled
}

// Value interface

type SuAstNode struct {
	ValueBase[SuAstNode]
}

func (SuAstNode) Type() types.Type {
	return types.AstNode
}

func (SuAstNode) Equal(any) bool {
	return false
}

func (SuAstNode) SetConcurrent() {
	// read-only so nothing to do
}

// children
type children struct {
	SuAstNode
	list []Value
}

func newChildren(node Node) children {
	var list []Value
	node.Children(func(child Node) Node {
		list = append(list, child.(Value))
		return child
	})
	return children{list: list}
}

func (a children) Get(_ *Thread, m Value) Value {
	if i, ok := m.ToInt(); ok {
		if i >= len(a.list) {
			return False
		}
		return a.list[i]
	}
	return nil
}

func NewChildren(list []Value) children {
	return children{list: list}
}
