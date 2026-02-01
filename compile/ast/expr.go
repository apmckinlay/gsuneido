// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"
	"math"
	"strings"

	"slices"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/str"
)

// Unraw is use by Where Simple to *not* evaluate raw
func Unraw(expr Node) Node {
	switch e := expr.(type) {
	case *Constant:
		e.Packed = ""
	case *In:
		e.evalRaw = false
	case *Unary:
		e.evalRaw = false
	case *Binary:
		e.evalRaw = false
	case *InRange:
		e.evalRaw = false
	case *Trinary:
		e.evalRaw = false
	case *Nary:
		e.evalRaw = false
	case *Call:
		e.RawEval = false
	}
	expr.Children(Unraw)
	return expr
}

// Note: Stored fields ignore rules.
// This seems "wrong" but it matches what jSuneido does.
// Normally this isn't an issue because the stored field has a value.

type Context interface {
	GetVal(id string) Value
	GetRaw(id string) string
	Thread() *Thread
}

type RowContext struct {
	Th   *Thread
	Tran *SuTran
	Hdr  *Header
	Row  Row
}

func (c *RowContext) GetVal(id string) Value {
	if ascii.IsUpper(id[0]) && !slices.Contains(c.Hdr.Columns, id) {
		return Global.GetName(c.Thread(), id)
	}
	return c.Row.GetVal(c.Hdr, id, c.Th, c.Tran)
}

func (c *RowContext) GetRaw(id string) string {
	return c.Row.GetRawVal(c.Hdr, id, c.Th, c.Tran)
}

func (c *RowContext) Thread() *Thread {
	return c.Th
}

// Constant ---------------------------------------------------------

func (a *Constant) CanEvalRaw([]string) bool {
	a.Packed = Pack(a.Val.(Packable))
	return true
}

func (a *Constant) EvalRaw(Context) string {
	return a.Packed
}

func (a *Constant) Eval(Context) Value {
	return a.Val
}

func (a *Constant) Columns() []string {
	return []string{}
}

// Ident ------------------------------------------------------------

func (a *Ident) CanEvalRaw(flds []string) bool {
	return slices.Contains(flds, a.Name) ||
		(strings.HasSuffix(a.Name, "_lower!") &&
			slices.Contains(flds, strings.TrimSuffix(a.Name, "_lower!")))
}

func IsField(e Expr, cols []string) (string, bool) {
	if id, ok := e.(*Ident); ok && (slices.Contains(cols, id.Name) ||
		(strings.HasSuffix(id.Name, "_lower!") &&
			slices.Contains(cols, strings.TrimSuffix(id.Name, "_lower!")))) {
		return id.Name, true
	}
	return "", false
}

func (a *Ident) EvalRaw(c Context) string {
	// need GetRawVal to handle PackForward
	return c.GetRaw(a.Name)
}

func (a *Ident) Eval(c Context) Value {
	return c.GetVal(a.Name)
}

func (a *Ident) Columns() []string {
	if str.Capitalized(a.Name) {
		return nil
	}
	return []string{a.Name}
}

// Unary ------------------------------------------------------------

func (a *Unary) CanEvalRaw(flds []string) bool {
	a.evalRaw = a.E.CanEvalRaw(flds) &&
		(a.Tok == tok.Not || a.Tok == tok.LParen)
	return a.evalRaw
}

func (a *Unary) EvalRaw(c Context) string {
	x := a.E.EvalRaw(c)
	if a.Tok == tok.Not {
		return PackBool(UnpackBool(x) == False)
	}
	if a.Tok == tok.LParen {
		return x
	}
	panic(assert.ShouldNotReachHere())
}

func (a *Unary) Eval(c Context) Value {
	if a.evalRaw {
		return Unpack(a.EvalRaw(c))
	}
	return a.eval(a.E.Eval(c))
}

func (a *Unary) eval(val Value) Value {
	switch a.Tok {
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
	}
	panic(assert.ShouldNotReachHere())
}

func (a *Unary) Columns() []string {
	return a.E.Columns()
}

// Binary -----------------------------------------------------------

// CanEvalRaw returns true if Eval doesn't need to unpack the values.
// It sets evalRaw and Packed which are used later by Eval.
func (a *Binary) CanEvalRaw(flds []string) bool {
	// depends on folder putting constant on the right
	a.evalRaw = a.RawOp() && a.Lhs.CanEvalRaw(flds) && a.Rhs.CanEvalRaw(flds)
	return a.evalRaw
}

func (a *Binary) RawOp() bool {
	switch a.Tok {
	case tok.Is, tok.Isnt, tok.Lt, tok.Lte, tok.Gt, tok.Gte:
		return true
	}
	return false
}

func (a *Binary) Eval(c Context) Value {
	if a.evalRaw {
		return UnpackBool(a.EvalRaw(c))
	}
	return a.eval(c.Thread(), a.Lhs.Eval(c), a.Rhs.Eval(c))
}

