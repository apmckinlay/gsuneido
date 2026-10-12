// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
)

// this implements the Value interface for AST nodes
// to expose them to the Suneido language
// See also: astContainer

// expressions ------------------------------------------------------

func (a *Ident) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Ident")
	case "name":
		return SuStr(a.Name)
	case "pos":
		return IntVal(a.GetPos())
	case "end":
		return IntVal(a.GetEnd())
	}
	return noChildren(m)
}

func noChildren(m Value) Value {
	if StrOrEmpty(m) == "children" {
		return children{}
	}
	return nil
}

func (a *Symbol) Get(th *Thread, m Value) Value {
	if StrOrEmpty(m) == "symbol" {
		return True
	}
	return a.Constant.Get(th, m)
}

func (a *Constant) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Constant")
	case "symbol":
		return False
	case "value":
		return a.Val
	case "pos":
		return IntVal(a.GetPos())
	case "end":
		return IntVal(a.GetEnd())
	case "children":
		return children{list: []Value{a.Val}}
	}
	return nil
}

func falsePos(a Node, m Value) Value {
	switch StrOrEmpty(m) {
	case "pos":
		return False
	case "end":
		return False
	case "children":
		return newChildren(a)
	}
	return nil
}

func (a *Unary) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Unary")
	case "op":
		return SuStr(a.Tok.String())
	case "expr":
		return a.E.(Value)
	case "children":
		return newChildren(a)
	case "pos":
		return IntVal(a.GetPos())
	case "end":
		return IntVal(a.GetEnd())
	}
	return nil
}

func (a *Binary) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Binary")
	case "op":
		return SuStr(a.Tok.String())
	case "lhs":
		return a.Lhs.(Value)
	case "rhs":
		return a.Rhs.(Value)
	}
	return falsePos(a, m)
}

func (a *Trinary) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Trinary")
	case "cond":
		return a.Cond.(Value)
	case "t":
		return a.T.(Value)
	case "f":
		return a.F.(Value)
	}
	return falsePos(a, m)
}

func (a *Nary) Get(_ *Thread, m Value) Value {
	if r := get(a, m); r != nil {
		return r
	}
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Nary")
	case "op":
		return SuStr(a.Tok.String())
	}
	return falsePos(a, m)
}

func (a *RangeTo) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("RangeTo")
	case "expr":
		return a.E.(Value)
	case "from":
		return nilToFalse(a.From)
	case "to":
		return nilToFalse(a.To)
	}
	return falsePos(a, m)
}

func (a *RangeLen) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("RangeLen")
	case "expr":
		return a.E.(Value)
	case "from":
		return nilToFalse(a.From)
	case "len":
		return nilToFalse(a.Len)
	}
	return falsePos(a, m)
}

func (a *InRange) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("InRange")
	case "expr":
		return a.E.(Value)
	case "from":
		return nilToFalse(a.Org)
	case "to":
		return nilToFalse(a.End)
	case "fromToken":
		return SuStr(strings.TrimSpace(tokEcho[a.OrgTok]))
	case "toToken":
		return SuStr(strings.TrimSpace(tokEcho[a.EndTok]))
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
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Mem")
	case "expr":
		return a.E.(Value)
	case "mem":
		return a.M.(Value)
	case "dotpos":
		return IntVal(a.DotPos)
	}
	return falsePos(a, m)
}

func (a *In) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("In")
	case "expr":
		return a.E.(Value)
	case "size":
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
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Call")
	case "func":
		return a.Fn.(Value)
	case "size":
		return IntVal(len(a.Args))
	case "end":
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
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Argument")
	case "name":
		if a.Name == nil {
			return False
		}
		return a.Name
	case "expr":
		return a.E.(Value)
	case "pos":
		return zeroToFalse(a.GetPos())
	case "end":
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
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Block")
	}
	return a.Function.Get(th, m)
}

func (a *Function) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Function")
	case "params":
		return Params{Params: a.Params}
	case "pos":
		return IntVal(a.GetPos())
	case "pos1":
		return IntVal(a.Pos1)
	case "pos2":
		return IntVal(a.Pos2)
	case "end":
		return IntVal(a.GetEnd())
	case "children":
		c := newChildren(a)
		for i := range a.Params {
			if a.Params[i].DefVal != nil {
				c.list = append(c.list, a.Params[i].DefVal)
			}
		}
		return c
	case "returnannotation":
		return SuStr(a.ReturnAnnotation)
	}
	return get(a, m)
}

func (a Params) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Params")
	case "size":
		return IntVal(len(a.Params))
	case "pos":
		return False
	case "end":
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
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Param")
	case "pos":
		return IntVal(a.GetPos())
	case "end":
		return IntVal(a.GetEnd())
	case "name":
		return SuStr(a.Name.Name)
	case "unused":
		return SuBool(a.Unused)
	case "hasdef":
		return SuBool(a.DefVal != nil)
	case "defval":
		return a.DefVal
	case "annotations":
		return SuStr(a.Annotations)
	}
	return nil
}

