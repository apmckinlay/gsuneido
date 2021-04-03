// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExprEval(t *testing.T) {
	th := NewThread()
	hdr, row := hdr_row()
	raw := false
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		// fmt.Println(expr)
		assert.That(expr.CanEvalRaw(hdr.GetFields()) == raw)
		result := expr.Eval(&ast.Context{T: th, Row: row, Hdr: hdr})
		assert.T(t).Msg(src).This(result.String()).Is(expected)
	}
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
	test("(s $ t).Size()", "6")
	test("[a: 123].a", "123")
	test("Object(s, x, :t)", `#("foo", 4, t: "bar")`)
	test("t[1::1]", "'a'")
	raw = true
	test("x in (3, 4, 5)", "true")
	test("t in (3, 4, 5)", "false")
	test("x < 9", "true")
	test("9 > x", "true")
	test("s is 'foo'", "true")
}

func hdr_row() (*Header, Row) {
	rb := RecordBuilder{}
	rb.Add(SuInt(4))
	rb.Add(SuInt(5))
	rb.Add(SuStr("foo"))
	rb.Add(SuStr("bar"))
	rec := rb.Build()
	dbrec := DbRec{Record: rec}
	row := Row{dbrec}
	hdr := NewHeader([][]string{{"x", "y", "s", "t"}},
		[]string{"x", "y", "s", "t"})
	return hdr, row
}

func BenchmarkEval(b *testing.B) {
	hdr, row := hdr_row()
	p := NewQueryParser("s is 'foo'")
	expr := p.Expression()
	ctx := &ast.Context{Row: row, Hdr: hdr}
	for i := 0; i < b.N; i++ {
		expr.Eval(ctx)
	}
}

func BenchmarkEval_raw(b *testing.B) {
	hdr, row := hdr_row()
	p := NewQueryParser("s is 'foo'")
	expr := p.Expression()
	expr.CanEvalRaw(hdr.GetFields())
	ctx := &ast.Context{Row: row, Hdr: hdr}
	for i := 0; i < b.N; i++ {
		expr.Eval(ctx)
	}
}

func TestExprColumns(t *testing.T) {
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src)
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
		p := NewQueryParser(src)
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