func (a *Binary) EvalRaw(c Context) string {
	lhs := a.Lhs.EvalRaw(c)
	rhs := a.Rhs.EvalRaw(c)
	switch a.Tok {
	case tok.Is:
		return PackBool(lhs == rhs)
	case tok.Isnt:
		return PackBool(lhs != rhs)
	case tok.Lt:
		return PackBool(packedCmp(lhs, rhs) < 0)
	case tok.Lte:
		return PackBool(packedCmp(lhs, rhs) <= 0)
	case tok.Gt:
		return PackBool(packedCmp(lhs, rhs) > 0)
	case tok.Gte:
		return PackBool(packedCmp(lhs, rhs) >= 0)
	}
	panic(assert.ShouldNotReachHere())
}

func packedCmp(x, y string) int {
	cmp := strings.Compare(x, y)
	if cmp != 0 && options.StrictCompareDb &&
		((x == "" && PackedOrd(y) < OrdStr) ||
			(y == "" && PackedOrd(x) < OrdStr)) {
		panic(fmt.Sprint("StrictCompareDb: ", PackedOrd(x), " ", Unpack(x),
			" <=> ", PackedOrd(y), " ", Unpack(y)))
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
		return Value(SuBool(strictCompare(lhs, rhs) < 0))
	case tok.Lte:
		return Value(SuBool(strictCompare(lhs, rhs) <= 0))
	case tok.Gt:
		return Value(SuBool(strictCompare(lhs, rhs) > 0))
	case tok.Gte:
		return Value(SuBool(strictCompare(lhs, rhs) >= 0))
	case tok.Mod:
		return OpMod(lhs, rhs)
	case tok.LShift:
		return OpLeftShift(lhs, rhs)
	case tok.RShift:
		return OpRightShift(lhs, rhs)
	}
	panic(assert.ShouldNotReachHere())
}

func strictCompare(x Value, y Value) int {
	cmp := x.Compare(y)
	if (cmp&3) == 2 && options.StrictCompareDb &&
		((x == EmptyStr && Order(y) < OrdStr) ||
			(y == EmptyStr && Order(x) < OrdStr)) {
		panic(fmt.Sprint("StrictCompareDb: ", x, " <=> ", y))
	}
	return cmp
}

func (a *Binary) Columns() []string {
	return set.Union(a.Lhs.Columns(), a.Rhs.Columns())
}

// Trinary ----------------------------------------------------------

func (a *Trinary) CanEvalRaw(flds []string) bool {
	c := a.Cond.CanEvalRaw(flds)
	t := a.T.CanEvalRaw(flds)
	f := a.F.CanEvalRaw(flds)
	a.evalRaw = c && t && f
	return a.evalRaw
}

func (a *Trinary) EvalRaw(c Context) string {
	if UnpackBool(a.Cond.EvalRaw(c)) == True {
		return a.T.EvalRaw(c)
	}
	return a.F.EvalRaw(c)
}

