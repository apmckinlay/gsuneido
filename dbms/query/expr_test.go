// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	_ "github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestExprEval(t *testing.T) {
	th := NewThread()
	b := RecordBuilder{}
	b.Add(SuInt(4))
	b.Add(SuInt(5))
	b.Add(SuStr("foo"))
	b.Add(SuStr("bar"))
	rec := b.Build()
	dbrec := DbRec{Record: rec}
	row := Row{dbrec}
	hdr := &Header{Columns: []string{"x", "y", "s", "t"},
		Fields: [][]string{{"x", "y", "s", "t"}}}
	hdr.EnsureMap()
	test := func(src string, expected string) {
		t.Helper()
		p := NewQueryParser(src)
		expr := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		// fmt.Println(expr)
		result := expr.Eval(&ast.Context{T: th, Row: row, Hdr: hdr})
		assert.T(t).This(result.String()).Is(expected)
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
	test("t[1::1]", `"a"`)
	test("x in (3, 4, 5)", "true")
	test("t in (3, 4, 5)", "false")
}
