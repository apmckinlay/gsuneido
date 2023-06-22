// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/options"
	. "github.com/apmckinlay/gsuneido/runtime"
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
	raw := false
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		// fmt.Println(expr)
		assert.This(expr.CanEvalRaw(hdr.Physical())).Is(raw)
		result := expr.Eval(&ast.Context{Th: th, Row: row, Hdr: hdr})
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
			expr.Eval(&ast.Context{Th: th, Row: row, Hdr: hdr})
		}).Panics(expected)
	}
	// not raw
	test("123", "123")
	test("x + 2", "6")
	test("1 + y", "6")
	test("x + y", "9")
	test("x - y", "-1")
	test("x * -y", "-20")
	test("(x >> 1) + (y << 1)", "12")
	test("1 + x * y / 10", "3")
	test("x + y is y + x", "true")
	test("x < y and y > x", "true")
	test("x > y and z", "false")
	test("x < y or z", "true")
	test("s $ t", `"foobar"`)
	test("s is t", "false")
	test("t[1::1]", "'a'")
	test("z > 0 and z < 10", "false")
	test("z >= '' and z < 'z'", "true")
	xtest("x + y < 'foo'", "StrictCompare")

	raw = true
	test("x in (3, 4, 5)", "true")
	test("t in (3, 4, 5)", "false")
	test("x < 9", "true")
	test("9 > x", "true")
	test("s is 'foo'", "true")
	test("123 is x", "false")
	// range
	test("x > 0 and 9 > x", "true")
	test("100 < x and x < 200", "false")
	test("0 < s and s < 10", "false") // wrong type
	xtest("s > 10", "StrictCompare")
}

func mkrow() (Row, *Header) {
	rb := RecordBuilder{}
	rb.Add(SuInt(4))
	rb.Add(SuInt(5))
	rb.Add(SuStr("foo"))
	rb.Add(SuStr("bar"))
	rec := rb.Build()
	dbrec := DbRec{Record: rec}
	row := Row{dbrec}
	flds := []string{"x", "y", "s", "t"}
	hdr := NewHeader([][]string{flds}, flds)
	return row, hdr
}

func BenchmarkEval(b *testing.B) {
	row, hdr := mkrow()
	p := NewQueryParser("x is 123.456", nil, nil)
	expr := p.Expression()
	ctx := &ast.Context{Row: row, Hdr: hdr}
	for i := 0; i < b.N; i++ {
		expr.Eval(ctx)
	}
}

func BenchmarkEval_raw(b *testing.B) {
	row, hdr := mkrow()
	p := NewQueryParser("x is 123.456", nil, nil)
	expr := p.Expression()
	assert.That(expr.CanEvalRaw(hdr.Columns))
	ctx := &ast.Context{Row: row, Hdr: hdr}
	for i := 0; i < b.N; i++ {
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
	from := []string{"b", "d"}
	to := []string{"B", "D"}
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		result := renameExpr(expr, from, to)
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
	from := []string{"x", "y", "z"}
	expr := NewQueryParser("5", nil, nil).Expression()
	to := []ast.Expr{expr, expr, expr}
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src, nil, nil)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		result := replaceExpr(expr, from, to)
		assert.T(t).Msg(src).This(result.Echo()).Is(expected)
		result = replaceExpr(expr, from, to)
		assert.T(t).Msg(src).This(result.Echo()).Is(expected)
	}
	test("x is a", "a is 5")
	test("x is 6", "false") // folded
	test("'=' $ x", `"=5"`)
	test("1 + x", `6`)
	test("x + y + z", `15`)
}