// statements -------------------------------------------------------

func stmtGet(a Statement, m Value) Value {
	switch StrOrEmpty(m) {
	case "pos":
		return IntVal(a.GetPos())
	case "end":
		return IntVal(a.GetEnd())
	}
	return get(a.(nodeVal), m)
}

func (a *ExprStmt) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("ExprStmt")
	case "expr":
		return a.E.(Value)
	}
	return stmtGet(a, m)
}

func (a *Compound) Get(_ *Thread, m Value) Value {
	if r := get(a, m); r != nil {
		return r
	}
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Compound")
	}
	return stmtGet(a, m)
}

func (a *If) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("If")
	case "cond":
		return a.Cond.(Value)
	case "t":
		return a.Then.(Value)
	case "f":
		if a.Else == nil {
			return False
		}
		return a.Else.(Value)
	case "elseend":
		return IntVal(a.ElseEnd)
	}
	return stmtGet(a, m)
}

func (a *Switch) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Switch")
	case "expr":
		return a.E.(Value)
	case "size":
		return IntVal(len(a.Cases))
	case "pos1":
		return IntVal(a.Pos1)
	case "pos2":
		return IntVal(a.Pos2)
	case "posdef":
		return IntVal(a.PosDef)
	case "def":
		if a.Default == nil {
			return False
		}
		return &Compound{Body: a.Default}
	case "children":
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
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Case")
	case "size":
		return IntVal(len(a.Exprs))
	case "body":
		return &Compound{Body: a.Body}
	case "pos":
		return IntVal(a.GetPos())
	case "end":
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
	if r := get(a, m); r != nil {
		return r
	}
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Return")
	case "expr":
		if len(a.Exprs) == 0 {
			return False
		}
		if len(a.Exprs) > 1 {
			panic("return.expr does not handle multiple expressions")
		}
		return a.Exprs[0].(Value)
	case "throw":
		return SuBool(a.ReturnThrow)
	case "spread":
		return SuBool(a.ReturnSpread)
	}
	return stmtGet(a, m)
}

func (a *MultiAssign) Get(_ *Thread, m Value) Value {
	if r := get(a, m); r != nil {
		return r
	}
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("MultiAssign")
	case "rhs":
		return a.Rhs.(Value)
	}
	return stmtGet(a, m)
}

func (a *Throw) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Throw")
	case "expr":
		return a.E.(Value)
	}
	return stmtGet(a, m)
}

func (a *TryCatch) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("TryCatch")
	case "try":
		return a.Try.(Value)
	case "catchend":
		return IntVal(a.CatchEnd)
	case "catch":
		return nilToFalse(a.Catch)
	case "catchvar":
		if a.Catch == nil || a.CatchVar.Name == "" {
			return False
		}
		return &a.CatchVar
	case "catchpat":
		if a.Catch == nil || a.CatchFilter == "" {
			return False
		}
		return SuStr(a.CatchFilter)
	}
	return stmtGet(a, m)
}

func (a *Forever) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Forever")
	case "body":
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *ForIn) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("ForIn")
	case "var":
		return SuStr(a.Var.Name)
	case "var2":
		return SuStr(a.Var2.Name)
	case "expr":
		return a.E.(Value)
	case "expr2":
		return a.E2.(Value)
	case "body":
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *For) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("For")
	case "init":
		ob := &SuObject{}
		for _, e := range a.Init {
			ob.Add(e.(Value))
		}
		return ob
	case "cond":
		return nilToFalse(a.Cond)
	case "inc":
		ob := &SuObject{}
		for _, e := range a.Inc {
			ob.Add(e.(Value))
		}
		return ob
	case "body":
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *While) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("While")
	case "cond":
		return a.Cond.(Value)
	case "body":
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *DoWhile) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("DoWhile")
	case "cond":
		return a.Cond.(Value)
	case "body":
		return a.Body.(Value)
	}
	return stmtGet(a, m)
}

func (a *Break) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Break")
	}
	return stmtGet(a, m)
}

func (a *Continue) Get(_ *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "type":
		return SuStr("Continue")
	}
	return stmtGet(a, m)
}

func (a *ExprPos) Get(t *Thread, m Value) Value {
	switch StrOrEmpty(m) {
	case "pos":
		return IntVal(a.Pos)
	case "end":
		return IntVal(a.End)
	}
	return a.Expr.Get(t, m)
}

//-------------------------------------------------------------------

type nodeVal interface {
	Node
	Value
}

func get[T nodeVal](node T, m Value) Value {
	switch StrOrEmpty(m) {
	case "size":
		n := 0
		node.Children(func(nd Node) Node {
			n++
			return nd
		})
		return IntVal(n)
	case "children":
		return node
	}
	if i, ok := m.ToInt(); ok {
		child := Value(False)
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
