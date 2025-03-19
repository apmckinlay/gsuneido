// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"

	"github.com/apmckinlay/gsuneido/compile/tokens"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

func TestIsField(t *testing.T) {
	flds := []string{"a", "b", "c"}
	test := func(id *Ident, expected bool) {
		t.Helper()
		_, ok := IsField(id, flds)
		assert.T(t).This(ok).Is(expected)
	}
	test(&Ident{Name: "a"}, true)
	test(&Ident{Name: "c"}, true)
	test(&Ident{Name: "x"}, false)
	test(&Ident{Name: "b_lower!"}, true)
	test(&Ident{Name: "x_lower!"}, false)
	test(&Ident{Name: "b_upper!"}, false)
}

func TestCanBeEmpty(t *testing.T) {
	test := func(expr Expr, expected bool) {
		t.Helper()
		actual := CanBeEmpty(expr)
		assert.T(t).This(actual).Is(expected)
	}
	id := &Ident{Name: "x"}
	one := &Constant{Val: One}
	empty := &Constant{Val: EmptyStr}
	foo := &Constant{Val: SuStr("foo")}
	binary := func(tok tok.Token, rhs Expr) Expr {
		return &Binary{Tok: tok, Lhs: id, Rhs: rhs}
	}
	test(binary(tok.Is, empty), true)    // is ""
	test(binary(tok.Is, one), false)     // is 1
	test(binary(tok.Isnt, empty), false) // isnt ""
	test(binary(tok.Isnt, foo), true)    // isnt "foo"
	test(binary(tok.Lt, one), true)      // < 1
	test(binary(tok.Lt, foo), true)      // < "foo"
	test(binary(tok.Lt, empty), false)   // < ""
	test(binary(tok.Lte, one), true)     // <= 1
	test(binary(tok.Lte, foo), true)     // <= "foo"
	test(binary(tok.Lte, empty), true)   // <= ""
	test(binary(tok.Gt, one), true)      // > 1
	test(binary(tok.Gt, foo), false)     // > "foo"
	test(binary(tok.Gt, empty), false)   // > ""
	test(binary(tok.Gte, foo), false)    // >= "foo"
	test(binary(tok.Gte, one), true)     // >= 1
	test(binary(tok.Gte, empty), true)   // >= ""
	test(&Binary{Lhs: foo}, true)
	test(&Binary{Tok: tok.Match, Lhs: id, Rhs: one}, true)
	test(binary(tok.Lt, id), true)

	in := func(exprs ...Expr) Expr {
		return &In{E: id, Exprs: exprs}
	}
	test(in(empty), true)               // in ("")
	test(in(empty, empty, empty), true) // in ("", "", "")
	test(in(one, empty, foo), true)     // in (1, "", "foo")
	test(in(one, foo), false)           // in (1, "foo")
	test(&In{E: one}, true)
	test(in(id), true)
	
	test(&Nary{Tok: tok.Add, Exprs: []Expr{id, id, id}}, true)
}

func BenchmarkEval(b *testing.B) {
	e, ctx := benchSetup()
	for b.Loop() {
		e.Eval(ctx)
	}
}

func BenchmarkEvalRaw(b *testing.B) {
	e, ctx := benchSetup()
	assert.That(e.CanEvalRaw(testFlds))
	for b.Loop() {
		e.Eval(ctx)
	}
}

func benchSetup() (Expr, *Context) {
	// a > b and a < c
	e1 := &Binary{Tok: tok.Gt, Lhs: &Ident{Name: "a"}, Rhs: &Ident{Name: "b"}}
	e2 := &Binary{Tok: tok.Lt, Lhs: &Ident{Name: "a"}, Rhs: &Ident{Name: "c"}}
	e := &Nary{Tok: tok.And, Exprs: []Expr{e1, e2}}
	rec := new(RecordBuilder).
		Add(SuDnum{Dnum: dnum.FromFloat(12.3)}). // a
		Add(SuDnum{Dnum: dnum.FromFloat(45.6)}). // b
		Add(SuDnum{Dnum: dnum.FromFloat(78.9)}). // c
		Build()
	ctx := &Context{
		Th:  &Thread{},
		Hdr: SimpleHeader(testFlds),
		Row: []DbRec{{Record: rec}},
	}
	return e, ctx
}

