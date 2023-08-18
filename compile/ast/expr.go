// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"
	"math"
	"strings"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
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

func (a *Constant) Eval(*Context) Value {
	return a.Val
}

func (a *Constant) Columns() []string {
	return []string{}
}

func (a *Ident) Eval(c *Context) Value {
	return c.Row.GetVal(c.Hdr, a.Name, c.Th, c.Tran)
}

func (a *Ident) Columns() []string {
	if str.Capitalized(a.Name) {
		return nil
	}
	return []string{a.Name}
}

func (a *Unary) Eval(c *Context) Value {
	return a.eval(a.E.Eval(c))
}

func (a *Unary) eval(val Value) Value {
	switch a.Tok {
	case tok.Add:
		return OpUnaryPlus(val)
	case tok.Sub:
		return OpUnaryMinus(val)
	case tok.Not:
		// TODO eval raw e.g. not Number?(x)
		return OpNot(val)
	case tok.BitNot:
		return OpBitNot(val)
	case tok.LParen:
		return val
	}
	panic(assert.ShouldNotReachHere())
}

func (a *Unary) Columns() []string {
	return a.E.Columns()
}

// Binary -----------------------------------------------------------

// CouldEvalRaw is used by replaceExpr to know when to copy
func (a *Binary) CouldEvalRaw() bool {
	// depends on folder putting constant on the right
	return a.rawOp() && isIdent(a.Lhs) && isConstant(a.Rhs)
}

// CanEvalRaw returns true if Eval doesn't need to unpack the values.
// It sets evalRaw and Packed which are used later by Eval.
func (a *Binary) CanEvalRaw(flds []string) bool {
	// depends on folder putting constant on the right
	if a.rawOp() && IsColumn(a.Lhs, flds) && isConstant(a.Rhs) {
		a.evalRaw = true
		c := a.Rhs.(*Constant)
		c.Packed = Pack(c.Val.(Packable))
		return true
	}
	a.evalRaw = false
	return false
}

func (a *Binary) rawOp() bool {
	switch a.Tok {
	case tok.Is, tok.Isnt, tok.Lt, tok.Lte, tok.Gt, tok.Gte:
		return true
	}
	return false
}

// IsColumn handles _lower!
func IsColumn(e Expr, cols []string) bool {
	if id, ok := e.(*Ident); ok && (slices.Contains(cols, id.Name) ||
		(strings.HasSuffix(id.Name, "_lower!") &&
			slices.Contains(cols, strings.TrimSuffix(id.Name, "_lower!")))) {
		return true
	}
	return false
}

func isConstant(e Expr) bool {
	_, ok := e.(*Constant)
	return ok
}

func isIdent(e Expr) bool {
	_, ok := e.(*Ident)
	return ok
}

func (a *Binary) Eval(c *Context) Value {
	if a.evalRaw {
		return a.rawEval(c)
	}
	return a.eval(c.Th, a.Lhs.Eval(c), a.Rhs.Eval(c))
}

func (a *Binary) rawEval(c *Context) Value {
	name := a.Lhs.(*Ident).Name
	lhs := c.Row.GetRaw(c.Hdr, name)
	rhs := a.Rhs.(*Constant).Packed
	switch a.Tok {
	case tok.Is:
		return SuBool(lhs == rhs)
	case tok.Isnt:
		return SuBool(lhs != rhs)
	case tok.Lt:
		return SuBool(packedCmp(lhs, rhs) < 0)
	case tok.Lte:
		return SuBool(packedCmp(lhs, rhs) <= 0)
	case tok.Gt:
		return SuBool(packedCmp(lhs, rhs) > 0)
	case tok.Gte:
		return SuBool(packedCmp(lhs, rhs) >= 0)
	}
	panic(assert.ShouldNotReachHere())
}

func packedCmp(x, y string) int {
	cmp := strings.Compare(x, y)
	if cmp != 0 && options.StrictCompareDb && PackedOrd(x) != PackedOrd(y) {
		panic(fmt.Sprint("StrictCompareDb: ", Unpack(x), " <=> ", Unpack(y)))
	}
	return cmp
}

func (a *Binary) eval(th *Thread, lhs, rhs Value) Value {
	switch a.Tok {
	case tok.Is:
		return OpIs(lhs, rhs)
	case tok.Isnt:
		return OpIsnt(lhs, rhs)
	case tok.Match:
		return OpMatch(th, lhs, rhs)
	case tok.MatchNot:
		return OpMatch(th, lhs, rhs).Not()
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
	}
	panic(assert.ShouldNotReachHere())
}

func (a *Binary) Columns() []string {
	return set.Union(a.Lhs.Columns(), a.Rhs.Columns())
}

func (a *Trinary) Eval(c *Context) Value {
	cond := a.Cond.Eval(c)
	if cond == True {
		return a.T.Eval(c)
	}
	return a.F.Eval(c)
}

