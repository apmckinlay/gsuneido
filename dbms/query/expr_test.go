// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExprEval(t *testing.T) {
	options.StrictCompare = true
	options.StrictCompareDb = true
	defer func() {
		options.StrictCompare = false
		options.StrictCompareDb = false
	}()
	th := &Thread{}
	row, hdr := mkrow()
	var raw bool
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		// fmt.Println(expr)
		assert.This(expr.CanEvalRaw(hdr.Physical())).Is(raw)
		result := expr.Eval(&ast.RowContext{Th: th, Row: row, Hdr: hdr})
		assert.T(t).Msg(src).This(result.String()).Is(expected)
	}
	xtest := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		// fmt.Println(expr)
		assert.This(expr.CanEvalRaw(hdr.Columns)).Is(raw)
		assert.This(func() {
			expr.Eval(&ast.RowContext{Th: th, Row: row, Hdr: hdr})
		}).Panics(expected)
	}
	raw = false
	test("x + 2", "6")
	test("1 + y", "6")
	test("x + y", "9")
	test("x - y", "-1")
	test("x * -y", "-20")
	test("(x >> 1) + (y << 1)", "12")
	test("1 + x * y / 10", "3")
	test("x + y is y + x", "true")
	test("x > y and z", "false")
	test("x < y or z", "true")
	test("s $ t", `"foobar"`)
	test("t[1::1]", "'a'")
	test("z > 0 and z < 10", "false")
	test("z >= '' and z < 'z'", "true")
	xtest("x + y < ''", "StrictCompare")

	raw = true
	test("123", "123")
	test("x in (3, 4, 5)", "true")
	test("t in (3, 4, 5)", "false")
	test("x < 9", "true")
	test("9 > x", "true")
	test("x < y and y > x", "true")
	test("s is t", "false")
	test("s is 'foo'", "true")
	test("123 is x", "false")
	// range
	test("x > 0 and 9 > x", "true")
	test("100 < x and x < 200", "false")
	test("0 < s and s < 10", "false") // wrong type
	xtest("x > ''", "StrictCompare")
}

func mkrow() (Row, *Header) {
	rb := RecordBuilder{}
	rb.Add(SuInt(4))     // x
	rb.Add(SuInt(5))     // y
	rb.Add(SuStr("foo")) // s
	rb.Add(SuStr("bar")) // t
	rec := rb.Build()
	dbrec := DbRec{Record: rec}
	row := Row{dbrec}
	hdr := SimpleHeader([]string{"x", "y", "s", "t"})
	return row, hdr
}

func BenchmarkEval(b *testing.B) {
	row, hdr := mkrow()
	p := NewQueryParser("x is 123.456", nil, nil)
	expr := p.Expression()
	ctx := &ast.RowContext{Row: row, Hdr: hdr}
	for b.Loop() {
		expr.Eval(ctx)
	}
}

func BenchmarkEval_raw(b *testing.B) {
	row, hdr := mkrow()
	p := NewQueryParser("x is 123.456", nil, nil)
	expr := p.Expression()
	assert.That(expr.CanEvalRaw(hdr.Columns))
	ctx := &ast.RowContext{Row: row, Hdr: hdr}
	for b.Loop() {
		expr.Eval(ctx)
	}
}

func TestExprColumns(t *testing.T) {
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		result := strings.Join(expr.Columns(), " ")
		assert.T(t).Msg(src).This(result).Is(expected)
	}
	test("123", "")
	test("foo", "foo")
	test("-x", "x")
	test("a < b", "a b")
	test("a * b / c", "a b c")
	test("a ? b : c", "a b c")
	test("a[b..c]", "a b c")
	test("a[b::c]", "a b c")
	test("a in (b, c)", "a b c")
	test("a.b", "a")
	test("a[b]", "a b")
	test("a(b, c)", "a b c")
}

func TestExprRename(t *testing.T) {
	r := &Rename{from: []string{"B", "D"}, to: []string{"b", "d"}}
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		result := renameExpr(expr, r)
		assert.T(t).Msg(src).This(result.Echo()).Is(expected)
	}
	test("123", "123")
	test("foo", "foo")
	test("b", "B")
	test("-x", "-x")
	test("-d", "-D")
	test("a < b", "a < B")
	test("a * b * c * d", "a * B * c * D")
	test("a ? b : c", "a ? B : c")
	test("a[b..c]", "a[B..c]")
	test("b[c::d]", "B[c::D]")
	test("a in (b, c)", "a in (B, c)")
	test("a.b", "a.b")
	test("b.d", "B.d")
	test("a[c]", "a[c]")
	test("b[d]", "B[D]")
	test("f(x, y)", "f(x, y)")
	test("b(c, d)", "B(c, D)")
}

func TestExprReplace(t *testing.T) {
	cols := []string{"x", "y", "z"}
	expr1 := NewQueryParser("5", nil, nil).Expression().(*ast.Constant) // x
	expr2 := NewQueryParser("x + 6", nil, nil).Expression().(*ast.Nary) // y
	expr3 := NewQueryParser("F(y)", nil, nil).Expression().(*ast.Call)  // z
	exprs := []ast.Expr{expr1, expr2, expr3}
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		result := replaceExpr(expr, cols, exprs, false)
		e2, e3 := *expr2, *expr3
		*expr2, *expr3 = ast.Nary{}, ast.Call{}
		assert.T(t).Msg(src).This(result.Echo()).Is(expected)
		*expr2, *expr3 = e2, e3
		result = replaceExpr(expr, cols, exprs, false)
		assert.T(t).Msg(src).This(result.Echo()).Is(expected)
	}
	test("a is x", "a is 5")
	test("x is a", "a is 5") // reversed
	test("x is 6", "false")  // folded
	test("'+' $ x", `"+5"`)
	test("1 + x", `6`)
	test("y", "11")
	test("z", "F(11)")
	test("x + y + z", `16 + F(11)`)

	// extend y = x, z = y where a is z => where a is x
	cols = []string{"y", "z"}
	e1 := NewQueryParser("x", nil, nil).Expression()
	e2 := NewQueryParser("y", nil, nil).Expression()
	exprs = []ast.Expr{e1, e2}
	test("a is z", "a is x")
}
