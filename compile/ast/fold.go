// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"

	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/regex"
)

// Fold traverses an AST and does constant propagation and folding
// modifying the AST
func PropFold(fn *Function) *Function {
	// Final variables (set once, not modified) are determined during parse
	propfold(fn, fn.Final)
	return fn
}

// propfold - constant propagation and folding
func propfold(fn *Function, vars map[string]int) {
	f := fold{vars: map[string]Value{}, srcpos: -1}
	for v := range vars {
		f.vars[v] = nil
	}
	defer func(f *fold) {
		if e := recover(); e != nil {
			panic(fmt.Sprintf("compile error @%d %s", f.srcpos, e))
		}
	}(&f)
	Traverse(fn, &f)
}

type fold struct {
	vars   map[string]Value
	srcpos int
}

func (f *fold) Before(node Node) bool {
	if stmt, ok := node.(Statement); ok {
		f.srcpos = stmt.Position() // for error reporting
	}
	return true
}

func (f *fold) After(node Node) Node {
	node = f.fold(node)
	node = f.findConst(node)
	return node
}

func (f *fold) fold(node Node) Node {
	switch node := node.(type) {
	case *Ident:
		return f.ident(node)
	case *Unary:
		return f.unary(node)
	case *Binary:
		return f.binary(node)
	case *Trinary:
		return f.trinary(node)
	case *In:
		return f.in(node)
	case *Nary:
		return f.nary(node)
	case *If:
		return f.ifStmt(node)
	}
	// TODO switch
	return node
}

func (f *fold) ident(node *Ident) Node {
	if val := f.vars[node.Name]; val != nil {
		return constant(val)
	}
	return node
}

func (f *fold) unary(unary *Unary) Node {
	if c, ok := unary.E.(*Constant); ok {
		val := c.Val
		switch unary.Tok {
		case tok.Add:
			val = UnaryPlus(val)
		case tok.Sub:
			val = UnaryMinus(val)
		case tok.Not:
			val = Not(val)
		case tok.BitNot:
			val = BitNot(val)
		case tok.LParen:

		default:
			return unary
		}
		return constant(val)
	}
	return unary
}

func (f *fold) binary(b *Binary) Node {
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
		val = Is(val, val2)
	case tok.Isnt:
		val = Isnt(val, val2)
	case tok.Match:
		pat := regex.Compile(ToStr(val2))
		val = Match(val, pat)
	case tok.MatchNot:
		pat := regex.Compile(ToStr(val2))
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
		return b
	}
	return constant(val)
}

func (f *fold) trinary(t *Trinary) Node {
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

func (f *fold) in(in *In) Node {
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
			return constTrue
		}
	}
	return constFalse
}

var allones Value = SuDnum{Dnum: dnum.FromInt(0xffffffff)}

func (f *fold) nary(n *Nary) Node {
	exprs := n.Exprs
	switch n.Tok {
	case tok.Add: // includes Sub
		exprs = commutative(exprs, Add, nil)
	case tok.Mul: // includes Div
		exprs = f.foldMul(exprs)
	case tok.BitOr:
		exprs = commutative(exprs, BitOr, allones)
	case tok.BitAnd:
		exprs = commutative(exprs, BitAnd, Zero)
	case tok.BitXor:
		exprs = commutative(exprs, BitXor, nil)
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
	return SuBool(Bool(x) || Bool(y))
}

func and(x, y Value) Value {
	return SuBool(Bool(x) && Bool(y))
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

func (f *fold) foldMul(exprs []Expr) []Expr {
	// extract and combine constants
	mul := One
	div := One
	dst := 0
	for _, e := range exprs {
		if ud := unaryDivConst(e); ud != nil {
			div = Mul(div, ud)
		} else if c, ok := e.(*Constant); ok {
			mul = Mul(mul, c.Val)
		} else {
			exprs[dst] = e
			dst++
		}
	}
	exprs = exprs[:dst]

	if !div.Equal(One) && (!mul.Equal(One) || len(exprs) == 0) {
		mul = Div(mul, div)
		div = One
	}
	if div.Equal(One) {
		if !mul.Equal(One) {
			exprs = append(exprs, constant(mul))
		}
	} else {
		exprs = append(exprs, &Unary{Tok: tok.Div, E: constant(div)})
	}
	if len(exprs) == 1 && !unaryDivOrConstant(exprs[0]) {
		// force an operation to preserve conversion
		exprs = append(exprs, constant(One))
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

var constTrue = &Constant{Val: True}
var constFalse = &Constant{Val: False}

func constant(val Value) *Constant {
	return &Constant{Val: val}
}

func (f *fold) findConst(node Node) Node {
	if node, ok := node.(*Binary); ok {
		if node.Tok == tok.Eq {
			if id, ok := node.Lhs.(*Ident); ok {
				if _, ok := f.vars[id.Name]; ok {
					if val, ok := node.Rhs.(*Constant); ok {
						f.vars[id.Name] = val.Val
						return node.Rhs // remove assignment
					}
				}
			}
		}
	}
	return node
}