func (a *Trinary) Eval(c Context) Value {
	if a.evalRaw {
		return Unpack(a.EvalRaw(c))
	}
	if ToBool(a.Cond.Eval(c)) {
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
	a.evalRaw = false
	if a.Tok == tok.Or || a.Tok == tok.And {
		a.evalRaw = true
		for _, e := range a.Exprs {
			if !e.CanEvalRaw(flds) {
				a.evalRaw = false // don't break, call on all exprs
			}
		}
	}
	return a.evalRaw
}

func (a *Nary) EvalRaw(c Context) string {
	if a.Tok == tok.Or {
		for _, e := range a.Exprs {
			if UnpackBool(e.EvalRaw(c)) == True {
				return PackedTrue
			}
		}
		return PackedFalse
	}
	if a.Tok == tok.And {
		for _, e := range a.Exprs {
			if UnpackBool(e.EvalRaw(c)) == False {
				return PackedFalse
			}
		}
		return PackedTrue
	}
	panic(assert.ShouldNotReachHere())
}

func (a *Nary) Eval(c Context) Value {
	if a.evalRaw {
		return UnpackBool(a.EvalRaw(c))
	}
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

func nary(exprs []Expr, c Context,
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

func muldiv(exprs []Expr, c Context) Value {
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

// RangeTo ----------------------------------------------------------

func (a *RangeTo) Eval(c Context) Value {
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

// RangeLen ---------------------------------------------------------

func (a *RangeLen) Eval(c Context) Value {
	e := a.E.Eval(c)
	from := evalOr(a.From, c, Zero)
	n := evalOr(a.Len, c, MaxInt)
	return e.RangeLen(ToIndex(from), ToInt(n))
}

func (a *RangeLen) Columns() []string {
	return set.Union(a.E.Columns(),
		set.Union(a.From.Columns(), a.Len.Columns()))
}

func evalOr(e Expr, c Context, v Value) Value {
	if e == nil {
		return v
	}
	return e.Eval(c)
}

// In ---------------------------------------------------------------

func (a *In) CanEvalRaw(flds []string) bool {
	a.evalRaw = false
	if a.E.CanEvalRaw(flds) {
		a.evalRaw = true
		for _, e := range a.Exprs {
			if !e.CanEvalRaw(flds) {
				a.evalRaw = false
			}
		}
	}
	return a.evalRaw
}

func (a *In) EvalRaw(c Context) string {
	x := a.E.EvalRaw(c)
	for _, e := range a.Exprs {
		if x == e.EvalRaw(c) {
			return PackedTrue
		}
	}
	return PackedFalse
}

func (a *In) Eval(c Context) Value {
	if a.evalRaw {
		return UnpackBool(a.EvalRaw(c))
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

func (a *InRange) CanEvalRaw(flds []string) bool {
	// InRange already ensures valid operators and constants
	a.evalRaw = a.E.CanEvalRaw(flds)
	a.Org.CanEvalRaw(flds)
	a.End.CanEvalRaw(flds)
	return a.evalRaw
}

func (a *InRange) EvalRaw(c Context) string {
	x := a.E.EvalRaw(c)
	org := a.Org.EvalRaw(c)
	if (a.OrgTok == tok.Gt && !(x > org)) || !(x >= org) {
		return PackedFalse
	}
	end := a.End.EvalRaw(c)
	if (a.EndTok == tok.Lt && !(x < end)) || !(x <= end) {
		return PackedFalse
	}
	return PackedTrue
}

func (a *InRange) Eval(c Context) Value {
	if a.evalRaw {
		return UnpackBool(a.EvalRaw(c))
	}
	return OpInRange(a.E.Eval(c), a.OrgTok, a.Org.Eval(c), a.EndTok, a.End.Eval(c))
}

func (a *InRange) Columns() []string {
	cols := a.E.Columns()
	cols = set.Union(cols, a.Org.Columns())
	cols = set.Union(cols, a.End.Columns())
	return cols
}

// Mem --------------------------------------------------------------

func (a *Mem) Eval(c Context) Value {
	e := a.E.Eval(c)
	m := a.M.Eval(c)
	result := e.Get(nil, m)
	if result == nil {
		MemberNotFound(m)
	}
	return result
}

func (a *Mem) Columns() []string {
	return set.Union(a.E.Columns(), a.M.Columns())
}

// Call -------------------------------------------------------------

func (a *Call) CanEvalRaw(flds []string) bool {
	fn, ok := a.Fn.(*Ident)
	a.RawEval = ok && len(a.Args) == 1 && a.Args[0].E.CanEvalRaw(flds) &&
		(fn.Name == "Number?" || fn.Name == "String?" || fn.Name == "Date?")
	for i := 1; i < len(a.Args); i++ {
		a.Args[i].E.CanEvalRaw(flds)
	}
	return a.RawEval
}

func (a *Call) EvalRaw(c Context) string {
	x := a.Args[0].E.EvalRaw(c)
	result := false
	switch a.Fn.(*Ident).Name {
	case "Number?":
		result = len(x) > 0 && (x[0] == PackMinus || x[0] == PackPlus)
	case "String?":
		result = x == "" || x[0] == PackString
	case "Date?":
		result = len(x) > 0 && x[0] == PackDate
	default:
		assert.ShouldNotReachHere()
	}
	return PackBool(result)
}

func (a *Call) Eval(c Context) Value {
	if a.RawEval {
		return UnpackBool(a.EvalRaw(c))
	}
	as := argspec(a.Args) //TODO cache
	args := make([]Value, len(a.Args))
	for i, a := range a.Args {
		args[i] = a.E.Eval(c)
	}
	var fn Value
	var this Value
	switch f := a.Fn.(type) {
	case *Mem:
		this = f.E.Eval(c)
		meth := f.M.Eval(c)
		fn = c.Thread().Lookup(this, ToStr(meth))
	default:
		fn = a.Fn.Eval(c)
	}
	return c.Thread().PushCall(fn, this, as, args...)
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

// Block ------------------------------------------------------------

func (a *Block) Eval(Context) Value {
	panic("queries do not support blocks")
}

func (a *Block) Columns() []string {
	panic("queries do not support blocks")
}

//-------------------------------------------------------------------

// CanBeEmpty returns true if the expression could match an empty string.
// It is conservative and handles both raw where "" is less than everything,
// and values where "" is between numbers and strings.
func CanBeEmpty(e Expr) bool {
	switch e := e.(type) {
	case *Binary:
		if _, ok := e.Lhs.(*Ident); !ok {
			return true
		}
		c, ok := e.Rhs.(*Constant)
		if !ok {
			return true
		}
		switch e.Tok {
		case tok.Is:
			return c.Val == EmptyStr
		case tok.Isnt:
			return c.Val != EmptyStr
		case tok.Lt:
			return c.Val != EmptyStr
		case tok.Lte:
			return true
		case tok.Gt:
			return c.Val.Compare(EmptyStr) < 0
		case tok.Gte:
			return c.Val.Compare(EmptyStr) <= 0
		default:
			return true
		}
	case *In:
		if _, ok := e.E.(*Ident); !ok {
			return true
		}
		for _, e := range e.Exprs {
			if c, ok := e.(*Constant); !ok || c.Val == EmptyStr {
				return true
			}
		}
		return false
	}
	return true
}