func (a *Trinary) Columns() []string {
	return set.Union(a.Cond.Columns(),
		set.Union(a.T.Columns(), a.F.Columns()))
}

// Nary -------------------------------------------------------------

func (a *Nary) CanEvalRaw(flds []string) bool {
	if a.Tok == tok.Or || a.Tok == tok.And {
		for _, e := range a.Exprs {
			if e.CanEvalRaw(flds) {
				return true
			}
		}
	}
	return false
}

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
	}
	panic(assert.ShouldNotReachHere())
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
		cols = set.Union(cols, e.Columns())
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
	cols := a.E.Columns()
	if a.From != nil {
		cols = set.Union(cols, a.From.Columns())
	}
	if a.To != nil {
		cols = set.Union(cols, a.To.Columns())
	}
	return cols
}

func (a *RangeLen) Eval(c *Context) Value {
	e := a.E.Eval(c)
	from := evalOr(a.From, c, Zero)
	n := evalOr(a.Len, c, MaxInt)
	return e.RangeLen(ToIndex(from), ToInt(n))
}

func (a *RangeLen) Columns() []string {
	return set.Union(a.E.Columns(),
		set.Union(a.From.Columns(), a.Len.Columns()))
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
func (a *In) CanEvalRaw(flds []string) bool {
	if !IsColumn(a.E, flds) {
		return false
	}
	packed := make([]string, 0, len(a.Exprs))
	for _, e := range a.Exprs {
		c, ok := e.(*Constant)
		if !ok {
			return false
		}
		packed = set.AddUnique(packed, Pack(c.Val.(Packable)))
	}
	a.Packed = packed
	return true
}

// CouldEvalRaw is used by replaceExpr to know when to copy
func (a *In) CouldEvalRaw() bool {
	if !isIdent(a.E) {
		return false
	}
	for _, e := range a.Exprs {
		if !isConstant(e) {
			return false
		}
	}
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
		cols = set.Union(cols, e.Columns())
	}
	return cols
}

// InRange ---------------------------------------------------------------

// CouldEvalRaw is used by replaceExpr to know when to copy
func (a *InRange) CouldEvalRaw() bool {
	return isIdent(a.E)
}

// CanEvalRaw returns true if Eval doesn't need to unpack the values.
// It sets Packed which is later used by Eval.
func (a *InRange) CanEvalRaw(flds []string) bool {
	// InRange already ensures valid operators and constants
	if !IsColumn(a.E, flds) {
		a.evalRaw = false
		return false
	}
	a.evalRaw = true
	c := a.Org.(*Constant)
	c.Packed = Pack(c.Val.(Packable))
	c = a.End.(*Constant)
	c.Packed = Pack(c.Val.(Packable))
	return true
}

func (a *InRange) Eval(c *Context) Value {
	if a.evalRaw {
		id := a.E.(*Ident)
		x := c.Row.GetRaw(c.Hdr, id.Name)
		org := a.Org.(*Constant).Packed
		if (a.OrgTok == tok.Gt && !(x > org)) || !(x >= org) {
			return False
		}
		end := a.End.(*Constant).Packed
		if (a.EndTok == tok.Lt && !(x < end)) || !(x <= end) {
			return False
		}
		return True
	}
	return OpInRange(a.E.Eval(c), a.OrgTok, a.Org.Eval(c), a.EndTok, a.End.Eval(c))
}

func (a *InRange) Columns() []string {
	cols := a.E.Columns()
	cols = set.Union(cols, a.Org.Columns())
	cols = set.Union(cols, a.End.Columns())
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
	return set.Union(a.E.Columns(), a.M.Columns())
}

func (a *Call) CanEvalRaw(flds []string) bool {
	fn, ok := a.Fn.(*Ident)
	a.evalRaw = ok && len(a.Args) == 1 && IsColumn(a.Args[0].E, flds) &&
		(fn.Name == "Number?" || fn.Name == "String?" || fn.Name == "Date?")
	return a.evalRaw
}

func (a *Call) Eval(c *Context) Value {
	if a.evalRaw {
		fn := a.Fn.(*Ident).Name
		id := a.Args[0].E.(*Ident).Name
		x := c.Row.GetRaw(c.Hdr, id)
		result := false
		switch fn {
		case "Number?":
			result = len(x) > 0 && (x[0] == PackMinus || x[0] == PackPlus)
		case "String?":
			result = x == "" || x[0] == PackString
		case "Date?":
			result = len(x) > 0 && x[0] == PackDate
		}
		return SuBool(result)
	}
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
		cols = set.Union(cols, e.E.Columns())
	}
	return cols
}

func (a *Block) Eval(*Context) Value {
	panic("queries do not support blocks")
}

func (a *Block) Columns() []string {
	panic("queries do not support blocks")
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
