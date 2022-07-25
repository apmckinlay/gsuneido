// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

// Folder implements constant folding for expressions.
// Doing the folding as the AST is built is implicitly bottom up
// without requiring an explicit tree traversal.
// It also means we only build the folded tree.
// Some methods are split so they can be used by propfold.
type Folder struct{}

var _ Builder = (*Folder)(nil)

func (f Folder) constant(val Value) Expr {
	return &Constant{Val: val}
}

func (f Folder) Symbol(s SuStr) Expr {
	return &Constant{Val: s}
}

func (f Folder) Unary(token tok.Token, expr Expr) Expr {
	return f.foldUnary(&Unary{Tok: token, E: expr})
}

func (f Folder) foldUnary(u *Unary) Expr {
	if c, ok := u.E.(*Constant); ok && u.Tok != tok.Div {
		return f.constant(u.eval(c.Val))
	}
	return u
}

func (f Folder) Binary(lhs Expr, token tok.Token, rhs Expr) Expr {
	return f.foldBinary(&Binary{Lhs: lhs, Tok: token, Rhs: rhs})
}

func (f Folder) foldBinary(b *Binary) Expr {
	rhs, ok := b.Rhs.(*Constant)
	if !ok {
		return b
	}
	lhs, ok := b.Lhs.(*Constant)
	if !ok {
		return b
	}
	return f.constant(b.eval(lhs.Val, rhs.Val))
}

func (f Folder) Trinary(cond Expr, e1 Expr, e2 Expr) Expr {
	return f.foldTrinary(&Trinary{Cond: cond, T: e1, F: e2})
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
	return f.foldIn(&In{E: e, Exprs: exprs})
}

func (f Folder) foldIn(in *In) Expr {
	if len(in.Exprs) == 0 {
		return f.constant(False)
	}
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
			return f.constant(True)
		}
	}
	return f.constant(False)
}

var allones Value = SuDnum{Dnum: dnum.FromInt(0xffffffff)}

func (f Folder) Nary(token tok.Token, exprs []Expr) Expr {
	return f.foldNary(&Nary{Tok: token, Exprs: exprs})
}

func (f Folder) foldNary(n *Nary) Expr {
	exprs := n.Exprs
	if len(exprs) == 1 {
		return exprs[0]
	}
	switch n.Tok {
	case tok.Add: // includes Sub
		exprs = commutative(exprs, OpAdd, nil, Zero)
	case tok.Mul: // includes Div
		exprs = f.foldMul(exprs)
	case tok.BitOr:
		exprs = commutative(exprs, OpBitOr, allones, Zero)
	case tok.BitAnd:
		exprs = commutative(exprs, OpBitAnd, Zero, allones)
	case tok.BitXor:
		exprs = commutative(exprs, OpBitXor, nil, Zero)
	case tok.Or:
		exprs = commutative(exprs, or, True, False)
		exprs = foldOrToIn(exprs)
	case tok.And:
		exprs = commutative(exprs, and, False, True)
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
// zero is a short circuit value e.g. false for and
func commutative(exprs []Expr, bop bopfn, zero, identity Value) []Expr {
	first := -1
	dst := 0
	for _, e := range exprs {
		if c, ok := e.(*Constant); !ok {
			exprs[dst] = e
			dst++
		} else {
			if c.Val.Equal(zero) {
				exprs[0] = c
				return exprs[:1]
			}
			if c.Val.Equal(identity) {
				continue
			}
			if first == -1 {
				first = dst
				exprs[dst] = e
				dst++
			} else {
				copy := *exprs[first].(*Constant)
				copy.Val = bop(copy.Val, c.Val)
				exprs[first] = &copy
				}
		}
	}
	if dst == 1 && first != -1 {
		// compile-time type check
		bop(identity, exprs[first].(*Constant).Val)
	}
	if dst <= 1 && first == -1 {
		// keep operation as run-time type check
		exprs[dst] = &Constant{Val: identity}
		dst++
	}
	return exprs[:dst]
}

func foldOrToIn(exprs []Expr) []Expr {
	var idPrev *Ident
	var in []Expr
	newExprs := make([]Expr, 0, len(exprs))
	for _, expr := range exprs {
		id, e := idIs(expr)
		if id != nil && idPrev != nil && idPrev.Name == id.Name {
			// accumulate
			in = append(in, e)
			if len(in) > 1 {
				continue
			}
		} else {
			if len(in) > 1 {
				// flush
				newExprs[len(newExprs)-1] = &In{E: idPrev, Exprs: in}
			}
			if id == nil {
				idPrev, in = nil, nil
			} else {
				idPrev = id
				in = []Expr{e}
			}
		}
		newExprs = append(newExprs, expr)
	}
	if len(in) > 1 {
		// flush
		newExprs[len(newExprs)-1] = &In{E: idPrev, Exprs: in}
	}
	return newExprs
}
func idIs(e Expr) (*Ident, Expr) {
	if bin, ok := e.(*Binary); ok && bin.Tok == tok.Is {
		if id, ok := bin.Lhs.(*Ident); ok {
			return id, bin.Rhs
		}
	}
	return nil, nil
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
		if !mul.Equal(One) || len(exprs) == 0 {
			exprs = append(exprs, f.constant(mul))
		}
	} else {
		exprs = append(exprs, f.Unary(tok.Div, f.constant(div)))
	}
	if len(exprs) == 1 && !unaryDivOrConstant(exprs[0]) {
		// force an operation to preserve conversion
		exprs = append(exprs, f.constant(One))
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
	first := -1
	dst := 0
	for _, e := range exprs {
		if c, ok := e.(*Constant); !ok {
			exprs[dst] = e
			dst++
			first = -1
		} else if first == -1 {
			first = dst
			exprs[dst] = e
			dst++
		} else {
			copy := *exprs[first].(*Constant)
			copy.Val = SuStr(AsStr(copy.Val) + AsStr(c.Val))
			exprs[first] = &copy
		}
	}
	return exprs[:dst]
}
