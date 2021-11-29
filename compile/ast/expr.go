// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"math"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/strs"
)

// Note: Stored fields ignore rules.
// This seems "wrong" but it matches what jSuneido does.
// Normally this isn't an issue because the stored field has a value.

type Context struct {
	Th   *Thread
	Tran *SuTran
	Hdr  *Header
	Row  Row
}

func (c *Constant) Eval(*Context) Value {
	return c.Val
}

func (c *Constant) Columns() []string {
	return []string{}
}

func (id *Ident) Eval(c *Context) Value {
	return c.Row.GetVal(c.Hdr, id.Name, c.Th, c.Tran)
}

func (id *Ident) Columns() []string {
	if str.Capitalized(id.Name) {
		return nil
	}
	return []string{id.Name}
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

func (u *Unary) Columns() []string {
	return u.E.Columns()
}

// Binary -----------------------------------------------------------

// CanEvalRaw returns true if Eval doesn't need to unpack the values.
// It sets b.evalRaw which is later used by Eval.
func (b *Binary) CanEvalRaw(cols []string) bool {
	if b.canEvalRaw2(cols) {
		b.evalRaw = true
		c := b.Rhs.(*Constant)
		c.Packed = Pack(c.Val.(Packable))
		return true
	}
	return false
}

func (b *Binary) canEvalRaw2(cols []string) bool {
	if !b.rawOp() {
		return false
	}
	if IsColumn(b.Lhs, cols) && isConstant(b.Rhs) {
		return true
	}
	if isConstant(b.Lhs) && IsColumn(b.Rhs, cols) {
		b.Lhs, b.Rhs = b.Rhs, b.Lhs // swap
		b.Tok = reverseBinary[b.Tok]
		return true
	}
	return false
}

func (b *Binary) rawOp() bool {
	switch b.Tok {
	case tok.Is, tok.Isnt, tok.Lt, tok.Lte, tok.Gt, tok.Gte:
		return true
	}
	return false
}

func IsColumn(e Expr, cols []string) bool {
	if id, ok := e.(*Ident); ok && strs.Contains(cols, id.Name) {
		return true
	}
	return false
}

func isConstant(e Expr) bool {
	_, ok := e.(*Constant)
	return ok
}

var reverseBinary = map[tok.Token]tok.Token{
	tok.Is:   tok.Isnt,
	tok.Isnt: tok.Is,
	tok.Lt:   tok.Gt,
	tok.Lte:  tok.Gte,
	tok.Gt:   tok.Lt,
	tok.Gte:  tok.Lte,
}

func (b *Binary) Eval(c *Context) Value {
	// NOTE: only Eval raw if b.evalRaw was set by CanEvalRaw
	if b.evalRaw {
		name := b.Lhs.(*Ident).Name
		lhs := c.Row.GetRaw(c.Hdr, name)
		rhs := b.Rhs.(*Constant).Packed
		switch b.Tok {
		case tok.Is:
			return SuBool(lhs == rhs)
		case tok.Isnt:
			return SuBool(lhs != rhs)
		case tok.Lt:
			return SuBool(lhs < rhs)
		case tok.Lte:
			return SuBool(lhs <= rhs)
		case tok.Gt:
			return SuBool(lhs > rhs)
		case tok.Gte:
			return SuBool(lhs >= rhs)
		}
	}
	return b.eval(b.Lhs.Eval(c), b.Rhs.Eval(c))
}

func (b *Binary) eval(lhs, rhs Value) Value {
	switch b.Tok {
	case tok.Is:
		return OpIs(lhs, rhs)
	case tok.Isnt:
		return OpIsnt(lhs, rhs)
	case tok.Match:
		return OpMatch(nil, lhs, rhs)
	case tok.MatchNot:
		return OpMatch(nil, lhs, rhs).Not()
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

func (b *Binary) Columns() []string {
	return sset.Union(b.Lhs.Columns(), b.Rhs.Columns())
}

func (tri *Trinary) Eval(c *Context) Value {
	cond := tri.Cond.Eval(c)
	if cond == True {
		return tri.T.Eval(c)
	}
	return tri.F.Eval(c)
}

func (tri *Trinary) Columns() []string {
	return sset.Union(tri.Cond.Columns(),
		sset.Union(tri.T.Columns(), tri.F.Columns()))
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

func (a *Nary) Columns() []string {
	cols := a.Exprs[0].Columns()
	for _, e := range a.Exprs[1:] {
		cols = sset.Union(cols, e.Columns())
	}
	return cols
}

// Range ------------------------------------------------------------

func (a *RangeTo) Eval(c *Context) Value {
	e := a.E.Eval(c)
	from := evalOr(a.From, c, Zero)
	to := evalOr(a.To, c, MaxInt)
	return e.RangeTo(ToIndex(from), ToInt(to))
}

func (a *RangeTo) Columns() []string {
	return sset.Union(a.E.Columns(),
		sset.Union(a.From.Columns(), a.To.Columns()))
}

func (a *RangeLen) Eval(c *Context) Value {
	e := a.E.Eval(c)
	from := evalOr(a.From, c, Zero)
	n := evalOr(a.Len, c, MaxInt)
	return e.RangeLen(ToIndex(from), ToInt(n))
}

func (a *RangeLen) Columns() []string {
	return sset.Union(a.E.Columns(),
		sset.Union(a.From.Columns(), a.Len.Columns()))
}

func evalOr(e Expr, c *Context, v Value) Value {
	if e == nil {
		return v
	}
	return e.Eval(c)
}

// In ---------------------------------------------------------------

// CanEvalRaw returns true if Eval doesn't need to unpack the values.
// It sets Packed which is later used by Eval.
func (a *In) CanEvalRaw(cols []string) bool {
	if !IsColumn(a.E, cols) {
		return false
	}
	packed := make([]string, 0, len(a.Exprs))
	for _, e := range a.Exprs {
		c, ok := e.(*Constant)
		if !ok {
			return false
		}
		packed = sset.AddUnique(packed, Pack(c.Val.(Packable)))
	}
	a.Packed = packed
	return true
}

func (a *In) Eval(c *Context) Value {
	if a.Packed != nil {
		id := a.E.(*Ident)
		e := c.Row.GetRaw(c.Hdr, id.Name)
		for _, p := range a.Packed {
			if e == p {
				return True
			}
		}
		return False
	}
	x := a.E.Eval(c)
	for _, e := range a.Exprs {
		y := e.Eval(c)
		if x.Equal(y) {
			return True
		}
	}
	return False
}

func (a *In) Columns() []string {
	cols := a.E.Columns()
	for _, e := range a.Exprs {
		cols = sset.Union(cols, e.Columns())
	}
	return cols
}

// ------------------------------------------------------------------

func (a *Mem) Eval(c *Context) Value {
	e := a.E.Eval(c)
	m := a.M.Eval(c)
	result := e.Get(nil, m)
	if result == nil {
		panic("uninitialized member: " + m.String())
	}
	return result
}

func (a *Mem) Columns() []string {
	return sset.Union(a.E.Columns(), a.M.Columns())
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
		fn = Global.GetName(c.Th, f.Name)
	case *Mem:
		this = f.E.Eval(c)
		meth := f.M.Eval(c)
		fn = c.Th.Lookup(this, ToStr(meth))
	default:
		fn = a.Fn.Eval(c)
	}
	return c.Th.PushCall(fn, this, as, args...)
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

func (a *Call) Columns() []string {
	cols := a.Fn.Columns()
	for _, e := range a.Args {
		cols = sset.Union(cols, e.E.Columns())
	}
	return cols
}

func (a *Block) Eval(*Context) Value {
	panic("queries do not support blocks")
}

func (a *Block) Columns() []string {
	panic("queries do not support blocks")
}

func (a *Function) Eval(*Context) Value {
	panic("queries do not support functions")
}

func (a *Function) Columns() []string {
	panic("queries do not support functions")
}

func CantBeEmpty(e Expr, cols []string) bool {
	if !e.CanEvalRaw(cols) {
		return false
	}
	switch e := e.(type) {
	case *Binary:
		c := e.Rhs.(*Constant)
		switch e.Tok {
		case tok.Is:
			return c.Val != EmptyStr
		case tok.Isnt:
			return c.Val == EmptyStr
		case tok.Lt:
			return c.Val.Compare(EmptyStr) <= 0
		case tok.Lte:
			return c.Val.Compare(EmptyStr) < 0
		case tok.Gt:
			return c.Val.Compare(EmptyStr) >= 0
		case tok.Gte:
			return c.Val.Compare(EmptyStr) > 0
		default:
			return false
		}
	case *In:
		for _, p := range e.Packed {
			if p == "" {
				return false
			}
		}
		return true
	default:
		return false
	}
}
