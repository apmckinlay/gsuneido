// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"math"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/regex"
)

type Context struct {
	T   *Thread
	Hdr *Header
	Row Row
}

func (c *Constant) Eval(*Context) Value {
	return c.Val
}

func (id *Ident) Eval(c *Context) Value {
	return c.Row.Get(c.Hdr, id.Name)
}

func (u *Unary) Eval(c *Context) Value {
	return u.eval(u.E.Eval(c))
}

func (u *Unary) eval(val Value) Value {
	switch u.Tok {
	case tok.Add:
		return OpUnaryPlus(val)
	case tok.Sub:
		return OpUnaryMinus(val)
	case tok.Not:
		return OpNot(val)
	case tok.BitNot:
		return OpBitNot(val)
	case tok.LParen:
		return val
	default:
		panic("unexpected unary operator " + u.Tok.String())
	}
}

func (b *Binary) Eval(c *Context) Value {
	return b.eval(b.Lhs.Eval(c), b.Rhs.Eval(c))
}

func (b *Binary) eval(lhs, rhs Value) Value {
	switch b.Tok {
	case tok.Is:
		return OpIs(lhs, rhs)
	case tok.Isnt:
		return OpIsnt(lhs, rhs)
	case tok.Match:
		pat := regex.Compile(ToStr(rhs)) //TODO cache
		return SuBool(pat.Matches(ToStr(lhs)))
	case tok.MatchNot:
		pat := regex.Compile(ToStr(rhs)) //TODO cache
		return SuBool(!pat.Matches(ToStr(lhs)))
	case tok.Lt:
		return OpLt(lhs, rhs)
	case tok.Lte:
		return OpLte(lhs, rhs)
	case tok.Gt:
		return OpGt(lhs, rhs)
	case tok.Gte:
		return OpGte(lhs, rhs)
	case tok.Mod:
		return OpMod(lhs, rhs)
	case tok.LShift:
		return OpLeftShift(lhs, rhs)
	case tok.RShift:
		return OpRightShift(lhs, rhs)
	default:
		panic("unexpected binary operator " + b.Tok.String())
	}
}

func (tri *Trinary) Eval(c *Context) Value {
	cond := tri.Cond.Eval(c)
	if cond == True {
		return tri.T.Eval(c)
	}
	return tri.F.Eval(c)
}

// Nary -------------------------------------------------------------

func (a *Nary) Eval(c *Context) Value {
	exprs := a.Exprs
	switch a.Tok {
	case tok.Add: // includes Sub
		return nary(exprs, c, OpAdd, nil)
	case tok.Mul: // includes Div
		return muldiv(exprs, c)
	case tok.BitOr:
		return nary(exprs, c, OpBitOr, allones)
	case tok.BitAnd:
		return nary(exprs, c, OpBitAnd, Zero)
	case tok.BitXor:
		return nary(exprs, c, OpBitXor, nil)
	case tok.Or:
		return nary(exprs, c, or, True)
	case tok.And:
		return nary(exprs, c, and, False)
	case tok.Cat:
		return nary(exprs, c, opCat, nil)
	default:
		panic("unexpected n-ary operator " + a.Tok.String())
	}
}

func opCat(x, y Value) Value {
	return OpCat(nil, x, y)
}

func nary(exprs []Expr, c *Context,
	op func(Value, Value) Value, zero Value) Value {
	result := exprs[0].Eval(c)
	for _, e := range exprs[1:] {
		if result.Equal(zero) {
			return zero
		}
		result = op(result, e.Eval(c))
	}
	return result
}

func muldiv(exprs []Expr, c *Context) Value {
	var divs []Expr
	result := exprs[0].Eval(c)
	for _, e := range exprs[1:] {
		if u, ok := e.(*Unary); ok && u.Tok == tok.Div {
			divs = append(divs, u.E)
		} else {
			result = OpMul(result, e.Eval(c))
		}
	}
	if len(divs) > 0 {
		x := divs[0].Eval(c)
		for _, e := range divs[1:] {
			x = OpMul(x, e.Eval(c))
		}
		result = OpDiv(result, x)
	}
	return result
}

// Range ------------------------------------------------------------

func (a *RangeTo) Eval(c *Context) Value {
	e := a.E.Eval(c)
	from := evalOr(a.From, c, Zero)
	to := evalOr(a.To, c, MaxInt)
	return e.RangeTo(ToIndex(from), ToInt(to))
}

func (a *RangeLen) Eval(c *Context) Value {
	e := a.E.Eval(c)
	from := evalOr(a.From, c, Zero)
	n := evalOr(a.Len, c, MaxInt)
	return e.RangeLen(ToIndex(from), ToInt(n))
}

func evalOr(e Expr, c *Context, v Value) Value {
	if e == nil {
		return v
	}
	return e.Eval(c)
}

// ------------------------------------------------------------------

func (a *In) Eval(c *Context) Value {
	x := a.E.Eval(c)
	for _, e := range a.Exprs {
		y := e.Eval(c)
		if x.Equal(y) {
			return True
		}
	}
	return False
}

func (a *Mem) Eval(c *Context) Value {
	e := a.E.Eval(c)
	m := a.M.Eval(c)
	result := e.Get(nil, m)
	if result == nil {
		panic("uninitialized member: " + m.String())
	}
	return result
}

func (a *Call) Eval(c *Context) Value {
	as := argspec(a.Args) //TODO cache
	args := make([]Value, len(a.Args))
	for i, a := range a.Args {
		args[i] = a.E.Eval(c)
	}
	var fn Callable
	var this Value
	switch f := a.Fn.(type) {
	case *Ident:
		fn = Global.GetName(c.T, f.Name)
	case *Mem:
		this = f.E.Eval(c)
		meth := f.M.Eval(c)
		fn = c.T.Lookup(this, ToStr(meth))
	default:
		fn = a.Fn.Eval(c)
	}
	return c.T.PushCall(fn, this, as, args...)
}

func argspec(args []Arg) *ArgSpec {
	if len(args) == 0 {
		return &ArgSpec0
	}
	if len(args) == 1 {
		if args[0].Name == SuStr("@") {
			return &ArgSpecEach0
		} else if args[0].Name == SuStr("@+1") {
			return &ArgSpecEach1
		}
	}
	assert.That(len(args) < math.MaxUint8)
	as := ArgSpec{Nargs: byte(len(args))}
	for _, arg := range args {
		if arg.Name != nil {
			as.Spec = append(as.Spec, byte(len(as.Names)))
			as.Names = append(as.Names, arg.Name)
		}
	}
	if as.Spec == nil {
		for i := 0; i <= AsEach1; i++ {
			if as.Equal(&StdArgSpecs[i]) {
				return &StdArgSpecs[i]
			}
		}
	}
	return &as
}

func (a *Block) Eval(*Context) Value {
	panic("queries do not support blocks")
}

func (a *Function) Eval(*Context) Value {
	panic("queries do not support functions")
}
