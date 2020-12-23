// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/regex"
)

// Folder implements constant folding for expressions.
// It is a "decorator" Factory that wraps another Factory (e.g. Builder)
// Doing the folding as the AST is built is implicitly bottom up
// without requiring an explicit tree traversal.
// It also means we only build the folded tree.
type Folder struct {
	Factory
}

func (f Folder) Unary(token tok.Token, expr Expr) Expr {
	return f.foldUnary(f.Factory.Unary(token, expr).(*Unary))
}

func (f Folder) foldUnary(u *Unary) Expr {
	if c, ok := u.E.(*Constant); ok && u.Tok != tok.Div {
		val := c.Val
		switch u.Tok {
		case tok.Add:
			val = OpUnaryPlus(val)
		case tok.Sub:
			val = OpUnaryMinus(val)
		case tok.Not:
			val = OpNot(val)
		case tok.BitNot:
			val = OpBitNot(val)
		case tok.LParen:
			break
		default:
			panic("folder unexpected unary operator " + u.Tok.String())
		}
		return f.Constant(val)
	}
	return u
}

func (f Folder) Binary(lhs Expr, token tok.Token, rhs Expr) Expr {
	return f.foldBinary(f.Factory.Binary(lhs, token, rhs).(*Binary))
}

func (f Folder) foldBinary(b *Binary) Expr {
	c1, ok := b.Lhs.(*Constant)
	if !ok {
		return b
	}
	c2, ok := b.Rhs.(*Constant)
	if !ok {
		return b
	}
	val := c1.Val
	val2 := c2.Val
	switch b.Tok {
	case tok.Is:
		val = OpIs(val, val2)
	case tok.Isnt:
		val = OpIsnt(val, val2)
	case tok.Match:
		pat := regex.Compile(ToStr(val2))
		val = SuBool(pat.Matches(ToStr(val)))
	case tok.MatchNot:
		pat := regex.Compile(ToStr(val2))
		val = SuBool(!pat.Matches(ToStr(val)))
	case tok.Lt:
		val = OpLt(val, val2)
	case tok.Lte:
		val = OpLte(val, val2)
	case tok.Gt:
		val = OpGt(val, val2)
	case tok.Gte:
		val = OpGte(val, val2)
	case tok.Mod:
		val = OpMod(val, val2)
	case tok.LShift:
		val = OpLeftShift(val, val2)
	case tok.RShift:
		val = OpRightShift(val, val2)
	default:
		panic("folder unexpected binary operator " + b.Tok.String())
	}
	return f.Constant(val)
}

func (f Folder) Trinary(cond Expr, e1 Expr, e2 Expr) Expr {
	return f.foldTrinary(f.Factory.Trinary(cond, e1, e2).(*Trinary))
}

func (f Folder) foldTrinary(t *Trinary) Expr {
	c, ok := t.Cond.(*Constant)
	if !ok {
		return t
	}
	if c.Val == True {
		return t.T
	}
	if c.Val == False {
		return t.F
	}
	panic("?: requires boolean")
}

func (f Folder) In(e Expr, exprs []Expr) Expr {
	return f.foldIn(f.Factory.In(e, exprs).(*In))
}

func (f Folder) foldIn(in *In) Expr {
	c, ok := in.E.(*Constant)
	if !ok {
		return in
	}
	for _, e := range in.Exprs {
		c2, ok := e.(*Constant)
		if !ok {
			return in
		}
		if c.Val.Equal(c2.Val) {
			return f.Constant(True)
		}
	}
	return f.Constant(False)
}

var allones Value = SuDnum{Dnum: dnum.FromInt(0xffffffff)}

func (f Folder) Nary(token tok.Token, exprs []Expr) Expr {
	return f.foldNary(f.Factory.Nary(token, exprs).(*Nary))
}

func (f Folder) foldNary(n *Nary) Expr {
	exprs := n.Exprs
	switch n.Tok {
	case tok.Add: // includes Sub
		exprs = commutative(exprs, OpAdd, nil)
	case tok.Mul: // includes Div
		exprs = f.foldMul(exprs)
	case tok.BitOr:
		exprs = commutative(exprs, OpBitOr, allones)
	case tok.BitAnd:
		exprs = commutative(exprs, OpBitAnd, Zero)
	case tok.BitXor:
		exprs = commutative(exprs, OpBitXor, nil)
	case tok.Or:
		exprs = commutative(exprs, or, True)
	case tok.And:
		exprs = commutative(exprs, and, False)
	case tok.Cat:
		exprs = foldCat(exprs)
	default:
		panic("fold unexpected n-ary operator " + n.Tok.String())
	}
	if len(exprs) == 1 {
		return exprs[0]
	}
	n.Exprs = exprs
	return n
}

func or(x, y Value) Value {
	return SuBool(OpBool(x) || OpBool(y))
}

func and(x, y Value) Value {
	return SuBool(OpBool(x) && OpBool(y))
}

type bopfn func(Value, Value) Value

// commutative folds constants in a list of expressions
// fold is a short circuit value e.g. zero for multiply
func commutative(exprs []Expr, bop bopfn, fold Value) []Expr {
	var first *Constant
	dst := 0
	for _, e := range exprs {
		if c, ok := e.(*Constant); !ok {
			exprs[dst] = e
			dst++
		} else {
			if c.Val.Equal(fold) {
				exprs[0] = c
				return exprs[:1]
			}
			if first == nil {
				first = c
				exprs[dst] = e
				dst++
			} else {
				first.Val = bop(first.Val, c.Val)
			}
		}
	}
	return exprs[:dst]
}

func (f Folder) foldMul(exprs []Expr) []Expr {
	// extract and combine constants
	mul := One
	div := One
	dst := 0
	for _, e := range exprs {
		if ud := unaryDivConst(e); ud != nil {
			div = OpMul(div, ud)
		} else if c, ok := e.(*Constant); ok {
			mul = OpMul(mul, c.Val)
		} else {
			exprs[dst] = e
			dst++
		}
	}
	exprs = exprs[:dst]

	if !div.Equal(One) && (!mul.Equal(One) || len(exprs) == 0) {
		mul = OpDiv(mul, div)
		div = One
	}
	if div.Equal(One) {
		if !mul.Equal(One) {
			exprs = append(exprs, f.Constant(mul))
		}
	} else {
		exprs = append(exprs, f.Unary(tok.Div, f.Constant(div)))
	}
	if len(exprs) == 1 && !unaryDivOrConstant(exprs[0]) {
		// force an operation to preserve conversion
		exprs = append(exprs, f.Constant(One))
	}
	return exprs
}

func unaryDivConst(e Expr) Value {
	if u, ok := e.(*Unary); ok {
		if c, ok := u.E.(*Constant); ok {
			return c.Val
		}
	}
	return nil
}

func unaryDivOrConstant(e Expr) bool {
	u, ok := e.(*Unary)
	if ok && u.Tok == tok.Div {
		return true
	}
	_, ok = e.(*Constant)
	return ok
}

// foldCat folds contiguous constants in a list of expressions
// cat is not commutative, so only combine contiguous constants
func foldCat(exprs []Expr) []Expr {
	var first *Constant
	dst := 0
	for _, e := range exprs {
		if c, ok := e.(*Constant); !ok {
			exprs[dst] = e
			dst++
			first = nil
		} else if first == nil {
			first = c
			exprs[dst] = e
			dst++
		} else {
			first.Val = SuStr(AsStr(first.Val) + AsStr(c.Val))
		}
	}
	return exprs[:dst]
}