func TestExpr_CanEvalRaw(t *testing.T) {
	assert := assert.T(t)
	flds := []string{"a", "b", "c"}
	test := func(e Expr) {
		assert.That(e.CanEvalRaw(flds))
	}

	c := &Constant{Val: SuStr("Hello")}
	test(c)
	assert.That(c.Packed != "")

	i := &Ident{Name: "b"}
	test(i)
	i = &Ident{Name: "x"}
	assert.False(i.CanEvalRaw(flds))
	i = &Ident{Name: "a_lower!"}
	test(i)

	u := &Unary{Tok: tok.Not, E: i}
	test(u)
	u = &Unary{Tok: tok.LParen, E: i}
	test(u)

	b := &Binary{Lhs: i, Rhs: c}
	for _, t := range []tokens.Token{tok.Is, tok.Isnt,
		tok.Lt, tok.Lte, tok.Gt, tok.Gte} {
		b.Tok = t
		test(b)
		assert.That(b.evalRaw)
	}

	n := &Nary{Tok: tok.And, Exprs: []Expr{i, c}}
	test(n)
	assert.That(n.evalRaw)
	n = &Nary{Tok: tok.Or, Exprs: []Expr{b, u}}
	test(n)
	assert.That(n.evalRaw)

	in := &In{E: i, Exprs: []Expr{c}}
	test(in)
	assert.That(in.evalRaw)

	ir := &InRange{E: i, Org: c, End: c}
	test(ir)
	assert.That(ir.evalRaw)

	tr := &Trinary{Cond: i, T: b, F: c}
	test(tr)
	assert.That(tr.evalRaw)

	for _, f := range []string{"Number?", "String?", "Date?"} {
		call := &Call{Fn: &Ident{Name: f}, Args: []Arg{{E: c}}}
		test(call)
		assert.That(call.RawEval)
	}
}

func TestRawExpr(t *testing.T) {
	Global.TestDef("Date?", &SuBuiltin1{Fn: func(arg Value) Value {
		return SuBool(arg.Type() == types.Date)
	}, BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1}})
	Global.TestDef("Number?", &SuBuiltin1{Fn: func(arg Value) Value {
		return SuBool(arg.Type() == types.Number)
	}, BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1}})
	Global.TestDef("String?", &SuBuiltin1{Fn: func(arg Value) Value {
		return SuBool(arg.Type() == types.String)
	}, BuiltinParams: BuiltinParams{ParamSpec: ParamSpec1}})

	n := 10000
	if testing.Short() {
		n = 1000
	}
	rec := new(RecordBuilder).
		Add(SuBool(true)).  // a
		Add(SuBool(false)). // b
		Add(IntVal(0)).     // c
		Add(IntVal(1)).     // d
		Add(SuStr("foo")).  // e
		Add(SuStr("bar")).  // f
		Build()
	ctx := &Context{
		Th:  &Thread{},
		Hdr: SimpleHeader(testFlds),
		Row: []DbRec{{Record: rec}},
	}
	var expr Expr
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(expr.Echo())
			fmt.Println(expr)
			panic(e)
		}
	}()
	for range n {
		expr = genRawExpr(0)
		// fmt.Println(e.Echo())
		x := expr.Eval(ctx)
		if !expr.CanEvalRaw(testFlds) {
			t.Error(expr.Echo())
		}
		y := expr.Eval(ctx) // should be raw this time
		assert.This(y).Is(x)
	}
}

var testFlds = []string{"a", "b", "c", "d", "e", "f"}

