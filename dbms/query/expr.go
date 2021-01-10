// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/runtime"
)

type Expr = ast.Expr

type Builder struct{}

var _ ast.Factory = Builder{}

func (Builder) Ident(name string, pos int32) Expr {
	return &Ident{ast.Ident{Name: name, Pos: pos}}
}
func (Builder) Constant(val runtime.Value) Expr {
	return &Constant{ast.Constant{Val: val}}
}
func (Builder) Unary(tok tok.Token, expr Expr) Expr {
	return &Unary{ast.Unary{Tok: tok, E: expr}}
}
func (Builder) Binary(lhs Expr, tok tok.Token, rhs Expr) Expr {
	return &Binary{ast.Binary{Lhs: lhs, Tok: tok, Rhs: rhs}}
}
func (Builder) Trinary(cond Expr, t Expr, f Expr) Expr {
	return &Trinary{ast.Trinary{Cond: cond, T: t, F: f}}
}
func (Builder) Nary(tok tok.Token, exprs []Expr) Expr {
	return &Nary{ast.Nary{Tok: tok, Exprs: exprs}}
}
func (Builder) Mem(e Expr, m Expr) Expr {
	return &Mem{ast.Mem{E: e, M: m}}
}
func (Builder) Call(fn Expr, args []ast.Arg) Expr {
	return &Call{ast.Call{Fn: fn, Args: args}}
}

func (Builder) In(e Expr, exprs []Expr) Expr {
	vals := make([]runtime.Value, len(exprs))
	for i, e := range exprs {
		c, ok := e.(*ast.Constant)
		if !ok {
			panic("query in values must be constants")
		}
		vals[i] = c.Val
	}
	return &In{ast.InVals{E: e, Vals: vals}}
}

type Constant struct {
	ast.Constant
}

type Ident struct {
	ast.Ident
}

type Unary struct {
	ast.Unary
}

type Binary struct {
	ast.Binary
}

type Trinary struct {
	ast.Trinary
}

type Nary struct {
	ast.Nary
}

type Mem struct {
	ast.Mem
}

type In struct {
	ast.InVals
}

type Call struct {
	ast.Call
}
