package ast

import (
	. "github.com/apmckinlay/gsuneido/lexer"
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

func (f Folder) Unary(tok Token, expr Expr) Expr {
	c, ok := expr.(*Constant)
	if !ok {
		return f.Factory.Unary(tok, expr)
	}
	val := c.Val
	switch tok {
	case ADD:
		val = Uplus(val)
	case SUB:
		val = Uminus(val)
	case DIV:
		val = Div(One, val)
	case NOT:
		val = Not(val)
	case BITNOT:
		val = BitNot(val)
	default:
		panic("folder unexpected unary operator " + tok.String())
	}
	return f.Constant(val)
}

func (f Folder) Binary(lhs Expr, tok Token, rhs Expr) Expr {
	c1, ok := lhs.(*Constant)
	if !ok {
		return f.Factory.Binary(lhs, tok, rhs)
	}
	c2, ok := rhs.(*Constant)
	if !ok {
		return f.Factory.Binary(lhs, tok, rhs)
	}
	val := c1.Val
	val2 := c2.Val
	switch tok {
	case IS:
		val = Is(val, val2)
	case ISNT:
		val = Isnt(val, val2)
	case MATCH:
		pat := regex.Compile(val2.ToStr())
		val = Match(val, pat)
	case MATCHNOT:
		pat := regex.Compile(val2.ToStr())
		val = Not(Match(val, pat))
	case LT:
		val = Lt(val, val2)
	case LTE:
		val = Lte(val, val2)
	case GT:
		val = Gt(val, val2)
	case GTE:
		val = Gte(val, val2)
	case MOD:
		val = Mod(val, val2)
	case LSHIFT:
		val = Lshift(val, val2)
	case RSHIFT:
		val = Rshift(val, val2)
	default:
		panic("folder unexpected unary operator " + tok.String())
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

func (f Folder) Nary(tok Token, exprs []Expr) Expr {
	switch tok {
	case ADD: // including SUB
		exprs = commutative(exprs, Add, Zero, nil)
	case MUL: // including DIV
		exprs = commutative(exprs, Mul, One, Zero)
	case BITOR:
		exprs = commutative(exprs, Bitor, Zero, allones)
	case BITAND:
		exprs = commutative(exprs, Bitand, allones, Zero)
	case BITXOR:
		exprs = commutative(exprs, Bitxor, Zero, nil)
	case OR:
		exprs = commutative(exprs, or, False, True)
	case AND:
		exprs = commutative(exprs, and, True, False)
	case CAT:
		exprs = foldCat(exprs)
	default:
		// 	panic("folder unexpected n-ary operator " + tok.String())
	}
	if len(exprs) == 1 {
		return exprs[0]
	}
	return f.Factory.Nary(tok, exprs)
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
			first.Val = SuStr(Cat(first.Val, c.Val).ToStr()) // flatten Concat
		}
	}
	return exprs[:dst]
}
