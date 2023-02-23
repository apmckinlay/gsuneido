// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
)

type Token = tokens.Token

type Builder interface {
	Symbol(s SuStr) Expr
	Unary(tok Token, expr Expr) Expr
	Binary(lhs Expr, tok Token, rhs Expr) Expr
	Trinary(cond Expr, t Expr, f Expr) Expr
	Nary(tok Token, exprs []Expr) Expr
	In(e Expr, exprs []Expr) Expr
	Call(fn Expr, args []Arg, end int32) Expr
}

// Factory is a simple pass through Builder
type Factory struct{}

var _ Builder = (*Factory)(nil)

func (Factory) Symbol(s SuStr) Expr {
	return &Constant{Val: s}
}
func (Factory) Unary(tok Token, expr Expr) Expr {
	if c, ok := expr.(*Constant); ok && tok == tokens.Sub {
		return &Constant{Val: OpUnaryMinus(c.Val)}
	}
	return &Unary{Tok: tok, E: expr}
}
func (Factory) Binary(lhs Expr, tok Token, rhs Expr) Expr {
	return &Binary{Lhs: lhs, Tok: tok, Rhs: rhs}
}
func (Factory) Trinary(cond Expr, t Expr, f Expr) Expr {
	return &Trinary{Cond: cond, T: t, F: f}
}
func (Factory) Nary(tok Token, exprs []Expr) Expr {
	return &Nary{Tok: tok, Exprs: exprs}
}
func (Factory) In(e Expr, exprs []Expr) Expr {
	return &In{E: e, Exprs: exprs}
}
func (Factory) Call(fn Expr, args []Arg, end int32) Expr {
	return &Call{Fn: fn, Args: args, End: end}
}