func genRawExpr(d int) (e Expr) {
	d++
	i := rand.Intn(9)
	if rand.Intn(d) > 1 {
		i = rand.Intn(2)
	}
	switch i {
	case 0:
		v := []Value{SuStr("foo"), SuStr("bar"), Zero, One}[rand.Intn(4)]
		return &Constant{Val: v}
	case 1:
		s := testFlds[rand.Intn(len(testFlds))]
		return &Ident{Name: s}
	case 2:
		if rand.Intn(2) == 0 {
			return &Unary{Tok: tok.Not, E: genRawBoolExpr(d)}
		}
		return &Unary{Tok: tok.LParen, E: genRawExpr(d)}
	case 3:
		t := []tokens.Token{tok.Is, tok.Isnt,
			tok.Lt, tok.Lte, tok.Gt, tok.Gte}[rand.Intn(6)]
		return &Binary{Lhs: genRawExpr(d), Tok: t, Rhs: genRawExpr(d)}
	case 4:
		t := []tokens.Token{tok.And, tok.Or}[rand.Intn(2)]
		var exprs []Expr
		for range 2 + rand.Intn(3) {
			exprs = append(exprs, genRawBoolExpr(d))
		}
		return &Nary{Tok: t, Exprs: exprs}
	case 5:
		var exprs []Expr
		for range 2 + rand.Intn(3) {
			exprs = append(exprs, genRawExpr(d))
		}
		return &In{E: genRawExpr(d), Exprs: exprs}
	case 6:
		return &InRange{E: genRawExpr(d), Org: genRawExpr(d), End: genRawExpr(d)}
	case 7:
		return &Trinary{Cond: genRawBoolExpr(d), T: genRawExpr(d), F: genRawExpr(d)}
	case 8:
		f := []string{"Number?", "String?", "Date?"}[rand.Intn(3)]
		return &Call{Fn: &Ident{Name: f}, Args: []Arg{{E: genRawExpr(d)}}}
	}
	panic(assert.ShouldNotReachHere())
}

func genRawBoolExpr(d int) (e Expr) {
	d++
	i := rand.Intn(9)
	if rand.Intn(d) > 1 {
		i = rand.Intn(2)
	}
	switch i {
	case 0:
		v := []Value{True, False}[rand.Intn(2)]
		return &Constant{Val: v}
	case 1:
		s := []string{"a", "b"}[rand.Intn(2)]
		return &Ident{Name: s}
	case 2:
		if rand.Intn(2) == 0 {
			return &Unary{Tok: tok.Not, E: genRawBoolExpr(d)}
		}
		return &Unary{Tok: tok.LParen, E: genRawBoolExpr(d)}
	case 3:
		t := []tokens.Token{tok.Is, tok.Isnt,
			tok.Lt, tok.Lte, tok.Gt, tok.Gte}[rand.Intn(6)]
		return &Binary{Lhs: genRawExpr(d), Tok: t, Rhs: genRawExpr(d)}
	case 4:
		t := []tokens.Token{tok.And, tok.Or}[rand.Intn(2)]
		var exprs []Expr
		for range 2 + rand.Intn(3) {
			exprs = append(exprs, genRawBoolExpr(d))
		}
		return &Nary{Tok: t, Exprs: exprs}
	case 5:
		var exprs []Expr
		for range 2 + rand.Intn(3) {
			exprs = append(exprs, genRawExpr(d))
		}
		return &In{E: genRawExpr(d), Exprs: exprs}
	case 6:
		return &InRange{E: genRawExpr(d), Org: genRawExpr(d), End: genRawExpr(d)}
	case 7:
		return &Trinary{Cond: genRawBoolExpr(d), T: genRawBoolExpr(d), F: genRawBoolExpr(d)}
	case 8:
		f := []string{"Number?", "String?", "Date?"}[rand.Intn(3)]
		return &Call{Fn: &Ident{Name: f}, Args: []Arg{{E: genRawExpr(d)}}}
	}
	panic(assert.ShouldNotReachHere())
}
