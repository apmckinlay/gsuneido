// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

// this implements the Value interface for AST nodes
// to expose them to the Suneido language

// expressions ------------------------------------------------------

func (a *Ident) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Ident")
	case SuStr("name"):
		return SuStr(a.Name)
	}
	return nil
}

func (a *Symbol) Get(t *Thread, m Value) Value {
	if m == SuStr("symbol") {
		return True
	}
	return a.Constant.Get(t, m)
}

func (a *Constant) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Constant")
	case SuStr("symbol"):
		return False
	case SuStr("value"):
		return a.Val
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
	return nil
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
	return nil
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
	return nil
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
	return nil
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
	return nil
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
	}
	return nil
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
	return nil
}

func (a *Call) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Call")
	case SuStr("func"):
		return a.Fn.(Value)
	case SuStr("size"):
		return IntVal(len(a.Args))
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(a.Args) <= i {
			return False
		}
		arg := a.Args[i]
		if arg.Name == nil {
			return NewSuObject([]Value{arg.E.(Value)})
		}
		return NewSuObject([]Value{arg.E.(Value), arg.Name})
	}
	return nil
}

func (a *Block) Get(t *Thread, m Value) Value {
	if m == SuStr("type") {
		return SuStr("Block")
	}
	return a.Function.Get(t, m)
}

func (a *Function) Get(_ *Thread, m Value) Value {
	if r := get(a, m); r != nil {
		return r
	}
	switch m {
	case SuStr("type"):
		return SuStr("Function")
	case SuStr("params"):
		ob := &SuObject{}
		for _, p := range a.Params {
			name := p.Name.Name
			if p.Unused && strings.TrimLeft(name, "@") != "unused" {
				name += "/*unused*/"
			}
			if p.DefVal == nil {
				ob.Add(NewSuObject([]Value{SuStr(name)}))
			} else {
				ob.Add(NewSuObject([]Value{SuStr(name), p.DefVal}))
			}
		}
		return ob
	}
	return nil
}

// statements -------------------------------------------------------

func (a *ExprStmt) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Expr")
	case SuStr("expr"):
		return a.E.(Value)
	}
	return nil
}

func (a *Compound) Get(_ *Thread, m Value) Value {
	if r := get(a, m); r != nil {
		return r
	}
	switch m {
	case SuStr("type"):
		return SuStr("Compound")
	}
	return nil
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
	}
	return nil
}

func (a *Switch) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Switch")
	case SuStr("expr"):
		return a.E.(Value)
	case SuStr("size"):
		return IntVal(len(a.Cases))
	case SuStr("def"):
		if a.Default == nil {
			return False
		}
		return &Compound{Body: a.Default}
	}
	if i, ok := m.ToInt(); ok {
		if i < 0 || len(a.Cases) <= i {
			return False
		}
		ob := &SuObject{}
		for _, e := range a.Cases[i].Exprs {
			ob.Add(e.(Value))
		}
		ob.Add(&Compound{Body: a.Cases[i].Body})
		return ob
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
	return nil
}

func (a *Throw) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Throw")
	case SuStr("expr"):
		return a.E.(Value)
	}
	return nil
}

func (a *TryCatch) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("TryCatch")
	case SuStr("try"):
		return a.Try.(Value)
	case SuStr("catch"):
		if a.Catch == nil {
			return False
		}
		ob := &SuObject{}
		ob.Add(a.Catch.(Value))
		if a.CatchVar.Name != "" {
			s := a.CatchVar.Name
			if a.CatchVarUnused && a.CatchVar.Name != "unused" {
				s += "/*unused*/"
			}
			ob.Add(SuStr(s))
			if a.CatchFilter != "" {
				ob.Add(SuStr(a.CatchFilter))
			}
		}
		return ob
	}
	return nil
}

func (a *Forever) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Forever")
	case SuStr("body"):
		return a.Body.(Value)
	}
	return nil
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
	return nil
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
	return nil
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
	return nil
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
	return nil
}

func (a *Break) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Break")
	}
	return nil
}

func (a *Continue) Get(_ *Thread, m Value) Value {
	switch m {
	case SuStr("type"):
		return SuStr("Continue")
	}
	return nil
}

func get(node Node, m Value) Value {
	if m == SuStr("size") {
		n := 0
		node.Children(func(nd Node) Node {
			n++
			return nd
		})
		return IntVal(n)
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

type AstNodeValue struct {
	CantConvert
}

func (AstNodeValue) Put(*Thread, Value, Value) {
	panic("AstNode does not support put")
}

func (AstNodeValue) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("AstNode does not support update")
}

func (AstNodeValue) RangeTo(int, int) Value {
	panic("AstNode does not support range")
}

func (AstNodeValue) RangeLen(int, int) Value {
	panic("AstNode does not support range")
}

func (AstNodeValue) Hash() uint32 {
	panic("AstNode hash not implemented")
}

func (AstNodeValue) Hash2() uint32 {
	panic("AstNode hash not implemented")
}

func (AstNodeValue) Compare(Value) int {
	panic("AstNode compare not implemented")
}

func (AstNodeValue) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call AstNode")
}

func (AstNodeValue) String() string {
	return "astNode"
}

func (AstNodeValue) Type() types.Type {
	return types.AstNode
}

func (a AstNodeValue) Equal(other interface{}) bool {
	a2, ok := other.(AstNodeValue)
	return ok && a == a2
}

func (AstNodeValue) Lookup(*Thread, string) Callable {
	return nil // no methods
}
