package ast

import (
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
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
	c, ok := expr.(*Constant)
	if !ok {
		return f.Factory.Unary(token, expr)
	}
	val := c.Val
	switch token {
	case tok.Add:
		val = UnaryPlus(val)
	case tok.Sub:
		val = UnaryMinus(val)
	case tok.Div:
		val = Div(One, val)
	case tok.Not:
		val = Not(val)
	case tok.BitNot:
		val = BitNot(val)
	case tok.LParen:
		break
	default:
		panic("folder unexpected unary operator " + token.String())
	}
	return f.Constant(val)
}

func (f Folder) Binary(lhs Expr, token tok.Token, rhs Expr) Expr {
	c1, ok := lhs.(*Constant)
	if !ok {
		return f.Factory.Binary(lhs, token, rhs)
	}
	c2, ok := rhs.(*Constant)
	if !ok {
		return f.Factory.Binary(lhs, token, rhs)
	}
	val := c1.Val
	val2 := c2.Val
	switch token {
	case tok.Is:
		val = Is(val, val2)
	case tok.Isnt:
		val = Isnt(val, val2)
	case tok.Match:
		pat := regex.Compile(IfStr(val2))
		val = Match(val, pat)
	case tok.MatchNot:
		pat := regex.Compile(IfStr(val2))
		val = Not(Match(val, pat))
	case tok.Lt:
		val = Lt(val, val2)
	case tok.Lte:
		val = Lte(val, val2)
	case tok.Gt:
		val = Gt(val, val2)
	case tok.Gte:
		val = Gte(val, val2)
	case tok.Mod:
		val = Mod(val, val2)
	case tok.LShift:
		val = LeftShift(val, val2)
	case tok.RShift:
		val = RightShift(val, val2)
	default:
		panic("folder unexpected unary operator " + token.String())
	}
	return f.Constant(val)
}

func (f Folder) Trinary(cond Expr, e1 Expr, e2 Expr) Expr {
	c, ok := cond.(*Constant)
	if !ok {
		return f.Factory.Trinary(cond, e1, e2)
	}
	if c.Val == True {
		return e1
	}
	if c.Val == False {
		return e2
	}
	panic("?: requires boolean")
}

func (f Folder) In(e Expr, exprs []Expr) Expr {
	c, ok := e.(*Constant)
	if !ok {
		return f.Factory.In(e, exprs)
	}
	for _, e := range exprs {
		c2, ok := e.(*Constant)
		if !ok {
			return f.Factory.In(e, exprs)
		}
		if c.Val.Equal(c2.Val) {
			return f.Constant(True)
		}
	}
	return f.Constant(False)
}

var allones Value = SuDnum{Dnum: dnum.FromInt(0xffffffff)}

func (f Folder) Nary(token tok.Token, exprs []Expr) Expr {
	switch token {
	case tok.Add: // including Sub
		exprs = commutative(exprs, Add, Zero, nil)
	case tok.Mul: // including Div
		exprs = commutative(exprs, Mul, One, Zero)
	case tok.BitOr:
		exprs = commutative(exprs, BitOr, Zero, allones)
	case tok.BitAnd:
		exprs = commutative(exprs, BitAnd, allones, Zero)
	case tok.BitXor:
		exprs = commutative(exprs, BitXor, Zero, nil)
	case tok.Or:
		exprs = commutative(exprs, or, False, True)
	case tok.And:
		exprs = commutative(exprs, and, True, False)
	case tok.Cat:
		exprs = foldCat(exprs)
	default:
		// 	panic("folder unexpected n-ary operator " + token.String())
	}
	if len(exprs) == 1 {
		return exprs[0]
	}
	return f.Factory.Nary(token, exprs)
}

func or(x, y Value) Value {
	return SuBool(Bool(x) || Bool(y))
}

func and(x, y Value) Value {
	return SuBool(Bool(x) && Bool(y))
}

type bopfn func(Value, Value) Value

// commutative folds constants in a list of expressions
// skip is the identity value e.g. zero for add
// fold is a short circuit value e.g. zero for multiply
func commutative(exprs []Expr, bop bopfn, skip, fold Value) []Expr {
	var first *Constant
	dst := 0
	for _, e := range exprs {
		if c, ok := e.(*Constant); !ok {
			exprs[dst] = e
			dst++
		} else if !skip.Equal(c.Val) {
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
			first.Val = SuStr(ToStr(first.Val) + ToStr(c.Val))
		}
	}
	return exprs[:dst]
}
